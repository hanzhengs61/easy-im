CREATE DATABASE IF NOT EXISTS easy_im CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE easy_im;

CREATE TABLE IF NOT EXISTS `users`
(
    `id`         BIGINT       NOT NULL AUTO_INCREMENT COMMENT '用户ID',
    `username`   VARCHAR(32)  NOT NULL COMMENT '登录用户名（唯一）',
    `password`   VARCHAR(128) NOT NULL COMMENT 'bcrypt 加密后的密码',
    `nickname`   VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '昵称',
    `avatar`     VARCHAR(256) NOT NULL DEFAULT '' COMMENT '头像 URL',
    `status`     TINYINT      NOT NULL DEFAULT 1 COMMENT '状态：1正常 2禁用',
    `created_at` BIGINT       NOT NULL DEFAULT 0 COMMENT '创建时间（Unix 毫秒）',
    `updated_at` BIGINT       NOT NULL DEFAULT 0 COMMENT '更新时间（Unix 毫秒）',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_username` (`username`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='用户表';