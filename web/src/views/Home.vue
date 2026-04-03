<template>
  <div class="home-page">
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
            <span>|</span>
            <a href="javascript:;">RSS订阅</a>
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

    <!-- 头条文章 -->
    <section class="headline-section" v-if="headlineArticle">
      <div class="container">
        <div class="headline-wrapper">
          <div class="headline-main">
            <h1 class="headline-title" @click="$router.push(`/article/${headlineArticle.id}`)">
              {{ headlineArticle.title }}
            </h1>
            <p class="headline-summary">{{ headlineArticle.summary || '点击阅读更多内容...' }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- 轮播图 -->
    <section class="carousel-section" v-if="carouselArticles.length > 0">
      <div class="container">
        <div class="carousel-wrapper">
          <div class="carousel-slide" v-for="(article, index) in carouselArticles" :key="article.id" v-show="index === currentSlide" @click="$router.push(`/article/${article.id}`)">
            <img v-if="article.cover" :src="article.cover" :alt="article.title" class="carousel-image" />
            <div v-else class="carousel-image" style="display: flex; align-items: center; justify-content: center; font-size: 48px; color: rgba(255,255,255,0.8);">📝</div>
            <div class="carousel-caption">
              <h3>{{ article.title }}</h3>
              <p>{{ article.summary || '' }}</p>
            </div>
          </div>
          <div class="carousel-controls">
            <button class="carousel-btn" @click.stop="prevSlide">‹</button>
            <button class="carousel-btn" @click.stop="nextSlide">›</button>
          </div>
          <div class="carousel-dots">
            <span v-for="(article, index) in carouselArticles" :key="index" class="carousel-dot" :class="{ active: index === currentSlide }" @click.stop="currentSlide = index"></span>
          </div>
        </div>
      </div>
    </section>

    <!-- 快捷入口 -->
    <section class="content-section">
      <div class="container">
        <div class="special-section">
          <div class="special-item" @click="$router.push('/articles?category=技术分享')">
            <span>💻 技术分享</span>
          </div>
          <div class="special-item" @click="$router.push('/articles?category=生活随笔')">
            <span>🌸 生活随笔</span>
          </div>
          <div class="special-item" @click="$router.push('/articles?category=学习笔记')">
            <span>📚 学习笔记</span>
          </div>
          <div class="special-item" @click="$router.push('/contact')">
            <span>💬 留言交流</span>
          </div>
        </div>
      </div>
    </section>

    <!-- 内容区域 -->
    <section class="content-section">
      <div class="container">
        <div class="content-wrapper">
          <!-- 主内容区 -->
          <div class="content-main">
            <!-- 文章列表 -->
            <div class="news-card">
              <div class="news-header">
                <h2>最新文章</h2>
                <div class="news-tabs">
                  <span class="news-tab" :class="{ active: activeTab === 'all' }" @click="activeTab = 'all'; fetchArticles()">全部</span>
                  <span class="news-tab" :class="{ active: activeTab === '技术分享' }" @click="activeTab = '技术分享'; fetchArticles()">技术分享</span>
                  <span class="news-tab" :class="{ active: activeTab === '生活随笔' }" @click="activeTab = '生活随笔'; fetchArticles()">生活随笔</span>
                  <span class="news-tab" :class="{ active: activeTab === '学习笔记' }" @click="activeTab = '学习笔记'; fetchArticles()">学习笔记</span>
                </div>
                <router-link to="/articles" class="news-more">更多 &gt;</router-link>
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
                <div v-if="articles.length === 0 && !loading" class="loading">
                  <p style="color: #999; padding: 40px 0; text-align: center;">暂无文章，去后台发布第一篇吧~</p>
                </div>
              </div>
            </div>
          </div>

          <!-- 侧边栏 -->
          <aside class="content-sidebar">
            <div class="quick-entry">
              <h3>文章分类</h3>
              <div class="quick-grid">
                <div class="quick-item" @click="$router.push('/articles')">
                  <div class="quick-icon blue">📚</div>
                  <span class="quick-text">全部文章</span>
                </div>
                <div class="quick-item" @click="$router.push('/contact')">
                  <div class="quick-icon red">💬</div>
                  <span class="quick-text">留言板</span>
                </div>
                <div class="quick-item" @click="$router.push('/articles?category=技术分享')">
                  <div class="quick-icon green">💻</div>
                  <span class="quick-text">技术分享</span>
                </div>
                <div class="quick-item" @click="$router.push('/articles?category=生活随笔')">
                  <div class="quick-icon orange">🌸</div>
                  <span class="quick-text">生活随笔</span>
                </div>
              </div>
            </div>
          </aside>
        </div>
      </div>
    </section>

    <!-- 底部导航 -->
    <section class="footer-nav">
      <div class="container">
        <div class="footer-grid">
          <div class="footer-col">
            <h4>文章分类</h4>
            <ul>
              <li><router-link to="/articles">最新文章</router-link></li>
              <li><router-link to="/articles?category=技术分享">技术分享</router-link></li>
              <li><router-link to="/articles?category=生活随笔">生活随笔</router-link></li>
            </ul>
          </div>
          <div class="footer-col">
            <h4>内容标签</h4>
            <ul>
              <li><router-link to="/articles?category=学习笔记">学习笔记</router-link></li>
              <li><a href="javascript:;">前端开发</a></li>
              <li><a href="javascript:;">后端技术</a></li>
            </ul>
          </div>
          <div class="footer-col">
            <h4>关于博客</h4>
            <ul>
              <li><a href="javascript:;">关于我</a></li>
              <li><a href="javascript:;">博客历程</a></li>
              <li><router-link to="/contact">联系我</router-link></li>
            </ul>
          </div>
          <div class="footer-col">
            <h4>互动交流</h4>
            <ul>
              <li><router-link to="/contact">留言板</router-link></li>
              <li><a href="javascript:;">友链申请</a></li>
              <li><a href="javascript:;">合作洽谈</a></li>
            </ul>
          </div>
          <div class="footer-col">
            <h4>更多</h4>
            <ul>
              <li><a href="javascript:;">RSS订阅</a></li>
              <li><a href="/api/sitemap.xml" target="_blank">网站地图</a></li>
              <li><a href="javascript:;">免责声明</a></li>
            </ul>
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
      <p>
        <a href="/api/sitemap.xml" target="_blank">网站地图</a>
        <a href="/contact">联系我</a>
      </p>
      <div v-if="friendLinks.length > 0" class="friend-links">
        <span>友情链接：</span>
        <a v-for="link in friendLinks" :key="link.id" :href="link.url" target="_blank" rel="noopener noreferrer">
          {{ link.name }}
        </a>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, onMounted, inject, watch } from 'vue'
import { useRouter } from 'vue-router'
import { getArticles, getFriendLinks } from '../api'

const router = useRouter()
const siteConfig = inject('siteConfig')

const articles = ref([])
const carouselArticles = ref([])
const headlineArticle = ref(null)
const loading = ref(false)
const currentSlide = ref(0)
const activeTab = ref('all')
const friendLinks = ref([])
const mobileMenuOpen = ref(false)
const searchKeyword = ref('')

let slideInterval = null

const fetchArticles = async () => {
  loading.value = true
  try {
    const params = { page: 1, page_size: 10 }
    if (activeTab.value !== 'all') {
      params.category = activeTab.value
    }
    const data = await getArticles(params)
    articles.value = data.items || []

    // 设置头条（第一条文章）
    if (articles.value.length > 0) {
      headlineArticle.value = articles.value[0]
    }

    // 设置轮播图（前5条有封面的文章）
    const withCovers = articles.value.filter(a => a.cover).slice(0, 5)
    if (withCovers.length > 0) {
      carouselArticles.value = withCovers
    } else if (articles.value.length > 0) {
      carouselArticles.value = articles.value.slice(0, 5)
    }
  } catch (error) {
    console.error('Failed to fetch articles:', error)
  } finally {
    loading.value = false
  }
}

const fetchFriendLinks = async () => {
  try {
    const data = await getFriendLinks()
    friendLinks.value = data || []
  } catch (e) {
    console.error('Failed to fetch friend links:', e)
  }
}

const prevSlide = () => {
  currentSlide.value = (currentSlide.value - 1 + carouselArticles.value.length) % carouselArticles.value.length
}

const nextSlide = () => {
  currentSlide.value = (currentSlide.value + 1) % carouselArticles.value.length
}

const handleSearch = () => {
  if (searchKeyword.value.trim()) {
    router.push(`/articles?keyword=${encodeURIComponent(searchKeyword.value.trim())}`)
  }
}

const formatDateDay = (dateStr) => {
  if (!dateStr) return '--'
  const date = new Date(dateStr)
  return String(date.getDate()).padStart(2, '0')
}

const formatDateMonth = (dateStr) => {
  if (!dateStr) return '--'
  const date = new Date(dateStr)
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`
}

onMounted(() => {
  fetchArticles()
  fetchFriendLinks()

  // 自动轮播
  if (carouselArticles.value.length > 1) {
    slideInterval = setInterval(() => {
      nextSlide()
    }, 5000)
  }
})

// 清理定时器
watch(carouselArticles, (val) => {
  if (slideInterval) {
    clearInterval(slideInterval)
  }
  if (val.length > 1) {
    slideInterval = setInterval(() => {
      nextSlide()
    }, 5000)
  }
})
</script>

<style scoped>
/* 组件特定样式已在全局 style.css 中定义 */
</style>
