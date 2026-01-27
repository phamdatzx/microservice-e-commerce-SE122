<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import ProductList from '@/components/ProductList.vue'
import RecentlyViewed from '@/components/RecentlyViewed.vue'
import { ArrowDownBold, ArrowUpBold, Filter, Search } from '@element-plus/icons-vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'
import { Loading } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()

interface Category {
  id: string
  name: string
  image: string
  product_count: number
}

const categories = ref<Category[]>([])

const locations = [
  'Ho Chi Minh City',
  'Binh Duong',
  'Dong Nai',
  'Ba Ria - Vung Tau',
  'Ha Noi',
  'Da Nang',
  'Can Tho',
  'Hai Phong',
]

const showAllCategories = ref(false)
const showAllLocations = ref(false)

const displayedCategories = computed(() =>
  showAllCategories.value ? categories.value : categories.value.slice(0, 4),
)
const displayedLocations = computed(() =>
  showAllLocations.value ? locations : locations.slice(0, 4),
)

interface Product {
  id: string
  name: string
  price: {
    min: number
    max: number
  }
  images: Array<{ url: string; order: number }>
  rating: number
  // Mock fields for UI compatibility if needed, or derived
  location: string
  discount: number
}

// Derived interface for ProductList component
interface DisplayProduct {
  id: string
  name: string
  minPrice: number
  maxPrice: number
  imageUrl: string
  rating: number
  location: string
  soldCount?: number
}

const products = ref<DisplayProduct[]>([])
const isLoading = ref(false)
const totalItems = ref(0)

const sortOptions = computed(() => [
  { label: 'Top Sales', value: 'top_sales', active: currentSort.value === 'top_sales' },
  { label: 'Rating', value: 'rating', active: currentSort.value === 'rating' },
])

const currentSort = ref('top_sales')

const pagination = computed(() => ({
  current: Number(route.query.page) || 1,
  total: totalItems.value,
  pageSize: Number(route.query.limit) || 20,
}))

const priceMin = ref<number | undefined>(undefined)
const priceMax = ref<number | undefined>(undefined)
const selectedCategory = ref<string[]>([])
const selectedLocation = ref<string | undefined>(undefined)
const selectedRating = ref<number | undefined>(undefined)

const clearAll = () => {
  priceMin.value = undefined
  priceMax.value = undefined
  selectedCategory.value = []
  selectedLocation.value = undefined
  selectedRating.value = undefined
  applyFilters()
}

const fetchProducts = async () => {
  isLoading.value = true
  try {
    const params = {
      page: route.query.page || 1,
      limit: route.query.limit || 20,
      search_query: route.query.q || '',
      min_price: route.query.min_price,
      max_price: route.query.max_price,
      min_rating: route.query.min_rating,
      category_ids: route.query.category_ids,
      sort_by: route.query.sort_by,
      sort_direction: route.query.sort_direction,
    }

    // Handle special sort cases
    if (currentSort.value === 'price_asc') {
      params.sort_by = 'price'
      params.sort_direction = 'asc'
    } else if (currentSort.value === 'price_desc') {
      params.sort_by = 'price'
      params.sort_direction = 'desc'
    } else if (currentSort.value === 'rating') {
      params.sort_by = 'rating'
      params.sort_direction = 'desc'
    } else if (currentSort.value === 'top_sales') {
      params.sort_by = 'sold_count'
      params.sort_direction = 'desc'
    }

    const response = await axios.get(`${import.meta.env.VITE_BE_API_URL}/product/public/search`, {
      params,
      headers: {
        'X-User-Id': localStorage.getItem('user_id'),
      },
    })

    const rawProducts: Product[] = response.data.products
    totalItems.value = response.data.pagination.total_items

    // Transform to display format
    products.value = rawProducts.map((p: any) => ({
      id: p.id,
      name: p.name,
      minPrice: p.price.min,
      maxPrice: p.price.max,
      imageUrl:
        p.images && p.images.length > 0
          ? p.images.sort((a: any, b: any) => a.order - b.order)[0].url
          : 'https://placehold.co/300x300?text=No+Image',
      rating: p.rating,
      location: 'Vietnam', // Backend doesn't have location yet
      soldCount: p.sold_count,
    }))
  } catch (error) {
    console.error('Failed to fetch products:', error)
    products.value = []
    totalItems.value = 0
  } finally {
    isLoading.value = false
  }
}

const fetchCategories = async () => {
  try {
    const response = await axios.get(`${import.meta.env.VITE_BE_API_URL}/product/public/category`)
    categories.value = response.data
  } catch (error) {
    console.error('Error fetching categories:', error)
  }
}

const applyFilters = () => {
  const query: any = { ...route.query, page: 1 } // Reset to page 1 on filter change

  if (priceMin.value) query.min_price = priceMin.value
  else delete query.min_price

  if (priceMax.value) query.max_price = priceMax.value
  else delete query.max_price

  if (selectedRating.value) query.min_rating = selectedRating.value
  else delete query.min_rating

  if (selectedCategory.value.length > 0) query.category_ids = selectedCategory.value.join(',')
  else delete query.category_ids

  router.push({ query })
}

const toggleCategory = (id: string, checked: boolean) => {
  if (checked) {
    if (!selectedCategory.value.includes(id)) {
      selectedCategory.value.push(id)
    }
  } else {
    selectedCategory.value = selectedCategory.value.filter((categoryId) => categoryId !== id)
  }
  applyFilters()
}

const handleSortChange = (val: string) => {
  currentSort.value = val
  if (val === 'price-asc') {
    currentSort.value = 'price_asc'
  } else if (val === 'price-desc') {
    currentSort.value = 'price_desc'
  }

  // Sync with URL
  const query: any = { ...route.query, sort_by: undefined, sort_direction: undefined, page: 1 }
  if (currentSort.value === 'price_asc') {
    query.sort_by = 'price'
    query.sort_direction = 'asc'
  } else if (currentSort.value === 'price_desc') {
    query.sort_by = 'price'
    query.sort_direction = 'desc'
  } else if (currentSort.value === 'rating') {
    query.sort_by = 'rating'
    query.sort_direction = 'desc'
  } else if (currentSort.value === 'top_sales') {
    query.sort_by = 'sold_count'
    query.sort_direction = 'desc'
  }

  router.push({ query })
}

const handlePageChange = (page: number) => {
  router.push({
    query: { ...route.query, page },
  })
}

// Watch for route changes to refetch
watch(
  () => route.query,
  () => {
    // Sync local filters with URL query on load/change
    priceMin.value = route.query.min_price ? Number(route.query.min_price) : undefined
    priceMax.value = route.query.max_price ? Number(route.query.max_price) : undefined
    selectedRating.value = route.query.min_rating ? Number(route.query.min_rating) : undefined
    selectedCategory.value = route.query.category_ids
      ? (route.query.category_ids as string).split(',')
      : []

    // Sync currentSort based on URL
    if (route.query.sort_by === 'price') {
      currentSort.value = route.query.sort_direction === 'asc' ? 'price_asc' : 'price_desc'
    } else if (route.query.sort_by === 'rating') {
      currentSort.value = 'rating'
    } else if (route.query.sort_by === 'sold_count') {
      currentSort.value = 'top_sales'
    } else {
      currentSort.value = 'top_sales'
    }

    fetchProducts()
  },
  { immediate: true },
)

onMounted(() => {
  fetchCategories()
})
</script>

<template>
  <div class="search-view-main">
    <ProductList
      :products="products"
      :sort-options="sortOptions"
      :pagination="pagination"
      :loading="isLoading"
      @sort-change="handleSortChange"
      @page-change="handlePageChange"
    >
      <!-- Top Slot -->
      <template #top>
        <div class="search-summary">
          <el-icon class="info-icon"><Search /></el-icon>
          <span v-if="route.query.q"
            >Search result for keyword '<span class="keyword">{{ route.query.q }}</span
            >'</span
          >
          <span v-else>All Products</span>
        </div>
      </template>

      <!-- Sidebar Slot -->
      <template #sidebar>
        <div class="filter-header">
          <el-icon><Filter /></el-icon>
          <span class="filter-title">SEARCH FILTER</span>
        </div>
        <!-- ... locations ... -->
        <!-- <div class="filter-group">
          <h3 class="group-title">Shipped From</h3>
          <div v-for="loc in displayedLocations" :key="loc" class="checkbox-item">
            <el-checkbox>{{ loc }}</el-checkbox>
          </div>
          <button class="more-btn" @click="showAllLocations = !showAllLocations">
            {{ showAllLocations ? 'Less' : 'More' }}
            <el-icon>
              <component :is="showAllLocations ? ArrowUpBold : ArrowDownBold" />
            </el-icon>
          </button>
        </div> -->

        <!-- By Category -->
        <div class="filter-group">
          <h3 class="group-title">By Category</h3>
          <div v-for="cat in displayedCategories" :key="cat.id" class="checkbox-item">
            <el-checkbox
              :model-value="selectedCategory.includes(cat.id)"
              @change="(val: string | number | boolean) => toggleCategory(cat.id, val as boolean)"
            >
              {{ cat.name }}
            </el-checkbox>
          </div>
          <button
            v-if="categories.length > 4"
            class="more-btn"
            @click="showAllCategories = !showAllCategories"
          >
            {{ showAllCategories ? 'Less' : 'More' }}
            <el-icon>
              <component :is="showAllCategories ? ArrowUpBold : ArrowDownBold" />
            </el-icon>
          </button>
        </div>

        <!-- Price Range -->
        <div class="filter-group">
          <h3 class="group-title">Price Range</h3>
          <div class="price-inputs">
            <el-input-number
              v-model="priceMin"
              size="large"
              :controls="false"
              placeholder="Min"
              class="price-input"
            />
            <span class="separator">-</span>
            <el-input-number
              v-model="priceMax"
              size="large"
              :controls="false"
              placeholder="Max"
              class="price-input"
            />
          </div>
          <button class="apply-btn" @click="applyFilters">APPLY</button>
        </div>

        <!-- Rating -->
        <div class="filter-group">
          <h3 class="group-title">Rating</h3>
          <div v-for="star in [5, 4, 3, 2, 1]" :key="star" class="rating-item">
            <el-radio
              v-model="selectedRating"
              :label="star"
              @change="applyFilters"
              style="height: auto; margin-right: 0"
            >
              <div class="stars">
                <span v-for="n in 5" :key="n" class="star" :class="{ filled: n <= star }">★</span>
                <span v-if="star < 5" class="and-up">& up</span>
              </div>
            </el-radio>
          </div>
        </div>

        <button class="clear-all-btn" @click="clearAll" style="margin-bottom: 30px">
          CLEAR ALL
        </button>
      </template>

      <template #bottom>
        <RecentlyViewed style="margin-top: 20px" />
      </template>
    </ProductList>
  </div>
</template>

<style scoped>
.filter-header {
  display: flex;
  align-items: center;
  gap: 10px;
  font-weight: 700;
  font-size: 16px;
  margin-bottom: 20px;
  color: #333;
}

.filter-group {
  margin-bottom: 25px;
  padding-bottom: 4px;
  border-bottom: 1px solid #e5e5e5;
}

.group-title {
  font-size: 14px;
  font-weight: 600;
  margin-bottom: 12px;
  color: #333;
}

.checkbox-item {
  margin-bottom: 8px;
}

:deep(.el-checkbox__label) {
  font-size: 14px;
  color: #333;
}

.more-btn {
  background: none;
  border: none;
  color: #555;
  font-size: 14px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 5px;
  padding: 6px 26px 10px;
  transition: color 0.2s;
}

.more-btn:hover {
  color: var(--main-color);
}

.price-inputs {
  display: flex;
  align-items: center;
  gap: 5px;
  margin-bottom: 10px;
}

.price-input {
  width: 100%;
}

:deep(.price-input .el-input__wrapper) {
  padding: 0 8px;
  border-radius: 2px;
  box-shadow: 0 0 0 1px #dbdbdb inset;
}

:deep(.price-input .el-input__inner) {
  font-size: 12px;
  text-align: left;
}

.separator {
  color: #bdbdbd;
}

.apply-btn {
  width: 100%;
  padding: 8px;
  background-color: var(--main-color);
  color: white;
  border: none;
  border-radius: 2px;
  font-weight: 600;
  cursor: pointer;
  font-size: 12px;
}

.rating-item {
  margin-bottom: 5px;
}

.stars {
  display: flex;
  align-items: center;
  gap: 2px;
}

.star {
  color: #ccc;
  font-size: 16px;
}

.star.filled {
  color: #ffce3d;
}

.and-up {
  font-size: 12px;
  margin-left: 5px;
  color: #333;
}

.clear-all-btn {
  width: 100%;
  padding: 10px;
  background-color: var(--main-color);
  color: white;
  border: none;
  border-radius: 2px;
  font-weight: 600;
  cursor: pointer;
  margin-top: 10px;
}

.search-summary {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 16px;
  margin-bottom: 0;
  color: #555;
}

.keyword {
  color: var(--main-color);
  text-transform: uppercase;
  font-weight: 600;
}
</style>
