<script setup lang="ts">
import CustomerChat from './CustomerChat.vue'
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

const router = useRouter()
const isLoggedIn = ref(false)
const cartCount = ref(0)
const API_URL = 'http://localhost:81/api/order/cart/count'

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

onMounted(() => {
  isLoggedIn.value = !!localStorage.getItem('access_token')
  if (isLoggedIn.value) {
    fetchCartCount()
  }
})

const handleLogout = () => {
  localStorage.removeItem('access_token')
  isLoggedIn.value = false
  router.push('/')
}

defineExpose({
  fetchCartCount,
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

          <nav class="nav">
            <a href="#" class="nav-link">
              HOMES
              <svg width="12" height="12" viewBox="0 0 12 12" fill="none">
                <path
                  d="M3 4.5L6 7.5L9 4.5"
                  stroke="currentColor"
                  stroke-width="1.5"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
              </svg>
            </a>
            <a href="#" class="nav-link">
              PAGES
              <svg width="12" height="12" viewBox="0 0 12 12" fill="none">
                <path
                  d="M3 4.5L6 7.5L9 4.5"
                  stroke="currentColor"
                  stroke-width="1.5"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
              </svg>
            </a>
            <a href="#" class="nav-link">
              PRODUCTS
              <svg width="12" height="12" viewBox="0 0 12 12" fill="none">
                <path
                  d="M3 4.5L6 7.5L9 4.5"
                  stroke="currentColor"
                  stroke-width="1.5"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
              </svg>
            </a>
            <a href="#" class="nav-link">CONTACT</a>
          </nav>

          <div class="header-actions">
            <button class="wishlist-btn" title="Wishlist">
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
                <path
                  d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"
                ></path>
              </svg>
            </button>

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
                <span class="badge" v-if="cartCount > 0">{{ cartCount }}</span>
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
                    src="https://ui-avatars.com/api/?name=User&background=22c55e&color=fff"
                    alt="Avatar"
                    class="avatar"
                  />
                </div>
                <div class="text-content">
                  <span class="label">MY ACCOUNT</span>
                  <span class="main-text">Hello, User</span>
                </div>
              </div>

              <!-- Dropdown Menu -->
              <div class="user-dropdown">
                <RouterLink to="/profile" class="dropdown-item">
                  <svg
                    width="18"
                    height="18"
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
                  My Profile
                </RouterLink>
                <RouterLink to="/orders" class="dropdown-item">
                  <svg
                    width="18"
                    height="18"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  >
                    <path d="M6 2 3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4Z" />
                    <path d="M3 6h18" />
                    <path d="M16 10a4 4 0 0 1-8 0" />
                  </svg>
                  My Orders
                </RouterLink>
                <RouterLink to="/profile?tab=my-voucher" class="dropdown-item">
                  <svg
                    width="18"
                    height="18"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  >
                    <path
                      d="M2 9a3 3 0 0 1 0 6v2a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2v-2a3 3 0 0 1 0-6V7a2 2 0 0 0-2-2H4a2 2 0 0 0-2 2Z"
                    ></path>
                    <path d="M13 5v2"></path>
                    <path d="M13 17v2"></path>
                    <path d="M13 11v2"></path>
                  </svg>
                  My Vouchers
                </RouterLink>
                <div class="dropdown-divider"></div>
                <button @click="handleLogout" class="dropdown-item logout-btn">
                  <svg
                    width="18"
                    height="18"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  >
                    <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" />
                    <polyline points="16 17 21 12 16 7" />
                    <line x1="21" y1="12" x2="9" y2="12" />
                  </svg>
                  Logout
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Search Bar -->
    <div class="search-section">
      <div class="container">
        <div class="search-bar">
          <select class="category-select">
            <option>All Categories</option>
          </select>
          <input type="text" placeholder="Search anything..." class="search-input" />
          <button class="search-btn">
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
      </div>
    </div>
  </header>

  <!-- Breadcrumb -->
  <div class="breadcrumb">
    <div class="container">
      <a href="#" class="breadcrumb-link">Home</a>
      <span class="breadcrumb-separator">/</span>
      <a href="#" class="breadcrumb-link">pages</a>
      <span class="breadcrumb-separator">/</span>
      <span class="breadcrumb-current">login</span>
    </div>
  </div>

  <CustomerChat />
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

.nav {
  display: flex;
  gap: 32px;
  align-items: center;
}

.nav-link {
  color: #3f3f46;
  text-decoration: none;
  font-weight: 600;
  font-size: 14px;
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 0;
  transition: all 0.2s ease;
}

.nav-link:hover {
  color: #22c55e;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 24px;
}

.wishlist-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 10px;
  color: #71717a;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 12px;
  transition: all 0.2s ease;
}

.wishlist-btn:hover {
  color: #ef4444;
  background-color: #fef2f2;
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

/* Search Section */
.search-section {
  background: linear-gradient(135deg, #22c55e, #16a34a);
  padding: 24px 0;
}

.search-bar {
  display: flex;
  max-width: 700px;
  margin: 0 auto;
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  transition: transform 0.2s ease;
}

.search-bar:focus-within {
  transform: translateY(-2px);
}

.category-select {
  padding: 0 24px;
  border: none;
  border-right: 1px solid #e4e4e7;
  background: white;
  font-size: 14px;
  font-weight: 600;
  color: #3f3f46;
  cursor: pointer;
  outline: none;
}

.search-input {
  flex: 1;
  padding: 16px 24px;
  border: none;
  font-size: 14px;
  outline: none;
  color: #18181b;
}

.search-input::placeholder {
  color: #a1a1aa;
}

.search-btn {
  background: #18181b;
  border: none;
  padding: 0 28px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background-color 0.2s ease;
}

.search-btn:hover {
  background: #27272a;
}

/* Breadcrumb */
.breadcrumb {
  padding: 16px 0;
  background: #fafafa;
  border-bottom: 1px solid #f4f4f5;
}

.breadcrumb .container {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 13px;
}

.breadcrumb-link {
  color: #71717a;
  text-decoration: none;
  font-weight: 500;
  transition: color 0.2s ease;
}

.breadcrumb-link:hover {
  color: #22c55e;
}

.breadcrumb-separator {
  color: #d4d4d8;
}

.breadcrumb-current {
  color: #18181b;
  font-weight: 600;
}
</style>
