import { defineConfig } from 'vite'
import uni from '@dcloudio/vite-plugin-uni'
import AutoImport from 'unplugin-auto-import/vite'
import { fileURLToPath, URL } from 'node:url'
import UnoCSS from 'unocss/vite'


export default defineConfig({
  transpileDependencies: ['z-paging','@dcloudio/uni-ui'],
  plugins: [
    AutoImport({
      imports: ['vue'] // vue api 自动导入
    }),
    uni.default(),
    UnoCSS()
  ],
  css: {
    preprocessorOptions: {
      scss: {
        // 取消sass废弃API的报警
        silenceDeprecations: ['legacy-js-api', 'color-functions', 'import']
      }
    }
  },
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  build: {
    minify: 'terser',
    terserOptions: {
      compress: {
        // drop_console: true, // 去除console
        // drop_debugger: true // 去除debugger
      }
    }
  }
})
