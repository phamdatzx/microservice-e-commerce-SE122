<script setup lang="ts">
const props = defineProps(['imageUrl', 'name'])
</script>

<template>
  <RouterLink
    to="/"
    style="margin-top: 10px; display: flex; flex-direction: column; align-items: center"
  >
    <img :src="props.imageUrl" style="width: 50%; height: 100%; object-fit: contain" />
    <p style="text-align: center; margin-top: 4px">{{ props.name }}</p>
  </RouterLink>
</template>
