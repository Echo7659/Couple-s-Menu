<template>
  <view class="container">
    <view class="header">
      <text class="title">菜品管理</text>
      <button class="add" @click="startAdd">新增菜品</button>
    </view>

    <!-- 菜谱选择 -->
    <view class="selectors">
      <view class="selector-item">
        <text class="selector-label">菜谱：</text>
        <picker :range="menuList" range-key="name" @change="onMenuChange">
          <view class="picker-view">
            <text>{{ selectedMenuName || '请选择菜谱' }}</text>
            <uni-icons type="arrowdown" size="16" color="#999"></uni-icons>
          </view>
        </picker>
      </view>
      <view class="selector-item" v-if="categoryList.length > 0">
        <text class="selector-label">分类：</text>
        <picker 
          :range="categoryFilterList" 
          range-key="name" 
          @change="onCategoryFilterChange"
          :value="categoryFilterPickerValue"
        >
          <view class="picker-view">
            <text :class="{ placeholder: selectedCategoryIdx < 0 }">
              {{ selectedCategoryFilterName }}
            </text>
            <uni-icons type="arrowdown" size="16" color="#999"></uni-icons>
          </view>
        </picker>
      </view>
    </view>

    <scroll-view 
      class="scroll-container" 
      scroll-y 
      @scrolltolower="onScrollToLower"
      :lower-threshold="100"
    >
      <view v-for="item in list" :key="item.ID" class="row">
        <view class="left">
          <image :src="item.img" class="thumb" mode="aspectFill" />
          <view class="meta">
            <text class="name">{{ item.name }}</text>
            <text class="desc">分类：{{ getCategoryName(item.categoryId) }}</text>
          </view>
        </view>
        <view class="ops">
          <button size="mini" class="op-btn edit" @click="startEdit(item)">编辑</button>
          <button size="mini" class="op-btn delete" @click="remove(item)">删除</button>
        </view>
      </view>

      <!-- 加载更多提示 -->
      <view v-if="list.length > 0" class="load-more">
        <uni-load-more 
          :status="loadMoreStatus" 
          :content-text="loadMoreText"
        ></uni-load-more>
      </view>

      <!-- 空状态 -->
      <view v-if="list.length === 0 && !loading && selectedMenuId" class="empty-state">
        <uni-icons type="folder-add" size="80" color="#ccc"></uni-icons>
        <text class="empty-text">该分类下还没有菜品</text>
      </view>
    </scroll-view>

    <uni-popup ref="popup" type="center">
      <view class="dialog">
        <text class="dialog-title">{{ current?.ID ? '编辑菜品' : '新增菜品' }}</text>
        <view class="form-content">
          <view class="form-item">
            <text class="form-label">菜品名称</text>
            <uni-easyinput 
              v-model="form.name" 
              placeholder="请输入菜品名称"
              :clearable="true"
            ></uni-easyinput>
          </view>
          <view class="form-item">
            <text class="form-label">菜品描述</text>
            <uni-easyinput 
              v-model="form.desc" 
              placeholder="请输入菜品描述（可选）"
              type="textarea"
              :clearable="true"
            ></uni-easyinput>
          </view>
          <view class="form-item">
            <text class="form-label">所属分类</text>
            <picker 
              :range="categoryList" 
              range-key="name" 
              @change="onCatChangeInDialog"
              :value="selectedCatIdx"
            >
              <view class="picker-view">
                <text :class="{ placeholder: selectedCatIdx < 0 }">
                  {{ pickCategoryName }}
                </text>
                <uni-icons type="arrowdown" size="16" color="#999"></uni-icons>
              </view>
            </picker>
          </view>
          <view class="form-item">
            <text class="form-label">菜品图片</text>
            <view class="uploader-box" @click="chooseImage">
              <image 
                v-if="form.img" 
                :src="form.img" 
                class="uploader-preview" 
                mode="aspectFill" 
              />
              <view v-else class="uploader-placeholder">
                <uni-icons type="camera" size="32" color="#999"></uni-icons>
                <text class="upload-text">点击上传图片</text>
              </view>
            </view>
            <view v-if="form.img" class="uploader-actions">
              <button size="mini" class="action-btn remove-btn" @click.stop="form.img = ''">移除</button>
              <button size="mini" class="action-btn change-btn" @click.stop="chooseImage">更换</button>
            </view>
          </view>
        </view>
        <view class="dialog-actions">
          <button class="cancel-btn" @click="close">取消</button>
          <button class="confirm-btn" @click="save" :disabled="!form.name.trim() || !form.categoryId">保存</button>
        </view>
      </view>
    </uni-popup>
    <t-cropper
      :image-url="cropperImageUrl"
      :width="600"
      :height="400"
      :max-width="1200"
      :max-height="1200"
      mode="free"
      @cancel="handleCropperCancel"
      @confirm="handleCropperConfirm"
    />
  </view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import api from '@/api'
import config from '@/config'
import TCropper from '@/components/t-cropper/t-cropper.vue'

const list = ref([])
const menuList = ref([])
const categoryList = ref([])
const selectedMenuId = ref(null)
const selectedMenuIdx = ref(-1)
const selectedCategoryId = ref(null)
const selectedCategoryIdx = ref(-1)
const popup = ref(null)
const current = ref(null)
const form = ref({ name: '', desc: '', img: '', categoryId: null })
const selectedCatIdx = ref(-1)
const cropperImageUrl = ref('')

// 分页相关
const page = ref(1)
const pageSize = ref(20)
const hasMore = ref(true)
const loading = ref(false)
const loadingMore = ref(false)

const selectedMenuName = computed(() => {
  if (selectedMenuIdx.value < 0) return ''
  return menuList.value[selectedMenuIdx.value]?.name || ''
})

const selectedCategoryName = computed(() => {
  if (selectedCategoryIdx.value < 0) return ''
  return categoryList.value[selectedCategoryIdx.value]?.name || ''
})

// 分类筛选列表（包含"全部"选项）
const categoryFilterList = computed(() => {
  return [{ ID: null, name: '全部' }, ...categoryList.value]
})

// 选中的分类筛选名称
const selectedCategoryFilterName = computed(() => {
  if (selectedCategoryIdx.value < 0) return '全部'
  const category = categoryList.value[selectedCategoryIdx.value]
  return category?.name || '全部'
})

// 分类筛选的 picker 值（需要包含"全部"选项的索引）
const categoryFilterPickerValue = computed(() => {
  if (selectedCategoryIdx.value < 0) return 0
  return selectedCategoryIdx.value + 1
})

// 加载更多状态
const loadMoreStatus = computed(() => {
  if (loadingMore.value) return 'loading'
  if (!hasMore.value) return 'noMore'
  return 'more'
})

// 加载更多文本
const loadMoreText = computed(() => {
  return {
    contentdown: '上拉加载更多',
    contentrefresh: '正在加载...',
    contentnomore: '没有更多了'
  }
})

const chooseImage = () => {
  uni.chooseImage({
    count: 1,
    sizeType: ['compressed'],
    sourceType: ['album', 'camera'],
    success: (res) => {
      const tempPath = res.tempFilePaths?.[0]
      if (!tempPath) return
      cropperImageUrl.value = tempPath
    }
  })
}

const handleCropperCancel = () => {
  cropperImageUrl.value = ''
}

const handleCropperConfirm = async (res) => {
  const croppedPath = res?.tempFilePath || res?.path
  if (!croppedPath) {
    handleCropperCancel()
    return
  }
  const previousImg = form.value.img
  form.value.img = croppedPath
  const uploadedUrl = await uploadFoodImage(croppedPath)
  if (uploadedUrl) {
    form.value.img = uploadedUrl
  } else {
    form.value.img = previousImg
  }
  handleCropperCancel()
}

const uploadFoodImage = async (filePath) => {
  try {
    uni.showLoading({ title: '上传中...' })
    const compressedPath = await compressImage(filePath)
    const fullUrl = (config.baseUrl || '').replace(/\/$/, '') + '/user/upload'
    const uploadRes = await new Promise((resolve, reject) => {
      uni.uploadFile({
        url: fullUrl,
        filePath: compressedPath,
        name: 'file',
        header: { 'Authorization': uni.getStorageSync('token') || '' },
        success: (r) => {
          try { resolve(JSON.parse(r.data)) } catch (e) { resolve({ code: -1 }) }
        },
        fail: reject
      })
    })
    if (uploadRes?.code === 0 && uploadRes.data?.file?.url) {
      return uploadRes.data.file.url
    }
    uni.showToast({ title: uploadRes?.msg || '上传失败', icon: 'none' })
    return ''
  } catch (e) {
    console.error('上传图片失败:', e)
    uni.showToast({ title: '上传失败', icon: 'none' })
    return ''
  } finally {
    uni.hideLoading()
  }
}

const compressImage = (path) => {
  return new Promise((resolve) => {
    uni.compressImage({
      src: path,
      quality: 70,
      success: (res) => resolve(res.tempFilePath || path),
      fail: (err) => {
        console.warn('图片压缩失败，使用原图:', err)
        resolve(path)
      }
    })
  })
}

// 加载菜谱列表
const loadMenus = async () => {
  try {
    const res = await api.getMenuList({ page: 1, pageSize: 100 })
    if (res?.code === 0) {
      menuList.value = res.data?.list || []
    }
  } catch (e) {
    console.error('加载菜谱列表失败:', e)
  }
}

// 加载分类列表
const loadCategories = async () => {
  if (!selectedMenuId.value) {
    categoryList.value = []
    return
  }
  try {
    const res = await api.getFoodCategoryList({ menuId: selectedMenuId.value, page: 1, pageSize: 100 })
    if (res?.code === 0) {
      categoryList.value = res.data?.list || []
      // 不再默认选择分类，分类在弹窗中选择
      selectedCategoryIdx.value = -1
      selectedCategoryId.value = null
    }
  } catch (e) {
    console.error('加载分类列表失败:', e)
  }
}

// 加载菜品列表
const loadFoods = async (reset = true) => {
  if (!selectedMenuId.value) {
    list.value = []
    hasMore.value = false
    return
  }
  
  if (reset) {
    page.value = 1
    list.value = []
    hasMore.value = true
    loading.value = true
  } else {
    if (!hasMore.value || loadingMore.value) return
    loadingMore.value = true
  }

  try {
    const params = { 
      page: page.value, 
      pageSize: pageSize.value, 
      menuId: selectedMenuId.value 
    }
    
    // 如果选择了分类，添加 categoryId 参数
    if (selectedCategoryId.value) {
      params.categoryId = selectedCategoryId.value
    }
    
    const res = await api.getFoodList(params)
    if (res?.code === 0) {
      const newList = res.data?.list || []
      const total = res.data?.total || 0

      if (reset) {
        list.value = newList
      } else {
        list.value = [...list.value, ...newList]
      }

      // 判断是否还有更多数据
      hasMore.value = list.value.length < total

      // 只有在成功加载数据后才增加页码
      if (!reset && newList.length > 0) {
        page.value += 1
      } else if (reset && hasMore.value) {
        page.value = 2  // 重置后如果有更多数据，下次加载第2页
      }
    }
  } catch (e) {
    console.error('加载菜品列表失败:', e)
    uni.showToast({ title: '加载失败', icon: 'none' })
  } finally {
    loading.value = false
    loadingMore.value = false
  }
}

// 菜谱选择变化
const onMenuChange = (e) => {
  selectedMenuIdx.value = Number(e.detail.value)
  selectedMenuId.value = menuList.value[selectedMenuIdx.value]?.ID || null
  selectedCategoryIdx.value = -1
  selectedCategoryId.value = null
  loadCategories()
  loadFoods(true)
}

// 分类筛选变化
const onCategoryFilterChange = (e) => {
  const index = Number(e.detail.value)
  if (index === 0) {
    // 选择"全部"
    selectedCategoryIdx.value = -1
    selectedCategoryId.value = null
  } else {
    // 选择具体分类
    selectedCategoryIdx.value = index - 1
    selectedCategoryId.value = categoryList.value[index - 1]?.ID || null
  }
  loadFoods(true)
}

// 触底加载更多（scroll-view 的 scrolltolower 事件）
const onScrollToLower = () => {
  if (hasMore.value && !loading.value && !loadingMore.value) {
    loadFoods(false)
  }
}

// 分类选择变化（顶部筛选）- 已移除，分类在弹窗中选择

// 从路由参数获取 menuId
const getMenuIdFromRoute = () => {
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1]
  const options = currentPage.options || {}
  return options.menuId ? Number(options.menuId) : null
}

onMounted(() => {
  loadMenus().then(() => {
    // 优先使用路由参数中的 menuId
    const routeMenuId = getMenuIdFromRoute()
    if (routeMenuId) {
      const menuIndex = menuList.value.findIndex(m => m.ID === routeMenuId)
      if (menuIndex >= 0) {
        selectedMenuIdx.value = menuIndex
        selectedMenuId.value = routeMenuId
        loadCategories()
        loadFoods(true)
      } else if (menuList.value.length > 0) {
        // 如果路由参数中的菜谱不存在，默认选择第一个
        selectedMenuIdx.value = 0
        selectedMenuId.value = menuList.value[0].ID
        loadCategories()
        loadFoods(true)
      }
    } else if (menuList.value.length > 0) {
      // 如果没有路由参数，默认选择第一个
      selectedMenuIdx.value = 0
      selectedMenuId.value = menuList.value[0].ID
      loadCategories()
      loadFoods(true)
    }
  })
})

const startAdd = () => {
  if (!selectedMenuId.value) {
    uni.showToast({ title: '请先选择菜谱', icon: 'none' })
    return
  }
  if (categoryList.value.length === 0) {
    uni.showToast({ title: '该菜谱下没有分类，请先创建分类', icon: 'none' })
    return
  }
  current.value = null
  form.value = { name: '', desc: '', img: '', categoryId: null }
  selectedCatIdx.value = -1
  popup.value.open()
}

const startEdit = (item) => {
  current.value = item
  form.value = { name: item.name, desc: item.desc, img: item.img, categoryId: item.categoryId }
  selectedCatIdx.value = categoryList.value.findIndex(c => c.ID === item.categoryId)
  popup.value.open()
}

const close = () => popup.value.close()

// 对话框中的分类选择
const onCatChangeInDialog = (e) => {
  selectedCatIdx.value = Number(e.detail.value)
  form.value.categoryId = categoryList.value[selectedCatIdx.value]?.ID || null
}

const pickCategoryName = computed(() => {
  if (selectedCatIdx.value < 0) return '未选择'
  return categoryList.value[selectedCatIdx.value]?.name || '未选择'
})

const getCategoryName = (id) => categoryList.value.find(c => c.ID === id)?.name || '—'

const save = async () => {
  if (!form.value.name.trim()) {
    uni.showToast({ title: '请输入菜品名称', icon: 'none' })
    return
  }
  if (!form.value.categoryId) {
    uni.showToast({ title: '请选择分类', icon: 'none' })
    return
  }
  try {
    if (current.value?.ID) {
      await api.updateFood({ ID: current.value.ID, ...form.value })
      uni.showToast({ title: '更新成功', icon: 'success' })
    } else {
      await api.addFood(form.value)
      uni.showToast({ title: '创建成功', icon: 'success' })
    }
    close()
    loadFoods(true)
  } catch (e) {
    console.error('保存菜品失败:', e)
    uni.showToast({ title: '操作失败', icon: 'none' })
  }
}

const remove = (item) => {
  uni.showModal({
    title: '删除确认',
    content: `确定删除"${item.name}"吗？`,
    success: async (res) => {
      if (res.confirm) {
        try {
          await api.deleteFood({ ID: item.ID })
          uni.showToast({ title: '删除成功', icon: 'success' })
          loadFoods(true)
        } catch (e) {
          console.error('删除菜品失败:', e)
          uni.showToast({ title: '删除失败', icon: 'none' })
        }
      }
    }
  })
}
</script>

<style lang="scss" scoped>
.container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: linear-gradient(135deg, #FFE4F0 0%, #F3E8FF 50%, #EDE7FF 100%);
}

.scroll-container {
  flex: 1;
  padding: 0 24rpx 24rpx;
  box-sizing: border-box;
  height: 0; /* 确保 flex: 1 正确计算高度 */
}
.header { 
  display: flex; 
  justify-content: space-between; 
  align-items: center; 
  padding: 24rpx 24rpx 0;
  flex-shrink: 0;
}
.title { font-weight: 700; font-size: 34rpx; color: #6A4C93; }
.add {margin-left: 220px; background: linear-gradient(135deg, #FF7EB3 0%, #FF758C 100%); color:#fff; border:none; border-radius: 28rpx; padding: 16rpx 26rpx; font-size: 26rpx; box-shadow: 0 6rpx 16rpx rgba(255, 117, 140, 0.35); }
.selectors { 
  background: rgba(255,255,255,0.9); 
  padding: 20rpx; 
  border-radius: 18rpx; 
  margin: 24rpx 24rpx 0;
  display: flex; 
  flex-direction: column; 
  gap: 16rpx; 
  box-shadow: 0 8rpx 24rpx rgba(140, 82, 255, 0.12);
  flex-shrink: 0;
}
.selector-item { display: flex; align-items: center; gap: 16rpx; }
.selector-label { font-size: 28rpx; color: #6A4C93; font-weight: 600; min-width: 120rpx; }
.picker-view { flex: 1; display: flex; align-items: center; justify-content: space-between; padding: 12rpx 20rpx; background: #faf7ff; border-radius: 12rpx; color: #3C2A4D; }
.row { background: rgba(255,255,255,0.9); padding: 26rpx; border-radius: 18rpx; display:flex; justify-content: space-between; align-items:center; margin-bottom: 18rpx; box-shadow: 0 8rpx 24rpx rgba(140, 82, 255, 0.12); }
.left { display:flex; align-items:center; }
.thumb { width: 120rpx; height: 120rpx; border-radius: 12rpx; margin-right: 16rpx; }
.meta { .name{display:block; font-size:28rpx; font-weight:600; color:#3C2A4D;} .desc{color:#8c8c8c; margin-top:6rpx;} }
.ops { display:flex; gap: 12rpx; }
.op-btn { border:none; padding: 12rpx 20rpx; border-radius: 22rpx; font-size: 24rpx; color:#fff; }
.op-btn.edit { background: linear-gradient(135deg, #8A63F3 0%, #B089FF 100%); box-shadow: 0 4rpx 12rpx rgba(138,99,243,.3); }
.op-btn.delete { background: linear-gradient(135deg, #FF8FA3 0%, #FF6B9D 100%); box-shadow: 0 4rpx 12rpx rgba(255,107,157,.3); }
.dialog { 
  background: linear-gradient(135deg, #fff 0%, #faf5ff 100%);
  padding: 50rpx 40rpx;
  border-radius: 24rpx;
  width: 680rpx;
  max-height: 85vh;
  overflow-y: auto;
  box-shadow: 0 8rpx 32rpx rgba(106, 76, 147, 0.15);
}
.dialog-title { 
  font-size: 38rpx;
  font-weight: 700;
  margin-bottom: 40rpx;
  display: block;
  color: #6A4C93;
  text-align: center;
}

.form-content {
  margin-bottom: 30rpx;
}

.form-item {
  margin-bottom: 30rpx;
  
  :deep(.uni-easyinput__content) {
    min-height: 88rpx;
  }
  
  :deep(.uni-easyinput__content-textarea) {
    min-height: 160rpx;
  }
}

.form-label {
  display: block;
  font-size: 28rpx;
  color: #6A4C93;
  margin-bottom: 12rpx;
  font-weight: 600;
}

.picker-view {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20rpx;
  background: #faf7ff;
  border-radius: 12rpx;
  color: #3C2A4D;
  border: 2rpx solid #e8dfff;
  
  .placeholder {
    color: #999;
  }
}

.uploader-box { 
  width: 100%; 
  height: 300rpx; 
  border: 3rpx dashed #d8c7ff; 
  border-radius: 16rpx; 
  display: flex;
  align-items: center;
  justify-content: center; 
  background: linear-gradient(135deg, #faf7ff 0%, #f3e8ff 100%);
  overflow: hidden;
  transition: all 0.3s ease;

  &:active {
    border-color: #FF6B9D;
    background: linear-gradient(135deg, #ffe4f0 0%, #faf7ff 100%);
  }
}

.uploader-placeholder { 
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12rpx;
  color: #999;
  
  .upload-text {
    font-size: 26rpx;
    color: #6A4C93;
  }
}

.uploader-preview { 
  width: 100%; 
  height: 100%; 
  object-fit: contain;
  border-radius: 12rpx;
}

.uploader-actions { 
  display: flex;
  gap: 12rpx;
  justify-content: flex-end;
  margin-top: 12rpx;
}

.action-btn {
  padding: 10rpx 20rpx;
  border-radius: 20rpx;
  font-size: 24rpx;
  border: none;
}

.remove-btn {
  background: #f5f5f5;
  color: #666;
}

.change-btn {
  background: linear-gradient(135deg, #8A63F3 0%, #B089FF 100%);
  color: #fff;
}

.dialog-actions { 
  display: flex;
  justify-content: center;
  gap: 20rpx;
  margin-top: 30rpx;
}

.cancel-btn,
.confirm-btn {
  flex: 1;
  padding: 20rpx 40rpx;
  border-radius: 30rpx;
  font-size: 30rpx;
  font-weight: 600;
  border: none;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.1);
}

.cancel-btn {
  background: #f5f5f5;
  color: #666;
}

.confirm-btn {
  background: linear-gradient(135deg, #FF6B9D 0%, #FF8E9E 100%);
  color: #fff;
  box-shadow: 0 6rpx 20rpx rgba(255, 107, 157, 0.3);

  &:disabled {
    background: #ccc;
    color: #999;
    box-shadow: none;
  }
}

.load-more {
  padding: 30rpx 0;
  display: flex;
  justify-content: center;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 100rpx 0;
  gap: 20rpx;

  .empty-text {
    font-size: 28rpx;
    color: #999;
  }
}
</style> 