package controllers

import (
	"net/http"
	"ochat/bootstrap"
	"ochat/comm"
	"ochat/comm/funcs"
	"ochat/service"
	"strconv"
)

func UserLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	mobile := r.PostForm.Get("mobile")
	password := r.PostForm.Get("password")

	userInfo, err := service.NewUserServ().Login(mobile, password)
	if err != nil {
		comm.ResFailure(w, 1001, err.Error())
		return
	}

	if userInfo.Id == 0 {
		comm.ResFailure(w, 1002, "failure: user data is empty")
		return
	}

	comm.ResSuccess(w, userInfo)
}

func UserRegister(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	mobile := r.Form.Get("mobile")
	username := r.Form.Get("username")
	avatar := r.Form.Get("avatar")
	nickname := r.Form.Get("nickname")
	passwd := r.Form.Get("password")
	sex, _ := strconv.Atoi(r.FormValue("sex"))

	if mobile == "" {
		comm.ResFailure(w, 1001, "register failure: mobile is empty!")
		return
	}

	if username == "" {
		comm.ResFailure(w, 1002, "register failure: username is empty!")
		return
	}

	if passwd == "" {
		comm.ResFailure(w, 1003, "register failure: password is empty!")
		return
	}

	if avatar == "" {
		avatar = bootstrap.HTTP_Avatar_URI
		switch sex {
		case 1:
			avatar += "avatar_boy_kid_person_icon.png"
		case 2:
			avatar += "child_girl_kid_person_icon.png"
		default:
			avatar += "avatar_boy_male_user_young_icon.png"
		}
	}

	if nickname == "" {
		nickname = funcs.RandStr(20, 5)
	}

	userInfo, err := service.NewUserServ().Register(mobile, username, avatar, nickname, passwd, sex)
	if err != nil {
		comm.ResFailure(w, 1001, "register failure: "+err.Error())
		return
	}

	if userInfo.Id == 0 {
		comm.ResFailure(w, 1001, "register failure: not insert data")
		return
	}

	comm.ResSuccess(w, userInfo)
}
