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
	"os"
	"path"
	"strconv"

	"github.com/skip2/go-qrcode"
)

func LoginQRCode(w http.ResponseWriter, r *http.Request) {
	htttpHost := bootstrap.HTTP_HOST
	ip := funcs.ClientIP(r)
	loginUrl := fmt.Sprintf(
		"%s/user/login_qrcode?ip=%s",
		htttpHost, ip)

	QRCodeConf := bootstrap.SystemConf.LoginQRCode
	filename := funcs.RandFileName(".png")
	filePath := path.Join(bootstrap.PROJECT_DIR, QRCodeConf.FileDir, filename)
	err := qrcode.WriteFile(loginUrl, qrcode.Medium, 256, filePath)
	// png, err := qrcode.Encode(loginUrl, qrcode.Medium, 256)
	if err != nil {
		log.Printf("client ip: %s create login qrcode failure", ip)
		comm.ResFailure(w, 1001, "create qrcode failure")
		return
	}

	QRCodeStream, err := os.Open(filePath)
	if err != nil {
		log.Printf("open login qrcode file failure: %s", filePath)
		comm.ResFailure(w, 1002, "get qrcode file failure")
		return
	}

	io.Copy(w, QRCodeStream)
}

func LoginQRCodeScan(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// ip := r.PostFormValue("ip")
	token := r.PostFormValue("token")
	useridStr := r.PostFormValue("user_id")
	if token == "" || useridStr == "" {
		comm.ResFailure(w, 1001, "the param are incorrect")
		return
	}
	userid, _ := strconv.ParseInt(useridStr, 10, 64)

	userInfo, err := service.NewUserServ().UserIdToUserInfo(userid)
	if err != nil || userInfo.Id == 0 {
		comm.ResFailure(w, 1002, "login failure")
		return
	}

	if token != userInfo.Token {
		comm.ResFailure(w, 1003, "wrongful login")
		return
	}

	// 更新token
	// userInfo, err = service.NewUserServ().UpToken(userid)

	comm.ResSuccess(w, userInfo)
}
