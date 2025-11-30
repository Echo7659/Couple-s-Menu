import { ref } from 'vue'

/**
 *
 * @param minutes 分钟数
 * @returns {{seconds: 剩余时间,单位为秒>, timeStr: 时间转换成hh:mm格式的字符串}}
 */
export const useCountDown = minutes => {
  let seconds = ref(minutes * 60)
  let timeStr = ref('')
  const intervalId = setInterval(function () {
    const minutesRemaining = Math.floor(seconds.value / 60)
    const secondsRemaining = seconds.value % 60

    // 将分钟和秒数格式化为字符串，确保单个数字前面有零
    const formattedMinutes = String(minutesRemaining).padStart(2, '0')
    const formattedSeconds = String(secondsRemaining).padStart(2, '0')

    timeStr.value = formattedMinutes + ':' + formattedSeconds

    if (seconds.value <= 0) {
      clearInterval(intervalId)
    } else {
      seconds.value--
    }
  }, 1000) // 每秒更新一次

  return {
    seconds,
    timeStr
  }
}

// 验证码倒计时
export const useSmsCodeCountDown = () => {
  const smsCodeBtnText = ref('获取验证码')
  const countdown = ref(60)
  let timer = null
  const isCountDown = ref(false) // 是否倒计时中

  const startCountDownInterval = () => {
    if (timer) return
    smsCodeBtnText.value = `${countdown.value}秒后重发`
    isCountDown.value = true
    timer = setInterval(() => {
      countdown.value--
      if (countdown.value === 0) {
        isCountDown.value = false
        smsCodeBtnText.value = '获取验证码'
        countdown.value = 60
        clearInterval(timer)
        timer = null
      } else {
        smsCodeBtnText.value = `${countdown.value}秒后重发`
      }
    }, 1000)
  }
  return {
    smsCodeBtnText,
    isCountDown,
    startCountDownInterval
  }
}
