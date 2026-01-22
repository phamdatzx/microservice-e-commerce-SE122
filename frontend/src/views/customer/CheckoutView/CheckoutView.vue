<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'
import {
  Location,
  Money,
  CreditCard,
  ArrowRight,
  Plus,
  Ticket,
  Van,
  WarnTriangleFilled,
} from '@element-plus/icons-vue'
import { ElMessage, ElLoading } from 'element-plus'
import { formatNumberWithDots } from '@/utils/formatNumberWithDots'

const router = useRouter()
const route = useRoute()

const checkoutItems = ref<any[]>([])
const addresses = ref<any[]>([])
const selectedAddressId = ref('')
const paymentMethod = ref('COD')
const voucherCode = ref('') // Keep for display if needed
const isLoading = ref(false)

// Shipping Service State
const shippingServices = ref<any[]>([])
const selectedService = ref<any>(null)
const sellerAddress = ref<any>(null)
const isCalculatingShipping = ref(false)

// Voucher State
const vouchers = ref<any[]>([])
const selectedVoucherId = ref('')
const showVoucherDialog = ref(false)

const subtotal = computed(() => {
  return checkoutItems.value.reduce((sum, item) => sum + item.price * (item.quantity || 1), 0)
})

const shippingFee = ref<number | null>(0)
const expectedDeliveryTime = ref<number | null>(null)

const selectedVoucher = computed(() => {
  return vouchers.value.find((v) => v.id === selectedVoucherId.value)
})

const discountAmount = computed(() => {
  if (!selectedVoucher.value) return 0
  const v = selectedVoucher.value.voucher

  // Check min spend
  if (subtotal.value < v.min_order_value) return 0

  let discount = 0
  if (v.discount_type === 'FIXED') {
    discount = v.discount_value
  } else {
    discount = subtotal.value * (v.discount_value / 100)
    if (v.max_discount_value && discount > v.max_discount_value) {
      discount = v.max_discount_value
    }
  }
  return Math.round(discount)
})

const total = computed(() => {
  const t = subtotal.value + (shippingFee.value || 0) - discountAmount.value
  return t > 0 ? t : 0
})

const selectedAddress = computed(() => {
  return addresses.value.find((a) => a.id === selectedAddressId.value)
})

const fetchVouchers = async () => {
  try {
    const response = await axios.get(`${import.meta.env.VITE_BE_API_URL}/product/saved-vouchers`, {
      headers: {
        Authorization: `Bearer ${localStorage.getItem('access_token')}`,
      },
    })
    if (response.data) {
      vouchers.value = response.data
    }
  } catch (error) {
    console.error('Failed to fetch vouchers:', error)
  }
}

const fetchAddresses = async () => {
  isLoading.value = true
  try {
    // Parallel fetch if possible, but sequential is fine
    await fetchVouchers()

    const response = await axios.get(`${import.meta.env.VITE_BE_API_URL}/user/addresses`, {
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

      // Auto fetch services once address is ready and seller address is ready
      if (selectedAddressId.value && sellerAddress.value) {
        fetchShippingServices()
      }
    }
  } catch (error) {
    console.error('Failed to fetch addresses:', error)
  } finally {
    isLoading.value = false
  }
}

const fetchShippingServices = async () => {
  // Require seller address
  if (!sellerAddress.value || !selectedAddress.value) {
    return
  }

  // Reset current selections
  shippingServices.value = []
  selectedService.value = null
  shippingFee.value = null
  expectedDeliveryTime.value = null

  isCalculatingShipping.value = true
  try {
    const fromDistrict = parseInt(sellerAddress.value.district_id)

    const response = await axios.post(
      'https://dev-online-gateway.ghn.vn/shiip/public-api/v2/shipping-order/available-services',
      {
        shop_id: 885,
        from_district: fromDistrict,
        to_district: parseInt(selectedAddress.value?.district_id),
      },
      {
        headers: {
          Token: import.meta.env.VITE_GHN_TOKEN,
        },
      },
    )

    if (response.data.code === 200) {
      shippingServices.value = response.data.data

      // Auto select first service and calculate details
      if (shippingServices.value.length > 0) {
        selectedService.value = shippingServices.value[0]
        // Trigger calculations sequentially
        await calculateShippingFee()
        await calculateLeadTime()
      }
    } else {
      ElMessage.error(response.data.message || 'Failed to fetch services')
    }
  } catch (error) {
    console.error('Failed to fetch services:', error)
    ElMessage.error('Failed to load shipping services')
  } finally {
    isCalculatingShipping.value = false
  }
}

const calculateShippingFee = async () => {
  if (!selectedService.value || !selectedAddress.value || !sellerAddress.value) {
    shippingFee.value = 0
    return
  }

  const toDistrictId = parseInt(selectedAddress.value.district_id)
  const toWardCode = selectedAddress.value.ward_code

  const fromDistrictId = parseInt(sellerAddress.value.district_id)
  const fromWardCode = sellerAddress.value.ward_code

  try {
    const response = await axios.post(
      'https://dev-online-gateway.ghn.vn/shiip/public-api/v2/shipping-order/fee',
      {
        shop_id: 885,
        from_district_id: fromDistrictId,
        from_ward_code: fromWardCode,
        service_id: selectedService.value.service_id,
        service_type_id: selectedService.value.service_type_id,
        to_district_id: toDistrictId,
        to_ward_code: toWardCode,
        height: 50,
        length: 20,
        weight: 20,
        width: 20,
      },
      {
        headers: {
          Token: import.meta.env.VITE_GHN_TOKEN,
        },
      },
    )

    if (response.data.code === 200) {
      shippingFee.value = response.data.data.total
    } else {
      console.error('Failed to calculate shipping fee:', response.data.message)
      // Optional: ElMessage.warning('Could not calculate shipping fee')
    }
  } catch (error) {
    console.error('Error calculating shipping fee:', error)
  }
}

const calculateLeadTime = async () => {
  if (!selectedService.value || !selectedAddress.value || !sellerAddress.value) {
    expectedDeliveryTime.value = null
    return
  }

  const toDistrictId = parseInt(selectedAddress.value.district_id)
  const toWardCode = selectedAddress.value.ward_code

  const fromDistrictId = parseInt(sellerAddress.value.district_id)
  const fromWardCode = sellerAddress.value.ward_code

  const body = {
    shop_id: 885,
    from_district_id: fromDistrictId,
    from_ward_code: fromWardCode,
    to_district_id: toDistrictId,
    to_ward_code: toWardCode,
    service_id: selectedService.value.service_id,
  }

  console.log('Calculating Lead Time Body:', JSON.stringify(body, null, 2))

  try {
    const response = await axios.post(
      'https://dev-online-gateway.ghn.vn/shiip/public-api/v2/shipping-order/leadtime',
      body,
      {
        headers: {
          Token: import.meta.env.VITE_GHN_TOKEN,
        },
      },
    )

    if (response.data.code === 200) {
      expectedDeliveryTime.value = response.data.data.leadtime
    }
  } catch (error) {
    console.error('Error calculating lead time:', error)
  }
}

// Watch for changes to trigger calculation
import { watch } from 'vue'

// When address changes, force refresh available services
watch(selectedAddress, (newVal) => {
  if (newVal && sellerAddress.value) {
    fetchShippingServices()
  }
})

const formatDate = (date: number | string | Date) => {
  if (!date) return ''
  const d = new Date(date)
  if (isNaN(d.getTime())) return ''
  return d.toLocaleDateString('en-US', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
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
    ward_code: selectedAddress.value.ward_code || '',
    province_id: selectedAddress.value.province_id ? String(selectedAddress.value.province_id) : '',
    district_id: selectedAddress.value.district_id ? String(selectedAddress.value.district_id) : '',
  }

  const isInstant = route.query.mode === 'instant'

  let payload: any = {}

  if (isInstant) {
    // Instant Checkout Payload Structure
    payload = {
      items: checkoutItems.value.map((item) => ({
        seller_id: item.sellerId,
        product_id: item.productId,
        variant_id: item.variantId,
        quantity: item.quantity,
      })),
      shipping_address: addressPayload,
      payment_method: paymentMethod.value,
      // voucher_id and delivery_service_id added below if exist
    }
  } else {
    // Standard Cart Checkout Payload
    payload = {
      cart_item_ids: checkoutItems.value.map((item) => item.id), // Assuming item.id is the cart_item_id
      shipping_address: addressPayload,
      payment_method: paymentMethod.value,
    }
  }

  if (selectedVoucher.value?.voucher?.id) {
    payload.voucher_id = selectedVoucher.value.voucher.id
  }

  if (selectedService.value) {
    payload.delivery_service_id = selectedService.value.service_id
  }

  console.log('Checkout Payload:', JSON.stringify(payload, null, 2))

  const loading = ElLoading.service({
    lock: true,
    text: 'Placing Order...',
    background: 'rgba(0, 0, 0, 0.7)',
  })

  try {
    // 1. Checkout API
    const endpoint = isInstant ? '/order/instant-checkout' : '/order/checkout'
    const response = await axios.post(`${import.meta.env.VITE_BE_API_URL}${endpoint}`, payload, {
      headers: {
        Authorization: `Bearer ${localStorage.getItem('access_token')}`,
      },
    })

    if (response.data) {
      const orderId = response.data.order_id

      if (paymentMethod.value === 'STRIPE') {
        // 2. Payment Session API
        await handleStripePayment(orderId)
      } else {
        ElMessage.success('Order placed successfully!')
        if (isInstant) {
          localStorage.removeItem('instant_checkout_item')
        } else {
          localStorage.removeItem('checkout_items')
          localStorage.removeItem('checkout_voucher_id')
        }
        router.push('/') // Go to home page
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
      `${import.meta.env.VITE_BE_API_URL}/order/${orderId}/payment`,
      {},
      {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('access_token')}`,
        },
      },
    )
    if (response.data && response.data.payment_url) {
      localStorage.removeItem('checkout_items')
      localStorage.removeItem('checkout_voucher_id')
      window.location.href = response.data.payment_url
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

const selectVoucher = (voucherId: string) => {
  if (selectedVoucherId.value === voucherId) {
    selectedVoucherId.value = '' // Deselect
  } else {
    selectedVoucherId.value = voucherId
  }
  showVoucherDialog.value = false
}

// Filter usable vouchers
const validVouchers = computed(() => {
  const now = new Date()
  return vouchers.value.filter((v) => {
    const isExpired = new Date(v.voucher.end_time) <= now
    const isFullyUsed = v.used_count >= v.max_uses_allowed
    // Also check if applicable to current cart value? Optional but good UX
    // const isApplicable = subtotal.value >= v.voucher.min_order_value
    return !isExpired && !isFullyUsed
  })
})

const fetchSellerAddress = async () => {
  try {
    const isInstant = route.query.mode === 'instant'

    if (isInstant) {
      // For instant checkout, we have the sellerId in the item
      const item = checkoutItems.value[0]
      if (!item || !item.sellerId) return

      const response = await axios.get(
        `${import.meta.env.VITE_BE_API_URL}/user/public/seller/${item.sellerId}`,
      )
      if (response.data && response.data.data && response.data.data.address) {
        sellerAddress.value = response.data.data.address
      }
      return
    }

    // Default flow: fetch from cart to get seller address
    const response = await axios.get(`${import.meta.env.VITE_BE_API_URL}/order/cart`, {
      headers: {
        Authorization: `Bearer ${localStorage.getItem('access_token')}`,
      },
    })

    if (response.data && response.data.cart_items && response.data.cart_items.length > 0) {
      // Find the first item that matches one of our checkout items (to ensure we get the right seller)
      // Assuming all checkout items belong to the same seller for now as per standard flow, or we just take the first.
      const firstItem = response.data.cart_items.find((item: any) =>
        checkoutItems.value.some((cItem) => cItem.id === item.id),
      )

      if (firstItem && firstItem.seller && firstItem.seller.address) {
        sellerAddress.value = firstItem.seller.address
      }
    }
  } catch (error) {
    console.error('Failed to fetch seller address:', error)
  }
}

onMounted(async () => {
  // Load checkout items
  const isInstant = route.query.mode === 'instant'
  const storedItems = localStorage.getItem(isInstant ? 'instant_checkout_item' : 'checkout_items')

  if (storedItems) {
    checkoutItems.value = JSON.parse(storedItems)
  }

  if (checkoutItems.value.length === 0) {
    ElMessage.info('No items to checkout')
    if (isInstant) {
      router.push('/')
    } else {
      router.push('/cart')
    }
    return
  }

  // Load persistence voucher from cart
  const savedVoucherId = localStorage.getItem('checkout_voucher_id')
  if (savedVoucherId) {
    selectedVoucherId.value = savedVoucherId
  }

  // Fetch seller address first
  await fetchSellerAddress()

  fetchAddresses()
})
</script>

<template>
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
                  <div class="item-name">{{ item.productName }}</div>
                  <div class="item-variant" v-if="item.productOption">
                    Variant: {{ item.productOption }}
                  </div>
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
          <div class="section-card summary-card" v-loading="isCalculatingShipping">
            <h3 style="margin-bottom: 10px">Order Summary</h3>
            <div class="summary-row">
              <span>Merchandise Subtotal</span>
              <span>{{ formatNumberWithDots(subtotal) }}đ</span>
            </div>

            <div class="summary-row">
              <span>Shipping Fee</span>
              <span>{{
                shippingFee !== null ? formatNumberWithDots(shippingFee) + 'đ' : '---'
              }}</span>
            </div>

            <!-- Voucher Section -->
            <div class="summary-row voucher-row" @click="showVoucherDialog = true">
              <div class="voucher-label">
                <el-icon><Ticket /></el-icon> <span>Voucher</span>
              </div>
              <div class="voucher-value" :class="{ 'has-discount': discountAmount > 0 }">
                <span v-if="discountAmount > 0">-{{ formatNumberWithDots(discountAmount) }}đ</span>
                <span v-else class="select-voucher-text">Select Voucher</span>
              </div>
            </div>

            <el-divider style="margin: 12px 0" />
            <div class="summary-row total-row">
              <span>Total Payment</span>
              <span class="total-amount">{{ formatNumberWithDots(total) }}đ</span>
            </div>

            <!-- Expected Delivery -->
            <div class="summary-row delivery-time-row">
              <div class="delivery-label">
                <span>Expected Delivery</span>
              </div>
              <div class="delivery-value">
                {{ expectedDeliveryTime ? formatDate(expectedDeliveryTime * 1000) : '---' }}
              </div>
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

    <!-- Voucher Dialog -->
    <el-dialog
      v-model="showVoucherDialog"
      title="Select Voucher"
      width="500px"
      class="voucher-dialog"
    >
      <div v-if="validVouchers.length === 0" class="no-vouchers">
        <el-empty description="No usable vouchers found" />
      </div>
      <div class="voucher-list-dialog" v-else>
        <div
          v-for="v in validVouchers"
          :key="v.id"
          class="voucher-item-dialog"
          :class="{
            selected: selectedVoucherId === v.id,
            disabled: subtotal < v.voucher.min_order_value,
          }"
          @click="subtotal >= v.voucher.min_order_value && selectVoucher(v.id)"
        >
          <div class="v-left">
            <div class="v-code">{{ v.voucher.code }}</div>
            <div class="v-name">{{ v.voucher.name }}</div>
            <div class="v-desc">
              Min. Spend {{ formatNumberWithDots(v.voucher.min_order_value) }}đ
            </div>
          </div>
          <div class="v-right">
            <div class="v-discount">
              <span v-if="v.voucher.discount_type === 'FIXED'"
                >-{{ formatNumberWithDots(v.voucher.discount_value) }}đ</span
              >
              <span v-else>-{{ v.voucher.discount_value }}%</span>
            </div>
            <el-radio
              :model-value="selectedVoucherId"
              :label="v.id"
              @change="selectVoucher(v.id)"
              :disabled="subtotal < v.voucher.min_order_value"
              >{{ '' }}</el-radio
            >
          </div>
        </div>
      </div>
    </el-dialog>

    <!-- Shipping Service Dialog REMOVED -->
  </div>
</template>

<style scoped>
.voucher-row {
  cursor: pointer;
  padding: 12px 12px;
  border-top: 1px dashed #eee;
  border-bottom: 1px dashed #eee;
  margin: 10px 0;
  align-items: center;
  transition: all 0.2s;
}

.voucher-row:hover {
  background-color: #fafafa;
}

.voucher-label {
  display: flex;
  align-items: center;
  gap: 5px;
  color: var(--main-color);
  font-weight: 500;
}

.voucher-value {
  color: #999;
  font-size: 14px;
}

.voucher-value.has-discount {
  color: #ef4444;
  font-weight: bold;
}

.select-voucher-text {
  color: #22c55e;
}

.delivery-time-row {
  border-radius: 8px;
  margin-top: 4px !important;
  margin-bottom: 0 !important;
  font-size: 13px;
}

.delivery-label {
  color: #666;
}

.delivery-value {
  color: var(--main-color);
  font-weight: 600;
}

/* Dialog Styles */
.voucher-list-dialog {
  display: flex;
  flex-direction: column;
  gap: 10px;
  max-height: 400px;
  overflow-y: auto;
}

.voucher-item-dialog {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border: 1px solid #eee;
  padding: 10px 15px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.voucher-item-dialog:hover {
  border-color: var(--main-color);
  background: #f0fdf4;
}

.voucher-item-dialog.selected {
  border-color: var(--main-color);
  background: #f0fdf4;
}

.voucher-item-dialog.disabled {
  opacity: 0.5;
  cursor: not-allowed;
  background: #f5f5f5;
}

.v-left {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.v-code {
  font-weight: bold;
  color: #333;
}

.v-name {
  font-size: 13px;
  color: #666;
}

.v-desc {
  font-size: 12px;
  color: #999;
}

.v-right {
  display: flex;
  align-items: center;
  gap: 10px;
}

.v-discount {
  font-weight: bold;
  color: var(--main-color);
  font-size: 16px;
}

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
  margin-bottom: 0 !important;
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
