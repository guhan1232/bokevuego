<template>
  <div class="article-detail-page">
    <!-- 顶部信息栏 -->
    <header class="top-bar">
      <div class="container">
        <div class="top-left">
          <div class="logo-wrapper">
            <div class="logo-emblem">
              <img v-if="siteConfig.site_logo" :src="siteConfig.site_logo" :alt="siteConfig.site_title || 'BokeUI 博客'" class="logo-image" />
              <span v-else>✍</span>
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
            <input v-model="searchKeyword" placeholder="搜索文章..." @keyup.enter="handleSearch" />
            <button @click="handleSearch">搜索</button>
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
        <router-link to="/articles">文章列表</router-link>
        <span>&gt;</span>
        <span>正文</span>
      </div>
    </div>

    <!-- 文章内容 -->
    <section class="article-page">
      <div class="container">
        <div class="article-wrapper" v-if="article">
          <div class="article-header">
            <h1>{{ article.title }}</h1>
            <div class="article-meta">
              <span>发布时间：{{ formatDate(article.created_at) }}</span>
              <span v-if="article.author">作者：{{ article.author }}</span>
              <span>阅读量：{{ article.views }}</span>
            </div>
          </div>
          <div class="article-content" v-html="renderedContent"></div>

          <!-- 文章标签 -->
          <div class="article-footer" v-if="article.category">
            <span class="article-tag">分类：{{ article.category }}</span>
          </div>
        </div>

        <div v-else-if="!loading" class="article-wrapper">
          <div class="article-header">
            <h1>文章不存在或已被删除</h1>
          </div>
        </div>

        <div v-if="loading" class="loading">
          <div class="loading-spinner"></div>
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
import { ref, computed, onMounted, inject } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getArticle, viewArticle } from '../api'
import { marked } from 'marked'

const route = useRoute()
const router = useRouter()
const siteConfig = inject('siteConfig')

const article = ref(null)
const loading = ref(true)
const mobileMenuOpen = ref(false)
const searchKeyword = ref('')

const renderedContent = computed(() => {
  if (!article.value?.content) return ''
  return marked(article.value.content, { breaks: true })
})

const formatDate = (dateStr) => {
  if (!dateStr) return '--'
  const date = new Date(dateStr)
  return `${date.getFullYear()}年${date.getMonth() + 1}月${date.getDate()}日`
}

const handleSearch = () => {
  if (searchKeyword.value.trim()) {
    router.push(`/articles?keyword=${encodeURIComponent(searchKeyword.value.trim())}`)
  }
}

onMounted(async () => {
  try {
    const id = route.params.id
    const [articleData] = await Promise.all([
      getArticle(id),
      viewArticle(id).catch(() => {}),
    ])
    article.value = articleData.article
  } catch {
    article.value = null
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.article-footer {
  margin-top: 30px;
  padding-top: 20px;
  border-top: 1px solid var(--border-color);
}

.article-tag {
  display: inline-block;
  padding: 4px 12px;
  background: #f0f0f0;
  color: var(--text-gray);
  border-radius: 3px;
  font-size: 13px;
}
</style>
