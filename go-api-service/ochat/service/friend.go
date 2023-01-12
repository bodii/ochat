package service

import (
	"net/http"
	"ochat/bootstrap"
	"ochat/comm/funcs"
	"ochat/models"
	"strconv"
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

// adds func: petitioner and responder both parties as friends
//
//	params:
//	 - pet: petitioner user info
//	 - res: responder user info
//	return:
//	 - ok: whether to add successfully
//	 - err: error message
func (f *FriendService) Adds(pet models.User, res models.User) (ok bool, err error) {
	friends := make([]*models.Friend, 2)
	friends[0] = &models.Friend{
		UserId:      res.Id,
		FriendId:    pet.Id,
		FriendAlias: pet.Nickname,
		AliasPrefix: funcs.StrPrefix(pet.Nickname, 1, 2),
		Status:      models.FRIEND_STATUS_FRIENDED,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	friends[1] = &models.Friend{
		UserId:      pet.Id,
		FriendId:    res.Id,
		FriendAlias: res.Nickname,
		AliasPrefix: funcs.StrPrefix(res.Nickname, 1, 2),
		Status:      models.FRIEND_STATUS_FRIENDED,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	num, err := f.DB.Insert(friends)
	if err != nil || num == 0 {
		return false, err
	}

	return true, nil
}

func (f *FriendService) Update(friend models.Friend, cols []string) (ok bool, err error) {
	canUpdateFields := map[string]bool{
		"friend_alias": true,
		"alias_prefix": true,
		"about":        true,
	}

	updateFields := make([]string, 0)
	for _, f := range cols {
		if canUpdateFields[f] {
			updateFields = append(updateFields, f)
		}
	}

	num, err := f.DB.Where("id = ?", friend.Id).Cols(updateFields...).Update(friend)
	if err != nil || num == 0 {
		return false, err
	}

	return true, nil
}

func (f *FriendService) UpdateStatus(r *http.Request, userId int64, status int) (ok bool, err error) {
	r.ParseForm()
	friendIdStr := r.PostFormValue("friend_id")
	friendId, err := strconv.ParseInt(friendIdStr, 10, 64)
	if err != nil {
		return false, err
	}

	friend, err := f.Info(userId, friendId)
	if err != nil {
		return false, err
	}

	friend.Status = status
	friend.UpdatedAt = time.Now()
	num, err := f.DB.Where("id = ?", friend.Id).Cols("status", "updated_at").Update(friend)
	if err != nil || num == 0 {
		return false, err
	}

	return true, nil
}

func (f *FriendService) Info(userId, friendId int64) (friend models.Friend, err error) {
	_, err = f.DB.Where(
		"user_id = ? and friend_id = ?", userId, friendId).Get(&friend)

	return
}
