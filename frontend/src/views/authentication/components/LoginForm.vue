<script setup lang="ts">
import { onMounted, ref, useTemplateRef } from 'vue'
import PasswordToggleBtn from './PasswordToggleBtn.vue'
import { handlePasswordToggle } from '@/utils/handlePasswordToggle'
import './assets/formStyle.css'
import { isValidEmail } from '@/utils/isValidEmail'
import axios from 'axios'

const BE_URL = 'http://localhost:8080/api/user/login'

const email = ref('')
const password = ref('')

const handleLoginFormSubmit = () => {
  if (!isValidEmail(email.value)) {
    alert('Please enter a valid email address')
    return
  }

  axios
    .post(BE_URL, {
      email: email.value,
      password: password.value,
    })
    .then((response) => {
      alert('Login successful: ' + response.data)
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
    <label for="email" class="form-label">Email Address</label>
    <input
      type="email"
      id="email"
      name="email"
      class="form-input"
      placeholder="Example@gmail.com"
      v-model="email"
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
