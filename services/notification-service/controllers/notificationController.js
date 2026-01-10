import { asyncHandler } from '../middleware/asyncHandler.js';
import { AppError } from '../middleware/errorHandler.js';
import * as notificationService from '../services/notificationService.js';

// @desc    Send a new notification
// @route   POST /api/notifications
// @access  Private
export const sendNotification = asyncHandler(async (req, res) => {
  const { userId, type, title, message, data } = req.body;

  // Validation
  if (!userId || !type || !message) {
    throw new AppError('userId, type, and message are required', 400);
  }

  const notification = await notificationService.sendNotification({
    userId,
    type,
    title,
    message,
    data,
  });

  res.status(201).json({
    success: true,
    data: notification,
  });
});

// @desc    Get all notifications
// @route   GET /api/notifications
// @access  Private
export const getNotifications = asyncHandler(async (req, res) => {
  const { userId, isRead, type, limit = 20, offset = 0 } = req.query;

  const notifications = await notificationService.getNotifications({
    userId,
    isRead: isRead === 'true' ? true : isRead === 'false' ? false : undefined,
    type,
    limit: parseInt(limit),
    offset: parseInt(offset),
  });

  res.status(200).json({
    success: true,
    count: notifications.length,
    data: notifications,
  });
});

// @desc    Get notification by ID
// @route   GET /api/notifications/:id
// @access  Private
export const getNotificationById = asyncHandler(async (req, res) => {
  const { id } = req.params;

  const notification = await notificationService.getNotificationById(id);

  if (!notification) {
    throw new AppError('Notification not found', 404);
  }

  res.status(200).json({
    success: true,
    data: notification,
  });
});

// @desc    Mark notification as read
// @route   PATCH /api/notifications/:id/read
// @access  Private
export const markAsRead = asyncHandler(async (req, res) => {
  const { id } = req.params;

  const notification = await notificationService.markAsRead(id);

  if (!notification) {
    throw new AppError('Notification not found', 404);
  }

  res.status(200).json({
    success: true,
    data: notification,
  });
});
