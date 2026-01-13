<script setup lang="ts">
import Header from '@/components/Header.vue'
import type { CheckboxValueType } from 'element-plus'
import { ref, useTemplateRef, watch, onMounted, computed } from 'vue'
import CartSellerProductWrapper from './components/CartSellerProductWrapper.vue'
import { Ticket, Loading, ShoppingCart } from '@element-plus/icons-vue'
import { ElNotification, ElMessageBox, ElMessage } from 'element-plus'

import { Splide, SplideSlide } from '@splidejs/vue-splide'
import axios from 'axios'
import { formatNumberWithDots } from '@/utils/formatNumberWithDots'
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

const productList = ref<any[]>([]) // Recently viewed

const headerRef = ref<InstanceType<typeof Header> | null>(null)

const checkAll = ref(false)
const checkedProducts = ref<Product[]>([])
const CartSellerProductWrapperRefs = ref<InstanceType<typeof CartSellerProductWrapper>[]>([])
const cartData = ref<CartSeller[]>([])
const isLoading = ref(false)

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

const loadRecentlyViewed = () => {
  const stored = localStorage.getItem('recently_viewed')
  if (stored) {
    productList.value = JSON.parse(stored)
  }
}

const allProducts = computed(() => cartData.value.map((seller) => seller.products).flat())

const handleCheckout = () => {
  if (checkedProducts.value.length === 0) {
    ElMessage.warning('Please select products to checkout')
    return
  }
  localStorage.setItem('checkout_items', JSON.stringify(checkedProducts.value))
  router.push('/checkout')
}

onMounted(() => {
  fetchCart()
  loadRecentlyViewed()
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
    headerRef.value?.fetchCartCount()
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
</script>

<template>
  <Header ref="headerRef" />

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
        <button class="btn">Remove Unactive Products</button>
        <button class="btn">Add To Favorites</button>
      </div>
      <div style="display: flex; justify-content: space-between; align-items: center">
        <el-button
          size="large"
          style="font-size: 18px; color: var(--main-color); padding: 28px 20px 24px"
        >
          <el-icon style="margin-right: 8px; font-size: 24px; position: relative; top: -2px"
            ><Ticket
          /></el-icon>
          Voucher: Apply Voucher
        </el-button>
        <div style="display: flex; align-items: center">
          <span style="font-size: 20px; margin-right: 12px"
            >Total ({{ checkedProducts.length }} products):
            <span style="font-size: 24px; color: var(--main-color)">
              {{
                formatNumberWithDots(
                  checkedProducts.reduce((sum, product) => {
                    return sum + product.price * (product.quantity ?? 1)
                  }, 0),
                )
              }}đ
            </span>
          </span>
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

    <div
      class="box-shadow border-radius recently-viewed"
      style="background-color: #fff; padding: 20px; margin-bottom: 20px; margin-top: 28px"
    >
      <div style="display: flex; margin-bottom: 12px">
        <h3 style="font-weight: bold">YOU RECENTLY VIEWED</h3>
        <RouterLink
          to="/"
          style="position: relative; top: 3px; margin-left: 20px; color: #999; font-size: 13px"
          >View All</RouterLink
        >
      </div>

      <Splide
        v-if="productList.length"
        :options="{
          rewind: true,
          perPage: 5,
          gap: '1rem',
          breakpoints: {
            1000: {
              perPage: 1,
            },
          },
        }"
      >
        <SplideSlide v-for="item in productList" :key="item.id">
          <ProductItem
            :image-url="item.imageUrl"
            :name="item.name"
            :price="item.price"
            :rating="item.rating"
            :location="item.location"
            :discount="item.discount"
            :id="item.id"
          />
        </SplideSlide>
      </Splide>
      <div v-else style="text-align: center; color: #999; padding: 40px">
        <el-icon size="40" style="margin-bottom: 12px; opacity: 0.5"><ShoppingCart /></el-icon>
        <p style="font-size: 14px">Your recently viewed list is empty.</p>
      </div>
    </div>
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
