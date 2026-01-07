<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Plus, Edit, Delete, Check } from '@element-plus/icons-vue'
import axios from 'axios'
import ProvinceSelect from './ProvinceSelect.vue'
import CommuneSelect from './CommuneSelect.vue'
import { ElMessage, ElMessageBox } from 'element-plus'

interface Address {
  id: string
  user_id: string
  full_name: string
  phone: string
  address_line: string
  ward: string
  province: string
  country: string
  latitude: number
  longitude: number
  default: boolean
}

const addresses = ref<Address[]>([])

const fetchAddresses = async () => {
  try {
    const response = await axios.get('http://localhost:81/api/user/addresses', {
      headers: {
        Authorization: `Bearer ${localStorage.getItem('access_token')}`,
      },
    })
    if (response.data.status === 200) {
      addresses.value = response.data.data.sort((a: Address, b: Address) => {
        if (a.default) return -1
        if (b.default) return 1
        return 0
      })
    }
  } catch (error) {
    console.error('Failed to fetch addresses:', error)
  }
}

const addressDialogVisible = ref(false)
const isEditing = ref(false)
const addressForm = ref({
  id: '',
  name: '',
  phone: '',
  provinceCode: '',
  provinceName: '',
  communeCode: '',
  communeName: '',
  addressDetail: '',
  isDefault: false,
})

const openAddAddress = () => {
  isEditing.value = false
  addressForm.value = {
    id: '',
    name: '',
    phone: '',
    provinceCode: '',
    provinceName: '',
    communeCode: '',
    communeName: '',
    addressDetail: '',
    isDefault: false,
  }
  addressDialogVisible.value = true
}

const editAddress = (address: Address) => {
  isEditing.value = true
  addressForm.value = {
    id: address.id,
    name: address.full_name,
    phone: address.phone,
    addressDetail: address.address_line,
    provinceCode: address.province, // Use name as code for initial display
    provinceName: address.province,
    communeCode: address.ward,
    communeName: address.ward,
    isDefault: address.default,
  }
  addressDialogVisible.value = true
}

const handleProvinceChange = (province: any) => {
  if (province && province.name === addressForm.value.provinceName) return
  addressForm.value.provinceName = province ? province.name : ''
  addressForm.value.communeCode = ''
  addressForm.value.communeName = ''
}

const handleCommuneChange = (commune: any) => {
  addressForm.value.communeName = commune ? commune.name : ''
}

const handlePhoneInput = (value: string) => {
  addressForm.value.phone = value.replace(/\D/g, '')
}

const handleSaveAddress = async () => {
  // Validation
  if (
    !addressForm.value.name ||
    !addressForm.value.phone ||
    !addressForm.value.provinceName ||
    !addressForm.value.communeName ||
    !addressForm.value.addressDetail
  ) {
    ElMessage.warning('Please complete all required fields')
    return
  }

  // Phone validation
  if (!/^0\d{9}$/.test(addressForm.value.phone)) {
    ElMessage.warning('Phone number must be 10 digits starting with 0')
    return
  }

  const payload = {
    full_name: addressForm.value.name,
    phone: addressForm.value.phone,
    address_line: addressForm.value.addressDetail,
    province: addressForm.value.provinceName,
    district: 'quan',
    ward: addressForm.value.communeName,
    country: 'viet nam',
    latitude: 1,
    longitude: 1,
    default: addressForm.value.isDefault,
  }

  try {
    const headers = {
      Authorization: `Bearer ${localStorage.getItem('access_token')}`,
    }
    if (isEditing.value) {
      await axios.put(`http://localhost:81/api/user/addresses/${addressForm.value.id}`, payload, {
        headers,
      })
    } else {
      await axios.post('http://localhost:81/api/user/addresses', payload, { headers })
    }
    await fetchAddresses()
    addressDialogVisible.value = false
  } catch (error) {
    console.error('Failed to save address:', error)
  }
}

const deleteAddress = async (id: string) => {
  try {
    await ElMessageBox.confirm('Are you sure you want to delete this address?', 'Warning', {
      confirmButtonText: 'OK',
      cancelButtonText: 'Cancel',
      type: 'warning',
    })

    await axios.delete(`http://localhost:81/api/user/addresses/${id}`, {
      headers: {
        Authorization: `Bearer ${localStorage.getItem('access_token')}`,
      },
    })
    await fetchAddresses()
    ElMessage.success('Address deleted successfully')
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete address:', error)
    }
  }
}

const setDefaultAddress = async (address: Address) => {
  const payload = {
    full_name: address.full_name,
    phone: address.phone,
    address_line: address.address_line,
    province: address.province,
    district: 'quan',
    ward: address.ward,
    country: 'viet nam',
    latitude: 1,
    longitude: 1,
    default: true,
  }

  try {
    await axios.put(`http://localhost:81/api/user/addresses/${address.id}`, payload, {
      headers: {
        Authorization: `Bearer ${localStorage.getItem('access_token')}`,
      },
    })
    await fetchAddresses()
  } catch (error) {
    console.error('Failed to set default address:', error)
  }
}

onMounted(fetchAddresses)
</script>

<template>
  <div class="my-address">
    <div style="display: flex; justify-content: space-between; align-items: center">
      <h2 class="section-title" style="margin-bottom: 0">My Address</h2>
      <el-button
        type="primary"
        @click="openAddAddress"
        size="large"
        :icon="Plus"
        style="padding: 22px 40px; font-size: 16px"
        >Add New Address
      </el-button>
    </div>
    <el-divider />

    <div class="address-list">
      <div v-for="address in addresses" :key="address.id" class="address-item">
        <div class="address-info">
          <div class="header">
            <span class="name">{{ address.full_name }}</span>
            <el-divider direction="vertical" />
            <span class="phone">{{ address.phone }}</span>
            <el-tag
              v-if="address.default"
              type="success"
              effect="dark"
              style="margin-left: 10px; position: relative; top: -2px"
              >Default</el-tag
            >
          </div>
          <div class="details">
            <p>{{ address.address_line }}</p>
            <p>{{ address.ward }}, {{ address.province }}</p>
          </div>
        </div>
        <div class="address-actions">
          <div class="top-actions">
            <el-button type="primary" :icon="Edit" @click="editAddress(address)" />
            <el-button
              type="danger"
              :icon="Delete"
              @click="deleteAddress(address.id)"
              :disabled="address.default"
            />
          </div>
          <el-button
            v-if="!address.default"
            size="default"
            class="set-default-btn"
            style="margin-left: 0; margin-top: 16px"
            @click="setDefaultAddress(address)"
          >
            Set as Default
          </el-button>
        </div>
      </div>
    </div>

    <!-- Address Dialog -->
    <el-dialog
      v-model="addressDialogVisible"
      :title="isEditing ? 'Update Address' : 'New Address'"
      width="500px"
      destroy-on-close
    >
      <el-form :model="addressForm" label-position="top">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="Full Name" required>
              <el-input v-model="addressForm.name" placeholder="Mark Cole" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Phone Number" required>
              <el-input
                v-model="addressForm.phone"
                placeholder="0123456789"
                @input="handlePhoneInput"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="City/Province" required>
              <ProvinceSelect
                v-model="addressForm.provinceCode"
                :initial-name="addressForm.provinceName"
                @change="handleProvinceChange"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Commune/Ward" required>
              <CommuneSelect
                v-model="addressForm.communeCode"
                :province-code="addressForm.provinceCode"
                :initial-name="addressForm.communeName"
                @change="handleCommuneChange"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="Specific Address" required>
          <el-input
            v-model="addressForm.addressDetail"
            type="textarea"
            rows="2"
            placeholder="House number, street name..."
          />
        </el-form-item>
        <el-form-item v-if="!isEditing">
          <el-checkbox v-model="addressForm.isDefault">Set as default address</el-checkbox>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="addressDialogVisible = false" size="large">Cancel</el-button>
          <el-button type="primary" @click="handleSaveAddress" size="large">
            {{ isEditing ? 'Save' : 'Add' }}
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.section-title {
  font-size: 28px;
  font-weight: 700;
  margin-bottom: 30px;
  color: #000;
}

.address-item {
  display: flex;
  justify-content: space-between;
  padding: 20px 0;
  border-bottom: 1px solid #f0f0f0;
}

.address-item:last-child {
  border-bottom: none;
}

.address-info .header {
  margin-bottom: 8px;
}

.address-info .name {
  font-weight: 700;
  font-size: 16px;
}

.address-info .phone {
  color: #666;
  font-size: 14px;
}

.address-info .details {
  color: #666;
  font-size: 14px;
  line-height: 1.6;
}

.address-actions {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  justify-content: center;
}

.top-actions {
  margin-bottom: 0;
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

.set-default-btn {
  border-color: #dcdfe6;
  color: #606266;
}

.set-default-btn:not(:disabled):hover {
  background-color: #f5f5f5;
  border-color: #c0c4cc;
}
</style>
