<script setup lang="ts">
import CarouselItem from '@/components/CarouselItem.vue'
import CategoryItem from '@/components/CategoryItem.vue'
import Header from '@/components/Header.vue'
import ProductItem from '@/components/ProductItem.vue'

import { Splide, SplideSlide } from '@splidejs/vue-splide'
import '@splidejs/vue-splide/css'
import { ref } from 'vue'

const categories = [
  'Laptops',
  'PC & Computers',
  'Cell Phones',
  'Tablets',
  'Gaming & VR',
  'Networking',
  'Cameras',
  'Sounds',
  'Office',
  'Storage, USB',
  'Accessories',
  'Clearance',
]

const carouselItems = [
  {
    imageUrl: '/src/assets/carousel-imgs/headphones.png',
    to: '/',
  },
  {
    imageUrl: '/src/assets/carousel-imgs/headphones.png',
    to: '/',
  },
  {
    imageUrl: '/src/assets/carousel-imgs/headphones.png',
    to: '/',
  },
  {
    imageUrl: '/src/assets/carousel-imgs/headphones.png',
    to: '/',
  },
]

const topCategoriesList = [
  {
    name: 'Laptops',
    imageUrl: '/src/assets/category-imgs/laptop.png',
  },
  {
    name: 'Laptops',
    imageUrl: '/src/assets/category-imgs/laptop.png',
  },
  {
    name: 'Laptops',
    imageUrl: '/src/assets/category-imgs/laptop.png',
  },
  {
    name: 'Laptops',
    imageUrl: '/src/assets/category-imgs/laptop.png',
  },
  {
    name: 'Laptops',
    imageUrl: '/src/assets/category-imgs/laptop.png',
  },
  {
    name: 'Laptops',
    imageUrl: '/src/assets/category-imgs/laptop.png',
  },
  {
    name: 'Laptops',
    imageUrl: '/src/assets/category-imgs/laptop.png',
  },
  {
    name: 'Laptops',
    imageUrl: '/src/assets/category-imgs/laptop.png',
  },
  {
    name: 'Laptops',
    imageUrl: '/src/assets/category-imgs/laptop.png',
  },
  {
    name: 'Laptops',
    imageUrl: '/src/assets/category-imgs/laptop.png',
  },
  {
    name: 'Laptops',
    imageUrl: '/src/assets/category-imgs/laptop.png',
  },
  {
    name: 'Laptops',
    imageUrl: '/src/assets/category-imgs/laptop.png',
  },
  {
    name: 'Laptops',
    imageUrl: '/src/assets/category-imgs/laptop.png',
  },
]

const productList = [
  {
    imageUrl: '/src/assets/product-imgs/product1.png',
    name: 'Laptop Dell 15 DC15250 71071928 (Core i5-1334U/16GB/512GB/Intel Graphics/15.6 inch FHD IPS/Win 11/Bạc) - Chính hãng',
    price: 16090000,
    rating: 5.0,
    location: 'Ha Noi',
    discount: 37,
  },
  {
    imageUrl: '/src/assets/product-imgs/product1.png',
    name: 'Laptop Dell 15 DC15250 71071928 (Core i5-1334U/16GB/512GB/Intel Graphics/15.6 inch FHD IPS/Win 11/Bạc) - Chính hãng',
    price: 16090000,
    rating: 5.0,
    location: 'Ha Noi',
    discount: 37,
  },
  {
    imageUrl: '/src/assets/product-imgs/product1.png',
    name: 'Laptop Dell 15 DC15250 71071928 (Core i5-1334U/16GB/512GB/Intel Graphics/15.6 inch FHD IPS/Win 11/Bạc) - Chính hãng',
    price: 16090000,
    rating: 5.0,
    location: 'Ha Noi',
    discount: 37,
  },
  {
    imageUrl: '/src/assets/product-imgs/product1.png',
    name: 'Laptop Dell 15 DC15250 71071928 (Core i5-1334U/16GB/512GB/Intel Graphics/15.6 inch FHD IPS/Win 11/Bạc) - Chính hãng',
    price: 16090000,
    rating: 5.0,
    location: 'Ha Noi',
    discount: 37,
  },
  {
    imageUrl: '/src/assets/product-imgs/product1.png',
    name: 'Laptop Dell 15 DC15250 71071928 (Core i5-1334U/16GB/512GB/Intel Graphics/15.6 inch FHD IPS/Win 11/Bạc) - Chính hãng',
    price: 16090000,
    rating: 5.0,
    location: 'Ha Noi',
    discount: 37,
  },
  {
    imageUrl: '/src/assets/product-imgs/product1.png',
    name: 'Laptop Dell 15 DC15250 71071928 (Core i5-1334U/16GB/512GB/Intel Graphics/15.6 inch FHD IPS/Win 11/Bạc) - Chính hãng',
    price: 16090000,
    rating: 5.0,
    location: 'Ha Noi',
    discount: 37,
  },
  {
    imageUrl: '/src/assets/product-imgs/product1.png',
    name: 'Laptop Dell 15 DC15250 71071928 (Core i5-1334U/16GB/512GB/Intel Graphics/15.6 inch FHD IPS/Win 11/Bạc) - Chính hãng',
    price: 16090000,
    rating: 5.0,
    location: 'Ha Noi',
    discount: 37,
  },
  {
    imageUrl: '/src/assets/product-imgs/product1.png',
    name: 'Laptop Dell 15 DC15250 71071928 (Core i5-1334U/16GB/512GB/Intel Graphics/15.6 inch FHD IPS/Win 11/Bạc) - Chính hãng',
    price: 16090000,
    rating: 5.0,
    location: 'Ha Noi',
    discount: 37,
  },
  {
    imageUrl: '/src/assets/product-imgs/product1.png',
    name: 'Laptop Dell 15 DC15250 71071928 (Core i5-1334U/16GB/512GB/Intel Graphics/15.6 inch FHD IPS/Win 11/Bạc) - Chính hãng',
    price: 16090000,
    rating: 5.0,
    location: 'Ha Noi',
    discount: 37,
  },
  {
    imageUrl: '/src/assets/product-imgs/product1.png',
    name: 'Laptop Dell 15 DC15250 71071928 (Core i5-1334U/16GB/512GB/Intel Graphics/15.6 inch FHD IPS/Win 11/Bạc) - Chính hãng',
    price: 16090000,
    rating: 5.0,
    location: 'Ha Noi',
    discount: 37,
  },
  {
    imageUrl: '/src/assets/product-imgs/product1.png',
    name: 'Laptop Dell 15 DC15250 71071928 (Core i5-1334U/16GB/512GB/Intel Graphics/15.6 inch FHD IPS/Win 11/Bạc) - Chính hãng',
    price: 16090000,
    rating: 5.0,
    location: 'Ha Noi',
    discount: 37,
  },
  {
    imageUrl: '/src/assets/product-imgs/product1.png',
    name: 'Laptop Dell 15 DC15250 71071928 (Core i5-1334U/16GB/512GB/Intel Graphics/15.6 inch FHD IPS/Win 11/Bạc) - Chính hãng',
    price: 16090000,
    rating: 5.0,
    location: 'Ha Noi',
    discount: 37,
  },
  {
    imageUrl: '/src/assets/product-imgs/product1.png',
    name: 'Laptop Dell 15 DC15250 71071928 (Core i5-1334U/16GB/512GB/Intel Graphics/15.6 inch FHD IPS/Win 11/Bạc) - Chính hãng',
    price: 16090000,
    rating: 5.0,
    location: 'Ha Noi',
    discount: 37,
  },
]

const activeTab = ref('best-seller')
</script>

<template>
  <Header />

  <main class="main-container">
    <el-row :gutter="20" style="padding: 20px 0">
      <el-col :span="6">
        <div
          class="box-shadow border-radius"
          style="width: 100%; height: 100%; background-color: #fff; padding: 20px 40px"
        >
          <h4 style="color: red; font-weight: bold; margin-bottom: 4px">SALE 40% OFF</h4>
          <ul>
            <li v-for="category in categories" :key="category">
              <RouterLink class="category-item" to="/">{{ category }}</RouterLink>
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
            <CarouselItem :imageUrl="item.imageUrl" :to="item.to" />
          </el-carousel-item>
        </el-carousel>
        <div></div>
      </el-col>
    </el-row>

    <div
      class="box-shadow border-radius top-categories"
      style="background-color: #fff; padding: 20px 20px 32px; margin-bottom: 20px"
    >
      <div style="display: flex; margin-bottom: 12px">
        <h3 style="font-weight: bold">TOP CATEGORIES</h3>
        <RouterLink
          to="/"
          style="position: relative; top: 3px; margin-left: 20px; color: #999; font-size: 13px"
          >View All</RouterLink
        >
      </div>

      <Splide
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
        <SplideSlide v-for="item in topCategoriesList" :key="item.name">
          <CategoryItem :image-url="item.imageUrl" :name="item.name" />
        </SplideSlide>
      </Splide>
    </div>

    <div
      class="box-shadow border-radius"
      style="background-color: #fff; padding: 20px 20px 0; margin-bottom: 20px"
    >
      <el-tabs v-model="activeTab">
        <el-tab-pane label="BEST SELLER" name="best-seller">
          <el-row :gutter="20" style="padding-bottom: 20px">
            <el-col
              :span="4.8"
              class="el-col-4-8"
              v-for="(item, index) in productList"
              :key="index"
            >
              <ProductItem
                :image-url="item.imageUrl"
                :name="item.name"
                :price="item.price"
                :rating="item.rating"
                :location="item.location"
                :discount="item.discount"
              />
            </el-col>
          </el-row>
        </el-tab-pane>
        <el-tab-pane label="NEW IN" name="new-in">Config</el-tab-pane>
        <el-tab-pane label="POPULAR" name="popular">Role</el-tab-pane>
      </el-tabs>
    </div>

    <div
      class="box-shadow border-radius recently-viewed"
      style="background-color: #fff; padding: 20px; margin-bottom: 20px"
    >
      <div style="display: flex; margin-bottom: 12px">
        <h3 style="font-weight: bold">YOU RECENTLY VIEWED</h3>
        <RouterLink
          to="/"
          style="position: relative; top: 3px; margin-left: 20px; color: #999; font-size: 13px"
          >View All</RouterLink
        >
      </div>

      <Splide
        :options="{
          rewind: true,
          perPage: 5,
          gap: '1rem',
          breakpoints: {
            1000: {
              perPage: 1,
            },
          },
        }"
      >
        <SplideSlide v-for="item in productList" :key="item.name">
          <ProductItem
            :image-url="item.imageUrl"
            :name="item.name"
            :price="item.price"
            :rating="item.rating"
            :location="item.location"
            :discount="item.discount"
          />
        </SplideSlide>
      </Splide>
    </div>
  </main>

  <div style="height: 1000px"></div>
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
  top: 30px;
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
