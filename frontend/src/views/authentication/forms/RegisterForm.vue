<script setup lang="ts">
import { handlePasswordToggle } from '@/utils/handlePasswordToggle'
import axios from 'axios'
import { ElLoading, ElNotification } from 'element-plus'
import { ref } from 'vue'
import PasswordToggleBtn from '../components/PasswordToggleBtn.vue'
import './assets/formStyle.css'

const USER_API_URL = import.meta.env.VITE_BE_API_URL + '/user/public'

const name = ref('')
const username = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const registerSuccess = ref(false)

const emits = defineEmits(['success'])

const handleFormSent = () => {
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
        registerSuccess.value = true
        emits('success')
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

defineExpose({
  handleFormSent,
})
</script>

<template>
  <p v-if="registerSuccess" class="info-text" style="margin-bottom: 140px">
    Please check your email to activate your account.
  </p>
  <div v-else>
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
  </div>
</template>
