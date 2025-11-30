// 订单状态枚举
export const orderStateEnums = [
  {
    value: -1,
    label: '已取消'
  },
  {
    value: 0,
    label: '代付款'
  },
  {
    value: 1,
    label: '代发货'
  },
  {
    value: 2,
    label: '代收货'
  },
  {
    value: 3,
    label: '交易完成'
  }
]

// 售后类型
export const afterSaleTypeEnums = [
  {
    value: 0,
    label: '我要退款(无需退货)',
    desc: '未收到货，或与商家协商之后申请退款'
  },
  {
    value: 1,
    label: '我要退货退款',
    desc: '已收到货，需要退还已收到的货品'
  }
]
