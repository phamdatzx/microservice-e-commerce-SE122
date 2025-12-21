<script setup lang="ts">
import Header from '@/components/Header.vue'
import type { CheckboxValueType } from 'element-plus'
import { ref, watch } from 'vue'
import CartSellerProductWrapper from './components/CartSellerProductWrapper.vue'

const checkAll = ref(false)
const checkedSellers = ref<CartSeller[]>([])

export interface CartSeller {
  sellerName: string
  products: Product[]
}

export interface Product {
  id: number
  imageUrl: string
  productName: string
  productOption: string
  price: number
  quantity?: number
  isSoldOut?: boolean
}

const cartSellers: CartSeller[] = [
  {
    sellerName: 'LOVITO OFFICIAL STORE',
    products: [
      {
        id: 1,
        imageUrl: '/src/assets/product-imgs/product1.png',
        productName:
          'Lovito Áo Nỉ Thường Ngày Trơn Có Khóa Kéo Mùa Xuân Và Mùa Hè Dành Cho Nữ L102AD477',
        productOption: 'Gray, XL',
        price: 23.99,
        quantity: 1,
      },
      {
        id: 2,
        imageUrl: '/src/assets/product-imgs/product1.png',
        productName:
          'Lovito Áo Nỉ Thường Ngày Trơn Có Khóa Kéo Mùa Xuân Và Mùa Hè Dành Cho Nữ L102AD477',
        productOption: 'Gray, XL',
        price: 23.99,
        quantity: 2,
      },
    ],
  },
  {
    sellerName: 'LOVITO OFFICIAL STORE2',
    products: [
      {
        id: 1,
        imageUrl: '/src/assets/product-imgs/product1.png',
        productName:
          'Lovito Áo Nỉ Thường Ngày Trơn Có Khóa Kéo Mùa Xuân Và Mùa Hè Dành Cho Nữ L102AD477',
        productOption: 'Gray, XL',
        price: 23.99,
        quantity: 1,
      },
      {
        id: 2,
        imageUrl: '/src/assets/product-imgs/product1.png',
        productName:
          'Lovito Áo Nỉ Thường Ngày Trơn Có Khóa Kéo Mùa Xuân Và Mùa Hè Dành Cho Nữ L102AD477',
        productOption: 'Gray, XL',
        price: 23.99,
        quantity: 2,
      },
    ],
  },
]

const handleCheckAllChange = (val: CheckboxValueType) => {
  checkedSellers.value = val ? cartSellers : []
}

const handleCheckedSellersChange = (value: CheckboxValueType[]) => {
  const checkedCount = value.length
  checkAll.value = checkedCount === cartSellers.length
}
</script>

<template>
  <Header />

  <div class="main-container">
    <div
      class="box-shadow border-radius"
      style="
        padding: 0 20px;
        margin: 20px 0;
        display: flex;
        align-items: center;
        margin-bottom: 12px;
        height: 64px;
      "
    >
      <el-checkbox
        v-model="checkAll"
        @change="handleCheckAllChange"
        size="large"
        class="cartCheckbox"
      />
      <h3 style="flex: 1">Product</h3>
      <h3 style="width: 174px; text-align: center">Price</h3>
      <h3 style="width: 169px; text-align: center">Quantity</h3>
      <h3 style="width: 114px; text-align: center">Total</h3>
      <h3 style="width: 139px; text-align: center">Operation</h3>
    </div>
  </div>

  <el-checkbox-group v-model="checkedSellers" @change="handleCheckedSellersChange">
    <div class="main-container" v-for="seller in cartSellers" :key="seller.sellerName">
      <CartSellerProductWrapper :seller="seller" />
    </div>
  </el-checkbox-group>
</template>

<style lang="css" scoped>
.product-name {
  --line-height: 16px;
  height: calc(var(--line-height) * 2);
  line-height: var(--line-height);

  color: #000000cc;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  line-clamp: 2;
  -webkit-line-clamp: 2;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>

<style>
.el-checkbox-group {
  font-size: 14px;
  line-height: 1.6;
}

.cartCheckbox {
  padding: 0 12px 0 20px;
  margin-right: 8px;
  transform: scale(1.4);
}
</style>
