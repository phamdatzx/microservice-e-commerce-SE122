import * as chatService from '../services/chatService.js';
import { SOCKET_EVENTS } from '../utils/constants.js';
import { getConversationRoom, emitToConversation, emitToUser } from '../utils/socketHelpers.js';

/**
 * Register chat-related socket event handlers
 * @param {import('socket.io').Server} io - Socket.io server instance
 * @param {import('socket.io').Socket} socket - Client socket
 */
export const registerChatHandlers = (io, socket) => {
  const userId = socket.userId;

  /**
   * Join a conversation room
   * Payload: { conversationId } or { sellerId } or { userId: customerId }
   */
  socket.on(SOCKET_EVENTS.JOIN_CONVERSATION, async (payload) => {
    try {
      let conversationId = payload.conversationId;

      // If conversationId not provided, find or create conversation
      if (!conversationId) {
        const { sellerId, userId: customerId } = payload;
        
        let finalUserId, finalSellerId;
        if (sellerId) {
          finalUserId = userId;
          finalSellerId = sellerId;
        } else if (customerId) {
          finalUserId = customerId;
          finalSellerId = userId;
        } else {
          socket.emit(SOCKET_EVENTS.ERROR, {
            message: 'Either conversationId, sellerId, or userId is required',
          });
          return;
        }

        const conversation = await chatService.findOrCreateConversation(finalUserId, finalSellerId);
        conversationId = conversation._id.toString();
      }

      // Verify user has access to this conversation
      const conversation = await chatService.getConversationById(conversationId, userId);
      if (!conversation) {
        socket.emit(SOCKET_EVENTS.ERROR, {
          message: 'Access denied to this conversation',
        });
        return;
      }

      // Join the conversation room
      const room = getConversationRoom(conversationId);
      socket.join(room);

      socket.emit(SOCKET_EVENTS.CONVERSATION_JOINED, {
        conversationId,
        conversation,
      });

      console.log(`👤 User ${userId} joined conversation ${conversationId}`);
    } catch (error) {
      console.error('Error in JOIN_CONVERSATION:', error);
      socket.emit(SOCKET_EVENTS.ERROR, {
        message: 'Failed to join conversation',
      });
    }
  });

  /**
   * Leave a conversation room
   * Payload: { conversationId }
   */
  socket.on(SOCKET_EVENTS.LEAVE_CONVERSATION, async (payload) => {
    try {
      const { conversationId } = payload;
      const room = getConversationRoom(conversationId);
      socket.leave(room);

      socket.emit(SOCKET_EVENTS.CONVERSATION_LEFT, {
        conversationId,
      });

      console.log(`👤 User ${userId} left conversation ${conversationId}`);
    } catch (error) {
      console.error('Error in LEAVE_CONVERSATION:', error);
    }
  });

  /**
   * Send a message
   * Payload: { conversationId, content, image }
   */
  socket.on(SOCKET_EVENTS.SEND_MESSAGE, async (payload) => {
    try {
      const { conversationId, content, image } = payload;

      if (!conversationId || !content) {
        socket.emit(SOCKET_EVENTS.ERROR, {
          message: 'conversationId and content are required',
        });
        return;
      }

      // Verify user has access to this conversation
      const conversation = await chatService.getConversationById(conversationId, userId);
      if (!conversation) {
        socket.emit(SOCKET_EVENTS.ERROR, {
          message: 'Access denied to this conversation',
        });
        return;
      }

      // Create the message
      const message = await chatService.createMessage(userId, conversationId, content, image);

      // Send confirmation to sender
      socket.emit(SOCKET_EVENTS.MESSAGE_SENT, message);

      // Determine the receiver ID
      const receiverId = conversation.userId === userId ? conversation.sellerId : conversation.userId;

      // Check if receiver is in the conversation room
      const room = getConversationRoom(conversationId);
      const socketsInRoom = await io.in(room).fetchSockets();
      const receiverInRoom = socketsInRoom.some(s => s.userId === receiverId);

      // Broadcast to all users in the conversation room
      emitToConversation(io, conversationId, SOCKET_EVENTS.NEW_MESSAGE, message);

      // If receiver is not in the room, create a notification
      if (!receiverInRoom) {
        console.log(`📲 Receiver ${receiverId} not in room, creating notification`);
        
        const { createNotification } = await import('../services/notificationService.js');
        const { NOTIFICATION_TYPES } = await import('../utils/constants.js');
        const { emitToUser } = await import('../utils/socketHelpers.js');
        
        try {
          const notification = await createNotification(
            receiverId,
            NOTIFICATION_TYPES.CHAT,
            'New Message',
            content,
            {
              conversationId,
              senderId: userId,
              messageId: message._id,
              hasImage: !!image,
            }
          );

          // Try to send notification to user if they're connected (but not in room)
          emitToUser(io, receiverId, SOCKET_EVENTS.NEW_NOTIFICATION, notification);
        } catch (notifError) {
          console.error('Failed to create notification:', notifError);
          // Don't fail the message send if notification creation fails
        }
      }

      console.log(`💬 Message sent in conversation ${conversationId}`);
    } catch (error) {
      console.error('Error in SEND_MESSAGE:', error);
      socket.emit(SOCKET_EVENTS.ERROR, {
        message: 'Failed to send message',
      });
    }
  });

  /**
   * Typing indicator
   * Payload: { conversationId, isTyping }
   */
  socket.on(SOCKET_EVENTS.TYPING, async (payload) => {
    try {
      const { conversationId, isTyping } = payload;

      // Verify user has access to this conversation
      const conversation = await chatService.getConversationById(conversationId, userId);
      if (!conversation) {
        return;
      }

      // Broadcast typing status to conversation room (except sender)
      const room = getConversationRoom(conversationId);
      socket.to(room).emit(SOCKET_EVENTS.USER_TYPING, {
        userId,
        conversationId,
        isTyping,
      });
    } catch (error) {
      console.error('Error in TYPING:', error);
    }
  });

  /**
   * Mark messages as read
   * Payload: { conversationId }
   */
  socket.on(SOCKET_EVENTS.MESSAGE_READ, async (payload) => {
    try {
      const { conversationId } = payload;

      // Verify user has access to this conversation
      const conversation = await chatService.getConversationById(conversationId, userId);
      if (!conversation) {
        socket.emit(SOCKET_EVENTS.ERROR, {
          message: 'Access denied to this conversation',
        });
        return;
      }

      // Mark messages as read
      const count = await chatService.markMessagesAsRead(conversationId, userId);

      // Notify all users in the conversation
      emitToConversation(io, conversationId, SOCKET_EVENTS.MESSAGES_UPDATED, {
        conversationId,
        markedAsRead: count,
      });

      console.log(`✅ Marked ${count} messages as read in conversation ${conversationId}`);
    } catch (error) {
      console.error('Error in MESSAGE_READ:', error);
      socket.emit(SOCKET_EVENTS.ERROR, {
        message: 'Failed to mark messages as read',
      });
    }
  });
};
