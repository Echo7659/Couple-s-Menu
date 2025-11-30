> 支持H5与小程序选择上传图片与视频,并且支持数量限制,图片,视频大小限制,自带图片视频的预览,也兼容了IOS在小程序中的预览
>



![](https://cdn.nlark.com/yuque/0/2025/png/2970129/1753429032538-27fd247a-c2de-4e4a-8bab-bec07e4e2acd.png)



### 使用
```vue
      <cu-upload-img-video v-model="fileList"></cu-upload-img-video>
```



### 属性props
| **属性名** | **说明** | **类型** | **默认值** | **可选值** |
| --- | --- | --- | --- | --- |
| maxLength | 最多上传多少个文件 | Number | 1 |  |
| videoMaxSize | 视频文件的最大大小,单位为M | Number | 100 |  |
| imgMaxSize | 图片文件的最大大小,单位为M | Number | 20 |  |
| v-model | 文件的字符串数组 | Array | string[] |  |
| onlyShow | 是否仅查看,为true时,不展示上传按钮 | Boolean | false | true |
| onlyImg | 是否只上传图片 | Boolean | false | true |


#### 注意事项,需要自己实现upload.js中的上传逻辑,因为每个人的上传逻辑不太一样,就不做封装了
```vue
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

```

