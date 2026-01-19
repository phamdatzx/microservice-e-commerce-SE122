<script setup lang="ts">
import { ref, onMounted } from 'vue'
import ProductItem from '@/components/ProductItem.vue'

interface Product {
  id: string
  name: string
  imageUrl: string
  price: number
  rating: number
  location: string
  discount: number
  soldCount: number
}

import axios from 'axios'

const recentlyViewed = ref<Product[]>([])

const loadRecentlyViewed = async () => {
  const token = localStorage.getItem('access_token')
  const headers = token ? { Authorization: `Bearer ${token}` } : {}
  try {
    const response = await axios.get(
      `${import.meta.env.VITE_BE_API_URL}/product/recently-viewed-products?limit=20`,
      { headers },
    )
    if (response.data && Array.isArray(response.data)) {
      recentlyViewed.value = response.data.map((p: any) => ({
        id: p.id,
        name: p.name,
        imageUrl:
          p.images.length > 0
            ? p.images.sort((a: any, b: any) => a.order - b.order)[0]?.url || ''
            : '',
        price: p.price.min,
        rating: p.rating,
        location: 'Vietnam',
        discount: 0,
        soldCount: p.sold_count,
      }))
    }
  } catch (error) {
    console.error('Error fetching recently viewed products:', error)
  }
}

onMounted(() => {
  loadRecentlyViewed()
})
</script>

<template>
  <div class="main-container">
    <div class="page-header">
      <div class="title-section">
        <h2 class="page-title">Recently Viewed</h2>
        <p class="subtitle">Products you've recently visited are saved here for easy access.</p>
      </div>
    </div>

    <div v-if="recentlyViewed.length > 0" class="products-grid">
      <el-row :gutter="20">
        <el-col
          v-for="item in recentlyViewed"
          :key="item.id"
          :xs="24"
          :sm="12"
          :md="6"
          :lg="4.8"
          class="el-col-4-8"
          style="margin-bottom: 20px"
        >
          <ProductItem
            :id="item.id"
            :name="item.name"
            :image-url="item.imageUrl"
            :price="item.price"
            :rating="item.rating"
            :location="item.location"
            :discount="item.discount"
            :sold-count="item.soldCount"
          />
        </el-col>
      </el-row>
    </div>

    <div v-else class="empty-state">
      <el-empty description="Your recently viewed history is empty">
        <el-button type="primary" color="var(--main-color)" @click="$router.push('/')">
          Go Shopping Now
        </el-button>
      </el-empty>
    </div>
  </div>
</template>

<style scoped>
.main-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 40px 20px 12px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  border-bottom: 1px solid #eee;
  padding-bottom: 16px;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  color: #1a1a1a;
  margin: 0 0 8px 0;
}

.subtitle {
  color: #666;
  font-size: 14px;
  margin: 0;
}

.empty-state {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
}

/* Custom class for 5 items per row if needed, since element-plus uses 24 grid */
@media (min-width: 1200px) {
  .el-col-4-8 {
    width: 20%;
    max-width: 20%;
    flex: 0 0 20%;
  }
}
</style>
