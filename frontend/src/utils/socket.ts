import { io, Socket } from 'socket.io-client'

const getSocketUrl = () => {
  if (import.meta.env.VITE_NOTIFICATION_SERVICE_URL) {
    return import.meta.env.VITE_NOTIFICATION_SERVICE_URL
  }
  const apiUrl = import.meta.env.VITE_BE_API_URL
  if (apiUrl) {
    return apiUrl.split('/api')[0]
  }
  return 'http://localhost:81'
}

const SOCKET_URL = getSocketUrl()

export const SOCKET_EVENTS = {
  CONNECTION: 'connection',
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

  JOIN_NOTIFICATIONS: 'join-notifications',
  NOTIFICATIONS_JOINED: 'notifications-joined',
  NEW_NOTIFICATION: 'new-notification',
  NOTIFICATION_READ: 'notification-read',
  NOTIFICATION_UPDATED: 'notification-updated',
}

class SocketService {
  private socket: Socket | null = null

  connect(token: string) {
    if (this.socket) {
      if (!this.socket.connected) {
        const socketAny = this.socket as any
        if (socketAny.auth) {
          socketAny.auth.token = token
        }
        this.socket.connect()
      }
      return
    }

    this.socket = io(SOCKET_URL, {
      auth: { token },
      transports: ['websocket', 'polling'],
    })

    this.socket.on('connect', () => {
      console.log('Connected to socket server')
      this.socket?.emit(SOCKET_EVENTS.JOIN_NOTIFICATIONS, {})
    })

    this.socket.on('connect_error', (error) => {
      console.error('Socket connection error:', error)
    })

    this.socket.on('disconnect', (reason) => {
      console.log('Disconnected from socket server:', reason)
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

  off(event: string, callback?: (...args: any[]) => void) {
    this.socket?.off(event, callback)
  }

  emit(event: string, data: any) {
    this.socket?.emit(event, data)
  }

  get isConnected() {
    return this.socket?.connected || false
  }
}

export const socketService = new SocketService()
