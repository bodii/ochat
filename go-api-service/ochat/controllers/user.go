package controllers

import (
	"log"
	"net/http"
	"ochat/comm"
	"ochat/service"
)

func Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	mobile := r.PostForm.Get("mobile")
	passwd := r.PostForm.Get("passwd")
	log.Printf("%v", r.PostForm)

	userServ := service.UserService{}
	userInfo, err := userServ.Login(mobile, passwd)
	if err != nil || userInfo.Id == 0 {
		comm.Res(w, 1001, err.Error(), nil)
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
	sex := r.Form.Get("sex")

	userServ := &service.UserService{}
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
