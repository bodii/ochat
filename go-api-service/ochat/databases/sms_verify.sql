drop table if exists sms_verify;
create table sms_verify (
    id bigint not null auto_increment primary key comment '手机短信验证表id',
	phone varchar(20) not null default 0 comment '手机号',
	code int not null default 0 comment '验证码',
	created_at datetime(6) not null comment '创建时间',
	key sms_verify_user_phone(phone, code)
) engine=innodb default charset=utf8mb4 collate=utf8mb4_unicode_ci comment '手机短信验证表';