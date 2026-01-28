<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { socketService, SOCKET_EVENTS } from '@/utils/socket'
import {
  Bell,
  ChatDotRound,
  Goods,
  InfoFilled,
  Wallet,
  Discount,
  Loading,
} from '@element-plus/icons-vue'

const router = useRouter()
const notifications = ref<any[]>([])
const unreadCount = ref(0)
const isNotificationOpen = ref(false)
const isLoadingNotifications = ref(false)
const notificationDropdownRef = ref<HTMLElement | null>(null)

const formatRelativeTime = (dateStr: string) => {
  const date = new Date(dateStr)
  const now = new Date()
  const diffInSeconds = Math.floor((now.getTime() - date.getTime()) / 1000)

  if (diffInSeconds < 60) return 'Just now'

  const diffInMinutes = Math.floor(diffInSeconds / 60)
  if (diffInMinutes < 60) return `${diffInMinutes}m ago`

  const diffInHours = Math.floor(diffInMinutes / 60)
  if (diffInHours < 24) return `${diffInHours}h ago`

  const diffInDays = Math.floor(diffInHours / 24)
  if (diffInDays < 7) return `${diffInDays}d ago`

  return date.toLocaleDateString()
}

const fetchNotifications = async () => {
  const token = localStorage.getItem('access_token')
  if (!token) return

  isLoadingNotifications.value = true
  try {
    const response = await axios.get(`${import.meta.env.VITE_BE_API_URL}/notification`, {
      params: { limit: 20 },
      headers: { Authorization: `Bearer ${token}` },
    })

    if (response.data && response.data.data) {
      notifications.value = response.data.data
    }
  } catch (error) {
    console.error('Error fetching notifications:', error)
  } finally {
    isLoadingNotifications.value = false
  }
}

const fetchUnreadCount = async () => {
  const token = localStorage.getItem('access_token')
  if (!token) return

  try {
    const response = await axios.get(
      `${import.meta.env.VITE_BE_API_URL}/notification/unread-count`,
      {
        headers: { Authorization: `Bearer ${token}` },
      },
    )
    unreadCount.value = response.data.data?.unreadCount || 0
  } catch (error) {
    console.error('Error fetching unread count:', error)
  }
}

const markAsRead = async (notification: any) => {
  if (notification.isRead) return

  const token = localStorage.getItem('access_token')
  if (!token) return

  try {
    notification.isRead = true
    unreadCount.value = Math.max(0, unreadCount.value - 1)

    await axios.patch(
      `${import.meta.env.VITE_BE_API_URL}/notification/${notification._id}/read`,
      {},
      { headers: { Authorization: `Bearer ${token}` } },
    )
  } catch (error) {
    console.error('Error marking notification as read:', error)
    notification.isRead = false
    unreadCount.value++
  }
}

const markAllAsRead = async () => {
  if (unreadCount.value === 0) return

  const token = localStorage.getItem('access_token')
  if (!token) return

  try {
    notifications.value.forEach((n) => (n.isRead = true))
    unreadCount.value = 0

    await axios.patch(
      `${import.meta.env.VITE_BE_API_URL}/notification/read-all`,
      {},
      { headers: { Authorization: `Bearer ${token}` } },
    )
  } catch (error) {
    console.error('Error marking all as read:', error)
    fetchNotifications()
    fetchUnreadCount()
  }
}

const handleNotificationClick = (notification: any) => {
  markAsRead(notification)

  // Context-aware redirection
  const isSeller = router.currentRoute.value.path.startsWith('/seller')

  if (notification.type === 'order') {
    if (isSeller) {
      router.push('/seller/order')
    } else {
      if (notification.data && notification.data.orderId) {
        router.push(`/order-tracking/${notification.data.orderId}`)
      } else {
        router.push('/profile/orders')
      }
    }
    isNotificationOpen.value = false
  } else if (notification.type === 'chat') {
    if (isSeller) {
      if (notification.data && notification.data.senderId) {
        router.push('/seller/chat?focus=' + notification.data.senderId)
      } else {
        router.push('/seller/chat')
      }
    } else {
      if (notification.data && notification.data.senderId) {
        const event = new CustomEvent('open-chat', {
          detail: { sellerId: notification.data.senderId },
        })
        window.dispatchEvent(event)
      }
    }
    isNotificationOpen.value = false
  }
}

const toggleNotifications = () => {
  isNotificationOpen.value = !isNotificationOpen.value
  if (isNotificationOpen.value) {
    fetchNotifications()
    fetchUnreadCount()
  }
}

const handleClickOutside = (event: MouseEvent) => {
  if (
    notificationDropdownRef.value &&
    !notificationDropdownRef.value.contains(event.target as Node) &&
    !(event.target as Element).closest('.notification-btn')
  ) {
    isNotificationOpen.value = false
  }
}

const handleNewNotification = (notification: any) => {
  notifications.value.unshift(notification)
  // Optimistic update
  unreadCount.value++
}

const handleNotificationUpdated = (payload: any) => {
  if (payload.allRead) {
    notifications.value.forEach((n) => (n.isRead = true))
    unreadCount.value = 0
  } else if (payload.notificationId && payload.isRead) {
    const notif = notifications.value.find((n) => n._id === payload.notificationId)
    if (notif) {
      if (!notif.isRead) {
        notif.isRead = true
        unreadCount.value = Math.max(0, unreadCount.value - 1)
      }
    } else {
      fetchUnreadCount()
    }
  }
}

const initSocket = () => {
  const token = localStorage.getItem('access_token')
  if (token) {
    socketService.connect(token)
    socketService.emit(SOCKET_EVENTS.JOIN_NOTIFICATIONS, {})
    socketService.on(SOCKET_EVENTS.NEW_NOTIFICATION, handleNewNotification)
    socketService.on(SOCKET_EVENTS.NOTIFICATION_UPDATED, handleNotificationUpdated)
  }
}

const cleanupSocket = () => {
  socketService.off(SOCKET_EVENTS.NEW_NOTIFICATION, handleNewNotification)
  socketService.off(SOCKET_EVENTS.NOTIFICATION_UPDATED, handleNotificationUpdated)
}

onMounted(() => {
  fetchUnreadCount()
  initSocket()
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  cleanupSocket()
  document.removeEventListener('click', handleClickOutside)
})

defineExpose({
  fetchNotifications,
  fetchUnreadCount,
})
</script>

<template>
  <div class="notification-wrapper">
    <button
      class="notification-btn"
      title="Notifications"
      @click="toggleNotifications"
      :class="{ active: isNotificationOpen }"
    >
      <div class="icon-wrapper">
        <svg
          width="22"
          height="22"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"></path>
          <path d="M13.73 21a2 2 0 0 1-3.46 0"></path>
        </svg>
        <span class="badge" v-if="unreadCount > 0">{{
          unreadCount > 99 ? '99+' : unreadCount
        }}</span>
      </div>
    </button>

    <div v-show="isNotificationOpen" class="notification-dropdown" ref="notificationDropdownRef">
      <div class="notification-header">
        <h3>Notifications</h3>
        <button v-if="unreadCount > 0" class="mark-all-btn" @click="markAllAsRead">
          Mark all as read
        </button>
      </div>

      <div class="notification-list">
        <div v-if="isLoadingNotifications" class="loading-state">
          <el-icon class="is-loading" :size="24"><Loading /></el-icon>
        </div>
        <div v-if="notifications.length === 0 && !isLoadingNotifications" class="empty-state">
          <el-icon :size="48" color="#dedede"><Bell /></el-icon>
          <p>No notifications yet</p>
        </div>

        <div
          v-for="notification in notifications"
          :key="notification._id"
          class="notification-item"
          :class="{ unread: !notification.isRead }"
          @click="handleNotificationClick(notification)"
        >
          <div class="notification-icon">
            <img
              v-if="notification.data?.image"
              :src="notification.data.image"
              class="notif-img"
              onerror="this.style.display='none'"
            />
            <div v-else class="notif-icon-placeholder" :class="notification.type">
              <el-icon v-if="notification.type === 'chat'"><ChatDotRound /></el-icon>
              <el-icon v-else-if="notification.type === 'order'"><Goods /></el-icon>
              <el-icon v-else-if="notification.type === 'payment'"><Wallet /></el-icon>
              <el-icon v-else-if="notification.type === 'promotion'"><Discount /></el-icon>
              <el-icon v-else><InfoFilled /></el-icon>
            </div>
          </div>
          <div class="notification-content">
            <p class="notif-title">{{ notification.title }}</p>
            <p class="notif-message">{{ notification.message }}</p>
            <span class="notif-time">{{ formatRelativeTime(notification.createdAt) }}</span>
          </div>
          <div class="unread-dot" v-if="!notification.isRead"></div>
        </div>
      </div>
      <div class="notification-footer"></div>
    </div>
  </div>
</template>

<style scoped>
.notification-wrapper {
  position: relative;
}

.notification-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 8px;
  color: #71717a;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 12px;
  transition: all 0.2s ease;
}

.notification-btn:hover,
.notification-btn.active {
  color: #22c55e;
  background-color: #f0fdf4;
}

.notification-btn:hover .icon-wrapper {
  color: #22c55e;
}

.notification-dropdown {
  position: absolute;
  top: calc(100% + 10px);
  right: -80px;
  width: 360px;
  background-color: white;
  border-radius: 12px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.15);
  border: 1px solid #f4f4f5;
  z-index: 2000;
  display: flex;
  flex-direction: column;
  max-height: 500px;
  overflow: hidden;
  animation: slideDown 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}

.notification-dropdown::before {
  content: '';
  position: absolute;
  top: -6px;
  right: 94px;
  width: 12px;
  height: 12px;
  background-color: white;
  border-top: 1px solid #f4f4f5;
  border-left: 1px solid #f4f4f5;
  transform: rotate(45deg);
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.notification-header {
  padding: 16px;
  border-bottom: 1px solid #f4f4f5;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.notification-header h3 {
  font-size: 16px;
  font-weight: 600;
  margin: 0;
  color: #18181b;
}

.mark-all-btn {
  background: none;
  border: none;
  color: #22c55e;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 4px;
}

.mark-all-btn:hover {
  background-color: #f0fdf4;
}

.notification-list {
  overflow-y: auto;
  flex: 1;
  max-height: 400px;
}

.notification-item {
  display: flex;
  gap: 12px;
  padding: 12px 16px;
  border-bottom: 1px solid #f4f4f5;
  cursor: pointer;
  transition: background-color 0.2s;
  position: relative;
}

.notification-item:hover {
  background-color: #fafafa;
}

.notification-item.unread {
  background-color: #f0fdf4;
}

.notification-icon {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  background-color: #f4f4f5;
  overflow: hidden;
  flex-shrink: 0;
}

.notif-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.notif-icon-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 20px;
}

.notif-icon-placeholder.chat {
  background-color: #3b82f6;
}
.notif-icon-placeholder.order {
  background-color: #22c55e;
}
.notif-icon-placeholder.payment {
  background-color: #f59e0b;
}
.notif-icon-placeholder.promotion {
  background-color: #ec4899;
}
.notif-icon-placeholder.system {
  background-color: #64748b;
}

.notification-content {
  flex: 1;
  min-width: 0;
}

.notif-title {
  font-size: 14px;
  font-weight: 600;
  color: #18181b;
  margin: 0 0 4px 0;
}

.notif-message {
  font-size: 13px;
  color: #52525b;
  margin: 0 0 6px 0;
  line-height: 1.4;
  display: -webkit-box;
  line-clamp: 2;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.notif-time {
  font-size: 11px;
  color: #a1a1aa;
}

.unread-dot {
  width: 8px;
  height: 8px;
  background-color: #22c55e;
  border-radius: 50%;
  align-self: center;
  flex-shrink: 0;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  color: #a1a1aa;
  gap: 10px;
}

.loading-state {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 400px;
  color: #22c55e;
}

.notification-footer {
  padding: 10px;
  border-top: 1px solid #f4f4f5;
  text-align: center;
}

.icon-wrapper {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #18181b;
}

.badge {
  position: absolute;
  top: -8px;
  right: -8px;
  background-color: #22c55e;
  color: white;
  font-size: 11px;
  font-weight: 700;
  min-width: 18px;
  height: 18px;
  padding: 0 4px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2px solid white;
}
</style>
