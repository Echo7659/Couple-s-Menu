<template>
  <view class="container">
    <view class="header">
      <view class="header-back" @click="goBack">
        <uni-icons type="left" size="20" color="#333"></uni-icons>
      </view>
      <text class="header-title">获取用户头像昵称</text>
    </view>

    <view class="avatarUrl">
      <button type="balanced" open-type="chooseAvatar" @chooseavatar="onChooseAvatar">
        <image :src="avatarUrl" class="refreshIcon"></image>
      </button>
    </view>

    <view class="userName">
      <text>昵称：</text>
      <input
        :clearable="false"
        type="nickname"
        class="weui-input"
        :value="userName"
        @blur="bindBlur"
        placeholder="请输入昵称"
        @input="bindInput"
      />
    </view>

    <view style="width: 100%;height: 1px; background: #EEE;"></view>

    <!-- 角色选择 -->
    <view class="role-section">
      <text class="role-label">选择你的角色：</text>
      <view class="role-options">
        <view 
          class="role-option" 
          :class="{ active: selectedRole === 1 }"
          @click="selectRole(1)"
        >
          <view class="role-icon feeder">
            <uni-icons type="heart-filled" size="32" color="#FFD700"></uni-icons>
          </view>
          <text class="role-name">饲养员</text>
          <text class="role-desc">负责点餐下单</text>
        </view>
        <view 
          class="role-option" 
          :class="{ active: selectedRole === 2 }"
          @click="selectRole(2)"
        >
          <view class="role-icon foodie">
            <uni-icons type="heart-filled" size="32" color="#FF6B9D"></uni-icons>
          </view>
          <text class="role-name">吃货</text>
          <text class="role-desc">享受美食时光</text>
        </view>
      </view>
    </view>

    <view style="width: 100%;height: 1px; background: #EEE;"></view>

    <view style="width: 700rpx; height: 20px; font-size: 13px; margin: auto; margin-top: 40rpx;">
      · 申请获取以下权限
    </view>
    <view style="width: 700rpx; height: 20px; font-size: 13px; margin: auto; color: #cbcbcb; margin-top: 25rpx;">
      · 获得你的信息(昵称、头像等)
    </view>

    <view class="btn">
      <button @click="onSubmit" class="submit-btn">保存</button>
    </view>
  </view>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { pathToBase64, base64ToPath } from 'image-tools'
import api from '@/api'
import { useUserStore } from '@/store/user.js'
import config from '@/config'

const userStore = useUserStore()

// 响应式数据
const avatarUrl = ref('https://menu-echo.oss-cn-chengdu.aliyuncs.com/avatar/def.jpg')
const userName = ref('')
const appUserId = ref('')
const selectedRole = ref(null) // 1=饲养员，2=吃货


// 输入事件处理
const bindBlur = (e) => {
  userName.value = e.detail.value // 获取微信昵称
}

const bindInput = (e) => {
  userName.value = e.detail.value // 获取微信昵称
}

// 头像选择事件
const onChooseAvatar = (e) => {
  const { avatarUrl: newAvatarUrl } = e.detail
  avatarUrl.value = newAvatarUrl
}

// 选择角色
const selectRole = (role) => {
  selectedRole.value = role
}

// 提交事件
const onSubmit = () => {
  uni.login({
    "provider": "weixin",
    "onlyAuthorize": true, // 微信登录仅请求授权认证
    success: async function(event) {
      const { code } = event
      const name = userName.value

      // 判断头像、昵称和角色都不为空 再上传
      if (userName.value && avatarUrl.value && selectedRole.value) {
        console.log('头像URL:', avatarUrl.value) // 返回的是临时图片地址

        // 检查是否是临时文件，如果是则需要先上传
        if (avatarUrl.value.startsWith('http://tmp/') || avatarUrl.value.startsWith('wxfile://')) {
          console.log('检测到临时文件，开始上传头像...')
          uploadAvatar(avatarUrl.value).then(uploadResult => {
            if (uploadResult.success) {
              console.log('头像上传成功:', uploadResult.data)
              // 使用上传后的URL保存用户信息
              save(uploadResult.data.url)
            } else {
              uni.showToast({
                title: uploadResult.message || '头像上传失败',
                icon: 'none'
              })
            }
          }).catch(error => {
            console.error('头像上传失败:', error)
            uni.showToast({
              title: '头像上传失败，请重试',
              icon: 'none'
            })
          })
        } else {
          // 已经是网络图片，直接保存
          save(avatarUrl.value)
        }
      } else {
        let errorMsg = '请'
        const missing = []
        if (!userName.value) missing.push('填写昵称')
        if (!avatarUrl.value) missing.push('上传头像')
        if (!selectedRole.value) missing.push('选择角色')
        errorMsg += missing.join('、')
        uni.showToast({
          icon: 'none',
          title: errorMsg
        })
        return false
      }
    },
    fail: function(err) {
      console.log('err', err)
    }
  })
}

/**
 * 上传头像到服务器
 */
const uploadAvatar = (tempFilePath) => {
  return new Promise((resolve, reject) => {
    uni.showLoading({
      title: '上传头像中...'
    })

    // 获取token
    const token = uni.getStorageSync('token')
    if (!token) {
      uni.hideLoading()
      reject(new Error('未找到登录token'))
      return
    }

    // 使用uni.uploadFile上传文件
    uni.uploadFile({
      url: `${config.baseUrl}/user/avatar`,
      filePath: tempFilePath,
      name: 'avatar',
      header: {
        'Authorization': `${token}`
      },
      success: (res) => {
        uni.hideLoading()
        console.log('上传响应:', res)

        try {
          const response = JSON.parse(res.data)
          if (response.code === 0) {
            resolve({
              success: true,
              data: response.data
            })
          } else {
            resolve({
              success: false,
              message: response.msg || '上传失败'
            })
          }
        } catch (error) {
          console.error('解析上传响应失败:', error)
          resolve({
            success: false,
            message: '上传响应解析失败'
          })
        }
      },
      fail: (error) => {
        uni.hideLoading()
        console.error('上传失败:', error)
        reject(error)
      }
    })
  })
}

/**
 * @description 保存事件，用于保存用户头像昵称
 */
const save = async (finalAvatarUrl) => {
  uni.showLoading({
    title: '保存中...'
  })

  const data = {
    nick_name: userName.value,
    avatar: finalAvatarUrl,
    role: selectedRole.value
  }

  console.log('保存用户信息:', data)
  const result = await api.updateUserInfo(data)

  if (result.code === 0) {
    // 更新成功后，获取最新用户信息并缓存
    try {
      console.log('用户信息更新成功，开始获取最新用户信息...')
      const userInfoResult = await api.getUserInfo()

      if (userInfoResult.code === 0 && userInfoResult.data) {
        console.log('✅ 获取最新用户信息成功:', userInfoResult.data)

        // 缓存用户信息到userInfo（确保使用正确的key）
        uni.setStorageSync('userInfo', userInfoResult.data)
        uni.setStorageSync('user', userInfoResult.data)

        // 更新本地用户信息
        userStore.updateUserInfo({
          nickname: userInfoResult.data.nick_name || userInfoResult.data.nickname,
          avatar: userInfoResult.data.avatar,
          iniv_code: userInfoResult.data.iniv_code,
          role: userInfoResult.data.role
        })
        
        // 同时更新 userStore 的完整用户信息
        if (userStore.setUserInfo) {
          userStore.setUserInfo({
            ...userStore.userInfo,
            nickname: userInfoResult.data.nick_name || userInfoResult.data.nickname,
            avatar: userInfoResult.data.avatar,
            userId: userInfoResult.data.ID,
            iniv_code: userInfoResult.data.iniv_code,
            role: userInfoResult.data.role
          })
        }

        uni.hideLoading()
        uni.showToast({
          title: '设置成功',
          icon: 'success'
        })

        // 设置完成后，尝试根据邀请码自动绑定情侣关系
        try {
          const inviteCode = uni.getStorageSync('INVITE_CODE')
          console.log('设置完成，检查邀请码:', inviteCode)
          if (inviteCode) {
            console.log('✅ 开始绑定情侣，邀请码:', inviteCode)
            const bindResult = await api.bindLover({ inivCode: inviteCode })
            console.log('绑定情侣接口返回:', bindResult)
            if (bindResult?.code === 0) {
              uni.removeStorageSync('INVITE_CODE')
              uni.showToast({
                title: '已为你绑定情侣',
                icon: 'success'
              })
            } else {
              console.warn('绑定情侣失败:', bindResult?.msg || '未知错误')
            }
          } else {
            console.log('⚠️ 设置完成，但没有邀请码，跳过绑定')
          }
        } catch (e) {
          console.error('绑定情侣失败:', e)
        }

        // 跳转到首页
        setTimeout(() => {
          uni.switchTab({
            url: '/pages/home/home'
          })
        }, 1500)
      } else {
        throw new Error(userInfoResult.msg || '获取用户信息失败')
      }
    } catch (error) {
      console.error('获取最新用户信息失败:', error)
      // 如果获取失败，仍然显示设置成功，但使用更新后的数据
      uni.showToast({
        title: '设置成功',
        icon: 'success'
      })

      // 使用更新后的数据缓存
      const cachedUser = uni.getStorageSync('userInfo') || {}
      const updatedUserInfo = {
        ...cachedUser,
        ID: cachedUser.ID || userStore.userInfo?.userId,
        nick_name: userName.value,
        nickname: userName.value,
        avatar: finalAvatarUrl,
        phone: cachedUser.phone || userStore.userInfo?.phone || '',
        iniv_code: cachedUser.iniv_code || userStore.userInfo?.iniv_code,
        role: selectedRole.value
      }
      uni.setStorageSync('userInfo', updatedUserInfo)
      uni.setStorageSync('user', updatedUserInfo)
      
      // 更新 userStore
      userStore.updateUserInfo({
        nickname: userName.value,
        avatar: finalAvatarUrl,
        role: selectedRole.value
      })
      
      if (userStore.setUserInfo) {
        userStore.setUserInfo({
          ...userStore.userInfo,
          nickname: userName.value,
          avatar: finalAvatarUrl,
          role: selectedRole.value
        })
      }

      uni.hideLoading()
      setTimeout(() => {
        uni.switchTab({
          url: '/pages/home/home'
        })
      }, 1500)
    }
  } else {
    uni.hideLoading()
    uni.showToast({
      title: result.msg || '保存失败',
      icon: 'none'
    })
  }
}

// 返回
const goBack = () => {
  const pages = getCurrentPages()
  if (pages.length > 1) {
    uni.navigateBack()
  } else {
    uni.switchTab({
      url: '/pages/profile/profile'
    })
  }
}

// 生命周期
onMounted(() => {
  // 初始化随机昵称
  userName.value = getRandomName(Math.floor(Math.random() * (6 - 2) + 3))
})
</script>

<style lang="scss" scoped>
.container {
  width: 100vw;
  height: 100vh;
  background: #fff;
  box-sizing: border-box;
  padding: 0 30rpx;

  .header {
    display: flex;
    align-items: center;
    padding: 80rpx 0 40rpx;
    background: #fff;
    position: relative;

    .header-back {
      position: absolute;
      left: 0;
      width: 60rpx;
      height: 60rpx;
      display: flex;
      align-items: center;
      justify-content: center;
      cursor: pointer;
    }

    .header-title {
      flex: 1;
      text-align: center;
      color: #333;
      font-weight: 400;
    }
  }

  .avatarUrl {
    padding: 80rpx 0 40rpx;
    background: #fff;

    button {
      background: #fff;
      line-height: 80rpx;
      height: auto;
      border: none !important;
      width: auto;
      margin: 0;
      display: flex;
      border: none;
      justify-content: center;
      align-items: center;

      &::after {
        border: none;
      }

      .refreshIcon {
        width: 160rpx;
        height: 160rpx;
        border-radius: 50%;
        background-color: #ccc;
      }

      .jt {
        width: 14rpx;
        height: 28rpx;
      }
    }
  }

  .userName {
    background: #fff;
    padding: 20rpx 30rpx 80rpx;
    display: flex;
    align-items: center;
    justify-content: center;

    .weui-input {
      padding-left: 60rpx;
    }
  }

  .role-section {
    padding: 40rpx 30rpx;
    background: #fff;
    
    .role-label {
      display: block;
      font-size: 28rpx;
      color: #333;
      font-weight: 600;
      margin-bottom: 30rpx;
    }
    
    .role-options {
      display: flex;
      gap: 30rpx;
      justify-content: space-between;
      
      .role-option {
        flex: 1;
        display: flex;
        flex-direction: column;
        align-items: center;
        padding: 40rpx 20rpx;
        border: 3rpx solid #e0e0e0;
        border-radius: 20rpx;
        background: #fafafa;
        transition: all 0.3s ease;
        cursor: pointer;
        
        &.active {
          border-color: #FF6B9D;
          background: linear-gradient(135deg, rgba(255, 107, 157, 0.05) 0%, rgba(255, 142, 158, 0.05) 100%);
          box-shadow: 0 4rpx 20rpx rgba(255, 107, 157, 0.2);
          transform: translateY(-4rpx);
        }
        
        .role-icon {
          width: 80rpx;
          height: 80rpx;
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;
          margin-bottom: 20rpx;
          
          &.feeder {
            background: linear-gradient(135deg, rgba(255, 215, 0, 0.2) 0%, rgba(255, 165, 0, 0.2) 100%);
          }
          
          &.foodie {
            background: linear-gradient(135deg, rgba(255, 107, 157, 0.2) 0%, rgba(255, 142, 158, 0.2) 100%);
          }
        }
        
        .role-name {
          font-size: 32rpx;
          font-weight: 600;
          color: #333;
          margin-bottom: 10rpx;
        }
        
        .role-desc {
          font-size: 24rpx;
          color: #999;
        }
      }
    }
  }

  .btn {
    margin-top: 60rpx;
    
    .submit-btn {
      width: 100%;
      height: 80rpx;
      background: linear-gradient(135deg, #FF6B9D 0%, #FF8E9E 100%);
      color: #fff;
      border: none;
      border-radius: 40rpx;
      font-size: 28rpx;
      font-weight: 500;
      cursor: pointer;
      transition: all 0.3s ease;
      box-shadow: 0 4rpx 20rpx rgba(255, 107, 157, 0.3);
      
      &:active {
        transform: scale(0.98);
        box-shadow: 0 2rpx 10rpx rgba(255, 107, 157, 0.4);
      }
    }
  }
}
</style>
