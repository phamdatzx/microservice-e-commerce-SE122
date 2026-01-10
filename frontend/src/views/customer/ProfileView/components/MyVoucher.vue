<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { Ticket, Timer, CopyDocument } from '@element-plus/icons-vue'
import { formatNumberWithDots } from '@/utils/formatNumberWithDots'
import { ElMessage } from 'element-plus'

interface SavedVoucher {
  id: string
  voucher: {
    code: string
    name: string
    description: string
    discount_type: 'FIXED' | 'PERCENTAGE'
    discount_value: number
    max_discount_value: number
    min_order_value: number
    end_time: string
    status: string
  }
}

const vouchers = ref<SavedVoucher[]>([])
const loading = ref(false)

const fetchSavedVouchers = async () => {
  loading.value = true
  try {
    const response = await axios.get('http://localhost:81/api/product/saved-vouchers', {
      headers: {
        Authorization: `Bearer ${localStorage.getItem('access_token')}`,
      },
    })
    if (response.data) {
      // The API returns an array directly as per the user's data example
      vouchers.value = response.data
    }
  } catch (error) {
    console.error('Failed to fetch saved vouchers:', error)
  } finally {
    loading.value = false
  }
}

const copyCode = (code: string) => {
  navigator.clipboard
    .writeText(code)
    .then(() => {
      ElMessage.success('Voucher code copied to clipboard')
    })
    .catch(() => {
      ElMessage.error('Failed to copy code')
    })
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('vi-VN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  })
}

onMounted(() => {
  fetchSavedVouchers()
})
</script>

<template>
  <div class="my-vouchers">
    <div class="header">
      <h3>My Saved Vouchers</h3>
    </div>

    <div v-loading="loading" class="voucher-list">
      <div v-if="vouchers.length === 0 && !loading" class="no-vouchers">
        <el-empty description="No saved vouchers found" />
      </div>

      <div v-for="item in vouchers" :key="item.id" class="voucher-card">
        <div class="voucher-left">
          <div class="icon-box">
            <el-icon :size="24"><Ticket /></el-icon>
          </div>
          <div class="voucher-info">
            <h4 class="voucher-name">{{ item.voucher.name }}</h4>
            <div class="voucher-code-wrapper">
              <span class="voucher-code">{{ item.voucher.code }}</span>
              <el-icon class="copy-icon" @click="copyCode(item.voucher.code)"
                ><CopyDocument
              /></el-icon>
            </div>
            <p class="voucher-desc">{{ item.voucher.description }}</p>
          </div>
        </div>

        <div class="voucher-right">
          <div class="discount-info">
            <span class="discount-amount">
              <template v-if="item.voucher.discount_type === 'FIXED'">
                {{ formatNumberWithDots(item.voucher.discount_value) }}đ off
              </template>
              <template v-else> {{ item.voucher.discount_value }}% off </template>
            </span>
            <span class="min-spend"
              >Min. Spend {{ formatNumberWithDots(item.voucher.min_order_value) }}đ</span
            >
          </div>
          <div class="expiry-info">
            <el-icon><Timer /></el-icon>
            <span>Exp: {{ formatDate(item.voucher.end_time) }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.my-vouchers {
  padding: 0 10px;
}

.header {
  margin-bottom: 20px;
  border-bottom: 1px solid #eee;
  padding-bottom: 10px;
}

.header h3 {
  font-size: 20px;
  color: #333;
  margin: 0;
}

.voucher-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.voucher-card {
  display: flex;
  background: #fff;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  overflow: hidden;
  transition: all 0.2s;
}

.voucher-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  border-color: var(--main-color);
}

.voucher-left {
  flex: 1;
  padding: 15px;
  display: flex;
  gap: 15px;
  align-items: center;
  border-right: 1px dashed #e0e0e0;
}

.icon-box {
  width: 50px;
  height: 50px;
  background-color: #fff0f0;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--main-color);
  flex-shrink: 0;
}

.voucher-info {
  flex: 1;
}

.voucher-name {
  margin: 0 0 5px;
  font-size: 16px;
  color: #333;
}

.voucher-code-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 5px;
}

.voucher-code {
  background-color: #f5f5f5;
  padding: 2px 8px;
  border-radius: 4px;
  font-family: monospace;
  font-weight: 600;
  color: #555;
  font-size: 14px;
}

.copy-icon {
  cursor: pointer;
  color: #888;
  font-size: 16px;
}

.copy-icon:hover {
  color: var(--main-color);
}

.voucher-desc {
  font-size: 13px;
  color: #888;
  margin: 0;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.voucher-right {
  width: 160px;
  padding: 15px;
  background-color: #fcfcfc;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: end;
  text-align: right;
}

.discount-info {
  display: flex;
  flex-direction: column;
  margin-bottom: 8px;
}

.discount-amount {
  font-size: 18px;
  font-weight: 700;
  color: var(--main-color);
}

.min-spend {
  font-size: 12px;
  color: #666;
}

.expiry-info {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #999;
}
</style>
