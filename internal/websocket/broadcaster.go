package websocket

import (
	"log"
	"sync"

	"github.com/google/uuid"

	"pulse/internal/store"
)

// Example usage from anywhere in the codebase:
//
//   import "pulse/internal/websocket"
//
//   // Broadcast to all connected clients
//   websocket.Broadcast("check_status_update", map[string]interface{}{
//       "check_id": "123",
//       "status": "up",
//   })
//
//   // Broadcast to a specific client
//   websocket.BroadcastTo("client-id", "check_status_update", map[string]interface{}{
//       "check_id": "123",
//       "status": "up",
//   })

var (
	globalHub *Hub
	hubOnce   sync.Once
	hubMu     sync.RWMutex
)

// InitGlobalHub initializes the global WebSocket hub
// This should be called once during application startup
func InitGlobalHub() *Hub {
	hubOnce.Do(func() {
		globalHub = NewHub()
		go globalHub.Run()
		log.Println("Global WebSocket hub initialized")
	})
	return globalHub
}

// GetGlobalHub returns the global WebSocket hub instance
// Returns nil if InitGlobalHub hasn't been called yet
func GetGlobalHub() *Hub {
	hubMu.RLock()
	defer hubMu.RUnlock()
	return globalHub
}

// Broadcast sends a message to all connected WebSocket clients
// This is a convenience function that uses the global hub
func Broadcast(messageType string, payload interface{}) error {
	hub := GetGlobalHub()
	if hub == nil {
		log.Println("Warning: Attempted to broadcast but global hub is not initialized")
		return nil
	}
	return hub.Broadcast(messageType, payload)
}

// BroadcastTo sends a message to a specific client by ID
// This is a convenience function that uses the global hub
func BroadcastTo(clientID string, messageType string, payload interface{}) error {
	hub := GetGlobalHub()
	if hub == nil {
		log.Println("Warning: Attempted to broadcast but global hub is not initialized")
		return nil
	}
	return hub.BroadcastTo(clientID, messageType, payload)
}

// GetClientCount returns the number of connected clients
func GetClientCount() int {
	hub := GetGlobalHub()
	if hub == nil {
		return 0
	}
	return hub.GetClientCount()
}

// GetClientIDs returns a list of all connected client IDs
func GetClientIDs() []string {
	hub := GetGlobalHub()
	if hub == nil {
		return []string{}
	}
	return hub.GetClientIDs()
}

// BroadcastToUser sends a message to all connections of a specific user
// If the user has multiple browser tabs/devices open, all will receive the message
func BroadcastToUser(userID uuid.UUID, messageType string, payload interface{}) error {
	hub := GetGlobalHub()
	if hub == nil {
		log.Println("Warning: Attempted to broadcast but global hub is not initialized")
		return nil
	}
	return hub.BroadcastToUser(userID, messageType, payload)
}

// BroadcastToProject sends a message to all users who are members of a project
// All members of the project will receive the message on all their active connections
// Requires a store instance to look up project members
func BroadcastToProject(s *store.Store, projectID uuid.UUID, messageType string, payload interface{}) error {
	hub := GetGlobalHub()
	if hub == nil {
		log.Println("Warning: Attempted to broadcast but global hub is not initialized")
		return nil
	}

	// Create a function to get project member IDs
	getProjectMemberIDs := func(pid uuid.UUID) ([]uuid.UUID, error) {
		return s.GetProjectMemberUserIDs(pid)
	}

	return hub.BroadcastToProject(getProjectMemberIDs, projectID, messageType, payload)
}

// BroadcastToTopic sends a message to all clients subscribed to a specific topic
// Only clients who have subscribed to the topic will receive the message
func BroadcastToTopic(topic string, messageType string, payload interface{}) error {
	hub := GetGlobalHub()
	if hub == nil {
		log.Println("Warning: Attempted to broadcast but global hub is not initialized")
		return nil
	}
	return hub.BroadcastToTopic(topic, messageType, payload)
}
