<script setup lang="ts">
import ChatWindow from '@/components/ChatWindow.vue'
</script>

<template>
  <div class="seller-chat-view">
    <ChatWindow :is-embedded="true" />
  </div>
</template>

<style scoped>
.seller-chat-view {
  height: calc(100vh - 64px - 40px); /* Adjust based on header and main padding */
  display: flex;
  flex-direction: column;
}
</style>
