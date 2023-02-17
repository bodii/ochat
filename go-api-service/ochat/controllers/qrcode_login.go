package controllers

import (
	"fmt"
	"log"
	"net/http"
	"ochat/bootstrap"
	"ochat/comm"
	"ochat/comm/funcs"
	"ochat/service"
)

func LoginQRCode(w http.ResponseWriter, r *http.Request) {
	ip := funcs.ClientIP(r)
	if ip == "" {
		comm.ResFailure(w, 1001, "ip get failure")
	}

	htttpHost := bootstrap.HTTP_HOST
	url := fmt.Sprintf(
		"%s/user/login_qrcode?ip=%s",
		htttpHost, ip)

	filename, err := funcs.QrCode(url, "login_qrcode")
	if err != nil {
		log.Printf("open login qrcode file failure: %s", url)
		comm.ResFailure(w, 1002, "get qrcode file failure")
		return
	}

	qrcodeUrl := funcs.GetImgUrl("login_qrcode", filename)
	comm.ResSuccess(w, comm.D{
		"qrcode_url": qrcodeUrl,
	})
}

func LoginQRCodeScan(w http.ResponseWriter, r *http.Request) {
	user, code, errStr := service.NewUserServ().CheckUserRequestLegal(r)
	if errStr != "" {
		comm.ResFailure(w, code, errStr)
		return
	}

	comm.ResSuccess(w, user)
}
