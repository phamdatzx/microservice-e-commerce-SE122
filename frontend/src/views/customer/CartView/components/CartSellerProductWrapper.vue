<script setup lang="ts">
import type { CheckboxValueType } from 'element-plus'
import { ref, useTemplateRef, watch } from 'vue'
import { type CartSeller, type Product } from '../CartView.vue'
import { ChatLineRound } from '@element-plus/icons-vue'
import CartProduct from './CartProduct.vue'

interface Props {
  seller: CartSeller
}

const props = withDefaults(defineProps<Props>(), {
  seller: () => ({
    sellerName: '',
    products: [],
  }),
})

const checkAll = ref(false)
const checkedProducts = ref<Product[]>([])

const handleCheckAllChange = (val: CheckboxValueType) => {
  checkedProducts.value = val ? props.seller.products : []
}

const handleCheckedProductsChange = (value: CheckboxValueType[]) => {
  const checkedCount = value.length
  checkAll.value = checkedCount === props.seller.products.length
}

watch([checkAll], () => {})
</script>

<template>
  <div class="box-shadow border-radius" style="margin: 20px 0">
    <div style="padding: 0 20px; display: flex; align-items: center; height: 64px">
      <el-checkbox
        ref="checkbox"
        :id="props.seller.sellerName"
        :value="props.seller"
        @change="handleCheckAllChange"
        size="large"
        class="cartCheckbox"
      />
      <h4 style="font-weight: 500">{{ props.seller.sellerName }}</h4>
      <el-button type="plain" circle style="margin-left: 8px">
        <el-icon size="large"><ChatLineRound /></el-icon>
      </el-button>
    </div>
    <el-divider style="margin: 0 0 15px 0" />
    <el-checkbox-group v-model="checkedProducts" @change="handleCheckedProductsChange">
      <CartProduct v-for="product in props.seller.products" :key="product.id" :product="product" />
    </el-checkbox-group>
  </div>
</template>
