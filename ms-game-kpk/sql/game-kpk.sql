create database bht_game;

-- 知识pk
CREATE TABLE `kpk_score` (
	`user_id` INT(10) NOT NULL,
	`score` INT(10) NOT NULL DEFAULT '0',
	`update_ts` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
	UNIQUE INDEX `idx_user_id` (`user_id`) USING BTREE
) COMMENT='知识pk积分统计表' COLLATE='utf8_general_ci' ENGINE=InnoDB;



-- 知识pk 积分记录表
CREATE TABLE IF NOT EXISTS  `kpk_record` (
    `id`          INT(10)      NOT NULL AUTO_INCREMENT,
    `user_id`     INT(10)      NOT NULL,
    `score`       INT(10)      NOT NULL DEFAULT 0 COMMENT '获取积分',
    `house_id`    VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '房间id',
    `ranking`     TINYINT(2)   NOT NULL DEFAULT 0 COMMENT '名次',
    `question_count` INT(10)   NOT NULL DEFAULT 0 COMMENT '题目总数',
    `answer_count`   INT(10)   NOT NULL DEFAULT 0 COMMENT '答题总数',
    `answer_correct_count` INT(10) NOT NULL DEFAULT 0 COMMENT '回答正确总数',
    `update_ts`   timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP  COMMENT '更新时间',
    PRIMARY KEY (`id`),
    INDEX `idx_user_id` (`user_id`) USING BTREE
) COMMENT '知识pk积分统计表' COLLATE='utf8_general_ci' ENGINE=InnoDB;
