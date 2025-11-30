<template>
  <view class="target-container">
    <!-- 头部统计 & 操作 -->
    <view class="header-card">
      <view class="header-main">
        <view class="title-block">
          <text class="title">打卡目标</text>
          <text class="sub-title">和 TA 一起完成的小目标</text>
        </view>
        <view class="stat-block">
          <view class="stat-item">
            <text class="stat-number">{{ ongoingCount }}</text>
            <text class="stat-label">进行中</text>
          </view>
          <view class="stat-item">
            <text class="stat-number">{{ finishedCount }}</text>
            <text class="stat-label">已完成</text>
          </view>
        </view>
      </view>
      <button class="create-btn" type="primary" size="mini" @click="openCreate">
        新建目标
      </button>
    </view>

    <!-- 筛选 Tabs -->
    <view class="tabs">
      <view
        class="tab-item"
        :class="{ active: !showCompleted }"
        @click="changeFilter(false)"
      >
        进行中
      </view>
      <view
        class="tab-item"
        :class="{ active: showCompleted }"
        @click="changeFilter(true)"
      >
        已完成
      </view>
    </view>

    <!-- 目标列表（网格） -->
    <scroll-view
      class="list-scroll"
      scroll-y
      @scrolltolower="loadMore"
      :lower-threshold="80"
    >
      <view v-if="targetList.length > 0">
        <view class="target-list grid">
          <view
            v-for="item in targetList"
            :key="item.ID"
            class="target-item"
            @click="goDetail(item)"
          >
            <view class="target-grid-card">
              <view class="icon-wrapper">
                <view class="icon-circle">
                  <uni-icons
                    type="heart-filled"
                    size="22"
                    color="#FF6B9D"
                  ></uni-icons>
                </view>
                <text
                  v-if="item.isCompleted"
                  class="status-tag status-done"
                >已完成</text>
              </view>
              <uni-icons
                v-if="!item.isCompleted"
                class="lock-alone"
                type="locked"
                size="18"
                color="#FF4B7A"
              ></uni-icons>
              <text class="grid-title" :class="{ done: item.isCompleted }">
                {{ item.title }}
              </text>
            </view>
          </view>
        </view>

        <view v-if="loadingMore" class="load-more">
          <text>加载中...</text>
        </view>
        <view v-else-if="!hasMore" class="load-more finished-text center-tip">
          <text>没有更多了</text>
        </view>
      </view>

      <!-- 空状态 -->
      <view v-else class="empty-state">
        <uni-icons type="flag" size="80" color="#ccc" />
        <text class="empty-title">还没有打卡目标</text>
        <text class="empty-desc">和 TA 制定一个小目标，开始打卡吧～</text>
        <button class="empty-btn" type="primary" @click="openCreate">
          立即创建
        </button>
      </view>
    </scroll-view>

    <!-- 新建目标弹窗 -->
    <uni-popup ref="createPopupRef" type="center">
      <view class="popup-card">
        <text class="popup-title">新建打卡目标</text>
        <view class="popup-form-item">
          <text class="label">标题</text>
          <uni-easyinput
            v-model="createForm.title"
            placeholder="例如：一起看一场电影"
          />
        </view>
        <view class="popup-form-item">
          <text class="label">描述</text>
          <uni-easyinput
            v-model="createForm.description"
            type="textarea"
            :maxlength="200"
            placeholder="简单描述一下你们的目标吧～"
            autoHeight
          />
        </view>
        <view class="popup-actions">
          <button class="btn-cancel" @click="closeCreate">取消</button>
          <button
            class="btn-confirm"
            type="primary"
            :disabled="!canSubmitCreate"
            :loading="creating"
            @click="submitCreate"
          >
            确认创建
          </button>
        </view>
      </view>
    </uni-popup>
  </view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { onPullDownRefresh, onShow } from '@dcloudio/uni-app'
import api from '@/api'
import { useUserStore } from '@/store/user.js'

const userStore = useUserStore()

// 列表 & 分页
const targetList = ref([])
const page = ref(1)
const pageSize = ref(10)
const hasMore = ref(true)
const loading = ref(false)
const loadingMore = ref(false)

// 过滤：是否查看已完成
const showCompleted = ref(false)

// 统计
const ongoingCount = computed(
  () => targetList.value.filter(t => !t.isCompleted).length
)
const finishedCount = computed(
  () => targetList.value.filter(t => t.isCompleted).length
)

// 创建弹窗
const createPopupRef = ref(null)
const createForm = ref({
  title: '',
  description: ''
})
const creating = ref(false)

const canSubmitCreate = computed(
  () => !!createForm.value.title.trim() && !creating.value
)

// 获取当前用户 ID
const getUserId = async () => {
  const c1 = uni.getStorageSync('userInfo')
  if (c1?.ID) return c1.ID
  const res = await api.getUserInfo()
  return res?.data?.ID || null
}

// 加载我的目标列表（推荐接口）
const loadTargets = async (reset = false) => {
  const userId = await getUserId()
  if (!userId) return

  if (loading.value || loadingMore.value) return

  try {
    if (reset) {
      loading.value = true
      page.value = 1
      hasMore.value = true
      targetList.value = []
    } else {
      loadingMore.value = true
    }

    const params = {
      userId,
      page: page.value,
      pageSize: pageSize.value,
      isCompleted: showCompleted.value,
      title: ''
    }

    const res = await api.getMyTargets(params)
    if (res?.code === 0 && res.data) {
      const list = Array.isArray(res.data.list) ? res.data.list : []
      targetList.value = reset
        ? list
        : targetList.value.concat(list || [])

      const total = typeof res.data.total === 'number' ? res.data.total : null
      if (total !== null) {
        hasMore.value = targetList.value.length < total
      } else {
        hasMore.value = list.length === pageSize.value
      }

      if (hasMore.value) {
        page.value += 1
      }
    } else {
      uni.showToast({
        title: res?.msg || '加载目标失败',
        icon: 'none'
      })
    }
  } catch (e) {
    console.error('加载打卡目标失败:', e)
    uni.showToast({
      title: '加载失败',
      icon: 'none'
    })
  } finally {
    loading.value = false
    loadingMore.value = false
    uni.stopPullDownRefresh()
  }
}

// 下拉刷新
onPullDownRefresh(() => {
  loadTargets(true)
})

// 页面展示时刷新一次
onShow(() => {
  loadTargets(true)
})

onMounted(() => {
  loadTargets(true)
})

// 滚动加载更多
const loadMore = () => {
  if (!hasMore.value || loading.value || loadingMore.value) return
  loadTargets(false)
}

// 切换进行中 / 已完成
const changeFilter = (completed) => {
  if (showCompleted.value === completed) return
  showCompleted.value = completed
  loadTargets(true)
}

// 打开 / 关闭新建弹窗
const openCreate = () => {
  // 检查登录状态
  if (!userStore.hasLogin) {
    uni.showToast({
      title: '请先登录',
      icon: 'none'
    })
    setTimeout(() => {
      uni.navigateTo({ url: '/pages/user/login' })
    }, 1500)
    return
  }
  
  createForm.value = {
    title: '',
    description: ''
  }
  createPopupRef.value && createPopupRef.value.open()
}

const closeCreate = () => {
  createPopupRef.value && createPopupRef.value.close()
}

// 提交新建目标
const submitCreate = async () => {
  if (!canSubmitCreate.value) return
  try {
    creating.value = true
    const userId = await getUserId()
    if (!userId) {
      uni.showToast({
        title: '请先登录',
        icon: 'none'
      })
      return
    }
    const payload = {
      title: createForm.value.title.trim(),
      description: createForm.value.description.trim(),
      creatorId: userId
    }
    const res = await api.createTarget(payload)
    if (res?.code === 0) {
      uni.showToast({
        title: '创建成功',
        icon: 'success'
      })
      closeCreate()
      // 重新加载列表
      loadTargets(true)
    } else {
      uni.showToast({
        title: res?.msg || '创建失败',
        icon: 'none'
      })
    }
  } catch (e) {
    console.error('创建目标失败:', e)
    uni.showToast({
      title: '创建失败',
      icon: 'none'
    })
  } finally {
    creating.value = false
  }
}

// 跳转详情
const goDetail = (item) => {
  if (!item?.ID) return
  uni.navigateTo({
    url: `/pages/target/detail?id=${item.ID}`
  })
}

// 格式化时间
const formatDate = (timeStr) => {
  if (!timeStr) return ''
  const d = new Date(timeStr)
  if (Number.isNaN(d.getTime())) return ''
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}
</script>

<style lang="scss" scoped>
.target-container {
  min-height: 100vh;
  background: #f5f5f5;
  padding: 20rpx;
  box-sizing: border-box;
}

.header-card {
  background: linear-gradient(135deg, #ff9fb8, #ff6b9d);
  border-radius: 20rpx;
  padding: 24rpx 24rpx 20rpx;
  color: #fff;
  box-shadow: 0 8rpx 20rpx rgba(255, 107, 157, 0.35);
  display: flex;
  flex-direction: column;
  gap: 16rpx;
}

.header-main {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
}

.title-block {
  display: flex;
  flex-direction: column;
}

.title {
  font-size: 34rpx;
  font-weight: 700;
}

.sub-title {
  margin-top: 8rpx;
  font-size: 24rpx;
  opacity: 0.9;
}

.stat-block {
  display: flex;
  flex-direction: row;
  gap: 20rpx;
}

.stat-item {
  min-width: 100rpx;
  align-items: flex-end;
  text-align: right;
}

.stat-number {
  font-size: 32rpx;
  font-weight: 700;
}

.stat-label {
  display: block;
  font-size: 22rpx;
  margin-top: 4rpx;
}

.create-btn {
  align-self: flex-end;
  margin-top: 4rpx;
  padding: 0 28rpx;
  height: 60rpx;
  line-height: 60rpx;
  border-radius: 999rpx;
  font-size: 26rpx;
  background: #fff;
  color: #ff6b9d;
}

.tabs {
  margin-top: 24rpx;
  display: flex;
  flex-direction: row;
  background: #ffffff;
  border-radius: 999rpx;
  padding: 6rpx;
  box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.04);
}

.tab-item {
  flex: 1;
  text-align: center;
  padding: 16rpx 0;
  font-size: 26rpx;
  color: #999;
  border-radius: 999rpx;
}

.tab-item.active {
  background: #ff6b9d;
  color: #fff;
  font-weight: 600;
}

.list-scroll {
  margin-top: 16rpx;
  height: calc(100vh - 260rpx);
}

.target-list.grid {
  margin-top: 16rpx;
  display: flex;
  flex-wrap: wrap;
  margin-left: -8rpx;
  margin-right: -8rpx;
}

.target-item {
  width: 25%;
  padding: 0 8rpx 16rpx;
  box-sizing: border-box;
}

.target-grid-card {
  background: #fff7f9;
  border-radius: 16rpx;
  padding: 18rpx 8rpx 14rpx;
  box-shadow: 0 4rpx 10rpx rgba(255, 107, 157, 0.12);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: space-between;
  height: 170rpx;
  position: relative;
}

.icon-wrapper {
  width: 100%;
  display: flex;
  flex-direction: row;
  justify-content: center;
  align-items: flex-start;
  position: relative;
}

.icon-circle {
  width: 60rpx;
  height: 60rpx;
  border-radius: 50%;
  background: #ffe3ed;
  display: flex;
  align-items: center;
  justify-content: center;
}

.lock-alone {
  position: absolute;
  top: 70rpx;
  right: 10rpx;
}

.status-tag {
  position: absolute;
  top: 0;
  right: 0;
  font-size: 18rpx;
  padding: 2rpx 8rpx;
  border-radius: 999rpx;
  background: rgba(153, 153, 153, 0.16);
  color: #999;
}

.status-done {
  background: rgba(153, 153, 153, 0.18);
  color: #999;
}

.grid-title {
  margin-top: 14rpx;
  font-size: 22rpx;
  color: #333;
  text-align: center;
  line-height: 1.4;
  max-height: 2.8em;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.grid-title.done {
  color: #999;
}

.load-more {
  text-align: center;
  padding: 20rpx 0 40rpx;
  font-size: 24rpx;
  color: #999;
}

.load-more.finished-text {
  color: #ccc;
}

.empty-state {
  margin-top: 80rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #999;
}

.empty-title {
  margin-top: 24rpx;
  font-size: 30rpx;
  font-weight: 600;
  color: #666;
}

.empty-desc {
  margin-top: 8rpx;
  font-size: 24rpx;
}

.empty-btn {
  margin-top: 32rpx;
  width: 320rpx;
  border-radius: 999rpx;
  background: #ff6b9d;
  color: #fff;
}

.popup-card {
  width: 620rpx;
  background: #fff;
  border-radius: 20rpx;
  padding: 24rpx 28rpx 20rpx;
  box-sizing: border-box;
}

.popup-title {
  font-size: 30rpx;
  font-weight: 600;
  margin-bottom: 12rpx;
  color: #333;
}

.popup-form-item {
  margin-top: 12rpx;
}

.label {
  display: block;
  font-size: 26rpx;
  color: #555;
  margin-bottom: 8rpx;
}

.popup-actions {
  margin-top: 20rpx;
  display: flex;
  flex-direction: row;
  justify-content: flex-end;
  gap: 16rpx;
}

.btn-cancel {
  min-width: 140rpx;
  border-radius: 999rpx;
  background: #f5f5f5;
  color: #666;
  font-size: 26rpx;
}

.btn-confirm {
  min-width: 180rpx;
  border-radius: 999rpx;
  background: #ff6b9d;
  font-size: 26rpx;
}
</style>