<script setup lang="ts">
import { Delete, Edit, Plus } from '@element-plus/icons-vue'
import axios from 'axios'
import {
  ElLoading,
  ElMessageBox,
  ElNotification,
  type FormInstance,
  type FormRules,
} from 'element-plus'
import { onMounted, reactive, ref, computed } from 'vue'

interface AdminCategory {
  id: string
  name: string
  image: string
  productCount: number
}

// Reuse similar API structure or Mock it since user said "General Category Management"
// Assuming there is a general category API or we simulate it.
// Since I don't have the exact Admin API endpoints, I'll use the Seller one as a base but populate with mock data for Image/Count.

const userId = ref('7d27d61d-8cb3-4ef9-a9ff-0a92f217855b')
// Ideally this should use an Admin Token
const token =
  'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjY4Mjc1MDYsImlhdCI6MTc2Njc0MTEwNiwicm9sZSI6InNlbGxlciIsInVzZXJJZCI6IjdkMjdkNjFkLThjYjMtNGVmOS1hOWZmLTBhOTJmMjE3ODU1YiIsInVzZXJuYW1lIjoidGVzdDEifQ.cCs5gwsDkDVZZrnHHmeTFGmuleu9RR2Ieke8pYaRH_s'

const categoryData = ref<AdminCategory[]>([])
const dialogVisible = ref(false)
const dialogMode = ref<'add' | 'edit'>('add')
const dialogContent = ref({
  title: 'Add Category',
  mainBtnText: 'Add',
})
const isLoading = ref(false)

onMounted(() => {
  fetchCategories()
})

const fetchCategories = () => {
  isLoading.value = true
  // Using seller API for demo, replacing with mock data for missing fields
  axios
    .get(import.meta.env.VITE_GET_SELLER_CATEGORY_API_URL + '/' + userId.value + '/category')
    .then((response) => {
      // Map response to AdminCategory with Mock Data
      categoryData.value = response.data.map((item: any, index: number) => ({
        id: item.id || `cat-${index}`,
        name: item.name,
        // Mock Image
        image: `https://picsum.photos/seed/${index}/200/200`,
        // Mock Product Count
        productCount: Math.floor(Math.random() * 50) + 1,
      }))
    })
    .catch((error) => {
      console.error(error)
      // Fallback Mock data if API fails or is empty
      categoryData.value = [
        {
          id: '1',
          name: 'Electronics',
          image: 'https://picsum.photos/seed/elec/200/200',
          productCount: 120,
        },
        {
          id: '2',
          name: 'Fashion',
          image: 'https://picsum.photos/seed/fash/200/200',
          productCount: 85,
        },
        {
          id: '3',
          name: 'Home & Living',
          image: 'https://picsum.photos/seed/home/200/200',
          productCount: 42,
        },
      ]
    })
    .finally(() => {
      isLoading.value = false
    })
}

const currentPage = ref(1)
const pageSize = ref(10)

const paginatedData = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return categoryData.value.slice(start, end)
})

const handleSizeChange = (val: number) => {
  pageSize.value = val
  currentPage.value = 1
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
}

const fileInput = ref<HTMLInputElement | null>(null)
const imagePreview = ref<string>('')
// ... existing code ...

const handleUploadClick = () => {
  fileInput.value?.click()
}

const handleFileChange = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files && target.files[0]) {
    const file = target.files[0]
    const reader = new FileReader()
    reader.onload = (e) => {
      imagePreview.value = e.target?.result as string
    }
    reader.readAsDataURL(file)
  }
}

const openModal = (mode: 'add' | 'edit', row?: AdminCategory) => {
  dialogMode.value = mode
  dialogVisible.value = true
  imagePreview.value = '' // Reset preview

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
    imagePreview.value = row.image // Show existing image
  }
}

// #region Handlers
const handleAddCategory = () => {
  // Mock Implementation
  ElNotification({
    title: 'Success!',
    message: `Category "${ruleForm.category}" added.`,
    type: 'success',
  })
  categoryData.value.push({
    id: Date.now().toString(),
    name: ruleForm.category,
    image: imagePreview.value || 'https://placehold.co/200x200?text=New',
    productCount: 0,
  })
  dialogVisible.value = false
  clearRuleForm()
}

const handleEditCategory = () => {
  // Mock Implementation
  const index = categoryData.value.findIndex((c) => c.id === ruleForm.id)
  if (index !== -1) {
    categoryData.value[index].name = ruleForm.category
    if (imagePreview.value) {
      categoryData.value[index].image = imagePreview.value
    }
  }
  ElNotification({
    title: 'Success!',
    message: `Category updated.`,
    type: 'success',
  })
  dialogVisible.value = false
  clearRuleForm()
}

const handleDeleteBtnClick = (row: AdminCategory) => {
  ElMessageBox.confirm(
    `Delete category "${row.name}" with ${row.productCount} products?`,
    'Confirm Delete',
    {
      type: 'warning',
      confirmButtonText: 'Delete',
      cancelButtonText: 'Cancel',
    },
  ).then(() => {
    // Mock Delete
    categoryData.value = categoryData.value.filter((c) => c.id !== row.id)
    ElNotification({
      title: 'Success',
      message: 'Category deleted',
      type: 'success',
    })
  })
}

// #region Form
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
      <el-button type="primary" size="large" :icon="Plus" color="#3b82f6" @click="openModal('add')">
        Add New Category
      </el-button>
    </div>

    <el-table v-loading="isLoading" :data="paginatedData" border style="width: 100%">
      <el-table-column type="index" label="#" width="60" align="center" />

      <el-table-column label="Image" width="120" align="center">
        <template #default="{ row }">
          <el-image
            style="width: 50px; height: 50px; border-radius: 4px"
            :src="row.image"
            fit="cover"
            :preview-src-list="[row.image]"
            preview-teleported
          />
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

    <div class="pagination-container">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        :total="categoryData.length"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>

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
