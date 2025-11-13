<script setup lang="ts">
import { useTemplateRef } from 'vue'
import '../assets/style.css'
import PasswordToggleBtn from '../components/PasswordToggleBtn.vue'
import { handlePasswordToggle } from '@/utils/handlePasswordToggle'

const password = useTemplateRef('password')
const email = useTemplateRef('email')

// Email validation helper
function isValidEmail(email: string) {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return emailRegex.test(email)
}

const handleLoginFormSubmit = (e: Event) => {
  if (!email.value || !password.value) return

  e.preventDefault()

  const emailValue = email.value.value
  const passwordValue = password.value.value

  // Basic validation
  if (!emailValue || !passwordValue) {
    alert('Please fill in all fields')
    return
  }

  if (!isValidEmail(emailValue)) {
    alert('Please enter a valid email address')
    return
  }

  // // Simulate login process
  // const originalText = loginBtn.value.textContent

  // loginBtn.value.textContent = 'LOGGING IN...'
  // loginBtn.value.disabled = true

  // setTimeout(() => {
  //   if (!loginBtn.value) return

  //   alert('Login functionality will be done here!')
  //   loginBtn.value.textContent = originalText
  //   loginBtn.value.disabled = false
  // }, 1500)
}
</script>

<template>
  <div class="form-group">
    <label for="email" class="form-label">Email Address</label>
    <input
      type="email"
      id="email"
      ref="email"
      name="email"
      class="form-input"
      placeholder="Example@gmail.com"
      required
    />
  </div>

  <div class="form-group">
    <label for="password" class="form-label">Password</label>
    <div style="position: relative">
      <input
        type="password"
        ref="password"
        id="password"
        name="password"
        class="form-input"
        placeholder="••••"
        required
      />
      <PasswordToggleBtn @click="handlePasswordToggle" />
    </div>
  </div>

  <div style="text-align: right; margin-bottom: 32px">
    <a href="#" class="forgot-link">Forget Password ?</a>
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
