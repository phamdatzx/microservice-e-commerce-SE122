import express from 'express';
import notificationRoutes from './notificationRoutes.js';
import chatRoutes from './chatRoutes.js';

const router = express.Router();

// Mount routes
router.use('/notification', notificationRoutes);
router.use('/notification/chat', chatRoutes);

export default router;
