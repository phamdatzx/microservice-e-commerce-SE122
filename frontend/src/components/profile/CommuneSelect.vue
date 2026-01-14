<script setup lang="ts">
import { ref, watch, onMounted, computed } from 'vue'
import axios from 'axios'
import { ElMessage } from 'element-plus'

interface Ward {
  WardCode: string
  DistrictID: number
  WardName: string
}

const props = defineProps<{
  modelValue: string
  districtId: number | string
  initialName?: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
  (e: 'change', ward: Ward | null): void
}>()

const wards = ref<Ward[]>([])
const loading = ref(false)
const token = import.meta.env.VITE_GHN_TOKEN

const hasInitialValue = computed(() => {
  return (
    props.modelValue &&
    props.initialName &&
    !wards.value.some((w) => w.WardCode === props.modelValue)
  )
})

const fetchWards = async (districtId: number | string) => {
  if (!districtId) {
    wards.value = []
    return
  }

  loading.value = true
  try {
    const response = await axios.get(
      `${import.meta.env.VITE_GHN_URL}/shiip/public-api/master-data/ward`,
      {
        headers: { token },
        params: { district_id: districtId },
      },
    )

    if (response.data && response.data.data) {
      wards.value = response.data.data
      if (wards.value.length === 0) {
        ElMessage.warning('No wards found for this district')
      }
    } else {
      wards.value = []
      ElMessage.warning('No wards found for this district')
    }

    // If we have an initial name and it's being used as a code, resolve the real code
    if (props.modelValue && props.initialName) {
      const match = wards.value.find((w) => w.WardName === props.initialName)
      if (match && match.WardCode !== props.modelValue) {
        emit('update:modelValue', match.WardCode)
        emit('change', match)
      }
    }
  } catch (error) {
    console.error('Failed to fetch wards:', error)
    ElMessage.error('Failed to fetch wards')
  } finally {
    loading.value = false
  }
}

const handleChange = (value: string) => {
  const selected = wards.value.find((w) => w.WardCode === value) || null
  emit('change', selected)
  emit('update:modelValue', value)
}

watch(
  () => props.districtId,
  (newId) => {
    fetchWards(newId)
  },
)

onMounted(() => {
  if (props.districtId) {
    fetchWards(props.districtId)
  }
})
</script>

<template>
  <el-select
    :model-value="modelValue"
    size="large"
    @update:model-value="handleChange"
    placeholder="Select Commune/Ward"
    filterable
    :loading="loading"
    :disabled="!districtId"
    style="width: 100%"
    no-data-text="Please select District first"
  >
    <el-option v-if="hasInitialValue" :key="modelValue" :label="initialName" :value="modelValue" />
    <el-option
      v-for="item in wards"
      :key="item.WardCode"
      :label="item.WardName"
      :value="item.WardCode"
    />
  </el-select>
</template>
