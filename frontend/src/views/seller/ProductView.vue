<script setup lang="ts">
import {
  Delete,
  Edit,
  Plus,
  ZoomIn,
  ZoomOut,
  Back,
  DArrowRight,
  Download,
  Refresh,
  RefreshLeft,
  RefreshRight,
  Right,
  Search,
} from '@element-plus/icons-vue'
import axios from 'axios'
import {
  ElLoading,
  ElMessage,
  ElMessageBox,
  ElNotification,
  type FormInstance,
  type FormRules,
  type UploadFile,
  type UploadFiles,
  type UploadProps,
  type UploadUserFile,
} from 'element-plus'
import { onMounted, reactive, ref, watch, computed } from 'vue'

interface ProductImage {
  id: string
  url: string
  order: number
}

interface OptionGroup {
  key: string
  values: string[]
}

interface ProductVariant {
  sku: string
  options: Record<string, string>
  price: number
  stock: number
  images: UploadUserFile[]
}

interface Product {
  id: string
  name: string
  description: string
  images: ProductImage[]
  status: string
  seller_id: string
  rating: number
  rate_count: number
  sold_count: number
  is_active: boolean
  created_at: Date
  updated_at: Date
  option_groups: OptionGroup[]
  variants: ProductVariant[]
  category_ids: string[]
  seller_category_ids: string[]
}

type DialogMode = 'add' | 'edit'

const userId = localStorage.getItem('user_id')
const token = localStorage.getItem('access_token')

const productData = ref<Product[]>([])
const isLoading = ref(false)
const searchQuery = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const sortBy = ref('name')
const sortDirection = ref('asc')
const filterStatus = ref('')

onMounted(() => {
  fetchProducts()
})

const fetchProducts = () => {
  isLoading.value = true
  axios
    .get(
      import.meta.env.VITE_BE_API_URL +
        '/product/public/products/seller/' +
        userId +
        `?status=${filterStatus.value}&search=${searchQuery.value}&sort_by=${sortBy.value}&sort_direction=${sortDirection.value}&page=${currentPage.value}&limit=${pageSize.value}`,
    )
    .then((response) => {
      productData.value = response.data.products
      total.value = response.data.pagination.total_items
    })
    .catch((error) => {
      ElNotification({
        title: 'Error!',
        message: 'Fetch products failed: ' + error.message,
        type: 'error',
      })
    })
    .finally(() => {
      isLoading.value = false
    })
}

const generalCategoryOptions = ref<{ value: string; label: string }[]>([])
const sellerCategoryOptions = ref<{ value: string; label: string }[]>([])

const dialogVisible = ref(false)
const dialogMode = ref<DialogMode>('add')
const dialogContent = ref({
  title: 'Add product',
  mainBtnText: 'Add',
})

const previewImageDialogVisible = ref(false)
const previewImageUrl = ref('')
// Categories will be fetched from API in a real app
onMounted(() => {
  fetchGeneralCategories()
  fetchSellerCategories()
})

const fetchGeneralCategories = () => {
  axios
    .get(import.meta.env.VITE_BE_API_URL + '/product/public/category')
    .then((res) => {
      generalCategoryOptions.value = res.data.map((c: any) => ({
        value: c.id,
        label: c.name,
      }))
    })
    .catch(() => {
      // Fallback
      generalCategoryOptions.value = [
        { value: '1', label: 'Fashion' },
        { value: '2', label: 'Electronics' },
      ]
    })
}

const fetchSellerCategories = () => {
  axios
    .get(import.meta.env.VITE_BE_API_URL + '/product/public/seller/' + userId + '/category')
    .then((res) => {
      sellerCategoryOptions.value = res.data.map((c: any) => ({
        value: c.id,
        label: c.name,
      }))
    })
    .catch(() => {
      // Fallback
      sellerCategoryOptions.value = [
        { value: '1', label: 'Shirts' },
        { value: '2', label: 'Pants' },
      ]
    })
}

//#region FORM
interface RuleForm {
  id?: string
  images: UploadUserFile[]
  name: string
  category_ids: string[]
  seller_category_ids: string[]
  description: string
  status: string
  is_active: boolean
  option_groups: OptionGroup[]
  variants: ProductVariant[]
}

const ruleFormRef = ref<FormInstance>()
const ruleForm = reactive<RuleForm>({
  id: undefined,
  images: [],
  name: '',
  category_ids: [],
  seller_category_ids: [],
  description: '',
  status: 'AVAILABLE',
  is_active: true,
  option_groups: [],
  variants: [],
})

const rules = reactive<FormRules<RuleForm>>({})

const submitForm = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  await formEl.validate((valid, fields) => {
    if (valid) {
      handleFormSubmit()
    } else {
      console.log('error submit!', fields)
    }
  })
}

const handleFormSubmit = () => {
  if (dialogMode.value === 'add') {
    handleAddProduct()
  } else if (dialogMode.value === 'edit') {
    handleEditProduct()
  }
}

const clearRuleForm = () => {
  ruleFormRef.value?.resetFields()
  ruleForm.id = undefined
  ruleForm.name = ''
  ruleForm.description = ''
  ruleForm.status = 'AVAILABLE'
  ruleForm.is_active = true
  ruleForm.option_groups = []
  ruleForm.variants = []
  ruleForm.category_ids = []
  ruleForm.seller_category_ids = []
  ruleForm.images = []
}
//#endregion

//#region IMAGES UPLOAD
const handleChange = (uploadFile: UploadFile, uploadFiles: UploadFiles) => {
  const rawFile = uploadFile.raw

  if (!rawFile) return

  const isImage = rawFile.type.startsWith('image/')
  if (!isImage) {
    ElMessage.error('Only .png, .jpg and .jpeg files are allowed')

    // Remove the invalid file
    const index = ruleForm.images.findIndex((f) => f.uid === uploadFile.uid)
    if (index !== -1) ruleForm.images.splice(index, 1)
    return
  }
}

const handleExceed: UploadProps['onExceed'] = (files, uploadFiles) => {
  ElMessage.warning(
    `The limit is 9, you selected ${files.length} files this time, add up to ${
      files.length + uploadFiles.length
    } totally`,
  )
}

const handleRemove: UploadProps['onRemove'] = (uploadFile, uploadFiles) => {
  console.log(uploadFile, uploadFiles)
}

const handleVariantImageChange = (uploadFile: UploadFile, row: ProductVariant) => {
  const rawFile = uploadFile.raw
  if (!rawFile) return

  const isImage = rawFile.type.startsWith('image/')
  if (!isImage) {
    ElMessage.error('Only .png, .jpg and .jpeg files are allowed')
    row.images = [] // Clear invalid file
    return
  }
}

const handlePictureCardPreview: UploadProps['onPreview'] = (uploadFile) => {
  previewImageUrl.value = uploadFile.url!
  previewImageDialogVisible.value = true
}

const uploadProductImages = async (productId: string, images: UploadUserFile[]) => {
  if (images.length === 0) return

  const formData = new FormData()
  images.forEach((img) => {
    if (img.raw) {
      formData.append('image', img.raw)
    }
  })

  await axios.post(`${import.meta.env.VITE_BE_API_URL}/product/${productId}/images`, formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
      Authorization: `Bearer ${token}`,
    },
  })
}

const uploadVariantImages = async (
  productId: string,
  variants: ProductVariant[],
  serverVariants: any[],
) => {
  const formData = new FormData()
  let hasImages = false

  variants.forEach((v) => {
    if (v.images && v.images.length > 0 && v.images[0].raw) {
      // Find server variant ID by SKU
      const serverV = serverVariants.find((sv) => sv.sku === v.sku)
      if (serverV) {
        formData.append(serverV.id, v.images[0].raw)
        hasImages = true
      }
    }
  })

  if (!hasImages) return

  await axios.post(
    `${import.meta.env.VITE_BE_API_URL}/product/${productId}/variants/images`,
    formData,
    {
      headers: {
        'Content-Type': 'multipart/form-data',
        Authorization: `Bearer ${token}`,
      },
    },
  )
}
//#endregion
const handleAddProduct = () => {
  const loading = ElLoading.service({
    lock: true,
    text: 'Adding product and uploading images...',
    background: 'rgba(0, 0, 0, 0.7)',
  })

  // Format the payload according to the user example
  const payload = {
    name: ruleForm.name,
    description: ruleForm.description,
    status: ruleForm.status,
    is_active: ruleForm.is_active,
    option_groups: ruleForm.option_groups.map((g) => ({
      key: g.key,
      values: g.values,
    })),
    variants: ruleForm.variants.map((v) => ({
      sku: v.sku,
      options: v.options,
      price: v.price,
      stock: v.stock,
    })),
    category_ids: ruleForm.category_ids,
    seller_category_ids: ruleForm.seller_category_ids,
  }

  axios
    .post(import.meta.env.VITE_BE_API_URL + '/product', payload, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
    .then(async (response: any) => {
      const productId = response.data.id
      const serverVariants = response.data.variants

      try {
        // Sequentially upload images
        await uploadProductImages(productId, ruleForm.images)
        await uploadVariantImages(productId, ruleForm.variants, serverVariants)

        ElNotification({
          title: 'Success!',
          message: 'Product added successfully with images',
          type: 'success',
        })
        dialogVisible.value = false
        fetchProducts()
      } catch (uploadError: any) {
        ElNotification({
          title: 'Partial Success',
          message: 'Product added, but image upload failed: ' + uploadError.message,
          type: 'warning',
        })
        dialogVisible.value = false
        fetchProducts()
      }
    })
    .catch((error) => {
      ElNotification({
        title: 'Error!',
        message: 'Add product failed: ' + error.message,
        type: 'error',
      })
    })
    .finally(() => {
      loading.close()
    })
}

const handleEditProduct = () => {
  const loading = ElLoading.service({
    lock: true,
    text: 'Saving changes and updating images...',
    background: 'rgba(0, 0, 0, 0.7)',
  })

  // Format the payload according to the user example
  const payload = {
    name: ruleForm.name,
    description: ruleForm.description,
    status: ruleForm.status,
    is_active: ruleForm.is_active,
    option_groups: ruleForm.option_groups.map((g) => ({
      key: g.key,
      values: g.values,
    })),
    variants: ruleForm.variants.map((v) => ({
      sku: v.sku,
      options: v.options,
      price: v.price,
      stock: v.stock,
    })),
    category_ids: ruleForm.category_ids,
    seller_category_ids: ruleForm.seller_category_ids,
  }

  const productId = ruleForm.id

  axios
    .put(import.meta.env.VITE_BE_API_URL + '/product/seller/products/' + productId, payload, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
    .then(async (response: any) => {
      const serverVariants = response.data.variants

      try {
        // Sequentially upload images
        // We only upload images that have a 'raw' property (newly selected files)
        const newImages = ruleForm.images.filter((img) => !!img.raw)
        if (newImages.length > 0) {
          await uploadProductImages(productId!, newImages)
        }

        await uploadVariantImages(productId!, ruleForm.variants, serverVariants)

        ElNotification({
          title: 'Success!',
          message: 'Product updated successfully',
          type: 'success',
        })
        dialogVisible.value = false
        fetchProducts()
      } catch (uploadError: any) {
        ElNotification({
          title: 'Partial Success',
          message: 'Product details updated, but image upload failed: ' + uploadError.message,
          type: 'warning',
        })
        dialogVisible.value = false
        fetchProducts()
      }
    })
    .catch((error) => {
      ElNotification({
        title: 'Error!',
        message: 'Update product failed: ' + error.message,
        type: 'error',
      })
    })
    .finally(() => {
      loading.close()
    })
}

const handleSortChange = ({ prop, order }: { prop: string; order: string }) => {
  if (order === 'ascending') {
    sortDirection.value = 'asc'
  } else if (order === 'descending') {
    sortDirection.value = 'desc'
  } else {
    // defaults
    sortBy.value = 'name'
    sortDirection.value = 'asc'
    fetchProducts()
    return
  }

  // Map prop to API sort keys if needed
  if (prop === 'quantity') {
    sortBy.value = 'stock'
  } else {
    sortBy.value = prop
  }

  fetchProducts()
}

const handleFilterChange = (filters: any) => {
  if (filters.status && filters.status.length > 0) {
    filterStatus.value = filters.status.join(',')
  } else {
    filterStatus.value = ''
  }
  fetchProducts()
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
  currentPage.value = 1
  fetchProducts()
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
  fetchProducts()
}

const getStatusTagContent = (status: string) => {
  if (status === 'AVAILABLE') return 'Available'
  else if (status === 'OUT_OF_STOCK') return 'Out of stock'
  else if (status === 'HIDDEN') return 'Hidden'
}

const getStatusTagType = (status: string) => {
  // if (status === 'AVAILABLE') return 'success'
  // else if (status === 'OUT_OF_STOCK') return 'danger'
  // else if (status === 'HIDDEN') return 'info'
  return ''
}

const filterTag = (value: string, row: any) => {
  return row.status === value
}

const getRowClassName = ({ row }: { row: Product }) => {
  if (!row.variants || row.variants.length <= 1) {
    return 'hide-expand'
  }
  return ''
}

//#region MODAL
const openModal = (mode: DialogMode, product?: Product) => {
  dialogMode.value = mode
  dialogVisible.value = true
  if (mode === 'add') {
    dialogContent.value = {
      title: 'Add product',
      mainBtnText: 'Add',
    }
    clearRuleForm()
  } else if (mode === 'edit' && product) {
    dialogContent.value = {
      title: 'Edit product',
      mainBtnText: 'Save',
    }
    ruleForm.id = product.id
    ruleForm.name = product.name
    ruleForm.description = product.description
    ruleForm.status = product.status
    ruleForm.is_active = product.is_active
    ruleForm.category_ids = product.category_ids
    ruleForm.seller_category_ids = product.seller_category_ids
    ruleForm.option_groups = (product.option_groups || []).map((g: any) => ({
      key: g.key || g.Key,
      values: g.values || g.Values,
    }))
    ruleForm.variants = (product.variants || []).map((v: any) => ({
      sku: v.sku,
      options: v.options,
      price: v.price,
      stock: v.stock,
      // Map the single image string from server to the images array for UI preview
      images: v.image
        ? [
            {
              name: v.sku,
              url: v.image,
            },
          ]
        : [],
    }))
    // Main product images
    ruleForm.images = (product.images || []).map((img) => ({
      name: img.id,
      url: img.url,
    }))
  }
}

const handleDeleteBtnClick = (id: string, name: string) => {
  ElMessageBox.confirm(`Delete product name:  "${name}"?`, 'Confirm delete', {
    type: 'warning',
    confirmButtonText: 'Delete',
    cancelButtonText: 'Cancel',
  }).then(async () => {
    handleDeleteProduct(id, name)
  })
}
const handleDeleteProduct = (id: string, name: string) => {}
//#endregion

const addOptionGroup = () => {
  ruleForm.option_groups.push({ key: '', values: [] })
}

const removeOptionGroup = (index: number) => {
  ruleForm.option_groups.splice(index, 1)
}

watch(
  [() => ruleForm.option_groups, () => ruleForm.name],
  () => {
    generateVariants()
  },
  { deep: true },
)

const generateVariants = () => {
  const groups = ruleForm.option_groups.filter((g) => g.key && g.values.length > 0)
  if (groups.length === 0) {
    ruleForm.variants = []
    return
  }

  const combinations = (idx: number, currentOptions: Record<string, string>): any[] => {
    if (idx === groups.length) {
      return [currentOptions]
    }

    const res: any[] = []
    const group = groups[idx]
    if (!group) return []
    for (const val of group.values) {
      res.push(...combinations(idx + 1, { ...currentOptions, [group.key]: val }))
    }
    return res
  }

  const allCombinations = combinations(0, {})

  // Preserve existing variant data (price, stock, images) if possible
  const oldVariants = [...ruleForm.variants]
  ruleForm.variants = allCombinations.map((options) => {
    const existing = oldVariants.find((v) => {
      return Object.entries(options).every(([key, val]) => v.options[key] === val)
    })

    const optionsPart = Object.values(options)
      .map((ov: any) => ov.toString().replace(/\s+/g, '-'))
      .join('-')

    // Generate SKU logic: Name-Val1-Val2
    const sku = (ruleForm.name.replace(/\s+/g, '-') || '') + '-' + optionsPart

    if (existing) {
      // If the existing SKU is empty or starts with a dash (meaning name was empty), update it
      if (!existing.sku || existing.sku.startsWith('-')) {
        existing.sku = sku
      }
      return existing
    }

    return {
      sku: sku,
      options: options,
      price: 0,
      stock: 0,
      images: [],
    }
  })
}

const download = (images: ProductImage[], index: number) => {
  if (!images[index]) return
  const url = images[index].url
  const suffix = url.slice(url.lastIndexOf('.'))
  const filename = Date.now() + suffix

  fetch(url)
    .then((response) => response.blob())
    .then((blob) => {
      const blobUrl = URL.createObjectURL(new Blob([blob]))
      const link = document.createElement('a')
      link.href = blobUrl
      link.download = filename
      document.body.appendChild(link)
      link.click()
      URL.revokeObjectURL(blobUrl)
      link.remove()
    })
}
</script>

<template>
  <div class="seller-product-view">
    <div class="toolbar">
      <h2>Product Management</h2>
      <div style="display: flex; gap: 12px">
        <el-input
          v-model="searchQuery"
          placeholder="Search products..."
          :prefix-icon="Search"
          style="width: 240px"
          clearable
          @change="fetchProducts"
          @clear="fetchProducts"
        />
        <el-button size="large" :icon="Plus" color="var(--main-color)" @click="openModal('add')">
          Add Product
        </el-button>
      </div>
    </div>

    <div class="table-container">
      <el-table
        v-loading="isLoading"
        :data="productData"
        border
        style="width: 100%"
        height="100%"
        :row-class-name="getRowClassName"
        @sort-change="handleSortChange"
        @filter-change="handleFilterChange"
      >
        <el-table-column width="1">
          <template #default="props">
            <span :class="`quantity-${props.row.variants?.length}`"></span>
          </template>
        </el-table-column>
        <el-table-column type="expand" class-name="expand-column">
          <template #default="props">
            <el-table
              :show-header="false"
              :data="props.row.variants"
              v-if="props.row.variants.length > 1"
            >
              <el-table-column width="48" />

              <el-table-column align="right" width="128">
                <template #default="scope">
                  <el-image
                    style="width: 52px; height: 52px; position: relative; top: 3px"
                    :src="scope.row.image ?? 'https://placehold.co/52x52'"
                    fit="contain"
                  />
                </template>
              </el-table-column>

              <el-table-column>
                <template #default="scope">
                  <p>
                    {{
                      Object.entries(scope.row.options)
                        .map(([key, value]) => {
                          return `${key}: ${value}`
                        })
                        .join(', ')
                    }}
                  </p>
                </template>
              </el-table-column>

              <el-table-column />

              <el-table-column label="Price" width="180">
                <template #default="scope"> {{ scope.row.price }}$ </template>
              </el-table-column>

              <el-table-column prop="stock" label="Stock" width="80" />

              <el-table-column width="120" />

              <el-table-column width="140" />
            </el-table>
          </template>
        </el-table-column>
        <el-table-column label="Product image" width="128" align="center">
          <template #default="scope">
            <el-image
              style="width: 56px; height: 56px; position: relative; top: 3px"
              :src="scope.row.images[0]?.url ?? 'https://placehold.co/56x56'"
              fit="contain"
              show-progress
              preview-teleported
              :preview-src-list="scope.row.images.map((image: ProductImage) => image.url)"
            >
              <template #toolbar="{ actions, prev, next, reset, activeIndex, setActiveItem }">
                <el-icon @click="prev"><Back /></el-icon>
                <el-icon @click="next"><Right /></el-icon>
                <el-icon @click="setActiveItem(scope.row.images.length - 1)">
                  <DArrowRight />
                </el-icon>
                <el-icon @click="actions('zoomOut')"><ZoomOut /></el-icon>
                <el-icon @click="actions('zoomIn', { enableTransition: false, zoomRate: 2 })">
                  <ZoomIn />
                </el-icon>
                <el-icon @click="actions('clockwise', { rotateDeg: 180, enableTransition: false })">
                  <RefreshRight />
                </el-icon>
                <el-icon @click="actions('anticlockwise')"><RefreshLeft /></el-icon>
                <el-icon @click="reset"><Refresh /></el-icon>
                <el-icon @click="download(scope.row.images, activeIndex)"><Download /></el-icon>
              </template>
            </el-image>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="Product name" sortable="custom" />
        <el-table-column prop="sold_count" label="Sold count" width="120" sortable="custom" />
        <el-table-column prop="price" label="Price" width="180" sortable="custom">
          <template #default="scope">
            <p v-if="scope.row.variants?.length > 1">
              {{
                Math.min(...scope.row.variants.map((x: ProductVariant) => x.price)) +
                '$ ~ ' +
                Math.max(...scope.row.variants.map((x: ProductVariant) => x.price))
              }}$
            </p>
            <p v-else>{{ scope.row.variants[0].price }}$</p>
          </template>
        </el-table-column>
        <el-table-column prop="quantity" label="Stock" width="100" sortable="custom">
          <template #default="scope">
            <p v-if="scope.row.variants?.length > 1">
              {{
                scope.row.variants.reduce(
                  (sum: number, item: ProductVariant) => sum + item.stock!,
                  0,
                )
              }}
            </p>
            <p v-else>{{ scope.row.variants[0].stock }}</p>
          </template>
        </el-table-column>
        <el-table-column
          prop="status"
          label="Status"
          column-key="status"
          align="center"
          width="120"
          :filters="[
            { text: 'Available', value: 'AVAILABLE' },
            { text: 'Out of stock', value: 'OUT_OF_STOCK' },
            { text: 'Hidden', value: 'HIDDEN' },
          ]"
        >
          <template #default="scope">
            <el-tag :type="getStatusTagType(scope.row.status)" disable-transitions>{{
              scope.row.status
            }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="center" label="Operations" width="140">
          <template #default="{ row }">
            <el-button type="primary" :icon="Edit" @click="openModal('edit', row)" />
            <el-button
              type="danger"
              :icon="Delete"
              @click="handleDeleteBtnClick(row.id, row.name)"
            />
          </template>
        </el-table-column>
      </el-table>
    </div>

    <div class="pagination-container">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        :total="total"
        size="large"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>

    <!-- Add/Edit dialog -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogContent.title"
      width="800"
      align-center
      style="height: 90%; overflow-y: auto"
    >
      <el-form label-width="auto" ref="ruleFormRef" :model="ruleForm" :rules="rules">
        <el-form-item prop="images" label="Image">
          <el-upload
            class="images-upload"
            v-model:file-list="ruleForm.images"
            accept="image/png, image/jpeg, image/jpg"
            :auto-upload="false"
            list-type="picture-card"
            :limit="9"
            multiple
            :on-exceed="handleExceed"
            :on-preview="handlePictureCardPreview"
            :on-remove="handleRemove"
            :on-change="handleChange"
          >
            <el-icon><Plus /></el-icon>
          </el-upload>
        </el-form-item>
        <el-form-item prop="name" label="Name">
          <el-input
            size="large"
            v-model="ruleForm.name"
            placeholder="Please input product name"
            clearable
          />
        </el-form-item>
        <el-form-item prop="description" label="Description">
          <el-input
            size="large"
            v-model="ruleForm.description"
            :rows="6"
            type="textarea"
            placeholder="Please input product description"
          />
        </el-form-item>
        <el-form-item label="General Categories">
          <el-select
            v-model="ruleForm.category_ids"
            multiple
            size="large"
            filterable
            placeholder="Select general categories"
            style="width: 100%"
          >
            <el-option
              v-for="item in generalCategoryOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="Seller Categories">
          <el-select
            v-model="ruleForm.seller_category_ids"
            multiple
            size="large"
            filterable
            placeholder="Select seller categories"
            style="width: 100%"
          >
            <el-option
              v-for="item in sellerCategoryOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>

        <div class="option-groups-section">
          <div
            v-for="(group, index) in ruleForm.option_groups"
            :key="index"
            class="option-group-card"
          >
            <div class="group-header">
              <span>Option Group {{ index + 1 }}</span>
              <el-button type="danger" link :icon="Delete" @click="removeOptionGroup(index)"
                >Remove</el-button
              >
            </div>
            <el-form-item label="Group Key">
              <el-input v-model="group.key" placeholder="e.g., Size, Color" />
            </el-form-item>
            <el-form-item label="Values">
              <el-select-v2
                v-model="group.values"
                :options="[]"
                placeholder="Press Enter to add values"
                allow-create
                default-first-option
                filterable
                multiple
                clearable
                popper-class="group_option_popper"
                :reserve-keyword="false"
              />
            </el-form-item>
          </div>
          <el-button
            type="primary"
            plain
            :icon="Plus"
            @click="addOptionGroup"
            v-if="ruleForm.option_groups.length < 2"
          >
            Add Option Group
          </el-button>
        </div>

        <el-divider v-if="ruleForm.variants.length > 0" content-position="left"
          >Product Variants</el-divider
        >

        <el-table
          v-if="ruleForm.variants.length > 0"
          :data="ruleForm.variants"
          border
          style="width: 100%; margin-top: 20px"
        >
          <el-table-column label="Image" width="100" align="center">
            <template #default="{ row }">
              <el-upload
                class="option-image-upload"
                v-model:file-list="row.images"
                accept="image/png, image/jpeg, image/jpg"
                :auto-upload="false"
                list-type="picture-card"
                :limit="1"
                :on-preview="handlePictureCardPreview"
                :on-change="(file: any) => handleVariantImageChange(file, row)"
              >
                <el-icon><Plus /></el-icon>
              </el-upload>
            </template>
          </el-table-column>

          <el-table-column label="Variant" min-width="150">
            <template #default="{ row }">
              <div class="variant-info">
                <span v-for="(val, key) in row.options" :key="key"> {{ key }}: {{ val }} </span>
              </div>
            </template>
          </el-table-column>

          <el-table-column label="Price" width="150">
            <template #default="{ row }">
              <el-input-number
                v-model="row.price"
                :min="0"
                controls-position="right"
                style="width: 100%"
              />
            </template>
          </el-table-column>

          <el-table-column label="Stock" width="120">
            <template #default="{ row }">
              <el-input-number
                v-model="row.stock"
                :min="0"
                controls-position="right"
                style="width: 100%"
              />
            </template>
          </el-table-column>

          <el-table-column label="SKU" width="180">
            <template #default="{ row }">
              <el-input v-model="row.sku" placeholder="SKU" />
            </template>
          </el-table-column>
        </el-table>
      </el-form>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dialogVisible = false">Cancel</el-button>
          <el-button type="primary" @click="submitForm(ruleFormRef)">
            {{ dialogContent.mainBtnText }}
          </el-button>
        </div>
      </template>
    </el-dialog>

    <!-- Image preview dialog -->
    <el-dialog v-model="previewImageDialogVisible">
      <img w-full :src="previewImageUrl" alt="Preview Image" />
    </el-dialog>
  </div>
</template>

<style scoped>
.seller-product-view {
  padding: 24px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
  display: flex;
  flex-direction: column;
  height: 100%;
  box-sizing: border-box;
  overflow: hidden;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  flex-shrink: 0;
}

.toolbar h2 {
  margin: 0;
  font-size: 20px;
  color: #1e293b;
}

.table-container {
  flex: 1;
  min-height: 0;
  overflow: hidden;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
  flex-shrink: 0;
}

.option-groups-section {
  padding: 16px;
  background: #f8fafc;
  border-radius: 8px;
  margin-bottom: 24px;
}

.option-group-card {
  background: white;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  padding: 16px;
  margin-bottom: 16px;
}

.group-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  font-weight: 600;
  color: #475569;
}

.variant-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.variant-info span {
  font-size: 13px;
  color: #64748b;
}

.images-upload :deep(.el-upload),
.images-upload :deep(.el-upload-list__item) {
  width: 120px !important;
  height: 120px !important;
}

.option-image-upload :deep(.el-upload),
.option-image-upload :deep(.el-upload-list__item) {
  width: 80px !important;
  height: 80px !important;
}

/* Hide group option popper */
.group_option_popper {
  display: none;
}

/* Disable adding multiple pictures to a product option */
.option-image-upload :deep(.el-upload-list__item + .el-upload) {
  display: none;
}

/* Make option image upload stays still */
.option-image-upload :deep(.el-upload-list__item) {
  animation: none !important;
  transition: none !important;
}

:deep(.hide-expand .el-table__expand-column .el-table__expand-icon) {
  display: none;
}
</style>

<style>
.seller-product-view .el-table__expand-icon svg {
  scale: 1.5;
}
</style>
