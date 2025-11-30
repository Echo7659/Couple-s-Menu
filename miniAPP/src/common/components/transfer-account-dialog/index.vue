<template>
  <!--    转账弹窗-->
  <cf-bottom-dialog ref="dialogRef" type="bottom">
    <view>
      <view class="content-box">
        <!--        转账的信息内容-->
        <transferContent v-if="form.state == 1" @success="changeState(2)"></transferContent>
        <!--        查询-->
        <loading v-if="form.state == 2" @result="changeState"></loading>
        <!--        转账成功-->
        <transferSuccess v-if="form.state == 3"></transferSuccess>
        <!--        转账失败-->
        <transferFail v-if="form.state == 4" @result="changeState"></transferFail>
      </view>
    </view>
  </cf-bottom-dialog>
</template>

<script setup name="TransferAccount">
import transferContent from './components/transfer-content.vue'
import loading from './components/loading.vue'
import transferSuccess from './components/transfer-success.vue'
import transferFail from './components/transfer-fail.vue'
const dialogRef = ref()

const form = ref({
  state: 1
})

const changeState = state => {
  dialogRef.value.open('转账结果')
  form.value.state = state
}

const confrimHanle = () => {
  uni.navigateTo({
    url: '/pages/payResult/payResult?category=success'
  })
}

const open = () => {
  dialogRef.value.open('请转账到')
}

defineExpose({
  open
})
</script>

<style scoped lang="scss">
.content-box {
  height: 680rpx;
}
</style>
