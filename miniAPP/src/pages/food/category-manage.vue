<template>
  <view class="container">
    <view class="header">
      <text class="title">分类管理</text>
      <button class="add" @click="startAdd">新增分类</button>
    </view>

    <!-- 菜谱选择 -->
    <view class="menu-selector">
      <text class="selector-label">选择菜谱：</text>
      <picker :range="menuList" range-key="name" @change="onMenuChange">
        <view class="picker-view">
          <text :class="{ placeholder: !selectedMenuName }">
            {{ selectedMenuName || '请选择菜谱' }}
          </text>
          <uni-icons type="arrowdown" size="16" color="#999"></uni-icons>
        </view>
      </picker>
    </view>

    <view v-for="item in list" :key="item.ID" class="row">
      <text class="name">{{ item.name }}</text>
      <view class="ops">
        <button size="mini" class="op-btn edit" @click="startEdit(item)">编辑</button>
        <button size="mini" class="op-btn delete" @click="remove(item)">删除</button>
      </view>
    </view>

    <uni-popup ref="popup" type="center">
      <view class="dialog">
        <text class="dialog-title">{{ current?.ID ? '编辑分类' : '新增分类' }}</text>
        <view class="form-content">
          <view class="form-item">
            <text class="form-label">分类名称</text>
            <uni-easyinput 
              v-model="form.name" 
              placeholder="请输入分类名称"
              :clearable="true"
            ></uni-easyinput>
          </view>
        </view>
        <view class="dialog-actions">
          <button class="cancel-btn" @click="close">取消</button>
          <button class="confirm-btn" @click="save" :disabled="!form.name.trim()">保存</button>
        </view>
      </view>
    </uni-popup>
  </view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import api from '@/api'

const list = ref([])
const menuList = ref([])
const selectedMenuId = ref(null)
const selectedMenuIdx = ref(-1)
const popup = ref(null)
const current = ref(null)
const form = ref({ name: '', menuId: null })

const selectedMenuName = computed(() => {
  if (selectedMenuIdx.value < 0) return ''
  return menuList.value[selectedMenuIdx.value]?.name || ''
})

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
const load = async () => {
  if (!selectedMenuId.value) {
    list.value = []
    return
  }
  try {
    const res = await api.getFoodCategoryList({ menuId: selectedMenuId.value, page: 1, pageSize: 100 })
    if (res?.code === 0) list.value = res.data?.list || []
  } catch (e) {
    console.error('加载分类列表失败:', e)
  }
}

// 菜谱选择变化
const onMenuChange = (e) => {
  selectedMenuIdx.value = Number(e.detail.value)
  selectedMenuId.value = menuList.value[selectedMenuIdx.value]?.ID || null
  load()
}

// 从路由参数获取 menuId
const getMenuIdFromRoute = () => {
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1]
  const options = currentPage.options || {}
  return options.menuId ? Number(options.menuId) : null
}

onMounted(async () => {
  await loadMenus()
  
  // 优先使用路由参数中的 menuId
  const routeMenuId = getMenuIdFromRoute()
  if (routeMenuId) {
    const menuIndex = menuList.value.findIndex(m => m.ID === routeMenuId)
    if (menuIndex >= 0) {
      selectedMenuIdx.value = menuIndex
      selectedMenuId.value = routeMenuId
      load()
      return
    }
  }
  
  // 如果没有路由参数，默认选择第一个
  if (menuList.value.length > 0) {
    selectedMenuIdx.value = 0
    selectedMenuId.value = menuList.value[0].ID
    load()
  }
})

const startAdd = () => {
  if (!selectedMenuId.value) {
    uni.showToast({ title: '请先选择菜谱', icon: 'none' })
    return
  }
  current.value = null
  form.value = { name: '', menuId: selectedMenuId.value }
  popup.value.open()
}

const startEdit = (item) => {
  current.value = item
  form.value = { name: item.name, menuId: item.menuId }
  popup.value.open()
}

const close = () => popup.value.close()

const save = async () => {
  if (!form.value.name.trim()) {
    uni.showToast({ title: '请输入分类名称', icon: 'none' })
    return
  }
  if (!form.value.menuId) {
    uni.showToast({ title: '请选择菜谱', icon: 'none' })
    return
  }
  try {
    if (current.value?.ID) {
      await api.updateFoodCategory({ ID: current.value.ID, name: form.value.name, menuId: form.value.menuId })
      uni.showToast({ title: '更新成功', icon: 'success' })
    } else {
      await api.addFoodCategory({ name: form.value.name, menuId: form.value.menuId })
      uni.showToast({ title: '创建成功', icon: 'success' })
    }
    close()
    load()
  } catch (e) {
    console.error('保存分类失败:', e)
    uni.showToast({ title: '操作失败', icon: 'none' })
  }
}

const remove = (item) => {
  uni.showModal({
    title: '删除确认',
    content: `确定删除分类"${item.name}"吗？`,
    success: async (res) => {
      if (res.confirm) {
        try {
          await api.deleteFoodCategory({ ID: item.ID })
          uni.showToast({ title: '删除成功', icon: 'success' })
          load()
        } catch (e) {
          console.error('删除分类失败:', e)
          uni.showToast({ title: '删除失败', icon: 'none' })
        }
      }
    }
  })
}
</script>

<style lang="scss" scoped>
.container {
  min-height: 100vh;
  padding: 24rpx;
  background: linear-gradient(135deg, #FFE4F0 0%, #F3E8FF 50%, #EDE7FF 100%);
}
.header {
  display: flex; justify-content: space-between; align-items: center; margin-bottom: 24rpx;
}
.title {
  font-weight: 700; font-size: 34rpx; color: #6A4C93;
}
.add {
  background: linear-gradient(135deg, #FF7EB3 0%, #FF758C 100%);
  color:#fff; border:none; border-radius: 28rpx; padding: 16rpx 26rpx; font-size: 26rpx;
  box-shadow: 0 6rpx 16rpx rgba(255, 117, 140, 0.35);
}
.menu-selector {
  background: rgba(255, 255, 255, 0.9);
  padding: 24rpx;
  border-radius: 18rpx;
  margin-bottom: 24rpx;
  box-shadow: 0 8rpx 24rpx rgba(140, 82, 255, 0.12);
  display: flex;
  align-items: center;
  gap: 16rpx;
  
  .selector-label {
    font-size: 28rpx;
    color: #6A4C93;
    font-weight: 600;
    white-space: nowrap;
  }
  
  .picker-view {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16rpx 20rpx;
    background: #faf7ff;
    border-radius: 12rpx;
    color: #3C2A4D;
    border: 2rpx solid #e8dfff;
    
    .placeholder {
      color: #999;
    }
  }
}
.row {
  background: rgba(255,255,255,0.9);
  padding: 26rpx; border-radius: 18rpx; display:flex; justify-content: space-between; align-items:center; margin-bottom: 18rpx;
  box-shadow: 0 8rpx 24rpx rgba(140, 82, 255, 0.12);
}
.name { font-size: 30rpx; color: #3C2A4D; font-weight: 600; }
.ops { display:flex; gap: 12rpx; }
.op-btn {
  border: none; padding: 12rpx 20rpx; border-radius: 22rpx; font-size: 24rpx; color:#fff;
}
.op-btn.edit {
  background: linear-gradient(135deg, #8A63F3 0%, #B089FF 100%);
  box-shadow: 0 4rpx 12rpx rgba(138, 99, 243, 0.3);
}
.op-btn.delete {
  background: linear-gradient(135deg, #FF8FA3 0%, #FF6B9D 100%);
  box-shadow: 0 4rpx 12rpx rgba(255, 107, 157, 0.3);
}
.dialog { 
  background: linear-gradient(135deg, #fff 0%, #faf5ff 100%);
  padding: 50rpx 40rpx;
  border-radius: 24rpx;
  width: 680rpx;
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
  margin-bottom: 20rpx;
  
  :deep(.uni-easyinput__content) {
    min-height: 88rpx;
  }
}

.form-label {
  display: block;
  font-size: 28rpx;
  color: #6A4C93;
  margin-bottom: 12rpx;
  font-weight: 600;
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
</style> 