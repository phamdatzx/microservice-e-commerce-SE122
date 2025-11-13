<script setup lang="ts">
import { ref, useTemplateRef, onMounted, onBeforeUnmount, watch } from 'vue'
import { RouterLink, RouterView } from 'vue-router'
import './assets/style.css'

const props = defineProps(['FormComponent', 'formType'])

const isLogingIn = ref(props.formType === 'login')
watch(
  () => props.formType,
  (newVal) => {
    isLogingIn.value = newVal === 'login'
  },
)

onMounted(() => {
  // Add smooth scrolling for anchor links
  const anchorLinks = document.querySelectorAll<HTMLAnchorElement>('a[href^="#"]')

  const handleAnchorClick = (e: Event) => {
    e.preventDefault()
    const target = document.querySelector(
      (e.currentTarget as HTMLAnchorElement).getAttribute('href') || '',
    )
    if (target) {
      target.scrollIntoView({ behavior: 'smooth' })
    }
  }

  anchorLinks.forEach((anchor) => {
    anchor.addEventListener('click', handleAnchorClick)
  })

  // Add input focus effects
  const formInputs = document.querySelectorAll<HTMLInputElement>('.form-input')

  const handleFocus = (e: Event) => {
    ;(e.currentTarget as HTMLElement).parentElement?.classList.add('focused')
  }

  const handleBlur = (e: Event) => {
    ;(e.currentTarget as HTMLElement).parentElement?.classList.remove('focused')
  }

  formInputs.forEach((input) => {
    input.addEventListener('focus', handleFocus)
    input.addEventListener('blur', handleBlur)
  })

  // Cleanup event listeners on unmount
  onBeforeUnmount(() => {
    anchorLinks.forEach((anchor) => {
      anchor.removeEventListener('click', handleAnchorClick)
    })
    formInputs.forEach((input) => {
      input.removeEventListener('focus', handleFocus)
      input.removeEventListener('blur', handleBlur)
    })
  })
})
</script>

<template>
  <main class="main-content">
    <div class="container">
      <div class="auth-container">
        <div class="illustration">
          <img
            style="max-width: 100%; height: auto; border-radius: 12px"
            src="./assets/authentication-illustration-image.svg"
            alt="Login Security Illustration"
            width="930"
            height="650"
          />
        </div>

        <div class="form-section">
          <div style="width: 100%; max-width: 400px">
            <h1 class="title">{{ isLogingIn ? 'Welcome Back' : 'Register' }}</h1>
            <p class="subtitle">{{ isLogingIn ? 'LOGIN TO CONTINUE' : 'JOIN TO US' }}</p>

            <form class="login-form" id="loginForm" ref="loginForm" @submit="handleLoginFormSubmit">
              <FormComponent />

              <button type="submit" class="login-btn" ref="loginBtn">
                {{ isLogingIn ? 'LOG IN' : 'REGISTER' }}
              </button>

              <div style="text-align: center; font-size: 14px">
                <span style="color: #9ca3af; margin-right: 8px">{{
                  isLogingIn ? 'NEW USER ?' : 'ALREADY USER ?'
                }}</span>
                <RouterLink :to="isLogingIn ? '/register' : '/login'" class="switch-link">{{
                  isLogingIn ? 'SIGN UP' : 'LOG IN'
                }}</RouterLink>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  </main>
</template>

<style scoped>
.main-content {
  padding: 60px 0;
}

.auth-container {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 80px;
  align-items: center;
  max-width: 1000px;
  margin: 0 auto;
}

.illustration {
  display: flex;
  justify-content: center;
}

.form-section {
  display: flex;
  justify-content: center;
}

.form-group {
  margin-bottom: 24px;
}

.form-label {
  display: block;
  margin-bottom: 8px;
  font-weight: 600;
  color: #333;
  font-size: 14px;
}

.form-input {
  width: 100%;
  padding: 15px 16px;
  border: 2px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
  transition: border-color 0.3s ease;

  &:focus {
    outline: none;
    border-color: #22c55e;
  }

  &::placeholder {
    color: #9ca3af;
  }
}

.login-btn {
  width: 100%;
  background: linear-gradient(135deg, #22c55e, #16a34a);
  color: white;
  border: none;
  padding: 16px;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: transform 0.2s ease;
  margin-bottom: 24px;
}

.login-btn:hover {
  transform: translateY(-1px);
}

/* Responsive Design */
@media (max-width: 768px) {
  .container {
    max-width: 768px;
    padding: 0 16px;
  }

  .header-content {
    flex-direction: column;
    gap: 20px;
  }

  .nav {
    gap: 20px;
  }

  .header-actions {
    gap: 20px;
  }

  .auth-container {
    grid-template-columns: 1fr;
    gap: 40px;
    text-align: center;
  }

  .illustration {
    order: 2;
  }

  .form-section {
    order: 1;
  }

  .search-bar {
    flex-direction: column;
  }

  .category-select {
    min-width: auto;
  }
}

@media (max-width: 480px) {
  .header-top .container {
    flex-direction: column;
    gap: 10px;
  }

  .nav {
    flex-wrap: wrap;
    justify-content: center;
  }

  .title {
    font-size: 28px;
  }

  .main-content {
    padding: 40px 0;
  }
}
</style>
