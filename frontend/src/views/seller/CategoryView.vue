<script setup lang="ts">
import { Delete, Edit, Plus, Search } from '@element-plus/icons-vue'
import axios from 'axios'
import {
  ElLoading,
  ElMessageBox,
  ElNotification,
  type FormInstance,
  type FormRules,
} from 'element-plus'
import { onMounted, reactive, ref, computed } from 'vue'

interface SellerCategory {
  id: string
  seller_id?: string
  name: string
  product_count: number
}

const userId = localStorage.getItem('user_id')
const token = localStorage.getItem('access_token')

const categoryData = ref<SellerCategory[]>([])
const dialogVisible = ref(false)
const dialogMode = ref<'add' | 'edit'>('add')
const dialogContent = ref({
  title: 'Add Category',
  mainBtnText: 'Add',
})
const isLoading = ref(false)
const searchQuery = ref('')
const currentPage = ref(1)
const pageSize = ref(10)

onMounted(() => {
  fetchCategories()
})

const fetchCategories = () => {
  isLoading.value = true
  axios
    .get(import.meta.env.VITE_BE_API_URL + '/product/public/seller/' + userId + '/category')
    .then((response) => {
      categoryData.value = response.data
    })
    .catch((error) => {
      ElNotification({
        title: 'Error!',
        message: 'Fetch categories failed: ' + error.message,
        type: 'error',
      })
    })
    .finally(() => {
      isLoading.value = false
    })
}

// #region Pagination & Filter
const filteredData = computed(() => {
  if (!searchQuery.value) return categoryData.value
  return categoryData.value.filter((category) =>
    category.name.toLowerCase().includes(searchQuery.value.toLowerCase()),
  )
})

const paginatedData = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredData.value.slice(start, end)
})

const handleSizeChange = (val: number) => {
  pageSize.value = val
  currentPage.value = 1
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
}
// #endregion

// #region Modal
const openModal = (mode: 'add' | 'edit', row?: SellerCategory) => {
  dialogMode.value = mode
  dialogVisible.value = true

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
  }
}
// #endregion

// #region API Actions
const handleAddCategory = () => {
  const loading = ElLoading.service({
    lock: true,
    text: 'Loading',
    background: 'rgba(0, 0, 0, 0.7)',
  })
  axios
    .post(
      import.meta.env.VITE_BE_API_URL + '/product/seller-category',
      {
        name: ruleForm.category,
      },
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      },
    )
    .then((response) => {
      ElNotification({
        title: 'Success!',
        message: `Category "${ruleForm.category}" was added successfully.`,
        type: 'success',
      })
      fetchCategories()
      dialogVisible.value = false
    })
    .catch((error) => {
      ElNotification({
        title: 'Error!',
        message: error.message,
        type: 'error',
      })
    })
    .finally(() => {
      loading.close()
    })
}

const handleEditCategory = () => {
  const loading = ElLoading.service({
    lock: true,
    text: 'Loading',
    background: 'rgba(0, 0, 0, 0.7)',
  })
  axios
    .put(
      import.meta.env.VITE_BE_API_URL + '/product/seller-category/' + ruleForm.id,
      {
        name: ruleForm.category,
      },
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      },
    )
    .then((response) => {
      ElNotification({
        title: 'Success!',
        message: `Category "${ruleForm.category}" was updated successfully.`,
        type: 'success',
      })
      fetchCategories()
      dialogVisible.value = false
    })
    .catch((error) => {
      ElNotification({
        title: 'Error!',
        message: error.message,
        type: 'error',
      })
    })
    .finally(() => {
      loading.close()
    })
}

const handleDeleteBtnClick = (row: SellerCategory) => {
  ElMessageBox.confirm(
    `Delete category "${row.name}" with ${row.product_count} products?`,
    'Confirm Delete',
    {
      type: 'warning',
      confirmButtonText: 'Delete',
      cancelButtonText: 'Cancel',
    },
  ).then(async () => {
    handleDeleteCategory(row.id, row.name)
  })
}

const handleDeleteCategory = (categoryId: string, categoryName: string) => {
  const loading = ElLoading.service({
    lock: true,
    text: 'Loading',
    background: 'rgba(0, 0, 0, 0.7)',
  })

  axios
    .delete(import.meta.env.VITE_BE_API_URL + '/product/seller-category/' + categoryId, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
    .then((response) => {
      if (response.data.code === 401) {
        ElNotification({
          title: 'Error!',
          message: response.data.message,
          type: 'error',
        })
      } else if (response.data.code === 500) {
        ElNotification({
          title: 'Error!',
          message: response.data.message,
          type: 'error',
        })
      } else {
        ElNotification({
          title: 'Success!',
          message: `Category "${categoryName}" was deleted successfully.`,
          type: 'success',
        })
        fetchCategories()
      }
    })
    .catch((error) => {
      ElNotification({
        title: 'Error!',
        message: error.message,
        type: 'error',
      })
    })
    .finally(() => {
      loading.close()
    })
}
// #endregion

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
}
// #endregion
</script>

<template>
  <div class="seller-category-view">
    <div class="toolbar">
      <h2>Category Management</h2>
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

    <div class="table-container">
      <el-table
        v-loading="isLoading"
        :data="paginatedData"
        border
        style="width: 100%; height: 100%"
      >
        <el-table-column type="index" label="#" width="60" align="center" />

        <el-table-column prop="name" label="Category Name" min-width="200" />

        <el-table-column prop="product_count" label="Products" width="120" align="center" sortable>
          <template #default="{ row }">
            <el-tag type="info" effect="plain" round>{{ row.product_count }}</el-tag>
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
    </div>

    <div class="pagination-container">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        size="large"
        layout="total, sizes, prev, pager, next, jumper"
        :total="filteredData.length"
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
.seller-category-view {
  padding: 24px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
  display: flex;
  flex-direction: column;
  height: 100%;
  box-sizing: border-box;
  overflow: hidden;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  flex-shrink: 0;
}

.toolbar h2 {
  margin: 0;
  font-size: 20px;
  color: #1e293b;
}

.table-container {
  flex: 1;
  overflow: hidden;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
  flex-shrink: 0;
}
</style>
