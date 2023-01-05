package controllers

import (
	"net/http"
	"ochat/comm"
	"ochat/service"
	"strconv"
)

// 好友 - 列表
func FriendList(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.FormValue("user_id") // 用户id
	token := r.FormValue("token")       // 用户token
	statusStr := r.FormValue("status")  // 获取指定状态的好友列表： -1:屏蔽;0:黑名单;1:好友;2:置顶'

	if userIdStr == "" || token == "" || statusStr == "" {
		comm.ResFailure(w, 1001, "the user params is empty")
		return
	}

	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	userInfo, err := service.NewUserServ().UserIdToUserInfo(userId)
	if err != nil {
		comm.ResFailure(w, 1002, "user are dose not exists")
		return
	}

	// 验证token是否合法
	if userInfo.Token != token {
		comm.ResFailure(w, 1003, "token parameter validation failed")
		return
	}

	status, _ := strconv.Atoi(statusStr)

	users, err := service.NewFriendServ().List(userId, status)
	if err != nil || len(users) == 0 {
		comm.ResFailure(w, 2001, "not exists")
		return
	}

	comm.ResSuccess(w, users)
}

// 好友 - 设置黑名单
func FriendToBlacklist(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	comm.ResSuccess(w, map[string]any{
		"path":   r.URL.Path,
		"method": r.Method,
		"params": r.Form,
	})
}

// 好友 - 设置屏蔽
func FriendToHide(w http.ResponseWriter, r *http.Request) {
}
