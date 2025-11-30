<template>
  <view class="diary-container">
    <!-- 日志列表 -->
    <view class="diary-list">
      <view 
        v-for="(diary, index) in diaryList" 
        :key="index" 
        class="diary-item"
        @click="viewDiary(diary)"
      >
        <view class="diary-header">
          <view class="author-info">
            <img :src="diary.avatar || diary.author?.avatar" class="author-avatar" alt="作者头像"></img>
            <view class="author-details">
              <text class="author-name">{{ diary.userName || diary.author?.name }}</text>
              <text class="diary-time">{{ formatTime(diary.createTime) }}</text>
            </view>
          </view>
          <view class="mood-indicator">
            <uni-icons :type="moodIcon(diary.mood)" size="20" :color="moodColor(diary.mood)"></uni-icons>
            <text class="mood-text">{{ moodText(diary.mood) }}</text>
          </view>
        </view>
        
        <view class="diary-content">
          <text class="diary-title">{{ diary.title }}</text>
          <text class="diary-text" :class="{ collapsed: !diary.expanded }">{{ diary.content }}</text>
          <view v-if="needToggle(diary.content)" class="toggle-more" @click.stop="toggleContent(diary)">
            <text class="toggle-text">{{ diary.expanded ? '收起' : '展开全部' }}</text>
            <uni-icons :type="diary.expanded ? 'arrowup' : 'arrowdown'" size="14" color="#999" />
          </view>
          
          <!-- 图片 -->
          <view v-if="diary.images && diary.images.length > 0" class="diary-images">
            <view 
              v-for="(img, imgIndex) in diary.images" 
              :key="imgIndex"
              class="image-item"
              @click.stop="previewImage(img, diary.images)"
            >
              <image
                :src="img"
                mode="aspectFill"
                class="diary-image"
              ></image>
            </view>
          </view>
        </view>
        
        <view class="diary-footer">
          <view class="diary-stats">
            <view class="stat-item" @click.stop="likeDiary(diary)">
              <uni-icons
                :type="diary.liked ? 'heart-filled' : 'heart'"
                size="16"
                :color="diary.liked ? '#FF6B9D' : '#999'"
              ></uni-icons>
              <text class="stat-text">{{ diary.likes || 0 }}</text>
            </view>
            <view class="stat-item">
              <uni-icons type="chat" size="16" color="#999"></uni-icons>
              <text class="stat-text">{{ diary.comments || 0 }}</text>
            </view>
          </view>
        </view>
      </view>
      
      <!-- 空状态 -->
      <view v-if="diaryList.length === 0" class="empty-state">
        <uni-icons type="compose" size="80" color="#ccc"></uni-icons>
        <text class="empty-text">还没有日志</text>
        <text class="empty-desc">记录你们的美好时光吧</text>
        <button
          @click="createDiary"
          class="create-btn"
        >开始写日志</button>
      </view>
    </view>

    <!-- 悬浮添加按钮 -->
    <view class="floating-add-btn" @click="createDiary">
      <uni-icons type="compose" size="24" color="#fff"></uni-icons>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { onShow, onPullDownRefresh, onReachBottom } from '@dcloudio/uni-app'
import api from '@/api'
import { useUserStore } from '@/store/user.js'

const userStore = useUserStore()

// 数据
const diaryList = ref([])
const selectedDate = ref('')
const diaryPage = ref(1)
const diaryPageSize = ref(10)
const diaryHasMore = ref(true)
const diaryLoading = ref(false)

// 统计
const totalDiaries = computed(() => diaryList.value.length)

// 工具
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

// 加载日记列表
const loadDiaries = async (reset = false) => {
  const userId = await getUserId()
  if (!userId || diaryLoading.value) return
  try {
    diaryLoading.value = true
    if (reset) {
      diaryPage.value = 1
      diaryHasMore.value = true
      diaryList.value = []
    }
    const res = await api.getDiaryList({ userId, date: selectedDate.value, page: diaryPage.value, pageSize: diaryPageSize.value })
    if (res?.code === 0) {
      const raw = res.data || {}
      const arr = Array.isArray(raw.list) ? raw.list : (Array.isArray(raw) ? raw : [])
      const list = arr.map(mapDiary)
      diaryList.value = diaryList.value.concat(list)
      const total = typeof raw.total === 'number' ? raw.total : undefined
      if (typeof total === 'number') {
        diaryHasMore.value = diaryList.value.length < total
      } else {
        diaryHasMore.value = list.length === diaryPageSize.value
      }
      if (diaryHasMore.value) diaryPage.value += 1
    }
  } finally {
    diaryLoading.value = false
  }
}

const mapDiary = (item) => {
  // 处理图片：img 字段逗号分隔
  let imgs = []
  if (Array.isArray(item.images)) imgs = item.images
  else if (typeof item.img === 'string') imgs = item.img.split(',').map(s => s.trim()).filter(Boolean)
  else if (typeof item.imgs === 'string') imgs = item.imgs.split(',').map(s => s.trim()).filter(Boolean)

  // mood 保持为数字，便于映射
  let moodNum = 9
  if (typeof item.mood === 'number') moodNum = item.mood
  else if (typeof item.mood === 'string') moodNum = parseInt(item.mood, 10) || 9

  return {
    id: item.id || item.ID,
    title: item.title || '无标题',
    content: item.content || item.text || '',
    userName: item.userName || item.nickname || item.author?.name || '好友',
    avatar: item.avatar || item.author?.avatar || '/static/images/default-avatar.png',
    author: {
      name: item.author?.name || item.nickname || '好友',
      avatar: item.author?.avatar || item.avatar || '/static/images/default-avatar.png'
    },
    createTime: item.createTime || item.createdAt || item.CreatedAt || new Date(),
    images: imgs,
    mood: moodNum,
    likes: item.likes || 0,
    comments: item.comments || 0,
    isLiked: false,
    expanded: false
  }
}

// UI 辅助
const formatTime = (time) => {
  const diaryTime = new Date(time)
  if (isNaN(diaryTime.getTime())) return ''
  const now = new Date()
  const diffDays = Math.floor((now - diaryTime) / (1000 * 60 * 60 * 24))
  if (diffDays === 0) return '今天'
  if (diffDays === 1) return '昨天'
  if (diffDays < 7) return `${diffDays}天前`
  return `${diaryTime.getFullYear()}-${String(diaryTime.getMonth()+1).padStart(2,'0')}-${String(diaryTime.getDate()).padStart(2,'0')}`
}

// mood: 1-开心 2-兴奋 3-幸福 4-思念 5-心累 6-疲倦 7-失落 8-生气 9-平静
const moodIcon = (mood) => ({
  1: 'heart',
  2: 'star',
  3: 'heart',
  4: 'paperplane',
  5: 'close',
  6: 'minus',
  7: 'sad',
  8: 'close',
  9: 'heart'
}[mood] || 'heart')

const moodColor = (mood) => ({
  1: '#FF6B9D',
  2: '#FFD700',
  3: '#FF8E9E',
  4: '#4A90E2',
  5: '#999999',
  6: '#A0A0A0',
  7: '#87CEEB',
  8: '#FF6347',
  9: '#666666'
}[mood] || '#FF6B9D')

const moodText = (mood) => ({
  1: '开心',
  2: '兴奋',
  3: '幸福',
  4: '思念',
  5: '心累',
  6: '疲倦',
  7: '失落',
  8: '生气',
  9: '平静'
}[mood] || '平静')

const viewDiary = (diary) => {
  uni.navigateTo({ url: `/pages/diary/detail?id=${diary.id}` })
}

const createDiary = () => {
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
  
  uni.navigateTo({ url: '/pages/diary/create' })
}

const likeDiary = (diary) => {
  diary.isLiked = !diary.isLiked
  diary.likes = Math.max(0, (diary.likes || 0) + (diary.isLiked ? 1 : -1))
}

const needToggle = (text) => {
  if (!text) return false
  return text.length > 60
}

const toggleContent = (diary) => {
  diary.expanded = !diary.expanded
}

const previewImage = (url, images) => {
  uni.previewImage({ urls: images, current: url })
}

// 生命周期
onMounted(() => {
  selectedDate.value = getToday()
  loadDiaries()
})

onShow(() => loadDiaries(true))

onPullDownRefresh(() => {
  loadDiaries(true).finally(() => uni.stopPullDownRefresh())
})

onReachBottom(() => {
  if (diaryHasMore.value && !diaryLoading.value) {
    loadDiaries(false)
  }
})
</script>

<style lang="scss" scoped>
.diary-container {
  min-height: 100vh;
  background: #f5f5f5;
  padding-bottom: 120rpx;
}

.header-stats {
  display: flex;
  align-items: center;
  justify-content: space-around;
  background: #fff;
  padding: 40rpx 30rpx;
  margin-bottom: 20rpx;
  
  .stats-item {
    text-align: center;
    
    .stats-number {
      display: block;
      font-size: 48rpx;
      font-weight: bold;
      color: #FF6B9D;
      margin-bottom: 10rpx;
    }
    
    .stats-label {
      font-size: 24rpx;
      color: #999;
    }
  }
  
  .stats-divider {
    width: 2rpx;
    height: 60rpx;
    background: #f0f0f0;
  }
}

.diary-list {
  padding: 0 30rpx;
  
  .diary-item {
    background: #fff;
    border-radius: 20rpx;
    padding: 30rpx;
    margin-bottom: 20rpx;
    box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.05);
    
    .diary-header {
      display: flex;
      align-items: center;
      margin-bottom: 20rpx;
      
      .author-info {
        display: flex;
        align-items: center;
        margin-right: 20rpx;
        
        .author-avatar {
          width: 40rpx;
          height: 40rpx;
          border-radius: 50%;
          border: 2rpx solid #fff;
          box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.1);
        }
        
        .author-details {
          margin-left: 20rpx;
          
          .author-name {
            display: block;
            font-size: 28rpx;
            font-weight: 500;
            color: #333;
            margin-bottom: 8rpx;
          }
          
          .diary-time {
            display: block;
            font-size: 24rpx;
            color: #999;
          }
        }
      }
      
      .mood-indicator {
        display: flex;
        align-items: center;
        background: #f8f8f8;
        border-radius: 50%;
        padding: 10rpx;
        
        .mood-text {
          margin-left: 10rpx;
          font-size: 24rpx;
          color: #666;
        }
      }
    }
    
    .diary-content {
      margin-bottom: 20rpx;
      
      .diary-title {
        display: block;
        font-size: 32rpx;
        font-weight: bold;
        color: #333;
        margin-bottom: 15rpx;
        line-height: 1.4;
      }
      
      .diary-text {
  display: block;
  font-size: 28rpx;
  color: #666;
  line-height: 1.6;
}
.diary-text.collapsed {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.toggle-more { margin-top: 8rpx; color:#999; display:flex; align-items:center; gap: 6rpx; }
.toggle-text { font-size: 24rpx; }
      
      .diary-images {
        display: flex;
        flex-wrap: wrap;
        gap: 15rpx;
        margin-top: 15rpx;
        
        .image-item {
          width: 120rpx;
          height: 120rpx;
          border-radius: 8rpx;
          overflow: hidden;
          
          .diary-image {
            width: 100%;
            height: 100%;
          }
        }
      }
    }
    
    .diary-footer {
      display: flex;
      justify-content: space-between;
      align-items: center;
      
      .diary-stats {
        display: flex;
        gap: 30rpx;
        
        .stat-item {
          display: flex;
          align-items: center;
          gap: 8rpx;
          
          .stat-text {
            font-size: 24rpx;
            color: #999;
          }
        }
      }
    }
  }
}

.empty-state {
  padding: 100rpx 0;
  text-align: center;
  
  .empty-text {
    display: block;
    font-size: 36rpx;
    color: #333;
    margin-top: 20rpx;
  }
  
  .empty-desc {
    display: block;
    font-size: 24rpx;
    color: #999;
    margin-top: 10rpx;
  }
  
  .empty-action {
    margin-top: 30rpx;
    
    .create-btn {
      background: linear-gradient(135deg, #FF6B9D 0%, #FF8E9E 100%);
      border: none;
      border-radius: 25rpx;
      padding: 20rpx 40rpx;
      font-size: 26rpx;
      font-weight: 500;
      color: #fff;
      cursor: pointer;
      transition: all 0.3s ease;
      box-shadow: 0 4rpx 15rpx rgba(255, 107, 157, 0.3);
      
      &:active {
        transform: scale(0.95);
        box-shadow: 0 2rpx 8rpx rgba(255, 107, 157, 0.4);
      }
    }
  }
}

.floating-add-btn {
  position: fixed;
  right: 40rpx;
  bottom: 120rpx;
  width: 100rpx;
  height: 100rpx;
  background: linear-gradient(135deg, #FF6B9D 0%, #FF8E9E 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8rpx 30rpx rgba(255, 107, 157, 0.3);
  z-index: 999;
  
  &:active {
    transform: scale(0.9);
  }
}
</style> 