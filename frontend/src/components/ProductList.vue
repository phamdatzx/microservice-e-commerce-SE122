<script setup lang="ts">
import ProductItem from '@/components/ProductItem.vue'
import { ArrowDownBold, ArrowLeft, ArrowRight, Loading } from '@element-plus/icons-vue'

interface Product {
  id: number | string
  name: string
  price: number
  imageUrl: string
  rating: number
  location: string
  discount: number
  soldCount?: number
}

interface SortOption {
  label: string
  value: string
  active?: boolean
}

defineProps<{
  products: Product[]
  sortOptions?: SortOption[]
  pagination?: {
    current: number
    total: number
    pageSize: number
  }
  loading?: boolean
}>()

defineEmits(['sort-change', 'page-change'])
</script>

<template>
  <div class="product-list-container">
    <div class="container">
      <div class="main-content">
        <!-- Sidebar Slot -->
        <aside class="sidebar">
          <slot name="sidebar"></slot>
        </aside>

        <!-- Main Area -->
        <main class="product-grid-area">
          <!-- Top Content Slot (e.g., search summary) -->
          <div v-if="$slots.top" class="product-list-top">
            <slot name="top"></slot>
          </div>

          <!-- Sort Bar -->
          <div class="sort-bar">
            <div class="sort-options">
              <span class="sort-label">Sort by:</span>
              <slot name="sort-options">
                <!-- Default Sort Options -->
                <button
                  v-for="opt in sortOptions"
                  :key="opt.value"
                  class="sort-btn"
                  :class="{ active: opt.active }"
                  @click="$emit('sort-change', opt.value)"
                >
                  {{ opt.label }}
                </button>

                <!-- Price Dropdown -->
                <el-dropdown trigger="click" @command="(val: string) => $emit('sort-change', val)">
                  <div class="price-dropdown">
                    <span style="font-family: Arial, Helvetica, sans-serif">Price</span>
                    <el-icon><ArrowDownBold /></el-icon>
                  </div>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item command="price-asc">Price: Low to High</el-dropdown-item>
                      <el-dropdown-item command="price-desc">Price: High to Low</el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </slot>
            </div>

            <div v-if="pagination" class="pagination-summary">
              <span class="page-current">{{ pagination.current }}</span
              >/<span class="page-total">{{
                Math.ceil(pagination.total / pagination.pageSize) || 1
              }}</span>
              <div class="page-btns">
                <button
                  class="page-btn"
                  :class="{ disabled: pagination.current === 1 }"
                  @click="$emit('page-change', pagination.current - 1)"
                >
                  <el-icon><ArrowLeft /></el-icon>
                </button>
                <button
                  class="page-btn"
                  :class="{
                    disabled:
                      pagination.current >= Math.ceil(pagination.total / pagination.pageSize),
                  }"
                  @click="$emit('page-change', pagination.current + 1)"
                >
                  <el-icon><ArrowRight /></el-icon>
                </button>
              </div>
            </div>
          </div>

          <!-- Loading State -->
          <div v-if="loading" class="loading-state">
            <el-icon class="is-loading"><Loading /></el-icon>
            <span>Loading products...</span>
          </div>

          <!-- Grid -->
          <div v-else-if="products.length > 0" class="product-grid">
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
              :soldCount="product.soldCount"
            />
          </div>

          <!-- Empty State -->
          <div v-else class="empty-state">
            <div class="empty-icon-wrapper">
              <svg
                width="64"
                height="64"
                viewBox="0 0 24 24"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  d="M21 21L15 15M17 10C17 13.866 13.866 17 10 17C6.13401 17 3 13.866 3 10C3 6.13401 6.13401 3 10 3C13.866 3 17 6.13401 17 10Z"
                  stroke="#e4e4e7"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
              </svg>
            </div>
            <p class="empty-text">No products found</p>
            <p class="empty-subtext">
              Try adjusting your search or filters to find what you're looking for.
            </p>
          </div>

          <!-- Bottom Pagination -->
          <div v-if="pagination && products.length > 0" class="pagination-wrapper">
            <el-pagination
              background
              layout="prev, pager, next"
              :current-page="pagination.current"
              :total="pagination.total"
              :page-size="pagination.pageSize"
              size="large"
              class="custom-pagination"
              @current-change="(val: number) => $emit('page-change', val)"
            />
          </div>
        </main>
      </div>
      <!-- Close for .main-content -->
      <slot name="bottom"></slot>
    </div>
    <!-- Close for .container -->
  </div>
  <!-- Close for .product-list-container -->
</template>

<style scoped>
.product-list-container {
  padding-top: 20px;
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

.product-list-top {
  margin-bottom: 20px;
}

/* Sidebar */
.sidebar {
  width: 190px;
  flex-shrink: 0;
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
  color: #333;
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

@media (max-width: 1200px) {
  .product-grid {
    grid-template-columns: repeat(4, 1fr);
  }
}

@media (max-width: 992px) {
  .product-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (max-width: 768px) {
  .product-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  padding: 25px 0 30px;
}

:deep(.custom-pagination .el-pager li.is-active) {
  background-color: var(--main-color) !important;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 0;
  text-align: center;
}

.empty-icon-wrapper {
  margin-bottom: 20px;
}

.empty-text {
  font-size: 18px;
  font-weight: 600;
  color: #3f3f46;
  margin-bottom: 8px;
}

.empty-subtext {
  font-size: 14px;
  color: #71717a;
  max-width: 400px;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 0;
  color: #71717a;
  gap: 15px;
}

.loading-state .el-icon {
  font-size: 24px;
}
</style>
