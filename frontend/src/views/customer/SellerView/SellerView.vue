<script setup lang="ts">
import SellerHeader from './components/SellerHeader.vue'
import SellerVouchers from './components/SellerVouchers.vue'
import SellerRecommendations from './components/SellerRecommendations.vue'
import ProductList from '@/components/ProductList.vue'

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
    name: 'Cargo Jeans for Men Wide Leg Relaxed Fit Loose Fashion Trendy',
    price: 175000,
    imageUrl: '/src/assets/product-imgs/product1.png',
    rating: 4.8,
    location: 'Ha Noi',
    discount: 32,
  },
]

const sortOptions = [
  { label: 'Popular', value: 'popular', active: true },
  { label: 'Latest', value: 'latest' },
  { label: 'Top Sales', value: 'top-sales' },
]

const pagination = {
  current: 1,
  total: 100,
  pageSize: 20,
}
</script>

<template>
  <div class="seller-page">
    <main class="seller-main">
      <SellerHeader />

      <div class="content-bg">
        <SellerVouchers />

        <div class="container">
          <SellerRecommendations />
        </div>

        <ProductList :products="products" :sort-options="sortOptions" :pagination="pagination">
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

:deep(.custom-pagination.is-background .el-pager li:not(.is-disabled):hover) {
  color: var(--main-color);
}
</style>
