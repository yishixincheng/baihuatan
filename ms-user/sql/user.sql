create database bht_user;

CREATE TABLE IS NOT EXSITS  `user` (
    `user_id`     INT(10)      NOT NULL AUTO_INCREMENT,
    `user_name`   VARCHAR(64)  NOT NULL DEFAULT '',
    `password`    VARCHAR(128) NOT NULL DEFAULT '',
    `birthday`    VARCHAR(20)  NOT NULL DEFAULT '',
    `sex`         TINYINT(2)   NOT NULL DEFAULT 1 COMMENT '1男2女',

    PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

insert into user (`user_name`, 'password') values ('user', 'password')
