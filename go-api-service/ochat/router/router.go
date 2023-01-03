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
}
