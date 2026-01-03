<script setup lang="ts">
import { ref } from 'vue'

const scrollContainer = ref<HTMLElement | null>(null)

const vouchers = [
  {
    discount: 'Save 10k₫',
    minSpend: 'Min. Spend 0₫',
    expiry: 'Exp: 31.01.2026',
    progress: null,
  },
  {
    discount: 'Save 10k₫',
    minSpend: 'Min. Spend 49k₫',
    expiry: 'Valid from: 15.01.2026',
    progress: null,
  },
  {
    discount: 'Save 15k₫',
    minSpend: 'Min. Spend 79k₫',
    expiry: 'Exp: 31.01.2026',
    progress: 68,
  },
  {
    discount: 'Save 20k₫',
    minSpend: 'Min. Spend 99k₫',
    expiry: 'Valid from: 15.01.2026',
    progress: null,
  },
]

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
        <div class="vouchers-wrapper" ref="scrollContainer">
          <div v-for="(voucher, index) in vouchers" :key="index" class="voucher-card">
            <div class="scalloped-edge"></div>
            <div class="voucher-left">
              <div class="voucher-info">
                <div class="voucher-discount">{{ voucher.discount }}</div>
                <div class="voucher-min-spend">{{ voucher.minSpend }}</div>
              </div>
              <div class="voucher-expiry">{{ voucher.expiry }}</div>
              <div v-if="voucher.progress !== null" class="voucher-progress-container">
                <div class="progress-bar">
                  <div class="progress-fill" :style="{ width: voucher.progress + '%' }"></div>
                </div>
                <span class="progress-text">Used {{ voucher.progress }}%</span>
              </div>
            </div>
            <div class="voucher-right">
              <button class="save-btn">Save</button>
            </div>
            <!-- Perforated line effect -->
            <div class="perforation"></div>
          </div>
        </div>

        <!-- Navigation Buttons -->
        <button class="nav-btn prev-btn" @click="scrollPrev">
          <svg viewBox="0 0 24 24" width="28"" height="28">
            <path
              fill="currentColor"
              d="M15.41 7.41L10.83 12l4.58 4.59L14 18l-6-6 6-6 1.41 1.41z"
            />
          </svg>
        </button>
        <button class="nav-btn next-btn" @click="scrollNext">
          <svg viewBox="0 0 24 24" width="28"" height="28">
            <path fill="currentColor" d="M8.59 16.59L13.17 12 8.59 7.41 10 6l6 6-6 6-1.41-1.41z" />
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.seller-vouchers {
  padding: 30px 0;
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
  border: none;
  padding: 8px 12px;
  border-radius: 2px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  width: 100%;
}

.save-btn:hover {
  filter: brightness(0.9);
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
