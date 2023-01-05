package controllers

import (
	"net/http"
	"ochat/comm"
	"ochat/comm/funcs"
	"ochat/service"
	"strconv"
	"time"
)

// get phone sms verify code
func PhoneSms(w http.ResponseWriter, r *http.Request) {
	phone := r.FormValue("phone")       // 手机号
	userIdStr := r.FormValue("user_id") // 用户id
	token := r.FormValue("token")       // 用户token

	if userIdStr == "" || token == "" {
		comm.ResFailure(w, 1001, "the user params is empty")
		return
	}

	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	userInfo, err := service.NewUserServ().UserIdToUserInfo(userId)
	if err != nil {
		comm.ResFailure(w, 1002, "user are dose not exists")
		return
	}

	// 验证token是否合法
	if userInfo.Token != token {
		comm.ResFailure(w, 1003, "token parameter validation failed")
		return
	}

	// 查看缓存中是否存在
	verifyCode, err := service.NewRedis().Get(service.REDIS_CTX, phone).Result()
	if err != nil {
		// 如果不存在，则生成一个新的
		verifyCode = funcs.RandStr(6, 3)
		// 添加到缓存(前端默认3分钟有效)
		service.NewRedis().SetEX(service.REDIS_CTX, phone, verifyCode, 4*time.Minute)
	}

	// 返回数据
	comm.ResSuccess(w, map[string]any{
		"verify_code": verifyCode,
	})
}

// verify that the mobile phone verification code is valid
func PhoneSmsVerify(w http.ResponseWriter, r *http.Request) {
	// verify user legal
	userInfo, code, errStr := service.NewUserServ().CheckUserRequestLegal(r)
	if errStr != "" {
		comm.ResFailure(w, code, errStr)
		return
	}

	phone := r.FormValue("phone")            // 手机号
	verifyCode := r.FormValue("verify_code") // 验证码

	if phone == "" || !funcs.IsMobile(phone) {
		comm.ResFailure(w, 1001, "the phone params is empty or is illegal")
	}

	if verifyCode == "" {
		comm.ResFailure(w, 1002, "the verify code params is empty")
		return
	}

	// 查看缓存中是否存在
	phoneCode, err := service.NewRedis().Get(service.REDIS_CTX, phone).Result()
	// 如果不存在
	if err != nil {
		comm.ResFailure(w, 1003, "verify code are dose not exists or defunct")
		return
	}

	// 验证码错误
	if verifyCode != phoneCode {
		comm.ResFailure(w, 1004, "verify code are incorrect")
		return
	}

	// 返回成功
	comm.ResSuccess(w, userInfo)
}
