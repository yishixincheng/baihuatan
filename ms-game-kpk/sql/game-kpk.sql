create database bht_game;

-- 知识pk
CREATE TABLE `kpk_score` (
	`user_id` INT(10) NOT NULL,
	`score` INT(10) NOT NULL DEFAULT '0',
	`pet_id` INT(10) NOT NULL DEFAULT '0' COMMENT '宠物',
    `road_id` INT(10) NOT NULL DEFAULT '0' COMMENT '跑道',
	`update_ts` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
	UNIQUE INDEX `user_id` (`user_id`) USING BTREE
) COMMENT='知识pk积分统计表' COLLATE='utf8_general_ci' ENGINE=InnoDB;



-- 知识pk 积分记录表
CREATE TABLE IF NOT EXISTS  `kpk_record` (
    `id`          INT(10)      NOT NULL AUTO_INCREMENT,
    `user_id`     INT(10)      NOT NULL,
    `score`       INT(10)      NOT NULL DEFAULT 0 COMMENT '获取积分',
    `room_id`    VARCHAR(128)  NOT NULL DEFAULT '' COMMENT '房间id',
    `ranking`     TINYINT(2)   NOT NULL DEFAULT 0 COMMENT '名次',
    `question_count` INT(10)   NOT NULL DEFAULT 0 COMMENT '题目总数',
    `answer_count`   INT(10)   NOT NULL DEFAULT 0 COMMENT '答题总数',
    `answer_correct_count` INT(10) NOT NULL DEFAULT 0 COMMENT '回答正确总数',
    `update_ts`   timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP  COMMENT '更新时间',
    PRIMARY KEY (`id`),
    INDEX `idx_user_id` (`user_id`) USING BTREE
) COMMENT '知识pk积分统计表' COLLATE='utf8_general_ci' ENGINE=InnoDB;

-- 知识pk 题库
CREATE TABLE IF NOT EXISTS `kpk_question` (
    `id`          INT(10)      NOT NULL AUTO_INCREMENT,
    `title`       VARCHAR(300) NOT NULL DEFAULT '' COMMENT '题目标题',
    `option_1`    VARCHAR(512) NOT NULL DEFAULT '' COMMENT '选项1',
    `option_2`    VARCHAR(512) NOT NULL DEFAULT '' COMMENT '选项2',
    `option_3`    VARCHAR(512) NOT NULL DEFAULT '' COMMENT '选项3',
    `option_4`    VARCHAR(512) NOT NULL DEFAULT '' COMMENT '选项4',
    `right_option` VARCHAR(8)  NOT NULL DEFAULT '' COMMENT '1,2,3,4',
    `annotation`   VARCHAR(512) NOT NULL　DEFAULT '' COMMENT '注释',
    `author_id`    INT(10)     NOT NULL DEFAULT 0 COMMENT '作者',
    `cate_id`      INT(10)     NOT NULL DEFAULT 0 COMMENT '分类',
    `update_ts`   timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP  COMMENT '更新时间',
    PRIMARY KEY (`id`),
    INDEX `idx_cate_id` (`cate_id`) USING BTREE
) COMMENT '知识pk题库表' COLLATE='utf8_general_ci' ENGINE=InnoDB;

-- 知识pk 题目分类表
CREATE TABLE IF NOT EXISTS `kpk_question_category` (
    `id`          INT(10)      NOT NULL AUTO_INCREMENT,
    `parent_id`   INT(10)      NOT NULL DEFAULT 0 COMMENT '父类id',
    `name`    VARCHAR(512) NOT NULL DEFAULT '' COMMENT '名称',
    `update_ts`   timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP  COMMENT '更新时间',
    PRIMARY KEY (`id`)
) COMMENT '知识pk题库分类' COLLATE='utf8_general_ci' ENGINE=InnoDB;