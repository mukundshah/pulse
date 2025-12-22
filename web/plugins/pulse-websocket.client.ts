/**
 * Global WebSocket plugin for Pulse
 * Provides a singleton WebSocket connection that can be used throughout the app
 */
export default defineNuxtPlugin({
  name: 'pulse-websocket',
  setup() {
    const config = useRuntimeConfig()
    const { token } = useAuth()

    // Get WebSocket URL from config or use default
    const wsUrl = config.public.wsUrl || 'ws://localhost:8080/ws'

    // Create global WebSocket connection
    const ws = usePulseWebSocket({
      url: wsUrl,
      immediate: false, // Don't auto-connect, let components decide
      autoReconnect: {
        retries: 5,
        delay: (retries) => Math.min(1000 * 2 ** (retries - 1), 30000), // Exponential backoff
        onFailed: () => {
          console.error('WebSocket connection failed after 5 retries')
        },
      },
      heartbeat: {
        interval: 30000,
        message: 'ping',
      },
      onAuthSuccess: () => {
        console.log('WebSocket authenticated')
      },
      onAuthError: (error) => {
        console.error('WebSocket authentication error:', error)
      },
      onError: (error) => {
        console.error('WebSocket error:', error)
      },
    })

    // Auto-connect when user is authenticated
    watch([token, ws.isOpen], ([newToken, isOpen]) => {
      if (newToken && !isOpen && !ws.isAuthenticated.value) {
        ws.open()
      }
    }, { immediate: true })

    // Disconnect on logout
    watch(() => ws.isAuthenticated.value, (isAuth) => {
      if (!isAuth && ws.isOpen.value) {
        ws.close()
      }
    })

    return {
      provide: {
        ws,
      },
    }
  },
})


