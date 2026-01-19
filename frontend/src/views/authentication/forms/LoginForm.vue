<script setup lang="ts">
import { handlePasswordToggle } from '@/utils/handlePasswordToggle'
import axios from 'axios'
import { ElLoading, ElNotification } from 'element-plus'
import { ref } from 'vue'
import PasswordToggleBtn from '../components/PasswordToggleBtn.vue'
import './assets/formStyle.css'

const USER_API_URL = import.meta.env.VITE_BE_API_URL + '/user/public'

const username = ref('')
const password = ref('')

const emits = defineEmits(['success'])

const handleFormSent = () => {
  const loading = ElLoading.service({
    lock: true,
    text: 'Logging in',
    background: 'rgba(0, 0, 0, 0.7)',
  })

  axios
    .post(`${USER_API_URL}/login`, {
      username: username.value,
      password: password.value,
    })
    .then((loginRes) => {
      if (loginRes.data.status === 200) {
        localStorage.setItem('access_token', loginRes.data.data.access_token)
        localStorage.setItem('user_id', loginRes.data.data.user_id)
        emits('success', loginRes.data.data)
      } else {
        ElNotification({
          title: 'Login failed!',
          message: 'Unable to login to your account.',
          type: 'error',
        })
      }
    })
    .catch((error) => {
      ElNotification({
        title: 'Login failed!',
        message: 'Error: ' + error.response.data.message,
        type: 'error',
      })
    })
    .finally(() => {
      loading.close()
    })
}

defineExpose({
  handleFormSent,
})
</script>

<template>
  <div class="form-group">
    <label for="username" class="form-label">Username</label>
    <input
      type="text"
      name="username"
      class="form-input"
      placeholder="Username"
      v-model="username"
      required
    />
  </div>

  <div class="form-group">
    <label for="password" class="form-label">Password</label>
    <div style="position: relative">
      <input
        type="password"
        id="password"
        name="password"
        class="form-input"
        placeholder="••••"
        v-model="password"
        required
      />
      <PasswordToggleBtn @click="handlePasswordToggle" />
    </div>
  </div>

  <div style="text-align: right; margin-bottom: 20px">
    <RouterLink class="forgot-link" to="/forgot-password">Forget Password ?</RouterLink>
  </div>
</template>

<style scoped>
.forgot-link {
  color: #9ca3af;
  text-decoration: none;
  font-size: 14px;

  &:hover {
    color: #22c55e;
  }
}
</style>
