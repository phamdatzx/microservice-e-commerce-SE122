<script setup lang="ts">
import { ref, watch, onMounted, computed } from 'vue'
import axios from 'axios'

interface Commune {
  code: string
  name: string
}

const props = defineProps<{
  modelValue: string
  provinceCode: string
  initialName?: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
  (e: 'change', commune: Commune | null): void
}>()

const communes = ref<Commune[]>([])
const loading = ref(false)

const hasInitialValue = computed(() => {
  return (
    props.modelValue &&
    props.initialName &&
    !communes.value.some((c) => c.code === props.modelValue)
  )
})

const fetchCommunes = async (code: string) => {
  if (!code) {
    communes.value = []
    return
  }
  loading.value = true
  try {
    const response = await axios.get(`/address-kit/2025-07-01/provinces/${code}/communes`)
    communes.value = response.data.communes

    // If we have an initial name and it's being used as a code, resolve the real code
    if (props.modelValue && props.initialName) {
      const match = communes.value.find((c) => c.name === props.initialName)
      if (match && match.code !== props.modelValue) {
        emit('update:modelValue', match.code)
        emit('change', match)
      }
    }
  } catch (error) {
    console.error('Failed to fetch communes:', error)
  } finally {
    loading.value = false
  }
}

const handleChange = (value: string) => {
  const selected = communes.value.find((c) => c.code === value) || null
  emit('change', selected)
}

watch(
  () => props.provinceCode,
  (newCode) => {
    fetchCommunes(newCode)
  },
)

onMounted(() => {
  if (props.provinceCode) {
    fetchCommunes(props.provinceCode)
  }
})
</script>

<template>
  <el-select
    :model-value="modelValue"
    size="large"
    @update:model-value="emit('update:modelValue', $event)"
    placeholder="Select Commune/Ward"
    filterable
    :loading="loading"
    :disabled="!provinceCode"
    style="width: 100%"
    @change="handleChange"
    no-data-text="Please select City/Province first"
  >
    <el-option v-if="hasInitialValue" :key="modelValue" :label="initialName" :value="modelValue" />
    <el-option v-for="item in communes" :key="item.code" :label="item.name" :value="item.code" />
  </el-select>
</template>
