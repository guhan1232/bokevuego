<template>
  <div class="login-container">
    <!-- 左侧品牌区域 -->
    <div class="login-brand">
      <div class="brand-content">
        <div class="brand-logo">
          <div class="logo-icon">
            <el-icon :size="48"><Edit /></el-icon>
          </div>
          <span class="logo-text">BokeUI</span>
        </div>
        <h1 class="brand-title">分享 · 记录 · 成长</h1>
        <p class="brand-desc">创作属于你的精彩故事</p>
        <div class="brand-features">
          <div class="feature-item">
            <el-icon><Document /></el-icon>
            <span>文章管理</span>
          </div>
          <div class="feature-item">
            <el-icon><ChatDotRound /></el-icon>
            <span>互动交流</span>
          </div>
          <div class="feature-item">
            <el-icon><DataAnalysis /></el-icon>
            <span>数据分析</span>
          </div>
        </div>
      </div>
      <div class="brand-footer">
        <span>© 2024 BokeUI Blog System</span>
      </div>
    </div>

    <!-- 右侧登录区域 -->
    <div class="login-form-wrapper">
      <div class="login-form-container">
        <div class="login-tabs">
          <span :class="['tab', { active: loginType === 'account' }]" @click="loginType = 'account'">账号登录</span>
          <span :class="['tab', { active: loginType === 'phone' }]" @click="loginType = 'phone'">手机登录</span>
          <div class="tab-indicator" :style="{ left: loginType === 'account' ? '0' : '50%' }"></div>
        </div>

        <!-- 账号密码登录 -->
        <el-form v-if="loginType === 'account'" :model="form" :rules="rules" ref="formRef" @submit.prevent="handleLogin">
          <el-form-item prop="username">
            <el-input v-model="form.username" placeholder="请输入用户名" size="large" prefix-icon="User" />
          </el-form-item>
          <el-form-item prop="password">
            <el-input v-model="form.password" type="password" placeholder="请输入密码" size="large" prefix-icon="Lock" show-password @keyup.enter="handleLogin" />
          </el-form-item>
          <div class="form-options">
            <el-checkbox v-model="rememberMe">记住我</el-checkbox>
            <a href="javascript:;" class="forgot-link">忘记密码?</a>
          </div>
          <el-form-item>
            <el-button type="primary" size="large" :loading="loading" class="login-btn" @click="handleLogin">
              {{ loading ? '登录中...' : '登 录' }}
            </el-button>
          </el-form-item>
        </el-form>

        <!-- 手机号登录 -->
        <el-form v-else :model="phoneForm" :rules="phoneRules" ref="phoneFormRef">
          <el-form-item prop="phone">
            <el-input v-model="phoneForm.phone" placeholder="请输入手机号" size="large" prefix-icon="Phone" />
          </el-form-item>
          <el-form-item prop="code">
            <div class="code-input">
              <el-input v-model="phoneForm.code" placeholder="请输入验证码" size="large" prefix-icon="Key" />
              <el-button size="large" :disabled="countdown > 0" @click="sendCode">
                {{ countdown > 0 ? `${countdown}s后重发` : '获取验证码' }}
              </el-button>
            </div>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" size="large" class="login-btn" @click="handlePhoneLogin">
              登 录
            </el-button>
          </el-form-item>
        </el-form>

        <div class="login-divider">
          <span>其他登录方式</span>
        </div>
        <div class="social-login">
          <a href="javascript:;" class="social-icon wechat" title="微信登录">
            <el-icon :size="20"><ChatDotRound /></el-icon>
          </a>
          <a href="javascript:;" class="social-icon weibo" title="微博登录">
            <el-icon :size="20"><Connection /></el-icon>
          </a>
          <a href="javascript:;" class="social-icon github" title="GitHub登录">
            <el-icon :size="20"><Link /></el-icon>
          </a>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { login } from '../api'
import { useUserStore } from '../stores/user'
import { ElMessage } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()
const formRef = ref()
const phoneFormRef = ref()
const loading = ref(false)
const loginType = ref('account')
const rememberMe = ref(false)
const countdown = ref(0)

const form = reactive({ username: '', password: '' })
const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
}

const phoneForm = reactive({ phone: '', code: '' })
const phoneRules = {
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '手机号格式不正确', trigger: 'blur' }
  ],
  code: [{ required: true, message: '请输入验证码', trigger: 'blur' }],
}

const handleLogin = async () => {
  await formRef.value.validate()
  loading.value = true
  try {
    const data = await login(form)
    userStore.setToken(data.token)
    await userStore.fetchProfile()
    ElMessage.success('登录成功')
    router.push('/')
  } catch {
    // error handled in interceptor
  } finally {
    loading.value = false
  }
}

const sendCode = () => {
  if (!phoneForm.phone || !/^1[3-9]\d{9}$/.test(phoneForm.phone)) {
    ElMessage.warning('请输入正确的手机号')
    return
  }
  countdown.value = 60
  const timer = setInterval(() => {
    countdown.value--
    if (countdown.value <= 0) clearInterval(timer)
  }, 1000)
  ElMessage.success('验证码已发送')
}

const handlePhoneLogin = () => {
  ElMessage.info('手机登录功能演示')
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  background: #f5f7fa;
}

/* 左侧品牌区域 */
.login-brand {
  flex: 1;
  background: linear-gradient(135deg, #1e3c72 0%, #2a5298 50%, #1e3c72 100%);
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: 40px;
  position: relative;
  overflow: hidden;
}

.login-brand::before {
  content: '';
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23ffffff' fill-opacity='0.05'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E");
  animation: float 20s linear infinite;
}

@keyframes float {
  0% { transform: translate(0, 0) rotate(0deg); }
  100% { transform: translate(-50px, -50px) rotate(360deg); }
}

.brand-content {
  position: relative;
  z-index: 1;
  text-align: center;
  color: #fff;
}

.brand-logo {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  margin-bottom: 32px;
}

.logo-icon {
  width: 72px;
  height: 72px;
  border-radius: 16px;
  background: linear-gradient(135deg, rgba(255,255,255,0.2), rgba(255,255,255,0.1));
  backdrop-filter: blur(10px);
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid rgba(255,255,255,0.2);
}

.logo-text {
  font-size: 32px;
  font-weight: 700;
  letter-spacing: 2px;
}

.brand-title {
  font-size: 28px;
  font-weight: 300;
  margin-bottom: 12px;
  letter-spacing: 4px;
}

.brand-desc {
  font-size: 16px;
  opacity: 0.8;
  margin-bottom: 48px;
}

.brand-features {
  display: flex;
  gap: 32px;
  justify-content: center;
}

.feature-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 16px 24px;
  background: rgba(255,255,255,0.1);
  border-radius: 12px;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255,255,255,0.15);
  transition: all 0.3s ease;
}

.feature-item:hover {
  background: rgba(255,255,255,0.2);
  transform: translateY(-4px);
}

.feature-item span {
  font-size: 14px;
  opacity: 0.9;
}

.brand-footer {
  position: absolute;
  bottom: 24px;
  font-size: 12px;
  opacity: 0.6;
}

/* 右侧登录区域 */
.login-form-wrapper {
  width: 480px;
  min-width: 480px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
  background: #fff;
}

.login-form-container {
  width: 100%;
  max-width: 380px;
}

.login-tabs {
  display: flex;
  position: relative;
  margin-bottom: 32px;
  border-bottom: 1px solid #e4e7ed;
}

.tab {
  flex: 1;
  text-align: center;
  padding: 16px 0;
  font-size: 16px;
  color: #909399;
  cursor: pointer;
  transition: all 0.3s ease;
  position: relative;
}

.tab.active {
  color: #409eff;
  font-weight: 500;
}

.tab-indicator {
  position: absolute;
  bottom: -1px;
  width: 50%;
  height: 2px;
  background: #409eff;
  transition: left 0.3s ease;
}

.form-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.forgot-link {
  color: #409eff;
  font-size: 14px;
  text-decoration: none;
}

.forgot-link:hover {
  text-decoration: underline;
}

.login-btn {
  width: 100%;
  height: 48px;
  font-size: 16px;
  border-radius: 8px;
}

.code-input {
  display: flex;
  gap: 12px;
}

.code-input .el-input {
  flex: 1;
}

.code-input .el-button {
  width: 120px;
}

.login-divider {
  display: flex;
  align-items: center;
  margin: 24px 0;
}

.login-divider::before,
.login-divider::after {
  content: '';
  flex: 1;
  height: 1px;
  background: #e4e7ed;
}

.login-divider span {
  padding: 0 16px;
  font-size: 13px;
  color: #909399;
}

.social-login {
  display: flex;
  justify-content: center;
  gap: 24px;
  margin-bottom: 24px;
}

.social-icon {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  transition: all 0.3s ease;
}

.social-icon:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 12px rgba(0,0,0,0.15);
}

.social-icon.wechat {
  background: #07c160;
}

.social-icon.weibo {
  background: #e6162d;
}

.social-icon.github {
  background: #24292e;
}

.login-tip {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 12px 16px;
  background: #fdf6ec;
  border-radius: 8px;
  font-size: 13px;
  color: #e6a23c;
}

/* 响应式 */
@media (max-width: 900px) {
  .login-brand {
    display: none;
  }
  
  .login-form-wrapper {
    width: 100%;
    min-width: auto;
  }
}
</style>
