<script setup lang="ts">
import {
  Search,
  CircleCheck,
  CircleClose,
  View,
  More,
  Lock,
  Delete,
  UserFilled,
  Warning,
} from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, onMounted, ref } from 'vue'

// #region Interfaces
interface ReportDetail {
  id: string
  reporterName: string
  reporterAvatar?: string
  violationType:
    | 'Fake Goods'
    | 'Prohibited Content'
    | 'Offensive Image'
    | 'Scam/Fake Price'
    | 'Other'
  orderId?: string
  description?: string
  createdAt: string
}

interface ReportedProduct {
  id: string
  productName: string
  productImage: string
  sellerName: string
  reportCount: number
  status: 'Pending' | 'Resolved' | 'Rejected'
  reports: ReportDetail[]
}
// #endregion

// #region State
const reportedProducts = ref<ReportedProduct[]>([])
const isLoading = ref(false)
const searchQuery = ref('')
const statusFilter = ref('')
const dialogVisible = ref(false)
const selectedProduct = ref<ReportedProduct | null>(null)
const currentPage = ref(1)
const pageSize = ref(10)
// #endregion

// #region Mock Data & Fetch
onMounted(() => {
  fetchReportedProducts()
})

const fetchReportedProducts = () => {
  isLoading.value = true
  // Simulate API delay
  setTimeout(() => {
    reportedProducts.value = Array.from({ length: 20 }).map((_, index) => {
      const id = `RP-${1000 + index}`
      const status = index % 3 === 0 ? 'Resolved' : index % 2 === 0 ? 'Rejected' : 'Pending'
      return {
        id,
        productName: `Product Sample ${index + 1} - High Quality Item`,
        productImage: `https://picsum.photos/seed/${index + 100}/100/100`,
        sellerName: `Seller ${index + 1}`,
        reportCount: Math.floor(Math.random() * 20) + 1,
        status: status as any,
        reports: Array.from({ length: Math.floor(Math.random() * 5) + 1 }).map((__, rIndex) => ({
          id: `R-${id}-${rIndex}`,
          reporterName: `User Reporter ${rIndex + 1}`,
          reporterAvatar: `https://i.pravatar.cc/150?u=${id}-${rIndex}`,
          violationType:
            rIndex % 2 === 0
              ? 'Fake Goods'
              : rIndex % 3 === 0
                ? 'Offensive Image'
                : 'Scam/Fake Price',
          orderId: rIndex % 2 === 0 ? `ORD-${Date.now()}-${rIndex}` : undefined,
          description: 'This item looks very suspicious and does not match description.',
          createdAt: new Date().toISOString().split('T')[0],
        })),
      }
    })
    isLoading.value = false
  }, 600)
}
// #endregion

// #region Computed
const filteredData = computed(() => {
  let data = reportedProducts.value

  // Search
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    data = data.filter(
      (item) =>
        item.productName.toLowerCase().includes(query) || item.id.toLowerCase().includes(query),
    )
  }

  // Filter
  if (statusFilter.value) {
    data = data.filter((item) => item.status === statusFilter.value)
  }

  // Sort by Report Count (Desc default)
  data.sort((a, b) => b.reportCount - a.reportCount)

  return data
})

const paginatedData = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredData.value.slice(start, end)
})
// #endregion

// #region Actions
const handleViewDetails = (row: ReportedProduct) => {
  selectedProduct.value = row
  dialogVisible.value = true
}

const handleAction = (command: string, row: ReportedProduct) => {
  if (command === 'remove') {
    ElMessageBox.confirm(`Remove product "${row.productName}"?`, 'Confirm Removal', {
      type: 'warning',
    }).then(() => {
      ElMessage.success('Product removed')
    })
  } else if (command === 'lock') {
    ElMessage.info('Product locked')
  } else if (command === 'ban') {
    ElMessageBox.confirm(`Ban seller "${row.sellerName}"?`, 'Confirm Ban', {
      type: 'error',
    }).then(() => {
      ElMessage.error('Seller banned')
    })
  } else if (command === 'warn') {
    ElMessage.warning('Warning sent to seller')
  } else if (command === 'reject') {
    row.status = 'Rejected'
    ElMessage.info('Reports rejected')
  } else if (command === 'resolve') {
    row.status = 'Resolved'
    ElMessage.success('Marked as resolved')
  }
}

const getStatusType = (status: string) => {
  if (status === 'Pending') return 'warning'
  if (status === 'Resolved') return 'success'
  if (status === 'Rejected') return 'info'
  return ''
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
  currentPage.value = 1
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
}
// #endregion
</script>

<template>
  <div class="report-management-view">
    <div class="toolbar">
      <h2>Report Management</h2>
      <div class="actions">
        <el-select
          v-model="statusFilter"
          placeholder="Filter by Status"
          clearable
          style="width: 160px"
        >
          <el-option label="Pending" value="Pending" />
          <el-option label="Resolved" value="Resolved" />
          <el-option label="Rejected" value="Rejected" />
        </el-select>
        <el-input
          v-model="searchQuery"
          placeholder="Search product or ID..."
          :prefix-icon="Search"
          style="width: 250px"
          clearable
        />
      </div>
    </div>

    <el-table v-loading="isLoading" :data="paginatedData" border style="width: 100%">
      <el-table-column prop="id" label="ID" width="100" sortable />

      <el-table-column label="Product" min-width="250">
        <template #default="{ row }">
          <div class="product-cell">
            <el-image
              style="width: 50px; height: 50px; border-radius: 4px"
              :src="row.productImage"
              fit="cover"
              preview-teleported
              :preview-src-list="[row.productImage]"
            />
            <div class="product-info">
              <span class="product-name">{{ row.productName }}</span>
              <span class="seller-name">Seller: {{ row.sellerName }}</span>
            </div>
          </div>
        </template>
      </el-table-column>

      <el-table-column prop="reportCount" label="Reports" width="120" sortable align="center">
        <template #default="{ row }">
          <el-tag type="danger" effect="plain" round>{{ row.reportCount }}</el-tag>
        </template>
      </el-table-column>

      <el-table-column prop="status" label="Status" width="120" align="center">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)">{{ row.status }}</el-tag>
        </template>
      </el-table-column>

      <el-table-column label="Actions" width="200" align="center">
        <template #default="{ row }">
          <div class="action-buttons">
            <el-button type="primary" link :icon="View" @click="handleViewDetails(row)">
              Details
            </el-button>
            <el-dropdown trigger="click" @command="(cmd: string) => handleAction(cmd, row)">
              <el-button type="primary" link style="margin-left: 10px">
                <el-icon size="large">
                  <More />
                </el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="resolve" :icon="CircleCheck"
                    >Mark Resolved</el-dropdown-item
                  >
                  <el-dropdown-item command="reject" :icon="CircleClose"
                    >Reject Reports</el-dropdown-item
                  >
                  <el-dropdown-item divided command="warn" :icon="Warning"
                    >Warn Seller</el-dropdown-item
                  >
                  <el-dropdown-item command="lock" :icon="Lock">Lock Product</el-dropdown-item>
                  <el-dropdown-item command="remove" :icon="Delete" style="color: red"
                    >Remove Product</el-dropdown-item
                  >
                  <el-dropdown-item command="ban" :icon="UserFilled" style="color: red"
                    >Ban Seller</el-dropdown-item
                  >
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination-container">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next, jumper"
        :total="filteredData.length"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>

    <!-- Details Dialog -->
    <el-dialog
      v-model="dialogVisible"
      title="Report Details"
      width="700px"
      align-center
      destroy-on-close
    >
      <div v-if="selectedProduct">
        <div class="report-header">
          <el-image class="report-product-img" :src="selectedProduct.productImage" fit="cover" />
          <div class="report-product-info">
            <h3>{{ selectedProduct.productName }}</h3>
            <p>ID: {{ selectedProduct.id }}</p>
            <el-tag :type="getStatusType(selectedProduct.status)">{{
              selectedProduct.status
            }}</el-tag>
          </div>
        </div>

        <el-divider content-position="left">User Reports</el-divider>

        <div class="reports-list">
          <div v-for="report in selectedProduct.reports" :key="report.id" class="report-item">
            <div class="report-user">
              <el-avatar :size="32" :src="report.reporterAvatar" />
              <div class="user-meta">
                <span class="username">{{ report.reporterName }}</span>
                <span class="date">{{ report.createdAt }}</span>
              </div>
            </div>
            <div class="report-content">
              <div class="tags">
                <el-tag size="small" effect="dark" type="danger">{{ report.violationType }}</el-tag>
                <el-tag v-if="report.orderId" size="small" type="info"
                  >Order: {{ report.orderId }}</el-tag
                >
              </div>
              <p class="description">{{ report.description }}</p>
            </div>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<style scoped>
.report-management-view {
  padding: 24px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.toolbar h2 {
  margin: 0;
  font-size: 20px;
  color: #1e293b;
}

.actions {
  display: flex;
  gap: 12px;
}

.product-cell {
  display: flex;
  gap: 12px;
  align-items: center;
}

.product-info {
  display: flex;
  flex-direction: column;
}

.product-name {
  font-weight: 500;
  line-height: 1.2;
  margin-bottom: 4px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.seller-name {
  font-size: 12px;
  color: #64748b;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.report-header {
  display: flex;
  gap: 16px;
  margin-bottom: 20px;
}

.report-product-img {
  width: 80px;
  height: 80px;
  border-radius: 6px;
}

.report-product-info h3 {
  margin: 0 0 4px 0;
  font-size: 16px;
}

.report-product-info p {
  color: #64748b;
  font-size: 13px;
  margin: 0 0 8px 0;
}

.reports-list {
  max-height: 400px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.report-item {
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  padding: 12px;
  background: #f8fafc;
}

.report-user {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 8px;
}

.user-meta {
  display: flex;
  flex-direction: column;
}

.username {
  font-weight: 600;
  font-size: 13px;
}

.date {
  font-size: 11px;
  color: #94a3b8;
}

.report-content .tags {
  display: flex;
  gap: 8px;
  margin-bottom: 8px;
}

.description {
  margin: 0;
  font-size: 13px;
  color: #334155;
  line-height: 1.5;
}
</style>
