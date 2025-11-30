<template>
  <view class="result-container">
    <view class="card">
      <view class="title">下单成功</view>
      <view class="desc">订单已创建，等待商家接单</view>
      <button class="btn" @click="goBack">返回点餐</button>
    </view>
  </view>
</template>

<script setup>
const goBack = () => {
  const pages = getCurrentPages()
  if (pages.length > 1) {
    uni.navigateBack({ delta: pages.length - 1 })
  } else {
    uni.switchTab({ url: '/pages/order/order' })
  }
}
</script>

<style lang="scss" scoped>
.result-container {
  min-height: 100vh; display:flex; align-items:center; justify-content:center; background:#f5f5f5;
}
.card { background:#fff; padding:40rpx; border-radius:20rpx; text-align:center; box-shadow:0 4rpx 20rpx rgba(0,0,0,.06); }
.title { font-size:36rpx; font-weight:700; color:#333; }
.desc { margin: 16rpx 0 24rpx; color:#666; }
.btn { background: linear-gradient(135deg, #FF6B9D 0%, #FF8E9E 100%); color:#fff; border:none; border-radius: 30rpx; padding: 16rpx 28rpx; font-size: 28rpx; }
</style> 