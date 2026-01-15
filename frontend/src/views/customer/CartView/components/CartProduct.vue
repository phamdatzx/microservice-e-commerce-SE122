<script setup lang="ts">
import { Delete } from '@element-plus/icons-vue'
import { type Product } from '../CartView.vue'
import { formatNumberWithDots } from '@/utils/formatNumberWithDots'

interface Props {
  product: Product
}

const props = withDefaults(defineProps<Props>(), {
  product: () => ({
    id: '',
    imageUrl: '',
    productName: '',
    productOption: '',
    price: 0,
    quantity: 0,
    isSoldOut: false,
    sellerName: '',
    sellerId: '',
    variantId: '',
    productId: '',
    stock: 0,
  }),
})

const emits = defineEmits(['delete'])

const handleDelete = () => {
  emits('delete', props.product.id)
}
</script>

<template>
  <div style="padding: 15px 20px; display: flex; align-items: center">
    <el-checkbox :value="props.product" class="cartCheckbox" />
    <RouterLink :to="`/product/${props.product.productId}`">
      <el-image style="width: 80px; height: 80px" fit="cover" :src="props.product.imageUrl" />
    </RouterLink>
    <div style="padding: 5px 20px 5px 10px; flex: 1">
      <RouterLink :to="`/product/${props.product.productId}`" class="product-name-link">
        <p class="two-line-ellipsis" style="margin-bottom: 5px">
          {{ props.product.productName }}
        </p>
      </RouterLink>
      <p
        style="line-height: 16px"
        :style="{
          color: props.product.isSoldOut ? 'var(--main-color)' : 'rgba(0, 0, 0, 0)',
        }"
      >
        This product variant is sold out, please select another option.
      </p>
    </div>
    <div style="margin-right: 10px; width: 150px">
      <template v-if="props.product.productOption">
        <span style="display: block">Variant:</span>
        <span style="display: block">{{ props.product.productOption }}</span>
      </template>
    </div>
    <span style="width: 174px; text-align: center"
      >{{ formatNumberWithDots(props.product.price) }}đ</span
    >
    <span style="width: 169px; text-align: center">
      <el-input-number
        v-model="props.product.quantity"
        :min="1"
        :max="props.product.stock"
        size="large"
        style="width: 140px"
      />
    </span>
    <span style="width: 114px; text-align: center"
      >{{ formatNumberWithDots(props.product.price * (props.product.quantity ?? 0)) }}đ</span
    >
    <span style="width: 139px; text-align: center">
      <el-button type="danger" plain :icon="Delete" @click="handleDelete" />
    </span>
  </div>
</template>

<style scoped>
.product-name-link {
  text-decoration: none;
  color: inherit;
  transition: color 0.2s;
}

.product-name-link:hover {
  color: var(--main-color);
}
</style>
