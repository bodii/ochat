drop table if exists friend;
create table friend (
    id bigint not null auto_increment primary key comment '好友表id',
	userid bigint not null default 0 comment '用户id',
	friend_id bigint not null default 0 comment '好友用户id',
	friend_alias varchar(25) not null default '' comment '好友别名',
	alias_prefix varchar(2) not null default '' comment '好友别名前缀(e.g 刘:L)',
	about varchar(250) not null default '' comment '描述',
	status tinyint not null default 0 comment '当前状态, -1:屏蔽;0:黑名单;1:好友;2:置顶',
	created_at datetime(6) not null comment '创建时间',
	updated_at datetime(6) not null comment '更新时间',
	key friend_user_status(userid, status),
    key friend_friend_id(friend_id),
	key friend_alias_prefix(alias_prefix)
) engine=innodb default charset=utf8mb4 collate=utf8mb4_unicode_ci comment '好友表';