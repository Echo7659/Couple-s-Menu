<template>
  <!--  外层的view flex-1之外,还需要加上overflow-hidden,否则在小程序中scroll-view无法继承父盒子的高度-->
  <view class="flex-1 h-full overflow-hidden">
    <scroll-view
      class="h-full w-full overflow-auto"
      :show-scrollbar="false"
      :scroll-y="true"
      :enhanced="true"
      :scroll-with-animation="true"
    >
      <slot></slot>
    </scroll-view>
  </view>
</template>
<!-- 自动获取剩余高度的滚动组件,去除了的IOS的滚动条 -->
<script setup>
defineProps({
  // 外层view的自定义样式
  customStyle: {
    type: Object,
    default: () => {
      return {}
    }
  }
})
</script>

<style scoped lang="scss"></style>
