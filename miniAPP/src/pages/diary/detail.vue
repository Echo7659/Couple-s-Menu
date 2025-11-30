<template>
  <view class="detail-container">
    <view v-if="loading" class="loading">
      <text>加载中...</text>
    </view>
    
    <view v-else-if="diary" class="diary-content">
      <!-- 作者信息 -->
      <view class="author-section">
        <image :src="diary.avatar || diary.author?.avatar || '/static/images/default-avatar.png'" class="author-avatar" mode="aspectFill"></image>
        <view class="author-info">
          <text class="author-name">{{ diary.userName || diary.author?.name || '用户' }}</text>
          <text class="diary-time">{{ formatTime(diary.createTime) }}</text>
        </view>
        <view class="mood-indicator">
          <uni-icons :type="moodIcon(diary.mood)" size="24" :color="moodColor(diary.mood)"></uni-icons>
          <text class="mood-text">{{ moodText(diary.mood) }}</text>
        </view>
      </view>

      <!-- 标题 -->
      <view class="title-section">
        <text class="diary-title">{{ diary.title }}</text>
      </view>

      <!-- 内容 -->
      <view class="content-section">
        <text class="diary-text">{{ diary.content }}</text>
      </view>

      <!-- 图片 -->
      <view v-if="diary.images && diary.images.length > 0" class="images-section">
        <view class="images-grid">
          <view 
            v-for="(img, index) in diary.images" 
            :key="index"
            class="image-item"
            @click="previewImage(img, diary.images)"
          >
            <image
              :src="img"
              mode="aspectFill"
              class="diary-image"
            ></image>
          </view>
        </view>
      </view>

      <!-- 操作按钮 -->
      <view class="actions-section">
        <button class="edit-btn" @click="editDiary">编辑</button>
        <button class="delete-btn" @click="deleteDiary">删除</button>
      </view>
    </view>

    <view v-else class="empty-state">
      <text class="empty-text">日志不存在</text>
    </view>
  </view>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/api'

const diary = ref(null)
const loading = ref(true)
const diaryId = ref('')

// 心情映射
const moodMap = {
  1: { icon: 'heart', color: '#FF6B9D', text: '开心' },
  2: { icon: 'star', color: '#FFD700', text: '兴奋' },
  3: { icon: 'heart', color: '#FF8E9E', text: '幸福' },
  4: { icon: 'paperplane', color: '#4A90E2', text: '思念' },
  5: { icon: 'close', color: '#999', text: '心累' },
  6: { icon: 'minus', color: '#A0A0A0', text: '疲倦' },
  7: { icon: 'sad', color: '#87CEEB', text: '失落' },
  8: { icon: 'close', color: '#FF6347', text: '生气' },
  9: { icon: 'heart', color: '#666', text: '平静' }
}

const moodIcon = (mood) => moodMap[mood]?.icon || 'heart'
const moodColor = (mood) => moodMap[mood]?.color || '#666'
const moodText = (mood) => moodMap[mood]?.text || '平静'

// 格式化时间
const formatTime = (timeStr) => {
  if (!timeStr) return ''
  const date = new Date(timeStr)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  return `${year}-${month}-${day} ${hours}:${minutes}`
}

// 加载日志详情
const loadDiary = async () => {
  try {
    loading.value = true
    const res = await api.getDiaryDetail(diaryId.value)
    
    if (res?.code === 0 && res.data) {
      // 处理图片字段（可能是逗号分隔的字符串）
      const imgStr = res.data.img || ''
      const images = imgStr ? imgStr.split(',').filter(url => url.trim()) : []
      
      diary.value = {
        ...res.data,
        images: images
      }
    } else {
      uni.showToast({
        title: res?.msg || '加载失败',
        icon: 'none'
      })
    }
  } catch (error) {
    console.error('加载日志详情失败:', error)
    uni.showToast({
      title: '加载失败',
      icon: 'none'
    })
  } finally {
    loading.value = false
  }
}

// 预览图片
const previewImage = (current, urls) => {
  uni.previewImage({
    current: current,
    urls: urls || [current]
  })
}

// 编辑日志
const editDiary = () => {
  if (!diary.value || !diary.value.ID) {
    uni.showToast({
      title: '日志信息不完整',
      icon: 'none'
    })
    return
  }
  uni.navigateTo({
    url: `/pages/diary/create?id=${diary.value.ID}`
  })
}

// 删除日志
const deleteDiary = () => {
  uni.showModal({
    title: '提示',
    content: '确定要删除这篇日志吗？',
    success: async (res) => {
      if (res.confirm) {
        try {
          const result = await api.deleteDiary({ id: diary.value.ID })
          if (result?.code === 0) {
            uni.showToast({
              title: '删除成功',
              icon: 'success'
            })
            setTimeout(() => {
              uni.navigateBack()
            }, 1500)
          } else {
            uni.showToast({
              title: result?.msg || '删除失败',
              icon: 'none'
            })
          }
        } catch (error) {
          console.error('删除日志失败:', error)
          uni.showToast({
            title: '删除失败',
            icon: 'none'
          })
        }
      }
    }
  })
}

// 生命周期
onMounted(() => {
  // 获取路由参数
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1]
  const options = currentPage.options || {}
  diaryId.value = options.id || ''
  
  if (diaryId.value) {
    loadDiary()
  } else {
    loading.value = false
    uni.showToast({
      title: '缺少日志ID',
      icon: 'none'
    })
  }
})
</script>

<style lang="scss" scoped>
.detail-container {
  min-height: 100vh;
  background: #f5f5f5;
  padding-bottom: 40rpx;
}

.loading {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  color: #999;
}

.diary-content {
  padding: 20rpx;
}

.author-section {
  display: flex;
  align-items: center;
  background: #fff;
  padding: 30rpx;
  border-radius: 16rpx;
  margin-bottom: 20rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.05);

  .author-avatar {
    width: 80rpx;
    height: 80rpx;
    border-radius: 50%;
    margin-right: 20rpx;
  }

  .author-info {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 8rpx;

    .author-name {
      font-size: 30rpx;
      font-weight: 600;
      color: #333;
    }

    .diary-time {
      font-size: 24rpx;
      color: #999;
    }
  }

  .mood-indicator {
    display: flex;
    align-items: center;
    gap: 8rpx;
    padding: 8rpx 16rpx;
    background: #f8f8f8;
    border-radius: 20rpx;

    .mood-text {
      font-size: 24rpx;
      color: #666;
    }
  }
}

.title-section {
  background: #fff;
  padding: 30rpx;
  border-radius: 16rpx;
  margin-bottom: 20rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.05);

  .diary-title {
    font-size: 36rpx;
    font-weight: bold;
    color: #333;
    line-height: 1.6;
  }
}

.content-section {
  background: #fff;
  padding: 30rpx;
  border-radius: 16rpx;
  margin-bottom: 20rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.05);

  .diary-text {
    font-size: 28rpx;
    color: #333;
    line-height: 1.8;
    white-space: pre-wrap;
    word-break: break-all;
  }
}

.images-section {
  background: #fff;
  padding: 30rpx;
  border-radius: 16rpx;
  margin-bottom: 20rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.05);

  .images-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 16rpx;
  }

  .image-item {
    width: 100%;
    aspect-ratio: 1;
    border-radius: 12rpx;
    overflow: hidden;
    background: #f5f5f5;

    .diary-image {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }
  }
}

.actions-section {
  display: flex;
  gap: 20rpx;
  padding: 0 20rpx;
  margin-top: 40rpx;

  .edit-btn,
  .delete-btn {
    flex: 1;
    height: 80rpx;
    border-radius: 40rpx;
    font-size: 28rpx;
    border: none;
    line-height: 80rpx;
    padding: 0;
    
    &::after {
      border: none;
    }
  }

  .edit-btn {
    background: linear-gradient(135deg, #FF6B9D 0%, #FF8E9E 100%);
    color: #fff;
  }

  .delete-btn {
    background: #f5f5f5;
    color: #666;
    border: 2rpx solid #e0e0e0;
    box-sizing: border-box;
  }
}

.empty-state {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;

  .empty-text {
    font-size: 28rpx;
    color: #999;
  }
}
</style>

