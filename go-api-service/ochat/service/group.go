package service

import (
	"errors"
	"fmt"
	"net/url"
	"ochat/bootstrap"
	"ochat/comm/funcs"
	"ochat/models"
	"time"

	"xorm.io/xorm"
)

type GroupService struct {
	DB *xorm.Engine
}

func NewGroupServ() *GroupService {
	return &GroupService{
		DB: bootstrap.DB_Engine,
	}
}

// 获取用户的所有有效的群列表信息
//
// params:
//   - userId: current user id
//
// return:
//   - groups [models.Group]: 群列表信息
//   - err [error]: 不成功时的出错信息
func (g *GroupService) UserList(userId int64) (groups []*models.Group, err error) {
	err = g.DB.Table("group").Alias("g").
		Join("left join", []string{"group_contact", "gc"}, "gc.group_id=g.id").
		Where("gc.user_id = ? and g.status = ?",
			userId, models.GROUP_STATUS_OPEN).
		In("gc.status", models.GROUP_CONTACT_STATUS_NORMAL, models.GROUP_CONTACT_STATUS_GROUP_TOP).
		Desc("gc.status").
		Asc("gc.id").
		Cols("g.*").
		Find(&groups)

	return groups, err
}

// 创建当前用户以及其它用户组成的新群
//
// params:
//   - master [models.User]: 发起者用户信息
//   - members [[]models.User]: 其它用户信息列表
//
// return:
//   - ok [bool]: 是否成功
//   - gInfo [models.Group]: 群信息
//   - err [error]: 不成功时的出错信息
func (g *GroupService) Create(master models.User, members ...models.User) (
	ok bool, gInfo models.Group, gs []*models.GroupContact, err error) {
	if master.Id == 0 {
		return false, gInfo, gs, errors.New("create group failure: creater user info is empty")
	}

	if len(members) == 0 {
		return false, gInfo, gs, errors.New("create group failure: group members info is empty")
	}

	group := models.Group{
		Name:      master.Nickname,
		Icon:      master.Avatar,
		Status:    models.GROUP_STATUS_OPEN,
		CreatedAt: time.Now(),
	}
	num, err := g.DB.InsertOne(&group)
	if err != nil || num == 0 {
		return false, group, gs, err
	}

	groupContacts := make([]*models.GroupContact, 1)
	groupContacts[0] = &models.GroupContact{
		UserId:       master.Id,
		GroupId:      group.Id,
		GroupAlias:   master.Nickname,
		Type:         models.GROUP_CONTACT_TYPE_MASTER,
		Nickname:     master.Nickname,
		NoticeStatus: models.GROUP_CONTACT_NOTICE_STATUS_NORMAL,
		Status:       models.GROUP_CONTACT_STATUS_NORMAL,
		CreatedAt:    time.Now(),
	}

	for _, u := range members {
		member := &models.GroupContact{
			UserId:       u.Id,
			GroupId:      group.Id,
			GroupAlias:   master.Nickname,
			Type:         models.GROUP_CONTACT_TYPE_MEMBER,
			Nickname:     u.Nickname,
			NoticeStatus: models.GROUP_CONTACT_NOTICE_STATUS_NORMAL,
			Status:       models.GROUP_CONTACT_STATUS_NORMAL,
			CreatedAt:    time.Now(),
		}

		groupContacts = append(groupContacts, member)
	}

	num, err = NewGroupContactServ().DB.Insert(groupContacts)
	if err != nil || num == 0 {
		return false, group, gs, err
	}

	return true, group, groupContacts, nil
}

// 查看群信息
//
// params:
//   - groupId [int64]: 群id
//
// return:
//   - group [models.Group]: 群信息
//   - err [error]: 不成功时的出错信息
func (g *GroupService) Info(groupId int64) (group models.Group, err error) {
	ok, err := g.DB.Where("id = ?", groupId).
		And("status = ?", models.GROUP_STATUS_OPEN).
		Get(&group)
	if err != nil || !ok {
		return group, errors.New("select failure")
	}

	return
}

// 创建群二维码
//
// params:
//   - group [*models.Group]: 群信息
//
// return:
//   - filename [string]: 生成的群二维码图片名
//   - err [error]: 不成功时的出错信息
func (g *GroupService) CreateQrCode(group *models.Group) (filename string, err error) {
	// 生成二维码
	qrCodeUrl := fmt.Sprintf("%s/group?group_id=%d", bootstrap.HTTP_HOST, group.Id)
	filename, err = funcs.QrCode(qrCodeUrl, "group_qrcode")
	if err != nil {
		return "", err
	}

	group.QrCode = funcs.GetImgUrl("group_qrcode", filename)
	group.UpdatedAt = time.Now()

	num, err := g.DB.Where("id =?", group.Id).Cols("qr_code", "updated_at").Update(group)
	if err != nil || num < 1 {
		return "", errors.New("update failure")
	}

	return
}

// 更新字段
//
// params:
//   - fields [url.Values]: 要更新字段传值
//   - userId [int64]: 更新者的用户id
//   - groupId [int64]: 要更新的群id
//
// return:
//   - group [models.Group]: 更新后的群信息
//   - err [error]: 不成功时的出错信息
func (g *GroupService) UpdateFields(fields url.Values, userId, groupId int64) (group models.Group, err error) {
	group, err = NewGroupServ().Info(groupId)
	if err != nil {
		return group, err
	}

	groupContact, err := NewGroupContactServ().Info(userId, group.Id)
	if err != nil {
		return group, err
	}

	if groupContact.Type < models.GROUP_CONTACT_TYPE_MANAGER {
		return group, errors.New("current user is group member, Unable to set")
	}

	canUpdateFields := []string{
		"name", "announcement", "about",
	}
	upFields := map[string]any{}
	for _, field := range canUpdateFields {
		if fields.Has(field) && fields.Get(field) != "" {
			upFields[field] = fields.Get(field)
		}
	}

	if len(upFields) == 0 {
		return group, errors.New("no update")
	}

	_, err = g.DB.Table("group").Where("id = ?", group.Id).Update(upFields)

	return
}
