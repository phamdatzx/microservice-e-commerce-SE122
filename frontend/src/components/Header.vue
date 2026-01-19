<script setup lang="ts">
import ChatWindow from './ChatWindow.vue'
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import IconUser from './icons/IconUser.vue'
import IconOrder from './icons/IconOrder.vue'
import IconVoucher from './icons/IconVoucher.vue'
import IconLock from './icons/IconLock.vue'
import IconLogout from './icons/IconLogout.vue'
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
const isLoggedIn = ref(false)
const cartCount = ref(0)
const userInfo = ref<any>(null)
const API_URL = `${import.meta.env.VITE_BE_API_URL}/order/cart/count`

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

  if (notification.type === 'order') {
    if (notification.data && notification.data.orderId) {
      router.push(`/order-tracking/${notification.data.orderId}`)
    } else {
      router.push('/profile/orders')
    }
    isNotificationOpen.value = false
  } else if (notification.type === 'chat') {
    if (notification.data && notification.data.senderId) {
      const event = new CustomEvent('open-chat', {
        detail: { sellerId: notification.data.senderId },
      })
      window.dispatchEvent(event)
    }
    isNotificationOpen.value = false
  } else {
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

const fetchCartCount = async () => {
  const token = localStorage.getItem('access_token')
  if (!token) return

  try {
    const response = await axios.get(API_URL, {
      headers: { Authorization: `Bearer ${token}` },
    })
    cartCount.value = response.data.count
  } catch (error) {
    console.error('Error fetching cart count:', error)
  }
}

const fetchUserInfo = async () => {
  const token = localStorage.getItem('access_token')
  if (!token) return

  try {
    const response = await axios.get(`${import.meta.env.VITE_BE_API_URL}/user/my-info`, {
      headers: { Authorization: `Bearer ${token}` },
    })
    userInfo.value = response.data.data
  } catch (error) {
    console.error('Error fetching user info:', error)
  }
}

const handleNewNotification = (notification: any) => {
  notifications.value.unshift(notification)

  fetchUnreadCount()
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

const handleNotificationsJoined = (data: any) => {}

const initSocket = () => {
  const token = localStorage.getItem('access_token')
  if (token) {
    socketService.connect(token)

    socketService.emit(SOCKET_EVENTS.JOIN_NOTIFICATIONS, {})

    socketService.on(SOCKET_EVENTS.NEW_NOTIFICATION, handleNewNotification)
    socketService.on(SOCKET_EVENTS.NOTIFICATION_UPDATED, handleNotificationUpdated)
    socketService.on(SOCKET_EVENTS.NOTIFICATIONS_JOINED, handleNotificationsJoined)
  }
}

const cleanupSocket = () => {
  socketService.off(SOCKET_EVENTS.NEW_NOTIFICATION, handleNewNotification)
  socketService.off(SOCKET_EVENTS.NOTIFICATION_UPDATED, handleNotificationUpdated)
  socketService.off(SOCKET_EVENTS.NOTIFICATIONS_JOINED, handleNotificationsJoined)
}

watch(isLoggedIn, (newVal) => {
  if (newVal) {
    initSocket()
  } else {
    cleanupSocket()
    socketService.disconnect()
  }
})

onMounted(() => {
  isLoggedIn.value = !!localStorage.getItem('access_token')
  if (isLoggedIn.value) {
    fetchCartCount()
    fetchUserInfo()
    fetchUnreadCount()

    initSocket()
  }

  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
  cleanupSocket()
})

const searchQuery = ref('')
const handleSearch = () => {
  if (searchQuery.value.trim()) {
    router.push({ path: '/search', query: { q: searchQuery.value } })
  }
}

const handleLogout = () => {
  localStorage.removeItem('access_token')
  localStorage.removeItem('user_id')
  isLoggedIn.value = false
  router.push('/')
}

defineExpose({
  fetchCartCount,
  fetchUserInfo,
  fetchNotifications,
  fetchUnreadCount,
})
</script>

<template>
  <!-- Header -->
  <header class="header">
    <div class="header-main">
      <div class="container">
        <div class="header-content">
          <RouterLink to="/" class="logo">
            <div class="logo-icon">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
                <path
                  d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                  stroke="white"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
              </svg>
            </div>
            <div class="logo-text">
              <span class="logo-name">SWOO</span>
              <span class="logo-tagline">TECH MART</span>
            </div>
          </RouterLink>

          <div class="search-bar">
            <input
              type="text"
              placeholder="Search anything..."
              class="search-input"
              v-model="searchQuery"
              @keyup.enter="handleSearch"
            />
            <button class="search-btn" @click="handleSearch">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="#fff"
                x="0px"
                y="0px"
                width="20"
                height="20"
                viewBox="0 0 50 50"
              >
                <path
                  d="M 21 3 C 11.621094 3 4 10.621094 4 20 C 4 29.378906 11.621094 37 21 37 C 24.710938 37 28.140625 35.804688 30.9375 33.78125 L 44.09375 46.90625 L 46.90625 44.09375 L 33.90625 31.0625 C 36.460938 28.085938 38 24.222656 38 20 C 38 10.621094 30.378906 3 21 3 Z M 21 5 C 29.296875 5 36 11.703125 36 20 C 36 28.296875 29.296875 35 21 35 C 12.703125 35 6 28.296875 6 20 C 6 11.703125 12.703125 5 21 5 Z"
                ></path>
              </svg>
            </button>
          </div>

          <div class="header-actions">
            <!-- Notification Section -->
            <div class="notification-wrapper" v-if="isLoggedIn">
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

              <!-- Notifications Dropdown -->
              <div
                v-show="isNotificationOpen"
                class="notification-dropdown"
                ref="notificationDropdownRef"
              >
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
                  <div
                    v-if="notifications.length === 0 && !isLoadingNotifications"
                    class="empty-state"
                  >
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
                        <el-icon v-else-if="notification.type === 'promotion'"
                          ><Discount
                        /></el-icon>
                        <el-icon v-else><InfoFilled /></el-icon>
                      </div>
                    </div>
                    <div class="notification-content">
                      <p class="notif-title">{{ notification.title }}</p>
                      <p class="notif-message">{{ notification.message }}</p>
                      <span class="notif-time">{{
                        formatRelativeTime(notification.createdAt)
                      }}</span>
                    </div>
                    <div class="unread-dot" v-if="!notification.isRead"></div>
                  </div>
                </div>

                <div class="notification-footer">
                  <!-- <RouterLink to="/notifications" class="view-all-link">View all</RouterLink> -->
                </div>
              </div>
            </div>

            <!-- Cart Section on the Left of User -->
            <div class="action-item cart-section" @click="$router.push('/cart')">
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
                  <circle cx="9" cy="21" r="1"></circle>
                  <circle cx="20" cy="21" r="1"></circle>
                  <path d="M1 1h4l2.68 13.39a2 2 0 0 0 2 1.61h9.72a2 2 0 0 0 2-1.61L23 6H6"></path>
                </svg>
                <span class="badge" v-if="cartCount > 0">{{
                  cartCount > 99 ? '99+' : cartCount
                }}</span>
              </div>
            </div>

            <!-- User Section -->
            <div v-if="!isLoggedIn" class="action-item user-section">
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
                  <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
                  <circle cx="12" cy="7" r="4"></circle>
                </svg>
              </div>
              <div class="text-content">
                <span class="label">WELCOME</span>
                <div class="auth-links">
                  <RouterLink to="/login" class="main-text">LOG IN</RouterLink>
                  <span class="separator">/</span>
                  <RouterLink to="/register" class="main-text">REGISTER</RouterLink>
                </div>
              </div>
            </div>

            <div v-else class="user-container">
              <div class="action-item user-profile">
                <div class="avatar-wrapper">
                  <img
                    :src="
                      userInfo?.image ||
                      `https://ui-avatars.com/api/?name=${userInfo?.name || 'User'}&background=22c55e&color=fff`
                    "
                    alt="Avatar"
                    class="avatar"
                  />
                </div>
                <div class="text-content">
                  <span class="label">MY ACCOUNT</span>
                  <span class="main-text">Hello, {{ userInfo?.name || 'User' }}</span>
                </div>
              </div>

              <!-- Dropdown Menu -->
              <div class="user-dropdown">
                <RouterLink to="/profile" class="dropdown-item">
                  <IconUser />
                  My Profile
                </RouterLink>
                <RouterLink to="/profile/orders" class="dropdown-item">
                  <IconOrder />
                  My Orders
                </RouterLink>
                <RouterLink to="/profile/vouchers" class="dropdown-item">
                  <IconVoucher />
                  My Vouchers
                </RouterLink>
                <RouterLink to="/profile/change-password" class="dropdown-item">
                  <IconLock />
                  Change Password
                </RouterLink>
                <div class="dropdown-divider"></div>
                <button @click="handleLogout" class="dropdown-item logout-btn">
                  <IconLogout />
                  Logout
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </header>

  <ChatWindow v-if="isLoggedIn" />
</template>

<style scoped>
/* Header Styles */
.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.header-main {
  background-color: white;
  padding: 20px 0;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
  text-decoration: none;
}

.logo-icon {
  width: 44px;
  height: 44px;
  background: linear-gradient(135deg, #22c55e, #16a34a);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 12px rgba(34, 197, 94, 0.3);
}

.logo-text {
  display: flex;
  flex-direction: column;
}

.logo-name {
  font-size: 26px;
  font-weight: 800;
  color: #1a1a1a;
  line-height: 1;
  letter-spacing: -0.5px;
}

.logo-tagline {
  font-size: 11px;
  color: #71717a;
  letter-spacing: 2px;
  font-weight: 600;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 24px;
}

/* Notification Styles */
.notification-wrapper {
  position: relative;
}

.notification-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 8px; /* Reduced padding slightly */
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

.notification-dropdown {
  position: absolute;
  top: calc(100% + 10px);
  right: -80px; /* Shift right slightly to align better with layout */
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
  right: 94px; /* Align arrow with bell icon */
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

.view-all-link {
  font-size: 13px;
  color: #22c55e;
  font-weight: 500;
  text-decoration: none;
}

.action-item {
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  padding: 8px 12px;
  border-radius: 12px;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
}

.action-item:hover {
  background-color: #f4f4f5;
}

.user-container {
  position: relative;
  display: flex;
  align-items: center;
}

.user-profile {
  display: flex;
  align-items: center;
  gap: 12px;
}

.avatar-wrapper {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  overflow: hidden;
  border: 2px solid #e4e4e7;
  transition: border-color 0.2s ease;
}

.user-container:hover .avatar-wrapper {
  border-color: #22c55e;
}

.avatar {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

/* User Dropdown */
.user-dropdown {
  position: absolute;
  top: 100%;
  right: 0;
  width: 200px;
  background-color: white;
  border-radius: 12px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
  border: 1px solid #f4f4f5;
  padding: 8px;
  margin-top: 8px;
  opacity: 0;
  visibility: hidden;
  transform: translateY(10px);
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  z-index: 1000;
}

.user-container:hover .user-dropdown {
  opacity: 1;
  visibility: visible;
  transform: translateY(0);
}

.dropdown-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  color: #3f3f46;
  text-decoration: none;
  font-size: 14px;
  font-weight: 500;
  border-radius: 8px;
  transition: all 0.2s ease;
  width: 100%;
  border: none;
  background: none;
  cursor: pointer;
  text-align: left;
}

.dropdown-item:hover {
  background-color: #f0fdf4;
  color: #22c55e;
}

.dropdown-item svg {
  color: #a1a1aa;
  transition: color 0.2s ease;
}

.dropdown-item:hover svg {
  color: #22c55e;
}

.dropdown-divider {
  height: 1px;
  background-color: #f4f4f5;
  margin: 6px 0;
}

.logout-btn:hover {
  background-color: #fef2f2;
  color: #ef4444;
}

.logout-btn:hover svg {
  color: #ef4444;
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

.text-content {
  display: flex;
  flex-direction: column;
}

.label {
  font-size: 10px;
  font-weight: 700;
  color: #a1a1aa;
  letter-spacing: 0.5px;
}

.main-text {
  font-size: 14px;
  font-weight: 700;
  color: #18181b;
  text-decoration: none;
  transition: color 0.2s ease;
}

.auth-links .main-text {
  padding: 4px 6px;
  border-radius: 8px;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.auth-links .main-text:hover {
  color: #22c55e;
}

.auth-links {
  display: flex;
  align-items: center;
  gap: 2px;
  margin-left: -6px; /* Offset padding to keep alignment */
}

.auth-links .separator {
  font-size: 12px;
  color: #e4e4e7;
  font-weight: 400;
  user-select: none;
}

.cart-section .main-text {
  color: #22c55e;
}

/* Search Bar */
.search-bar {
  display: flex;
  flex: 1;
  max-width: 540px;
  margin: 0 40px;
  background: #f4f4f5;
  border-radius: 7px;
  overflow: hidden;
  transition: all 0.2s ease;
  border: 2px solid #e4e4e7;
}

.search-bar:focus-within {
  background: white;
  border-color: #22c55e;
  box-shadow: 0 4px 12px rgba(34, 197, 94, 0.1);
}

.search-input {
  flex: 1;
  padding: 14px 16px;
  border: none;
  font-size: 14px;
  outline: none;
  color: #18181b;
  background: transparent;
}

.search-input::placeholder {
  color: #a1a1aa;
}

.search-btn {
  background: var(--main-color);
  border: none;
  padding: 0 20px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background-color 0.2s ease;
}

.search-btn:hover {
  background: #22c55e;
}
</style>
