<template>
  <view>
    <view class="warning-info">
      <i
        class="iconfont icon-a-Lineartishi text-[32rpx] mr-[19rpx]"
        style="color: var(--color-function-warning)"
      ></i>

      <view class="text-[26rpx]" style="color: var(--color-function-warning)">
        请使用张清橙（*8888）的银行账户向以下账户自行转账，为确保订单正常发货，转账金额需保持一致！
      </view>
    </view>

    <view class="form-box">
      <view class="row-item" v-for="(item, index) in formItems" :key="index">
        <text class="item-label">{{ item.label }}</text>
        <view class="flex items-center">
          <text class="mr-[15rpx]" :style="item.style">{{ item.showText }}</text>
          <view class="copy-item" @click="copy(item.value)">复制</view>
        </view>
      </view>

      <view class="primary-button-box h-[96rpx] mt-[44rpx]" @click="confirmHandle">
        我已完成转账
      </view>
    </view>
  </view>
</template>

<script setup>
const formItems = ref([
  {
    label: '转账金额',
    showText: '￥120,999.00',
    style: 'fontWeight:bold',
    value: ''
  },
  {
    label: '收款公司名称',
    showText: '江楠鲜品科技有限公司',
    value: ''
  },
  {
    label: '收款账户',
    showText: '8888 8888 8888 8888',
    value: ''
  },
  {
    label: '开户行名称',
    showText: '招商银行广州白云支行',
    value: ''
  }
])

const copy = value => {
  uni.setClipboardData({
    data: value
  })
}

const emits = defineEmits(['success'])
const confirmHandle = () => {
  emits('success')
}
</script>

<style scoped lang="scss">
.form-box {
  padding: 0 24rpx;
}
.warning-info {
  height: 112rpx;
  background: #fff7f2;
  display: flex;
  align-items: center;
  padding: 0 28rpx;
}
.row-item {
  height: 96rpx;
  margin: 0 24rpx;
  border-bottom: 1rpx solid var(--color-division-line);
  display: flex;
  justify-content: space-between;
  align-items: center;
  .item-label {
    font-size: 32rpx;
    color: var(--color-text-main);
  }
}

.copy-item {
  width: 91rpx;
  text-align: center;
  border-left: 1rpx solid var(--color-division-line);
  font-size: 26rpx;
  color: var(--color-text-secondary);
}
</style>
