import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  { path: '/', component: () => import('../views/Home.vue') },
  { path: '/articles', component: () => import('../views/Articles.vue') },
  { path: '/article/:id', component: () => import('../views/ArticleDetail.vue') },
  { path: '/contact', component: () => import('../views/Contact.vue') },
  { path: '/:pathMatch(.*)*', redirect: '/' },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  }
})

export default router
