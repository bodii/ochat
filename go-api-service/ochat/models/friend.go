package models

import "time"

const (
	// 好友状态：屏蔽
	FRIEND_STATUS_HIDE int = iota - 1
	// 好友状态：黑名单
	FRIEND_STATUS_BLACKLIST
	// 好友状态：好友
	FRIEND_STATUS_FRIENDED
	// 好友状态：置顶
	FRIEND_STATUS_TOP
)

// 好友表
type Friend struct {
	Id        int64     `xorm:"pk autoincr bigint not null comment('好友表id')" form:"id" json:"id"`
	UserId    string    `xorm:"bigint index('friend_user_status') not null default 0 comment('用户id')" form:"userid" json:"userid"`
	FriendId  int64     `xorm:"bigint index('friend_id') not null default 0 comment('好友用户id')" form:"friend_id" json:"friend_id"`
	Status    int       `xorm:"tinyint index('friend_user_status') not null default 1 comment('当前状态, -1:屏蔽;0:黑名单;1:好友;2:置顶')" form:"status" json:"status"`
	CreatedAt time.Time `xorm:"datetime not null comment('创建时间')" form:"created_at" json:"created_at"`
	UpdatedAt time.Time `xorm:"datetime not null comment('更新时间')" form:"updated_at" json:"updated_at"`
}
