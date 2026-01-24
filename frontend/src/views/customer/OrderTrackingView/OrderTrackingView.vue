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
import NotFoundView from '@/components/NotFoundView.vue'

const router = useRouter()
const route = useRoute()
const isLoggedIn = ref(false)
const isLoading = ref(false)
const isNotFound = ref(false)

const order = ref<any>(null)

// Report Dialog State
const isReportDialogOpen = ref(false)
const isSubmittingReport = ref(false)
const selectedReportProduct = ref<any>(null)
const reportForm = ref({
  reason: '',
  description: '',
})

const reportReasons = [
  { value: 'fake_product', label: 'Fake Product' },
  { value: 'prohibited_items', label: 'Prohibited Items' },
  { value: 'counterfeit', label: 'Counterfeit Goods' },
  { value: 'scam', label: 'Scam/Fraud' },
  { value: 'other', label: 'Other' },
]

const orderInfo = computed(() => ({
  id: order.value?.id ? order.value.id.toUpperCase() : '',
  status: order.value?.status || '',
}))

const steps = ref([
  { title: 'Order Placed', time: '', completed: false, icon: Tickets, key: 'ready_to_pick' },
  { title: 'Order Paid', time: '', completed: false, icon: Wallet, key: 'paid' },
  { title: 'Handing over to carrier', time: '', completed: false, icon: Van, key: 'picked' },
  { title: 'Out for delivery', time: '', completed: false, icon: Box, key: 'delivering' },
  { title: 'Completed', time: '', completed: false, icon: CircleCheck, key: 'completed' },
])

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
    product_id: item.product_id,
    variant_id: item.variant_id,
    name: item.product_name,
    variant: item.variant_name,
    price: item.price,
    qty: item.quantity,
    img: item.image || 'https://placehold.co/80x80',
  }))
})

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

const shippingStatusMessage = computed(() => {
  const status = order.value?.status
  if (status === 'TO_PAY') {
    return 'Waiting for payment... Your order will be prepared once payment is confirmed.'
  }
  if (status === 'TO_PICKUP') {
    return 'The seller is preparing your order. Tracking info will be available soon.'
  }
  return 'Your order is being processed. Shipping information will be updated shortly.'
})

const updateStepsFromInternalStatus = (orderData: any) => {
  const createdDate = new Date(orderData.created_at).toLocaleString()

  if (steps.value[0]) {
    steps.value[0].time = createdDate
    steps.value[0].completed = true
  }

  if (steps.value[1]) {
    if (orderData.payment_method !== 'STRIPE') {
      steps.value[1].title = 'Payment Information Confirmed'
    } else {
      steps.value[1].title = 'Order Paid'
    }

    if (orderData.status !== 'TO_PAY') {
      steps.value[1].completed = true
    }

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
}

const fetchOrder = async () => {
  const orderId = route.params.id
  if (!orderId) return

  isLoading.value = true
  isNotFound.value = false
  try {
    let orderData = null

    if (history.state.orderData) {
      try {
        orderData = JSON.parse(history.state.orderData)
      } catch (e) {
        console.error('Failed to parse order state', e)
      }
    }

    if (!orderData) {
      const token = localStorage.getItem('access_token')
      const orderResponse = await axios.get(`${import.meta.env.VITE_BE_API_URL}/order/${orderId}`, {
        headers: { Authorization: `Bearer ${token}` },
      })
      orderData = orderResponse.data.data || orderResponse.data
    }

    order.value = orderData
    updateStepsFromInternalStatus(orderData)
  } catch (error: any) {
    if (error.response?.status === 404) {
      isNotFound.value = true
    }
    console.error('Error fetching order:', error)
  } finally {
    isLoading.value = false
  }
}

const handleTrackGHN = () => {
  const code = order.value?.delivery_code || order.value?.shipping_code
  if (code) {
    window.open(`https://tracking.ghn.dev/?order_code=${code}`, '_blank')
  }
}

const handleBack = () => {
  router.back()
}

const openReportDialog = (product: any) => {
  selectedReportProduct.value = product
  reportForm.value = {
    reason: '',
    description: '',
  }
  isReportDialogOpen.value = true
}

const submitReport = async () => {
  if (!reportForm.value.reason || !reportForm.value.description) {
    ElMessage.warning('Please fill in all fields')
    return
  }

  isSubmittingReport.value = true
  try {
    const payload = {
      product_id: selectedReportProduct.value.product_id,
      variant_id: selectedReportProduct.value.variant_id,
      reason: reportForm.value.reason,
      description: reportForm.value.description,
    }

    await axios.post(`${import.meta.env.VITE_BE_API_URL}/product/report`, payload, {
      headers: {
        Authorization: `Bearer ${localStorage.getItem('access_token')}`,
      },
    })

    ElMessage.success('Report submitted successfully')
    isReportDialogOpen.value = false
  } catch (error: any) {
    console.error('Failed to submit report:', error)
    ElMessage.error(error.response?.data?.message || 'Failed to submit report')
  } finally {
    isSubmittingReport.value = false
  }
}

onMounted(() => {
  isLoggedIn.value = !!localStorage.getItem('access_token')
  fetchOrder()
})
</script>

<template>
  <div v-loading="isLoading">
    <div v-if="order && !isNotFound" class="order-tracking-page">
      <div class="main-container">
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
                    ((steps.filter((s) => s.completed).length - 1) / (steps.length - 1)) * 100 +
                    '%',
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
          <!-- <el-button plain class="sec-btn">Return/Refund Request</el-button> -->
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

          <div class="delivery-card card">
            <h3 class="section-title">Shipping Information</h3>
            <div class="delivery-body">
              <div v-if="order.delivery_code || order.shipping_code" class="ghn-link-container">
                <p class="delivery-code-label">
                  GHN Tracking Code:
                  <span class="code">{{ order.delivery_code || order.shipping_code }}</span>
                </p>
                <el-button type="primary" plain :icon="Van" @click="handleTrackGHN">
                  Track on GHN Express
                </el-button>
              </div>
              <p v-else class="no-code">{{ shippingStatusMessage }}</p>
            </div>
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
              <RouterLink :to="`/product/${product.product_id}`">
                <el-image :src="product.img" class="p-img" fit="cover" />
              </RouterLink>
              <div class="p-info">
                <RouterLink :to="`/product/${product.product_id}`" class="product-link">
                  <h4 class="p-name">{{ product.name }}</h4>
                </RouterLink>
                <p class="p-variant">Variant: {{ product.variant }}</p>
                <p class="p-qty">x{{ product.qty }}</p>
                <el-button
                  v-if="order && order.status === 'COMPLETED'"
                  type="danger"
                  link
                  size="small"
                  @click="openReportDialog(product)"
                  >Report</el-button
                >
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

    <NotFoundView
      v-else-if="!isLoading && isNotFound"
      title="Order Not Found"
      message="The order you are looking for does not exist or has been removed."
    />

    <!-- Report Dialog -->
    <el-dialog v-model="isReportDialogOpen" title="Report Product" width="500px">
      <el-form label-position="top">
        <el-form-item label="Reason">
          <el-select v-model="reportForm.reason" placeholder="Select a reason" style="width: 100%">
            <el-option
              v-for="item in reportReasons"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="Description">
          <el-input
            v-model="reportForm.description"
            type="textarea"
            :rows="3"
            placeholder="Please describe the issue..."
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="isReportDialogOpen = false">Cancel</el-button>
          <el-button type="danger" :loading="isSubmittingReport" @click="submitReport">
            Submit Report
          </el-button>
        </div>
      </template>
    </el-dialog>
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

.delivery-card {
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

.delivery-card {
  border-left: 2px dashed #f5f5f5;
}

.ghn-link-container {
  display: flex;
  flex-direction: column;
  gap: 12px;
  align-items: flex-start;
}

.delivery-code-label {
  font-size: 14px;
  color: #666;
}

.delivery-code-label .code {
  font-weight: bold;
  color: var(--main-color);
  margin-left: 4px;
}

.no-code {
  color: #999;
  font-style: italic;
}

.product-link {
  text-decoration: none;
  color: inherit;
  display: block;
}

.product-link:hover .p-name {
  color: var(--main-color);
}
</style>
