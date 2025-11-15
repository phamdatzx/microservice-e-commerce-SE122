<script setup lang="ts">
import PasswordToggleBtn from './PasswordToggleBtn.vue'
import { handlePasswordToggle } from '@/utils/handlePasswordToggle'
import './assets/formStyle.css'
import { ref } from 'vue'
import axios from 'axios'
import { ElLoading, ElNotification } from 'element-plus'

const USER_API_URL = import.meta.env.VITE_USER_API_URL

const name = ref('')
const username = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')

const handleRegisterFormSubmit = () => {
  const loading = ElLoading.service({
    lock: true,
    text: 'Registering',
    background: 'rgba(0, 0, 0, 0.7)',
  })

  axios
    .post(`${USER_API_URL}/register`, {
      name: name.value,
      username: username.value,
      email: email.value,
      password: password.value,
      role: 'customer',
    })
    .then((response) => {
      if (response.data.status === 200) {
        ElNotification({
          title: 'Register successful!',
          message: 'Please check your email to activate your account.',
          type: 'success',
        })
      } else {
        ElNotification({
          title: 'Register failed!',
          message: 'Unable to register your account.',
          type: 'error',
        })
      }
    })
    .catch((error) => {
      ElNotification({
        title: 'Register failed!',
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
  handleRegisterFormSubmit,
})
</script>

<template>
  <div class="form-group">
    <label class="form-label" for="name">Your name</label>
    <input
      class="form-input"
      type="text"
      name="name"
      placeholder="Jhon Deo"
      required
      v-model="name"
    />
  </div>

  <div class="form-group">
    <label class="form-label" for="username">Username</label>
    <input
      class="form-input"
      type="text"
      name="username"
      v-model="username"
      placeholder="Username"
      required
    />
  </div>

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
    <label class="form-label" for="password">Password</label>
    <div class="input-wrapper" style="position: relative">
      <input
        class="form-input"
        type="password"
        id="password"
        name="password"
        placeholder="***"
        v-model="password"
        required
      />

      <PasswordToggleBtn @click="handlePasswordToggle" />
    </div>
  </div>

  <div class="form-group">
    <label class="form-label" for="confirmPassword">Confirm Password</label>
    <div class="input-wrapper" style="position: relative">
      <input
        class="form-input"
        type="password"
        id="confirmPassword"
        name="confirmPassword"
        v-model="confirmPassword"
        placeholder="***"
        required
      />

      <PasswordToggleBtn @click="handlePasswordToggle" />
    </div>
  </div>
</template>
