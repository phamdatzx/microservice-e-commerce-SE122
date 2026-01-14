<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import axios from 'axios'
import { Loading } from '@element-plus/icons-vue'
import { ElNotification } from 'element-plus'

const props = defineProps<{
  name: string
  email: string
  phone: string
  avatar?: string
}>()

const emit = defineEmits<{
  (e: 'update:name', value: string): void
  (e: 'update:email', value: string): void
  (e: 'update:phone', value: string): void
  (e: 'update:avatar', value: string): void
}>()

const localName = ref(props.name)
const localEmail = ref(props.email)
const localPhone = ref(props.phone)
const avatarPreview = ref(props.avatar || '')
const selectedFile = ref<File | null>(null)
const isLoading = ref(false)

watch(localName, (newVal) => emit('update:name', newVal))
watch(localEmail, (newVal) => emit('update:email', newVal))
watch(localPhone, (newVal) => emit('update:phone', newVal))
watch(avatarPreview, (newVal) => emit('update:avatar', newVal))

const fetchUserInfo = async () => {
  const token = localStorage.getItem('access_token')
  if (!token) return

  isLoading.value = true

  try {
    const response = await axios.get(`${import.meta.env.VITE_BE_API_URL}/user/my-info`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })

    const userData = response.data.data

    localName.value = userData.name || ''
    localEmail.value = userData.email || ''
    localPhone.value = userData.phone || ''
    if (userData.image) {
      avatarPreview.value = userData.image
    }

    // Sync with parent
    emit('update:name', localName.value)
    emit('update:email', localEmail.value)
    emit('update:phone', localPhone.value)
  } catch (error) {
    console.error('Failed to fetch user info:', error)
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchUserInfo()
})

const handleAvatarChange = async (file: any) => {
  const isJPG =
    file.raw.type === 'image/jpeg' || file.raw.type === 'image/png' || file.raw.type === 'image/jpg'

  if (!isJPG) {
    ElNotification({
      title: 'Error',
      message: 'Avatar must be JPG, JPEG or PNG format!',
      type: 'error',
    })
    return
  }

  selectedFile.value = file.raw

  // triggers upload immediately
  await uploadAvatar()
}

const uploadAvatar = async () => {
  const token = localStorage.getItem('access_token')
  if (selectedFile.value && token) {
    isLoading.value = true
    try {
      const formData = new FormData()
      formData.append('image', selectedFile.value)

      await axios.post(`${import.meta.env.VITE_BE_API_URL}/user/upload-image`, formData, {
        headers: {
          Authorization: `Bearer ${token}`,
          'Content-Type': 'multipart/form-data',
        },
      })

      ElNotification({
        title: 'Success',
        message: 'Avatar updated successfully',
        type: 'success',
      })

      // Clear selected file after upload
      selectedFile.value = null

      // Refresh info to get new avatar URL if backend returns it or just rely on preview
      fetchUserInfo()
    } catch (error) {
      console.error('Failed to upload image:', error)
      ElNotification({
        title: 'Error',
        message: 'Failed to upload avatar',
        type: 'error',
      })
      isLoading.value = false
    }
  }
}

const handleSave = () => {
  // Update other info
  emit('update:name', localName.value)
  emit('update:email', localEmail.value)
  emit('update:phone', localPhone.value)
}
</script>

<template>
  <div class="account-info">
    <h2 class="section-title">Account Info</h2>

    <div v-if="isLoading" class="loading-state">
      <el-icon class="is-loading"><Loading /></el-icon>
    </div>

    <el-row v-else :gutter="40">
      <el-col :span="16">
        <el-form label-position="top">
          <el-form-item label="Name" required>
            <el-input v-model="localName" size="large" placeholder="Enter your name" />
          </el-form-item>

          <el-form-item label="Email Address" required>
            <el-input v-model="localEmail" size="large" placeholder="Enter your email" />
          </el-form-item>

          <el-form-item label="Phone Number (Optional)">
            <el-input v-model="localPhone" size="large" placeholder="Enter your phone number" />
          </el-form-item>

          <el-form-item>
            <el-button
              type="primary"
              size="large"
              @click="handleSave"
              style="padding: 22px 40px; font-size: 16px"
              >SAVE</el-button
            >
          </el-form-item>
        </el-form>
      </el-col>

      <el-col :span="8">
        <div class="avatar-upload-container">
          <div class="avatar-preview-wrapper">
            <img :src="avatarPreview" alt="User Avatar" />
          </div>
          <el-upload
            class="avatar-uploader"
            action="#"
            :auto-upload="false"
            :show-file-list="false"
            :on-change="handleAvatarChange"
            accept=".jpg,.jpeg,.png"
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

.loading-state {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 100px 0;
  font-size: 40px;
  color: #888;
}
</style>
