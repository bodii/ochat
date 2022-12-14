package models

import "time"

const (
	SEX_WOMAN   = 'w'
	SEX_MAN     = 'm'
	SEX_UNKNOWN = 'u'
)

// 接入鉴权/用户信息表
type User struct {
	Id         int64     `xorm:"pk autoincr bigint not null comment('用户id')" form:"id" json:"id"`
	Mobile     string    `xorm:"varchar(20) not null default '' comment('手机号')" form:"mobile" json:"mobile"`
	Nickname   string    `xorm:"varchar(20) not null default '' comment('用户id')" form:"nickname" json:"nickname"`
	Passwd     string    `xorm:"varchar(40) not null default '' comment('密码')" form:"passwd" json:"passwd"`
	Memo       string    `xorm:"varchar(250) not null default '' comment('简单描述')" form:"memo" json:"memo"`
	Token      string    `xorm:"varchar(50) not null default '' comment('用户SSO登录的token')" form:"token" json:"token"`
	Salt       string    `xorm:"varchar(40) not null default '' comment('盐')" form:"" json:""`
	Avatar     string    `xorm:"varchar(160) not null default '' comment('头像')" form:"avatar" json:"avatar"`
	Sex        string    `xorm:"varchar(2) not null default 'w' comment('性别,w:女;m:男')" form:"sex" json:"sex"`
	Online     int64     `xorm:"varchar(2) not null default 0 comment('在线时长')" form:"online" json:"online"`
	Stat       string    `xorm:"tinyint not null default 1 comment('状态是否可用, 1:可用;0:冻结')" form:"" json:""`
	Created_at time.Time `xorm:"datetime not null comment('创建时间')" form:"created_at" json:"created_at"`
}
