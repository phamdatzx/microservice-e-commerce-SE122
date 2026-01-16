import * as notificationService from '../services/notificationService.js';
import { PAGINATION, ERROR_MESSAGES } from '../utils/constants.js';

/**
 * Create a new notification (called by other services)
 * POST /api/notifications
 * Body: { userId, type, title, message, data }
 */
export const createNotification = async (req, res) => {
  try {
    const { userId, type, title, message, data } = req.body;

    if (!userId || !type || !title || !message) {
      return res.status(400).json({
        success: false,
        error: 'userId, type, title, and message are required',
      });
    }

    const notification = await notificationService.createNotification(
      userId,
      type,
      title,
      message,
      data
    );

    // Emit notification via WebSocket if io instance is available
    if (req.app.locals.io) {
      const { emitToNotificationRoom } = await import('../utils/socketHelpers.js');
      const { SOCKET_EVENTS } = await import('../utils/constants.js');
      emitToNotificationRoom(req.app.locals.io, userId, SOCKET_EVENTS.NEW_NOTIFICATION, notification);
    }

    return res.status(201).json({
      success: true,
      data: notification,
    });
  } catch (error) {
    console.error('Error in createNotification:', error);
    return res.status(500).json({
      success: false,
      error: ERROR_MESSAGES.SERVER_ERROR,
    });
  }
};

/**
 * Get user's notifications
 * GET /api/notifications
 * Query params: page, limit, isRead
 */
export const getUserNotifications = async (req, res) => {
  try {
    const userId = req.headers['x-user-id'];
    const page = parseInt(req.query.page) || PAGINATION.DEFAULT_PAGE;
    const limit = parseInt(req.query.limit) || PAGINATION.DEFAULT_LIMIT;
    const isRead = req.query.isRead === 'true' ? true : req.query.isRead === 'false' ? false : null;

    if (!userId) {
      return res.status(401).json({
        success: false,
        error: ERROR_MESSAGES.UNAUTHORIZED,
      });
    }

    const result = await notificationService.getUserNotifications(userId, page, limit, isRead);

    return res.status(200).json({
      success: true,
      data: result.notifications,
      pagination: result.pagination,
    });
  } catch (error) {
    console.error('Error in getUserNotifications:', error);
    return res.status(500).json({
      success: false,
      error: ERROR_MESSAGES.SERVER_ERROR,
    });
  }
};

/**
 * Mark a notification as read
 * PATCH /api/notifications/:notificationId/read
 */
export const markNotificationAsRead = async (req, res) => {
  try {
    const userId = req.headers['x-user-id'];
    const { notificationId } = req.params;

    if (!userId) {
      return res.status(401).json({
        success: false,
        error: ERROR_MESSAGES.UNAUTHORIZED,
      });
    }

    const notification = await notificationService.markAsRead(notificationId, userId);

    if (!notification) {
      return res.status(404).json({
        success: false,
        error: ERROR_MESSAGES.NOT_FOUND,
      });
    }

    // Emit update via WebSocket if io instance is available
    if (req.app.locals.io) {
      const { emitToNotificationRoom } = await import('../utils/socketHelpers.js');
      const { SOCKET_EVENTS } = await import('../utils/constants.js');
      emitToNotificationRoom(req.app.locals.io, userId, SOCKET_EVENTS.NOTIFICATION_UPDATED, {
        notificationId,
        isRead: true,
      });
    }

    return res.status(200).json({
      success: true,
      data: notification,
    });
  } catch (error) {
    console.error('Error in markNotificationAsRead:', error);
    return res.status(500).json({
      success: false,
      error: ERROR_MESSAGES.SERVER_ERROR,
    });
  }
};

/**
 * Mark all notifications as read
 * PATCH /api/notifications/read-all
 */
export const markAllNotificationsAsRead = async (req, res) => {
  try {
    const userId = req.headers['x-user-id'];

    if (!userId) {
      return res.status(401).json({
        success: false,
        error: ERROR_MESSAGES.UNAUTHORIZED,
      });
    }

    const count = await notificationService.markAllAsRead(userId);

    return res.status(200).json({
      success: true,
      data: {
        markedAsRead: count,
      },
    });
  } catch (error) {
    console.error('Error in markAllNotificationsAsRead:', error);
    return res.status(500).json({
      success: false,
      error: ERROR_MESSAGES.SERVER_ERROR,
    });
  }
};

/**
 * Get unread notification count
 * GET /api/notifications/unread-count
 */
export const getUnreadCount = async (req, res) => {
  try {
    const userId = req.headers['x-user-id'];

    if (!userId) {
      return res.status(401).json({
        success: false,
        error: ERROR_MESSAGES.UNAUTHORIZED,
      });
    }

    const count = await notificationService.getUnreadCount(userId);

    return res.status(200).json({
      success: true,
      data: {
        unreadCount: count,
      },
    });
  } catch (error) {
    console.error('Error in getUnreadCount:', error);
    return res.status(500).json({
      success: false,
      error: ERROR_MESSAGES.SERVER_ERROR,
    });
  }
};
