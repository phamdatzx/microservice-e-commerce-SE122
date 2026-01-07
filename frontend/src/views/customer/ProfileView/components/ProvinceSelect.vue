<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import axios from 'axios'

interface Province {
  code: string
  name: string
}

const props = defineProps<{
  modelValue: string
  initialName?: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
  (e: 'change', province: Province | null): void
}>()

const provinces = ref<Province[]>([])
const loading = ref(false)

const hasInitialValue = computed(() => {
  return (
    props.modelValue &&
    props.initialName &&
    !provinces.value.some((p) => p.code === props.modelValue)
  )
})

const fetchProvinces = async () => {
  loading.value = true
  try {
    const response = await axios.get('/address-kit/2025-07-01/provinces')
    provinces.value = response.data.provinces

    // If we have an initial name and it's being used as a code, resolve the real code
    if (props.modelValue && props.initialName) {
      const match = provinces.value.find((p) => p.name === props.initialName)
      if (match && match.code !== props.modelValue) {
        emit('update:modelValue', match.code)
        emit('change', match)
      }
    }
  } catch (error) {
    console.error('Failed to fetch provinces:', error)
  } finally {
    loading.value = false
  }
}

const handleChange = (value: string) => {
  const selected = provinces.value.find((p) => p.code === value) || null
  emit('change', selected)
}

onMounted(fetchProvinces)
</script>

<template>
  <el-select
    size="large"
    :model-value="modelValue"
    @update:model-value="emit('update:modelValue', $event)"
    placeholder="Select Province/City"
    filterable
    :loading="loading"
    style="width: 100%"
    @change="handleChange"
  >
    <el-option v-if="hasInitialValue" :key="modelValue" :label="initialName" :value="modelValue" />
    <el-option v-for="item in provinces" :key="item.code" :label="item.name" :value="item.code" />
  </el-select>
</template>
