package models

// 消息体核心model
type Message struct {
	Id      int64  `xorm:"pk autoincr bigint not null comment('消息id')" json:"id,omitempty" form:"id"`          // 消息ID
	Userid  int64  `xorm:"bigint not null default 0 comment('发送者用户id')" json:"userid,omitempty" form:"userid"` // 谁发的
	Cmd     int    `xorm:"int not null default 0 comment('操作类型,1:聊天;2:朋友圈')" json:"cmd,omitempty" form:"cmd"`  // 群聊还是私聊
	Dstid   int64  `xorm:"bigint not null default 0 comment('接收方(群)id')" json:"dstid,omitempty" form:"dstid"`  // 对端id/群id
	Media   int    `xorm:"int not null default 0 comment('消息内容样式, 1:')" json:"media,omitempty" form:"media"`   // 消息样式
	Content string `xorm:"mediumtext not null comment('消息内容')" json:"content,omitempty" form:"content"`        // 消息的内容
	Pic     string `xorm:"varchar(220) not null comment('预览图片')" json:"pic,omitempty" form:"pic"`              // 预览图片
	Url     string `xorm:"varchar(220) not null comment('服务的url')" json:"url,omitempty" form:"url"`            // 服务的url
	Memo    string `xorm:"varchar(220) not null comment('简单描述')" json:"memo,omitempty" form:"memo"`            // 简单描述
	Amount  int    `xorm:"int not null default 0 comment('和数字相关')" json:"amount,omitempty" form:"amount"`      // 和数字相关的
}
