<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import axios from 'axios'
import ProfileLayout from '@/components/ProfileLayout.vue'

const router = useRouter()
const route = useRoute()

const name = ref('')
const email = ref('')
const phone = ref('')
const avatar = ref('')
const isLoading = ref(false)

const fetchUserInfo = async () => {
  const token = localStorage.getItem('access_token')
  if (!token) return

  isLoading.value = true
  try {
    const response = await axios.get(`${import.meta.env.VITE_BE_API_URL}/user/my-info`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })

    const userData = response.data.data

    name.value = userData.name || ''
    email.value = userData.email || ''
    phone.value = userData.phone || ''
    if (userData.image) {
      avatar.value = userData.image
    }
  } catch (error) {
    console.error('Failed to fetch user info:', error)
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchUserInfo()
})

const menuItems = [
  { id: 'account-info', label: 'Account info', path: '/profile/account-info' },
  { id: 'my-order', label: 'My order', path: '/profile/orders' },
  { id: 'my-voucher', label: 'My vouchers', path: '/profile/vouchers' },
  { id: 'my-address', label: 'My address', path: '/profile/address' },
  { id: 'change-password', label: 'Change password', path: '/profile/change-password' },
]

const handleMenuClick = (item: any) => {
  if (item.path) {
    router.push(item.path)
  }
}

const activeItem = computed(() => {
  const active = menuItems.find((item) => route.path.includes(item.path))
  return active ? active.id : 'account-info'
})

const userInfo = computed(() => ({
  name: name.value,
  email: email.value,
  avatar: avatar.value,
}))
</script>

<template>
  <ProfileLayout
    :loading="isLoading"
    :user="userInfo"
    :menu-items="menuItems"
    :active-item="activeItem"
    @menu-click="handleMenuClick"
  >
    <router-view v-slot="{ Component }">
      <component
        :is="Component"
        v-model:name="name"
        v-model:email="email"
        v-model:phone="phone"
        v-model:avatar="avatar"
      />
    </router-view>
  </ProfileLayout>
</template>

<style scoped>
/* Scoped styles removed as they are now in ProfileLayout */
</style>
