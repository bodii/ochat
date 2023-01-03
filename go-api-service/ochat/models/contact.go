package models

import "time"

const (
	// 联系人类型: 单聊好友
	CONTCAT_TYPE_USER int = iota + 1
	// 联系人类型: 群聊的群
	CONTCAT_TYPE_GROUP
)

const (
	// 联系人状态: 拉黑(屏蔽)
	CONTACT_STATUS_SHIELD int = iota - 1
	// 联系人状态: 无效
	CONTACT_STATUS_INVALID
	// 联系人状态: 生效
	CONTACT_STATUS_VALID
)

// 联系人表
type Contact struct {
	Id           int64     `xorm:"pk autoincr bigint not null comment('联系人表id')" form:"id" json:"id"`
	UserId       int64     `xorm:"bigint index('contact_user_id_status') not null default 0 comment('用户id')" form:"user_id" json:"user_id"`
	FriendUserId int64     `xorm:"bigint index('contact_friend_user_id') not null default 0 comment('关联用户id')" form:"friend_user_id" json:"firend_user_id"`
	Nickname     string    `xorm:"varchar(30) not null default '' comment('称呼, 1:好友称呼;2:在群中的别称')" form:"nickname" json:"nickname"`
	Abstract     string    `xorm:"varchar(120) not null default 0 comment('简介')" form:"abstract" json:"abstract"`
	Status       int       `xorm:"tinyint index('contact_user_id_status') not null default 1 comment('是否有效,-1:拉黑(屏蔽);0:无效;1:有效;')" form:"status" json:"status"`
	Type         int       `xorm:"tinyint index('contact_user_id_status') not null default 0 comment('联系人类型')" form:"type" json:"type"`
	CreatedAt    time.Time `xorm:"datetime(6) not null comment('创建时间')" form:"created_at" json:"created_at"`
	UpdatedAt    time.Time `xorm:"datetime(6) not null comment('更新时间')" form:"updated_at" json:"updated_at"`
}
