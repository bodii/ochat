package models

const (
	// 消息类型: 文字
	MESSAGE_TYPE_TEXT int = iota + 1
	// 消息类型: 图片
	MESSAGE_TYPE_IMAGE
	// 消息类型: 表情
	MESSAGE_TYPE_ENOJI
	// 消息类型: 音频
	MESSAGE_TYPE_RECORD
	// 消息类型: 名片
	MESSAGE_TYPE_NAME_CARD
	// 消息类型: 红包
	MESSAGE_TYPE_RED_ENVELOPES
	// 消息类型: 音频文件
	MESSAGE_TYPE_AUDIO_FILE
	// 消息类型: 视频文件
	MESSAGE_TYPE_VIDEO_FILE
	// 消息类型: 接龙
	MESSAGE_TYPE_SOLITAIRE
	// 消息类型: 代码
	MESSAGE_TYPE_CODE
)

const (
	// 消息发送者状态: 已撤回
	MESSAGE_SENDER_STATUS_WITHDRAWN int = iota - 1
	// 消息发送者状态: 已删除
	MESSAGE_SENDER_STATUS_DELETE
	// 消息发送者状态: 成功（默认）
	MESSAGE_SENDER_STATUS_DEFAULT
	// 消息状态: 锁定（不可再撤回）
	MESSAGE_STATUS_FORMAL
)

const (
	// 消息接收者状态: 已删除
	MESSAGE_RECEIVER_STATUS_DELETE int = iota - 1
	// 消息接收状态: 未读
	MESSAGE_RECEIVER_STATUS_UNREAD
	// 消息接收者状态: 已读
	MESSAGE_RECEIVER_STATUS_READED
)

const (
	// 聊天模式：单聊
	MESSAGE_MODE_SINGLE int = iota + 1
	// 聊天模式: 群聊
	MESSAGE_MODE_GROUP
)

// 消息体核心model
type Message struct {
	Id                int64  `xorm:"pk autoincr bigint not null comment('消息id')" json:"id,omitempty" form:"id"`
	SenderId          int64  `xorm:"bigint index('message_sender_id') not null default 0 comment('发送用户id')" json:"sender_id" form:"sender_id"`
	ReceiverId        int64  `xorm:"bigint index('message_receiver_id_mode') not null default 0 comment('接收方id, [mode=1]:对方id,[mode=2]:群id')" json:"receiver_id,omitempty" form:"receiver_id"`
	Mode              int    `xorm:"tinyint index('message_receiver_id_mode') not null default 1 comment('模式,1:单聊;2:群聊')" json:"mode" form:"mode"`
	Type              int    `xorm:"tinyint not null default 1 comment('消息内容类型,1:文字;2:图片;3:表情;4:录音;5:名片;6:红包;7:音频文件;8:视频文件;9:接龙;')" json:"type" form:"type"`
	Content           string `xorm:"mediumtext not null comment('消息内容')" json:"content,omitempty" form:"content"`
	Pic               string `xorm:"varchar(220) not null comment('预览图片')" json:"pic,omitempty" form:"pic"`
	Url               string `xorm:"varchar(220) not null comment('服务的url')" json:"url,omitempty" form:"url"`
	About             string `xorm:"varchar(220) not null default '' comment '简单描述'" json:"about,omitempty" form:"about"`
	Amount            int    `xorm:"int not null default 0 comment('金额')" json:"amount,omitempty" form:"amount"`
	SenderStatus      int    `xorm:"tinyint index('message_sender_status') not null default 1 comment('发送者状态,-1:撤回;0:删除;1:成功（默认）;2:锁定（不可再撤回）;')" json:"sender_status,omitempty" form:"sender_status"`
	ReceiverStatus    int    `xorm:"tinyint index('message_receiver_status') not null default 0 comment('接收者状态,-1:删除;0:未读;1:已读')" json:"receiver_status,omitempty" form:"receiver_status"`
	CreatedAt         int    `xorm:"datetime(6) not null comment('创建时间')" json:"created_at,omitempty" form:"created_at"`
	SenderUpdatedAt   int    `xorm:"datetime(6) not null comment('发送者更新时间')" json:"sender_updated_at,omitempty" form:"sender_updated_at"`
	ReceiverUpdatedAt int    `xorm:"datetime(6) not null comment('接收者更新时间')" json:"receiver_updated_at,omitempty" form:"receiver_updated_at"`
}
