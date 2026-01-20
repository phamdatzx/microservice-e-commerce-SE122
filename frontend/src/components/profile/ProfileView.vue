<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Loading, ArrowRight } from '@element-plus/icons-vue'
import axios from 'axios'
import AccountInfo from '@/components/profile/AccountInfo.vue'
import ChangePassword from '@/components/profile/ChangePassword.vue'
import MyAddress from '@/components/profile/MyAddress.vue'

const props = defineProps<{
  role: 'admin' | 'seller' | 'customer'
}>()

const router = useRouter()
const route = useRoute()

const name = ref('')
const email = ref('')
const phone = ref('')
const avatar = ref('')
const isLoading = ref(false)
const activeTab = ref('account-info')

const fetchUserInfo = async () => {
  const token = localStorage.getItem('access_token')
  if (!token) return

  isLoading.value = true
  try {
    const response = await axios.get(`${import.meta.env.VITE_BE_API_URL}/user/my-info`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })

    const userData = response.data.data
    name.value = userData.name || ''
    email.value = userData.email || ''
    phone.value = userData.phone || ''
    if (userData.image) {
      avatar.value = userData.image
    }
  } catch (error) {
    console.error('Failed to fetch user info:', error)
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchUserInfo()
  if (props.role !== 'customer' && route.query.tab) {
    activeTab.value = route.query.tab as string
  }
})

watch(
  () => route.query.tab,
  (newTab) => {
    if (props.role !== 'customer' && newTab) {
      activeTab.value = newTab as string
    }
  },
)

const menuItems = computed(() => {
  if (props.role === 'admin') {
    return [
      { id: 'account-info', label: 'Account info' },
      { id: 'change-password', label: 'Change password' },
    ]
  } else if (props.role === 'seller') {
    return [
      { id: 'account-info', label: 'Account info' },
      { id: 'my-address', label: 'My address' },
      { id: 'change-password', label: 'Change password' },
    ]
  } else {
    // Customer
    return [
      { id: 'account-info', label: 'Account info', path: '/profile/account-info' },
      { id: 'my-order', label: 'My order', path: '/profile/orders' },
      { id: 'my-voucher', label: 'My vouchers', path: '/profile/vouchers' },
      { id: 'my-address', label: 'My address', path: '/profile/address' },
      { id: 'change-password', label: 'Change password', path: '/profile/change-password' },
    ]
  }
})

const activeItemId = computed(() => {
  if (props.role === 'customer') {
    const active = menuItems.value.find((item) => route.path.includes(item.path!))
    return active ? active.id : 'account-info'
  }
  return activeTab.value
})

const handleMenuClick = (item: any) => {
  if (props.role === 'customer') {
    router.push(item.path)
  } else {
    activeTab.value = item.id
    router.push({ query: { ...route.query, tab: item.id } })
  }
}

const userInfo = computed(() => ({
  name: name.value,
  email: email.value,
  avatar: avatar.value,
}))
</script>

<template>
  <div class="main-container" :style="role !== 'customer' ? { marginTop: 0, marginBottom: 0 } : {}">
    <el-row :gutter="28">
      <!-- Sidebar -->
      <el-col :span="6">
        <div class="sidebar box-shadow border-radius">
          <div class="user-profile">
            <div v-if="isLoading" class="loading-state">
              <el-icon class="is-loading"><Loading /></el-icon>
            </div>
            <div v-else>
              <div class="avatar-wrapper">
                <img :src="avatar || 'https://placehold.co/140x140'" alt="User Avatar" />
              </div>
              <h3 class="user-name">{{ name }}</h3>
              <p class="user-email">{{ email }}</p>
            </div>
          </div>

          <div class="sidebar-menu">
            <div
              v-for="item in menuItems"
              :key="item.id"
              class="menu-item"
              :class="{ active: activeItemId === item.id }"
              @click="handleMenuClick(item)"
            >
              <span>{{ item.label }}</span>
              <el-icon><ArrowRight /></el-icon>
            </div>
          </div>
        </div>
      </el-col>

      <!-- Main Content -->
      <el-col :span="18">
        <div v-if="role !== 'customer'" class="content-wrapper box-shadow border-radius">
          <AccountInfo
            v-if="activeTab === 'account-info'"
            v-model:name="name"
            v-model:email="email"
            v-model:phone="phone"
            v-model:avatar="avatar"
          />
          <ChangePassword v-else-if="activeTab === 'change-password'" />
          <MyAddress v-else-if="activeTab === 'my-address' && role === 'seller'" />
          <div v-else>
            <h2 class="section-title">Coming Soon</h2>
            <div style="padding: 40px; text-align: center; color: #999">
              This section is under development.
            </div>
          </div>
        </div>

        <router-view v-else v-slot="{ Component }">
          <component
            :is="Component"
            v-model:name="name"
            v-model:email="email"
            v-model:phone="phone"
            v-model:avatar="avatar"
          />
        </router-view>
      </el-col>
    </el-row>
  </div>
</template>

<style scoped>
.main-container {
  margin-top: 20px;
  margin-bottom: 20px;
}

.sidebar {
  background-color: #fff;
  padding-bottom: 20px;
  overflow: hidden;
}

.user-profile {
  padding: 30px 20px;
  text-align: center;
  border-bottom: 1px solid #f0f0f0;
  background-color: #f9f9f9;
}

.loading-state {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 214px;
  font-size: 24px;
  color: #888;
}

.avatar-wrapper {
  width: 140px;
  height: 140px;
  margin: 0 auto 16px;
  border-radius: 12px;
  overflow: hidden;
  background-color: #e0e0e0;
}

.avatar-wrapper img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.user-name {
  font-size: 20px;
  font-weight: 700;
  color: #000;
  margin-bottom: 4px;
}

.user-email {
  font-size: 14px;
  color: #888;
}

.sidebar-menu {
  padding: 10px 15px;
}

.menu-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  margin: 8px 0;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
  font-weight: 500;
  color: #444;
}

.menu-item:hover {
  background-color: #f5f5f5;
  color: var(--main-color);
}

.menu-item.active {
  background-color: var(--main-color);
  color: #fff;
}

.content-wrapper {
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
  min-height: 500px;
}

.section-title {
  font-size: 28px;
  font-weight: 700;
  margin-bottom: 30px;
  color: #000;
}
</style>
