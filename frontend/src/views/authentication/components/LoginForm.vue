<script setup lang="ts">
import { onMounted, ref, useTemplateRef } from 'vue'
import PasswordToggleBtn from './PasswordToggleBtn.vue'
import { handlePasswordToggle } from '@/utils/handlePasswordToggle'
import './assets/formStyle.css'
import axios from 'axios'
import { ElLoading } from 'element-plus'

const USER_API_URL = import.meta.env.VITE_USER_API_URL

const username = ref('')
const password = ref('')

const handleLoginFormSubmit = () => {
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
        console.log('Login successful!')
        localStorage.setItem('access_token', loginRes.data.data.access_token)
        //       // window.location.href = '/'
      } else {
        console.log('Login failed: ' + loginRes.data)
      }
    })
    .catch((error) => {
      alert('Login failed: ' + error.response.data.message)
    })
    .finally(() => {
      loading.close()
    })
}

// EXPOSES
defineExpose({
  handleLoginFormSubmit,
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

  <div style="text-align: right; margin-bottom: 32px">
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
