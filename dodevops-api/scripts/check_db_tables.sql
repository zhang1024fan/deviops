-- 查看数据库gin-api中的所有表
USE `gin-api`;
SHOW TABLES;

-- 查看是否有jenkins相关的表
SHOW TABLES LIKE '%jenkins%';

-- 查看应用相关的表
SHOW TABLES LIKE 'app_%';

-- 如果app_jenkins_envs表不存在，创建它
CREATE TABLE IF NOT EXISTS `app_jenkins_envs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `app_id` bigint unsigned NOT NULL COMMENT '应用ID',
  `env_name` varchar(50) NOT NULL COMMENT '环境名称',
  `jenkins_server_id` bigint unsigned DEFAULT NULL COMMENT 'Jenkins服务器ID',
  `job_name` varchar(255) DEFAULT '' COMMENT 'Jenkins任务名称',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_app_jenkins_envs_app_id` (`app_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 检查创建的表结构
DESCRIBE `app_jenkins_envs`;