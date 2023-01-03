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
func (a *ApplyService) Add(userId, addUserId int64, comment string, addType int) (models.Apply, error) {
	addData := models.Apply{
		Petitioner: userId,
		Responder:  addUserId,
		Status:     0,
		Type:       1,
		FriendId:   addUserId,
		Comment:    comment,
		CreatedAt:  time.Now(),
	}

	if num, err := a.DB.InsertOne(&addData); err != nil || num <= 0 {
		errStr := "apply info insert database failure"
		return addData, errors.New(errStr)
	}

	return addData, nil
}

// 查看向当前用户发起申请的状态的好友列表
//
// [param]
//
//	userId: current user id
//	status: 应答者是否同意,-1:拒绝;0:未查看;1:已查看;2:同意
//	applyType: 申请类型,1:好友;2:群
//
// [return] user list or error
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
// [param]
//
//	id: apply id
//	status: 应答者是否同意,-1:拒绝;0:未查看;1:已查看;2:同意
//
// [return] whether to set up successfully
func (a *ApplyService) Set(id int64, status int) (bool, error) {
	updateData := models.Apply{
		Status:    status,
		UpdatedAt: time.Now(),
	}

	_, err := a.DB.Where("id = ?", id).
		Cols("status", "UpdatedAt").
		Update(&updateData)
	if err != nil {
		return false, err
	}

	return true, nil
}
