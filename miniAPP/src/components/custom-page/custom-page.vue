<template>
  <view :class="{ iphone_mb: !isTabPage }" class="page-view" :style="statusBarHeight()">
    <view class="page-container relative" :class="systemStore.commonThemeName" :style="pageStyle()">
      <slot></slot>
    </view>

    <!--    底部安全区域的颜色,解决IOS下导航栏与底部安全要显示不同颜色的问题-->
    <div
      v-if="showBottomSafeBox"
      class="bottom-safe-box"
      :style="{ backgroundColor: safeBottomBgColor }"
    ></div>

    <view class="loading-box" v-if="systemStore.showLoading">
      <image class="h-[120rpx]" mode="aspectFit" :src="systemStore.shopInfo.global_logo"></image>
      <div class="color-text-secondary text-[26rpx] mt-[18rpx]">加载中</div>
    </view>
  </view>
</template>

<script setup>
import { useSystemStore } from '@/store/system.js'
const systemStore = useSystemStore()

const props = defineProps({
  customStyle: {
    type: Object,
    default() {
      return {}
    }
  },
  customStatusStyle: {
    type: Object,
    default() {
      return {}
    }
  },
  isFullPage: {
    //   是否全屏,包括手机的状态栏
    type: Boolean,
    default: false
  },
  isShoeHead: {
    //   是否显示头部
    type: Boolean,
    default: true
  },
  isTabPage: {
    //   是否是tab页面,tab页面时,会默认给全面屏加一个底部的安全区域,就不需要手动添加,
    type: Boolean,
    default: false
  },
  // IOS底部安全区域的背景色
  safeBottomBgColor: {
    type: String,
    default: '#fff'
  },
  // 是否需要添加底部安全区域
  showBottomSafeBox: {
    type: Boolean,
    default: false
  }
})

const statusBarHeight = () => {
  if (!props.isFullPage) {
    // 非全屏时,增加状态栏高度
    return {
      paddingTop: `${systemStore.statusBarHeight}px`,
      ...props.customStatusStyle
    }
  }
}

const pageStyle = () => {
  // 页面内容的样式
  return {
    ...props.customStyle
  }
}

onBeforeMount(() => {
  systemStore.setShowLoading(false)
})

onUnmounted(() => {
  systemStore.setShowLoading(false)
})
</script>

<style scoped lang="scss">
.iphone_mb {
  padding-bottom: env(safe-area-inset-bottom);
}

.page-view {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  box-sizing: border-box;
}

.loading-box {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  position: absolute;
  flex-direction: column;
  top: 0;
  z-index: 9999999;
}

.bottom-safe-box {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: #fff;
  z-index: 2;
  height: env(safe-area-inset-bottom);
}
</style>
