<script setup lang="ts">
import { ref } from 'vue'
import { Plus } from '@element-plus/icons-vue'

interface Address {
  id: number
  name: string
  phone: string
  addressDetail: string
  city: string
  isDefault: boolean
  type: 'Home' | 'Office'
}

const addresses = ref<Address[]>([
  {
    id: 1,
    name: 'Mark Cole',
    phone: '(+84) 345 967 735',
    addressDetail: '1/14/7, Binh Duong 2 Street, An Binh Ward',
    city: 'Di An City, Binh Duong Province',
    isDefault: true,
    type: 'Home',
  },
  {
    id: 2,
    name: 'Mark Cole',
    phone: '(+84) 345 967 735',
    addressDetail: '3 Road No. 2, Linh Trung Ward',
    city: 'Thu Duc City, Ho Chi Minh City',
    isDefault: false,
    type: 'Office',
  },
])

const addressDialogVisible = ref(false)
const isEditing = ref(false)
const addressForm = ref({
  name: '',
  phone: '',
  city: '',
  addressDetail: '',
  type: 'Home' as 'Home' | 'Office',
  isDefault: false,
})

const openAddAddress = () => {
  isEditing.value = false
  addressForm.value = {
    name: '',
    phone: '',
    city: '',
    addressDetail: '',
    type: 'Home',
    isDefault: false,
  }
  addressDialogVisible.value = true
}

const editAddress = (address: Address) => {
  isEditing.value = true
  addressForm.value = { ...address }
  addressDialogVisible.value = true
}
</script>

<template>
  <div class="my-address">
    <div style="display: flex; justify-content: space-between; align-items: center">
      <h2 class="section-title" style="margin-bottom: 0">My Address</h2>
      <el-button type="success" @click="openAddAddress" size="large" :icon="Plus"
        >Add New Address</el-button
      >
    </div>
    <el-divider />

    <div class="address-list">
      <div v-for="address in addresses" :key="address.id" class="address-item">
        <div class="address-info">
          <div class="header">
            <span class="name">{{ address.name }}</span>
            <el-divider direction="vertical" />
            <span class="phone">{{ address.phone }}</span>
          </div>
          <div class="details">
            <p>{{ address.addressDetail }}</p>
            <p>{{ address.city }}</p>
          </div>
          <div class="tags">
            <span v-if="address.isDefault" class="tag default">Default</span>
          </div>
        </div>
        <div class="address-actions">
          <div class="top-actions">
            <el-button type="primary" link @click="editAddress(address)">Update</el-button>
            <el-button v-if="!address.isDefault" type="primary" link>Delete</el-button>
          </div>
          <el-button
            :disabled="address.isDefault"
            size="default"
            class="set-default-btn"
            @click="address.isDefault = true"
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
              <el-input v-model="addressForm.phone" placeholder="(+84) 345 967 735" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="City/Province, District, Ward" required>
          <el-select
            v-model="addressForm.city"
            placeholder="Select your location"
            style="width: 100%"
          >
            <el-option label="Ho Chi Minh City" value="Ho Chi Minh City" />
            <el-option label="Ha Noi" value="Ha Noi" />
            <el-option label="Binh Duong" value="Binh Duong" />
          </el-select>
        </el-form-item>

        <el-form-item label="Specific Address" required>
          <el-input
            v-model="addressForm.addressDetail"
            type="textarea"
            rows="2"
            placeholder="House number, street name..."
          />
        </el-form-item>

        <el-form-item label="Address Type">
          <el-radio-group v-model="addressForm.type">
            <el-radio-button
              style="--el-radio-button-checked-bg-color: var(--main-color)"
              label="Home"
              value="Home"
            />
            <el-radio-button
              style="--el-radio-button-checked-bg-color: var(--main-color)"
              label="Office"
              value="Office"
            />
          </el-radio-group>
        </el-form-item>

        <el-form-item>
          <el-checkbox
            v-model="addressForm.isDefault"
            style="
              --el-checkbox-checked-text-color: var(--main-color);
              --el-checkbox-checked-bg-color: var(--main-color);
            "
            >Set as default address</el-checkbox
          >
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="addressDialogVisible = false" size="large">Cancel</el-button>
          <el-button type="success" @click="addressDialogVisible = false" size="large">
            Complete
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

.tags {
  margin-top: 8px;
}

.tag {
  display: inline-block;
  padding: 2px 8px;
  font-size: 12px;
  border-radius: 4px;
}

.tag.default {
  color: var(--main-color);
  border: 1px solid var(--main-color);
}

.address-actions {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  justify-content: space-between;
}

.top-actions {
  margin-bottom: 10px;
}

.set-default-btn {
  border-color: #dcdfe6;
  color: #606266;
}

.set-default-btn:not(:disabled):hover {
  background-color: #f5f5f5;
  border-color: #c0c4cc;
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
</style>
