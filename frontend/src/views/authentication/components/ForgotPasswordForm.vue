<script setup lang="ts">
import { isValidEmail } from '@/utils/isValidEmail'
import './assets/formStyle.css'
import { ref } from 'vue'

const email = ref('')

const handleForgotPasswordFormSubmit = () => {
  // Basic validation
  if (!email.value) {
    alert('Please fill in all fields')
    return
  }

  if (!isValidEmail(email.value)) {
    alert('Please enter a valid email address')
    return
  }

  // // Simulate login process
  // const originalText = loginBtn.value.textContent

  // loginBtn.value.textContent = 'LOGGING IN...'
  // loginBtn.value.disabled = true

  setTimeout(() => {
    alert('Forgot password functionality will be done here!')
    // loginBtn.value.textContent = originalText
    // loginBtn.value.disabled = false
  }, 500)
}

// EXPOSES
defineExpose({
  handleForgotPasswordFormSubmit,
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
</template>
