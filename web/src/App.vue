<template>
  <div class="app" :style="bgStyle">
    <router-view />
  </div>
</template>

<script setup>
import { ref, provide, onMounted, computed } from 'vue'
import { getConfigs } from './api'

const siteConfig = ref({})
provide('siteConfig', siteConfig)

const bgStyle = computed(() => {
  const bgImage = siteConfig.value.site_bg_image
  if (!bgImage) return {}
  return {
    backgroundImage: `url(${bgImage})`,
    backgroundSize: 'cover',
    backgroundAttachment: 'fixed',
    backgroundPosition: 'center',
    backgroundRepeat: 'no-repeat',
  }
})

onMounted(async () => {
  try {
    siteConfig.value = await getConfigs()
    document.title = siteConfig.value.site_title || 'BokeUI 博客'
    
    // 设置 favicon
    if (siteConfig.value.site_favicon) {
      let link = document.querySelector("link[rel~='icon']")
      if (!link) {
        link = document.createElement('link')
        link.rel = 'icon'
        document.getElementsByTagName('head')[0].appendChild(link)
      }
      link.href = siteConfig.value.site_favicon
    }
  } catch {}
})
</script>

<style>
/* 全局样式已移至 assets/style.css */
</style>
