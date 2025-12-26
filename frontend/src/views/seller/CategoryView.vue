<script setup lang="ts">
import { Delete, Edit } from '@element-plus/icons-vue'
import axios from 'axios'
import {
  ElLoading,
  ElMessageBox,
  ElNotification,
  type FormInstance,
  type FormRules,
} from 'element-plus'
import { onMounted, reactive, ref } from 'vue'

interface SellerCategory {
  id: string
  seller_id: string
  name: string
}

type DialogMode = 'add' | 'edit'

const userId = ref('7d27d61d-8cb3-4ef9-a9ff-0a92f217855b')
// const token = localStorage.getItem('access_token')
const token =
  'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjY4Mjc1MDYsImlhdCI6MTc2Njc0MTEwNiwicm9sZSI6InNlbGxlciIsInVzZXJJZCI6IjdkMjdkNjFkLThjYjMtNGVmOS1hOWZmLTBhOTJmMjE3ODU1YiIsInVzZXJuYW1lIjoidGVzdDEifQ.cCs5gwsDkDVZZrnHHmeTFGmuleu9RR2Ieke8pYaRH_s'

const categoryData = ref<SellerCategory[]>([])
const dialogVisible = ref(false)
const dialogMode = ref<DialogMode>('add')
const dialogContent = ref({
  title: 'Add category',
  placeholder: 'Please enter a category',
  mainBtnText: 'Add',
})
const isLoading = ref(false)

onMounted(() => {
  fetchCategories()
})

const fetchCategories = () => {
  isLoading.value = true
  axios
    .get(import.meta.env.VITE_GET_SELLER_CATEGORY_API_URL + '/' + userId.value + '/category')
    .then((response) => {
      categoryData.value = response.data
      console.log(categoryData.value)
    })
    .catch((error) => {
      console.error(error)
    })
    .finally(() => {
      isLoading.value = false
    })
}

const openModal = (mode: DialogMode, category?: string, id?: string) => {
  dialogMode.value = mode
  dialogVisible.value = true

  if (mode === 'add') {
    dialogContent.value = {
      title: 'Add category',
      placeholder: 'Please enter a category',
      mainBtnText: 'Add',
    }
    clearRuleForm()
  } else if (mode === 'edit') {
    dialogContent.value = {
      title: 'Edit category',
      placeholder: 'Please edit the category',
      mainBtnText: 'Save',
    }
    ruleForm.category = category ?? ''
    ruleForm.id = id ?? ''
  }
}

const handleAddCategory = () => {
  const loading = ElLoading.service({
    lock: true,
    text: 'Loading',
    background: 'rgba(0, 0, 0, 0.7)',
  })
  if (ruleForm.category.trim() === '') {
    ElNotification({
      title: 'Error!',
      message: 'Category name cannot be empty.',
      type: 'error',
    })
    return
  } else {
    axios
      .post(
        import.meta.env.VITE_MANAGE_SELLER_CATEGORY_API_URL,
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
        // Success
        ElNotification({
          title: 'Success!',
          message: `Category "${ruleForm.category}" was added successfully.`,
          type: 'success',
        })
        fetchCategories()
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
        dialogVisible.value = false
        ruleForm.category = ''
      })
  }
}

const handleEditCategory = () => {
  const loading = ElLoading.service({
    lock: true,
    text: 'Loading',
    background: 'rgba(0, 0, 0, 0.7)',
  })
  axios
    .put(
      import.meta.env.VITE_MANAGE_SELLER_CATEGORY_API_URL + '/' + ruleForm.id,
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
      // Success
      ElNotification({
        title: 'Success!',
        message: `Category "${ruleForm.category}" was updated successfully.`,
        type: 'success',
      })
      fetchCategories()
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
      dialogVisible.value = false
      clearRuleForm()
    })
}

const handleDeleteBtnClick = (categoryId: string, categoryName: string) => {
  ElMessageBox.confirm(`Delete category name: "${categoryName}"?`, 'Confirm delete', {
    type: 'warning',
    confirmButtonText: 'Delete',
    cancelButtonText: 'Cancel',
  }).then(async () => {
    handleDeleteCategory(categoryId, categoryName)
  })
}
const handleDeleteCategory = (categoryId: string, categoryName: string) => {
  const loading = ElLoading.service({
    lock: true,
    text: 'Loading',
    background: 'rgba(0, 0, 0, 0.7)',
  })

  axios
    .delete(import.meta.env.VITE_MANAGE_SELLER_CATEGORY_API_URL + '/' + categoryId, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
    .then((response) => {
      // Failed: wrong token
      if (response.data.code === 401) {
        ElNotification({
          title: 'Error!',
          message: response.data.message,
          type: 'error',
        })
        return
      }
      // Failed: wrong id
      else if (response.data.code === 500) {
        ElNotification({
          title: 'Error!',
          message: response.data.message,
          type: 'error',
        })
      }
      // Success
      else {
        ElNotification({
          title: 'Success!',
          message: `Category "${categoryName}" was deleted successfully.`,
          type: 'success',
        })
        fetchCategories()
      }
    })
    // API failed
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
  category: [{ required: true, message: 'Please input category', trigger: 'change' }],
})

const submitForm = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  await formEl.validate((valid, fields) => {
    if (valid) {
      handleFormSubmit()
    } else {
      console.log('error submit!', fields)
    }
  })
}

const handleFormSubmit = () => {
  if (dialogMode.value === 'add') {
    handleAddCategory()
  } else if (dialogMode.value === 'edit') {
    handleEditCategory()
  }
}

const clearRuleForm = () => {
  ruleFormRef.value?.resetFields()
}
// #endregion
</script>

<template>
  <el-button
    type="primary"
    size="large"
    style="margin-bottom: 12px"
    color="#1aae51"
    @click="openModal('add')"
  >
    Add category
  </el-button>
  <el-table v-loading="isLoading" :data="categoryData" border style="width: 100%">
    <el-table-column type="index" label="Index" width="100" />
    <el-table-column prop="name" label="Category" />
    <el-table-column align="center" label="Operations" width="140">
      <template #default="{ row }">
        <el-button type="primary" :icon="Edit" @click="openModal('edit', row.name, row.id)" />
        <el-button type="danger" :icon="Delete" @click="handleDeleteBtnClick(row.id, row.name)" />
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
      label-width="auto"
      style="max-width: 600px"
      ref="ruleFormRef"
      :model="ruleForm"
      :rules="rules"
    >
      <el-form-item prop="category">
        <el-input
          v-model="ruleForm.category"
          size="large"
          :placeholder="dialogContent.placeholder"
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
</template>
