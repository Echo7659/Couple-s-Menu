<template>
  <view class="create-container">
    <view class="form-card">
      <view class="form-item">
        <text class="label">标题</text>
        <uni-easyinput v-model="form.title" placeholder="给这篇日记取个标题吧"></uni-easyinput>
      </view>

      <view class="form-item">
        <text class="label">内容</text>
        <uni-easyinput v-model="form.content" type="textarea" :maxlength="500" placeholder="记录此刻的心情..." autoHeight></uni-easyinput>
        <text class="count">{{ form.content?.length || 0 }}/500</text>
      </view>

      <view class="form-item">
        <text class="label">心情</text>
        <view class="mood-list">
          <view v-for="m in moods" :key="m.value" class="mood-item" :class="{ active: form.mood === m.value }" @click="form.mood = m.value">
            <uni-icons :type="m.icon" size="20" :color="m.color" />
            <text class="mood-text">{{ m.label }}</text>
          </view>
        </view>
      </view>

      <view class="form-item">
        <text class="label">图片</text>
        <view class="uploader">
          <view v-for="(url, idx) in images" :key="idx" class="img-item">
            <image :src="url" class="img" mode="aspectFill" @click="preview(url)"></image>
            <view class="remove" @click.stop="removeImage(idx)">
              <uni-icons type="close" size="16" color="#fff" />
            </view>
          </view>
          <view v-if="images.length < 9" class="add" @click="chooseImages">
            <uni-icons type="camera" size="28" color="#999" />
          </view>
        </view>
      </view>
    </view>

    <view class="footer">
      <button class="submit" :loading="submitting" :disabled="!canSubmit" @click="submit">{{ isEdit ? '更新日记' : '发布日记' }}</button>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import api from '@/api'
import config from '@/config'

const form = ref({ title: '', content: '', mood: 9 })
const images = ref([])
const submitting = ref(false)
const diaryId = ref('')
const isEdit = computed(() => !!diaryId.value)

const moods = [
  { value: 1, label: '开心', icon: 'heart', color: '#FF6B9D' },
  { value: 2, label: '兴奋', icon: 'star', color: '#FFD700' },
  { value: 3, label: '幸福', icon: 'heart', color: '#FF8E9E' },
  { value: 4, label: '思念', icon: 'paperplane', color: '#4A90E2' },
  { value: 5, label: '心累', icon: 'close', color: '#999' },
  { value: 6, label: '疲倦', icon: 'minus', color: '#A0A0A0' },
  { value: 7, label: '失落', icon: 'sad', color: '#87CEEB' },
  { value: 8, label: '生气', icon: 'close', color: '#FF6347' },
  { value: 9, label: '平静', icon: 'heart', color: '#666' }
]

const canSubmit = computed(() => form.value.title.trim() && form.value.content.trim() && !submitting.value)

const getUserId = async () => {
  const c1 = uni.getStorageSync('userInfo')
  if (c1?.ID) return c1.ID
  const res = await api.getUserInfo()
  return res?.data?.ID || null
}

const getUserNameAvatar = () => {
  const cached = uni.getStorageSync('userInfo') || {}
  const userName = cached.nick_name || cached.nickname || undefined
  const avatar = cached.avatar || undefined
  const result = {}
  if (userName) result.userName = userName
  if (avatar) result.avatar = avatar
  return result
}

const chooseImages = () => {
  uni.chooseImage({
    count: 9 - images.value.length,
    sizeType: ['compressed'], // 只选择压缩后的图片
    sourceType: ['album', 'camera'],
    success: async (res) => {
      const files = res.tempFilePaths || []
      for (const filePath of files) {
        // 进一步压缩图片
        await compressAndUpload(filePath)
      }
    }
  })
}

// 压缩并上传图片
const compressAndUpload = async (filePath) => {
  try {
    // 先压缩图片
    const compressedPath = await new Promise((resolve, reject) => {
      uni.compressImage({
        src: filePath,
        quality: 70, // 压缩质量 0-100，70 是较好的平衡点
        success: (res) => {
          resolve(res.tempFilePath)
        },
        fail: (err) => {
          console.warn('图片压缩失败，使用原图:', err)
          resolve(filePath) // 压缩失败则使用原图
        }
      })
    })
    
    // 上传压缩后的图片
    await uploadOne(compressedPath)
  } catch (error) {
    console.error('压缩图片失败:', error)
    // 如果压缩失败，直接上传原图
    await uploadOne(filePath)
  }
}

const uploadOne = async (tempPath) => {
  try {
    uni.showLoading({ title: '上传中' })
    const fullUrl = (config.baseUrl || '').replace(/\/$/, '') + '/user/upload'
    const uploadRes = await new Promise((resolve, reject) => {
      uni.uploadFile({
        url: fullUrl,
        filePath: tempPath,
        name: 'file',
        header: { 'Authorization': uni.getStorageSync('token') || '' },
        success: (r) => { try { resolve(JSON.parse(r.data)) } catch { resolve({ code: -1 }) } },
        fail: reject
      })
    })
    if (uploadRes?.code === 0 && uploadRes.data?.file?.url) {
      images.value.push(uploadRes.data.file.url)
    } else {
      images.value.push(tempPath)
      uni.showToast({ title: uploadRes?.msg || '上传失败', icon: 'none' })
    }
  } catch (e) {
    images.value.push(tempPath)
    uni.showToast({ title: '上传失败', icon: 'none' })
  } finally {
    uni.hideLoading()
  }
}

const removeImage = (idx) => images.value.splice(idx, 1)
const preview = (url) => uni.previewImage({ urls: images.value, current: url })

// 加载日志详情（编辑模式）
const loadDiary = async (id) => {
  try {
    uni.showLoading({ title: '加载中...' })
    const res = await api.getDiaryDetail(id)
    
    if (res?.code === 0 && res.data) {
      // 填充表单数据
      form.value = {
        title: res.data.title || '',
        content: res.data.content || '',
        mood: res.data.mood || 9
      }
      
      // 处理图片字段（可能是逗号分隔的字符串）
      const imgStr = res.data.img || ''
      images.value = imgStr ? imgStr.split(',').filter(url => url.trim()) : []
    } else {
      uni.showToast({
        title: res?.msg || '加载失败',
        icon: 'none'
      })
    }
  } catch (error) {
    console.error('加载日志详情失败:', error)
    uni.showToast({
      title: '加载失败',
      icon: 'none'
    })
  } finally {
    uni.hideLoading()
  }
}

const submit = async () => {
  if (!canSubmit.value) return
  try {
    submitting.value = true
    const userId = await getUserId()
    const payload = {
      userId,
      title: form.value.title.trim(),
      content: form.value.content.trim(),
      mood: form.value.mood,
      img: images.value.join(',')
    }
    // 可选的 userName / avatar
    Object.assign(payload, getUserNameAvatar())

    let res
    if (isEdit.value) {
      // 编辑模式：更新日志
      payload.id = diaryId.value
      res = await api.updateDiary(payload)
    } else {
      // 新建模式：添加日志
      res = await api.addDiary(payload)
    }

    if (res?.code === 0) {
      uni.showToast({ 
        title: isEdit.value ? '更新成功' : '发布成功', 
        icon: 'success' 
      })
      setTimeout(() => {
        uni.navigateBack()
      }, 600)
    } else {
      uni.showToast({ title: res?.msg || (isEdit.value ? '更新失败' : '发布失败'), icon: 'none' })
    }
  } catch (e) {
    uni.showToast({ title: isEdit.value ? '更新失败' : '发布失败', icon: 'none' })
  } finally {
    submitting.value = false
  }
}

// 生命周期
onMounted(() => {
  // 获取路由参数
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1]
  const options = currentPage.options || {}
  diaryId.value = options.id || ''
  
  // 如果是编辑模式，加载日志详情
  if (diaryId.value) {
    loadDiary(diaryId.value)
  }
})
</script>

<style lang="scss" scoped>
.create-container { min-height: 100vh; background: #f5f5f5; padding-bottom: 120rpx; }
.form-card { background:#fff; margin: 20rpx; padding: 24rpx; border-radius: 16rpx; box-shadow: 0 4rpx 16rpx rgba(0,0,0,.05); }
.form-item { margin-bottom: 20rpx; position: relative; }
.label { display:block; font-size: 28rpx; color:#333; margin-bottom: 12rpx; font-weight: 600; }
.count { position:absolute; right: 0; bottom: -24rpx; font-size: 22rpx; color:#999; }

.mood-list { display:flex; flex-wrap: wrap; gap: 12rpx; }
.mood-item { display:flex; align-items:center; gap: 8rpx; padding: 10rpx 14rpx; border-radius: 20rpx; background:#f8f8f8; color:#666; }
.mood-item.active { background: #FFE3EE; color:#FF6B9D; }
.mood-text { font-size: 24rpx; }

.uploader { display:flex; flex-wrap: wrap; gap: 16rpx; }
.img-item { width: 160rpx; height: 160rpx; position: relative; }
.img { width: 100%; height: 100%; border-radius: 10rpx; object-fit: cover; }
.remove { position:absolute; right: 6rpx; top: 6rpx; width: 36rpx; height: 36rpx; border-radius: 50%; background: rgba(0,0,0,.35); display:flex; align-items:center; justify-content:center; }
.add { width: 160rpx; height: 160rpx; border: 2rpx dashed #ddd; border-radius: 10rpx; display:flex; align-items:center; justify-content:center; }

.footer { position: fixed; left: 0; right: 0; bottom: 0; background:#fff; padding: 16rpx 20rpx env(safe-area-inset-bottom); box-shadow: 0 -4rpx 16rpx rgba(0,0,0,.06); }
.submit { width: 100%; background: linear-gradient(135deg, #FF6B9D 0%, #FF8E9E 100%); color:#fff; border:none; border-radius: 30rpx; padding: 16rpx; font-size: 28rpx; }
</style> 