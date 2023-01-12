package models

import "time"

const (
	// 申请的类型: 好友
	APPLY_TYPE_USER int = iota + 1
	// 申请的类型: 群
	APPLY_TYPE_GROUP
)

const (
	// 申请后应答者的状态: 拒绝
	APPLY_STATUS_REFUSE int = iota - 1
	// 申请后应答者的状态: 未读
	APPLY_STATUS_UNREAD
	// 申请后应答者的状态: 已读
	APPLY_STATUS_READ
	// 申请后应答者的状态: 同意
	APPLY_STATUS_AGREE
)

// 申请好友(加群)表
type Apply struct {
	Id         int64     `xorm:"pk autoincr bigint not null comment('申请好友表id')" form:"id" json:"id"`
	Petitioner int64     `xorm:"bigint index('petitioner') not null default 0 comment('申请者userid')" form:"petitioner" json:"petitioner"`
	Responder  int64     `xorm:"bigint index('apply_responder_status_type') not null default 0 comment('应答者userid(如果是群，则是群id)')" form:"responder" json:"responder"`
	Status     int       `xorm:"tinyint index('apply_responder_status_type') not null default 1 comment('应答者是否同意,-1:拒绝;0:未查看;1:已查看;2:同意')" form:"status" json:"status"`
	Type       int       `xorm:"tinyint index('apply_responder_status_type') not null default 0 comment('申请类型,1:好友;2:群')" form:"type" json:"type"`
	Comment    string    `xorm:"varchar(250) not null default '' comment('申请时的留言')" form:"comment" json:"comment"`
	CreatedAt  time.Time `xorm:"datetime(6) not null comment('创建时间')" form:"created_at" json:"created_at"`
	UpdatedAt  time.Time `xorm:"datetime(6) not null comment('更新时间')" form:"updated_at" json:"updated_at"`
}
