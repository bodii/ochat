package service

import (
	"errors"
	"net/url"
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

func (g *GroupContactService) ChangeStatus(gc *models.GroupContact, status int) (
	err error) {

	num, err := g.DB.Where("id = ?", gc.Id).
		Cols("status", "updated_at").
		Update(map[string]any{
			"status":     status,
			"updated_at": time.Now(),
		})
	if err != nil || num == 0 {
		return err
	}

	gc.Status = status

	return nil
}

func (g *GroupContactService) ChangeType(gc *models.GroupContact, typeVal int) (
	err error) {

	num, err := g.DB.Where("id = ?", gc.Id).
		Cols("type", "updated_at").
		Update(map[string]any{
			"type":       typeVal,
			"updated_at": time.Now(),
		})
	if err != nil || num == 0 {
		return err
	}

	gc.Type = typeVal

	return nil
}

func (g *GroupContactService) UpdateFields(
	userId int64, groupId int64, fields url.Values) error {

	groupContact, err := g.Info(userId, groupId)
	if err != nil {
		return err
	}

	canUpFields := []string{
		"group_alias",
		"nickname",
		"notice_status",
	}

	// gcRef := reflect.ValueOf(groupContact)

	upFields := map[string]string{}
	for _, field := range canUpFields {
		if fields.Has(field) && fields.Get(field) != "" {
			upFields[field] = fields.Get(field)
		}
	}

	if len(upFields) == 0 {
		return errors.New("no update")
	}

	_, err = g.DB.Table("group_contact").Where("id = ?", groupContact.Id).Update(upFields)
	if err != nil {
		return err
	}

	return nil
}
