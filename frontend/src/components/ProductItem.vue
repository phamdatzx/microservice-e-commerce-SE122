<script setup lang="ts">
import { formatNumberWithDots } from '@/utils/formatNumberWithDots'
import LocationIcon from './icons/LocationIcon.vue'
import StarIcon from './icons/StarIcon.vue'

const props = defineProps(['imageUrl', 'name', 'price', 'rating', 'location', 'discount', 'id'])
</script>

<template>
  <RouterLink :to="`/product`" class="wrapper" style="display: block">
    <div
      :style="{
        background: `url(${props.imageUrl}) no-repeat center`,
      }"
      style="width: 100%; padding-bottom: 100%; background-size: cover"
    />
    <div style="padding: 8px">
      <p class="product-name">{{ props.name }}</p>
      <div style="display: flex; align-items: center">
        <p class="product-price">{{ formatNumberWithDots(props.price) }}₫</p>
        <span
          style="
            --line-height: 14px;
            line-height: var(--line-height);
            height: var(--line-height);

            display: block;
            background-color: #feeeea;
            margin-left: 4px;
            padding: 0 2px;
            color: #ee4d2d;
            font-size: 10px;
          "
          >-{{ props.discount }}%</span
        >
      </div>
      <div style="display: flex; justify-content: space-between">
        <div>
          <StarIcon style="margin-right: 4px" />
          <span>{{ props.rating }}</span>
        </div>
        <div>
          <LocationIcon style="margin-right: 4px" />
          <span>{{ props.location }}</span>
        </div>
      </div>
    </div>
  </RouterLink>
</template>

<style lang="css" scoped>
.wrapper {
  box-shadow: rgba(99, 99, 99, 0.2) 0px 2px 8px 0px;
  border-radius: 2px;
  margin: 20px 0;
  transition: box-shadow 0.2s;
  &:hover {
    cursor: pointer;
    box-shadow: rgba(149, 157, 165, 0.2) 0px 8px 24px;
  }
}

.product-name {
  --line-height: 20px;
  height: calc(var(--line-height) * 2);
  line-height: var(--line-height);

  font-size: 16px;
  font-weight: 500;
  color: #000000cc;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  line-clamp: 2;
  -webkit-line-clamp: 2;
  overflow: hidden;
  text-overflow: ellipsis;
}

.product-price {
  font-size: 18px;
  color: var(--main-color);
  margin: 5px;
}
</style>
