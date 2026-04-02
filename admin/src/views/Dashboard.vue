<template>
  <div class="dashboard">
    <!-- 统计卡片 -->
    <div class="stats-grid">
      <div class="stat-card" v-for="item in statCards" :key="item.label">
        <div class="stat-content">
          <div class="stat-info">
            <span class="stat-label">{{ item.label }}</span>
            <div class="stat-value-row">
              <span class="stat-value">{{ item.value }}</span>
              <span class="stat-trend" :class="item.trend > 0 ? 'up' : 'down'">
                <el-icon v-if="item.trend > 0"><Top /></el-icon>
                <el-icon v-else><Bottom /></el-icon>
                {{ Math.abs(item.trend) }}%
              </span>
            </div>
            <span class="stat-desc">较昨日</span>
          </div>
          <div class="stat-icon" :style="{ background: item.bg }">
            <el-icon :size="26"><component :is="item.icon" /></el-icon>
          </div>
        </div>
      </div>
    </div>

    <!-- 主内容区 -->
    <div class="dashboard-main">
      <!-- 左侧：图表和快捷操作 -->
      <div class="dashboard-left">
        <!-- Tab 切换 -->
        <div class="card tabs-card">
          <div class="card-header">
            <div class="tabs-list">
              <span 
                v-for="tab in tabs" 
                :key="tab.key" 
                :class="['tab-item', { active: activeTab === tab.key }]"
                @click="activeTab = tab.key"
              >
                {{ tab.label }}
              </span>
            </div>
          </div>
          <div class="card-body">
            <!-- 文章列表 -->
            <div class="data-list" v-if="recentArticles.length > 0">
              <div class="list-item" v-for="article in recentArticles" :key="article.id">
                <div class="item-icon blue">
                  <el-icon><Document /></el-icon>
                </div>
                <div class="item-content">
                  <div class="item-title">{{ article.title }}</div>
                  <div class="item-meta">
                    <span>{{ article.category || '未分类' }}</span>
                    <span>{{ formatDate(article.created_at) }}</span>
                  </div>
                </div>
                <div class="item-status" :class="article.status === 1 ? 'published' : 'draft'">
                  {{ article.status === 1 ? '已发布' : '草稿' }}
                </div>
              </div>
            </div>
            <el-empty v-else description="暂无文章数据" :image-size="80" />
          </div>
        </div>

        <!-- 快捷操作 -->
        <div class="card quick-card">
          <div class="card-header">
            <h3>快捷操作</h3>
          </div>
          <div class="card-body">
            <div class="quick-grid">
              <div class="quick-item" @click="$router.push('/articles/new')">
                <div class="quick-icon purple">
                  <el-icon :size="24"><EditPen /></el-icon>
                </div>
                <span>写文章</span>
              </div>
              <div class="quick-item" @click="$router.push('/articles')">
                <div class="quick-icon blue">
                  <el-icon :size="24"><Document /></el-icon>
                </div>
                <span>文章列表</span>
              </div>
              <div class="quick-item" @click="$router.push('/messages')">
                <div class="quick-icon green">
                  <el-icon :size="24"><ChatDotRound /></el-icon>
                </div>
                <span>查看留言</span>
              </div>
              <div class="quick-item" @click="$router.push('/settings')">
                <div class="quick-icon orange">
                  <el-icon :size="24"><Setting /></el-icon>
                </div>
                <span>站点设置</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧：留言和系统信息 -->
      <div class="dashboard-right">
        <!-- 最新留言 -->
        <div class="card messages-card">
          <div class="card-header">
            <h3>最新留言</h3>
            <el-button type="primary" link @click="$router.push('/messages')">查看全部</el-button>
          </div>
          <div class="card-body">
            <div class="message-list" v-if="recentMessages.length > 0">
              <div class="message-item" v-for="msg in recentMessages" :key="msg.id">
                <el-avatar :size="36" class="msg-avatar">{{ msg.name?.[0] || 'U' }}</el-avatar>
                <div class="msg-content">
                  <div class="msg-header">
                    <span class="msg-name">{{ msg.name }}</span>
                    <span class="msg-time">{{ formatDate(msg.created_at) }}</span>
                  </div>
                  <div class="msg-text">{{ msg.content }}</div>
                </div>
              </div>
            </div>
            <el-empty v-else description="暂无留言" :image-size="60" />
          </div>
        </div>

        <!-- 系统信息 -->
        <div class="card system-card">
          <div class="card-header">
            <h3>系统信息</h3>
          </div>
          <div class="card-body">
            <div class="system-list">
              <div class="system-item">
                <span class="system-label">系统版本</span>
                <span class="system-value">v1.0.0</span>
              </div>
              <div class="system-item">
                <span class="system-label">运行环境</span>
                <span class="system-value">Production</span>
              </div>
              <div class="system-item">
                <span class="system-label">数据库</span>
                <span class="system-value">SQLite</span>
              </div>
              <div class="system-item">
                <span class="system-label">后端框架</span>
                <span class="system-value">Go + Gin</span>
              </div>
              <div class="system-item">
                <span class="system-label">前端框架</span>
                <span class="system-value">Vue 3 + Element Plus</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { getStats, getArticles, getMessages } from '../api'
import { Document, ChatDotRound, View, EditPen, Top, Bottom } from '@element-plus/icons-vue'

const stats = ref({})
const recentArticles = ref([])
const recentMessages = ref([])
const activeTab = ref('all')

const tabs = [
  { key: 'all', label: '全部' },
  { key: 'published', label: '已发布' },
  { key: 'draft', label: '草稿' },
]

const statCards = computed(() => [
  { 
    label: '文章总数', 
    value: stats.value.total_articles || 0, 
    icon: Document, 
    bg: 'linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%)',
    trend: 12
  },
  { 
    label: '已发布', 
    value: stats.value.published_articles || 0, 
    icon: EditPen, 
    bg: 'linear-gradient(135deg, #10b981 0%, #34d399 100%)',
    trend: 8
  },
  { 
    label: '总浏览量', 
    value: stats.value.total_views || 0, 
    icon: View, 
    bg: 'linear-gradient(135deg, #f59e0b 0%, #fbbf24 100%)',
    trend: 24
  },
  { 
    label: '留言数', 
    value: stats.value.total_messages || 0, 
    icon: ChatDotRound, 
    bg: 'linear-gradient(135deg, #ec4899 0%, #f472b6 100%)',
    trend: -5
  },
])

const formatDate = (dateStr) => {
  if (!dateStr) return '--'
  const date = new Date(dateStr)
  return `${date.getMonth() + 1}/${date.getDate()}`
}

onMounted(async () => {
  try {
    stats.value = await getStats()
    
    // 获取最新文章
    const articlesData = await getArticles({ page: 1, page_size: 5 })
    recentArticles.value = articlesData.items || []
    
    // 获取最新留言
    const messagesData = await getMessages({ page: 1, page_size: 4 })
    recentMessages.value = messagesData.items || []
  } catch (e) {
    console.error('Failed to fetch dashboard data:', e)
  }
})
</script>

<style scoped>
.dashboard {
  padding: 0;
}

/* 统计卡片 */
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
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  border: 1px solid #e2e8f0;
  transition: all 0.3s;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.08);
}

.stat-content {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.stat-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.stat-label {
  font-size: 14px;
  color: #64748b;
  font-weight: 500;
}

.stat-value-row {
  display: flex;
  align-items: center;
  gap: 10px;
}

.stat-value {
  font-size: 32px;
  font-weight: 700;
  color: #1e293b;
  line-height: 1;
}

.stat-trend {
  display: flex;
  align-items: center;
  gap: 2px;
  font-size: 12px;
  font-weight: 600;
  padding: 2px 6px;
  border-radius: 4px;
}

.stat-trend.up {
  color: #10b981;
  background: rgba(16, 185, 129, 0.1);
}

.stat-trend.down {
  color: #ef4444;
  background: rgba(239, 68, 68, 0.1);
}

.stat-desc {
  font-size: 12px;
  color: #94a3b8;
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  flex-shrink: 0;
}

/* 主内容区布局 */
.dashboard-main {
  display: flex;
  gap: 20px;
}

.dashboard-left {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.dashboard-right {
  width: 360px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* 卡片样式 */
.card {
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  border: 1px solid #e2e8f0;
  overflow: hidden;
}

.card-header {
  padding: 20px 24px;
  border-bottom: 1px solid #f1f5f9;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.card-header h3 {
  font-size: 16px;
  font-weight: 600;
  color: #1e293b;
}

.card-body {
  padding: 20px 24px;
}

/* Tab 切换 */
.tabs-list {
  display: flex;
  gap: 8px;
}

.tab-item {
  padding: 8px 16px;
  font-size: 14px;
  font-weight: 500;
  color: #64748b;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.tab-item:hover {
  color: #6366f1;
  background: rgba(99, 102, 241, 0.05);
}

.tab-item.active {
  color: #6366f1;
  background: rgba(99, 102, 241, 0.1);
}

/* 数据列表 */
.data-list {
  display: flex;
  flex-direction: column;
}

.list-item {
  display: flex;
  align-items: center;
  padding: 14px 0;
  border-bottom: 1px solid #f1f5f9;
  gap: 14px;
  transition: background 0.2s;
}

.list-item:last-child {
  border-bottom: none;
}

.list-item:hover {
  background: #fafafa;
  margin: 0 -24px;
  padding: 14px 24px;
}

.item-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.item-icon.blue {
  background: rgba(99, 102, 241, 0.1);
  color: #6366f1;
}

.item-content {
  flex: 1;
  min-width: 0;
}

.item-title {
  font-size: 14px;
  font-weight: 500;
  color: #1e293b;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.item-meta {
  font-size: 12px;
  color: #94a3b8;
  margin-top: 4px;
  display: flex;
  gap: 12px;
}

.item-status {
  font-size: 12px;
  font-weight: 500;
  padding: 4px 10px;
  border-radius: 6px;
}

.item-status.published {
  background: rgba(16, 185, 129, 0.1);
  color: #10b981;
}

.item-status.draft {
  background: rgba(245, 158, 11, 0.1);
  color: #f59e0b;
}

/* 快捷操作 */
.quick-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.quick-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  padding: 20px 16px;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s;
  background: #f8fafc;
}

.quick-item:hover {
  background: #f1f5f9;
  transform: translateY(-4px);
}

.quick-icon {
  width: 52px;
  height: 52px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.quick-icon.purple { background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%); }
.quick-icon.blue { background: linear-gradient(135deg, #3b82f6 0%, #60a5fa 100%); }
.quick-icon.green { background: linear-gradient(135deg, #10b981 0%, #34d399 100%); }
.quick-icon.orange { background: linear-gradient(135deg, #f59e0b 0%, #fbbf24 100%); }

.quick-item span {
  font-size: 13px;
  font-weight: 500;
  color: #475569;
}

/* 留言列表 */
.message-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.message-item {
  display: flex;
  gap: 12px;
}

.msg-avatar {
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  color: #fff;
  font-size: 14px;
  font-weight: 600;
  flex-shrink: 0;
}

.msg-content {
  flex: 1;
  min-width: 0;
}

.msg-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.msg-name {
  font-size: 14px;
  font-weight: 600;
  color: #1e293b;
}

.msg-time {
  font-size: 12px;
  color: #94a3b8;
}

.msg-text {
  font-size: 13px;
  color: #64748b;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* 系统信息 */
.system-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.system-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid #f1f5f9;
}

.system-item:last-child {
  border-bottom: none;
}

.system-label {
  font-size: 13px;
  color: #64748b;
}

.system-value {
  font-size: 13px;
  font-weight: 500;
  color: #1e293b;
}

/* 响应式 */
@media (max-width: 1200px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .dashboard-main {
    flex-direction: column;
  }
  
  .dashboard-right {
    width: 100%;
  }
  
  .quick-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
  
  .quick-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
