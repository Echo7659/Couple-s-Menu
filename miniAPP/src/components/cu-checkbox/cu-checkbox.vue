<template>
  <view class="checkbox-container" :class="checkboxClassHandle" @tap.stop="changeCheckbox">
    <i class="iconfont icon-a-Lineargou text-[#fff]" v-if="checkboxSelected"></i>
  </view>

  <!--  需要放大的热区,跟父元素一样的大小-->
  <div class="scale-box" v-if="showScale" @tap.stop="changeCheckbox"></div>
</template>

<script setup>
const props = defineProps({
  modelValue: {
    // v-model的值,一般用于只有一个复选框时
    type: Boolean
  },
  selected: {
    // 当前组件是否选中状态
    type: Boolean
  },
  disabled: {
    type: Boolean,
    default: false
  },
  showScale: {
    type: Boolean,
    default: false
  }
})

const emits = defineEmits(['update:modelValue', 'change', 'customClick'])

const checkboxSelected = computed({
  get() {
    return props.modelValue || props.selected
  },
  set(value) {
    emits('update:modelValue', value)
    emits('change', value)
  }
})

const checkboxClassHandle = computed(() => {
  let classNames = props.disabled ? 'disabled ' : ''
  classNames += checkboxSelected.value ? 'selected' : ''
  return classNames
})

const changeCheckbox = () => {
  if (props.disabled) return
  checkboxSelected.value = !checkboxSelected.value

  emits('customClick', checkboxSelected.value)
}
</script>

<style scoped lang="scss">
.checkbox-container {
  width: 40rpx;
  height: 40rpx;
  border-radius: 50%;
  background: #ffffff;
  border: 2.4rpx solid #dcdcdc;
  display: flex;
  justify-content: center;
  align-items: center;
  position: relative;
  &.selected {
    background: #2e69ff;
    border: none;
  }

  &.disabled {
    background: rgba(0, 0, 0, 0.1);
  }
}

.scale-box {
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0;
  top: 0;
}
</style>
