// isdev 为 true 表示开发环境 false 表示发布环境
const isdev = false;

// 开发环境使用代理，生产环境使用完整URL
const baseUrl = isdev ? 'https://463bf33e.r10.cpolar.top' : 'https://menu.hnlc5588.cn';// 开发环境代理 & 生产环境

const shareUrl = isdev ? 'https://h5.gwkjxb.com/' : 'http://test_h5.gwkjxb.com/';

const config = {
	appName: '小呆呆的私人菜谱',
	baseUrl,
	appVersion: '1.0.0',
	developer: '小呆呆的私人菜谱',
	shareUrl,
	appID:'wx8ed262fbd9eaaf74',
	isdev
}
export default config
