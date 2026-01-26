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
import { Bell, Loading, Clock } from '@element-plus/icons-vue'
import NotificationDropdown from './NotificationDropdown.vue'
import { eventBus } from '@/utils/eventBus'

const router = useRouter()
const isLoggedIn = ref(false)
const cartCount = ref(0)
const userInfo = ref<any>(null)
const API_URL = `${import.meta.env.VITE_BE_API_URL}/order/cart/count`

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

const initSocket = () => {
  const token = localStorage.getItem('access_token')
  if (token) {
    socketService.connect(token)
    socketService.emit(SOCKET_EVENTS.JOIN_NOTIFICATIONS, {})
  }
}

const cleanupSocket = () => {
  // socket cleanup handled by service or dropdown
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

    initSocket()

    eventBus.on('update_user_info', () => {
      fetchUserInfo()
    })

    eventBus.on('user_logged_out', () => {
      isLoggedIn.value = false
      userInfo.value = null
      cleanupSocket()
      socketService.disconnect()
    })
  }
})

onUnmounted(() => {
  cleanupSocket()
})

const searchQuery = ref('')
const searchHistory = ref<any[]>([])
const isSearchHistoryOpen = ref(false)
const searchInputRef = ref<HTMLInputElement | null>(null)

const fetchSearchHistory = async () => {
  const token = localStorage.getItem('access_token')
  const userId = localStorage.getItem('user_id')
  if (!token || !userId) return

  try {
    const response = await axios.get(`${import.meta.env.VITE_BE_API_URL}/product/search-history`, {
      params: { limit: 20 },
      headers: {
        Authorization: `Bearer ${token}`,
        'X-User-Id': userId,
      },
    })
    if (response.data && Array.isArray(response.data)) {
      const uniqueQueries = new Set()
      searchHistory.value = response.data.filter((item: any) => {
        if (uniqueQueries.has(item.query)) return false
        uniqueQueries.add(item.query)
        return true
      })
    }
  } catch (error) {
    console.error('Error fetching search history:', error)
  }
}

const handleSearchInputFocus = () => {
  if (isLoggedIn.value) {
    fetchSearchHistory()
    isSearchHistoryOpen.value = true
  }
}

const handleSearchInputBlur = () => {
  setTimeout(() => {
    isSearchHistoryOpen.value = false
  }, 200)
}

const selectHistoryItem = (query: string) => {
  searchQuery.value = query
  handleSearch()
}

const handleSearch = () => {
  if (searchQuery.value.trim()) {
    isSearchHistoryOpen.value = false
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
})
</script>

<template>
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

          <div class="search-bar" style="position: relative; z-index: 2002">
            <input
              type="text"
              placeholder="Search anything..."
              class="search-input"
              v-model="searchQuery"
              @keyup.enter="handleSearch"
              @focus="handleSearchInputFocus"
              @blur="handleSearchInputBlur"
              ref="searchInputRef"
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

            <div v-if="isSearchHistoryOpen && searchHistory.length > 0" class="search-dropdown">
              <div class="search-history-header">Recent Searches</div>
              <div class="search-history-list">
                <div
                  v-for="(item, index) in searchHistory"
                  :key="index"
                  class="search-history-item"
                  @mousedown="selectHistoryItem(item.query)"
                >
                  <el-icon style="margin-right: 8px; color: #999"><Clock /></el-icon>
                  <span>{{ item.query }}</span>
                </div>
              </div>
            </div>
          </div>

          <div class="header-actions">
            <NotificationDropdown v-if="isLoggedIn" />

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
.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.header-main {
  background-color: white;
  padding: 20px 0;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  position: relative;
  z-index: 2001;
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

.search-dropdown {
  position: absolute;
  top: 100%;
  left: 0;
  width: 100%;
  background-color: white;
  border: 1px solid #e4e4e7;
  border-radius: 0 0 8px 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  z-index: 2000;
  margin-top: 4px;
  overflow: hidden;
}

.search-history-header {
  padding: 8px 16px;
  background-color: #f8fafc;
  color: #64748b;
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  border-bottom: 1px solid #e4e4e7;
}

.search-history-list {
  max-height: 212px;
  overflow-y: auto;
}

.search-history-list::-webkit-scrollbar {
  width: 6px;
}

.search-history-list::-webkit-scrollbar-track {
  background: transparent;
}

.search-history-list::-webkit-scrollbar-thumb {
  background: #e4e4e7;
  border-radius: 3px;
}

.search-history-list::-webkit-scrollbar-thumb:hover {
  background: #d4d4d8;
}

.search-history-item {
  padding: 10px 16px;
  display: flex;
  align-items: center;
  cursor: pointer;
  transition: background-color 0.2s;
  color: #1a1a1a;
  font-size: 14px;
}

.search-history-item:hover {
  background-color: #f1f5f9;
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
  margin-left: -6px;
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

.search-bar {
  display: flex;
  flex: 1;
  max-width: 540px;
  margin: 0 40px;
  background: #f4f4f5;
  border-radius: 7px;
  transition: all 0.2s ease;
  border: 2px solid #e4e4e7;
  position: relative;
  z-index: 2000;
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
  border-top-left-radius: 5px;
  border-bottom-left-radius: 5px;
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
  border-top-right-radius: 5px;
  border-bottom-right-radius: 5px;
}

.search-btn:hover {
  background: #22c55e;
}
</style>
