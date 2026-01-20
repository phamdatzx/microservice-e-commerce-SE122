<script setup lang="ts">
import { useRouter } from 'vue-router'
import { HomeFilled } from '@element-plus/icons-vue'
const props = withDefaults(
  defineProps<{
    title?: string
    message?: string
    hasHeader?: boolean
  }>(),
  {
    title: 'Page Not Found',
    message: "Oops! The page you're looking for doesn't exist or has been moved.",
    hasHeader: true,
  },
)

const router = useRouter()

const goHome = () => {
  router.push('/')
}
</script>

<template>
  <div class="not-found-container">
    <div class="not-found-content">
      <h1 class="error-code">404</h1>
      <h2 class="error-title">{{ props.title }}</h2>
      <p class="error-message">
        {{ props.message }}
      </p>
      <el-button type="primary" size="large" :icon="HomeFilled" @click="goHome">
        Back to Home
      </el-button>
    </div>
  </div>
</template>

<style scoped>
.not-found-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: v-bind('props.hasHeader ? "calc(100vh - 96px)" : "100vh"');
  background-color: #f8f9fa;
  padding: 20px;
}

.not-found-content {
  text-align: center;
  max-width: 500px;
  background: white;
  padding: 60px 40px;
  border-radius: 12px;
  box-shadow: 0 8px 24px rgba(149, 157, 165, 0.1);
}

.error-code {
  font-size: 120px;
  font-weight: 900;
  color: var(--main-color);
  line-height: 1;
  margin: 0;
  margin-bottom: 20px;
}

.error-title {
  font-size: 28px;
  font-weight: 700;
  color: #333;
  margin-bottom: 16px;
}

.error-message {
  font-size: 16px;
  color: #666;
  margin-bottom: 30px;
  line-height: 1.6;
}

@media (max-width: 480px) {
  .error-code {
    font-size: 80px;
  }

  .error-title {
    font-size: 24px;
  }

  .not-found-content {
    padding: 40px 20px;
  }
}
</style>
