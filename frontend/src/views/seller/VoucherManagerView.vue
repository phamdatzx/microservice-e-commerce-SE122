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
import { onMounted, reactive, ref, computed, watch } from 'vue'
import { formatNumberWithDots } from '@/utils/formatNumberWithDots'

interface Voucher {
  id: string
  seller_id: string
  code: string
  name: string
  description: string
  discount_type: 'FIXED' | 'PERCENTAGE'
  discount_value: number
  max_discount_value: number | null
  min_order_value: number
  apply_scope: 'ALL' | 'CATEGORY'
  apply_seller_category_ids: string[]
  total_quantity: number
  used_quantity: number
  usage_limit_per_user: number
  start_time: string
  end_time: string
  status: 'ACTIVE' | 'INACTIVE'
  created_at?: string
  updated_at?: string
}

interface SellerCategory {
  id: string
  seller_id: string
  name: string
}

const voucherData = ref<Voucher[]>([])
const categoryList = ref<SellerCategory[]>([])
const isLoading = ref(false)
const dialogVisible = ref(false)
const dialogMode = ref<'add' | 'edit'>('add')
const userId = ref(localStorage.getItem('user_id') || '')
const token = localStorage.getItem('access_token') || ''

const ruleFormRef = ref<FormInstance>()
const ruleForm = reactive({
  id: '',
  code: '',
  name: '',
  description: '',
  discount_type: 'FIXED',
  discount_value: 0,
  max_discount_value: null as number | null,
  min_order_value: 0,
  apply_scope: 'ALL',
  apply_seller_category_ids: [] as string[],
  total_quantity: 1,
  usage_limit_per_user: 1,
  dateRange: null as [Date, Date] | null,
  start_time: '',
  end_time: '',
  status: 'ACTIVE',
})

const rules = reactive<FormRules>({
  code: [{ required: true, message: 'Please input code', trigger: 'blur' }],
  name: [{ required: true, message: 'Please input name', trigger: 'blur' }],
  discount_value: [
    { required: true, message: 'Required', trigger: 'blur' },
    {
      validator: (rule: any, value: any, callback: any) => {
        if (value <= 0) {
          callback(new Error('Must be > 0'))
        } else if (ruleForm.discount_type === 'PERCENTAGE' && value > 100) {
          callback(new Error('Percentage must be <= 100'))
        } else {
          callback()
        }
      },
      trigger: 'blur',
    },
  ],
  min_order_value: [{ required: true, message: 'Required', trigger: 'blur' }],
  total_quantity: [{ required: true, message: 'Required', trigger: 'blur' }],
  dateRange: [{ required: true, message: 'Please select date range', trigger: 'change' }],
  apply_seller_category_ids: [
    {
      validator: (rule: any, value: any, callback: any) => {
        if (ruleForm.apply_scope === 'CATEGORY' && (!value || value.length === 0)) {
          callback(new Error('Please select at least one category'))
        } else {
          callback()
        }
      },
      trigger: 'change',
    },
  ],
})

const VOUCHER_API_URL = 'http://localhost:81/api/product/vouchers'

const isPercentage = computed(() => ruleForm.discount_type === 'PERCENTAGE')

watch(
  () => ruleForm.discount_type,
  () => {
    if (ruleFormRef.value) {
      ruleFormRef.value.validateField('discount_value')
    }
  },
)

const fetchVouchers = async () => {
  isLoading.value = true
  try {
    const response = await axios.get(VOUCHER_API_URL, {
      headers: { Authorization: `Bearer ${token}` },
    })
    voucherData.value = response.data
  } catch (error) {
    console.error('Error fetching vouchers:', error)
    ElNotification.error('Failed to load vouchers')
  } finally {
    isLoading.value = false
  }
}

const fetchCategories = async () => {
  if (!userId.value) return
  const url = `${import.meta.env.VITE_BE_API_URL}/product/public/seller/${userId.value}/category`
  try {
    const response = await axios.get(url)
    categoryList.value = response.data
  } catch (error) {
    console.error('Error fetching categories:', error)
  }
}

onMounted(() => {
  fetchVouchers()
  fetchCategories()
})

const openDialog = (mode: 'add' | 'edit', row?: Voucher) => {
  dialogMode.value = mode
  dialogVisible.value = true

  if (mode === 'add') {
    resetForm()
  } else if (row) {
    // Populate form
    ruleForm.id = row.id
    ruleForm.code = row.code
    ruleForm.name = row.name
    ruleForm.description = row.description
    ruleForm.discount_type = row.discount_type
    ruleForm.discount_value = row.discount_value
    ruleForm.max_discount_value = row.max_discount_value
    ruleForm.min_order_value = row.min_order_value
    ruleForm.apply_scope = row.apply_scope
    ruleForm.apply_seller_category_ids = row.apply_seller_category_ids || []
    ruleForm.total_quantity = row.total_quantity
    ruleForm.usage_limit_per_user = row.usage_limit_per_user
    ruleForm.status = row.status

    // Handle Dates
    if (row.start_time && row.end_time) {
      ruleForm.dateRange = [new Date(row.start_time), new Date(row.end_time)]
    } else {
      ruleForm.dateRange = null
    }
  }
}

const resetForm = () => {
  if (ruleFormRef.value) ruleFormRef.value.resetFields()
  ruleForm.id = ''
  ruleForm.code = ''
  ruleForm.name = ''
  ruleForm.description = ''
  ruleForm.discount_type = 'FIXED'
  ruleForm.discount_value = 0
  ruleForm.max_discount_value = null
  ruleForm.min_order_value = 0
  ruleForm.apply_scope = 'ALL'
  ruleForm.apply_seller_category_ids = []
  ruleForm.total_quantity = 100
  ruleForm.usage_limit_per_user = 1
  ruleForm.dateRange = null
  ruleForm.status = 'ACTIVE'
}

const handleSubmit = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  await formEl.validate(async (valid) => {
    if (valid) {
      const loading = ElLoading.service({
        lock: true,
        text: 'Saving...',
        background: 'rgba(0, 0, 0, 0.7)',
      })

      try {
        const payload = {
          code: ruleForm.code,
          name: ruleForm.name,
          description: ruleForm.description,
          discount_type: ruleForm.discount_type,
          discount_value: ruleForm.discount_value,
          max_discount_value: ruleForm.max_discount_value,
          min_order_value: ruleForm.min_order_value,
          apply_scope: ruleForm.apply_scope,
          apply_seller_category_ids:
            ruleForm.apply_scope === 'CATEGORY' ? ruleForm.apply_seller_category_ids : [],
          total_quantity: ruleForm.total_quantity,
          usage_limit_per_user: ruleForm.usage_limit_per_user,
          start_time: ruleForm.dateRange ? ruleForm.dateRange[0].toISOString() : null,
          end_time: ruleForm.dateRange ? ruleForm.dateRange[1].toISOString() : null,
          status: ruleForm.status,
        }

        if (dialogMode.value === 'add') {
          await axios.post(VOUCHER_API_URL, payload, {
            headers: { Authorization: `Bearer ${token}` },
          })
          ElNotification.success('Voucher created successfully')
        } else {
          await axios.put(`${VOUCHER_API_URL}/${ruleForm.id}`, payload, {
            headers: { Authorization: `Bearer ${token}` },
          })
          ElNotification.success('Voucher updated successfully')
        }

        dialogVisible.value = false
        fetchVouchers()
      } catch (error: any) {
        ElNotification.error(error.message || 'Failed to save voucher')
      } finally {
        loading.close()
      }
    }
  })
}

const deleteVoucher = (id: string) => {
  ElMessageBox.confirm('Are you sure you want to delete this voucher?', 'Warning', {
    confirmButtonText: 'Delete',
    cancelButtonText: 'Cancel',
    type: 'warning',
  }).then(async () => {
    try {
      await axios.delete(`${VOUCHER_API_URL}/${id}`, {
        headers: { Authorization: `Bearer ${token}` },
      })
      ElNotification.success('Voucher deleted')
      fetchVouchers()
    } catch (error) {
      ElNotification.error('Failed to delete voucher')
    }
  })
}
</script>
<template>
  <div class="voucher-manager">
    <div class="header">
      <h2>Voucher Management</h2>
      <el-button
        size="large"
        type="primary"
        :icon="Plus"
        color="var(--main-color)"
        @click="openDialog('add')"
        >Create Voucher</el-button
      >
    </div>

    <el-table
      :data="voucherData"
      v-loading="isLoading"
      style="width: 100%; margin-top: 20px"
      border
    >
      <el-table-column prop="code" label="Code" width="120" />
      <el-table-column prop="name" label="Name" min-width="150" />
      <el-table-column label="Discount" width="150">
        <template #default="{ row }">
          <span v-if="row.discount_type === 'FIXED'"
            >{{ formatNumberWithDots(row.discount_value) }}đ</span
          >
          <span v-else
            >{{ row.discount_value }}% (Max:
            {{ row.max_discount_value ? formatNumberWithDots(row.max_discount_value) : '' }}đ)</span
          >
        </template>
      </el-table-column>
      <el-table-column label="Scope" width="150">
        <template #default="{ row }">
          <el-tag v-if="row.apply_scope === 'ALL'">All Products</el-tag>
          <el-popover
            v-else-if="row.apply_scope === 'CATEGORY'"
            placement="top"
            title="Selected Categories"
            :width="200"
            trigger="hover"
          >
            <template #reference>
              <el-tag type="warning" style="cursor: pointer">
                {{ row.apply_seller_category_ids ? row.apply_seller_category_ids.length : 0 }}
                {{ (row.apply_seller_category_ids?.length || 0) === 1 ? 'Category' : 'Categories' }}
              </el-tag>
            </template>
            <ul style="padding-left: 20px; margin: 0">
              <li v-for="id in row.apply_seller_category_ids" :key="id">
                {{ categoryList.find((c) => c.id === id)?.name || 'Unknown' }}
              </li>
            </ul>
          </el-popover>
          <el-tag v-else>{{ row.apply_scope }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="Usage" width="120">
        <template #default="{ row }"> {{ row.used_quantity }} / {{ row.total_quantity }} </template>
      </el-table-column>
      <el-table-column label="Validity" width="200">
        <template #default="{ row }">
          <div>
            {{ new Date(row.start_time).toLocaleDateString() }} -
            {{ new Date(row.end_time).toLocaleDateString() }}
          </div>
        </template>
      </el-table-column>
      <el-table-column label="Status" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 'ACTIVE' ? 'success' : 'info'">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="Actions" width="150" fixed="right" align="center">
        <template #default="{ row }">
          <el-button type="primary" :icon="Edit" @click="openDialog('edit', row)" />
          <el-button type="danger" :icon="Delete" @click="deleteVoucher(row.id)" />
        </template>
      </el-table-column>
    </el-table>

    <el-dialog
      v-model="dialogVisible"
      :title="dialogMode === 'add' ? 'Create Voucher' : 'Edit Voucher'"
      width="820px"
      destroy-on-close
      align-center
      class="voucher-dialog"
    >
      <el-form
        ref="ruleFormRef"
        :model="ruleForm"
        :rules="rules"
        label-width="140px"
        label-position="left"
      >
        <el-form-item label="Voucher Code" prop="code">
          <el-input v-model="ruleForm.code" placeholder="e.g., SAVE20" />
        </el-form-item>
        <el-form-item label="Name" prop="name">
          <el-input v-model="ruleForm.name" placeholder="e.g., Summer Sale" />
        </el-form-item>
        <el-form-item label="Description" prop="description">
          <el-input type="textarea" v-model="ruleForm.description" />
        </el-form-item>

        <el-divider content-position="left">Discount Settings</el-divider>

        <el-form-item label="Type" prop="discount_type">
          <el-radio-group v-model="ruleForm.discount_type">
            <el-radio label="FIXED">Fixed Amount</el-radio>
            <el-radio label="PERCENTAGE">Percentage</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="Discount Value" prop="discount_value">
          <el-input-number v-model="ruleForm.discount_value" :min="0" style="width: 100%" />
        </el-form-item>

        <el-form-item v-if="isPercentage" label="Max Discount" prop="max_discount_value">
          <el-input-number
            v-model="ruleForm.max_discount_value"
            :min="0"
            placeholder="Optional limit"
            style="width: 100%"
          />
        </el-form-item>

        <el-form-item label="Min Order Value" prop="min_order_value">
          <el-input-number v-model="ruleForm.min_order_value" :min="0" style="width: 100%" />
        </el-form-item>

        <el-divider content-position="left">Usage & Scope</el-divider>

        <el-form-item label="Total Quantity" prop="total_quantity">
          <el-input-number v-model="ruleForm.total_quantity" :min="1" style="width: 100%" />
        </el-form-item>
        <el-form-item label="Limit Per User" prop="usage_limit_per_user">
          <el-input-number v-model="ruleForm.usage_limit_per_user" :min="1" style="width: 100%" />
        </el-form-item>

        <el-form-item label="Validity Period" prop="dateRange">
          <el-date-picker
            v-model="ruleForm.dateRange"
            type="datetimerange"
            range-separator="To"
            start-placeholder="Start date"
            end-placeholder="End date"
            style="width: 100%"
            popper-class="main-color-popper"
            :teleported="false"
          />
        </el-form-item>

        <el-form-item label="Apply Scope" prop="apply_scope">
          <el-select
            v-model="ruleForm.apply_scope"
            style="width: 100%"
            popper-class="main-color-popper"
            :teleported="false"
          >
            <el-option label="All Products" value="ALL" />
            <el-option label="Specific Categories" value="CATEGORY" />
          </el-select>
        </el-form-item>

        <el-form-item
          v-if="ruleForm.apply_scope === 'CATEGORY'"
          label="Select Categories"
          prop="apply_seller_category_ids"
        >
          <el-select
            v-model="ruleForm.apply_seller_category_ids"
            multiple
            placeholder="Select categories"
            style="width: 100%"
            popper-class="main-color-popper"
            :teleported="false"
          >
            <el-option
              v-for="cat in categoryList"
              :key="cat.id"
              :label="cat.name"
              :value="cat.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">Cancel</el-button>
          <el-button type="primary" color="var(--main-color)" @click="handleSubmit(ruleFormRef)"
            >Confirm</el-button
          >
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<style>
.voucher-manager,
.main-color-popper {
  --el-color-primary: var(--main-color);
  --el-color-primary-light-9: #e8f5e9; /* Approximate light shade for hover backgrounds if needed, or let Element calculate */
  /* Re-declare main color loop if Element depends on light-x vars, but usually primary is enough for text/border */
}
/* Ensure the radio button inside the dialog (if teleported or not) gets it */
.el-dialog,
.el-message-box {
  --el-color-primary: var(--main-color);
}
.voucher-dialog .el-dialog__body {
  max-height: 65vh;
  overflow-y: auto;
  padding-top: 10px;
  padding-bottom: 10px;
  padding-right: 16px;
}
</style>

<style scoped>
.voucher-manager {
  padding: 20px;
  background: #fff;
  min-height: 100%;
}
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}
</style>
