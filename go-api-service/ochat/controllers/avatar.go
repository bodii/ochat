package controllers

import (
	"net/http"
	"ochat/comm"
	"ochat/comm/funcs"
	"ochat/service"
)

// 上传头像图片
func AvatarUpload(w http.ResponseWriter, r *http.Request) {
	// verify user legal
	_, code, errStr := service.NewUserServ().CheckUserRequestLegal(r)
	if errStr != "" {
		comm.ResFailure(w, code, errStr)
		return
	}

	filename, _, err := funcs.UploadFile(r, "avatar_image", "avatar")
	if err != nil {
		comm.ResFailure(w, 1001, "upload avatar file failure")
	}

	url := funcs.GetImgUrl("avatar", filename)
	comm.ResSuccess(w, comm.D{
		"avatar_url": url,
	})
}
