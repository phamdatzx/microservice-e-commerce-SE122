<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import axios from 'axios'
import SellerHeader from './components/SellerHeader.vue'
import SellerVouchers from './components/SellerVouchers.vue'
import SellerRecommendations from './components/SellerRecommendations.vue'
import ProductList from '@/components/ProductList.vue'

const route = useRoute()
const sellerId = computed(() => route.params.sellerId as string)

interface SellerInfo {
  id: string
  name: string
  image: string
  address: {
    province: string
  }
  sale_info: {
    follow_count: number
    rating_count: number
    rating_average: number
    product_count: number // Added if available or will derive
    is_following: boolean
  }
}

const sellerInfo = ref<SellerInfo | null>(null)
const products = ref<any[]>([])
const vouchers = ref<any[]>([])
const savedVoucherIds = ref<string[]>([])
const isLoading = ref(false)
const isFollowLoading = ref(false)
const totalItems = ref(0)
const currentPage = ref(1)
const pageSize = ref(15)
const sortBy = ref('sold_count')
const sortDirection = ref('desc')

const categories = [
  'Cargo Jeans ⭐',
  'Skinny Jeans ⭐',
  'WASH RETRO ⭐',
  'Shorts ⭐',
  'Trousers ⭐',
  'Jackets ⭐',
  'Elastic Cuff Pants ⭐',
  'BAGGY ⭐',
  'Windbreakers ⭐',
  'Thermals ⭐',
]

const fetchSellerInfo = async () => {
  if (!sellerId.value) return
  const token = localStorage.getItem('access_token')
  const headers = token ? { Authorization: `Bearer ${token}` } : {}
  try {
    const response = await axios.get(
      `${import.meta.env.VITE_BE_API_URL}/user/public/seller/${sellerId.value}`,
      { headers },
    )
    if (response.data && response.data.data) {
      sellerInfo.value = response.data.data
    }
  } catch (error) {
    console.error('Error fetching seller info:', error)
  }
}

interface SellerProduct {
  id: string
  name: string
  price: {
    min: number
    max: number
  }
  images: Array<{
    url: string
    order: number
  }>
  rating: number
  sold_count: number
}

const fetchSellerProducts = async () => {
  if (!sellerId.value) return
  isLoading.value = true
  try {
    const response = await axios.get(
      `${import.meta.env.VITE_BE_API_URL}/product/public/products/seller/${sellerId.value}`,
      {
        params: {
          page: currentPage.value,
          limit: pageSize.value,
          sort_by: sortBy.value,
          sort_direction: sortDirection.value,
        },
      },
    )
    if (response.data) {
      products.value = response.data.products.map((p: SellerProduct) => ({
        id: p.id,
        name: p.name,
        price: p.price.min,
        imageUrl:
          p.images?.sort((a, b) => a.order - b.order)[0]?.url ||
          'https://placehold.co/300x300?text=No+Image',
        rating: p.rating,
        location: sellerInfo.value?.address?.province || 'Vietnam',
        discount: 0,
        soldCount: p.sold_count,
      }))
      totalItems.value = response.data.pagination.total_items
    }
  } catch (error) {
    console.error('Error fetching seller products:', error)
  } finally {
    isLoading.value = false
  }
}

import { ElMessage } from 'element-plus'

const fetchSellerVouchers = async () => {
  if (!sellerId.value) return
  try {
    const response = await axios.get(
      `${import.meta.env.VITE_BE_API_URL}/product/public/vouchers/seller/${sellerId.value}`,
    )
    if (response.data) {
      vouchers.value = response.data
    }
  } catch (error) {
    console.error('Error fetching seller vouchers:', error)
  }
}

const fetchSavedVouchers = async () => {
  const token = localStorage.getItem('access_token')
  if (!token) return
  try {
    const response = await axios.get(`${import.meta.env.VITE_BE_API_URL}/product/saved-vouchers`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
    if (response.data) {
      savedVoucherIds.value = response.data.map((item: any) => item.voucher_id)
    }
  } catch (error) {
    console.error('Error fetching saved vouchers:', error)
  }
}

const handleFollowToggle = async () => {
  const token = localStorage.getItem('access_token')
  if (!token) {
    ElMessage.warning('Please login to follow sellers')
    return
  }

  if (!sellerInfo.value) return

  const isCurrentlyFollowing = sellerInfo.value.sale_info.is_following
  const url = `${import.meta.env.VITE_BE_API_URL}/user/follow/${sellerId.value}`

  isFollowLoading.value = true
  try {
    if (isCurrentlyFollowing) {
      await axios.delete(url, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
      ElMessage.success('Unfollowed seller')
    } else {
      await axios.post(
        url,
        {},
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        },
      )
      ElMessage.success('Followed seller')
    }
    // Refresh seller info to update follow count and status
    await fetchSellerInfo()
  } catch (error: any) {
    const errorMsg = error.response?.data?.message || 'Failed to update follow status'
    ElMessage.error(errorMsg)
  } finally {
    isFollowLoading.value = false
  }
}

const handleSaveVoucher = async (voucherId: string) => {
  const token = localStorage.getItem('access_token')
  if (!token) {
    ElMessage.warning('Please login to save vouchers')
    return
  }

  try {
    await axios.post(
      `${import.meta.env.VITE_BE_API_URL}/product/saved-vouchers`,
      {
        voucher_id: voucherId,
      },
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      },
    )
    ElMessage.success('Voucher saved successfully!')
    // Refresh vouchers and saved list to update UI
    fetchSellerVouchers()
    fetchSavedVouchers()
  } catch (error: any) {
    const errorMsg = error.response?.data?.message || 'Failed to save voucher'
    ElMessage.error(errorMsg)
  }
}

const handlePageChange = (page: number) => {
  currentPage.value = page
  fetchSellerProducts()
}

const handleSortChange = (value: string) => {
  if (value === 'price-asc') {
    sortBy.value = 'price'
    sortDirection.value = 'asc'
  } else if (value === 'price-desc') {
    sortBy.value = 'price'
    sortDirection.value = 'desc'
  } else {
    sortBy.value = value
    sortDirection.value = 'desc'
  }
  currentPage.value = 1
  fetchSellerProducts()
}

onMounted(() => {
  fetchSellerInfo()
  fetchSellerProducts()
  fetchSellerVouchers()
  fetchSavedVouchers()
})

const sortOptions = computed(() => [
  { label: 'Top Sales', value: 'sold_count', active: sortBy.value === 'sold_count' },
  { label: 'Rating', value: 'rating', active: sortBy.value === 'rating' },
])

const pagination = computed(() => ({
  current: currentPage.value,
  total: totalItems.value,
  pageSize: pageSize.value,
}))
</script>

<template>
  <div class="seller-page">
    <main class="seller-main">
      <SellerHeader
        v-if="sellerInfo"
        :seller-info="sellerInfo"
        :is-follow-loading="isFollowLoading"
        @toggle-follow="handleFollowToggle"
      />

      <div class="content-bg">
        <SellerVouchers
          :vouchers="vouchers"
          :saved-voucher-ids="savedVoucherIds"
          @save="handleSaveVoucher"
        />

        <div class="container">
          <SellerRecommendations />
        </div>

        <ProductList
          :products="products"
          :sort-options="sortOptions"
          :pagination="pagination"
          :loading="isLoading"
          @page-change="handlePageChange"
          @sort-change="handleSortChange"
        >
          <template #sidebar>
            <div class="sidebar-header">
              <span class="sidebar-icon">≡</span>
              <span class="sidebar-title">Categories</span>
            </div>
            <nav class="category-nav">
              <a href="#" class="category-item active">All Products</a>
              <a v-for="cat in categories" :key="cat" href="#" class="category-item">
                {{ cat }}
              </a>
            </nav>
          </template>
        </ProductList>
      </div>
    </main>
  </div>
</template>

<style scoped>
.seller-page {
  background-color: #f5f5f5;
  min-height: 100vh;
}

.content-bg {
  background-color: #f5f5f5;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 15px;
}

.sidebar-header {
  display: flex;
  align-items: center;
  font-size: 16px;
  font-weight: 700;
  padding-bottom: 15px;
  border-bottom: 1px solid #e5e5e5;
  margin-bottom: 15px;
}

.sidebar-icon {
  margin-right: 10px;
  font-size: 18px;
}

.category-nav {
  display: flex;
  flex-direction: column;
}

.category-item {
  text-decoration: none;
  color: #333;
  font-size: 14px;
  padding: 8px 0;
  display: flex;
  align-items: center;
  transition: color 0.2s;
}

.category-item:hover {
  color: var(--main-color);
}

.category-item.active {
  color: var(--main-color);
  font-weight: 600;
}

.category-item.active::before {
  content: '▶';
  font-size: 10px;
  margin-right: 8px;
}

:deep(.custom-pagination.is-background .el-pager li:not(.is-disabled).is-active) {
  background-color: var(--main-color) !important;
  color: #fff;
}

:deep(.custom-pagination.is-background .el-pager li:not(.is-disabled):not(.is-active):hover) {
  color: var(--main-color);
}
</style>
