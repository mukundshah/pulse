package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Message represents a WebSocket message
type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// Client represents a WebSocket client connection
type Client struct {
	ID            string
	UserID        uuid.UUID // User ID of the authenticated user (uuid.Nil if not authenticated)
	Authenticated bool      // Whether the client has been authenticated
	Send          chan []byte
	Hub           *Hub
	conn          *websocket.Conn
	mu            sync.Mutex
}

// Hub maintains the set of active clients and broadcasts messages to them
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Clients indexed by user ID (one user can have multiple connections)
	clientsByUser map[uuid.UUID]map[*Client]bool

	// Pub/Sub: Topics to clients mapping
	subscriptions map[string]map[*Client]bool

	// Inbound messages from clients
	broadcast chan []byte

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	mu sync.RWMutex
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		clients:       make(map[*Client]bool),
		clientsByUser: make(map[uuid.UUID]map[*Client]bool),
		broadcast:     make(chan []byte, 256),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			// Index by user ID (only if authenticated)
			if client.Authenticated && client.UserID != uuid.Nil {
				if h.clientsByUser[client.UserID] == nil {
					h.clientsByUser[client.UserID] = make(map[*Client]bool)
				}
				h.clientsByUser[client.UserID][client] = true
			}
			h.mu.Unlock()
			if client.Authenticated {
				log.Printf("WebSocket client connected: %s (user: %s, total: %d)", client.ID, client.UserID, len(h.clients))
			} else {
				log.Printf("WebSocket client connected: %s (awaiting auth, total: %d)", client.ID, len(h.clients))
			}

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				// Remove from user index (only if authenticated)
				if client.Authenticated && client.UserID != uuid.Nil {
					if userClients, ok := h.clientsByUser[client.UserID]; ok {
						delete(userClients, client)
						if len(userClients) == 0 {
							delete(h.clientsByUser, client.UserID)
						}
					}
				}
				// Remove from all topic subscriptions
				for topic, subscribers := range h.subscriptions {
					if subscribers[client] {
						delete(subscribers, client)
						if len(subscribers) == 0 {
							delete(h.subscriptions, topic)
						}
					}
				}
				close(client.Send)
			}
			h.mu.Unlock()
			if client.Authenticated {
				log.Printf("WebSocket client disconnected: %s (user: %s, total: %d)", client.ID, client.UserID, len(h.clients))
			} else {
				log.Printf("WebSocket client disconnected: %s (total: %d)", client.ID, len(h.clients))
			}

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Broadcast sends a message to all connected clients
func (h *Hub) Broadcast(messageType string, payload interface{}) error {
	msg := Message{
		Type:    messageType,
		Payload: payload,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	select {
	case h.broadcast <- data:
	default:
		log.Printf("Warning: WebSocket broadcast channel full, dropping message")
	}

	return nil
}

// BroadcastTo sends a message to a specific client by ID
func (h *Hub) BroadcastTo(clientID string, messageType string, payload interface{}) error {
	msg := Message{
		Type:    messageType,
		Payload: payload,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client.ID == clientID {
			select {
			case client.Send <- data:
			default:
				log.Printf("Warning: Client %s send channel full, dropping message", clientID)
			}
			return nil
		}
	}

	return nil
}

// GetClientCount returns the number of connected clients
func (h *Hub) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// GetClientIDs returns a list of all connected client IDs
func (h *Hub) GetClientIDs() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	ids := make([]string, 0, len(h.clients))
	for client := range h.clients {
		ids = append(ids, client.ID)
	}
	return ids
}

// BroadcastToUser sends a message to all connections of a specific user
func (h *Hub) BroadcastToUser(userID uuid.UUID, messageType string, payload interface{}) error {
	msg := Message{
		Type:    messageType,
		Payload: payload,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	userClients, ok := h.clientsByUser[userID]
	if !ok {
		// User has no active connections
		return nil
	}

	sent := 0
	for client := range userClients {
		select {
		case client.Send <- data:
			sent++
		default:
			log.Printf("Warning: Client %s (user %s) send channel full, dropping message", client.ID, userID)
		}
	}

	if sent > 0 {
		log.Printf("Broadcasted to user %s: %d connection(s)", userID, sent)
	}

	return nil
}

// BroadcastToProject sends a message to all users who are members of a project
// This requires a function to get project member user IDs
func (h *Hub) BroadcastToProject(getProjectMemberIDs func(uuid.UUID) ([]uuid.UUID, error), projectID uuid.UUID, messageType string, payload interface{}) error {
	// Get all user IDs that are members of this project
	userIDs, err := getProjectMemberIDs(projectID)
	if err != nil {
		return err
	}

	if len(userIDs) == 0 {
		// No members in project
		return nil
	}

	msg := Message{
		Type:    messageType,
		Payload: payload,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	totalSent := 0
	for _, userID := range userIDs {
		userClients, ok := h.clientsByUser[userID]
		if !ok {
			continue
		}

		for client := range userClients {
			select {
			case client.Send <- data:
				totalSent++
			default:
				log.Printf("Warning: Client %s (user %s) send channel full, dropping message", client.ID, userID)
			}
		}
	}

	if totalSent > 0 {
		log.Printf("Broadcasted to project %s: %d connection(s) across %d user(s)", projectID, totalSent, len(userIDs))
	}

	return nil
}

// Subscribe adds a client to a topic subscription
func (h *Hub) Subscribe(client *Client, topic string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.subscriptions[topic] == nil {
		h.subscriptions[topic] = make(map[*Client]bool)
	}
	h.subscriptions[topic][client] = true
	log.Printf("Client %s subscribed to topic: %s", client.ID, topic)
}

// Unsubscribe removes a client from a topic subscription
func (h *Hub) Unsubscribe(client *Client, topic string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if subscribers, ok := h.subscriptions[topic]; ok {
		delete(subscribers, client)
		if len(subscribers) == 0 {
			delete(h.subscriptions, topic)
		}
		log.Printf("Client %s unsubscribed from topic: %s", client.ID, topic)
	}
}

// BroadcastToTopic sends a message to all clients subscribed to a specific topic
func (h *Hub) BroadcastToTopic(topic string, messageType string, payload interface{}) error {
	msg := Message{
		Type:    messageType,
		Payload: payload,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	subscribers, ok := h.subscriptions[topic]
	if !ok {
		// No subscribers for this topic
		return nil
	}

	sent := 0
	for client := range subscribers {
		// Only send to authenticated clients
		if !client.Authenticated {
			continue
		}
		select {
		case client.Send <- data:
			sent++
		default:
			log.Printf("Warning: Client %s (topic %s) send channel full, dropping message", client.ID, topic)
		}
	}

	if sent > 0 {
		log.Printf("Broadcasted to topic %s: %d subscriber(s)", topic, sent)
	}

	return nil
}
