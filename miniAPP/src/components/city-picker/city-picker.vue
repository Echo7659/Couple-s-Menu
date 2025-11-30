<template>
  <view class="pupop">
    <view class="popup-box" :animation="animationData">
      <view class="pupop-btn">
        <view @tap="cancel">取消</view>
        <view @tap="confirm" style="color: var(--color-function-success)">确定</view>
      </view>
      <picker-view
        v-if="addressList.length"
        :value="value"
        :indicator-style="indicatorStyle"
        @change="onChange"
        class="picker-view"
      >
        <picker-view-column>
          <view class="item" v-for="(item, index) in provinceList" :key="index">
            {{ item.name }}
          </view>
        </picker-view-column>
        <picker-view-column>
          <view class="item" v-for="(item, index) in cityList" :key="index">{{ item.name }}</view>
        </picker-view-column>
        <picker-view-column>
          <view class="item" v-for="(item, index) in areaList" :key="index">{{ item.name }}</view>
        </picker-view-column>
        <picker-view-column v-if="column === 4">
          <view class="item" v-for="(item, index) in streetList" :key="index">{{ item.name }}</view>
        </picker-view-column>
      </picker-view>
    </view>

    <view
      @tap="close"
      @touchmove.stop.prevent
      :class="visible ? 'pupop-model' : 'pupop-models'"
    ></view>
  </view>
</template>

<script setup>
import { ref, watch, nextTick, computed } from 'vue'
import addressData from './pcas-code.json'
// import { addressList as addressData } from './cityData.js'

let addressList = reactive([])
const getAddressList = async () => {
  // const { data } = await addressApi.getPcasJson()
  // 这里建议放到项目中时,把地图的json数据放到云端,避免小程序主包过大
  addressList = addressData
  init()
}

onMounted(() => {
  getAddressList()
})

const props = defineProps({
  column: {
    // 显示的列数，如果是3的话，就只有省市区，没有街道，以此类推
    type: Number,
    default: 4
  },
  defaultValue: {
    // 传入中文的省市区街道数组,则会自动回显
    type: [String, Array],
    default: ''
  },
  visible: {
    type: Boolean,
    default: false
  },
  maskCloseAble: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['confirm', 'cancel'])

const value = ref([])
const animationData = ref({})
const indicatorStyle = 'height: 50px;'

const provinceList = ref([])
const cityList = ref([])
const areaList = ref([])
const streetList = ref([])

const provinceIndex = ref(0)
const cityIndex = ref(0)
const areaIndex = ref(0)
const streetIndex = ref(0)

watch(
  () => props.visible,
  () => {
    triggerAnimation()
  }
)

watch(
  () => props.defaultValue,
  () => {
    init()
  }
)

function init() {
  if (!addressList.length) return

  provinceList.value = addressList.map(item => ({
    name: item.name,
    code: item.code
  }))

  console.log('provinceList.value -=---', provinceList.value)
  if (!props.defaultValue) {
    cityList.value = addressList[0].children
    areaList.value = cityList.value[0].children
    streetList.value = areaList.value[0]?.children || []
  } else if (Array.isArray(props.defaultValue)) {
    const [pName, cName, aName, sName] = props.defaultValue

    provinceIndex.value = addressList.findIndex(p => p.name === pName)
    if (provinceIndex.value < 0) provinceIndex.value = 0
    const province = addressList[provinceIndex.value]
    cityList.value = province.children || []
    cityIndex.value = cityList.value.findIndex(c => c.name === cName)
    if (cityIndex.value < 0) cityIndex.value = 0
    const city = cityList.value[cityIndex.value]

    areaList.value = city.children || []
    areaIndex.value = areaList.value.findIndex(a => a.name === aName)
    if (areaIndex.value < 0) areaIndex.value = 0
    const area = areaList.value[areaIndex.value]

    streetList.value = area.children || []
    if (props.column === 4) {
      streetIndex.value = streetList.value.findIndex(s => s.name === sName)
      if (streetIndex.value < 0) streetIndex.value = 0
    }

    nextTick(() => {
      value.value = [provinceIndex.value, cityIndex.value, areaIndex.value]
      if (props.column === 4) value.value.push(streetIndex.value)
    })
  }
  triggerAnimation()
}

function triggerAnimation() {
  const animation = uni.createAnimation({
    duration: 400,
    timingFunction: 'linear'
  })
  animation.bottom(props.visible ? 0 : '-350px').step()
  animationData.value = animation.export()
}

function onChange(e) {
  if (!addressList.length) return

  const v = e.detail.value

  const [p = 0, c = 0, a = 0, s = 0] = v

  if (p !== provinceIndex.value) {
    provinceIndex.value = p
    cityIndex.value = 0
    areaIndex.value = 0
    streetIndex.value = 0
    cityList.value = addressList[p].children
    areaList.value = cityList.value[0]?.children || []
    streetList.value = areaList.value[0]?.children || []
  } else if (c !== cityIndex.value) {
    cityIndex.value = c
    areaIndex.value = 0
    streetIndex.value = 0

    areaList.value = cityList.value[c]?.children || []
    streetList.value = areaList.value[0]?.children || []
  } else if (a !== areaIndex.value) {
    areaIndex.value = a
    streetIndex.value = 0

    streetList.value = areaList.value[a]?.children || []
  } else if (s !== streetIndex.value && props.column === 4) {
    streetIndex.value = s
  }

  value.value = [provinceIndex.value, cityIndex.value, areaIndex.value]
  if (props.column === 4) value.value.push(streetIndex.value)
}

function confirm() {
  const province = addressList[provinceIndex.value]
  const city = cityList.value[cityIndex.value]
  const area = areaList.value[areaIndex.value]
  const street = streetList.value[streetIndex.value]

  const data = {
    code: [province?.code, city?.code, area?.code],
    name: province.name + city.name + area.name,
    provinceName: province.name,
    cityName: city.name,
    areaName: area.name
  }

  if (props.column === 4 && street) {
    data.code = [...data.code, street.code]
    data.name += street.name
    data.streetName = street.name
  }
  emit('confirm', data)
}

function cancel() {
  emit('cancel')
}

function close() {
  if (props.maskCloseAble) {
    cancel()
  }
}
</script>

<style scoped lang="scss">
.pupop {
  .popup-box {
    position: fixed;
    left: 0;
    bottom: -315px;
    z-index: 99999;
    background: #fff;
    padding-bottom: 50px;
    .pupop-btn {
      height: 40px;
      display: flex;
      align-items: center;
      padding: 0 30upx;
      justify-content: space-between;
    }
  }
  .pupop-model {
    position: fixed;
    left: 0;
    top: 0;
    width: 100%;
    height: 100%;
    z-index: 9999;
    background: rgba(0, 0, 0, 0.5);
  }
  .pupop-models {
    display: none;
  }
  .picker-view {
    width: 750rpx;
    height: 225px;
    margin-top: 20rpx;
  }
  .item {
    height: 50px;
    align-items: center;
    justify-content: center;
    text-align: center;
    line-height: 50px;
  }
}
</style>
