package service

import (
	"errors"
	"ochat/bootstrap"
	"ochat/models"
	"time"

	"xorm.io/xorm"
)

type GroupContactService struct {
	DB *xorm.Engine
}

func NewGroupContactServ() *GroupContactService {
	return &GroupContactService{
		DB: bootstrap.DB_Engine,
	}
}

func (g *GroupContactService) Add(user models.User, group models.Group) (
	gs *models.GroupContact, err error) {

	gs = &models.GroupContact{
		UserId:       user.Id,
		GroupId:      group.Id,
		GroupAlias:   group.Name,
		Type:         models.GROUP_CONTACT_TYPE_MEMBER,
		Nickname:     user.Nickname,
		NoticeStatus: models.GROUP_CONTACT_NOTICE_STATUS_NORMAL,
		Status:       models.GROUP_CONTACT_NOTICE_STATUS_NORMAL,
		CreatedAt:    time.Now(),
	}
	num, err := g.DB.InsertOne(gs)
	if err != nil || num == 0 {
		return gs, errors.New("join group faiure")
	}

	return gs, nil
}

func (g *GroupContactService) Adds(group models.Group, members ...models.User) (
	gs []*models.GroupContact, err error) {

	for _, u := range members {
		gsInfo := &models.GroupContact{
			UserId:       u.Id,
			GroupId:      group.Id,
			GroupAlias:   group.Name,
			Type:         models.GROUP_CONTACT_TYPE_MEMBER,
			Nickname:     u.Nickname,
			NoticeStatus: models.GROUP_CONTACT_NOTICE_STATUS_NORMAL,
			Status:       models.GROUP_CONTACT_NOTICE_STATUS_NORMAL,
			CreatedAt:    time.Now(),
		}

		gs = append(gs, gsInfo)
	}

	num, err := g.DB.Insert(gs)
	if err != nil || num == 0 {
		return gs, errors.New("more member join group failure")
	}

	return gs, nil
}

func (g *GroupContactService) Info(userId int64, groupId int64) (gc models.GroupContact, err error) {
	_, err = g.DB.Where("user_id = ? and group_id = ?", userId, groupId).
		In("status", models.GROUP_CONTACT_STATUS_NORMAL, models.GROUP_CONTACT_STATUS_GROUP_TOP).
		Get(&gc)

	return gc, err
}

func (g *GroupContactService) TypeInfo(groupId int64, typeVal int) (gc models.GroupContact, err error) {
	_, err = g.DB.Where("group_id = ? and type = ?", groupId, typeVal).
		In("status", models.GROUP_CONTACT_STATUS_NORMAL, models.GROUP_CONTACT_STATUS_GROUP_TOP).
		Asc("id").
		Limit(0, 1).
		Get(&gc)

	return gc, err
}

func (g *GroupContactService) ChangeStatus(groupContactId int64, status int) (
	gc models.GroupContact, err error) {

	gc = models.GroupContact{}
	ok, err := g.DB.
		Where("id = ?", groupContactId).
		Get(&gc)
	if err != nil || !ok {
		return gc, err
	}

	num, err := g.DB.Where("id = ?", groupContactId).
		Cols("status", "updated_at").
		Update(map[string]any{
			"status":     status,
			"updated_at": time.Now(),
		})
	if err != nil || num == 0 {
		return gc, err
	}

	return gc, nil
}
