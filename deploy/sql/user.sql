CREATE TABLE IF NOT EXISTS `user`
(
    `id`          bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '用户ID',
    `username`    varchar(64)     NOT NULL UNIQUE COMMENT '用户名',
    `password`    varchar(255)    NOT NULL COMMENT '密码（bcrypt 加密）',
    `nickname`    varchar(64)     NOT NULL DEFAULT '' COMMENT '昵称',
    `avatar`      varchar(255)    NOT NULL DEFAULT '' COMMENT '头像',
    `gender`      tinyint         NOT NULL DEFAULT 0 COMMENT '性别 0:未知 1:男 2:女',
    `status`      tinyint         NOT NULL DEFAULT 1 COMMENT '状态 1:正常 0:禁用',
    `create_time` timestamp       NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_username` (`username`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='用户表';