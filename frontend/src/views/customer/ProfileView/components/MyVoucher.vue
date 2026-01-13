<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import axios from 'axios'
import { Ticket, Timer, CopyDocument, Delete } from '@element-plus/icons-vue'
import { formatNumberWithDots } from '@/utils/formatNumberWithDots'
import { ElMessage, ElMessageBox } from 'element-plus'

interface SavedVoucher {
  id: string
  voucher_id: string
  used_count: number
  voucher: {
    id: string
    code: string
    name: string
    description: string
    discount_type: 'FIXED' | 'PERCENTAGE'
    discount_value: number
    max_discount_value: number
    min_order_value: number
    end_time: string
    status: string
    usage_limit_per_user: number
    apply_scope: 'ALL' | 'CATEGORY'
  }
}

const vouchers = ref<SavedVoucher[]>([])
const loading = ref(false)
const activeTab = ref('all') // 'all' | 'available'

const fetchSavedVouchers = async () => {
  loading.value = true
  try {
    const response = await axios.get(`${import.meta.env.VITE_BE_API_URL}/product/saved-vouchers`, {
      headers: {
        Authorization: `Bearer ${localStorage.getItem('access_token')}`,
      },
    })
    if (response.data) {
      vouchers.value = response.data
    }
  } catch (error) {
    console.error('Failed to fetch saved vouchers:', error)
  } finally {
    loading.value = false
  }
}

const filteredVouchers = computed(() => {
  const now = new Date()
  return vouchers.value.filter((item) => {
    const isExpired = new Date(item.voucher.end_time) <= now
    const isFullyUsed = item.used_count >= item.voucher.usage_limit_per_user

    if (activeTab.value === 'all') {
      return true
    } else {
      // available/usable
      return !isExpired && !isFullyUsed
    }
  })
})

const unsaveVoucher = async (voucherId: string) => {
  try {
    await ElMessageBox.confirm('Are you sure you want to remove this voucher?', 'Warning', {
      confirmButtonText: 'Remove',
      cancelButtonText: 'Cancel',
      type: 'warning',
    })

    const response = await axios.delete(
      `${import.meta.env.VITE_BE_API_URL}/product/saved-vouchers/${voucherId}`,
      {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('access_token')}`,
        },
      },
    )
    if (response.status === 200) {
      ElMessage.success('Voucher unsaved successfully')
      fetchSavedVouchers() // Refresh list
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to unsave voucher:', error)
      ElMessage.error('Failed to unsave voucher')
    }
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
  return new Date(dateString).toLocaleDateString('en-US', {
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

    <!-- Filter Tabs -->
    <div class="tabs">
      <div class="tab-item" :class="{ active: activeTab === 'all' }" @click="activeTab = 'all'">
        All
      </div>
      <div
        class="tab-item"
        :class="{ active: activeTab === 'available' }"
        @click="activeTab = 'available'"
      >
        Usable
      </div>
    </div>

    <div v-loading="loading" class="voucher-list">
      <div v-if="filteredVouchers.length === 0 && !loading" class="no-vouchers">
        <el-empty description="No vouchers found in this category" />
      </div>

      <div
        v-for="item in filteredVouchers"
        :key="item.id"
        class="voucher-card"
        :class="{
          'is-used':
            item.used_count >= item.voucher.usage_limit_per_user ||
            new Date(item.voucher.end_time) <= new Date(),
        }"
      >
        <div class="voucher-main-content">
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
              <div class="apply-scope">
                <el-tag
                  size="small"
                  :type="item.voucher.apply_scope === 'ALL' ? 'success' : 'warning'"
                >
                  {{ item.voucher.apply_scope === 'ALL' ? 'All Products' : 'Specific Categories' }}
                </el-tag>
              </div>
              <div class="usage-info">
                Used: {{ item.used_count }} / {{ item.voucher.usage_limit_per_user }}
                <span v-if="item.used_count >= item.voucher.usage_limit_per_user" class="used-badge"
                  >(Fully Used)</span
                >
                <span v-else-if="new Date(item.voucher.end_time) <= new Date()" class="used-badge"
                  >(Expired)</span
                >
              </div>
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

        <!-- Hover Remove Overlay -->
        <div class="remove-overlay" @click.stop="unsaveVoucher(item.voucher_id)">
          <el-icon :size="20"><Delete /></el-icon>
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

.tabs {
  display: flex;
  gap: 20px;
  margin-bottom: 20px;
  border-bottom: 2px solid #f0f0f0;
}

.tab-item {
  padding: 10px 5px;
  cursor: pointer;
  font-weight: 500;
  color: #666;
  position: relative;
  transition: all 0.2s;
}

.tab-item:hover {
  color: var(--main-color);
}

.tab-item.active {
  color: var(--main-color);
  font-weight: 700;
}

.tab-item.active::after {
  content: '';
  position: absolute;
  bottom: -2px;
  left: 0;
  width: 100%;
  height: 2px;
  background-color: var(--main-color);
}

.voucher-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.voucher-card {
  position: relative;
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

.voucher-main-content {
  display: flex;
  width: 100%;
  flex: 1;
}

.voucher-card:hover .voucher-right {
  transform: translateX(-50px);
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
  margin: 0 0 5px;
  display: -webkit-box;
  line-clamp: 1;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.apply-scope {
  margin-bottom: 5px;
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
  transition: transform 0.3s ease;
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

.remove-overlay {
  position: absolute;
  top: 0;
  right: -50px; /* Hidden by default */
  width: 50px;
  height: 100%;
  background-color: #ef4444;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  cursor: pointer;
  transition: all 0.3s ease;
  z-index: 5;
}

.voucher-card:hover .remove-overlay {
  right: 0;
}

.remove-overlay:hover {
  background-color: #dc2626;
}

.voucher-card.is-used .voucher-main-content {
  filter: grayscale(1);
  opacity: 0.7;
}

.voucher-card.is-used:hover {
  box-shadow: none;
}

.usage-info {
  font-size: 13px;
  color: #666;
  margin-top: 5px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.used-badge {
  color: #ef4444;
  font-weight: bold;
  font-size: 12px;
}
</style>
