<script setup lang="ts">
import { Delete, Edit } from '@element-plus/icons-vue'

const tableData = [
  {
    image_url: 'https://placehold.co/40x40',
    name: 'Electronics',
    sold_count: 1,
    price: 140000,
    quantity: 21,
    status: 'AVAILABLE',
  },
  {
    image_url: 'https://placehold.co/80x80',
    name: 'Clothing',
    sold_count: 2,
    price: 130000,
    quantity: 23,
    status: 'OUT_OF_STOCK',
  },
  {
    image_url: 'https://placehold.co/120x40',
    name: 'Home & Kitchen',
    sold_count: 3,
    price: 120000,
    quantity: 22,
    status: 'AVAILABLE',
  },
  {
    image_url: 'https://placehold.co/40x120',
    name: 'Books',
    sold_count: 4,
    price: 110000,
    quantity: 24,
    status: 'HIDDEN',
  },
]

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
</script>

<template>
  <el-button type="primary" size="large" style="margin-bottom: 12px" color="#1aae51"
    >Add Product</el-button
  >
  <el-table :data="tableData" border style="width: 100%">
    <el-table-column prop="image_url" label="Product image" width="82px">
      <template #default="scope">
        <el-image style="width: 56px; height: 56px" :src="scope.row.image_url" fit="contain" />
      </template>
    </el-table-column>
    <el-table-column prop="name" label="Product name" />
    <el-table-column prop="sold_count" label="Sold count" />
    <el-table-column prop="price" label="Price" />
    <el-table-column prop="quantity" label="Stock" />
    <el-table-column
      prop="status"
      label="Status"
      align="center"
      width="110"
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
        <el-button type="danger" :icon="Delete" @click="openDeleteModal(row.category)" />
      </template>
    </el-table-column>
  </el-table>
</template>
