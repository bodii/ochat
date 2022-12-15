package controllers

import (
	"net/http"
	"ochat/bootstrap"
	"ochat/comm"
	"ochat/service"
	"strconv"
)

func Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	userServ := &service.UserService{
		DB: bootstrap.DB_Engine,
	}

	mobile := r.PostForm.Get("mobile")
	password := r.PostForm.Get("password")

	userInfo, err := userServ.Login(mobile, password)
	if err != nil {
		comm.Res(w, 1001, err.Error(), nil)
		return
	}

	if userInfo.Id == 0 {
		comm.Res(w, 1002, "failure: user data is empty", nil)
		return
	}

	comm.Res(w, 200, "success", userInfo)
}

func Register(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	mobile := r.Form.Get("mobile")
	avatar := r.Form.Get("avatar")
	nickname := r.Form.Get("nickname")
	passwd := r.Form.Get("passwd")
	sex, _ := strconv.Atoi(r.FormValue("sex"))

	userServ := &service.UserService{
		DB: bootstrap.DB_Engine,
	}

	userInfo, err := userServ.Register(mobile, avatar, nickname, passwd, sex)
	if err != nil {
		comm.Res(w, 1001, "register failure: "+err.Error(), nil)
		return
	}

	if userInfo.Id == 0 {
		comm.Res(w, 1001, "register failure: not insert data", nil)
		return
	}

	comm.Res(w, 200, "register success", userInfo)
}
