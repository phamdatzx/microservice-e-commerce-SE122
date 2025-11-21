<script setup lang="ts">
import { Delete, Edit, Plus } from '@element-plus/icons-vue'
import {
  ElMessage,
  type CascaderNode,
  type FormInstance,
  type FormRules,
  type UploadFile,
  type UploadFiles,
  type UploadProps,
  type UploadUserFile,
} from 'element-plus'
import { reactive, ref, watch } from 'vue'

const tableData = [
  {
    image_url: 'https://down-vn.img.susercontent.com/file/sg-11134201-825af-mgbp824qs1e5b0_tn.webp',
    name: 'Electronics',
    sold_count: 1,
    quantity: 21,
    status: 'AVAILABLE',
    product_options: [
      {
        image_url:
          'https://down-vn.img.susercontent.com/file/sg-11134201-825af-mgbp824qs1e5b0_tn.webp',
        option1: 'Red',
        option2: 'XL',
        price: 200000,
        stock: 24,
      },
      {
        image_url:
          'https://down-vn.img.susercontent.com/file/sg-11134201-825af-mgbp824qs1e5b0_tn.webp',
        option1: 'Blue',
        option2: 'XXL',
        price: 220000,
        stock: 24,
      },
    ],
  },
  {
    image_url: 'https://down-vn.img.susercontent.com/file/sg-11134201-821fy-mghfw9c1o1sb29_tn.webp',
    name: 'Clothing',
    sold_count: 2,
    price: 130000,
    quantity: 23,
    status: 'OUT_OF_STOCK',
  },
  {
    image_url: 'https://down-vn.img.susercontent.com/file/sg-11134201-821cz-mghfwcncl0y6e6_tn.webp',
    name: 'Home & Kitchen',
    sold_count: 3,
    price: 120000,
    quantity: 22,
    status: 'AVAILABLE',
  },
  {
    image_url: 'https://down-vn.img.susercontent.com/file/sg-11134201-821f4-mghdelyvkzka36_tn.webp',
    name: 'Books',
    sold_count: 4,
    price: 110000,
    quantity: 24,
    status: 'HIDDEN',
  },
]

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

const addEditDialogVisible = ref(false)
const deleteDialogVisible = ref(false)
const previewImageDialogVisible = ref(false)
const previewImageUrl = ref('')
const title = ref('Add product')
const mainBtnText = ref('Add')
const currentAction = ref('add')
const hasGroup1 = ref(false)
const hasGroup2 = ref(false)
const deletingName = ref('')

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
      // handleAddEditCategory(currentAction.value)
    } else {
      console.log('error submit!', fields)
    }
  })
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

const getStatusTagContent = (status: string) => {
  if (status === 'AVAILABLE') return 'Available'
  else if (status === 'OUT_OF_STOCK') return 'Out of stock'
  else if (status === 'HIDDEN') return 'Hidden'
}

const getStatusTagType = (status: string) => {
  if (status === 'AVAILABLE') return 'success'
  else if (status === 'OUT_OF_STOCK') return 'danger'
  else if (status === 'HIDDEN') return 'info'
}

//#region MODAL

//#region CATEGORY CASCADER
const categoryFilterMethod = (node: CascaderNode, keyword: string) => {
  return node.label.toLowerCase().includes(keyword.toLowerCase())
}
//#endregion
const openAddModal = () => {
  addEditDialogVisible.value = true

  title.value = 'Add product'
  mainBtnText.value = 'Add'
  currentAction.value = 'add'
  // ruleForm.category = ''
}

const openEditModal = (category: string) => {
  addEditDialogVisible.value = true

  title.value = 'Edit category'
  mainBtnText.value = 'Save'
  currentAction.value = 'edit'
  // ruleForm.category = category
}

const openDeleteModal = (name: string) => {
  deleteDialogVisible.value = true
  deletingName.value = name
}
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
</script>

<template>
  <el-button
    type="primary"
    size="large"
    style="margin-bottom: 12px"
    color="#1aae51"
    @click="openAddModal()"
    >Add Product</el-button
  >
  <el-table :data="tableData" style="width: 100%">
    <el-table-column width="1">
      <template #default="props">
        <span :class="`quantity-${props.row.product_options?.length}`"></span>
      </template>
    </el-table-column>
    <el-table-column type="expand" class-name="expand-column">
      <template #default="props">
        <el-table
          :show-header="false"
          :data="props.row.product_options"
          v-if="props.row.product_options.length > 0"
        >
          <el-table-column width="48" />
          <el-table-column prop="image_url" align="right" width="128">
            <template #default="scope">
              <el-image
                style="width: 52px; height: 52px; position: relative; top: 3px"
                :src="scope.row.image_url"
                fit="contain"
              />
            </template>
          </el-table-column>
          <el-table-column>
            <template #default="scope">
              <p style="text-align: right; position: relative; right: -13px">
                {{ scope.row.option1 }},
              </p>
            </template>
          </el-table-column>
          <el-table-column>
            <template #default="scope">
              <p style="text-align: left; position: relative; left: -8px">
                {{ scope.row.option2 }}
              </p>
            </template>
          </el-table-column>
          <el-table-column prop="price" label="Price" />
          <el-table-column prop="stock" label="Stock" />
          <el-table-column width="120" />
          <el-table-column width="140" />
        </el-table>
      </template>
    </el-table-column>
    <el-table-column prop="image_url" label="Product image" width="128" align="center">
      <template #default="scope">
        <el-image
          style="width: 56px; height: 56px; position: relative; top: 3px"
          :src="scope.row.image_url"
          fit="contain"
        />
      </template>
    </el-table-column>
    <el-table-column prop="name" label="Product name" />
    <el-table-column prop="sold_count" label="Sold count" />
    <el-table-column prop="price" label="Price">
      <template #default="scope">
        <p v-if="scope.row.product_options?.length > 0">
          {{ Math.min(...scope.row.product_options.map((x: ProductOption) => x.price)) }} ~
          {{ Math.max(...scope.row.product_options.map((x: ProductOption) => x.price)) }}
        </p>
        <p v-else>{{ scope.row.price }}</p>
      </template>
    </el-table-column>
    <el-table-column prop="quantity" label="Stock">
      <template #default="scope">
        <p v-if="scope.row.product_options?.length > 0">
          {{
            scope.row.product_options.reduce(
              (sum: number, item: ProductOption) => sum + item.stock!,
              0,
            )
          }}
        </p>
        <p v-else>{{ scope.row.quantity }}</p>
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
          getStatusTagContent(scope.row.status)
        }}</el-tag>
      </template>
    </el-table-column>
    <el-table-column align="center" label="Operations" width="140">
      <template #default="{ row }">
        <el-button type="primary" :icon="Edit" @click="openEditModal(row.category)" />
        <el-button type="danger" :icon="Delete" @click="openDeleteModal(row.name)" />
      </template>
    </el-table-column>
  </el-table>

  <!-- Add/Edit dialog -->
  <el-dialog
    v-model="addEditDialogVisible"
    :title="title"
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
      <el-form-item prop="group2_options" label="Options" v-show="hasGroup2">
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
    </el-form>

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
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="addEditDialogVisible = false">Cancel</el-button>
        <el-button type="primary" @click="submitForm(ruleFormRef)">
          {{ mainBtnText }}
        </el-button>
      </div>
    </template>
  </el-dialog>

  <!-- Delete dialog -->
  <el-dialog v-model="deleteDialogVisible" title="Delete category" width="500" align-center>
    <span>Are you sure you want to delete product name: "{{ deletingName }}"?</span>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="deleteDialogVisible = false">Cancel</el-button>
        <el-button type="danger" @click="handleDeleteCategory()"> Delete </el-button>
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

.expand-column:has(.quantity-0) {
  display: none;
}

.el-table__cell:has(.quantity-undefined) ~ .expand-column .el-table__expand-icon {
  visibility: hidden;
}
</style>
