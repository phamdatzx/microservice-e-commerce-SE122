# Notification Service

A microservice for handling notifications (email, SMS, push notifications) in the e-commerce platform.

## Features

- ✉️ Email notifications
- 📱 SMS notifications
- 🔔 Push notifications
- 📲 In-app notifications
- 📊 Notification history and tracking

## Tech Stack

- **Runtime**: Node.js
- **Framework**: Express.js
- **Architecture**: RESTful API

## Project Structure

```
notification-service/
├── controllers/          # Request handlers
├── services/            # Business logic
├── routes/              # API routes
├── middleware/          # Custom middleware
├── utils/               # Utility functions
├── config/              # Configuration files
├── index.js             # Application entry point
└── package.json         # Dependencies
```

## Getting Started

### Prerequisites

- Node.js (v18 or higher)
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

3. Configure environment variables in `.env`

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

### Notifications
- `POST /api/notifications` - Send a notification
- `GET /api/notifications` - Get notifications (with filters)
- `GET /api/notifications/:id` - Get specific notification
- `PATCH /api/notifications/:id/read` - Mark as read

## Environment Variables

See `.env.example` for all available configuration options.

## Development

### Adding New Notification Types

1. Add type to `utils/constants.js`
2. Create provider in `services/providers.js`
3. Update service logic in `services/notificationService.js`

## TODO

- [ ] Integrate with actual email service (SendGrid, Mailgun)
- [ ] Integrate with SMS service (Twilio)
- [ ] Integrate with push notification service (Firebase)
- [ ] Add database persistence (MongoDB/PostgreSQL)
- [ ] Add message queue (RabbitMQ/Redis)
- [ ] Add authentication middleware
- [ ] Add rate limiting
- [ ] Add notification templates
- [ ] Add batch notification support
- [ ] Add notification scheduling
- [ ] Add WebSocket support for real-time notifications

## License

ISC
