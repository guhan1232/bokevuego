import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue'),
  },
  {
    path: '/',
    component: () => import('../views/Layout.vue'),
    redirect: '/dashboard',
    children: [
      { path: 'dashboard', name: 'Dashboard', component: () => import('../views/Dashboard.vue'), meta: { title: '仪表盘' } },
      { path: 'articles', name: 'Articles', component: () => import('../views/Articles.vue'), meta: { title: '文章管理' } },
      { path: 'articles/new', name: 'ArticleEdit', component: () => import('../views/ArticleEdit.vue'), meta: { title: '新建文章' } },
      { path: 'articles/:id', name: 'ArticleUpdate', component: () => import('../views/ArticleEdit.vue'), meta: { title: '编辑文章' } },
      { path: 'categories', name: 'Categories', component: () => import('../views/Categories.vue'), meta: { title: '分类管理' } },
      { path: 'messages', name: 'Messages', component: () => import('../views/Messages.vue'), meta: { title: '留言管理' } },
      { path: 'visits', name: 'Visits', component: () => import('../views/Visits.vue'), meta: { title: '访问记录' } },
      { path: 'friend-links', name: 'FriendLinks', component: () => import('../views/FriendLinks.vue'), meta: { title: '友情链接' } },
      { path: 'settings', name: 'Settings', component: () => import('../views/Settings.vue'), meta: { title: '站点设置' } },
    ],
  },
]

const router = createRouter({
  history: createWebHistory('/admin'),
  routes,
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (to.path !== '/login' && !token) {
    next('/login')
  } else if (to.path === '/login' && token) {
    next('/')
  } else {
    next()
  }
})

export default router
