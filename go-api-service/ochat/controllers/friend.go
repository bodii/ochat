package controllers

import (
	"net/http"
	"ochat/comm"
	"ochat/models"
	"ochat/service"
	"strconv"
)

// 好友 - 列表
func FriendList(w http.ResponseWriter, r *http.Request) {
	// verify user legal
	user, code, errStr := service.NewUserServ().CheckUserRequestLegal(r)
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

	users, err := service.NewFriendServ().List(user.UserId, status)
	if err != nil || len(users) == 0 {
		comm.ResFailure(w, 2001, "not exists")
		return
	}

	comm.ResSuccess(w, users)
}

func FriendAdd(w http.ResponseWriter, r *http.Request) {
	// verify user legal
	user, code, errStr := service.NewUserServ().CheckUserRequestLegal(r)
	if errStr != "" {
		comm.ResFailure(w, code, errStr)
		return
	}

	r.ParseForm()
	friendIdStr := r.PostFormValue("friend_id") // 好友id
	// friendAlias := r.PostFormValue("friend_alias") // 好友别称
	// about := r.PostFormValue("about")              // 描述

	friendId, _ := strconv.ParseInt(friendIdStr, 10, 64)
	if friendId == 0 {
		comm.ResFailure(w, 1001, "friend id field failure")
		return
	}

	friend, err := service.NewUserServ().UserIdTouser(friendId)
	if err != nil {
		comm.ResFailure(w, 1201, "friend info is exists")
		return
	}

	ok, err := service.NewFriendServ().Adds(user, friend)
	if err != nil || !ok {
		comm.ResFailure(w, 2001, "add friend failute")
		return
	}

	comm.ResSuccess(w, comm.D{
		"user_info":   user,
		"friend_info": friend,
	})
}

// 好友 - 设置黑名单
func FriendToBlacklist(w http.ResponseWriter, r *http.Request) {
	// verify user legal
	user, code, errStr := service.NewUserServ().CheckUserRequestLegal(r)
	if errStr != "" {
		comm.ResFailure(w, code, errStr)
		return
	}

	// 如果好友的状态是普通
	ok, err := service.NewFriendServ().UpdateStatus(
		r, user.UserId, models.FRIEND_STATUS_BLACKLIST)
	if err != nil || !ok {
		comm.ResFailure(w, 1101, err.Error())
		return
	}

	comm.ResSuccess(w, nil)
}

// 好友 - 设置屏蔽
func FriendToHide(w http.ResponseWriter, r *http.Request) {
	// verify user legal
	user, code, errStr := service.NewUserServ().CheckUserRequestLegal(r)
	if errStr != "" {
		comm.ResFailure(w, code, errStr)
		return
	}

	// 如果好友的状态是普通
	ok, err := service.NewFriendServ().UpdateStatus(
		r, user.UserId, models.FRIEND_STATUS_HIDE)
	if err != nil || !ok {
		comm.ResFailure(w, 1101, err.Error())
		return
	}

	comm.ResSuccess(w, nil)
}

// 好友 - 设置置顶
func FriendToTop(w http.ResponseWriter, r *http.Request) {
	// verify user legal
	user, code, errStr := service.NewUserServ().CheckUserRequestLegal(r)
	if errStr != "" {
		comm.ResFailure(w, code, errStr)
		return
	}

	// 如果好友的状态是普通
	ok, err := service.NewFriendServ().UpdateStatus(
		r, user.UserId, models.FRIEND_STATUS_TOP)
	if err != nil || !ok {
		comm.ResFailure(w, 1101, err.Error())
		return
	}

	comm.ResSuccess(w, nil)
}

// 好友 - 更新信息
func FriendUpdate(w http.ResponseWriter, r *http.Request) {
	// verify user legal
	user, code, errStr := service.NewUserServ().CheckUserRequestLegal(r)
	if errStr != "" {
		comm.ResFailure(w, code, errStr)
		return
	}

	r.ParseForm()
	friendIdStr := r.PostFormValue("friend_id") // 好友id
	// friendAlias := r.PostFormValue("friend_alias") // 好友别称
	// about := r.PostFormValue("about")              // 描述

	friendId, _ := strconv.ParseInt(friendIdStr, 10, 64)
	if friendId == 0 {
		comm.ResFailure(w, 1001, "friend id field failure")
		return
	}

	friend, err := service.NewUserServ().UserIdTouser(friendId)
	if err != nil {
		comm.ResFailure(w, 1201, "friend info is exists")
		return
	}

	ok, err := service.NewFriendServ().Adds(user, friend)
	if err != nil || !ok {
		comm.ResFailure(w, 2001, "add friend failute")
		return
	}

	comm.ResSuccess(w, comm.D{
		"user_info":   user,
		"friend_info": friend,
	})
}
