<script setup lang="ts">
import Header from '@/components/Header.vue'
import { ref } from 'vue'
import { ArrowRight } from '@element-plus/icons-vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()

const name = ref('Mark Cole')
const email = ref('swoo@gmail.com')
const phone = ref('+1 0231 4554 452')

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
  <Header />

  <div class="main-container">
    <el-row :gutter="28">
      <!-- Sidebar -->
      <el-col :span="6">
        <div class="sidebar box-shadow border-radius">
          <div class="user-profile">
            <div class="avatar-wrapper">
              <img src="/src/assets/avatar.jpg" alt="User Avatar" />
            </div>
            <h3 class="user-name">{{ name }}</h3>
            <p class="user-email">swoo@gmail.com</p>
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
            />
          </router-view>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<style scoped>
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
