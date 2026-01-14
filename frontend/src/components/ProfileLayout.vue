<script setup lang="ts">
import { Loading, ArrowRight } from '@element-plus/icons-vue'

interface MenuItem {
  id: string
  label: string
}

defineProps<{
  loading: boolean
  user: {
    name: string
    email: string
    avatar: string
  }
  menuItems: MenuItem[]
  activeItem: string
}>()

const emit = defineEmits<{
  (e: 'menu-click', item: MenuItem): void
}>()

const handleMenuClick = (item: MenuItem) => {
  emit('menu-click', item)
}
</script>

<template>
  <div class="main-container">
    <el-row :gutter="28">
      <!-- Sidebar -->
      <el-col :span="6">
        <div class="sidebar box-shadow border-radius">
          <div class="user-profile">
            <div v-if="loading" class="loading-state">
              <el-icon class="is-loading"><Loading /></el-icon>
            </div>
            <div v-else>
              <div class="avatar-wrapper">
                <img :src="user.avatar || '/src/assets/avatar.jpg'" alt="User Avatar" />
              </div>
              <h3 class="user-name">{{ user.name }}</h3>
              <p class="user-email">{{ user.email }}</p>
            </div>
          </div>

          <div class="sidebar-menu">
            <div
              v-for="item in menuItems"
              :key="item.id"
              class="menu-item"
              :class="{ active: activeItem === item.id }"
              @click="handleMenuClick(item)"
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
          <slot></slot>
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
</style>
