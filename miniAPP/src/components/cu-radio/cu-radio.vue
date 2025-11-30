<template>
  <view class="checkbox-container" :class="checkboxClassHandle" @click="changeCheckbox">
    <view class="circle" v-if="checkboxSelected"></view>
  </view>
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
  }
})

const emits = defineEmits(['update:modelValue', 'change'])

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
}
</script>

<style scoped lang="scss">
.checkbox-container {
  width: 40rpx;
  height: 40rpx;
  border-radius: 50%;
  background: #ffffff;
  border: 2.4rpx solid #dcdcdc;
  &.selected {
    border-color: var(--color-function-success);
  }
  .circle {
    margin: 8rpx 0 0 8rpx;
    width: 20rpx;
    height: 20rpx;
    border-radius: 50%;
    background: var(--color-function-success);
  }
}
</style>
