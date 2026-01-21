<script setup lang="ts">
import { ref } from 'vue'
import { Picture } from '@element-plus/icons-vue'

const props = defineProps(['imageUrl', 'name'])
const imageError = ref(false)

const handleImageError = () => {
  imageError.value = true
}
</script>

<template>
  <div
    style="
      margin-top: 10px;
      display: flex;
      flex-direction: column;
      align-items: center;
      text-decoration: none;
      cursor: pointer;
    "
  >
    <div
      v-if="!props.imageUrl || imageError"
      class="placeholder-container"
      style="
        width: 80px;
        height: 80px;
        background-color: #f4f4f5;
        border-radius: 8px;
        display: flex;
        align-items: center;
        justify-content: center;
        color: #a1a1aa;
      "
    >
      <el-icon :size="32"><Picture /></el-icon>
    </div>
    <img
      v-else
      :src="props.imageUrl"
      style="width: 80px; height: 80px; object-fit: cover; display: block; border-radius: 8px"
      @error="handleImageError"
    />
    <p style="text-align: center; margin-top: 8px; font-size: 14px; color: #18181b">
      {{ props.name }}
    </p>
  </div>
</template>
