/**
 * WebSocket authentication middleware
 * Validates user authentication and attaches userId to socket
 */
export const socketAuth = (socket, next) => {
  try {
    // Get userId from handshake query or auth header
    const userId = socket.handshake.query.userId || socket.handshake.auth?.userId;

    if (!userId) {
      return next(new Error('Authentication error: userId is required'));
    }

    // Attach userId to socket for use in event handlers
    socket.userId = userId;
    
    console.log(`✅ User ${userId} authenticated on socket ${socket.id}`);
    next();
  } catch (error) {
    console.error('Socket authentication error:', error);
    next(new Error('Authentication failed'));
  }
};

export default socketAuth;
