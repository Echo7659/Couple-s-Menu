<template>
  <view class="flex flex-col items-center justify-between pl-[24rpx] pr-[24rpx] w-full h-full">
    <view class="flex flex-col items-center">
      <i
        class="text-[133rpx] mt-[93rpx] iconfont icon-a-Fillcuowu"
        style="color: var(--color-function-error); line-height: 1"
      ></i>

      <text class="color-text-main font-bold text-[36rpx] mt-[40rpx]">暂未收到转账信息</text>

      <text class="fail-hint">
        我们正在处理您的转账，由于系统延迟，目前尚未确认确认到来自张清橙
        (*8888)的转账。我们会持续更新进度，请稍后查看。感谢您的理解与支持！
      </text>
    </view>

    <view class="primary-button-box h-[96rpx] mt-[44rpx] mb-[24rpx] w-full" @click="confirmHandle">
      重新转账
    </view>
  </view>
</template>

<script setup>
const emits = defineEmits(['result'])
const confirmHandle = () => {
  emits('result', 1)
}
</script>

<style scoped lang="scss">
.fail-hint {
  font-size: 28rpx;
  color: var(--color-text-secondary);
  margin-top: 8rpx;
}
</style>
