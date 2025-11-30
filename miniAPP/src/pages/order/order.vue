<template>
  <view class="order-container">
    <!-- 顶部搜索 -->
    <view class="search-section">
      <view class="search-bar">
      <uni-easyinput
        v-model="searchText"
          placeholder="搜索菜品名称"
        prefixIcon="search"
        trim="both"
          @confirm="doSearch"
      ></uni-easyinput>
        <view class="search-btn" @click="doSearch">
          <text>搜索</text>
        </view>
      </view>
    </view>

    <!-- 主体：左侧分类 + 右侧菜品列表 -->
    <view class="content">
      <scroll-view class="category-pane" scroll-y>
        <view :class="['category-item', { active: selectedCategoryId === 0 }]" @click="selectCategory(0)">
          <text class="name">全部</text>
        </view>
      <view 
          v-for="cat in categoryList"
          :key="cat.ID"
          :class="['category-item', { active: selectedCategoryId === cat.ID }]"
          @click="selectCategory(cat.ID)"
        >
          <text class="name">{{ cat.name }}</text>
        </view>
        <view class="category-manage" @click="openCategoryManage">
          <uni-icons type="gear" size="16" color="#999"></uni-icons>
          <text class="manage-text">分类管理</text>
        </view>
        <view class="food-manage" @click="openFoodManage">
          <uni-icons type="list" size="16" color="#999"></uni-icons>
          <text class="manage-text">菜品管理</text>
            </view>
      </scroll-view>

      <scroll-view class="food-pane" scroll-y @scrolltolower="loadMore">
        <view v-if="filteredFoodList.length === 0" class="empty-holder">
          <uni-icons type="folder-add" size="80" color="#ccc"></uni-icons>
          <text class="empty-text">这里空空如也~</text>
            </view>

        <view v-for="item in filteredFoodList" :key="item.ID" class="food-card">
          <image :src="item.img" mode="aspectFill" class="cover"></image>
          <view class="info">
            <text class="title">{{ item.name }}</text>
            <text class="desc">{{ item.desc || '美味等你来尝~' }}</text>
            <view class="actions">
              <view v-if="getQty(item.ID) > 0" class="qty">
                <view class="btn minus" @click="decrease(item)">
                  <uni-icons type="minus" size="24" color="#FF6B9D"></uni-icons>
          </view>
                <text class="num">{{ getQty(item.ID) }}</text>
                <view class="btn plus" @click="increase(item)">
                  <uni-icons type="plus" size="24" color="#fff"></uni-icons>
          </view>
        </view>
              <view v-else class="add-btn" @click="add(item)">
                <text>添加</text>
            </view>
          </view>
        </view>
      </view>
      
        <view class="load-more" v-if="loadingFoods || !foodHasMore">
          <!-- 剧中显示 -->
          <text class="tip" style="text-align: center;">{{ loadingFoods ? '加载中...' : '没有更多了' }}</text>
      </view>
      </scroll-view>
    </view>

    <!-- 底部购物车栏 -->
    <view class="cart-bar" v-if="totalQty > 0">
      <view class="left" @click="openCart">
        <view class="cart-icon">
          <uni-icons type="cart" size="22" color="#FF6B9D"></uni-icons>
          <view class="badge">{{ totalQty }}</view>
        </view>
        <text class="summary">菜品数：{{ totalQty }}</text>
      </view>
      <view class="right">
        <view class="clear" @click.stop="clearCart">
          <uni-icons type="trash" size="18" color="#999"></uni-icons>
        </view>
        <button class="submit" @click="submitOrder">去下单</button>
      </view>
    </view>

    <!-- 购物车弹层 -->
    <uni-popup ref="cartPopup" type="bottom" :mask-click="true" @change="onPopupChange">
      <view class="cart-sheet">
        <view class="sheet-header">
          <text class="title">已选商品</text>
          <view class="actions">
            <view class="clear" @click="clearCart">
              <uni-icons type="trash" size="18" color="#999"></uni-icons>
              <text class="txt">清空</text>
            </view>
            <view class="close" @click="closeCart">
              <uni-icons type="close" size="22" color="#999"></uni-icons>
            </view>
          </view>
        </view>
        <scroll-view class="sheet-body" scroll-y>
          <view v-for="it in currentCart" :key="it.id" class="row">
            <image :src="it.data.img" class="thumb" mode="aspectFill"></image>
            <view class="meta">
              <text class="name">{{ it.data.name }}</text>
            </view>
            <view class="qty">
              <view class="btn minus" @click="decrease(it.data)">
                <uni-icons type="minus" size="24" color="#FF6B9D"></uni-icons>
              </view>
              <text class="num">{{ it.qty }}</text>
              <view class="btn plus" @click="increase(it.data)">
                <uni-icons type="plus" size="24" color="#fff"></uni-icons>
              </view>
            </view>
          </view>
        </scroll-view>
      </view>
    </uni-popup>
  </view>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { onPullDownRefresh, onShow } from '@dcloudio/uni-app'
import api from '@/api'

// 搜索与筛选
const searchText = ref('')
const menuList = ref([])
const selectedMenuId = ref(null)
const categoryList = ref([])
const selectedCategoryId = ref(0) // 0 表示全部

// 分页
const categoryPage = ref(1)
const categoryPageSize = ref(50)
const foodPage = ref(1)
const foodPageSize = ref(10)
const foodHasMore = ref(true)
const loadingFoods = ref(false)

// 常量
const CART_STORAGE_KEY = 'menu_cart_map'

const serializeCartMap = (val) => {
  try {
    return JSON.parse(JSON.stringify(val || {}))
  } catch {
    return {}
  }
}

// 菜品与购物车
const foodList = ref([]) // 接口返回的原始菜品列表
const cartMap = ref({}) // { [menuId]: [{ id, qty, data }] }
const currentCart = computed(() => {
  const id = selectedMenuId.value
  if (!id) return []
  return cartMap.value[id] || []
})

const cartPopup = ref(null)
const cartVisible = ref(false)
const openCart = () => { if (currentCart.value.length && cartPopup.value?.open) cartPopup.value.open() }
const closeCart = () => { if (cartPopup.value?.close) cartPopup.value.close() }
const onPopupChange = (e) => { cartVisible.value = e.show }

const ensureCart = () => {
  const id = selectedMenuId.value
  if (!id) return []
  if (!cartMap.value[id]) {
    cartMap.value = {
      ...cartMap.value,
      [id]: []
    }
  }
  return cartMap.value[id]
}

// 从路由参数获取 menuId
const getMenuIdFromRoute = () => {
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1]
  const options = currentPage.options || {}
  return options.menuId || null
}

// 加载菜谱的完整数据（推荐接口）
const loadMenuData = async () => {
  if (!selectedMenuId.value) return
  try {
    const res = await api.getMenuWithCategories(selectedMenuId.value)
    if (res?.code === 0 && res.data) {
      // 更新分类列表
      categoryList.value = res.data.categories || []
      // 更新菜品列表
      foodList.value = res.data.foods || []
      foodHasMore.value = false // 一次性加载完所有数据
    }
  } catch (e) {
    console.error('加载菜谱数据失败:', e)
    // 如果推荐接口失败，回退到分步加载
    loadCategories()
    loadFoods(true)
  }
}

// 分步加载分类（备用方案）
const loadCategories = async (reset = false) => {
  if (!selectedMenuId.value) return
  try {
    if (reset) categoryPage.value = 1
    const res = await api.getFoodCategoryList({ 
      menuId: selectedMenuId.value,
      page: categoryPage.value, 
      pageSize: categoryPageSize.value 
    })
    if (res?.code === 0) {
      const list = res.data?.list || []
      categoryList.value = list
    }
  } catch (e) {
    console.error('加载分类失败:', e)
  }
}

// 分步加载菜品（备用方案）
const loadFoods = async (reset = false) => {
  if (loadingFoods.value || !selectedMenuId.value) return
  try {
    loadingFoods.value = true
    if (reset) {
      foodPage.value = 1
      foodHasMore.value = true
      foodList.value = []
    }
    const kw = (searchText.value || '').trim()
    const params = {
      menuId: selectedMenuId.value,
      page: foodPage.value,
      pageSize: foodPageSize.value
    }
    // 支持 name 和 keyword 两种搜索参数（根据API指南）
    if (kw) {
      params.name = kw
      params.keyword = kw
    }
    if (selectedCategoryId.value && selectedCategoryId.value !== 0) {
      params.categoryId = selectedCategoryId.value
    }
    const res = await api.getFoodList(params)
    if (res?.code === 0) {
      const list = res.data?.list || []
      if (reset) {
        foodList.value = list
      } else {
        foodList.value = foodList.value.concat(list)
      }
      const total = res.data?.total ?? list.length
      const loaded = foodList.value.length
      foodHasMore.value = loaded < total
      if (foodHasMore.value) foodPage.value += 1
    }
  } catch (e) {
    console.error('加载菜品失败:', e)
  } finally {
    loadingFoods.value = false
  }
}

// 初始化
const init = async () => {
  // 从路由参数获取 menuId
  const routeMenuId = getMenuIdFromRoute()
  if (routeMenuId) {
    selectedMenuId.value = Number(routeMenuId)
    loadMenuData()
  } else {
    // 如果没有 menuId 参数，返回菜谱列表页
    uni.showToast({
      title: '请先选择菜谱',
      icon: 'none'
    })
    setTimeout(() => {
      uni.switchTab({
        url: '/pages/order/index'
      })
    }, 1500)
  }
}

onMounted(() => {
  // 初始化本地购物车缓存
  try {
    const stored = uni.getStorageSync(CART_STORAGE_KEY)
    if (stored && typeof stored === 'object') {
      cartMap.value = stored
    }
  } catch (e) {
    console.warn('读取购物车缓存失败', e)
  }
  // 获取用户ID
  fetchUserId()
  init()
})

// 持久化购物车到本地缓存
watch(
  cartMap,
  (val) => {
    try {
      const plain = serializeCartMap(val)
      uni.setStorageSync(CART_STORAGE_KEY, plain)
    } catch (e) {
      console.warn('写入购物车缓存失败', e)
    }
  },
  { deep: true }
)

watch(selectedMenuId, () => {
  // 确保切换菜谱时存在对应购物车
  ensureCart()
  closeCart()
})

onShow(() => {
  // 如果已选择菜谱，刷新数据
  if (selectedMenuId.value) {
    loadMenuData()
  } else {
    init()
  }
})

onPullDownRefresh(() => {
  if (selectedMenuId.value) {
    loadMenuData().finally(() => {
      uni.stopPullDownRefresh()
    })
  } else {
    init().finally(() => {
      uni.stopPullDownRefresh()
    })
  }
})

const selectCategory = (id) => {
  selectedCategoryId.value = id
  // 如果使用推荐接口一次性加载，则在前端过滤
  // 否则重新请求接口
  if (foodList.value.length > 0 && !searchText.value) {
    // 前端过滤
    return
  }
  loadFoods(true)
}

const doSearch = () => {
  // 搜索时需要重新请求接口
  loadFoods(true)
}

const openCategoryManage = () => {
  if (!selectedMenuId.value) {
    uni.showToast({
      title: '请先选择菜谱',
      icon: 'none'
    })
    return
  }
  uni.navigateTo({ 
    url: `/pages/food/category-manage?menuId=${selectedMenuId.value}` 
  })
}

const openFoodManage = () => {
  if (!selectedMenuId.value) {
    uni.showToast({
      title: '请先选择菜谱',
      icon: 'none'
    })
    return
  }
  uni.navigateTo({ 
    url: `/pages/food/food-manage?menuId=${selectedMenuId.value}` 
  })
}

// 过滤菜品列表（根据选择的分类和搜索关键词）
const filteredFoodList = computed(() => {
  let list = foodList.value
  
  // 按分类过滤
  if (selectedCategoryId.value && selectedCategoryId.value !== 0) {
    list = list.filter(item => item.categoryId === selectedCategoryId.value)
  }
  
  // 按搜索关键词过滤
  const kw = (searchText.value || '').trim()
  if (kw) {
    list = list.filter(item => 
      item.name && item.name.includes(kw)
    )
  }
  
  return list
})

const loadMore = () => {
  if (foodHasMore.value && !loadingFoods.value) {
    loadFoods(false)
  }
}

// 购物车逻辑
const getQty = (id) => {
  const list = ensureCart()
  const found = list.find(i => i.id === id)
  return found ? found.qty : 0
}

const add = (item) => {
  const list = ensureCart()
  const found = list.find(i => i.id === item.ID)
  if (found) {
    found.qty++
  } else {
    list.push({ id: item.ID, qty: 1, data: item })
  }
}

const increase = (item) => {
  add(item)
}

const decrease = (item) => {
  const list = ensureCart()
  const idx = list.findIndex(i => i.id === item.ID)
  if (idx !== -1) {
    if (list[idx].qty > 1) {
      list[idx].qty--
    } else {
      list.splice(idx, 1)
    }
  }
}

const clearCart = () => {
  if (!selectedMenuId.value) return
  uni.showModal({
    title: '清空购物车',
    content: '确定要清空已选菜品吗？',
    success: (res) => {
      if (res.confirm) {
        cartMap.value = {
          ...cartMap.value,
          [selectedMenuId.value]: []
        }
      }
    }
  })
}

const totalQty = computed(() => currentCart.value.reduce((sum, it) => sum + it.qty, 0))
const totalPrice = computed(() => currentCart.value.reduce((sum, it) => sum + ((it.data.price || 0) * it.qty), 0))

const userId = ref(null)
const fetchUserId = async () => {
  try {
    const cached = uni.getStorageSync('userInfo')
    if (cached?.ID || cached?.id) {
      userId.value = cached.ID || cached.id
      return
    }
    const res = await api.getUserInfo()
    if (res?.code === 0 && res.data) {
      userId.value = res.data.ID || res.data.id || null
    }
  } catch (e) {
    console.warn('获取用户ID失败', e)
  }
}

const submitOrder = async () => {
  const list = currentCart.value
  if (list.length === 0) {
    uni.showToast({ title: '请先选择菜品', icon: 'none' })
    return
  }
  
  // 确保已获取用户ID
  if (!userId.value) {
    await fetchUserId()
  }
  
  // 如果还是没有用户ID，提示登录
  if (!userId.value) {
    uni.showToast({
      title: '请先登录',
      icon: 'none'
    })
    setTimeout(() => {
      uni.navigateTo({ url: '/pages/user/login' })
    }, 1500)
    return
  }
  
  const items = list.map(it => ({
    id: it.id,
    name: it.data.name,
    img: it.data.img,
    qty: it.qty,
    price: it.data.price || 0
  }))
  const ids = items.map(i => i.id).join(',')

  uni.navigateTo({
    url: '/pages/order/confirm',
    success: (res) => {
      res.eventChannel.emit('initData', {
        items,
        ids,
        userId: userId.value
      })
    }
  })
}
</script>

<style lang="scss" scoped>
.order-container {
  min-height: 100vh;
  background: #f5f5f5;
  padding-bottom: 140rpx;
}

.search-section {
  background: #fff;
  padding: 20rpx 24rpx;

  .search-bar {
    display: flex;
    align-items: center;
    gap: 16rpx;

    .search-btn {
      padding: 18rpx 26rpx;
      background: linear-gradient(135deg, #FF6B9D 0%, #FF8E9E 100%);
      color: #fff;
      border-radius: 16rpx;
      font-size: 26rpx;
    }
  }
}

.content {
  display: flex;
  height: calc(100vh - 220rpx);
}

.category-pane {
  width: 180rpx;
  background: #fafafa;
  border-right: 1rpx solid #eee;
      
      .category-item {
    padding: 28rpx 20rpx;
    color: #666;
          font-size: 28rpx;
    border-left: 6rpx solid transparent;
        
        &.active {
      color: #FF6B9D;
      background: #fff;
      border-left-color: #FF6B9D;
      font-weight: 600;
    }
  }

  .category-manage, .food-manage {
    padding: 24rpx 20rpx;
    display: flex;
    align-items: center;
    gap: 10rpx;
    color: #999;
    font-size: 24rpx;
  }
}

.food-pane {
  flex: 1;
  padding: 20rpx;

  .empty-holder {
    text-align: center;
    padding: 140rpx 0;

    .empty-text {
      display: block;
      margin-top: 16rpx;
      color: #999;
      font-size: 26rpx;
    }
  }

  .food-card {
    display: flex;
    background: #fff;
    border-radius: 20rpx;
    padding: 20rpx;
    margin-bottom: 20rpx;
    box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.05);
    
    .cover {
      width: 160rpx;
      height: 160rpx;
      border-radius: 16rpx;
      object-fit: cover;
      margin-right: 20rpx;
      flex-shrink: 0;
    }
    
    .info {
      flex: 1;
      
      .title {
        display: block;
        font-size: 32rpx;
        font-weight: 700;
        color: #333;
      }
      
      .desc {
        display: block;
        margin-top: 8rpx;
        font-size: 24rpx;
        color: #999;
      }
      
      .actions {
        margin-top: 16rpx;
        display: flex;
        justify-content: flex-end;
        align-items: center;
        gap: 16rpx;

        .add-btn {
          padding: 12rpx 22rpx;
          background: #FF6B9D;
          color: #fff;
          border-radius: 22rpx;
          font-size: 24rpx;
          box-shadow: 0 4rpx 12rpx rgba(255, 107, 157, 0.25);
        }

        .qty {
          display: flex;
          align-items: center;
          gap: 14rpx;

          .btn {
            width: 48rpx;
            height: 48rpx;
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;
            box-shadow: 0 2rpx 8rpx rgba(255, 107, 157, 0.2);
        }
        
          .minus {
            background: #fff;
          border: 2rpx solid #FF6B9D;
          }
          
          .plus {
            background: #FF6B9D;
          }

          .num {
            min-width: 40rpx;
            text-align: center;
            font-size: 28rpx;
            color: #333;
          }
        }
      }
    }
  }
}

.cart-bar {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  background: #fff;
  border-radius: 20rpx 20rpx 0 0;
  box-shadow: 0 -4rpx 20rpx rgba(0, 0, 0, 0.08);
  padding: 20rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  z-index: 999;
  
  .left {
    display: flex;
    align-items: center;
    gap: 14rpx;
    
    .cart-icon {
      position: relative;
      
      .badge {
        position: absolute;
        top: -8rpx;
        right: -10rpx;
        background: #FF6B9D;
        color: #fff;
        font-size: 20rpx;
        padding: 2rpx 8rpx;
        border-radius: 12rpx;
      }
    }

    .summary {
      font-size: 28rpx;
      color: #FF6B9D;
      font-weight: 600;
    }
  }

  .right {
    display: flex;
    align-items: center;
    gap: 16rpx;

    .clear {
      padding: 10rpx;
    }

    .submit {
      background: linear-gradient(135deg, #FF6B9D 0%, #FF8E9E 100%);
      color: #fff;
      border: none;
      border-radius: 28rpx;
      padding: 16rpx 28rpx;
      font-size: 28rpx;
    }
  }
}

/* 购物车弹层样式 */
.cart-sheet {
  background: #fff;
  border-radius: 20rpx 20rpx 0 0;
  overflow: hidden;
  /* 预留底部空间，避免被底部栏遮挡 */
  padding-bottom: calc(150rpx + constant(safe-area-inset-bottom));
  padding-bottom: calc(150rpx + env(safe-area-inset-bottom));
}

.sheet-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24rpx;
  border-bottom: 1rpx solid #f0f0f0;

  .title {
    font-size: 30rpx;
    font-weight: 600;
    color: #333;
  }

  .actions {
    display: flex;
    align-items: center;
    gap: 10rpx;

    .clear {
      display: flex;
      align-items: center;
      gap: 6rpx;
      padding: 6rpx 10rpx;
      border-radius: 12rpx;
      background: #f7f7f7;
    }
  }
}

.sheet-body {
  max-height: 60vh;
  padding: 10rpx 24rpx 0;

  .row {
    display: flex;
    align-items: center;
    padding: 16rpx 0;
    border-bottom: 1rpx solid #f5f5f5;

    .thumb {
      width: 100rpx;
      height: 100rpx;
      border-radius: 12rpx;
      object-fit: cover;
      margin-right: 16rpx;
    }

    .meta {
      flex: 1;

      .name {
        display: block;
        font-size: 28rpx;
        color: #333;
      }

      .price {
        display: block;
        margin-top: 6rpx;
        color: #FF6B9D;
        font-weight: 600;
      }
    }

    .qty {
      display: flex;
      align-items: center;
      gap: 14rpx;
      margin-right: 40rpx;

      .btn {
        width: 44rpx;
        height: 44rpx;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        box-shadow: 0 2rpx 8rpx rgba(255, 107, 157, 0.2);
      }

      .minus { background: #fff; border: 2rpx solid #FF6B9D; }
      .plus  { background: #FF6B9D; }

      .num {
        min-width: 40rpx;
        text-align: center;
        font-size: 28rpx;
        color: #333;
      }
    }
  }
}

.sheet-footer {
  padding: 20rpx 24rpx 30rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  position: sticky;
  bottom: 0;
  background: #fff;

  .amount { font-size: 28rpx; color: #333; }
  .price { color: #FF6B9D; font-weight: 700; }

  .submit {
    background: linear-gradient(135deg, #FF6B9D 0%, #FF8E9E 100%);
    color: #fff;
    border: none;
    border-radius: 30rpx;
    padding: 16rpx 28rpx;
      font-size: 28rpx;
  }
}
</style> 