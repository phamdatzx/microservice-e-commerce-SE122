<script setup lang="ts">
import type { CheckboxValueType } from 'element-plus'
import { ref, useTemplateRef, watch, onMounted } from 'vue'
import { type CartSeller, type Product } from '../CartView.vue'
import { ChatDotRound, Shop } from '@element-plus/icons-vue'
import CartProduct from './CartProduct.vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const isLoggedIn = ref(false)

onMounted(() => {
  isLoggedIn.value = !!localStorage.getItem('access_token')
})

interface Props {
  seller: CartSeller
}
const props = withDefaults(defineProps<Props>(), {
  seller: () => ({
    sellerId: '',
    sellerName: '',
    products: [],
  }),
})

const emits = defineEmits(['checkedProductsChange', 'deleteProduct', 'updateQuantity'])

const checkAll = ref(false)
const checkedProducts = ref<Product[]>([])

watch([checkedProducts], () => {
  emits('checkedProductsChange', checkedProducts.value, props.seller.sellerName)
})

const handleCheckAllChange = (val: CheckboxValueType) => {
  checkedProducts.value = val ? props.seller.products : []
}

const handleCheckedProductsChange = (value: CheckboxValueType[]) => {
  const checkedCount = value.length
  checkAll.value = checkedCount === props.seller.products.length
}

const handleCheckAll = () => {
  checkAll.value = true
  checkedProducts.value = props.seller.products
}

const handleUncheckAll = () => {
  checkAll.value = false
  checkedProducts.value = []
}

const openChat = () => {
  if (!isLoggedIn.value) return
  const event = new CustomEvent('open-chat', {
    detail: { sellerId: props.seller.sellerId },
  })
  window.dispatchEvent(event)
}

defineExpose({ checkedProducts, handleCheckAll, handleUncheckAll })
</script>

<template>
  <div class="cart-box-shadow border-radius" style="margin: 24px 0; background-color: #fff">
    <div style="padding: 0 20px; display: flex; align-items: center; height: 64px">
      <el-checkbox
        v-model="checkAll"
        @change="handleCheckAllChange"
        size="large"
        class="cartCheckbox"
      />
      <h4 style="font-weight: 500">{{ props.seller.sellerName }}</h4>
      <div style="margin-left: 12px; display: flex; align-items: center">
        <el-button
          v-if="isLoggedIn"
          :icon="ChatDotRound"
          size="small"
          link
          class="shop-action-btn"
          @click="openChat"
          >Chat</el-button
        >
        <el-divider direction="vertical" />
        <el-button
          :icon="Shop"
          size="small"
          link
          class="shop-action-btn"
          @click="router.push(`/seller-page/${props.seller.sellerId}`)"
        >
          View Shop
        </el-button>
      </div>
    </div>
    <el-divider style="margin: 0 0 15px 0" />
    <el-checkbox-group v-model="checkedProducts" @change="handleCheckedProductsChange">
      <CartProduct
        v-for="product in props.seller.products"
        :key="product.id"
        :product="product"
        @delete="(id) => emits('deleteProduct', id)"
        @change="(id, val) => emits('updateQuantity', id, val)"
      />
    </el-checkbox-group>
  </div>
</template>

<style scoped>
.shop-action-btn {
  color: #666;
}

.shop-action-btn:hover {
  color: var(--main-color) !important;
}
</style>
