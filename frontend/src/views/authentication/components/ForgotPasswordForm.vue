<script setup lang="ts">
import { isValidEmail } from '@/utils/isValidEmail'
import './assets/formStyle.css'
import { useTemplateRef } from 'vue'

const email = useTemplateRef('email')

const handleLoginFormSubmit = (e: Event) => {
  if (!email.value) return

  e.preventDefault()

  const emailValue = email.value.value

  // Basic validation
  if (!emailValue) {
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
</template>
