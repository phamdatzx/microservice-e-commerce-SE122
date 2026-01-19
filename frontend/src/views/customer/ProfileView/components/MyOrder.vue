<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { Search, ChatDotRound, Shop } from '@element-plus/icons-vue'
import { formatNumberWithDots } from '@/utils/formatNumberWithDots'
import { formatStatus } from '@/utils/formatStatus'
import axios from 'axios'

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
      console.log('Orders Response:', response.data)
      orders.value = response.data.orders
      console.log('Populated Orders:', orders.value)
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
              :src="item.image || 'https://via.placeholder.com/80'"
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
            <span class="total-label">Order Total:</span>
            <span class="total-amount">{{ formatNumberWithDots(order.total) }}đ</span>
          </div>
          <div class="order-actions">
            <template v-if="order.status === 'TO_PAY'">
              <el-button type="primary" size="large">Pay Order</el-button>
              <el-button size="large">Cancel Order</el-button>
              <el-button
                v-if="isLoggedIn"
                size="large"
                @click.stop="openChat(order.seller?.id || order.seller?._id)"
                >Contact Seller</el-button
              >
            </template>
            <template v-else-if="order.status === 'TO_CONFIRM'">
              <el-button size="large">Cancel Order</el-button>
              <el-button
                v-if="isLoggedIn"
                size="large"
                @click.stop="openChat(order.seller?.id || order.seller?._id)"
                >Contact Seller</el-button
              >
            </template>
            <template v-else-if="order.status === 'TO_PICKUP'">
              <el-button size="large">Cancel Order</el-button>
              <el-button
                v-if="isLoggedIn"
                size="large"
                @click.stop="openChat(order.seller?.id || order.seller?._id)"
                >Contact Seller</el-button
              >
            </template>
            <template v-else-if="order.status === 'SHIPPING'">
              <el-button type="primary" size="large">Order Received</el-button>
              <el-button size="large" @click.stop="openChat(order.seller?.id || order.seller?._id)"
                >Contact Seller</el-button
              >
            </template>
            <template v-else-if="order.status === 'COMPLETED'">
              <el-button type="primary" size="large">Buy Again</el-button>
              <el-button size="large">Rate</el-button>
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
</style>
