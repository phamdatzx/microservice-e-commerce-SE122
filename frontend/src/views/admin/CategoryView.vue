<script setup lang="ts">
import { Delete, Edit, Loading, Picture, Plus, Search } from '@element-plus/icons-vue'
import axios from 'axios'
import {
  ElLoading,
  ElMessageBox,
  ElNotification,
  type FormInstance,
  type FormRules,
} from 'element-plus'
import { onMounted, reactive, ref, computed, watch } from 'vue'

interface AdminCategory {
  id: string
  name: string
  image: string
  productCount: number
}

const token = localStorage.getItem('access_token') || ''

const categoryData = ref<AdminCategory[]>([])
const dialogVisible = ref(false)
const dialogMode = ref<'add' | 'edit'>('add')
const dialogContent = ref({
  title: 'Add Category',
  mainBtnText: 'Add',
})
const isLoading = ref(false)
const searchQuery = ref('')
const selectedFile = ref<File | null>(null)

onMounted(() => {
  fetchCategories()
})

const fetchCategories = () => {
  isLoading.value = true
  axios
    .get(import.meta.env.VITE_BE_API_URL + '/product/public/category', {
      params: { name: searchQuery.value },
    })
    .then((response) => {
      categoryData.value = response.data.map((item: any) => ({
        id: item.id,
        name: item.name,
        image: item.image,
        productCount: item.product_count || 0,
      }))
    })
    .catch((error) => {
      console.error('Failed to fetch categories:', error)
      ElNotification({
        title: 'Error',
        message: 'Failed to load categories',
        type: 'error',
      })
    })
    .finally(() => {
      isLoading.value = false
    })
}

let searchTimeout: any = null
watch(searchQuery, () => {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    fetchCategories()
  }, 300)
})

const filteredData = computed(() => {
  return categoryData.value
})

const fileInput = ref<HTMLInputElement | null>(null)
const imagePreview = ref<string>('')

const handleUploadClick = () => {
  fileInput.value?.click()
}

const handleFileChange = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files && target.files[0]) {
    selectedFile.value = target.files[0]
    const reader = new FileReader()
    reader.onload = (e) => {
      imagePreview.value = e.target?.result as string
    }
    reader.readAsDataURL(selectedFile.value)
  }
}

const openModal = (mode: 'add' | 'edit', row?: AdminCategory) => {
  dialogMode.value = mode
  dialogVisible.value = true
  imagePreview.value = ''
  selectedFile.value = null

  if (mode === 'add') {
    dialogContent.value = {
      title: 'Add Category',
      mainBtnText: 'Add',
    }
    clearRuleForm()
  } else if (mode === 'edit' && row) {
    dialogContent.value = {
      title: 'Edit Category',
      mainBtnText: 'Save',
    }
    ruleForm.category = row.name
    ruleForm.id = row.id
    imagePreview.value = row.image
  }
}

const handleAddCategory = () => {
  const loading = ElLoading.service({ text: 'Adding category...' })
  const formData = new FormData()
  formData.append('name', ruleForm.category)
  if (selectedFile.value) {
    formData.append('image', selectedFile.value)
  }

  axios
    .post(import.meta.env.VITE_BE_API_URL + '/product/category', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
        Authorization: `Bearer ${token}`,
      },
    })
    .then(() => {
      ElNotification({
        title: 'Success!',
        message: `Category "${ruleForm.category}" added.`,
        type: 'success',
      })
      fetchCategories()
      dialogVisible.value = false
      clearRuleForm()
    })
    .catch((error) => {
      console.error('Failed to add category:', error)
      ElNotification({
        title: 'Error',
        message: error.response?.data?.message || 'Failed to add category',
        type: 'error',
      })
    })
    .finally(() => {
      loading.close()
    })
}

const handleEditCategory = () => {
  const loading = ElLoading.service({ text: 'Updating category...' })
  const formData = new FormData()
  formData.append('name', ruleForm.category)
  if (selectedFile.value) {
    formData.append('image', selectedFile.value)
  }

  axios
    .put(import.meta.env.VITE_BE_API_URL + `/product/category/${ruleForm.id}`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
        Authorization: `Bearer ${token}`,
      },
    })
    .then(() => {
      ElNotification({
        title: 'Success!',
        message: `Category updated.`,
        type: 'success',
      })
      fetchCategories()
      dialogVisible.value = false
      clearRuleForm()
    })
    .catch((error) => {
      console.error('Failed to update category:', error)
      ElNotification({
        title: 'Error',
        message: error.response?.data?.message || 'Failed to update category',
        type: 'error',
      })
    })
    .finally(() => {
      loading.close()
    })
}

const handleDeleteBtnClick = (row: AdminCategory) => {
  ElMessageBox.confirm(
    `Are you sure you want to delete category "${row.name}"? This action cannot be undone.`,
    'Confirm Delete',
    {
      type: 'warning',
      confirmButtonText: 'Delete',
      cancelButtonText: 'Cancel',
    },
  ).then(() => {
    const loading = ElLoading.service({ text: 'Deleting category...' })
    axios
      .delete(import.meta.env.VITE_BE_API_URL + `/product/category/${row.id}`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
      .then(() => {
        ElNotification({
          title: 'Success',
          message: 'Category deleted',
          type: 'success',
        })
        fetchCategories()
      })
      .catch((error) => {
        console.error('Failed to delete category:', error)
        ElNotification({
          title: 'Error',
          message: error.response?.data?.message || 'Failed to delete category',
          type: 'error',
        })
      })
      .finally(() => {
        loading.close()
      })
  })
}

interface RuleForm {
  category: string
  id: string
}
const ruleFormRef = ref<FormInstance>()
const ruleForm = reactive<RuleForm>({
  category: '',
  id: '',
})
const rules = reactive<FormRules<RuleForm>>({
  category: [{ required: true, message: 'Please input category name', trigger: 'change' }],
})

const submitForm = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  await formEl.validate((valid) => {
    if (valid) {
      if (dialogMode.value === 'add') handleAddCategory()
      else handleEditCategory()
    }
  })
}

const clearRuleForm = () => {
  ruleFormRef.value?.resetFields()
  ruleForm.category = ''
  ruleForm.id = ''
  imagePreview.value = ''
  if (fileInput.value) fileInput.value.value = ''
}
// #endregion
</script>

<template>
  <div class="admin-category-view">
    <div class="toolbar">
      <h2>General Category Management</h2>
      <div style="display: flex; gap: 12px">
        <el-input
          v-model="searchQuery"
          placeholder="Search categories..."
          :prefix-icon="Search"
          style="width: 240px"
          clearable
        />
        <el-button size="large" :icon="Plus" color="var(--main-color)" @click="openModal('add')">
          Add New Category
        </el-button>
      </div>
    </div>

    <el-table v-loading="isLoading" :data="filteredData" border style="width: 100%">
      <el-table-column type="index" label="#" width="60" align="center" />

      <el-table-column label="Image" width="120" align="center">
        <template #default="{ row }">
          <el-image
            style="width: 50px; height: 50px; border-radius: 4px"
            :src="row.image || 'https://placehold.co/200x200?text=No+Image'"
            fit="cover"
            :preview-src-list="[row.image || 'https://placehold.co/200x200?text=No+Image']"
            preview-teleported
          >
            <template #placeholder>
              <div
                class="image-slot"
                style="
                  background: #f5f7fa;
                  display: flex;
                  justify-content: center;
                  align-items: center;
                  width: 100%;
                  height: 100%;
                  border-radius: 4px;
                "
              >
                <el-icon class="is-loading"><Loading /></el-icon>
              </div>
            </template>
            <template #error>
              <div
                class="image-slot"
                style="
                  background: #f5f7fa;
                  color: #909399;
                  font-size: 20px;
                  display: flex;
                  justify-content: center;
                  align-items: center;
                  width: 100%;
                  height: 100%;
                  border-radius: 4px;
                "
              >
                <el-icon><Picture /></el-icon>
              </div>
            </template>
          </el-image>
        </template>
      </el-table-column>

      <el-table-column prop="name" label="Category Name" min-width="200" />

      <el-table-column prop="productCount" label="Products" width="120" align="center" sortable>
        <template #default="{ row }">
          <el-tag type="info" effect="plain" round>{{ row.productCount }}</el-tag>
        </template>
      </el-table-column>

      <el-table-column align="center" label="Actions" width="180">
        <template #default="{ row }">
          <el-button type="primary" link :icon="Edit" @click="openModal('edit', row)"
            >Edit</el-button
          >
          <el-button type="danger" link :icon="Delete" @click="handleDeleteBtnClick(row)"
            >Delete</el-button
          >
        </template>
      </el-table-column>
    </el-table>

    <!-- Dialog -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogContent.title"
      width="500"
      align-center
      destroy-on-close
    >
      <el-form
        label-width="120px"
        ref="ruleFormRef"
        :model="ruleForm"
        :rules="rules"
        label-position="top"
      >
        <el-form-item label="Category Name" prop="category">
          <el-input
            v-model="ruleForm.category"
            size="large"
            placeholder="Enter category name"
            clearable
          />
        </el-form-item>

        <!-- Mock Image Upload Field -->
        <el-form-item label="Category Image">
          <div class="upload-container" @click="handleUploadClick">
            <div v-if="imagePreview" class="image-preview-wrapper">
              <img :src="imagePreview" class="image-preview" />
              <div class="overlay">
                <el-icon><Edit /></el-icon>
              </div>
            </div>
            <div v-else class="upload-placeholder">
              <el-icon class="upload-icon"><Plus /></el-icon>
              <span>Upload Image</span>
            </div>
            <input
              type="file"
              ref="fileInput"
              style="display: none"
              accept="image/*"
              @change="handleFileChange"
            />
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dialogVisible = false">Cancel</el-button>
          <el-button type="primary" @click="submitForm(ruleFormRef)">
            {{ dialogContent.mainBtnText }}
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.admin-category-view {
  padding: 24px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.toolbar h2 {
  margin: 0;
  font-size: 20px;
  color: #1e293b;
}

.upload-placeholder {
  width: 100px;
  height: 100px;
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  color: #8c939d;
  font-size: 12px;
}

.upload-placeholder:hover {
  border-color: #409eff;
  color: #409eff;
}

.upload-icon {
  font-size: 20px;
  margin-bottom: 8px;
}
.upload-container {
  display: inline-block;
}

.image-preview-wrapper {
  position: relative;
  width: 100px;
  height: 100px;
  border-radius: 6px;
  overflow: hidden;
  cursor: pointer;
  border: 1px solid #d9d9d9;
}

.image-preview {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.image-preview-wrapper:hover .overlay {
  opacity: 1;
}

.overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  opacity: 0;
  transition: opacity 0.3s;
  color: white;
  font-size: 20px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
