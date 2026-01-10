import express from 'express';
import {
  sendNotification,
  getNotifications,
  getNotificationById,
  markAsRead,
} from '../controllers/notificationController.js';

const router = express.Router();

// POST /api/notifications - Send a new notification
router.post('/', sendNotification);

// GET /api/notifications - Get all notifications (with filters)
router.get('/', getNotifications);

// GET /api/notifications/:id - Get a specific notification
router.get('/:id', getNotificationById);

// PATCH /api/notifications/:id/read - Mark notification as read
router.patch('/:id/read', markAsRead);

export default router;
