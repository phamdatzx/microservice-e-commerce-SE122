import express from 'express';
import * as notificationController from '../controllers/notificationController.js';

const router = express.Router();

// Create notification (for other services to call)
router.post('/', notificationController.createNotification);

// Get user's notifications
router.get('/', notificationController.getUserNotifications);

// Get unread count (specific route before :notificationId)
router.get('/unread-count', notificationController.getUnreadCount);

// Mark all as read
router.patch('/read-all', notificationController.markAllNotificationsAsRead);

// Mark specific notification as read
router.patch('/:notificationId/read', notificationController.markNotificationAsRead);

export default router;
