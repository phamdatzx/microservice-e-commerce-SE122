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

              <div style="text-align: center; font-size: 14px">
                <span style="color: #9ca3af; margin-right: 8px">{{ propsData.switchText }}</span>
                <RouterLink :to="propsData.switchLinkUrl ?? ''" class="switch-link">{{
                  propsData.switchLinkText
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
