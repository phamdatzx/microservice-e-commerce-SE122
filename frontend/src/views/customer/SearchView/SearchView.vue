<script setup lang="ts">
import { ref, computed } from 'vue'
import ProductList from '@/components/ProductList.vue'
import Header from '@/components/Header.vue'
import RecentlyViewed from '@/components/RecentlyViewed.vue'
import { ArrowDownBold, ArrowUpBold, Filter, Search } from '@element-plus/icons-vue'

const categories = [
  'Computer Components',
  'Printers & Scanners',
  'Computer Accessories',
  'Network Devices',
  'Laptops',
  'Software',
  'Tablets',
  'Storage Devices',
  'Monitors',
  'Input Devices',
]

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
  showAllCategories.value ? categories : categories.slice(0, 4),
)
const displayedLocations = computed(() =>
  showAllLocations.value ? locations : locations.slice(0, 4),
)

const ratings = [5, 4, 3, 2, 1]
// ... existing products ...
const products = [
  {
    id: 1,
    name: 'Jeans for Men Skinny Fit White Plain Stretchable Office Stylish Fashion',
    price: 149600,
    imageUrl: '/src/assets/product-imgs/product1.png',
    rating: 4.8,
    location: 'Ho Chi Minh City',
    discount: 28,
  },
  {
    id: 2,
    name: 'Cargo Jeans for Men Wide Leg Relaxed Fit Loose Fashion Trendy',
    price: 175000,
    imageUrl: '/src/assets/product-imgs/product1.png',
    rating: 4.8,
    location: 'Ha Noi',
    discount: 32,
  },
  {
    id: 3,
    name: 'Skinny Jeans Men Black Washed Casual Comfort Fit Elastic',
    price: 125000,
    imageUrl: '/src/assets/product-imgs/product1.png',
    rating: 4.8,
    location: 'Da Nang',
    discount: 37,
  },
  {
    id: 4,
    name: 'WASH RETRO Jeans Men Dark Grey Loose Fit Vintage Style',
    price: 174999,
    imageUrl: '/src/assets/product-imgs/product1.png',
    rating: 4.8,
    location: 'Ho Chi Minh City',
    discount: 24,
  },
  {
    id: 5,
    name: 'Elastic Cuff Pants for Men Casual Streetwear Jogger Style',
    price: 92900,
    imageUrl: '/src/assets/product-imgs/product1.png',
    rating: 4.9,
    location: 'Binh Duong',
    discount: 45,
  },
  {
    id: 6,
    name: 'Premium Cotton T-Shirt for Men Soft Breathable Casual Wear',
    price: 85000,
    imageUrl: '/src/assets/product-imgs/product1.png',
    rating: 4.7,
    location: 'Ho Chi Minh City',
    discount: 15,
  },
  {
    id: 7,
    name: 'Digital Camera 4K Ultra HD Professional Video Vlogging',
    price: 4500000,
    imageUrl: '/src/assets/product-imgs/product1.png',
    rating: 4.9,
    location: 'Ho Chi Minh City',
    discount: 10,
  },
  {
    id: 8,
    name: 'Wireless Bluetooth Headphones Noise Cancelling Over-Ear',
    price: 1250000,
    imageUrl: '/src/assets/product-imgs/product1.png',
    rating: 4.8,
    location: 'Ha Noi',
    discount: 20,
  },
  {
    id: 9,
    name: 'Smart Watch Series 7 Waterproof Fitness Tracker Heart Rate',
    price: 850000,
    imageUrl: '/src/assets/product-imgs/product1.png',
    rating: 4.6,
    location: 'Da Nang',
    discount: 15,
  },
  {
    id: 10,
    name: 'Portable Power Bank 20000mAh Fast Charging Dual USB',
    price: 450000,
    imageUrl: '/src/assets/product-imgs/product1.png',
    rating: 4.7,
    location: 'Binh Duong',
    discount: 30,
  },
  {
    id: 11,
    name: 'Mechanical Gaming Keyboard RGB Backlit Wired Blue Switch',
    price: 750000,
    imageUrl: '/src/assets/product-imgs/product1.png',
    rating: 4.8,
    location: 'Ho Chi Minh City',
    discount: 25,
  },
  {
    id: 12,
    name: 'Optical Wireless Mouse Ergonmic Design 2400 DPI',
    price: 150000,
    imageUrl: '/src/assets/product-imgs/product1.png',
    rating: 4.5,
    location: 'Ha Noi',
    discount: 40,
  },
  {
    id: 13,
    name: 'Laptop Bag 15.6 inch Waterproof Shoulder Messenger Bag',
    price: 350000,
    imageUrl: '/src/assets/product-imgs/product1.png',
    rating: 4.8,
    location: 'Da Nang',
    discount: 10,
  },
  {
    id: 14,
    name: 'USB-C Hub Multiport Adapter with 4K HDMI USB 3.0',
    price: 650000,
    imageUrl: '/src/assets/product-imgs/product1.png',
    rating: 4.7,
    location: 'Ho Chi Minh City',
    discount: 15,
  },
  {
    id: 15,
    name: 'External Hard Drive 2TB USB 3.0 Portable Storage',
    price: 2100000,
    imageUrl: '/src/assets/product-imgs/product1.png',
    rating: 4.9,
    location: 'Ha Noi',
    discount: 5,
  },
  {
    id: 16,
    name: 'Webcam 1080P Full HD with Microphone for PC Laptop',
    price: 550000,
    imageUrl: '/src/assets/product-imgs/product1.png',
    rating: 4.6,
    location: 'Binh Duong',
    discount: 20,
  },
  {
    id: 17,
    name: 'Bluetooth Speaker Portable Waterproof Outdoor Wireless',
    price: 890000,
    imageUrl: '/src/assets/product-imgs/product1.png',
    rating: 4.8,
    location: 'Dong Nai',
    discount: 35,
  },
  {
    id: 18,
    name: 'LED Desk Lamp with USB Charging Port Touch Control',
    price: 250000,
    imageUrl: '/src/assets/product-imgs/product1.png',
    rating: 4.7,
    location: 'Ba Ria - Vung Tau',
    discount: 45,
  },
  {
    id: 19,
    name: 'Monitor Stand Riser with Desk Organizer Storage',
    price: 180000,
    imageUrl: '/src/assets/product-imgs/product1.png',
    rating: 4.5,
    location: 'Ho Chi Minh City',
    discount: 50,
  },
  {
    id: 20,
    name: 'Cat 7 Ethernet Cable 10Gbps High Speed RJ45 Network',
    price: 120000,
    imageUrl: '/src/assets/product-imgs/product1.png',
    rating: 4.9,
    location: 'Ha Noi',
    discount: 10,
  },
]

const sortOptions = [
  { label: 'Relevance', value: 'relevance', active: true },
  { label: 'Latest', value: 'latest' },
  { label: 'Top Sales', value: 'top-sales' },
]

const pagination = {
  current: 1,
  total: 340,
  pageSize: 20,
}

const priceMin = ref<number>()
const priceMax = ref<number>()

const clearAll = () => {
  priceMin.value = undefined
  priceMax.value = undefined
}
</script>

<template>
  <Header />
  <div class="search-view-main">
    <ProductList :products="products" :sort-options="sortOptions" :pagination="pagination">
      <!-- Top Slot -->
      <template #top>
        <div class="search-summary">
          <el-icon class="info-icon"><Search /></el-icon>
          <span>Search result for keyword '<span class="keyword">laptop</span>'</span>
        </div>
      </template>

      <!-- Sidebar Slot -->
      <template #sidebar>
        <div class="filter-header">
          <el-icon><Filter /></el-icon>
          <span class="filter-title">SEARCH FILTER</span>
        </div>
        <!-- ... locations ... -->
        <div class="filter-group">
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
        </div>

        <!-- By Category -->
        <div class="filter-group">
          <h3 class="group-title">By Category</h3>
          <div v-for="cat in displayedCategories" :key="cat" class="checkbox-item">
            <el-checkbox>{{ cat }}</el-checkbox>
          </div>
          <button class="more-btn" @click="showAllCategories = !showAllCategories">
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
          <button class="apply-btn">APPLY</button>
        </div>

        <!-- Rating -->
        <div class="filter-group">
          <h3 class="group-title">Rating</h3>
          <div v-for="star in ratings" :key="star" class="rating-item">
            <el-checkbox>
              <div class="stars">
                <span v-for="n in 5" :key="n" class="star" :class="{ filled: n <= star }">★</span>
                <span v-if="star < 5" class="and-up">& up</span>
              </div>
            </el-checkbox>
          </div>
        </div>

        <button class="clear-all-btn" @click="clearAll" style="margin-bottom: 30px">
          CLEAR ALL
        </button>
      </template>
    </ProductList>

    <RecentlyViewed style="margin-top: 40px" />
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
