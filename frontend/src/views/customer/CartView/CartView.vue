<script setup lang="ts">
import type { CheckboxValueType } from 'element-plus'
import { ref, useTemplateRef, watch, onMounted, computed, inject } from 'vue'
import CartSellerProductWrapper from './components/CartSellerProductWrapper.vue'
import { Ticket, Loading, ShoppingCart } from '@element-plus/icons-vue'
import { ElNotification, ElMessageBox, ElMessage } from 'element-plus'

import { Splide, SplideSlide } from '@splidejs/vue-splide'
import axios from 'axios'
import { formatNumberWithDots } from '@/utils/formatNumberWithDots'
import RecentlyViewed from '@/components/RecentlyViewed.vue'
import { useRouter } from 'vue-router'

const router = useRouter()

export interface VariantAPI {
  id: string
  sku: string
  options: Record<string, string>
  price: number
  stock: number
  image: string
}

export interface CartItemAPI {
  id: string
  user_id: string
  seller: {
    id: string
    name: string
    username: string
  }
  product: {
    id: string
    name: string
    seller_id: string
    seller_category_ids: string[]
  }
  variant_id: string
  quantity: number
  variant: VariantAPI
  created_at: string
  updated_at: string
}

export interface Product {
  id: string
  productId: string
  imageUrl: string
  productName: string
  productOption: string
  price: number
  quantity: number
  sellerName: string
  sellerId: string
  variantId: string
  stock: number
  isSoldOut?: boolean
}

export interface CartSeller {
  sellerName: string
  sellerId: string
  products: Product[]
}

const fetchCartCount = inject('fetchCartCount') as () => void

const checkAll = ref(false)
const checkedProducts = ref<Product[]>([])
const CartSellerProductWrapperRefs = ref<InstanceType<typeof CartSellerProductWrapper>[]>([])
const cartData = ref<CartSeller[]>([])
const isLoading = ref(false)

// Voucher State
const vouchers = ref<any[]>([])
const selectedVoucherId = ref('')
const showVoucherDialog = ref(false)

const subtotal = computed(() => {
  return checkedProducts.value.reduce((sum, item) => sum + item.price * (item.quantity || 1), 0)
})

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

const totalPayment = computed(() => {
  const t = subtotal.value - discountAmount.value
  return t > 0 ? t : 0
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

const selectVoucher = (voucherId: string) => {
  if (selectedVoucherId.value === voucherId) {
    selectedVoucherId.value = '' // Deselect
  } else {
    selectedVoucherId.value = voucherId
  }
  showVoucherDialog.value = false
}

const validVouchers = computed(() => {
  const now = new Date()
  return vouchers.value.filter((v) => {
    const isExpired = new Date(v.voucher.end_time) <= now
    const isFullyUsed = v.used_count >= v.max_uses_allowed
    return !isExpired && !isFullyUsed
  })
})

const fetchCart = async () => {
  const token = localStorage.getItem('access_token')
  if (!token) return

  isLoading.value = true
  try {
    const response = await axios.get(`${import.meta.env.VITE_BE_API_URL}/order/cart`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })

    const items: CartItemAPI[] = response.data.cart_items

    // Group by seller.id
    const grouped = items.reduce((acc: Record<string, CartSeller>, item) => {
      const sellerId = item.seller.id
      if (!acc[sellerId]) {
        acc[sellerId] = {
          sellerId: sellerId,
          sellerName: item.seller.name,
          products: [],
        }
      }

      acc[sellerId].products.push({
        id: item.id,
        productId: item.product.id,
        imageUrl: item.variant.image || 'https://placehold.co/100x100?text=No+Image',
        productName: item.product.name,
        productOption: Object.values(item.variant.options).join(', '),
        price: item.variant.price,
        quantity: item.quantity,
        stock: item.variant.stock,
        sellerName: item.seller.name,
        sellerId: sellerId,
        variantId: item.variant_id,
      })

      return acc
    }, {})

    cartData.value = Object.values(grouped)
  } catch (error) {
    console.error('Error fetching cart:', error)
  } finally {
    isLoading.value = false
  }
}

const allProducts = computed(() => cartData.value.map((seller) => seller.products).flat())

const handleCheckout = () => {
  if (checkedProducts.value.length === 0) {
    ElMessage.warning('Please select products to checkout')
    return
  }
  localStorage.setItem('checkout_items', JSON.stringify(checkedProducts.value))
  localStorage.setItem('checkout_voucher_id', selectedVoucherId.value)
  router.push('/checkout')
}

onMounted(() => {
  fetchCart()
  fetchVouchers()
})

const handleCheckAll = () => {
  CartSellerProductWrapperRefs.value.forEach((ref) => {
    ref.handleCheckAll()
  })
}

const handleUncheckAll = () => {
  CartSellerProductWrapperRefs.value.forEach((ref) => {
    ref.handleUncheckAll()
  })
}

const handleCheckAllChange = (val: CheckboxValueType) => {
  checkedProducts.value = val ? allProducts.value : []

  if (val) {
    handleCheckAll()
  } else {
    handleUncheckAll()
  }
}

const handleCheckedProductsChange = (checkedSellerProducts: Product[], sellerName: string) => {
  checkedProducts.value = checkedProducts.value.filter(
    (product) => product.sellerName !== sellerName,
  )

  if (checkedSellerProducts.length > 0) {
    checkedProducts.value.push(...checkedSellerProducts)
  }

  checkAll.value = checkedProducts.value.length === allProducts.value.length
}

const removeFromCart = async (cartItemId: string, showNotification = true) => {
  const token = localStorage.getItem('access_token')
  if (!token) return

  if (showNotification) {
    try {
      await ElMessageBox.confirm(
        'Are you sure you want to remove this item from your cart?',
        'Confirm Remove',
        {
          confirmButtonText: 'Remove',
          cancelButtonText: 'Cancel',
          type: 'warning',
        },
      )
    } catch (error) {
      return // User cancelled
    }
  }

  try {
    await axios.delete(`${import.meta.env.VITE_BE_API_URL}/order/cart/${cartItemId}`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })

    // Remove from cartData
    cartData.value = cartData.value
      .map((seller) => ({
        ...seller,
        products: seller.products.filter((p) => p.id !== cartItemId),
      }))
      .filter((seller) => seller.products.length > 0)

    // Remove from checkedProducts
    checkedProducts.value = checkedProducts.value.filter((p) => p.id !== cartItemId)

    if (showNotification) {
      ElNotification({
        title: 'Success',
        message: 'Item removed from cart',
        type: 'success',
      })
    }
    fetchCartCount?.()
  } catch (error) {
    console.error('Error removing from cart:', error)
    if (showNotification) {
      ElNotification({
        title: 'Error',
        message: 'Failed to remove item from cart',
        type: 'error',
      })
    }
  }
}

const handleFooterSelectAll = () => {
  checkAll.value = !checkAll.value
  handleCheckAllChange(checkAll.value)
}

const handleDeleteSelected = async () => {
  if (checkedProducts.value.length === 0) return

  try {
    await ElMessageBox.confirm(
      `Are you sure you want to delete ${checkedProducts.value.length} items?`,
      'Confirm Delete',
      {
        confirmButtonText: 'Delete',
        cancelButtonText: 'Cancel',
        type: 'warning',
      },
    )

    const idsToDelete = checkedProducts.value.map((p) => p.id)

    // Using Promise.all to delete each item
    // Note: If backend supports batch delete, it would be better to use that.
    await Promise.all(idsToDelete.map((id) => removeFromCart(id, false)))

    ElNotification({
      title: 'Success',
      message: `Successfully removed ${idsToDelete.length} items from cart.`,
      type: 'success',
    })
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Error in bulk delete:', error)
    }
  }
}

const updateCartQuantity = async (cartItemId: string, quantity: number) => {
  const token = localStorage.getItem('access_token')
  if (!token) return

  try {
    await axios.put(
      `${import.meta.env.VITE_BE_API_URL}/order/cart/${cartItemId}`,
      { quantity },
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      },
    )
    fetchCartCount?.()
  } catch (error) {
    console.error('Error updating cart quantity:', error)
    ElNotification({
      title: 'Error',
      message: 'Failed to update cart quantity',
      type: 'error',
    })
    // Optional: Refresh cart to revert local change if API fails
    fetchCart()
  }
}
</script>

<template>
  <div class="main-container">
    <div
      class="cart-box-shadow border-radius"
      style="
        padding: 0 20px;
        margin: 20px 0;
        display: flex;
        align-items: center;
        margin-bottom: 12px;
        height: 64px;
        background-color: #fff;
      "
    >
      <el-checkbox
        v-model="checkAll"
        @change="handleCheckAllChange"
        size="large"
        class="cartCheckbox"
      />
      <h3 style="flex: 1">Product</h3>
      <h3 style="width: 174px; text-align: center">Price</h3>
      <h3 style="width: 169px; text-align: center">Quantity</h3>
      <h3 style="width: 114px; text-align: center">Total</h3>
      <h3 style="width: 139px; text-align: center">Operation</h3>
    </div>

    <div v-if="isLoading" style="padding: 100px; text-align: center">
      <el-icon class="is-loading" size="40"><Loading /></el-icon>
    </div>

    <template v-else-if="cartData.length > 0">
      <CartSellerProductWrapper
        v-for="seller in cartData"
        :key="seller.sellerId"
        ref="CartSellerProductWrapperRefs"
        :seller="seller"
        @checked-products-change="handleCheckedProductsChange"
        @delete-product="removeFromCart"
        @update-quantity="updateCartQuantity"
      />
    </template>

    <div v-else style="padding: 100px; text-align: center; background: #fff" class="border-radius">
      <el-icon size="64" style="margin-bottom: 20px; opacity: 0.2"><ShoppingCart /></el-icon>
      <h3>Your shopping cart is empty</h3>
      <p style="color: #999; margin: 12px 0 24px">
        Adding items to your cart that you like to buy!
      </p>
      <el-button type="primary" color="var(--main-color)" size="large" @click="$router.push('/')">
        Go Shopping Now
      </el-button>
    </div>

    <div
      class="border-radius"
      style="
        background-color: #fff;
        border: 1px solid #ccc;
        padding: 16px;
        position: sticky;
        bottom: 0;
        z-index: 1000;
      "
    >
      <div style="display: flex; padding-bottom: 12px">
        <button class="btn" @click="handleFooterSelectAll">
          {{ checkAll ? 'Unselect All' : 'Select All' }} ({{ allProducts.length }})
        </button>
        <button class="btn" @click="handleDeleteSelected">Delete</button>
      </div>
      <div style="display: flex; justify-content: space-between; align-items: center">
        <el-button
          size="large"
          style="font-size: 18px; color: var(--main-color); padding: 28px 20px 24px"
          @click="showVoucherDialog = true"
        >
          <el-icon style="margin-right: 8px; font-size: 24px; position: relative; top: -2px"
            ><Ticket
          /></el-icon>
          Voucher:
          <span v-if="selectedVoucher" style="margin-left: 8px; color: #ef4444"
            >{{ selectedVoucher.voucher.code }} (-{{ formatNumberWithDots(discountAmount) }}đ)</span
          >
          <span v-else style="margin-left: 8px">Apply Voucher</span>
        </el-button>
        <div style="display: flex; align-items: center">
          <div
            style="display: flex; flex-direction: column; align-items: flex-end; margin-right: 12px"
          >
            <span
              v-if="discountAmount > 0"
              style="font-size: 14px; color: #999; text-decoration: line-through"
            >
              Subtotal: {{ formatNumberWithDots(subtotal) }}đ
            </span>
            <span style="font-size: 20px">
              Total ({{ checkedProducts.length }} products):
              <span style="font-size: 24px; color: var(--main-color)">
                {{ formatNumberWithDots(totalPayment) }}đ
              </span>
            </span>
          </div>
          <el-button
            style="
              font-size: 24px;
              padding: 32px 40px;
              color: var(--main-color);
              margin-left: 4px;
              margin-right: 4px;
            "
            @click="handleCheckout"
            >Checkout</el-button
          >
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

    <RecentlyViewed style="margin-top: 20px" />
  </div>
</template>

<style lang="css" scoped>
.btn {
  background-color: #fff;
  border: none;
  padding: 12px 20px;
  font-size: 18px;
  border: 1px solid #ccc;
  margin-right: 8px;
  transition:
    color 0.2s,
    border-color 0.2s;

  &:hover {
    color: #22c55e;
    border-color: currentColor;
  }
}

.voucher-list-dialog {
  display: flex;
  flex-direction: column;
  gap: 12px;
  max-height: 400px;
  overflow-y: auto;
  padding: 4px;
}

.voucher-item-dialog {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border: 1px solid #e4e4e7;
  padding: 16px;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.voucher-item-dialog:hover:not(.disabled) {
  border-color: #22c55e;
  background-color: #f0fdf4;
}

.voucher-item-dialog.selected {
  border-color: #22c55e;
  background-color: #f0fdf4;
  box-shadow: 0 0 0 1px #22c55e;
}

.voucher-item-dialog.disabled {
  opacity: 0.5;
  cursor: not-allowed;
  background-color: #f9fafb;
}

.v-left {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.v-code {
  font-weight: 700;
  color: #18181b;
  font-size: 16px;
}

.v-name {
  font-size: 14px;
  color: #52525b;
}

.v-desc {
  font-size: 12px;
  color: #a1a1aa;
}

.v-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.v-discount {
  font-weight: 700;
  color: #ef4444;
  font-size: 18px;
}

.no-vouchers {
  padding: 40px 0;
}
</style>

<style>
.el-checkbox-group {
  font-size: 14px;
  line-height: 1.6;
}

.cartCheckbox {
  padding: 0 12px 0 20px;
  margin-right: 8px;
  transform: scale(1.4);
}

.cart-box-shadow {
  border: 1px solid #ccc;
}
</style>
