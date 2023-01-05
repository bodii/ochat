package service

import (
	"ochat/bootstrap"
	"ochat/models"

	"xorm.io/xorm"
)

type FriendService struct {
	DB *xorm.Engine
}

func NewFriendServ() *FriendService {
	return &FriendService{
		DB: bootstrap.DB_Engine,
	}
}

func (f *FriendService) List(userId int64, status int) (users []models.User, err error) {

	err = f.DB.Table("friend").Alias("f").Join("left", "user as u", "f.friend_id=u.id").
		Where("f.userid = ? and f.status = ? and u.status = ?",
						userId, status, models.USER_STATUS_VALID).
		Asc("u.nickname_prefix"). // 根据用户的昵称前缀和朋友表的id排序
		Asc("f.id").
		Find(&users)

	return users, err
}
