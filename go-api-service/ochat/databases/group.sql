drop table if exists `group`;
create table `group` (
    id bigint not null auto_increment primary key comment '群表id',
	name varchar(60) not null default '' comment '群名称',
	icon varchar(160) not null default 0 comment '群logo',
	qr_code  varchar(160) not null default 0 comment '群二维码',
	announcement varchar(220) not null default 0 comment '群公告', 
	about varchar(120) not null default 0 comment '简介',
	type tinyint not null default 0 comment '类型',
	status tinyint not null default 1 comment '是否有效,1:有效;0:无效',
	created_at datetime(6) not null comment '创建时间',
	updated_at datetime(6) not null comment '更新时间',
	key group_name(name),
    key group_status(status)
) engine=innodb default charset=utf8mb4 collate=utf8mb4_unicode_ci comment '群信息表';