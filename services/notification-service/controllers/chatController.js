import * as chatService from '../services/chatService.js';
import { PAGINATION, ERROR_MESSAGES } from '../utils/constants.js';

/**
 * Get all conversations for a user
 * GET /api/chat/conversations
 */
export const getConversations = async (req, res) => {
  try {
    const userId = req.headers['x-user-id'];
    const page = parseInt(req.query.page) || PAGINATION.DEFAULT_PAGE;
    const limit = parseInt(req.query.limit) || PAGINATION.DEFAULT_LIMIT;

    if (!userId) {
      return res.status(401).json({
        success: false,
        error: ERROR_MESSAGES.UNAUTHORIZED,
      });
    }

    const result = await chatService.getUserConversations(userId, page, limit);

    return res.status(200).json({
      success: true,
      data: result.conversations,
      pagination: result.pagination,
    });
  } catch (error) {
    console.error('Error in getConversations:', error);
    return res.status(500).json({
      success: false,
      error: ERROR_MESSAGES.SERVER_ERROR,
    });
  }
};

/**
 * Get or create a conversation
 * POST /api/chat/conversations
 * Body: { sellerId } or { userId } depending on who's calling
 */
export const getOrCreateConversation = async (req, res) => {
  try {
    const currentUserId = req.headers['x-user-id'];
    const { sellerId, userId: customerId } = req.body;

    if (!currentUserId) {
      return res.status(401).json({
        success: false,
        error: ERROR_MESSAGES.UNAUTHORIZED,
      });
    }

    // Determine userId and sellerId based on who's calling
    let finalUserId, finalSellerId;
    
    if (sellerId) {
      // Current user is customer, creating conversation with seller
      finalUserId = currentUserId;
      finalSellerId = sellerId;
    } else if (customerId) {
      // Current user is seller, creating conversation with customer
      finalUserId = customerId;
      finalSellerId = currentUserId;
    } else {
      return res.status(400).json({
        success: false,
        error: 'Either sellerId or userId is required',
      });
    }

    const conversation = await chatService.findOrCreateConversation(finalUserId, finalSellerId);

    return res.status(200).json({
      success: true,
      data: conversation,
    });
  } catch (error) {
    console.error('Error in getOrCreateConversation:', error);
    return res.status(500).json({
      success: false,
      error: ERROR_MESSAGES.SERVER_ERROR,
    });
  }
};

/**
 * Get conversation by ID
 * GET /api/chat/conversations/:conversationId
 */
export const getConversationById = async (req, res) => {
  try {
    const userId = req.headers['x-user-id'];
    const { conversationId } = req.params;

    if (!userId) {
      return res.status(401).json({
        success: false,
        error: ERROR_MESSAGES.UNAUTHORIZED,
      });
    }

    const conversation = await chatService.getConversationById(conversationId, userId);

    if (!conversation) {
      return res.status(404).json({
        success: false,
        error: ERROR_MESSAGES.NOT_FOUND,
      });
    }

    return res.status(200).json({
      success: true,
      data: conversation,
    });
  } catch (error) {
    console.error('Error in getConversationById:', error);
    return res.status(500).json({
      success: false,
      error: ERROR_MESSAGES.SERVER_ERROR,
    });
  }
};

/**
 * Get messages for a conversation
 * GET /api/chat/conversations/:conversationId/messages
 */
export const getMessages = async (req, res) => {
  try {
    const userId = req.headers['x-user-id'];
    const { conversationId } = req.params;
    const page = parseInt(req.query.page) || PAGINATION.DEFAULT_PAGE;
    const limit = parseInt(req.query.limit) || PAGINATION.DEFAULT_LIMIT;

    if (!userId) {
      return res.status(401).json({
        success: false,
        error: ERROR_MESSAGES.UNAUTHORIZED,
      });
    }

    // Verify user has access to this conversation
    const conversation = await chatService.getConversationById(conversationId, userId);
    if (!conversation) {
      return res.status(403).json({
        success: false,
        error: 'Access denied to this conversation',
      });
    }

    const result = await chatService.getConversationMessages(conversationId, page, limit);

    return res.status(200).json({
      success: true,
      data: result.messages,
      pagination: result.pagination,
    });
  } catch (error) {
    console.error('Error in getMessages:', error);
    return res.status(500).json({
      success: false,
      error: ERROR_MESSAGES.SERVER_ERROR,
    });
  }
};

/**
 * Send a message (HTTP fallback)
 * POST /api/chat/messages
 * Body: { conversationId, content, image }
 */
export const sendMessage = async (req, res) => {
  try {
    const senderId = req.headers['x-user-id'];
    const { conversationId, content, image } = req.body;

    if (!senderId) {
      return res.status(401).json({
        success: false,
        error: ERROR_MESSAGES.UNAUTHORIZED,
      });
    }

    if (!conversationId || !content) {
      return res.status(400).json({
        success: false,
        error: 'conversationId and content are required',
      });
    }

    // Verify user has access to this conversation
    const conversation = await chatService.getConversationById(conversationId, senderId);
    if (!conversation) {
      return res.status(403).json({
        success: false,
        error: 'Access denied to this conversation',
      });
    }

    const message = await chatService.createMessage(senderId, conversationId, content, image);

    return res.status(201).json({
      success: true,
      data: message,
    });
  } catch (error) {
    console.error('Error in sendMessage:', error);
    return res.status(500).json({
      success: false,
      error: ERROR_MESSAGES.SERVER_ERROR,
    });
  }
};

/**
 * Mark messages as read in a conversation
 * PATCH /api/chat/conversations/:conversationId/read
 */
export const markConversationAsRead = async (req, res) => {
  try {
    const userId = req.headers['x-user-id'];
    const { conversationId } = req.params;

    if (!userId) {
      return res.status(401).json({
        success: false,
        error: ERROR_MESSAGES.UNAUTHORIZED,
      });
    }

    // Verify user has access to this conversation
    const conversation = await chatService.getConversationById(conversationId, userId);
    if (!conversation) {
      return res.status(403).json({
        success: false,
        error: 'Access denied to this conversation',
      });
    }

    const count = await chatService.markMessagesAsRead(conversationId, userId);

    return res.status(200).json({
      success: true,
      data: {
        messagesMarkedAsRead: count,
      },
    });
  } catch (error) {
    console.error('Error in markConversationAsRead:', error);
    return res.status(500).json({
      success: false,
      error: ERROR_MESSAGES.SERVER_ERROR,
    });
  }
};
