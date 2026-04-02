<template>
  <div class="page-container">
    <div class="page-header">
      <h2>分类管理</h2>
    </div>

    <el-row :gutter="20">
      <el-col :xs="24" :md="12">
        <el-card shadow="hover">
          <template #header>
            <span class="card-title">所有分类</span>
          </template>
          <div class="category-list" v-loading="loading">
            <div v-for="cat in categories" :key="cat" class="category-item">
              <el-icon><FolderOpened /></el-icon>
              <span>{{ cat }}</span>
              <el-tag size="small" type="info">{{ getCategoryCount(cat) }} 篇</el-tag>
            </div>
            <el-empty v-if="categories.length === 0" description="暂无分类" />
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :md="12">
        <el-card shadow="hover">
          <template #header>
            <span class="card-title">所有标签</span>
          </template>
          <div class="tag-cloud" v-loading="tagLoading">
            <el-tag
              v-for="tag in tags"
              :key="tag"
              size="large"
              class="tag-item"
            >
              # {{ tag }}
            </el-tag>
            <el-empty v-if="tags.length === 0" description="暂无标签" />
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getCategories, getTags, getArticles } from '../api'
import { FolderOpened } from '@element-plus/icons-vue'

const categories = ref([])
const tags = ref([])
const articles = ref([])
const loading = ref(false)
const tagLoading = ref(false)

const getCategoryCount = (cat) => {
  return articles.value.filter(a => a.category === cat).length
}

onMounted(async () => {
  loading.value = true
  tagLoading.value = true
  try {
    categories.value = (await getCategories()).items || []
    tags.value = (await getTags()).items || []
    const data = await getArticles({ page_size: 100 })
    articles.value = data.items || []
  } finally {
    loading.value = false
    tagLoading.value = false
  }
})
</script>

<style scoped>
.card-title {
  font-weight: 600;
  font-size: 15px;
}

.category-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.category-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  background: #f5f7fa;
  border-radius: 8px;
  font-size: 14px;
}

.tag-cloud {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.tag-item {
  cursor: default;
  border-radius: 20px;
}
</style>
