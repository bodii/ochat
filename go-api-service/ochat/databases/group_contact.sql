drop table if exists group_contact;
create table group_contact (
    id bigint not null auto_increment primary key comment '群联系人表id',
    user_id bigint not null default 0 comment '用户id',
    group_id bigint not null default 0 comment '群id',
    group_alias varchar(30) not null default '' comment '群别称',
    type tinyint not null default 1 comment '联系人类型,1:普通成员;2:管理员;3:群主',
    nickname varchar(30) not null default '' comment '联系人在群里的昵称',
    notice_status tinyint not null default 1 comment '通知状态,0:免打扰;1:正常;',
    status tinyint not null default 1 comment '状态,-1:被踢出;0:退出;1:正常;2:群置顶',
    created_at datetime(6) not null comment '创建时间',
    updated_at datetime(6) not null comment '更新时间',
    key group_contact_user_id_type(user_id, type),
    key group_contact_status(status),
    key group_contact_group_id(group_id)
) engine=innodb default charset=utf8mb4 collate=utf8mb4_unicode_ci comment '群联系人表';