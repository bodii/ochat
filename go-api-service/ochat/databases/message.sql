create table message (
    id bigint not null auto_increment primary key comment '消息id',
    sender_id int not null default 0 comment '发送用户id',
    type tinyint not null default 1 comment '操作类型,1:聊天;2:朋友圈',
    cmd int not null default 1 comment '类型,1:一对一;2:一对多',
    dstid bigint not null default 1 comment '接收方(群)id',
    media int not null default 1 comment '消息内容样式, 1:',
    content mediumtext not null  comment '消息内容',
    pic varchar(220) not null default '' comment '预览图片',
    url varchar(220) not null default '' comment '服务的url',
    memo varchar(220) not null default '' comment '简单描述',
    amount int not null default 0 comment '和数字相关'
) engine=innodb default charset=utf8mb4 comment '消息体核心';