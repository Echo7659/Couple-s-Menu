import config from '../config'

const LOGIN_PAGE = '/pages/user/login'

// -------- Token 刷新控制 --------
let isRefreshingToken = false
let refreshWaitQueue = [] // 存放等待刷新完成后继续的回调

const enqueueRefresh = (resolver) => refreshWaitQueue.push(resolver)
const resolveRefreshQueue = (token) => {
  refreshWaitQueue.forEach((resolve) => resolve(token))
  refreshWaitQueue = []
}

async function refreshTokenSilently() {
  // 如果本地没有 token，说明用户未登录，不执行刷新
  const currentToken = getToken()
  if (!currentToken) {
    console.log('未登录，跳过 token 刷新')
    return null
  }
  
  if (isRefreshingToken) {
    // 已在刷新流程，返回一个等待 Promise
    return new Promise((resolve) => enqueueRefresh(resolve))
  }
  isRefreshingToken = true
  try {
    // 调用你后端提供的刷新接口（无需参数）
    let newToken = ''
    try {
      const res = await baseRequest({ url: '/user/refreshToken', method: 'POST', data: {}, checkLogin: false })
      // 该接口返回形如：{ code:0, data: '<new_token>', msg:'成功' }
      if (res?.code === 0 && typeof res.data === 'string') {
        newToken = res.data
      }
    } catch (e) {
      // ignore
    }

    if (newToken) {
      uni.setStorageSync('token', newToken)
      resolveRefreshQueue(newToken)
      return newToken
    }
    resolveRefreshQueue(null)
    return null
  } finally {
    isRefreshingToken = false
  }
}

/**
 * 获取本地 token
 */
function getToken() {
  return uni.getStorageSync('token') || ''
}

/**
 * 统一处理未登录或被踢下线
 */
function handleAuthFail(checkLogin = true) {
  console.log('开始清理登录信息并跳转登录页')
  uni.removeStorageSync('token')
  uni.removeStorageSync('tokenExpiresAt')
  uni.removeStorageSync('userInfo')
  // if (checkLogin) {
  //   console.log('准备跳转到登录页:', LOGIN_PAGE)
  //   // 使用 reLaunch 强制跳转到登录页，避免用户返回到已失效的页面
  //   uni.reLaunch({ 
  //     url: LOGIN_PAGE,
  //     success: () => {
  //       console.log('跳转登录页成功')
  //     },
  //     fail: (err) => {
  //       console.log('跳转登录页失败:', err)
  //     }
  //   })
  // }
}

/**
 * 构建请求头
 */
function buildHeaders(customHeaders = {}) {
  return {
    'Content-Type': 'application/json',
    'Authorization': getToken(),
    ...customHeaders
  }
}

/**
 * 通用请求方法
 */
function baseRequest({ url, method = 'GET', data = {}, checkLogin = true, headers = {} }) {
  return new Promise((resolve, reject) => {
    const fullUrl = config.baseUrl.replace(/\/$/, '') + '/' + url.replace(/^\//, '')
    console.log('请求URL:', fullUrl)

    const MAX_REFRESH_RETRY = 3

    const send = async (left) => {
      uni.request({
        url: fullUrl,
        method,
        data,
        header: buildHeaders(headers),
        success: async (res) => {
          const { statusCode, data: resData } = res

          // 统一处理 token 过期：401 或 200+code=9
          if (statusCode === 401 || (statusCode === 200 && resData && resData.code === 9)) {
            // 如果本地没有 token，说明用户未登录，直接返回错误，不执行刷新
            const currentToken = getToken()
            if (!currentToken) {
              resolve(resData || { code: 9, msg: '未登录' })
              return
            }
            
            if (left <= 0) {
              handleAuthFail(checkLogin)
              resolve(resData || { code: 9, msg: 'Token已过期' })
              return
            }
            try {
              const newToken = await refreshTokenSilently()
              if (newToken) {
                // 递归重试，减少次数
                send(left - 1)
                return
              }
            } catch (e) {}
            handleAuthFail(checkLogin)
            resolve(resData || { code: 9, msg: 'Token已过期' })
            return
          }

          if (statusCode === 200 && resData) {
            // 账号在其他设备登录
            if (resData.code === 409) {
              uni.showModal({
                title: '提示',
                content: '您的账号已在其他设备登录，已强制下线!',
                showCancel: false,
                success: () => handleAuthFail(checkLogin)
              })
              resolve(resData)
              return
            }
            // 用户认证失败
            if (resData.code === 5) {
              handleAuthFail(checkLogin)
              resolve(resData)
              return
            }
            resolve(resData)
          } else {
            reject({ code: statusCode, msg: res.errMsg || '网络错误', data: null })
          }
        },
        fail: (err) => {
          console.log('请求失败:', err)
          console.log('请求失败详情:', JSON.stringify(err, null, 2))
          
          if (err.statusCode === 401 || err.code === 401) {
            // 如果本地没有 token，说明用户未登录，直接返回错误，不执行刷新
            const currentToken = getToken()
            if (!currentToken) {
              reject({ code: 7, msg: '未登录', data: null })
              return
            }
            
            if (left <= 0) {
              handleAuthFail(true)
              reject({ code: 7, msg: '授权已过期', data: null })
              return
            }
            // 尝试刷新后再发
            refreshTokenSilently()
              .then((tk) => {
                if (tk) send(left - 1)
                else {
                  handleAuthFail(true)
                  reject({ code: 7, msg: '授权已过期', data: null })
                }
              })
              .catch(() => {
                handleAuthFail(true)
                reject({ code: 7, msg: '授权已过期', data: null })
              })
            return
          }
          
          // 统一错误格式
          const errorMsg = err.errMsg || err.message || '网络请求失败'
          reject({ 
            code: err.statusCode || err.code || -1, 
            msg: errorMsg, 
            data: null 
          })
        }
      })
    }

    send(MAX_REFRESH_RETRY)
  })
}

/**
 * 导出常用请求方法
 */
const request = {
  get(url, data = {}, checkLogin = true, headers = {}) {
    return baseRequest({ url, method: 'GET', data, checkLogin, headers })
  },
  post(url, data = {}, checkLogin = true, headers = {}) {
    return baseRequest({ url, method: 'POST', data, checkLogin, headers })
  },
  put(url, data = {}, checkLogin = true, headers = {}) {
    return baseRequest({ url, method: 'PUT', data, checkLogin, headers })
  },
  delete(url, data = {}, checkLogin = true, headers = {}) {
    return baseRequest({ url, method: 'DELETE', data, checkLogin, headers })
  },
  
  // 兼容原有的http.request方法
  request: ({ url, method = 'GET', data = {}, headers = {} }) => {
    return baseRequest({ url, method, data, checkLogin: true, headers })
  }
}

export default request