import { SOCKET_EVENTS } from '../utils/constants.js';
import { getUserNotificationRoom } from '../utils/socketHelpers.js';

/**
 * Register notification-related socket event handlers
 * @param {import('socket.io').Server} io - Socket.io server instance
 * @param {import('socket.io').Socket} socket - Client socket
 */
export const registerNotificationHandlers = (io, socket) => {
  const userId = socket.userId;

  /**
   * Join user's notification room
   * This allows the server to send notifications to this specific user
   */
  socket.on(SOCKET_EVENTS.JOIN_NOTIFICATIONS, async () => {
    try {
      const room = getUserNotificationRoom(userId);
      socket.join(room);

      socket.emit(SOCKET_EVENTS.NOTIFICATIONS_JOINED, {
        userId,
        message: 'Successfully joined notifications room',
      });

      console.log(`🔔 User ${userId} joined notifications room`);
    } catch (error) {
      console.error('Error in JOIN_NOTIFICATIONS:', error);
      socket.emit(SOCKET_EVENTS.ERROR, {
        message: 'Failed to join notifications room',
      });
    }
  });

  /**
   * Mark notification as read (via WebSocket)
   * Payload: { notificationId }
   * Note: This is handled via HTTP API primarily, but can be done via socket too
   */
  socket.on(SOCKET_EVENTS.NOTIFICATION_READ, async (payload) => {
    try {
      const { notificationId } = payload;

      if (!notificationId) {
        socket.emit(SOCKET_EVENTS.ERROR, {
          message: 'notificationId is required',
        });
        return;
      }

      // Import notification service dynamically to avoid circular dependencies
      const { markAsRead } = await import('../services/notificationService.js');
      
      const notification = await markAsRead(notificationId, userId);

      if (!notification) {
        socket.emit(SOCKET_EVENTS.ERROR, {
          message: 'Notification not found or access denied',
        });
        return;
      }

      // Send confirmation
      socket.emit(SOCKET_EVENTS.NOTIFICATION_UPDATED, {
        notificationId,
        isRead: true,
      });

      console.log(`✅ Notification ${notificationId} marked as read`);
    } catch (error) {
      console.error('Error in NOTIFICATION_READ:', error);
      socket.emit(SOCKET_EVENTS.ERROR, {
        message: 'Failed to mark notification as read',
      });
    }
  });
};
