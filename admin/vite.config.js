import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  base: '/admin',
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:9088',
        changeOrigin: true,
      }
    }
  },
  build: {
    outDir: 'dist',
  }
})
