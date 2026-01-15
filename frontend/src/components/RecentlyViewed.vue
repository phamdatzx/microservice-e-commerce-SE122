<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Splide, SplideSlide } from '@splidejs/vue-splide'
import { ShoppingCart } from '@element-plus/icons-vue'
import ProductItem from './ProductItem.vue'

const recentlyViewedProducts = ref<any[]>([])

const loadRecentlyViewed = () => {
  const stored = localStorage.getItem('recently_viewed')
  if (stored) {
    recentlyViewedProducts.value = JSON.parse(stored)
  }
}

onMounted(() => {
  loadRecentlyViewed()
})

defineExpose({
  loadRecentlyViewed,
})
</script>

<template>
  <div
    class="box-shadow border-radius recently-viewed-section"
    style="background-color: #fff; padding: 20px 20px 28px; margin-bottom: 20px"
  >
    <div style="display: flex; margin-bottom: 12px">
      <h3 style="font-weight: bold">YOU RECENTLY VIEWED</h3>
      <RouterLink to="/recently-viewed" class="view-all-btn"> View All </RouterLink>
    </div>

    <div v-if="recentlyViewedProducts.length">
      <Splide
        :options="{
          rewind: true,
          perPage: 5,
          gap: '1rem',
          breakpoints: {
            1024: { perPage: 4 },
            768: { perPage: 2 },
            480: { perPage: 1 },
          },
        }"
      >
        <SplideSlide v-for="item in recentlyViewedProducts" :key="item.id">
          <ProductItem
            :image-url="item.imageUrl"
            :name="item.name"
            :price="item.price"
            :rating="item.rating"
            :location="item.location"
            :discount="item.discount"
            :sold-count="item.soldCount"
            :id="item.id"
            style="padding-bottom: 12px"
          />
        </SplideSlide>
      </Splide>
    </div>

    <div v-else style="text-align: center; color: #999; padding: 40px">
      <el-icon size="40" style="margin-bottom: 12px; opacity: 0.5"><ShoppingCart /></el-icon>
      <p style="font-size: 14px">Your recently viewed list is empty.</p>
      <p style="font-size: 13px; color: #bbb">
        Products you visit will appear here for easy access.
      </p>
    </div>
  </div>
</template>

<style scoped>
.recently-viewed-section {
  box-shadow: rgba(0, 0, 0, 0.05) 0px 2px 10px;
  border-radius: 12px;
}

h3 {
  font-size: 18px;
  color: #1a1a1a;
}

:deep(.splide__track) {
  padding: 20px 20px !important;
  margin: -20px -20px;
}

:deep(.splide__slide) {
  padding: 4px;
}

.view-all-btn {
  position: relative;
  top: 4px;
  margin-left: 20px;
  color: #999;
  font-size: 14px;
  text-decoration: none;
  transition: color 0.2s;
}

.view-all-btn:hover {
  color: var(--main-color);
}
</style>

<style>
.recently-viewed-section .splide__pagination__page {
  top: 8px;
}
</style>
