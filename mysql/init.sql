-- 创建用户表
CREATE TABLE `users`
(
	`id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
	`username`      VARCHAR(64)     NOT NULL UNIQUE,
	`email`         VARCHAR(128)             DEFAULT NULL UNIQUE,
	`password_hash` VARCHAR(256)    NOT NULL,
	`phone`         VARCHAR(20)              DEFAULT NULL UNIQUE,
	`status`        TINYINT         NOT NULL DEFAULT 1,
	`created_at`    DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`updated_at`    DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	`deleted_at`    DATETIME                 DEFAULT NULL,
	INDEX idx_username (`username`),
	INDEX idx_email (`email`),
	INDEX idx_phone (`phone`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;
