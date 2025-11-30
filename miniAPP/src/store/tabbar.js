import { defineStore } from 'pinia'

export const useTabBarStore = defineStore({
  id: 'tabBar',
  state: () => {
    return {
      // 底部tab菜单
      tabBarList: [
        {
          pagePath: '/pages/home/home',
          name: 'home',
          text: '首页',
          icon: 'home',
          selectIcon: ''
        },
        {
          name: 'components',
          text: '组件',
          icon: 'list',
          selectIcon: '',
          pagePath: '/pages/components/components'
        },
        {
          name: 'functions',
          text: '功能',
          icon: 'grid',
          selectIcon: '',
          pagePath: '/pages/functions/functions'
        }
      ],
      currentValue: 'home'
    }
  },
  actions: {
    setValue(name) {
      this.currentValue = name
    }
  }
})
