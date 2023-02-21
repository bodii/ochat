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

// 群 - 查看群信息
func Group(w http.ResponseWriter, r *http.Request) {
	// verify user legal
	_, code, errStr := service.NewUserServ().CheckUserRequestLegal(r)
	if errStr != "" {
		comm.ResFailure(w, code, errStr)
		return
	}

	r.ParseForm()
	groupIdStr := r.PostFormValue("group_id")
	if groupIdStr == "" || !funcs.IsNumber(groupIdStr) {
		comm.ResFailure(w, 1001, "group_id param failure")
		return
	}
	groupId, _ := strconv.ParseInt(groupIdStr, 10, 64)

	group, err := service.NewGroupServ().Info(groupId)
	if err != nil {
		comm.ResFailure(w, 2001, err.Error())
		return
	}

	comm.ResSuccess(w, group)
}

// 群 - 创建
func GroupCreate(w http.ResponseWriter, r *http.Request) {
	// verify user legal
	user, code, errStr := service.NewUserServ().CheckUserRequestLegal(r)
	if errStr != "" {
		comm.ResFailure(w, code, errStr)
		return
	}

	r.ParseForm()
	usersIdStr := r.PostFormValue("users_id")

	users := make([]models.User, 0)
	if usersIdStr != "" {
		usersIds := strings.Split(usersIdStr, ",")
		for _, u := range usersIds {
			if !funcs.IsNumber(u) {
				continue
			}

			userid, _ := strconv.ParseInt(u, 10, 64)
			user, err := service.NewUserServ().UserIdTouser(userid)
			if err != nil {
				continue
			}

			users = append(users, user)
		}
	}

	ok, groupInfo, contacts, err := service.NewGroupServ().Create(user, users...)
	if err != nil || !ok {
		comm.ResFailure(w, 2101, err.Error())
		return
	}

	comm.ResSuccess(w, map[string]any{
		"group_info": groupInfo,
		"contacts":   contacts,
	})
}

// 群 - 查看用户的所有群信息
func GroupList(w http.ResponseWriter, r *http.Request) {
	// verify user legal
	user, code, errStr := service.NewUserServ().CheckUserRequestLegal(r)
	if errStr != "" {
		comm.ResFailure(w, code, errStr)
		return
	}

	groups, err := service.NewGroupServ().UserList(user.UserId)
	if err != nil {
		comm.ResFailure(w, 2101, err.Error())
		return
	}

	comm.ResSuccess(w, groups)
}

// 群 - 修改群信息
func GroupUpFiled(w http.ResponseWriter, r *http.Request) {
	user, code, errStr := service.NewUserServ().CheckUserRequestLegal(r)
	if errStr != "" {
		comm.ResFailure(w, code, errStr)
		return
	}

	r.ParseForm()
	groupIdStr := r.PostFormValue("group_id")
	if groupIdStr == "" || !funcs.IsNumber(groupIdStr) {
		comm.ResFailure(w, 1001, "group_id param failure")
		return
	}
	groupId, _ := strconv.ParseInt(groupIdStr, 10, 64)

	group, err := service.NewGroupServ().UpdateFields(r.PostForm, user.UserId, groupId)
	if err != nil {
		comm.ResFailure(w, 2001, err.Error())
		return
	}

	comm.ResSuccess(w, group)
}

// 群 - 二维码
func GroupQrCode(w http.ResponseWriter, r *http.Request) {
	// verify user legal
	_, code, errStr := service.NewUserServ().CheckUserRequestLegal(r)
	if errStr != "" {
		comm.ResFailure(w, code, errStr)
		return
	}

	r.ParseForm()
	groupIdStr := r.PostFormValue("group_id")
	if groupIdStr == "" || !funcs.IsNumber(groupIdStr) {
		comm.ResFailure(w, 1001, "group_id param failure")
		return
	}
	groupId, _ := strconv.ParseInt(groupIdStr, 10, 64)

	group, err := service.NewGroupServ().Info(groupId)
	if err != nil {
		comm.ResFailure(w, 2001, err.Error())
		return
	}

	if group.QrCode == "" {
		_, err = service.NewGroupServ().CreateQrCode(&group)
		if err != nil {
			comm.ResFailure(w, 2002, "create QR Code failure")
			return
		}
	}

	comm.ResSuccess(w, group)
}
