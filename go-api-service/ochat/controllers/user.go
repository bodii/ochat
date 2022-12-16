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
	username := r.Form.Get("username")
	avatar := r.Form.Get("avatar")
	nickname := r.Form.Get("nickname")
	passwd := r.Form.Get("password")
	sex, _ := strconv.Atoi(r.FormValue("sex"))

	if mobile == "" {
		comm.Res(w, 1001, "register failure: mobile is empty!", nil)
		return
	}

	if username == "" {
		comm.Res(w, 1002, "register failure: username is empty!", nil)
		return
	}

	if passwd == "" {
		comm.Res(w, 1003, "register failure: password is empty!", nil)
		return
	}

	if avatar == "" {
		switch sex {
		case 1:
			avatar = "avatar_boy_kid_person_icon.png"
		case 2:
			avatar = "child_girl_kid_person_icon.png"
		default:
			avatar = "avatar_boy_male_user_young_icon.png"
		}
	}

	if nickname == "" {
		nickname = comm.RandStr(20, 5)
	}

	userServ := &service.UserService{
		DB: bootstrap.DB_Engine,
	}

	userInfo, err := userServ.Register(mobile, username, avatar, nickname, passwd, sex)
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
