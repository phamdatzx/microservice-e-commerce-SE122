<script setup lang="ts">
import { ArrowDown, Goods, Star, User, ChatDotRound } from '@element-plus/icons-vue'
import { formatNumberWithDots } from '@/utils/formatNumberWithDots'

defineProps<{
  sellerInfo: {
    id: string
    name: string
    image: string
    address: {
      province: string
    }
    sale_info: {
      follow_count: number
      rating_count: number
      rating_average: number
      product_count: number
      is_following: boolean
    }
  }
}>()

defineEmits(['toggle-follow'])
</script>

<template>
  <div class="seller-header">
    <div class="container seller-header-content">
      <!-- Shop Info Card -->
      <div class="shop-info-card">
        <div class="shop-banner-overlay"></div>
        <div class="shop-info-main">
          <div class="shop-logo-wrapper">
            <img
              :src="sellerInfo.image || '/src/assets/avatar.jpg'"
              alt="Shop Logo"
              class="shop-logo"
            />
          </div>
          <div class="shop-details">
            <h1 class="shop-name">{{ sellerInfo.name }}</h1>
            <p class="shop-status">Online 2 minutes ago</p>
            <div class="shop-actions">
              <button class="btn btn-outline-white" @click="$emit('toggle-follow')">
                <span class="icon">{{ sellerInfo.sale_info.is_following ? '✓' : '+' }}</span>
                {{ sellerInfo.sale_info.is_following ? 'Following' : 'Follow' }}
              </button>
              <button class="btn btn-outline-white">
                <el-icon><ChatDotRound /></el-icon> Chat
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Shop Stats -->
      <div class="shop-stats">
        <div class="stat-item">
          <div class="stat-value">
            {{ sellerInfo.sale_info.product_count || 0 }} <el-icon><Goods /></el-icon>
          </div>
          <div class="stat-label">Products</div>
        </div>
        <div class="stat-divider"></div>
        <div class="stat-item">
          <div class="stat-value">
            {{ formatNumberWithDots(sellerInfo.sale_info.follow_count) }}
            <el-icon><User /></el-icon>
          </div>
          <div class="stat-label">Followers</div>
        </div>
        <div class="stat-divider"></div>
        <div class="stat-item">
          <div class="stat-value highlight">
            {{ sellerInfo.sale_info.rating_average.toFixed(1) }} <el-icon><Star /></el-icon>
          </div>
          <div class="stat-label">
            Rating
            <span class="rating-count">({{ sellerInfo.sale_info.rating_count }} Ratings)</span>
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- Shop Navigation Tabs -->
  <div class="shop-nav">
    <div class="container">
      <el-row class="nav-tabs">
        <el-col :span="4" class="nav-tab active">
          <a href="#"><span class="one-line-ellipsis">Home</span></a>
        </el-col>
        <el-col :span="4" class="nav-tab">
          <a href="#"><span class="one-line-ellipsis">ALL PRODUCTS</span></a>
        </el-col>
        <el-col :span="4" class="nav-tab">
          <a href="#"><span class="one-line-ellipsis">Cargo Jeans ⭐</span></a>
        </el-col>
        <el-col :span="4" class="nav-tab">
          <a href="#"><span class="one-line-ellipsis">Skinny Jeans ⭐</span></a>
        </el-col>
        <el-col :span="4" class="nav-tab">
          <a href="#"><span class="one-line-ellipsis">WASH RETRO ...</span></a>
        </el-col>
        <el-col :span="4" class="nav-tab">
          <a href="#"><span class="one-line-ellipsis">Shorts ⭐</span></a>
        </el-col>
        <el-col :span="4" class="nav-tab dropdown-tab">
          <el-dropdown class="dropdown" trigger="hover">
            <div class="dropdown-trigger">
              <span class="one-line-ellipsis">More</span>
              <el-icon class="el-icon--right"><ArrowDown /></el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu class="custom-dropdown-menu">
                <el-dropdown-item>Top Sales</el-dropdown-item>
                <el-dropdown-item>New Arrivals</el-dropdown-item>
                <el-dropdown-item>Summer Collection</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<style scoped>
.seller-header {
  background-color: #fff;
  padding: 20px 0;
  border-bottom: 1px solid #f5f5f5;
}

.seller-header-content {
  display: flex;
  gap: 40px;
  align-items: center;
}

.shop-info-card {
  position: relative;
  width: 390px;
  height: 145px;
  border-radius: 4px;
  overflow: hidden;
  background-image: url('https://down-vn.img.susercontent.com/file/vn-11134207-7ras8-m1mzl888q4t0bf_tn'); /* Just as a fallback/background */
  background-size: cover;
  background-position: center;
  flex-shrink: 0;
}

.shop-banner-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.1);
}

.shop-info-main {
  position: relative;
  z-index: 1;
  display: flex;
  padding: 20px;
  height: 100%;
  align-items: center;
  gap: 15px;
  background-color: #e6f7ed; /* More visible Light Green */
}

.shop-logo-wrapper {
  position: relative;
  width: 80px;
  height: 80px;
  flex-shrink: 0;
}

.shop-logo {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  border: 4px solid rgba(255, 255, 255, 0.3);
  object-fit: cover;
}

.shop-details {
  color: #333;
  flex: 1;
}

.shop-name {
  font-size: 20px;
  font-weight: 500;
  margin: 0;
  color: #222;
}

.shop-status {
  font-size: 12px;
  color: #666;
  margin: 4px 0 10px 0;
}

.shop-actions {
  display: flex;
  gap: 10px;
}

.btn {
  height: 30px;
  padding: 0 12px;
  border-radius: 2px;
  font-size: 12px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  transition: all 0.2s;
  box-sizing: border-box;
}

.btn-outline-white {
  background: #fff;
  border: 1px solid var(--main-color);
  color: var(--main-color);
}

.btn-outline-white:hover {
  background: #f0f7f4;
}

.shop-stats {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0;
  padding-left: 20px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  padding: 0 40px;
}

.stat-item:first-child {
  padding-left: 0;
}

.stat-divider {
  width: 1px;
  height: 40px;
  background-color: #eee;
}

.stat-label {
  font-size: 13px;
  color: #999;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-top: 4px;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: var(--main-color);
  line-height: 1;
  display: flex;
  align-items: center;
  gap: 4px;
}

.rating-count {
  font-size: 14px;
  font-weight: 400;
  color: #999;
}

/* Shop Nav */
.shop-nav {
  background-color: #fff;
  border-bottom: 1px solid #e5e5e5;
  position: sticky;
  top: 0;
  z-index: 10;
}

.nav-tabs {
  display: flex;
  height: 60px;
  align-items: center;
}

.nav-tab {
  flex: 1;
  height: 100%;
  display: flex;
  align-items: center;
  text-decoration: none;
  color: #222;
  font-size: 14px;
  text-transform: uppercase;
  border-bottom: 3px solid transparent;
  transition: all 0.2s;
  display: flex;
  justify-content: center;

  a,
  .dropdown {
    display: block;
    width: 100%;
    height: 100%;
    transition: all 0.2s;
    display: flex;
    justify-content: center;
    align-items: center;

    padding: 0 20px;

    &:hover {
      color: var(--main-color);
    }
  }
}

.nav-tab.active {
  color: var(--main-color);
  border-bottom-color: var(--main-color);
}

.dropdown-tab {
  cursor: pointer;
  outline: none;
}

.dropdown-trigger {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  padding: 0 20px;
  color: #222;
  transition: all 0.2s;

  &:hover {
    color: var(--main-color);
  }
}

:deep(.custom-dropdown-menu) {
  padding: 10px 0;
  border-radius: 4px;
  border: none;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);

  .el-dropdown-menu__item {
    font-size: 14px;
    padding: 10px 25px;
    color: #333;

    &:hover {
      background-color: #f5f5f5;
      color: var(--main-color);
    }
  }
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 15px;
}
</style>
