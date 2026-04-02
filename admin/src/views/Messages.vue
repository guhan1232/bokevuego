<template>
  <div class="page-container">
    <div class="page-header">
      <h2>留言管理</h2>
    </div>

    <el-table :data="messages" stripe v-loading="loading" style="width: 100%">
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="name" label="昵称" width="120" />
      <el-table-column prop="email" label="邮箱" width="180" />
      <el-table-column prop="content" label="内容" min-width="200" show-overflow-tooltip />
      <el-table-column label="状态" width="90">
        <template #default="{ row }">
          <el-tag :type="row.reply ? 'success' : 'warning'" size="small">
            {{ row.reply ? '已回复' : '待回复' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="时间" width="170">
        <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="160" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link size="small" @click="openReply(row)">回复</el-button>
          <el-popconfirm title="确定删除？" @confirm="handleDelete(row.id)">
            <template #reference>
              <el-button type="danger" link size="small">删除</el-button>
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination-wrap">
      <el-pagination
        v-model:current-page="page"
        :page-size="10"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="fetchList"
      />
    </div>

    <el-dialog v-model="replyVisible" title="回复留言" width="480px">
      <el-alert v-if="replyTarget" :title="replyTarget.content" type="info" :closable="false" style="margin-bottom: 16px" />
      <el-input v-model="replyContent" type="textarea" :rows="4" placeholder="输入回复内容" />
      <template #footer>
        <el-button @click="replyVisible = false">取消</el-button>
        <el-button type="primary" @click="handleReply">确定回复</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getMessages, replyMessage, deleteMessage } from '../api'
import { ElMessage } from 'element-plus'

const messages = ref([])
const total = ref(0)
const page = ref(1)
const loading = ref(false)
const replyVisible = ref(false)
const replyTarget = ref(null)
const replyContent = ref('')

const fetchList = async () => {
  loading.value = true
  try {
    const data = await getMessages({ page: page.value })
    messages.value = data.items || []
    total.value = data.total || 0
  } finally {
    loading.value = false
  }
}

const openReply = (row) => {
  replyTarget.value = row
  replyContent.value = row.reply || ''
  replyVisible.value = true
}

const handleReply = async () => {
  if (!replyContent.value.trim()) {
    ElMessage.warning('请输入回复内容')
    return
  }
  await replyMessage({ id: replyTarget.value.id, reply: replyContent.value })
  ElMessage.success('回复成功')
  replyVisible.value = false
  fetchList()
}

const handleDelete = async (id) => {
  await deleteMessage(id)
  ElMessage.success('删除成功')
  fetchList()
}

const formatDate = (d) => d ? new Date(d).toLocaleString('zh-CN') : '-'

onMounted(fetchList)
</script>
