<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import { useRoute } from 'vue-router'
import { ArrowLeft, ArrowRight } from '@element-plus/icons-vue'
import axios from 'axios'
import ProfileLayout from '@/components/ProfileLayout.vue'
import AccountInfo from '@/views/customer/ProfileView/components/AccountInfo.vue'
import ChangePassword from '@/views/customer/ProfileView/components/ChangePassword.vue'
import MyAddress from '@/views/customer/ProfileView/components/MyAddress.vue'

const name = ref('')
const email = ref('')
const phone = ref('')
const avatar = ref('')
const isLoading = ref(false)

const route = useRoute()
const activeMenu = ref('account-info')

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
  if (route.query.tab) {
    activeMenu.value = route.query.tab as string
  }
})

watch(
  () => route.query.tab,
  (newTab) => {
    if (newTab) {
      activeMenu.value = newTab as string
    }
  },
)

const menuItems = [
  { id: 'account-info', label: 'Account info' },
  { id: 'my-address', label: 'My address' },
  { id: 'change-password', label: 'Change password' },
]

const handleMenuClick = (item: any) => {
  activeMenu.value = item.id
}

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
    :active-item="activeMenu"
    @menu-click="handleMenuClick"
    style="margin-top: 0; margin-bottom: 0"
  >
    <div class="content-wrapper box-shadow border-radius">
      <AccountInfo
        v-if="activeMenu === 'account-info'"
        v-model:name="name"
        v-model:email="email"
        v-model:phone="phone"
        v-model:avatar="avatar"
      />

      <ChangePassword v-else-if="activeMenu === 'change-password'" />

      <MyAddress v-else-if="activeMenu === 'my-address'" />

      <!-- Fallback -->
      <div v-else>
        <h2 class="section-title">Coming Soon</h2>
        <div style="padding: 40px; text-align: center; color: #999">
          This section is under development.
        </div>
      </div>
    </div>
  </ProfileLayout>
</template>

<style scoped>
.content-wrapper {
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
  min-height: 500px;
}

.section-title {
  font-size: 28px;
  font-weight: 700;
  margin-bottom: 30px;
  color: #000;
}
</style>
