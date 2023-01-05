package service

import (
	"ochat/bootstrap"
	"ochat/comm/funcs"
	"ochat/models"
	"time"

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
		Asc("f.Alias_prefix"). // 根据用户的昵称前缀和朋友表的id排序
		Asc("f.id").
		Find(&users)

	return users, err
}

func (f *FriendService) Add(userId, friendId int64, friendAlias, about string) (friendInfo models.Friend, err error) {
	// 获取昵称前缀
	aliasPrefix := funcs.StrPrefix(friendAlias, 1, 2)
	// 如果前缀首个字符不是英文或中文，则返回＃
	if !funcs.IsEnglish(aliasPrefix) {
		aliasPrefix = "#"
	}
	friendInfo = models.Friend{
		UserId:      userId,
		FriendId:    friendId,
		FriendAlias: friendAlias,
		AliasPrefix: aliasPrefix,
		About:       about,
		Status:      1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return friendInfo, nil
}
