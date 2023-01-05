package controllers

import (
	"net/http"
	"ochat/comm"
	"ochat/service"
	"strconv"
)

// 好友 - 列表
func FriendList(w http.ResponseWriter, r *http.Request) {
	// verify user legal
	userInfo, code, errStr := service.NewUserServ().CheckUserRequestLegal(r)
	if errStr != "" {
		comm.ResFailure(w, code, errStr)
		return
	}

	statusStr := r.FormValue("status") // 获取指定状态的好友列表： -1:屏蔽;0:黑名单;1:好友;2:置顶'

	if statusStr == "" {
		comm.ResFailure(w, 1001, "the status params is empty")
		return
	}

	status, _ := strconv.Atoi(statusStr)

	users, err := service.NewFriendServ().List(userInfo.Id, status)
	if err != nil || len(users) == 0 {
		comm.ResFailure(w, 2001, "not exists")
		return
	}

	comm.ResSuccess(w, users)
}

func FriendAdd(w http.ResponseWriter, r *http.Request) {
	// verify user legal
	userInfo, code, errStr := service.NewUserServ().CheckUserRequestLegal(r)
	if errStr != "" {
		comm.ResFailure(w, code, errStr)
		return
	}

	r.ParseForm()
	friendIdStr := r.PostFormValue("friend_id")    // 好友id
	friendAlias := r.PostFormValue("friend_alias") // 好友别称
	about := r.PostFormValue("about")              // 描述

	friendId, _ := strconv.ParseInt(friendIdStr, 10, 64)
	friendInfo, err := service.NewFriendServ().
		Add(userInfo.Id, friendId, friendAlias, about)
	if err != nil {
		comm.ResFailure(w, 2001, "add friend failute")
		return
	}

	comm.ResSuccess(w, comm.D{
		"user_info":   userInfo,
		"friend_info": friendInfo,
	})
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
