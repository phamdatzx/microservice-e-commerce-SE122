import { ROOM_PREFIXES } from './constants.js';

// Store user socket mappings (userId -> socketId)
const userSockets = new Map();

/**
 * Register a user's socket connection
 * @param {string} userId - User ID
 * @param {string} socketId - Socket ID
 */
export const registerUserSocket = (userId, socketId) => {
  userSockets.set(userId, socketId);
  console.log(`📌 Registered socket ${socketId} for user ${userId}`);
};

/**
 * Unregister a user's socket connection
 * @param {string} userId - User ID
 */
export const unregisterUserSocket = (userId) => {
  const socketId = userSockets.get(userId);
  userSockets.delete(userId);
  console.log(`📌 Unregistered socket ${socketId} for user ${userId}`);
};

/**
 * Get socket ID for a specific user
 * @param {string} userId - User ID
 * @returns {string|undefined} Socket ID or undefined if not connected
 */
export const getUserSocketId = (userId) => {
  return userSockets.get(userId);
};

/**
 * Get conversation room name
 * @param {string} conversationId - Conversation ID
 * @returns {string} Room name
 */
export const getConversationRoom = (conversationId) => {
  return `${ROOM_PREFIXES.CONVERSATION}${conversationId}`;
};

/**
 * Get user notification room name
 * @param {string} userId - User ID
 * @returns {string} Room name
 */
export const getUserNotificationRoom = (userId) => {
  return `${ROOM_PREFIXES.USER_NOTIFICATIONS}${userId}`;
};

/**
 * Emit event to a specific user
 * @param {import('socket.io').Server} io - Socket.io server instance
 * @param {string} userId - User ID
 * @param {string} event - Event name
 * @param {any} data - Data to send
 */
export const emitToUser = (io, userId, event, data) => {
  const socketId = getUserSocketId(userId);
  if (socketId) {
    io.to(socketId).emit(event, data);
    console.log(`📤 Emitted ${event} to user ${userId} (socket ${socketId})`);
  } else {
    console.log(`⚠️ User ${userId} is not connected`);
  }
};

/**
 * Emit event to a conversation room
 * @param {import('socket.io').Server} io - Socket.io server instance
 * @param {string} conversationId - Conversation ID
 * @param {string} event - Event name
 * @param {any} data - Data to send
 */
export const emitToConversation = (io, conversationId, event, data) => {
  const room = getConversationRoom(conversationId);
  io.to(room).emit(event, data);
  console.log(`📤 Emitted ${event} to conversation ${conversationId}`);
};

/**
 * Emit event to user's notification room
 * @param {import('socket.io').Server} io - Socket.io server instance
 * @param {string} userId - User ID
 * @param {string} event - Event name
 * @param {any} data - Data to send
 */
export const emitToNotificationRoom = (io, userId, event, data) => {
  const room = getUserNotificationRoom(userId);
  io.to(room).emit(event, data);
  console.log(`📤 Emitted ${event} to notification room for user ${userId}`);
};
