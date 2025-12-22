# WebSocket Server

This package provides a WebSocket server that can broadcast status updates to connected frontend clients with authentication and scoped broadcasting support.

## Features

- **Authentication**: All WebSocket connections require valid JWT authentication
- **Pub/Sub**: Clients can subscribe to topics and only receive messages for subscribed topics
- **User-scoped broadcasting**: Send messages to all connections of a specific user
- **Project-scoped broadcasting**: Send messages to all members of a project
- **Topic-based broadcasting**: Send messages to all clients subscribed to a specific topic
- **Global broadcasting**: Send messages to all connected clients

## Authentication

All WebSocket connections require authentication via JWT token. Authentication is done **after** the WebSocket connection is established by sending an authentication message.

1. **Connect** to the WebSocket endpoint (no auth required at this stage)
2. **Send an auth message** with your JWT token:

   ```json
   {
     "type": "auth",
     "token": "<your_jwt_token>"
   }
   ```

3. **Receive confirmation** - you'll get an `auth_success` message if authentication succeeds

**Authentication timeout**: You have 30 seconds after connection to authenticate. If you don't authenticate within this time, the connection will be closed.

**Error handling**: If authentication fails, you'll receive an `error` message with details.

## Usage

### Broadcasting to All Clients

```go
import "pulse/internal/websocket"

// Broadcast to all connected clients
err := websocket.Broadcast("check_status_update", map[string]interface{}{
    "check_id": "123",
    "status": "up",
    "timestamp": time.Now(),
})
```

### User-Scoped Broadcasting

When an event is scoped to a user, it will be broadcasted to **all active connections** of that user. This means if a user has multiple browser tabs or devices open, they will all receive the message.

```go
import (
    "pulse/internal/websocket"
    "github.com/google/uuid"
)

userID := uuid.MustParse("user-id-here")

// Broadcast to all connections of a specific user
err := websocket.BroadcastToUser(userID, "notification", map[string]interface{}{
    "title": "New alert",
    "message": "Your check is down",
})
```

**Use cases:**

- User-specific notifications
- Account updates
- Personal alerts

### Project-Scoped Broadcasting

When an event is scoped to a project, it will be broadcasted to **all users who are members of that project**. Each member will receive the message on all their active connections.

```go
import (
    "pulse/internal/websocket"
    "github.com/google/uuid"
    "pulse/internal/store"
)

projectID := uuid.MustParse("project-id-here")
store := // your store instance

// Broadcast to all members of a project
err := websocket.BroadcastToProject(store, projectID, "check_run_completed", map[string]interface{}{
    "project_id": projectID.String(),
    "check_id": "check-123",
    "status": "down",
    "timestamp": time.Now(),
})
```

**Use cases:**

- Check run status updates
- Project-wide alerts
- Team notifications
- Real-time dashboard updates

**How it works:**

1. The system looks up all users who are members of the project
2. For each member, it finds all their active WebSocket connections
3. The message is sent to all those connections

### Example: Broadcasting Check Run Status

```go
// In your check run handler or worker
import (
    "pulse/internal/websocket"
    "github.com/google/uuid"
)

// After a check run completes, broadcast to all project members
websocket.BroadcastToProject(store, check.ProjectID, "check_run_completed", map[string]interface{}{
    "project_id": check.ProjectID.String(),
    "check_id": check.ID.String(),
    "run_id": run.ID.String(),
    "status": run.Status,
    "duration": run.Duration,
    "timestamp": run.CreatedAt,
})
```

### Example: Broadcasting User Notification

```go
// In your notification handler
import (
    "pulse/internal/websocket"
    "github.com/google/uuid"
)

// Send a notification to a specific user
websocket.BroadcastToUser(userID, "notification", map[string]interface{}{
    "type": "alert",
    "title": "Check Alert",
    "message": "Your check 'API Health' is down",
    "check_id": checkID.String(),
})
```

### Example: Broadcasting Alert

```go
// In your alerter
import (
    "pulse/internal/websocket"
    "github.com/google/uuid"
)

// Broadcast alert to all project members
websocket.BroadcastToProject(store, alert.ProjectID, "alert_created", map[string]interface{}{
    "alert_id": alert.ID.String(),
    "check_id": alert.CheckID.String(),
    "status": alert.Status,
    "timestamp": alert.CreatedAt,
})
```

## Message Format

All messages are sent in the following JSON format:

```json
{
  "type": "message_type",
  "payload": { ... }
}
```

## Frontend Connection

Connect from your frontend and authenticate using a message:

```javascript
// Get token from your auth system
const token = getAuthToken() // Your function to get JWT token

// Connect to WebSocket (no auth required at connection time)
const ws = new WebSocket('ws://localhost:8080/ws')

ws.onopen = () => {
  console.log('WebSocket connected')

  // Authenticate immediately after connection
  ws.send(JSON.stringify({
    type: 'auth',
    token
  }))
}

ws.onmessage = (event) => {
  const message = JSON.parse(event.data)

  // Handle authentication response
  if (message.type === 'auth_success') {
    console.log('Authentication successful')

    // Subscribe to topics you're interested in
    ws.send(JSON.stringify({
      type: 'subscribe',
      topic: `project:${currentProjectId}`
    }))
  } else if (message.type === 'subscribed') {
    console.log('Subscribed to topic:', message.payload.topic)
  } else if (message.type === 'unsubscribed') {
    console.log('Unsubscribed from topic:', message.payload.topic)
  } else if (message.type === 'error') {
    console.error('Error:', message.payload.message)
  } else {
    // Handle other message types (broadcasts)
    console.log('Received:', message.type, message.payload)

    switch (message.type) {
      case 'check_run_completed':
        // Update UI with check run result
        updateCheckStatus(message.payload)
        break
      case 'alert_created':
        // Show alert notification
        showAlert(message.payload)
        break
      case 'notification':
        // Show user notification
        showNotification(message.payload)
        break
    }
  }
}

ws.onerror = (error) => {
  console.error('WebSocket error:', error)
}

ws.onclose = () => {
  console.log('WebSocket disconnected')
  // Optionally reconnect
}
```

## Architecture

- **Hub**: Manages all WebSocket connections and handles broadcasting
- **Client**: Represents a single WebSocket connection with associated user ID
- **Global Broadcaster**: Provides easy access to broadcast from anywhere in the codebase
- **User Index**: Tracks clients by user ID for efficient user-scoped broadcasting
- **Project Broadcasting**: Uses store to look up project members and broadcasts to all their connections

The hub runs in a separate goroutine and handles:

- Client registration/unregistration
- User indexing for efficient lookups
- Message broadcasting (global, user-scoped, project-scoped)
- Connection management

## Connection Behavior

- **Multiple connections per user**: Users can have multiple active connections (multiple tabs, devices, etc.)
- **User-scoped events**: All connections of a user receive the message
- **Project-scoped events**: All members of a project receive the message on all their connections
- **Automatic cleanup**: Disconnected clients are automatically removed from the hub
