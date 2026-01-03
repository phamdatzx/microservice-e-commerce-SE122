<script setup lang="ts">
import ProductItem from '@/components/ProductItem.vue'

const recommendedProducts = [
  {
    id: 101,
    name: 'Jeans for Men Skinny Fit White Plain Stretchable Office Stylish Fashion',
    price: 149600,
    imageUrl: '/src/assets/product-imgs/product1.png',
    rating: 4.8,
    location: 'Ho Chi Minh City',
    discount: 28,
  },
  {
    id: 102,
    name: 'Cargo Jeans for Men Wide Leg Relaxed Fit Loose Fashion Trendy',
    price: 175000,
    imageUrl: '/src/assets/product-imgs/product1.png',
    rating: 4.8,
    location: 'Ha Noi',
    discount: 32,
  },
  {
    id: 103,
    name: 'Skinny Jeans Men Black Washed Casual Comfort Fit Elastic',
    price: 125000,
    imageUrl: '/src/assets/product-imgs/product1.png',
    rating: 4.8,
    location: 'Da Nang',
    discount: 37,
  },
  {
    id: 104,
    name: 'WASH RETRO Jeans Men Dark Grey Loose Fit Vintage Style',
    price: 174999,
    imageUrl: '/src/assets/product-imgs/product1.png',
    rating: 4.8,
    location: 'Ho Chi Minh City',
    discount: 24,
  },
  {
    id: 105,
    name: 'Elastic Cuff Pants for Men Casual Streetwear Jogger Style',
    price: 92900,
    imageUrl: '/src/assets/product-imgs/product1.png',
    rating: 4.9,
    location: 'Binh Duong',
    discount: 45,
  },
  {
    id: 106,
    name: 'Skinny Jeans Men Black Washed Casual Comfort Fit Elastic',
    price: 125000,
    imageUrl: '/src/assets/product-imgs/product1.png',
    rating: 4.8,
    location: 'Da Nang',
    discount: 37,
  },
]
</script>

<template>
  <div class="recommendations">
    <div class="section-header">
      <h2 class="section-title">RECOMMENDED FOR YOU</h2>
      <a href="#" class="view-all">View All</a>
    </div>

    <div class="product-row">
      <div v-for="product in recommendedProducts" :key="product.id" class="product-col">
        <ProductItem
          :id="product.id"
          :name="product.name"
          :imageUrl="product.imageUrl"
          :price="product.price"
          :rating="product.rating"
          :location="product.location"
          :discount="product.discount"
        />
      </div>
    </div>
  </div>
</template>

<style scoped>
.recommendations {
  margin-top: 20px;
  margin-bottom: 30px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 5px;
}

.section-title {
  font-size: 16px;
  font-weight: 500;
  color: #757575;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.view-all {
  font-size: 14px;
  color: var(--main-color);
  text-decoration: none;
}

.product-row {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 10px;
}

@media (max-width: 1200px) {
  .product-row {
    grid-template-columns: repeat(4, 1fr);
  }
}

@media (max-width: 768px) {
  .product-row {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
