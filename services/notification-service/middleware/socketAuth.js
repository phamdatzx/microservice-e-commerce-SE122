import { userServiceClient } from '../clients/userServiceClient.js';

/**
 * WebSocket authentication middleware
 * Validates user authentication and attaches userId to socket
 */
export const socketAuth = async (socket, next) => {
  try {
    // Get token from handshake query or auth header
    const token = socket.handshake.query.token || socket.handshake.auth?.token;
    
    if (!token) {
      return next(new Error('Authentication error: token is required'));
    }

    // Verify token with user-service
    const userInfo = await userServiceClient.verifyToken(token);
    
    if (!userInfo.success || !userInfo.userId) {
      return next(new Error('Authentication error: invalid token'));
    }

    // Attach user info to socket for use in event handlers
    socket.userId = userInfo.userId;
    socket.username = userInfo.username;
    socket.userRole = userInfo.userRole;
    socket.token = token;
    
    console.log(`✅ User ${userInfo.userId} (${userInfo.username}) authenticated on socket ${socket.id}`);
    next();
  } catch (error) {
    console.error('Socket authentication error:', error.message);
    next(new Error(`Authentication failed: ${error.message}`));
  }
};

export default socketAuth;
