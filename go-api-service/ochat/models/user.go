package models

import "time"

const (
	// 用户性别：未知
	USER_SEX_UNKNOWN int = iota
	// 用户性别: 男
	USER_SEX_MAN
	// 用户性别: 女
	USER_SEX_WOMAN
)

const (
	// 用户状态: 失效（已不可用）
	USER_STATUS_INVALID int = iota
	// 用户状态: 有效（可用）
	USER_STATUS_VALID
)

// 接入鉴权/用户信息表
type User struct {
	UserId    int64     `xorm:"pk autoincr bigint not null comment('用户id,user_id')" form:"user_id" json:"user_id"`
	Mobile    string    `xorm:"varchar(20) unique not null default '' comment('手机号')" form:"mobile" json:"mobile"`
	Username  string    `xorm:"varchar(25) unique not null default '' comment('用户名')" form:"username" json:"username"`
	Nickname  string    `xorm:"varchar(30) not null default '' comment('用户昵称')" form:"nickname" json:"nickname"`
	QrCode    string    `xorm:"qr_code varchar(200) not null default '' comment('二维码')" form:"qr_code" json:"qr_code"`
	Password  string    `xorm:"varchar(40) not null default '' comment('密码')" form:"-" json:"-"`
	About     string    `xorm:"varchar(250) not null default '' comment('简单描述')" form:"about" json:"about"`
	Token     string    `xorm:"varchar(50) not null default '' comment('用户的token')" form:"token" json:"token"`
	Salt      string    `xorm:"varchar(40) not null default '' comment('盐')" form:"-" json:"-"`
	Avatar    string    `xorm:"varchar(160) not null default '' comment('头像')" form:"avatar" json:"avatar"`
	Sex       int       `xorm:"tinyint not null default 1 comment('性别,0:无;1:男;2:女')" form:"sex" json:"sex"`
	Birthday  time.Time `xorm:"date not null comment('生日')" form:"birthday" json:"birthday"`
	Online    int64     `xorm:"bigint not null default 0 comment('在线时长')" form:"online" json:"online"`
	Status    int       `xorm:"tinyint index('user_status') not null default 1 comment('状态是否可用, 1:可用;0:冻结')" form:"-" json:"-"`
	CreatedAt time.Time `xorm:"datetime(6) not null comment('创建时间')" form:"created_at" json:"created_at"`
	UpdatedAt time.Time `xorm:"datetime(6) not null comment('更新时间')" form:"updated_at" json:"updated_at"`
}
