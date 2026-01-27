<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { Plus, Search, ChatDotRound, Shop } from '@element-plus/icons-vue'
import { formatNumberWithDots } from '@/utils/formatNumberWithDots'
import { formatStatus } from '@/utils/formatStatus'
import axios from 'axios'
import { ElMessageBox, ElNotification, ElLoading } from 'element-plus'

const router = useRouter()
const isLoggedIn = ref(false)

onMounted(() => {
  isLoggedIn.value = !!localStorage.getItem('access_token')
})

const handleOrderClick = (order: any) => {
  router.push({
    path: `/order-tracking/${order.id}`,
    state: { orderData: JSON.stringify(order) },
  })
}

const openChat = (sellerId: string) => {
  if (!isLoggedIn.value) return
  const event = new CustomEvent('open-chat', {
    detail: { sellerId },
  })
  window.dispatchEvent(event)
}

const activeOrderTab = ref('all')
const orders = ref<any[]>([])
const loading = ref(false)

// Pagination state
const currentPage = ref(1)
const pageSize = ref(10)
const totalItems = ref(0)

const statusMapping: Record<string, string> = {
  all: '',
  'to-pay': 'TO_PAY',
  'to-confirm': 'TO_CONFIRM',
  'to-pickup': 'TO_PICKUP',
  shipping: 'SHIPPING',
  completed: 'COMPLETED',
  cancelled: 'CANCELLED',
}

const fetchOrders = async () => {
  loading.value = true
  try {
    const params: any = {
      page: currentPage.value,
      limit: pageSize.value,
      sort_by: 'created_at',
      sort_order: 'descending',
    }

    if (activeOrderTab.value !== 'all') {
      params.status =
        statusMapping[activeOrderTab.value] || activeOrderTab.value.toUpperCase().replace(/-/g, '_')
    }

    const response = await axios.get(`${import.meta.env.VITE_BE_API_URL}/order`, {
      params,
      headers: {
        Authorization: `Bearer ${localStorage.getItem('access_token')}`,
      },
    })

    if (response.data && response.data.orders) {
      orders.value = response.data.orders

      totalItems.value = response.data.total_count
    } else {
      orders.value = []
      totalItems.value = 0
    }
  } catch (error) {
    console.error('Failed to fetch orders:', error)
  } finally {
    loading.value = false
  }
}

const handlePageChange = (page: number) => {
  currentPage.value = page
  fetchOrders()
}

watch(activeOrderTab, () => {
  currentPage.value = 1
  fetchOrders()
})

onMounted(() => {
  fetchOrders()
})

const filteredOrders = computed(() => orders.value)

const orderTabs = [
  { label: 'All', value: 'all' },
  { label: 'To Pay', value: 'to-pay' },
  { label: 'To Confirm', value: 'to-confirm' },
  { label: 'To Pickup', value: 'to-pickup' },
  { label: 'Shipping', value: 'shipping' },
  { label: 'Completed', value: 'completed' },
  { label: 'Cancelled', value: 'cancelled' },
]

const handleCancelOrder = (order: any) => {
  ElMessageBox.confirm('Are you sure you want to cancel this order?', 'Cancel Order', {
    confirmButtonText: 'Yes, Cancel',
    cancelButtonText: 'No',
    type: 'warning',
  })
    .then(async () => {
      try {
        await axios.put(
          `${import.meta.env.VITE_BE_API_URL}/order/${order.id}`,
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
        fetchOrders()
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

const handlePayOrder = async (order: any) => {
  const loadingInstance = ElLoading.service({
    lock: true,
    text: 'Processing Payment...',
    background: 'rgba(0, 0, 0, 0.7)',
  })
  try {
    const response = await axios.post(
      `${import.meta.env.VITE_BE_API_URL}/order/${order.id}/payment`,
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

const handleOrderReceived = (order: any) => {
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
          `${import.meta.env.VITE_BE_API_URL}/order/${order.id}`,
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
        fetchOrders()
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

const handleBuyAgain = async (order: any) => {
  const loadingInstance = ElLoading.service({
    lock: true,
    text: 'Adding items to cart...',
    background: 'rgba(0, 0, 0, 0.7)',
  })

  try {
    const variantIds: string[] = []
    const promises = order.items.map(async (item: any) => {
      // Assuming item.variant_id is available. If the item was deleted, this might fail, ideally backend handles checks.
      try {
        await axios.post(
          `${import.meta.env.VITE_BE_API_URL}/order/cart`,
          {
            seller_id: order.seller.id,
            product_id: item.product_id,
            variant_id: item.variant_id,
            quantity: 1, // Default to 1 for buy again
          },
          {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('access_token')}`,
            },
          },
        )
        variantIds.push(item.variant_id)
      } catch (err) {
        // Find existing cart item to still select it if it was already there?
        // Actually if post fails (e.g. already in cart?), likely we still want to select it.
        // But backend usually returns error if stock issue, or adds quantity if exists.
        // If simply adding quantity, it succeeds.
        // If error, we might still want to try selecting it (e.g. max quantity reached).
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

const handleRateClick = (order: any) => {
  selectedOrderForRating.value = order
  ratingForms.value = {}
  // Initialize form for each item
  order.items.forEach((item: any) => {
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

const handleReportClick = (order: any) => {
  selectedOrderForReport.value = order
  reportForms.value = {}
  order.items.forEach((item: any) => {
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
</script>

<template>
  <div class="order-tab">
    <el-tabs v-model="activeOrderTab" class="order-tabs">
      <el-tab-pane
        v-for="tab in orderTabs"
        :key="tab.value"
        :label="tab.label"
        :name="tab.value"
      ></el-tab-pane>
    </el-tabs>

    <div class="search-orders">
      <el-input
        placeholder="You can search by Shop name, Order ID or Product name"
        :prefix-icon="Search"
        size="large"
      />
    </div>

    <div v-loading="loading" class="order-list">
      <div v-for="order in filteredOrders" :key="order.id" class="order-card box-shadow">
        <div class="order-header">
          <div class="shop-info">
            <span class="shop-name">{{ order.seller?.name || 'Shop' }}</span>
            <el-button
              v-if="isLoggedIn"
              :icon="ChatDotRound"
              size="small"
              link
              class="shop-action-btn"
              @click.stop="openChat(order.seller?.id || order.seller?._id)"
              >Chat</el-button
            >
            <el-divider direction="vertical" />
            <el-button
              :icon="Shop"
              size="small"
              link
              class="shop-action-btn"
              @click="router.push(`/seller-page/${order.seller.id}`)"
              >View Shop</el-button
            >
          </div>
          <div class="order-status">
            <span class="shipping-status"
              >Order ID: {{ order.id.split('-')[0].toUpperCase() }}</span
            >
            <el-divider direction="vertical" />
            <span class="status-text">{{ formatStatus(order.status) }}</span>
          </div>
        </div>

        <div class="order-items">
          <div
            v-for="item in order.items"
            :key="item.variant_id"
            class="order-item"
            @click="handleOrderClick(order)"
          >
            <img
              :src="item.image || 'https://placehold.co/100'"
              :alt="item.product_name"
              class="item-image"
            />
            <div class="item-details">
              <h4 class="item-name">{{ item.product_name }}</h4>
              <p
                v-if="item.variant_name && item.variant_name.toLowerCase() !== 'default'"
                class="item-variant"
              >
                Variant: {{ item.variant_name }}
              </p>
              <p class="item-quantity">x{{ item.quantity }}</p>
            </div>
            <div class="item-price">
              <span class="current-price">{{ formatNumberWithDots(item.price) }}đ</span>
            </div>
          </div>
        </div>

        <div class="order-footer">
          <div class="total-section">
            <template v-if="order.voucher">
              <span class="total-label">Voucher Discount:</span>
              <span class="total-amount" style="color: #26aa99; margin-right: 20px">
                -{{
                  order.voucher.discount_type === 'FIXED'
                    ? formatNumberWithDots(order.voucher.discount_value) + 'đ'
                    : order.voucher.discount_value + '%'
                }}
              </span>
            </template>
            <span class="total-label">Order Total:</span>
            <span class="total-amount">{{ formatNumberWithDots(order.total) }}đ</span>
          </div>
          <div class="order-actions">
            <template v-if="order.status === 'TO_PAY'">
              <el-button type="primary" size="large" @click="handlePayOrder(order)"
                >Pay Order</el-button
              >
              <el-button size="large" @click="handleCancelOrder(order)">Cancel Order</el-button>
              <el-button
                v-if="isLoggedIn"
                size="large"
                @click.stop="openChat(order.seller?.id || order.seller?._id)"
                >Contact Seller</el-button
              >
            </template>
            <template v-else-if="order.status === 'TO_CONFIRM'">
              <el-button size="large" @click="handleCancelOrder(order)">Cancel Order</el-button>
              <el-button
                v-if="isLoggedIn"
                size="large"
                @click.stop="openChat(order.seller?.id || order.seller?._id)"
                >Contact Seller</el-button
              >
            </template>
            <template v-else-if="order.status === 'TO_PICKUP'">
              <el-button size="large" @click="handleCancelOrder(order)">Cancel Order</el-button>
              <el-button
                v-if="isLoggedIn"
                size="large"
                @click.stop="openChat(order.seller?.id || order.seller?._id)"
                >Contact Seller</el-button
              >
            </template>
            <template v-else-if="order.status === 'SHIPPING'">
              <el-button type="primary" size="large" @click="handleOrderReceived(order)"
                >Order Received</el-button
              >
              <el-button size="large" @click.stop="openChat(order.seller?.id || order.seller?._id)"
                >Contact Seller</el-button
              >
            </template>
            <template v-else-if="order.status === 'COMPLETED'">
              <el-button type="primary" size="large" @click="handleBuyAgain(order)"
                >Buy Again</el-button
              >
              <el-button v-if="!order.is_rated" size="large" @click="handleRateClick(order)"
                >Rate</el-button
              >
              <el-button v-if="!order.is_reported" size="large" @click="handleReportClick(order)"
                >Report</el-button
              >
              <el-button size="large" @click.stop="openChat(order.seller?.id || order.seller?._id)"
                >Contact Seller</el-button
              >
            </template>
          </div>
        </div>
      </div>

      <!-- Pagination -->
      <div v-if="totalItems > pageSize" class="pagination-container">
        <el-pagination
          background
          layout="prev, pager, next"
          :total="totalItems"
          :page-size="pageSize"
          size="large"
          v-model:current-page="currentPage"
          @current-change="handlePageChange"
        />
      </div>

      <div v-if="filteredOrders.length === 0 && !loading" class="no-orders text-center">
        <el-empty description="No orders found" />
      </div>
    </div>

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
.order-tabs :deep(.el-tabs__nav-wrap::after) {
  height: 1px;
}

.order-tabs :deep(.el-tabs__item) {
  padding: 0 24px;
  height: 50px;
  line-height: 50px;
  font-size: 15px;
}

.order-tabs :deep(.el-tabs__item.is-active) {
  color: var(--main-color);
}

.order-tabs :deep(.el-tabs__active-bar) {
  background-color: var(--main-color);
}

.search-orders {
  margin: 15px 0;
}

.order-card {
  background: #fff;
  margin-bottom: 15px;
  padding: 20px;
  border-radius: 4px;
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 12px;
  border-bottom: 1px solid #f5f5f5;
  margin-bottom: 15px;
}

.shop-name {
  font-weight: 600;
  font-size: 14px;
  margin-right: 12px;
}

.order-status {
  display: flex;
  align-items: center;
}

.shipping-status {
  color: #26aa99;
  font-size: 13px;
}

.status-text {
  color: var(--main-color);
  font-weight: 500;
  font-size: 13px;
}

.order-item {
  display: flex;
  padding: 12px 0;
  border-bottom: 1px solid #f5f5f5;
  cursor: pointer;
  transition: background-color 0.2s;
}

.order-item:hover {
  background-color: #fafafa;
}

.item-image {
  width: 80px;
  height: 80px;
  object-fit: cover;
  border: 1px solid #e8e8e8;
  border-radius: 2px;
  margin-right: 15px;
}

.item-details {
  flex: 1;
}

.item-name {
  font-size: 15px;
  font-weight: 400;
  color: #333;
  margin-bottom: 5px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  line-clamp: 2;
  overflow: hidden;
}

.item-variant {
  color: #888;
  font-size: 13px;
  margin-bottom: 5px;
}

.item-quantity {
  font-size: 13px;
  color: #333;
}

.item-price {
  text-align: right;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.old-price {
  text-decoration: line-through;
  color: #888;
  font-size: 13px;
  margin-bottom: 2px;
}

.current-price {
  color: var(--main-color);
  font-size: 15px;
}

.order-footer {
  padding-top: 15px;
  text-align: right;
}

.total-section {
  margin-bottom: 20px;
}

.total-label {
  font-size: 14px;
  margin-right: 10px;
}

.total-amount {
  font-size: 24px;
  color: var(--main-color);
  font-weight: 500;
}

.order-actions .el-button {
  min-width: 150px;
  padding: 0 20px;
}

.no-orders {
  padding: 60px 0;
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 28px;
  padding-bottom: 8px;
}

.pagination-container :deep(.el-pagination.is-background .el-pager li:not(.is-active)),
.pagination-container :deep(.el-pagination.is-background .btn-prev),
.pagination-container :deep(.el-pagination.is-background .btn-next) {
  background-color: #fff !important;
  border: 1px solid #e0e0e0;
}

.pagination-container :deep(.el-pagination.is-background .el-pager li.is-active) {
  background-color: var(--main-color) !important;
}

.shop-action-btn {
  color: #666;
}

.shop-action-btn:hover {
  color: var(--main-color) !important;
}

.rate-item {
  margin-bottom: 20px;
}

.rate-item-header {
  display: flex;
  margin-bottom: 15px;
}

.rate-item-image {
  width: 60px;
  height: 60px;
  object-fit: cover;
  margin-right: 15px;
  border: 1px solid #e0e0e0;
}

.rate-item-info {
  flex: 1;
}

.rate-item-name {
  font-weight: 500;
  margin-bottom: 4px;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  line-clamp: 1;
  overflow: hidden;
}

.rate-item-variant {
  font-size: 13px;
  color: #888;
}

.rate-stars {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 15px;
}

.rate-input {
  margin-bottom: 15px;
}

.rate-upload {
  margin-bottom: 15px;
}

.rate-actions {
  text-align: right;
}

/* Report Dialog Styles */
.report-item {
  margin-bottom: 20px;
}

.report-item-header {
  display: flex;
  align-items: center;
  margin-bottom: 15px;
}

.report-item-image {
  width: 60px;
  height: 60px;
  object-fit: cover;
  margin-right: 15px;
  border: 1px solid #e0e0e0;
}

.report-item-info {
  flex: 1;
  padding-right: 10px;
}

.report-item-name {
  font-weight: 500;
  margin-bottom: 4px;
  font-size: 14px;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  line-clamp: 1;
  overflow: hidden;
}

.report-item-variant {
  font-size: 13px;
  color: #888;
}

.report-form {
  background-color: #f9f9f9;
  padding: 15px;
  border-radius: 4px;
  margin-bottom: 15px;
}

.report-input-group {
  margin-bottom: 15px;
}

.report-label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  font-size: 14px;
}

.report-select {
  width: 100%;
}

.report-actions {
  text-align: right;
}
</style>
