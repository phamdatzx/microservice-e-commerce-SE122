<script setup lang="ts">
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
  CircleCheck,
  Loading,
} from '@element-plus/icons-vue'
import { ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { formatNumberWithDots } from '@/utils/formatNumberWithDots'
import axios from 'axios'
import { ElMessage } from 'element-plus'

const router = useRouter()
const route = useRoute()
const isLoggedIn = ref(false)
const isLoading = ref(false)

const order = ref<any>(null)
const trackingInfo = ref<any>(null)

const orderInfo = computed(() => ({
  id: order.value?.id ? order.value.id.toUpperCase() : '',
  status: order.value?.status || '',
}))

const steps = ref([
  { title: 'Order Placed', time: '', completed: false, icon: Tickets, key: 'ready_to_pick' },
  { title: 'Order Paid', time: '', completed: false, icon: Wallet, key: 'paid' }, // 'paid' isn't a GHN status, handled separately
  { title: 'Handing over to carrier', time: '', completed: false, icon: Van, key: 'picked' },
  { title: 'Out for delivery', time: '', completed: false, icon: Box, key: 'delivering' },
  { title: 'Completed', time: '', completed: false, icon: CircleCheck, key: 'completed' }, // Internal status
])

const timeline = ref<any[]>([])

const address = computed(() => ({
  name: order.value?.shipping_address?.full_name || '',
  phone: order.value?.shipping_address?.phone || '',
  detail: order.value?.shipping_address
    ? `${order.value.shipping_address.address_line}, ${order.value.shipping_address.ward}, ${order.value.shipping_address.district}, ${order.value.shipping_address.province}`
    : '',
}))

const products = computed(() => {
  if (!order.value?.items) return []
  return order.value.items.map((item: any) => ({
    id: item.variant_id,
    name: item.product_name,
    variant: item.variant_name,
    price: item.price,
    qty: item.quantity,
    img: item.image || 'https://placehold.co/80x80',
  }))
})

const costs = computed(() => [])

const paymentMethod = computed(() => order.value?.payment_method || '')
const total = computed(() => `${formatNumberWithDots(order.value?.total || 0)}đ`)

const openChat = () => {
  if (!isLoggedIn.value || !order.value?.seller) return
  const sellerId = order.value.seller.id || order.value.seller._id
  if (sellerId) {
    const event = new CustomEvent('open-chat', {
      detail: { sellerId },
    })
    window.dispatchEvent(event)
  }
}

const fetchOrder = async () => {
  const orderId = route.params.id
  if (!orderId) return

  isLoading.value = true
  try {
    let orderData = null

    // Check for state passed from MyOrder
    if (history.state.orderData) {
      try {
        orderData = JSON.parse(history.state.orderData)
        console.log('Using order data from state:', orderData)
      } catch (e) {
        console.error('Failed to parse order state', e)
      }
    }

    if (!orderData) {
      // Fallback: Fetch internal order details
      const token = localStorage.getItem('access_token')
      const orderResponse = await axios.get(`${import.meta.env.VITE_BE_API_URL}/order/${orderId}`, {
        headers: { Authorization: `Bearer ${token}` },
      })
      orderData = orderResponse.data.data || orderResponse.data
    }

    order.value = orderData

    // Update steps based on internal status first
    updateStepsFromInternalStatus(orderData)

    // 2. Fetch GHN Tracking if shipping_code exists
    const shippingCode = orderData.shipping_code || orderData.delivery_code

    if (shippingCode) {
      await fetchGHNTracking(shippingCode)
    }
  } catch (error) {
    console.error('Error fetching order:', error)
    ElMessage.error('Failed to load order details')
  } finally {
    isLoading.value = false
  }
}

const fetchGHNTracking = async (orderCode: string) => {
  try {
    const response = await axios.post(
      'https://dev-online-gateway.ghn.vn/shiip/public-api/v2/shipping-order/detail',
      { order_code: orderCode },
      {
        headers: {
          token: import.meta.env.VITE_GHN_TOKEN,
          'Content-Type': 'application/json',
        },
      },
    )

    if (response.data.code === 200) {
      const data = response.data.data
      console.log('GHN Data:', data)
      trackingInfo.value = data

      // Update timeline from log
      let ghnTimeline: any[] = []
      if (data.log) {
        ghnTimeline = data.log
          .map((log: any) => ({
            time: new Date(log.updated_date).toLocaleString(),
            title: formatGHNStatus(log.status),
            desc: getStatusDescription(log.status),
            type: log.status === 'delivered' || log.status === 'returned' ? 'success' : 'info',
          }))
          .reverse()
      }

      // Append internal events (Order Placed, Payment)
      const internalEvents = []

      // Payment - only for STRIPE
      if (order.value.payment_method === 'STRIPE') {
        if (order.value.status !== 'TO_PAY') {
          internalEvents.push({
            time: '',
            title: 'Order Paid',
            desc: 'Payment has been confirmed',
            type: 'success',
          })
        } else {
          internalEvents.push({
            time: new Date(order.value.created_at).toLocaleString(),
            title: 'To Pay',
            desc: 'Order is waiting for payment',
            type: 'warning',
          })
        }
      }

      // Order Placed
      internalEvents.push({
        time: new Date(order.value.created_at).toLocaleString(),
        title: 'Order Placed',
        desc: 'Order has been created',
        type: 'info',
      })

      timeline.value = [...ghnTimeline, ...internalEvents]

      // Update steps based on GHN status
      updateStepsFromGHN(data)
    }
  } catch (error) {
    console.error('Error fetching tracking info:', error)
  }
}

const formatGHNStatus = (status: string) => {
  const map: Record<string, string> = {
    ready_to_pick: 'Ready to Pick',
    picking: 'Picking',
    cancel: 'Cancelled',
    picked: 'Picked',
    storing: 'Storing',
    transporting: 'Transporting',
    sorting: 'Sorting',
    delivering: 'Delivering',
    delivered: 'Delivered',
    delivery_fail: 'Delivery Failed',
    waiting_to_return: 'Waiting to Return',
    return: 'Returning',
    returned: 'Returned',
  }
  return map[status] || status
}

const getStatusDescription = (status: string) => {
  const descs: Record<string, string> = {
    ready_to_pick: 'Shipping order has just been created',
    picking: 'Shipper is coming to pick up the goods',
    cancel: 'Shipping order has been cancelled',
    picked: 'Shipper is picked the goods',
    storing: 'The goods has been shipped to GHN sorting hub',
    transporting: 'The goods are being rotated',
    sorting: 'The goods are being classified',
    delivering: 'Shipper is delivering the goods to customer',
    delivered: 'The goods has been delivered to customer',
    delivery_fail: "The goods hasn't been delivered to customer",
    waiting_to_return: 'The goods are pending delivery',
    return: 'The goods are waiting to return to seller',
    returned: 'The goods has been returned to seller',
  }
  return descs[status] || 'Status updated'
}

const updateStepsFromInternalStatus = (orderData: any) => {
  const createdDate = new Date(orderData.created_at).toLocaleString()

  if (steps.value[0]) {
    steps.value[0].time = createdDate
    steps.value[0].completed = true
  }

  // Update step title based on payment method
  if (steps.value[1]) {
    if (orderData.payment_method !== 'STRIPE') {
      steps.value[1].title = 'Payment Information Confirmed'
    } else {
      steps.value[1].title = 'Order Paid'
    }

    if (orderData.status !== 'TO_PAY') {
      steps.value[1].completed = true
    }

    // Highlight based on user request
    if (orderData.status === 'TO_PICKUP') {
      if (steps.value[2]) steps.value[2].completed = true
    } else if (orderData.status === 'SHIPPING') {
      if (steps.value[2]) steps.value[2].completed = true
      if (steps.value[3]) steps.value[3].completed = true
    } else if (orderData.status === 'COMPLETED') {
      if (steps.value[2]) steps.value[2].completed = true
      if (steps.value[3]) steps.value[3].completed = true
      if (steps.value[4]) steps.value[4].completed = true
    }
  }

  // Initialize timeline with internal events
  const internalEvents = []

  // Payment - only for STRIPE
  if (orderData.payment_method === 'STRIPE') {
    if (orderData.status !== 'TO_PAY') {
      internalEvents.push({
        time: '',
        title: 'Order Paid',
        desc: 'Payment has been confirmed',
        type: 'success',
      })
    } else {
      internalEvents.push({
        time: createdDate,
        title: 'To Pay',
        desc: 'Order is waiting for payment',
        type: 'warning',
      })
    }
  }

  internalEvents.push({
    time: createdDate,
    title: 'Order Placed',
    desc: 'Order has been created',
    type: 'info',
  })

  // Set timeline initially (GHN will prepend to this if fetched)
  timeline.value = internalEvents
}

const updateStepsFromGHN = (ghnData: any) => {
  // Map GHN logs/status to step completion
  const logs = ghnData.log || []
  const hasLog = (status: string) => logs.some((l: any) => l.status === status)

  // Handed over (picked/storing)
  const pickedLog = logs.find((l: any) => ['picked', 'storing', 'transporting'].includes(l.status))
  if (pickedLog && steps.value[2]) {
    steps.value[2].completed = true
    steps.value[2].time = new Date(pickedLog.updated_date).toLocaleString()
  }

  // Out for delivery (delivering)
  const deliveringLog = logs.find((l: any) => l.status === 'delivering')
  if (deliveringLog && steps.value[3]) {
    steps.value[3].completed = true
    steps.value[3].time = new Date(deliveringLog.updated_date).toLocaleString()
  }

  // Completed/Rating (delivered)
  const deliveredLog = logs.find((l: any) => l.status === 'delivered')
  if (deliveredLog) {
    if (steps.value[4]) steps.value[4].completed = true
    if (steps.value[3]) steps.value[3].completed = true
    if (steps.value[2]) steps.value[2].completed = true
  }
}

const handleBack = () => {
  router.back()
}

onMounted(() => {
  isLoggedIn.value = !!localStorage.getItem('access_token')
  fetchOrder()
})
</script>

<template>
  <div class="order-tracking-page" v-loading="isLoading">
    <div class="main-container" v-if="order">
      <!-- Top Action Bar -->
      <div class="tracking-top-bar">
        <el-button link :icon="ArrowLeft" @click="handleBack">BACK</el-button>
        <div class="order-meta">
          <span>ORDER ID: {{ orderInfo.id }}</span>
          <span class="status-divider">|</span>
          <span class="order-status-text">{{ orderInfo.status.replace('_', ' ') }}</span>
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
        <el-button v-if="isLoggedIn" plain class="sec-btn" @click="openChat"
          >Contact Seller</el-button
        >
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
        </div>
      </div>

      <!-- Products Section -->
      <div class="products-section card">
        <div class="shop-header">
          <div class="shop-info">
            <!-- Removals: "Favorite" badge removed -->
            <span class="shop-name">{{ order.seller?.name || 'Shop' }}</span>
            <el-button
              v-if="isLoggedIn"
              link
              :icon="ChatDotRound"
              class="chat-btn"
              @click="openChat"
              >Chat</el-button
            >
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
              <span class="current-price">{{ formatNumberWithDots(product.price) }}₫</span>
            </div>
          </div>
        </div>

        <div class="order-summary-footer">
          <div class="summary-table">
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
