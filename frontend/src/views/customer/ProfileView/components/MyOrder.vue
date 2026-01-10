<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { Search, ChatDotRound, Shop } from '@element-plus/icons-vue'
import { formatNumberWithDots } from '@/utils/formatNumberWithDots'
import axios from 'axios'

const activeOrderTab = ref('all')

const orders = ref<any[]>([])

const statusMapping: Record<string, string> = {
  all: '',
  'to-pay': 'TO_STATUS_PAY', // Backend doesn't support 'TO_PAY' filter directly if payment status is separate?
  // Wait, UserService CreateOrder sets status=TO_PAY for Stripe.
  // Let's assume backend filters by status field directly.
  'to-ship': 'TO_CONFIRM',
  'to-receive': 'SHIPPING', // or TO_PICKUP? User example used TO_CONFIRM. Backend update logic has TO_PICKUP, SHIPPING.
  // Let's use standard Shopee-like flow: To Ship -> TO_CONFIRM/TO_PICKUP, To Receive -> SHIPPING
  // For this project, listing "To Ship" usually maps to TO_CONFIRM.
  completed: 'COMPLETED',
  cancelled: 'CANCELLED',
  'return-refund': 'RETURNED',
}

const fetchOrders = async () => {
  try {
    const params: any = {
      page: 1,
      limit: 10,
      sort_by: 'created_at',
      sort_order: 'descending',
    }

    if (activeOrderTab.value !== 'all') {
      // Map frontend tab to backend status
      let status = activeOrderTab.value.toUpperCase().replace(/-/g, '_')

      // Override with specific mapping if needed
      if (activeOrderTab.value === 'to-ship') status = 'TO_CONFIRM'
      if (activeOrderTab.value === 'to-receive') status = 'SHIPPING' // or TO_PICKUP
      if (activeOrderTab.value === 'return-refund') status = 'RETURNED'
      if (activeOrderTab.value === 'to-pay') status = 'TO_PAY'

      params.status = status
    }

    const response = await axios.get('http://localhost:81/api/order', {
      params,
      headers: {
        Authorization: `Bearer ${localStorage.getItem('access_token')}`,
      },
    })

    if (response.data && response.data.orders) {
      orders.value = response.data.orders.map((order: any) => ({
        ...order,
        // Since API doesn't return items/shop info, we mock or handle empty for now
        // UI expects: shopName, shippingUpdate, statusText, items[], totalAmount
        shopName: 'Official Store', // Placeholder
        shippingUpdate: 'Updated recently',
        statusText: order.status,
        totalAmount: order.total,
        items: [], // API doesn't return items yet
      }))
    } else {
      orders.value = []
    }
  } catch (error) {
    console.error('Failed to fetch orders:', error)
  }
}

watch(activeOrderTab, () => {
  fetchOrders()
})

onMounted(() => {
  fetchOrders()
})

// orders are now fetched filtered
const filteredOrders = computed(() => orders.value)

const orderTabs = [
  { label: 'All', value: 'all' },
  { label: 'To Pay', value: 'to-pay' },
  { label: 'To Ship', value: 'to-ship' },
  { label: 'To Receive', value: 'to-receive' },
  { label: 'Completed', value: 'completed' },
  { label: 'Cancelled', value: 'cancelled' },
  { label: 'Return/Refund', value: 'return-refund' },
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

    <div class="order-list">
      <div v-for="order in filteredOrders" :key="order.id" class="order-card box-shadow">
        <div class="order-header">
          <div class="shop-info">
            <span class="shop-name">{{ order.shopName }}</span>
            <el-button :icon="ChatDotRound" size="small" link>Chat</el-button>
            <el-divider direction="vertical" />
            <el-button :icon="Shop" size="small" link>View Shop</el-button>
          </div>
          <div class="order-status">
            <span class="shipping-status">{{ order.shippingUpdate }}</span>
            <el-divider direction="vertical" />
            <span class="status-text">{{ order.statusText }}</span>
          </div>
        </div>

        <div class="order-items">
          <div v-for="item in order.items" :key="item.id" class="order-item">
            <img :src="item.image" :alt="item.name" class="item-image" />
            <div class="item-details">
              <h4 class="item-name">{{ item.name }}</h4>
              <p class="item-variant">Variant: {{ item.variant }}</p>
              <p class="item-quantity">x{{ item.quantity }}</p>
            </div>
            <div class="item-price">
              <span v-if="item.oldPrice" class="old-price"
                >{{ formatNumberWithDots(item.oldPrice) }}đ</span
              >
              <span class="current-price">{{ formatNumberWithDots(item.price) }}đ</span>
            </div>
          </div>
        </div>

        <div class="order-footer">
          <div class="total-section">
            <span class="total-label">Order Total:</span>
            <span class="total-amount">{{ formatNumberWithDots(order.totalAmount) }}đ</span>
          </div>
          <div class="order-actions">
            <template v-if="order.status === 'TO_RECEIVE'">
              <el-button type="primary" size="large">Order Received</el-button>
              <el-button size="large">Return/Refund</el-button>
              <el-button size="large">Contact Seller</el-button>
            </template>
            <template v-else-if="order.status === 'COMPLETED'">
              <el-button type="primary" size="large">Buy Again</el-button>
              <el-button size="large">Rate</el-button>
              <el-button size="large">Contact Seller</el-button>
            </template>
          </div>
        </div>
      </div>

      <div v-if="filteredOrders.length === 0" class="no-orders text-center">
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
  text-transform: uppercase;
  font-size: 13px;
}

.order-item {
  display: flex;
  padding: 12px 0;
  border-bottom: 1px solid #f5f5f5;
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
</style>
