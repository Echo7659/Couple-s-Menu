import { defineStore } from 'pinia'

export const useSystemStore = defineStore({
  id: 'system',
  state: () => {
    return {
      statusBarHeight: 0, // 状态栏的高度
      phoneInfo: {},
      showLoading: false, // 全局的loading
      commonThemeName: 'default-theme' // 主题颜色
    }
  },
  actions: {
    setStatusBarHeight(height) {
      this.statusBarHeight = height
    },
    setThemeName(name) {
      this.commonThemeName = name
    },

    setShowLoading(value) {
      this.showLoading = value
    },
    init() {
      uni.getSystemInfo({
        success: e => {
          this.phoneInfo = e
          this.statusBarHeight = e.statusBarHeight
        }
      })
    }
  }
})
