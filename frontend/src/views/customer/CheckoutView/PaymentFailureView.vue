<script setup lang="ts">
import { useRouter } from 'vue-router'
import { Close, House, RefreshRight } from '@element-plus/icons-vue'

const router = useRouter()

const retryPayment = () => {
  router.push('/checkout')
}

const goHome = () => {
  router.push('/')
}
</script>

<template>
  <div class="payment-failure-page">
    <div class="failure-card">
      <div class="icon-wrapper">
        <el-icon class="failure-icon"><Close /></el-icon>
      </div>
      <h2>Payment Failed</h2>
      <p>
        Something went wrong with your payment. Please try again or contact support if the issue
        persists.
      </p>

      <div class="actions">
        <el-button type="danger" size="large" @click="retryPayment">
          <el-icon style="margin-right: 8px"><RefreshRight /></el-icon> Try Again
        </el-button>
        <el-button size="large" @click="goHome">
          <el-icon style="margin-right: 8px"><House /></el-icon> Back to Home
        </el-button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.payment-failure-page {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f5f5f5;
  padding: 20px;
}

.failure-card {
  background: white;
  padding: 40px;
  border-radius: 12px;
  text-align: center;
  box-shadow:
    0 4px 6px -1px rgba(0, 0, 0, 0.1),
    0 2px 4px -1px rgba(0, 0, 0, 0.06);
  max-width: 500px;
  width: 100%;
}

.icon-wrapper {
  width: 80px;
  height: 80px;
  background-color: #fee2e2; /* Red-100 */
  border-radius: 50%;
  display: flex;
  justify-content: center;
  align-items: center;
  margin: 0 auto 20px;
}

.failure-icon {
  font-size: 40px;
  color: #dc2626; /* Red-600 */
}

h2 {
  color: #991b1b; /* Red-800 */
  margin-bottom: 10px;
  font-size: 24px;
}

p {
  color: #6b7280;
  margin-bottom: 30px;
  line-height: 1.5;
}

.actions {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

@media (min-width: 640px) {
  .actions {
    flex-direction: row;
    justify-content: center;
  }
}
</style>
