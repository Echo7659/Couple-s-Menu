<template>
  <root-portal>
    <uni-popup
      ref="popup"
      type="bottom"
      :mask-click="false"
    >
      <view class="popup-content" :style="customStyle">
        <view class="dialog-header">
          <text>{{ title }}</text>
          <i class="iconfont icon-a-Linearguanbi text-[36rpx] close-icon" @tap="close"></i>
        </view>
        <slot></slot>
      </view>
    </uni-popup>
  </root-portal>
</template>

<script setup>
import { ref } from 'vue'

defineProps({
  customStyle: {
    type: Object,
    default: () => {
      return {}
    }
  }
})

const title = ref('提示')
const popup = ref(null)

const open = titleStr => {
  title.value = titleStr
  popup.value.open()
}

const emits = defineEmits(['close'])
const close = () => {
  popup.value.close()
  emits('close')
}

defineExpose({
  open,
  close
})
</script>
<style scoped lang="scss">
.dialog-header {
  text-align: center;
  height: 112rpx;
  line-height: 112rpx;
  color: var(--color-text-main);
  font-size: 32rpx;
  position: relative;
  font-weight: 600;
}

.close-icon {
  position: absolute;
  right: 34rpx;
  top: 0rpx;
  color: var(--color-text-disable);
}
.popup-content {
  background: #fff;
  border-radius: 24rpx 24rpx 0 0;
}
</style>
