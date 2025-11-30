<template>
  <view class="custom-header">
    <view class="header-left" v-if="isBack">
      <slot name="left">
        <view @tap="BackPage" v-if="!isTopPage" class="center nav-slot">
          <uni-icons type="left" size="19"></uni-icons>
          <view class="divider"></view>
          <uni-icons type="home" size="20"></uni-icons>
        </view>
      </slot>
    </view>
    <view class="header-title">
      <slot name="title">
        <span class="title-text">{{ titleText }}</span>
      </slot>
    </view>
    <view class="header-right">
      <slot name="right"></slot>
    </view>
  </view>
</template>

<script setup>
import { ref, onMounted } from 'vue'

// 判断当前页面是否为最顶层页面
const isTopPage = ref(false)

onMounted(() => {
  checkIsTopPage()
})

const checkIsTopPage = () => {
  const pages = getCurrentPages()

  console.log('pages---', pages)
  isTopPage.value = pages.length <= 1
}

const props = defineProps({
  isBack: {
    // 是否显示左侧的返回键
    type: Boolean,
    default: true
  },
  tapLeft: {
    // 左上角的自定回调函数
    type: Function
  },
  backSuccess: {
    // 返回上一页成功的回调函数
    type: Function
  },
  titleText: {
    type: String,
    default: ''
  }
})

const BackPage = () => {
  if (props.tapLeft) {
    props.tapLeft()
    return
  }
  uni.navigateBack({
    delta: 1,
    success: function () {
      setTimeout(() => {
        if (props.backSuccess) uni.$emit('backSuccess')
      }, 100)
    }
  })
}
</script>

<style scoped lang="scss">
.custom-header {
  width: 750rpx;
  height: 88rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  background: #fff;
}

.header-left {
  position: absolute;
  left: 24rpx;
  //padding: 5rpx;
}

.header-right {
  position: absolute;
  right: 24rpx;
}

.nav-slot {
  border-radius: 100rpx;
  border: 1rpx solid #ddd;
  padding: 3rpx 7rpx;
  opacity: 0.8;
}

.title-text {
  font-weight: 600;
  font-size: 36rpx;
  color: var(--color-text-main);
  text-align: center;
}

.divider {
  width: 1rpx;
  height: 16rpx;
  background: #ddd;
  margin: 0 8rpx;
}
</style>
