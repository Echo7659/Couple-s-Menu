<template>
  <view class="orders-container">
    <view class="summary">
      <text>一共记录了{{ total }}个订单</text>
    </view>

    <view class="calendar-box">
      <uni-calendar
        ref="calendarRef"
        :insert="true"
        :lunar="false"
        :start-date="'2000-01-01'"
        :end-date="'2100-12-31'"
        :date="selectedDate"
        @change="onDateChange"
        @monthSwitch="onMonthSwitch"
      />
    </view>

    <view class="orders-list">
      <view v-for="order in orders" :key="order.id" class="order-card">
        <view class="card-header">
          <text class="order-date">{{ formatDate(order.date || selectedDate) }} {{ order.time || '' }}</text>
          <view class="card-actions">
            <view class="act-btn" @click.stop="onEvaluate(order)">
              <uni-icons type="compose" size="18" color="#999" />
              <text class="act-text">评价</text>
            </view>
            <view class="act-btn" @click.stop="onDelete(order)">
              <uni-icons type="trash" size="18" color="#999" />
              <text class="act-text">删除</text>
            </view>
          </view>
        </view>

        <view v-if="isExpanded(order.id)" class="list">
          <view v-for="(food, idx) in order.Foods || []" :key="food.ID || idx" class="list-item">
            <image :src="food.img" mode="aspectFill" class="list-thumb" />
            <view class="list-info">
              <text class="name">{{ food.name }}</text>
              <text class="qty">x1</text>
            </view>
          </view>
        </view>
        <view v-else class="grid">
          <view v-for="(food, idx) in order.Foods || []" :key="food.ID || idx" class="grid-item">
            <image :src="food.img" mode="aspectFill" class="grid-thumb" />
            <text class="grid-name">{{ food.name }}</text>
          </view>
        </view>

        <view class="toggle" @click="toggleExpand(order.id)">
          <uni-icons :type="isExpanded(order.id) ? 'arrowup' : 'arrowdown'" size="18" color="#999" />
        </view>
      </view>

      <view v-if="orders.length === 0" class="empty">
        <uni-icons type="paperclip" size="80" color="#ddd" />
        <text class="empty-text">这里空空如也~</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { onPullDownRefresh, onShow } from '@dcloudio/uni-app'
import api from '@/api'

const orders = ref([])
const selectedDate = ref('')
const calendarRef = ref(null)

const total = computed(() => orders.value.length)

const expandedSet = ref(new Set())
const isExpanded = (id) => expandedSet.value.has(id)
const toggleExpand = (id) => {
  const s = new Set(expandedSet.value)
  if (s.has(id)) s.delete(id)
  else s.add(id)
  expandedSet.value = s
}
const onEvaluate = () => { uni.showToast({ title: '评价功能开发中', icon: 'none' }) }

const getToday = () => {
  const d = new Date()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${d.getFullYear()}-${m}-${day}`
}

const getUserId = async () => {
  const c1 = uni.getStorageSync('userInfo')
  if (c1?.ID) return c1.ID
  const c2 = uni.getStorageSync('user')
  if (c2?.ID) return c2.ID
  const res = await api.getUserInfo()
  return res?.data?.ID || null
}

const loadOrders = async () => {
  const userId = await getUserId()
  if (!userId) return
  const res = await api.getOrderList({ userId, date: selectedDate.value })
  if (res?.code === 0) orders.value = res.data || []
}

const onDateChange = (e) => {
  selectedDate.value = e.fulldate || e.value || selectedDate.value
  loadOrders()
}
const onMonthSwitch = () => {}

onMounted(() => {
  selectedDate.value = getToday()
  loadOrders()
})

onShow(() => {
  // 返回页面时确保数据最新
  loadOrders()
})

onPullDownRefresh(() => {
  loadOrders().finally(() => uni.stopPullDownRefresh())
})

const onDelete = (order) => {
  uni.showModal({
    title: '删除确认',
    content: '确定删除该订单吗？',
    success: async (res) => {
      if (res.confirm) {
        const r = await api.deleteOrder({ id: order.id })
        if (r?.code === 0) {
          uni.showToast({ title: '删除成功', icon: 'success' })
          loadOrders()
        } else {
          uni.showToast({ title: r?.msg || '删除失败', icon: 'none' })
        }
      }
    }
  })
}

const formatDate = (s) => s?.replace(/T.*/, '') || ''
</script>

<style lang="scss" scoped>
.orders-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #FFE4E1 0%, #F3E8FF 50%, #E6E6FA 100%);
  padding-bottom: 40rpx;
}
.summary { padding: 20rpx 24rpx; color: #666; }
.calendar-box { margin: 0 16rpx 20rpx; background: #fff; border-radius: 16rpx; overflow: hidden; box-shadow: 0 6rpx 18rpx rgba(140,82,255,.08); }

.orders-list { padding: 0 20rpx; }
.order-card {
  background: #fff; border-radius: 16rpx; padding: 16rpx; box-shadow: 0 6rpx 18rpx rgba(0,0,0,.05); margin: 16rpx 0;
  .card-header { display:flex; align-items:center; justify-content:space-between; padding-bottom: 12rpx; border-bottom: 1rpx solid #f0f0f0; }
  .order-date { color:#666; font-size: 26rpx; }
  .card-actions { display:flex; gap: 16rpx; }
  .act-btn { display:flex; align-items:center; gap: 6rpx; }
  .act-text { color:#999; font-size: 22rpx; }

  .grid { display:flex; gap: 16rpx; padding: 14rpx 0; }
  .grid-item { text-align:center; }
  .grid-thumb { width: 140rpx; height: 140rpx; border-radius: 12rpx; object-fit: cover; }
  .grid-name { display:block; margin-top: 8rpx; font-size: 24rpx; color:#333; }

  .list { padding: 6rpx 0; }
  .list-item { display:flex; align-items:center; padding: 16rpx 0; border-bottom: 1rpx dashed #eee; }
  .list-item:last-child { border-bottom: none; }
  .list-thumb { width: 120rpx; height: 120rpx; border-radius: 12rpx; margin-right: 14rpx; object-fit: cover; }
  .list-info { display:flex; align-items:center; justify-content:space-between; flex:1; }
  .name { font-size: 28rpx; color:#333; }
  .qty { color:#999; font-size: 24rpx; }

  .toggle { text-align:center; padding: 6rpx 0; }
}

.empty { text-align: center; padding: 100rpx 0; color:#aaa; .empty-text{ display:block; margin-top: 12rpx; } }
</style> 