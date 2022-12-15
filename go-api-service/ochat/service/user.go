package service

import (
	"errors"
	"ochat/comm"
	"ochat/models"
	"time"

	"xorm.io/xorm"
)

type UserService struct {
	DB *xorm.Engine
}

func (s *UserService) Register(
	mobile, avatar, nickname, password string, sex int) (user models.User, err error) {

	userInfo := models.User{}
	_, err = s.DB.Where("mobile = ?", mobile).Get(&userInfo)

	if err == nil && userInfo.Id > 0 {
		errStr := "the user to which the current mobile phone number belongs exists"
		return userInfo, errors.New(errStr)
	}

	salt := comm.RandStr(6, comm.Rand_Str_Level_5)

	token := comm.GenerateToken(password + salt)

	userInfo = models.User{
		Mobile:     mobile,
		Avatar:     avatar,
		Nickname:   nickname,
		Sex:        sex,
		Password:   comm.GeneratePasswd(password, salt),
		Salt:       salt,
		Created_at: time.Now(),
		Token:      token,
	}

	if num, err := s.DB.InsertOne(&userInfo); err != nil || num <= 0 {
		errStr := "user data insert database failure"
		return userInfo, errors.New(errStr)
	}

	return userInfo, nil
}

func (s *UserService) Login(mobile, password string) (user models.User, err error) {
	_, err = s.DB.Where("mobile = ?", mobile).Get(&user)
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

func (s *UserService) UpToken(user_id int) (user models.User, err error) {
	_, err = s.DB.Where("id = ?", user_id).Get(&user)
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
