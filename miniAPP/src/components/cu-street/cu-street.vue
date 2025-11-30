<template>
  <view class="pupop">
    <view class="popup-box" :animation="animationData">
      <view class="pupop-btn">
        <view @tap="close">取消</view>
        <view @tap="confirm" style="color: var(--color-function-success)">确定</view>
      </view>
      <picker-view
        :value="value"
        @change="bindChange"
        :indicator-style="indicatorStyle"
        class="picker-view"
      >
        <picker-view-column>
          <view class="item" v-for="(item, index) in options" :key="index">
            {{ item.n }}
          </view>
        </picker-view-column>
      </picker-view>
    </view>

    <view
      @tap="close"
      @touchmove.stop.prevent
      :class="visible ? 'pupop-model' : 'pupop-models'"
    ></view>
  </view>
</template>

<script setup>
// 街道选择器,用于省市区跟街道分开的选择器
const props = defineProps({
  areaCode: {
    type: String,
    require: true
  },
  streetData: {
    type: Array,
    default() {
      return []
    }
  }
})

const value = ref([])

const indicatorStyle = `height: 50px;`
const emits = defineEmits(['confirm', 'cancel'])

const options = computed(() => {
  return props.streetData.filter(e => e.areaCode === props.areaCode)
})

const selectedStreet = ref(options.value[0])

const bindChange = e => {
  selectedStreet.value = options.value[e.detail.value[0]]
}

const visible = ref(false)

const animationData = ref('')

const confirm = () => {
  visible.value = false
  changeActive()
  emits('confirm', selectedStreet.value)
  emitCancel()
}
const changeActive = () => {
  let active = '-315px'
  if (visible.value) {
    active = 0
  }
  let animation = uni.createAnimation({
    duration: 400,
    timingFunction: 'linear'
  })
  animation.bottom(active).step()
  animationData.value = animation.export()
}

const emitCancel = () => {
  setTimeout(() => {
    emits('cancel')
  }, 400)
}
// 点击模态框
const close = () => {
  visible.value = false
  changeActive()
  emitCancel()
}

onMounted(() => {
  visible.value = true
  changeActive()
})
</script>

<style scoped lang="scss">
.pupop {
  .popup-box {
    position: fixed;
    left: 0;
    bottom: -315px;
    z-index: 99999;
    background: #fff;
    padding-bottom: 50px;
    .pupop-btn {
      height: 40px;
      display: flex;
      align-items: center;
      padding: 0 30upx;
      justify-content: space-between;
    }
  }
  .pupop-model {
    position: fixed;
    left: 0;
    top: 0;
    width: 100%;
    height: 100%;
    z-index: 9999;
    background: rgba(0, 0, 0, 0.5);
  }
  .pupop-models {
    display: none;
  }
  .picker-view {
    width: 750rpx;
    height: 225px;
    margin-top: 20rpx;
  }
  .item {
    height: 50px;
    align-items: center;
    justify-content: center;
    text-align: center;
    line-height: 50px;
  }
}
</style>
