<template>
  <view class="anniversary-container">
    <!-- 自定义头部 -->
    <view class="custom-header">
      <view class="status-bar"></view>
      <view class="header-content">
        <view class="back-btn" @click="goBack">
          <uni-icons type="left" size="24" color="#fff"></uni-icons>
        </view>
        <view class="header-title">纪念日</view>
        <view class="header-actions">
          <!-- 移除右上角的新增按钮 -->
        </view>
      </view>
    </view>

    <!-- 主要内容区域 -->
    <view class="main-content">
      <!-- 纪念日统计 -->
      <view class="stats-section">
        <view class="stat-item">
          <text class="stat-number">{{ anniversaryList.length }}</text>
          <text class="stat-label">总纪念日</text>
        </view>
        <view class="stat-item">
          <text class="stat-number">{{ upcomingCount }}</text>
          <text class="stat-label">即将到来</text>
        </view>
        <view class="stat-item">
          <text class="stat-number">{{ pastCount }}</text>
          <text class="stat-label">已过纪念日</text>
        </view>
      </view>

      <!-- 纪念日列表 -->
      <view class="anniversary-list">
        <view
          v-for="(item, index) in anniversaryList"
          :key="item.ID || index"
          class="anniversary-item"
          @click="viewAnniversary(item)"
        >
          <view class="anniversary-icon">
            <uni-icons type="calendar" size="32" color="#FF6B9D"></uni-icons>
          </view>
          <view class="anniversary-info">
            <text class="anniversary-name">{{ item.title }}</text>
            <text class="anniversary-date">{{ formatDate(item.date) }}</text>
          </view>
          <view class="anniversary-countdown">
            <text class="countdown-text" :class="getCountdownClass(item)">{{ getCountdownText(item) }}</text>
            <text class="countdown-days">{{ getCountdownDays(item) }}</text>
          </view>
          <view class="anniversary-actions">
            <view class="action-btn edit-btn" @click.stop="editAnniversary(item)">
              <uni-icons type="compose" size="18" color="#2196F3"></uni-icons>
              <text class="action-text">编辑</text>
            </view>
            <view class="action-btn delete-btn" @click.stop="deleteAnniversary(item)">
              <uni-icons type="trash" size="18" color="#F44336"></uni-icons>
              <text class="action-text">删除</text>
            </view>
          </view>
        </view>
      </view>

      <!-- 空状态 -->
      <view v-if="anniversaryList.length === 0" class="empty-state">
        <uni-icons type="calendar" size="80" color="#ccc"></uni-icons>
        <text class="empty-text">还没有添加纪念日</text>
        <text class="empty-desc">点击右下角添加第一个纪念日吧</text>
      </view>
    </view>

    <!-- 悬浮新增按钮 -->
    <view class="floating-add-btn" v-if="!showDialog" @click="showAddDialog">
      <uni-icons type="plus" size="32" color="#fff"></uni-icons>
    </view>

    <!-- 添加/编辑纪念日弹窗 -->
    <uni-popup ref="popup" type="center" :mask-click="false">
      <view class="dialog-content">
        <view class="dialog-header">
          <text class="dialog-title">{{ isEdit ? '编辑纪念日' : '添加纪念日' }}</text>
          <view class="close-btn" @click="closeDialog">
            <uni-icons type="close" size="24" color="#999"></uni-icons>
          </view>
        </view>

        <view class="dialog-body">
          <view class="form-item">
            <text class="form-label">纪念日名称</text>
            <uni-easyinput
              v-model="formData.title"
              placeholder="请输入纪念日名称"
              class="form-input"
            ></uni-easyinput>
          </view>

          <view class="form-item">
            <text class="form-label">纪念日类型</text>
            <view class="type-group">
              <view class="type-item" :class="{active: formData.type === 1}" @click="formData.type = 1">恋爱纪念日</view>
              <view class="type-item" :class="{active: formData.type === 2}" @click="formData.type = 2">生日</view>
              <view class="type-item" :class="{active: formData.type === 3}" @click="formData.type = 3">其他</view>
            </view>
          </view>

          <view class="form-item">
            <text class="form-label">纪念日期</text>
            <view class="date-input" @click="openCalendar">
              <text class="date-text">{{ formData.date || '请选择日期' }}</text>
              <uni-icons type="calendar" size="20" color="#666"></uni-icons>
            </view>
            <uni-calendar
              ref="calendarRef"
              :insert="false"
              :lunar="true"
              :start-date="'1900-01-01'"
              :end-date="'2100-12-31'"
              :date="formData.date"
              @confirm="onCalendarConfirm"
            ></uni-calendar>
          </view>
        </view>

        <view class="dialog-footer">
          <button @click="closeDialog" class="cancel-btn">取消</button>
          <button @click="saveAnniversary" class="save-btn">保存</button>
        </view>
      </view>
    </uni-popup>
  </view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { onPullDownRefresh } from '@dcloudio/uni-app'
import api from '@/api'

// 响应式数据
const anniversaryList = ref([])
const showDialog = ref(false)
const isEdit = ref(false)
const popup = ref(null)
const calendarRef = ref(null)
const formData = ref({
  id: null,
  user_id: null,
  lover_id: null,
  title: '',
  type: 3,
  date: ''
})

const currentUserId = ref(null)
const currentLoverId = ref(null)

const fetchIds = async () => {
  try {
    const cached = uni.getStorageSync('userInfo')
    if (cached?.ID) currentUserId.value = cached.ID
    else {
      const me = await api.getUserInfo()
      if (me?.code === 0 && me.data?.ID) currentUserId.value = me.data.ID
    }
  } catch (e) {}
  try {
    const loverRes = await api.getLoverInfo()
    if (loverRes?.code === 0 && loverRes.data?.loverId) currentLoverId.value = loverRes.data.loverId
  } catch (e) {}
}

// 计算属性
const upcomingCount = computed(() => {
  return anniversaryList.value.filter(item => getCountdownDays(item) > 0).length
})

const pastCount = computed(() => {
  return anniversaryList.value.filter(item => getCountdownDays(item) <= 0).length
})

// 生命周期
onMounted(async () => {
  await fetchIds()
  loadAnniversaryList()
})

// 下拉刷新
onPullDownRefresh(() => {
  loadAnniversaryList()
    .finally(() => {
      uni.stopPullDownRefresh()
    })
})

// 方法
// 加载纪念日列表
const loadAnniversaryList = async () => {
  try {
    const userId = currentUserId.value || uni.getStorageSync('userInfo')?.ID
    const res = await api.getMemorialDayList({ userId })
    if (res?.code === 0 && res.data?.list) {
      anniversaryList.value = res.data.list || []
      // 按日期排序
      anniversaryList.value.sort((a, b) => {
        const dateA = new Date(a.date)
        const dateB = new Date(b.date)
        return dateA - dateB
      })
    }
  } catch (error) {
    console.error('加载纪念日列表失败:', error)
    uni.showToast({
      title: '加载失败',
      icon: 'none'
    })
  }
}

const showAddDialog = () => {
  console.log('显示添加弹窗')
  isEdit.value = false
  const currentDate = getCurrentDate()

  // 设置表单数据
  formData.value = {
    id: null,
    user_id: currentUserId.value,
    lover_id: currentLoverId.value,
    title: '',
    type: 3,
    date: currentDate // 默认设置为当前日期
  }

  showDialog.value = true
  popup.value.open()
  // 悬浮按钮由 v-show 受 showDialog 控制自动隐藏
  console.log('显示弹窗状态:', showDialog.value)
  console.log('设置的默认日期:', currentDate)
}

const editAnniversary = (item) => {
  console.log('编辑纪念日:', item)
  isEdit.value = true
  formData.value = {
    id: item.ID,
    user_id: item.user_id,
    lover_id: item.lover_id,
    title: item.title,
    type: Number(item.type) || 3,
    date: item.date
  }

  console.log('设置表单数据:', formData.value)
  showDialog.value = true
  popup.value.open()
  console.log('显示弹窗状态:', showDialog.value)
}

const closeDialog = () => {
  showDialog.value = false
  popup.value.close()
}

const saveAnniversary = async () => {
  console.log('保存前的表单数据:', formData.value)

  if (!formData.value.title.trim()) {
    uni.showToast({
      title: '请输入纪念日名称',
      icon: 'none'
    })
    return
  }

  if (!formData.value.date) {
    uni.showToast({
      title: '请选择纪念日期',
      icon: 'none'
    })
    return
  }

  try {
    let res
    if (isEdit.value) {
      // 编辑时需要传递id
      console.log('提交编辑数据:', formData.value)
      res = await api.updateMemorialDay(formData.value)
    } else {
      // 新增时不需要传递id
      const { id, ...addData } = formData.value
      console.log('提交新增数据:', addData)
      res = await api.addMemorialDay(addData)
    }

    if (res?.code === 0) {
      uni.showToast({
        title: isEdit.value ? '更新成功' : '添加成功',
        icon: 'success'
      })
      closeDialog()
      loadAnniversaryList()
    } else {
      uni.showToast({
        title: res?.msg || '操作失败',
        icon: 'none'
      })
    }
  } catch (error) {
    console.error('保存纪念日失败:', error)
    uni.showToast({
      title: '操作失败',
      icon: 'none'
    })
  }
}

const deleteAnniversary = (item) => {
  uni.showModal({
    title: '确认删除',
    content: `确定要删除"${item.title}"吗？`,
    success: async (res) => {
      if (res.confirm) {
        try {
          const result = await api.deleteMemorialDay({ id: item.ID })
          if (result?.code === 0) {
            uni.showToast({
              title: '删除成功',
              icon: 'success'
            })
            loadAnniversaryList()
          } else {
            uni.showToast({
              title: result?.msg || '删除失败',
              icon: 'none'
            })
          }
        } catch (error) {
          console.error('删除纪念日失败:', error)
          uni.showToast({
            title: '删除失败',
            icon: 'none'
          })
        }
      }
    }
  })
}

const viewAnniversary = (item) => {
  uni.navigateTo({
    url: `/pages/anniversary/detail?id=${item.ID}`
  })
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return `${date.getFullYear()}年${date.getMonth() + 1}月${date.getDate()}日`
}

const getCountdownDays = (item) => {
  if (!item.date) return 0
  const today = new Date()
  const anniversary = new Date(item.date)

  // 设置今年的纪念日
  const thisYear = new Date(today.getFullYear(), anniversary.getMonth(), anniversary.getDate())

  // 如果今年的纪念日已过，计算明年的
  if (thisYear < today) {
    thisYear.setFullYear(thisYear.getFullYear() + 1)
  }

  const diffTime = thisYear.getTime() - today.getTime()
  return Math.ceil(diffTime / (1000 * 60 * 60 * 24))
}

const getCountdownText = (item) => {
  const days = getCountdownDays(item)
  if (days === 0) return '就是今天'
  if (days < 0) return '已过'
  if (days === 1) return '明天'
  if (days <= 7) return `还有${days}天`
  if (days <= 30) return `还有${Math.ceil(days / 7)}周`
  return `还有${Math.ceil(days / 30)}个月`
}

const getCountdownClass = (item) => {
  const days = getCountdownDays(item)
  if (days === 0) return 'countdown-today'
  if (days <= 7) return 'countdown-soon'
  if (days <= 30) return 'countdown-coming'
  return 'countdown-far'
}

const goBack = () => {
  const pages = getCurrentPages()
  if (pages.length > 1) {
    uni.navigateBack()
  } else {
    uni.switchTab({ url: '/pages/home/home' })
  }
}

const openCalendar = () => {
  if (calendarRef.value && calendarRef.value.open) {
    calendarRef.value.open()
  }
}

const getCurrentDate = () => {
  const today = new Date()
  const year = today.getFullYear()
  const month = String(today.getMonth() + 1).padStart(2, '0')
  const day = String(today.getDate()).padStart(2, '0')
  const currentDate = `${year}-${month}-${day}`
  console.log('获取当前日期:', currentDate)
  return currentDate
}

const onCalendarConfirm = (e) => {
  formData.value.date = e.fulldate || e.value || ''
}
</script>

<style lang="scss" scoped>
.anniversary-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #FFE4E1 0%, #FFB6C1 50%, #E6E6FA 100%);
}

.custom-header {
  background: linear-gradient(135deg, #FF6B9D 0%, #FF8E9E 100%);
  padding: 0 20rpx;

  .status-bar {
    height: var(--status-bar-height);
  }

  .header-content {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 20rpx 0;
    position: relative;

    .back-btn {
      width: 60rpx;
      height: 60rpx;
      display: flex;
      align-items: center;
      justify-content: center;
      z-index: 1;
    }

    .header-title {
      color: #fff;
      font-size: 36rpx;
      font-weight: bold;
      position: absolute;
      left: 50%;
      transform: translateX(-50%);
    }

    .header-actions {
      width: 60rpx;
    }
  }
}

.main-content {
  padding: 40rpx 30rpx;
}

.stats-section {
  display: flex;
  justify-content: space-around;
  background: #fff;
  border-radius: 20rpx;
  padding: 40rpx 20rpx;
  margin-bottom: 40rpx;
  box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.05);

  .stat-item {
    text-align: center;

    .stat-number {
      display: block;
      font-size: 48rpx;
      font-weight: bold;
      color: #FF6B9D;
      margin-bottom: 10rpx;
    }

    .stat-label {
      font-size: 24rpx;
      color: #999;
    }
  }
}

.anniversary-list {
  .anniversary-item {
    display: flex;
    align-items: center;
    background: #fff;
    padding: 30rpx;
    border-radius: 20rpx;
    margin-bottom: 20rpx;
    box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.05);

    .anniversary-icon {
      margin-right: 30rpx;
    }

    .anniversary-info {
      flex: 1;

      .anniversary-name {
        display: block;
        font-size: 32rpx;
        font-weight: bold;
        color: #333;
        margin-bottom: 10rpx;
      }

      .anniversary-date {
        display: block;
        font-size: 26rpx;
        color: #666;
        margin-bottom: 8rpx;
      }
    }

    .anniversary-countdown {
      text-align: center;
      margin-right: 20rpx;

      .countdown-text {
        display: block;
        font-size: 26rpx;
        font-weight: bold;
        margin-bottom: 8rpx;

        &.countdown-today {
          color: #FF6B9D;
        }

        &.countdown-soon {
          color: #FF6B9D;
        }

        &.countdown-coming {
          color: #FFA500;
        }

        &.countdown-far {
          color: #999;
        }
      }

      .countdown-days {
        font-size: 22rpx;
        color: #999;
      }
    }

    .anniversary-actions {
      display: flex;
      gap: 16rpx;

      .action-btn {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        padding: 16rpx 12rpx;
        border-radius: 16rpx;
        min-width: 80rpx;
        transition: all 0.2s ease;
        box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.1);
        cursor: pointer;
        position: relative;
        overflow: hidden;
        user-select: none;
        -webkit-tap-highlight-color: transparent;

        &::before {
          content: '';
          position: absolute;
          top: 0;
          left: -100%;
          width: 100%;
          height: 100%;
          background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
          transition: left 0.5s ease;
        }

        &:active::before {
          left: 100%;
        }

        &:active {
          transform: scale(0.92);
        }

        .action-text {
          font-size: 20rpx;
          margin-top: 8rpx;
          font-weight: 500;
          transition: all 0.2s ease;
          pointer-events: none;
        }

        &.edit-btn {
          background: linear-gradient(135deg, #E3F2FD 0%, #BBDEFB 100%);
          border: 1rpx solid #2196F3;

          &:active {
            background: linear-gradient(135deg, #BBDEFB 0%, #90CAF9 100%);
            box-shadow: 0 1rpx 4rpx rgba(0, 0, 0, 0.15);
            border-color: #1976D2;
          }

          &:active .action-text {
            color: #1976D2;
          }
        }

        &.delete-btn {
          background: linear-gradient(135deg, #FFEBEE 0%, #FFCDD2 100%);
          border: 1rpx solid #F44336;

          &:active {
            background: linear-gradient(135deg, #FFCDD2 0%, #EF9A9A 100%);
            box-shadow: 0 1rpx 4rpx rgba(0, 0, 0, 0.15);
            border-color: #D32F2F;
          }

          &:active .action-text {
            color: #D32F2F;
          }
        }
      }
    }
  }
}

.empty-state {
  text-align: center;
  padding: 100rpx 0;

  .empty-text {
    display: block;
    font-size: 32rpx;
    color: #999;
    margin: 30rpx 0 20rpx;
  }

  .empty-desc {
    font-size: 26rpx;
    color: #ccc;
  }
}

.dialog-content {
  background: #fff;
  border-radius: 20rpx;
  overflow: hidden;
  position: relative;
  z-index: 1000;

  .dialog-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 40rpx 40rpx 20rpx;
    border-bottom: 1rpx solid #f0f0f0;

    .dialog-title {
      font-size: 36rpx;
      font-weight: bold;
      color: #333;
    }

    .close-btn {
      width: 60rpx;
      height: 60rpx;
      display: flex;
      align-items: center;
      justify-content: center;
    }
  }

  .dialog-body {
    padding: 40rpx;

    .form-item {
      margin-bottom: 40rpx;

      .form-label {
        display: block;
        font-size: 28rpx;
        color: #333;
        margin-bottom: 20rpx;
        font-weight: bold;
      }

              .form-input {
          background: #f8f8f8;
          border-radius: 10rpx;
          padding: 20rpx;
        }

        .type-group {
          display: flex;
          gap: 20rpx;
        }

        .type-item {
          padding: 14rpx 22rpx;
          border-radius: 12rpx;
          border: 1rpx solid #ddd;
          color: #666;
          background: #f8f8f8;
        }

        .type-item.active {
          color: #fff;
          border-color: transparent;
          background: linear-gradient(135deg, #FF6B9D 0%, #FF8E9E 100%);
        }

        .date-input {
          display: flex;
          align-items: center;
          justify-content: space-between;
          background: #f8f8f8;
          border-radius: 10rpx;
          padding: 20rpx;
        }

        .date-text {
          font-size: 28rpx;
          color: #333;
        }

        .debug-info {
        font-size: 20rpx;
        color: #999;
        margin-top: 10rpx;
        display: block;
      }

      .test-date-btn {
        margin-top: 10rpx;
        padding: 10rpx 20rpx;
        background: #007AFF;
        color: #fff;
        border-radius: 10rpx;
        font-size: 20rpx;
        text-align: center;
        display: inline-block;
      }
    }
  }

  .dialog-footer {
    display: flex;
    gap: 20rpx;
    padding: 20rpx 40rpx 40rpx;
    
    .cancel-btn, .save-btn {
      flex: 1;
      height: 80rpx;
      border-radius: 40rpx;
      border: none;
      font-size: 28rpx;
      font-weight: 500;
      cursor: pointer;
      transition: all 0.3s ease;
    }
    
    .cancel-btn {
      background: #f5f5f5;
      color: #666;
      
      &:active {
        background: #e8e8e8;
        transform: scale(0.98);
      }
    }
    
    .save-btn {
      background: linear-gradient(135deg, #FF6B9D 0%, #FF8E9E 100%);
      color: #fff;
      
      &:active {
        transform: scale(0.98);
        box-shadow: 0 4rpx 15rpx rgba(255, 107, 157, 0.4);
      }
    }
  }
}

.floating-add-btn {
  position: fixed;
  bottom: 120rpx;
  right: 40rpx;
  width: 100rpx;
  height: 100rpx;
  background: linear-gradient(135deg, #FF6B9D 0%, #FF8E9E 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8rpx 30rpx rgba(255, 107, 157, 0.3);
  z-index: 999;
  transition: all 0.3s ease;

  &:active {
    transform: scale(0.95);
    box-shadow: 0 4rpx 15rpx rgba(255, 107, 157, 0.4);
  }

  animation: pulse 2s infinite;
}

@keyframes pulse {
  0% {
    box-shadow: 0 8rpx 30rpx rgba(255, 107, 157, 0.3);
  }
  50% {
    box-shadow: 0 8rpx 30rpx rgba(255, 107, 157, 0.6);
  }
  100% {
    box-shadow: 0 8rpx 30rpx rgba(255, 107, 157, 0.3);
  }
}
</style>
