构建分表sql文件的工具

# 背景

在业务中经常涉及到分表的创建和表变更, 每次都要写个脚本来生成sql语句, 这个工具能让使用者无需写脚本, 即可快速创建分表或表变更的sql文件.

# 安装

`go install github.com/zlyuancn/stt/stf@latest`

# 创建分表

首先准备一个sql文件, 它表示第一个分表的语句, 本工具会更根据这个语句将剩余的分表语句生成好并放在一个文件中.

比如准备一个用户表文件. `user_.sql`.

注意, 文件名必须与语句中的表名相同, column 也不要包含表名字符, 因为工具是根据文件名做全量替换的

```sql
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
```

通过以下命令创建分表

`stf -t user_.sql`

现在可以看到已经生成了一个文件 `user_.out.sql`

```sql
CREATE TABLE `user_0`
(
    `id`     int unsigned  NOT NULL AUTO_INCREMENT,
    `openid` varchar(40)   NOT NULL DEFAULT '',
    `extend` varchar(2048) NOT NULL DEFAULT '{}' COMMENT '扩展数据, json',
    `ctime`  datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `utime`  datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `openid_index` (`openid`)
) COMMENT ='user table';

CREATE index ctime_index on user_0 (ctime);


CREATE TABLE `user_1`
(
    `id`     int unsigned  NOT NULL AUTO_INCREMENT,
    `openid` varchar(40)   NOT NULL DEFAULT '',
    `extend` varchar(2048) NOT NULL DEFAULT '{}' COMMENT '扩展数据, json',
    `ctime`  datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `utime`  datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `openid_index` (`openid`)
) COMMENT ='user table';

CREATE index ctime_index on user_1 (ctime);
```

分表编号默认从 `0` 开始, 可以通过命令 `-s` 修改

# 创建变更

随着业务发展, 原有的表已经不满足需求了, 需要增加或者修改一些字段.
可以再创建一个 `user_.alter.sql` 文件(必须以`.alter.xxx`作为后缀), 然后写入要变更的语句. 如下:

```sql
alter table user_
    change openid uid varchar(40) default '' not null;

alter table user_
    add avatar varchar(128) default '' not null comment '头像';
```

通过以下命令创建分表变更语句

`stf -t user_.alter.sql`

现在可以看到已经生成了一个文件 `user_.alter.out.sql`

```sql
alter table user_0
    change openid uid varchar(40) default '' not null;

alter table user_0
    add avatar varchar(128) default '' not null comment '头像';

alter table user_1
    change openid uid varchar(40) default '' not null;

alter table user_1
    add avatar varchar(128) default '' not null comment '头像';
```

# 命令说明

```text
Usage of stf:
  -c int                                    
        分表数量 (default 1)                
  -i int                                    
        第一个分表的编号                    
  -o string                                 
        输出文件名后缀 (default ".out.sql") 
  -t string                                 
        要生成的表模板文件用半角逗号分隔
```
