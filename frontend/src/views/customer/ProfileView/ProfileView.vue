<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ArrowRight, Loading } from '@element-plus/icons-vue'
import { useRouter, useRoute } from 'vue-router'
import axios from 'axios'

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

const handleMenuClick = (path: string) => {
  router.push(path)
}

const isActive = (path: string) => {
  return route.path.includes(path)
}
</script>

<template>
  <div class="main-container">
    <el-row :gutter="28">
      <!-- Sidebar -->
      <el-col :span="6">
        <div class="sidebar box-shadow border-radius">
          <div class="user-profile">
            <div v-if="isLoading" class="loading-state">
              <el-icon class="is-loading"><Loading /></el-icon>
            </div>
            <div v-else>
              <div class="avatar-wrapper">
                <img :src="avatar" alt="User Avatar" />
              </div>
              <h3 class="user-name">{{ name }}</h3>
              <p class="user-email">{{ email }}</p>
            </div>
          </div>

          <div class="sidebar-menu">
            <div
              v-for="item in menuItems"
              :key="item.id"
              class="menu-item"
              :class="{ active: isActive(item.path) }"
              @click="handleMenuClick(item.path)"
            >
              <span>{{ item.label }}</span>
              <el-icon><ArrowRight /></el-icon>
            </div>
          </div>
        </div>
      </el-col>

      <!-- Main Content -->
      <el-col :span="18">
        <div class="content-area">
          <router-view v-slot="{ Component }">
            <component
              :is="Component"
              v-model:name="name"
              v-model:email="email"
              v-model:phone="phone"
              v-model:avatar="avatar"
            />
          </router-view>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<style scoped>
.main-container {
  margin-top: 20px;
  margin-bottom: 20px;
}

.sidebar {
  background-color: #fff;
  padding-bottom: 20px;
  overflow: hidden;
}

.user-profile {
  padding: 30px 20px;
  text-align: center;
  border-bottom: 1px solid #f0f0f0;
  background-color: #f9f9f9;
}

.loading-state {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 214px;
  font-size: 24px;
  color: #888;
}

.avatar-wrapper {
  width: 140px;
  height: 140px;
  margin: 0 auto 16px;
  border-radius: 12px;
  overflow: hidden;
  background-color: #e0e0e0;
}

.avatar-wrapper img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.user-name {
  font-size: 20px;
  font-weight: 700;
  color: #000;
  margin-bottom: 4px;
}

.user-email {
  font-size: 14px;
  color: #888;
}

.sidebar-menu {
  padding: 10px 15px;
}

.menu-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  margin: 8px 0;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
  font-weight: 500;
  color: #444;
}

.menu-item:hover {
  background-color: #f5f5f5;
  color: var(--main-color);
}

.menu-item.active {
  background-color: var(--main-color);
  color: #fff;
}

.content-area {
  background-color: #fff;
  padding: 20px;
  margin-bottom: 40px;
}

.section-title {
  font-size: 28px;
  font-weight: 700;
  margin-bottom: 30px;
  color: #000;
}
</style>
