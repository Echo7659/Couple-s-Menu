<template>
  <view class="home-container">
    <!-- 自定义头部 -->
    <view class="custom-header">
      <view class="status-bar"></view>
      <view class="header-content">
        <view class="header-title">情侣菜谱</view>
        <view class="header-actions">
          <uni-icons type="more-filled" size="20" color="#fff"></uni-icons>
        </view>
      </view>
    </view>

    <!-- 主要内容区域 -->
    <view class="main-content">
      <!-- 情侣头像和恋爱天数 -->
      <view class="couple-info">
        <view class="avatar-section">
          <view class="avatar-container left-avatar">
            <view class="avatar">
              <img
                :src="userInfo.avatar"
                class="avatar-img"
                :key="userInfo.avatar"
                alt="用户头像"
              ></img>
            </view>
            <view class="nickname">
              <text class="nickname-text">{{ userInfo.nickname || '我' }}</text>
              <text 
                v-if="userInfo.role" 
                class="role-badge" 
                :class="getRoleClass(userInfo.role)"
                @click.stop="switchRole"
              >
                {{ getRoleText(userInfo.role) }}
              </text>
            </view>
          </view>
          
          <view class="heart-icon">
            <uni-icons type="heart-filled" size="30" color="#FF6B9D"></uni-icons>
          </view>
          
          <view class="avatar-container right-avatar">
            <view class="avatar">
              <!-- 未绑定情侣时，右侧头像作为分享按钮，点击直接触发邀请分享 -->
              <button
                v-if="!hasPartner"
                class="avatar-share-btn"
                open-type="share"
                :data-shareinfo="shareInfo"
              >
                <view class="avatar-img placeholder"></view>
                <view class="add-partner">
                  <uni-icons type="plus" size="24" color="#FF6B9D"></uni-icons>
                </view>
              </button>
              <!-- 已绑定情侣时正常展示对方头像 -->
              <img 
                v-else
                :src="partnerInfo.avatar" 
                class="avatar-img"
                alt="伴侣头像"
              ></img>
            </view>
            <view class="nickname">
              <text class="nickname-text">{{ partnerInfo.nickname || (hasPartner ? 'Ta' : '邀请Ta') }}</text>
            </view>
          </view>
        </view>

        <view class="relationship-info" @click="handleUpdateStartDate" v-if="hasPartner">
          <text class="relationship-text">我们在一起已经</text>
          <view class="days-count">
            <text class="days-number">{{ loveDays }}</text>
            <text class="days-unit">天</text>
          </view>
          <text class="start-date">{{ startDate }}</text>
        </view>
        <view class="relationship-info" v-else>
          <text class="relationship-text">我们在一起已经</text>
          <view class="days-count">
            <text class="days-number">{{ loveDays }}</text>
            <text class="days-unit">天</text>
          </view>
          <text class="start-date">{{ startDate }}</text>
        </view>
      </view>

      <!-- 快捷功能卡片 -->
      <view class="feature-cards">
        <view class="feature-card diary-card" @click="goToTarget">
          <view class="card-content">
            <text class="card-title">去打卡</text>
            <uni-icons type="flag" size="32" color="#fff"></uni-icons>
          </view>
        </view>
        <view class="feature-card order-card" @click="goToOrder">
          <view class="card-content">
            <text class="card-title">点餐下单</text>
            <uni-icons type="cart" size="32" color="#fff"></uni-icons>
          </view>
        </view>
      </view>

      <!-- 纪念日列表 -->
      <view class="anniversary-section">
        <view class="section-header">
          <text class="section-title">重要纪念日</text>
          <view class="add-anniversary" @click="addAnniversary">
            <uni-icons type="plus" size="20" color="#FF6B9D"></uni-icons>
          </view>
        </view>

        <view class="anniversary-list">
          <view
            v-for="(item, index) in anniversaryList"
            :key="item.ID || index"
            class="anniversary-item"
            @click="viewAnniversary(item)"
          >
            <view class="anniversary-icon no-circle" :data-type="item.type">
              <view class="icon-wrapper">
                <!-- 生日类型显示图片，其他类型显示图标 -->
                <img 
                  v-if="item.type === 2" 
                  src="/static/images/birthday.png" 
                  class="birthday-icon"
                  alt="生日"
                />
                <uni-icons 
                  v-else
                  :type="getIconByType(item.type)" 
                  size="26" 
                  :color="getIconColorByType(item.type)"
                ></uni-icons>
              </view>
            </view>
            <view class="anniversary-info">
              <text class="anniversary-name">{{ item.title }}</text>
              <view class="anniversary-meta">
                <!-- <text class="anniversary-type" :data-type="item.type">{{ getTypeTextByType(item.type) }}</text> -->
                <text class="anniversary-date">{{ item.date }}</text>
              </view>
            </view>
            <view class="anniversary-countdown">
              <text class="countdown-text">{{ getCountdownText(item) }}</text>
            </view>
          </view>

          <!-- 空状态 -->
          <view v-if="anniversaryList.length === 0" class="empty-anniversary">
            <text class="empty-text">还没有添加纪念日</text>
            <text class="empty-desc">点击右上角添加第一个纪念日吧</text>
          </view>
        </view>
      </view>
    </view>
    
    <!-- 日期选择器弹窗 -->
    <uni-popup ref="datePickerPopup" :show="showDatePicker" type="bottom" @change="onPopupChange">
      <view class="date-picker-popup">
        <view class="popup-header">
          <text class="popup-title">选择恋爱开始日期</text>
          <view class="popup-actions">
            <text class="popup-btn cancel-btn" @click="cancelUpdateDate">取消</text>
            <text class="popup-btn confirm-btn" @click="confirmUpdateDate">确定</text>
          </view>
        </view>
        <view class="popup-content">
          <picker
            mode="date"
            :value="selectedDate"
            @change="onDateChange"
            :start="'1900-01-01'"
            :end="new Date().toISOString().split('T')[0]"
          >
            <view class="picker-view">
              <text class="picker-text">{{ selectedDate || '请选择日期' }}</text>
              <uni-icons type="calendar" size="20" color="#FF6B9D"></uni-icons>
            </view>
          </picker>
        </view>
      </view>
    </uni-popup>
  </view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useUserStore } from '@/store/user.js'
import api from '@/api'
import { onPullDownRefresh, onShow, onShareAppMessage } from '@dcloudio/uni-app'

const userStore = useUserStore()

// 响应式数据
const userInfo = ref({
  avatar: '',
  nickname: '',
  role: null // 1=饲养员，2=吃货
})

const partnerInfo = ref({
  avatar: '',
  nickname: ''
})

// 是否已经绑定情侣
const hasPartner = ref(false)

const startDate = ref('')
const loveDays = ref(0)
const loverRecordId = ref(null) // 保存情侣关系记录的ID（用于更新开始日期）

// 日期选择器相关
const showDatePicker = ref(false)
const selectedDate = ref('')

const anniversaryList = ref([])

// 计算属性
const computedLoveDays = computed(() => loveDays.value)

// API: 初始化用户与情侣信息
const initHomeData = async () => {
  try {
    // 1. 优先从缓存拿本人信息（登录后已缓存到 userInfo）
    const cached = uni.getStorageSync('userInfo')
    if (cached) {
      userInfo.value.avatar = cached.avatar || ''
      userInfo.value.nickname = cached.nick_name || cached.nickname || ''
      userInfo.value.role = cached.role || null
    }

    // 2. 获取当前用户ID（用于判断谁是恋人）
    const currentUserId = cached?.ID || cached?.id || cached?.userId || null
    
    // 3. 调用情侣关系接口，拿到 loverId、userId、ID 与 startDate
    const loverRes = await api.getLoverInfo()
    if (loverRes?.code === 0 && loverRes.data) {
      const { loverId, userId, startDate: s, ID, id } = loverRes.data
      
      // 保存情侣关系记录的ID（用于更新开始日期）
      loverRecordId.value = ID || id || null
      
      // 判断真正的恋人ID
      // 如果数据库中只有一条记录 userId:1, loverId:2
      // 当用户1打开时：返回的 userId=1, loverId=2，应该取 loverId=2
      // 当用户2打开时：返回的 userId=1, loverId=2，但当前用户是2，应该取 userId=1
      let actualLoverId = null
      if (loverId && userId) {
        // 如果当前用户ID等于返回的loverId，说明当前用户是被绑定的一方，应该取userId作为恋人ID
        if (currentUserId && String(currentUserId) === String(loverId)) {
          actualLoverId = userId
          console.log('当前用户是被绑定方，使用 userId 作为恋人ID:', actualLoverId)
        } 
        // 如果当前用户ID等于返回的userId，说明当前用户是主动绑定的一方，应该取loverId作为恋人ID
        else if (currentUserId && String(currentUserId) === String(userId)) {
          actualLoverId = loverId
          console.log('当前用户是主动绑定方，使用 loverId 作为恋人ID:', actualLoverId)
        }
        // 如果没有当前用户ID，默认使用loverId（兼容旧逻辑）
        else {
          actualLoverId = loverId
          console.log('无法判断用户身份，默认使用 loverId:', actualLoverId)
        }
      } else if (loverId) {
        // 如果接口只返回loverId，使用loverId
        actualLoverId = loverId
        console.log('接口只返回 loverId，使用:', actualLoverId)
      }
      
      hasPartner.value = !!actualLoverId
      // 设置开始日期（用于恋爱天数）
      startDate.value = s || ''

      // 获取最新的用户信息（包括邀请码）
      // 即使缓存中有信息，也要调用一次以确保邀请码是最新的
      const meRes = await api.getUserInfo()
      if (meRes?.code === 0 && meRes.data) {
        userInfo.value.avatar = meRes.data.avatar || userInfo.value.avatar
        userInfo.value.nickname = meRes.data.nick_name || meRes.data.nickname || userInfo.value.nickname
        userInfo.value.role = meRes.data.role || userInfo.value.role

        // 更新缓存（包括邀请码）
        uni.setStorageSync('userInfo', meRes.data)
        console.log('已更新用户信息到缓存，邀请码:', meRes.data.iniv_code)
      }

      // 4. 根据实际恋人ID获取另一半资料
      if (actualLoverId) {
        console.log('开始获取恋人信息，恋人ID:', actualLoverId)
        const loverInfoRes = await api.getUserInfoById(actualLoverId)
        if (loverInfoRes?.code === 0 && loverInfoRes.data) {
          partnerInfo.value.avatar = loverInfoRes.data.avatar || ''
          partnerInfo.value.nickname = loverInfoRes.data.nick_name || loverInfoRes.data.nickname || ''
          console.log('获取恋人信息成功:', partnerInfo.value)
        } else {
          console.warn('获取恋人信息失败:', loverInfoRes)
        }
      } else {
        console.log('没有有效的恋人ID，跳过获取恋人信息')
      }

      // 5. 计算恋爱天数
      computeLoveDays(startDate.value)
    }

    // 5. 加载纪念日列表
    await loadAnniversaryList()
  } catch (e) {
    console.error('初始化首页数据失败:', e)
  }
}

const computeLoveDays = (dateStr) => {
  if (!dateStr) {
    loveDays.value = 0
    return
  }
  // 支持 yyyy-MM-dd 或 yyyy/MM/dd
  const norm = dateStr.replace(/\./g, '-').replace(/\//g, '-')
  const start = new Date(norm)
  if (isNaN(start.getTime())) {
    loveDays.value = 0
    return
  }
  const today = new Date()
  const diffTime = today.getTime() - start.getTime()
  loveDays.value = Math.max(0, Math.ceil(diffTime / (1000 * 60 * 60 * 24)))
}

// 加载纪念日列表
const loadAnniversaryList = async () => {
  try {
    const cached = uni.getStorageSync('userInfo')
    const userId = cached?.ID
    const res = await api.getMemorialDayList({ userId })
    if (res?.code === 0 && res.data?.list) {
      anniversaryList.value = res.data.list || []
      // 按日期排序，只显示最近的几个
      anniversaryList.value.sort((a, b) => {
        const dateA = new Date(a.date)
        const dateB = new Date(b.date)
        return dateA - dateB
      })
      // 只显示前3个
      anniversaryList.value = anniversaryList.value.slice(0, 3)
    }
  } catch (error) {
    console.error('加载纪念日列表失败:', error)
  }
}

// 获取纪念日倒计时文本
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
  return `还有${days}天`
}

// 获取纪念日图标
const getIconByType = (type) => {
  const iconMap = {
    1: 'heart-filled',    // 恋爱纪念日 - 使用 heart-filled
    2: 'star-filled',     // 生日 - 使用 star-filled
    3: 'calendar-filled'  // 其他 - 使用 calendar-filled
  }
  return iconMap[type] || 'calendar-filled'
}

// 获取纪念日图标颜色
const getIconColorByType = (type) => {
  const colorMap = {
    1: '#FF6B9D',         // 恋爱纪念日 - 粉色
    2: '#FFD700',         // 生日 - 金色
    3: '#4CAF50'          // 其他 - 绿色
  }
  return colorMap[type] || '#FF6B9D'
}

// 获取纪念日类型文字
const getTypeTextByType = (type) => {
  const textMap = {
    1: '恋爱纪念日',
    2: '生日',
    3: '其他'
  }
  return textMap[type] || '其他'
}

// 获取角色文字
const getRoleText = (role) => {
  const roleMap = {
    1: '饲养员',
    2: '吃货'
  }
  return roleMap[role] || ''
}

// 获取角色样式类
const getRoleClass = (role) => {
  return role === 1 ? 'role-feeder' : 'role-foodie'
}

// 切换角色
const switchRole = async () => {
  if (!userStore.hasLogin) {
    uni.showToast({
      title: '请先登录',
      icon: 'none'
    })
    return
  }
  
  const currentRole = userInfo.value.role
  if (!currentRole) {
    uni.showToast({
      title: '角色信息获取失败',
      icon: 'none'
    })
    return
  }
  
  // 切换角色：1变2，2变1
  const newRole = currentRole === 1 ? 2 : 1
  const newRoleText = getRoleText(newRole)
  
  // 先弹出确认提示
  uni.showModal({
    title: '切换角色',
    content: `确定要切换为${newRoleText}角色吗？`,
    confirmText: '确定',
    cancelText: '取消',
    success: async (res) => {
      if (res.confirm) {
        // 用户确认后执行切换
        try {
          uni.showLoading({ title: '切换中...' })
          const result = await api.updateUserRole({ role: newRole })
          if (result?.code === 0) {
            userInfo.value.role = newRole
            // 更新缓存
            const cached = uni.getStorageSync('userInfo') || {}
            cached.role = newRole
            uni.setStorageSync('userInfo', cached)
            
            uni.showToast({
              title: `已切换为${newRoleText}`,
              icon: 'success'
            })
          } else {
            uni.showToast({
              title: result?.msg || '切换失败',
              icon: 'none'
            })
          }
        } catch (error) {
          console.error('切换角色失败:', error)
          uni.showToast({
            title: '切换失败',
            icon: 'none'
          })
        } finally {
          uni.hideLoading()
        }
      }
    }
  })
}

// 方法
const invitePartner = () => {
  // 由于小程序不允许直接用代码弹出分享面板，这里主要作为兜底提示
  uni.showToast({
    title: '点击右侧头像即可邀请Ta一起使用~',
    icon: 'none'
  })
}

// 处理更新恋爱开始日期
const handleUpdateStartDate = () => {
  if (!hasPartner.value) {
    return
  }
  
  if (!loverRecordId.value) {
    uni.showToast({
      title: '无法获取关系记录ID',
      icon: 'none'
    })
    return
  }
  
  // 如果有当前日期，设置为默认值
  if (startDate.value) {
    selectedDate.value = startDate.value
  } else {
    // 默认选择今天
    const today = new Date()
    const year = today.getFullYear()
    const month = String(today.getMonth() + 1).padStart(2, '0')
    const day = String(today.getDate()).padStart(2, '0')
    selectedDate.value = `${year}-${month}-${day}`
  }
  
  // 显示日期选择器弹窗
  if (datePickerPopup.value) {
    datePickerPopup.value.open()
  } else {
    showDatePicker.value = true
  }
}

// 日期选择器变化事件
const onDateChange = (e) => {
  selectedDate.value = e.detail.value
}

// 确认更新日期
const confirmUpdateDate = async () => {
  if (!selectedDate.value) {
    uni.showToast({
      title: '请选择日期',
      icon: 'none'
    })
    return
  }
  
  // 验证日期格式
  const dateRegex = /^\d{4}-\d{2}-\d{2}$/
  if (!dateRegex.test(selectedDate.value)) {
    uni.showToast({
      title: '日期格式不正确',
      icon: 'none'
    })
    return
  }
  
  try {
    uni.showLoading({ title: '更新中...' })
    
    const result = await api.updateLoverDay({
      ID: loverRecordId.value,
      startDate: selectedDate.value
    })
    
    if (result?.code === 0) {
      // 更新成功，刷新数据
      startDate.value = selectedDate.value
      computeLoveDays(selectedDate.value)
      
      uni.showToast({
        title: '更新成功',
        icon: 'success'
      })
      
      // 关闭弹窗
      if (datePickerPopup.value) {
        datePickerPopup.value.close()
      } else {
        showDatePicker.value = false
      }
    } else {
      uni.showToast({
        title: result?.msg || '更新失败',
        icon: 'none'
      })
    }
  } catch (error) {
    console.error('更新恋爱开始日期失败:', error)
    uni.showToast({
      title: '更新失败',
      icon: 'none'
    })
  } finally {
    uni.hideLoading()
  }
}

// 取消更新日期
const cancelUpdateDate = () => {
  if (datePickerPopup.value) {
    datePickerPopup.value.close()
  } else {
    showDatePicker.value = false
  }
  selectedDate.value = ''
}

// 弹窗状态变化
const onPopupChange = (e) => {
  showDatePicker.value = e.show
  if (!e.show) {
    selectedDate.value = ''
  }
}

// 日期选择器引用
const datePickerPopup = ref(null)

// 通用登录校验：未登录时跳转到登录页
const ensureLogin = () => {
  if (!userStore.hasLogin) {
    uni.navigateTo({ url: '/pages/user/login' })
    return false
  }
  return true
}

// 读取当前用户的邀请码（从缓存的 userInfo 中获取）
const inviteCode = computed(() => {
  const cached = uni.getStorageSync('userInfo') || {}
  return cached.iniv_code || ''
})

// 分享信息（使用 computed 确保邀请码实时更新）
// 注意：computed 值在绑定到 data 属性时会自动获取其值
const shareInfo = computed(() => {
  // 从缓存中获取邀请码（确保实时获取）
  const cached = uni.getStorageSync('userInfo') || {}
  const code = cached.iniv_code || ''
  
  console.log('分享时邀请码:', code)
  
  // 小程序分享路径：必须以 / 开头，指向 pages.json 中配置的页面路径
  // 注意：不能分享到 tabBar 页面，所以这里分享到登录页
  const sharePath = code ? `/pages/user/login?invite=${encodeURIComponent(code)}` : '/pages/user/login'
  const shareTitle = '邀请你一起用情侣菜谱，记录我们的每一餐~'
  
  console.log('分享配置:', JSON.stringify({ title: shareTitle, path: sharePath }, null, 2))
  
  // 返回普通对象，确保能正确序列化到 data 属性
  return {
    title: shareTitle,
    path: sharePath,
    imageUrl: '', // 分享图片（空字符串使用默认）
    desc: '',
    content: ''
  }
})



const goToDiary = () => {
  if (!ensureLogin()) return
  uni.switchTab({ url: '/pages/diary/diary' })
}

const goToTarget = () => {
  if (!ensureLogin()) return
  uni.switchTab({ url: '/pages/target/index' })
}
const goToOrder = () => {
  if (!ensureLogin()) return
  uni.switchTab({ url: '/pages/order/index' })
}

const addAnniversary = () => {
  if (!ensureLogin()) return
  uni.navigateTo({
    url: '/pages/anniversary/anniversary'
  })
}

const viewAnniversary = (item) => {
  if (!ensureLogin()) return
  uni.navigateTo({
    url: '/pages/anniversary/anniversary'
  })
}

// 生命周期
onMounted(() => {
  initHomeData()
})

// 页面显示时刷新用户信息（从缓存或接口）
onShow(() => {
      // 从缓存读取最新的用户信息
      const cachedUserInfo = uni.getStorageSync('userInfo')
      if (cachedUserInfo) {
        userInfo.value = {
          avatar: cachedUserInfo.avatar || userInfo.value.avatar,
          nickname: cachedUserInfo.nick_name || cachedUserInfo.nickname || userInfo.value.nickname,
          role: cachedUserInfo.role || userInfo.value.role
        }
      }
  
  // 刷新用户信息（静默刷新，不显示loading）
  refreshUserInfo()
})

// 刷新用户信息（静默）
const refreshUserInfo = async () => {
  try {
    const result = await api.getUserInfo()
    if (result?.code === 0 && result.data) {
      userInfo.value = {
        avatar: result.data.avatar || userInfo.value.avatar,
        nickname: result.data.nick_name || result.data.nickname || userInfo.value.nickname,
        role: result.data.role || userInfo.value.role
      }
      // 更新缓存
      uni.setStorageSync('userInfo', result.data)
    }
  } catch (error) {
    console.error('刷新用户信息失败:', error)
    // 失败时使用缓存数据
    const cachedUserInfo = uni.getStorageSync('userInfo')
    if (cachedUserInfo) {
        userInfo.value = {
          avatar: cachedUserInfo.avatar || userInfo.value.avatar,
          nickname: cachedUserInfo.nick_name || cachedUserInfo.nickname || userInfo.value.nickname,
          role: cachedUserInfo.role || userInfo.value.role
        }
    }
  }
}

// 下拉刷新
onPullDownRefresh(() => {
  // 同时刷新用户信息和纪念日列表
  Promise.all([
    initHomeData(),
    loadAnniversaryList()
  ]).finally(() => {
    // 结束下拉刷新动画
    uni.stopPullDownRefresh()
  })
})

// 配置分享功能 - 每次分享时动态获取最新的邀请码
onShareAppMessage(async (options) => {
  console.log('触发分享，options:', options)

  try {
    // 从缓存获取当前用户ID
    const cached = uni.getStorageSync('userInfo') || {}
    const userId = cached.ID || cached.id || cached.userId

    if (!userId) {
      console.warn('无法获取用户ID，使用缓存的邀请码')
      const cachedInviteCode = cached.iniv_code || ''
      const sharePath = cachedInviteCode
        ? `/pages/user/login?invite=${encodeURIComponent(cachedInviteCode)}`
        : '/pages/user/login'

      return {
        title: '邀请你一起用情侣菜谱，记录我们的每一餐~',
        path: sharePath,
        imageUrl: ''
      }
    }

    // 动态获取最新的用户信息，确保邀请码是最新的
    console.log('开始获取最新用户信息，用户ID:', userId)
    const userInfoRes = await api.getUserInfo()

    let inviteCode = ''
    if (userInfoRes?.code === 0 && userInfoRes.data) {
      inviteCode = userInfoRes.data.iniv_code || ''
      console.log('获取到最新邀请码:', inviteCode)

      // 更新缓存中的邀请码
      cached.iniv_code = inviteCode
      uni.setStorageSync('userInfo', cached)
    } else {
      console.warn('获取用户信息失败，使用缓存的邀请码')
      inviteCode = cached.iniv_code || ''
    }

    // 构建分享路径
    const sharePath = inviteCode
      ? `/pages/user/login?invite=${encodeURIComponent(inviteCode)}`
      : '/pages/user/login'

    console.log('分享配置 - 标题: 邀请你一起用情侣菜谱，记录我们的每一餐~')
    console.log('分享配置 - 路径:', sharePath)
    console.log('分享配置 - 邀请码:', inviteCode)

    return {
      title: '邀请你一起用情侣菜谱，记录我们的每一餐~',
      path: sharePath,
      imageUrl: '' // 使用默认截图
    }
  } catch (error) {
    console.error('分享时获取邀请码失败:', error)

    // 出错时使用缓存的邀请码
    const cached = uni.getStorageSync('userInfo') || {}
    const cachedInviteCode = cached.iniv_code || ''
    const sharePath = cachedInviteCode
      ? `/pages/user/login?invite=${encodeURIComponent(cachedInviteCode)}`
      : '/pages/user/login'

    return {
      title: '邀请你一起用情侣菜谱，记录我们的每一餐~',
      path: sharePath,
      imageUrl: ''
    }
  }
})

</script>

<style lang="scss" scoped>
.home-container {
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

    .header-title {
      color: #fff;
      font-size: 36rpx;
      font-weight: bold;
    }

    .header-actions {
      display: flex;
      gap: 20rpx;
    }
  }
}

.main-content {
  padding: 40rpx 30rpx 60rpx;
}

.couple-info {
  text-align: center;
  margin-bottom: 60rpx;

  .avatar-section {
    display: flex;
    justify-content: center;
    align-items: center;
    margin-bottom: 40rpx;
    gap: 40rpx;

    .avatar-container {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 15rpx;

      .avatar {
        position: relative;
        
          .avatar-share-btn {
            width: 160rpx;
            height: 160rpx;
            border-radius: 50%;
            padding: 0;
            margin: 0;
            border: none;
            // background: transparent;
            display: block;
            overflow: hidden;
            
            &::after {
              border: none;
            }
          }

        .avatar-img {
          width: 160rpx;
          height: 160rpx;
          border-radius: 50%;
          border: 4rpx solid #fff;
          box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.15);
          object-fit: cover;
          background: #f0f0f0;
          transition: all 0.3s ease;
          cursor: pointer;
          
          &:hover {
            transform: scale(1.05);
            box-shadow: 0 6rpx 25rpx rgba(0, 0, 0, 0.2);
          }

            &.placeholder {
              background: transparent;
              border: 4rpx dashed rgba(255, 255, 255, 0.6);
            }
        }

        .add-partner {
          position: absolute;
          top: 50%;
          left: 50%;
          transform: translate(-50%, -50%);
          background: rgba(255, 255, 255, 0.9);
          border-radius: 50%;
          width: 40rpx;
          height: 40rpx;
          display: flex;
          align-items: center;
          justify-content: center;
          box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.1);
        }
      }

      .nickname {
        text-align: center;
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 8rpx;
        
        .nickname-text {
          font-size: 28rpx;
          font-weight: 600;
          color: #fff;
          text-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.3);
          background: linear-gradient(135deg, rgba(255, 255, 255, 0.2), rgba(255, 255, 255, 0.1));
          padding: 8rpx 20rpx;
          border-radius: 20rpx;
          backdrop-filter: blur(10rpx);
          border: 1rpx solid rgba(255, 255, 255, 0.3);
          box-shadow: 0 4rpx 15rpx rgba(0, 0, 0, 0.1);
          transition: all 0.3s ease;
          
          &:hover {
            transform: translateY(-2rpx);
            box-shadow: 0 6rpx 20rpx rgba(0, 0, 0, 0.15);
            background: linear-gradient(135deg, rgba(255, 255, 255, 0.3), rgba(255, 255, 255, 0.2));
          }
        }
        
        .role-badge {
          font-size: 22rpx;
          padding: 4rpx 12rpx;
          border-radius: 12rpx;
          font-weight: 500;
          cursor: pointer;
          transition: all 0.3s ease;
          
          &.role-feeder {
            background: linear-gradient(135deg, #FFD700 0%, #FFA500 100%);
            color: #fff;
            box-shadow: 0 2rpx 8rpx rgba(255, 215, 0, 0.3);
          }
          
          &.role-foodie {
            background: linear-gradient(135deg, #FF6B9D 0%, #FF8E9E 100%);
            color: #fff;
            box-shadow: 0 2rpx 8rpx rgba(255, 107, 157, 0.3);
          }
          
          &:active {
            transform: scale(0.95);
          }
        }
      }
    }

    .heart-icon {
      animation: heartbeat 1.5s ease-in-out infinite;
      filter: drop-shadow(0 2rpx 8rpx rgba(255, 107, 157, 0.3));
      border-radius: 50%;
      padding: 15rpx;
      margin-top: -20rpx; // 稍微向上调整，与头像中心对齐
    }
  }

  .relationship-info {
    margin-bottom: 30rpx;
    cursor: pointer;
    transition: all 0.3s ease;
    
    &:active {
      opacity: 0.8;
      transform: scale(0.98);
    }

    .relationship-text {
      display: block;
      color: #ff6b9d;
      font-size: 31rpx;
      margin-bottom: 20rpx;
    }

    .days-count {
      margin-bottom: 20rpx;

      .days-number {
        font-size: 80rpx;
        font-weight: bold;
        color: #FF6B9D;
      }

      .days-unit {
        font-size: 32rpx;
        color: #FF6B9D;
        margin-left: 10rpx;
      }
    }

    .start-date {
      color: #ff6b9d;
      font-size: 30rpx;
    }
  }

  .invite-button {
    display: inline-flex;
    align-items: center;
    gap: 10rpx;
    background: linear-gradient(135deg, #FF6B9D 0%, #FF8E9E 100%);
    color: #fff;
    padding: 20rpx 40rpx;
    border-radius: 50rpx;
    font-size: 28rpx;

    .invite-text {
      color: #fff;
    }
  }
}

.feature-cards {
  display: flex;
  gap: 30rpx;
  margin-bottom: 60rpx;

  .feature-card {
    flex: 1;
    height: 160rpx;
    border-radius: 20rpx;
    overflow: hidden;

    .card-content {
      height: 100%;
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 0 30rpx;

      .card-title {
        color: #fff;
        font-size: 32rpx;
        font-weight: bold;
      }
    }

    &.diary-card {
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    }

    &.order-card {
      background: linear-gradient(135deg, #4CAF50 0%, #45a049 100%);
    }
  }
}

.anniversary-section {
  background: linear-gradient(135deg, #E6E6FA 0%, #DDA0DD 50%, #D8BFD8 100%);
  border-radius: 25rpx;
  padding: 30rpx;
  margin-top: 20rpx;
  box-shadow: 0 6rpx 25rpx rgba(221, 160, 221, 0.3);
  animation: fadeInScale 0.8s ease-out;
  
  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 30rpx;

    .section-title {
      font-size: 32rpx;
      font-weight: bold;
      color: #333;
      position: relative;
      
      &::after {
        content: '';
        position: absolute;
        bottom: -8rpx;
        left: 0;
        width: 60rpx;
        height: 3rpx;
        background: linear-gradient(90deg, #FF6B9D, #FF8E9E);
        border-radius: 2rpx;
      }
    }

    .add-anniversary {
      width: 60rpx;
      height: 60rpx;
      background: #fff;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.1);
      transition: all 0.3s ease;
      cursor: pointer;
      
      &:active {
        transform: scale(0.95);
        box-shadow: 0 2rpx 10rpx rgba(0, 0, 0, 0.15);
      }
      
      &:hover {
        transform: scale(1.05);
        box-shadow: 0 6rpx 25rpx rgba(0, 0, 0, 0.15);
      }
    }
  }

  .anniversary-list {
    .anniversary-item {
      display: flex;
      align-items: center;
      background: rgba(255, 255, 255, 0.95);
      padding: 35rpx;
      border-radius: 25rpx;
      margin-bottom: 25rpx;
      box-shadow: 0 8rpx 30rpx rgba(221, 160, 221, 0.2);
      backdrop-filter: blur(15rpx);
      border: 2rpx solid rgba(255, 255, 255, 0.6);
      transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
      position: relative;
      overflow: hidden;
      animation: slideInUp 0.6s ease-out;
      
      @for $i from 1 through 5 {
        &:nth-child(#{$i}) {
          animation-delay: #{$i * 0.1}s;
        }
      }
      
      &::before {
        content: '';
        position: absolute;
        top: 0;
        left: -100%;
        width: 100%;
        height: 100%;
        background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.3), transparent);
        transition: left 0.6s ease;
      }
      
      &::after {
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: linear-gradient(135deg, rgba(255, 255, 255, 0.1) 0%, transparent 50%, rgba(255, 255, 255, 0.1) 100%);
        opacity: 0;
        transition: opacity 0.3s ease;
      }
      
      &:last-child {
        margin-bottom: 0;
      }
      
      &:active {
        transform: translateY(3rpx) scale(0.98);
        box-shadow: 0 4rpx 20rpx rgba(221, 160, 221, 0.3);
        
        &::before {
          left: 100%;
        }
        
        &::after {
          opacity: 1;
        }
      }
      
      &:hover {
        transform: translateY(-2rpx);
        box-shadow: 0 12rpx 40rpx rgba(221, 160, 221, 0.3);
      }

      .anniversary-icon {
        margin-right: 35rpx;
        border-radius: 50%;
        padding: 20rpx;
        position: relative;
        z-index: 1;
        transition: all 0.3s ease;
        
        &[data-type="1"],
        &[data-type="2"],
        &[data-type="3"] {
          background: transparent;
          box-shadow: none;
        }
        
        &::after {
          display: none;
        }
        
        .icon-wrapper {
          background: transparent;
          border-radius: 0;
          padding: 0;
          box-shadow: none;
          
          .birthday-icon {
            width: 50rpx;
            height: 50rpx;
            object-fit: contain;
            display: block;
          }
        }
        
        &:hover {
          transform: scale(1.1);
          
          &[data-type="1"] {
            box-shadow: 0 6rpx 20rpx rgba(255, 107, 157, 0.4);
          }
          
          &[data-type="2"] {
            box-shadow: 0 6rpx 20rpx rgba(255, 215, 0, 0.4);
          }
          
          &[data-type="3"] {
            box-shadow: 0 6rpx 20rpx rgba(76, 175, 80, 0.4);
          }
          
          &::after {
            opacity: 0.5;
            transform: scale(1.1);
          }
        }
      }

      .anniversary-info {
        flex: 1;
        position: relative;
        z-index: 1;

        .anniversary-name {
          display: block;
          font-size: 32rpx;
          font-weight: 700;
          color: #2c3e50;
          margin-bottom: 12rpx;
          line-height: 1.3;
          text-shadow: 0 1rpx 2rpx rgba(0, 0, 0, 0.05);
        }

        .anniversary-meta {
          display: flex;
          align-items: center;
          gap: 10rpx;
          margin-bottom: 12rpx;

          .anniversary-type {
            font-size: 24rpx;
            font-weight: 600;
            padding: 6rpx 12rpx;
            border-radius: 12rpx;
            border: 1rpx solid;
            position: relative;
            overflow: hidden;
            
            &[data-type="1"] {
              color: #FF6B9D;
              background: rgba(255, 107, 157, 0.1);
              border-color: rgba(255, 107, 157, 0.3);
            }
            
            &[data-type="2"] {
              color: #FFD700;
              background: rgba(255, 215, 0, 0.1);
              border-color: rgba(255, 215, 0, 0.3);
            }
            
            &[data-type="3"] {
              color: #4CAF50;
              background: rgba(76, 175, 80, 0.1);
              border-color: rgba(76, 175, 80, 0.3);
            }
          }

          .anniversary-date {
            font-size: 26rpx;
            color: #7f8c8d;
            font-weight: 500;
            background: rgba(255, 255, 255, 0.8);
            padding: 6rpx 12rpx;
            border-radius: 15rpx;
            display: inline-block;
            box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.08);
          }
        }
      }

      .anniversary-countdown {
        position: relative;
        z-index: 1;
        
        .countdown-text {
          color: #FF6B9D;
          font-size: 28rpx;
          font-weight: 700;
          background: linear-gradient(135deg, rgba(255, 255, 255, 0.9) 0%, rgba(255, 255, 255, 0.95) 100%);
          padding: 12rpx 20rpx;
          border-radius: 25rpx;
          box-shadow: 0 4rpx 15rpx rgba(0, 0, 0, 0.1);
          border: 2rpx solid rgba(255, 107, 157, 0.2);
          position: relative;
          overflow: hidden;
          
          &::before {
            content: '';
            position: absolute;
            top: 0;
            left: -100%;
            width: 100%;
            height: 100%;
            background: linear-gradient(90deg, transparent, rgba(255, 107, 157, 0.1), transparent);
            transition: left 0.6s ease;
          }
          
          &:hover::before {
            left: 100%;
          }
        }
      }
    }

    .empty-anniversary {
      text-align: center;
      padding: 80rpx 0;
      color: #666;
      font-size: 28rpx;
      background: rgba(255, 255, 255, 0.7);
      border-radius: 20rpx;
      backdrop-filter: blur(10rpx);
      border: 1rpx solid rgba(255, 255, 255, 0.4);

      .empty-text {
        display: block;
        margin-bottom: 10rpx;
        color: #333;
        font-weight: 500;
      }

      .empty-desc {
        font-size: 24rpx;
        color: #888;
      }
    }
  }
}

@keyframes heartbeat {
  0% { 
    transform: scale(1); 
    filter: drop-shadow(0 2rpx 8rpx rgba(255, 107, 157, 0.3));
  }
  25% { 
    transform: scale(1.1); 
    filter: drop-shadow(0 4rpx 12rpx rgba(255, 107, 157, 0.4));
  }
  50% { 
    transform: scale(1.15); 
    filter: drop-shadow(0 6rpx 16rpx rgba(255, 107, 157, 0.5));
  }
  75% { 
    transform: scale(1.1); 
    filter: drop-shadow(0 4rpx 12rpx rgba(255, 107, 157, 0.4));
  }
  100% { 
    transform: scale(1); 
    filter: drop-shadow(0 2rpx 8rpx rgba(255, 107, 157, 0.3));
  }
}

@keyframes slideInUp {
  from {
    opacity: 0;
    transform: translateY(50px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes fadeInScale {
  from {
    opacity: 0;
    transform: scale(0.9);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

// 日期选择器弹窗样式
.date-picker-popup {
  background: #fff;
  border-radius: 30rpx 30rpx 0 0;
  padding: 40rpx 30rpx;
  max-height: 60vh;
  
  .popup-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 40rpx;
    padding-bottom: 30rpx;
    border-bottom: 1rpx solid #f0f0f0;
    
    .popup-title {
      font-size: 32rpx;
      font-weight: bold;
      color: #333;
    }
    
    .popup-actions {
      display: flex;
      gap: 30rpx;
      
      .popup-btn {
        font-size: 28rpx;
        padding: 10rpx 20rpx;
        border-radius: 10rpx;
        cursor: pointer;
        transition: all 0.3s ease;
        
        &.cancel-btn {
          color: #666;
          
          &:active {
            background: #f5f5f5;
          }
        }
        
        &.confirm-btn {
          color: #FF6B9D;
          font-weight: 600;
          
          &:active {
            background: rgba(255, 107, 157, 0.1);
          }
        }
      }
    }
  }
  
  .popup-content {
    .picker-view {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 30rpx;
      background: linear-gradient(135deg, #FFE4E1 0%, #FFB6C1 100%);
      border-radius: 20rpx;
      cursor: pointer;
      transition: all 0.3s ease;
      
      &:active {
        transform: scale(0.98);
        opacity: 0.9;
      }
      
      .picker-text {
        font-size: 32rpx;
        color: #333;
        font-weight: 500;
      }
    }
  }
}
</style>
