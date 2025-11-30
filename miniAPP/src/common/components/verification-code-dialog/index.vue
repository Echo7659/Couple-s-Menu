<template>
  <!--    验证码弹窗-->
  <cf-middle-dialog ref="dialogRef">
    <view>
      <view class="pay-name color-text-main">{{ form.payName }}</view>

      <view class="font-bold text-center mt-[8rpx]">
        <text class="text-[48rpx]">￥</text>
        <text class="text-[64rpx]">{{ form.price }}</text>
      </view>

      <view class="code-form-box">
        <uni-easyinput
          type="number"
          v-model="form.code"
          :inputBorder="false"
          placeholder="请输入验证码"
          :placeholderStyle="placeholderStyle"
          primaryColor="var(--color-text-disable)"
        ></uni-easyinput>

        <view class="count-down">
          {{ smsCodeBtnText }}
        </view>
      </view>

      <view class="text-center color-text-secondary text-[28rpx] mt-[32rpx]">
        验证码已发送至: 132****8206
      </view>

      <view class="primary-button-box h-[92rpx] mt-[32rpx]" @click="confrimHanle">确认</view>
    </view>
  </cf-middle-dialog>
</template>

<script setup name="verificationCode">
import { useSmsCodeCountDown } from '@/hooks/useGlobalUtil'

const { smsCodeBtnText, isCountDown, startCountDownInterval } = useSmsCodeCountDown()

const dialogRef = ref()
const placeholderStyle = 'font-size:32rpx;color:var(--color-text-disable)'

const form = ref({
  payName: '钱包支付',
  price: '1299.00',
  code: ''
})

const confrimHanle = () => {
  uni.navigateTo({
    url: '/pages/payResult/payResult?category=success'
  })
}

const open = () => {
  dialogRef.value.open('请输入验证码')
  startCountDownInterval()
}

defineExpose({
  open
})
</script>

<style scoped lang="scss">
.pay-name {
  margin-top: 31rpx;
  text-align: center;
  font-size: 28rpx;
}
.code-form-box {
  height: 96rpx;
  border-top: 1rpx solid var(--color-division-line);
  border-bottom: 1rpx solid var(--color-division-line);
  margin-top: 64rpx;
  display: flex;
  align-items: center;
}

.count-down {
  width: 194rpx;
  height: 48rpx;
  font-size: 32rpx;
  color: var(--color-text-secondary);
  line-height: 48rpx;
  text-align: right;
  border-left: 1rpx solid var(--color-division-line);
}
</style>
