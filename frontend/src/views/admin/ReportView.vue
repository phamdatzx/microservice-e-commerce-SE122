<script setup lang="ts">
import {
  ArrowLeft,
  CircleCheck,
  CircleClose,
  Delete,
  Lock,
  Unlock,
  Warning,
} from '@element-plus/icons-vue'
import axios from 'axios'
import { ElLoading, ElMessage, ElMessageBox } from 'element-plus'
import { computed, onMounted, ref, watch } from 'vue'

interface User {
  id: string
  name: string
  email: string
  avatar: string
}

interface Product {
  id: string
  name: string
  images: { id: string; url: string }[]
  seller_id: string
  is_disabled: boolean
}

interface Report {
  id: string
  product_id: string
  variant_id: string
  user: User
  reason: string
  description: string
  status: 'pending' | 'reviewed' | 'resolved' | 'rejected'
  created_at: string
}

const token = localStorage.getItem('access_token') || ''
const viewMode = ref<'list' | 'detail'>('list')
const reports = ref<Report[]>([])
const isLoading = ref(false)
const statusFilter = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const totalReports = ref(0)

// Detail View State
const selectedReport = ref<Report | null>(null)
const productInfo = ref<Product | null>(null)
const relatedReports = ref<Report[]>([])
const isDetailLoading = ref(false)
const relatedCurrentPage = ref(1)
const relatedPageSize = ref(10)
const relatedTotalReports = ref(0)

onMounted(() => {
  fetchReports()
})

const fetchReports = () => {
  isLoading.value = true
  axios
    .get(import.meta.env.VITE_BE_API_URL + '/product/report', {
      params: {
        page: currentPage.value,
        limit: pageSize.value,
        status: statusFilter.value || undefined,
      },
      headers: { Authorization: `Bearer ${token}` },
    })
    .then((res) => {
      reports.value = res.data.reports || []
      totalReports.value = res.data.total || 0
    })
    .catch((err) => {
      console.error('Failed to fetch reports:', err)
      ElMessage.error('Failed to load reports')
    })
    .finally(() => {
      isLoading.value = false
    })
}

watch([currentPage, pageSize, statusFilter], () => {
  fetchReports()
})

const handleViewDetail = async (report: Report) => {
  selectedReport.value = report
  viewMode.value = 'detail'
  isDetailLoading.value = true

  try {
    // Fetch product info
    const productRes = await axios.get(
      import.meta.env.VITE_BE_API_URL + `/product/public/${report.product_id}`,
    )
    productInfo.value = productRes.data

    // Fetch related reports
    await fetchRelatedReports()
  } catch (err) {
    console.error('Failed to fetch detail info:', err)
    ElMessage.error('Failed to load product or related reports')
  } finally {
    isDetailLoading.value = false
  }
}

const fetchRelatedReports = async () => {
  if (!selectedReport.value) return
  try {
    const relatedRes = await axios.get(
      import.meta.env.VITE_BE_API_URL +
        `/product/report/product/${selectedReport.value.product_id}`,
      {
        params: {
          page: relatedCurrentPage.value,
          limit: relatedPageSize.value,
        },
        headers: { Authorization: `Bearer ${token}` },
      },
    )
    relatedReports.value = (relatedRes.data.reports || []).filter(
      (r: Report) => r.id !== selectedReport.value?.id,
    )
    relatedTotalReports.value = relatedRes.data.total || 0
  } catch (err) {
    console.error('Failed to fetch related reports:', err)
    ElMessage.error('Failed to load related reports')
  }
}

watch([relatedCurrentPage, relatedPageSize], () => {
  if (viewMode.value === 'detail') {
    fetchRelatedReports()
  }
})

const handleStatusUpdate = (reportId: string, status: string) => {
  const loading = ElLoading.service({ text: 'Updating status...' })
  axios
    .put(
      import.meta.env.VITE_BE_API_URL + `/product/report/${reportId}/status`,
      { status },
      { headers: { Authorization: `Bearer ${token}` } },
    )
    .then(() => {
      ElMessage.success(`Report ${status} successfully`)
      if (selectedReport.value && selectedReport.value.id === reportId) {
        selectedReport.value.status = status as any
      }
      fetchReports()
    })
    .catch((err) => {
      console.error('Failed to update status:', err)
      ElMessage.error('Failed to update status')
    })
    .finally(() => {
      loading.close()
    })
}

const handleToggleProductDisable = () => {
  if (!productInfo.value) return
  const isDisabled = !productInfo.value.is_disabled
  const action = isDisabled ? 'disable' : 'enable'

  ElMessageBox.prompt(`Please provide a reason to ${action} this product:`, `Confirm ${action}`, {
    confirmButtonText: action.charAt(0).toUpperCase() + action.slice(1),
    cancelButtonText: 'Cancel',
    inputPattern: /.+/,
    inputErrorMessage: 'Reason is required',
  }).then(({ value }) => {
    const loading = ElLoading.service({ text: `${action}ing product...` })
    axios
      .put(
        import.meta.env.VITE_BE_API_URL + `/product/${productInfo.value!.id}/disable`,
        { is_disabled: isDisabled, reason: value },
        { headers: { Authorization: `Bearer ${token}` } },
      )
      .then(() => {
        ElMessage.success(`Product ${action}d successfully`)
        productInfo.value!.is_disabled = isDisabled
      })
      .catch((err) => {
        console.error(`Failed to ${action} product:`, err)
        ElMessage.error(`Failed to ${action} product`)
      })
      .finally(() => {
        loading.close()
      })
  })
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return 'N/A'
  return new Date(dateStr).toLocaleString()
}

const getReasonLabel = (reason: string) => {
  const labels: Record<string, string> = {
    fake_product: 'Fake Product',
    quality_issue: 'Quality Issue',
    misleading_description: 'Misleading Description',
    counterfeit: 'Counterfeit',
    damaged_product: 'Damaged Product',
    other: 'Other',
  }
  return labels[reason] || reason
}

const getStatusType = (status: string) => {
  const types: Record<string, string> = {
    pending: 'warning',
    reviewed: 'info',
    resolved: 'success',
    rejected: 'danger',
  }
  return types[status] || 'info'
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
  currentPage.value = 1
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
}

const goBack = () => {
  viewMode.value = 'list'
  selectedReport.value = null
  productInfo.value = null
  relatedReports.value = []
}
</script>

<template>
  <div class="report-management-view">
    <!-- List View -->
    <div v-if="viewMode === 'list'">
      <div class="toolbar">
        <h2>Report Management</h2>
        <div class="actions">
          <el-select
            v-model="statusFilter"
            placeholder="Filter by Status"
            clearable
            style="width: 160px"
          >
            <el-option label="Pending" value="pending" />
            <el-option label="Resolved" value="resolved" />
            <el-option label="Rejected" value="rejected" />
          </el-select>
        </div>
      </div>

      <el-table v-loading="isLoading" :data="reports" border style="width: 100%">
        <el-table-column prop="id" label="Report ID" min-width="120" sortable />

        <el-table-column prop="product_id" label="Product ID" min-width="120" />

        <el-table-column label="Reporter" min-width="180">
          <template #default="{ row }">
            <div class="user-cell">
              <el-avatar :size="32" :src="row.user.avatar" />
              <span>{{ row.user.name }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="Reason" width="180">
          <template #default="{ row }">
            <el-tag type="danger" effect="plain">{{ getReasonLabel(row.reason) }}</el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="status" label="Status" width="120" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ row.status.toUpperCase() }}</el-tag>
          </template>
        </el-table-column>

        <el-table-column label="Actions" width="150" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewDetail(row)"> Details </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="totalReports"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>

    <!-- Detail View -->
    <div v-else-if="viewMode === 'detail'" v-loading="isDetailLoading">
      <div class="detail-header">
        <el-button :icon="ArrowLeft" @click="goBack">Back to List</el-button>
        <h2>Report Detail</h2>
      </div>

      <div v-if="selectedReport" class="detail-container">
        <!-- Product Section -->
        <div class="detail-section product-section">
          <h3>Product Information</h3>
          <div v-if="productInfo" class="product-card">
            <el-image
              class="product-image"
              :src="productInfo.images?.[0]?.url"
              fit="cover"
              :preview-src-list="[productInfo.images?.[0]?.url]"
              preview-teleported
            />
            <div class="product-meta">
              <div class="name-row">
                <h4>{{ productInfo.name }}</h4>
                <el-tag :type="productInfo.is_disabled ? 'danger' : 'success'">
                  {{ productInfo.is_disabled ? 'DISABLED' : 'ACTIVE' }}
                </el-tag>
              </div>
              <p>ID: {{ productInfo.id }}</p>
              <p>Seller ID: {{ productInfo.seller_id }}</p>

              <div class="product-actions">
                <el-button
                  :type="productInfo.is_disabled ? 'success' : 'danger'"
                  :icon="productInfo.is_disabled ? Unlock : Lock"
                  @click="handleToggleProductDisable"
                >
                  {{ productInfo.is_disabled ? 'Enable Product' : 'Disable Product' }}
                </el-button>
              </div>
            </div>
          </div>
          <el-skeleton v-else animated />
        </div>

        <div class="detail-grid">
          <!-- Current Report section -->
          <div class="detail-section">
            <h3>Current Report</h3>
            <div class="report-card main-report">
              <div class="report-user-info">
                <el-avatar :size="48" :src="selectedReport.user.avatar" />
                <div class="user-text">
                  <strong>{{ selectedReport.user.name }}</strong>
                  <span>{{ selectedReport.user.email }}</span>
                </div>
                <div class="report-date">{{ formatDate(selectedReport.created_at) }}</div>
              </div>

              <div class="report-body">
                <div class="report-tag">
                  <strong>Reason:</strong>
                  <el-tag type="danger">{{ getReasonLabel(selectedReport.reason) }}</el-tag>
                </div>
                <div class="report-description">
                  <strong>Description:</strong>
                  <p>{{ selectedReport.description || 'No description provided' }}</p>
                </div>
              </div>

              <div class="report-footer">
                <div class="status-actions">
                  <span style="margin-right: 12px">Current Status:</span>
                  <el-tag :type="getStatusType(selectedReport.status)">{{
                    selectedReport.status.toUpperCase()
                  }}</el-tag>

                  <div class="action-buttons">
                    <el-button
                      type="success"
                      :icon="CircleCheck"
                      @click="handleStatusUpdate(selectedReport!.id, 'resolved')"
                      :disabled="selectedReport.status === 'resolved'"
                    >
                      Resolve
                    </el-button>
                    <el-button
                      type="danger"
                      :icon="CircleClose"
                      @click="handleStatusUpdate(selectedReport!.id, 'rejected')"
                      :disabled="selectedReport.status === 'rejected'"
                    >
                      Reject
                    </el-button>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Other reports for this product -->
          <div class="detail-section">
            <h3>Other Reports for this Product</h3>
            <div class="related-reports">
              <div v-if="relatedReports.length === 0" class="no-data">
                No other reports for this product.
              </div>
              <div
                v-for="report in relatedReports"
                :key="report.id"
                class="report-card related-report"
              >
                <div class="report-user-info mini">
                  <el-avatar :size="24" :src="report.user.avatar" />
                  <strong>{{ report.user.name }}</strong>
                  <div class="report-date small">{{ formatDate(report.created_at) }}</div>
                </div>
                <div class="report-body mini">
                  <el-tag size="small" type="danger" effect="plain">{{
                    getReasonLabel(report.reason)
                  }}</el-tag>
                  <p class="desc-small">{{ report.description }}</p>
                </div>
                <div class="report-footer mini">
                  <el-tag size="small" :type="getStatusType(report.status)">{{
                    report.status
                  }}</el-tag>
                  <el-button type="primary" link size="small" @click="handleViewDetail(report)">
                    View
                  </el-button>
                </div>
              </div>

              <div class="related-pagination">
                <el-pagination
                  v-model:current-page="relatedCurrentPage"
                  v-model:page-size="relatedPageSize"
                  :page-sizes="[10, 20, 50]"
                  layout="prev, pager, next"
                  :total="relatedTotalReports"
                  small
                />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.report-management-view {
  padding: 24px;
  background: white;
  min-height: calc(100vh - 120px);
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.user-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.pagination-container {
  margin-top: 24px;
  display: flex;
  justify-content: flex-end;
}

/* Detail View Styles */
.detail-header {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 24px;
  border-bottom: 1px solid #eee;
  padding-bottom: 16px;
}

.detail-header h2 {
  margin: 0;
}

.detail-container {
  max-width: 1200px;
  margin: 0 auto;
}

.detail-section {
  margin-bottom: 32px;
}

.detail-section h3 {
  margin-bottom: 16px;
  color: #1a1a1a;
  font-size: 18px;
}

.product-card {
  display: flex;
  gap: 24px;
  padding: 20px;
  border: 1px solid #ebeef5;
  border-radius: 8px;
  background-color: #fcfcfc;
}

.product-image {
  width: 120px;
  height: 120px;
  border-radius: 8px;
  border: 1px solid #efefef;
  flex-shrink: 0;
}

.product-meta {
  flex: 1;
}

.name-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 8px;
}

.name-row h4 {
  margin: 0;
  font-size: 20px;
}

.product-meta p {
  margin: 4px 0;
  color: #666;
  font-size: 14px;
}

.product-actions {
  margin-top: 16px;
}

.detail-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 32px;
}

@media (max-width: 900px) {
  .detail-grid {
    grid-template-columns: 1fr;
  }
}

.report-card {
  border: 1px solid #ebeef5;
  border-radius: 8px;
  padding: 20px;
  background-color: white;
}

.report-card.main-report {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
}

.report-user-info {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
  position: relative;
}

.user-text {
  display: flex;
  flex-direction: column;
}

.user-text strong {
  font-size: 16px;
}

.user-text span {
  color: #888;
  font-size: 14px;
}

.report-date {
  margin-left: auto;
  color: #999;
  font-size: 13px;
}

.report-body {
  margin-bottom: 24px;
  padding: 16px;
  background-color: #f9f9f9;
  border-radius: 6px;
}

.report-tag {
  margin-bottom: 12px;
}

.report-tag strong,
.report-description strong {
  display: block;
  margin-bottom: 8px;
  color: #333;
}

.report-description p {
  margin: 0;
  line-height: 1.6;
}

.report-footer {
  border-top: 1px solid #eee;
  padding-top: 20px;
}

.status-actions {
  display: flex;
  align-items: center;
}

.action-buttons {
  margin-left: auto;
  display: flex;
  gap: 12px;
}

.related-reports {
  display: flex;
  flex-direction: column;
  gap: 16px;
  max-height: 500px;
  overflow-y: auto;
  padding-right: 8px;
}

.related-report {
  padding: 12px;
}

.report-user-info.mini {
  margin-bottom: 8px;
  gap: 8px;
}

.report-user-info.mini strong {
  font-size: 14px;
}

.report-date.small {
  font-size: 11px;
}

.report-body.mini {
  margin-bottom: 8px;
  padding: 8px;
}

.desc-small {
  font-size: 12px;
  margin: 4px 0 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.report-footer.mini {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 8px;
}

.related-pagination {
  margin-top: 16px;
  display: flex;
  justify-content: center;
}

.no-data {
  text-align: center;
  padding: 40px;
  color: #999;
  background: #fcfcfc;
  border: 1px dashed #dcdfe6;
  border-radius: 8px;
}
</style>
