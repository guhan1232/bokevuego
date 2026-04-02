<template>
  <div class="page-container">
    <div class="page-header">
      <h2>访问记录</h2>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon" style="background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%)">
          <el-icon :size="24"><View /></el-icon>
        </div>
        <div class="stat-info">
          <span class="stat-value">{{ stats.today || 0 }}</span>
          <span class="stat-label">今日访问</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: linear-gradient(135deg, #10b981 0%, #34d399 100%)">
          <el-icon :size="24"><Calendar /></el-icon>
        </div>
        <div class="stat-info">
          <span class="stat-value">{{ stats.yesterday || 0 }}</span>
          <span class="stat-label">昨日访问</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: linear-gradient(135deg, #f59e0b 0%, #fbbf24 100%)">
          <el-icon :size="24"><Timer /></el-icon>
        </div>
        <div class="stat-info">
          <span class="stat-value">{{ stats.week || 0 }}</span>
          <span class="stat-label">近7天访问</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: linear-gradient(135deg, #ec4899 0%, #f472b6 100%)">
          <el-icon :size="24"><DataLine /></el-icon>
        </div>
        <div class="stat-info">
          <span class="stat-value">{{ stats.total || 0 }}</span>
          <span class="stat-label">总访问量</span>
        </div>
      </div>
    </div>

    <!-- 热门页面 -->
    <el-card shadow="hover" style="margin-bottom: 20px" v-if="stats.top_pages && stats.top_pages.length > 0">
      <template #header>
        <span style="font-weight: 600">热门页面（近7天）</span>
      </template>
      <div class="top-pages">
        <div class="page-item" v-for="(page, index) in stats.top_pages" :key="index">
          <span class="page-rank">{{ index + 1 }}</span>
          <span class="page-path">{{ page.path }}</span>
          <span class="page-count">{{ page.count }} 次</span>
        </div>
      </div>
    </el-card>

    <!-- 访问列表 -->
    <el-card shadow="hover">
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center">
          <span style="font-weight: 600">访问日志</span>
          <el-input
            v-model="searchPath"
            placeholder="搜索路径"
            style="width: 200px"
            clearable
            @clear="fetchVisits"
            @keyup.enter="fetchVisits"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </div>
      </template>

      <el-table :data="visits" stripe style="width: 100%">
        <el-table-column prop="path" label="访问路径" min-width="200">
          <template #default="{ row }">
            <el-link type="primary" :href="row.path" target="_blank">{{ row.path }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="ip" label="IP 地址 / 位置" min-width="220">
          <template #default="{ row }">
            <div>
              <div style="font-weight: 500">{{ row.ip }}</div>
              <div v-if="row.region" style="font-size: 12px; color: #909399; margin-top: 2px">
                {{ row.region }}
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="user_agent" label="User Agent" min-width="200" show-overflow-tooltip />
        <el-table-column prop="created_at" label="访问时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrap">
        <el-pagination
          v-model:current-page="page"
          :page-size="20"
          :total="total"
          layout="total, prev, pager, next"
          @current-change="fetchVisits"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import request from '../api/request'
import { View, Calendar, Timer, DataLine, Search } from '@element-plus/icons-vue'

const visits = ref([])
const stats = ref({})
const page = ref(1)
const total = ref(0)
const searchPath = ref('')

const fetchVisits = async () => {
  try {
    const res = await request.get('/visits', {
      params: {
        page: page.value,
        path: searchPath.value
      }
    })
    visits.value = res.items || []
    total.value = res.total || 0
  } catch (e) {
    console.error('Failed to fetch visits:', e)
  }
}

const fetchStats = async () => {
  try {
    const res = await request.get('/visits/stats')
    stats.value = res
  } catch (e) {
    console.error('Failed to fetch stats:', e)
  }
}

const formatDate = (dateStr) => {
  if (!dateStr) return '--'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

onMounted(() => {
  fetchVisits()
  fetchStats()
})
</script>

<style scoped>
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  margin-bottom: 24px;
}

.stat-card {
  background: #fff;
  border-radius: 16px;
  padding: 24px;
  display: flex;
  align-items: center;
  gap: 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  border: 1px solid #e2e8f0;
}

.stat-icon {
  width: 52px;
  height: 52px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  flex-shrink: 0;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #1e293b;
  display: block;
}

.stat-label {
  font-size: 13px;
  color: #64748b;
  margin-top: 4px;
  display: block;
}

.top-pages {
  display: flex;
  flex-direction: column;
}

.page-item {
  display: flex;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid #f1f5f9;
}

.page-item:last-child {
  border-bottom: none;
}

.page-rank {
  width: 24px;
  height: 24px;
  border-radius: 6px;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
  margin-right: 12px;
  flex-shrink: 0;
}

.page-path {
  flex: 1;
  font-size: 14px;
  color: #1e293b;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.page-count {
  font-size: 14px;
  font-weight: 600;
  color: #6366f1;
  margin-left: 16px;
}

.pagination-wrap {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

@media (max-width: 1200px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
}
</style>
