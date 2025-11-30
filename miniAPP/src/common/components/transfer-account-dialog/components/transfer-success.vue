<template>
  <view class="flex flex-col items-center justify-between pl-[24rpx] pr-[24rpx] w-full h-full">
    <view class="flex flex-col items-center">
      <i
        class="text-[133rpx] mt-[93rpx] iconfont icon-a-Fillchenggong"
        style="color: var(--color-primary); line-height: 1"
      ></i>

      <text class="color-text-main font-bold text-[36rpx] mt-[40rpx]">转账成功</text>
    </view>

    <view class="primary-button-box h-[96rpx] mt-[44rpx] w-full mb-[24rpx]" @click="success">
      完成
    </view>
  </view>
</template>

<script setup>
const success = () => {}
</script>

<style scoped lang="scss"></style>
