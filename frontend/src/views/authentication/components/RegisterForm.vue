<script setup lang="ts">
import PasswordToggleBtn from './PasswordToggleBtn.vue'
import { handlePasswordToggle } from '@/utils/handlePasswordToggle'
import './assets/formStyle.css'
import { ref } from 'vue'
import axios from 'axios'

const BE_URL = 'http://localhost:8080/api/user/register'

const name = ref('')
const username = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')

const handleRegisterFormSubmit = () => {
  axios
    .post(BE_URL, {
      name: name.value,
      username: username.value,
      email: email.value,
      password: password.value,
      role: 'customer',
    })
    .then((response) => {
      if (response.data.status === 200) {
        alert('Register successful: ' + response.data)
      } else {
        alert('Register failed: ' + response.data)
      }
    })
    .catch((error) => {
      alert('Login failed: ' + error)
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
