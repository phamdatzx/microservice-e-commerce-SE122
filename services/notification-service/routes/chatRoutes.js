import express from 'express';
import multer from 'multer';
import * as chatController from '../controllers/chatController.js';

const router = express.Router();

// Configure multer for image uploads
const upload = multer({
  storage: multer.memoryStorage(),
  limits: {
    fileSize: 5 * 1024 * 1024, // 5MB max file size
  },
  fileFilter: (req, file, cb) => {
    // Accept images only
    if (!file.mimetype.startsWith('image/')) {
      return cb(new Error('Only image files are allowed'), false);
    }
    cb(null, true);
  },
});

// Conversation routes
router.get('/conversations', chatController.getConversations);
router.post('/conversations', chatController.getOrCreateConversation);
router.get('/conversations/:conversationId', chatController.getConversationById);
router.get('/conversations/:conversationId/messages', chatController.getMessages);
router.patch('/conversations/:conversationId/read', chatController.markConversationAsRead);

// Message routes
router.post('/messages', chatController.sendMessage);
router.post('/messages/with-image', upload.single('image'), chatController.sendMessageWithImage);

export default router;
