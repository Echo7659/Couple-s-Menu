<template>
  <view class="login-container">
    <!-- 背景装饰 -->
    <view class="bg-decoration">
      <view class="circle circle-1"></view>
      <view class="circle circle-2"></view>
      <view class="circle circle-3"></view>
    </view>

    <!-- 主要内容 -->
    <view class="main-content">
      <!-- Logo和标题 -->
      <view class="header-section">
        <view class="app-logo">
          <uni-icons type="heart-filled" size="80" color="#FF6B9D"></uni-icons>
        </view>
        <!-- <text class="app-title">专属于小呆呆的菜谱</text> -->
        <text class="app-title">情侣菜单-记录我们的</text>
        <text class="app-subtitle">记录我们的美食时光</text>
      <br>
        <text class="app-subtitle">爱就是要在一起吃很多很多饭呀</text>
      </view>

      <!-- 功能介绍 -->
      <view class="features-section">
        <view class="feature-item">
          <uni-icons type="compose" size="24" color="#52c41a"></uni-icons>
          <text class="feature-text">记录美好时光</text>
        </view>
        <view class="feature-item">
          <uni-icons type="flag" size="24" color="#1890ff"></uni-icons>
          <text class="feature-text">打卡目标完成</text>
        </view>
      </view>

      <!-- 登录按钮 -->
      <view class="login-section">
        <view class="login-tips">
          <text class="tips-text">登录后即可使用全部功能</text>
        </view>
        
        <button 
          class="wechat-login-btn"
          @click="onGetUserInfo"
          :loading="isLogging"
          :disabled="isLogging"
        >
          <view class="btn-content">
            <text class="btn-text">{{ isLogging ? '登录中...' : '微信一键登录' }}</text>
          </view>
        </button>

        <!-- 暂不登录按钮 -->
        <view class="skip-login-section">
          <text class="skip-login-btn" @click="skipLogin">暂不登录，先逛逛</text>
        </view>

        <view class="privacy-tips">
          <label class="agreement-checkbox">
            <checkbox 
              :checked="agreedToTerms" 
              @tap="toggleAgreement"
              color="#FF6B9D"
            />
            <text class="privacy-text">登录即表示同意</text>
            <text class="privacy-link" @click.stop="viewUserAgreement">《用户协议》</text>
            <text class="privacy-text">和</text>
            <text class="privacy-link" @click.stop="viewPrivacyPolicy">《隐私政策》</text>
          </label>
        </view>
      </view>
    </view>

    <!-- 底部装饰 -->
    <view class="bottom-decoration">
      <view class="wave wave-1"></view>
      <view class="wave wave-2"></view>
    </view>

    <!-- 用户协议弹窗 -->
    <view v-if="showAgreementPopup" class="popup-mask" @click="closeAgreement">
      <view class="agreement-popup" @click.stop>
        <view class="popup-header">
          <text class="popup-title">{{ agreementTitle }}</text>
          <view class="close-btn" @click="closeAgreement">
            <uni-icons type="close" size="20" color="#666"></uni-icons>
          </view>
        </view>

        <scroll-view class="popup-content" scroll-y>
          <view class="agreement-text-wrapper">
            <text class="agreement-text">{{ agreementContent }}</text>
          </view>
        </scroll-view>

        <view class="popup-footer">
          <button class="popup-confirm-btn" @click="closeAgreement">我知道了</button>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useUserStore } from '@/store/user.js'
import api from '@/api'
import { onLoad } from '@dcloudio/uni-app'

const userStore = useUserStore()

// 响应式数据
const isLogging = ref(false)

// 是否同意协议
const agreedToTerms = ref(false)

// 协议弹窗相关
const agreementTitle = ref('')
const agreementContent = ref('')
const showAgreementPopup = ref(false)

// 邀请码（通过分享链接携带的 iniv_code）
const inviteCode = ref('')

// 解析页面入参中的邀请码，并暂存到本地，供登录成功后绑定情侣使用
onLoad((options) => {
  console.log('登录页 onLoad 被调用')
  console.log('登录页 onLoad options 参数:', JSON.stringify(options, null, 2))
  console.log('登录页 onLoad options 类型:', typeof options)
  console.log('登录页 onLoad options 是否为对象:', options && typeof options === 'object')
  
  // 获取当前页面实例，检查页面参数
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1]
  console.log('当前页面实例:', currentPage)
  console.log('当前页面 route:', currentPage?.route)
  console.log('当前页面 options:', JSON.stringify(currentPage?.options, null, 2))
  
  let invite = null
  
  // 方式1: 直接从 onLoad 的 options 参数中获取（最常见的方式）
  if (options && typeof options === 'object') {
    invite = options.invite || options.inviteCode
    console.log('方式1 - 从 onLoad options 获取:', invite)
  }
  
  // 方式2: 从当前页面实例的 options 中获取
  if (!invite && currentPage && currentPage.options) {
    invite = currentPage.options.invite || currentPage.options.inviteCode
    console.log('方式2 - 从页面实例 options 获取:', invite)
  }
  
  // 方式3: 从 URL 路径中解析（如果参数在路径中）
  if (!invite && currentPage && currentPage.route) {
    const route = currentPage.route
    const fullPath = currentPage.$page?.fullPath || ''
    console.log('页面 route:', route)
    console.log('页面 fullPath:', fullPath)
    
    // 尝试从 fullPath 中解析
    if (fullPath) {
      const match = fullPath.match(/[?&]invite=([^&]+)/)
      if (match) {
        invite = decodeURIComponent(match[1])
        console.log('方式3 - 从 fullPath 解析:', invite)
      }
    }
  }
  
  // 如果找到了邀请码，保存
  if (invite) {
    inviteCode.value = invite
    uni.setStorageSync('INVITE_CODE', inviteCode.value)
    console.log('✅ 成功获取到邀请码:', inviteCode.value)
  } else {
    // 如果页面参数没有，尝试从本地存储读取（可能从 App.vue 的 onLaunch 中保存的）
    const cachedInvite = uni.getStorageSync('INVITE_CODE')
    if (cachedInvite) {
      inviteCode.value = cachedInvite
      console.log('✅ 从本地存储获取到邀请码:', inviteCode.value)
    } else {
      console.log('⚠️ 未找到邀请码，所有方式都尝试过了')
    }
  }
})

// 登录成功后，尝试根据邀请码自动绑定情侣关系
const tryBindLoverByInvite = async () => {
  // 重新从多个来源获取邀请码
  let code = inviteCode.value || uni.getStorageSync('INVITE_CODE')
  
  // 如果还是没有，尝试从当前页面参数获取
  if (!code) {
    const pages = getCurrentPages()
    const currentPage = pages[pages.length - 1]
    if (currentPage && currentPage.options) {
      code = currentPage.options.invite || currentPage.options.inviteCode
    }
  }
  
  console.log('尝试绑定情侣，邀请码:', code)
  
  if (!code) {
    console.log('⚠️ 没有邀请码，跳过绑定')
    return
  }
  
  try {
    console.log('开始调用绑定情侣接口，邀请码:', code)
    const res = await api.bindLover({ inivCode: code })
    console.log('绑定情侣接口返回:', res)
    
    if (res?.code === 0) {
      // 绑定成功后清理本地邀请码，避免重复绑定
      uni.removeStorageSync('INVITE_CODE')
      inviteCode.value = ''
      uni.showToast({
        title: '已为你绑定情侣',
        icon: 'success'
      })
    } else {
      console.warn('绑定情侣失败:', res?.msg || '未知错误')
    }
  } catch (e) {
    console.error('根据邀请码绑定情侣失败:', e)
  }
}

// 启动时检查：如果已登录且 token 未过期，静默更新用户信息（不强制跳转）
onMounted(async () => {
  try {
    const token = uni.getStorageSync('token')
    if (!token) return
    // 尝试获取一次用户信息，以验证 token 可用
    const res = await api.getUserInfo()
    if (res?.code === 0 && res.data) {
      // 静默更新用户信息，但不强制跳转
      uni.setStorageSync('userInfo', res.data)
      userStore.setUserInfo({
        userId: res.data.ID || res.data.userId,
        nickname: res.data.nick_name || res.data.nickname,
        avatar: res.data.avatar,
        phone: res.data.phone,
        iniv_code: res.data.iniv_code
      })
    }
  } catch (e) {
    // token 失效则停留在登录页，不做处理
  }
})

// 查看用户协议
const viewUserAgreement = () => {
  console.log('查看用户协议')
  agreementTitle.value = '用户协议'
  agreementContent.value = `生效日期：2024年1月1日

欢迎使用"情侣菜谱"小程序！在使用本服务前，请您仔细阅读并同意本用户协议。

一、协议范围
本协议是您与"情侣菜谱"小程序（以下简称"本服务"）之间关于使用本服务的法律协议。

二、服务内容
本服务是一款专为情侣设计的应用，提供以下功能：
- 情侣信息绑定与管理
- 日志记录与分享
- 菜谱管理与点餐服务
- 打卡目标管理
- 纪念日提醒

三、用户行为规范
1. 账户安全：您有责任保管好您的账户信息，不得将账户转让、出售或出借给他人使用。
2. 内容规范：您在使用本服务时，应遵守相关法律法规，不得发布违法违规、虚假、侵权、骚扰、低俗等内容。
3. 禁止行为：禁止利用本服务从事任何违法违规活动。

四、隐私保护
我们重视您的隐私保护，具体内容请参见《隐私政策》。

五、服务变更与终止
1. 我们有权根据业务发展需要调整、变更或终止部分或全部服务内容。
2. 如您违反本协议，我们有权立即终止向您提供服务。

六、免责声明
1. 本服务按"现状"提供，我们不对服务的及时性、准确性、完整性作任何保证。
2. 因不可抗力、网络故障等原因导致的服务中断或数据丢失，我们不承担责任。

七、知识产权
本服务的所有知识产权归我们所有，未经授权不得使用。

八、协议修改
我们有权随时修改本协议，修改后的协议将在小程序内公布。如您继续使用本服务，视为接受修改后的协议。

九、联系我们
如您对本协议有任何疑问，可通过小程序内的反馈功能联系我们。

感谢您使用"情侣菜谱"小程序！`

  // 显示弹窗
  showAgreementPopup.value = true
}

// 查看隐私政策
const viewPrivacyPolicy = () => {
  agreementTitle.value = '隐私政策'
  agreementContent.value = `生效日期：2025年1月1日

"情侣菜谱"小程序（以下简称"我们"）非常重视您的隐私保护。本隐私政策说明了我们如何收集、使用、存储和保护您的个人信息。

一、信息收集
我们可能收集以下信息：
1. 账户信息：微信昵称、头像、用户ID等基础信息
2. 使用信息：您在使用服务过程中产生的日志、菜谱、打卡记录等
3. 设备信息：设备型号、操作系统版本等（用于优化服务体验）

二、信息使用
我们使用收集的信息用于：
1. 提供服务：为您提供日志记录、菜谱管理、打卡等功能
2. 改善服务：分析使用情况，优化产品功能
3. 安全保障：检测、防范安全风险，保护账户安全
4. 法律合规：遵守相关法律法规要求

三、信息存储
1. 您的信息将存储在安全的服务器上，我们会采取合理的技术手段保护您的信息安全。
2. 信息存储期限：在您使用服务期间及法律法规要求的期限内保存。

四、信息分享
我们承诺：
1. 不会出售您的个人信息给第三方
2. 不会未经授权向第三方提供您的个人信息
3. 仅在以下情况下可能分享信息：
   - 获得您的明确同意
   - 法律法规要求
   - 保护我们的合法权益

五、您的权利
您有权：
1. 访问：查看我们收集的您的个人信息
2. 更正：修改不准确的个人信息
3. 删除：要求删除您的个人信息（法律法规要求保留的除外）
4. 撤回同意：撤回您对信息处理的同意

六、未成年人保护
我们非常重视未成年人的个人信息保护。如果您是未成年人，请在监护人同意和指导下使用本服务。

七、隐私政策更新
我们可能适时更新本隐私政策，更新后的政策将在小程序内公布。如您继续使用本服务，视为接受更新后的政策。

八、联系我们
如您对本隐私政策有任何疑问或需要行使相关权利，可通过小程序内的反馈功能联系我们。

我们致力于保护您的隐私安全！`

  // 显示弹窗
  showAgreementPopup.value = true
}

// 关闭协议弹窗
const closeAgreement = () => {
  console.log('关闭协议弹窗')
  showAgreementPopup.value = false
  // 弹窗关闭时重置状态
  agreementTitle.value = ''
  agreementContent.value = ''
}

// 切换协议勾选状态
const toggleAgreement = () => {
  agreedToTerms.value = !agreedToTerms.value
}

// 方法
const onGetUserInfo = () => {
  console.log('点击登录按钮')
  
  // 检查是否勾选了协议
  if (!agreedToTerms.value) {
    uni.showToast({
      title: '请先同意用户协议和隐私政策',
      icon: 'none',
      duration: 2000
    })
    return
  }
  
  // 用户已勾选协议，继续登录流程
  doLogin()
}

// 执行登录流程
const doLogin = async () => {
  console.log('用户点击登录按钮，开始登录流程')
  
  try {
    isLogging.value = true
    
    // 获取微信登录code
    const loginRes = await uni.login({
      provider: 'weixin'
    })
    
    console.log('微信登录结果:', loginRes)
    
    if (loginRes.errMsg !== 'login:ok') {
      throw new Error('微信登录失败')
    }

    // 调用登录接口，同时透传邀请码（如果有）
    // 优先使用响应式变量，其次从本地存储读取
    const cachedInviteCode = inviteCode.value || uni.getStorageSync('INVITE_CODE') || ''
    console.log('准备登录，邀请码:', cachedInviteCode)
    
    const loginResult = await api.login({
      code: loginRes.code,
      inivCode: cachedInviteCode || undefined // 如果没有邀请码，不传该字段
    })

    
    // 检查登录结果 - 处理可能被错误包装的成功响应
    if (!loginResult) {
      throw new Error('登录接口返回数据为空')
    }

    // 如果loginResult本身就是一个包含成功数据的对象，直接使用
    if (loginResult.code === 0) {
      console.log('登录成功，开始处理用户信息')
      
      // 检查必要的数据字段
      if (!loginResult.data || !loginResult.data.token || !loginResult.data.Info) {
        throw new Error('登录返回数据格式错误')
      }
      
      // 登录成功，保存token到本地缓存
      const token = loginResult.data.token
      console.log('获取到的token:', token)
      uni.setStorageSync('token', token)
      
      userStore.setToken(token)

      uni.showToast({
        title: '登录成功',
        icon: 'success'
      })

      // 判断是否首次登录
      if (loginResult.data.isFirstLogin) {
        // 首次登录，跳转到设置页面
        // 注意：首次登录时，邀请码绑定应该在设置完成后进行，因为此时用户可能还没有完善信息
        // 邀请码会保留在本地存储中，等待设置完成后绑定
        console.log('首次登录，邀请码将在设置完成后绑定，当前邀请码:', uni.getStorageSync('INVITE_CODE'))
        setTimeout(() => {
          uni.navigateTo({
            url: '/pages/user/setup'
          })
        }, 1500)
      } else {
        // 非首次登录，立即尝试绑定情侣
        await tryBindLoverByInvite()
        // 非首次登录，获取最新用户信息并缓存
        try {
          console.log('非首次登录，开始获取最新用户信息...')
          const userInfoResult = await api.getUserInfo()
          
          if (userInfoResult.code === 0 && userInfoResult.data) {
            console.log('获取用户信息成功:', userInfoResult.data)
            
            // 缓存用户信息到userInfo
            uni.setStorageSync('userInfo', userInfoResult.data)
            
            // 更新用户store
            userStore.setUserInfo({
              userId: userInfoResult.data.ID || userInfoResult.data.userId,
              nickname: userInfoResult.data.nick_name || userInfoResult.data.nickname,
              avatar: userInfoResult.data.avatar,
              phone: userInfoResult.data.phone,
              iniv_code: userInfoResult.data.iniv_code
            })
            
            console.log('用户信息已缓存到userInfo:', userInfoResult.data)
          } else {
            console.warn('获取用户信息失败，使用登录返回的基础信息')
            // 如果获取失败，使用登录返回的基础信息
            const basicUserInfo = {
              ID: loginResult.data.Info.ID,
              nick_name: loginResult.data.Info.nick_name,
              avatar: loginResult.data.Info.avatar,
              phone: loginResult.data.Info.phone
            }
            uni.setStorageSync('userInfo', basicUserInfo)
          }
        } catch (error) {
          console.error('获取用户信息失败:', error)
          // 如果获取失败，使用登录返回的基础信息
          const basicUserInfo = {
            ID: loginResult.data.Info.ID,
            nick_name: loginResult.data.Info.nick_name,
            avatar: loginResult.data.Info.avatar,
            phone: loginResult.data.Info.phone
          }
          uni.setStorageSync('userInfo', basicUserInfo)
        }
        
        // 跳转到首页
        setTimeout(() => {
          uni.switchTab({
            url: '/pages/home/home'
          })
        }, 1500)
      }
    } else {
      console.log('登录失败，错误信息:', loginResult?.msg)
      throw new Error(loginResult?.msg || '登录失败')
    }

  } catch (error) {
    // 检查是否是被错误包装的成功响应
    if (error && typeof error === 'object' && error.code === 0) {
      console.log('🔍 检测到被错误包装的成功响应，正在重新处理...')
      
      try {
        // 重新处理成功响应
        const loginResult = error
        
        if (loginResult.data && loginResult.data.token && loginResult.data.Info) {
          console.log('✅ 成功解析登录数据，开始处理用户信息')
          
          // 登录成功，保存token到本地缓存
          const token = loginResult.data.token
          console.log('🔑 获取到的token:', token)
          uni.setStorageSync('token', token)
        
          userStore.setToken(token)

          uni.showToast({
            title: '登录成功',
            icon: 'success'
          })

          // 判断是否首次登录
          if (loginResult.data.isFirstLogin) {
            console.log('🎯 首次登录，准备跳转设置页面')
            // 首次登录时，邀请码绑定应该在设置完成后进行
            console.log('首次登录，邀请码将在设置完成后绑定，当前邀请码:', uni.getStorageSync('INVITE_CODE'))
            setTimeout(() => {
              uni.navigateTo({
                url: '/pages/user/setup'
              })
            }, 1500)
          } else {
            // 非首次登录，立即尝试绑定情侣
            await tryBindLoverByInvite()
            console.log('🏠 非首次登录，准备获取最新用户信息并跳转首页')
            
            // 非首次登录，获取最新用户信息并缓存
            try {
              const userInfoResult = await api.getUserInfo()
              
              if (userInfoResult.code === 0 && userInfoResult.data) {
                console.log('✅ 获取用户信息成功:', userInfoResult.data)
                
                // 缓存用户信息到userInfo
                uni.setStorageSync('userInfo', userInfoResult.data)
                
                // 更新用户store
                userStore.setUserInfo({
                  userId: userInfoResult.data.ID || userInfoResult.data.userId,
                  nickname: userInfoResult.data.nick_name || userInfoResult.data.nickname,
                  avatar: userInfoResult.data.avatar,
                  phone: userInfoResult.data.phone,
                  iniv_code: userInfoResult.data.iniv_code
                })
                
                console.log('🎉 用户信息已缓存到userInfo:', userInfoResult.data)
              } else {
                console.warn('⚠️ 获取用户信息失败，使用登录返回的基础信息')
                // 如果获取失败，使用登录返回的基础信息
                const basicUserInfo = {
                  ID: loginResult.data.Info.ID,
                  nick_name: loginResult.data.Info.nick_name,
                  avatar: loginResult.data.Info.avatar,
                  phone: loginResult.data.Info.phone
                }
                uni.setStorageSync('userInfo', basicUserInfo)
              }
            } catch (error) {
              console.error('❌ 获取用户信息失败:', error)
              // 如果获取失败，使用登录返回的基础信息
              const basicUserInfo = {
                ID: loginResult.data.Info.ID,
                nick_name: loginResult.data.Info.nick_name,
                avatar: loginResult.data.Info.avatar,
                phone: loginResult.data.Info.phone
              }
              uni.setStorageSync('userInfo', basicUserInfo)
            }
            
            // 跳转到首页
            setTimeout(() => {
              uni.switchTab({
                url: '/pages/home/home'
              })
            }, 1500)
          }
          
          // 标记为已处理，避免后续错误日志
          error._handled = true
          console.log('🎉 登录流程处理完成')
          return // 成功处理，退出函数
        } else {
          console.warn('⚠️ 登录数据格式不完整，无法处理')
        }
      } catch (innerError) {
        console.error('❌ 重新处理成功响应时出错:', innerError)
      }
    }
    
    // 真正的错误处理 - 只处理未被标记的错误
    if (!error._handled) {
      console.error('❌ 登录遇到真实错误:', error)
      console.error('🚨 错误类型:', error?.constructor?.name || typeof error)
      console.error('💬 错误对象:', JSON.stringify(error, null, 2))
      
      // 提取错误信息：优先使用 error.message，其次 error.msg，最后 error.errMsg
      let errorMsg = ''
      if (error && typeof error === 'object') {
        errorMsg = error.message || error.msg || error.errMsg || ''
      } else if (typeof error === 'string') {
        errorMsg = error
      }
      
      // 如果是 API 返回的错误响应，尝试提取 msg
      if (!errorMsg && error?.data) {
        errorMsg = error.data.msg || error.data.message || ''
      }
      
      // 如果还是没有错误信息，使用默认提示
      if (!errorMsg || errorMsg === '登录失败，请重试') {
        errorMsg = '登录失败，请重试'
      }
      
      // 限制错误信息长度，避免显示过长
      if (errorMsg.length > 30) {
        errorMsg = errorMsg.substring(0, 30) + '...'
      }
      
      uni.showToast({
        title: errorMsg,
        icon: 'none',
        duration: 3000
      })
    }
  } finally {
    isLogging.value = false
  }
}


// 暂不登录，返回首页
const skipLogin = () => {
  uni.switchTab({
    url: '/pages/home/home'
  })
}
</script>

<style lang="scss" scoped>
.login-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #FFE4E1 0%, #FFB6C1 50%, #E6E6FA 100%);
  position: relative;
  overflow: hidden;
}

.bg-decoration {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  pointer-events: none;
  
  .circle {
    position: absolute;
    border-radius: 50%;
    background: rgba(255, 255, 255, 0.1);
    
    &.circle-1 {
      width: 200rpx;
      height: 200rpx;
      top: 10%;
      right: -50rpx;
      animation: float 6s ease-in-out infinite;
    }
    
    &.circle-2 {
      width: 150rpx;
      height: 150rpx;
      top: 30%;
      left: -30rpx;
      animation: float 8s ease-in-out infinite reverse;
    }
    
    &.circle-3 {
      width: 100rpx;
      height: 100rpx;
      top: 60%;
      right: 20%;
      animation: float 10s ease-in-out infinite;
    }
  }
}

.main-content {
  position: relative;
  z-index: 1;
  padding: 100rpx 60rpx 60rpx;
  text-align: center;
}

.header-section {
  margin-bottom: 80rpx;
  
  .app-logo {
    margin-bottom: 30rpx;
    animation: bounce 2s ease-in-out infinite;
  }
  
  .app-title {
    display: block;
    font-size: 48rpx;
    font-weight: bold;
    color: #333;
    margin-bottom: 20rpx;
  }
  
  .app-subtitle {
    font-size: 28rpx;
    color: #666;
    line-height: 1.4;
  }
}

.features-section {
  margin-bottom: 80rpx;
  
  .feature-item {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 20rpx;
    margin-bottom: 30rpx;
    
    .feature-text {
      font-size: 28rpx;
      color: #333;
    }
  }
}

  .login-section {
    .login-tips {
      margin-bottom: 40rpx;
      
      .tips-text {
        font-size: 26rpx;
        color: #999;
      }
    }
    
    .skip-login-section {
      margin-top: 30rpx;
      text-align: center;
      
      .skip-login-btn {
        font-size: 28rpx;
        color: #999;
        text-decoration: underline;
        padding: 20rpx;
        display: inline-block;
      }
    }
    
    .wechat-login-btn {
    width: 100%;
    height: 100rpx;
    background: linear-gradient(135deg, #07c160 0%, #10ad6a 100%);
    border: none;
    border-radius: 50rpx;
    margin-bottom: 40rpx;
    box-shadow: 0 8rpx 30rpx rgba(7, 193, 96, 0.3);
    transition: all 0.3s;
    
    &:active {
      transform: translateY(2rpx);
      box-shadow: 0 4rpx 15rpx rgba(7, 193, 96, 0.3);
    }
    
    .btn-content {
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 15rpx;
      height: 100%;
      
      .btn-text {
        color: #fff;
        font-size: 32rpx;
        font-weight: bold;
      }
    }
  }
  
  .privacy-tips {
    text-align: center;
    line-height: 1.6;
    margin-top: 30rpx;
    
    .agreement-checkbox {
      display: flex;
      align-items: center;
      justify-content: center;
      flex-wrap: wrap;
      gap: 8rpx;
      
      checkbox {
        transform: scale(0.9);
      }
    }
    
    .privacy-text {
      font-size: 24rpx;
      color: #999;
    }
    
    .privacy-link {
      font-size: 24rpx;
      color: #FF6B9D;
      text-decoration: underline;
    }
  }
}

// 弹窗遮罩层样式
.popup-mask {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 999;
  display: flex;
  align-items: center;
  justify-content: center;

  // 协议弹窗样式
  .agreement-popup {
    width: 680rpx !important;
    max-width: 680rpx !important;
    max-height: 80vh !important;
    background: #fff;
    border-radius: 20rpx;
    overflow: hidden;
    display: flex;
    flex-direction: column;
    position: relative;
    z-index: 1000;
    box-sizing: border-box;
    
    .popup-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 30rpx;
      border-bottom: 1rpx solid #f0f0f0;
      
      .popup-title {
        font-size: 32rpx;
        font-weight: bold;
        color: #333;
      }
      
      .close-btn {
        width: 40rpx;
        height: 40rpx;
        display: flex;
        align-items: center;
        justify-content: center;
        cursor: pointer;
      }
    }
    
    .popup-content {
      flex: 1;
      padding: 0;
      max-height: 50vh;
      min-height: 200rpx;
      
      .agreement-text-wrapper {
        padding: 30rpx;
        
        .agreement-text {
          font-size: 26rpx;
          line-height: 1.8;
          color: #333;
          white-space: pre-wrap;
          word-break: break-word;
        }
      }
    }
    
    .popup-footer {
      padding: 20rpx 30rpx 30rpx;
      border-top: 1rpx solid #f0f0f0;
      
      .popup-confirm-btn {
        width: 100%;
        height: 80rpx;
        background: linear-gradient(135deg, #FF6B9D 0%, #FF8E9E 100%);
        color: #fff;
        border: none;
        border-radius: 40rpx;
        font-size: 28rpx;
        font-weight: 500;
      }
    }
  }
}

.bottom-decoration {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 200rpx;
  pointer-events: none;
  
  .wave {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    height: 100rpx;
    background: rgba(255, 255, 255, 0.1);
    border-radius: 50% 50% 0 0;
    
    &.wave-1 {
      bottom: 0;
      animation: wave 4s ease-in-out infinite;
    }
    
    &.wave-2 {
      bottom: 20rpx;
      animation: wave 6s ease-in-out infinite reverse;
    }
  }
}

@keyframes float {
  0%, 100% { transform: translateY(0px); }
  50% { transform: translateY(-20px); }
}

@keyframes bounce {
  0%, 20%, 50%, 80%, 100% { transform: translateY(0); }
  40% { transform: translateY(-10px); }
  60% { transform: translateY(-5px); }
}

@keyframes wave {
  0%, 100% { transform: scaleX(1); }
  50% { transform: scaleX(1.1); }
}
</style>
