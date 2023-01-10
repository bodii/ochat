package controllers

import (
	"net/http"
	"ochat/comm"
	"ochat/comm/funcs"
	"ochat/models"
	"ochat/service"
	"strconv"
	"time"
)

// 群联系人 - 列表
func GroupContactList(w http.ResponseWriter, r *http.Request) {
	// verify user legal
	user, code, errStr := service.NewUserServ().CheckUserRequestLegal(r)
	if errStr != "" {
		comm.ResFailure(w, code, errStr)
		return
	}

	comm.ResSuccess(w, user)
}

// 群联系人 - 查看
func GroupContact(w http.ResponseWriter, r *http.Request) {
	// verify user legal
	user, code, errStr := service.NewUserServ().CheckUserRequestLegal(r)
	if errStr != "" {
		comm.ResFailure(w, code, errStr)
		return
	}

	groupIdStr := r.FormValue("group_id")
	userIdStr := r.FormValue("user_id")
	groupId, _ := strconv.ParseInt(groupIdStr, 10, 64)
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	groupContact, err := service.NewGroupContactServ().Info(userId, groupId)
	if err != nil {
		comm.ResFailure(w, 2101, "get failure")
	}

	comm.ResSuccess(w, map[string]any{
		"group_contact_info": groupContact,
		"is_self":            userId == user.Id,
	})
}

// 群联系人 - 群主/管理员踢人
func GroupContactKickOut(w http.ResponseWriter, r *http.Request) {
	// verify user legal
	user, code, errStr := service.NewUserServ().CheckUserRequestLegal(r)
	if errStr != "" {
		comm.ResFailure(w, code, errStr)
		return
	}

	r.ParseForm()
	groupIdStr := r.PostFormValue("group_id")
	kUserIdStr := r.PostFormValue("kick_out_user_id")
	if groupIdStr == "" || funcs.IsNumber(groupIdStr) ||
		kUserIdStr == "" || funcs.IsNumber(kUserIdStr) {

		comm.ResFailure(w, 1001, "params failure")
		return
	}
	groupId, _ := strconv.ParseInt(groupIdStr, 10, 64)
	kUserId, _ := strconv.ParseInt(kUserIdStr, 10, 64)

	groupContact, err := service.NewGroupContactServ().Info(user.Id, groupId)
	if err != nil {
		comm.ResFailure(w, 2001, "change failure")
		return
	}

	kickGroupContact, err := service.NewGroupContactServ().Info(kUserId, groupId)
	if err != nil {
		comm.ResFailure(w, 2002, "change failure")
		return
	}

	if kickGroupContact.Type >= groupContact.Type {
		comm.ResFailure(w, 2003, "operation without permission")
		return
	}

	kickGroupContact, err = service.NewGroupContactServ().ChangeStatus(
		kickGroupContact.Id, models.GROUP_CONTACT_STATUS_KICK_OUT)
	if err != nil {
		comm.ResFailure(w, 2103, "change failure")
		return
	}

	comm.ResSuccess(w, kickGroupContact)
}

// 群联系人 - 退出
func GroupContactExit(w http.ResponseWriter, r *http.Request) {
	// verify user legal
	user, code, errStr := service.NewUserServ().CheckUserRequestLegal(r)
	if errStr != "" {
		comm.ResFailure(w, code, errStr)
		return
	}

	r.ParseForm()
	groupIdStr := r.PostFormValue("group_id")
	if groupIdStr == "" || funcs.IsNumber(groupIdStr) {
		comm.ResFailure(w, 1001, "params failure")
		return
	}
	groupId, _ := strconv.ParseInt(groupIdStr, 10, 64)

	groupContact, err := service.NewGroupContactServ().Info(user.Id, groupId)
	if err != nil {
		comm.ResFailure(w, 2001, "change failure")
		return
	}

	// TODO: 如果当前用户是群主，则退出后群主由第一个管理员接手
	if groupContact.Type == models.GROUP_CONTACT_TYPE_MASTER {
		gc, err := service.NewGroupContactServ().TypeInfo(groupId, models.GROUP_CONTACT_TYPE_MANAGER)
		if err != nil {
			comm.ResFailure(w, 2101, "query failure")
			return
		}
		if gc.Id == 0 {
			gc, err = service.NewGroupContactServ().TypeInfo(groupId, models.GROUP_CONTACT_TYPE_MEMBER)
			if err != nil {
				comm.ResFailure(w, 2101, "query failure")
				return
			}
		}

		gc.Type = models.GROUP_CONTACT_TYPE_MASTER
		service.NewGroupContactServ().DB.Where("id = ?", gc.Id).Cols("type", "updated_at").Update(&gc)
		service.NewGroupServ().DB.Where("id = ?", gc.GroupId).Cols("manager_id", "updated_at").
			Update(map[string]any{
				"manager_id": gc.UserId,
				"updated_at": time.Now(),
			})

	}
	// TODO：如果当前群联系人只剩一个人时，退出后，群解散

	groupContact, err = service.NewGroupContactServ().ChangeStatus(
		groupContact.Id, models.GROUP_CONTACT_STATUS_EXIT)
	if err != nil {
		comm.ResFailure(w, 2103, "change failure")
		return
	}

	comm.ResSuccess(w, groupContact)
}

// 群联系人 - 置顶群
func GroupContactTop(w http.ResponseWriter, r *http.Request) {
	// verify user legal
	user, code, errStr := service.NewUserServ().CheckUserRequestLegal(r)
	if errStr != "" {
		comm.ResFailure(w, code, errStr)
		return
	}

	r.ParseForm()
	groupIdStr := r.PostFormValue("group_id")
	if groupIdStr == "" || funcs.IsNumber(groupIdStr) {
		comm.ResFailure(w, 1001, "params failure")
		return
	}
	groupId, _ := strconv.ParseInt(groupIdStr, 10, 64)

	groupContact, err := service.NewGroupContactServ().Info(user.Id, groupId)
	if err != nil {
		comm.ResFailure(w, 2001, "change failure")
		return
	}

	groupContact, err = service.NewGroupContactServ().ChangeStatus(
		groupContact.Id, models.GROUP_CONTACT_STATUS_GROUP_TOP)
	if err != nil {
		comm.ResFailure(w, 2103, "change failure")
		return
	}

	comm.ResSuccess(w, groupContact)
}
