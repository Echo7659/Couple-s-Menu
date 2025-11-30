import { useUserStore } from '@/store/user.js'
export const handleAuth = () => {
  // 用户登录鉴权
  const userInfo = useUserStore()
  if (!userInfo.token) {
    uni.navigateTo({
      url: '/pages/wxLogin/wxLogin'
    })
    return false
  }

  return true
}
