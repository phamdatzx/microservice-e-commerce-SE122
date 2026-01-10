import mongoose from 'mongoose';

const conversationSchema = new mongoose.Schema(
  {
    sellerId: {
      type: String,
      required: true,
      index: true,
    },
    userId: {
      type: String,
      required: true,
      index: true,
    },
    lastMessage: {
      type: String,
      default: '',
    },
    lastUpdated: {
      type: Date,
      default: Date.now,
    },
  },
  {
    timestamps: true,
  }
);

// Compound index for finding conversations between specific users
conversationSchema.index({ sellerId: 1, userId: 1 }, { unique: true });

// Index for sorting by last updated
conversationSchema.index({ lastUpdated: -1 });

const Conversation = mongoose.model('Conversation', conversationSchema);

export default Conversation;
