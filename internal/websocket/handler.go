package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"pulse/internal/auth/token"
	"pulse/internal/config"
	"pulse/internal/store"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period (must be less than pongWait)
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: restrict this in production
		return true
	},
}

// HandleWebSocket handles WebSocket connections
// Authentication is done via a message after connection is established
func (h *Hub) HandleWebSocket(cfg *config.Config, s *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Upgrade HTTP connection to WebSocket (no auth required at this stage)
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("WebSocket upgrade error: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upgrade connection"})
			return
		}

		// Create unauthenticated client
		client := &Client{
			ID:            uuid.New().String(),
			UserID:        uuid.Nil,
			Authenticated: false,
			Hub:           h,
			Send:          make(chan []byte, 256),
			conn:          conn,
		}

		log.Printf("WebSocket client connected: %s (awaiting authentication)", client.ID)

		// Register client (will be indexed by user ID after authentication)
		h.register <- client

		// Start goroutines for reading and writing
		// Authentication will happen in readPump when auth message is received
		go client.writePump()
		go client.readPump(cfg, s)
	}
}

// AuthMessage represents an authentication message from the client
type AuthMessage struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}

// readPump pumps messages from the WebSocket connection to the hub
func (c *Client) readPump(cfg *config.Config, s *store.Store) {
	jwtGenerator := token.NewJWTTokenGenerator(token.TokenConfig{
		Secret:   cfg.JWTSecret,
		Validity: 24 * time.Hour,
	})

	// Set authentication timeout (30 seconds)
	authTimeout := time.NewTimer(30 * time.Second)
	authDone := make(chan bool, 1)

	defer func() {
		authTimeout.Stop()
		c.Hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// Handle authentication timeout
	go func() {
		select {
		case <-authTimeout.C:
			if !c.Authenticated {
				log.Printf("WebSocket client %s failed to authenticate within timeout", c.ID)
				c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "Authentication timeout"))
				c.conn.Close()
			}
		case <-authDone:
			authTimeout.Stop()
		}
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}

		// Parse message as JSON
		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Failed to parse message from client %s: %v", c.ID, err)
			continue
		}

		// Handle subscription messages (only if authenticated)
		if c.Authenticated {
			if msgType, ok := msg["type"].(string); ok {
				switch msgType {
				case "subscribe":
					topic, ok := msg["topic"].(string)
					if !ok || topic == "" {
						c.sendError("Invalid subscribe message: topic required")
						continue
					}
					c.Hub.Subscribe(c, topic)
					// Send confirmation
					successMsg := Message{
						Type:    "subscribed",
						Payload: map[string]interface{}{"topic": topic},
					}
					data, _ := json.Marshal(successMsg)
					c.Send <- data
					continue

				case "unsubscribe":
					topic, ok := msg["topic"].(string)
					if !ok || topic == "" {
						c.sendError("Invalid unsubscribe message: topic required")
						continue
					}
					c.Hub.Unsubscribe(c, topic)
					// Send confirmation
					successMsg := Message{
						Type:    "unsubscribed",
						Payload: map[string]interface{}{"topic": topic},
					}
					data, _ := json.Marshal(successMsg)
					c.Send <- data
					continue
				}
			}
		}

		// Handle authentication message
		if msgType, ok := msg["type"].(string); ok && msgType == "auth" {
			if !c.Authenticated {
				tokenStr, ok := msg["token"].(string)
				if !ok || tokenStr == "" {
					c.sendError("Invalid auth message: token required")
					continue
				}

				// Validate token
				claims, err := jwtGenerator.Validate(tokenStr)
				if err != nil {
					c.sendError("Invalid or expired token")
					continue
				}

				// Validate session
				if claims.ID == "" {
					c.sendError("Token missing session identifier")
					continue
				}

				_, err = s.GetSessionByJTI(claims.ID)
				if err != nil {
					c.sendError("Invalid or expired session")
					continue
				}

				// Authenticate client
				c.mu.Lock()
				c.UserID = claims.UserID
				c.Authenticated = true
				c.mu.Unlock()

				// Re-register to update user index
				c.Hub.unregister <- c
				c.Hub.register <- c

				// Update session activity (non-blocking)
				go func() {
					_ = s.UpdateSessionActivity(claims.ID)
				}()

				// Send success message
				successMsg := Message{
					Type:    "auth_success",
					Payload: map[string]interface{}{"message": "Authentication successful"},
				}
				data, _ := json.Marshal(successMsg)
				c.Send <- data

				log.Printf("WebSocket client authenticated: %s (user: %s)", c.ID, claims.UserID)
				authDone <- true
			} else {
				c.sendError("Already authenticated")
			}
			continue
		}

		// Only process other messages if authenticated
		if !c.Authenticated {
			c.sendError("Not authenticated. Send auth message first")
			continue
		}

		// Process other message types here if needed
		log.Printf("Received message from client %s (user: %s): %s", c.ID, c.UserID, string(message))
	}
}

// sendError sends an error message to the client
func (c *Client) sendError(message string) {
	errorMsg := Message{
		Type:    "error",
		Payload: map[string]interface{}{"message": message},
	}
	data, err := json.Marshal(errorMsg)
	if err != nil {
		return
	}
	select {
	case c.Send <- data:
	default:
		log.Printf("Warning: Failed to send error to client %s", c.ID)
	}
}

// writePump pumps messages from the hub to the WebSocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current websocket message
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
