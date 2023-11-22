alter table user_
    change openid uid varchar(40) default '' not null;

alter table user_
    add avatar varchar(128) default '' not null comment '头像';
