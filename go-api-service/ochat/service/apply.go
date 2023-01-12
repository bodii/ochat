package service

import (
	"errors"
	"ochat/bootstrap"
	"ochat/models"
	"time"

	"xorm.io/xorm"
)

type ApplyService struct {
	DB *xorm.Engine
}

func NewApplyServ() *ApplyService {
	return &ApplyService{
		DB: bootstrap.DB_Engine,
	}
}

// 添加好友/群申请
//
// params:
//   - userId: current user id
//   - addUserId: add the user id
//   - comment:
//   - applyType: 申请类型,1:好友;2:群
//
// return: user list or error
func (a *ApplyService) Add(userId, addUserId int64, comment string, addType int) (models.Apply, error) {
	// TODO: 添加时，查看被添加者是否需要验证 user.friend_verify

	addData := models.Apply{
		Petitioner: userId,
		Responder:  addUserId,
		Status:     models.APPLY_STATUS_UNREAD,
		Type:       models.APPLY_TYPE_USER,
		Comment:    comment,
		CreatedAt:  time.Now(),
	}

	if num, err := a.DB.InsertOne(&addData); err != nil || num == 0 {
		errStr := "apply info insert database failure"
		return addData, errors.New(errStr)
	}

	return addData, nil
}

// 查看向当前用户发起申请的状态的好友列表
//
// params:
//   - userId: current user id
//   - status: 应答者是否同意,-1:拒绝;0:未查看;1:已查看;2:同意
//   - applyType: 申请类型,1:好友;2:群
//
// return: user list or error
func (a *ApplyService) List(userId int64, status, applyType int) ([]models.User, error) {
	userInfos := make([]models.User, 0)

	err := a.DB.Table("apply").Join("left", "user", "apply.petitioner = user.id").
		Where(
			"apply.responder = ? and apply.status = ? and apply.type = ?",
			userId, status, applyType).
		Asc("id").
		Find(&userInfos)

	return userInfos, err
}

// 设置/更新申请状态
//
// params:
//   - id: apply id
//   - status: 应答者是否同意,-1:拒绝;0:未查看;1:已查看;2:同意
//   - user: 处理设置的用户
//
// return: whether to set up successfully
func (a *ApplyService) SetStatus(id int64, status int, user models.User) (bool, error) {
	updateData := models.Apply{
		Status:    status,
		UpdatedAt: time.Now(),
	}

	// 查询申请信息
	apply := models.Apply{}
	exists, err := a.DB.Where("id = ? and responder = ?", id, user.Id).
		And("type = ?", models.APPLY_TYPE_USER).
		Get(&apply)
	if err != nil || !exists {
		return false, errors.New("get apply data failure")
	}

	// 如果已拒绝，则无法更改
	if apply.Status == models.APPLY_STATUS_REFUSE ||
		apply.Status == models.APPLY_STATUS_AGREE {

		return false, errors.New("operation failure")
	}

	// 如果是添加，先查询是否已存在
	if status == models.APPLY_STATUS_AGREE {
		if apply.Type == models.APPLY_TYPE_USER {
			firend, err := NewFriendServ().Info(user.Id, id)
			if err == nil && firend.Id > 0 {
				return false, errors.New("current friend exists")
			}
		} else if apply.Type == models.APPLY_TYPE_GROUP {
			groupContact, err := NewGroupContactServ().Info(user.Id, id)
			if err == nil && groupContact.Id > 0 {
				return false, errors.New("current group contact exists")
			}
		}
	}

	// 添加好友信息
	if status == models.APPLY_STATUS_AGREE {
		if apply.Type == models.APPLY_TYPE_USER {
			friend, err := NewUserServ().UserIdToUserInfo(apply.Petitioner)
			if err != nil {
				return false, err
			}

			_, err = NewFriendServ().Adds(user, friend)
			if err != nil {
				return false, err
			}

		} else if apply.Type == models.APPLY_TYPE_GROUP { // 添加群
			// 查看群信息
			group, err := NewGroupServ().Info(id)
			if err != nil {
				return false, err
			}

			// 添加群成员
			_, err = NewGroupContactServ().Add(user, group)
			if err != nil {
				return false, err
			}
		}
	}

	// 更新
	_, err = a.DB.Where("id = ?", apply.Id).
		Cols("status", "UpdatedAt").
		Update(&updateData)
	if err != nil {
		return false, err
	}

	return true, nil
}
