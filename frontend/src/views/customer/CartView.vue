<script setup lang="ts">
import Header from '@/components/Header.vue'

const tableData = [
  {
    date: '2016-05-02',
    name: 'John Smith',
    address: 'No. 1518,  Jinshajiang Road, Putuo District, Shanghai',
  },
  {
    date: '2016-05-04',
    name: 'Jim Green',
    address: 'No. 1518,  Jinshajiang Road, Putuo District, Shanghai',
  },
  {
    date: '2016-05-01',
    name: 'Joe Black',
    address: 'No. 1518,  Jinshajiang Road, Putuo District, Shanghai',
  },
  {
    date: '2016-05-03',
    name: 'Jon Snow',
    address: 'No. 1518,  Jinshajiang Road, Putuo District, Shanghai',
  },
]
</script>

<template>
  <Header />

  <div class="main-container">
    <div class="box-shadow border-radius" style="padding: 20px; margin: 20px 0">
      <el-table :data="tableData" style="width: 100%">
        <el-table-column label="Product">
          <template #default="scope">
            <!--  -->
          </template>
        </el-table-column>
        <el-table-column prop="price" label="Price" />
        <el-table-column label="Quantity">
          <template #default="scope">
            <!--  -->
          </template>
        </el-table-column>
        <el-table-column label="Total">
          <template #default="scope">
            <!--  -->
          </template>
        </el-table-column>
        <el-table-column label="Operations">
          <template #default="scope">
            <!--  -->
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>
