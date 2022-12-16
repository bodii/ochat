drop table if exists message;
create table message (
    id bigint not null auto_increment primary key comment '消息id',
    mode tinyint not null default 1 comment '模式,1:单聊;2:群聊',
    sender_id int not null default 0 comment '发送用户id',
    receiver_id bigint not null default 0 comment '接收方id, [mode=1]:对方id,[mode=2]:群id',
    type tinyint not null default 1 comment '消息内容类型,1:文字;2:图片;3:表情;4:录音;5:名片;6:红包;7:音频文件;8:视频文件;9:接龙;',
    content mediumtext not null  comment '消息内容，或实体(如图片)链接',
    pic varchar(220) not null default '' comment '预览图片',
    url varchar(220) not null default '' comment '服务的url',
    about varchar(220) not null default '' comment '简单描述',
    amount decimal not null default 0.00 comment '金额',
    status tinyint not null default 0 comment '是否已读,-2:删除;-1:撤销;0:有效;1:已读',
    sender_status   tinyint  not null default 1 comment '发送者状态,-1:撤回;0:删除;1:成功（默认）;2:锁定（不可再撤回）;',
	receiver_status tinyint not null default 0 comment '接收者状态,-1:删除;0:未读;1:已读',
    created_at datetime(6) not null  comment '创建时间',
    sender_updated_at datetime(6) not null comment '撤销时间',
    receiver_updated_at datetime(6) not null comment '收接者已读时间',
    key message_receiver_id_mode(receiver_id,mode),
    key message_sender_id(sender_id),
    key message_sender_status(sender_status),
    key message_receiver_status(receiver_status)
) engine=innodb default charset=utf8mb4 collate=utf8mb4_unicode_ci comment '消息表';