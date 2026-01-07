<script setup lang="ts">
import Header from '@/components/Header.vue'
import {
  ArrowLeft,
  Location,
  ChatDotRound,
  Shop,
  Tickets,
  Wallet,
  Van,
  Box,
  Star,
} from '@element-plus/icons-vue'
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// #region Mock Data
const orderInfo = {
  id: '251229N1GDEMAU',
  status: 'DELIVERED SUCCESSFULLY',
}

const steps = [
  { title: 'Order Placed', time: '22:59 29-12-2025', completed: true, icon: Tickets },
  { title: 'Order Paid (95,285₫)', time: '23:00 29-12-2025', completed: true, icon: Wallet },
  { title: 'Handed over to carrier', time: '15:38 30-12-2025', completed: true, icon: Van },
  { title: 'Out for delivery', time: '', completed: true, icon: Box },
  { title: 'Rating', time: '', completed: false, icon: Star },
]

const timeline = [
  { time: '16:48 03-01-2026', title: 'Delivered', desc: 'Delivered successfully', type: 'success' },
  {
    time: '13:47 03-01-2026',
    title: 'In transit',
    desc: 'Your order will be delivered soon, please stay reachable.',
    type: 'info',
  },
  {
    time: '12:08 02-01-2026',
    title: 'In transit',
    desc: 'The order has arrived at the local sorting facility and will be delivered within the next 12 hours.',
    type: 'info',
  },
  {
    time: '06:37 31-12-2025',
    title: 'In transit',
    desc: 'The order has arrived at the local sorting facility and will be delivered within the next 12 hours.',
    type: 'info',
  },
  {
    time: '20:19 30-12-2025',
    title: 'In transit',
    desc: 'The order has arrived at the warehouse.',
    type: 'info',
  },
  {
    time: '18:55 30-12-2025',
    title: 'In transit',
    desc: 'The order has arrived at the post office.',
    type: 'info',
  },
]

const address = {
  name: 'Johnathan Doe',
  phone: '(+84) 912 345 678',
  detail: '123 Nguyen Hue Street, Ben Nghe Ward, District 1, Ho Chi Minh City',
}

const products = [
  {
    id: 1,
    name: 'Tempered Glass Screen Protector - Anti-Blue Light for iPhone 17 air 16 15 14 13 12 11 Pro max',
    variant: 'HD, 12/12Pro',
    price: 28800,
    oldPrice: 50000,
    qty: 1,
    img: 'https://placehold.co/80x80',
  },
  {
    id: 2,
    name: 'Woven Wrist Strap for Phone Case - Compatible with iPhone OPPO Samsung Xiaomi',
    variant: 'Pink',
    price: 32800,
    oldPrice: 57600,
    qty: 1,
    img: 'https://placehold.co/80x80',
  },
  {
    id: 3,
    name: '【FAST SHIPPING】Simple Clear Case, Non-yellowing for iPhone 17 air 16 15 14 13 12 11 Pro max',
    variant: 'Transparent, 12Pro',
    price: 52500,
    oldPrice: 80000,
    qty: 1,
    img: 'https://placehold.co/80x80',
  },
]

const costs = [
  { label: 'Merchandise Subtotal', value: '114,100₫' },
  { label: 'Shipping Fee', value: '22,200₫' },
  { label: 'Shipping Discount', value: '-22,200₫', icon: true },
  { label: 'Platform Voucher', value: '-16,815₫' },
  { label: 'Shop Voucher', value: '-2,000₫' },
]

const paymentMethod = 'Bank Account'
const total = '95,285₫'
// #endregion

const handleBack = () => {
  router.back()
}
</script>

<template>
  <Header />

  <div class="order-tracking-page">
    <div class="main-container">
      <!-- Top Action Bar -->
      <div class="tracking-top-bar">
        <el-button link :icon="ArrowLeft" @click="handleBack">BACK</el-button>
        <div class="order-meta">
          <span>ORDER ID. {{ orderInfo.id }}</span>
          <span class="status-divider">|</span>
          <span class="order-status-text">{{ orderInfo.status }}</span>
        </div>
      </div>

      <!-- Stepper Section -->
      <div class="stepper-section card">
        <div class="stepper-container">
          <!-- Background progress bar -->
          <div class="stepper-progress-bg">
            <div
              class="stepper-progress-fill"
              :style="{
                width:
                  ((steps.filter((s) => s.completed).length - 1) / (steps.length - 1)) * 100 + '%',
              }"
            ></div>
          </div>

          <div
            v-for="(step, index) in steps"
            :key="index"
            class="step-item"
            :class="{ active: step.completed }"
          >
            <div class="step-icon">
              <el-icon><component :is="step.icon" /></el-icon>
            </div>
            <div class="step-content">
              <div class="step-title">{{ step.title }}</div>
              <div class="step-time">{{ step.time }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Action Banner -->
      <div class="action-banner card">
        <p class="banner-text">
          Please check all products in the order carefully before clicking "Order Received".
        </p>
        <div class="banner-actions">
          <el-button type="primary" class="received-btn">Order Received</el-button>
        </div>
      </div>

      <!-- Secondary Actions -->
      <div class="secondary-actions">
        <el-button plain class="sec-btn">Return/Refund Request</el-button>
        <el-button plain class="sec-btn">Contact Seller</el-button>
      </div>

      <!-- Delivery / Address Info -->
      <div class="info-grid">
        <div class="address-card card">
          <h3 class="section-title">Delivery Address</h3>
          <div class="address-body">
            <p class="user-name">{{ address.name }}</p>
            <p class="user-phone">{{ address.phone }}</p>
            <p class="address-text">{{ address.detail }}</p>
          </div>
        </div>

        <div class="timeline-card card">
          <div class="timeline-header">
            <span class="carrier-name">SPX Express</span>
            <span class="tracking-id">SPXVN05293562230C</span>
          </div>
          <el-timeline>
            <el-timeline-item
              v-for="(item, index) in timeline"
              :key="index"
              :timestamp="item.time"
              placement="top"
              :type="item.type"
              :hollow="index !== 0"
            >
              <div class="timeline-content">
                <span class="timeline-title">{{ item.title }}</span>
                <p class="timeline-desc">{{ item.desc }}</p>
                <!-- Removals: "View delivery image" link removed -->
              </div>
            </el-timeline-item>
          </el-timeline>
          <div class="timeline-footer">
            <el-button link type="primary">View more</el-button>
          </div>
        </div>
      </div>

      <!-- Products Section -->
      <div class="products-section card">
        <div class="shop-header">
          <div class="shop-info">
            <!-- Removals: "Favorite" badge removed -->
            <span class="shop-name">Coolsen Electronics</span>
            <el-button link :icon="ChatDotRound" class="chat-btn">Chat</el-button>
            <el-button link :icon="Shop" class="view-shop-btn">View Shop</el-button>
          </div>
        </div>

        <div class="product-list">
          <div v-for="product in products" :key="product.id" class="product-item">
            <el-image :src="product.img" class="p-img" fit="cover" />
            <div class="p-info">
              <h4 class="p-name">{{ product.name }}</h4>
              <p class="p-variant">Variant: {{ product.variant }}</p>
              <p class="p-qty">x{{ product.qty }}</p>
            </div>
            <div class="p-prices">
              <span class="old-price">{{ product.oldPrice.toLocaleString() }}₫</span>
              <span class="current-price">{{ product.price.toLocaleString() }}₫</span>
            </div>
          </div>
        </div>

        <div class="order-summary-footer">
          <div class="summary-table">
            <div v-for="cost in costs" :key="cost.label" class="summary-line">
              <span class="summary-label">
                {{ cost.label }}
                <el-icon v-if="cost.icon" class="info-icon"><Location /></el-icon>
              </span>
              <span class="summary-value">{{ cost.value }}</span>
            </div>
            <div class="summary-line total-line">
              <span class="summary-label">Total Payment</span>
              <span class="total-value">{{ total }}</span>
            </div>
            <div class="summary-line payment-line">
              <span class="summary-label">Payment Method</span>
              <span class="payment-value">{{ paymentMethod }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.order-tracking-page {
  background-color: #f5f5f5;
  padding: 20px 0 60px;
  min-height: 100vh;
}

.card {
  background: #fff;
  border-radius: 4px;
  padding: 24px;
  margin-bottom: 20px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}

.tracking-top-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 0 20px;
  font-size: 14px;
}

.order-meta {
  color: #333;
}

.status-divider {
  margin: 0 10px;
  color: #ddd;
}

.order-status-text {
  color: var(--main-color);
  font-weight: 500;
}

/* Stepper */
.stepper-section {
  padding: 40px 20px;
}

.stepper-container {
  display: flex;
  justify-content: space-between;
  position: relative;
}

.stepper-progress-bg {
  position: absolute;
  top: 24px;
  left: 0;
  right: 0;
  height: 4px;
  background-color: #f0f0f0;
  margin: 0 10%; /* Adjust based on step item width to center between first and last icons */
  z-index: 0;
}

.stepper-progress-fill {
  height: 100%;
  background-color: #26aa99;
  transition: width 0.3s ease;
}

.step-item {
  flex: 1;
  text-align: center;
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  z-index: 1;
}

.step-icon {
  width: 56px;
  height: 56px;
  border: 4px solid #f0f0f0;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #fff;
  color: #ccc;
  font-size: 28px;
  margin-bottom: 12px;
  transition: all 0.3s;
}

.step-item.active .step-icon {
  border-color: #26aa99;
  color: #26aa99;
}

.step-title {
  font-size: 14px;
  color: #333;
  margin-bottom: 4px;
}

.step-time {
  font-size: 12px;
  color: #999;
}

/* Action Banner */
.action-banner {
  background-color: #fffcf5;
  border: 1px solid #ffdecb;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 40px;
}

.banner-text {
  color: #888;
  font-size: 14px;
}

.received-btn {
  background-color: var(--main-color);
  border-color: var(--main-color);
  padding: 12px 30px;
  font-weight: 500;
}

.secondary-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-bottom: 20px;
}

.sec-btn {
  flex: 0 0 auto;
  min-width: 140px;
}

/* Info Grid */
.info-grid {
  display: grid;
  grid-template-columns: 1fr 2fr;
  gap: 20px;
}

.section-title {
  font-size: 20px;
  font-weight: 500;
  margin-bottom: 24px;
  color: #333;
}

.address-body p {
  margin: 8px 0;
  line-height: 1.5;
}

.user-name {
  font-weight: 500;
}
.address-text {
  color: #888;
}

.timeline-card {
  border-left: 2px dashed #f5f5f5;
}

.timeline-header {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-bottom: 20px;
  color: #888;
  font-size: 13px;
}

.timeline-content {
  padding-left: 10px;
}

.timeline-title {
  font-weight: 500;
  color: #26aa99;
}

.timeline-desc {
  color: #888;
  font-size: 13px;
  margin: 4px 0;
}

/* Products Section */
.shop-header {
  border-bottom: 1px solid #f5f5f5;
  padding-bottom: 15px;
  margin-bottom: 15px;
}

.shop-name {
  font-weight: bold;
  margin-right: 15px;
}

.chat-btn,
.view-shop-btn {
  border: 1px solid #ddd;
  padding: 4px 12px;
  margin-left: 8px;
}

.product-item {
  display: flex;
  padding: 15px 0;
  border-bottom: 1px solid #f9f9f9;
}

.p-img {
  width: 80px;
  height: 80px;
  border: 1px solid #eee;
}

.p-info {
  flex: 1;
  padding: 0 15px;
}

.p-name {
  font-size: 14px;
  line-height: 1.4;
  margin-bottom: 8px;
  font-weight: normal;
}

.p-variant,
.p-qty {
  color: #999;
  font-size: 12px;
  margin-bottom: 4px;
}

.p-prices {
  text-align: right;
  width: 120px;
}

.old-price {
  color: #ccc;
  text-decoration: line-through;
  font-size: 12px;
  margin-right: 8px;
}

.current-price {
  color: var(--main-color);
}

/* Cost Summary */
.order-summary-footer {
  background-color: #fffcf5;
  margin: 20px -24px -24px;
  padding: 24px;
  border-top: 1px dashed #ffdecb;
}

.summary-table {
  max-width: 400px;
  margin-left: auto;
}

.summary-line {
  display: flex;
  justify-content: space-between;
  margin-bottom: 15px;
  font-size: 13px;
  color: #888;
}

.total-line {
  border-top: 1px solid #eee;
  padding-top: 15px;
  font-size: 16px;
  color: #333;
}

.total-value {
  color: var(--main-color);
  font-size: 24px;
  font-weight: bold;
}

.payment-line {
  margin-top: 10px;
}

.info-icon {
  font-size: 14px;
  vertical-align: middle;
  color: var(--main-color);
}
</style>
