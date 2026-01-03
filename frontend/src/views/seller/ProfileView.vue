<script setup lang="ts">
import { ref } from 'vue'
import { ArrowRight } from '@element-plus/icons-vue'
import AccountInfo from '@/views/customer/ProfileView/components/AccountInfo.vue'
import ChangePassword from '@/views/customer/ProfileView/components/ChangePassword.vue'
import MyAddress from '@/views/customer/ProfileView/components/MyAddress.vue'

const name = ref('Mark Cole')
const email = ref('swoo@gmail.com')
const phone = ref('+1 0231 4554 452')

const activeMenu = ref('account-info')

const menuItems = [
  { id: 'account-info', label: 'Account info' },
  { id: 'my-address', label: 'My address' },
  { id: 'change-password', label: 'Change password' },
]
</script>

<template>
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
              :class="{ active: activeMenu === item.id }"
              @click="activeMenu = item.id"
            >
              <span>{{ item.label }}</span>
              <el-icon><ArrowRight /></el-icon>
            </div>
          </div>
        </div>
      </el-col>

      <!-- Main Content -->
      <el-col :span="18">
        <div class="content-area box-shadow border-radius">
          <AccountInfo
            v-if="activeMenu === 'account-info'"
            v-model:name="name"
            v-model:email="email"
            v-model:phone="phone"
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
      </el-col>
    </el-row>
  </div>
</template>

<style scoped>
.main-container {
  padding: 20px;
}

.sidebar {
  background-color: #fff;
  padding-bottom: 20px;
  overflow: hidden;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
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
