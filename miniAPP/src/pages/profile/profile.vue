<template>
  <view class="profile-container">
    <!-- 用户信息卡片 -->
    <view class="user-card">
      <view class="user-header">
        <view class="user-avatar">
          <img :src="userInfo.avatar" :key="avatarKey" class="avatar-img" alt="用户头像"></img>
          <view class="edit-avatar" @click="editAvatar">
            <uni-icons type="camera" size="20" color="#fff"></uni-icons>
          </view>
        </view>
        <view class="user-info">
          <text class="user-name">{{ userInfo.nickname }}</text>
          <text class="user-id">ID: {{ userInfo.userId }}</text>
        </view>
        <view class="user-actions">
          <view class="action-btn" @click="editProfile">
            <uni-icons type="compose" size="20" color="#FF6B9D"></uni-icons>
          </view>
        </view>
      </view>

    </view>

    <!-- 快捷功能 -->
    <view class="quick-functions">
      <view class="function-group">
        <view class="group-title">我的功能</view>
        <view class="function-list">
          <view class="function-item" @click="goToMyRecipes">
            <view class="function-icon">
              <uni-icons type="cart" size="24" color="#FF6B9D"></uni-icons>
            </view>
            <text class="function-name">我的订单</text>
            <uni-icons type="right" size="16" color="#ccc"></uni-icons>
          </view>
          <!-- <view class="function-item" @click="goToMyDiaries">
            <view class="function-icon">
              <uni-icons type="compose" size="24" color="#52c41a"></uni-icons>
            </view>
            <text class="function-name">我的日志</text>
            <uni-icons type="right" size="16" color="#ccc"></uni-icons>
          </view> -->
          <view class="function-item" @click="goToMyGoals">
            <view class="function-icon">
              <uni-icons type="flag" size="24" color="#1890ff"></uni-icons>
            </view>
            <text class="function-name">打卡目标</text>
            <uni-icons type="right" size="16" color="#ccc"></uni-icons>
          </view>
          <!-- 手动订阅订单通知服务 -->
           <view class="function-item" @click="goToOrderNotify">
            <view class="function-icon">
              <uni-icons type="notification" size="24" color="#13c2c2"></uni-icons>
            </view>
            <text class="function-name">订单通知</text>
            <text class="quota">剩余{{ notifyLeft }}次</text>
            <uni-icons type="right" size="16" color="#ccc"></uni-icons>
             <text class="func-desc">手动订阅订单状态提醒</text>
           </view>
        </view>
      </view>

    </view>

    <!-- 退出登录按钮 -->
    <view class="logout-section">
      <button
        @click="logout"
        type="warn"
        class="logout-btn"
      >
        退出登录
      </button>
    </view>
  </view>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useUserStore } from '@/store/user.js'
import api from '@/api'
import config from '@/config'
import { onShow, onPullDownRefresh } from '@dcloudio/uni-app'
// #ifdef MP-WEIXIN
const TEMPLATE_ID = 'qwh5Cymj3VYvIfPMjI6zw-h-Pstg1saDBlbDysbODYU'
// 本地维护一个粗略的订阅剩余计数（仅作为展示，实际剩余次数需服务端统计/记录）
// 可将该值持久化在本地缓存
// #endif

const userStore = useUserStore()

const notifyLeft = ref(Number(uni.getStorageSync('notifyLeft') || 0))

// 响应式数据
const userInfo = ref({
  nickname: '',
  userId: '',
  avatar: '/static/images/default-avatar.png',
  level: 'LV.1 新手',
  recipeCount: 0,
  diaryCount: 0,
  checkinCount: 0
})

// 头像key，用于强制刷新
const avatarKey = ref(0)

// 获取用户信息
const getUserInfo = async () => {
  try {
    uni.showLoading({
      title: '加载中...'
    })

    const result = await api.getUserInfo()

    if (result.code === 0 && result.data) {
      console.log('获取用户信息成功:', result.data)

      // 更新用户信息
      userInfo.value = {
        nickname: result.data.nick_name || result.data.nickname || '用户',
        userId: result.data.ID || result.data.userId || '',
        avatar: result.data.avatar || '/static/images/default-avatar.png',
        level: getUserLevel(result.data.level || 1),
        recipeCount: result.data.recipe_count || result.data.recipeCount || 0,
        diaryCount: result.data.diary_count || result.data.diaryCount || 0,
        checkinCount: result.data.checkin_count || result.data.checkinCount || 0
      }

      // 同时更新两个缓存key
      uni.setStorageSync('user', result.data)
      uni.setStorageSync('userInfo', result.data)

    } else {
      throw new Error(result.msg || '获取用户信息失败')
    }
  } catch (error) {
    console.error('获取用户信息失败:', error)
    uni.showToast({
      title: error.message || '获取用户信息失败',
      icon: 'none'
    })

    // 如果获取失败，尝试从缓存读取
    const cachedUserInfo = uni.getStorageSync('userInfo') || uni.getStorageSync('user')
    if (cachedUserInfo) {
      userInfo.value = {
        nickname: cachedUserInfo.nick_name || cachedUserInfo.nickname || '用户',
        userId: cachedUserInfo.ID || cachedUserInfo.userId || '',
        avatar: cachedUserInfo.avatar || '/static/images/default-avatar.png',
        level: getUserLevel(cachedUserInfo.level || 1),
        recipeCount: cachedUserInfo.recipe_count || cachedUserInfo.recipeCount || 0,
        diaryCount: cachedUserInfo.diary_count || cachedUserInfo.diaryCount || 0,
        checkinCount: cachedUserInfo.checkin_count || cachedUserInfo.checkinCount || 0
      }
    }
  } finally {
    uni.hideLoading()
  }
}

// 获取用户等级
const getUserLevel = (level) => {
  const levelMap = {
    1: 'LV.1 新手',
    2: 'LV.2 初学者',
    3: 'LV.3 进阶者',
    4: 'LV.4 熟练者',
    5: 'LV.5 美食达人',
    6: 'LV.6 厨艺大师',
    7: 'LV.7 传奇厨师',
    8: 'LV.8 食神'
  }
  return levelMap[level] || 'LV.1 新手'
}

// 编辑头像
const editAvatar = () => {
  uni.chooseImage({
    count: 1,
    sizeType: ['compressed'],
    sourceType: ['album', 'camera'],
    success: (res) => {
      const tempFilePath = res.tempFilePaths[0]
      console.log('选择的头像:', tempFilePath)

      // 上传头像
      uploadAvatar(tempFilePath)
    }
  })
}

// 上传头像
const uploadAvatar = (tempFilePath) => {
  uni.showLoading({
    title: '上传头像中...'
  })

  // 获取token
  const token = uni.getStorageSync('token')
  if (!token) {
    uni.hideLoading()
    uni.showToast({
      title: '未找到登录token',
      icon: 'none'
    })
    return
  }

  uni.uploadFile({
    url: `${config.baseUrl}/user/avatar`,
    filePath: tempFilePath,
    name: 'avatar',
    header: {
      'Authorization': `${token}`
    },
    success: async (res) => {
      uni.hideLoading()
      console.log('头像上传响应:', res)

      try {
        const response = JSON.parse(res.data)
        if (response.code === 0) {
          const newAvatarUrl = response.data.url
          
          // 更新头像
          userInfo.value.avatar = newAvatarUrl
          // 强制刷新头像（通过改变key）
          avatarKey.value = Date.now()

          const updateRes = await api.updateUserInfo({
            avatar: newAvatarUrl
          })
          console.log('更新用户信息响应:', updateRes)

          // 更新缓存 - 同时更新两个缓存key
          const cachedUserInfo = uni.getStorageSync('userInfo') || {}
          cachedUserInfo.avatar = newAvatarUrl
          uni.setStorageSync('userInfo', cachedUserInfo)
          
          const cachedUser = uni.getStorageSync('user') || {}
          if (cachedUser.userInfo) {
            cachedUser.userInfo.avatar = newAvatarUrl
            uni.setStorageSync('user', cachedUser)
          } else {
            cachedUser.avatar = newAvatarUrl
            uni.setStorageSync('user', cachedUser)
          }

          // 更新状态管理
          if (userStore.setUserInfo) {
            userStore.setUserInfo({
              ...userStore.userInfo,
              avatar: newAvatarUrl
            })
          }

          uni.showToast({
            title: '头像更新成功',
            icon: 'success'
          })
        } else {
          throw new Error(response.msg || '头像上传失败')
        }
      } catch (error) {
        console.error('头像上传失败:', error)
        uni.showToast({
          title: error.message || '头像上传失败',
          icon: 'none'
        })
      }
    },
    fail: (error) => {
      uni.hideLoading()
      console.error('头像上传失败:', error)
      uni.showToast({
        title: '头像上传失败，请重试',
        icon: 'none'
      })
    }
  })
}

const editProfile = () => {
  uni.navigateTo({
    url: '/pages/user/setup'
  })
}

const goToMyRecipes = () => {
  // 访问个人订单前需要登录
  if (!userStore.hasLogin) {
    uni.navigateTo({ url: '/pages/user/login' })
    return
  }
  uni.navigateTo({
    url: '/pages/order/my-orders'
  })
}

const goToMyDiaries = () => {
  // 查看我的日志前需要登录
  if (!userStore.hasLogin) {
    uni.navigateTo({ url: '/pages/user/login' })
    return
  }
  uni.switchTab({
    url: '/pages/diary/diary'
  })
}

const goToMyGoals = () => {
  // 查看打卡目标前需要登录
  if (!userStore.hasLogin) {
    uni.navigateTo({ url: '/pages/user/login' })
    return
  }
  uni.switchTab({
    url: '/pages/target/index'
  })
}

const goToOrderNotify = () => {
  // 仅微信小程序支持
  // #ifdef MP-WEIXIN
     uni.requestSubscribeMessage({
     tmplIds: [TEMPLATE_ID],
     success: (res) => {
       if (res[TEMPLATE_ID] === 'accept') {
         notifyLeft.value = Math.max(0, notifyLeft.value + 1)
         uni.setStorageSync('notifyLeft', notifyLeft.value)
         uni.showToast({ title: '订阅成功', icon: 'success' })
       } else if (res[TEMPLATE_ID] === 'reject') {
         uni.showToast({ title: '已拒绝授权', icon: 'none' })
       } else {
         uni.showToast({ title: '状态：' + (res[TEMPLATE_ID] || '未知'), icon: 'none' })
       }
     },
     fail: (err) => {
       console.error('订阅失败', err)
       uni.showToast({ title: '订阅失败', icon: 'none' })
     }
   })
  // #endif
  // #ifndef MP-WEIXIN
  uni.showToast({ title: '仅微信小程序支持订阅', icon: 'none' })
  // #endif
}

const goToFavorites = () => {
  uni.showToast({
    title: '收藏功能开发中',
    icon: 'none'
  })
}

const goToSettings = () => {
  uni.showToast({
    title: '设置功能开发中',
    icon: 'none'
  })
}

const goToPrivacy = () => {
  uni.showToast({
    title: '隐私设置功能开发中',
    icon: 'none'
  })
}

const goToNotification = () => {
  // 兼容保留
  goToOrderNotify()
}

const goToHelp = () => {
  uni.showToast({
    title: '帮助中心功能开发中',
    icon: 'none'
  })
}

const goToFeedback = () => {
  uni.showToast({
    title: '意见反馈功能开发中',
    icon: 'none'
  })
}

const goToAbout = () => {
  uni.showToast({
    title: '关于我们功能开发中',
    icon: 'none'
  })
}

const logout = () => {
  uni.showModal({
    title: '提示',
    content: '确定要退出登录吗？',
    success: (res) => {
      if (res.confirm) {
        // 清除用户数据
        userStore.logout()
        uni.removeStorageSync('user')
        uni.removeStorageSync('token')
        uni.removeStorageSync('INVITE_CODE')
        uni.removeStorageSync('notifyLeft')
        uni.removeStorageSync('userInfo')

        uni.showToast({
          title: '退出成功',
          icon: 'success'
        })

        // 跳转到登录页
        setTimeout(() => {
          uni.reLaunch({
            url: '/pages/user/login'
          })
        }, 1500)
      }
    }
  })
}

// 页面是否已初始化
const isInitialized = ref(false)

// 生命周期
onMounted(() => {
  // 页面加载时获取用户信息
  getUserInfo().finally(() => {
    isInitialized.value = true
  })
})

// 页面显示时刷新用户信息
onShow(() => {
  console.log('profile页面 onShow 触发, isInitialized:', isInitialized.value)

  // 未登录时，直接跳转到登录页
  if (!userStore.hasLogin) {
    uni.navigateTo({ url: '/pages/user/login' })
    return
  }
  
  // 从缓存读取最新的用户信息（快速显示）
  const cachedUserInfo = uni.getStorageSync('userInfo') || uni.getStorageSync('user')
  if (cachedUserInfo) {
    console.log('从缓存读取用户信息:', cachedUserInfo)
    // 直接赋值整个对象，确保响应式更新
    userInfo.value = {
      nickname: cachedUserInfo.nick_name || cachedUserInfo.nickname || userInfo.value.nickname || '用户',
      userId: cachedUserInfo.ID || cachedUserInfo.userId || userInfo.value.userId || '',
      avatar: cachedUserInfo.avatar || userInfo.value.avatar || '/static/images/default-avatar.png',
      level: getUserLevel(cachedUserInfo.level || 1),
      recipeCount: cachedUserInfo.recipe_count || cachedUserInfo.recipeCount || userInfo.value.recipeCount || 0,
      diaryCount: cachedUserInfo.diary_count || cachedUserInfo.diaryCount || userInfo.value.diaryCount || 0,
      checkinCount: cachedUserInfo.checkin_count || cachedUserInfo.checkinCount || userInfo.value.checkinCount || 0
    }
    console.log('更新后的 userInfo:', userInfo.value)
  }
  
  // 如果页面已初始化，则静默刷新（避免与 onMounted 冲突）
  if (isInitialized.value) {
    refreshUserInfo()
  }
})

// 下拉刷新
onPullDownRefresh(() => {
  console.log('profile页面下拉刷新触发')
  // 强制刷新用户信息
  getUserInfo().finally(() => {
    uni.stopPullDownRefresh()
  })
})

// 静默刷新用户信息
const refreshUserInfo = async () => {
  try {
    console.log('开始静默刷新用户信息...')
    const result = await api.getUserInfo()
    if (result?.code === 0 && result.data) {
      console.log('接口返回的用户信息:', result.data)
      // 直接赋值整个对象，确保响应式更新
      userInfo.value = {
        nickname: result.data.nick_name || result.data.nickname || '用户',
        userId: result.data.ID || result.data.userId || '',
        avatar: result.data.avatar || '/static/images/default-avatar.png',
        level: getUserLevel(result.data.level || 1),
        recipeCount: result.data.recipe_count || result.data.recipeCount || 0,
        diaryCount: result.data.diary_count || result.data.diaryCount || 0,
        checkinCount: result.data.checkin_count || result.data.checkinCount || 0
      }
      // 强制刷新头像
      avatarKey.value = Date.now()
      // 更新缓存
      uni.setStorageSync('userInfo', result.data)
      uni.setStorageSync('user', result.data)
      console.log('用户信息刷新成功，更新后的 userInfo:', userInfo.value)
    }
  } catch (error) {
    console.error('刷新用户信息失败:', error)
    // 失败时使用缓存数据
    const cachedUserInfo = uni.getStorageSync('userInfo') || uni.getStorageSync('user')
    if (cachedUserInfo) {
      userInfo.value = {
        nickname: cachedUserInfo.nick_name || cachedUserInfo.nickname || '用户',
        userId: cachedUserInfo.ID || cachedUserInfo.userId || '',
        avatar: cachedUserInfo.avatar || '/static/images/default-avatar.png',
        level: getUserLevel(cachedUserInfo.level || 1),
        recipeCount: cachedUserInfo.recipe_count || cachedUserInfo.recipeCount || 0,
        diaryCount: cachedUserInfo.diary_count || cachedUserInfo.diaryCount || 0,
        checkinCount: cachedUserInfo.checkin_count || cachedUserInfo.checkinCount || 0
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.profile-container {
  min-height: 100vh;
  background: #f5f5f5;
  padding-bottom: 40rpx;
}

.user-card {
  background: linear-gradient(135deg, #FF6B9D 0%, #FF8E9E 100%);
  margin: 20rpx;
  border-radius: 20rpx;
  padding: 40rpx 30rpx;
  color: #fff;

  .user-header {
    display: flex;
    align-items: center;
    margin-bottom: 40rpx;

    .user-avatar {
      position: relative;
      margin-bottom: 30rpx;
      
      .avatar-img {
        width: 120rpx;
        height: 120rpx;
        border-radius: 50%;
        border: 4rpx solid #fff;
        box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.1);
      }
      
      .edit-avatar {
        position: absolute;
        bottom: 0;
        right: 0;
        width: 40rpx;
        height: 40rpx;
        background: #FF6B9D;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        border: 3rpx solid #fff;
        box-shadow: 0 2rpx 10rpx rgba(0, 0, 0, 0.2);
      }
    }

          .user-info {
        flex: 1;
        margin-left: 20rpx;

        .user-name {
          display: block;
          font-size: 36rpx;
          font-weight: bold;
          margin-bottom: 8rpx;
        }

        .user-id {
          display: block;
          font-size: 24rpx;
          opacity: 0.8;
        }
      }

    .user-actions {
      .action-btn {
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

  .user-stats {
    display: flex;
    align-items: center;
    justify-content: space-around;

    .stat-item {
      text-align: center;

      .stat-number {
        display: block;
        font-size: 40rpx;
        font-weight: bold;
        margin-bottom: 8rpx;
      }

      .stat-label {
        font-size: 24rpx;
        opacity: 0.9;
      }
    }

    .stat-divider {
      width: 2rpx;
      height: 60rpx;
      background: rgba(255, 255, 255, 0.3);
    }
  }
}

.quick-functions {
  margin: 0 20rpx;

  .function-group {
    background: #fff;
    border-radius: 20rpx;
    margin-bottom: 20rpx;
    overflow: hidden;

    .group-title {
      padding: 30rpx 30rpx 20rpx;
      font-size: 28rpx;
      font-weight: bold;
      color: #333;
      border-bottom: 2rpx solid #f5f5f5;
    }

    .function-list {
      .function-item {
        display: flex;
        align-items: center;
        padding: 30rpx;
        border-bottom: 2rpx solid #f5f5f5;

        &:last-child {
          border-bottom: none;
        }

        .function-icon {
          width: 60rpx;
          height: 60rpx;
          background: #f8f8f8;
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;
          margin-right: 20rpx;
        }

        .function-name {
          flex: 1;
          font-size: 30rpx;
          color: #333;
        }

        .quota { color:#FF6B9D; font-size: 24rpx; margin-right: 10rpx; }
        .func-desc { color:#999; font-size: 24rpx; margin-left: 10rpx; }
      }
    }
  }
}

.logout-section {
  margin: 40rpx 30rpx;
  
  .logout-btn {
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
