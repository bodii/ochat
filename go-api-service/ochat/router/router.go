package router

import (
	"net/http"
	"ochat/controllers"
)

func Init() {
	// 用户登录
	http.HandleFunc("/user/login", controllers.UserLogin)
	// 用户注册
	http.HandleFunc("/user/signup", controllers.UserRegister)
	// 登录二维码
	http.HandleFunc("/user/login/qrcode", controllers.LoginQRCode)
	// 扫描登录二维码
	http.HandleFunc("/user/login/scan_qrcode", controllers.LoginQRCodeScan)
	// 获取手机号的验证码
	http.HandleFunc("/user/sms", controllers.PhoneSms)
	// 验证手机验证码是否正确
	http.HandleFunc("/user/sms/verify", controllers.PhoneSmsVerify)

	// 头像显示
	http.HandleFunc("/avatar/show", controllers.AvatarShow)
	// 头像上传
	http.HandleFunc("/avatar/upload", controllers.AvatarUpload)

	// 申请好友 - 查找
	http.HandleFunc("/apply/find", controllers.ApplyFind)
	// 申请好友 - 添加
	http.HandleFunc("/apply/add", controllers.ApplyAdd)
	// 申请好友 - 查看
	http.HandleFunc("/apply/list", controllers.ApplyList)
	// 申请好友 - 操作
	http.HandleFunc("/apply/dispose", controllers.ApplyDispose)

	// 好友 - 设置黑名单
	http.HandleFunc("/friend/blacklist", controllers.FriendToBlacklist)
}
