<template>
  <div class="page-container">
    <div class="page-header">
      <h2>{{ isEdit ? '编辑文章' : '新建文章' }}</h2>
      <div>
        <el-button @click="$router.back()">取消</el-button>
        <el-button type="info" @click="handleSave(0)" :loading="saving">存为草稿</el-button>
        <el-button type="primary" @click="handleSave(1)" :loading="saving">发布</el-button>
      </div>
    </div>

    <el-card shadow="hover">
      <el-form :model="form" label-width="80px" label-position="top">
        <el-form-item label="标题">
          <el-input v-model="form.title" placeholder="请输入文章标题" size="large" />
        </el-form-item>

        <el-row :gutter="20">
          <el-col :xs="24" :md="12">
            <el-form-item label="分类">
              <el-input v-model="form.category" placeholder="输入分类" />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :md="12">
            <el-form-item label="封面图">
              <el-input v-model="form.cover" placeholder="封面图片URL" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="摘要">
          <el-input v-model="form.summary" type="textarea" :rows="3" placeholder="文章摘要（可选）" />
        </el-form-item>

        <el-form-item label="标签">
          <div class="tags-input">
            <el-tag
              v-for="tag in form.tags"
              :key="tag"
              closable
              @close="removeTag(tag)"
              style="margin-right: 8px; margin-bottom: 4px"
            >
              {{ tag }}
            </el-tag>
            <el-input
              v-if="tagInputVisible"
              ref="tagInputRef"
              v-model="tagInputValue"
              size="small"
              style="width: 100px"
              @keyup.enter="addTag"
              @blur="addTag"
            />
            <el-button v-else size="small" @click="tagInputVisible = true">+ 添加标签</el-button>
          </div>
        </el-form-item>

        <el-form-item label="内容">
          <div class="editor-wrap">
            <el-input
              v-model="form.content"
              type="textarea"
              :rows="18"
              placeholder="支持 Markdown 格式内容"
              resize="vertical"
            />
          </div>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getArticle, createArticle, updateArticle } from '../api'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()
const isEdit = computed(() => !!route.params.id)
const saving = ref(false)

import { computed } from 'vue'

const form = reactive({
  title: '',
  summary: '',
  content: '',
  cover: '',
  category: '',
  tags: [],
})

const tagInputVisible = ref(false)
const tagInputValue = ref('')
const tagInputRef = ref()

const addTag = () => {
  const v = tagInputValue.value.trim()
  if (v && !form.tags.includes(v)) {
    form.tags.push(v)
  }
  tagInputVisible.value = false
  tagInputValue.value = ''
}

const removeTag = (tag) => {
  form.tags = form.tags.filter(t => t !== tag)
}

const handleSave = async (status) => {
  if (!form.title.trim()) {
    ElMessage.warning('请输入标题')
    return
  }
  saving.value = true
  try {
    const data = { ...form, status }
    if (isEdit.value) {
      await updateArticle(route.params.id, data)
      ElMessage.success('更新成功')
    } else {
      await createArticle(data)
      ElMessage.success('创建成功')
    }
    router.push('/articles')
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  if (route.params.id) {
    try {
      const data = await getArticle(route.params.id)
      Object.assign(form, data.article)
      form.tags = data.tags || []
    } catch {}
  }
})
</script>

<style scoped>
.tags-input {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
}

.editor-wrap {
  width: 100%;
}

.editor-wrap :deep(.el-textarea__inner) {
  font-family: 'Menlo', 'Monaco', 'Courier New', monospace;
  line-height: 1.6;
}
</style>
