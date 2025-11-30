import Request from 'luch-request'
const http = new Request()

export const useUpload = async file => {
  function getRandom(num) {
    let random = Math.floor(
      (Math.random() + Math.floor(Math.random() * 9 + 1)) * Math.pow(10, num - 1)
    )

    return random
  }

  const fileUri = `${new Date().getTime()}${getRandom(10)}${file.name || ''}`

  const uploadFileFn = async arrayBuffer => {
    uni.showLoading({
      title: '上传中'
    })

    uni.hideLoading()

    return arrayBuffer

    // await http
    //   .request({
    //     url: uploadUrl.url,
    //     method: 'PUT',
    //     data: arrayBuffer // ArrayBuffer数据
    //   })
    //   .catch(err => {
    //     console.error('上传过程中发生错误', err)
    //     uni.showToast({
    //       title: '上传异常',
    //       icon: 'none'
    //     })
    //     uni.hideLoading()
    //   })
    //
    // uni.hideLoading()
    //
    // return fileConfig.fileUrl + fileUri
  }

  return new Promise((resolve, reject) => {
    // #ifdef MP-WEIXIN
    uni.showLoading({
      title: '读取中'
    })

    uni.getFileSystemManager().readFile({
      filePath: file.path,
      success: async res => {
        console.log('success=====', res)
        const url = await uploadFileFn(res.data)
        resolve(url)
      },
      fail: e => {
        console.log('fail----', e)
      }
    })
    // #endif

    // #ifdef H5
    // const reader = new FileReader()
    const url = URL.createObjectURL(file)
    resolve(url)
    console.log('file---', file)
    // reader.readAsArrayBuffer(file)
    // reader.onload = async () => {
    //   const arrayBuffer = reader.result
    //   const url = await uploadFileFn(arrayBuffer)
    //   resolve(url)
    // }
    // #endif
  })
}
