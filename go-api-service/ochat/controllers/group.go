package controllers

import (
	"net/http"
	"ochat/comm"
	"ochat/service"
)

// 群 - 查看群信息
func Group(w http.ResponseWriter, r *http.Request) {
	// verify user legal
	_, code, errStr := service.NewUserServ().CheckUserRequestLegal(r)
	if errStr != "" {
		comm.ResFailure(w, code, errStr)
		return
	}

}

// 群 - 查看用户的所有群信息
func GroupList(w http.ResponseWriter, r *http.Request) {
	// verify user legal
	user, code, errStr := service.NewUserServ().CheckUserRequestLegal(r)
	if errStr != "" {
		comm.ResFailure(w, code, errStr)
		return
	}

	groups, err := service.NewGroupServ().UserList(user.Id)
	if err != nil {
		comm.ResFailure(w, 2101, err.Error())
		return
	}

	comm.ResSuccess(w, groups)
}
