<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { Search, Download, ChatDotRound } from '@element-plus/icons-vue'
import { formatNumberWithDots } from '@/utils/formatNumberWithDots'
import { formatStatus } from '@/utils/formatStatus'
import axios from 'axios'
import { ElNotification, ElLoading } from 'element-plus'

interface User {
  id: string
  username: string
  name: string
  email: string
}

interface ProductItem {
  product_id: string
  variant_id: string
  product_name: string
  variant_name: string
  sku: string
  price: number
  image: string
  quantity: number
}

interface ShippingAddress {
  full_name: string
  phone: string
  address_line: string
  ward: string
  district: string
  province: string
  country: string
}

interface Order {
  id: string
  status: string
  user: User
  payment_method: string
  payment_status: string
  items: ProductItem[]
  created_at: string
  updated_at: string
  total: number
  shipping_address: ShippingAddress
}

const activeTab = ref('all')
const activeSubTab = ref('all')
const searchKeyword = ref('')
const dateRange = ref('')
const sortBy = ref('created_at')
const sortOrder = ref('descending')

const orders = ref<Order[]>([])
const loading = ref(false)
const totalCount = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const token = localStorage.getItem('access_token') || ''
const BE_API_URL = import.meta.env.VITE_BE_API_URL

const fetchOrders = async () => {
  loading.value = true
  try {
    const params: any = {
      page: currentPage.value,
      limit: pageSize.value,
      sort_by: sortBy.value,
      sort_order: sortOrder.value,
    }

    if (searchKeyword.value) {
      params.keyword = searchKeyword.value
    }

    if (activeTab.value !== 'all') {
      const statusMap: Record<string, string> = {
        to_confirm: 'TO_CONFIRM',
        to_pickup: 'TO_PICKUP',
        shipping: 'SHIPPING',
        completed: 'COMPLETED',
        cancelled: 'CANCELLED',
        return: 'RETURNED',
      }
      params.status = statusMap[activeTab.value]
    }

    const response = await axios.get(`${BE_API_URL}/order/seller`, {
      params,
      headers: { Authorization: `Bearer ${token}` },
    })

    orders.value = response.data.orders
    totalCount.value = response.data.total_count
  } catch (error) {
    console.error('Error fetching orders:', error)
    ElNotification.error('Failed to fetch orders')
  } finally {
    loading.value = false
  }
}

watch([activeTab, currentPage, pageSize, sortBy, sortOrder, searchKeyword], () => {
  fetchOrders()
  // Refresh counts when data might have changed or when switching tabs
  fetchAllTabCounts()
})

const tabCounts = reactive({
  all: 0,
  to_confirm: 0,
  to_pickup: 0,
  shipping: 0,
  completed: 0,
  cancelled: 0,
  return: 0,
})

const fetchAllTabCounts = async () => {
  const statuses = [
    { key: 'all', status: undefined },
    { key: 'to_confirm', status: 'TO_CONFIRM' },
    { key: 'to_pickup', status: 'TO_PICKUP' },
    { key: 'shipping', status: 'SHIPPING' },
    { key: 'completed', status: 'COMPLETED' },
    { key: 'cancelled', status: 'CANCELLED' },
    { key: 'return', status: 'RETURNED' },
  ]

  try {
    const promises = statuses.map((item) => {
      const params: any = { page: 1, limit: 1 }
      if (item.status) params.status = item.status
      return axios.get(`${BE_API_URL}/order/seller`, {
        params,
        headers: { Authorization: `Bearer ${token}` },
      })
    })

    const results = await Promise.all(promises)
    results.forEach((res, index) => {
      const statusItem = statuses[index]
      if (statusItem) {
        const key = statusItem.key as keyof typeof tabCounts
        tabCounts[key] = res.data.total_count
      }
    })
  } catch (error) {
    console.error('Error fetching tab counts:', error)
  }
}

onMounted(() => {
  fetchOrders()
  fetchAllTabCounts()
})

const handleExport = () => {
  const headers = ['Order ID', 'Customer', 'Status', 'Price', 'Payment Method', 'Created At']
  const rows = orders.value.map((order) => [
    order.id,
    order.user.name,
    order.status,
    order.total,
    order.payment_method,
    order.created_at,
  ])

  const csvContent = [
    headers.join(','),
    ...rows.map((row) =>
      row
        .map((cell) => {
          const str = String(cell ?? '')
          return str.includes(',') ? `"${str}"` : str
        })
        .join(','),
    ),
  ].join('\n')

  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
  const link = document.createElement('a')
  link.href = URL.createObjectURL(blob)
  link.download = `orders_export_${new Date().toISOString().slice(0, 10)}.csv`
  link.click()
  URL.revokeObjectURL(link.href)
}

const handleBatchDelivery = () => {
  console.log('Batch delivery clicked')
}

const handleUpdateStatus = async (orderId: string, status: string) => {
  try {
    console.log(`Updating order ${orderId} to status: ${status}`)
    const loading = ElLoading.service({
      lock: true,
      text: `Moving to ${status.replace('_', ' ')}...`,
      background: 'rgba(0, 0, 0, 0.7)',
    })
    const response = await axios.put(
      `${BE_API_URL}/order/${orderId}`,
      { status },
      { headers: { Authorization: `Bearer ${token}` } },
    )
    console.log('Update response:', response.data)
    ElNotification.success(`Order status updated to ${formatStatus(status)}`)
    fetchOrders()
    fetchAllTabCounts()
    loading.close()
  } catch (error: any) {
    console.error('Error updating order status:', error)
    ElNotification.error(error.response?.data?.message || 'Failed to update order status')
  }
}
</script>

<template>
  <div class="order-manager">
    <!-- Main Tabs -->
    <div class="main-tabs">
      <el-tabs v-model="activeTab" class="custom-tabs">
        <el-tab-pane :label="`All (${tabCounts.all})`" name="all" />
        <el-tab-pane :label="`To Confirm (${tabCounts.to_confirm})`" name="to_confirm" />
        <el-tab-pane :label="`To Pickup (${tabCounts.to_pickup})`" name="to_pickup" />
        <el-tab-pane :label="`Shipping (${tabCounts.shipping})`" name="shipping" />
        <el-tab-pane :label="`Completed (${tabCounts.completed})`" name="completed" />
        <el-tab-pane :label="`Cancelled (${tabCounts.cancelled})`" name="cancelled" />
        <el-tab-pane :label="`Return/Refund (${tabCounts.return})`" name="return" />
      </el-tabs>
    </div>

    <!-- Filters -->
    <div class="filter-section">
      <div class="filter-row">
        <el-input
          v-model="searchKeyword"
          placeholder="Search Orders"
          :prefix-icon="Search"
          class="search-input"
        />
        <div class="date-filter">
          <span>Order Date</span>
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="-"
            start-placeholder="Start Date"
            end-placeholder="End Date"
            class="date-picker"
          />
          <el-button :icon="Download" @click="handleExport">Export</el-button>
        </div>
      </div>

      <div class="sub-tabs-row">
        <div class="sub-tabs">
          <el-radio-group v-model="activeSubTab" size="small">
            <el-radio-button size="large" label="all">All ({{ totalCount }})</el-radio-button>
            <el-radio-button size="large" label="unprocessed">Unprocessed</el-radio-button>
            <el-radio-button size="large" label="processed">Processed</el-radio-button>
          </el-radio-group>
        </div>
      </div>
    </div>

    <!-- Orders List Section -->
    <div class="orders-section" v-loading="loading">
      <div class="orders-header-bar">
        <h3>{{ totalCount }} Orders Found</h3>
        <div class="orders-actions">
          <span class="sort-label">Sort by</span>
          <el-select v-model="sortOrder" placeholder="Created Date" style="width: 200px">
            <el-option label="Order Created Date: Oldest to Newest" value="ascending" />
            <el-option label="Order Created Date: Newest to Oldest" value="descending" />
          </el-select>
          <el-button
            type="primary"
            color="var(--main-color)"
            class="batch-btn"
            @click="handleBatchDelivery"
            >Mass Delivery</el-button
          >
        </div>
      </div>

      <!-- Column Headers -->
      <div class="table-header-row">
        <div class="col-product">Product(s)</div>
        <div class="col-price">Price</div>
        <div class="col-status">Status</div>
        <div class="col-countdown">Created At</div>
        <div class="col-action">Actions</div>
      </div>

      <!-- Order List -->
      <div class="order-list">
        <div v-for="order in orders" :key="order.id" class="order-card">
          <!-- Order Header -->
          <div class="order-card-header">
            <div class="user-info">
              <el-avatar :size="24" icon="User" />
              <span class="username">{{ order.user.name }}</span>
              <el-icon class="chat-icon"><ChatDotRound /></el-icon>
            </div>
            <div class="order-id">Order ID: {{ order.id }}</div>
          </div>

          <!-- Order Content -->
          <div class="order-card-content">
            <!-- Products Column -->
            <div class="col-product product-list">
              <div
                v-for="item in order.items"
                :key="item.product_id + item.variant_id"
                class="product-item"
              >
                <img
                  :src="item.image || 'https://placehold.co/60x60?text=Product'"
                  class="product-img"
                  alt="Product"
                />
                <div class="product-info">
                  <div class="product-name">{{ item.product_name }}</div>
                  <div
                    v-if="item.variant_name && item.variant_name.toLowerCase() !== 'default'"
                    class="product-variant"
                  >
                    {{ item.variant_name }}
                  </div>
                </div>
                <div class="product-qty">x{{ item.quantity }}</div>
              </div>
            </div>

            <!-- Other Columns (Vertical alignment wrapper) -->
            <div class="col-price border-left flex-center-col">
              <span class="price-text">₫{{ formatNumberWithDots(order.total) }}</span>
              <span class="sub-text">{{ order.payment_method }}</span>
            </div>

            <div class="col-status border-left flex-center-col">
              <span class="status-strong">{{ formatStatus(order.status) }}</span>
            </div>

            <div class="col-countdown border-left flex-center-col">
              <span class="countdown-text">{{ new Date(order.created_at).toLocaleString() }}</span>
            </div>

            <div class="col-action border-left flex-center-col">
              <el-button
                v-if="order.status === 'TO_CONFIRM'"
                type="primary"
                color="var(--main-color)"
                size="small"
                @click="handleUpdateStatus(order.id, 'TO_PICKUP')"
                >Confirm</el-button
              >
              <el-button
                v-if="order.status === 'TO_PICKUP'"
                type="primary"
                color="var(--main-color)"
                size="small"
                @click="handleUpdateStatus(order.id, 'SHIPPING')"
                >Ship</el-button
              >
              <el-button
                v-if="order.status === 'SHIPPING'"
                type="primary"
                color="var(--main-color)"
                size="small"
                @click="handleUpdateStatus(order.id, 'COMPLETED')"
                >Delivered</el-button
              >
              <span
                v-if="['COMPLETED', 'CANCELLED', 'RETURNED'].includes(order.status)"
                class="sub-text"
                >No Actions</span
              >
            </div>
          </div>
        </div>
      </div>

      <!-- Pagination -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="totalCount"
          @size-change="fetchOrders"
          @current-change="fetchOrders"
        />
      </div>
    </div>
  </div>
</template>

<style scoped>
.order-manager {
  padding: 20px;
  background-color: #f6f6f6;
  min-height: 100vh;
  font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif;
  color: #333;
}

/* Tabs */
.main-tabs {
  background: white;
  margin-bottom: 20px;
  padding: 0 20px;
  border-radius: 2px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}
:deep(.custom-tabs .el-tabs__header) {
  margin-bottom: 0;
}
:deep(.custom-tabs .el-tabs__item) {
  height: 50px;
  line-height: 50px;
  font-size: 15px;
}
:deep(.custom-tabs .el-tabs__item.is-active) {
  color: var(--main-color);
}
:deep(.custom-tabs .el-tabs__active-bar) {
  background-color: var(--main-color);
}

/* Filters */
.filter-section {
  background: white;
  padding: 20px;
  border-radius: 2px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
  margin-bottom: 20px;
}

.filter-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  gap: 20px;
}

.search-input {
  max-width: 400px;
}

.date-filter {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 14px;
  color: #666;
}

.date-filter span {
  white-space: nowrap;
}

.sub-tabs-row {
  display: flex;
  align-items: center;
}

/* Orders Section */
.orders-section {
  background: transparent;
}

.orders-header-bar {
  background: #fdfdfd; /* Slightly different to distinguish from card */
  padding: 15px 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #eaeaea;
  background: white;
  margin-bottom: 10px;
  border-radius: 2px;
}

.orders-header-bar h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 500;
}

.orders-actions {
  display: flex;
  align-items: center;
  gap: 15px;
}

.sort-label {
  color: #666;
  font-size: 14px;
}

.batch-btn {
  font-weight: 500;
}

/* Custom Table Layout */
.table-header-row {
  display: flex;
  background: #fafafa;
  padding: 12px 20px;
  color: #888;
  font-size: 13px;
  font-weight: 500;
  border: 1px solid #e5e5e5;
  border-bottom: none;
  margin-bottom: 12px;
}

.col-product {
  flex: 4;
}
.col-total {
  flex: 1.5;
  text-align: center;
}
.col-status {
  flex: 1.5;
  text-align: center;
}
.col-countdown {
  flex: 1.5;
  text-align: center;
}

/* Order Card */
.order-card {
  background: white;
  border: 1px solid #e5e5e5;
  margin-bottom: 12px;
  border-radius: 2px;
  overflow: hidden;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}

.order-card-header {
  padding: 12px 20px;
  border-bottom: 1px solid #f0f0f0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #fff;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 500;
  font-size: 14px;
}
.chat-icon {
  color: var(--main-color);
  cursor: pointer;
}

.order-id {
  font-size: 13px;
  color: #666;
}

.order-card-content {
  display: flex;
  min-height: 100px;
}

.product-list {
  padding: 0;
}

.product-item {
  display: flex;
  padding: 15px 20px;
  border-bottom: 1px solid #f5f5f5;
  align-items: flex-start;
  gap: 15px;
}
.product-item:last-child {
  border-bottom: none;
}

.product-img {
  width: 60px;
  height: 60px;
  object-fit: cover;
  border: 1px solid #f5f5f5;
}

.product-info {
  flex: 1;
}

.product-name {
  font-size: 14px;
  color: #333;
  line-height: 1.4;
  margin-bottom: 4px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.product-variant {
  font-size: 12px;
  color: #888;
}

.product-qty {
  font-size: 14px;
  color: #333;
}

.border-left {
  border-left: 1px solid #f5f5f5;
}

.flex-center-col {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 15px;
  text-align: center;
}

.price-text {
  font-weight: 500;
  color: #333;
  font-size: 15px;
}

.sub-text {
  font-size: 12px;
  color: #888;
  margin-top: 4px;
}

.status-strong {
  font-weight: 500;
  color: #333;
  margin-bottom: 4px;
}

.countdown-text {
  color: var(--main-color); /* Alert color for deadline */
  font-size: 13px;
}
.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 20px;
  padding: 20px;
  background: white;
  border-radius: 2px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}
</style>
