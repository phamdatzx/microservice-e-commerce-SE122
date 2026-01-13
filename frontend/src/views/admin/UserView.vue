<script setup lang="ts">
import { Delete, Edit, Plus, Search } from '@element-plus/icons-vue'
import {
  ElLoading,
  ElMessageBox,
  ElNotification,
  type FormInstance,
  type FormRules,
} from 'element-plus'
import { onMounted, reactive, ref, computed } from 'vue'

interface AdminUser {
  id: string
  email: string
  username: string
  fullName: string
  role: string
  isActive: boolean
}

// Mock Data
const userData = ref<AdminUser[]>([])
const dialogVisible = ref(false)
const dialogContent = ref({
  title: 'Edit User',
  mainBtnText: 'Save',
})
const isLoading = ref(false)
const searchQuery = ref('')
const currentPage = ref(1)
const pageSize = ref(10)

onMounted(() => {
  fetchUsers()
})

const fetchUsers = () => {
  isLoading.value = true
  // Simulate API call
  setTimeout(() => {
    userData.value = Array.from({ length: 50 }).map((_, index) => ({
      id: `user-${index + 1}`,
      email: `user${index + 1}@example.com`,
      username: `user_${index + 1}`,
      fullName: `User ${index + 1}`,
      role: index % 3 === 0 ? 'Admin' : 'Customer',
      isActive: index % 5 !== 0, // Some inactive
    }))
    isLoading.value = false
  }, 500)
}

const filteredData = computed(() => {
  if (!searchQuery.value) return userData.value
  const query = searchQuery.value.toLowerCase()
  return userData.value.filter(
    (user) =>
      user.email.toLowerCase().includes(query) ||
      user.username.toLowerCase().includes(query) ||
      user.fullName.toLowerCase().includes(query),
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

const openModal = (row: AdminUser) => {
  dialogVisible.value = true
  Object.assign(ruleForm, row) // Populate form
}

const handleEditUser = () => {
  // Mock
  const index = userData.value.findIndex((u) => u.id === ruleForm.id)
  if (index !== -1) {
    Object.assign(userData.value[index], ruleForm)
  }
  ElNotification({
    title: 'Success',
    message: 'User updated successfully',
    type: 'success',
  })
  dialogVisible.value = false
  clearRuleForm()
}

const handleDeleteBtnClick = (row: AdminUser) => {
  ElMessageBox.confirm(`Delete user "${row.username}"?`, 'Confirm Delete', {
    type: 'warning',
    confirmButtonText: 'Delete',
    cancelButtonText: 'Cancel',
  }).then(() => {
    userData.value = userData.value.filter((u) => u.id !== row.id)
    ElNotification({
      title: 'Success',
      message: 'User deleted',
      type: 'success',
    })
  })
}

const ruleFormRef = ref<FormInstance>()
const ruleForm = reactive({
  id: '',
  email: '',
  username: '',
  fullName: '',
  role: 'Customer',
  isActive: true,
})
const rules = reactive<FormRules>({
  email: [
    { required: true, message: 'Please input email', trigger: 'blur' },
    { type: 'email', message: 'Please input correct email address', trigger: ['blur', 'change'] },
  ],
  username: [{ required: true, message: 'Please input username', trigger: 'blur' }],
  fullName: [{ required: true, message: 'Please input full name', trigger: 'blur' }],
  role: [{ required: true, message: 'Please select role', trigger: 'change' }],
})

const submitForm = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  await formEl.validate((valid) => {
    if (valid) {
      handleEditUser()
    }
  })
}

const clearRuleForm = () => {
  ruleFormRef.value?.resetFields()
  ruleForm.id = ''
  ruleForm.email = ''
  ruleForm.username = ''
  ruleForm.fullName = ''
  ruleForm.role = 'Customer'
  ruleForm.isActive = true
}
// #endregion
</script>

<template>
  <div class="user-management-view">
    <div class="toolbar">
      <h2>User Management</h2>
      <div style="display: flex; gap: 12px">
        <el-input
          v-model="searchQuery"
          placeholder="Search users..."
          :prefix-icon="Search"
          style="width: 300px"
          clearable
        />
      </div>
    </div>

    <el-table v-loading="isLoading" :data="paginatedData" border style="width: 100%">
      <el-table-column type="index" label="#" width="60" align="center" />
      <el-table-column prop="email" label="Email" min-width="200" />
      <el-table-column prop="username" label="Username" width="150" />
      <el-table-column prop="fullName" label="Full Name" width="180" />
      <el-table-column prop="role" label="Role" width="120" align="center">
        <template #default="{ row }">
          <el-tag :type="row.role === 'Admin' ? 'danger' : 'success'">{{ row.role }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="isActive" label="Status" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="row.isActive ? 'success' : 'info'" effect="plain">
            {{ row.isActive ? 'Active' : 'Inactive' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column align="center" label="Actions" width="180">
        <template #default="{ row }">
          <el-button type="primary" link :icon="Edit" @click="openModal(row)">Edit</el-button>
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
        <el-form-item label="Email" prop="email">
          <el-input v-model="ruleForm.email" placeholder="Enter email" />
        </el-form-item>
        <el-form-item label="Username" prop="username">
          <el-input v-model="ruleForm.username" placeholder="Enter username" />
        </el-form-item>
        <el-form-item label="Full Name" prop="fullName">
          <el-input v-model="ruleForm.fullName" placeholder="Enter full name" />
        </el-form-item>
        <el-form-item label="Role" prop="role">
          <el-select v-model="ruleForm.role" placeholder="Select role" style="width: 100%">
            <el-option label="Admin" value="Admin" />
            <el-option label="Customer" value="Customer" />
          </el-select>
        </el-form-item>
        <el-form-item label="Status" prop="isActive">
          <el-switch v-model="ruleForm.isActive" active-text="Active" inactive-text="Inactive" />
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
.user-management-view {
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

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
