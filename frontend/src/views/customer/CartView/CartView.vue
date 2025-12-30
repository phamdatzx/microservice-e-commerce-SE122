<script setup lang="ts">
import Header from '@/components/Header.vue'
import type { CheckboxValueType } from 'element-plus'
import { ref, useTemplateRef, watch } from 'vue'
import CartSellerProductWrapper from './components/CartSellerProductWrapper.vue'
import { Ticket } from '@element-plus/icons-vue'

import { Splide, SplideSlide } from '@splidejs/vue-splide'

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
  sellerName: string
}

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

const checkAll = ref(false)
const checkedProducts = ref<Product[]>([])
const CartSellerProductWrapperRefs = ref<InstanceType<typeof CartSellerProductWrapper>[]>([])
const cartData: CartSeller[] = [
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
        sellerName: 'LOVITO OFFICIAL STORE',
      },
      {
        id: 2,
        imageUrl: '/src/assets/product-imgs/product1.png',
        productName:
          'Lovito Áo Nỉ Thường Ngày Trơn Có Khóa Kéo Mùa Xuân Và Mùa Hè Dành Cho Nữ L102AD477',
        productOption: 'Gray, XL',
        price: 23.99,
        quantity: 2,
        sellerName: 'LOVITO OFFICIAL STORE',
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
        sellerName: 'LOVITO OFFICIAL STORE2',
      },
      {
        id: 2,
        imageUrl: '/src/assets/product-imgs/product1.png',
        productName:
          'Lovito Áo Nỉ Thường Ngày Trơn Có Khóa Kéo Mùa Xuân Và Mùa Hè Dành Cho Nữ L102AD477',
        productOption: 'Gray, XL',
        price: 23.99,
        quantity: 2,
        sellerName: 'LOVITO OFFICIAL STORE2',
      },
    ],
  },
]
const allProducts = cartData.map((seller) => seller.products).flat()

const handleCheckAll = () => {
  CartSellerProductWrapperRefs.value.forEach((ref) => {
    ref.handleCheckAll()
  })
}

const handleUncheckAll = () => {
  CartSellerProductWrapperRefs.value.forEach((ref) => {
    ref.handleUncheckAll()
  })
}

const handleCheckAllChange = (val: CheckboxValueType) => {
  checkedProducts.value = val ? allProducts : []

  if (val) {
    handleCheckAll()
  } else {
    handleUncheckAll()
  }
}

const handleCheckedProductsChange = (checkedSellerProducts: Product[], sellerName: string) => {
  checkedProducts.value = checkedProducts.value.filter(
    (product) => product.sellerName !== sellerName,
  )

  if (checkedSellerProducts.length > 0) {
    checkedProducts.value.push(...checkedSellerProducts)
  }

  checkAll.value = checkedProducts.value.length === allProducts.length
}
</script>

<template>
  <Header />

  <div class="main-container">
    <div
      class="cart-box-shadow border-radius"
      style="
        padding: 0 20px;
        margin: 20px 0;
        display: flex;
        align-items: center;
        margin-bottom: 12px;
        height: 64px;
        background-color: #fff;
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

    <CartSellerProductWrapper
      v-for="seller in cartData"
      :key="seller.sellerName"
      ref="CartSellerProductWrapperRefs"
      :seller="seller"
      @checked-products-change="handleCheckedProductsChange"
    />

    <div
      class="border-radius"
      style="
        background-color: #fff;
        border: 1px solid #ccc;
        padding: 16px;
        position: sticky;
        bottom: 0;
        z-index: 1000;
      "
    >
      <div style="display: flex; padding-bottom: 12px">
        <button class="btn">Select All ({{ allProducts.length }})</button>
        <button class="btn">Delete</button>
        <button class="btn">Remove Unactive Products</button>
        <button class="btn">Add To Favorites</button>
      </div>
      <div style="display: flex; justify-content: space-between; align-items: center">
        <el-button
          size="large"
          style="font-size: 18px; color: var(--main-color); padding: 28px 20px 24px"
        >
          <el-icon style="margin-right: 8px; font-size: 24px; position: relative; top: -2px"
            ><Ticket
          /></el-icon>
          Voucher: Apply Voucher
        </el-button>
        <div style="display: flex; align-items: center">
          <span style="font-size: 20px; margin-right: 12px"
            >Total ({{ checkedProducts.length }} products):
            <span style="font-size: 24px; color: var(--main-color)">
              {{
                checkedProducts.reduce((sum, product) => {
                  return sum + product.price * (product.quantity ?? 1)
                }, 0)
              }}$
            </span>
          </span>
          <el-button style="font-size: 24px; padding: 32px 40px; color: var(--main-color)"
            >Checkout</el-button
          >
        </div>
      </div>
    </div>

    <div
      class="box-shadow border-radius recently-viewed"
      style="background-color: #fff; padding: 20px; margin-bottom: 20px; margin-top: 28px"
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
  </div>
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

.btn {
  background-color: #fff;
  border: none;
  padding: 12px 20px;
  font-size: 18px;
  border: 1px solid #ccc;
  margin-right: 8px;
  transition:
    color 0.2s,
    border-color 0.2s;

  &:hover {
    color: #22c55e;
    border-color: currentColor;
  }
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

.cart-box-shadow {
  border: 1px solid #ccc;
}
</style>
