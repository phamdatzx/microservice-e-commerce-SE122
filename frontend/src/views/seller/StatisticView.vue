<script setup lang="ts">
import { ref, computed } from 'vue'
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

const timeRange = ref('today')
const orderType = ref('confirmed')

const salesData = ref({
  value: '12,500,000',
  comparison: 12.5,
  timeRange: 'vs 00:00-15:00 yesterday',
})

const ordersData = ref({
  value: 45,
  comparison: -5.2,
  timeRange: 'vs 00:00-15:00 yesterday',
})

const chartData = computed(() => ({
  labels: ['00:00', '03:00', '06:00', '09:00', '12:00', '15:00', '18:00', '21:00'],
  datasets: [
    {
      label: 'Sales',
      borderColor: '#409EFF',
      backgroundColor: '#409EFF',
      data: [1500000, 2200000, 1800000, 3500000, 4200000, 3800000, 6500000, 5000000],
      tension: 0.4,
    },
    {
      label: 'Orders',
      borderColor: '#E6A23C',
      backgroundColor: '#E6A23C',
      data: [5, 8, 6, 12, 15, 13, 22, 18],
      tension: 0.4,
      yAxisID: 'y1',
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
  scales: {
    y: {
      beginAtZero: true,
      grid: {},
    },
    x: {
      grid: {
        display: false,
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
          <el-option label="Today until 15:00 today (GMT+7)" value="today" />
          <el-option label="Yesterday" value="yesterday" />
          <el-option label="Last 7 days" value="week" />
          <el-option label="Last 30 days" value="month" />
        </el-select>
      </div>

      <div class="filter-item">
        <span class="label">Order Type</span>
        <el-tooltip content="Filter by order status" placement="top">
          <el-icon class="info-icon"><InfoFilled /></el-icon>
        </el-tooltip>
        <el-select v-model="orderType" style="width: 200px">
          <el-option label="Confirmed orders" value="confirmed" />
          <el-option label="All orders" value="all" />
        </el-select>
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
            <div class="comparison">
              <span class="text">{{ salesData.timeRange }}</span>
              <span class="percent">{{ salesData.comparison.toFixed(2) }}%</span>
            </div>
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
            <div class="comparison">
              <span class="text">{{ ordersData.timeRange }}</span>
              <span class="percent">{{ ordersData.comparison.toFixed(2) }}%</span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Chart Section -->
    <div class="chart-header">
      <h3 class="section-title">Chart</h3>
      <span class="selected-count">Selected 2/4</span>
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
  margin-bottom: 16px;
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
