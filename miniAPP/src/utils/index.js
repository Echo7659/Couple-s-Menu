/**
 * 树形数组扁平
 * @param {*} tree
 * @returns
 */
export function treeToArray(tree) {
  return tree.reduce((res, item) => {
    const { children, ...i } = item
    return res.concat(i, children && children.length ? treeToArray(children) : [])
  }, [])
}

/**
 * 反向递归获取树状数据列表,递归获取当前数据所有的上级,并且按照层级返回一个扁平化的数组
 * @param {Array} treeArray 树状数据
 * @param {Object} targetValue 目标值
 * @param {String}  key 判断的key值
 * @param {String}  childrenKey 子集的key值
 * @returns {Array} 递归出来的扁平化数组,里面的元素是对象
 * */
export function getPathArrByTree(treeArray, targetValue, key, childrenKey = 'children') {
  function backwardRecursion(arr, target, currentPath = []) {
    for (const node of arr) {
      const newPath = [...currentPath, node]
      // 如果找到目标值，返回当前路径
      if (node[key] === target[key]) {
        return newPath
      }

      // 如果当前节点有子节点，递归在子节点中查找
      if (node[childrenKey] && node[childrenKey].length > 0) {
        const result = backwardRecursion(node[childrenKey], target, newPath)
        if (result) {
          return result
        }
      }
    }
    // 如果在当前分支未找到目标值，返回null
    return null
  }
  return backwardRecursion(treeArray, targetValue)
}

/**
 * 正向递归获取树状数据列表,递归获取当前数据所有的下级,并且按照层级返回一个扁平化的数组
 * @param {Array} treeArray 树状数据
 * @param {String}  childrenKey 子集的key值
 * @returns {Array} 递归出来的扁平化数组,里面的元素是对象
 * */
export function recursion(treeArray, childrenKey = 'children') {
  function downRecursion(node) {
    let objects = []
    if (typeof node === 'object') {
      objects.push(node)
      if (Array.isArray(node[childrenKey])) {
        node[childrenKey].forEach(child => {
          objects = objects.concat(downRecursion(child))
        })
      }
    }

    return objects
  }

  return treeArray.reduce((acc, curr) => acc.concat(downRecursion(curr)), [])
}

/**
 * 是否图片
 * @param {*} str
 * @returns
 */
export const isImage = str => {
  if (str && typeof str === 'string') {
    // 定义图片文件扩展名的正则表达式
    const imageExtensions = /\.(jpeg|jpg|gif|png|bmp|webp)$/i
    return imageExtensions.test(str)
  }
  return false
}

/**
 * @description 文件转base64工具函数
 * @param {Object} file 文件对象
 * @returns {Promise}
 * */
export function fileToBase64(file) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => resolve(reader.result.split(',')[1])
    reader.onerror = error => reject(error)
    reader.readAsDataURL(file)
  })
}

/**
 * @description 根据文件对象获取后缀名
 * @param {Object} file 文件对象
 * @returns {String}
 * */
export function getFileExtension(file) {
  const fileName = file.name
  const lastDotIndex = fileName.lastIndexOf('.')
  const fileExtension = fileName.substring(lastDotIndex + 1)
  return fileExtension
}

// 生成随机唯一值
export function generateRandomValue() {
  // 生成一个随机的浮点数，然后将其转换为字符串
  const randomValue = 'D' + Math.random().toString().substring(2, 5)

  // 获取当前时间戳，并转换为字符串
  const timestamp = new Date().getTime().toString()

  // 将随机值和时间戳拼接在一起
  const uniqueRandomValue = randomValue + timestamp
  return uniqueRandomValue
}

/**
 * 将对象参数转变成url的参数
 * @param obj
 * @returns {string}
 */
export function objectToUrlParams(obj) {
  return Object.keys(obj)
    .map(key => `${key}=${obj[key]}`)
    .join('&')
}

/**
 * 电话号码脱敏
 * @param phone
 * @returns {String}
 */
export function phoneDesensitization(phone) {
  if (!phone) return

  // 定义手机号正则表达式
  let reg = /^(1[3-9][0-9])\d{4}(\d{4}$)/
  // 判断手机号是否能够通过正则校验
  let isMobile = reg.test(phone)
  console.log(isMobile)
  // 将手机号中间4位用*号进行显示
  let hiddenMobile = phone.replace(reg, '$1****$2')
  return hiddenMobile
}

/**
 *
 * @param money {number} 后台使用的是分,前端展示一般都是元,所以/100
 * @returns {number}
 * @constructor
 */
export const moneyConversion = money => {
  return money || money === 0 ? (money / 100).toFixed(2) : 0
}

/**
 * @param money {number} 带小数点的金额
 * @return {Object} {integer:整数,decimals:小数点}
 */
export const moneySplit = money => {
  const arr = money.toString().split('.')
  return {
    integer: arr[0],
    decimals: arr[1]
  }
}

export const moneyConversionAndSplit = money => {
  return moneySplit(moneyConversion(money))
}

/**
 *
 * @param money {number} 后台使用的是分,前端表单一般都是元,所以 * 100
 * @returns {number}
 * @constructor
 */
export const submitMoneyConversion = money => {
  return +(money * 100).toFixed(2)
}

/**
 * 防抖函数
 * @param fn
 * @param wait
 * @returns {(function(): void)|*}
 */
export const debounce = (func, wait) => {
  // @TODO：实现逻辑
  let timeout
  return function () {
    if (timeout) {
      clearTimeout(timeout)
    }
    timeout = setTimeout(() => {
      //...arguments用过获取参数
      func.call(this, ...arguments)
    }, wait)
  }
}

export const contactService = () => {
  uni.makePhoneCall({
    phoneNumber: '4008882292' //仅为示例
  })
}

export function rpxToPx(rpx) {
  const systemInfo = uni.getSystemInfoSync()
  const screenWidth = systemInfo.screenWidth // 获取屏幕宽度
  return (rpx / 750) * screenWidth
}
