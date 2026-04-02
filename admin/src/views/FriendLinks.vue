<template>
  <div class="page-container">
    <div class="page-header">
      <h2>友情链接</h2>
      <el-button type="primary" @click="openDialog()">
        <el-icon><Plus /></el-icon>
        添加链接
      </el-button>
    </div>

    <el-card shadow="hover">
      <el-table :data="links" stripe style="width: 100%">
        <el-table-column prop="name" label="名称" width="180" />
        <el-table-column prop="url" label="链接" min-width="250">
          <template #default="{ row }">
            <el-link :href="row.url" target="_blank" type="primary">{{ row.url }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="logo" label="Logo" width="120">
          <template #default="{ row }">
            <el-image v-if="row.logo" :src="row.logo" style="width: 24px; height: 24px" fit="contain" />
            <span v-else style="color: #909399">无</span>
          </template>
        </el-table-column>
        <el-table-column prop="sort" label="排序" width="80" />
        <el-table-column prop="visible" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.visible ? 'success' : 'info'" size="small">
              {{ row.visible ? '显示' : '隐藏' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="140">
          <template #default="{ row }">
            <el-button type="primary" link @click="openDialog(row)">编辑</el-button>
            <el-popconfirm title="确定删除?" @confirm="handleDelete(row.id)">
              <template #reference>
                <el-button type="danger" link>删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 编辑弹窗 -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑链接' : '添加链接'" width="500px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="名称" required>
          <el-input v-model="form.name" placeholder="网站名称" />
        </el-form-item>
        <el-form-item label="链接" required>
          <el-input v-model="form.url" placeholder="https://example.com" />
        </el-form-item>
        <el-form-item label="Logo">
          <el-input v-model="form.logo" placeholder="Logo 图片 URL（可选）" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort" :min="0" :max="999" />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="form.visible" active-text="显示" inactive-text="隐藏" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import request from '../api/request'

const links = ref([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const form = ref({
  id: null,
  name: '',
  url: '',
  logo: '',
  sort: 0,
  visible: true
})

const fetchLinks = async () => {
  try {
    const res = await request.get('/friend-links')
    links.value = res || []
  } catch (e) {
    console.error('Failed to fetch links:', e)
  }
}

const openDialog = (row = null) => {
  if (row) {
    isEdit.value = true
    form.value = { ...row }
  } else {
    isEdit.value = false
    form.value = { id: null, name: '', url: '', logo: '', sort: 0, visible: true }
  }
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!form.value.name || !form.value.url) {
    ElMessage.warning('请填写名称和链接')
    return
  }
  try {
    if (isEdit.value) {
      await request.put(`/friend-links/${form.value.id}`, form.value)
      ElMessage.success('更新成功')
    } else {
      await request.post('/friend-links', form.value)
      ElMessage.success('添加成功')
    }
    dialogVisible.value = false
    fetchLinks()
  } catch (e) {
    ElMessage.error('操作失败')
  }
}

const handleDelete = async (id) => {
  try {
    await request.delete(`/friend-links/${id}`)
    ElMessage.success('删除成功')
    fetchLinks()
  } catch (e) {
    ElMessage.error('删除失败')
  }
}

onMounted(() => {
  fetchLinks()
})
</script>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}
</style>
