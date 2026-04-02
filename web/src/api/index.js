import request from './request'

export const getArticles = (params) => request.get('/articles', { params })
export const getArticle = (id) => request.get(`/articles/${id}`)
export const viewArticle = (id) => request.post(`/articles/${id}/view`)
export const getCategories = () => request.get('/categories')
export const getTags = () => request.get('/tags')
export const getConfigs = () => request.get('/configs')
export const createMessage = (data) => request.post('/messages', data)
export const getFriendLinks = () => request.get('/friend-links')
