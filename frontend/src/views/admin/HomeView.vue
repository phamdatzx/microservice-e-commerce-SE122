<script setup lang="ts">
import { Menu, User, Warning } from '@element-plus/icons-vue'
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import UserAvatarDropdown from '@/components/UserAvatarDropdown.vue'

const route = useRoute()

const activeIndex = computed(() => {
  const path = route.path
  if (path.includes('/admin/category')) return '1'
  if (path.includes('/admin/users')) return '2'
  if (path.includes('/admin/report')) return '3'
  if (path.includes('/admin/profile')) return '4'
  return '1'
})

const handleOpen = (key: string, keyPath: string[]) => {}
const handleClose = (key: string, keyPath: string[]) => {}
</script>

<template>
  <div class="common-layout admin-layout">
    <el-container>
      <el-header
        class="header"
        style="background-color: white; box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1); height: 64px"
      >
        <div class="container">
          <div class="header-content">
            <!-- Admin Logo Style (Blue/Purple gradient) -->
            <div class="logo">
              <div class="logo-icon">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
                  <path
                    d="M12 11c0 3.517-1.009 6.799-2.753 9.571m-3.44-2.04l.054-.09A13.916 13.916 0 008 11a4 4 0 118 0c0 1.017-.07 2.019-.203 3m-2.118 6.844A21.88 21.88 0 0015.171 17m3.839 1.132c.645-2.266.99-4.659.99-7.132A8 8 0 008 4.07M3 15.364c.64-1.319 1-2.8 1-4.364 0-1.457.2-2.858.59-4.18"
                    stroke="white"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  />
                </svg>
              </div>
              <div class="logo-text">
                <span class="logo-name">SWOO</span>
                <span class="logo-tagline">ADMIN PANEL</span>
              </div>
            </div>

            <UserAvatarDropdown role="admin" :show-address="false" style="margin-right: 20px" />
          </div>
        </div>
      </el-header>
      <el-container style="height: calc(100vh - 64px)">
        <el-aside width="200px" style="margin-top: 3px">
          <el-menu
            :default-active="activeIndex"
            @open="handleOpen"
            @close="handleClose"
            style="height: 100%"
          >
            <RouterLink to="/admin/category">
              <el-menu-item index="1">
                <el-icon><Menu /></el-icon>
                <span>Category</span>
              </el-menu-item>
            </RouterLink>
            <RouterLink to="/admin/users">
              <el-menu-item index="2">
                <el-icon><User /></el-icon>
                <span>User</span>
              </el-menu-item>
            </RouterLink>
            <RouterLink to="/admin/report">
              <el-menu-item index="3">
                <el-icon><Warning /></el-icon>
                <span>Report</span>
              </el-menu-item>
            </RouterLink>
            <RouterLink to="/admin/profile">
              <el-menu-item index="4">
                <el-icon><User /></el-icon>
                <span>Profile</span>
              </el-menu-item>
            </RouterLink>
          </el-menu>
        </el-aside>
        <el-main style="background-color: #f6f6f6">
          <RouterView />
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<style scoped>
.admin-layout {
  --el-color-primary: #409eff;
  --el-color-primary-light-3: #79bbff;
  --el-color-primary-light-5: #a0cfff;
  --el-color-primary-light-7: #c6e2ff;
  --el-color-primary-light-8: #d9ecff;
  --el-color-primary-light-9: #ecf5ff;
  --el-color-primary-dark-2: #337ecc;

  /* Also override main-color if used by custom components */
  --main-color: #409eff;
}

a {
  text-decoration: none;
}

/* Header Styles */
.container {
  margin: 0 auto;
}

.header {
  padding: 0 0;
}

.header-link {
  color: #666;
  text-decoration: none;
  font-size: 14px;
}

.header-link:hover {
  color: #333;
}

.header-content {
  display: flex;
  align-items: center;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-left: 20px;
}

.logo-icon {
  width: 40px;
  height: 40px;
  background: linear-gradient(135deg, #16a34a, #22c55e); /* Green gradient for Admin */
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-text {
  display: flex;
  flex-direction: column;
}

.logo-name {
  font-size: 24px;
  font-weight: 800;
  color: #333;
  line-height: 1;
}

.logo-tagline {
  font-size: 12px;
  color: #666;
  letter-spacing: 1px;
}
</style>
