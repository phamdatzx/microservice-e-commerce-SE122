<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { Search, Download, ChatDotRound } from '@element-plus/icons-vue'
import { formatNumberWithDots } from '@/utils/formatNumberWithDots'

// --- Interfaces ---
interface Product {
  id: string
  name: string
  image: string
  variant: string
  quantity: number
  price: number
}

interface Order {
  id: string
  customerName: string
  customerAvatar?: string
  status: 'TO_PICKUP' | 'TO_CONFIRM' | 'SHIPPING' | 'COMPLETED' | 'CANCELLED' | 'RETURNED'
  totalAmount: number
  paymentMethod: string
  countdown: string
  confirmedDate: string // ISO 8601

  isProcessed: boolean
  products: Product[]
}

// --- State ---
const activeTab = ref('to_pickup')
const activeSubTab = ref('all')
const searchKeyword = ref('')
const dateRange = ref('')
const sortBy = ref('confirmed_date_desc')

// Mock Data
const orders = ref<Order[]>([
  {
    id: '210803JQCW7ZM3',
    customerName: 'tuyenmrxv',
    customerAvatar: 'https://placehold.co/30x30/orange/white?text=T',
    status: 'TO_PICKUP',
    totalAmount: 251500,
    paymentMethod: 'COD',
    countdown: 'Prepare by 23/08/2025',
    confirmedDate: '2025-08-20T10:00:00',

    isProcessed: false,
    products: [
      {
        id: 'p1',
        name: '[WHOLESALE] Infrared Thermometer for Baby, Milk, Food',
        image: 'https://placehold.co/60x60?text=Thermo',
        variant: 'Model 1502',
        quantity: 1,
        price: 251500,
      },
    ],
  },
  {
    id: '210731C6RXE3VN',
    customerName: 'myny1997',
    customerAvatar: 'https://placehold.co/30x30/gray/white?text=M',
    status: 'TO_PICKUP',
    totalAmount: 177658,
    paymentMethod: 'ShopeePay',
    countdown: 'Prepare by 20/08/2025',
    confirmedDate: '2025-08-18T14:30:00',
    isProcessed: true,
    products: [
      {
        id: 'p2',
        name: '[GOOD] Electronic Kitchen Scale for Baking, Diet',
        image: 'https://placehold.co/60x60?text=Scale',
        variant: 'Electronic Scale',
        quantity: 1,
        price: 90000,
      },
      {
        id: 'p3',
        name: '[WARRANTY 6M] Silent Computer Mouse Bluetooth Q5',
        image: 'https://placehold.co/60x60?text=Mouse',
        variant: 'Black Q1',
        quantity: 1,
        price: 87658,
      },
      {
        id: 'p4',
        name: 'Gift - Ceramic Knife',
        image: 'https://placehold.co/60x60?text=Knife',
        variant: 'Free Gift',
        quantity: 1,
        price: 0,
      },
    ],
  },
  {
    id: '210725RGKPVV5N',
    customerName: 'xe2_pzc2a1',
    customerAvatar: 'https://placehold.co/30x30/blue/white?text=X',
    status: 'TO_PICKUP',
    totalAmount: 124300,
    paymentMethod: 'COD',
    countdown: 'Prepare by 13/08/2025',
    confirmedDate: '2025-08-10T09:15:00',
    isProcessed: false,
    products: [
      {
        id: 'p5',
        name: '[GENUINE] Kemei Hair Straightener 4 Mode',
        image: 'https://placehold.co/60x60?text=Hair',
        variant: 'Kemei KM',
        quantity: 1,
        price: 124300,
      },
    ],
  },
  // ... existing mock data ...
  {
    id: 'ORDER_CONFIRM_1',
    customerName: 'new_buyer',
    customerAvatar: 'https://placehold.co/30x30/green/white?text=N',
    status: 'TO_CONFIRM',
    totalAmount: 500000,
    paymentMethod: 'COD',
    countdown: 'Confirm by 24/08/2025',
    confirmedDate: '2025-08-21T08:00:00',
    isProcessed: false,
    products: [
      {
        id: 'p6',
        name: 'Gaming Headset 7.1 Surround',
        image: 'https://placehold.co/60x60?text=Headset',
        variant: 'Black/Red',
        quantity: 1,
        price: 500000,
      },
    ],
  },
  {
    id: 'ORDER_SHIP_1',
    customerName: 'loyal_cust',
    customerAvatar: 'https://placehold.co/30x30/purple/white?text=L',
    status: 'SHIPPING',
    totalAmount: 150000,
    paymentMethod: 'ShopeePay',
    countdown: 'Estimated 25/08/2025',
    confirmedDate: '2025-08-19T11:45:00',
    isProcessed: true,
    products: [
      {
        id: 'p7',
        name: 'Wireless Charger Pad',
        image: 'https://placehold.co/60x60?text=Charger',
        variant: 'White',
        quantity: 2,
        price: 75000,
      },
    ],
  },
  {
    id: 'ORDER_COMP_1',
    customerName: 'happy_user',
    customerAvatar: 'https://placehold.co/30x30/blue/white?text=H',
    status: 'COMPLETED',
    totalAmount: 1200000,
    paymentMethod: 'Credit Card',
    countdown: 'Completed 20/08/2025',
    confirmedDate: '2025-08-15T16:20:00',
    isProcessed: true,
    products: [
      {
        id: 'p8',
        name: 'Mechanical Keyboard Blue Switch',
        image: 'https://placehold.co/60x60?text=Keyboard',
        variant: 'RGB',
        quantity: 1,
        price: 1200000,
      },
    ],
  },
  {
    id: 'ORDER_CANC_1',
    customerName: 'cancel_user',
    customerAvatar: 'https://placehold.co/30x30/red/white?text=C',
    status: 'CANCELLED',
    totalAmount: 45000,
    paymentMethod: 'COD',
    countdown: 'Cancelled by User',
    confirmedDate: '2025-08-12T13:10:00',
    isProcessed: true,
    products: [
      {
        id: 'p9',
        name: 'Phone Case iPhone 13',
        image: 'https://placehold.co/60x60?text=Case',
        variant: 'Transparent',
        quantity: 1,
        price: 45000,
      },
    ],
  },
  {
    id: 'ORDER_RET_1',
    customerName: 'return_guy',
    customerAvatar: 'https://placehold.co/30x30/yellow/white?text=R',
    status: 'RETURNED',
    totalAmount: 300000,
    paymentMethod: 'COD',
    countdown: 'Return requested',
    confirmedDate: '2025-08-22T09:50:00',
    isProcessed: false,
    products: [
      {
        id: 'p10',
        name: 'Running Shoes Size 42',
        image: 'https://placehold.co/60x60?text=Shoes',
        variant: 'Blue',
        quantity: 1,
        price: 300000,
      },
    ],
  },
])

const filteredOrders = computed(() => {
  let result = orders.value

  // 1. Filter by Tab
  if (activeTab.value !== 'all') {
    const statusMap: Record<string, string> = {
      to_confirm: 'TO_CONFIRM',
      to_pickup: 'TO_PICKUP',
      shipping: 'SHIPPING',
      completed: 'COMPLETED',
      cancelled: 'CANCELLED',
      return: 'RETURNED',
    }
    result = result.filter((o) => o.status === statusMap[activeTab.value])
  }

  // 2. Filter by Search Keyword
  if (searchKeyword.value) {
    const kw = searchKeyword.value.toLowerCase()
    result = result.filter(
      (o) =>
        o.id.toLowerCase().includes(kw) ||
        o.customerName.toLowerCase().includes(kw) ||
        o.products.some((p) => p.name.toLowerCase().includes(kw)),
    )
  }

  // 3. Filter by Date Range
  if (dateRange.value && Array.isArray(dateRange.value) && dateRange.value.length === 2) {
    const start = new Date(dateRange.value[0])
    const end = new Date(dateRange.value[1])
    // Normalize end date to end of day if needed, or rely on picker behavior
    // For simplicity, comparing timestamps
    result = result.filter((o) => {
      const orderDate = new Date(o.confirmedDate)
      return orderDate >= start && orderDate <= end
    })
  }

  // 4. Filter by Sub-tab
  if (activeSubTab.value === 'unprocessed') {
    result = result.filter((o) => !o.isProcessed)
  } else if (activeSubTab.value === 'processed') {
    result = result.filter((o) => o.isProcessed)
  }

  // 5. Sort
  result = [...result].sort((a, b) => {
    const dateA = new Date(a.confirmedDate).getTime()
    const dateB = new Date(b.confirmedDate).getTime()
    if (sortBy.value === 'confirmed_date_asc') {
      return dateA - dateB
    } else {
      return dateB - dateA
    }
  })

  return result
})

const subTabCounts = computed(() => {
  // Base list filtered by Tab, Search, Date (same as filteredOrders before Sub-tab step)
  let result = orders.value

  // 1. Tab
  if (activeTab.value !== 'all') {
    const statusMap: Record<string, string> = {
      to_confirm: 'TO_CONFIRM',
      to_pickup: 'TO_PICKUP',
      shipping: 'SHIPPING',
      completed: 'COMPLETED',
      cancelled: 'CANCELLED',
      return: 'RETURNED',
    }
    result = result.filter((o) => o.status === statusMap[activeTab.value])
  }

  // 2. Search
  if (searchKeyword.value) {
    const kw = searchKeyword.value.toLowerCase()
    result = result.filter(
      (o) =>
        o.id.toLowerCase().includes(kw) ||
        o.customerName.toLowerCase().includes(kw) ||
        o.products.some((p) => p.name.toLowerCase().includes(kw)),
    )
  }

  // 3. Date
  if (dateRange.value && Array.isArray(dateRange.value) && dateRange.value.length === 2) {
    const start = new Date(dateRange.value[0])
    const end = new Date(dateRange.value[1])
    result = result.filter((o) => {
      const orderDate = new Date(o.confirmedDate)
      return orderDate >= start && orderDate <= end
    })
  }

  const allCount = result.length
  const unprocessedCount = result.filter((o) => !o.isProcessed).length
  const processedCount = result.filter((o) => o.isProcessed).length

  return { all: allCount, unprocessed: unprocessedCount, processed: processedCount }
})

const tabCounts = computed(() => {
  const counts = {
    all: orders.value.length,
    to_confirm: 0,
    to_pickup: 0,
    shipping: 0,
    completed: 0,
    cancelled: 0,
    return: 0,
  }

  orders.value.forEach((order) => {
    switch (order.status) {
      case 'TO_CONFIRM':
        counts.to_confirm++
        break
      case 'TO_PICKUP':
        counts.to_pickup++
        break
      case 'SHIPPING':
        counts.shipping++
        break
      case 'COMPLETED':
        counts.completed++
        break
      case 'CANCELLED':
        counts.cancelled++
        break
      case 'RETURNED':
        counts.return++
        break
    }
  })
  return counts
})

const handleExport = () => {
  const headers = [
    'Order ID',
    'Customer',
    'Status',
    'Total Amount',
    'Payment Method',
    'Confirmed Date',
  ]
  const rows = filteredOrders.value.map((order) => [
    order.id,
    order.customerName,
    order.status,
    order.totalAmount,
    order.paymentMethod,
    order.confirmedDate,
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
</script>

<template>
  <div class="shipping-manager">
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
            <el-radio-button size="large" label="all">All ({{ subTabCounts.all }})</el-radio-button>
            <el-radio-button size="large" label="unprocessed"
              >Unprocessed ({{ subTabCounts.unprocessed }})</el-radio-button
            >
            <el-radio-button size="large" label="processed"
              >Processed ({{ subTabCounts.processed }})</el-radio-button
            >
          </el-radio-group>
        </div>
      </div>
    </div>

    <!-- Orders List Section -->
    <div class="orders-section">
      <div class="orders-header-bar">
        <h3>{{ filteredOrders.length }} Orders</h3>
        <div class="orders-actions">
          <span class="sort-label">Sort by</span>
          <el-select v-model="sortBy" placeholder="Confirmed Date" style="width: 200px">
            <el-option label="Order Confirmed Date: Oldest to Newest" value="confirmed_date_asc" />
            <el-option label="Order Confirmed Date: Newest to Oldest" value="confirmed_date_desc" />
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
        <div class="col-total">Total Amount</div>
        <div class="col-status">Status</div>
        <div class="col-countdown">Countdown</div>
      </div>

      <!-- Order List -->
      <div class="order-list">
        <div v-for="order in filteredOrders" :key="order.id" class="order-card">
          <!-- Order Header -->
          <div class="order-card-header">
            <div class="user-info">
              <el-avatar :size="24" :src="order.customerAvatar" />
              <span class="username">{{ order.customerName }}</span>
              <el-icon class="chat-icon"><ChatDotRound /></el-icon>
            </div>
            <div class="order-id">Order ID: {{ order.id }}</div>
          </div>

          <!-- Order Content -->
          <div class="order-card-content">
            <!-- Products Column -->
            <div class="col-product product-list">
              <div v-for="product in order.products" :key="product.id" class="product-item">
                <img :src="product.image" class="product-img" alt="Product" />
                <div class="product-info">
                  <div class="product-name">{{ product.name }}</div>
                  <div class="product-variant">{{ product.variant }}</div>
                </div>
                <div class="product-qty">x{{ product.quantity }}</div>
              </div>
            </div>

            <!-- Other Columns (Vertical alignment wrapper) -->
            <div class="col-total border-left flex-center-col">
              <span class="price-text">₫{{ formatNumberWithDots(order.totalAmount) }}</span>
              <span class="sub-text">{{ order.paymentMethod }}</span>
            </div>

            <div class="col-status border-left flex-center-col">
              <span class="status-strong">{{
                order.status === 'TO_PICKUP' ? 'To Pickup' : order.status
              }}</span>
            </div>

            <div class="col-countdown border-left flex-center-col">
              <span class="countdown-text">{{ order.countdown }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.shipping-manager {
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
  border-bottom: none; /* Merges visually with first card if close, but we have margin */
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
</style>
