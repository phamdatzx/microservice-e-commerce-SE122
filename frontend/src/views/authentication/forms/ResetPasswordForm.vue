<script setup lang="ts">
import './assets/formStyle.css'
import { ref } from 'vue'
import axios from 'axios'
import { ElLoading, ElNotification } from 'element-plus'
import { useRoute } from 'vue-router'
import { handlePasswordToggle } from '@/utils/handlePasswordToggle'
import PasswordToggleBtn from '../components/PasswordToggleBtn.vue'

const USER_API_URL = import.meta.env.VITE_USER_API_URL

const route = useRoute()
const token = route.query.token

const newPassword = ref('')
const resetSuccess = ref(false)

const emits = defineEmits(['reset-success'])

const handleResetPasswordFormSubmit = () => {
  const loading = ElLoading.service({
    lock: true,
    text: 'Resetting your password',
    background: 'rgba(0, 0, 0, 0.7)',
  })

  axios
    .post(`${USER_API_URL}/reset-password?token=${token}`, {
      password: newPassword.value,
    })
    .then((response) => {
      if (response.data.status === 200) {
        ElNotification({
          title: 'Reset successful!',
          message: 'Your password has been reset successfully.',
          type: 'success',
        })
        resetSuccess.value = true
        emits('reset-success')
      } else {
        ElNotification({
          title: 'Reset failed!',
          message: 'Unable to reset your password.',
          type: 'error',
        })
      }
    })
    .catch((error) => {
      ElNotification({
        title: 'Reset failed!',
        message: 'Error: ' + error.response.data.message,
        type: 'error',
      })
    })
    .finally(() => {
      loading.close()
    })
}

// EXPOSES
defineExpose({
  handleResetPasswordFormSubmit,
})
</script>

<template>
  <p v-if="resetSuccess" class="info-text" style="margin-bottom: 140px">
    Your password has been reset successfully. You can now log in with your new password.
  </p>
  <div v-else class="form-group">
    <label for="password" class="form-label">New password</label>
    <div style="position: relative">
      <input
        type="password"
        name="password"
        class="form-input"
        placeholder="••••"
        v-model="newPassword"
        required
      />
      <PasswordToggleBtn @click="handlePasswordToggle" />
    </div>
  </div>
</template>
