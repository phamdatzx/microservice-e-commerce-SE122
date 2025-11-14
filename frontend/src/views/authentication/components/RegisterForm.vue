<script setup lang="ts">
import PasswordToggleBtn from './PasswordToggleBtn.vue'
import { handlePasswordToggle } from '@/utils/handlePasswordToggle'
import './assets/formStyle.css'
</script>

<template>
  <div class="form-group">
    <label class="form-label" for="name">Your name</label>
    <input class="form-input" type="text" id="name" name="name" placeholder="Jhon Deo" required />
  </div>

  <div class="form-group">
    <label class="form-label" for="email">Email Address</label>
    <input
      class="form-input"
      type="email"
      id="email"
      name="email"
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
        placeholder="***"
        required
      />

      <PasswordToggleBtn @click="handlePasswordToggle" />
    </div>
  </div>
</template>
