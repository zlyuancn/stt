CREATE TABLE `user_`
(
    `id`     int unsigned  NOT NULL AUTO_INCREMENT,
    `openid` varchar(40)   NOT NULL DEFAULT '',
    `extend` varchar(2048) NOT NULL DEFAULT '{}' COMMENT '扩展数据, json',
    `ctime`  datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `utime`  datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `openid_index` (`openid`)
) COMMENT ='user table';

CREATE index ctime_index on user_ (ctime);
