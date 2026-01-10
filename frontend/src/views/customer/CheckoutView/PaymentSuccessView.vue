<script setup lang="ts">
import { useRouter } from 'vue-router'
import { Check, View, House } from '@element-plus/icons-vue'

const router = useRouter()

const viewOrder = () => {
  router.push('/profile') // Redirect to order history
}

const goHome = () => {
  router.push('/')
}
</script>

<template>
  <div class="payment-success-page">
    <div class="success-card">
      <div class="icon-wrapper">
        <el-icon class="success-icon"><Check /></el-icon>
      </div>
      <h2>Payment Successful!</h2>
      <p>Thank you for your purchase. Your order has been placed successfully.</p>

      <div class="actions">
        <el-button type="primary" size="large" @click="viewOrder">
          <el-icon style="margin-right: 8px"><View /></el-icon> View Order
        </el-button>
        <el-button size="large" @click="goHome">
          <el-icon style="margin-right: 8px"><House /></el-icon> Back to Home
        </el-button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.payment-success-page {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f5f5f5;
  padding: 20px;
}

.success-card {
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
  background-color: #dcfce7; /* Green-100 */
  border-radius: 50%;
  display: flex;
  justify-content: center;
  align-items: center;
  margin: 0 auto 20px;
}

.success-icon {
  font-size: 40px;
  color: #16a34a; /* Green-600 */
}

h2 {
  color: #166534; /* Green-800 */
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
