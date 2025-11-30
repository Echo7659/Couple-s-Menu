<script setup>
import { useSystemStore } from '@/store/system.js'
import { useUserStore } from '@/store/user.js'
import { onLaunch, onShow } from '@dcloudio/uni-app'

const systemStore = useSystemStore()
const userStore = useUserStore()

onLaunch((options) => {
  // 初始化系统信息，并同步一次本地登录状态（不做强制跳转）
  systemStore.init()
  userStore.checkLoginStatus()
  
  console.log('App onLaunch 参数:', JSON.stringify(options, null, 2))
  
  // 检查启动参数中是否有邀请码（从分享链接进入时）
  // 小程序从分享链接进入时，参数可能在 options.query 中
  let hasInviteCode = false
  if (options) {
    let invite = null
    
    // 方式1: 从 query 对象中获取
    if (options.query && typeof options.query === 'object') {
      invite = options.query.invite || options.query.inviteCode
    }
    
    // 方式2: 从路径中解析（如果 query 是字符串）
    if (!invite && options.path) {
      const pathMatch = options.path.match(/[?&]invite=([^&]+)/)
      if (pathMatch) {
        invite = decodeURIComponent(pathMatch[1])
      }
    }
    
    // 方式3: 直接从 options 中获取（某些情况下参数可能在顶层）
    if (!invite) {
      invite = options.invite || options.inviteCode
    }
    
    if (invite) {
      console.log('✅ 从启动参数获取到邀请码:', invite)
      uni.setStorageSync('INVITE_CODE', invite)
      hasInviteCode = true
    } else {
      console.log('⚠️ 启动参数中未找到邀请码')
    }
  }
  
  // 兼容 scene 场景值（小程序码、分享等）
  if (options && options.scene) {
    console.log('小程序启动场景值:', options.scene)
  }
  
  // 如果从分享链接进入且有邀请码，但当前页面是首页，则跳转到登录页
  // 注意：这里需要延迟执行，确保页面已经加载
  if (hasInviteCode) {
    setTimeout(() => {
      const pages = getCurrentPages()
      const currentPage = pages[pages.length - 1]
      if (currentPage && currentPage.route === 'pages/home/home') {
        console.log('检测到从分享链接进入首页，跳转到登录页')
        const invite = uni.getStorageSync('INVITE_CODE')
        uni.navigateTo({
          url: `/pages/user/login?invite=${encodeURIComponent(invite)}`
        })
      }
    }, 500)
  }
})

onShow((options) => {
  // 每次显示时仅同步一次本地登录状态，不再强制跳转登录页
  userStore.checkLoginStatus()
  
  // 检查显示参数中是否有邀请码（从分享链接进入时）
  if (options && options.query) {
    const invite = options.query.invite || options.query.inviteCode
    if (invite) {
      console.log('从显示参数获取到邀请码:', invite)
      uni.setStorageSync('INVITE_CODE', invite)
    }
  }
})
</script>

<style lang="scss">
/*每个页面公共css */
@import 'styles/global.scss';
@import 'static/font/iconfont.css';
//主题scss
.default-theme {
  @import 'styles/themes/default-theme.scss';
}
//主题scss
.red-theme {
  @import 'styles/themes/default-theme.scss';
}
</style>

<style></style>
