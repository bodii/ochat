package service

import (
	"errors"
	"fmt"
	"ochat/bootstrap"
	"ochat/comm"
	"ochat/models"
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

	if err == nil && userInfo.Id > 0 {
		errStr := "the user to which the current mobile phone number belongs exists"
		return userInfo, errors.New(errStr)
	}

	salt := comm.RandStr(12, comm.Rand_Str_Level_5)
	token := comm.GenerateToken(password + salt)

	avatar = fmt.Sprintf("%s%s", bootstrap.HTTP_Avatar_URI, avatar)

	userInfo = models.User{
		Mobile:     mobile,
		Username:   username,
		Avatar:     avatar,
		Nickname:   nickname,
		Sex:        sex,
		Password:   comm.GeneratePasswd(password, salt),
		Salt:       salt,
		Token:      token,
		About:      "",
		Status:     1,
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}

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

	if !comm.VaildataPasswd(password, user.Salt, user.Password) {
		return models.User{}, errors.New("password vaildate failute")
	}

	return user, nil
}

func (s *UserService) UpToken(user_id int64) (user models.User, err error) {
	user, err = s.UserIdToUserInfo(user_id)
	if err != nil {
		return user, err
	}

	if user.Id == 0 {
		return models.User{}, err
	}

	token := comm.GenerateToken(user.Password + user.Salt)

	user.Token = token
	num, err := s.DB.ID(user_id).Cols("token").Update(&user)
	if err != nil || num < 1 {
		return user, errors.New("update failure")
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
