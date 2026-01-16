import { SOCKET_EVENTS } from '../utils/constants.js';
import { registerUserSocket, unregisterUserSocket } from '../utils/socketHelpers.js';
import { registerChatHandlers } from './chatSocket.js';
import { registerNotificationHandlers } from './notificationSocket.js';

/**
 * Initialize all socket event handlers
 * @param {import('socket.io').Server} io - Socket.io server instance
 */
export const initializeSocketHandlers = (io) => {
  io.on(SOCKET_EVENTS.CONNECTION, (socket) => {
    const userId = socket.userId;
    console.log(`🔌 User ${userId} connected (socket ${socket.id})`);

    // Register user socket mapping
    registerUserSocket(userId, socket.id);

    // Register event handlers
    registerChatHandlers(io, socket);
    registerNotificationHandlers(io, socket);

    // Handle disconnection
    socket.on(SOCKET_EVENTS.DISCONNECT, () => {
      console.log(`🔌 User ${userId} disconnected (socket ${socket.id})`);
      unregisterUserSocket(userId);
    });

    // Handle errors
    socket.on(SOCKET_EVENTS.ERROR, (error) => {
      console.error(`❌ Socket error for user ${userId}:`, error);
    });
  });

  console.log('✅ Socket event handlers initialized');
};

export default initializeSocketHandlers;
