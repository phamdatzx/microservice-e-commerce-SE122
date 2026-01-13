<script setup lang="ts">
import Header from '@/components/Header.vue'
import { ref, provide } from 'vue'

const headerRef = ref<InstanceType<typeof Header> | null>(null)

const fetchCartCount = () => {
  headerRef.value?.fetchCartCount()
}

provide('fetchCartCount', fetchCartCount)
</script>

<template>
  <Header ref="headerRef" />
  <router-view />
</template>
