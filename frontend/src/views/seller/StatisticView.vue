<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import axios from 'axios'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  type ChartOptions,
} from 'chart.js'
import { Line } from 'vue-chartjs'
import { InfoFilled, ArrowRight } from '@element-plus/icons-vue'

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend)

const timeRange = ref('week')
const viewType = ref('day') // Default to 'day' view
const customDateRange = ref<[Date, Date]>()
const isLoading = ref(false)

const salesData = ref({
  value: '0',
  comparison: 0,
  timeRange: 'vs previous period',
})

const ordersData = ref({
  value: 0,
  comparison: 0,
  timeRange: 'vs previous period',
})

const chartLabels = ref<string[]>([])
const chartRevenueData = ref<number[]>([])
const chartOrderData = ref<number[]>([])

const getDateRange = (range: string) => {
  const now = new Date()
  const to = now.toISOString()
  let from = new Date()

  switch (range) {
    case 'today':
      from.setHours(0, 0, 0, 0)
      break
    case 'yesterday':
      from.setDate(from.getDate() - 1)
      from.setHours(0, 0, 0, 0)
      const endYesterday = new Date(from)
      endYesterday.setHours(23, 59, 59, 999)
      return { from: from.toISOString(), to: endYesterday.toISOString() }
    case 'week':
      from.setDate(from.getDate() - 7)
      break
    case 'month':
      from.setDate(from.getDate() - 30)
      break
    case 'last_6_months':
      from.setMonth(from.getMonth() - 6)
      break
    case 'last_12_months':
      from.setFullYear(from.getFullYear() - 1)
      break
    case 'custom':
      if (customDateRange.value && customDateRange.value.length === 2) {
        return {
          from: customDateRange.value[0].toISOString(),
          to: customDateRange.value[1].toISOString(),
        }
      }
      return { from: to, to } // Fallback if no date selected
  }
  return { from: from.toISOString(), to }
}

const generateTimeLabels = (start: Date, end: Date, type: 'day' | 'month') => {
  const labels: string[] = []
  const current = new Date(start)
  const endDate = new Date(end)

  if (type === 'day') {
    // Reset time components to compare dates correctly
    current.setHours(0, 0, 0, 0)
    endDate.setHours(23, 59, 59, 999)

    while (current <= endDate) {
      const year = current.getFullYear()
      const month = String(current.getMonth() + 1).padStart(2, '0')
      const day = String(current.getDate()).padStart(2, '0')
      labels.push(`${year}-${month}-${day}`)
      current.setDate(current.getDate() + 1)
    }
  } else {
    // Month type
    current.setDate(1) // Start at the first of the month
    // We only compare year and month for the end condition
    const endYear = endDate.getFullYear()
    const endMonth = endDate.getMonth()

    while (
      current.getFullYear() < endYear ||
      (current.getFullYear() === endYear && current.getMonth() <= endMonth)
    ) {
      const year = current.getFullYear()
      const month = String(current.getMonth() + 1).padStart(2, '0')
      labels.push(`${year}-${month}`)
      current.setMonth(current.getMonth() + 1)
    }
  }
  return labels
}

const fetchStatistics = async () => {
  const token = localStorage.getItem('access_token')
  if (!token) return

  // Prevent fetch if custom range is selected but no dates are picked
  if (
    timeRange.value === 'custom' &&
    (!customDateRange.value || customDateRange.value.length !== 2)
  ) {
    return
  }

  isLoading.value = true
  try {
    const { from, to } = getDateRange(timeRange.value)

    // Use selected view type
    const type = viewType.value

    const response = await axios.get(`${import.meta.env.VITE_BE_API_URL}/order/seller/statistics`, {
      params: { from, to, type },
      headers: { Authorization: `Bearer ${token}` },
    })

    if (response.data) {
      const data = response.data

      // Update Summary Cards
      salesData.value.value = new Intl.NumberFormat('vi-VN').format(data.total_revenue || 0)
      ordersData.value.value = data.order_count || 0

      // Update Chart Data
      const allLabels = generateTimeLabels(new Date(from), new Date(to), type as 'day' | 'month')
      const revenueMap = new Map()
      const orderMap = new Map()

      if (data.breakdown && Array.isArray(data.breakdown)) {
        data.breakdown.forEach((item: any) => {
          revenueMap.set(item.period, item.revenue)
          orderMap.set(item.period, item.order_count)
        })
      }

      chartLabels.value = allLabels
      chartRevenueData.value = allLabels.map((label) => revenueMap.get(label) || 0)
      chartOrderData.value = allLabels.map((label) => orderMap.get(label) || 0)
    }
  } catch (error) {
    console.error('Error fetching statistics:', error)
  } finally {
    isLoading.value = false
  }
}

watch(timeRange, (newVal) => {
  if (newVal !== 'custom') {
    fetchStatistics()
  }
})

watch(customDateRange, () => {
  if (timeRange.value === 'custom') {
    fetchStatistics()
  }
})

watch(viewType, (newVal) => {
  // Reset time range when switching view types
  if (newVal === 'month') {
    timeRange.value = 'last_6_months'
  } else {
    timeRange.value = 'week'
  }
  // fetchStatistics will be triggered by timeRange watch or we can call it here if timeRange doesn't change
  // Actually, setting timeRange triggers the watcher above.
  // But if timeRange was already 'today' and we switch to 'day' (no change), it won't trigger.
  // However, usually we switch types.
  // Exception: if I am on 'custom' in 'month' view and switch to 'day' view, I might want to keep 'custom' or reset.
  // The logic "timeRange.value = ..." triggers the watcher.
  // Let's ensure fetch happens if the range didn't change (unlikely with this logic but possible if defaults match).
  // Wait, if I switch Day->Month, I set 'this_year'. 'today' != 'this_year', so watcher triggers.
  // If I switch Month->Day, I set 'today'. 'this_year' != 'today', watcher triggers.
})

onMounted(() => {
  fetchStatistics()
})

const chartData = computed(() => ({
  labels: chartLabels.value,
  datasets: [
    {
      label: 'Sales',
      borderColor: '#409EFF',
      backgroundColor: '#409EFF',
      data: chartRevenueData.value,
      tension: 0.4,
      yAxisID: 'y1', // Sales on Right
    },
    {
      label: 'Orders',
      borderColor: '#E6A23C',
      backgroundColor: '#E6A23C',
      data: chartOrderData.value,
      tension: 0.4,
      yAxisID: 'y', // Orders on Left
    },
  ],
}))

const chartOptions: ChartOptions<'line'> = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      position: 'bottom',
    },
  },
  elements: {
    point: {
      radius: 4,
      hitRadius: 30,
      hoverRadius: 6,
    },
  },
  scales: {
    y: {
      type: 'linear',
      display: true,
      position: 'left',
      beginAtZero: true,
      title: {
        display: true,
        text: 'Orders',
      },
      grid: {},
    },
    y1: {
      type: 'linear',
      display: true,
      position: 'right',
      beginAtZero: true,
      title: {
        display: true,
        text: 'Revenue',
      },
      grid: {
        drawOnChartArea: false, // only want the grid lines for one axis to show up
      },
    },
    x: {
      grid: {
        display: false,
      },
      ticks: {
        autoSkip: true,
        maxTicksLimit: 10,
        maxRotation: 0,
      },
    },
  },
}
</script>

<template>
  <div class="statistic-view">
    <!-- Filter Bar -->
    <div class="filter-bar">
      <div class="filter-item">
        <span class="label">Time Range</span>
        <el-select v-model="timeRange" style="width: 280px">
          <template v-if="viewType === 'day'">
            <el-option label="Last 7 days" value="week" />
            <el-option label="Last 30 days" value="month" />
          </template>
          <template v-else>
            <el-option label="Last 6 Months" value="last_6_months" />
            <el-option label="Last 12 Months" value="last_12_months" />
          </template>
          <el-option label="Custom Period" value="custom" />
        </el-select>

        <el-date-picker
          v-if="timeRange === 'custom'"
          v-model="customDateRange"
          type="daterange"
          range-separator="To"
          start-placeholder="Start date"
          end-placeholder="End date"
          style="width: 300px"
          :clearable="false"
        />

        <el-radio-group v-model="viewType">
          <el-radio-button label="Day" value="day" />
          <el-radio-button label="Month" value="month" />
        </el-radio-group>
      </div>
    </div>

    <!-- Overview Section -->
    <h3 class="section-title">Overview</h3>
    <el-row :gutter="20" class="overview-cards">
      <el-col :span="12">
        <el-card shadow="hover" class="stat-card">
          <template #header>
            <div class="card-header">
              <span>Sales</span>
              <el-tooltip content="Total sales revenue" placement="top">
                <el-icon><InfoFilled /></el-icon>
              </el-tooltip>
            </div>
          </template>
          <div class="card-content">
            <div class="main-value">₫ {{ salesData.value }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="hover" class="stat-card">
          <template #header>
            <div class="card-header">
              <span>Orders</span>
              <el-tooltip content="Total number of orders" placement="top">
                <el-icon><InfoFilled /></el-icon>
              </el-tooltip>
            </div>
          </template>
          <div class="card-content">
            <div class="main-value">{{ ordersData.value }}</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Chart Section -->
    <div class="chart-header">
      <h3 class="section-title">Chart</h3>
    </div>

    <el-card shadow="never" class="chart-card">
      <div class="chart-container">
        <Line :data="chartData" :options="chartOptions" />
      </div>
    </el-card>
  </div>
</template>

<style scoped>
.statistic-view {
  padding: 20px;
  background-color: #f6f6f6; /* Matching the gray background of the main area if needed, though usually handled by parent */
  min-height: 100%;
}

.filter-bar {
  display: flex;
  gap: 20px;
  background: white;
  padding: 16px;
  border-radius: 4px;
  margin-bottom: 24px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
  align-items: center;
}

.filter-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.label {
  color: #606266;
  font-size: 14px;
  font-weight: 500;
}

.info-icon {
  color: #909399;
  margin-right: 8px;
  cursor: help;
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 16px;
  margin-top: 0;
}

.overview-cards {
  margin-bottom: 24px;
}

.stat-card {
  height: 100%;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #606266;
  font-size: 14px;
}

.card-content {
  margin-top: 8px;
}

.main-value {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 8px;
}

.comparison {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #909399;
}

.percent {
  color: #303133;
  font-weight: 500;
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0px;
}

.selected-count {
  font-size: 12px;
  color: #909399;
}

.chart-card {
  background: white;
}

.chart-container {
  height: 400px;
  width: 100%;
}
</style>
