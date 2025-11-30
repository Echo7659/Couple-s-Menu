<template>
  <!-- 上传start -->
  <view style="display: flex; flex-wrap: wrap">
    <view class="update-file">
      <!--图片-->
      <view v-for="(item, index) in fileList" :key="index" class="mr-[12rpx]">
        <view class="upload-box">
          <image
            class="preview-file"
            v-if="item.type == 0"
            :src="item.url"
            @tap="previewImage(item.url)"
            mode="aspectFill"
          />
          <video v-else-if="item.type == 1" class="preview-file" :src="item.videoUrl"></video>
          <view class="remove-icon" @tap="deleteHandle(index)" v-if="!onlyShow">
            <i class="iconfont icon-a-Linearguanbi text-[20rpx]"></i>
          </view>
        </view>
      </view>
      <!--按钮-->
      <view v-if="VideoOfImagesShow && !onlyShow" @tap="chooseVideoImage" class="upload-btn">
        <view class="add-box">
          <i class="iconfont icon-a-Linearxiangji text-[50rpx]"></i>
          <view class="color-text-disable text-[24rpx]">
            <view class="text-center">上传</view>
            <view class="text-center">最多{{ maxLength }}张</view>
          </view>
        </view>
      </view>
    </view>
  </view>
  <!-- 上传 end -->
</template>
<script setup>
import { useUpload } from './upload'
// 设置初始数据和响应式状态
const sourceType = ref(['拍摄', '相册', '拍摄或相册'])
const sourceTypeIndex = ref(2)
const cameraList = reactive([
  { value: 'back', name: '后置摄像头', checked: 'true' },
  { value: 'front', name: '前置摄像头' }
])
const cameraIndex = ref(0)
const props = defineProps({
  maxLength: {
    type: Number,
    default: 1
  },
  videoMaxSize: {
    // 单位为M
    type: Number,
    default: 100
  },
  imgMaxSize: {
    // 单位为M
    type: Number,
    default: 20
  },
  modelValue: {
    type: Array,
    default: () => []
  },
  onlyShow: {
    // 是否仅查看
    type: Boolean,
    default: false
  },
  onlyImg: {
    // 是否只上传图片
    type: Boolean,
    default: false
  }
})

const VideoOfImagesShow = computed(() => {
  return fileList.value.length < props.maxLength
})

const emits = defineEmits(['update:modelValue'])

const fileList = computed({
  get() {
    props.modelValue.map(t => {
      if (t.type == 1 && /http/.test(t.url) && !t.videoUrl) {
        //   如果是视频流,ios需要先下载下来,才能播放,不能直接播放二进制的视频流
        uni.downloadFile({
          url: t.url,
          success: res => {
            t.videoUrl = res.tempFilePath
          }
        })
      }
    })

    return props.modelValue
  },
  set(newValue) {
    emits('update:modelValue', newValue)
  }
})

// 方法
function chooseVideoImage() {
  if (props.onlyImg) {
    chooseImages()
  } else {
    uni.showActionSheet({
      title: '选择上传类型',
      itemList: ['图片', '视频'],
      success: res => {
        if (res.tapIndex == 0) {
          chooseImages()
        } else {
          chooseVideo()
        }
      }
    })
  }
}

function chooseImages() {
  uni.chooseImage({
    count: props.maxLength - fileList.value.length, // 默认9张
    sizeType: ['original', 'compressed'],
    sourceType: ['album', 'camera'],
    success: res => {
      const validRes = res.tempFiles.every(item => {
        if (item.size > props.imgMaxSize * 1024 * 1000) {
          uni.showToast({
            title: `图片大小不能超过${props.imgMaxSize}M`,
            icon: 'error'
          })
          return false
        }
        return true
      })

      if (!validRes) return

      res.tempFiles.map(async t => {
        const url = await useUpload(t)
        console.log('url----', url)
        fileList.value.push({
          type: 0,
          url: url
        })
      })
    }
  })
}

function chooseVideo() {
  uni.chooseVideo({
    maxDuration: 60,
    count: 9,
    camera: cameraList[cameraIndex.value].value,
    // sourceType: sourceType.value[sourceTypeIndex.value],
    success: async res => {
      let videoFile = res.tempFile || res

      // #ifdef MP-WEIXIN
      videoFile.path = videoFile.tempFilePath
      // #endif

      if (videoFile.size > props.videoMaxSize * 1024 * 1000) {
        uni.showToast({
          title: `视频大小不能超过${props.videoMaxSize}M`,
          icon: 'error'
        })
        return
      }
      const url = await useUpload(videoFile)

      console.log('videoFile---', videoFile.path)

      fileList.value.push({
        type: 1,
        url: url,
        videoUrl: videoFile.path
      })
    },
    fail: error => {
      uni.hideLoading()
      uni.showToast({
        title: error,
        icon: 'none'
      })
    }
  })
}

function previewImage(item) {
  console.log('预览图片', item)
  uni.previewImage({
    current: item,
    urls: fileList.value.filter(t => t.type == 0).map(t2 => t2.url)
  })
}

function deleteHandle(index) {
  uni.showModal({
    title: '提示',
    content: '是否要删除?',
    success: res => {
      if (res.confirm) {
        fileList.value.splice(index, 1)
      }
    }
  })
}
</script>

<style scoped lang="scss">
// 上传
.update-file {
  margin-left: 10rpx;
  height: auto;
  display: flex;
  //justify-content: space-between;
  flex-wrap: wrap;
  margin-bottom: 5rpx;

  .del-icon {
    width: 44rpx;
    height: 44rpx;
    position: absolute;
    right: 10rpx;
    top: 12rpx;
  }

  .btn-text {
    color: #606266;
  }

  .preview-file {
    width: 160rpx;
    height: 160rpx;
    border: 1px solid #e0e0e0;
    border-radius: 10rpx;
    object-fit: contain;
  }

  .upload-box {
    position: relative;
    width: 160rpx;
    height: 160rpx;
    margin: 0 10rpx 20rpx 0;
  }

  .remove-icon {
    position: absolute;
    right: 10rpx;
    top: 10rpx;
    z-index: 100;
    width: 30rpx;
    height: 30rpx;
    background: var(--color-function-error);
    border-radius: 50%;
    text-align: center;
    line-height: 30rpx;
    color: #fff;
    font-size: 16rpx;
    //padding: 5rpx;
  }

  .upload-btn {
    width: 160rpx;
    height: 160rpx;
    border-radius: 10rpx;
    background-color: #f4f5f6;
    display: flex;
    justify-content: center;
    align-items: center;
  }
}
.guide-view {
  margin-top: 30rpx;
  display: flex;

  .guide-text {
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    padding-left: 20rpx;

    .guide-text-price {
      padding-bottom: 10rpx;
      color: #ff0000;
      font-weight: bold;
    }
  }
}
.service-process {
  background-color: #ffffff;
  padding: 20rpx;
  padding-top: 30rpx;
  margin-top: 30rpx;
  border-radius: 10rpx;
  padding-bottom: 30rpx;

  .title {
    text-align: center;
    margin-bottom: 20rpx;
  }
}
.form-view-parent {
  border-radius: 20rpx;
  background-color: #ffffff;
  padding: 0rpx 20rpx;

  .form-view {
    background-color: #ffffff;
    margin-top: 20rpx;
  }

  .form-view-textarea {
    margin-top: 20rpx;
    padding: 20rpx 0rpx;

    .upload-hint {
      margin-top: 10rpx;
      margin-bottom: 10rpx;
    }
  }
}

.bottom-class {
  margin-bottom: 60rpx;
}

.bottom-btn-class {
  padding-bottom: 1%;

  .bottom-hint {
    display: flex;
    justify-content: center;
    padding-bottom: 20rpx;
  }
}
.add-box {
  width: 160rpx;
  height: 160rpx;
  background: #f7f8f9;
  border: 1rpx dashed #c0c4cc;
  border-radius: 8rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}
</style>
