import { Conversation, Message } from '../models/index.js';
import { PAGINATION } from '../utils/constants.js';
import { userServiceClient } from '../clients/userServiceClient.js';

/**
 * Find or create a conversation between a user and a seller
 * @param {string} userId - Customer user ID
 * @param {string} sellerId - Seller user ID
 * @returns {Promise<Object>} Conversation object with user details
 */
export const findOrCreateConversation = async (userId, sellerId) => {
  try {
    // Try to find existing conversation
    let conversation = await Conversation.findOne({
      userId,
      sellerId,
    }).lean();

    // Create new conversation if it doesn't exist
    if (!conversation) {
      conversation = await Conversation.create({
        userId,
        sellerId,
        lastMessage: '',
        lastUpdated: new Date(),
        unreadCount: 0,
      });
      conversation = conversation.toObject();
      console.log(`✅ Created new conversation between user ${userId} and seller ${sellerId}`);
    }

    // Fetch customer and seller information
    const [customerInfo, sellerInfo] = await Promise.all([
      userServiceClient.getUserById(userId).catch(err => {
        console.error(`Failed to fetch customer ${userId}:`, err.message);
        return null;
      }),
      userServiceClient.getUserById(sellerId).catch(err => {
        console.error(`Failed to fetch seller ${sellerId}:`, err.message);
        return null;
      }),
    ]);

    // Return conversation with user details
    return {
      ...conversation,
      customer: customerInfo,
      seller: sellerInfo,
    };
  } catch (error) {
    console.error('Error in findOrCreateConversation:', error);
    throw error;
  }
};

/**
 * Get all conversations for a user with pagination
 * @param {string} userId - User ID (can be customer or seller)
 * @param {number} page - Page number
 * @param {number} limit - Items per page
 * @returns {Promise<Object>} Paginated conversations with user details
 */
export const getUserConversations = async (userId, page = PAGINATION.DEFAULT_PAGE, limit = PAGINATION.DEFAULT_LIMIT) => {
  try {
    const skip = (page - 1) * limit;
    const effectiveLimit = Math.min(limit, PAGINATION.MAX_LIMIT);

    // Find conversations where user is either the customer or seller
    const [conversations, total] = await Promise.all([
      Conversation.find({
        $or: [{ userId }, { sellerId: userId }],
      })
        .sort({ lastUpdated: -1 })
        .skip(skip)
        .limit(effectiveLimit)
        .lean(),
      Conversation.countDocuments({
        $or: [{ userId }, { sellerId: userId }],
      }),
    ]);

    // Fetch user info for all conversations in parallel
    const enrichedConversations = await Promise.all(
      conversations.map(async (conversation) => {
        const [customerInfo, sellerInfo] = await Promise.all([
          userServiceClient.getUserById(conversation.userId).catch(err => {
            console.error(`Failed to fetch customer ${conversation.userId}:`, err.message);
            return null;
          }),
          userServiceClient.getUserById(conversation.sellerId).catch(err => {
            console.error(`Failed to fetch seller ${conversation.sellerId}:`, err.message);
            return null;
          }),
        ]);

        return {
          ...conversation,
          customer: customerInfo,
          seller: sellerInfo,
        };
      })
    );

    return {
      conversations: enrichedConversations,
      pagination: {
        page,
        limit: effectiveLimit,
        total,
        pages: Math.ceil(total / effectiveLimit),
      },
    };
  } catch (error) {
    console.error('Error in getUserConversations:', error);
    throw error;
  }
};

/**
 * Get conversation by ID
 * @param {string} conversationId - Conversation ID
 * @param {string} userId - User ID (for authorization)
 * @returns {Promise<Object|null>} Conversation object or null
 */
export const getConversationById = async (conversationId, userId) => {
  try {
    const conversation = await Conversation.findOne({
      _id: conversationId,
      $or: [{ userId }, { sellerId: userId }],
    }).lean();

    return conversation;
  } catch (error) {
    console.error('Error in getConversationById:', error);
    throw error;
  }
};

/**
 * Get paginated messages for a conversation
 * @param {string} conversationId - Conversation ID
 * @param {number} page - Page number
 * @param {number} limit - Items per page
 * @returns {Promise<Object>} Paginated messages
 */
export const getConversationMessages = async (conversationId, page = PAGINATION.DEFAULT_PAGE, limit = PAGINATION.DEFAULT_LIMIT) => {
  try {
    const skip = (page - 1) * limit;
    const effectiveLimit = Math.min(limit, PAGINATION.MAX_LIMIT);

    const [messages, total] = await Promise.all([
      Message.find({ conversationId })
        .sort({ createdAt: -1 })
        .skip(skip)
        .limit(effectiveLimit)
        .lean(),
      Message.countDocuments({ conversationId }),
    ]);

    // Reverse to get chronological order (oldest to newest)
    messages.reverse();

    return {
      messages,
      pagination: {
        page,
        limit: effectiveLimit,
        total,
        pages: Math.ceil(total / effectiveLimit),
      },
    };
  } catch (error) {
    console.error('Error in getConversationMessages:', error);
    throw error;
  }
};

/**
 * Create a new message
 * @param {string} senderId - Sender user ID
 * @param {string} conversationId - Conversation ID
 * @param {string} content - Message content
 * @param {string} image - Image URL (optional)
 * @returns {Promise<Object>} Created message
 */
export const createMessage = async (senderId, conversationId, content, image = null) => {
  try {
    const message = await Message.create({
      senderId,
      conversationId,
      content,
      image,
      isRead: false,
    });

    // Update conversation's last message and timestamp
    await updateConversationLastMessage(conversationId, content);

    console.log(`✅ Created message in conversation ${conversationId}`);
    return message;
  } catch (error) {
    console.error('Error in createMessage:', error);
    throw error;
  }
};

/**
 * Update conversation's last message and timestamp
 * @param {string} conversationId - Conversation ID
 * @param {string} lastMessage - Last message content
 */
export const updateConversationLastMessage = async (conversationId, lastMessage) => {
  try {
    await Conversation.findByIdAndUpdate(conversationId, {
      lastMessage,
      lastUpdated: new Date(),
      $inc: { unreadCount: 1 },
    });
  } catch (error) {
    console.error('Error in updateConversationLastMessage:', error);
    throw error;
  }
};

/**
 * Mark messages as read in a conversation
 * @param {string} conversationId - Conversation ID
 * @param {string} userId - User ID who is reading the messages
 * @returns {Promise<number>} Number of messages marked as read
 */
export const markMessagesAsRead = async (conversationId, userId) => {
  try {
    // Mark all unread messages in the conversation where user is NOT the sender
    const result = await Message.updateMany(
      {
        conversationId,
        senderId: { $ne: userId },
        isRead: false,
      },
      {
        $set: { isRead: true },
      }
    );

    // Reset unread count for this conversation
    await Conversation.findByIdAndUpdate(conversationId, {
      unreadCount: 0,
    });

    console.log(`✅ Marked ${result.modifiedCount} messages as read in conversation ${conversationId}`);
    return result.modifiedCount;
  } catch (error) {
    console.error('Error in markMessagesAsRead:', error);
    throw error;
  }
};
