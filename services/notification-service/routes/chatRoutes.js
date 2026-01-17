import express from 'express';
import * as chatController from '../controllers/chatController.js';

const router = express.Router();

// Conversation routes
router.get('/conversations', chatController.getConversations);
router.post('/conversations', chatController.getOrCreateConversation);
router.get('/conversations/:conversationId', chatController.getConversationById);
router.get('/conversations/:conversationId/messages', chatController.getMessages);
router.patch('/conversations/:conversationId/read', chatController.markConversationAsRead);

// Message routes
router.post('/messages', chatController.sendMessage);

export default router;
