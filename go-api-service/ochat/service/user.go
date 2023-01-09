package service

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"ochat/bootstrap"
	"ochat/comm/funcs"
	"ochat/models"
	"strconv"
	"time"

	"xorm.io/xorm"
)

type UserService struct {
	DB *xorm.Engine
}

func NewUserServ() *UserService {
	return &UserService{
		DB: bootstrap.DB_Engine,
	}
}

func (s *UserService) Register(
	mobile, username, avatar, nickname, password string,
	sex int) (user models.User, err error) {

	userInfo, err := s.MobileToUserInfo(mobile)
	// 信息验证
	if err == nil && userInfo.Id > 0 {
		errStr := "the user to which the current mobile phone number belongs exists"
		return userInfo, errors.New(errStr)
	}

	salt := funcs.RandStr(12, funcs.Rand_Str_Level_5)
	token := funcs.GenerateToken(password + salt)

	// 构建用户数据
	userInfo = models.User{
		Mobile:     mobile,
		Username:   username,
		Avatar:     avatar,
		Nickname:   nickname,
		Sex:        sex,
		Password:   funcs.GeneratePasswd(password, salt),
		Salt:       salt,
		Token:      token,
		Status:     1,
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}

	// 保存数据
	if num, err := s.DB.InsertOne(&userInfo); err != nil || num <= 0 {
		errStr := "user data insert database failure"
		return userInfo, errors.New(errStr)
	}

	return userInfo, nil
}

func (s *UserService) Login(mobile, password string) (user models.User, err error) {
	user, err = s.MobileToUserInfo(mobile)
	if err != nil {
		return
	}

	if user.Id == 0 {
		return
	}

	if !funcs.VaildataPasswd(password, user.Salt, user.Password) {
		return models.User{}, errors.New("password vaildate failute")
	}

	// 更新token
	user, err = s.UpToken(user)
	if err != nil {
		return models.User{}, errors.New("failure: user data get failure")
	}

	return user, nil
}

func (s *UserService) UpToken(user models.User) (models.User, error) {
	if user.Id == 0 || user.Token == "" {
		return models.User{}, errors.New("update failure")
	}

	token := funcs.GenerateToken(user.Password + user.Salt)
	if token != user.Token {
		user.Token = token
		user.Updated_at = time.Now()
		num, err := s.DB.ID(user.Id).Cols("token", "updated_at").Update(&user)
		if err != nil || num < 1 {
			return user, errors.New("update failure")
		}
	}

	return user, nil
}

func (s *UserService) CheckToken(user_id int64, token string) bool {
	user, err := s.UserIdToUserInfo(user_id)
	if err != nil || user.Id == 0 || user.Token == "" {
		return false
	}

	if user.Token != token {
		return false
	}

	return true
}

func (s *UserService) MobileToUserInfo(mobile string) (models.User, error) {
	var user models.User
	_, err := s.DB.Where("mobile = ?", mobile).Get(&user)

	return user, err
}

func (s *UserService) UserIdToUserInfo(userId int64) (models.User, error) {
	var user models.User
	_, err := s.DB.Where("id = ?", userId).Get(&user)

	return user, err
}

func (s *UserService) UsernameToUserInfo(username string) (models.User, error) {
	var user models.User
	_, err := s.DB.Where("username = ?", username).Get(&user)

	return user, err
}

func (s *UserService) CheckUserRequestLegal(r *http.Request) (userInfo models.User, code int, errStr string) {
	userIdStr := r.FormValue("user_id") // 用户id
	token := r.FormValue("token")       // 用户token

	if userIdStr == "" || !funcs.IsNumber(userIdStr) {
		return userInfo, 101, "the user id params is empty or is illegal"
	}

	if token == "" {
		return userInfo, 102, "the user token params is empty"
	}

	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		return userInfo, 103, "the user id params is illegal"
	}

	userInfo, err = s.UserIdToUserInfo(userId)
	if err != nil || userInfo.Id == 0 {
		return userInfo, 104, "user are dose not exists"
	}

	// 验证token是否合法
	if userInfo.Token != token {
		return userInfo, 105, "token parameter validation failed"
	}

	return
}

func (s *UserService) CreateQrCode(user models.User) (filename string, err error) {
	// 生成二维码
	qrCodeUrl := fmt.Sprintf("%s/user?user_id=%d", bootstrap.HTTP_HOST, user.Id)
	filename, err = funcs.QrCode(qrCodeUrl, "user_qrcode")
	if err != nil {
		return "", err
	}

	user.QrCode = funcs.GetImgUrl("user_qrcode", filename)
	user.Updated_at = time.Now()

	num, err := s.DB.ID(user.Id).Cols("qr_code", "updated_at").Update(&user)
	if err != nil || num < 1 {
		return "", errors.New("update failure")
	}

	return
}

func (s *UserService) UpdateFields(fileds url.Values, userId int64, resetData bool) (user models.User, err error) {
	canUpFailds := []string{
		"mobile",   // 手机号
		"nickname", // 用户昵称
		"password", // 密码
		"about",    // 简单描述
		"avatar",   // 头像
		"sex",      // 性别,0:无;1:男;2:女;
		"birthday", // 生日
	}

	upFields := map[string]string{}
	for _, field := range canUpFailds {
		if fileds.Has(field) && fileds.Get(field) != "" {
			upFields[field] = fileds.Get(field)
		}
	}

	if len(upFields) == 0 {
		return models.User{}, errors.New("no update")

	}

	upFields["updated_at"] = funcs.UpdateTime()

	_, err = s.DB.Table("user").Where("id = ?", userId).Update(upFields)
	if err != nil {
		return models.User{}, err
	}

	if !resetData {
		return models.User{}, nil
	}

	user, err = s.UserIdToUserInfo(userId)
	if err != nil {
		return user, errors.New("get user info failure")
	}

	return user, nil
}
