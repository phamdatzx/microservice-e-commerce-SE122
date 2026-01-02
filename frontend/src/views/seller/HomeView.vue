<script setup lang="ts">
import { ArrowDown, ArrowRight, Goods, Menu, ChatDotRound, User } from '@element-plus/icons-vue'

const handleOpen = (key: string, keyPath: string[]) => {
  console.log(key, keyPath)
}
const handleClose = (key: string, keyPath: string[]) => {
  console.log(key, keyPath)
}
</script>

<template>
  <div class="common-layout">
    <el-container>
      <el-header
        class="header"
        style="background-color: white; box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1); height: 64px"
      >
        <div class="container">
          <div class="header-content">
            <div class="logo">
              <div class="logo-icon">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
                  <path
                    d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                    stroke="white"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  />
                </svg>
              </div>
              <div class="logo-text">
                <span class="logo-name">SWOO</span>
                <span class="logo-tagline">TECH MART</span>
              </div>
            </div>

            <div class="breadcrumb">
              <div class="container">
                <a href="#" class="breadcrumb-link">Home</a>
                <span class="breadcrumb-separator">
                  <el-icon><ArrowRight /></el-icon>
                </span>
                <a href="#" class="breadcrumb-link">pages</a>
                <span class="breadcrumb-separator">
                  <el-icon><ArrowRight /></el-icon>
                </span>
                <span class="breadcrumb-current">login</span>
              </div>
            </div>

            <div class="user-section">
              <img
                src="https://placehold.co/40x40"
                alt="User Avatar"
                class="user-avatar"
                style="border-radius: 100%"
              />
              <span>John Doe</span>
              <el-icon><ArrowDown /></el-icon>
            </div>
          </div>
        </div>
      </el-header>
      <el-container style="height: calc(100vh - 64px)">
        <el-aside width="200px" style="margin-top: 3px">
          <el-menu default-active="0" @open="handleOpen" @close="handleClose" style="height: 100%">
            <el-sub-menu index="1">
              <template #title>
                <el-icon><location /></el-icon>
                <span>Navigator One</span>
              </template>
              <el-menu-item-group title="Group One">
                <el-menu-item index="1-1">item one</el-menu-item>
                <el-menu-item index="1-2">item two</el-menu-item>
              </el-menu-item-group>
              <el-menu-item-group title="Group Two">
                <el-menu-item index="1-3">item three</el-menu-item>
              </el-menu-item-group>
              <el-sub-menu index="1-4">
                <template #title>item four</template>
                <el-menu-item index="1-4-1">item one</el-menu-item>
              </el-sub-menu>
            </el-sub-menu>
            <RouterLink to="/seller/category">
              <el-menu-item index="2">
                <el-icon><Menu /></el-icon>
                <span>Category</span>
              </el-menu-item>
            </RouterLink>
            <RouterLink to="/seller/product">
              <el-menu-item index="3">
                <el-icon><Goods /></el-icon>
                <span>Product</span>
              </el-menu-item>
            </RouterLink>
            <RouterLink to="/seller/chat">
              <el-menu-item index="4">
                <el-icon><ChatDotRound /></el-icon>
                <span>Chat</span>
              </el-menu-item>
            </RouterLink>
            <RouterLink to="/seller/profile">
              <el-menu-item index="5">
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
  background: linear-gradient(135deg, #22c55e, #16a34a);
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

/* Breadcrumb */
.breadcrumb {
  background: white;
  margin-left: 24px;
}

.breadcrumb .container {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 18px;
}

.breadcrumb-link {
  color: #666;
  text-decoration: none;
}

.breadcrumb-link:hover {
  color: #22c55e;
}

.breadcrumb-separator {
  color: #999;
  margin: 0 4px;
}

.breadcrumb-current {
  color: #333;
  font-weight: 600;
}

/* User section */
.user-section {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  padding: 12px 16px;

  &:hover {
    background-color: #f5f5f5;
  }
}
</style>
