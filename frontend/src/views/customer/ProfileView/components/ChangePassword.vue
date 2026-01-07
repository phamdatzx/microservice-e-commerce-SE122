<script setup lang="ts">
import { ref } from 'vue'

const currentPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')

const handleSave = () => {
  // Logic to save password would go here
  console.log('Password saved', {
    currentPassword: currentPassword.value,
    newPassword: newPassword.value,
    confirmPassword: confirmPassword.value,
  })
}
</script>

<template>
  <div class="change-password">
    <h2 class="section-title">Change Password</h2>

    <el-form label-position="top">
      <el-form-item label="Current Password" required>
        <el-input
          v-model="currentPassword"
          type="password"
          size="large"
          placeholder="********"
          show-password
        />
      </el-form-item>

      <el-form-item label="New Password" required>
        <el-input
          v-model="newPassword"
          type="password"
          size="large"
          placeholder="Set a new password"
          show-password
        />
      </el-form-item>

      <el-form-item label="Confirm New Password" required>
        <el-input
          v-model="confirmPassword"
          type="password"
          size="large"
          placeholder="Re-enter your new password"
          show-password
        />
      </el-form-item>

      <el-form-item>
        <el-button
          type="primary"
          size="large"
          @click="handleSave"
          style="padding: 22px 40px; font-size: 16px; margin-top: 4px"
          >SAVE</el-button
        >
      </el-form-item>
    </el-form>
  </div>
</template>

<style scoped>
.section-title {
  font-size: 28px;
  font-weight: 700;
  margin-bottom: 30px;
  color: #000;
}

:deep(.el-form-item__label) {
  font-weight: 600;
  color: #333;
  margin-bottom: 8px !important;
}

:deep(.el-input__wrapper) {
  padding: 8px 15px;
  border-radius: 8px;
  box-shadow: 0 0 0 1px #dcdfe6 inset;
}

:deep(.el-input__inner) {
  height: 48px;
}
</style>
