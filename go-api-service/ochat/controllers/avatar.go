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
	user, code, errStr := service.NewUserServ().CheckUserRequestLegal(r)
	if errStr != "" {
		comm.ResFailure(w, code, errStr)
		return
	}

	filename, _, err := funcs.UploadFile(r, "avatar_image", "user_avatar")
	if err != nil {
		comm.ResFailure(w, 1001, "upload avatar file failure")
	}

	url, _ := funcs.GetImgUrl("user_avatar", filename)

	r.PostForm.Add("avatar", url)
	user.Avatar = url
	_, err = service.NewUserServ().UpdateFields(r.PostForm, user.UserId, false)
	if err != nil {
		comm.ResFailure(w, 1201, err.Error())
		return
	}

	comm.ResSuccess(w, comm.D{
		"avatar_url": url,
	})
}
