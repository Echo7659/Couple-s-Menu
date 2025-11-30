<template>
  <uni-tabbar :value="tabBar.currentValue" :fixed="true" :placeholder="true">
    <uni-tabbar-item
      @click="toSwitchPage(item)"
      v-for="item in tabBar.tabBarList"
      :key="item.name"
      :name="item.name"
      :text="item.text"
      :icon="item.icon"
    ></uni-tabbar-item>
  </uni-tabbar>
</template>

<script setup>
import { useTabBarStore } from '@/store/tabbar.js'
import { onShow } from '@dcloudio/uni-app'
const tabBar = useTabBarStore()

onShow(() => {
  uni.hideTabBar({
    animation: false
  })
})

const toSwitchPage = item => {
  uni.switchTab({
    url: item.pagePath
  })
  tabBar.setValue(item.name)
}
</script>

<style scoped lang="scss"></style>
