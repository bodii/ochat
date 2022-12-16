drop table if exists user;
create table user (
    id bigint not null auto_increment primary key comment '用户id,user_id',
    mobile varchar(20) unique not null default '' comment '手机号',
    username varchar(25) unique not null default '' comment '用户名',
    nickname varchar(30) not null default '' comment '称呼',
    password varchar(40) not null default '' comment '密码',
    about varchar(250) not null default '' comment '简单描述',
    token varchar(50) not null default '' comment '用户的token',
    salt varchar(40) not null default '' comment '盐',
    avatar varchar(160) not null default '' comment '头像',
    sex tinyint not null default 0 comment '性别,0:无;1:男;2:女;',
    online bigint not null default 0 comment '在线时长',
    status tinyint not null default 1 comment '状态是否可用, 1:可用;0:冻结',
    created_at datetime(6) not null comment '创建时间',
    updated_at datetime(6) not null comment '更新时间',
    key user_status(status) 
) engine=innodb default charset=utf8mb4 collate=utf8mb4_unicode_ci comment '接入鉴权/用户信息表';