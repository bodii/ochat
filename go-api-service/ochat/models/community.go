package models

import "time"

const (
	// 群状态：关闭
	COMMUNITY_STATUS_CLOSE int = iota
	// 群状态：开放
	COMMUNITY_STATUS_OPEN
)

const (
	// 群内大最成员数
	COMMUNITY_MAX_NUMBERS int = 500
)

// 群信息表
type Community struct {
	Id        int64     `xorm:"pk autoincr bigint not null comment('群表id')" form:"id" json:"id"`
	Name      string    `xorm:"varchar(60) index 'community_name' not null default '' comment('群名称')" form:"name" json:"name"`
	ManagerId int64     `xorm:"bigint index 'community_manager_id' not null default 0 comment('群主user_id')" form:"manager_id" json:"manager_id"`
	Icon      string    `xorm:"varchar(160) not null default 0 comment('群logo')" form:"icon" json:"icon"`
	Abstract  string    `xorm:"varchar(120) not null default 0 comment('简介')" form:"abstract" json:"abstract"`
	Type      int       `xorm:"tinyint not null default 0 comment('类型')" form:"type" json:"type"`
	Status    int       `xorm:"tinyint index 'community_status' not null default 1 comment('是否有效,1:有效;0:无效')" form:"status" json:"status"`
	CreatedAt time.Time `xorm:"datetime not null comment('创建时间')" form:"created_at" json:"created_at"`
	UpdatedAt time.Time `xorm:"datetime not null comment('更新时间')" form:"updated_at" json:"updated_at"`
}
