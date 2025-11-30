<template>
  <view class="upload-demo-container">
    <view class="demo-section">
      <text class="demo-title">文件上传示例</text>
      
      <!-- 图片上传 -->
      <view class="upload-section">
        <text class="upload-label">图片上传</text>
        <uni-file-picker
          v-model="imageFiles"
          file-mediatype="image"
          mode="grid"
          :limit="9"
          @select="onImageSelect"
          @delete="onImageDelete"
        ></uni-file-picker>
      </view>
      
      <!-- 视频上传 -->
      <view class="upload-section">
        <text class="upload-label">视频上传</text>
        <uni-file-picker
          v-model="videoFiles"
          file-mediatype="video"
          mode="list"
          :limit="3"
          @select="onVideoSelect"
          @delete="onVideoDelete"
        ></uni-file-picker>
      </view>
      
      <!-- 文件列表 -->
      <view class="file-list">
        <text class="list-title">已选择的文件</text>
        <view v-for="(file, index) in allFiles" :key="index" class="file-item">
          <text class="file-name">{{ file.name }}</text>
          <text class="file-size">{{ formatFileSize(file.size) }}</text>
        </view>
      </view>
    </view>
  </view>
</template>
<script setup>
import { ref, computed } from 'vue'

// 响应式数据
const imageFiles = ref([])
const videoFiles = ref([])

// 计算属性
const allFiles = computed(() => {
  return [...imageFiles.value, ...videoFiles.value]
})

// 方法
const onImageSelect = (e) => {
  console.log('选择图片:', e)
}

const onImageDelete = (e) => {
  console.log('删除图片:', e)
}

const onVideoSelect = (e) => {
  console.log('选择视频:', e)
}

const onVideoDelete = (e) => {
  console.log('删除视频:', e)
}

const formatFileSize = (size) => {
  if (size < 1024) {
    return size + 'B'
  } else if (size < 1024 * 1024) {
    return (size / 1024).toFixed(2) + 'KB'
  } else {
    return (size / (1024 * 1024)).toFixed(2) + 'MB'
  }
}
</script>

<style lang="scss" scoped>
.upload-demo-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  padding: 40rpx 30rpx;
}

.demo-section {
  background: #fff;
  border-radius: 20rpx;
  padding: 40rpx;
  box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.05);
  
  .demo-title {
    font-size: 36rpx;
    font-weight: bold;
    color: #333;
    text-align: center;
    margin-bottom: 40rpx;
  }
}

.upload-section {
  margin-bottom: 40rpx;
  
  .upload-label {
    display: block;
    font-size: 28rpx;
    color: #333;
    margin-bottom: 20rpx;
    font-weight: 500;
  }
}

.file-list {
  margin-top: 40rpx;
  padding-top: 30rpx;
  border-top: 1rpx solid #f0f0f0;
  
  .list-title {
    display: block;
    font-size: 28rpx;
    color: #333;
    margin-bottom: 20rpx;
    font-weight: 500;
  }
  
  .file-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 20rpx;
    background: #f8f8f8;
    border-radius: 10rpx;
    margin-bottom: 15rpx;
    
    .file-name {
      font-size: 26rpx;
      color: #333;
      flex: 1;
      margin-right: 20rpx;
    }
    
    .file-size {
      font-size: 24rpx;
      color: #999;
    }
  }
}
</style>
