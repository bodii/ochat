package models

import "time"

const (
	// 联系人类型: 单聊好友
	CONTCAT_TYPE_USER int = iota + 1
	// 联系人类型: 群聊的群
	CONTCAT_TYPE_COMMUNITY
)

const (
	// 联系人状态: 无效
	CONTACT_STATUS_INVALID int = iota
	// 联系人状态: 生效
	CONTACT_STATUS_VALID
)

type Contact struct {
	Id           int64     `xorm:"pk autoincr bigint not null comment('联系人表id')" form:"id" json:"id"`
	UserId       int64     `xorm:"bigint not null default 0 comment('用户id')" form:"user_id" json:"user_id"`
	FriendUserId int64     `xorm:"bigint not null default 0 comment('关联用户id')" form:"friend_user_id" json:"firend_user_id"`
	Type         int       `xorm:"tinyint not null default 0 comment('联系人类型')" form:"type" json:"type"`
	Nickname     string    `xorm:"varchar(30) not null default '' comment('称呼, 1:好友称呼;2:在群中的别称')" form:"nickname" json:"nickname"`
	Abstract     string    `xorm:"varchar(120) not null default 0 comment('简介')" form:"abstract" json:"abstract"`
	Status       int       `xorm:"tinyint not null default 1 comment('是否有效,1:有效;0:无效')" form:"status" json:"status"`
	CreatedAt    time.Time `xorm:"datetime(6) not null comment('创建时间')" form:"created_at" json:"created_at"`
	UpdatedAt    time.Time `xorm:"datetime(6) not null comment('更新时间')" form:"updated_at" json:"updated_at"`
}
