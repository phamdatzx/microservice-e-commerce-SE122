<script setup lang="ts">
import { Delete, Edit } from '@element-plus/icons-vue'
import { ElNotification, type FormInstance, type FormRules } from 'element-plus'
import { reactive, ref } from 'vue'

const tableData = [
  {
    index: 1,
    category: 'Electronics',
  },
  {
    index: 2,
    category: 'Clothing',
  },
  {
    index: 3,
    category: 'Home & Kitchen',
  },
  {
    index: 4,
    category: 'Books',
  },
]

const addEditDialogVisible = ref(false)
const deleteDialogVisible = ref(false)
const title = ref('Add category')
const placeholder = ref('Please enter a category')
const mainBtnText = ref('Add')
const currentAction = ref('add')
const deletingCategory = ref('')

const openAddModal = () => {
  addEditDialogVisible.value = true

  title.value = 'Add category'
  placeholder.value = 'Please enter a category'
  mainBtnText.value = 'Add'
  currentAction.value = 'add'
  ruleForm.category = ''
}

const openEditModal = (category: string) => {
  addEditDialogVisible.value = true

  title.value = 'Edit category'
  placeholder.value = 'Please edit the category'
  mainBtnText.value = 'Save'
  currentAction.value = 'edit'
  ruleForm.category = category
}

//---------- FORM ----------------
interface RuleForm {
  category: string
}

const ruleFormRef = ref<FormInstance>()
const ruleForm = reactive<RuleForm>({
  category: '',
})

const rules = reactive<FormRules<RuleForm>>({
  category: [{ required: true, message: 'Please input category', trigger: 'change' }],
})

const submitForm = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  await formEl.validate((valid, fields) => {
    if (valid) {
      handleAddEditCategory(currentAction.value)
    } else {
      console.log('error submit!', fields)
    }
  })
}
// --------------------------------------------------------

const openDeleteModal = (category: string) => {
  deleteDialogVisible.value = true
  deletingCategory.value = category
}

const handleAddEditCategory = (currentAction: string) => {
  if (currentAction === 'add') {
    if (ruleForm.category.trim() === '') {
      ElNotification({
        title: 'Error!',
        message: 'Category name cannot be empty.',
        type: 'error',
      })
      return
    }

    ElNotification({
      title: 'Success!',
      message: `Category "${ruleForm.category}" was added successfully.`,
      type: 'success',
    })

    addEditDialogVisible.value = false
    ruleForm.category = ''
  } else if (currentAction === 'edit') {
    ElNotification({
      title: 'Success!',
      message: `Category was updated to "${ruleForm.category}" successfully.`,
      type: 'success',
    })

    addEditDialogVisible.value = false
    ruleForm.category = ''
  }
}

const handleDeleteCategory = () => {
  deleteDialogVisible.value = false

  ElNotification({
    title: 'Success!',
    message: `Category "${deletingCategory.value}" was deleted successfully.`,
    type: 'success',
  })
}
</script>

<template>
  <el-button
    type="primary"
    size="large"
    style="margin-bottom: 12px"
    color="#1aae51"
    @click="openAddModal()"
    >Add category</el-button
  >
  <el-table :data="tableData" border style="width: 100%">
    <el-table-column prop="index" label="Index" width="100" />
    <el-table-column prop="category" label="Category" />
    <el-table-column align="center" label="Operations" width="140">
      <template #default="{ row }">
        <el-button type="primary" :icon="Edit" @click="openEditModal(row.category)" />
        <el-button type="danger" :icon="Delete" @click="openDeleteModal(row.category)" />
      </template>
    </el-table-column>
  </el-table>

  <!-- Add/Edit dialog -->
  <el-dialog v-model="addEditDialogVisible" :title="title" width="500" align-center>
    <el-form
      label-width="auto"
      style="max-width: 600px"
      ref="ruleFormRef"
      :model="ruleForm"
      :rules="rules"
    >
      <el-form-item prop="category">
        <el-input v-model="ruleForm.category" size="large" :placeholder="placeholder" clearable />
      </el-form-item>
    </el-form>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="addEditDialogVisible = false">Cancel</el-button>
        <el-button type="primary" @click="submitForm(ruleFormRef)">
          {{ mainBtnText }}
        </el-button>
      </div>
    </template>
  </el-dialog>

  <!-- Delete dialog -->
  <el-dialog v-model="deleteDialogVisible" title="Delete category" width="500" align-center>
    <span>Are you sure you want to delete category "{{ deletingCategory }}"?</span>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="deleteDialogVisible = false">Cancel</el-button>
        <el-button type="danger" @click="handleDeleteCategory()"> Delete </el-button>
      </div>
    </template>
  </el-dialog>
</template>
