<template>
  <view class="t-cropper" v-show="props.imageUrl" disable-scroll @touchmove.stop.prevent>
    <!-- #ifdef MP-WEIXIN || MP-ALIPAY-->
    <canvas type="2d" canvas-id="canvas-cropper" id="canvas-cropper" class="canvas"></canvas>
    <!-- #endif -->

    <!-- #ifdef APP-PLUS || H5 -->
    <canvas
      :canvas-id="canvasId"
      :style="{ width: canvasWidth + 'px', height: canvasHeight + 'px' }"
    ></canvas>
    <!-- #endif -->
    <!-- 裁剪区域 -->
    <view class="t-preview-container" :class="{ showPage: inInit }">
      <view
        class="preview-body"
        @touchstart="e => touchStart(e, 'body')"
        @touchmove="touchMove"
        @touchend="touchEnd"
        @touchcancel="touchCancel"
      >
        <!-- 完整背景图片 -->
        <view class="image-wrap" :class="{ transit: transit }" :style="comImageStyle('image-wrap')">
          <image
            class="image"
            :class="{ transit: transit }"
            :src="imageUrl"
            @load="imageLoad"
            :style="comImageStyle('image')"
          />
        </view>
        <view class="mask-model"></view>
        <view class="frame-box" :class="{ transit: transit }" :style="comImageStyle('frame-box')">
          <view class="rect" :style="cirStyle">
            <!-- 裁剪框图片 -->
            <view
              class="image-rect"
              :class="{ transit: transit }"
              :style="comImageStyle('image-rect')"
            >
              <image
                class="rect-img"
                :class="{ transit: transit }"
                :src="imageUrl"
                :style="comImageStyle('image')"
              />
            </view>
          </view>
          <!-- 矩阵框线条 -->
          <view v-if="props.radius < 50" class="line-box">
            <view class="line-one"></view>
            <view class="line-two"></view>
            <view class="line-three"></view>
            <view class="line-four"></view>
          </view>
          <view class="frame-left-top" @touchstart.stop="e => touchStart(e, 'left-top')"></view>
          <view
            class="frame-left-bottom"
            @touchstart.stop="e => touchStart(e, 'left-bottom')"
          ></view>
          <view class="frame-right-top" @touchstart.stop="e => touchStart(e, 'right-top')"></view>
          <view
            class="frame-right-bottom"
            @touchstart.stop="e => touchStart(e, 'right-bottom')"
          ></view>
        </view>
      </view>
      <!-- 底部工具栏 -->
      <view class="toolbar">
        <view @tap.stop="onCancle" class="btn-cancel">取消</view>
        <view v-if="props.isRotateBtn" @tap.stop="onAngle" class="btn-rotate">
          <image
            src="https://cfzy-durian-front.tos-cn-guangzhou.volces.com/zm-mall-mini-app/rotate.svg"
          />
        </view>
        <view @tap.stop="onConfirm" class="btn-confirm">确定</view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { reactive, ref, watchEffect, computed, getCurrentInstance, toRefs } from 'vue'

/**
 * @property {String} mode 模式
 * @value fixed 固定模式，裁剪出固定大小
 * @value ratio 等比模式，宽高等比缩放
 * @value free 自由模式，不限制宽高比
 * @property {String} imageUrl 图片路径
 * @property {Number} width 宽度
 * @property {Number} height 高度
 * @property {Number} maxWidth 最大宽度
 * @property {Number} minHeight 最大高度
 * @property {Number} scaleRatio 裁剪比列缩放
 * @property {Number} minRatio 最小缩放
 * @property {Number} maxRatio 最大缩放
 * @property {Boolean} isRotateBtn 是否显示旋转
 * @property {Number} radius 裁剪图片圆角半径，单位px
 */
const props = defineProps({
  mode: {
    type: String,
    default: 'ratio'
  },
  imageUrl: {
    type: String,
    default: ''
  },
  width: {
    type: Number,
    default: 200
  },
  height: {
    type: Number,
    default: 200
  },
  maxWidth: {
    type: Number,
    default: 1024
  },
  maxHeight: {
    type: Number,
    default: 1024
  },
  scaleRatio: {
    type: Number,
    default: 0.7
  },
  minRatio: {
    type: Number,
    default: 1
  },
  maxRatio: {
    type: Number,
    default: 3
  },
  isRotateBtn: {
    type: Boolean,
    default: true
  },
  radius: {
    type: Number,
    default: 0
  }
})

// 获取组件实例
const instance = getCurrentInstance()
const inInit = ref(false)
const $emit = defineEmits(['confirm', 'cancel'])
const canvasId = ref(Math.random().toString(36).slice(-6))

//获取原始图片的大小
const origin = reactive({
  imageUrl: '', //设置图片url
  real: {
    //原始图片宽高
    width: 100,
    height: 100
  },
  body: {
    //页面宽高
    width: 100,
    height: 100
  },
  frame: {
    //矩形框
    left: 50,
    top: 50,
    width: 200,
    height: 300
  },
  bgImage: {
    //背景框
    left: 20,
    top: 20,
    width: 300,
    height: 400
  },
  rotate: 0,
  transit: false,
  canvasWidth: 100,
  canvasHeight: 100
})

const { transit, rotate, imageUrl, canvasWidth, canvasHeight } = toRefs(origin)

//监听原始图片的变动
watchEffect(() => {
  if (props.imageUrl) {
    uni.showLoading({
      title: '请稍候...',
      mask: true
    })
    origin.imageUrl = props.imageUrl //设置图片url}
  }
})

//检测图片加载完成
const imageLoad = e => {
  const { width, height } = e.detail
  origin.real.width = width
  origin.real.height = height
  let query = uni.createSelectorQuery().in(instance)
  query
    .select('.preview-body')
    .boundingClientRect(data => {
      origin.body.width = data.width
      origin.body.height = data.height
      inInit.value = true
      imageReset() //重置图片
    })
    .exec()
}

//重置图片
const imageReset = () => {
  origin.rotate = 0
  let frameRate = props.width / props.height //裁剪比列
  let frameHeight = origin.body.height * props.scaleRatio //裁剪框图片高度 * 缩小0.7倍
  let frameWidth = origin.body.width * props.scaleRatio //裁剪框图片宽度 * 缩小0.7倍

  // 裁剪框的位置：
  // 缩放后的宽度/高度 > 组件裁剪的比例，就要对宽度重写
  if (frameWidth / frameHeight > frameRate) {
    frameWidth = frameHeight * frameRate
  } else {
    frameHeight = frameWidth / frameRate
  }

  // 裁剪框左边距=页面宽度-裁剪框宽度 / 2
  let frameleft = (origin.body.width - frameWidth) / 2
  let frameTop = (origin.body.height - frameHeight) / 2
  origin.frame = {
    left: frameleft,
    top: frameTop,
    width: frameWidth,
    height: frameHeight
  }

  // 背景图片位置：
  let bgRate = origin.real.width / origin.real.height //背景比列
  let bgWidth = frameWidth
  let bgHeight = frameHeight

  // 裁剪框图片宽度/裁剪框图片高度>实际的图片比例
  if (bgWidth / bgHeight > bgRate) {
    bgHeight = bgWidth / bgRate
  } else {
    bgWidth = bgHeight * bgRate
  }

  // 背景左边距
  let bgLeft = (frameWidth - bgWidth) / 2 + origin.frame.left
  // 背景右边距
  let bgTop = (frameHeight - bgHeight) / 2 + origin.frame.top

  origin.bgImage = {
    left: bgLeft,
    top: bgTop,
    width: bgWidth,
    height: bgHeight
  }
  uni.hideLoading()
}

//公共图片背景位置
const comImageStyle = computed(() => source => {
  const { left, top, width, height } = origin.bgImage
  if (source == 'image-wrap') {
    //计算背景盒子样式
    return `
      top: ${top}px;
      left: ${left}px;
      width: ${width}px;
      height: ${height}px;
    `
  } else if (source == 'image') {
    // 计算背景图片样式
    let left = 0
    let top = 0
    let width = origin.bgImage.width
    let height = origin.bgImage.height

    if (origin.rotate % 180 != 0) {
      width = origin.bgImage.height
      height = origin.bgImage.width
      top = width / 2 - height / 2
      left = height / 2 - width / 2
    }
    const style = {}
    style.left = left + 'px'
    style.top = top + 'px'
    style.width = width + 'px'
    style.height = height + 'px'
    style.transform = `rotate(${origin.rotate}deg)`
    return style
  } else if (source == 'image-rect') {
    // 计算裁剪框上显示的图片位置
    return `
      top: ${top - origin.frame.top}px;
      left: ${left - origin.frame.left}px;
      width: ${width}px;
      height: ${height}px;
    `
  } else if (source == 'frame-box') {
    return `
      left: ${origin.frame.left}px;
      top: ${origin.frame.top}px;
      width: ${origin.frame.width}px;
      height: ${origin.frame.height}px;
    `
  }
})

//计算圆角比列
const cirStyle = computed(() => {
  const { width } = origin.frame
  let scale = props.width / width
  let radius = props.radius / scale
  return `
  border-radius: ${radius}px
  `
})

// 取消按钮
const onCancle = () => {
  $emit('cancel')
}
// 旋转按钮
const onAngle = () => {
  origin.rotate -= 90
  let width = origin.bgImage.height
  let height = origin.bgImage.width
  let left = origin.bgImage.left
  let top = origin.bgImage.top
  let rate = width / height
  if (width < origin.frame.width) {
    width = origin.frame.width
    height = width / rate
  }
  if (height < origin.frame.height) {
    height = origin.frame.height
    width = height * rate
  }
  if (left > origin.frame.left) {
    left = origin.frame.left
  }
  if (top > origin.frame.top) {
    top = origin.frame.top
  }
  if (left + width < origin.frame.left + origin.frame.width) {
    left = origin.frame.left + origin.frame.width - width
  }
  if (top + height < origin.frame.top + origin.frame.height) {
    top = origin.frame.top + origin.frame.height - height
  }
  origin.bgImage = {
    left: left,
    top: top,
    width: width,
    height: height
  }
  origin.transit = true
  setTimeout(() => {
    origin.transit = false
  }, 300)
}
// 确定按钮
const onConfirm = async () => {
  // #ifdef H5 || APP-PLUS
  confirmH5()
  // #endif

  // #ifdef MP-WEIXIN || MP-ALIPAY
  confirmWx()
  // #endif
}

// h5裁剪确认
const confirmH5 = async () => {
  let mx = computeMatrix() //获取画布的信息
  // 设备画布的宽高
  origin.canvasWidth = mx.tw
  origin.canvasHeight = mx.th
  uni.showLoading({
    title: '处理中'
  })

  // 设置宽高后等待更新完成在获取Canvas节点
  await new Promise(resolve => {
    setTimeout(() => {
      resolve()
    }, 200)
  })

  // 创建一个画布
  let ctx = uni.createCanvasContext(canvasId.value, instance)
  // 矩形内清除指定的像素
  await ctx.clearRect(0, 0, mx.tw, mx.th)
  //剪切圆角
  await drawClipImage(ctx, mx)
  // 旋转
  ctx.rotate((origin.rotate * Math.PI) / 180)
  //绘制图片
  ctx.drawImage(props.imageUrl, mx.sx, mx.sy, mx.sw, mx.sh, mx.dx, mx.dy, mx.dw, mx.dh)
  ctx.restore()

  ctx.draw(false, () => {
    uni.canvasToTempFilePath(
      {
        canvasId: canvasId.value,
        success: res => {
          let path = res.tempFilePath
          // #ifdef H5
          path = parseBlob(path)
          $emit('confirm', {
            errMsg: 'parseBlob:ok',
            tempFilePath: path
          })
          // #endif

          // #ifdef APP-PLUS
          $emit('confirm', res)
          // #endif
        },
        fail: err => {
          console.log(err)
        },
        complete: () => {
          uni.hideLoading() //关闭loading
        }
      },
      instance
    )
  })
}
// 微信确认
const confirmWx = () => {
  uni.showLoading({
    title: '处理中'
  })
  let mx = computeMatrix() //获取画布的信息
  uni
    .createSelectorQuery()
    // #ifndef MP-ALIPAY
    .in(instance)
    // #endif
    .select('#canvas-cropper')
    .fields({
      node: true,
      size: true
    })
    .exec(async res => {
      const textCanvas = res[0].node
      const ctx = textCanvas.getContext('2d')

      // 初始化画布大小
      textCanvas.width = mx.tw
      textCanvas.height = mx.th
      // //绘制canvas
      await drawClipImage(ctx, mx) //剪切圆角
      await ctx.rotate((origin.rotate * Math.PI) / 180)
      await createImage(origin.imageUrl, mx, ctx, textCanvas)

      // #ifdef MP-WEIXIN
      wx.canvasToTempFilePath({
        canvas: textCanvas,
        canvasId: 'canvasID',
        success: res => {
          $emit('confirm', res)
        },
        fail: err => {
          console.log(err)
        }
      })
      // #endif

      // #ifdef MP-ALIPAY
      textCanvas.toTempFilePath({
        success: res => {
          $emit('confirm', res)
        },
        fail: err => {
          console.log(err)
        }
      })
      // #endif

      uni.hideLoading() //关闭loading
    })
}
//剪切图片圆角
const drawClipImage = (ctx, mx) => {
  let { radius } = props
  if (radius > 0) {
    const w = Math.round(mx.tw) //画布宽度
    const h = Math.round(mx.th) //画布高度

    // 被缩放后的尺寸且裁剪的宽高相同
    if (props.width === props.height && w == h && props.width < w) {
      radius = (w / props.width) * radius
    }
    // 被缩放后的尺寸且裁剪的宽高不相同
    if (props.width != props.height && w != h && props.width < w) {
      console.log('11111111')
      if (w > h) {
        radius = (w / h) * radius
      } else {
        radius = (h / w) * radius
      }
    }
    // 如果裁剪是圆形
    if (w === h && radius >= w / 2) {
      // 圆形
      ctx.arc(w / 2, h / 2, w / 2, 0, 2 * Math.PI)
    } else {
      // 圆角矩形
      if (w !== h) {
        // 限制圆角半径不能超过短边的一半
        radius = Math.min(w / 2, h / 2, radius)
      }
      // 其他形状
      ctx.moveTo(radius, 0)
      ctx.arcTo(w, 0, w, h, radius)
      ctx.arcTo(w, h, 0, h, radius)
      ctx.arcTo(0, h, 0, 0, radius)
      ctx.arcTo(0, 0, w, 0, radius)
      ctx.closePath()
    }
    ctx.clip() //从原始画布剪切任意形状和尺寸的区域
    ctx.restore() //返回之前保存过的路径状态和属性
  }
}
// 创建一个图片
const createImage = (img, mx, ctx, textCanvas) => {
  return new Promise(resolve => {
    const headerImg = textCanvas.createImage()
    headerImg.src = img
    headerImg.onload = () => {
      let result = ctx.drawImage(headerImg, mx.sx, mx.sy, mx.sw, mx.sh, mx.dx, mx.dy, mx.dw, mx.dh)
      resolve(result)
    }
  })
}

//计算矩阵
const computeMatrix = () => {
  let width = props.width
  let height = props.height
  let mul = origin.bgImage.width / origin.real.width
  if (origin.rotate % 180 != 0) {
    mul = origin.bgImage.height / origin.real.width
  }
  if (props.mode != 'fixed') {
    width = origin.frame.width / mul
    height = origin.frame.height / mul
  }
  let rate = width / height
  if (width > props.maxWidth) {
    width = props.maxWidth
    height = width / rate
  }
  if (height > props.maxHeight) {
    height = props.maxHeight
    width = height * rate
  }
  let sx = (origin.frame.left - origin.bgImage.left) / mul
  let sy = (origin.frame.top - origin.bgImage.top) / mul
  let sw = origin.frame.width / mul
  let sh = origin.frame.height / mul
  let ox = sx + sw / 2
  let oy = sy + sh / 2
  if (origin.rotate % 180 != 0) {
    let temp = sw
    sw = sh
    sh = temp
  }
  let angle = origin.rotate % 360
  if (angle < 0) {
    angle += 360
  }
  if (angle == 270) {
    var x = origin.real.width - oy
    var y = ox
    ox = x
    oy = y
  }
  if (angle == 180) {
    var x = origin.real.width - ox
    var y = origin.real.height - oy
    ox = x
    oy = y
  }
  if (angle == 90) {
    var x = oy
    var y = origin.real.height - ox
    ox = x
    oy = y
  }
  sx = ox - sw / 2
  sy = oy - sh / 2
  let dr = { x: 0, y: 0, w: width, h: height }
  dr = parseRect(dr, -origin.rotate)
  return {
    tw: width,
    th: height,
    sx: sx,
    sy: sy,
    sw: sw,
    sh: sh,
    dx: dr.x,
    dy: dr.y,
    dw: dr.w,
    dh: dr.h
  }
}
// 计算矩阵
const parseRect = (rect, angle) => {
  let x1 = rect.x
  let y1 = rect.y
  let x2 = rect.x + rect.w
  let y2 = rect.y + rect.h
  let p1 = parsePoint({ x: x1, y: y1 }, angle)
  let p2 = parsePoint({ x: x2, y: y2 }, angle)
  let result = {}
  result.x = Math.min(p1.x, p2.x)
  result.y = Math.min(p1.y, p2.y)
  result.w = Math.abs(p2.x - p1.x)
  result.h = Math.abs(p2.y - p1.y)
  return result
}
// 计算x、y
const parsePoint = (point, angle) => {
  let result = {}
  result.x =
    point.x * Math.cos((angle * Math.PI) / 180) - point.y * Math.sin((angle * Math.PI) / 180)
  result.y =
    point.y * Math.cos((angle * Math.PI) / 180) + point.x * Math.sin((angle * Math.PI) / 180)
  return result
}
// base64转blob
const parseBlob = base64 => {
  let arr = base64.split(',')
  let mime = arr[0].match(/:(.*?);/)[1]
  let bstr = atob(arr[1])
  let n = bstr.length
  let u8arr = new Uint8Array(n)
  for (let i = 0; i < n; i++) {
    u8arr[i] = bstr.charCodeAt(i)
  }
  let url = URL || webkitURL
  return url.createObjectURL(new Blob([u8arr], { type: mime }))
}
/**
 * 绘制手指触摸
 */
const move = reactive({
  touchType: 'body', //移动的类型
  start: {
    //开始的位置
    frame: {
      left: 0,
      top: 0,
      width: 0,
      height: 0
    },
    bgImage: {
      left: 0,
      top: 0,
      width: 0,
      height: 0
    }
  }
})

const touch = reactive({
  touchStart: [] //记录开始触摸
})

//  手指触摸动作开始
const touchStart = (e, touchType = '') => {
  // #ifdef APP-PLUS || H5
  if (e.preventDefault) {
    e.preventDefault()
  }
  if (e.stopPropagation) {
    e.stopPropagation()
  }
  // #endif

  switch (touchType) {
    case 'body':
      move.touchType = 'body'
      break
    case 'left-top':
      move.touchType = 'left-top'
      break
    case 'left-bottom':
      move.touchType = 'left-bottom'
      break
    case 'right-top':
      move.touchType = 'right-top'
      break
    case 'right-bottom':
      move.touchType = 'right-bottom'
      break
    default:
      move.touchType = 'body'
      break
  }

  touch.touchStart = e.touches
  // 裁剪框的开始位置
  move.start.frame.left = origin.frame.left
  move.start.frame.top = origin.frame.top
  move.start.frame.width = origin.frame.width
  move.start.frame.height = origin.frame.height

  // 背景图片的开始位置
  move.start.bgImage.left = origin.bgImage.left
  move.start.bgImage.top = origin.bgImage.top
  move.start.bgImage.width = origin.bgImage.width
  move.start.bgImage.height = origin.bgImage.height
  return false
}

// 手指触摸后移动
const touchMove = e => {
  // #ifdef APP-PLUS || H5
  if (e.preventDefault) {
    e.preventDefault()
  }
  if (e.stopPropagation) {
    e.stopPropagation()
  }
  // #endif
  // 单手操作
  if (touch.touchStart.length == 1) {
    if (move.touchType == 'body') {
      // 移动图片
      moveImage(touch.touchStart[0], e.touches[0])
    } else {
      if (props.mode != 'fixed') {
        // 放大缩小图片
        scaleFrame(touch.touchStart[0], e.touches[0])
      }
    }
  } else if (touch.touchStart.length == 2 && e.touches.length == 2) {
    let ta = touch.touchStart[0]
    let tb = touch.touchStart[1]
    let tc = e.touches[0]
    let td = e.touches[1]
    if (ta.identifier != tc.identifier) {
      let temp = tc
      tc = td
      td = temp
    }
    scaleImage(ta, tb, tc, td)
  }
}

// 移动图片
const moveImage = (ta, tb) => {
  let mx = tb.clientX - ta.clientX
  let my = tb.clientY - ta.clientY
  // 开始移动图片背景 + 移动的左边距离
  origin.bgImage.left = move.start.bgImage.left + mx
  origin.bgImage.top = move.start.bgImage.top + my

  let frameLeft = origin.frame.left
  let frameTop = origin.frame.top
  let frameWidth = origin.frame.width
  let frameHeight = origin.frame.height

  // 左边触边
  if (origin.bgImage.left > frameLeft) {
    origin.bgImage.left = frameLeft
  }
  // 头部触边
  if (origin.bgImage.top > frameTop) {
    origin.bgImage.top = frameTop
  }
  // 右边触边
  if (origin.bgImage.left + origin.bgImage.width < frameLeft + frameWidth) {
    origin.bgImage.left = frameLeft + frameWidth - origin.bgImage.width
  }
  // 底部触边：背景距离头部+背景的高度 < 背景高度 + 矩形框顶部
  if (origin.bgImage.top + origin.bgImage.height < frameTop + frameHeight) {
    origin.bgImage.top = frameTop + frameHeight - origin.bgImage.height
  }
}

//缩放裁剪框
const scaleFrame = (ta, tb) => {
  let mx = tb.clientX - ta.clientX
  let my = tb.clientY - ta.clientY
  let x1 = move.start.frame.left
  let y1 = move.start.frame.top
  let x2 = move.start.frame.left + move.start.frame.width
  let y2 = move.start.frame.top + move.start.frame.height
  let cx1 = false
  let cy1 = false
  let cx2 = false
  let cy2 = false
  var mix = 30
  var mix = 30
  let rate = origin.frame.width / origin.frame.height

  //左上角
  if (move.touchType == 'left-top') {
    x1 += mx
    y1 += my
    cx1 = true
    cy1 = true
    // 左下角
  } else if (move.touchType == 'left-bottom') {
    x1 += mx
    y2 += my
    cx1 = true
    cy2 = true
    // 右上角
  } else if (move.touchType == 'right-top') {
    x2 += mx
    y1 += my
    cx2 = true
    cy1 = true
    // 右下角
  } else if (move.touchType == 'right-bottom') {
    x2 += mx
    y2 += my
    cx2 = true
    cy2 = true
  }

  if (x1 < origin.bgImage.left) {
    x1 = origin.bgImage.left
  }
  if (y1 < origin.bgImage.top) {
    y1 = origin.bgImage.top
  }
  if (x2 > origin.bgImage.left + origin.bgImage.width) {
    x2 = origin.bgImage.left + origin.bgImage.width
  }
  if (y2 > origin.bgImage.top + origin.bgImage.height) {
    y2 = origin.bgImage.top + origin.bgImage.height
  }
  if (cx1) {
    if (x1 > x2 - mix) {
      x1 = x2 - mix
    }
  }
  if (cy1) {
    if (y1 > y2 - mix) {
      y1 = y2 - mix
    }
  }
  if (cx2) {
    if (x2 < x1 + mix) {
      x2 = x1 + mix
    }
  }
  if (cy2) {
    if (y2 < y1 + mix) {
      y2 = y1 + mix
    }
  }
  if (cx1) {
    if (props.mode != 'free') {
      var val = x2 - rate * (y2 - y1)
      if (x1 < val) {
        x1 = val
      }
    }
  }
  if (cy1) {
    if (props.mode != 'free') {
      var val = y2 - (x2 - x1) / rate
      if (y1 < val) {
        y1 = val
      }
    }
  }
  if (cx2) {
    if (props.mode != 'free') {
      var val = rate * (y2 - y1) + x1
      if (x2 > val) {
        x2 = val
      }
    }
  }
  if (cy2) {
    if (props.mode != 'free') {
      var val = (x2 - x1) / rate + y1
      if (y2 > val) {
        y2 = val
      }
    }
  }

  origin.frame.left = x1
  origin.frame.top = y1
  origin.frame.width = x2 - x1
  origin.frame.height = y2 - y1
}

//双指缩放图片
const scaleImage = (ta, tb, tc, td) => {
  let x1 = ta.clientX
  let y1 = ta.clientY
  let x2 = tb.clientX
  let y2 = tb.clientY
  let x3 = tc.clientX
  let y3 = tc.clientY
  let x4 = td.clientX
  let y4 = td.clientY
  let ol = Math.sqrt((x1 - x2) * (x1 - x2) + (y1 - y2) * (y1 - y2))
  let el = Math.sqrt((x3 - x4) * (x3 - x4) + (y3 - y4) * (y3 - y4))
  let ocx = (x1 + x2) / 2
  let ocy = (y1 + y2) / 2
  let ecx = (x3 + x4) / 2
  let ecy = (y3 + y4) / 2
  let ax = ecx - ocx
  let ay = ecy - ocy
  let scale = el / ol

  // 检测放大/缩小倍率
  if (move.start.bgImage.width * scale < origin.frame.width) {
    scale = origin.frame.width / move.start.bgImage.width
  }

  if (move.start.bgImage.height * scale < origin.frame.height) {
    scale = origin.frame.height / move.start.bgImage.height
  }

  // 最大缩放
  if (move.start.bgImage.width * scale > origin.frame.width * props.maxRatio) {
    return
  }

  origin.bgImage.left = move.start.bgImage.left + ax - (ocx - move.start.bgImage.left) * (scale - 1)

  origin.bgImage.top = move.start.bgImage.top + ay - (ocy - move.start.bgImage.top) * (scale - 1)

  origin.bgImage.width = move.start.bgImage.width * scale
  origin.bgImage.height = move.start.bgImage.height * scale

  if (origin.bgImage.left > origin.frame.left) {
    origin.bgImage.left = origin.frame.left
  }
  if (origin.bgImage.top > origin.frame.top) {
    origin.bgImage.top = origin.frame.top
  }
  if (origin.bgImage.left + origin.bgImage.width < origin.frame.left + origin.frame.width) {
    origin.bgImage.left = origin.frame.left + origin.frame.width - origin.bgImage.width
  }
  if (origin.bgImage.top + origin.bgImage.height < origin.frame.top + origin.frame.height) {
    origin.bgImage.top = origin.frame.top + origin.frame.height - origin.bgImage.height
  }
}

// 手指触摸动作结束
const touchEnd = () => {
  touch.touchStart = []
  resizeImage() //调整图片的大小
}

// 手指触摸动作被打断，如来电提醒，弹窗
const touchCancel = () => {
  touch.touchStart = []
}

// 调整图片的大小
const resizeImage = () => {
  let rate = origin.frame.width / origin.frame.height //比列
  let width = origin.body.width * props.scaleRatio //裁剪框图片高度 * 缩小0.7倍
  let height = origin.body.height * props.scaleRatio //裁剪框图片宽度 * 缩小0.7倍
  // 图片的位置：
  // 缩放后的宽度/高度 > 组件裁剪的比例，就要对宽度重写
  if (width / height > rate) {
    width = height * rate
  } else {
    height = width / rate
  }
  let left = (origin.body.width - width) / 2
  let top = (origin.body.height - height) / 2
  let mul = width / origin.frame.width
  let ox = origin.frame.left - origin.bgImage.left
  let oy = origin.frame.top - origin.bgImage.top

  origin.frame = {
    left: left,
    top: top,
    width: width,
    height: height
  }

  width = origin.bgImage.width * mul
  height = origin.bgImage.height * mul
  left = origin.frame.left - ox * mul
  top = origin.frame.top - oy * mul
  origin.bgImage = {
    left: left,
    top: top,
    width: width,
    height: height
  }
  if (mul != 1) {
    origin.transit = true
    setTimeout(() => {
      origin.transit = false
    }, 300)
  }
}

/**
 * 绘制手指触摸结束
 */
</script>

<style>
page {
  overflow: hidden;
  overscroll-behavior: none;
}
</style>
<style lang="scss" scoped>
@import './t-cropper';
</style>
