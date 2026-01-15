<script setup lang="ts">
import { ref, computed } from 'vue'
import { formatNumberWithDots } from '@/utils/formatNumberWithDots'

interface Voucher {
  id: string
  code: string
  name: string
  description: string
  discount_type: 'FIXED' | 'PERCENTAGE'
  discount_value: number
  max_discount_value: number | null
  min_order_value: number
  total_quantity: number
  used_quantity: number
  start_time: string
  end_time: string
}

const props = defineProps<{
  vouchers: Voucher[]
  savedVoucherIds: string[]
}>()

const emit = defineEmits(['save'])

const isSaved = (voucherId: string) => {
  return props.savedVoucherIds.includes(voucherId)
}

const activeVouchers = computed(() => {
  const now = new Date()
  return props.vouchers.filter((v) => new Date(v.end_time) > now)
})

const scrollContainer = ref<HTMLElement | null>(null)

const formatDiscount = (v: Voucher) => {
  if (v.discount_type === 'FIXED') {
    return `Save ${formatNumberWithDots(v.discount_value)}₫`
  }
  return `Save ${v.discount_value}%`
}

const formatExpiry = (endTime: string) => {
  const date = new Date(endTime)
  return `Exp: ${date.toLocaleDateString('vi-VN')}`
}

const getUsageProgress = (v: Voucher) => {
  if (v.total_quantity === 0) return 100
  return Math.round((v.used_quantity / v.total_quantity) * 100)
}

const scrollNext = () => {
  if (scrollContainer.value) {
    scrollContainer.value.scrollBy({ left: 300, behavior: 'smooth' })
  }
}

const scrollPrev = () => {
  if (scrollContainer.value) {
    scrollContainer.value.scrollBy({ left: -300, behavior: 'smooth' })
  }
}
</script>

<template>
  <div class="seller-vouchers">
    <div class="container">
      <div class="vouchers-relative-container">
        <template v-if="activeVouchers.length > 0">
          <div class="vouchers-wrapper" ref="scrollContainer">
            <div v-for="voucher in activeVouchers" :key="voucher.id" class="voucher-card">
              <div class="scalloped-edge"></div>
              <div class="voucher-left">
                <div class="voucher-info">
                  <div class="voucher-discount">{{ formatDiscount(voucher) }}</div>
                  <div class="voucher-min-spend">
                    Min. Spend {{ formatNumberWithDots(voucher.min_order_value) }}₫
                  </div>
                </div>
                <div class="voucher-expiry">{{ formatExpiry(voucher.end_time) }}</div>
                <div class="voucher-progress-container">
                  <div class="progress-bar">
                    <div
                      class="progress-fill"
                      :style="{ width: getUsageProgress(voucher) + '%' }"
                    ></div>
                  </div>
                  <span class="progress-text">Used {{ getUsageProgress(voucher) }}%</span>
                </div>
              </div>
              <div class="voucher-right">
                <button
                  class="save-btn"
                  :class="{ 'is-saved': isSaved(voucher.id) }"
                  :disabled="isSaved(voucher.id)"
                  @click="emit('save', voucher.id)"
                >
                  {{ isSaved(voucher.id) ? 'Saved' : 'Save' }}
                </button>
              </div>
              <!-- Perforated line effect -->
              <div class="perforation"></div>
            </div>
          </div>

          <!-- Navigation Buttons - only show if there are enough vouchers to scroll -->
          <template v-if="activeVouchers.length > 4">
            <button class="nav-btn prev-btn" @click="scrollPrev">
              <svg viewBox="0 0 24 24" width="28" height="28">
                <path
                  fill="currentColor"
                  d="M15.41 7.41L10.83 12l4.58 4.59L14 18l-6-6 6-6 1.41 1.41z"
                />
              </svg>
            </button>
            <button class="nav-btn next-btn" @click="scrollNext">
              <svg viewBox="0 0 24 24" width="28" height="28">
                <path
                  fill="currentColor"
                  d="M8.59 16.59L13.17 12 8.59 7.41 10 6l6 6-6 6-1.41-1.41z"
                />
              </svg>
            </button>
          </template>
        </template>
        <div v-else class="empty-vouchers">
          <div class="no-vouchers-message">
            <svg
              viewBox="0 0 24 24"
              width="48"
              height="48"
              fill="none"
              stroke="currentColor"
              stroke-width="1.5"
              class="empty-icon"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M16.5 6v.75m0 3v.75m0 3v.75m0 3V18m-9-5.25h5.25M7.5 15h3M3.375 5.25c-.621 0-1.125.504-1.125 1.125v11.25c0 .621.504 1.125 1.125 1.125h17.25c.621 0 1.125-.504 1.125-1.125V6.375c0-.621-.504-1.125-1.125-1.125H3.375z"
              />
            </svg>
            <p>No active vouchers available at the moment</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.seller-vouchers {
  padding: 30px 0 10px;
  background-color: #f5f5f5;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 15px;
}

.vouchers-relative-container {
  position: relative;
  background: white;
  border-radius: 4px;
  min-height: 100px;
  display: flex;
  flex-direction: column;
}

.empty-vouchers {
  width: 100%;
  padding: 40px 0;
  display: flex;
  justify-content: center;
  align-items: center;
}

.no-vouchers-message {
  display: flex;
  flex-direction: column;
  align-items: center;
  color: #999;
  gap: 12px;
}

.no-vouchers-message p {
  margin: 0;
  font-size: 14px;
}

.empty-icon {
  color: #dbdbdb;
}

.vouchers-wrapper {
  display: flex;
  gap: 15px;
  padding: 20px;
  overflow-x: auto;
  scrollbar-width: none; /* Hide scrollbar for Firefox */
}

.vouchers-wrapper::-webkit-scrollbar {
  display: none; /* Hide scrollbar for Chrome, Safari, Opera */
}

.voucher-card {
  position: relative;
  display: flex;
  min-width: 280px;
  background-color: #f0f9f4;
  border: 1px solid var(--main-color);
  border-radius: 2px;
  padding-left: 8px;
}

.scalloped-edge {
  position: absolute;
  top: 0;
  bottom: 0;
  left: -4px;
  width: 8px;
  background-image:
    radial-gradient(
      circle at 0px 5px,
      transparent 3px,
      var(--main-color) 2px,
      var(--main-color) 4px,
      transparent 4px
    ),
    radial-gradient(circle at 4px 5px, #fff 3px, transparent 3px);
  background-position:
    4px 0,
    0 0;
  background-size:
    4px 10px,
    8px 10px;
  background-repeat: repeat-y;
  z-index: 1;
}

.voucher-left {
  flex: 1;
  padding: 15px 12px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.voucher-info {
  margin-bottom: 12px;
}

.voucher-discount {
  color: var(--main-color);
  font-size: 18px;
  font-weight: 500;
  margin-bottom: 2px;
}

.voucher-min-spend {
  font-size: 15px;
  color: var(--main-color);
  margin-bottom: 0;
}

.voucher-expiry {
  font-size: 13px;
  color: #757575;
}

.voucher-progress-container {
  margin-top: 8px;
}

.progress-bar {
  height: 4px;
  background-color: #dcf1e7;
  border-radius: 2px;
  overflow: hidden;
  margin-bottom: 4px;
}

.progress-fill {
  height: 100%;
  background-color: var(--main-color);
}

.progress-text {
  font-size: 10px;
  color: var(--main-color);
}

.voucher-right {
  width: 70px;
  background: white;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 10px;
  border-radius: 2px;
}

.save-btn {
  background-color: var(--main-color);
  color: white;
  border: 1px solid transparent;
  padding: 8px 4px;
  border-radius: 2px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  width: 100%;
  text-align: center;
  display: flex;
  align-items: center;
  justify-content: center;
  box-sizing: border-box;
}

.save-btn:hover {
  filter: brightness(0.9);
}

.save-btn.is-saved {
  background-color: #f5f5f5;
  color: #999;
  border: 1px solid #dbdbdb;
  cursor: not-allowed;
}

.save-btn:disabled {
  opacity: 0.8;
}

.perforation {
  position: absolute;
  top: 0;
  right: 69px;
  bottom: 0;
  width: 1px;
  border-right: 1px dashed #c0e6d6;
}

.nav-btn {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: white;
  border: none;
  box-shadow: 0 1px 12px rgba(0, 0, 0, 0.12);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: #757575;
  z-index: 2;
  transition: all 0.2s;
}

.prev-btn {
  left: -20px;
}

.next-btn {
  right: -20px;
}

.nav-btn:hover {
  background: #fafafa;
  color: var(--main-color);
  box-shadow: 0 2px 15px rgba(0, 0, 0, 0.15);
}
</style>
