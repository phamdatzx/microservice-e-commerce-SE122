<script setup lang="ts">
import { ArrowDown } from '@element-plus/icons-vue'
import ProductSellerInfo from '@/components/ProductSellerInfo.vue'

defineProps<{
  sellerInfo: {
    id: string
    name: string
    image: string
    address: {
      province: string
    }
    sale_info: {
      follow_count: number
      rating_count: number
      rating_average: number
      product_count: number
      is_following: boolean
    }
  }
}>()

defineEmits(['toggle-follow'])
</script>

<template>
  <div class="seller-header">
    <div class="container">
      <ProductSellerInfo
        layout="header"
        :seller-info="sellerInfo"
        :show-follow-button="true"
        :show-view-shop="false"
        @toggle-follow="$emit('toggle-follow')"
      />
    </div>
  </div>

  <!-- Shop Navigation Tabs -->
  <div class="shop-nav">
    <div class="container">
      <el-row class="nav-tabs">
        <el-col :span="4" class="nav-tab active">
          <a href="#"><span class="one-line-ellipsis">Home</span></a>
        </el-col>
        <el-col :span="4" class="nav-tab">
          <a href="#"><span class="one-line-ellipsis">ALL PRODUCTS</span></a>
        </el-col>
        <el-col :span="4" class="nav-tab">
          <a href="#"><span class="one-line-ellipsis">Cargo Jeans ⭐</span></a>
        </el-col>
        <el-col :span="4" class="nav-tab">
          <a href="#"><span class="one-line-ellipsis">Skinny Jeans ⭐</span></a>
        </el-col>
        <el-col :span="4" class="nav-tab">
          <a href="#"><span class="one-line-ellipsis">WASH RETRO ...</span></a>
        </el-col>
        <el-col :span="4" class="nav-tab">
          <a href="#"><span class="one-line-ellipsis">Shorts ⭐</span></a>
        </el-col>
        <el-col :span="4" class="nav-tab dropdown-tab">
          <el-dropdown class="dropdown" trigger="hover">
            <div class="dropdown-trigger">
              <span class="one-line-ellipsis">More</span>
              <el-icon class="el-icon--right"><ArrowDown /></el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu class="custom-dropdown-menu">
                <el-dropdown-item>Top Sales</el-dropdown-item>
                <el-dropdown-item>New Arrivals</el-dropdown-item>
                <el-dropdown-item>Summer Collection</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<style scoped>
.seller-header {
  background-color: #fff;
  padding: 20px 0;
  border-bottom: 1px solid #f5f5f5;
}

/* Shop Nav */
.shop-nav {
  background-color: #fff;
  border-bottom: 1px solid #e5e5e5;
  position: sticky;
  top: 0;
  z-index: 10;
}

.nav-tabs {
  display: flex;
  height: 60px;
  align-items: center;
}

.nav-tab {
  flex: 1;
  height: 100%;
  display: flex;
  align-items: center;
  text-decoration: none;
  color: #222;
  font-size: 14px;
  text-transform: uppercase;
  border-bottom: 3px solid transparent;
  transition: all 0.2s;
  display: flex;
  justify-content: center;

  a,
  .dropdown {
    display: block;
    width: 100%;
    height: 100%;
    transition: all 0.2s;
    display: flex;
    justify-content: center;
    align-items: center;

    padding: 0 20px;

    &:hover {
      color: var(--main-color);
    }
  }
}

.nav-tab.active {
  color: var(--main-color);
  border-bottom-color: var(--main-color);
}

.dropdown-tab {
  cursor: pointer;
  outline: none;
}

.dropdown-trigger {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  padding: 0 20px;
  color: #222;
  transition: all 0.2s;

  &:hover {
    color: var(--main-color);
  }
}

:deep(.custom-dropdown-menu) {
  padding: 10px 0;
  border-radius: 4px;
  border: none;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);

  .el-dropdown-menu__item {
    font-size: 14px;
    padding: 10px 25px;
    color: #333;

    &:hover {
      background-color: #f5f5f5;
      color: var(--main-color);
    }
  }
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 15px;
}
</style>
