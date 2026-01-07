<script setup lang="ts">
import axios from 'axios'
import { ElLoading, ElNotification } from 'element-plus'
import { ref } from 'vue'
import './assets/formStyle.css'

const USER_API_URL = import.meta.env.VITE_BE_API_URL + '/user/public'

const username = ref('')
const sendResetSuccess = ref(false)

const emits = defineEmits(['success'])

const handleFormSent = () => {
  const loading = ElLoading.service({
    lock: true,
    text: 'Sending request',
    background: 'rgba(0, 0, 0, 0.7)',
  })

  axios
    .post(`${USER_API_URL}/send-reset-password`, {
      username: username.value,
    })
    .then((response) => {
      if (response.data.status === 200) {
        ElNotification({
          title: 'Request sended!',
          message: 'Please check your email to reset your password.',
          type: 'success',
        })
        sendResetSuccess.value = true
        emits('success')
      } else {
        ElNotification({
          title: 'Request failed!',
          message: 'Unable to send your password reset request.',
          type: 'error',
        })
        alert('Sending password request failed: ' + response.data)
      }
    })
    .catch((error) => {
      ElNotification({
        title: 'Request failed!',
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
  <p v-if="sendResetSuccess" class="info-text" style="margin-bottom: 140px">
    Please check your email to reset your password.
  </p>
  <div v-else class="form-group">
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
