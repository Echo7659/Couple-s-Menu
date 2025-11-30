<template>
  <uni-popup ref="popupRef" v-bind="$attrs" mask-background-color="var(--color-bg-mask)">
    <view class="popup-content" :style="{ width: `${props.width}rpx` }">
      <view class="dialog-header">
        <text>{{ title }}</text>
        <i class="iconfont icon-a-Linearguanbi text-[36rpx] close-icon" @tap="close"></i>
      </view>
      <slot></slot>
    </view>
  </uni-popup>
</template>

<script setup>
const props = defineProps({
  width: {
    type: Number,
    default: 630
  }
})

const title = ref('提示')

const popupRef = ref()
const open = titleStr => {
  title.value = titleStr
  popupRef.value.open()
}
const close = () => {
  popupRef.value.close()
}

defineExpose({
  open,
  close
})
</script>
<style scoped lang="scss">
.dialog-header {
  text-align: center;
  height: 119rpx;
  line-height: 119rpx;
  color: var(--color-text-main);
  font-size: 36rpx;
  position: relative;
  font-weight: 500;
  border-bottom: 1rpx solid var(--color-division-line);
}

.close-icon {
  position: absolute;
  right: 0;
  top: 0rpx;
  color: var(--color-text-disable);
}
.popup-content {
  background: #fff;
  border-radius: 20rpx;
  padding: 0 48rpx 48rpx 48rpx;
}
</style>
