<script setup lang="ts">
import Header from '@/components/Header.vue'
import SellerHeader from './components/SellerHeader.vue'
import SellerVouchers from './components/SellerVouchers.vue'
import SellerRecommendations from './components/SellerRecommendations.vue'
import SellerProductList from './components/SellerProductList.vue'
</script>

<template>
  <div class="seller-page">
    <Header />

    <main class="seller-main">
      <SellerHeader />

      <div class="content-bg">
        <div class="container">
          <SellerVouchers />

          <SellerRecommendations />

          <SellerProductList />

          <div class="pagination-wrapper">
            <el-pagination
              background
              layout="prev, pager, next"
              :total="100"
              :page-size="20"
              size="large"
              class="custom-pagination"
            />
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<style scoped>
.seller-page {
  background-color: #f5f5f5;
  min-height: 100vh;
}

.content-bg {
  background-color: #f5f5f5;
  padding-bottom: 20px;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 15px;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  padding-bottom: 20px;
}

:deep(.custom-pagination.is-background .el-pager li:not(.is-disabled).is-active) {
  background-color: var(--main-color) !important;
  color: #fff;
}

:deep(.custom-pagination.is-background .el-pager li:not(.is-disabled):hover) {
  color: var(--main-color);
}
</style>
