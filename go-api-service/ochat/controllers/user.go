package controllers

import (
	"net/http"
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
	// 更新token
	userInfo, err = service.NewUserServ().UpToken(userInfo.Id)
	if err != nil {
		comm.ResFailure(w, 1003, "failure: user data get failure")
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
		avatarFilename := funcs.DefaultAvatar(sex)
		avatar = funcs.GetImgUrl("avatar", avatarFilename)
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

func UserQrCode(w http.ResponseWriter, r *http.Request) {
	// verify user legal
	userInfo, code, errStr := service.NewUserServ().CheckUserRequestLegal(r)
	if errStr != "" {
		comm.ResFailure(w, code, errStr)
		return
	}

	// 如果二维码不存在，则创建
	if userInfo.QrCode == "" {
		filename, err := service.NewUserServ().CreateQrCode(userInfo)
		if err != nil {
			comm.ResFailure(w, 3001, "create qr code failure")
			return
		}
		userInfo.QrCode = funcs.GetImgUrl("user_qrcode", filename)
	}

	comm.ResSuccess(w, comm.D{
		"user_info": userInfo,
	})
}
