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
  Plus,
} from '@element-plus/icons-vue'
import { ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { formatNumberWithDots } from '@/utils/formatNumberWithDots'
import axios from 'axios'
import { ElMessage, ElMessageBox, ElNotification, ElLoading } from 'element-plus'
import NotFoundView from '@/components/NotFoundView.vue'

const router = useRouter()
const route = useRoute()
const isLoggedIn = ref(false)
const isLoading = ref(false)
const isNotFound = ref(false)

const order = ref<any>(null)

// Action Handlers
const handleCancelOrder = (orderData: any) => {
  ElMessageBox.confirm('Are you sure you want to cancel this order?', 'Cancel Order', {
    confirmButtonText: 'Yes, Cancel',
    cancelButtonText: 'No',
    type: 'warning',
  })
    .then(async () => {
      try {
        await axios.put(
          `${import.meta.env.VITE_BE_API_URL}/order/${orderData.id}`,
          { status: 'CANCELLED' },
          {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('access_token')}`,
            },
          },
        )
        ElNotification({
          title: 'Success',
          message: 'Order cancelled successfully',
          type: 'success',
        })
        fetchOrder()
      } catch (error) {
        console.error('Failed to cancel order:', error)
        ElNotification({
          title: 'Error',
          message: 'Failed to cancel order',
          type: 'error',
        })
      }
    })
    .catch(() => {})
}

const handlePayOrder = async (orderData: any) => {
  const loadingInstance = ElLoading.service({
    lock: true,
    text: 'Processing Payment...',
    background: 'rgba(0, 0, 0, 0.7)',
  })
  try {
    const response = await axios.post(
      `${import.meta.env.VITE_BE_API_URL}/order/${orderData.id}/payment`,
      {},
      {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('access_token')}`,
        },
      },
    )

    if (response.data && response.data.payment_url) {
      window.location.href = response.data.payment_url
    } else {
      loadingInstance.close()
      ElNotification({
        title: 'Error',
        message: 'Failed to retrieve payment URL',
        type: 'error',
      })
    }
  } catch (error) {
    loadingInstance.close()
    console.error('Failed to initiate payment:', error)
    ElNotification({
      title: 'Error',
      message: 'Failed to initiate payment',
      type: 'error',
    })
  }
}

const handleOrderReceived = (orderData: any) => {
  ElMessageBox.confirm(
    'Are you sure you have received this order? This action cannot be undone.',
    'Confirm Receipt',
    {
      confirmButtonText: 'Yes, Received',
      cancelButtonText: 'Cancel',
      type: 'warning',
    },
  )
    .then(async () => {
      const loadingInstance = ElLoading.service({
        lock: true,
        text: 'Updating Order...',
        background: 'rgba(0, 0, 0, 0.7)',
      })
      try {
        await axios.put(
          `${import.meta.env.VITE_BE_API_URL}/order/${orderData.id}`,
          { status: 'COMPLETED' },
          {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('access_token')}`,
            },
          },
        )
        loadingInstance.close()
        ElNotification({
          title: 'Success',
          message: 'Order completed successfully',
          type: 'success',
        })
        fetchOrder()
      } catch (error) {
        loadingInstance.close()
        console.error('Failed to update order status:', error)
        ElNotification({
          title: 'Error',
          message: 'Failed to update order status',
          type: 'error',
        })
      }
    })
    .catch(() => {})
}

const handleBuyAgain = async (orderData: any) => {
  const loadingInstance = ElLoading.service({
    lock: true,
    text: 'Adding items to cart...',
    background: 'rgba(0, 0, 0, 0.7)',
  })

  try {
    const variantIds: string[] = []
    const promises = orderData.items.map(async (item: any) => {
      try {
        await axios.post(
          `${import.meta.env.VITE_BE_API_URL}/order/cart`,
          {
            seller_id: orderData.seller.id || orderData.seller._id,
            product_id: item.product_id,
            variant_id: item.variant_id,
            quantity: 1,
          },
          {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('access_token')}`,
            },
          },
        )
        variantIds.push(item.variant_id)
      } catch (err) {
        variantIds.push(item.variant_id)
      }
    })

    await Promise.all(promises)

    loadingInstance.close()

    if (variantIds.length > 0) {
      router.push({
        path: '/cart',
        query: { selected_variants: variantIds.join(',') },
      })
    } else {
      ElNotification({
        title: 'Warning',
        message: 'Could not add any items to cart (product might be unavailable).',
        type: 'warning',
      })
    }
  } catch (error) {
    loadingInstance.close()
    console.error('Failed to buy again:', error)
    ElNotification({
      title: 'Error',
      message: 'Failed to process Buy Again request',
      type: 'error',
    })
  }
}

// Rating Logic
const showRateDialog = ref(false)
const selectedOrderForRating = ref<any>(null)
const ratingForms = ref<Record<string, { star: number; content: string; fileList: any[] }>>({})

const handleRateClick = (orderData: any) => {
  selectedOrderForRating.value = orderData
  ratingForms.value = {}
  // Initialize form for each item
  orderData.items.forEach((item: any) => {
    ratingForms.value[item.variant_id] = {
      star: 5,
      content: '',
      fileList: [],
    }
  })
  showRateDialog.value = true
}

const handleRateImageChange = (uploadFile: any, variantId: string) => {
  const isJPG = uploadFile.raw.type === 'image/jpeg' || uploadFile.raw.type === 'image/jpg'
  const isPNG = uploadFile.raw.type === 'image/png'

  if (!isJPG && !isPNG) {
    ElNotification({
      title: 'Error',
      message: 'Image must be JPG, JPEG or PNG format!',
      type: 'error',
    })

    if (ratingForms.value[variantId]) {
      const list = ratingForms.value[variantId].fileList
      const index = list.indexOf(uploadFile)
      if (index > -1) {
        list.splice(index, 1)
      } else {
        const idx = list.findIndex((f) => f.uid === uploadFile.uid)
        if (idx > -1) list.splice(idx, 1)
      }
    }
    return
  }
}

const submitRating = async (item: any) => {
  const form = ratingForms.value[item.variant_id]
  if (!form) return

  if (!form.content.trim()) {
    ElNotification({
      title: 'Warning',
      message: 'Please write some review content.',
      type: 'warning',
    })
    return
  }

  const loadingInstance = ElLoading.service({
    lock: true,
    text: 'Submitting Review...',
    background: 'rgba(0, 0, 0, 0.7)',
  })

  try {
    const formData = new FormData()
    formData.append('product_id', item.product_id)
    formData.append('variant_id', item.variant_id)
    formData.append('content', form.content)
    formData.append('star', form.star.toString())

    form.fileList.forEach((file) => {
      if (file.raw) {
        formData.append('image', file.raw)
      }
    })

    await axios.post(`${import.meta.env.VITE_BE_API_URL}/product/rating`, formData, {
      headers: {
        Authorization: `Bearer ${localStorage.getItem('access_token')}`,
        'Content-Type': 'multipart/form-data',
      },
    })

    ElNotification({
      title: 'Success',
      message: 'Review submitted successfully!',
      type: 'success',
    })

    showRateDialog.value = false
  } catch (error: any) {
    console.error('Failed to submit rating:', error)
    ElNotification({
      title: 'Error',
      message: error.response?.data?.message || 'Failed to submit review',
      type: 'error',
    })
  } finally {
    loadingInstance.close()
  }
}

// Report Logic
const showReportDialog = ref(false)
const selectedOrderForReport = ref<any>(null)
const reportForms = ref<Record<string, { reason: string; description: string }>>({})

const reportReasons = [
  { label: 'Fake Product', value: 'fake_product' },
  { label: 'Quality Issue', value: 'quality_issue' },
  { label: 'Misleading Description', value: 'misleading_description' },
  { label: 'Counterfeit', value: 'counterfeit' },
  { label: 'Damaged Product', value: 'damaged_product' },
  { label: 'Other', value: 'other' },
]

const handleReportClick = (orderData: any) => {
  selectedOrderForReport.value = orderData
  reportForms.value = {}
  orderData.items.forEach((item: any) => {
    reportForms.value[item.variant_id] = {
      reason: '',
      description: '',
    }
  })
  showReportDialog.value = true
}

const submitReport = async (item: any) => {
  const form = reportForms.value[item.variant_id]
  if (!form) return

  if (!form.reason) {
    ElNotification({
      title: 'Warning',
      message: 'Please select a reason for reporting.',
      type: 'warning',
    })
    return
  }

  const loadingInstance = ElLoading.service({
    lock: true,
    text: 'Submitting Report...',
    background: 'rgba(0, 0, 0, 0.7)',
  })

  try {
    await axios.post(
      `${import.meta.env.VITE_BE_API_URL}/product/report`,
      {
        product_id: item.product_id,
        variant_id: item.variant_id,
        reason: form.reason,
        description: form.description,
      },
      {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('access_token')}`,
        },
      },
    )

    ElNotification({
      title: 'Success',
      message: 'Report submitted successfully',
      type: 'success',
    })

    form.reason = ''
    form.description = ''
    showReportDialog.value = false
  } catch (error: any) {
    console.error('Failed to submit report:', error)
    ElNotification({
      title: 'Error',
      message: error.response?.data?.message || 'Failed to submit report',
      type: 'error',
    })
  } finally {
    loadingInstance.close()
  }
}

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

const viewShop = () => {
  const sellerId = order.value?.seller?.id || order.value?.seller?._id
  if (sellerId) {
    router.push(`/seller-page/${sellerId}`)
  }
}

const handleBack = () => {
  router.back()
}
// Old report handlers removed

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

        <!-- Action Banner with dynamic buttons -->
        <div class="action-banner card" v-if="order.status !== 'CANCELLED'">
          <div class="order-actions-container">
            <template v-if="order.status === 'TO_PAY'">
              <el-button type="primary" size="large" @click="handlePayOrder(order)"
                >Pay Order</el-button
              >
              <el-button size="large" @click="handleCancelOrder(order)">Cancel Order</el-button>
              <el-button v-if="isLoggedIn" size="large" @click="openChat">Contact Seller</el-button>
            </template>
            <template v-else-if="order.status === 'TO_CONFIRM'">
              <el-button size="large" @click="handleCancelOrder(order)">Cancel Order</el-button>
              <el-button v-if="isLoggedIn" size="large" @click="openChat">Contact Seller</el-button>
            </template>
            <template v-else-if="order.status === 'TO_PICKUP'">
              <el-button size="large" @click="handleCancelOrder(order)">Cancel Order</el-button>
              <el-button v-if="isLoggedIn" size="large" @click="openChat">Contact Seller</el-button>
            </template>
            <template v-else-if="order.status === 'SHIPPING'">
              <el-button type="primary" size="large" @click="handleOrderReceived(order)"
                >Order Received</el-button
              >
              <el-button size="large" @click="openChat">Contact Seller</el-button>
            </template>
            <template v-else-if="order.status === 'COMPLETED'">
              <el-button type="primary" size="large" @click="handleBuyAgain(order)"
                >Buy Again</el-button
              >
              <el-button size="large" @click="handleRateClick(order)">Rate</el-button>
              <el-button size="large" @click="handleReportClick(order)">Report</el-button>
              <el-button size="large" @click="openChat">Contact Seller</el-button>
            </template>
            <!-- Fallback for other statuses or if logic needs default -->
            <template v-else>
              <el-button v-if="isLoggedIn" size="large" @click="openChat">Contact Seller</el-button>
            </template>
          </div>
        </div>

        <!-- Secondary Actions replaced by above but if we want to keep structure -->
        <!-- We removed secondary-actions separate div and merged into one main action area -->

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
              <el-button link :icon="Shop" class="view-shop-btn" @click="viewShop"
                >View Shop</el-button
              >
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

    <!-- Rating Dialog -->
    <el-dialog v-model="showRateDialog" title="Rate Product" width="600px" destroy-on-close>
      <div v-if="selectedOrderForRating" class="rate-product-list">
        <div v-for="item in selectedOrderForRating.items" :key="item.variant_id" class="rate-item">
          <div class="rate-item-header">
            <img
              :src="item.image || 'https://placehold.co/100'"
              alt="Product Image"
              class="rate-item-image"
            />
            <div class="rate-item-info">
              <div class="rate-item-name">{{ item.product_name }}</div>
              <div
                class="rate-item-variant"
                v-if="item.variant_name && item.variant_name.toLowerCase() !== 'default'"
              >
                Variant: {{ item.variant_name }}
              </div>
            </div>
          </div>

          <div class="rate-form" v-if="ratingForms[item.variant_id]">
            <div class="rate-stars">
              <span>Product Quality</span>
              <el-rate
                v-model="ratingForms[item.variant_id]!.star"
                size="large"
                show-text
                :texts="['Very Bad', 'Bad', 'Average', 'Good', 'Excellent']"
              />
            </div>

            <div class="rate-input">
              <el-input
                v-model="ratingForms[item.variant_id]!.content"
                type="textarea"
                :rows="3"
                placeholder="Share your experience about this product..."
              />
            </div>

            <div class="rate-upload">
              <el-upload
                v-model:file-list="ratingForms[item.variant_id]!.fileList"
                action="#"
                list-type="picture-card"
                :auto-upload="false"
                :on-change="(file: any) => handleRateImageChange(file, item.variant_id)"
                multiple
                :limit="5"
                accept=".jpg,.jpeg,.png"
              >
                <el-icon><Plus /></el-icon>
              </el-upload>
            </div>

            <div class="rate-actions">
              <el-button type="primary" @click="submitRating(item)">Submit Review</el-button>
            </div>
          </div>
          <el-divider
            v-if="item !== selectedOrderForRating.items[selectedOrderForRating.items.length - 1]"
          />
        </div>
      </div>
    </el-dialog>

    <!-- Report Dialog -->
    <el-dialog v-model="showReportDialog" title="Report Product" width="600px" destroy-on-close>
      <div v-if="selectedOrderForReport" class="report-product-list">
        <div
          v-for="item in selectedOrderForReport.items"
          :key="item.variant_id"
          class="report-item"
        >
          <div class="report-item-header">
            <img
              :src="item.image || 'https://placehold.co/100'"
              alt="Product Image"
              class="report-item-image"
            />
            <div class="report-item-info">
              <div class="report-item-name">{{ item.product_name }}</div>
              <div
                class="report-item-variant"
                v-if="item.variant_name && item.variant_name.toLowerCase() !== 'default'"
              >
                Variant: {{ item.variant_name }}
              </div>
            </div>
          </div>

          <div class="report-form" v-if="reportForms[item.variant_id]">
            <div class="report-input-group">
              <label class="report-label">Reason</label>
              <el-select
                v-model="reportForms[item.variant_id]!.reason"
                placeholder="Select a reason"
                class="report-select"
              >
                <el-option
                  v-for="option in reportReasons"
                  :key="option.value"
                  :label="option.label"
                  :value="option.value"
                />
              </el-select>
            </div>

            <div class="report-input-group">
              <label class="report-label">Description</label>
              <el-input
                v-model="reportForms[item.variant_id]!.description"
                type="textarea"
                :rows="3"
                placeholder="Please describe the issue..."
              />
            </div>

            <div class="report-actions">
              <el-button type="danger" @click="submitReport(item)">Submit Report</el-button>
            </div>
          </div>
          <el-divider
            v-if="item !== selectedOrderForReport.items[selectedOrderForReport.items.length - 1]"
          />
        </div>
      </div>
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
  justify-content: flex-end; /* Align right */
  align-items: center;
  padding: 20px 40px;
}

.order-actions-container {
  display: flex;
  gap: 12px;
  width: 100%;
  justify-content: flex-end;
}

.order-actions-container .el-button {
  min-width: 150px;
}

.received-btn {
  background-color: var(--main-color);
  border-color: var(--main-color);
  padding: 12px 30px;
  font-weight: 500;
}

/* Secondary actions removed but styles kept for reference if needed */
.secondary-actions {
  display: none; /* Hide if empty div remains */
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
