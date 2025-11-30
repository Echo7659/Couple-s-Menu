import net from './request.js';

// import store from '../store';

const API = {
  /*
   用户相关
   */
  login: data => {
    return net.post('/user/login', data, false).catch(error => {
      // 统一处理登录接口的错误，确保错误格式一致
      console.error('登录接口错误:', error)
      // 如果是网络错误或请求失败，返回统一格式
      if (error && typeof error === 'object') {
        // 如果已经是标准格式，直接返回
        if (error.code !== undefined || error.msg !== undefined) {
          return Promise.reject(error)
        }
        // 否则包装成标准格式
        return Promise.reject({
          code: error.code || -1,
          msg: error.msg || error.message || error.errMsg || '登录失败，请检查网络连接',
          data: null
        })
      }
      // 字符串错误
      if (typeof error === 'string') {
        return Promise.reject({
          code: -1,
          msg: error,
          data: null
        })
      }
      // 其他情况
      return Promise.reject({
        code: -1,
        msg: '登录失败，请重试',
        data: null
      })
    })
  }, // 登录
  getUserInfo: () => {
    // 检查是否有 token，如果没有则不发起请求
    const token = uni.getStorageSync('token')
    if (!token) {
      return Promise.resolve({ code: 9, msg: '未登录', data: null })
    }
    return net.get('/user/info', {}, true)
  }, // 获取用户信息
  updateUserInfo: data => net.put('/user/info', data, true), // 更新用户信息
  upload: data => net.post('/user/upload', data, true), // 通用上传接口
  getUserInfoById: (id) => net.get(`/user/info/${id}`, {}, true), // 根据ID获取用户信息
  getLoverInfo: () => net.get('/user/lover', {}, true), // 获取情侣关系信息（包含 startDate 与 loverId）
  bindLover: data => net.post('/user/lover/bind', data, true), // 使用邀请码绑定情侣关系
  updateLoverDay: data => net.put('/user/lover/day', data, true), // 更新情侣恋爱开始日期
  updateUserRole: data => net.put('/user/role', data, true), // 更新用户角色（1=饲养员，2=吃货）
  refreshToken: () => net.post('/user/refreshToken', {}, true), // 刷新token

  /* 
    纪念日相关
  */
  getMemorialDay: (id) => net.get(`/anniversary/${id}`, {}, true), // 获取纪念日详情
  getMemorialDayList: (params = {}) => net.get('/anniversary/list', params, true), // 获取纪念日列表（需 userId）
  addMemorialDay: data => net.post('/anniversary', data, true), // 添加纪念日
  updateMemorialDay: data => net.put('/anniversary', data, true), // 更新纪念日
  deleteMemorialDay: data => net.delete('/anniversary', data, true), // 删除纪念日

  /* 
     菜谱相关接口
  */
  getMenuList: (params = {}) => net.get('/menu/list', params, true), // 获取菜谱列表（支持 page/pageSize/keyword）
  getMenuById: (id) => net.get(`/menu/${id}`, {}, true), // 获取菜谱详情
  getMenuWithCategories: (id) => net.get(`/menu/${id}/categories`, {}, true), // 获取菜谱及其分类和菜品（推荐）
  createMenu: data => net.post('/menu', data, true), // 创建菜谱
  updateMenu: data => net.put('/menu', data, true), // 更新菜谱
  deleteMenu: data => net.delete('/menu', data, true), // 删除菜谱

  /* 
     菜品分类相关接口
  */
  getFoodCategoryList: (params = {}) => net.get('/food/category/list', params, true), // 获取分类列表（支持 menuId/page/pageSize/keyword）
  getFoodCategoryById: (id) => net.get(`/food/category/${id}`, {}, true), // 获取分类详情
  addFoodCategory: data => net.post('/food/category', data, true), // 添加分类（需要 menuId）
  updateFoodCategory: data => net.put('/food/category', data, true), // 更新分类
  deleteFoodCategory: data => net.delete('/food/category', data, true), // 删除分类

  /* 
     菜品相关接口
  */
  getFoodList: (params = {}) => net.get('/food/list', params, true), // 获取菜品列表（支持 menuId/categoryId/name/keyword/page/pageSize）
  getFoodDetail: (id) => net.get(`/food/${id}`, {}, true), // 获取菜品详情
  addFood: data => net.post('/food', data, true), // 添加菜品（需要 categoryId）
  updateFood: data => net.put('/food', data, true), // 更新菜品
  deleteFood: data => net.delete('/food', data, true), // 删除菜品

  // 订单相关接口
  createOrder: data => net.post('/order', data, true), // 创建订单
  getOrderList: (params = {}) => net.get('/order/list', params, true), // 获取订单列表,只需要userId和date参数
  getOrderDetail: (id) => net.get(`/order/${id}`, {}, true), // 获取订单详情
  deleteOrder: data => net.delete('/order', data, true), // 删除订单

  /* 
      日记相关接口
  */
  getDiaryList: (params = {}) => net.get('/diary/list', params, true), // 获取日记列表,只需要userId和date参数
  getDiaryDetail: (id) => net.get(`/diary/${id}`, {}, true), // 获取日记详情
  addDiary: data => net.post('/diary', data, true), // 添加日记
  updateDiary: data => net.put('/diary', data, true), // 更新日记
  deleteDiary: data => net.delete('/diary', data, true), // 删除日记

  /* 
      打卡目标（Target）相关接口
      后端文档基础路径为 /api/app/target，这里使用完整路径以兼容现有接口
  */
  // 创建打卡目标
  createTarget: data => net.post('/target', data, true),
  // 更新打卡目标
  updateTarget: data => net.put('/target', data, true),
  // 删除打卡目标
  deleteTarget: data => net.delete('/target', data, true),
  // 打卡目标列表（通用查询）
  getTargetList: (params = {}) => net.get('/target/list', params, true),
  // 获取“我的目标”（推荐）
  getMyTargets: (params = {}) => net.get('/target/my', params, true),
  // 获取打卡目标详情（含完整 records）
  getTargetDetail: (id) => net.get(`/target/${id}`, {}, true),
  // 创建打卡记录（打卡）
  createTargetRecord: data => net.post('/target/record', data, true),
  // 获取某个目标的打卡记录
  getTargetRecords: (id) => net.get(`/target/${id}/records`, {}, true),
  // 完成打卡目标
  completeTarget: data => net.post('/target/complete', data, true)
}

export default API