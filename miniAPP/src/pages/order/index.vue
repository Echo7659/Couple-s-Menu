<template>
  <view class="menu-index-container" @tap="handleBackgroundTap">
    <view class="menu-list">
      <view class="menu-grid">
        <view 
          v-for="menu in menuList" 
          :key="menu.ID" 
          class="menu-card"
          :class="{ 'is-active': activeMenu?.ID === menu.ID }"
          @longpress="handleLongPress(menu)"
          @tap.stop="handleCardTap(menu.ID)"
        >
          <view class="menu-cover-wrapper">
            <image 
              v-if="menu.cover"
              :src="menu.cover" 
              class="menu-cover" 
              mode="aspectFill"
            ></image>
            <view v-else class="menu-cover-placeholder">
              <text class="cover-name">{{ menu.name || '未命名菜谱' }}</text>
            </view>
          </view>
          <view class="menu-info">
            <text class="menu-name">{{ menu.name }}</text>
            <text class="menu-desc" v-if="menu.description">{{ menu.description }}</text>
          </view>
        </view>

        <!-- 新增菜谱卡片 -->
        <view class="menu-card add-menu-card" @tap.stop="createMenu">
          <view class="add-icon">
            <uni-icons type="plus" size="44" color="#FF6B9D"></uni-icons>
          </view>
          <text class="add-text">新增菜谱</text>
        </view>
      </view>

      <!-- 空状态 -->
      <view v-if="menuList.length === 0 && !loading" class="empty-state">
        <uni-icons type="folder-add" size="80" color="#ccc"></uni-icons>
        <text class="empty-text">还没有菜谱</text>
        <text class="empty-desc">点击下方按钮创建第一个菜谱吧</text>
      </view>
    </view>

    <!-- 新增菜谱弹窗 -->
    <uni-popup ref="createPopup" type="center">
      <view class="create-dialog">
        <text class="dialog-title">{{ dialogTitle }}</text>
        <uni-easyinput 
          v-model="menuForm.name" 
          placeholder="请输入菜谱名称"
          :clearable="true"
        ></uni-easyinput>
        <uni-easyinput 
          style="margin-top: 20rpx;"
          v-model="menuForm.description" 
          placeholder="请输入菜谱描述（可选）"
          type="textarea"
          :clearable="true"
        ></uni-easyinput>
        <view class="upload-section">
          <text class="upload-label">封面图片（可选）</text>
          <view class="upload-box" @click="chooseCover">
            <image 
              v-if="menuForm.cover" 
              :src="menuForm.cover" 
              class="cover-preview" 
              mode="aspectFill"
            ></image>
            <view v-else class="upload-placeholder">
              <uni-icons type="camera" size="28" color="#999"></uni-icons>
              <text class="upload-text">点击上传</text>
            </view>
          </view>
        </view>
        <view class="dialog-actions">
          <button class="cancel-btn" @click="closeCreate">取消</button>
          <button class="confirm-btn" @click="saveMenu" :disabled="!menuForm.name.trim()">
            {{ isEditing ? '保存' : '创建' }}
          </button>
        </view>
      </view>
    </uni-popup>

    <t-cropper
      v-if="cropperImageUrl"
      :image-url="cropperImageUrl"
      :width="600"
      :height="400"
      :max-width="1200"
      :max-height="1200"
      mode="free"
      @cancel="handleCropperCancel"
      @confirm="handleCropperConfirm"
      class="fullscreen-cropper"
    />

    <view v-if="activeMenu" class="menu-overlay" @tap="clearActiveMenu"></view>

    <view 
      v-if="activeMenu" 
      class="menu-bottom-bar"
      @tap.stop
    >
      <view class="bottom-bar-inner">
        <view class="bottom-menu-title">
          <text class="label">已选菜谱：</text>
          <text class="value">{{ activeMenu.name || '未命名菜谱' }}</text>
        </view>
        <view class="bottom-actions">
          <button class="bottom-btn edit" @tap="handleBottomEdit">
            <uni-icons type="compose" size="16" color="#fff"></uni-icons>
            <text>编辑</text>
          </button>
          <button class="bottom-btn delete" @tap="handleBottomDelete">
            <uni-icons type="trash" size="16" color="#fff"></uni-icons>
            <text>删除</text>
          </button>
          <button class="bottom-btn cancel" @tap="clearActiveMenu">
            <text>取消</text>
          </button>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { onShow, onPullDownRefresh } from '@dcloudio/uni-app'
import api from '@/api'
import config from '@/config'
import { useUserStore } from '@/store/user.js'
import TCropper from '@/components/t-cropper/t-cropper.vue'

const userStore = useUserStore()

const menuList = ref([])
const loading = ref(false)
const createPopup = ref(null)
const menuForm = ref({
  id: null,
  name: '',
  description: '',
  cover: ''
})
const editingMenuId = ref(null)
const isEditing = computed(() => Boolean(editingMenuId.value))
const dialogTitle = computed(() => (isEditing.value ? '编辑菜谱' : '新增菜谱'))
const cropperImageUrl = ref('')
const cropperTarget = ref(null)
const activeMenu = ref(null)
const isCropping = computed(() => !!cropperImageUrl.value) // 是否正在裁剪图片

// 加载菜谱列表
const loadMenus = async () => {
  try {
    loading.value = true
    const res = await api.getMenuList({ page: 1, pageSize: 100 })
    if (res?.code === 0) {
      menuList.value = res.data?.list || []
    }
  } catch (error) {
    console.error('加载菜谱列表失败:', error)
    uni.showToast({
      title: '加载失败',
      icon: 'none'
    })
  } finally {
    loading.value = false
  }
}

const handleCardTap = (menuId) => {
  if (activeMenu.value) {
    activeMenu.value = null
    return
  }
  selectMenu(menuId)
}

// 选择菜谱
const selectMenu = (menuId) => {
  uni.navigateTo({
    url: `/pages/order/order?menuId=${menuId}`
  })
}

const handleBackgroundTap = () => {
  if (activeMenu.value) {
    clearActiveMenu()
  }
}

// 长按菜谱
const handleLongPress = (menu) => {
  activeMenu.value = menu
}

// 打开创建弹窗
const createMenu = () => {
  // 检查登录状态
  if (!userStore.hasLogin) {
    uni.showToast({
      title: '请先登录',
      icon: 'none'
    })
    setTimeout(() => {
      uni.navigateTo({ url: '/pages/user/login' })
    }, 1500)
    return
  }
  
  menuForm.value = {
    id: null,
    name: '',
    description: '',
    cover: ''
  }
  editingMenuId.value = null
  createPopup.value.open()
}

// 关闭创建弹窗
const closeCreate = () => {
  createPopup.value.close()
  editingMenuId.value = null
}

// 编辑菜谱
const editMenu = (menu) => {
  menuForm.value = {
    id: menu.ID,
    name: menu.name || '',
    description: menu.description || '',
    cover: menu.cover || ''
  }
  editingMenuId.value = menu.ID
  createPopup.value.open()
}

const clearActiveMenu = () => {
  activeMenu.value = null
}

const handleBottomEdit = () => {
  if (!activeMenu.value) return
  editMenu(activeMenu.value)
}

const handleBottomDelete = () => {
  if (!activeMenu.value) return
  confirmDeleteMenu(activeMenu.value)
}

// 处理滑动操作
// 删除菜谱确认
const confirmDeleteMenu = (menu) => {
  uni.showModal({
    title: '删除菜谱',
    content: `确定删除「${menu.name || '该菜谱'}」吗？`,
    confirmText: '删除',
    confirmColor: '#FF4D4F',
    success: (res) => {
      if (res.confirm) {
        deleteMenuById(menu.ID)
      }
    }
  })
}

// 删除菜谱
const deleteMenuById = async (menuId) => {
  try {
    uni.showLoading({ title: '删除中...' })
    const res = await api.deleteMenu({ id: menuId })
    if (res?.code === 0) {
      uni.showToast({
        title: '删除成功',
        icon: 'success'
      })
      loadMenus()
    } else {
      uni.showToast({
        title: res?.msg || '删除失败',
        icon: 'none'
      })
    }
  } catch (error) {
    console.error('删除菜谱失败:', error)
    uni.showToast({
      title: '删除失败',
      icon: 'none'
    })
  } finally {
    uni.hideLoading()
  }
}

// 选择封面图片
const chooseCover = () => {
  uni.chooseImage({
    count: 1,
    sizeType: ['compressed'],
    sourceType: ['album', 'camera'],
    success: (res) => {
      const tempPath = res.tempFilePaths[0]
      if (!tempPath) return
      cropperTarget.value = 'menuCover'
      cropperImageUrl.value = tempPath
      // 禁用下拉刷新
      uni.pageScrollTo({
        scrollTop: 0,
        duration: 0
      })
    }
  })
}

const handleCropperCancel = () => {
  cropperImageUrl.value = ''
  cropperTarget.value = null
  // 重新启用下拉刷新
}

const handleCropperConfirm = async (res) => {
  const croppedPath = res?.tempFilePath || res?.path
  if (!croppedPath) {
    handleCropperCancel()
    return
  }
  await uploadImageByTarget(croppedPath, cropperTarget.value)
  handleCropperCancel()
}

const uploadImageByTarget = async (path, target) => {
  if (!target) return
  const url = await uploadFileToServer(path)
  if (url && target === 'menuCover') {
    menuForm.value.cover = url
  }
}

const uploadFileToServer = async (filePath) => {
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
          try { resolve(JSON.parse(r.data)) } catch { resolve({ code: -1 }) }
        },
        fail: reject
      })
    })
    if (uploadRes?.code === 0 && uploadRes.data?.file?.url) {
      return uploadRes.data.file.url
    }
    uni.showToast({ title: uploadRes?.msg || '上传失败', icon: 'none' })
    return ''
  } catch (error) {
    console.error('上传失败:', error)
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

// 保存菜谱
const saveMenu = async () => {
  if (!menuForm.value.name.trim()) {
    uni.showToast({
      title: '请输入菜谱名称',
      icon: 'none'
    })
    return
  }
  try {
    uni.showLoading({ title: isEditing.value ? '保存中...' : '创建中...' })
    const payload = {
      name: menuForm.value.name.trim(),
      description: menuForm.value.description.trim() || undefined,
      cover: menuForm.value.cover || undefined
    }
    if (isEditing.value && editingMenuId.value) {
      payload.id = editingMenuId.value
    }
    const requestFn = isEditing.value ? api.updateMenu : api.createMenu
    const res = await requestFn(payload)
    if (res?.code === 0) {
      uni.showToast({
        title: isEditing.value ? '更新成功' : '创建成功',
        icon: 'success'
      })
      closeCreate()
      loadMenus()
    } else {
      uni.showToast({
        title: res?.msg || (isEditing.value ? '更新失败' : '创建失败'),
        icon: 'none'
      })
    }
  } catch (error) {
    console.error('保存菜谱失败:', error)
    uni.showToast({
      title: isEditing.value ? '更新失败' : '创建失败',
      icon: 'none'
    })
  } finally {
    uni.hideLoading()
  }
}

// 生命周期
onMounted(() => {
  loadMenus()
})

onShow(() => {
  // 页面显示时刷新列表
  loadMenus()
})

onPullDownRefresh(() => {
  // 如果正在裁剪图片，不执行下拉刷新
  if (isCropping.value) {
    uni.stopPullDownRefresh()
    return
  }
  loadMenus().finally(() => {
    uni.stopPullDownRefresh()
  })
})
</script>

<style lang="scss" scoped>
.menu-index-container {
  min-height: 100vh;
  background: #f5f5f5;
  padding: 24rpx;
  box-sizing: border-box;
}

.menu-list {
  display: flex;
  flex-direction: column;
  gap: 20rpx;
}

.menu-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  column-gap: 28px;
  row-gap: 32px;
  max-width: 700rpx;
  width: 100%;
  margin: 0 auto;
  padding: 0 10px 40px;
  box-sizing: border-box;
  justify-items: center;
}

.menu-card {
  background: #fff;
  border-radius: 18rpx;
  padding: 14rpx;
  display: flex;
  flex-direction: column;
  align-items: stretch;
  box-shadow: none;
  transition: transform 0.15s ease;
  width: 100%;
  min-height: 360rpx;

  &:active {
    transform: translateY(4rpx);
  }

  &.is-active {
    box-shadow: 0 0 0 3rpx rgba(255, 107, 157, 0.35);
  }

  .menu-cover-wrapper {
    width: 100%;
    height: 280rpx;
    border-radius: 18rpx;
    overflow: hidden;
    background: #f0f0f0;
    margin-bottom: 18rpx;

    .menu-cover {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }

    .menu-cover-placeholder {
      width: 100%;
      height: 100%;
      border-radius: 18rpx;
      background: linear-gradient(135deg, #f8d7ff 0%, #e3f0ff 100%);
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 30rpx;
      text-align: center;

      .cover-name {
        font-size: 28rpx;
        color: #6A4C93;
        font-weight: 600;
        line-height: 1.5;
      }
    }
  }

  .menu-info {
    display: flex;
    flex-direction: column;
    gap: 10rpx;
    align-items: center;
    text-align: center;

    .menu-name {
      font-size: 30rpx;
      font-weight: 600;
      color: #2E2E3A;
    }

    .menu-desc {
      font-size: 22rpx;
      color: #9c9ca5;
      line-height: 1.4;
      width: 100%;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
  }
}

.add-menu-card {
  border: none;
  background: rgba(255, 255, 255, 0.95);
  justify-content: center;
  align-items: center;
  flex-direction: column;
  gap: 18rpx;
  padding: 40rpx 20rpx;
  width: 100%;
  min-height: 360rpx;

  .add-icon {
    width: 92rpx;
    height: 92rpx;
    border-radius: 50%;
    background: rgba(255, 107, 157, 0.1);
    display: flex;
    align-items: center;
    justify-content: center;
  }

.add-text {
    font-size: 28rpx;
    color: #FF6B9D;
    font-weight: 500;
    text-align: center;
  }
}

.menu-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.18);
  z-index: 15;
}

.menu-bottom-bar {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  padding: env(safe-area-inset-bottom, 0);
  background: linear-gradient(180deg, rgba(255, 255, 255, 0) 0%, #ffffff 40%);
  box-shadow: 0 -6rpx 20rpx rgba(0, 0, 0, 0.06);
  z-index: 20;
}

.bottom-bar-inner {
  padding: 16rpx 24rpx 24rpx;
}

.bottom-menu-title {
  display: flex;
  align-items: baseline;
  margin-bottom: 16rpx;

  .label {
    font-size: 24rpx;
    color: #999;
  }

  .value {
    margin-left: 8rpx;
    font-size: 26rpx;
    color: #333;
    font-weight: 600;
  }
}

.bottom-actions {
  display: flex;
  gap: 12rpx;
}

.bottom-btn {
  flex: 1;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 4rpx;
  padding: 12rpx 0;
  border-radius: 999rpx;
  border: none;
  font-size: 24rpx;
  font-weight: 500;
}

.bottom-btn.edit {
  background: linear-gradient(135deg, #8A63F3 0%, #B089FF 100%);
  color: #fff;
}

.bottom-btn.delete {
  background: linear-gradient(135deg, #FF8FA3 0%, #FF6B9D 100%);
  color: #fff;
}

.bottom-btn.cancel {
  background: #f0f0f0;
  color: #666;
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

  .empty-desc {
    font-size: 24rpx;
    color: #ccc;
  }
}

// 全屏裁剪组件样式
.fullscreen-cropper {
  position: fixed !important;
  top: 0 !important;
  left: 0 !important;
  right: 0 !important;
  bottom: 0 !important;
  width: 100vw !important;
  height: 100vh !important;
  z-index: 9999 !important;
  margin: 0 !important;
  padding: 0 !important;
}

.create-dialog {
  background: linear-gradient(135deg, #fff 0%, #faf5ff 100%);
  padding: 50rpx 40rpx;
  border-radius: 24rpx;
  width: 680rpx;
  max-height: 85vh;
  overflow-y: auto;
  box-shadow: 0 8rpx 32rpx rgba(106, 76, 147, 0.15);

  .dialog-title {
    font-size: 38rpx;
    font-weight: 700;
    color: #6A4C93;
    margin-bottom: 40rpx;
    display: block;
    text-align: center;
  }

  .upload-section {
    margin-top: 30rpx;

    .upload-label {
      font-size: 28rpx;
      color: #6A4C93;
      margin-bottom: 16rpx;
      display: block;
      font-weight: 600;
    }

    .upload-box {
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

      .cover-preview {
        width: 100%;
        height: 100%;
        object-fit: contain;
        border-radius: 12rpx;
      }

      .upload-placeholder {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 16rpx;
        color: #999;

        .upload-text {
          font-size: 26rpx;
          color: #6A4C93;
        }
      }
    }
  }

  .dialog-actions {
    display: flex;
    gap: 20rpx;
    margin-top: 40rpx;
    justify-content: center;

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
  }
}
</style>

