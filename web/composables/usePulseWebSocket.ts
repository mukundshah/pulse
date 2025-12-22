import type { UseWebSocketReturn } from '@vueuse/core'
import type { Ref } from 'vue'
import { useWebSocket } from '@vueuse/core'
import { useAuth } from '~/modules/auth/runtime/composable'

/**
 * WebSocket message types
 */
export type WebSocketMessageType
  = 'auth'
    | 'auth_success'
    | 'subscribe'
    | 'unsubscribe'
    | 'subscribed'
    | 'unsubscribed'
    | 'error'
    | string // Allow custom message types

/**
 * WebSocket message structure
 */
export interface WebSocketMessage<T = unknown> {
  type: WebSocketMessageType
  payload: T
}

/**
 * WebSocket connection status
 */
export type WebSocketStatus = 'CONNECTING' | 'OPEN' | 'CLOSING' | 'CLOSED'

/**
 * Message handler function
 */
export type MessageHandler<T = unknown> = (message: WebSocketMessage<T>) => void

/**
 * Options for usePulseWebSocket
 */
export interface UsePulseWebSocketOptions {
  /**
   * WebSocket URL (defaults to ws://localhost:8080/ws)
   */
  url?: string | Ref<string>

  /**
   * Auto-connect on mount (default: true)
   */
  immediate?: boolean

  /**
   * Auto-reconnect on disconnect (default: true)
   */
  autoReconnect?: boolean | {
    retries?: number
    delay?: number | ((retries: number) => number)
    onFailed?: () => void
  }

  /**
   * Heartbeat interval in ms (default: 30000)
   * Set to false to disable
   */
  heartbeat?: boolean | {
    interval?: number
    message?: string
  }

  /**
   * Topics to auto-subscribe to on connection
   */
  autoSubscribe?: string[]

  /**
   * Callback for authentication success
   */
  onAuthSuccess?: () => void

  /**
   * Callback for authentication error
   */
  onAuthError?: (error: string) => void

  /**
   * Callback for subscription success
   */
  onSubscribed?: (topic: string) => void

  /**
   * Callback for unsubscription success
   */
  onUnsubscribed?: (topic: string) => void

  /**
   * Callback for errors
   */
  onError?: (error: string) => void

  /**
   * Callback for any message
   */
  onMessage?: <T = unknown>(message: WebSocketMessage<T>) => void
}

/**
 * Return type for usePulseWebSocket
 */
export interface UsePulseWebSocketReturn {
  /**
   * WebSocket connection status
   */
  status: Ref<WebSocketStatus>

  /**
   * Last received message
   */
  data: Ref<WebSocketMessage | null>

  /**
   * Whether the connection is open
   */
  isOpen: Ref<boolean>

  /**
   * Whether the client is authenticated
   */
  isAuthenticated: Ref<boolean>

  /**
   * Set of subscribed topics
   */
  subscriptions: Ref<Set<string>>

  /**
   * Open the WebSocket connection
   */
  open: () => void

  /**
   * Close the WebSocket connection
   */
  close: () => void

  /**
   * Send a message
   */
  send: (message: WebSocketMessage) => void

  /**
   * Authenticate with token
   */
  authenticate: (token: string) => void

  /**
   * Subscribe to a topic
   */
  subscribe: (topic: string) => void

  /**
   * Unsubscribe from a topic
   */
  unsubscribe: (topic: string) => void

  /**
   * Subscribe to multiple topics
   */
  subscribeMany: (topics: string[]) => void

  /**
   * Unsubscribe from multiple topics
   */
  unsubscribeMany: (topics: string[]) => void

  /**
   * Check if subscribed to a topic
   */
  isSubscribed: (topic: string) => boolean

  /**
   * Clear all subscriptions
   */
  clearSubscriptions: () => void

  /**
   * Register a handler for a specific message type
   */
  on: <T = unknown>(type: string, handler: MessageHandler<T>) => () => void

  /**
   * Unregister a handler
   */
  off: (type: string, handler: MessageHandler) => void

  /**
   * Register a one-time handler
   */
  once: <T = unknown>(type: string, handler: MessageHandler<T>) => void

  /**
   * Clear all handlers for a type (or all handlers if no type specified)
   */
  clearHandlers: (type?: string) => void
}

/**
 * Composable for Pulse WebSocket connection with authentication, pub/sub, and message handling
 *
 * @example
 * ```ts
 * const ws = usePulseWebSocket({
 *   autoSubscribe: ['project:123'],
 *   onMessage: (msg) => {
 *     console.log('Received:', msg.type, msg.payload)
 *   }
 * })
 *
 * // Register handlers
 * ws.on('check_run_completed', (msg) => {
 *   console.log('Check run:', msg.payload)
 * })
 *
 * // Subscribe to topics
 * ws.subscribe('project:456')
 * ```
 */
export function usePulseWebSocket(options: UsePulseWebSocketOptions = {}): UsePulseWebSocketReturn {
  const {
    url = 'ws://localhost:8080/ws',
    immediate = true,
    autoReconnect = true,
    heartbeat = true,
    autoSubscribe = [],
    onAuthSuccess,
    onAuthError,
    onSubscribed,
    onUnsubscribed,
    onError,
    onMessage,
  } = options

  // Get auth token from auth module
  const { token } = useAuth()

  // Reactive state
  const isAuthenticated = ref(false)
  const subscriptions = ref<Set<string>>(new Set())
  const messageHandlers = new Map<string, Set<MessageHandler>>()

  // Configure heartbeat
  const heartbeatConfig = heartbeat === true
    ? { interval: 30000, message: 'ping' }
    : heartbeat === false
      ? false
      : { interval: heartbeat.interval ?? 30000, message: heartbeat.message ?? 'ping' }

  // Configure auto-reconnect
  const reconnectConfig = autoReconnect === true
    ? { retries: 5, delay: (retries: number) => Math.min(1000 * 2 ** (retries - 1), 30000) }
    : autoReconnect === false
      ? false
      : {
          retries: autoReconnect.retries ?? 5,
          delay: typeof autoReconnect.delay === 'function'
            ? autoReconnect.delay
            : (autoReconnect.delay ?? 1000),
          onFailed: autoReconnect.onFailed,
        }

  // Create WebSocket connection using VueUse
  const ws = useWebSocket(url, {
    immediate,
    autoReconnect: reconnectConfig as any, // Type assertion needed due to VueUse type complexity
    heartbeat: heartbeatConfig,
    onConnected: () => {
      // Auto-authenticate if token is available
      const authToken = token
      if (authToken) {
        authenticate(authToken)
      }
    },
    onDisconnected: () => {
      isAuthenticated.value = false
      subscriptions.value.clear()
    },
  }) as UseWebSocketReturn<WebSocketMessage>

  // Process incoming messages
  watch(ws.data, (message) => {
    if (!message) return

    try {
      const parsed = typeof message === 'string' ? JSON.parse(message) : message

      // Handle system messages
      switch (parsed.type) {
        case 'auth_success': {
          isAuthenticated.value = true
          onAuthSuccess?.()

          // Auto-subscribe to topics after authentication
          if (autoSubscribe.length > 0 && isAuthenticated.value) {
            nextTick(() => {
              subscribeMany(autoSubscribe)
            })
          }
          break
        }

        case 'error': {
          const errorMsg = parsed.payload?.message || 'Unknown error'
          onError?.(errorMsg)
          if (errorMsg.includes('auth') || errorMsg.includes('token')) {
            onAuthError?.(errorMsg)
          }
          break
        }

        case 'subscribed': {
          const subTopic = parsed.payload?.topic
          if (subTopic) {
            subscriptions.value.add(subTopic)
            onSubscribed?.(subTopic)
          }
          break
        }

        case 'unsubscribed': {
          const unsubTopic = parsed.payload?.topic
          if (unsubTopic) {
            subscriptions.value.delete(unsubTopic)
            onUnsubscribed?.(unsubTopic)
          }
          break
        }
      }

      // Call registered message handlers
      const typeHandlers = messageHandlers.get(parsed.type)
      if (typeHandlers) {
        typeHandlers.forEach((handler) => {
          try {
            handler(parsed)
          } catch (error) {
            console.error(`Error in message handler for type "${parsed.type}":`, error)
          }
        })
      }

      // Call wildcard handlers (*)
      const wildcardHandlers = messageHandlers.get('*')
      if (wildcardHandlers) {
        wildcardHandlers.forEach((handler) => {
          try {
            handler(parsed)
          } catch (error) {
            console.error('Error in wildcard message handler:', error)
          }
        })
      }

      // Call custom message handler
      onMessage?.(parsed)
    } catch (error) {
      console.error('Failed to parse WebSocket message:', error)
    }
  })

  /**
   * Authenticate with JWT token
   */
  function authenticate(authToken: string) {
    if (!ws.status.value || ws.status.value !== 'OPEN') {
      console.warn('WebSocket is not open. Cannot authenticate.')
      return
    }

    ws.send(JSON.stringify({
      type: 'auth',
      token: authToken,
    }))
  }

  /**
   * Subscribe to a topic
   */
  function subscribe(topic: string) {
    if (!isAuthenticated.value) {
      console.warn('Not authenticated. Cannot subscribe to topic:', topic)
      return
    }

    if (!ws.status.value || ws.status.value !== 'OPEN') {
      console.warn('WebSocket is not open. Cannot subscribe to topic:', topic)
      return
    }

    ws.send(JSON.stringify({
      type: 'subscribe',
      topic,
    }))
  }

  /**
   * Unsubscribe from a topic
   */
  function unsubscribe(topic: string) {
    if (!ws.status.value || ws.status.value !== 'OPEN') {
      console.warn('WebSocket is not open. Cannot unsubscribe from topic:', topic)
      return
    }

    ws.send(JSON.stringify({
      type: 'unsubscribe',
      topic,
    }))
  }

  /**
   * Subscribe to multiple topics
   */
  function subscribeMany(topics: string[]) {
    topics.forEach(topic => subscribe(topic))
  }

  /**
   * Unsubscribe from multiple topics
   */
  function unsubscribeMany(topics: string[]) {
    topics.forEach(topic => unsubscribe(topic))
  }

  /**
   * Check if subscribed to a topic
   */
  function isSubscribed(topic: string): boolean {
    return subscriptions.value.has(topic)
  }

  /**
   * Clear all subscriptions
   */
  function clearSubscriptions() {
    const topics = Array.from(subscriptions.value)
    unsubscribeMany(topics)
  }

  /**
   * Send a custom message
   */
  function send(message: WebSocketMessage) {
    if (!ws.status.value || ws.status.value !== 'OPEN') {
      console.warn('WebSocket is not open. Cannot send message.')
      return
    }

    ws.send(JSON.stringify(message))
  }

  /**
   * Register a handler for a specific message type
   */
  function on<T = unknown>(type: string, handler: MessageHandler<T>): () => void {
    if (!messageHandlers.has(type)) {
      messageHandlers.set(type, new Set())
    }

    const handlerWrapper: MessageHandler = (message) => {
      handler(message as WebSocketMessage<T>)
    }
    messageHandlers.get(type)!.add(handlerWrapper)

    // Return unsubscribe function
    return () => off(type, handlerWrapper)
  }

  /**
   * Unregister a handler
   */
  function off(type: string, handler: MessageHandler): void {
    const typeHandlers = messageHandlers.get(type)
    if (typeHandlers) {
      typeHandlers.delete(handler)
      if (typeHandlers.size === 0) {
        messageHandlers.delete(type)
      }
    }
  }

  /**
   * Register a one-time handler
   */
  function once<T = unknown>(type: string, handler: MessageHandler<T>): void {
    const wrappedHandler: MessageHandler<T> = (message) => {
      handler(message)
      off(type, wrappedHandler as MessageHandler)
    }
    on(type, wrappedHandler)
  }

  /**
   * Clear all handlers for a type (or all handlers if no type specified)
   */
  function clearHandlers(type?: string): void {
    if (type) {
      messageHandlers.delete(type)
    } else {
      messageHandlers.clear()
    }
  }

  // Watch for token changes and re-authenticate if needed
  // Token is a getter, so we need to watch it via computed
  const tokenRef = computed(() => token)
  watch(tokenRef, (newToken) => {
    if (newToken && ws.status.value === 'OPEN' && !isAuthenticated.value) {
      authenticate(newToken)
    }
  })

  // Cleanup on unmount
  onUnmounted(() => {
    clearSubscriptions()
    clearHandlers()
    ws.close()
  })

  return {
    status: ws.status,
    data: ws.data as Ref<WebSocketMessage | null>,
    isOpen: computed(() => ws.status.value === 'OPEN'),
    isAuthenticated: readonly(isAuthenticated),
    subscriptions,
    open: ws.open,
    close: ws.close,
    send,
    authenticate,
    subscribe,
    unsubscribe,
    subscribeMany,
    unsubscribeMany,
    isSubscribed,
    clearSubscriptions,
    on,
    off,
    once,
    clearHandlers,
  }
}
