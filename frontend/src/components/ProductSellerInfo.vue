<script setup lang="ts">
import { ChatDotRound, Shop, Goods, User, Star, Plus, Check } from '@element-plus/icons-vue'
import { formatNumberWithDots } from '@/utils/formatNumberWithDots'
import { useRouter } from 'vue-router'
import { ref, onMounted } from 'vue'

interface SellerInfo {
  id: string
  name: string
  image: string
  sale_info: {
    follow_count: number
    rating_count: number
    rating_average: number
    product_count: number
    is_following: boolean
  }
}

const props = withDefaults(
  defineProps<{
    sellerInfo: SellerInfo
    layout?: 'card' | 'header'
    showFollowButton?: boolean
    showViewShop?: boolean
    isFollowLoading?: boolean
  }>(),
  {
    layout: 'card',
    showFollowButton: false,
    showViewShop: true,
    isFollowLoading: false,
  },
)

const emit = defineEmits(['toggle-follow'])
const router = useRouter()
const isLoggedIn = ref(false)

onMounted(() => {
  isLoggedIn.value = !!localStorage.getItem('access_token')
})

const openChat = () => {
  if (!isLoggedIn.value) {
    // Optionally redirect to login or show message
    return
  }
  const event = new CustomEvent('open-chat', {
    detail: { sellerId: props.sellerInfo.id },
  })
  window.dispatchEvent(event)
}
</script>

<template>
  <div
    class="seller-info-section"
    :class="{
      'box-shadow border-radius': layout === 'card',
      'layout-header': layout === 'header',
    }"
  >
    <div class="seller-info-content">
      <div class="seller-info-left" :class="{ 'header-info-left': layout === 'header' }">
        <div class="seller-avatar-wrapper">
          <img
            :src="sellerInfo.image || '/src/assets/avatar.jpg'"
            class="seller-avatar"
            :class="{ 'header-avatar': layout === 'header' }"
          />
        </div>
        <div class="seller-main">
          <h3 class="seller-name">{{ sellerInfo.name }}</h3>
          <p class="seller-status">Active 2 minutes ago</p>
          <div class="seller-actions">
            <template v-if="showFollowButton">
              <el-button
                :color="'var(--main-color)'"
                plain
                :icon="sellerInfo.sale_info.is_following ? Check : Plus"
                size="default"
                :loading="isFollowLoading"
                @click="emit('toggle-follow')"
              >
                {{ sellerInfo.sale_info.is_following ? 'Following' : 'Follow' }}
              </el-button>
            </template>
            <el-button
              v-if="isLoggedIn"
              color="var(--main-color)"
              plain
              :icon="ChatDotRound"
              size="default"
              @click="openChat"
              >Chat Now</el-button
            >
            <el-button
              v-if="showViewShop"
              plain
              :icon="Shop"
              size="default"
              @click="router.push(`/seller-page/${sellerInfo.id}`)"
              >View Shop</el-button
            >
          </div>
        </div>
      </div>

      <div class="shop-stats">
        <div class="stat-item">
          <div class="stat-value">
            {{ sellerInfo.sale_info.product_count || 0 }} <el-icon><Goods /></el-icon>
          </div>
          <div class="stat-label">Products</div>
        </div>
        <div class="stat-divider"></div>
        <div class="stat-item">
          <div class="stat-value">
            {{ formatNumberWithDots(sellerInfo.sale_info.follow_count) }}
            <el-icon><User /></el-icon>
          </div>
          <div class="stat-label">Followers</div>
        </div>
        <div class="stat-divider"></div>
        <div class="stat-item">
          <div class="stat-value">
            {{ Number(sellerInfo.sale_info.rating_average || 0).toFixed(1) }}
            <el-icon><Star /></el-icon>
          </div>
          <div class="stat-label">
            Rating
            <span class="rating-count">({{ sellerInfo.sale_info.rating_count }} Ratings)</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="css" scoped>
.seller-info-section {
  background-color: #fff;
  margin-bottom: 20px;
}

.seller-info-section.layout-header {
  background-color: #edf7f1;
  margin-bottom: 0;
  border-radius: 4px;
}

.seller-info-content {
  display: flex;
  padding: 25px;
  gap: 30px;
  align-items: center;
}

.layout-header .seller-info-content {
  padding: 20px;
  width: 100%;
}

.seller-info-left {
  display: flex;
  gap: 20px;
  padding-right: 30px;
  border-right: 1px solid #eee;
  min-width: 380px;
}

.header-info-left {
  border-right: none;
  padding-right: 0;
  min-width: 390px;
}

.seller-avatar-wrapper {
  width: 80px;
  height: 80px;
  position: relative;
  flex-shrink: 0;
}

.seller-avatar {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  object-fit: cover;
  border: 1px solid #eee;
}

.header-avatar {
  border: 4px solid rgba(255, 255, 255, 0.3);
}

.seller-name {
  font-size: 16px;
  margin: 0 0 4px 0;
  font-weight: 500;
  color: #000000cc;
}

.layout-header .seller-name {
  font-size: 20px;
  color: #222;
}

.seller-status {
  font-size: 12px;
  color: #666;
  margin-bottom: 12px;
}

.seller-actions {
  display: flex;
  gap: 10px;
}

.shop-stats {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0;
  padding-left: 20px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  padding: 0 40px;
}

.stat-item:first-child {
  padding-left: 0;
}

.stat-divider {
  width: 1px;
  height: 40px;
  background-color: #eee;
}

.stat-label {
  font-size: 13px;
  color: #999;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-top: 4px;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: var(--main-color);
  line-height: 1;
  display: flex;
  align-items: center;
  gap: 4px;
}

.rating-count {
  font-size: 14px;
  font-weight: 400;
  color: #999;
}
</style>
