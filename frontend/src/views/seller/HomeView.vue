<script setup lang="ts">
import {
  ArrowRight,
  Goods,
  Menu,
  ChatDotRound,
  User,
  Ticket,
  Van,
  TrendCharts,
} from '@element-plus/icons-vue'
import { onMounted, ref, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import UserAvatarDropdown from '@/components/UserAvatarDropdown.vue'
import NotificationDropdown from '@/components/NotificationDropdown.vue'
import { socketService, SOCKET_EVENTS } from '@/utils/socket'

const router = useRouter()
const unreadChatCount = ref(0)

const fetchUnreadChatCount = async () => {
  const token = localStorage.getItem('access_token')
  if (!token) return

  try {
    const response = await axios.get(
      `${import.meta.env.VITE_BE_API_URL}/notification/chat/conversations`,
      {
        params: { limit: 1000 },
        headers: { Authorization: `Bearer ${token}` },
      },
    )

    if (response.data && response.data.data) {
      unreadChatCount.value = response.data.data.reduce(
        (total: number, conv: any) => total + (conv.unreadCount || 0),
        0,
      )
    }
  } catch (error) {
    console.error('Error fetching unread chat count:', error)
  }
}

const setupSocketListeners = () => {
  socketService.on(SOCKET_EVENTS.NEW_MESSAGE, fetchUnreadChatCount)
  socketService.on(SOCKET_EVENTS.MESSAGE_READ, fetchUnreadChatCount)
  socketService.on(SOCKET_EVENTS.MESSAGES_UPDATED, fetchUnreadChatCount)
  socketService.on(SOCKET_EVENTS.NEW_NOTIFICATION, (notification: any) => {
    if (notification.type === 'chat') {
      fetchUnreadChatCount()
    }
  })
}

const cleanupSocketListeners = () => {
  socketService.off(SOCKET_EVENTS.NEW_MESSAGE, fetchUnreadChatCount)
  socketService.off(SOCKET_EVENTS.MESSAGE_READ, fetchUnreadChatCount)
  socketService.off(SOCKET_EVENTS.MESSAGES_UPDATED, fetchUnreadChatCount)
}

const handleOpen = (key: string, keyPath: string[]) => {}
const handleClose = (key: string, keyPath: string[]) => {}

const initSocket = () => {
  const token = localStorage.getItem('access_token')
  if (token) {
    socketService.connect(token)
    socketService.emit(SOCKET_EVENTS.JOIN_NOTIFICATIONS, {})
  }
}

const handleBadgeUpdate = () => {
  fetchUnreadChatCount()
}

onMounted(() => {
  initSocket()
  setupSocketListeners()
  fetchUnreadChatCount()
  window.addEventListener('trigger-chat-badge-update', handleBadgeUpdate)
})

onUnmounted(() => {
  cleanupSocketListeners()
  window.removeEventListener('trigger-chat-badge-update', handleBadgeUpdate)
})
</script>

<template>
  <div class="common-layout">
    <el-container>
      <el-header class="header">
        <div class="container">
          <div class="header-content">
            <div class="logo">
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
            </div>

            <div class="breadcrumb">
              <div class="container">
                <a href="#" class="breadcrumb-link">Home</a>
                <span class="breadcrumb-separator">
                  <el-icon><ArrowRight /></el-icon>
                </span>
                <a href="#" class="breadcrumb-link">pages</a>
                <span class="breadcrumb-separator">
                  <el-icon><ArrowRight /></el-icon>
                </span>
                <span class="breadcrumb-current">login</span>
              </div>
            </div>

            <div class="header-right">
              <NotificationDropdown />
              <UserAvatarDropdown role="seller" />
            </div>
          </div>
        </div>
      </el-header>
      <el-container style="height: calc(100vh - 64px)">
        <el-aside width="200px" style="margin-top: 3px">
          <el-menu default-active="0" @open="handleOpen" @close="handleClose" style="height: 100%">
            <RouterLink to="/seller/statistic">
              <el-menu-item index="8">
                <el-icon><TrendCharts /></el-icon>
                <span>Statistics</span>
              </el-menu-item>
            </RouterLink>
            <RouterLink to="/seller/category">
              <el-menu-item index="2">
                <el-icon><Menu /></el-icon>
                <span>Category</span>
              </el-menu-item>
            </RouterLink>
            <RouterLink to="/seller/voucher">
              <el-menu-item index="6">
                <el-icon><Ticket /></el-icon>
                <span>Voucher</span>
              </el-menu-item>
            </RouterLink>
            <RouterLink to="/seller/order">
              <el-menu-item index="7">
                <el-icon><Van /></el-icon>
                <span>Order</span>
              </el-menu-item>
            </RouterLink>
            <RouterLink to="/seller/product">
              <el-menu-item index="3">
                <el-icon><Goods /></el-icon>
                <span>Product</span>
              </el-menu-item>
            </RouterLink>
            <RouterLink to="/seller/chat">
              <el-menu-item index="4" class="menu-item-with-badge">
                <el-icon><ChatDotRound /></el-icon>
                <span>Chat</span>
                <el-badge
                  :value="unreadChatCount"
                  :max="99"
                  :hidden="unreadChatCount === 0"
                  class="badge-inline"
                />
              </el-menu-item>
            </RouterLink>
            <RouterLink to="/seller/profile">
              <el-menu-item index="5">
                <el-icon><User /></el-icon>
                <span>Profile</span>
              </el-menu-item>
            </RouterLink>
          </el-menu>
        </el-aside>
        <el-main style="background-color: #f6f6f6">
          <RouterView />
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<style scoped>
a {
  text-decoration: none;
}

/* Header Styles */
.container {
  margin: 0 auto;
  width: 100%;
}

.header {
  padding: 0 0;
  background-color: white;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.08);
  height: 64px;
  position: relative;
  z-index: 10;
}

.header-link {
  color: #666;
  text-decoration: none;
  font-size: 14px;
}

.header-link:hover {
  color: #333;
}

.header-content {
  display: flex;
  align-items: center;
  width: 100%;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-left: auto;
  margin-right: 20px;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-left: 20px;
}

.logo-icon {
  width: 40px;
  height: 40px;
  background: linear-gradient(135deg, #22c55e, #16a34a);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-text {
  display: flex;
  flex-direction: column;
}

.logo-name {
  font-size: 24px;
  font-weight: 800;
  color: #333;
  line-height: 1;
}

.logo-tagline {
  font-size: 12px;
  color: #666;
  letter-spacing: 1px;
}

/* Breadcrumb */
.breadcrumb {
  background: white;
  margin-left: 24px;
}

.breadcrumb .container {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 18px;
}

.breadcrumb-link {
  color: #666;
  text-decoration: none;
}

.breadcrumb-link:hover {
  color: #22c55e;
}

.breadcrumb-separator {
  color: #999;
  margin: 0 4px;
}

.breadcrumb-current {
  color: #333;
  font-weight: 600;
}

.menu-item-with-badge {
  display: flex !important;
  align-items: center !important;
}

.badge-inline {
  margin-left: 8px;
  top: -18px;
}

:deep(.badge-inline .el-badge__content) {
  border: none;
  box-shadow: none;
  background-color: var(--main-color) !important; /* Set to main color */
}

.el-menu-item span {
  margin-left: 8px;
}
</style>
