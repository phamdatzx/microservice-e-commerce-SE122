import { io, Socket } from 'socket.io-client'

export const SOCKET_EVENTS = {
  CONNECT: 'connect',
  DISCONNECT: 'disconnect',
  ERROR: 'error',
  JOIN_CONVERSATION: 'join-conversation',
  LEAVE_CONVERSATION: 'leave-conversation',
  CONVERSATION_JOINED: 'conversation-joined',
  CONVERSATION_LEFT: 'conversation-left',
  SEND_MESSAGE: 'send-message',
  NEW_MESSAGE: 'new-message',
  MESSAGE_SENT: 'message-sent',
  TYPING: 'typing',
  USER_TYPING: 'user-typing',
  MESSAGE_READ: 'message-read',
  MESSAGES_UPDATED: 'messages-updated',
  NEW_NOTIFICATION: 'new-notification',
  JOIN_NOTIFICATIONS: 'join-notifications',
  NOTIFICATIONS_JOINED: 'notifications-joined',
  NOTIFICATION_READ: 'notification-read',
  NOTIFICATION_UPDATED: 'notification-updated',
}

class SocketService {
  private socket: Socket | null = null

  connect(token: string) {
    if (this.socket?.connected) return

    const apiUrl = import.meta.env.VITE_BE_API_URL || 'http://localhost:8080'
    let socketUrl = import.meta.env.VITE_NOTIFICATION_SERVICE_URL
    if (!socketUrl) {
      try {
        const url = new URL(apiUrl)
        socketUrl = `${url.protocol}//${url.hostname}:81`
      } catch (e) {
        socketUrl = apiUrl.replace(/:\d+.*$/, ':81')
      }
    }

    console.log('Connecting to socket at:', socketUrl)

    this.socket = io(socketUrl, {
      auth: { token },
      transports: ['websocket', 'polling'],
    })

    this.socket.on(SOCKET_EVENTS.CONNECT, () => {
      console.log('✅ Connected to notification socket')
    })

    this.socket.on('connect_error', (error) => {
      console.error('❌ Socket connection error:', error)
    })

    this.socket.on(SOCKET_EVENTS.DISCONNECT, () => {
      console.log('❌ Disconnected from notification socket')
    })
  }

  disconnect() {
    if (this.socket) {
      this.socket.disconnect()
      this.socket = null
    }
  }

  on(event: string, callback: (...args: any[]) => void) {
    this.socket?.on(event, callback)
  }

  off(event: string, callback: (...args: any[]) => void) {
    this.socket?.off(event, callback)
  }

  emit(event: string, data: any) {
    this.socket?.emit(event, data)
  }
}

export const socketService = new SocketService()
