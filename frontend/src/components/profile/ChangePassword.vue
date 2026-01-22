<script setup lang="ts">
import { ref, computed } from 'vue'
import axios from 'axios'
import { ElNotification } from 'element-plus'

const currentPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')
const isLoading = ref(false)

const touched = ref({
  current: false,
  new: false,
  confirm: false,
})

const handleBlur = (field: 'current' | 'new' | 'confirm') => {
  touched.value[field] = true
}

const isValid = computed(() => {
  return (
    currentPassword.value &&
    newPassword.value &&
    confirmPassword.value &&
    newPassword.value === confirmPassword.value
  )
})

const handleSave = async () => {
  if (!currentPassword.value || !newPassword.value || !confirmPassword.value) {
    ElNotification({
      title: 'Error',
      message: 'Please fill in all fields',
      type: 'error',
    })
    return
  }

  if (newPassword.value !== confirmPassword.value) {
    ElNotification({
      title: 'Error',
      message: 'Confirm password does not match',
      type: 'error',
    })
    return
  }

  const token = localStorage.getItem('access_token')
  if (!token) return

  isLoading.value = true

  try {
    // 1. Get username first (needed for login check)
    const infoResponse = await axios.get(`${import.meta.env.VITE_BE_API_URL}/user/my-info`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
    const username = infoResponse.data.data.username

    // 2. Verify current password by attempting login
    try {
      await axios.post(`${import.meta.env.VITE_BE_API_URL}/user/public/login`, {
        username: username,
        password: currentPassword.value,
      })
    } catch (e) {
      throw new Error('Current password is incorrect')
    }

    // 3. Change password using reset-password endpoint with access token
    await axios.post(
      `${import.meta.env.VITE_BE_API_URL}/user/public/reset-password`,
      {
        password: newPassword.value,
      },
      {
        params: {
          token: token,
        },
      },
    )

    ElNotification({
      title: 'Success',
      message: 'Password changed successfully',
      type: 'success',
    })

    // Reset form
    currentPassword.value = ''
    newPassword.value = ''
    confirmPassword.value = ''
    touched.value = {
      current: false,
      new: false,
      confirm: false,
    }
  } catch (error: any) {
    console.error('Failed to change password:', error)
    ElNotification({
      title: 'Error',
      message: error.message || error.response?.data?.message || 'Failed to change password',
      type: 'error',
    })
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="change-password">
    <h2 class="section-title">Change Password</h2>

    <el-form label-position="top">
      <el-form-item
        label="Current Password"
        required
        :error="touched.current && !currentPassword ? 'Current password is required' : ''"
      >
        <el-input
          v-model="currentPassword"
          type="password"
          size="large"
          placeholder="********"
          show-password
          @blur="handleBlur('current')"
        />
      </el-form-item>

      <el-form-item
        label="New Password"
        required
        :error="touched.new && !newPassword ? 'New password is required' : ''"
      >
        <el-input
          v-model="newPassword"
          type="password"
          size="large"
          placeholder="Set a new password"
          show-password
          @blur="handleBlur('new')"
        />
      </el-form-item>

      <el-form-item
        label="Confirm New Password"
        required
        :error="
          touched.confirm
            ? !confirmPassword
              ? 'Confirm password is required'
              : newPassword !== confirmPassword
                ? 'Passwords do not match'
                : ''
            : ''
        "
      >
        <el-input
          v-model="confirmPassword"
          type="password"
          size="large"
          placeholder="Re-enter your new password"
          show-password
          @blur="handleBlur('confirm')"
        />
      </el-form-item>

      <el-form-item>
        <el-button
          type="primary"
          size="large"
          @click="handleSave"
          :loading="isLoading"
          :disabled="!isValid"
          class="save-btn"
          style="padding: 22px 40px; font-size: 16px; margin-top: 4px"
          >SAVE</el-button
        >
      </el-form-item>
    </el-form>
  </div>
</template>

<style scoped>
.section-title {
  font-size: 28px;
  font-weight: 700;
  margin-bottom: 30px;
  color: #000;
}

:deep(.el-form-item__label) {
  font-weight: 600;
  color: #333;
  margin-bottom: 8px !important;
}

:deep(.el-input__wrapper) {
  padding: 8px 15px;
  border-radius: 8px;
  box-shadow: 0 0 0 1px #dcdfe6 inset;
}

:deep(.el-input__inner) {
  height: 48px;
}

.save-btn {
  --el-button-disabled-bg-color: #dcdfe6;
  --el-button-disabled-border-color: #dcdfe6;
}
</style>
