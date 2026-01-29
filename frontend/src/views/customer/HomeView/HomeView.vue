<script setup lang="ts">
import ProductItem from '@/components/ProductItem.vue'
import RecentlyViewed from '@/components/RecentlyViewed.vue'

import { Splide, SplideSlide } from '@splidejs/vue-splide'
import '@splidejs/vue-splide/css'
import { ref, onMounted } from 'vue'
import { ShoppingCart } from '@element-plus/icons-vue'
import CategoryItem from './components/CategoryItem.vue'
import SellerFollowItem from './components/SellerFollowItem.vue'
import axios from 'axios'
import { eventBus } from '@/utils/eventBus'

const fetchedCategories = ref<any[]>([])
const bestSellers = ref<any[]>([])

const loadingCategories = ref(false)
const loadingBestSellers = ref(false)
const loadingFollowedSellers = ref(false)
const isLoggedIn = ref(false)

const fetchBestSellers = async () => {
  loadingBestSellers.value = true
  try {
    const response = await axios.get(
      `${import.meta.env.VITE_BE_API_URL}/product/public/search?page=1&limit=15&sort_by=sold_count&sort_direction=desc`,
    )
    if (response.data && response.data.products) {
      bestSellers.value = response.data.products.map((p: any) => ({
        id: p.id,
        name: p.name,
        imageUrl:
          p.images.length > 0
            ? p.images.sort((a: any, b: any) => a.order - b.order)[0]?.url || ''
            : '',
        minPrice: p.price.min,
        maxPrice: p.price.max,
        rating: p.rating,
        location: 'Vietnam',
        soldCount: p.sold_count,
      }))
    }
  } catch (error) {
    console.error('Error fetching best sellers:', error)
  } finally {
    loadingBestSellers.value = false
  }
}

const fetchCategories = async () => {
  loadingCategories.value = true
  try {
    const response = await axios.get(`${import.meta.env.VITE_BE_API_URL}/product/public/category`)
    fetchedCategories.value = response.data
  } catch (error) {
    console.error('Error fetching categories:', error)
  } finally {
    loadingCategories.value = false
  }
}

const carouselItems = [
  {
    imageUrl: new URL('@/assets/carousel-imgs/banner1.jpg', import.meta.url).href,
  },
  {
    imageUrl: new URL('@/assets/carousel-imgs/banner2.jpg', import.meta.url).href,
  },
  {
    imageUrl: new URL('@/assets/carousel-imgs/banner3.jpg', import.meta.url).href,
  },
  {
    imageUrl: new URL('@/assets/carousel-imgs/banner4.jpg', import.meta.url).href,
  },
]

const followedSellers = ref<any[]>([])

const fetchFollowedSellers = async () => {
  const token = localStorage.getItem('access_token')
  if (!token) return

  loadingFollowedSellers.value = true
  try {
    const response = await axios.get(`${import.meta.env.VITE_BE_API_URL}/user/follow`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
    if (response.data && response.data.data) {
      followedSellers.value = response.data.data
    }
  } catch (error) {
    console.error('Error fetching followed sellers:', error)
  } finally {
    loadingFollowedSellers.value = false
  }
}

onMounted(() => {
  isLoggedIn.value = !!localStorage.getItem('access_token')
  fetchCategories()
  fetchFollowedSellers()
  fetchBestSellers()

  eventBus.on('user_logged_out', () => {
    isLoggedIn.value = false
  })
})
</script>

<template>
  <main class="main-container">
    <el-row :gutter="20" style="padding: 20px 0; height: 440px">
      <el-col :span="6">
        <div
          class="box-shadow border-radius"
          style="width: 100%; height: 100%; background-color: #fff; padding: 20px 40px"
          v-loading="loadingCategories"
        >
          <h4 style="color: var(--main-color); font-weight: bold; margin-bottom: 4px">
            CATEGORIES
          </h4>
          <ul v-if="fetchedCategories.length > 0">
            <li v-for="category in fetchedCategories" :key="category.id">
              <RouterLink
                class="category-item"
                :to="{ path: '/search', query: { category_ids: category.id } }"
                >{{ category.name }}</RouterLink
              >
            </li>
          </ul>
        </div>
      </el-col>
      <el-col :span="18">
        <el-carousel
          class="box-shadow border-radius"
          height="100%"
          style="width: 100%; height: 100%; overflow: hidden; position: relative"
        >
          <el-carousel-item v-for="item in carouselItems" :key="item">
            <div style="width: 100%; height: 100%; font-size: 0px">
              <img
                :src="item.imageUrl"
                style="width: 100%; height: 100%; object-fit: cover; display: block"
              />
            </div>
          </el-carousel-item>
        </el-carousel>
        <div></div>
      </el-col>
    </el-row>

    <!-- Followed Sellers Section -->
    <div
      v-if="followedSellers.length > 0"
      class="box-shadow border-radius followed-sellers"
      style="background-color: #fff; padding: 20px 20px 32px; margin-bottom: 20px"
      v-loading="loadingFollowedSellers"
    >
      <div style="display: flex; margin-bottom: 24px">
        <h3 style="font-weight: bold">FOLLOWED SELLERS</h3>
      </div>

      <Splide
        :options="{
          rewind: true,
          perPage: 8,
          gap: '2rem',
          pagination: false,
          arrows: followedSellers.length > 8,
          breakpoints: {
            1200: { perPage: 6 },
            992: { perPage: 4 },
            768: { perPage: 3 },
            576: { perPage: 2 },
          },
        }"
      >
        <SplideSlide v-for="seller in followedSellers" :key="seller.id">
          <SellerFollowItem
            :image-url="seller.image"
            :name="seller.full_name || seller.username"
            :id="seller.user_id"
          />
        </SplideSlide>
      </Splide>
    </div>

    <div
      class="box-shadow border-radius top-categories"
      style="background-color: #fff; padding: 20px 20px 40px; margin-bottom: 20px"
      v-loading="loadingCategories"
    >
      <div style="display: flex; margin-bottom: 20px">
        <h3 style="font-weight: bold">TOP CATEGORIES</h3>
      </div>

      <Splide
        v-if="fetchedCategories.length > 0"
        :options="{
          rewind: true,
          perPage: 6,
          gap: '1rem',
          padding: '40px',
          breakpoints: {
            1000: {
              perPage: 1,
            },
          },
        }"
      >
        <SplideSlide v-for="item in fetchedCategories" :key="item.id">
          <RouterLink
            :to="{ path: '/search', query: { category_ids: item.id } }"
            style="text-decoration: none; color: inherit"
          >
            <CategoryItem :image-url="item.image" :name="item.name" />
          </RouterLink>
        </SplideSlide>
      </Splide>
    </div>

    <div
      class="box-shadow border-radius"
      style="background-color: #fff; padding: 20px; margin-bottom: 20px"
      v-loading="loadingBestSellers"
    >
      <div style="display: flex; margin-bottom: 24px">
        <h3 style="font-weight: bold">BEST SELLER</h3>
      </div>

      <el-row :gutter="20">
        <el-col
          :span="4.8"
          class="el-col-4-8"
          v-for="(item, index) in bestSellers"
          :key="index"
          style="margin-bottom: 20px"
        >
          <ProductItem
            :image-url="item.imageUrl"
            :name="item.name"
            :min-price="item.minPrice"
            :max-price="item.maxPrice"
            :rating="item.rating"
            :location="item.location"
            :sold-count="item.soldCount"
            :id="item.id"
          />
        </el-col>
      </el-row>
    </div>

    <RecentlyViewed v-if="isLoggedIn" />
  </main>
</template>

<style lang="css" scoped>
.category-item {
  transition: color 0.2s;
  font-weight: 500;
  display: block;
  padding: 4px 0;

  &:hover {
    color: #999;
  }
}
</style>

<style>
.top-categories .splide__pagination__page {
  position: relative !important;
  top: 34px;
}

.recently-viewed .splide__pagination__page {
  position: relative !important;
  top: 16px;
}

.splide__pagination__page.is-active {
  background-color: var(--main-color);
  transform: scale(1.2);
}
</style>
