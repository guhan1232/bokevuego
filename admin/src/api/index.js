import request from './request'

// 认证
export const login = (data) => request.post('/login', data)
export const getProfile = () => request.get('/profile')

// 文章
export const getArticles = (params) => request.get('/articles', { params })
export const getArticle = (id) => request.get(`/articles/${id}`)
export const createArticle = (data) => request.post('/articles', data)
export const updateArticle = (id, data) => request.put(`/articles/${id}`, data)
export const deleteArticle = (id) => request.delete(`/articles/${id}`)
export const getTrashArticles = (params) => request.get('/articles/trash', { params })
export const restoreArticle = (id) => request.post(`/articles/${id}/restore`)
export const hardDeleteArticle = (id) => request.delete(`/articles/${id}/hard`)

// 分类/标签
export const getCategories = () => request.get('/categories')
export const getTags = () => request.get('/tags')

// 站点配置
export const getConfigs = () => request.get('/configs')
export const updateConfigs = (data) => request.put('/configs', data)

// 留言
export const getMessages = (params) => request.get('/messages', { params })
export const replyMessage = (data) => request.post('/messages/reply', data)
export const deleteMessage = (id) => request.delete(`/messages/${id}`)

// 统计
export const getStats = () => request.get('/stats')
