<script setup lang="ts">
import axios from 'axios'
import { ref, onMounted, onBeforeUnmount, watch, reactive } from 'vue'
import { RouterLink } from 'vue-router'
import LoginForm from './components/LoginForm.vue'
import RegisterForm from './components/RegisterForm.vue'
import ForgotPasswordForm from './components/ForgotPasswordForm.vue'
import Header from '@/components/Header.vue'

const props = defineProps(['formType'])

const loginForm = ref()
const registerForm = ref()
const forgotPasswordForm = ref()

const contents = reactive({
  title: '',
  subtitle: '',
  mainBtnText: '',
  switchText: '',
  switchLinkUrl: '',
  switchLinkText: '',
})
const setContents = (formType: string) => {
  if (formType === 'login') {
    contents.title = 'Welcome Back'
    contents.subtitle = 'LOGIN TO CONTINUE'
    contents.mainBtnText = 'LOG IN'
    contents.switchText = 'NEW USER ?'
    contents.switchLinkUrl = '/register'
    contents.switchLinkText = 'SIGN UP'
  } else if (formType === 'register') {
    contents.title = 'Register'
    contents.subtitle = 'JOIN TO US'
    contents.mainBtnText = 'REGISTER'
    contents.switchText = 'ALREADY USER ?'
    contents.switchLinkUrl = '/login'
    contents.switchLinkText = 'LOG IN'
  } else if (formType === 'forgot-password') {
    contents.title = 'Forgot Password'
    contents.subtitle = 'RESET YOUR PASSWORD'
    contents.mainBtnText = 'RESET PASSWORD'
    contents.switchText = 'REMEMBERED YOUR PASSWORD ?'
    contents.switchLinkUrl = '/login'
    contents.switchLinkText = 'LOG IN'
  }
}
watch(
  () => props.formType,
  (newVal) => {
    setContents(newVal)
  },
)

onMounted(() => {
  setContents(props.formType)

  // // Add smooth scrolling for anchor links
  // const anchorLinks = document.querySelectorAll<HTMLAnchorElement>('a[href^="#"]')

  // const handleAnchorClick = (e: Event) => {
  //   e.preventDefault()
  //   const target = document.querySelector(
  //     (e.currentTarget as HTMLAnchorElement).getAttribute('href') || '',
  //   )
  //   if (target) {
  //     target.scrollIntoView({ behavior: 'smooth' })
  //   }
  // }

  // anchorLinks.forEach((anchor) => {
  //   anchor.addEventListener('click', handleAnchorClick)
  // })

  // // Add input focus effects
  // const formInputs = document.querySelectorAll<HTMLInputElement>('.form-input')

  // const handleFocus = (e: Event) => {
  //   ;(e.currentTarget as HTMLElement).parentElement?.classList.add('focused')
  // }

  // const handleBlur = (e: Event) => {
  //   ;(e.currentTarget as HTMLElement).parentElement?.classList.remove('focused')
  // }

  // formInputs.forEach((input) => {
  //   input.addEventListener('focus', handleFocus)
  //   input.addEventListener('blur', handleBlur)
  // })

  // // Cleanup event listeners on unmount
  // onBeforeUnmount(() => {
  //   anchorLinks.forEach((anchor) => {
  //     anchor.removeEventListener('click', handleAnchorClick)
  //   })
  //   formInputs.forEach((input) => {
  //     input.removeEventListener('focus', handleFocus)
  //     input.removeEventListener('blur', handleBlur)
  //   })
  // })
})

const handleMainBtnClick = (formType: string) => {
  if (formType === 'login') {
    loginForm.value.handleLoginFormSubmit()
  } else if (formType === 'register') {
    registerForm.value.handleRegisterFormSubmit()
  } else if (formType === 'forgot-password') {
    forgotPasswordForm.value.handleForgotPasswordFormSubmit()
  }
}
</script>

<template>
  <Header />

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
            <h1 class="title">{{ contents.title }}</h1>
            <p class="subtitle">{{ contents.subtitle }}</p>

            <form
              class="login-form"
              id="loginForm"
              @submit.prevent="handleMainBtnClick(props.formType)"
            >
              <LoginForm v-if="props.formType === 'login'" ref="loginForm" />
              <RegisterForm v-else-if="props.formType === 'register'" ref="registerForm" />
              <ForgotPasswordForm
                v-else="props.formType === 'forgot-password'"
                ref="forgotPasswordForm"
              />

              <button type="submit" class="login-btn" ref="loginBtn">
                {{ contents.mainBtnText }}
              </button>

              <div style="text-align: center; font-size: 14px">
                <span style="color: #9ca3af; margin-right: 8px">{{ contents.switchText }}</span>
                <RouterLink :to="contents.switchLinkUrl" class="switch-link">{{
                  contents.switchLinkText
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

.title {
  font-size: 36px;
  font-weight: 700;
  color: #22c55e;
  margin-bottom: 8px;
}

.subtitle {
  color: #999;
  font-size: 14px;
  letter-spacing: 1px;
  margin-bottom: 40px;
}

.form-section {
  display: flex;
  justify-content: center;
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

.switch-link {
  color: #22c55e;
  text-decoration: none;
  font-weight: 600;
}

.switch-link:hover {
  text-decoration: underline;
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
