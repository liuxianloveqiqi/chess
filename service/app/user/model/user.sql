-- 用户表（User）
CREATE TABLE `user`
(
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT '用户ID',
    `password` varchar(50)     NOT NULL DEFAULT '' COMMENT '用户密码，MD5加密',
    `user_Nick`     varchar(100)    NOT NULL DEFAULT '' COMMENT '用户昵称',
    `User_Sex`      tinyint(1)      NOT NULL DEFAULT 0 COMMENT '用户性别：0男，1女，2保密',
    `user_Email`    varchar(255)    NOT NULL DEFAULT ''COMMENT '用户邮箱',
    `user_Phone`    varchar(11)     NOT NULL DEFAULT '' COMMENT '手机号',
    `create_time`  datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time`  datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `delete_time` datetime       COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `userPhone` (`user_Phone`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT ='用户表';