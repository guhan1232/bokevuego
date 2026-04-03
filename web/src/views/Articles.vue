<template>
  <div class="articles-page">
    <!-- 顶部信息栏 -->
    <header class="top-bar">
      <div class="container">
        <div class="top-left">
          <div class="logo-wrapper">
            <div class="logo-emblem">
              <img v-if="siteConfig.site_logo" :src="siteConfig.site_logo" :alt="siteConfig.site_title || 'BokeUI 博客'" class="logo-image" />
              <span v-else>墨</span>
            </div>
            <div class="logo-text">
              <h1>{{ siteConfig.site_title || 'BokeUI 博客' }}</h1>
              <p>{{ siteConfig.site_subtitle || '记录生活，分享技术' }}</p>
            </div>
          </div>
        </div>
        <div class="top-right">
          <div class="top-tools">
            <a href="javascript:;">关于我</a>
            <span>|</span>
            <a href="javascript:;">友链</a>
          </div>
          <div class="top-search">
            <input v-model="keyword" placeholder="搜索文章..." @keyup.enter="search" />
            <button @click="search">搜索</button>
          </div>
        </div>
      </div>
    </header>

    <!-- 主导航栏 -->
    <nav class="main-nav">
      <div class="container">
        <ul class="nav-list" :class="{ open: mobileMenuOpen }">
          <li class="nav-item"><router-link to="/">首页</router-link></li>
          <li class="nav-item"><router-link to="/articles">全部文章</router-link></li>
          <li class="nav-item"><router-link to="/articles?category=技术分享">技术分享</router-link></li>
          <li class="nav-item"><router-link to="/articles?category=生活随笔">生活随笔</router-link></li>
          <li class="nav-item"><router-link to="/contact">留言板</router-link></li>
        </ul>
        <button class="mobile-menu-btn" @click="mobileMenuOpen = !mobileMenuOpen">
          {{ mobileMenuOpen ? '✕' : '☰' }}
        </button>
      </div>
    </nav>

    <!-- 面包屑 -->
    <div class="container">
      <div class="breadcrumb">
        <router-link to="/">首页</router-link>
        <span>&gt;</span>
        <span>文章列表</span>
      </div>
    </div>

    <!-- 内容区域 -->
    <section class="content-section">
      <div class="container">
        <div class="news-card">
          <div class="news-header">
            <h2>{{ currentCategory || '全部文章' }}</h2>
            <div class="news-tabs">
              <span class="news-tab" :class="{ active: !currentCategory }" @click="filterCategory('')">全部</span>
              <span v-for="cat in categories" :key="cat" class="news-tab" :class="{ active: currentCategory === cat }" @click="filterCategory(cat)">{{ cat }}</span>
            </div>
          </div>
          <div class="news-list">
            <div class="news-item" v-for="article in articles" :key="article.id">
              <div class="news-date">
                <span class="day">{{ formatDateDay(article.created_at) }}</span>
                <span class="month">{{ formatDateMonth(article.created_at) }}</span>
              </div>
              <div class="news-title">
                <router-link :to="`/article/${article.id}`">{{ article.title }}</router-link>
              </div>
            </div>
          </div>

          <!-- 分页 -->
          <div class="pagination" v-if="totalPages > 1">
            <button class="page-btn" :disabled="page <= 1" @click="changePage(page - 1)">上一页</button>
            <span class="page-info">第 {{ page }} 页 / 共 {{ totalPages }} 页</span>
            <button class="page-btn" :disabled="page >= totalPages" @click="changePage(page + 1)">下一页</button>
          </div>
        </div>
      </div>
    </section>

    <!-- 页脚 -->
    <footer class="footer">
      <p>{{ siteConfig.site_footer || '© 2026 BokeUI 博客 · 用心记录每一刻' }}</p>
      <p v-if="siteConfig.site_icp" class="footer-icp">
        <a href="https://beian.miit.gov.cn/" target="_blank" rel="noopener noreferrer">{{ siteConfig.site_icp }}</a>
      </p>
    </footer>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, inject, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getArticles, getCategories } from '../api'

const route = useRoute()
const router = useRouter()
const siteConfig = inject('siteConfig')

const articles = ref([])
const categories = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = 15
const loading = ref(false)
const keyword = ref('')
const currentCategory = ref('')
const mobileMenuOpen = ref(false)

const totalPages = computed(() => Math.ceil(total.value / pageSize))

const fetchArticles = async () => {
  loading.value = true
  try {
    const params = { page: page.value, page_size: pageSize }
    if (keyword.value) params.keyword = keyword.value
    if (currentCategory.value) params.category = currentCategory.value

    const data = await getArticles(params)
    articles.value = data.items || []
    total.value = data.total || 0
  } finally {
    loading.value = false
  }
}

const fetchCategories = async () => {
  try {
    const data = await getCategories()
    categories.value = data.items || []
  } catch {}
}

const filterCategory = (cat) => {
  currentCategory.value = cat
  page.value = 1
  router.push({ path: '/articles', query: cat ? { category: cat } : {} })
  fetchArticles()
}

const search = () => {
  page.value = 1
  router.push({ path: '/articles', query: { keyword: keyword.value } })
  fetchArticles()
}

const changePage = (newPage) => {
  page.value = newPage
  fetchArticles()
  window.scrollTo(0, 0)
}

const formatDateDay = (dateStr) => {
  if (!dateStr) return '--'
  return String(new Date(dateStr).getDate()).padStart(2, '0')
}

const formatDateMonth = (dateStr) => {
  if (!dateStr) return '--'
  const date = new Date(dateStr)
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`
}

// 监听路由参数
watch(() => route.query, (query) => {
  if (query.category) currentCategory.value = query.category
  if (query.keyword) keyword.value = query.keyword
  page.value = 1
  fetchArticles()
}, { immediate: true })

onMounted(() => {
  fetchCategories()
})
</script>

<style scoped>
.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 24px;
  padding: 28px;
  border-top: 2px dashed var(--border);
}

.page-btn {
  padding: 10px 24px;
  background: var(--ink-black);
  color: var(--paper-white);
  border: 2px solid var(--ink-black);
  border-radius: var(--radius);
  cursor: pointer;
  font-size: 14px;
  font-family: "Noto Serif SC", serif;
  letter-spacing: 1px;
  transition: all 0.3s ease;
}

.page-btn:hover:not(:disabled) {
  background: var(--seal-red);
  border-color: var(--seal-red);
  transform: translateY(-2px);
  box-shadow: var(--shadow);
}

.page-btn:disabled {
  background: var(--paper-dark);
  border-color: var(--paper-dark);
  cursor: not-allowed;
  transform: none;
}

.page-info {
  font-size: 14px;
  color: var(--text-secondary);
  letter-spacing: 1px;
}
</style>
