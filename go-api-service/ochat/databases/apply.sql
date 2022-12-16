drop table if exists apply;
create table apply (
    id bigint not null auto_increment primary key comment '申请好友表id',
	petitioner bigint not null default 0 comment '申请者userid',
	responder bigint not null default 0 comment '应答者userid(如果是群，则是管理者userid)',
	type tinyint not null default 1 comment '申请类型,1:好友;2:群',
	friend_id bigint not null default 0 comment '被申请好友的用户id(如果是群，则是群id)',
	comment varchar(250) not null default '' comment '申请时的留言',
	status tinyint not null default 0 comment '应答者是否同意,-1:拒绝;0:未查看;1:已查看;2:同意',
	created_at datetime(6) not null comment '创建时间',
	updated_at datetime(6) not null comment '更新时间',
	key apply_responder_status_type(responder, status, type),
    key community_petitioner(petitioner)
) engine=innodb default charset=utf8mb4 collate=utf8mb4_unicode_ci comment '申请好友(加群)表';