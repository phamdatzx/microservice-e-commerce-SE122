<script setup lang="ts">
import Header from '@/components/Header.vue'
import { computed, ref, watch } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const propsData = computed(() => route.meta)

const formComponent = ref()
const mainBtn = ref()

const handleSuccess = (userData: any) => {
  if (route.path === '/login' || route.path === '/') {
    if (userData?.role === 'seller') {
      router.push('/seller')
    } else if (userData?.role === 'admin') {
      router.push('/admin')
    } else {
      router.push('/')
    }
  }
}

// Fix: Main btn still hidden when back to other Forms
watch(
  () => route.meta,
  (newVal) => {
    mainBtn.value.style.display = 'block'
  },
)
</script>

<template>
  <main class="main-content">
    <RouterLink to="/" class="back-home-link">
      <svg
        width="18"
        height="18"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <path d="m3 9 9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" />
        <polyline points="9 22 9 12 15 12 15 22" />
      </svg>
      Back to Home
    </RouterLink>
    <div class="container">
      <div class="auth-container">
        <div class="illustration">
          <img
            style="max-width: 100%; height: auto; border-radius: 12px"
            src="../../assets/authentication-illustration-image.svg"
            alt="Authenticate Security Illustration"
            width="930"
            height="650"
          />
        </div>

        <div class="form-section">
          <div style="width: 100%; max-width: 400px">
            <h1 class="title">{{ propsData.title }}</h1>
            <p class="subtitle">{{ propsData.subtitle }}</p>

            <form
              class="login-form"
              id="loginForm"
              @submit.prevent="formComponent.handleFormSent()"
            >
              <RouterView v-slot="{ Component }">
                <component
                  :is="Component"
                  ref="formComponent"
                  @success="
                    (userData: any) => {
                      mainBtn.style.display = 'none'
                      handleSuccess(userData)
                    }
                  "
                />
              </RouterView>

              <button type="submit" class="main-btn" ref="mainBtn">
                {{ propsData.mainBtnText }}
              </button>

              <div class="auth-footer">
                <div class="switch-container">
                  <span class="switch-label">{{ propsData.switchText }}</span>
                  <RouterLink :to="propsData.switchLinkUrl ?? ''" class="switch-link">
                    {{ propsData.switchLinkText }}
                  </RouterLink>
                </div>

                <div v-if="route.path === '/login'" class="switch-container">
                  <span class="switch-label">NEW SELLER ?</span>
                  <RouterLink to="/register/seller" class="switch-link">
                    SIGN UP AS SELLER
                  </RouterLink>
                </div>
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
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
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

.main-btn {
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

.main-btn:hover {
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

.auth-footer {
  margin-top: 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
}

.switch-container {
  font-size: 14px;
}

.switch-label {
  color: #9ca3af;
  margin-right: 8px;
}

.divider {
  width: 100%;
  height: 1px;
  background-color: #f3f4f6;
}

.back-home-link {
  position: fixed;
  top: 40px;
  left: 40px;
  display: flex;
  align-items: center;
  gap: 8px;
  color: #6b7280;
  text-decoration: none;
  font-size: 14px;
  font-weight: 600;
  transition: all 0.2s ease;
  padding: 10px 20px;
  border-radius: 12px;
  background-color: white;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  border: 1px solid #f4f4f5;
  z-index: 100;
}

.back-home-link:hover {
  color: #22c55e;
  border-color: #22c55e;
  transform: translateX(-4px);
  box-shadow: 0 4px 15px rgba(34, 197, 94, 0.15);
}

.back-home-link svg {
  color: #9ca3af;
  transition: color 0.2s ease;
}

.back-home-link:hover svg {
  color: #22c55e;
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
