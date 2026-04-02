<template>
  <div class="contact-page">
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
        <span>留言板</span>
      </div>
    </div>

    <!-- 留言表单 -->
    <section class="content-section">
      <div class="container">
        <div class="news-card">
          <div class="news-header">
            <h2>留言板</h2>
          </div>
          <div style="padding: 30px;">
            <div class="form-group">
              <label>昵称 <span style="color: #c00;">*</span></label>
              <input type="text" v-model="form.name" placeholder="请输入您的昵称" />
            </div>
            <div class="form-group">
              <label>网站/博客</label>
              <input type="text" v-model="form.phone" placeholder="您的个人网站或博客地址（选填）" />
            </div>
            <div class="form-group">
              <label>电子邮箱 <span style="color: #c00;">*</span></label>
              <input type="email" v-model="form.email" placeholder="请输入电子邮箱（用于接收回复通知）" />
            </div>
            <div class="form-group">
              <label>留言内容 <span style="color: #c00;">*</span></label>
              <textarea v-model="form.content" placeholder="想说点什么...欢迎交流讨论~" rows="6"></textarea>
            </div>
            <div class="form-actions">
              <button class="submit-btn" @click="handleSubmit" :disabled="submitting">
                {{ submitting ? '提交中...' : '提交留言' }}
              </button>
              <button class="reset-btn" @click="resetForm">重置</button>
            </div>
          </div>
        </div>

        <!-- 留言须知 -->
        <div class="news-card" style="margin-top: 20px;">
          <div class="news-header">
            <h2>留言须知</h2>
          </div>
          <div style="padding: 20px; line-height: 2;">
            <p>1. 欢迎在留言板与我交流，我会尽快回复您的留言。</p>
            <p>2. 请留下您的邮箱，方便我回复后通知您。</p>
            <p>3. 请文明留言，勿发表违反法律法规的言论。</p>
            <p>4. 如有技术问题，也可以在 GitHub 上提 issue 交流。</p>
          </div>
        </div>
      </div>
    </section>

    <!-- 提示消息 -->
    <div v-if="toastMessage" class="toast-message">{{ toastMessage }}</div>

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
import { ref, reactive, inject } from 'vue'
import { useRouter } from 'vue-router'
import { createMessage } from '../api'

const router = useRouter()
const siteConfig = inject('siteConfig')

const form = reactive({
  name: '',
  phone: '',
  email: '',
  content: ''
})

const submitting = ref(false)
const toastMessage = ref('')
const mobileMenuOpen = ref(false)
const searchKeyword = ref('')

const handleSubmit = async () => {
  if (!form.name.trim()) {
    showToast('请输入昵称')
    return
  }
  if (!form.email.trim()) {
    showToast('请输入电子邮箱')
    return
  }
  if (!form.content.trim()) {
    showToast('请输入留言内容')
    return
  }

  submitting.value = true
  try {
    await createMessage({ ...form })
    showToast('留言提交成功，感谢您的留言！')
    resetForm()
  } catch {
    showToast('提交失败，请稍后重试')
  } finally {
    submitting.value = false
  }
}

const resetForm = () => {
  form.name = ''
  form.phone = ''
  form.email = ''
  form.content = ''
}

const showToast = (msg) => {
  toastMessage.value = msg
  setTimeout(() => {
    toastMessage.value = ''
  }, 3000)
}

const handleSearch = () => {
  if (searchKeyword.value.trim()) {
    router.push(`/articles?keyword=${encodeURIComponent(searchKeyword.value.trim())}`)
  }
}
</script>

<style scoped>
.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-dark);
  margin-bottom: 8px;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  font-size: 14px;
  font-family: inherit;
  transition: border-color 0.2s;
  outline: none;
}

.form-group input:focus,
.form-group textarea:focus {
  border-color: var(--primary-blue);
}

.form-group textarea {
  resize: vertical;
  min-height: 120px;
}

.form-actions {
  display: flex;
  gap: 15px;
  margin-top: 30px;
}

.submit-btn,
.reset-btn {
  padding: 10px 30px;
  font-size: 14px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
}

.submit-btn {
  background: var(--primary-blue);
  color: #fff;
}

.submit-btn:hover:not(:disabled) {
  background: var(--nav-hover);
}

.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.reset-btn {
  background: #f5f5f5;
  color: var(--text-gray);
  border: 1px solid var(--border-color);
}

.reset-btn:hover {
  background: #e8e8e8;
}

.toast-message {
  position: fixed;
  top: 100px;
  left: 50%;
  transform: translateX(-50%);
  padding: 12px 24px;
  background: var(--primary-blue);
  color: #fff;
  border-radius: 4px;
  font-size: 14px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  z-index: 1000;
  animation: fadeInDown 0.3s ease;
}

@keyframes fadeInDown {
  from {
    opacity: 0;
    transform: translate(-50%, -20px);
  }
  to {
    opacity: 1;
    transform: translate(-50%, 0);
  }
}
</style>
