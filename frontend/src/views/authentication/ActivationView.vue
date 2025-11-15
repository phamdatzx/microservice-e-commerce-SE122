<script setup lang="ts">
import axios from 'axios'
import { onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'

const USER_API_URL = import.meta.env.VITE_USER_API_URL
const route = useRoute()
const token = route.query.token

const count = ref(3)

const activationStatus = ref('success')

onMounted(() => {
  setInterval(() => {
    if (activationStatus.value === 'success') {
      if (count.value > 0) {
        count.value -= 1
      } else {
        window.location.href = '/login'
      }
    }
  }, 1000)

  axios
    .post(`${USER_API_URL}/activate?token=${token}`)
    .then((response) => {
      console.log(response.data)
      if (response.data.status === 200) {
        activationStatus.value = 'success'
      } else {
        activationStatus.value = 'failed'
      }
    })
    .catch((error) => {
      activationStatus.value = 'failed'
    })
})
</script>

<template>
  <div
    :style="{
      width: '100vw',
      height: '100vh',
      display: 'flex',
      justifyContent: 'center',
      alignItems: 'center',
      backgroundColor: activationStatus === 'success' ? '#cdfacf' : '#f9bfbf',
    }"
  >
    <div class="container">
      <div
        class="checkmark"
        :style="{
          backgroundColor: activationStatus === 'success' ? '#10b981' : '#ef4444',
        }"
      >
        <svg v-if="activationStatus === 'success'" viewBox="0 0 52 52">
          <path d="M14 27l10 10 20-20" />
        </svg>
        <svg v-else viewBox="0 0 52 52">
          <line x1="16" y1="16" x2="36" y2="36" />
          <line x1="36" y1="16" x2="16" y2="36" />
        </svg>
      </div>
      <h1 v-if="activationStatus === 'success'">Your account has been activated</h1>
      <h1 v-else>Error activating your account</h1>
      <p v-if="activationStatus === 'success'">You can safely close this page now.</p>
      <p v-else>Please contact support for further assistance.</p>
      <p
        v-if="activationStatus === 'success'"
        style="margin-top: 20px; font-size: 14px; color: #374151"
      >
        Redirecting to login page in {{ count }} seconds...
      </p>
    </div>
  </div>
</template>

<style scoped>
.container {
  background: white;
  border-radius: 16px;
  padding: 80px 40px;
  text-align: center;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  max-width: 600px;
  width: 100%;
  animation: fadeIn 0.6s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.checkmark {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  margin: 0 auto 30px;
  display: flex;
  align-items: center;
  justify-content: center;
  animation: scaleIn 0.5s ease-out 0.2s both;
}

@keyframes scaleIn {
  from {
    transform: scale(0);
  }
  to {
    transform: scale(1);
  }
}

.checkmark svg {
  width: 50px;
  height: 50px;
  stroke: white;
  stroke-width: 3;
  fill: none;
  stroke-linecap: round;
  stroke-linejoin: round;
  stroke-dasharray: 100;
  stroke-dashoffset: 100;
  animation: drawCheck 0.8s ease-out 0.5s forwards;
}

@keyframes drawCheck {
  to {
    stroke-dashoffset: 0;
  }
}

h1 {
  color: #1f2937;
  font-size: 28px;
  margin-bottom: 16px;
  font-weight: 600;
}

p {
  color: #6b7280;
  font-size: 16px;
  line-height: 1.6;
}
</style>
