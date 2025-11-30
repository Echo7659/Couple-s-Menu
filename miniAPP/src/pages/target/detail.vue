<template>
  <view class="detail-container">
    <view v-if="loading" class="loading">
      <text>加载中...</text>
    </view>

    <view v-else-if="target" class="content">
      <!-- 顶部信息 + 操作 -->
      <view class="top-section">
        <!-- 目标信息卡片 -->
        <view class="target-card">
          <view class="title-row">
            <text class="title">{{ target.title }}</text>
            <text v-if="target.isCompleted" class="badge finished">已完成</text>
            <text v-else class="badge ongoing">进行中</text>
          </view>
          <text v-if="target.description" class="desc">
            {{ target.description }}
          </text>
          <view class="meta-row">
            <text class="meta-text">
              创建时间：{{ formatDate(target.CreatedAt || target.createdAt) }}
            </text>
          </view>
        </view>

        <!-- 操作区域 -->
        <view class="actions-card">
          <button
            class="btn-complete"
            v-if="!target.isCompleted"
            type="default"
            size="mini"
            :loading="completing"
            @click="markCompleted"
          >
            标记为已完成
          </button>
          <button
            class="btn-delete"
            type="warn"
            size="mini"
            :loading="deleting"
            @click="deleteTarget"
          >
            删除目标
          </button>
        </view>
      </view>

      <!-- 打卡记录 & 打卡表单 -->
      <view class="records-section">
        <view class="section-header">
          <text class="section-title">打卡记录</text>
          <text class="section-sub" v-if="records.length">
            共 {{ records.length }} 条
          </text>
        </view>

        <view class="record-entry">
          <button
            class="btn-record"
            type="primary"
            :loading="recording"
            :disabled="target.isCompleted"
            @click="handleGoRecord"
          >
            {{ target.isCompleted ? '目标已完成，无法打卡' : '去打卡' }}
          </button>
          <text v-if="target.isCompleted" class="record-tip">
            该目标已完成，不能再新增打卡记录
          </text>
          <text v-else class="record-tip">
            完成后记得来这里打卡～
          </text>
        </view>

        <!-- 内嵌打卡表单 -->
        <view v-if="showRecordForm" class="record-form-card">
          <view class="form-item">
            <text class="label">备注</text>
            <uni-easyinput
              v-model="recordForm.recordNote"
              type="textarea"
              :maxlength="200"
              placeholder="简单记录一下此刻的想法～"
              autoHeight
            />
          </view>
          <view class="form-item">
            <text class="label">图片（可选）</text>
            <view class="img-row">
              <image
                v-if="recordForm.recordImg"
                :src="recordForm.recordImg"
                class="preview-img"
                mode="aspectFill"
                @click="previewImage(recordForm.recordImg)"
              />
              <view v-else class="img-placeholder" @click="chooseImage">
                <uni-icons type="camera" size="30" color="#999" />
                <text class="img-tip">点击选择图片</text>
              </view>
            </view>
          </view>
          <view class="form-actions">
            <button class="btn-cancel" @click="cancelRecord">取消</button>
            <button
              class="btn-confirm"
              type="primary"
              :loading="recording"
              :disabled="!canSubmitRecord"
              @click="submitRecord"
            >
              确认打卡
            </button>
          </view>
        </view>

        <view v-if="records.length > 0" class="record-list">
          <view v-for="item in records" :key="item.ID" class="record-item">
            <view class="record-header">
              <image
                :src="item.user?.avatar || '/static/images/profile.png'"
                class="avatar"
                mode="aspectFill"
              />
              <view class="user-info">
                <text class="name">
                  {{ item.user?.nick_name || item.user?.nickname || '用户' }}
                </text>
                <text class="time">
                  {{ formatDateTime(item.CreatedAt || item.createdAt) }}
                </text>
              </view>
            </view>
            <view class="record-body">
              <text v-if="item.recordNote" class="note">
                {{ item.recordNote }}
              </text>
              <image
                v-if="item.recordImg"
                :src="item.recordImg"
                class="record-img"
                mode="aspectFill"
                @click="previewImage(item.recordImg)"
              />
            </view>
          </view>
        </view>

        <view v-else class="empty-record">
          <text class="empty-text">还没有打卡记录</text>
          <text class="empty-desc">完成后记得来打卡哦～</text>
        </view>
      </view>
    </view>

    <view v-else class="empty-state">
      <text class="empty-text">目标不存在或已被删除</text>
    </view>
  </view>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import api from '@/api'
import config from '@/config'

const target = ref(null)
const records = ref([])
const targetId = ref('')
const loading = ref(true)
const completing = ref(false)
const deleting = ref(false)
const recording = ref(false)
const showRecordForm = ref(false)
const recordForm = ref({
  recordNote: '',
  recordImg: ''
})

const canSubmitRecord = computed(() => {
  return !!recordForm.value.recordNote.trim() && !recording.value
})

// 获取当前用户 ID
const getUserId = async () => {
  const c1 = uni.getStorageSync('userInfo')
  if (c1?.ID) return c1.ID
  const res = await api.getUserInfo()
  return res?.data?.ID || null
}

// 加载目标详情（含记录）
const loadDetail = async () => {
  if (!targetId.value) return
  try {
    loading.value = true
    const res = await api.getTargetDetail(targetId.value)
    if (res?.code === 0 && res.data) {
      target.value = res.data
      // 文档中 records 可能包含完整记录
      if (Array.isArray(res.data.records)) {
        records.value = res.data.records
      } else {
        // 兜底再拉一次记录列表
        const recRes = await api.getTargetRecords(targetId.value)
        if (recRes?.code === 0 && Array.isArray(recRes.data)) {
          records.value = recRes.data
        } else {
          records.value = []
        }
      }
    } else {
      uni.showToast({
        title: res?.msg || '加载失败',
        icon: 'none'
      })
    }
  } catch (e) {
    console.error('加载打卡目标详情失败:', e)
    uni.showToast({
      title: '加载失败',
      icon: 'none'
    })
  } finally {
    loading.value = false
  }
}

// 完成目标
const markCompleted = () => {
  if (!target.value?.ID || completing.value) return
  uni.showModal({
    title: '提示',
    content: '确定将该目标标记为已完成吗？',
    success: async (res) => {
      if (!res.confirm) return
      try {
        completing.value = true
        const result = await api.completeTarget({ id: target.value.ID })
        if (result?.code === 0) {
          uni.showToast({
            title: '已标记为完成',
            icon: 'success'
          })
          // 更新本地状态
          target.value.isCompleted = true
        } else {
          uni.showToast({
            title: result?.msg || '操作失败',
            icon: 'none'
          })
        }
      } catch (e) {
        console.error('标记完成失败:', e)
        uni.showToast({
          title: '操作失败',
          icon: 'none'
        })
      } finally {
        completing.value = false
      }
    }
  })
}

// 删除目标
const deleteTarget = () => {
  if (!target.value?.ID || deleting.value) return
  uni.showModal({
    title: '提示',
    content: '确定要删除这个目标吗？删除后无法恢复。',
    success: async (res) => {
      if (!res.confirm) return
      try {
        deleting.value = true
        const result = await api.deleteTarget({ ID: target.value.ID })
        if (result?.code === 0) {
          uni.showToast({
            title: '删除成功',
            icon: 'success'
          })
          setTimeout(() => {
            uni.navigateBack()
          }, 600)
        } else {
          uni.showToast({
            title: result?.msg || '删除失败',
            icon: 'none'
          })
        }
      } catch (e) {
        console.error('删除目标失败:', e)
        uni.showToast({
          title: '删除失败',
          icon: 'none'
        })
      } finally {
        deleting.value = false
      }
    }
  })
}

// 选择图片并上传
const chooseImage = () => {
  if (recording.value) return
  uni.chooseImage({
    count: 1,
    sizeType: ['compressed'],
    sourceType: ['album', 'camera'],
    success: async (res) => {
      const filePath = (res.tempFilePaths || [])[0]
      if (!filePath) return
      await uploadImage(filePath)
    }
  })
}

const uploadImage = async (tempPath) => {
  try {
    uni.showLoading({ title: '上传中' })
    const fullUrl = (config.baseUrl || '').replace(/\/$/, '') + '/user/upload'
    const uploadRes = await new Promise((resolve, reject) => {
      uni.uploadFile({
        url: fullUrl,
        filePath: tempPath,
        name: 'file',
        header: { Authorization: uni.getStorageSync('token') || '' },
        success: (r) => {
          try {
            resolve(JSON.parse(r.data))
          } catch (e) {
            resolve({ code: -1 })
          }
        },
        fail: reject
      })
    })
    if (uploadRes?.code === 0 && uploadRes.data?.file?.url) {
      recordForm.value.recordImg = uploadRes.data.file.url
    } else {
      uni.showToast({
        title: uploadRes?.msg || '上传失败',
        icon: 'none'
      })
    }
  } catch (e) {
    console.error('上传图片失败:', e)
    uni.showToast({
      title: '上传失败',
      icon: 'none'
    })
  } finally {
    uni.hideLoading()
  }
}

// 提交打卡记录
const submitRecord = async () => {
  if (!canSubmitRecord.value || !target.value?.ID) return
  try {
    recording.value = true
    const userId = await getUserId()
    if (!userId) {
      uni.showToast({
        title: '请先登录',
        icon: 'none'
      })
      return
    }
    const payload = {
      targetId: target.value.ID,
      userId,
      recordImg: recordForm.value.recordImg || '',
      recordNote: recordForm.value.recordNote.trim()
    }
    const res = await api.createTargetRecord(payload)
    if (res?.code === 0) {
      uni.showToast({
        title: '打卡成功',
        icon: 'success'
      })
      showRecordForm.value = false
      // 重新刷新详情与记录
      loadDetail()
    } else {
      uni.showToast({
        title: res?.msg || '打卡失败',
        icon: 'none'
      })
    }
  } catch (e) {
    console.error('打卡失败:', e)
    uni.showToast({
      title: '打卡失败',
      icon: 'none'
    })
  } finally {
    recording.value = false
  }
}

// 打卡表单控制
const openRecordForm = () => {
  recordForm.value = {
    recordNote: '',
    recordImg: ''
  }
  showRecordForm.value = true
}

const cancelRecord = () => {
  showRecordForm.value = false
}

// 点击“去打卡”
const handleGoRecord = () => {
  if (target.value?.isCompleted) return
  openRecordForm()
}

// 工具方法
const formatDate = (timeStr) => {
  if (!timeStr) return ''
  const d = new Date(timeStr)
  if (Number.isNaN(d.getTime())) return ''
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}

const formatDateTime = (timeStr) => {
  if (!timeStr) return ''
  const d = new Date(timeStr)
  if (Number.isNaN(d.getTime())) return ''
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const hh = String(d.getHours()).padStart(2, '0')
  const mm = String(d.getMinutes()).padStart(2, '0')
  return `${y}-${m}-${day} ${hh}:${mm}`
}

const previewImage = (url) => {
  if (!url) return
  uni.previewImage({
    urls: [url],
    current: url
  })
}

onMounted(() => {
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1]
  const options = currentPage.options || {}
  targetId.value = options.id || ''

  if (!targetId.value) {
    loading.value = false
    uni.showToast({
      title: '缺少目标ID',
      icon: 'none'
    })
    return
  }

  loadDetail()
})
</script>

<style lang="scss" scoped>
.detail-container {
  min-height: 100vh;
  background: #f5f5f5;
  padding: 20rpx;
  box-sizing: border-box;
}

.loading {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  color: #999;
}

.content {
  display: flex;
  flex-direction: column;
  gap: 20rpx;
}

.top-section {
  display: flex;
  flex-direction: column;
  gap: 12rpx;
}

.target-card {
  background: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
  box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.05);
}

.title-row {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  gap: 12rpx;
}

.title {
  flex: 1;
  font-size: 32rpx;
  font-weight: 600;
  color: #333;
}

.badge {
  font-size: 22rpx;
  padding: 4rpx 14rpx;
  border-radius: 999rpx;
}

.badge.ongoing {
  background: rgba(255, 107, 157, 0.1);
  color: #ff6b9d;
}

.badge.finished {
  background: rgba(153, 153, 153, 0.12);
  color: #999;
}

.desc {
  margin-top: 12rpx;
  font-size: 26rpx;
  color: #555;
  line-height: 1.6;
}

.meta-row {
  margin-top: 14rpx;
  text-align: right;
}

.meta-text {
  font-size: 22rpx;
  color: #aaa;
}

.actions-card {
  display: flex;
  flex-direction: row;
  justify-content: flex-end;
  gap: 16rpx;
}

.btn-complete {
  border-radius: 999rpx;
  padding: 0 32rpx;
}

.btn-delete {
  border-radius: 999rpx;
  padding: 0 32rpx;
}

.records-section {
  background: #fff;
  border-radius: 16rpx;
  padding: 20rpx 20rpx 16rpx;
  box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.05);
}

.section-header {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10rpx;
}

.section-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #333;
}

.btn-record {
  border-radius: 999rpx;
  padding: 0 26rpx;
  font-size: 24rpx;
}

.record-entry {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8rpx;
  margin-bottom: 10rpx;
}

.record-tip {
  font-size: 22rpx;
  color: #999;
}

.record-form-card {
  margin-top: 10rpx;
  padding: 12rpx 8rpx 4rpx;
  border-radius: 12rpx;
  background: #fff7f9;
}

.form-item {
  margin-top: 6rpx;
}

.record-list {
  margin-top: 4rpx;
}

.record-item {
  padding: 18rpx 0;
  border-bottom: 1px solid #f2f2f2;
}

.record-item:last-child {
  border-bottom-width: 0;
}

.record-header {
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 12rpx;
}

.avatar {
  width: 64rpx;
  height: 64rpx;
  border-radius: 50%;
}

.user-info {
  display: flex;
  flex-direction: column;
}

.name {
  font-size: 26rpx;
  color: #333;
}

.time {
  margin-top: 4rpx;
  font-size: 22rpx;
  color: #aaa;
}

.record-body {
  margin-top: 10rpx;
}

.note {
  font-size: 26rpx;
  color: #555;
  line-height: 1.6;
}

.record-img {
  margin-top: 10rpx;
  width: 100%;
  height: 300rpx;
  border-radius: 12rpx;
}

.empty-record {
  padding: 40rpx 0 20rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  color: #999;
}

.empty-text {
  font-size: 26rpx;
}

.empty-desc {
  margin-top: 6rpx;
  font-size: 22rpx;
}

.empty-state {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #999;
  font-size: 26rpx;
}

.popup-card {
  width: 620rpx;
  background: #fff;
  border-radius: 20rpx;
  padding: 24rpx 28rpx 20rpx;
  box-sizing: border-box;
}

.popup-title {
  font-size: 30rpx;
  font-weight: 600;
  margin-bottom: 12rpx;
  color: #333;
}

.popup-form-item {
  margin-top: 12rpx;
}

.label {
  display: block;
  font-size: 26rpx;
  color: #555;
  margin-bottom: 8rpx;
}

.img-row {
  display: flex;
  flex-direction: row;
  align-items: center;
}

.preview-img {
  width: 200rpx;
  height: 200rpx;
  border-radius: 12rpx;
}

.img-placeholder {
  width: 200rpx;
  height: 200rpx;
  border-radius: 12rpx;
  border: 1px dashed #ddd;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #999;
}

.img-tip {
  margin-top: 8rpx;
  font-size: 22rpx;
}
.form-actions {
  margin-top: 16rpx;
  display: flex;
  flex-direction: row;
  justify-content: flex-end;
  gap: 16rpx;
}

.btn-cancel {
  min-width: 140rpx;
  border-radius: 999rpx;
  background: #f5f5f5;
  color: #666;
  font-size: 26rpx;
}

.btn-confirm {
  min-width: 180rpx;
  border-radius: 999rpx;
  background: #ff6b9d;
  font-size: 26rpx;
}
</style>


