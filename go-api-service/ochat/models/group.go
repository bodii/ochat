package models

import "time"

const (
	// 群状态：关闭
	GROUP_STATUS_CLOSE int = iota
	// 群状态：开放
	GROUP_STATUS_OPEN
)

const (
	// 群内大最成员数
	GROUP_MAX_NUMBERS int = 500
)

// 群信息表
type Group struct {
	Id           int64     `xorm:"pk autoincr bigint not null comment('群表id')" form:"id" json:"id"`
	Name         string    `xorm:"varchar(60) index('group_name') not null default '' comment('群名称')" form:"name" json:"name"`
	Icon         string    `xorm:"varchar(160) not null default 0 comment('群logo')" form:"icon" json:"icon"`
	QrCode       string    `xorm:"varchar(160) not null default 0 comment('群二维码')" form:"qr_code" json:"qr_code"`
	Announcement string    `xorm:"varchar(220) not null default 0 comment('公告')" form:"announcement" json:"announcement"`
	About        string    `xorm:"varchar(120) not null default 0 comment('简介')" form:"about" json:"about"`
	Type         int       `xorm:"tinyint not null default 0 comment('类型')" form:"type" json:"type"`
	Status       int       `xorm:"tinyint index('group_status') not null default 1 comment('是否有效,1:有效;0:无效')" form:"status" json:"status"`
	CreatedAt    time.Time `xorm:"datetime(6) not null comment('创建时间')" form:"created_at" json:"created_at"`
	UpdatedAt    time.Time `xorm:"datetime(6) not null comment('更新时间')" form:"updated_at" json:"updated_at"`
}
