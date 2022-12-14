package service

import (
	"errors"
	"ochat/comm"
	"ochat/models"
	"time"
)

type UserService struct {
}

func (s *UserService) Register(
	mobile, avatar, nickname, passwd, sex string) (user models.User, err error) {

	userInfo := models.User{}
	_, err = DB.Where("mobile = ?", mobile).Get(&userInfo)

	if err == nil && userInfo.Id > 0 {
		errStr := "the user to which the current mobile phone number belongs exists"
		return userInfo, errors.New(errStr)
	}

	salt := comm.RandStr(6, comm.RandStrlevel5)

	token := comm.GenerateToken(passwd + salt)

	userInfo = models.User{
		Mobile:     mobile,
		Avatar:     avatar,
		Nickname:   nickname,
		Sex:        sex,
		Passwd:     comm.GeneratePasswd(passwd, salt),
		Salt:       salt,
		Created_at: time.Now(),
		Token:      token,
	}

	if num, err := DB.InsertOne(&userInfo); err != nil || num <= 0 {
		errStr := "user data insert database failure"
		return userInfo, errors.New(errStr)
	}

	return userInfo, nil
}

func (s *UserService) Login(mobile, passwd string) (user models.User, err error) {
	userInfo := models.User{}
	if _, err = DB.Where("mobile = ?", mobile).Get(&userInfo); err != nil {
		return userInfo, err
	}

	if userInfo.Id == 0 {
		return userInfo, err
	}

	if !comm.VaildataPasswd(passwd, userInfo.Salt, user.Passwd) {
		return models.User{}, errors.New("passwd vaildate failute")
	}

	return userInfo, nil
}

func (s *UserService) UpToken(user_id int) (user models.User, err error) {
	userInfo := models.User{}
	if _, err = DB.Where("id = ?", user_id).Get(&userInfo); err != nil {
		return userInfo, err
	}

	if userInfo.Id == 0 {
		return models.User{}, err
	}

	token := comm.GenerateToken(userInfo.Passwd + userInfo.Salt)

	userInfo.Token = token
	num, err := DB.ID(user_id).Cols("token").Update(&userInfo)
	if err != nil || num < 1 {
		return userInfo, errors.New("update failure")
	}

	return userInfo, nil
}
