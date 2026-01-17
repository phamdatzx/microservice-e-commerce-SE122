import { Server } from 'socket.io';

/**
 * Initialize and configure Socket.io server
 * @param {import('http').Server} httpServer - HTTP server instance
 * @returns {Server} Socket.io server instance
 */
export const initializeSocket = (httpServer) => {
  const io = new Server(httpServer, {
    cors: {
      origin: process.env.CORS_ORIGIN || '*',
      methods: ['GET', 'POST'],
      credentials: true,
    },
    pingTimeout: 60000,
    pingInterval: 25000,
    maxHttpBufferSize: 1e6, // 1MB
    transports: ['websocket', 'polling'],
  });

  console.log('✅ Socket.io server initialized');
  return io;
};

export default initializeSocket;
