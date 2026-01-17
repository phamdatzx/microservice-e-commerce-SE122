// Socket event names
export const SOCKET_EVENTS = {
  // Connection events
  CONNECTION: 'connection',
  DISCONNECT: 'disconnect',
  ERROR: 'error',

  // Chat events
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

  // Notification events
  JOIN_NOTIFICATIONS: 'join-notifications',
  NOTIFICATIONS_JOINED: 'notifications-joined',
  NEW_NOTIFICATION: 'new-notification',
  NOTIFICATION_READ: 'notification-read',
  NOTIFICATION_UPDATED: 'notification-updated',
};

// Notification types
export const NOTIFICATION_TYPES = {
  ORDER: 'order',
  PAYMENT: 'payment',
  PRODUCT: 'product',
  SYSTEM: 'system',
  PROMOTION: 'promotion',
};

// Room prefixes
export const ROOM_PREFIXES = {
  CONVERSATION: 'conversation_',
  USER_NOTIFICATIONS: 'user_notifications_',
};

// Pagination defaults
export const PAGINATION = {
  DEFAULT_PAGE: 1,
  DEFAULT_LIMIT: 20,
  MAX_LIMIT: 100,
};

// Error messages
export const ERROR_MESSAGES = {
  UNAUTHORIZED: 'Unauthorized',
  INVALID_CONVERSATION: 'Invalid conversation ID',
  INVALID_MESSAGE: 'Invalid message data',
  INVALID_NOTIFICATION: 'Invalid notification data',
  NOT_FOUND: 'Resource not found',
  SERVER_ERROR: 'Internal server error',
};
