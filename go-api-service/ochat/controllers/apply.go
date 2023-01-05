package controllers

import (
	"net/http"
	"ochat/comm"
	"ochat/comm/funcs"
	"ochat/models"
	"ochat/service"
	"strconv"
	"strings"
)

// 通过用户名或手机号查找好友
func ApplyFind(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	usernameOrMobile := r.PostFormValue("username|mobile")
	usernameOrMobile = strings.TrimSpace(usernameOrMobile)
	if usernameOrMobile == "" {
		comm.ResFailure(w, 1001, "username or moblie is empty")
		return
	}

	var user models.User
	var condition string
	if funcs.IsMobile(usernameOrMobile) {
		// 通过手机号查找
		condition = "mobile = ?"
	} else {
		// 通过用户名查找
		condition = "username = ?"
	}
	// 查找
	service.NewUserServ().DB.Where(condition, usernameOrMobile).Get(&user)

	if user.Id == 0 {
		comm.ResFailure(
			w,
			1002,
			"the user info accociated with the username or mobile does not exist")
		return
	}

	// 返回用户信息
	comm.ResSuccess(w, user)
}

// 发起申请添加好友
func ApplyAdd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userIdStr := r.PostFormValue("userid")        // 当前发起申请的用户
	addUserIdStr := r.PostFormValue("add_userid") // 被申请的用户
	comment := r.PostFormValue("comment")         // 留言

	if userIdStr == "" || funcs.IsNumber(userIdStr) {
		comm.ResFailure(w, 1001, "the user parameters to be added are incorrect")
		return
	}

	if addUserIdStr == "" || funcs.IsNumber(addUserIdStr) {
		comm.ResFailure(w, 1002, "the user parameters to be added are incorrect")
		return
	}

	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	addUserId, _ := strconv.ParseInt(addUserIdStr, 10, 64)
	// 查找
	_, err := service.NewUserServ().UserIdToUserInfo(userId)
	if err != nil {
		comm.ResFailure(w, 1003, "user does not exists")
		return
	}
	_, err = service.NewUserServ().UserIdToUserInfo(addUserId)
	if err != nil {
		comm.ResFailure(w, 1004, "user does not exists")
		return
	}

	// 添加用户申请
	addData, err := service.NewApplyServ().
		Add(userId, addUserId, comment, models.APPLY_TYPE_USER)
	if err != nil {
		comm.ResFailure(w, 1005, "add failure")
		return
	}

	comm.ResSuccess(w, addData)
}

// 查看向我申请好友的列表信息
func ApplyList(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userIdStr := r.PostFormValue("userid") // 当前用户id

	if userIdStr == "" || !funcs.IsNumber(userIdStr) {
		comm.ResFailure(w, 1001, "the user parameters to be added are incorrect")
		return
	}

	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	// status应答者是否同意,-1:拒绝;0:未查看;1:已查看;2:同意
	users, err := service.NewApplyServ().List(userId, 0, 1)
	if err != nil || len(users) == 0 {
		comm.ResFailure(w, 1002, "the query is incorrect or does not exist")
		return
	}

	comm.ResSuccess(w, users)
}

// 处理向我申请好友
func ApplyDispose(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	idStr := r.PostFormValue("id")                // 当前用户id
	disposeIdStr := r.PostFormValue("dispose_id") // 操作id 对应models.apply.APPLY_STATUS_REFUSE...

	if idStr == "" || !funcs.IsNumber(idStr) {
		comm.ResFailure(w, 1001, "the apply id parameters are incorrect")
		return
	}

	id, _ := strconv.ParseInt(idStr, 10, 64)

	if disposeIdStr == "" {
		comm.ResFailure(w, 1002, "the dispose parameters dose not exists")
		return
	}

	disposeId, err := strconv.Atoi(disposeIdStr)
	if err != nil || disposeId < -1 || disposeId > 2 {
		comm.ResFailure(w, 1003, "the dispose parameters are incorrect")
		return
	}

	// status应答者是否同意,-1:拒绝;0:未查看;1:已查看;2:同意
	ok, err := service.NewApplyServ().Set(id, disposeId)
	if !ok || err != nil {
		comm.ResFailure(w, 1004, "dispose are incorrect")
		return
	}

	comm.ResSuccess(w, nil)
}
