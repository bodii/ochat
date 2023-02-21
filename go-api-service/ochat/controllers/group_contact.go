package controllers

import (
	"net/http"
	"ochat/comm"
	"ochat/comm/funcs"
	"ochat/models"
	"ochat/service"
	"strconv"
)

// 群联系人 - 列表
func GroupContactList(w http.ResponseWriter, r *http.Request) {
	// verify user legal
	_, code, errStr := service.NewUserServ().CheckUserRequestLegal(r)
	if errStr != "" {
		comm.ResFailure(w, code, errStr)
		return
	}

	groupIdStr := r.FormValue("group_id")
	groupId, _ := strconv.ParseInt(groupIdStr, 10, 64)

	contactList, err := service.NewGroupContactServ().GetMembers(groupId)
	if err != nil {
		comm.ResFailure(w, 2101, "get failure")
	}

	comm.ResSuccess(w, contactList)
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
		"is_self":            userId == user.UserId,
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

	groupContact, err := service.NewGroupContactServ().Info(user.UserId, groupId)
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

	err = service.NewGroupContactServ().ChangeStatus(
		&kickGroupContact, models.GROUP_CONTACT_STATUS_KICK_OUT)
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

	groupContact, err := service.NewGroupContactServ().Info(user.UserId, groupId)
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

	}
	// TODO：如果当前群联系人只剩一个人时，退出后，群解散

	err = service.NewGroupContactServ().ChangeStatus(
		&groupContact, models.GROUP_CONTACT_STATUS_EXIT)
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

	groupContact, err := service.NewGroupContactServ().Info(user.UserId, groupId)
	if err != nil {
		comm.ResFailure(w, 2001, "change failure")
		return
	}

	err = service.NewGroupContactServ().ChangeStatus(
		&groupContact, models.GROUP_CONTACT_STATUS_GROUP_TOP)
	if err != nil {
		comm.ResFailure(w, 2103, "change failure")
		return
	}

	comm.ResSuccess(w, groupContact)
}

// 群联系人 - 设置管理员
func GroupContactManager(w http.ResponseWriter, r *http.Request) {
	// verify user legal
	user, code, errStr := service.NewUserServ().CheckUserRequestLegal(r)
	if errStr != "" {
		comm.ResFailure(w, code, errStr)
		return
	}

	r.ParseForm()

	groupIdStr := r.PostFormValue("group_id")
	memberIdStr := r.PostFormValue("member_id")
	setType := r.PostFormValue("type") // cancel、upgrade
	if groupIdStr == "" || !funcs.IsNumber(groupIdStr) {
		comm.ResFailure(w, 1101, "group_id param failure")
		return
	}

	if memberIdStr == "" || !funcs.IsNumber(memberIdStr) {
		comm.ResFailure(w, 1102, "manager_id param failure")
		return
	}

	if setType == "" || (setType != "cancel" && setType != "upgrade") {
		comm.ResFailure(w, 1103, "type param failure")
		return
	}

	groupId, _ := strconv.ParseInt(groupIdStr, 10, 64)
	memberId, _ := strconv.ParseInt(memberIdStr, 10, 64)
	if user.UserId == memberId {
		comm.ResFailure(w, 2101, "current user and input member is the same person")
		return
	}

	userGroupContact, err := service.NewGroupContactServ().Info(user.UserId, groupId)
	if err != nil {
		comm.ResFailure(w, 2101, "current user the group contact info is exists")
		return
	}
	if userGroupContact.Type < models.GROUP_CONTACT_TYPE_MANAGER {
		comm.ResFailure(w, 2102, "current user is not group master or manager")
		return
	}

	memberGroupContact, err := service.NewGroupContactServ().Info(memberId, groupId)
	if err != nil {
		comm.ResFailure(w, 2103, "current member the group contact info is exists")
		return
	}

	if setType == "upgrade" {
		if memberGroupContact.Type > models.GROUP_CONTACT_TYPE_MEMBER {
			comm.ResFailure(w, 2104, "current member is not ordinary member")
			return
		}

		service.NewGroupContactServ().ChangeType(&memberGroupContact,
			models.GROUP_CONTACT_TYPE_MANAGER)
	}

	if setType == "cancel" {
		if memberGroupContact.Type < models.GROUP_CONTACT_TYPE_MANAGER {
			comm.ResFailure(w, 2105, "current member is not ordinary member")
			return
		}

		service.NewGroupContactServ().ChangeType(&memberGroupContact,
			models.GROUP_CONTACT_TYPE_MEMBER)
	}

	comm.ResSuccess(w, memberGroupContact)
}

// 群联系人 - 更新信息
func GroupContactUpField(w http.ResponseWriter, r *http.Request) {
	// verify user legal
	user, code, errStr := service.NewUserServ().CheckUserRequestLegal(r)
	if errStr != "" {
		comm.ResFailure(w, code, errStr)
		return
	}

	r.ParseForm()
	groupIdstr := r.PostFormValue("group_id")
	if groupIdstr == "" || !funcs.IsNumber(groupIdstr) {
		comm.ResFailure(w, 1001, "group id param failure")
		return
	}

	groupId, _ := strconv.ParseInt(groupIdstr, 10, 64)
	err := service.NewGroupContactServ().UpdateFields(user.UserId, groupId, r.PostForm)
	if err != nil {
		comm.ResFailure(w, 1001, err.Error())
		return
	}

	comm.ResSuccess(w, nil)

}
