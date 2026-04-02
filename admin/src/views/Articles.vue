<template>
  <div class="page-container">
    <div class="page-header">
      <h2>文章管理</h2>
      <div>
        <el-button v-if="activeTab === 'list'" type="primary" :icon="Plus" @click="$router.push('/articles/new')">新建文章</el-button>
        <el-button v-if="activeTab === 'trash'" type="primary" :icon="Back" @click="activeTab = 'list'">返回文章列表</el-button>
      </div>
    </div>

    <!-- 标签页切换 -->
    <div class="tabs-wrapper">
      <el-radio-group v-model="activeTab" @change="handleTabChange">
        <el-radio-button value="list">
          文章列表
          <el-badge v-if="stats.published_articles > 0 || stats.draft_articles > 0" :value="stats.published_articles + stats.draft_articles" class="tab-badge" />
        </el-radio-button>
        <el-radio-button value="trash">
          回收站
          <el-badge v-if="stats.trash_articles > 0" :value="stats.trash_articles" class="tab-badge" type="warning" />
        </el-radio-button>
      </el-radio-group>
    </div>

    <!-- 文章列表 -->
    <div v-show="activeTab === 'list'">
      <div class="filter-bar">
        <el-input v-model="query.keyword" placeholder="搜索标题..." clearable style="width: 240px" @clear="fetchList" @keyup.enter="fetchList">
          <template #prefix><el-icon><Search /></el-icon></template>
        </el-input>
        <el-select v-model="query.status" placeholder="状态" clearable style="width: 120px" @change="fetchList">
          <el-option label="草稿" :value="0" />
          <el-option label="已发布" :value="1" />
        </el-select>
        <el-select v-model="query.category" placeholder="分类" clearable style="width: 140px" @change="fetchList">
          <el-option v-for="c in categories" :key="c" :label="c" :value="c" />
        </el-select>
      </div>

      <el-table :data="articles" stripe v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="title" label="标题" min-width="200" show-overflow-tooltip />
        <el-table-column prop="category" label="分类" width="120">
          <template #default="{ row }">
            <el-tag v-if="row.category" size="small">{{ row.category }}</el-tag>
            <span v-else class="text-muted">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
              {{ row.status === 1 ? '已发布' : '草稿' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="views" label="浏览" width="80" />
        <el-table-column prop="created_at" label="创建时间" width="170">
          <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="$router.push(`/articles/${row.id}`)">编辑</el-button>
            <el-popconfirm title="确定移入回收站？" @confirm="handleDelete(row.id)">
              <template #reference>
                <el-button type="danger" link size="small">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrap">
        <el-pagination
          v-model:current-page="query.page"
          :page-size="query.page_size"
          :total="total"
          layout="total, prev, pager, next"
          @current-change="fetchList"
        />
      </div>
    </div>

    <!-- 回收站 -->
    <div v-show="activeTab === 'trash'">
      <div class="filter-bar">
        <el-input v-model="trashQuery.keyword" placeholder="搜索标题..." clearable style="width: 240px" @clear="fetchTrashList" @keyup.enter="fetchTrashList">
          <template #prefix><el-icon><Search /></el-icon></template>
        </el-input>
      </div>

      <el-alert type="warning" :closable="false" style="margin-bottom: 16px">
        回收站中的文章将在 30 天后自动清理，如需保留请及时恢复。
      </el-alert>

      <el-table :data="trashArticles" stripe v-loading="trashLoading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="title" label="标题" min-width="200" show-overflow-tooltip />
        <el-table-column prop="category" label="分类" width="120">
          <template #default="{ row }">
            <el-tag v-if="row.category" size="small" type="info">{{ row.category }}</el-tag>
            <span v-else class="text-muted">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="原状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
              {{ row.status === 1 ? '已发布' : '草稿' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="deleted_at" label="删除时间" width="170">
          <template #default="{ row }">{{ formatDate(row.deleted_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="success" link size="small" @click="handleRestore(row.id)">恢复</el-button>
            <el-popconfirm title="确定彻底删除？此操作不可恢复！" @confirm="handleHardDelete(row.id)">
              <template #reference>
                <el-button type="danger" link size="small">彻底删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrap">
        <el-pagination
          v-model:current-page="trashQuery.page"
          :page-size="trashQuery.page_size"
          :total="trashTotal"
          layout="total, prev, pager, next"
          @current-change="fetchTrashList"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getArticles, deleteArticle, getCategories, getTrashArticles, restoreArticle, hardDeleteArticle, getStats } from '../api'
import { Plus, Search, Back } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const activeTab = ref('list')
const articles = ref([])
const categories = ref([])
const total = ref(0)
const loading = ref(false)
const query = reactive({ page: 1, page_size: 10, keyword: '', status: null, category: '' })

const trashArticles = ref([])
const trashTotal = ref(0)
const trashLoading = ref(false)
const trashQuery = reactive({ page: 1, page_size: 10, keyword: '' })

const stats = ref({ published_articles: 0, draft_articles: 0, trash_articles: 0 })

const fetchList = async () => {
  loading.value = true
  try {
    const params = { ...query }
    if (params.status === null || params.status === '') delete params.status
    const data = await getArticles(params)
    articles.value = data.items || []
    total.value = data.total || 0
  } finally {
    loading.value = false
  }
}

const fetchTrashList = async () => {
  trashLoading.value = true
  try {
    const params = { ...trashQuery }
    const data = await getTrashArticles(params)
    trashArticles.value = data.items || []
    trashTotal.value = data.total || 0
  } finally {
    trashLoading.value = false
  }
}

const fetchStats = async () => {
  const data = await getStats()
  stats.value = {
    published_articles: data.published_articles || 0,
    draft_articles: data.draft_articles || 0,
    trash_articles: data.trash_articles || 0
  }
}

const handleTabChange = (tab) => {
  if (tab === 'list') {
    fetchList()
  } else if (tab === 'trash') {
    fetchTrashList()
  }
}

const handleDelete = async (id) => {
  await deleteArticle(id)
  ElMessage.success('已移入回收站')
  fetchList()
  fetchStats()
}

const handleRestore = async (id) => {
  await restoreArticle(id)
  ElMessage.success('恢复成功')
  fetchTrashList()
  fetchStats()
}

const handleHardDelete = async (id) => {
  await hardDeleteArticle(id)
  ElMessage.success('已彻底删除')
  fetchTrashList()
  fetchStats()
}

const formatDate = (d) => {
  if (!d) return '-'
  return new Date(d).toLocaleString('zh-CN')
}

onMounted(() => {
  fetchList()
  fetchStats()
  getCategories().then(d => categories.value = d.items || [])
})
</script>

<style scoped>
.text-muted {
  color: #c0c4cc;
}

.tabs-wrapper {
  margin-bottom: 16px;
}

.tab-badge {
  margin-left: 6px;
}

.tab-badge :deep(.el-badge__content) {
  transform: scale(0.85);
}
</style>
