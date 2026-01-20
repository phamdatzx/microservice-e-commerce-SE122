<script setup lang="ts">
import { ref, onMounted, computed, watch, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import axios from 'axios'
import SellerHeader from './components/SellerHeader.vue'
import SellerVouchers from './components/SellerVouchers.vue'
import SellerRecommendations from './components/SellerRecommendations.vue'
import ProductList from '@/components/ProductList.vue'
import NotFoundView from '@/components/NotFoundView.vue'

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
    product_count: number
    is_following: boolean
  }
}

interface SellerCategory {
  id: string
  seller_id: string
  name: string
  product_count: number
}

const sellerInfo = ref<SellerInfo | null>(null)
const products = ref<any[]>([])
const vouchers = ref<any[]>([])
const savedVoucherIds = ref<string[]>([])
const isLoading = ref(false)
const isNotFound = ref(false)
const isFollowLoading = ref(false)
const totalItems = ref(0)
const currentPage = ref(1)
const pageSize = ref(15)
const sortBy = ref('sold_count')
const sortDirection = ref('desc')
const sellerCategories = ref<SellerCategory[]>([])
const selectedCategoryId = ref<string | null>(null)
const activeTab = ref<string>('home')
const productListRef = ref<any>(null)

const fetchSellerInfo = async () => {
  if (!sellerId.value) return
  const token = localStorage.getItem('access_token')
  const headers = token ? { Authorization: `Bearer ${token}` } : {}
  isNotFound.value = false
  try {
    const response = await axios.get(
      `${import.meta.env.VITE_BE_API_URL}/user/public/seller/${sellerId.value}`,
      { headers },
    )
    if (response.data && response.data.data) {
      sellerInfo.value = response.data.data
    }
  } catch (error: any) {
    if (error.response?.status === 404 || error.response?.status === 500) {
      isNotFound.value = true
    }
    console.error('Error fetching seller info:', error)
  }
}

const fetchSellerCategories = async () => {
  if (!sellerId.value) return
  try {
    const response = await axios.get(
      `${import.meta.env.VITE_BE_API_URL}/product/public/seller/${sellerId.value}/category`,
    )
    if (response.data) {
      sellerCategories.value = response.data
    }
  } catch (error) {
    console.error('Error fetching seller categories:', error)
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
          seller_category: selectedCategoryId.value || undefined,
        },
      },
    )
    if (response.data && response.data.products) {
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
    } else {
      products.value = []
      totalItems.value = 0
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

const handleCategoryChange = async (
  categoryId: string | null,
  scrollTarget: 'top' | 'products' = 'products',
) => {
  const isCategoryChanging = selectedCategoryId.value !== categoryId
  selectedCategoryId.value = categoryId
  currentPage.value = 1

  if (categoryId === null) {
    activeTab.value = scrollTarget === 'top' ? 'home' : 'all'
  } else {
    activeTab.value = categoryId
  }

  // Only refetch if filter changed OR we are navigating to the product list
  if (isCategoryChanging || scrollTarget === 'products') {
    await fetchSellerProducts()
  }

  await nextTick()
  if (scrollTarget === 'top') {
    window.scrollTo({ top: 0, behavior: 'smooth' })
  } else if (productListRef.value) {
    productListRef.value.scrollIntoView({ behavior: 'smooth' })
  }
}

onMounted(() => {
  fetchSellerInfo()
  fetchSellerProducts()
  fetchSellerVouchers()
  fetchSavedVouchers()
  fetchSellerCategories()
})

watch(sellerId, () => {
  selectedCategoryId.value = null
  currentPage.value = 1
  fetchSellerInfo()
  fetchSellerProducts()
  fetchSellerVouchers()
  fetchSavedVouchers()
  fetchSellerCategories()
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
    <main v-if="!isNotFound" class="seller-main">
      <SellerHeader
        v-if="sellerInfo"
        :seller-info="sellerInfo"
        :categories="sellerCategories"
        :selected-category-id="selectedCategoryId"
        :active-tab="activeTab"
        :is-follow-loading="isFollowLoading"
        @toggle-follow="handleFollowToggle"
        @category-change="handleCategoryChange"
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

        <div ref="productListRef">
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
                <a
                  href="javascript:void(0)"
                  class="category-item"
                  :class="{ active: !selectedCategoryId }"
                  @click="handleCategoryChange(null)"
                  >All Products</a
                >
                <a
                  v-for="cat in sellerCategories"
                  :key="cat.id"
                  href="javascript:void(0)"
                  class="category-item"
                  :class="{ active: selectedCategoryId === cat.id }"
                  @click="handleCategoryChange(cat.id)"
                >
                  {{ cat.name }} ({{ cat.product_count }})
                </a>
              </nav>
            </template>
          </ProductList>
        </div>
      </div>
    </main>

    <NotFoundView
      v-else
      title="Seller Not Found"
      message="The seller you are looking for does not exist or has been deactivated."
    />
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
