DROP DATABASE IF EXISTS sqlopt_demo;
CREATE DATABASE sqlopt_demo;
USE sqlopt_demo;

-- 创建大表
CREATE TABLE users
(
	id         BIGINT PRIMARY KEY AUTO_INCREMENT,
	name       VARCHAR(255),
	email      VARCHAR(255),
	age        INT,
	created_at DATETIME,
	INDEX idx_email (email),
	INDEX idx_age_created (age, created_at)
);

-- 插入大量测试数据（可通过 Go 程序生成）
