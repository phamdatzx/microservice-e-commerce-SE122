<script setup lang="ts">
import CustomerChat from './CustomerChat.vue'
</script>

<template>
  <!-- Header -->
  <header class="header">
    <div class="header-top">
      <div class="container">
        <div class="header-top-left"></div>
        <div class="header-top-right">
          <a href="#" class="header-link">Sell on Swoo</a>
          <a href="#" class="header-link">Order Track</a>
        </div>
      </div>
    </div>

    <div class="header-main">
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
            <button class="wishlist-btn">
              <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
                <path
                  d="M3.172 5.172a4 4 0 015.656 0L10 6.343l1.172-1.171a4 4 0 115.656 5.656L10 17.657l-6.828-6.829a4 4 0 010-5.656z"
                  stroke="currentColor"
                  stroke-width="1.5"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
              </svg>
            </button>

            <div class="user-section">
              <span class="welcome">WELCOME</span>
              <a href="#" class="login-register">LOG IN / REGISTER</a>
            </div>

            <div class="cart">
              <span class="cart-label">CART</span>
              <div class="cart-info">
                <span class="cart-count">5</span>
                <span class="cart-total">$1,689.00</span>
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

.header-top {
  background-color: #e9ecef;
  padding: 8px 0;
  font-size: 14px;
}

.header-top .container {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-top-left {
  display: flex;
  gap: 15px;
  align-items: center;
}

.hotline {
  color: #666;
}

.phone {
  font-weight: 600;
  color: #333;
}

.header-top-right {
  display: flex;
  gap: 20px;
}

.header-link {
  color: #666;
  text-decoration: none;
  font-size: 14px;
}

.header-link:hover {
  color: #333;
}

.header-main {
  background-color: white;
  padding: 15px 0;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
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

.nav {
  display: flex;
  gap: 40px;
  align-items: center;
}

.nav-link {
  color: #333;
  text-decoration: none;
  font-weight: 600;
  font-size: 14px;
  display: flex;
  align-items: center;
  gap: 5px;
  padding: 10px 0;
}

.nav-link:hover {
  color: #22c55e;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 30px;
}

.wishlist-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 8px;
  color: #666;
}

.wishlist-btn:hover {
  color: #22c55e;
}

.user-section {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}

.welcome {
  font-size: 12px;
  color: #666;
  margin-bottom: 2px;
}

.login-register {
  color: #333;
  text-decoration: none;
  font-weight: 600;
  font-size: 14px;
}

.login-register:hover {
  color: #22c55e;
}

.cart {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}

.cart-label {
  font-size: 12px;
  color: #666;
  margin-bottom: 2px;
}

.cart-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.cart-count {
  background-color: #22c55e;
  color: white;
  border-radius: 50%;
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
}

.cart-total {
  font-weight: 600;
  color: #333;
}

/* Search Section */
.search-section {
  background: linear-gradient(135deg, #22c55e, #16a34a);
  padding: 20px 0;
}

.search-bar {
  display: flex;
  max-width: 600px;
  margin: 0 auto;
  background: white;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.category-select {
  padding: 15px 20px;
  border: none;
  background: white;
  font-size: 14px;
  color: #333;
  cursor: pointer;
  min-width: 140px;
}

.search-input {
  flex: 1;
  padding: 15px 20px;
  border: none;
  font-size: 14px;
  outline: none;
}

.search-input::placeholder {
  color: #999;
}

.search-btn {
  background: #333;
  border: none;
  padding: 15px 20px;
  cursor: pointer;
  color: white;
}

.search-btn:hover {
  background: #555;
}

/* Breadcrumb */
.breadcrumb {
  padding: 20px 0;
  background: white;
}

.breadcrumb .container {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
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
}

.breadcrumb-current {
  color: #333;
  font-weight: 600;
}
</style>
