import { defineStore } from 'pinia'

export const useUserStore = defineStore({
  id: 'user',
  state: () => {
    return {
      token: uni.getStorageSync('token') || '', // 从缓存读取token
      userInfo: {}, // 用户信息
      appInfo: {}, // 应用信息
      // 按钮权限编码
      permissionCodes: [],
      // 权限菜单
      menuList: [],
      // 权限路由
      permissionRoutes: [],
      // 是否已登录
      isLoggedIn: !!uni.getStorageSync('token')
    }
  },
  unistorage: true,
  getters: {
    // 获取用户昵称
    nickname: (state) => state.userInfo.nickname || '未登录',
    
    // 获取用户头像
    avatar: (state) => state.userInfo.avatar || '/static/images/default-avatar.png',
    
    // 获取用户ID
    userId: (state) => state.userInfo.userId || '',
    
    // 检查是否已登录
    hasLogin: (state) => state.isLoggedIn && !!state.token,
    
    // 检查是否首次登录
    isFirstLogin: (state) => state.userInfo.isFirstLogin || false
  },
  actions: {
    // 设置用户信息
    setUserInfo(userInfo) {
      this.userInfo = userInfo
      this.isLoggedIn = true
    },
    
    // 设置token
    setToken(token) {
      this.token = token
      // 同时保存到本地缓存
      uni.setStorageSync('token', token)
    },
    
    // 登录
    async login(loginData) {
      try {
        // 这里可以调用登录API
        // const result = await userApi.login(loginData)
        
        // 模拟登录成功
        const mockUserInfo = {
          userId: '10086',
          nickname: '小美',
          avatar: '/static/images/default-avatar.png',
          level: 'LV.5 美食达人'
        }
        
        this.setUserInfo(mockUserInfo)
        this.setToken('mock_token_' + Date.now())
        
        return {
          success: true,
          data: mockUserInfo
        }
      } catch (error) {
        console.error('登录失败:', error)
        return {
          success: false,
          message: error.message || '登录失败'
        }
      }
    },
    
    // 登出
    logout() {
      this.token = ''
      this.userInfo = {}
      this.isLoggedIn = false
      this.permissionCodes = []
      this.menuList = []
      this.permissionRoutes = []
      // 清除本地缓存
      uni.removeStorageSync('token')
    },
    
    // 更新用户信息
    updateUserInfo(userInfo) {
      this.userInfo = { ...this.userInfo, ...userInfo }
      // 如果更新了昵称，清除首次登录状态
      if (userInfo.nickname) {
        this.userInfo.isFirstLogin = false
      }
    },
    
    // 检查登录状态
    checkLoginStatus() {
      // 从缓存重新读取token
      const cachedToken = uni.getStorageSync('token')
      if (cachedToken && cachedToken !== this.token) {
        this.token = cachedToken
        this.isLoggedIn = true
      }
      return this.hasLogin
    },
    
    // 清除首次登录状态
    clearFirstLogin() {
      this.userInfo.isFirstLogin = false
    }
  }
})
