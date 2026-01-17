<script setup lang="ts">
import { formatNumberWithDots } from '@/utils/formatNumberWithDots'
import { quantityFormatNumber } from '@/utils/quantityFormatNumber'
import LocationIcon from './icons/LocationIcon.vue'
import StarIcon from './icons/StarIcon.vue'
import { Picture } from '@element-plus/icons-vue'

const props = defineProps([
  'imageUrl',
  'name',
  'price',
  'rating',
  'location',
  'discount',
  'id',
  'soldCount',
])
</script>

<template>
  <RouterLink :to="`/product/${props.id}`" class="wrapper" style="display: block">
    <el-image
      :src="props.imageUrl"
      fit="cover"
      style="width: 100%; aspect-ratio: 1/1; display: block"
      loading="lazy"
    >
      <template #error>
        <div class="image-placeholder">
          <el-icon><Picture /></el-icon>
        </div>
      </template>
    </el-image>
    <div style="padding: 12px 8px">
      <p class="product-name two-line-ellipsis">{{ props.name }}</p>
      <div style="display: flex; align-items: center; margin-top: 8px">
        <p class="product-price">{{ formatNumberWithDots(props.price) }}₫</p>
        <span
          v-if="props.discount > 0"
          style="
            line-height: 14px;
            height: 14px;
            display: block;
            background-color: #f0fdf4;
            margin-left: 4px;
            padding: 0 4px;
            color: var(--main-color);
            font-size: 10px;
            border-radius: 2px;
          "
          >-{{ props.discount }}%</span
        >
      </div>
      <div
        style="display: flex; justify-content: space-between; align-items: center; margin-top: 8px"
      >
        <div style="display: flex; align-items: center; gap: 4px">
          <div style="display: flex; align-items: center">
            <StarIcon style="margin-right: 2px; width: 10px; height: 10px" />
            <span style="font-size: 10px; color: #333">{{ Number(props.rating).toFixed(1) }}</span>
          </div>
          <div
            v-if="props.soldCount"
            style="
              font-size: 10px;
              color: #757575;
              border-left: 1px solid #dbdbdb;
              padding-left: 4px;
            "
          >
            Sold
            {{ quantityFormatNumber(props.soldCount) }}
          </div>
        </div>
        <div style="display: flex; align-items: center">
          <LocationIcon style="margin-right: 2px; width: 10px; height: 10px" />
          <span style="font-size: 10px; color: #757575">{{ props.location }}</span>
        </div>
      </div>
    </div>
  </RouterLink>
</template>

<style lang="css" scoped>
.wrapper {
  background: #fff;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.1);
  border-radius: 2px;
  margin: 0;
  height: 100%;
  transition: transform 0.1s ease;
  &:hover {
    cursor: pointer;
    box-shadow: 0 1px 20px 0 rgba(0, 0, 0, 0.05);
    transform: translateY(-1px);
    z-index: 1;
  }
}

.image-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f5f5f5;
  color: #ccc;
  font-size: 24px;
}

.product-name {
  color: #000000cc;
  font-size: 14px;
  line-height: 20px;
  height: 40px;
  margin: 0;
  word-break: break-all;
}

.two-line-ellipsis {
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  overflow: hidden;
}

.product-price {
  font-size: 16px;
  font-weight: 500;
  color: var(--main-color);
  margin: 0;
}
</style>
