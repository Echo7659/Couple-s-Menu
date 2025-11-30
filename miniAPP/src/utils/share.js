const shareSuccess = () => {
	uni.showToast({
		title: '分享成功'
	})
}

const shareFail = () => {
	uni.showToast({
		title: '分享失败',
		icon: 'none'
	})
}

export default {
	data() {
		return {
			// 默认的分享参数
			share: {
				title: '',
				path: '/pages/home/home', // 默认分享路径
				imageUrl: '', // 默认分享图片
				desc: '',
				content: ''
			}
		}
	},
	onShareAppMessage(res) {
		console.log('全局分享', res)
		let shareInfo = res.target ? res.target.dataset.shareinfo : this.share
		return {
			...shareInfo,
			success(res) {
				shareSuccess()
			},
			fail(res) {
				shareFail()
			}
		}
	},
	// 分享到朋友圈
	onShareTimeline(res) {
		let shareInfo = res.target ? res.target.dataset.shareinfo : this.share
		return {
			...shareInfo,
			success(res) {
				shareSuccess()
			},
			fail(res) {
				shareFail()
			}
		}
	},
}