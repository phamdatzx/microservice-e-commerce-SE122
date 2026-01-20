<script setup lang="ts">
import { Check, CircleClose, Delete, Edit, Lock, Unlock } from '@element-plus/icons-vue'
import axios from 'axios'
import { ElLoading, ElMessageBox, ElNotification } from 'element-plus'
import { onMounted, ref, watch } from 'vue'

interface AdminUser {
  id: string
  email: string
  username: string
  name: string
  role: string
  isBanned: boolean
  isActive: boolean
}

const token = localStorage.getItem('access_token') || ''
const userData = ref<AdminUser[]>([])
const isLoading = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)
const totalUsers = ref(0)
const roleFilter = ref('')
const statusFilter = ref('')

onMounted(() => {
  fetchUsers()
})

const fetchUsers = () => {
  isLoading.value = true
  axios
    .get(import.meta.env.VITE_BE_API_URL + '/user/admin/users', {
      params: {
        page: currentPage.value,
        limit: pageSize.value,
        role: roleFilter.value || undefined,
        is_banned: statusFilter.value === '' ? undefined : statusFilter.value === 'banned',
      },
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
    .then((response) => {
      const data = response.data.data
      userData.value = (data.users || []).map((user: any) => ({
        id: user.id,
        email: user.email,
        username: user.username,
        name: user.name,
        role: user.role,
        isBanned: user.is_banned,
        isActive: user.is_active,
      }))
      totalUsers.value = data.total
    })
    .catch((error) => {
      console.error('Failed to fetch users:', error)
      ElNotification({
        title: 'Error',
        message: 'Failed to load users',
        type: 'error',
      })
    })
    .finally(() => {
      isLoading.value = false
    })
}

// Watch filters and pagination
watch([currentPage, pageSize, roleFilter, statusFilter], () => {
  fetchUsers()
})

const handleSizeChange = (val: number) => {
  pageSize.value = val
  currentPage.value = 1
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
}

const clearFilters = () => {
  roleFilter.value = ''
  statusFilter.value = ''
  currentPage.value = 1
}

const handleToggleBan = (row: AdminUser) => {
  const action = row.isBanned ? 'unban' : 'ban'
  ElMessageBox.confirm(
    `Are you sure you want to ${action} user "${row.username}"?`,
    'Confirm Action',
    {
      type: 'warning',
      confirmButtonText: action.charAt(0).toUpperCase() + action.slice(1),
      cancelButtonText: 'Cancel',
    },
  ).then(() => {
    const loading = ElLoading.service({
      text: `${action === 'ban' ? 'Banning' : 'Unbanning'} user...`,
    })
    axios
      .put(
        import.meta.env.VITE_BE_API_URL + `/user/admin/users/${row.id}/ban`,
        { is_banned: !row.isBanned },
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        },
      )
      .then(() => {
        ElNotification({
          title: 'Success',
          message: `User ${action}ned successfully`,
          type: 'success',
        })
        fetchUsers()
      })
      .catch((error) => {
        console.error(`Failed to ${action} user:`, error)
        ElNotification({
          title: 'Error',
          message: error.response?.data?.message || `Failed to ${action} user`,
          type: 'error',
        })
      })
      .finally(() => {
        loading.close()
      })
  })
}
</script>

<template>
  <div class="user-management-view">
    <div class="toolbar">
      <h2>User Management</h2>
      <div class="filters">
        <el-select v-model="roleFilter" placeholder="Filter by Role" clearable style="width: 150px">
          <el-option label="Admin" value="admin" />
          <el-option label="Seller" value="seller" />
          <el-option label="Customer" value="customer" />
        </el-select>
        <el-select
          v-model="statusFilter"
          placeholder="Filter by Status"
          clearable
          style="width: 150px"
        >
          <el-option label="Active" value="active" />
          <el-option label="Banned" value="banned" />
        </el-select>
        <el-button @click="clearFilters">Clear filters</el-button>
      </div>
    </div>

    <el-table v-loading="isLoading" :data="userData" border style="width: 100%">
      <el-table-column type="index" label="#" width="60" align="center" />
      <el-table-column prop="email" label="Email" min-width="200" />
      <el-table-column prop="username" label="Username" width="150" />
      <el-table-column prop="name" label="Full Name" width="180" />
      <el-table-column prop="role" label="Role" width="120" align="center">
        <template #default="{ row }">
          <el-tag
            :type="row.role === 'admin' ? 'danger' : row.role === 'seller' ? 'warning' : 'success'"
          >
            {{ row.role.toUpperCase() }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="isBanned" label="Status" width="120" align="center">
        <template #default="{ row }">
          <el-tag :type="row.isBanned ? 'danger' : 'success'" effect="dark">
            {{ row.isBanned ? 'BANNED' : 'ACTIVE' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column align="center" label="Actions" width="150">
        <template #default="{ row }">
          <el-button
            :type="row.isBanned ? 'success' : 'danger'"
            link
            :icon="row.isBanned ? Unlock : Lock"
            @click="handleToggleBan(row)"
          >
            {{ row.isBanned ? 'Unban' : 'Ban' }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination-container">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        :total="totalUsers"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>
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

.filters {
  display: flex;
  gap: 12px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
