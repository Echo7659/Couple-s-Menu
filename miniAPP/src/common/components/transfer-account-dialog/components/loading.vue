<template>
  <view class="flex flex-col items-center">
    <image class="w-[133rpx] h-[133rpx] mt-[93rpx]" src="@/static/images/loading.gif"></image>

    <text class="color-text-main font-bold text-[36rpx] mt-[40rpx]">正在查询支付结果…</text>
  </view>
</template>

<script setup>
const emit = defineEmits(['result'])

setTimeout(() => {
  emit('result', 4)
}, 2000)
</script>

<style scoped lang="scss"></style>
