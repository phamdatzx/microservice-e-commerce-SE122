<script setup lang="ts">
const props = defineProps(['imageUrl', 'name', 'id'])
</script>

<template>
  <RouterLink
    :to="`/seller-page/${props.id}`"
    style="display: flex; flex-direction: column; align-items: center; transition: transform 0.2s"
    class="seller-item"
  >
    <div class="avatar-container">
      <img
        :src="props.imageUrl || 'https://placehold.co/100x100?text=Seller'"
        class="seller-avatar"
      />
    </div>
    <p class="seller-name">{{ props.name }}</p>
  </RouterLink>
</template>

<style scoped>
.seller-item:hover {
  transform: translateY(-5px);
}

.avatar-container {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  overflow: hidden;
  border: 4px solid #f0fdf4;
  box-shadow: 0 4px 6px -1px rgb(0 0 0 / 0.1);
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #fff;
}

.seller-avatar {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.seller-name {
  text-align: center;
  margin-top: 8px;
  font-weight: 600;
  color: #333;
  font-size: 14px;
  width: 120px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
