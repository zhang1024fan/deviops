-- 检查表是否存在
USE `gin-api`;

-- 显示所有表
SHOW TABLES;

-- 检查jenkins_envs表是否存在
SHOW TABLES LIKE 'jenkins_envs';

-- 如果表不存在，手动创建jenkins_envs表
CREATE TABLE IF NOT EXISTS `jenkins_envs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `app_id` bigint unsigned NOT NULL COMMENT '应用ID',
  `env_name` varchar(50) NOT NULL COMMENT '环境名称',
  `jenkins_server_id` bigint unsigned DEFAULT NULL COMMENT 'Jenkins服务器ID',
  `job_name` varchar(255) DEFAULT NULL COMMENT 'Jenkins任务名称',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_jenkins_envs_app_id` (`app_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 检查quick_deployments表是否存在
SHOW TABLES LIKE 'quick_deployments';

-- 如果表不存在，手动创建quick_deployments表
CREATE TABLE IF NOT EXISTS `quick_deployments` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL COMMENT '发布标题',
  `business_group_id` bigint unsigned NOT NULL COMMENT '业务组ID',
  `business_dept_id` bigint unsigned NOT NULL COMMENT '业务部门ID',
  `description` text COMMENT '发布描述',
  `status` tinyint DEFAULT '1' COMMENT '发布状态',
  `task_count` int NOT NULL DEFAULT '0' COMMENT '任务数量',
  `creator_id` bigint unsigned NOT NULL COMMENT '创建人ID',
  `creator_name` varchar(100) DEFAULT NULL COMMENT '创建人姓名',
  `start_time` datetime(3) DEFAULT NULL COMMENT '开始发布时间',
  `end_time` datetime(3) DEFAULT NULL COMMENT '结束发布时间',
  `duration` int DEFAULT NULL COMMENT '发布耗时(秒)',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 检查quick_deployment_tasks表是否存在
SHOW TABLES LIKE 'quick_deployment_tasks';

-- 如果表不存在，手动创建quick_deployment_tasks表
CREATE TABLE IF NOT EXISTS `quick_deployment_tasks` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `deployment_id` bigint unsigned NOT NULL COMMENT '发布ID',
  `app_id` bigint unsigned NOT NULL COMMENT '应用ID',
  `app_name` varchar(255) DEFAULT NULL COMMENT '应用名称',
  `app_code` varchar(255) DEFAULT NULL COMMENT '应用编码',
  `environment` varchar(50) DEFAULT NULL COMMENT '环境名称',
  `jenkins_env_id` bigint unsigned NOT NULL COMMENT 'Jenkins环境配置ID',
  `jenkins_job_url` varchar(500) DEFAULT NULL COMMENT 'Jenkins任务URL',
  `build_number` int DEFAULT NULL COMMENT '构建编号',
  `status` tinyint DEFAULT '1' COMMENT '任务状态',
  `execute_order` int NOT NULL COMMENT '执行顺序',
  `start_time` datetime(3) DEFAULT NULL COMMENT '任务开始时间',
  `end_time` datetime(3) DEFAULT NULL COMMENT '任务结束时间',
  `duration` int DEFAULT NULL COMMENT '任务耗时(秒)',
  `error_message` text COMMENT '错误信息',
  `log_url` varchar(500) DEFAULT NULL COMMENT '日志URL',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_quick_deployment_tasks_deployment_id` (`deployment_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 最后再次显示所有表，确认创建成功
SHOW TABLES;