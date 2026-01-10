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
} from '@element-plus/icons-vue'
import axios from 'axios'
import {
  ElMessage,
  ElMessageBox,
  type CascaderNode,
  type FormInstance,
  type FormRules,
  type UploadFile,
  type UploadFiles,
  type UploadProps,
  type UploadUserFile,
} from 'element-plus'
import { onMounted, reactive, ref, watch } from 'vue'
import { formatNumberWithDots } from '@/utils/formatNumberWithDots'

interface ProductImage {
  id: string
  url: string
  order: number
}

interface OptionGroup {
  Key: string
  Values: string[]
}

interface ProductVariant {
  id: string
  sku: string
  options: Record<string, string>
  price: number
  stock: number
  image: string
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

const userId = ref('7d27d61d-8cb3-4ef9-a9ff-0a92f217855b')
// const token = localStorage.getItem('access_token')
const token =
  'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjY4Mjc1MDYsImlhdCI6MTc2Njc0MTEwNiwicm9sZSI6InNlbGxlciIsInVzZXJJZCI6IjdkMjdkNjFkLThjYjMtNGVmOS1hOWZmLTBhOTJmMjE3ODU1YiIsInVzZXJuYW1lIjoidGVzdDEifQ.cCs5gwsDkDVZZrnHHmeTFGmuleu9RR2Ieke8pYaRH_s'

const productData = ref<Product[]>([])
const isLoading = ref(false)
const currentPage = ref(1)
const pageSize = ref(20)
const totalPages = ref(10)

onMounted(() => {
  fetchProducts()
})

const fetchProducts = () => {
  isLoading.value = true
  axios
    .get(
      import.meta.env.VITE_GET_SELLER_PRODUCT_API_URL +
        '/' +
        'seller-1' +
        `?sort_by=name&sort_direction=asc&page=${currentPage.value}&limit=${pageSize.value}`,
    )
    .then((response) => {
      productData.value = response.data.products
      totalPages.value = response.data.pagination.total_pages
    })
    .catch((error) => {
      console.error(error)
    })
    .finally(() => {
      isLoading.value = false
    })
}

const categoryOptions = [
  {
    value: 'guide',
    label: 'Guide',
    children: [
      {
        value: 'disciplines',
        label: 'Disciplines',
        children: [
          {
            value: 'consistency',
            label: 'Consistency',
          },
          {
            value: 'feedback',
            label: 'Feedback',
          },
          {
            value: 'efficiency',
            label: 'Efficiency',
          },
          {
            value: 'controllability',
            label: 'Controllability',
          },
        ],
      },
      {
        value: 'navigation',
        label: 'Navigation',
        children: [
          {
            value: 'side nav',
            label: 'Side Navigation',
          },
          {
            value: 'top nav',
            label: 'Top Navigation',
          },
        ],
      },
    ],
  },
  {
    value: 'component',
    label: 'Component',
    children: [
      {
        value: 'basic',
        label: 'Basic',
        children: [
          {
            value: 'layout',
            label: 'Layout',
          },
          {
            value: 'color',
            label: 'Color',
          },
          {
            value: 'typography',
            label: 'Typography',
          },
          {
            value: 'icon',
            label: 'Icon',
          },
          {
            value: 'button',
            label: 'Button',
          },
        ],
      },
      {
        value: 'form',
        label: 'Form',
        children: [
          {
            value: 'radio',
            label: 'Radio',
          },
          {
            value: 'checkbox',
            label: 'Checkbox',
          },
          {
            value: 'input',
            label: 'Input',
          },
          {
            value: 'input-number',
            label: 'InputNumber',
          },
          {
            value: 'select',
            label: 'Select',
          },
          {
            value: 'cascader',
            label: 'Cascader',
          },
          {
            value: 'switch',
            label: 'Switch',
          },
          {
            value: 'slider',
            label: 'Slider',
          },
          {
            value: 'time-picker',
            label: 'TimePicker',
          },
          {
            value: 'date-picker',
            label: 'DatePicker',
          },
          {
            value: 'datetime-picker',
            label: 'DateTimePicker',
          },
          {
            value: 'upload',
            label: 'Upload',
          },
          {
            value: 'rate',
            label: 'Rate',
          },
          {
            value: 'form',
            label: 'Form',
          },
        ],
      },
      {
        value: 'data',
        label: 'Data',
        children: [
          {
            value: 'table',
            label: 'Table',
          },
          {
            value: 'tag',
            label: 'Tag',
          },
          {
            value: 'progress',
            label: 'Progress',
          },
          {
            value: 'tree',
            label: 'Tree',
          },
          {
            value: 'pagination',
            label: 'Pagination',
          },
          {
            value: 'badge',
            label: 'Badge',
          },
        ],
      },
      {
        value: 'notice',
        label: 'Notice',
        children: [
          {
            value: 'alert',
            label: 'Alert',
          },
          {
            value: 'loading',
            label: 'Loading',
          },
          {
            value: 'message',
            label: 'Message',
          },
          {
            value: 'message-box',
            label: 'MessageBox',
          },
          {
            value: 'notification',
            label: 'Notification',
          },
        ],
      },
      {
        value: 'navigation',
        label: 'Navigation',
        children: [
          {
            value: 'menu',
            label: 'Menu',
          },
          {
            value: 'tabs',
            label: 'Tabs',
          },
          {
            value: 'breadcrumb',
            label: 'Breadcrumb',
          },
          {
            value: 'dropdown',
            label: 'Dropdown',
          },
          {
            value: 'steps',
            label: 'Steps',
          },
        ],
      },
      {
        value: 'others',
        label: 'Others',
        children: [
          {
            value: 'dialog',
            label: 'Dialog',
          },
          {
            value: 'tooltip',
            label: 'Tooltip',
          },
          {
            value: 'popover',
            label: 'Popover',
          },
          {
            value: 'card',
            label: 'Card',
          },
          {
            value: 'carousel',
            label: 'Carousel',
          },
          {
            value: 'collapse',
            label: 'Collapse',
          },
        ],
      },
    ],
  },
  {
    value: 'resource',
    label: 'Resource',
    children: [
      {
        value: 'axure',
        label: 'Axure Components',
      },
      {
        value: 'sketch',
        label: 'Sketch Templates',
      },
      {
        value: 'docs',
        label: 'Design Documentation',
      },
    ],
  },
]

const optionListData = ref<ProductOption[]>([])

const dialogVisible = ref(false)
const dialogMode = ref<DialogMode>('add')
const dialogContent = ref({
  title: 'Add product',
  mainBtnText: 'Add',
})

const previewImageDialogVisible = ref(false)
const previewImageUrl = ref('')
const hasGroup1 = ref(false)
const hasGroup2 = ref(false)

//#region FORM
type ProductOption = {
  image: File | undefined
  option1: string
  option2: string
  price: number | undefined
  stock: number | undefined
}

interface RuleForm {
  images: UploadUserFile[]
  name: string
  category: string // or might be category_id
  description: string
  visibility: boolean
  price: number | undefined
  stock: number | undefined
  group1: string
  group2: string
  group1_options: string[]
  group2_options: string[]
  product_options: ProductOption[]
}

const ruleFormRef = ref<FormInstance>()
const ruleForm = reactive<RuleForm>({
  images: [],
  name: '',
  category: '',
  description: '',
  visibility: true,
  price: undefined,
  stock: undefined,
  group1: '',
  group2: '',
  group1_options: [],
  group2_options: [],
  product_options: [],
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

const handlePictureCardPreview: UploadProps['onPreview'] = (uploadFile) => {
  previewImageUrl.value = uploadFile.url!
  previewImageDialogVisible.value = true
}
//#endregion

//#region CATEGORY CASCADER
const categoryFilterMethod = (node: CascaderNode, keyword: string) => {
  return node.label.toLowerCase().includes(keyword.toLowerCase())
}
//#endregion

const handleAddProduct = () => {}
const handleEditProduct = () => {}

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

//#region MODAL
const openModal = (mode: DialogMode) => {
  dialogMode.value = mode
  dialogVisible.value = true
  if (mode === 'add') {
    dialogContent.value = {
      title: 'Add product',
      mainBtnText: 'Add',
    }
    clearRuleForm()
  } else if (mode === 'edit') {
    dialogContent.value = {
      title: 'Edit product',
      mainBtnText: 'Save',
    }
    // TODO: set ruleForm values here!!!
    // ruleForm.category = category ?? ''
    // ruleForm.id = id ?? ''
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

watch(
  [() => ruleForm.group1_options, () => ruleForm.group2_options],
  ([new_group1_options, new_group2_options]) => {
    let newOptionListData: ProductOption[] = []
    if (new_group2_options.length) {
      new_group1_options.forEach((group1_option) => {
        new_group2_options.forEach((group2_option) => {
          newOptionListData.push({
            image: undefined,
            option1: group1_option,
            option2: group2_option,
            price: undefined,
            stock: undefined,
          })
        })
      })

      optionListData.value.forEach((option) => {
        if (
          option.image !== undefined ||
          option.price !== undefined ||
          option.stock !== undefined
        ) {
          newOptionListData[
            newOptionListData.findIndex(
              (newOption) =>
                option.option1 === newOption.option1 && option.option2 === newOption.option2,
            )
          ] = option
        }

        optionListData.value = newOptionListData
      })
    } else {
      new_group1_options.forEach((group1_option) => {
        newOptionListData.push({
          image: undefined,
          option1: group1_option,
          option2: '',
          price: undefined,
          stock: undefined,
        })
      })

      optionListData.value.forEach((option) => {
        if (
          option.image !== undefined ||
          option.price !== undefined ||
          option.stock !== undefined
        ) {
          newOptionListData[
            newOptionListData.findIndex((newOption) => option.option1 === newOption.option1)
          ] = option
        }

        optionListData.value = newOptionListData
      })

      optionListData.value = newOptionListData
    }
  },
)

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
  <el-button
    type="primary"
    size="large"
    style="margin-bottom: 12px"
    color="#1aae51"
    @click="openModal('add')"
    >Add Product</el-button
  >
  <el-table :data="productData" style="width: 100%">
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
            <template #default="scope"> {{ formatNumberWithDots(scope.row.price) }}đ </template>
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
    <el-table-column prop="name" label="Product name" />
    <el-table-column prop="sold_count" label="Sold count" width="100" />
    <el-table-column prop="price" label="Price" width="180">
      <template #default="scope">
        <p v-if="scope.row.variants?.length > 1">
          {{
            formatNumberWithDots(
              Math.min(...scope.row.variants.map((x: ProductVariant) => x.price)),
            ) +
            'đ ~ ' +
            formatNumberWithDots(
              Math.max(...scope.row.variants.map((x: ProductVariant) => x.price)),
            )
          }}đ
        </p>
        <p v-else>{{ formatNumberWithDots(scope.row.variants[0].price) }}đ</p>
      </template>
    </el-table-column>
    <el-table-column prop="quantity" label="Stock" width="80">
      <template #default="scope">
        <p v-if="scope.row.variants?.length > 1">
          {{
            scope.row.variants.reduce((sum: number, item: ProductVariant) => sum + item.stock!, 0)
          }}
        </p>
        <p v-else>{{ scope.row.variants[0].stock }}</p>
      </template>
    </el-table-column>
    <el-table-column
      prop="status"
      label="Status"
      align="center"
      width="120"
      :filters="[
        { text: 'Available', value: 'AVAILABLE' },
        { text: 'Out of stock', value: 'OUT_OF_STOCK' },
        { text: 'Hidden', value: 'HIDDEN' },
      ]"
      :filter-method="filterTag"
      filter-placement="bottom-end"
    >
      <template #default="scope">
        <el-tag :type="getStatusTagType(scope.row.status)" disable-transitions>{{
          scope.row.status
        }}</el-tag>
      </template>
    </el-table-column>
    <el-table-column align="center" label="Operations" width="140">
      <template #default="{ row }">
        <el-button type="primary" :icon="Edit" @click="openModal('edit')" />
        <el-button type="danger" :icon="Delete" @click="handleDeleteBtnClick(row.id, row.name)" />
      </template>
    </el-table-column>
  </el-table>

  <el-pagination
    background
    layout="prev, pager, next"
    :total="totalPages"
    v-model:page-size="pageSize"
    v-model:current-page="currentPage"
    :page-sizes="[10, 20, 30, 40]"
    @size-change="fetchProducts"
    @current-change="fetchProducts"
    size="large"
    style="
      display: flex;
      justify-content: center;
      margin-top: 20px;
      --el-pagination-button-bg-color: #fff;
    "
  />

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
      <el-form-item prop="category" label="Category">
        <el-cascader
          size="large"
          placeholder="Please select a category"
          :options="categoryOptions"
          filterable
          :filter-method="categoryFilterMethod"
          clearable
          style="width: 100%"
        />
      </el-form-item>
      <el-button
        type="primary"
        size="large"
        style="margin-left: 82px; margin-bottom: 16px; margin-top: 12px"
        @click="
          () => {
            if (hasGroup1) {
              if (hasGroup2) {
                hasGroup2 = false
                // Pass data from group2 to group1
              } else {
                hasGroup1 = false
              }
            } else {
              hasGroup1 = true
            }
          }
        "
      >
        {{ hasGroup1 ? 'Remove product group 1' : 'Add product group' }}
      </el-button>
      <el-form-item prop="price" label="Price" v-show="!hasGroup1">
        <el-input-number size="large" v-model="ruleForm.price" :min="1000" :max="1000000000000">
          <template #prefix>
            <span>đ</span>
          </template>
        </el-input-number>
      </el-form-item>
      <el-form-item prop="stock" label="Stock" v-show="!hasGroup1">
        <el-input-number size="large" v-model="ruleForm.stock" :min="0" :max="1000000000000">
        </el-input-number>
      </el-form-item>

      <el-form-item prop="group1" label="Group 1" v-show="hasGroup1">
        <el-input
          size="large"
          v-model="ruleForm.group1"
          placeholder="Please input group name"
          clearable
        />
      </el-form-item>
      <el-form-item prop="group1" label="Options" v-show="hasGroup1">
        <el-select-v2
          v-model="ruleForm.group1_options"
          :options="[]"
          placeholder="Please input an option then press enter"
          style="margin-right: 16px; vertical-align: middle"
          allow-create
          default-first-option
          filterable
          multiple
          clearable
          popper-class="group_option_popper"
          :reserve-keyword="false"
        />
      </el-form-item>

      <el-button
        v-show="hasGroup1"
        type="primary"
        size="large"
        style="margin-left: 82px; margin-bottom: 16px; margin-top: 12px"
        @click="hasGroup2 = !hasGroup2"
      >
        {{ hasGroup2 ? 'Remove product group 2' : 'Add product group 2' }}
      </el-button>
      <el-form-item prop="group2" label="Group 2" v-show="hasGroup2">
        <el-input
          size="large"
          v-model="ruleForm.group2"
          placeholder="Please input group name"
          clearable
        />
      </el-form-item>
      <el-form-item
        prop="group2_options"
        label="Options"
        v-show="hasGroup2"
        style="margin-bottom: 32px"
      >
        <el-select-v2
          v-model="ruleForm.group2_options"
          :options="[]"
          placeholder="Please input an option then press enter"
          style="margin-right: 16px; vertical-align: middle"
          allow-create
          default-first-option
          filterable
          multiple
          clearable
          popper-class="group_option_popper"
          :reserve-keyword="false"
        />
      </el-form-item>
      <!-- HAS ONE GROUP -->
      <el-form-item prop="" label="Option list" v-show="hasGroup1 && !hasGroup2">
        <el-table :data="optionListData" border style="width: 100%">
          <el-table-column
            prop="option1"
            :label="ruleForm.group1 === '' ? 'Group 1' : ruleForm.group1"
            width="100"
          />
          <el-table-column prop="image" label="Image" align="center" style="width: 80px">
            <template #default="{ row }">
              <el-upload
                class="option-image-upload"
                v-model:file-list="row.image"
                accept="image/png, image/jpeg, image/jpg"
                :auto-upload="false"
                list-type="picture-card"
                :limit="1"
                :on-preview="handlePictureCardPreview"
                :on-remove="handleRemove"
                :on-change="handleChange"
              >
                <el-icon><Plus /></el-icon>
              </el-upload>
            </template>
          </el-table-column>
          <el-table-column prop="price" label="Price">
            <template #default="{ row }">
              <el-input v-model="row.price" placeholder="Input price" />
            </template>
          </el-table-column>
          <el-table-column prop="stock" label="Stock">
            <template #default="{ row }">
              <el-input v-model="row.stock" placeholder="Input stock" />
            </template>
          </el-table-column>
        </el-table>
      </el-form-item>

      <!-- HAS BOTH GROUPS -->
      <el-form-item prop="" label="Option list" v-show="hasGroup1 && hasGroup2">
        <el-table :data="optionListData" border style="width: 100%">
          <el-table-column
            prop="option1"
            :label="ruleForm.group1 === '' ? 'Group 1' : ruleForm.group1"
            width="100"
          />
          <el-table-column
            prop="option2"
            :label="ruleForm.group2 === '' ? 'Group 2' : ruleForm.group2"
            width="100"
          />
          <el-table-column prop="image" label="Image" align="center" style="width: 80px">
            <template #default="{ row }">
              <el-upload
                class="option-image-upload"
                v-model:file-list="row.image"
                accept="image/png, image/jpeg, image/jpg"
                :auto-upload="false"
                list-type="picture-card"
                :limit="1"
                :on-preview="handlePictureCardPreview"
                :on-remove="handleRemove"
                :on-change="handleChange"
              >
                <el-icon><Plus /></el-icon>
              </el-upload>
            </template>
          </el-table-column>
          <el-table-column prop="price" label="Price">
            <template #default="{ row }">
              <el-input v-model="row.price" placeholder="Input price" />
            </template>
          </el-table-column>
          <el-table-column prop="stock" label="Stock">
            <template #default="{ row }">
              <el-input v-model="row.stock" placeholder="Input stock" />
            </template>
          </el-table-column>
        </el-table>
      </el-form-item>
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
</template>

<style>
.images-upload {
  .el-upload,
  .el-upload-list__item {
    width: 120px !important;
    height: 120px !important;
  }
}
.el-upload,
.el-upload-list__item {
  width: 120px !important;
  height: 120px !important;
}

.option-image-upload {
  .el-upload,
  .el-upload-list__item {
    width: 80px !important;
    height: 80px !important;
  }
}

/* Hide group option popper */
.group_option_popper {
  display: none;
}

/* Disable adding multiple pictures to a product option */
.option-image-upload .el-upload-list__item + .el-upload {
  display: none;
}

/* Make option image upload stays still */
.option-image-upload .el-upload-list__item {
  animation: none !important;
  transition: none !important;
}

.expand-column:has(.quantity-1) {
  display: none;
}

.el-table__cell:has(.quantity-1) ~ .expand-column .el-table__expand-icon {
  visibility: hidden;
}

.el-table__expand-icon svg {
  scale: 1.5;
}
</style>
