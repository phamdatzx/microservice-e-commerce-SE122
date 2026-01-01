<script setup lang="ts">
import ProductItem from '@/components/ProductItem.vue'
import { ArrowDownBold, ArrowLeft, ArrowRight } from '@element-plus/icons-vue'

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
]
</script>

<template>
  <div class="product-list-section">
    <div class="container main-content">
      <!-- Sidebar -->
      <aside class="sidebar">
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
      </aside>

      <!-- Main Grid Area -->
      <main class="product-grid-area">
        <!-- Sort Bar -->
        <div class="sort-bar">
          <div class="sort-options">
            <span class="sort-label">Sort by:</span>
            <button class="sort-btn active">Popular</button>
            <button class="sort-btn">Latest</button>
            <button class="sort-btn">Best Selling</button>
            <div class="price-dropdown" style="">
              Price <el-icon><ArrowDownBold /></el-icon>
            </div>
          </div>
          <div class="pagination-summary">
            <span class="page-current">1</span>/<span class="page-total">2</span>
            <div class="page-btns">
              <button class="page-btn disabled">
                <el-icon><ArrowLeft /></el-icon>
              </button>
              <button class="page-btn">
                <el-icon><ArrowRight /></el-icon>
              </button>
            </div>
          </div>
        </div>

        <!-- Grid -->
        <div class="product-grid">
          <ProductItem
            v-for="product in products"
            :key="product.id"
            :id="product.id"
            :name="product.name"
            :imageUrl="product.imageUrl"
            :price="product.price"
            :rating="product.rating"
            :location="product.location"
            :discount="product.discount"
          />
        </div>
      </main>
    </div>
  </div>
</template>

<style scoped>
.product-list-section {
  padding: 12px 0 20px;
  background-color: #f5f5f5;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 15px;
}

.main-content {
  display: flex;
  gap: 20px;
}

/* Sidebar */
.sidebar {
  width: 190px;
  flex-shrink: 0;
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

/* Main Area */
.product-grid-area {
  flex: 1;
}

.sort-bar {
  background-color: #ededed;
  padding: 13px 20px;
  border-radius: 2px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.sort-options {
  display: flex;
  align-items: center;
  gap: 10px;
}

.sort-label {
  font-size: 14px;
  color: #555;
  margin-right: 5px;
}

.sort-btn {
  background: white;
  border: none;
  padding: 7px 15px;
  border-radius: 2px;
  font-size: 14px;
  cursor: pointer;
  color: #333;
  line-height: 20px;
}

.sort-btn.active {
  background-color: var(--main-color);
  color: white;
}

.price-dropdown {
  background: white;
  padding: 7px 15px;
  border-radius: 2px;
  font-size: 14px;
  cursor: pointer;
  min-width: 150px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  line-height: 20px;
}

.pagination-summary {
  display: flex;
  align-items: center;
  gap: 15px;
  font-size: 14px;
}

.page-btns {
  display: flex;
  gap: 1px;
}

.page-btn {
  width: 36px;
  height: 34px;
  background: white;
  border: 1px solid #e5e5e5;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}

.page-btn.disabled {
  background: #f9f9f9;
  color: #ccc;
  cursor: not-allowed;
}

/* Grid */
.product-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 10px;
}

@media (max-width: 1024px) {
  .product-grid {
    grid-template-columns: repeat(4, 1fr);
  }
}

@media (max-width: 768px) {
  .product-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
