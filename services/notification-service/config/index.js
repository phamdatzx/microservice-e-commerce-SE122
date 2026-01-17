import dotenv from 'dotenv';
dotenv.config();

export default {
  port: process.env.PORT || 3000,
  nodeEnv: process.env.NODE_ENV || 'development',
  
  // Database configuration
  database: {
    mongodb: process.env.MONGO_URI,
  },

  // Redis configuration
  redis: {
    host: process.env.REDIS_HOST || 'localhost',
    port: process.env.REDIS_PORT || 6379,
  },

  // Email service configuration
  email: {
    apiKey: process.env.EMAIL_SERVICE_API_KEY,
    from: process.env.EMAIL_FROM || 'noreply@example.com',
  },

  // SMS service configuration
  sms: {
    accountSid: process.env.TWILIO_ACCOUNT_SID,
    authToken: process.env.TWILIO_AUTH_TOKEN,
    phoneNumber: process.env.TWILIO_PHONE_NUMBER,
  },

  // Push notification configuration
  push: {
    firebaseServerKey: process.env.FIREBASE_SERVER_KEY,
  },

  // API configuration
  api: {
    secretKey: process.env.API_SECRET_KEY,
  },
};
