drop table if exists contact;
create table contact (
    id bigint not null auto_increment primary key comment '联系人表id',
    user_id bigint not null default 0 comment '用户id',
    friend_user_id bigint not null default 0 comment '关联用户id',
    type tinyint not null default 1 comment '联系人类型,1:好友;2:群',
    nickname varchar(30) not null default '' comment '称呼，1:好友称呼;2:在群中的别称',
    abstract varchar(250) not null default '' comment '简介',
    status tinyint not null default 1 comment '是否有效,-1:拉黑(屏蔽);0:无效;1:有效;',
    created_at datetime(6) not null comment '创建时间',
    updated_at datetime(6) not null comment '更新时间',
    key contact_user_id_status(user_id, status, type),
    key contact_friend_user_id(friend_user_id)
) engine=innodb default charset=utf8mb4 collate=utf8mb4_unicode_ci comment '联系人表';