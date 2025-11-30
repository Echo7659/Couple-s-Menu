# 省市区选择
> 省市区街道选择组件,不依赖任何UI框架,直接拷贝也可以使用(**注意JSON文件在小程序中,建议最好放到云端**)
>



![](https://cdn.nlark.com/yuque/0/2025/png/2970129/1753425272786-64072c4f-40c2-4512-bd67-9d743298707a.png)

### 使用
```vue
<city-picker
  :visible="cityPickerVisible"
  @cancel="cityPickerVisible = false">
</city-picker>
```



### 属性说明props


| **属性名** | **说明** | **类型** | **默认值** | **可选值** |
| --- | --- | --- | --- | --- |
| column | 显示的列数,如果是3的话,那就只有省市区 | Number | 4 | 1|2|3|4 |
| defaultValue | 传入中文的省市区街道数组,会自动回显 | Array | [] | - |
| visible | 是否显示组件 | Boolean | false | true |
| maskCloseAble | 点击蒙版是否关闭 | Boolean | true | false |


### 回调事件
| **事件名** | **说明** | **回调参数** |
| --- | --- | --- |
| cancel | 点击取消按钮回调事件 | <font style="color:rgb(44, 62, 80);">无</font> |
| confirm | 点击确认回调事件 | {<br/>code:[], // 选择地址的省市区<br/>name: '', // 选择地址的名称拼接好的字符串<br/>provinceName:'', // 省名称<br/>cityName:'', // 市名称<br/>areaName:'', // 区名称<br/>streetName:'', // 街道名称<br/>} |


