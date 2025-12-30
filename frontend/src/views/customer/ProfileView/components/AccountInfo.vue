<script setup lang="ts">
import { ref } from 'vue'

const props = defineProps<{
  name: string
  email: string
  phone: string
}>()

const emit = defineEmits<{
  (e: 'update:name', value: string): void
  (e: 'update:email', value: string): void
  (e: 'update:phone', value: string): void
}>()

const localName = ref(props.name)
const localEmail = ref(props.email)
const localPhone = ref(props.phone)
const avatarPreview = ref('')

const handleAvatarChange = (file: any) => {
  const reader = new FileReader()
  reader.onload = (e: any) => {
    avatarPreview.value = e.target.result
  }
  reader.readAsDataURL(file.raw)
}

const handleSave = () => {
  emit('update:name', localName.value)
  emit('update:email', localEmail.value)
  emit('update:phone', localPhone.value)
}
</script>

<template>
  <div class="account-info">
    <h2 class="section-title">Account Info</h2>

    <el-row :gutter="40">
      <el-col :span="16">
        <el-form label-position="top">
          <el-form-item label="Name" required>
            <el-input v-model="localName" size="large" placeholder="Mark Cole" />
          </el-form-item>

          <el-form-item label="Email Address" required>
            <el-input v-model="localEmail" size="large" placeholder="swoo@gmail.com" />
          </el-form-item>

          <el-form-item label="Phone Number (Optional)">
            <el-input v-model="localPhone" size="large" placeholder="+1 0231 4554 452" />
          </el-form-item>

          <el-form-item>
            <el-button type="success" size="large" class="save-btn" @click="handleSave"
              >SAVE</el-button
            >
          </el-form-item>
        </el-form>
      </el-col>

      <el-col :span="8">
        <div class="avatar-upload-container">
          <div class="avatar-preview-wrapper">
            <img :src="avatarPreview || '/src/assets/avatar.jpg'" alt="User Avatar" />
          </div>
          <el-upload
            class="avatar-uploader"
            action="#"
            :auto-upload="false"
            :show-file-list="false"
            :on-change="handleAvatarChange"
          >
            <el-button type="default" size="default" class="select-image-btn"
              >Select Image</el-button
            >
          </el-upload>
          <div class="upload-tip">
            <p>Maximum file size 1 MB</p>
            <p>Format: .JPEG, .PNG</p>
          </div>
        </div>
      </el-col>
    </el-row>
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

.save-btn {
  background-color: var(--main-color);
  border-color: var(--main-color);
  padding: 0 40px;
  height: 50px;
  font-size: 16px;
  font-weight: 700;
  border-radius: 8px;
  margin-top: 20px;
}

.save-btn:hover {
  background-color: #16a34a;
  border-color: #16a34a;
}

.avatar-upload-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding-left: 20px;
  border-left: 1px solid #f0f0f0;
}

.avatar-preview-wrapper {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  overflow: hidden;
  margin-bottom: 20px;
  border: 1px solid #eee;
  background-color: #f9f9f9;
}

.avatar-preview-wrapper img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.select-image-btn {
  border-color: #dcdfe6;
  color: #606266;
  padding: 8px 16px;
  font-weight: 500;
}

.select-image-btn:hover {
  background-color: #f5f5f5;
  border-color: #c0c4cc;
}

.upload-tip {
  margin-top: 15px;
  text-align: center;
  color: #999;
  font-size: 13px;
  line-height: 1.5;
}

.upload-tip p {
  margin: 0;
}
</style>
