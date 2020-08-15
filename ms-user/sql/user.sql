create database bht_user;

CREATE TABLE IF NOT EXISTS  `user` (
	`user_id` INT(10) NOT NULL AUTO_INCREMENT,
	`user_name` VARCHAR(64) NOT NULL DEFAULT '' COLLATE 'utf8_general_ci',
	`password` VARCHAR(128) NOT NULL DEFAULT '' COLLATE 'utf8_general_ci',
	`birthday` VARCHAR(20) NOT NULL DEFAULT '' COLLATE 'utf8_general_ci',
	`sex` TINYINT(2) NOT NULL DEFAULT '1' COMMENT '1男2女',
	`avatar` VARCHAR(200) NOT NULL DEFAULT '' COMMENT '头像' COLLATE 'utf8_general_ci',
	`city` VARCHAR(10) NOT NULL DEFAULT '' COMMENT '城市' COLLATE 'utf8_general_ci',
	`district` VARCHAR(10) NOT NULL DEFAULT '' COMMENT '区域' COLLATE 'utf8_general_ci',
	`introduction` VARCHAR(1000) NOT NULL DEFAULT '' COMMENT '介绍' COLLATE 'utf8_general_ci',
	`role_id` TINYINT(2) NOT NULL DEFAULT '0' COMMENT '角色',
	`regist_ts` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

insert into user (`user_name`, `password`) values ('user', 'password')
