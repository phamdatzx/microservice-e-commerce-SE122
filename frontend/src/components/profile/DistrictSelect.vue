<script setup lang="ts">
import { ref, watch, onMounted, computed } from 'vue'
import axios from 'axios'

interface District {
  DistrictID: number
  ProvinceID: number
  DistrictName: string
}

const props = defineProps<{
  modelValue: number | string
  provinceId: number | string
  initialName?: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: number): void
  (e: 'change', district: District | null): void
}>()

const districts = ref<District[]>([])
const loading = ref(false)
const token = import.meta.env.VITE_GHN_TOKEN

const hasInitialValue = computed(() => {
  return (
    props.modelValue &&
    props.initialName &&
    !districts.value.some((d) => d.DistrictID === Number(props.modelValue))
  )
})

const fetchDistricts = async (provinceId: number | string) => {
  if (!provinceId) {
    districts.value = []
    return
  }
  loading.value = true
  try {
    const response = await axios.get(
      `${import.meta.env.VITE_GHN_URL}/shiip/public-api/master-data/district`,
      {
        headers: { token },
        params: { province_id: provinceId },
      },
    )

    if (response.data && response.data.data) {
      districts.value = response.data.data
    }

    if (props.initialName) {
      const match = districts.value.find((d) => d.DistrictName === props.initialName)
      if (match && match.DistrictID !== Number(props.modelValue)) {
        emit('update:modelValue', match.DistrictID)
        emit('change', match)
      }
    }
  } catch (error) {
    console.error('Failed to fetch districts:', error)
  } finally {
    loading.value = false
  }
}

const handleChange = (value: number) => {
  const selected = districts.value.find((d) => d.DistrictID === value) || null
  emit('change', selected)
  emit('update:modelValue', value)
}

watch(
  () => props.provinceId,
  (newId) => {
    fetchDistricts(newId)
  },
)

onMounted(() => {
  if (props.provinceId) {
    fetchDistricts(props.provinceId)
  }
})
</script>

<template>
  <el-select
    :model-value="modelValue"
    size="large"
    @update:model-value="handleChange"
    placeholder="Select District"
    filterable
    :loading="loading"
    :disabled="!provinceId"
    style="width: 100%"
    no-data-text="Please select City/Province first"
  >
    <el-option v-if="hasInitialValue" :key="modelValue" :label="initialName" :value="modelValue" />
    <el-option
      v-for="item in districts"
      :key="item.DistrictID"
      :label="item.DistrictName"
      :value="item.DistrictID"
    />
  </el-select>
</template>
