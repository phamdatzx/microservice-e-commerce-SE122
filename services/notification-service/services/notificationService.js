import { Notification } from '../models/index.js';
import { PAGINATION, NOTIFICATION_TYPES } from '../utils/constants.js';

/**
 * Create a new notification
 * @param {string} userId - User ID
 * @param {string} type - Notification type
 * @param {string} title - Notification title
 * @param {string} message - Notification message
 * @param {Object} data - Additional metadata
 * @returns {Promise<Object>} Created notification
 */
export const createNotification = async (userId, type, title, message, data = {}) => {
  try {
    // Validate notification type
    if (!Object.values(NOTIFICATION_TYPES).includes(type)) {
      throw new Error(`Invalid notification type: ${type}`);
    }

    const notification = await Notification.create({
      userId,
      type,
      title,
      message,
      data,
      isRead: false,
    });

    console.log(`✅ Created notification for user ${userId}: ${title}`);
    return notification;
  } catch (error) {
    console.error('Error in createNotification:', error);
    throw error;
  }
};

/**
 * Get user's notifications with pagination and optional filtering
 * @param {string} userId - User ID
 * @param {number} page - Page number
 * @param {number} limit - Items per page
 * @param {boolean} isRead - Filter by read status (optional)
 * @returns {Promise<Object>} Paginated notifications
 */
export const getUserNotifications = async (userId, page = PAGINATION.DEFAULT_PAGE, limit = PAGINATION.DEFAULT_LIMIT, isRead = null) => {
  try {
    const skip = (page - 1) * limit;
    const effectiveLimit = Math.min(limit, PAGINATION.MAX_LIMIT);

    // Build query
    const query = { userId };
    if (isRead !== null) {
      query.isRead = isRead;
    }

    const [notifications, total] = await Promise.all([
      Notification.find(query)
        .sort({ createdAt: -1 })
        .skip(skip)
        .limit(effectiveLimit)
        .lean(),
      Notification.countDocuments(query),
    ]);

    return {
      notifications,
      pagination: {
        page,
        limit: effectiveLimit,
        total,
        pages: Math.ceil(total / effectiveLimit),
      },
    };
  } catch (error) {
    console.error('Error in getUserNotifications:', error);
    throw error;
  }
};

/**
 * Mark a notification as read
 * @param {string} notificationId - Notification ID
 * @param {string} userId - User ID (for authorization)
 * @returns {Promise<Object|null>} Updated notification or null
 */
export const markAsRead = async (notificationId, userId) => {
  try {
    const notification = await Notification.findOneAndUpdate(
      {
        _id: notificationId,
        userId,
      },
      {
        $set: { isRead: true },
      },
      { new: true }
    );

    if (notification) {
      console.log(`✅ Marked notification ${notificationId} as read`);
    }

    return notification;
  } catch (error) {
    console.error('Error in markAsRead:', error);
    throw error;
  }
};

/**
 * Mark all notifications as read for a user
 * @param {string} userId - User ID
 * @returns {Promise<number>} Number of notifications marked as read
 */
export const markAllAsRead = async (userId) => {
  try {
    const result = await Notification.updateMany(
      {
        userId,
        isRead: false,
      },
      {
        $set: { isRead: true },
      }
    );

    console.log(`✅ Marked ${result.modifiedCount} notifications as read for user ${userId}`);
    return result.modifiedCount;
  } catch (error) {
    console.error('Error in markAllAsRead:', error);
    throw error;
  }
};

/**
 * Get unread notification count for a user
 * @param {string} userId - User ID
 * @returns {Promise<number>} Unread count
 */
export const getUnreadCount = async (userId) => {
  try {
    const count = await Notification.countDocuments({
      userId,
      isRead: false,
    });

    return count;
  } catch (error) {
    console.error('Error in getUnreadCount:', error);
    throw error;
  }
};

/**
 * Delete old notifications
 * @param {number} daysOld - Delete notifications older than this many days
 * @returns {Promise<number>} Number of notifications deleted
 */
export const deleteOldNotifications = async (daysOld = 90) => {
  try {
    const cutoffDate = new Date();
    cutoffDate.setDate(cutoffDate.getDate() - daysOld);

    const result = await Notification.deleteMany({
      createdAt: { $lt: cutoffDate },
      isRead: true,
    });

    console.log(`🗑️ Deleted ${result.deletedCount} old notifications`);
    return result.deletedCount;
  } catch (error) {
    console.error('Error in deleteOldNotifications:', error);
    throw error;
  }
};
