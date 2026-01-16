import express from 'express';
import notificationRoutes from './notificationRoutes.js';
import chatRoutes from './chatRoutes.js';

const router = express.Router();

// Mount routes
router.use('/notifications', notificationRoutes);
router.use('/chat', chatRoutes);

export default router;
