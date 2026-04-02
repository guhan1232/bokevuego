<template>
  <div class="page-container">
    <div class="page-header">
      <h2>站点设置</h2>
    </div>

    <el-tabs v-model="activeTab" class="settings-tabs">
      <!-- 基础设置 -->
      <el-tab-pane label="基础设置" name="basic">
        <el-card shadow="hover" v-loading="loading">
          <el-form :model="form" label-width="120px" label-position="top">
            <el-row :gutter="20">
              <el-col :xs="24" :md="12">
                <el-form-item label="站点标题">
                  <el-input v-model="form.site_title" />
                </el-form-item>
              </el-col>
              <el-col :xs="24" :md="12">
                <el-form-item label="站点副标题">
                  <el-input v-model="form.site_subtitle" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="20">
              <el-col :xs="24" :md="12">
                <el-form-item label="Logo 图片">
                  <el-input v-model="form.site_logo" placeholder="Logo 图片 URL（建议 50x50）" />
                  <div class="form-tip">支持 PNG/JPG/SVG 格式，建议尺寸 50x50 像素</div>
                </el-form-item>
              </el-col>
              <el-col :xs="24" :md="12">
                <el-form-item label="Favicon 图标">
                  <el-input v-model="form.site_favicon" placeholder="网站图标 URL（建议 32x32）" />
                  <div class="form-tip">浏览器标签图标，建议使用 ICO 或 PNG 格式</div>
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="20">
              <el-col :xs="24" :md="12">
                <el-form-item label="页脚文本">
                  <el-input v-model="form.site_footer" />
                </el-form-item>
              </el-col>
              <el-col :xs="24" :md="12">
                <el-form-item label="备案号">
                  <el-input v-model="form.site_icp" placeholder="如：京ICP备XXXXXXXX号" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="20">
              <el-col :xs="24" :md="12">
                <el-form-item label="导航栏背景色">
                  <el-color-picker v-model="form.nav_bg" />
                </el-form-item>
              </el-col>
              <el-col :xs="24" :md="12">
                <el-form-item label="首页横幅背景">
                  <el-input v-model="form.hero_bg" placeholder="CSS 渐变或图片URL" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-form-item label="网站背景图">
              <el-input v-model="form.site_bg_image" placeholder="背景图片URL（留空则使用默认背景）" />
              <div class="form-tip">支持 JPG/PNG/WebP 格式，建议尺寸 1920x1080 以上</div>
            </el-form-item>
          </el-form>
        </el-card>

        <el-card shadow="hover" style="margin-top: 20px">
          <template #header><span style="font-weight: 600">前台预览</span></template>
          <div class="preview-bar" :style="{ background: form.hero_bg || '#667eea' }">
            <span style="color: #fff; font-size: 18px; font-weight: 600">{{ form.site_title }}</span>
            <span style="color: rgba(255,255,255,0.8); font-size: 14px; margin-left: 16px">{{ form.site_subtitle }}</span>
          </div>
        </el-card>
      </el-tab-pane>

      <!-- 邮件配置 -->
      <el-tab-pane label="邮件配置" name="email">
        <el-card shadow="hover" v-loading="loading">
          <el-alert type="info" :closable="false" style="margin-bottom: 20px">
            配置邮件服务后，当有新留言时会自动发送通知邮件
          </el-alert>
          <el-form :model="form" label-width="120px" label-position="top">
            <el-row :gutter="20">
              <el-col :xs="24" :md="12">
                <el-form-item label="SMTP 服务器">
                  <el-input v-model="form.smtp_host" placeholder="如：smtp.qq.com" />
                </el-form-item>
              </el-col>
              <el-col :xs="24" :md="12">
                <el-form-item label="SMTP 端口">
                  <el-input v-model="form.smtp_port" placeholder="如：587 或 465" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="20">
              <el-col :xs="24" :md="12">
                <el-form-item label="邮箱账号">
                  <el-input v-model="form.smtp_user" placeholder="发件邮箱地址" />
                </el-form-item>
              </el-col>
              <el-col :xs="24" :md="12">
                <el-form-item label="邮箱密码/授权码">
                  <el-input v-model="form.smtp_pass" type="password" show-password placeholder="邮箱密码或授权码" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-form-item label="收件邮箱">
              <el-input v-model="form.smtp_to" placeholder="接收通知的邮箱地址（多个用逗号分隔）" />
            </el-form-item>
          </el-form>
        </el-card>
      </el-tab-pane>

      <!-- SEO 设置 -->
      <el-tab-pane label="SEO 设置" name="seo">
        <el-card shadow="hover" v-loading="loading">
          <el-alert type="info" :closable="false" style="margin-bottom: 20px">
            IndexNow 是一种通知搜索引擎内容更新的协议，支持 Bing、Yandex 等搜索引擎
          </el-alert>
          <el-form :model="form" label-width="120px" label-position="top">
            <el-form-item label="IndexNow Key">
              <el-input v-model="form.indexnow_key" placeholder="输入 IndexNow Key（至少8位字符）" />
              <div class="form-tip">
                生成方式：随机生成一个字符串（如 UUID），或访问 
                <a href="https://www.bing.com/indexnow" target="_blank">Bing IndexNow</a> 获取
              </div>
            </el-form-item>
            <el-form-item label="启用 IndexNow">
              <el-switch v-model="indexnowEnabled" />
              <span style="margin-left: 10px; color: #909399; font-size: 13px">
                开启后发布/更新文章会自动推送至搜索引擎
              </span>
            </el-form-item>
          </el-form>

          <div class="seo-actions">
            <el-button type="primary" @click="generateSitemap">
              <el-icon><Link /></el-icon>
              生成站点地图
            </el-button>
            <el-button type="success" @click="pushIndexNow" :loading="pushing">
              <el-icon><Promotion /></el-icon>
              推送到 IndexNow
            </el-button>
          </div>

          <el-card shadow="never" style="margin-top: 20px" v-if="sitemapUrl">
            <template #header><span>站点地图</span></template>
            <div class="sitemap-preview">
              <el-link :href="sitemapUrl" target="_blank" type="primary">{{ sitemapUrl }}</el-link>
            </div>
          </el-card>
        </el-card>
      </el-tab-pane>
    </el-tabs>

    <div class="save-bar">
      <el-button type="primary" size="large" @click="handleSave" :loading="saving">
        保存设置
      </el-button>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref, computed, onMounted } from 'vue'
import { getConfigs, updateConfigs } from '../api'
import { ElMessage } from 'element-plus'
import request from '../api/request'

const loading = ref(false)
const saving = ref(false)
const pushing = ref(false)
const activeTab = ref('basic')
const sitemapUrl = ref('')

const form = reactive({
  site_title: '',
  site_subtitle: '',
  site_footer: '',
  site_icp: '',
  site_bg_image: '',
  site_logo: '',
  site_favicon: '',
  nav_bg: '#1a1a2e',
  hero_bg: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
  smtp_host: '',
  smtp_port: '587',
  smtp_user: '',
  smtp_pass: '',
  smtp_to: '',
  indexnow_key: '',
  indexnow_enabled: 'false',
})

const indexnowEnabled = computed({
  get: () => form.indexnow_enabled === 'true',
  set: (val) => { form.indexnow_enabled = val ? 'true' : 'false' }
})

const fetchConfigs = async () => {
  loading.value = true
  try {
    const data = await getConfigs()
    Object.assign(form, data)
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  saving.value = true
  try {
    await updateConfigs({ ...form })
    ElMessage.success('保存成功')
  } finally {
    saving.value = false
  }
}

const generateSitemap = () => {
  sitemapUrl.value = `${window.location.origin}/api/sitemap.xml`
  ElMessage.success('站点地图已生成')
}

const pushIndexNow = async () => {
  if (!form.indexnow_key) {
    ElMessage.warning('请先配置 IndexNow Key')
    return
  }
  
  pushing.value = true
  try {
    const res = await request.post('/indexnow/push', {
      base_url: window.location.origin,
      key: form.indexnow_key
    })
    ElMessage.success(`推送成功，共 ${res.count} 个URL`)
  } catch (e) {
    ElMessage.error('推送失败：' + (e.response?.data?.error || e.message))
  } finally {
    pushing.value = false
  }
}

onMounted(fetchConfigs)
</script>

<style scoped>
.settings-tabs {
  background: #fff;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  border: 1px solid #e2e8f0;
}

.settings-tabs :deep(.el-tabs__header) {
  margin-bottom: 24px;
}

.settings-tabs :deep(.el-tabs__item) {
  font-size: 14px;
  font-weight: 500;
}

.settings-tabs :deep(.el-tabs__item.is-active) {
  color: #6366f1;
}

.preview-bar {
  padding: 32px 24px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  min-height: 100px;
}

.form-tip {
  margin-top: 6px;
  font-size: 12px;
  color: #909399;
}

.form-tip a {
  color: #6366f1;
}

.seo-actions {
  display: flex;
  gap: 12px;
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #f1f5f9;
}

.sitemap-preview {
  padding: 12px;
  background: #f8fafc;
  border-radius: 8px;
}

.save-bar {
  position: fixed;
  bottom: 0;
  left: 240px;
  right: 0;
  padding: 16px 24px;
  background: #fff;
  border-top: 1px solid #e2e8f0;
  display: flex;
  justify-content: center;
  z-index: 100;
  transition: left 0.3s;
}

@media (max-width: 768px) {
  .save-bar {
    left: 0;
  }
}
</style>
