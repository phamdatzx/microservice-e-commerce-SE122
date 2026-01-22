<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ArrowDown } from '@element-plus/icons-vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import IconUser from '@/components/icons/IconUser.vue'
import IconAddress from '@/components/icons/IconAddress.vue'
import IconLock from '@/components/icons/IconLock.vue'
import IconLogout from '@/components/icons/IconLogout.vue'
import { eventBus } from '@/utils/eventBus'

const props = defineProps({
  showAddress: {
    type: Boolean,
    default: true,
  },
  role: {
    type: String,
    required: true,
    validator: (value: string) => ['admin', 'seller'].includes(value),
  },
})

const router = useRouter()
const userName = ref('User')
const avatarUrl = ref('https://placehold.co/40x40')

const fetchUserInfo = async () => {
  const token = localStorage.getItem('access_token')
  if (!token) return

  try {
    const response = await axios.get(`${import.meta.env.VITE_BE_API_URL}/user/my-info`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
    const userData = response.data.data
    userName.value = userData.name || 'User'
    if (userData.image) {
      avatarUrl.value = userData.image
    }
  } catch (error) {
    console.error('Failed to fetch user info:', error)
  }
}

onMounted(() => {
  fetchUserInfo()

  eventBus.on('update_user_info', () => {
    fetchUserInfo()
  })
})

const handleLogout = () => {
  localStorage.removeItem('access_token')
  router.push('/login')
}
</script>

<template>
  <el-dropdown trigger="hover" class="user-dropdown-trigger">
    <div class="user-section">
      <img :src="avatarUrl" alt="User Avatar" class="user-avatar" style="border-radius: 100%" />
      <span>{{ userName }}</span>
      <el-icon><ArrowDown /></el-icon>
    </div>
    <template #dropdown>
      <el-dropdown-menu class="user-dropdown-menu">
        <el-dropdown-item
          class="dropdown-item"
          @click="$router.push(`/${role}/profile?tab=account-info`)"
        >
          <IconUser />
          <span>My Profile</span>
        </el-dropdown-item>
        <el-dropdown-item
          v-if="showAddress"
          class="dropdown-item"
          @click="$router.push(`/${role}/profile?tab=my-address`)"
        >
          <IconAddress />
          <span>My Address</span>
        </el-dropdown-item>
        <el-dropdown-item
          class="dropdown-item"
          @click="$router.push(`/${role}/profile?tab=change-password`)"
        >
          <IconLock />
          <span>Change Password</span>
        </el-dropdown-item>
        <el-dropdown-item divided class="dropdown-item logout-item" @click="handleLogout">
          <IconLogout />
          <span>Log out</span>
        </el-dropdown-item>
      </el-dropdown-menu>
    </template>
  </el-dropdown>
</template>

<style scoped>
/* User section */
.user-dropdown-trigger {
  margin-left: auto;
  outline: none;
}

.user-section {
  display: flex;
  align-items: center;
  justify-content: space-around;
  gap: 12px;
  cursor: pointer;
  padding: 12px 20px;
  outline: none;
  transition: background-color 0.2s;
  min-width: 160px;
}

.user-section:hover {
  background-color: #f5f5f5;
}

.user-avatar {
  width: 40px;
  height: 40px;
  object-fit: cover;
}

/* User Dropdown Styling */
:deep(.user-dropdown-menu) {
  padding: 8px;
  border-radius: 12px;
  border: 1px solid #f4f4f5;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
  min-width: 160px;
}

:deep(.dropdown-item) {
  display: flex !important;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  color: #3f3f46;
  font-size: 14px;
  font-weight: 500;
  border-radius: 8px;
  transition: all 0.2s ease;
  line-height: 1.2;
}

:deep(.dropdown-item:hover) {
  background-color: #f0fdf4 !important;
  color: #22c55e !important;
}

:deep(.dropdown-item svg) {
  color: #a1a1aa;
  transition: color 0.2s ease;
}

:deep(.dropdown-item:hover svg) {
  color: #22c55e;
}

:deep(.logout-item:hover) {
  background-color: #fef2f2 !important;
  color: #ef4444 !important;
}

:deep(.logout-item:hover svg) {
  color: #ef4444 !important;
}
</style>
