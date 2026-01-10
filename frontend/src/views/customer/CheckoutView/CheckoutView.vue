<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import Header from '@/components/Header.vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { Location, Money, CreditCard, ArrowRight, Plus, Loading } from '@element-plus/icons-vue'
import { ElMessage, ElLoading } from 'element-plus'
import { formatNumberWithDots } from '@/utils/formatNumberWithDots'

const router = useRouter()

// --- State ---
const checkoutItems = ref<any[]>([])
const addresses = ref<any[]>([])
const selectedAddressId = ref('')
const paymentMethod = ref('COD')
const voucherCode = ref('') // logic for voucher if needed later
const isLoading = ref(false)

// --- Computed ---
const subtotal = computed(() => {
  return checkoutItems.value.reduce((sum, item) => sum + item.price * (item.quantity || 1), 0)
})

const shippingFee = computed(() => {
  return 30000 // Fixed for now
})

const total = computed(() => {
  return subtotal.value + shippingFee.value
})

const selectedAddress = computed(() => {
  return addresses.value.find((a) => a.id === selectedAddressId.value)
})

// --- Methods ---

const fetchAddresses = async () => {
  isLoading.value = true
  try {
    const response = await axios.get('http://localhost:81/api/user/addresses', {
      headers: {
        Authorization: `Bearer ${localStorage.getItem('access_token')}`,
      },
    })
    if (response.data.status === 200) {
      addresses.value = response.data.data
      const defaultAddr = addresses.value.find((a) => a.default)
      if (defaultAddr) {
        selectedAddressId.value = defaultAddr.id
      } else if (addresses.value.length > 0) {
        selectedAddressId.value = addresses.value[0].id
      }
    }
  } catch (error) {
    console.error('Failed to fetch addresses:', error)
  } finally {
    isLoading.value = false
  }
}

const handlePlaceOrder = async () => {
  if (!selectedAddressId.value) {
    ElMessage.warning('Please select a shipping address')
    return
  }

  // Prepare payload
  const addressPayload = {
    full_name: selectedAddress.value.full_name,
    phone: selectedAddress.value.phone,
    address_line: selectedAddress.value.address_line,
    ward: selectedAddress.value.ward,
    district: selectedAddress.value.district || 'quan',
    province: selectedAddress.value.province,
    country: selectedAddress.value.country || 'viet nam',
    latitude: selectedAddress.value.latitude || 0,
    longitude: selectedAddress.value.longitude || 0,
  }

  const payload = {
    cart_item_ids: checkoutItems.value.map((item) => item.id), // Assuming item.id is the cart_item_id
    shipping_address: addressPayload,
    payment_method: paymentMethod.value,
  }

  const loading = ElLoading.service({
    lock: true,
    text: 'Placing Order...',
    background: 'rgba(0, 0, 0, 0.7)',
  })

  try {
    // 1. Checkout API
    const response = await axios.post('http://localhost:81/api/order/checkout', payload, {
      headers: {
        Authorization: `Bearer ${localStorage.getItem('access_token')}`,
      },
    })

    if (response.data) {
      const orderId = response.data.id // Adjust based on actual response structure

      if (paymentMethod.value === 'STRIPE') {
        // 2. Payment Session API
        await handleStripePayment(orderId)
      } else {
        ElMessage.success('Order placed successfully!')
        router.push('/profile') // Go to order history
      }
    }
  } catch (error: any) {
    console.error('Checkout failed:', error)
    ElMessage.error(error.response?.data?.message || 'Failed to place order')
  } finally {
    loading.close()
  }
}

const handleStripePayment = async (orderId: string) => {
  try {
    const response = await axios.post(
      `http://localhost:81/api/order/${orderId}/payment`,
      {},
      {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('access_token')}`,
        },
      },
    )
    if (response.data && response.data.url) {
      window.location.href = response.data.url
    } else {
      ElMessage.warning('Payment session created but no URL provided')
      router.push('/profile')
    }
  } catch (error) {
    console.error('Payment session failed:', error)
    ElMessage.error('Failed to initiate payment')
    router.push('/profile') // Redirect anyway so user can try paying again from history
  }
}

const navigateToAddress = () => {
  router.push('/profile') // Shortcut to manage addresses
}

onMounted(() => {
  // Load checkout items
  const storedItems = localStorage.getItem('checkout_items')
  if (storedItems) {
    checkoutItems.value = JSON.parse(storedItems)
  }

  if (checkoutItems.value.length === 0) {
    ElMessage.info('No items to checkout')
    router.push('/cart')
    return
  }

  fetchAddresses()
})
</script>

<template>
  <Header />
  <div class="checkout-page">
    <div class="container">
      <h2 class="page-title">
        <el-icon><Money /></el-icon> Checkout
      </h2>

      <div class="checkout-grid">
        <!-- Left Column: Address & Payment -->
        <div class="left-col">
          <!-- Address Section -->
          <div class="section-card address-section" v-loading="isLoading">
            <div class="section-header">
              <h3>
                <el-icon><Location /></el-icon> Delivery Address
              </h3>
              <el-button link type="primary" @click="navigateToAddress">Manage Addresses</el-button>
            </div>

            <div v-if="addresses.length === 0" class="no-address">
              <p>You have no saved addresses.</p>
              <el-button type="primary" @click="navigateToAddress">Add Address</el-button>
            </div>

            <div class="address-list">
              <el-radio-group v-model="selectedAddressId" class="addr-group">
                <el-radio
                  v-for="addr in addresses"
                  :key="addr.id"
                  :label="addr.id"
                  border
                  class="addr-radio"
                >
                  <div class="addr-content">
                    <span class="addr-name">{{ addr.full_name }} ({{ addr.phone }})</span>
                    <span class="addr-detail"
                      >{{ addr.address_line }}, {{ addr.ward }}, {{ addr.province }}</span
                    >
                  </div>
                </el-radio>
              </el-radio-group>
            </div>
          </div>

          <!-- Products Section -->
          <div class="section-card product-section">
            <h3>Products Ordered</h3>
            <div class="product-list">
              <div v-for="item in checkoutItems" :key="item.id" class="product-item">
                <img :src="item.imageUrl" alt="Product" class="item-img" />
                <div class="item-info">
                  <div class="item-name">{{ item.name }}</div>
                  <div class="item-variant">Variant: {{ item.productOption }}</div>
                </div>
                <div class="item-meta">
                  <div class="item-price">{{ formatNumberWithDots(item.price) }}đ</div>
                  <div class="item-qty">x{{ item.quantity }}</div>
                  <div class="item-total">
                    {{ formatNumberWithDots(item.price * item.quantity) }}đ
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Payment Method -->
          <div class="section-card payment-section">
            <h3>Payment Method</h3>
            <el-radio-group v-model="paymentMethod">
              <el-radio label="COD" border>Cash on Delivery</el-radio>
              <el-radio label="STRIPE" border>Stripe / Credit Card</el-radio>
            </el-radio-group>
          </div>
        </div>

        <!-- Right Column: Summary -->
        <div class="right-col">
          <div class="section-card summary-card">
            <h3>Order Summary</h3>
            <div class="summary-row">
              <span>Merchandise Subtotal</span>
              <span>{{ formatNumberWithDots(subtotal) }}đ</span>
            </div>
            <div class="summary-row">
              <span>Shipping Total</span>
              <span>{{ formatNumberWithDots(shippingFee) }}đ</span>
            </div>
            <el-divider />
            <div class="summary-row total-row">
              <span>Total Payment</span>
              <span class="total-amount">{{ formatNumberWithDots(total) }}đ</span>
            </div>
            <el-button
              type="primary"
              class="place-order-btn"
              size="large"
              :disabled="isLoading || !selectedAddressId || addresses.length === 0"
              @click="handlePlaceOrder"
            >
              Place Order
            </el-button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.checkout-page {
  background-color: #f5f5f5;
  min-height: 100vh;
  padding-bottom: 40px;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.page-title {
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  gap: 10px;
  color: var(--main-color);
}

.checkout-grid {
  display: grid;
  grid-template-columns: 1fr 350px;
  gap: 20px;
}

.section-card {
  background: white;
  border-radius: 4px;
  padding: 20px;
  margin-bottom: 20px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.section-header h3 {
  margin: 0;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 18px;
}

/* Address List */
.addr-group {
  display: flex;
  flex-direction: column;
  gap: 10px;
  width: 100%;
}

.addr-radio {
  width: 100%;
  height: auto;
  padding: 15px;
  margin-right: 0;
}

:deep(.addr-radio .el-radio__label) {
  width: 100%;
  white-space: normal;
}

.addr-content {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.addr-name {
  font-weight: bold;
}

.addr-detail {
  color: #666;
  font-size: 14px;
}

/* Product List */
.product-item {
  display: flex;
  gap: 15px;
  margin-bottom: 15px;
  padding-bottom: 15px;
  border-bottom: 1px solid #f0f0f0;
}

.product-item:last-child {
  border-bottom: none;
  margin-bottom: 0;
  padding-bottom: 0;
}

.item-img {
  width: 60px;
  height: 60px;
  object-fit: cover;
  border: 1px solid #eee;
}

.item-info {
  flex: 1;
}

.item-name {
  font-weight: 500;
  margin-bottom: 4px;
}

.item-variant {
  font-size: 12px;
  color: #888;
}

.item-meta {
  text-align: right;
}

.item-total {
  font-weight: bold;
  color: var(--main-color);
}

/* Summary */
.summary-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 10px;
  color: #666;
}

.total-row {
  color: #333;
  margin-top: 15px;
  align-items: center;
}

.total-row span:first-child {
  font-size: 16px;
}

.total-amount {
  font-size: 24px;
  color: var(--main-color);
  font-weight: bold;
}

.place-order-btn {
  width: 100%;
  margin-top: 20px;
  background-color: var(--main-color);
  border-color: var(--main-color);
}
</style>
