# Pulse WebSocket Composables

A modular and scalable WebSocket client implementation for Pulse using VueUse's `useWebSocket`.

## Features

- ✅ Authentication with JWT tokens
- ✅ Pub/Sub topic subscriptions
- ✅ Auto-reconnect with exponential backoff
- ✅ Heartbeat/ping support
- ✅ Type-safe message handling
- ✅ Reactive state management
- ✅ Modular design for easy extension

## Installation

The composables are automatically available in your Nuxt app. No installation needed!

## Basic Usage

### Simple Connection

```vue
<script setup lang="ts">
const ws = usePulseWebSocket({
  onMessage: (message) => {
    console.log('Received:', message.type, message.payload)
  }
})
</script>

<template>
  <div>
    <p>Status: {{ ws.status }}</p>
    <p>Authenticated: {{ ws.isAuthenticated }}</p>
  </div>
</template>
```

### With Auto-Subscription

```vue
<script setup lang="ts">
const projectId = 'project-123'

const ws = usePulseWebSocket({
  autoSubscribe: [`project:${projectId}`],
  onMessage: (message) => {
    if (message.type === 'check_run_completed') {
      // Handle check run update
      console.log('Check run:', message.payload)
    }
  }
})
</script>
```

## Advanced Usage

### Topic-Specific Subscription

```vue
<script setup lang="ts">
const { isSubscribed, lastMessage } = usePulseWebSocketTopic({
  topic: computed(() => `project:${projectId.value}`),
  onMessage: (message) => {
    // Only receives messages for this topic
    console.log('Project update:', message.payload)
  }
})
</script>
```

### Message Handler Registry

```vue
<script setup lang="ts">
const ws = usePulseWebSocket()
const handler = usePulseWebSocketMessageHandler({ ws })

// Register handlers for specific message types
handler.on('check_run_completed', (msg) => {
  console.log('Check run:', msg.payload)
})

handler.on('alert_created', (msg) => {
  console.log('Alert:', msg.payload)
})

// One-time handler
handler.once('notification', (msg) => {
  console.log('Notification:', msg.payload)
})

// Wildcard handler (all messages)
handler.on('*', (msg) => {
  console.log('Any message:', msg)
})
</script>
```

### Manual Subscription Management

```vue
<script setup lang="ts">
const ws = usePulseWebSocket()

// Subscribe to topics
ws.subscribe('project:123')
ws.subscribe('check:456')

// Subscribe to multiple topics
ws.subscribeMany(['project:123', 'check:456'])

// Check subscription status
if (ws.isSubscribed('project:123')) {
  console.log('Subscribed to project:123')
}

// Unsubscribe
ws.unsubscribe('project:123')

// Clear all subscriptions
ws.clearSubscriptions()
</script>
```

## Using the Global Plugin

The plugin provides a global WebSocket connection accessible via `$ws`:

```vue
<script setup lang="ts">
const { $ws } = useNuxtApp()

// Use the global connection
$ws.subscribe('project:123')

// Or create a handler
const handler = usePulseWebSocketMessageHandler({ ws: $ws })
handler.on('check_run_completed', (msg) => {
  console.log('Check run:', msg.payload)
})
</script>
```

## API Reference

### `usePulseWebSocket(options)`

Main WebSocket composable.

**Options:**
- `url?: string | Ref<string>` - WebSocket URL (default: `ws://localhost:8080/ws`)
- `immediate?: boolean` - Auto-connect on mount (default: `true`)
- `autoReconnect?: boolean | object` - Auto-reconnect config
- `heartbeat?: boolean | object` - Heartbeat config
- `autoSubscribe?: string[]` - Topics to auto-subscribe
- `onAuthSuccess?: () => void` - Auth success callback
- `onAuthError?: (error: string) => void` - Auth error callback
- `onSubscribed?: (topic: string) => void` - Subscription callback
- `onUnsubscribed?: (topic: string) => void` - Unsubscription callback
- `onError?: (error: string) => void` - Error callback
- `onMessage?: (message: WebSocketMessage) => void` - Message callback

**Returns:**
- `status: Ref<WebSocketStatus>` - Connection status
- `data: Ref<WebSocketMessage | null>` - Last message
- `isOpen: Ref<boolean>` - Whether connection is open
- `isAuthenticated: Ref<boolean>` - Whether authenticated
- `subscriptions: Ref<Set<string>>` - Active subscriptions
- `open()` - Open connection
- `close()` - Close connection
- `send(message)` - Send message
- `authenticate(token)` - Authenticate with token
- `subscribe(topic)` - Subscribe to topic
- `unsubscribe(topic)` - Unsubscribe from topic
- `subscribeMany(topics)` - Subscribe to multiple topics
- `unsubscribeMany(topics)` - Unsubscribe from multiple topics
- `isSubscribed(topic)` - Check if subscribed
- `clearSubscriptions()` - Clear all subscriptions

### `usePulseWebSocketTopic(options)`

Topic-specific subscription composable.

**Options:**
- `topic: string | Ref<string>` - Topic to subscribe to
- `autoSubscribe?: boolean` - Auto-subscribe (default: `true`)
- `onMessage?: (message) => void` - Message callback
- `onSubscribed?: () => void` - Subscription callback
- `onUnsubscribed?: () => void` - Unsubscription callback

**Returns:**
- `isSubscribed: Ref<boolean>` - Subscription status
- `lastMessage: Ref<WebSocketMessage | null>` - Last message
- `subscribe()` - Subscribe
- `unsubscribe()` - Unsubscribe

### `usePulseWebSocketMessageHandler(options)`

Message handler registry composable.

**Options:**
- `ws?: UsePulseWebSocketReturn` - WebSocket instance

**Returns:**
- `on(type, handler)` - Register handler
- `off(type, handler)` - Unregister handler
- `once(type, handler)` - Register one-time handler
- `clear(type?)` - Clear handlers

## Configuration

Add WebSocket URL to your `nuxt.config.ts`:

```ts
export default defineNuxtConfig({
  runtimeConfig: {
    public: {
      wsUrl: process.env.WS_URL || 'ws://localhost:8080/ws'
    }
  }
})
```

## Examples

### Project Dashboard

```vue
<script setup lang="ts">
const route = useRoute()
const projectId = computed(() => route.params.projectId as string)

const ws = usePulseWebSocket({
  autoSubscribe: [`project:${projectId.value}`],
})

const handler = usePulseWebSocketMessageHandler({ ws })

// Handle check run updates
handler.on('check_run_completed', (msg) => {
  // Update UI with check run data
  updateCheckStatus(msg.payload)
})

// Handle alerts
handler.on('alert_created', (msg) => {
  // Show alert notification
  showAlert(msg.payload)
})
</script>
```

### Real-time Notifications

```vue
<script setup lang="ts">
const notifications = ref([])

const ws = usePulseWebSocket({
  autoSubscribe: ['notifications'],
})

watch(ws.data, (message) => {
  if (message?.type === 'notification') {
    notifications.value.unshift(message.payload)
  }
})
</script>
```

## Best Practices

1. **Use topic subscriptions** instead of global broadcasts when possible
2. **Clean up subscriptions** in `onUnmounted` hooks
3. **Handle errors** gracefully with error callbacks
4. **Use message handlers** for type-safe message processing
5. **Leverage auto-reconnect** for better reliability
6. **Monitor connection status** to show UI feedback

## Type Safety

All composables are fully typed with TypeScript. Message types can be extended:

```ts
interface CheckRunMessage {
  check_id: string
  status: string
  timestamp: string
}

handler.on<CheckRunMessage>('check_run_completed', (msg) => {
  // msg.payload is typed as CheckRunMessage
  console.log(msg.payload.check_id)
})
```


