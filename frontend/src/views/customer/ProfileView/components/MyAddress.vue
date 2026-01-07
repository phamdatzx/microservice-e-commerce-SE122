<script setup lang="ts">
import { ref, onMounted, nextTick, watch } from 'vue'
import { Plus, Edit, Delete, Check, Search } from '@element-plus/icons-vue'
import axios from 'axios'
import ProvinceSelect from './ProvinceSelect.vue'
import CommuneSelect from './CommuneSelect.vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'

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
  latitude: 10.762622,
  longitude: 106.660172,
})

const mapContainer = ref<HTMLElement | null>(null)
let map: L.Map | null = null
let marker: L.Marker | null = null
const searchQuery = ref('')

const searchLocation = async () => {
  if (!searchQuery.value) return

  try {
    const response = await fetch(
      `https://nominatim.openstreetmap.org/search?format=json&q=${encodeURIComponent(
        searchQuery.value,
      )}&countrycodes=vn`,
    )
    const data = await response.json()

    if (data && data.length > 0) {
      const { lat, lon } = data[0]
      const newLat = parseFloat(lat)
      const newLng = parseFloat(lon)

      addressForm.value.latitude = newLat
      addressForm.value.longitude = newLng

      if (map && marker) {
        map.setView([newLat, newLng], 16)
        marker.setLatLng([newLat, newLng])
      }
      ElMessage.success('Location found')
    } else {
      ElMessage.warning('Location not found')
    }
  } catch (error) {
    console.error('Search failed:', error)
    ElMessage.error('Failed to search location')
  }
}

const querySearchAsync = async (queryString: string, cb: (arg: any) => void) => {
  if (!queryString) {
    cb([])
    return
  }

  try {
    const response = await fetch(
      `https://nominatim.openstreetmap.org/search?format=json&q=${encodeURIComponent(
        queryString,
      )}&limit=5&countrycodes=vn`,
    )
    const data = await response.json()
    const results = data.map((item: any) => ({
      value: item.display_name,
      lat: item.lat,
      lon: item.lon,
    }))
    cb(results)
  } catch (error) {
    console.error('Autocomplete search failed:', error)
    cb([])
  }
}

const handleSelect = (item: any) => {
  const newLat = parseFloat(item.lat)
  const newLng = parseFloat(item.lon)

  addressForm.value.latitude = newLat
  addressForm.value.longitude = newLng

  if (map && marker) {
    map.setView([newLat, newLng], 16)
    marker.setLatLng([newLat, newLng])
  }
}

const initMap = () => {
  if (!mapContainer.value) return

  if (map) {
    map.remove()
  }

  const lat = addressForm.value.latitude || 10.762622
  const lng = addressForm.value.longitude || 106.660172

  map = L.map(mapContainer.value).setView([lat, lng], 13)
  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution:
      '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
  }).addTo(map)

  // Fix marker icon issue in Webpack/Vite
  const icon = L.icon({
    iconUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png',
    shadowUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-shadow.png',
    iconSize: [25, 41],
    iconAnchor: [12, 41],
  })

  marker = L.marker([lat, lng], { draggable: true, icon }).addTo(map)

  marker.on('dragend', () => {
    if (marker) {
      const position = marker.getLatLng()
      addressForm.value.latitude = position.lat
      addressForm.value.longitude = position.lng
    }
  })

  map.on('click', (e) => {
    if (marker) {
      marker.setLatLng(e.latlng)
      addressForm.value.latitude = e.latlng.lat
      addressForm.value.longitude = e.latlng.lng
    }
  })
}

watch(addressDialogVisible, (newVal) => {
  if (newVal) {
    nextTick(() => {
      initMap()
    })
  }
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
    latitude: 10.762622,
    longitude: 106.660172,
  }
  searchQuery.value = ''
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
    latitude: address.latitude || 10.762622,
    longitude: address.longitude || 106.660172,
  }
  searchQuery.value = ''
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
    latitude: addressForm.value.latitude,
    longitude: addressForm.value.longitude,
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
      ElMessage.success('Address updated successfully')
    } else {
      await axios.post('http://localhost:81/api/user/addresses', payload, { headers })
      ElMessage.success('Address added successfully')
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
      width="600px"
      align-center
      destroy-on-close
    >
      <div class="dialog-content">
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

          <el-form-item label="Location" required>
            <div style="margin-bottom: 10px; width: 100%">
              <el-autocomplete
                v-model="searchQuery"
                :fetch-suggestions="querySearchAsync"
                placeholder="Search for a location..."
                @select="handleSelect"
                :trigger-on-focus="false"
                style="width: 100%"
              />
            </div>
            <div ref="mapContainer" style="height: 300px; width: 100%"></div>
            <div style="margin-top: 10px; font-size: 12px; color: #666">
              Click on map or drag marker to select location
            </div>
          </el-form-item>

          <el-form-item v-if="!isEditing">
            <el-checkbox v-model="addressForm.isDefault">Set as default address</el-checkbox>
          </el-form-item>
        </el-form>
      </div>
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

.dialog-content {
  max-height: 64vh;
  overflow-y: auto;
  overflow-x: hidden;
  padding-left: 25px;
  padding-right: 30px;
}
</style>
