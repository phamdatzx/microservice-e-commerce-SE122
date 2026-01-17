# Notification Service

A microservice for handling real-time chat and notifications in the e-commerce platform.

## Features

### 💬 Chat
- Real-time messaging between customers and sellers
- Conversation management
- Message persistence in MongoDB
- Typing indicators
- Read receipts
- Message history with pagination

### 🔔 Notifications
- Real-time system notifications
- Triggered by other microservices
- Notification persistence in MongoDB
- Read/unread status tracking
- Notification types: order, payment, product, system, promotion
- Unread count tracking

## Tech Stack

- **Runtime**: Node.js
- **Framework**: Express.js
- **Real-time**: Socket.io
- **Database**: MongoDB with Mongoose
- **Architecture**: RESTful API + WebSocket

## Project Structure

```
notification-service/
├── config/              # Configuration files
│   ├── database.js      # MongoDB connection
│   └── socket.js        # Socket.io configuration
├── controllers/         # Request handlers
│   ├── chatController.js
│   └── notificationController.js
├── middleware/          # Custom middleware
│   ├── asyncHandler.js
│   ├── errorHandler.js
│   └── socketAuth.js    # WebSocket authentication
├── models/              # Mongoose models
│   ├── Conversation.js
│   ├── Message.js
│   ├── Notification.js
│   └── index.js
├── routes/              # API routes
│   ├── chatRoutes.js
│   ├── notificationRoutes.js
│   └── index.js
├── services/            # Business logic
│   ├── chatService.js
│   └── notificationService.js
├── sockets/             # WebSocket event handlers
│   ├── chatSocket.js
│   ├── notificationSocket.js
│   └── index.js
├── utils/               # Utility functions
│   ├── constants.js
│   ├── socketHelpers.js
│   └── validators.js
├── index.js             # Application entry point
├── test-client.html     # WebSocket test client
└── package.json         # Dependencies
```

## Getting Started

### Prerequisites

- Node.js (v18 or higher)
- MongoDB instance
- npm or yarn

### Installation

1. Install dependencies:
```bash
npm install
```

2. Create `.env` file:
```bash
cp .env.example .env
```

3. Configure environment variables in `.env`:
```env
PORT=8085
NODE_ENV=development
MONGO_URI=<your-mongodb-connection-string>
CORS_ORIGIN=*
```

### Running the Service

Development mode (with auto-reload):
```bash
npm run dev
```

Production mode:
```bash
npm start
```

## API Endpoints

### Health Check
- `GET /health` - Check service health

### Chat Endpoints
- `GET /api/chat/conversations` - Get all conversations for a user
- `POST /api/chat/conversations` - Create or get conversation
- `GET /api/chat/conversations/:conversationId` - Get conversation by ID
- `GET /api/chat/conversations/:conversationId/messages` - Get messages (paginated)
- `POST /api/chat/messages` - Send message (HTTP fallback)
- `PATCH /api/chat/conversations/:conversationId/read` - Mark messages as read

### Notification Endpoints
- `POST /api/notifications` - Create notification (for other services)
- `GET /api/notifications` - Get user's notifications (paginated)
- `GET /api/notifications/unread-count` - Get unread count
- `PATCH /api/notifications/:notificationId/read` - Mark as read
- `PATCH /api/notifications/read-all` - Mark all as read

## WebSocket Events

### Connection
Connect to the WebSocket server by passing `userId` in the query parameters:
```javascript
const socket = io('http://localhost:8085', {
  query: { userId: 'user123' }
});
```

### Chat Events

#### Client → Server
- `join-conversation` - Join a conversation room
  ```json
  { "conversationId": "conv_id" }
  // or
  { "sellerId": "seller_id" }
  ```
- `leave-conversation` - Leave a conversation room
  ```json
  { "conversationId": "conv_id" }
  ```
- `send-message` - Send a message
  ```json
  { "conversationId": "conv_id", "content": "Hello!", "image": "url" }
  ```
- `typing` - Typing indicator
  ```json
  { "conversationId": "conv_id", "isTyping": true }
  ```
- `message-read` - Mark messages as read
  ```json
  { "conversationId": "conv_id" }
  ```

#### Server → Client
- `conversation-joined` - Conversation joined successfully
- `conversation-left` - Conversation left successfully
- `new-message` - New message received
- `message-sent` - Message sent confirmation
- `user-typing` - Another user is typing
- `messages-updated` - Messages read status updated

### Notification Events

#### Client → Server
- `join-notifications` - Join user's notification room
- `notification-read` - Mark notification as read
  ```json
  { "notificationId": "notif_id" }
  ```

#### Server → Client
- `notifications-joined` - Joined notification room
- `new-notification` - New notification received
- `notification-updated` - Notification status updated

## Testing

### Using the Test Client

1. Start the notification service:
```bash
npm run dev
```

2. Open `test-client.html` in two different browser windows

3. In window 1:
   - Set User ID to `user123`
   - Click "Connect"
   - Click "Join Notifications"
   - Set Seller ID to `seller456`
   - Click "Join Conversation"

4. In window 2:
   - Set User ID to `seller456`
   - Click "Connect"
   - Set Seller ID (actually customer ID) to `user123`
   - Click "Join Conversation"

5. Send messages back and forth to test real-time chat

6. In window 1, create a test notification to see real-time delivery

### Using cURL

**Create a conversation:**
```bash
curl -X POST http://localhost:8085/api/chat/conversations \
  -H "X-User-Id: user123" \
  -H "Content-Type: application/json" \
  -d '{"sellerId": "seller456"}'
```

**Send a message:**
```bash
curl -X POST http://localhost:8085/api/chat/messages \
  -H "X-User-Id: user123" \
  -H "Content-Type: application/json" \
  -d '{"conversationId": "CONV_ID", "content": "Hello!"}'
```

**Create a notification:**
```bash
curl -X POST http://localhost:8085/api/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "user123",
    "type": "order",
    "title": "Order Confirmed",
    "message": "Your order #12345 has been confirmed",
    "data": {"orderId": "12345"}
  }'
```

**Get notifications:**
```bash
curl -X GET "http://localhost:8085/api/notifications?page=1&limit=20" \
  -H "X-User-Id: user123"
```

## Authentication

### HTTP Requests
Use the `X-User-Id` header for authentication:
```
X-User-Id: user123
```

### WebSocket Connection
Pass `userId` in the connection query parameters:
```javascript
const socket = io('http://localhost:8085', {
  query: { userId: 'user123' }
});
```

## Integration with Other Services

Other services can trigger notifications by calling:
```bash
POST http://notification-service:8085/api/notifications
Content-Type: application/json

{
  "userId": "user123",
  "type": "order|payment|product|system|promotion",
  "title": "Notification Title",
  "message": "Notification message",
  "data": {
    "orderId": "12345",
    // ... additional metadata
  }
}
```

The notification will be:
1. Saved to MongoDB
2. Delivered in real-time via WebSocket to the connected user

## Database Models

### Conversation
```javascript
{
  userId: String,        // Customer ID
  sellerId: String,      // Seller ID
  lastMessage: String,   // Last message content
  lastUpdated: Date,     // Last update timestamp
  unreadCount: Number,   // Unread message count
  createdAt: Date,
  updatedAt: Date
}
```

### Message
```javascript
{
  senderId: String,         // Sender user ID
  conversationId: ObjectId, // Reference to conversation
  content: String,          // Message content
  image: String,            // Image URL (optional)
  isRead: Boolean,          // Read status
  createdAt: Date,
  updatedAt: Date
}
```

### Notification
```javascript
{
  userId: String,    // Recipient user ID
  type: String,      // Notification type
  title: String,     // Notification title
  message: String,   // Notification message
  data: Object,      // Additional metadata
  isRead: Boolean,   // Read status
  createdAt: Date,
  updatedAt: Date
}
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `3000` |
| `NODE_ENV` | Environment | `development` |
| `MONGO_URI` | MongoDB connection string | Required |
| `CORS_ORIGIN` | CORS allowed origins | `*` |

## License

ISC
