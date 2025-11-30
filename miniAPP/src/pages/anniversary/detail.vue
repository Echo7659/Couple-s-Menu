<template>
  <view class="detail-container">
    <!-- 自定义头部 -->
    <view class="custom-header">
      <view class="status-bar"></view>
      <view class="header-content">
        <view class="back-btn" @click="goBack">
          <uni-icons type="left" size="24" color="#fff"></uni-icons>
        </view>
        <view class="header-title">纪念日详情</view>
        <view class="header-actions">
          <view class="edit-btn" @click="toggleEdit">
            <uni-icons :type="isEditing ? 'checkmarkempty' : 'compose'" size="24" color="#fff"></uni-icons>
          </view>
        </view>
      </view>
    </view>

    <!-- 主要内容区域 -->
    <view class="main-content">
      <!-- 纪念日卡片 -->
      <view class="anniversary-card">
        <view class="card-header">
          <view class="anniversary-icon">
            <uni-icons :type="getIconByType(anniversary.type || 'love')" size="60" color="#FF6B9D"></uni-icons>
          </view>
          <view class="anniversary-title">
            <text class="title-text">{{ anniversary.title }}</text>
            <text class="title-desc">{{ anniversary.description || '暂无描述' }}</text>
          </view>
        </view>

        <view class="card-body">
          <view class="info-item">
            <text class="info-label">纪念日期</text>
            <text class="info-value">{{ formatDate(anniversary.date) }}</text>
          </view>

          <view class="info-item">
            <text class="info-label">纪念日类型</text>
            <text class="info-value">{{ getTypeLabel(anniversary.type) }}</text>
          </view>

          <view class="info-item">
            <text class="info-label">倒计时</text>
            <text class="info-value countdown">{{ getCountdownText(anniversary) }}</text>
          </view>
        </view>
      </view>

      <!-- 编辑表单 -->
      <view v-if="isEditing" class="edit-form">
        <view class="form-item">
          <text class="form-label">纪念日名称</text>
          <uni-easyinput
            v-model="editForm.title"
            placeholder="请输入纪念日名称"
            class="form-input"
          ></uni-easyinput>
        </view>

        <view class="form-item">
          <text class="form-label">纪念日类型</text>
          <uni-data-checkbox
            v-model="editForm.type"
            :localdata="anniversaryTypes"
            mode="button"
            class="checkbox-group"
          ></uni-data-checkbox>
        </view>

        <view class="form-item">
          <text class="form-label">纪念日期</text>
          <view class="date-input" @click="openCalendar">
            <text class="date-text">{{ editForm.date || '请选择日期' }}</text>
            <uni-icons type="calendar" size="20" color="#666"></uni-icons>
          </view>
          <uni-calendar
            ref="calendarRef"
            :insert="false"
            :lunar="true"
            :start-date="'1900-01-01'"
            :end-date="'2100-12-31'"
            :date="editForm.date"
            @confirm="onCalendarConfirm"
          ></uni-calendar>
        </view>

        <view class="form-item">
          <text class="form-label">描述</text>
          <uni-easyinput
            v-model="editForm.description"
            placeholder="请输入描述（可选）"
            type="textarea"
            :maxlength="100"
            class="form-textarea"
          ></uni-easyinput>
        </view>

        <view class="form-actions">
          <button @click="cancelEdit" class="cancel-btn">取消</button>
          <button @click="saveChanges" class="save-btn">保存</button>
        </view>
      </view>

      <!-- 操作按钮 -->
      <view v-if="!isEditing" class="action-buttons">
        <button @click="deleteAnniversary" type="warn" class="delete-btn">删除纪念日</button>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/api'

// 响应式数据
const anniversary = ref({})
const isEditing = ref(false)
const showCalendar = ref(false)
const editForm = ref({})
const calendarRef = ref(null)

// 纪念日类型
const anniversaryTypes = [
  { value: 'love', label: '恋爱纪念' },
  { value: 'birthday', label: '生日' },
  { value: 'wedding', label: '结婚纪念' },
  { value: 'other', label: '其他' }
]

// 生命周期
onMounted(() => {
  // 从页面参数获取纪念日ID
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1]
  const anniversaryId = currentPage.options?.id

  if (anniversaryId) {
    loadAnniversaryDetail(anniversaryId)
  }
})

// 方法
const loadAnniversaryDetail = async (id) => {
  try {
    const res = await api.getMemorialDay(id)
    if (res?.code === 0 && res.data) {
      anniversary.value = res.data
      editForm.value = { ...res.data }
    } else {
      uni.showToast({
        title: '纪念日不存在',
        icon: 'none'
      })
      setTimeout(() => {
        goBack()
      }, 1500)
    }
  } catch (error) {
    console.error('加载纪念日详情失败:', error)
    uni.showToast({
      title: '加载失败',
      icon: 'none'
    })
    setTimeout(() => {
      goBack()
    }, 1500)
  }
}

const toggleEdit = () => {
  if (isEditing.value) {
    saveChanges()
  } else {
    isEditing.value = true
    editForm.value = {
      id: anniversary.value.ID,
      title: anniversary.value.title,
      type: anniversary.value.type || 'love',
      date: anniversary.value.date,
      description: anniversary.value.description || ''
    }
  }
}

const cancelEdit = () => {
  isEditing.value = false
  editForm.value = {
    id: anniversary.value.ID,
    title: anniversary.value.title,
    type: anniversary.value.type || 'love',
    date: anniversary.value.date,
    description: anniversary.value.description || ''
  }
}

const openCalendar = () => {
  if (calendarRef.value && calendarRef.value.open) {
    calendarRef.value.open()
  }
}

const onCalendarConfirm = (e) => {
  editForm.value.date = e.fulldate || e.value || ''
}

const saveChanges = async () => {
  if (!editForm.value.title.trim()) {
    uni.showToast({
      title: '请输入纪念日名称',
      icon: 'none'
    })
    return
  }

  if (!editForm.value.date) {
    uni.showToast({
      title: '请选择纪念日期',
      icon: 'none'
    })
    return
  }

  try {
    const res = await api.updateMemorialDay(editForm.value)
    if (res?.code === 0) {
      uni.showToast({
        title: '更新成功',
        icon: 'success'
      })
      anniversary.value = { ...editForm.value }
      isEditing.value = false
    } else {
      uni.showToast({
        title: res?.msg || '更新失败',
        icon: 'none'
      })
    }
  } catch (error) {
    console.error('更新纪念日失败:', error)
    uni.showToast({
      title: '更新失败',
      icon: 'none'
    })
  }
}

const deleteAnniversary = () => {
  uni.showModal({
    title: '确认删除',
    content: `确定要删除"${anniversary.value.title}"吗？`,
    success: async (res) => {
      if (res.confirm) {
        try {
          const result = await api.deleteMemorialDay({ id: anniversary.value.ID })
          if (result?.code === 0) {
            uni.showToast({
              title: '删除成功',
              icon: 'success'
            })
            setTimeout(() => {
              goBack()
            }, 1500)
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

const goBack = () => {
  const pages = getCurrentPages()
  if (pages.length > 1) {
    uni.navigateBack()
  } else {
    uni.switchTab({ url: '/pages/home/home' })
  }
}

const getIconByType = (type) => {
  const iconMap = {
    love: 'heart-fill',
    birthday: 'gift-fill',
    wedding: 'star-fill',
    other: 'calendar-fill'
  }
  return iconMap[type] || 'calendar-fill'
}

const getTypeLabel = (type) => {
  const found = anniversaryTypes.find(t => t.value === type)
  return found ? found.label : '其他'
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return `${date.getFullYear()}年${date.getMonth() + 1}月${date.getDate()}日`
}

const getCountdownText = (item) => {
  if (!item.date) return ''
  const today = new Date()
  const anniversary = new Date(item.date)

  // 设置今年的纪念日
  const thisYear = new Date(today.getFullYear(), anniversary.getMonth(), anniversary.getDate())

  // 如果今年的纪念日已过，计算明年的
  if (thisYear < today) {
    thisYear.setFullYear(thisYear.getFullYear() + 1)
  }

  const diffTime = thisYear.getTime() - today.getTime()
  const days = Math.ceil(diffTime / (1000 * 60 * 60 * 24))

  if (days === 0) return '就是今天'
  if (days < 0) return '已过'
  if (days === 1) return '明天'
  if (days <= 7) return `还有${days}天`
  if (days <= 30) return `还有${Math.ceil(days / 7)}周`
  return `还有${Math.ceil(days / 30)}个月`
}
</script>

<style lang="scss" scoped>
.detail-container {
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

    .back-btn {
      width: 60rpx;
      height: 60rpx;
      display: flex;
      align-items: center;
      justify-content: center;
    }

    .header-title {
      color: #fff;
      font-size: 36rpx;
      font-weight: bold;
    }

    .header-actions {
      .edit-btn {
        width: 60rpx;
        height: 60rpx;
        background: rgba(255, 255, 255, 0.2);
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
      }
    }
  }
}

.main-content {
  padding: 40rpx 30rpx;
}

.anniversary-card {
  background: #fff;
  border-radius: 20rpx;
  padding: 40rpx;
  margin-bottom: 40rpx;
  box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.05);

  .card-header {
    display: flex;
    align-items: center;
    margin-bottom: 40rpx;

    .anniversary-icon {
      margin-right: 30rpx;
    }

    .anniversary-title {
      flex: 1;

      .title-text {
        display: block;
        font-size: 36rpx;
        font-weight: bold;
        color: #333;
        margin-bottom: 10rpx;
      }

      .title-desc {
        font-size: 26rpx;
        color: #999;
      }
    }
  }

  .card-body {
    .info-item {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 20rpx 0;
      border-bottom: 1rpx solid #f0f0f0;

      &:last-child {
        border-bottom: none;
      }

      .info-label {
        font-size: 28rpx;
        color: #666;
      }

      .info-value {
        font-size: 28rpx;
        color: #333;
        font-weight: bold;

        &.countdown {
          color: #FF6B9D;
        }
      }
    }
  }
}

.edit-form {
  background: #fff;
  border-radius: 20rpx;
  padding: 40rpx;
  margin-bottom: 40rpx;
  box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.05);

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

    .form-textarea {
      background: #f8f8f8;
      border-radius: 10rpx;
      padding: 20rpx;
      min-height: 120rpx;
    }

    .checkbox-group {
      display: flex;
      flex-wrap: wrap;
      gap: 20rpx;

      .checkbox-item {
        flex: 1;
        min-width: 120rpx;
      }
    }

    .date-picker {
      display: flex;
      justify-content: space-between;
      align-items: center;
      background: #f8f8f8;
      border-radius: 10rpx;
      padding: 20rpx;

      .date-text {
        font-size: 28rpx;
        color: #333;
      }
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
  }

  .form-actions {
    display: flex;
    gap: 20rpx;

    .cancel-btn, .save-btn {
      flex: 1;
      height: 80rpx;
      border-radius: 40rpx;
    }

    .cancel-btn {
      background: #f5f5f5;
      color: #666;
    }

    .save-btn {
      background: linear-gradient(135deg, #FF6B9D 0%, #FF8E9E 100%);
      color: #fff;
    }
  }
}

.action-buttons {
  display: flex;
  gap: 20rpx;
  margin-top: 40rpx;
  
  .cancel-btn, .save-btn {
    flex: 1;
    height: 80rpx;
    border: none;
    border-radius: 40rpx;
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

.delete-section {
  margin-top: 30rpx;
  
  .delete-btn {
    width: 100%;
    height: 80rpx;
    background: linear-gradient(135deg, #FF6B6B 0%, #FF8E8E 100%);
    color: #fff;
    border: none;
    border-radius: 40rpx;
    font-size: 28rpx;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.3s ease;
    box-shadow: 0 4rpx 20rpx rgba(255, 107, 107, 0.3);
    
    &:active {
      transform: scale(0.98);
      box-shadow: 0 2rpx 10rpx rgba(255, 107, 107, 0.4);
    }
  }
}
</style>
