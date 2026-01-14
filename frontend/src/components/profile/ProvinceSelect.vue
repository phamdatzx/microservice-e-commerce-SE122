<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import axios from 'axios'

interface Province {
  ProvinceID: number
  ProvinceName: string
  NameExtension: string[]
}

const props = defineProps<{
  modelValue: number | string
  initialName?: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: number): void
  (e: 'change', province: Province | null): void
}>()

const provinces = ref<Province[]>([])
const loading = ref(false)
const token = import.meta.env.VITE_GHN_TOKEN

const hasInitialValue = computed(() => {
  return (
    props.modelValue &&
    props.initialName &&
    !provinces.value.some((p) => p.ProvinceID === Number(props.modelValue))
  )
})

const fetchProvinces = async () => {
  loading.value = true
  try {
    const response = await axios.get(
      `${import.meta.env.VITE_GHN_URL}/shiip/public-api/master-data/province`,
      {
        headers: { token },
      },
    )
    // Filter provinces that have NameExtension as per user instruction
    if (response.data && response.data.data) {
      provinces.value = response.data.data.filter(
        (p: any) => p.NameExtension && p.NameExtension.length > 0,
      )
    }

    // If we have an initial name, try to find matching ID if modelValue is missing or mismatch
    if (props.initialName) {
      const match = provinces.value.find((p) => p.ProvinceName === props.initialName)
      if (match && match.ProvinceID !== Number(props.modelValue)) {
        emit('update:modelValue', match.ProvinceID)
        emit('change', match)
      }
    }
  } catch (error) {
    console.error('Failed to fetch provinces:', error)
  } finally {
    loading.value = false
  }
}

const handleChange = (value: number) => {
  const selected = provinces.value.find((p) => p.ProvinceID === value) || null
  emit('change', selected)
  emit('update:modelValue', value)
}

onMounted(fetchProvinces)
</script>

<template>
  <el-select
    size="large"
    :model-value="modelValue"
    @update:model-value="handleChange"
    placeholder="Select Province/City"
    filterable
    :loading="loading"
    style="width: 100%"
  >
    <el-option v-if="hasInitialValue" :key="modelValue" :label="initialName" :value="modelValue" />
    <el-option
      v-for="item in provinces"
      :key="item.ProvinceID"
      :label="item.ProvinceName"
      :value="item.ProvinceID"
    />
  </el-select>
</template>
