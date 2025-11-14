<script setup lang="ts">
import PasswordToggleBtn from './PasswordToggleBtn.vue'
import { handlePasswordToggle } from '@/utils/handlePasswordToggle'
import './assets/formStyle.css'
import { ref } from 'vue'
import { isValidEmail } from '@/utils/isValidEmail'

const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const name = ref('')

const handleRegisterFormSubmit = () => {
  if (!isValidEmail(email.value)) {
    alert('Please enter a valid email address')
    return
  }

  setTimeout(() => {
    alert('Register functionality will be done here!')
    // loginBtn.value.textContent = originalText
    // loginBtn.value.disabled = false
  }, 500)
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
      id="name"
      name="name"
      placeholder="Jhon Deo"
      required
      v-model="name"
    />
  </div>

  <div class="form-group">
    <label class="form-label" for="email">Email Address</label>
    <input
      class="form-input"
      type="email"
      id="email"
      name="email"
      v-model="email"
      placeholder="Example@gmail.com"
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
