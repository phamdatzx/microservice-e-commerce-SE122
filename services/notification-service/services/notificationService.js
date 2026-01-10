// In-memory storage (replace with database in production)
let notifications = [];
let notificationIdCounter = 1;

export const sendNotification = async (notificationData) => {
  const notification = {
    id: String(notificationIdCounter++),
    ...notificationData,
    isRead: false,
    createdAt: new Date().toISOString(),
  };

  notifications.push(notification);

  // TODO: Implement actual notification sending logic
  // - Send email via email service
  // - Send SMS via SMS service
  // - Send push notification via push service
  // - Publish to message queue for async processing

  console.log('📧 Notification sent:', notification);

  return notification;
};

export const getNotifications = async (filters) => {
  let result = [...notifications];

  // Apply filters
  if (filters.userId) {
    result = result.filter((n) => n.userId === filters.userId);
  }

  if (filters.isRead !== undefined) {
    result = result.filter((n) => n.isRead === filters.isRead);
  }

  if (filters.type) {
    result = result.filter((n) => n.type === filters.type);
  }

  // Apply pagination
  const start = filters.offset || 0;
  const end = start + (filters.limit || 20);
  result = result.slice(start, end);

  return result;
};

export const getNotificationById = async (id) => {
  return notifications.find((n) => n.id === id);
};

export const markAsRead = async (id) => {
  const notification = notifications.find((n) => n.id === id);

  if (notification) {
    notification.isRead = true;
    notification.readAt = new Date().toISOString();
  }

  return notification;
};
