import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')
  const apiTarget = env.VITE_PROXY_TARGET || 'http://[::1]:8080'
  const wsTarget = env.VITE_WS_PROXY_TARGET || apiTarget.replace(/^http/i, 'ws')

  return {
    plugins: [vue()],
    server: {
      port: 5173,
      proxy: {
        '/api': {
          target: apiTarget,
          changeOrigin: true
        },
        '/ws': {
          target: wsTarget,
          ws: true,
          changeOrigin: true
        }
      }
    }
  }
})
