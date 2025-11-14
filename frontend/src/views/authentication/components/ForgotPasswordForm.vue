<script setup lang="ts">
import './assets/formStyle.css'
import { ref } from 'vue'
import axios from 'axios'

const BE_URL = 'http://localhost:8080/api/user/send-reset-password'

const username = ref('')

const handleForgotPasswordFormSubmit = () => {
  axios
    .post(BE_URL, {
      username: username.value,
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
  handleForgotPasswordFormSubmit,
})
</script>

<template>
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
</template>
