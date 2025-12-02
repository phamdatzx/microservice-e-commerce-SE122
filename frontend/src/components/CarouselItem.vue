<script setup lang="ts">
const props = defineProps(['imageUrl', 'to'])
</script>

<template>
  <img :src="props.imageUrl" style="width: 100%; height: 100%; object-fit: cover" />
  <router-link
    class="button"
    :to="props.to"
    style="
      position: absolute;
      left: 24px;
      bottom: 20px;
      border-radius: 4px;
      background-color: #fff;
      padding: 5px 18px;
    "
    >Buy Now</router-link
  >
</template>

<style lang="css" scoped>
.button {
  transition: opacity 0.2s;
  box-shadow: rgba(149, 157, 165, 0.2) 0px 8px 24px;

  &:hover {
    opacity: 0.8;
  }
}
</style>
