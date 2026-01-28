<script setup lang="ts">
import { handlePasswordToggle } from '@/utils/handlePasswordToggle'
import axios from 'axios'
import { ElLoading, ElNotification } from 'element-plus'
import { ref, watch } from 'vue'
import PasswordToggleBtn from '../components/PasswordToggleBtn.vue'
import './assets/formStyle.css'

const USER_API_URL = import.meta.env.VITE_BE_API_URL + '/user/public'

const name = ref('')
const username = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const registerSuccess = ref(false)

const usernameStatus = ref<'available' | 'taken' | 'none'>('none')
const isCheckingUsername = ref(false)
let debounceTimer: any = null

const checkUsername = async (val: string) => {
  if (!val) {
    usernameStatus.value = 'none'
    return
  }

  isCheckingUsername.value = true
  try {
    const response = await axios.post(`${USER_API_URL}/check-username`, {
      username: val,
    })
    // Based on the response structure { data: { exists: boolean } }
    if (response.data.data.exists) {
      usernameStatus.value = 'taken'
    } else {
      usernameStatus.value = 'available'
    }
  } catch (error) {
    console.error('Check username failed:', error)
    usernameStatus.value = 'none'
  } finally {
    isCheckingUsername.value = false
  }
}

watch(username, (newVal) => {
  if (debounceTimer) clearTimeout(debounceTimer)
  usernameStatus.value = 'none'

  if (newVal.length >= 3) {
    debounceTimer = setTimeout(() => {
      checkUsername(newVal)
    }, 500)
  }
})

const emits = defineEmits(['success'])

const handleFormSent = () => {
  if (usernameStatus.value === 'taken') return

  if (isCheckingUsername.value) {
    ElNotification({
      title: 'Validation in progress',
      message: 'Please wait while we check the username.',
      type: 'info',
    })
    return
  }

  const loading = ElLoading.service({
    lock: true,
    text: 'Registering as Seller',
    background: 'rgba(0, 0, 0, 0.7)',
  })

  axios
    .post(`${USER_API_URL}/register`, {
      name: name.value,
      username: username.value,
      email: email.value,
      password: password.value,
      role: 'seller', // Hardcoded role
    })
    .then((response) => {
      if (response.data.status === 200) {
        ElNotification({
          title: 'Register successful!',
          message: 'Please check your email to activate your seller account.',
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
    Please check your email to activate your seller account.
  </p>
  <div v-else class="seller-register-form">
    <div class="form-group">
      <label class="form-label" for="name">Shop/Seller Name</label>
      <input
        class="form-input"
        type="text"
        name="name"
        placeholder="My Awesome Shop"
        required
        v-model="name"
      />
    </div>

    <div class="form-group">
      <label class="form-label" for="username">Username</label>
      <input
        class="form-input"
        :class="{
          'input-error': usernameStatus === 'taken',
          'input-success': usernameStatus === 'available',
        }"
        type="text"
        name="username"
        v-model="username"
        placeholder="Username"
        @blur="checkUsername(username)"
        required
      />
      <div v-if="isCheckingUsername" class="check-status status-checking">
        Checking availability...
      </div>
      <div
        v-else-if="usernameStatus === 'available' && username"
        class="check-status status-available"
      >
        Username available
      </div>
      <div v-else-if="usernameStatus === 'taken'" class="check-status status-taken">
        Username already exists
      </div>
    </div>

    <div class="form-group">
      <label for="email" class="form-label">Business Email</label>
      <input
        type="email"
        id="email"
        name="email"
        class="form-input"
        placeholder="business@example.com"
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

<style scoped>
.seller-register-form .form-input:focus {
  border-color: #22c55e;
  background-color: white;
  box-shadow: 0 0 0 3px rgba(34, 197, 94, 0.1);
}

.form-label {
  color: #374151; /* Darker text */
  font-weight: 600;
}

.check-status {
  font-size: 12px;
  margin-top: 4px;
}

.status-checking {
  color: #64748b;
}

.status-available {
  color: #10b981;
  font-weight: 500;
}

.status-taken {
  color: #ef4444;
  font-weight: 500;
}

.input-error {
  border-color: #ef4444 !important;
}

.input-error:focus {
  box-shadow: 0 0 0 2px rgba(239, 68, 68, 0.1) !important;
}

.input-success {
  border-color: #10b981 !important;
}

.input-success:focus {
  box-shadow: 0 0 0 2px rgba(16, 185, 129, 0.1) !important;
}
</style>
