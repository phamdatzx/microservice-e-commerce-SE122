<script setup lang="ts">
import { onMounted, ref, useTemplateRef } from 'vue'
import PasswordToggleBtn from './PasswordToggleBtn.vue'
import { handlePasswordToggle } from '@/utils/handlePasswordToggle'
import './assets/formStyle.css'
import axios from 'axios'

const BE_URL = 'http://localhost:8080/api/user/login'

const username = ref('')
const password = ref('')

const handleLoginFormSubmit = () => {
  axios
    .post(BE_URL, {
      username: username.value,
      password: password.value,
    })
    .then((response) => {
      if (response.data.status === 200) {
        alert('Login successful: ' + response.data)
      } else {
        alert('Login failed: ' + response.data)
      }
    })
    .catch((error) => {
      alert('Login failed: ' + error)
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
