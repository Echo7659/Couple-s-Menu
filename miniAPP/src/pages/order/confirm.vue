<template>
  <view class="confirm-container">
    <view class="card">
      <view class="card-title">已选菜品</view>
      <view v-for="item in items" :key="item.id" class="row">
        <image :src="item.img" class="thumb" mode="aspectFill" />
        <view class="meta">
          <text class="name">{{ item.name }}</text>
          <text class="price">¥{{ item.price || 0 }}</text>
        </view>
        <text class="qty">×{{ item.qty }}</text>
      </view>
      <view class="total">合计：<text class="price">¥{{ totalPrice }}</text></view>
    </view>

    <view class="card">
      <view class="card-title">留言</view>
      <uni-easyinput v-model="remark" placeholder="请输入备注（可选）"></uni-easyinput>
    </view>

    <view class="submit-bar">
      <text class="sum">合计：<text class="price">¥{{ totalPrice }}</text></text>
      <button class="submit" :loading="submitting" @click="create">提交订单</button>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import api from '@/api'

const items = ref([])
const ids = ref('')
const userId = ref(null)
const remark = ref('')
const submitting = ref(false)

onMounted(() => {
  const channel = getCurrentPages().pop()?.getOpenerEventChannel?.()
  channel?.on('initData', (payload) => {
    items.value = payload.items || []
    ids.value = payload.ids || ''
    userId.value = payload.userId || null
  })
})

const totalPrice = computed(() => items.value.reduce((s, it) => s + (it.price || 0) * it.qty, 0))

const create = async () => {
  if (!ids.value) {
    uni.showToast({ title: '未选择菜品', icon: 'none' })
    return
  }
  
  // 如果用户ID为空，尝试从缓存获取
  if (!userId.value) {
    try {
      const cached = uni.getStorageSync('userInfo')
      if (cached?.ID || cached?.id) {
        userId.value = cached.ID || cached.id
      } else {
        const res = await api.getUserInfo()
        if (res?.code === 0 && res.data) {
          userId.value = res.data.ID || res.data.id || null
        }
      }
    } catch (e) {
      console.warn('获取用户ID失败', e)
    }
  }
  
  // 如果还是没有用户ID，提示登录
  if (!userId.value) {
    uni.showToast({
      title: '请先登录',
      icon: 'none'
    })
    setTimeout(() => {
      uni.navigateTo({ url: '/pages/user/login' })
    }, 1500)
    return
  }
  
  try {
    submitting.value = true
    const payload = {
      food_ids: ids.value,
      user_id: userId.value,
      remark: remark.value
    }
    const res = await api.createOrder(payload)
    if (res?.code === 0) {
      uni.showToast({ title: '下单成功', icon: 'success' })
      setTimeout(() => {
        uni.redirectTo({ url: `/pages/order/result?id=${res.data?.ID || ''}` })
      }, 600)
    } else {
      uni.showToast({ title: res?.msg || '下单失败', icon: 'none' })
    }
  } catch (e) {
    console.error('create order failed', e)
    uni.showToast({ title: '下单失败', icon: 'none' })
  } finally {
    submitting.value = false
  }
}
</script>

<style lang="scss" scoped>
.confirm-container {
  min-height: 100vh;
  background: #f5f5f5;
  padding-bottom: 140rpx;
}

.card {
  background: #fff;
  margin: 20rpx;
  padding: 20rpx;
  border-radius: 16rpx;

  .card-title {
    font-size: 28rpx;
    font-weight: 600;
    margin-bottom: 16rpx;
  }

  .row {
    display: flex;
    align-items: center;
    padding: 16rpx 0;
    border-bottom: 1rpx solid #f5f5f5;

    &:last-child { border-bottom: none; }

    .thumb { width: 120rpx; height: 120rpx; border-radius: 12rpx; margin-right: 16rpx; }
    .meta { flex: 1; .name { display:block; font-size: 28rpx; } .price { color:#FF6B9D; margin-top: 6rpx; display:block; } }
    .qty { color: #666; }
  }

  .total { text-align: right; margin-top: 10rpx; .price{ color:#FF6B9D; font-weight: 700; } }
}

.submit-bar {
  position: fixed;
  left: 0; right: 0; bottom: 0;
  background: #fff;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20rpx;
  box-shadow: 0 -4rpx 20rpx rgba(0,0,0,0.06);

  .sum { font-size: 28rpx; .price{ color:#FF6B9D; font-weight:700; }}
  .submit { 
    background: linear-gradient(135deg, #FF6B9D 0%, #FF8E9E 100%); 
    color:#fff; 
    border:none; 
    border-radius: 30rpx; 
    padding: 16rpx 28rpx; 
    font-size: 28rpx; 
  }
}
</style> 