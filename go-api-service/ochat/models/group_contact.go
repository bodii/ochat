package models

import "time"

const (
	// 群联系人类型：普通成员
	GROUP_CONTACT_TYPE_MEMBER int = iota
	// 群联系人类型：管理员
	GROUP_CONTACT_TYPE_MANAGER
	// 群联系人类型：群主
	GROUP_CONTACT_TYPE_MASTER
)

const (
	// 群联系人状态: 被踢出
	GROUP_CONTACT_STATUS_KICK_OUT int = iota - 1
	// 群联系人状态: 退出
	GROUP_CONTACT_STATUS_EXIT
	// 群联系人状态: 正常
	GROUP_CONTACT_STATUS_NORMAL
	// 群联系人状态: 群置顶
	GROUP_CONTACT_STATUS_GROUP_TOP
)

const (
	// 群联系人消息状态设置: 免打扰
	GROUP_CONTACT_NOTICE_STATUS_NOT_DISTURB int = iota
	// 群联系人消息状态设置: 正常
	GROUP_CONTACT_NOTICE_STATUS_NORMAL
)

// 群联系人表
type GroupContact struct {
	Id           int64     `xorm:"pk autoincr bigint not null comment('群联系人表id')" from:"id" json:"id"`
	UserId       int64     `xorm:"bigint index('group_contact_user_id_type') not null comment('用户id')" from:"user_id" json:"user_id"`
	GroupId      int64     `xorm:"bigint index('group_contact_group_id') not null comment('群id')" from:"group_id" json:"group_id"`
	GroupAlias   string    `xorm:"varchar(30) not null default '' comment('群别称')" form:"group_alias" json:"group_alias"`
	Type         int       `xorm:"tinyint index('group_contact_user_id_type') not null default 1 comment('联系人类型,1:普通成员;2:管理员;3:群主')" form:"type" json:"type"`
	Nickname     string    `xorm:"varchar(30) not null default '' comment('联系人在群里的昵称')" form:"nickname" json:"nickname"`
	NoticeStatus int       `xorm:"tinyint not null default 1 comment('通知状态,0:免打扰;1:正常;')" form:"notice_status" json:"notice_status"`
	Status       int       `xorm:"tinyint index('group_contact_status') not null default 1 comment('状态,-1:被踢出;0:退出;1:正常;2:群置顶')" form:"status" json:"status"`
	CreatedAt    time.Time `xorm:"datetime(6) not null comment('创建时间')" form:"created_at" json:"created_at"`
	UpdatedAt    time.Time `xorm:"datetime(6) not null comment('更新时间')" form:"updated_at" json:"updated_at"`
}
