-- 会话表（记录两个用户之间的会话）
CREATE TABLE IF NOT EXISTS `conversations`
(
    `id`           BIGINT  NOT NULL AUTO_INCREMENT,
    `owner_uid`    BIGINT  NOT NULL COMMENT '会话拥有者',
    `target_id`    BIGINT  NOT NULL COMMENT '单聊：对方UID，群聊：群ID',
    `chat_type`    TINYINT NOT NULL DEFAULT 1 COMMENT '1单聊 2群聊',
    `last_msg_id`  BIGINT  NOT NULL DEFAULT 0 COMMENT '最新消息ID',
    `last_msg_seq` BIGINT  NOT NULL DEFAULT 0 COMMENT '最新消息序号',
    `unread_count` INT     NOT NULL DEFAULT 0 COMMENT '未读数',
    `updated_at`   BIGINT  NOT NULL DEFAULT 0,
    `created_at`   BIGINT  NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_owner_target` (`owner_uid`, `target_id`, `chat_type`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='会话表';

-- 消息索引表（MySQL 存索引，MongoDB 存正文）
CREATE TABLE IF NOT EXISTS `messages`
(
    `id`         BIGINT  NOT NULL AUTO_INCREMENT COMMENT '全局消息ID',
    `seq`        BIGINT  NOT NULL DEFAULT 0 COMMENT '客户端序号',
    `chat_type`  TINYINT NOT NULL COMMENT '1单聊 2群聊',
    `from_uid`   BIGINT  NOT NULL COMMENT '发送方UID',
    `to_id`      BIGINT  NOT NULL COMMENT '接收方ID',
    `msg_type`   INT     NOT NULL COMMENT '消息类型',
    `status`     TINYINT NOT NULL DEFAULT 1 COMMENT '1正常 2撤回',
    `send_time`  BIGINT  NOT NULL DEFAULT 0 COMMENT '发送时间（毫秒）',
    `created_at` BIGINT  NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    KEY `idx_from_to` (`from_uid`, `to_id`, `send_time`),
    KEY `idx_to_time` (`to_id`, `send_time`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='消息索引表';

-- 群组表
CREATE TABLE IF NOT EXISTS `groups`
(
    `id`           BIGINT       NOT NULL AUTO_INCREMENT COMMENT '群组ID',
    `name`         VARCHAR(64)  NOT NULL COMMENT '群名称',
    `avatar`       VARCHAR(256) NOT NULL DEFAULT '' COMMENT '群头像',
    `description`  VARCHAR(256) NOT NULL DEFAULT '' COMMENT '群描述',
    `owner_uid`    BIGINT       NOT NULL COMMENT '群主UID',
    `member_count` INT          NOT NULL DEFAULT 0 COMMENT '成员数量',
    `max_member`   INT          NOT NULL DEFAULT 500 COMMENT '最大成员数',
    `status`       TINYINT      NOT NULL DEFAULT 1 COMMENT '1正常 2解散',
    `created_at`   BIGINT       NOT NULL DEFAULT 0,
    `updated_at`   BIGINT       NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    KEY `idx_owner` (`owner_uid`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='群组表';

-- 群成员表
CREATE TABLE IF NOT EXISTS `group_members`
(
    `id`         BIGINT  NOT NULL AUTO_INCREMENT,
    `group_id`   BIGINT  NOT NULL COMMENT '群组ID',
    `uid`        BIGINT  NOT NULL COMMENT '成员UID',
    `role`       TINYINT NOT NULL DEFAULT 3 COMMENT '1群主 2管理员 3普通成员',
    `muted`      TINYINT NOT NULL DEFAULT 0 COMMENT '是否被禁言 0否 1是',
    `joined_at`  BIGINT  NOT NULL DEFAULT 0 COMMENT '加入时间',
    `created_at` BIGINT  NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_group_uid` (`group_id`, `uid`),
    KEY `idx_uid` (`uid`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='群成员表';