package controllers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"ochat/bootstrap"
	"ochat/comm"
	"ochat/comm/funcs"
	"ochat/service"
)

func LoginQRCode(w http.ResponseWriter, r *http.Request) {
	htttpHost := bootstrap.HTTP_HOST
	ip := funcs.ClientIP(r)
	url := fmt.Sprintf(
		"%s/user/login_qrcode?ip=%s",
		htttpHost, ip)

	QRCodeStream, err := funcs.QrCode(url)
	if err != nil {
		log.Printf("open login qrcode file failure: %s", url)
		comm.ResFailure(w, 1002, "get qrcode file failure")
		return
	}

	io.Copy(w, QRCodeStream)
}

func LoginQRCodeScan(w http.ResponseWriter, r *http.Request) {
	userInfo, code, errStr := service.NewUserServ().CheckUserRequestLegal(r)
	if errStr != "" {
		comm.ResFailure(w, code, errStr)
		return
	}

	comm.ResSuccess(w, userInfo)
}
