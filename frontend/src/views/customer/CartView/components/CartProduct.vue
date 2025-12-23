<script setup lang="ts">
import { Delete } from '@element-plus/icons-vue'
import { type Product } from '../CartView.vue'

interface Props {
  product: Product
}

const props = withDefaults(defineProps<Props>(), {
  product: () => ({
    id: -1,
    imageUrl: '',
    productName: '',
    productOption: '',
    price: 0,
    quantity: 0,
    isSoldOut: false,
    sellerName: '',
  }),
})

const quantity = defineModel<number>()
</script>

<template>
  <div style="padding: 15px 20px; display: flex; align-items: center">
    <el-checkbox :value="props.product" class="cartCheckbox" />
    <el-image style="width: 80px; height: 80px" fit="cover" :src="props.product.imageUrl" />
    <div style="padding: 5px 20px 5px 10px; flex: 1">
      <p class="product-name" style="margin-bottom: 5px">
        {{ props.product.productName }}
      </p>
      <p
        style="line-height: 16px"
        :style="{
          color: props.product.isSoldOut ? 'rgb(238, 77, 45)' : 'rgba(0, 0, 0, 0)',
        }"
      >
        This product option is sold out, please select another option.
      </p>
    </div>
    <div style="margin-right: 10px">
      <span style="display: block">Product Option:</span>
      <span style="display: block">{{ props.product.productOption }}</span>
    </div>
    <span style="width: 174px; text-align: center">{{ props.product.price }}$</span>
    <span style="width: 169px; text-align: center">
      <el-input-number
        v-model="props.product.quantity"
        :min="1"
        :max="10"
        size="large"
        style="width: 140px"
      />
    </span>
    <span style="width: 114px; text-align: center"
      >{{ props.product.price * (props.product.quantity ?? 0) }}$</span
    >
    <span style="width: 139px; text-align: center">
      <el-button type="danger" plain :icon="Delete" />
    </span>
  </div>
</template>
