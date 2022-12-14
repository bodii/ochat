create table user (
    id bigint not null auto_increment primary key comment '用户id',
    mobile varchar(20) not null default '' comment '手机号',
    nickname varchar(20) not null default '' comment '称呼',
    passwd varchar(40) not null default '' comment '密码',
    memo varchar(250) not null default '' comment '简单描述',
    token varchar(50) not null default '' comment '用户SSO登录的token',
    salt varchar(40) not null default '' comment '盐',
    avatar varchar(160) not null default '' comment '头像',
    sex varchar(2) not null default 'w' comment '性别,w:女;m:男',
    online int not null default 0 comment '在线时长',
    stat tinyint not null default 1 comment '状态是否可用, 1:可用;0:冻结',
    created_at datetime not null comment '创建时间',
    key user_mobile(mobile),
    key user_stat(stat) 
) engine=innodb default charset=utf8mb4 comment '接入鉴权/用户信息表';