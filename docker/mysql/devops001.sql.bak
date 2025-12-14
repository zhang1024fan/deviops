/*
 Navicat Premium Dump SQL

 Source Server         : 8.130.14.34
 Source Server Type    : MySQL
 Source Server Version : 80043 (8.0.43-0ubuntu0.24.04.1)
 Source Host           : 8.130.14.34:3306
 Source Schema         : gin-api

 Target Server Type    : MySQL
 Target Server Version : 80043 (8.0.43-0ubuntu0.24.04.1)
 File Encoding         : 65001

 Date: 16/10/2025 11:01:05
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for app_application
-- ----------------------------
DROP TABLE IF EXISTS `app_application`;
CREATE TABLE `app_application` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `code` varchar(255) NOT NULL,
  `business_group_id` bigint unsigned NOT NULL,
  `business_dept_id` bigint unsigned NOT NULL,
  `description` text,
  `repo_url` varchar(500) DEFAULT NULL,
  `dev_owners` json DEFAULT NULL,
  `test_owners` json DEFAULT NULL,
  `ops_owners` json DEFAULT NULL,
  `programming_lang` varchar(100) DEFAULT NULL,
  `start_command` text,
  `stop_command` text,
  `health_api` varchar(500) DEFAULT NULL,
  `domains` json DEFAULT NULL,
  `hosts` json DEFAULT NULL,
  `databases` json DEFAULT NULL,
  `other_res` json DEFAULT NULL,
  `status` tinyint DEFAULT '1',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_app_application_code` (`code`),
  KEY `idx_app_application_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of app_application
-- ----------------------------
BEGIN;
INSERT INTO `app_application` (`id`, `name`, `code`, `business_group_id`, `business_dept_id`, `description`, `repo_url`, `dev_owners`, `test_owners`, `ops_owners`, `programming_lang`, `start_command`, `stop_command`, `health_api`, `domains`, `hosts`, `databases`, `other_res`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (4, 's3-api', 's3-api', 12, 2, 's3-api', 'git@gitee.com:zhang_fan1024/ansible-playbook.git', '[101, 99]', '[99, 101]', '[89]', 'Java', 'systemctl start api', 'systemctl stop api', '/api/status', NULL, NULL, NULL, '{}', 2, '2025-09-16 13:35:12.995', '2025-09-16 18:09:09.739', NULL);
INSERT INTO `app_application` (`id`, `name`, `code`, `business_group_id`, `business_dept_id`, `description`, `repo_url`, `dev_owners`, `test_owners`, `ops_owners`, `programming_lang`, `start_command`, `stop_command`, `health_api`, `domains`, `hosts`, `databases`, `other_res`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (5, 's3-web', 's3-web', 18, 2, 's3-web', 'git@gitee.com:zhang_fan1024/ansible-playbook.git', '[101]', '[99]', '[89]', 'Node.js', 'systemctl start nginx', 'systemctl stop nginx', '/api/status', NULL, NULL, NULL, '{}', 2, '2025-09-16 15:02:40.356', '2025-09-16 18:08:57.912', NULL);
INSERT INTO `app_application` (`id`, `name`, `code`, `business_group_id`, `business_dept_id`, `description`, `repo_url`, `dev_owners`, `test_owners`, `ops_owners`, `programming_lang`, `start_command`, `stop_command`, `health_api`, `domains`, `hosts`, `databases`, `other_res`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (6, 's4-api', 's4-api', 14, 13, 's4-api', '', '[101]', '[99]', '[89]', 'Go', '', '', '', NULL, NULL, NULL, '{}', 2, '2025-09-16 16:11:24.562', '2025-09-16 18:08:09.769', NULL);
INSERT INTO `app_application` (`id`, `name`, `code`, `business_group_id`, `business_dept_id`, `description`, `repo_url`, `dev_owners`, `test_owners`, `ops_owners`, `programming_lang`, `start_command`, `stop_command`, `health_api`, `domains`, `hosts`, `databases`, `other_res`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (7, 's4-web', 's4-web', 14, 13, 's4-web', '', '[101]', '[99]', '[89]', 'Node.js', '', '', '', NULL, NULL, NULL, '{}', 2, '2025-09-16 16:11:50.838', '2025-09-16 18:07:54.723', NULL);
INSERT INTO `app_application` (`id`, `name`, `code`, `business_group_id`, `business_dept_id`, `description`, `repo_url`, `dev_owners`, `test_owners`, `ops_owners`, `programming_lang`, `start_command`, `stop_command`, `health_api`, `domains`, `hosts`, `databases`, `other_res`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (8, 's5-web', 's5-web', 4, 12, 's5-web', '', '[101]', '[99]', '[98]', 'Python', '', '', '', NULL, NULL, NULL, '{}', 2, '2025-09-16 16:12:42.467', '2025-09-22 11:48:39.215', NULL);
INSERT INTO `app_application` (`id`, `name`, `code`, `business_group_id`, `business_dept_id`, `description`, `repo_url`, `dev_owners`, `test_owners`, `ops_owners`, `programming_lang`, `start_command`, `stop_command`, `health_api`, `domains`, `hosts`, `databases`, `other_res`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (9, 's5-api', 's5-api', 9, 12, '', '', '[99]', '[89]', '[98]', 'PHP', '', '', '', NULL, NULL, NULL, '{}', 2, '2025-09-16 16:13:34.939', '2025-09-18 10:30:14.069', NULL);
INSERT INTO `app_application` (`id`, `name`, `code`, `business_group_id`, `business_dept_id`, `description`, `repo_url`, `dev_owners`, `test_owners`, `ops_owners`, `programming_lang`, `start_command`, `stop_command`, `health_api`, `domains`, `hosts`, `databases`, `other_res`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (10, 's4-web', 's4-web-001', 14, 13, 's4-web1111', '', '[101]', '[99]', '[89]', 'Node.js', '', '', '', NULL, NULL, NULL, '{}', 2, '2025-09-16 18:07:55.683', '2025-10-01 15:18:35.353', NULL);
COMMIT;

-- ----------------------------
-- Table structure for app_jenkins_env
-- ----------------------------
DROP TABLE IF EXISTS `app_jenkins_env`;
CREATE TABLE `app_jenkins_env` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `app_id` bigint unsigned NOT NULL,
  `env_name` varchar(50) NOT NULL,
  `jenkins_server_id` bigint unsigned DEFAULT NULL,
  `job_name` varchar(255) DEFAULT '',
  `job_url` varchar(500) DEFAULT NULL,
  `build_params` json DEFAULT NULL,
  `deploy_config` json DEFAULT NULL,
  `notification` json DEFAULT NULL,
  `is_active` tinyint DEFAULT '1',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_app_jenkins_env_app_id` (`app_id`),
  KEY `idx_app_jenkins_env_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_app_application_jenkins_envs` FOREIGN KEY (`app_id`) REFERENCES `app_application` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=26 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of app_jenkins_env
-- ----------------------------
BEGIN;
INSERT INTO `app_jenkins_env` (`id`, `app_id`, `env_name`, `jenkins_server_id`, `job_name`, `job_url`, `build_params`, `deploy_config`, `notification`, `is_active`, `created_at`, `updated_at`, `deleted_at`) VALUES (4, 4, 'prod', 10, 's3-api', NULL, NULL, NULL, NULL, 1, '2025-09-16 13:35:13.203', '2025-09-16 14:56:29.213', NULL);
INSERT INTO `app_jenkins_env` (`id`, `app_id`, `env_name`, `jenkins_server_id`, `job_name`, `job_url`, `build_params`, `deploy_config`, `notification`, `is_active`, `created_at`, `updated_at`, `deleted_at`) VALUES (5, 4, 'test', 11, 's3-api', NULL, NULL, NULL, NULL, 1, '2025-09-16 13:35:13.401', '2025-09-16 14:58:16.597', NULL);
INSERT INTO `app_jenkins_env` (`id`, `app_id`, `env_name`, `jenkins_server_id`, `job_name`, `job_url`, `build_params`, `deploy_config`, `notification`, `is_active`, `created_at`, `updated_at`, `deleted_at`) VALUES (6, 4, 'dev', 9, 's3-api', NULL, NULL, NULL, NULL, 1, '2025-09-16 13:35:13.601', '2025-09-16 14:54:49.866', NULL);
INSERT INTO `app_jenkins_env` (`id`, `app_id`, `env_name`, `jenkins_server_id`, `job_name`, `job_url`, `build_params`, `deploy_config`, `notification`, `is_active`, `created_at`, `updated_at`, `deleted_at`) VALUES (7, 5, 'prod', 10, 's3-web', NULL, NULL, NULL, NULL, 1, '2025-09-16 15:02:40.586', '2025-09-16 15:27:02.176', NULL);
INSERT INTO `app_jenkins_env` (`id`, `app_id`, `env_name`, `jenkins_server_id`, `job_name`, `job_url`, `build_params`, `deploy_config`, `notification`, `is_active`, `created_at`, `updated_at`, `deleted_at`) VALUES (8, 5, 'test', 9, 's3-web', NULL, NULL, NULL, NULL, 1, '2025-09-16 15:02:40.810', '2025-09-16 15:27:12.567', NULL);
INSERT INTO `app_jenkins_env` (`id`, `app_id`, `env_name`, `jenkins_server_id`, `job_name`, `job_url`, `build_params`, `deploy_config`, `notification`, `is_active`, `created_at`, `updated_at`, `deleted_at`) VALUES (9, 5, 'dev', 9, 's3-web', NULL, NULL, NULL, NULL, 1, '2025-09-16 15:02:41.025', '2025-09-16 15:26:43.689', NULL);
INSERT INTO `app_jenkins_env` (`id`, `app_id`, `env_name`, `jenkins_server_id`, `job_name`, `job_url`, `build_params`, `deploy_config`, `notification`, `is_active`, `created_at`, `updated_at`, `deleted_at`) VALUES (10, 5, 'test1', 9, 's3-web1', NULL, NULL, NULL, NULL, 1, '2025-09-16 15:36:43.585', '2025-09-16 15:36:43.585', NULL);
INSERT INTO `app_jenkins_env` (`id`, `app_id`, `env_name`, `jenkins_server_id`, `job_name`, `job_url`, `build_params`, `deploy_config`, `notification`, `is_active`, `created_at`, `updated_at`, `deleted_at`) VALUES (11, 6, 'prod', 9, 's4-api', NULL, NULL, NULL, NULL, 1, '2025-09-16 16:11:24.773', '2025-09-16 16:14:44.414', NULL);
INSERT INTO `app_jenkins_env` (`id`, `app_id`, `env_name`, `jenkins_server_id`, `job_name`, `job_url`, `build_params`, `deploy_config`, `notification`, `is_active`, `created_at`, `updated_at`, `deleted_at`) VALUES (12, 6, 'test', 10, 's4-api', NULL, NULL, NULL, NULL, 1, '2025-09-16 16:11:24.994', '2025-09-16 16:14:21.179', NULL);
INSERT INTO `app_jenkins_env` (`id`, `app_id`, `env_name`, `jenkins_server_id`, `job_name`, `job_url`, `build_params`, `deploy_config`, `notification`, `is_active`, `created_at`, `updated_at`, `deleted_at`) VALUES (13, 6, 'dev', 9, 's4-api', NULL, NULL, NULL, NULL, 1, '2025-09-16 16:11:25.216', '2025-09-16 16:14:06.669', NULL);
INSERT INTO `app_jenkins_env` (`id`, `app_id`, `env_name`, `jenkins_server_id`, `job_name`, `job_url`, `build_params`, `deploy_config`, `notification`, `is_active`, `created_at`, `updated_at`, `deleted_at`) VALUES (14, 7, 'prod', NULL, '', NULL, NULL, NULL, NULL, 1, '2025-09-16 16:11:51.051', '2025-09-16 16:11:51.051', NULL);
INSERT INTO `app_jenkins_env` (`id`, `app_id`, `env_name`, `jenkins_server_id`, `job_name`, `job_url`, `build_params`, `deploy_config`, `notification`, `is_active`, `created_at`, `updated_at`, `deleted_at`) VALUES (15, 7, 'test', 10, 's4-web', NULL, NULL, NULL, NULL, 1, '2025-09-16 16:11:51.263', '2025-09-16 16:15:19.182', NULL);
INSERT INTO `app_jenkins_env` (`id`, `app_id`, `env_name`, `jenkins_server_id`, `job_name`, `job_url`, `build_params`, `deploy_config`, `notification`, `is_active`, `created_at`, `updated_at`, `deleted_at`) VALUES (16, 7, 'dev', 9, 's4-web', NULL, NULL, NULL, NULL, 1, '2025-09-16 16:11:51.476', '2025-09-16 16:15:07.702', NULL);
INSERT INTO `app_jenkins_env` (`id`, `app_id`, `env_name`, `jenkins_server_id`, `job_name`, `job_url`, `build_params`, `deploy_config`, `notification`, `is_active`, `created_at`, `updated_at`, `deleted_at`) VALUES (17, 8, 'prod', NULL, '', NULL, NULL, NULL, NULL, 1, '2025-09-16 16:12:42.683', '2025-09-16 16:12:42.683', NULL);
INSERT INTO `app_jenkins_env` (`id`, `app_id`, `env_name`, `jenkins_server_id`, `job_name`, `job_url`, `build_params`, `deploy_config`, `notification`, `is_active`, `created_at`, `updated_at`, `deleted_at`) VALUES (18, 8, 'test', NULL, '', NULL, NULL, NULL, NULL, 1, '2025-09-16 16:12:42.902', '2025-09-16 16:12:42.902', NULL);
INSERT INTO `app_jenkins_env` (`id`, `app_id`, `env_name`, `jenkins_server_id`, `job_name`, `job_url`, `build_params`, `deploy_config`, `notification`, `is_active`, `created_at`, `updated_at`, `deleted_at`) VALUES (19, 8, 'dev', 11, 's5-web', NULL, NULL, NULL, NULL, 1, '2025-09-16 16:12:43.120', '2025-09-22 11:48:38.869', NULL);
INSERT INTO `app_jenkins_env` (`id`, `app_id`, `env_name`, `jenkins_server_id`, `job_name`, `job_url`, `build_params`, `deploy_config`, `notification`, `is_active`, `created_at`, `updated_at`, `deleted_at`) VALUES (20, 9, 'prod', NULL, '', NULL, NULL, NULL, NULL, 1, '2025-09-16 16:13:35.146', '2025-09-16 16:13:35.146', NULL);
INSERT INTO `app_jenkins_env` (`id`, `app_id`, `env_name`, `jenkins_server_id`, `job_name`, `job_url`, `build_params`, `deploy_config`, `notification`, `is_active`, `created_at`, `updated_at`, `deleted_at`) VALUES (21, 9, 'test', NULL, '', NULL, NULL, NULL, NULL, 1, '2025-09-16 16:13:35.361', '2025-09-16 16:13:35.361', NULL);
INSERT INTO `app_jenkins_env` (`id`, `app_id`, `env_name`, `jenkins_server_id`, `job_name`, `job_url`, `build_params`, `deploy_config`, `notification`, `is_active`, `created_at`, `updated_at`, `deleted_at`) VALUES (22, 9, 'dev', 10, 's4-api', NULL, NULL, NULL, NULL, 1, '2025-09-16 16:13:35.575', '2025-09-18 10:30:13.709', NULL);
INSERT INTO `app_jenkins_env` (`id`, `app_id`, `env_name`, `jenkins_server_id`, `job_name`, `job_url`, `build_params`, `deploy_config`, `notification`, `is_active`, `created_at`, `updated_at`, `deleted_at`) VALUES (23, 10, 'prod', NULL, '', NULL, NULL, NULL, NULL, 1, '2025-09-16 18:07:55.904', '2025-09-16 18:07:55.904', NULL);
INSERT INTO `app_jenkins_env` (`id`, `app_id`, `env_name`, `jenkins_server_id`, `job_name`, `job_url`, `build_params`, `deploy_config`, `notification`, `is_active`, `created_at`, `updated_at`, `deleted_at`) VALUES (24, 10, 'test', NULL, '', NULL, NULL, NULL, NULL, 1, '2025-09-16 18:07:56.121', '2025-09-16 18:07:56.121', NULL);
INSERT INTO `app_jenkins_env` (`id`, `app_id`, `env_name`, `jenkins_server_id`, `job_name`, `job_url`, `build_params`, `deploy_config`, `notification`, `is_active`, `created_at`, `updated_at`, `deleted_at`) VALUES (25, 10, 'dev', 9, 's4-web', NULL, NULL, NULL, NULL, 1, '2025-09-16 18:07:56.342', '2025-09-18 10:29:40.413', NULL);
COMMIT;

-- ----------------------------
-- Table structure for cmdb_group
-- ----------------------------
DROP TABLE IF EXISTS `cmdb_group`;
CREATE TABLE `cmdb_group` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `parent_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '''父级分组ID''',
  `name` longtext NOT NULL COMMENT '''分组名称''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `remark` longtext COMMENT '''备注''',
  `update_time` datetime(3) DEFAULT NULL COMMENT '''更新时间''',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=30 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of cmdb_group
-- ----------------------------
BEGIN;
INSERT INTO `cmdb_group` (`id`, `parent_id`, `name`, `create_time`, `remark`, `update_time`) VALUES (1, 0, '默认业务组', '2025-07-10 11:02:07.226', NULL, NULL);
INSERT INTO `cmdb_group` (`id`, `parent_id`, `name`, `create_time`, `remark`, `update_time`) VALUES (2, 0, 'saas3业务线', '2025-07-10 11:03:36.622', NULL, NULL);
INSERT INTO `cmdb_group` (`id`, `parent_id`, `name`, `create_time`, `remark`, `update_time`) VALUES (3, 0, 'saas4业务线', '2025-07-10 11:03:56.083', NULL, NULL);
INSERT INTO `cmdb_group` (`id`, `parent_id`, `name`, `create_time`, `remark`, `update_time`) VALUES (4, 11, '公寓业务', '2025-07-10 11:05:45.431', NULL, NULL);
INSERT INTO `cmdb_group` (`id`, `parent_id`, `name`, `create_time`, `remark`, `update_time`) VALUES (9, 11, '智能窗帘业务', '2025-07-17 10:17:52.017', NULL, NULL);
INSERT INTO `cmdb_group` (`id`, `parent_id`, `name`, `create_time`, `remark`, `update_time`) VALUES (11, 0, 'saas5业务线', '2025-07-17 10:21:44.700', NULL, NULL);
INSERT INTO `cmdb_group` (`id`, `parent_id`, `name`, `create_time`, `remark`, `update_time`) VALUES (12, 2, '智能安防业务', '2025-07-17 11:07:19.279', NULL, NULL);
INSERT INTO `cmdb_group` (`id`, `parent_id`, `name`, `create_time`, `remark`, `update_time`) VALUES (14, 3, '智慧门锁业务', '2025-07-17 11:19:14.090', NULL, NULL);
INSERT INTO `cmdb_group` (`id`, `parent_id`, `name`, `create_time`, `remark`, `update_time`) VALUES (15, 0, 'saas6业务组', '2025-07-17 12:31:33.429', NULL, NULL);
INSERT INTO `cmdb_group` (`id`, `parent_id`, `name`, `create_time`, `remark`, `update_time`) VALUES (16, 15, '智能家居业务', '2025-07-22 10:32:14.837', NULL, NULL);
INSERT INTO `cmdb_group` (`id`, `parent_id`, `name`, `create_time`, `remark`, `update_time`) VALUES (17, 0, 'ops业务组', '2025-07-23 10:39:44.565', NULL, NULL);
INSERT INTO `cmdb_group` (`id`, `parent_id`, `name`, `create_time`, `remark`, `update_time`) VALUES (18, 2, '智控照明业务', '2025-07-23 10:40:12.787', NULL, NULL);
INSERT INTO `cmdb_group` (`id`, `parent_id`, `name`, `create_time`, `remark`, `update_time`) VALUES (19, 17, '运维组', '2025-07-23 20:21:08.618', NULL, NULL);
INSERT INTO `cmdb_group` (`id`, `parent_id`, `name`, `create_time`, `remark`, `update_time`) VALUES (20, 17, '安全组', '2025-07-23 20:21:33.222', NULL, NULL);
INSERT INTO `cmdb_group` (`id`, `parent_id`, `name`, `create_time`, `remark`, `update_time`) VALUES (24, 0, 'lockin鹿客云', '2025-09-07 17:49:17.318', NULL, NULL);
INSERT INTO `cmdb_group` (`id`, `parent_id`, `name`, `create_time`, `remark`, `update_time`) VALUES (26, 24, '国内鹿客云', '2025-09-07 18:04:56.683', NULL, NULL);
INSERT INTO `cmdb_group` (`id`, `parent_id`, `name`, `create_time`, `remark`, `update_time`) VALUES (27, 24, '新加坡鹿客云', '2025-09-07 18:05:13.713', NULL, NULL);
INSERT INTO `cmdb_group` (`id`, `parent_id`, `name`, `create_time`, `remark`, `update_time`) VALUES (28, 24, '北美鹿客云', '2025-09-07 18:05:31.846', NULL, NULL);
INSERT INTO `cmdb_group` (`id`, `parent_id`, `name`, `create_time`, `remark`, `update_time`) VALUES (29, 17, 'ai运维组', '2025-10-01 15:17:35.617', NULL, NULL);
COMMIT;

-- ----------------------------
-- Table structure for cmdb_host
-- ----------------------------
DROP TABLE IF EXISTS `cmdb_host`;
CREATE TABLE `cmdb_host` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `host_name` longtext NOT NULL COMMENT '''名称''',
  `group_id` bigint unsigned NOT NULL COMMENT '''分组ID''',
  `private_ip` longtext COMMENT '''私网IP''',
  `public_ip` longtext COMMENT '''公网IP''',
  `ssh_name` longtext COMMENT '''SSH用户名''',
  `ssh_key_id` bigint unsigned DEFAULT NULL COMMENT '''SSH凭据ID''',
  `ssh_port` bigint DEFAULT '22' COMMENT '''SSH端口''',
  `remark` longtext COMMENT '''备注''',
  `vendor` bigint DEFAULT NULL COMMENT '''1->自建,2->阿里云,3->腾讯云''',
  `region` longtext COMMENT '''区域''',
  `instance_id` longtext COMMENT '''实例ID''',
  `os` longtext COMMENT '''操作系统''',
  `status` bigint DEFAULT NULL COMMENT '''状态:1->认证成功,2->未认证,3->认证失败''',
  `cpu` longtext COMMENT '''CPU信息''',
  `memory` longtext COMMENT '''内存信息''',
  `disk` longtext COMMENT '''磁盘信息''',
  `billing_type` longtext COMMENT '''计费方式''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `expire_time` datetime(3) DEFAULT NULL COMMENT '''到期时间''',
  `update_time` datetime(3) DEFAULT NULL COMMENT '''更新时间''',
  `ssh_ip` longtext NOT NULL COMMENT '''SSH连接IP''',
  `name` longtext NOT NULL COMMENT '''ecs主机名称''',
  `ssh_gateway_id` bigint unsigned DEFAULT NULL COMMENT '''中转网关凭据ID''',
  PRIMARY KEY (`id`),
  KEY `fk_cmdb_group_hosts` (`group_id`),
  CONSTRAINT `fk_cmdb_group_hosts` FOREIGN KEY (`group_id`) REFERENCES `cmdb_group` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=501 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of cmdb_host
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for cmdb_sql
-- ----------------------------
DROP TABLE IF EXISTS `cmdb_sql`;
CREATE TABLE `cmdb_sql` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `type` int NOT NULL,
  `account_id` bigint unsigned NOT NULL,
  `group_id` bigint unsigned NOT NULL,
  `tags` varchar(255) DEFAULT NULL,
  `description` varchar(500) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of cmdb_sql
-- ----------------------------
BEGIN;
INSERT INTO `cmdb_sql` (`id`, `name`, `type`, `account_id`, `group_id`, `tags`, `description`, `created_at`, `updated_at`) VALUES (1, 'saas3-mysql', 1, 1, 1, 'prod', '1111', '2025-07-29 21:23:17.309', '2025-07-29 21:23:17.309');
INSERT INTO `cmdb_sql` (`id`, `name`, `type`, `account_id`, `group_id`, `tags`, `description`, `created_at`, `updated_at`) VALUES (2, 'saas3-redis-1', 3, 3, 1, 'prod', '1111', '2025-07-29 21:24:57.985', '2025-09-06 15:12:22.605');
INSERT INTO `cmdb_sql` (`id`, `name`, `type`, `account_id`, `group_id`, `tags`, `description`, `created_at`, `updated_at`) VALUES (4, 'saas3-pgsql', 2, 1, 1, 'prod', '1111', '2025-07-29 21:36:11.147', '2025-09-06 15:12:44.586');
INSERT INTO `cmdb_sql` (`id`, `name`, `type`, `account_id`, `group_id`, `tags`, `description`, `created_at`, `updated_at`) VALUES (5, 's3-mongo', 4, 1, 18, 'test', '123', '2025-07-30 11:16:27.710', '2025-10-01 15:17:09.321');
COMMIT;

-- ----------------------------
-- Table structure for cmdb_sql_log
-- ----------------------------
DROP TABLE IF EXISTS `cmdb_sql_log`;
CREATE TABLE `cmdb_sql_log` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `instance_id` varchar(64) NOT NULL,
  `database` varchar(128) NOT NULL,
  `operation_type` varchar(32) NOT NULL,
  `sql_content` text NOT NULL,
  `exec_user` varchar(64) NOT NULL,
  `ip` varchar(64) NOT NULL,
  `scanned_rows` bigint DEFAULT '0',
  `affected_rows` bigint DEFAULT '0',
  `execution_time` bigint DEFAULT '0',
  `returned_rows` bigint DEFAULT '0',
  `result` varchar(32) NOT NULL,
  `query_time` datetime(3) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_cmdb_sql_log_query_time` (`query_time`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of cmdb_sql_log
-- ----------------------------
BEGIN;
INSERT INTO `cmdb_sql_log` (`id`, `instance_id`, `database`, `operation_type`, `sql_content`, `exec_user`, `ip`, `scanned_rows`, `affected_rows`, `execution_time`, `returned_rows`, `result`, `query_time`) VALUES (7, '8.130.14.34', 'gin-api', 'SELECT', 'select * from  cmdb_host;', 'admin', '127.0.0.1', 0, 0, 59, 3, 'SUCCESS', '2025-08-25 10:17:24.672');
COMMIT;

-- ----------------------------
-- Table structure for cmdb_sql_records
-- ----------------------------
DROP TABLE IF EXISTS `cmdb_sql_records`;
CREATE TABLE `cmdb_sql_records` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `instance_id` varchar(64) NOT NULL,
  `database` varchar(128) NOT NULL,
  `operation_type` varchar(32) NOT NULL,
  `sql_content` text NOT NULL,
  `exec_user` varchar(64) NOT NULL,
  `scanned_rows` bigint DEFAULT '0',
  `affected_rows` bigint DEFAULT '0',
  `execution_time` bigint DEFAULT '0',
  `returned_rows` bigint DEFAULT '0',
  `result` varchar(32) NOT NULL,
  `query_time` datetime(3) NOT NULL,
  `name` varchar(64) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_cmdb_sql_records_query_time` (`query_time`)
) ENGINE=InnoDB AUTO_INCREMENT=19 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of cmdb_sql_records
-- ----------------------------
BEGIN;
INSERT INTO `cmdb_sql_records` (`id`, `instance_id`, `database`, `operation_type`, `sql_content`, `exec_user`, `scanned_rows`, `affected_rows`, `execution_time`, `returned_rows`, `result`, `query_time`, `name`) VALUES (1, '8.130.14.34:3306', 'gin-api', 'SELECT', 'select * from cmdb_group;', 'anonymous', 14, 0, 403, 14, 'SUCCESS', '2025-07-29 11:14:13.686', '');
INSERT INTO `cmdb_sql_records` (`id`, `instance_id`, `database`, `operation_type`, `sql_content`, `exec_user`, `scanned_rows`, `affected_rows`, `execution_time`, `returned_rows`, `result`, `query_time`, `name`) VALUES (2, '8.130.14.30', 'saas3-mysql', 'SELECT', 'select * from cmdb_group;', '', 0, 0, 50, 10, 'SUCCESS', '2025-07-30 13:29:24.409', '');
INSERT INTO `cmdb_sql_records` (`id`, `instance_id`, `database`, `operation_type`, `sql_content`, `exec_user`, `scanned_rows`, `affected_rows`, `execution_time`, `returned_rows`, `result`, `query_time`, `name`) VALUES (3, '8.130.14.34', 'gin-api', 'SELECT', 'select * from cmdb_group;', '', 0, 0, 54, 14, 'SUCCESS', '2025-07-30 13:58:13.386', '');
INSERT INTO `cmdb_sql_records` (`id`, `instance_id`, `database`, `operation_type`, `sql_content`, `exec_user`, `scanned_rows`, `affected_rows`, `execution_time`, `returned_rows`, `result`, `query_time`, `name`) VALUES (4, '8.130.14.34', 'saas3-mysql', 'INSERT', 'UPDATE `cmdb_group`SET `name` = \'sql测试组0000\' WHERE `parent_id` = 17AND `name` = \'sql测试组\';', '', 0, 1, 80, 0, 'SUCCESS', '2025-07-30 14:00:44.370', '');
INSERT INTO `cmdb_sql_records` (`id`, `instance_id`, `database`, `operation_type`, `sql_content`, `exec_user`, `scanned_rows`, `affected_rows`, `execution_time`, `returned_rows`, `result`, `query_time`, `name`) VALUES (5, '8.130.14.34', 'saas3-mysql', 'INSERT', 'UPDATE `cmdb_group` SET `name` = \'test123\' WHERE `id` = 22;', '', 0, 1, 80, 0, 'SUCCESS', '2025-07-30 14:04:30.684', '');
INSERT INTO `cmdb_sql_records` (`id`, `instance_id`, `database`, `operation_type`, `sql_content`, `exec_user`, `scanned_rows`, `affected_rows`, `execution_time`, `returned_rows`, `result`, `query_time`, `name`) VALUES (6, '8.130.14.34', 'saas3-mysql', 'INSERT', 'UPDATE `cmdb_group` SET `name` = \'test123111\' WHERE `id` = 22;', '', 0, 1, 80, 0, 'SUCCESS', '2025-07-30 14:06:34.692', '');
INSERT INTO `cmdb_sql_records` (`id`, `instance_id`, `database`, `operation_type`, `sql_content`, `exec_user`, `scanned_rows`, `affected_rows`, `execution_time`, `returned_rows`, `result`, `query_time`, `name`) VALUES (7, '8.130.14.34', 'gin-api', 'SELECT', 'UPDATE `cmdb_group` SET `name` = \'test123111\' WHERE `id` = 22;', '', 0, 0, 122, 0, 'SUCCESS', '2025-07-30 14:07:44.151', '');
INSERT INTO `cmdb_sql_records` (`id`, `instance_id`, `database`, `operation_type`, `sql_content`, `exec_user`, `scanned_rows`, `affected_rows`, `execution_time`, `returned_rows`, `result`, `query_time`, `name`) VALUES (8, '8.130.14.34', 'gin-api', 'SELECT', 'UPDATE `cmdb_group` SET `name` = \'test001\' WHERE `id` = 22;', '', 0, 0, 55, 0, 'SUCCESS', '2025-07-30 14:11:42.626', '');
INSERT INTO `cmdb_sql_records` (`id`, `instance_id`, `database`, `operation_type`, `sql_content`, `exec_user`, `scanned_rows`, `affected_rows`, `execution_time`, `returned_rows`, `result`, `query_time`, `name`) VALUES (9, '8.130.14.34', 'gin-api', 'SELECT', 'select * from cmdb_group;', '', 0, 0, 67, 15, 'SUCCESS', '2025-07-30 15:24:57.109', '');
INSERT INTO `cmdb_sql_records` (`id`, `instance_id`, `database`, `operation_type`, `sql_content`, `exec_user`, `scanned_rows`, `affected_rows`, `execution_time`, `returned_rows`, `result`, `query_time`, `name`) VALUES (10, '8.130.14.34', 'gin-api', 'SELECT', 'select * from cmdb_host;', '', 0, 0, 114, 75, 'SUCCESS', '2025-07-30 15:25:19.542', '');
INSERT INTO `cmdb_sql_records` (`id`, `instance_id`, `database`, `operation_type`, `sql_content`, `exec_user`, `scanned_rows`, `affected_rows`, `execution_time`, `returned_rows`, `result`, `query_time`, `name`) VALUES (11, '8.130.14.34', 'gin-api', 'SELECT', 'select * from cmdb_group;', '', 0, 0, 56, 15, 'SUCCESS', '2025-07-30 15:29:16.507', '');
INSERT INTO `cmdb_sql_records` (`id`, `instance_id`, `database`, `operation_type`, `sql_content`, `exec_user`, `scanned_rows`, `affected_rows`, `execution_time`, `returned_rows`, `result`, `query_time`, `name`) VALUES (12, '8.130.14.34', 'saas3-mysql', 'EXECUTE', 'create databases  db;', '', 0, 1, 100, 0, 'SUCCESS', '2025-07-30 15:34:31.246', '');
INSERT INTO `cmdb_sql_records` (`id`, `instance_id`, `database`, `operation_type`, `sql_content`, `exec_user`, `scanned_rows`, `affected_rows`, `execution_time`, `returned_rows`, `result`, `query_time`, `name`) VALUES (13, '8.130.14.34', 'gin-api', 'SELECT', 'select id,name from cmdb_group;', '', 0, 0, 51, 15, 'SUCCESS', '2025-07-30 17:18:52.168', '');
INSERT INTO `cmdb_sql_records` (`id`, `instance_id`, `database`, `operation_type`, `sql_content`, `exec_user`, `scanned_rows`, `affected_rows`, `execution_time`, `returned_rows`, `result`, `query_time`, `name`) VALUES (14, '8.130.14.34', 'gin-api', 'SELECT', 'select id,name from cmdb_group;', '', 0, 0, 53, 15, 'SUCCESS', '2025-07-30 17:29:36.465', '');
INSERT INTO `cmdb_sql_records` (`id`, `instance_id`, `database`, `operation_type`, `sql_content`, `exec_user`, `scanned_rows`, `affected_rows`, `execution_time`, `returned_rows`, `result`, `query_time`, `name`) VALUES (15, '8.130.14.34', 'gin-api', 'SELECT', 'select id,name from  cmdb_group;', '', 0, 0, 56, 15, 'SUCCESS', '2025-07-30 21:06:07.136', '');
INSERT INTO `cmdb_sql_records` (`id`, `instance_id`, `database`, `operation_type`, `sql_content`, `exec_user`, `scanned_rows`, `affected_rows`, `execution_time`, `returned_rows`, `result`, `query_time`, `name`) VALUES (16, '8.130.14.34', 'gin-api', 'SELECT', 'select id,name from  cmdb_group;', '', 0, 0, 64, 15, 'SUCCESS', '2025-07-30 21:12:04.886', '');
INSERT INTO `cmdb_sql_records` (`id`, `instance_id`, `database`, `operation_type`, `sql_content`, `exec_user`, `scanned_rows`, `affected_rows`, `execution_time`, `returned_rows`, `result`, `query_time`, `name`) VALUES (17, '8.130.14.34', 'gin-api', 'SELECT', 'select id,name from  cmdb_group;', 'admin', 0, 0, 51, 15, 'SUCCESS', '2025-07-30 21:23:55.532', '');
INSERT INTO `cmdb_sql_records` (`id`, `instance_id`, `database`, `operation_type`, `sql_content`, `exec_user`, `scanned_rows`, `affected_rows`, `execution_time`, `returned_rows`, `result`, `query_time`, `name`) VALUES (18, '8.130.14.34', 'gin-api', 'SELECT', 'select id,name from  cmdb_group;', 'zhangsan', 0, 0, 54, 15, 'SUCCESS', '2025-07-30 21:26:51.642', '');
COMMIT;

-- ----------------------------
-- Table structure for config_account
-- ----------------------------
DROP TABLE IF EXISTS `config_account`;
CREATE TABLE `config_account` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `alias` varchar(128) NOT NULL,
  `host` varchar(128) NOT NULL,
  `port` bigint NOT NULL,
  `name` varchar(128) NOT NULL,
  `password` text NOT NULL,
  `type` bigint NOT NULL,
  `remark` text,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of config_account
-- ----------------------------
BEGIN;
INSERT INTO `config_account` (`id`, `alias`, `host`, `port`, `name`, `password`, `type`, `remark`, `created_at`, `updated_at`) VALUES (11, '阿里云jenkins', '180.76.231.65', 8080, 'root', 'KioCHW/kEA093zr+cBHECf+1xxfxTGgwy6KJXCKcHjo=', 4, '阿里云jenkins', '2025-09-16 14:57:22.089', '2025-09-16 14:57:22.089');
COMMIT;

-- ----------------------------
-- Table structure for config_ecsauth
-- ----------------------------
DROP TABLE IF EXISTS `config_ecsauth`;
CREATE TABLE `config_ecsauth` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `name` longtext NOT NULL COMMENT '''凭证名称''',
  `type` bigint NOT NULL COMMENT '''认证类型:1->密码,2->私钥,3->公钥(免认证)''',
  `username` longtext COMMENT '''用户名''',
  `password` longtext COMMENT '''密码(type=1时使用)''',
  `public_key` text COMMENT '''私钥内容(type=2时使用，字段名历史原因)''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `remark` longtext COMMENT '''备注''',
  `port` bigint DEFAULT '22' COMMENT '''端口号''',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of config_ecsauth
-- ----------------------------
BEGIN;
INSERT INTO `config_ecsauth` (`id`, `name`, `type`, `username`, `password`, `public_key`, `create_time`, `remark`, `port`) VALUES (13, 'ubuntu-1', 2, 'root', '', '-----BEGIN OPENSSH PRIVATE KEY-----\nb3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAABlwAAAAdzc2gtcn\nNhAAAAAwEAAQAAAYEAqQG9d/UFU0eaa9rMd5l31mX1OLV8PKfuwiAy1vu4MDsIyWcQdMjx\nn2Sy/Oj91RKUgnb/LNLQNpM7ARYlUJdyU5X7/PZVUF7bKV7RSdN5bcBTW36uOaX5K4mWNW\nVWm/V0zFxllMb5tK5t3OlSyekEcFuKTSvsmZFfSaNBLkRY3qAC64RVRcxDiw0yHR68kipv\n3QvJfj2uFTt9X+lKXYsHebTvxgIVgM+KtkEi5R7MZNFiqeTTWGOidwyooVD+wJPkxFvi9u\nmIPkm7/ExR4A/WwQpIbyiA4r8k2hYFGRfbn2J2I7xq9t59xl57TtYRQ0VWbGSj+xxHFIOX\novkk79fea1zgCUxbiTLESDm9rSIbQHWPC3uGhgoiwGn6Xb4VI+Xb6hieOdcbtftZq176Ko\nYivJrw8eNGuMZf7/VbsVtQZhNovHTeXkVDBLzb45qU3oyhM/Bi0OkizR0NEDQGNZvCs3P3\nkY6hLwSyC9FbOSPaG63PZpnUP/F2lNjbDx+w5qHBAAAFiKrffV+q331fAAAAB3NzaC1yc2\nEAAAGBAKkBvXf1BVNHmmvazHeZd9Zl9Ti1fDyn7sIgMtb7uDA7CMlnEHTI8Z9ksvzo/dUS\nlIJ2/yzS0DaTOwEWJVCXclOV+/z2VVBe2yle0UnTeW3AU1t+rjml+SuJljVlVpv1dMxcZZ\nTG+bSubdzpUsnpBHBbik0r7JmRX0mjQS5EWN6gAuuEVUXMQ4sNMh0evJIqb90LyX49rhU7\nfV/pSl2LB3m078YCFYDPirZBIuUezGTRYqnk01hjoncMqKFQ/sCT5MRb4vbpiD5Ju/xMUe\nAP1sEKSG8ogOK/JNoWBRkX259idiO8avbefcZee07WEUNFVmxko/scRxSDl6L5JO/X3mtc\n4AlMW4kyxEg5va0iG0B1jwt7hoYKIsBp+l2+FSPl2+oYnjnXG7X7Wate+iqGIrya8PHjRr\njGX+/1W7FbUGYTaLx03l5FQwS82+OalN6MoTPwYtDpIs0dDRA0BjWbwrNz95GOoS8EsgvR\nWzkj2hutz2aZ1D/xdpTY2w8fsOahwQAAAAMBAAEAAAGAPMbzbAAhOiG11rOIdDosbl2cIh\nih5O4/XnBV2BoN0spoDoWR1W2t9rQv6eOf5NAZuBEzEtt9JSMtRObB5ImSA50cUYwEgjKa\nffYrLXFvaZiQYYbdAL4/LPj9m5Yl31AWAKf5n9cuVVulBpyhuRqjDgTdZ4M0QsWUjta2yf\nyqOOUyJ6AfSkp6v7avHum+xpGcSNHEVDG6OMh7/dHXfGxS8+GtmHeHZBs+tuwVsG+y8U6b\nPJ2do1uxQT68JccK1hZYyOTCPyK/2jQxSXFJEG6FY1jndxkaSvsKFe39jzdxRic1D6R1bN\nbeZASTGrBzr734BaWtjQkzpO9kJ8kTLjbP8MVQYY4wCTfcnf1eAhFMlJJY4YGohrZ6FhAX\nuOMrP+v6n0m8dXsrIMAgRWiYh1K8ZAK7tIIlM5SPyh7E8YbCGHuQNF2XARebLk0o13u8ff\n5U58bVoEJBnXrOuiJdbvR/7QRcWcZUTV9kdRisF3ooo+TKWd7VOK2BAXHsam56PXwBAAAA\nwQC43Px4FX19yPFXAu7zhw5IeH1XdfAoPzToYFe1VtvHYKAG3ivcVmNqo2EzUG1TJIIY7f\nclepUAZqUj4c3R/GafwPY8QW3QbpMimpx/hVuyA3yDC0qjeaQ335yvdTtVu/Jv4yXAbjRi\n2kgGDfoPQoHnNlNcfhcZI5c1e8yYqNy+cnPi90dnzvLACsIIrcw82IEN6cjf0XomplTdvM\nR63qxvOAWOQZxcoUbZuTUPIqb8P1M/vzOS0mfvMqj/dqW8mv4AAADBANNix6vGRxWQoPcK\nuY29BV2brwV3SKlFlTUtj8+YXIWdGVA4k94QJGeSr+mikiO8PySjeBZB/Ki8qdh+vPMKZ4\nCZM9Uk38qkOdoAZTlJ2zCTIQoWGt1j4eSiAy38LBEiruAJPtXXOWTKSM0g/WaaX/8EhyTp\nEuODV5UVJlsg6mXeH4OgcCo9B2Fmn5l8AKiexCucITGuRdzVhN3KNph/EIzd0Kn7zL2dMW\nT0rtYL3Ja1y70UAqwDdVUyJjQFtrbbiQAAAMEAzK02BABVLbgcnEGRWeNaJQ7lA36hU3x1\n6XgmESJWSPN4Yp5gNQBHsuFxB6ksZQkX76penNQG1cJPu7xURP4lrCQdyoCfdzcO+PpNKX\nfPCBNLNGCefWKjcBdp961bkhNazH1B2p5g4lytZbS2NJduB3W+L7dcF18mnBCIlNdLxBm7\nZ/PDjQfAhurtKZHvhWjg9qX32xqETldjePZUNHdt1TQ0Ih0Ei0J23trsKxBXkP0SzAIXiR\n/z4Lox2qK0BG55AAAADnJvb3RAdWJ1bnR1LTAxAQIDBA==\n-----END OPENSSH PRIVATE KEY-----', '2025-07-28 16:11:42.476', '这是自建虚拟机', 22);
INSERT INTO `config_ecsauth` (`id`, `name`, `type`, `username`, `password`, `public_key`, `create_time`, `remark`, `port`) VALUES (20, '免密码认证', 3, 'root', '', '', '2025-10-14 21:06:27.749', '', 22);
COMMIT;

-- ----------------------------
-- Table structure for config_keymanage
-- ----------------------------
DROP TABLE IF EXISTS `config_keymanage`;
CREATE TABLE `config_keymanage` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `key_type` bigint NOT NULL,
  `key_id` text NOT NULL,
  `key_secret` text NOT NULL,
  `remark` text,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of config_keymanage
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for config_sync_schedule
-- ----------------------------
DROP TABLE IF EXISTS `config_sync_schedule`;
CREATE TABLE `config_sync_schedule` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `cron_expr` varchar(100) NOT NULL,
  `key_types` text NOT NULL,
  `status` bigint NOT NULL DEFAULT '1',
  `last_run_time` datetime(3) DEFAULT NULL,
  `next_run_time` datetime(3) DEFAULT NULL,
  `remark` text,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `sync_log` text,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of config_sync_schedule
-- ----------------------------
BEGIN;
INSERT INTO `config_sync_schedule` (`id`, `name`, `cron_expr`, `key_types`, `status`, `last_run_time`, `next_run_time`, `remark`, `created_at`, `updated_at`, `sync_log`) VALUES (3, '阿里云定时同步', '*/3 * * * *', '[1]', 0, '2025-10-01 13:33:00.001', '2025-10-01 13:36:00.000', '', '2025-09-29 18:41:10.257', '2025-10-01 13:34:41.439', '[2025-10-01 13:33:00] 开始同步\n- 阿里云: 同步成功\n[2025-10-01 13:33:00] 同步完成，耗时: 827.31718ms\n');
COMMIT;

-- ----------------------------
-- Table structure for k8s_cluster
-- ----------------------------
DROP TABLE IF EXISTS `k8s_cluster`;
CREATE TABLE `k8s_cluster` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键ID''',
  `name` varchar(100) NOT NULL COMMENT '''集群名称''',
  `version` varchar(50) NOT NULL COMMENT '''集群版本''',
  `status` bigint NOT NULL DEFAULT '1' COMMENT '''集群状态:1-创建中,2-运行中,3-离线''',
  `credential` text COMMENT '''集群凭证(kubeconfig)''',
  `description` text COMMENT '''集群描述''',
  `cluster_type` bigint NOT NULL DEFAULT '1' COMMENT '''集群类型:1-自建,2-导入''',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `node_count` bigint DEFAULT '0' COMMENT '''节点数量''',
  `ready_nodes` bigint DEFAULT '0' COMMENT '''就绪节点数''',
  `master_nodes` bigint DEFAULT '0' COMMENT '''Master节点数''',
  `worker_nodes` bigint DEFAULT '0' COMMENT '''Worker节点数''',
  `last_sync_at` datetime(3) DEFAULT NULL COMMENT '''最后同步时间''',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_k8s_cluster_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=32 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of k8s_cluster
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for monitor_agent
-- ----------------------------
DROP TABLE IF EXISTS `monitor_agent`;
CREATE TABLE `monitor_agent` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `host_id` bigint unsigned NOT NULL COMMENT '''主机ID''',
  `host_name` longtext COMMENT '''主机名称''',
  `version` varchar(191) DEFAULT '1.0.0' COMMENT '''Agent版本''',
  `status` bigint DEFAULT NULL COMMENT '''状态:1->部署中,2->部署失败,3->运行中,4->已停止''',
  `install_path` longtext COMMENT '''安装路径''',
  `port` bigint DEFAULT '9100' COMMENT '''监听端口''',
  `pid` bigint DEFAULT NULL COMMENT '''进程ID''',
  `last_heartbeat` datetime(3) DEFAULT NULL COMMENT '''最后心跳时间''',
  `update_time` datetime(3) DEFAULT NULL COMMENT '''更新时间''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `error_msg` text COMMENT '''错误信息''',
  `install_progress` bigint DEFAULT '0' COMMENT '''安装进度(0-100)''',
  PRIMARY KEY (`id`),
  KEY `idx_monitor_agent_host_id` (`host_id`)
) ENGINE=InnoDB AUTO_INCREMENT=76 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of monitor_agent
-- ----------------------------
BEGIN;
INSERT INTO `monitor_agent` (`id`, `host_id`, `host_name`, `version`, `status`, `install_path`, `port`, `pid`, `last_heartbeat`, `update_time`, `create_time`, `error_msg`, `install_progress`) VALUES (75, 1, 'go-ops', '1.0.0', 3, '/opt/agent', 9100, 0, NULL, '2025-10-10 10:57:20.911', '2025-10-10 10:55:49.579', '', 100);
COMMIT;

-- ----------------------------
-- Table structure for quick_deployment_tasks
-- ----------------------------
DROP TABLE IF EXISTS `quick_deployment_tasks`;
CREATE TABLE `quick_deployment_tasks` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `deployment_id` bigint unsigned NOT NULL,
  `app_id` bigint unsigned NOT NULL,
  `app_name` varchar(255) DEFAULT NULL,
  `app_code` varchar(255) DEFAULT NULL,
  `environment` varchar(50) DEFAULT NULL,
  `jenkins_env_id` bigint unsigned NOT NULL,
  `jenkins_job_url` varchar(500) DEFAULT NULL,
  `build_number` bigint DEFAULT NULL,
  `status` tinyint DEFAULT '1',
  `execute_order` bigint NOT NULL,
  `start_time` datetime(3) DEFAULT NULL,
  `end_time` datetime(3) DEFAULT NULL,
  `duration` bigint DEFAULT NULL,
  `error_message` text,
  `log_url` varchar(500) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_quick_deployment_tasks_deployment_id` (`deployment_id`),
  KEY `fk_quick_deployment_tasks_jenkins_env` (`jenkins_env_id`),
  KEY `fk_quick_deployment_tasks_application` (`app_id`),
  CONSTRAINT `fk_quick_deployment_tasks_application` FOREIGN KEY (`app_id`) REFERENCES `app_application` (`id`),
  CONSTRAINT `fk_quick_deployment_tasks_jenkins_env` FOREIGN KEY (`jenkins_env_id`) REFERENCES `app_jenkins_env` (`id`),
  CONSTRAINT `fk_quick_deployments_tasks` FOREIGN KEY (`deployment_id`) REFERENCES `quick_deployments` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=24 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of quick_deployment_tasks
-- ----------------------------
BEGIN;
INSERT INTO `quick_deployment_tasks` (`id`, `deployment_id`, `app_id`, `app_name`, `app_code`, `environment`, `jenkins_env_id`, `jenkins_job_url`, `build_number`, `status`, `execute_order`, `start_time`, `end_time`, `duration`, `error_message`, `log_url`, `created_at`, `updated_at`) VALUES (10, 6, 4, 's3-api', 's3-api', 'test', 5, '', 5, 4, 1, '2025-09-17 11:46:15.519', '2025-09-17 11:46:26.570', 11, 'Jenkins构建失败', 'http://180.76.231.65:8080/job/s3-api/5/console', '2025-09-17 11:46:00.369', '2025-09-17 11:46:26.621');
INSERT INTO `quick_deployment_tasks` (`id`, `deployment_id`, `app_id`, `app_name`, `app_code`, `environment`, `jenkins_env_id`, `jenkins_job_url`, `build_number`, `status`, `execute_order`, `start_time`, `end_time`, `duration`, `error_message`, `log_url`, `created_at`, `updated_at`) VALUES (11, 7, 4, 's3-api', 's3-api', 'test', 5, '', 6, 4, 1, '2025-09-17 11:57:23.584', '2025-09-17 11:57:34.613', 11, 'Jenkins构建失败', 'http://180.76.231.65:8080/job/s3-api/6/console', '2025-09-17 11:57:04.011', '2025-09-17 11:57:34.665');
INSERT INTO `quick_deployment_tasks` (`id`, `deployment_id`, `app_id`, `app_name`, `app_code`, `environment`, `jenkins_env_id`, `jenkins_job_url`, `build_number`, `status`, `execute_order`, `start_time`, `end_time`, `duration`, `error_message`, `log_url`, `created_at`, `updated_at`) VALUES (12, 8, 4, 's3-api', 's3-api', 'test', 5, '', 6, 3, 1, '2025-09-17 13:32:00.216', '2025-09-17 13:32:04.600', 4, '', 'http://180.76.231.65:8080/job/s3-api/6/console', '2025-09-17 13:31:45.833', '2025-09-17 13:32:04.652');
INSERT INTO `quick_deployment_tasks` (`id`, `deployment_id`, `app_id`, `app_name`, `app_code`, `environment`, `jenkins_env_id`, `jenkins_job_url`, `build_number`, `status`, `execute_order`, `start_time`, `end_time`, `duration`, `error_message`, `log_url`, `created_at`, `updated_at`) VALUES (13, 9, 5, 's3-web', 's3-web', 'test', 8, '', 2, 3, 1, '2025-09-17 13:39:32.334', '2025-09-17 13:41:21.430', 109, '', 'http://180.76.231.65:8080/job/s3-web/2/console', '2025-09-17 13:39:12.959', '2025-09-17 13:41:21.479');
INSERT INTO `quick_deployment_tasks` (`id`, `deployment_id`, `app_id`, `app_name`, `app_code`, `environment`, `jenkins_env_id`, `jenkins_job_url`, `build_number`, `status`, `execute_order`, `start_time`, `end_time`, `duration`, `error_message`, `log_url`, `created_at`, `updated_at`) VALUES (14, 10, 5, 's3-web', 's3-web', 'test', 8, '', 5, 3, 1, '2025-09-17 14:26:36.991', '2025-09-17 14:26:55.134', 18, '', 'http://180.76.231.65:8080/job/s3-web/5/console', '2025-09-17 14:26:06.943', '2025-09-17 14:26:55.185');
INSERT INTO `quick_deployment_tasks` (`id`, `deployment_id`, `app_id`, `app_name`, `app_code`, `environment`, `jenkins_env_id`, `jenkins_job_url`, `build_number`, `status`, `execute_order`, `start_time`, `end_time`, `duration`, `error_message`, `log_url`, `created_at`, `updated_at`) VALUES (15, 10, 4, 's3-api', 's3-api', 'test', 5, '', 8, 3, 2, '2025-09-17 14:26:44.716', '2025-09-17 14:28:33.960', 109, '', 'http://180.76.231.65:8080/job/s3-api/8/console', '2025-09-17 14:26:07.352', '2025-09-17 14:28:34.008');
INSERT INTO `quick_deployment_tasks` (`id`, `deployment_id`, `app_id`, `app_name`, `app_code`, `environment`, `jenkins_env_id`, `jenkins_job_url`, `build_number`, `status`, `execute_order`, `start_time`, `end_time`, `duration`, `error_message`, `log_url`, `created_at`, `updated_at`) VALUES (16, 11, 4, 's3-api', 's3-api', 'dev', 6, '', 9, 3, 1, '2025-09-17 15:21:26.032', '2025-09-17 15:23:17.544', 111, '', 'http://180.76.231.65:8080/job/s3-api/9/console', '2025-09-17 15:18:44.619', '2025-09-17 15:23:17.599');
INSERT INTO `quick_deployment_tasks` (`id`, `deployment_id`, `app_id`, `app_name`, `app_code`, `environment`, `jenkins_env_id`, `jenkins_job_url`, `build_number`, `status`, `execute_order`, `start_time`, `end_time`, `duration`, `error_message`, `log_url`, `created_at`, `updated_at`) VALUES (17, 11, 5, 's3-web', 's3-web', 'dev', 9, '', 6, 3, 2, '2025-09-17 15:21:35.936', '2025-09-17 15:21:56.142', 20, '', 'http://180.76.231.65:8080/job/s3-web/6/console', '2025-09-17 15:18:45.079', '2025-09-17 15:21:56.190');
INSERT INTO `quick_deployment_tasks` (`id`, `deployment_id`, `app_id`, `app_name`, `app_code`, `environment`, `jenkins_env_id`, `jenkins_job_url`, `build_number`, `status`, `execute_order`, `start_time`, `end_time`, `duration`, `error_message`, `log_url`, `created_at`, `updated_at`) VALUES (18, 12, 4, 's3-api', 's3-api', 'dev', 6, '', 10, 3, 1, '2025-09-17 15:43:37.641', '2025-09-17 15:45:26.874', 109, '', 'http://180.76.231.65:8080/job/s3-api/10/console', '2025-09-17 15:43:23.726', '2025-09-17 15:45:26.927');
INSERT INTO `quick_deployment_tasks` (`id`, `deployment_id`, `app_id`, `app_name`, `app_code`, `environment`, `jenkins_env_id`, `jenkins_job_url`, `build_number`, `status`, `execute_order`, `start_time`, `end_time`, `duration`, `error_message`, `log_url`, `created_at`, `updated_at`) VALUES (19, 12, 5, 's3-web', 's3-web', 'dev', 9, '', 7, 3, 2, '2025-09-17 15:45:27.567', '2025-09-17 15:45:45.629', 18, '', 'http://180.76.231.65:8080/job/s3-web/7/console', '2025-09-17 15:43:24.121', '2025-09-17 15:45:45.681');
INSERT INTO `quick_deployment_tasks` (`id`, `deployment_id`, `app_id`, `app_name`, `app_code`, `environment`, `jenkins_env_id`, `jenkins_job_url`, `build_number`, `status`, `execute_order`, `start_time`, `end_time`, `duration`, `error_message`, `log_url`, `created_at`, `updated_at`) VALUES (20, 13, 6, 's4-api', 's4-api', 'test', 12, '', 1, 3, 1, '2025-09-21 18:11:04.606', '2025-09-21 18:12:54.230', 109, '', 'http://180.76.231.65:8080/job/s4-api/1/console', '2025-09-21 18:10:45.631', '2025-09-21 18:12:54.287');
INSERT INTO `quick_deployment_tasks` (`id`, `deployment_id`, `app_id`, `app_name`, `app_code`, `environment`, `jenkins_env_id`, `jenkins_job_url`, `build_number`, `status`, `execute_order`, `start_time`, `end_time`, `duration`, `error_message`, `log_url`, `created_at`, `updated_at`) VALUES (21, 13, 7, 's4-web', 's4-web', 'test', 15, '', 1, 3, 2, '2025-09-21 18:11:12.478', '2025-09-21 18:13:04.108', 111, '', 'http://180.76.231.65:8080/job/s4-web/1/console', '2025-09-21 18:10:46.164', '2025-09-21 18:13:04.175');
INSERT INTO `quick_deployment_tasks` (`id`, `deployment_id`, `app_id`, `app_name`, `app_code`, `environment`, `jenkins_env_id`, `jenkins_job_url`, `build_number`, `status`, `execute_order`, `start_time`, `end_time`, `duration`, `error_message`, `log_url`, `created_at`, `updated_at`) VALUES (22, 14, 9, 's5-api', 's5-api', 'dev', 22, '', 2, 3, 1, '2025-09-22 11:49:53.750', '2025-09-22 11:51:45.360', 111, '', 'http://180.76.231.65:8080/job/s4-api/2/console', '2025-09-22 11:49:40.811', '2025-09-22 11:51:45.408');
INSERT INTO `quick_deployment_tasks` (`id`, `deployment_id`, `app_id`, `app_name`, `app_code`, `environment`, `jenkins_env_id`, `jenkins_job_url`, `build_number`, `status`, `execute_order`, `start_time`, `end_time`, `duration`, `error_message`, `log_url`, `created_at`, `updated_at`) VALUES (23, 14, 8, 's5-web', 's5-web', 'dev', 19, '', 1, 3, 2, '2025-09-22 11:51:46.452', '2025-09-22 11:53:35.660', 109, '', 'http://180.76.231.65:8080/job/s5-web/1/console', '2025-09-22 11:49:41.228', '2025-09-22 11:53:35.714');
COMMIT;

-- ----------------------------
-- Table structure for quick_deployments
-- ----------------------------
DROP TABLE IF EXISTS `quick_deployments`;
CREATE TABLE `quick_deployments` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL,
  `business_group_id` bigint unsigned NOT NULL,
  `business_dept_id` bigint unsigned NOT NULL,
  `description` text,
  `status` tinyint DEFAULT '1',
  `task_count` bigint NOT NULL DEFAULT '0',
  `creator_id` bigint unsigned NOT NULL,
  `creator_name` varchar(100) DEFAULT NULL,
  `start_time` datetime(3) DEFAULT NULL,
  `end_time` datetime(3) DEFAULT NULL,
  `duration` bigint DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `execution_mode` tinyint DEFAULT '1',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of quick_deployments
-- ----------------------------
BEGIN;
INSERT INTO `quick_deployments` (`id`, `title`, `business_group_id`, `business_dept_id`, `description`, `status`, `task_count`, `creator_id`, `creator_name`, `start_time`, `end_time`, `duration`, `created_at`, `updated_at`, `execution_mode`) VALUES (6, '测试任务004', 12, 2, '123', 4, 1, 1, '管理员', '2025-09-17 11:46:14.976', '2025-09-17 11:46:26.783', 11, '2025-09-17 11:45:59.884', '2025-09-17 11:46:26.899', 1);
INSERT INTO `quick_deployments` (`id`, `title`, `business_group_id`, `business_dept_id`, `description`, `status`, `task_count`, `creator_id`, `creator_name`, `start_time`, `end_time`, `duration`, `created_at`, `updated_at`, `execution_mode`) VALUES (7, '测试任务005', 18, 2, '测试任务005', 4, 1, 1, '管理员', '2025-09-17 11:57:23.071', '2025-09-17 11:57:34.819', 11, '2025-09-17 11:57:03.363', '2025-09-17 11:57:34.866', 1);
INSERT INTO `quick_deployments` (`id`, `title`, `business_group_id`, `business_dept_id`, `description`, `status`, `task_count`, `creator_id`, `creator_name`, `start_time`, `end_time`, `duration`, `created_at`, `updated_at`, `execution_mode`) VALUES (8, '测试任务007', 12, 2, '测试任务007', 3, 1, 1, '管理员', '2025-09-17 13:31:59.690', '2025-09-17 13:32:04.330', 4, '2025-09-17 13:31:45.407', '2025-09-17 13:32:04.377', 1);
INSERT INTO `quick_deployments` (`id`, `title`, `business_group_id`, `business_dept_id`, `description`, `status`, `task_count`, `creator_id`, `creator_name`, `start_time`, `end_time`, `duration`, `created_at`, `updated_at`, `execution_mode`) VALUES (9, '测试任务006', 18, 2, '测试任务006', 3, 1, 1, '管理员', '2025-09-17 13:39:31.848', '2025-09-17 13:39:39.887', 7, '2025-09-17 13:39:12.576', '2025-09-17 13:39:39.935', 1);
INSERT INTO `quick_deployments` (`id`, `title`, `business_group_id`, `business_dept_id`, `description`, `status`, `task_count`, `creator_id`, `creator_name`, `start_time`, `end_time`, `duration`, `created_at`, `updated_at`, `execution_mode`) VALUES (10, '双任务测试', 18, 2, '双任务测试', 3, 2, 1, '管理员', '2025-09-17 14:26:36.482', '2025-09-17 14:26:52.382', 15, '2025-09-17 14:26:06.527', '2025-09-17 14:26:52.431', 1);
INSERT INTO `quick_deployments` (`id`, `title`, `business_group_id`, `business_dept_id`, `description`, `status`, `task_count`, `creator_id`, `creator_name`, `start_time`, `end_time`, `duration`, `created_at`, `updated_at`, `execution_mode`) VALUES (11, '测试任务执行时间统计', 12, 2, '测试任务执行时间统计', 3, 2, 1, '管理员', '2025-09-17 15:21:25.500', '2025-09-17 15:21:45.684', 0, '2025-09-17 15:18:44.191', '2025-09-17 15:21:45.861', 1);
INSERT INTO `quick_deployments` (`id`, `title`, `business_group_id`, `business_dept_id`, `description`, `status`, `task_count`, `creator_id`, `creator_name`, `start_time`, `end_time`, `duration`, `created_at`, `updated_at`, `execution_mode`) VALUES (12, '测试串行任务', 12, 2, '测试串行任务', 3, 2, 1, '管理员', '2025-09-17 15:43:37.100', '2025-09-17 15:45:46.056', 127, '2025-09-17 15:43:23.318', '2025-09-17 15:45:46.216', 2);
INSERT INTO `quick_deployments` (`id`, `title`, `business_group_id`, `business_dept_id`, `description`, `status`, `task_count`, `creator_id`, `creator_name`, `start_time`, `end_time`, `duration`, `created_at`, `updated_at`, `execution_mode`) VALUES (13, 'saas4新版本迭代', 14, 2, 'saas4新版本迭代', 3, 2, 1, '管理员', '2025-09-21 18:11:04.046', '2025-09-21 18:13:04.593', 220, '2025-09-21 18:10:45.148', '2025-09-21 18:13:04.789', 1);
INSERT INTO `quick_deployments` (`id`, `title`, `business_group_id`, `business_dept_id`, `description`, `status`, `task_count`, `creator_id`, `creator_name`, `start_time`, `end_time`, `duration`, `created_at`, `updated_at`, `execution_mode`) VALUES (14, 'saas5-新服务上线', 9, 5, 'saas5-新服务上线', 3, 2, 1, '管理员', '2025-09-22 11:49:53.183', '2025-09-22 11:53:36.090', 220, '2025-09-22 11:49:40.389', '2025-09-22 11:53:36.251', 2);
COMMIT;

-- ----------------------------
-- Table structure for sys_admin
-- ----------------------------
DROP TABLE IF EXISTS `sys_admin`;
CREATE TABLE `sys_admin` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
  `post_id` int DEFAULT NULL COMMENT '岗位id',
  `dept_id` int DEFAULT NULL COMMENT '部门id',
  `username` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL COMMENT '账号',
  `password` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL COMMENT '密码',
  `nickname` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '昵称',
  `icon` varchar(500) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '头像',
  `email` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '邮箱',
  `phone` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '手机',
  `note` varchar(500) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '备注信息',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `status` int NOT NULL DEFAULT '1' COMMENT '帐号启用状态：1->启用,2->禁用',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `username` (`username`) USING BTREE,
  KEY `create_time` (`create_time`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=103 DEFAULT CHARSET=utf8mb3 ROW_FORMAT=DYNAMIC COMMENT='后台管理员表';

-- ----------------------------
-- Records of sys_admin
-- ----------------------------
BEGIN;
INSERT INTO `sys_admin` (`id`, `post_id`, `dept_id`, `username`, `password`, `nickname`, `icon`, `email`, `phone`, `note`, `create_time`, `status`) VALUES (89, 1, 15, 'admin', 'e10adc3949ba59abbe56e057f20f883e', 'admin', 'http://10.7.16.22:8080/api/v1/upload/20251010/702910000.png', '123456789@qq.com', '13754354536', '后端研发', '2023-05-23 22:15:50', 1);
INSERT INTO `sys_admin` (`id`, `post_id`, `dept_id`, `username`, `password`, `nickname`, `icon`, `email`, `phone`, `note`, `create_time`, `status`) VALUES (98, 12, 15, 'lisi', '827ccb0eea8a706c4c34a16891f84e7b', '李四', 'http://127.0.0.1:8000/upload/20250629/500333800.png', '123@qq.com', '13826541566', 'ops', '2025-06-20 11:54:02', 1);
INSERT INTO `sys_admin` (`id`, `post_id`, `dept_id`, `username`, `password`, `nickname`, `icon`, `email`, `phone`, `note`, `create_time`, `status`) VALUES (99, 10, 6, 'zhangsan', 'e10adc3949ba59abbe56e057f20f883e', '大白同学', 'http://127.0.0.1:8000/upload/20250629/341737900.jpg', 'zhangfan@lockin.com', '13826541566', '123', '2025-06-29 15:02:48', 1);
INSERT INTO `sys_admin` (`id`, `post_id`, `dept_id`, `username`, `password`, `nickname`, `icon`, `email`, `phone`, `note`, `create_time`, `status`) VALUES (101, 11, 5, 'wangwu', 'e10adc3949ba59abbe56e057f20f883e', '王五', 'http://127.0.0.1:8000/upload/20250629/341737900.jpg', 'zf@123.com', '13826541566', '123', '2025-07-01 14:34:18', 1);
INSERT INTO `sys_admin` (`id`, `post_id`, `dept_id`, `username`, `password`, `nickname`, `icon`, `email`, `phone`, `note`, `create_time`, `status`) VALUES (102, 12, 14, 'test', 'e10adc3949ba59abbe56e057f20f883e', 'test', 'http://10.7.16.22:8080/api/v1/upload/20251009/699174000.png', 'zfwh1024@163.com', '13826541511', '游客', '2025-09-24 12:49:06', 1);
COMMIT;

-- ----------------------------
-- Table structure for sys_admin_role
-- ----------------------------
DROP TABLE IF EXISTS `sys_admin_role`;
CREATE TABLE `sys_admin_role` (
  `admin_id` int NOT NULL COMMENT '管理员id',
  `role_id` int NOT NULL COMMENT '角色id'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 ROW_FORMAT=DYNAMIC COMMENT='管理员和角色关系表';

-- ----------------------------
-- Records of sys_admin_role
-- ----------------------------
BEGIN;
INSERT INTO `sys_admin_role` (`admin_id`, `role_id`) VALUES (89, 1);
INSERT INTO `sys_admin_role` (`admin_id`, `role_id`) VALUES (98, 13);
INSERT INTO `sys_admin_role` (`admin_id`, `role_id`) VALUES (101, 10);
INSERT INTO `sys_admin_role` (`admin_id`, `role_id`) VALUES (99, 12);
INSERT INTO `sys_admin_role` (`admin_id`, `role_id`) VALUES (102, 13);
COMMIT;

-- ----------------------------
-- Table structure for sys_dept
-- ----------------------------
DROP TABLE IF EXISTS `sys_dept`;
CREATE TABLE `sys_dept` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
  `parent_id` int NOT NULL COMMENT '父id',
  `dept_type` int NOT NULL COMMENT '部门类型（1->公司, 2->中心，3->部门）',
  `dept_name` varchar(30) CHARACTER SET utf8mb3 COLLATE utf8mb3_bin NOT NULL COMMENT '部门名称',
  `dept_status` int NOT NULL DEFAULT '1' COMMENT '部门状态（1->正常 2->停用）',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `dept_name` (`dept_name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin ROW_FORMAT=DYNAMIC COMMENT='部门表';

-- ----------------------------
-- Records of sys_dept
-- ----------------------------
BEGIN;
INSERT INTO `sys_dept` (`id`, `parent_id`, `dept_type`, `dept_name`, `dept_status`, `create_time`) VALUES (1, 0, 1, '神舟科技有限公司', 1, '2023-06-14 17:53:23');
INSERT INTO `sys_dept` (`id`, `parent_id`, `dept_type`, `dept_name`, `dept_status`, `create_time`) VALUES (2, 1, 2, '深圳研发中心', 1, '2023-06-14 17:53:55');
INSERT INTO `sys_dept` (`id`, `parent_id`, `dept_type`, `dept_name`, `dept_status`, `create_time`) VALUES (3, 2, 3, '架构设计部门', 1, '2023-06-14 17:54:15');
INSERT INTO `sys_dept` (`id`, `parent_id`, `dept_type`, `dept_name`, `dept_status`, `create_time`) VALUES (5, 2, 3, '后端研发部门', 1, '2023-06-14 17:55:25');
INSERT INTO `sys_dept` (`id`, `parent_id`, `dept_type`, `dept_name`, `dept_status`, `create_time`) VALUES (6, 2, 3, '系统测试部门', 1, '2023-06-14 17:55:31');
INSERT INTO `sys_dept` (`id`, `parent_id`, `dept_type`, `dept_name`, `dept_status`, `create_time`) VALUES (12, 1, 2, '北京研发中心', 1, '2025-06-28 23:42:46');
INSERT INTO `sys_dept` (`id`, `parent_id`, `dept_type`, `dept_name`, `dept_status`, `create_time`) VALUES (13, 1, 2, '重庆研发中心', 1, '2025-06-28 23:43:15');
INSERT INTO `sys_dept` (`id`, `parent_id`, `dept_type`, `dept_name`, `dept_status`, `create_time`) VALUES (14, 12, 3, '运维1部', 1, '2025-06-28 23:43:34');
INSERT INTO `sys_dept` (`id`, `parent_id`, `dept_type`, `dept_name`, `dept_status`, `create_time`) VALUES (15, 13, 3, '运维2部', 1, '2025-06-28 23:44:15');
INSERT INTO `sys_dept` (`id`, `parent_id`, `dept_type`, `dept_name`, `dept_status`, `create_time`) VALUES (16, 13, 3, '重庆研发部-001', 1, '2025-07-04 13:10:58');
COMMIT;

-- ----------------------------
-- Table structure for sys_login_info
-- ----------------------------
DROP TABLE IF EXISTS `sys_login_info`;
CREATE TABLE `sys_login_info` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT '' COMMENT '用户账号',
  `ip_address` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT '' COMMENT '登录IP地址',
  `login_location` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT '' COMMENT '登录地点',
  `browser` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT '' COMMENT '浏览器类型',
  `os` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT '' COMMENT '操作系统',
  `login_status` int DEFAULT NULL COMMENT '登录状态（1-成功 2-失败）',
  `message` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT '' COMMENT '提示消息',
  `login_time` datetime DEFAULT NULL COMMENT '访问时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=305 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC COMMENT='登录日志记录';

-- ----------------------------
-- Records of sys_login_info
-- ----------------------------
BEGIN;
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (1, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-06-26 15:10:07');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (2, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 2, '验证码不正确', '2025-06-27 10:11:15');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (3, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-06-27 10:11:22');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (4, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-06-27 10:35:36');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (5, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-06-27 10:39:36');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (6, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 2, '验证码已过期', '2025-06-27 10:51:17');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (7, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-06-27 10:51:22');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (8, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-06-27 11:06:40');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (9, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-06-27 13:56:34');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (10, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-06-27 14:56:28');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (11, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-06-27 19:19:50');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (12, 'admin', '192.168.3.40', '局域网', 'Chrome/137.0.0.0', 'Windows 10', 2, '密码不正确', '2025-06-28 14:29:35');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (13, 'admin', '192.168.3.40', '局域网', 'Chrome/137.0.0.0', 'Windows 10', 1, '登录成功', '2025-06-28 14:29:40');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (14, 'admin', '192.168.3.40', '局域网', 'Chrome/137.0.0.0', 'Windows 10', 1, '登录成功', '2025-06-28 15:57:37');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (15, 'admin', '192.168.3.40', '局域网', 'Chrome/137.0.0.0', 'Windows 10', 1, '登录成功', '2025-06-28 15:58:50');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (16, 'admin', '192.168.3.40', '局域网', 'Chrome/137.0.0.0', 'Windows 10', 1, '登录成功', '2025-06-28 16:01:59');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (17, 'admin', '192.168.3.40', '局域网', 'Chrome/137.0.0.0', 'Windows 10', 1, '登录成功', '2025-06-28 19:06:18');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (18, 'admin', '192.168.3.40', '局域网', 'Chrome/137.0.0.0', 'Windows 10', 1, '登录成功', '2025-06-28 23:36:22');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (23, 'lisi', '192.168.3.40', '局域网', 'Chrome/137.0.0.0', 'Windows 10', 1, '登录成功', '2025-06-29 15:23:15');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (24, 'test', '192.168.3.40', '局域网', 'Chrome/137.0.0.0', 'Windows 10', 2, '密码不正确', '2025-06-29 15:23:55');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (25, 'test', '192.168.3.40', '局域网', 'Chrome/137.0.0.0', 'Windows 10', 1, '登录成功', '2025-06-29 15:24:00');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (26, 'admin', '192.168.3.40', '局域网', 'Chrome/137.0.0.0', 'Windows 10', 1, '登录成功', '2025-06-29 15:24:42');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (27, 'lisi', '192.168.3.40', '局域网', 'Chrome/137.0.0.0', 'Windows 10', 1, '登录成功', '2025-06-29 15:25:26');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (28, 'admin', '192.168.3.40', '局域网', 'Chrome/137.0.0.0', 'Windows 10', 2, '密码不正确', '2025-06-29 15:26:54');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (31, 'zhangfan', '192.168.3.40', '局域网', 'Chrome/137.0.0.0', 'Windows 10', 1, '登录成功', '2025-06-29 16:16:15');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (32, 'admin', '192.168.3.40', '局域网', 'Chrome/137.0.0.0', 'Windows 10', 1, '登录成功', '2025-06-30 00:06:13');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (33, 'admin', '192.168.3.40', '局域网', 'Chrome/137.0.0.0', 'Windows 10', 1, '登录成功', '2025-06-30 00:53:23');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (34, 'lisi', '192.168.3.40', '局域网', 'Chrome/137.0.0.0', 'Windows 10', 1, '登录成功', '2025-06-30 00:53:59');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (35, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-06-30 10:28:30');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (36, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-01 10:30:54');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (37, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-02 11:58:08');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (38, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-03 11:59:20');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (39, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-03 14:12:11');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (40, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-03 14:27:21');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (41, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-03 16:53:19');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (42, 'zhangfan', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 2, '密码不正确', '2025-07-03 17:26:51');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (43, 'zhangfan', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 2, '密码不正确', '2025-07-03 17:26:57');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (44, 'zhangfan', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 2, '密码不正确', '2025-07-03 17:27:02');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (45, 'zhangfan', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 2, '密码不正确', '2025-07-03 17:27:11');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (46, 'zhangsan', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 2, '密码不正确', '2025-07-03 17:27:18');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (47, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-03 17:27:29');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (48, 'zhangfan', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-03 17:28:28');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (49, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-03 17:29:34');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (50, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-03 17:47:19');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (51, 'zhangfan', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-03 17:52:11');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (52, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-03 17:54:00');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (53, 'lisi', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-03 18:54:26');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (54, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-03 18:55:35');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (55, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-04 13:01:03');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (56, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-04 13:08:12');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (57, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-04 16:46:54');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (58, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-04 17:14:45');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (59, 'zhangfan', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-04 17:54:44');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (60, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-04 17:55:09');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (61, 'zhangfan', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-04 17:55:39');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (62, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-07 15:30:48');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (63, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-08 17:57:42');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (64, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-09 11:38:15');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (65, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-10 10:39:06');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (66, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-10 13:07:58');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (67, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-10 14:22:56');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (68, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 2, '验证码已过期', '2025-07-10 14:52:40');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (69, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 2, '验证码已过期', '2025-07-10 14:52:46');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (70, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-10 14:52:52');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (71, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-14 11:40:24');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (72, 'admin', '127.0.0.1', '服务器登录', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-15 11:47:28');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (73, 'admin', '127.0.0.1', '服务器登录', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-15 14:24:52');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (74, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-16 14:36:29');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (75, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-16 15:51:44');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (76, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-17 10:55:44');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (78, 'admin', '10.7.16.22', '局域网', 'Chrome/137.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-17 20:45:35');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (79, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-18 14:46:12');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (80, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-18 14:49:31');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (81, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-18 15:07:07');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (82, 'zhangsan', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-18 15:32:43');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (83, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-18 15:33:28');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (84, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-18 16:59:52');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (85, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-22 10:20:45');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (86, 'zhangfan', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Mac OS X 10_15_7', 2, '密码不正确', '2025-07-22 10:29:25');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (87, 'zhangfan', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Mac OS X 10_15_7', 2, '密码不正确', '2025-07-22 10:29:36');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (88, 'zhangsan', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-22 10:29:45');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (89, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-22 10:30:24');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (90, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-22 14:22:55');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (91, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-22 16:09:14');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (92, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Windows 10', 1, '登录成功', '2025-07-22 22:45:37');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (93, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-23 10:26:42');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (94, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-24 10:28:20');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (95, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Windows 10', 1, '登录成功', '2025-07-24 14:39:29');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (96, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Windows 10', 1, '登录成功', '2025-07-24 14:40:05');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (97, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Windows 10', 1, '登录成功', '2025-07-24 14:42:30');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (98, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-24 15:21:36');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (99, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-24 19:42:47');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (130, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-31 12:41:48');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (131, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-31 12:45:40');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (132, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-31 12:46:34');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (133, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-31 13:06:53');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (134, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-31 13:20:24');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (135, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-31 13:21:43');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (136, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Mac OS X 10_15_7', 2, '验证码不正确', '2025-07-31 19:45:11');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (137, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-31 19:45:19');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (138, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-07-31 21:19:11');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (139, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-01 21:53:34');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (140, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-02 00:16:46');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (141, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-02 16:41:00');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (142, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-03 16:51:00');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (143, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-04 10:13:37');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (144, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-04 19:35:28');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (145, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Windows 10', 1, '登录成功', '2025-08-04 22:34:40');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (146, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-05 09:59:29');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (147, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-05 14:40:59');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (148, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-05 14:41:09');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (149, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-05 14:41:28');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (150, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Windows 10', 1, '登录成功', '2025-08-05 23:29:26');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (151, 'admin', '127.0.0.1', '服务器登录', 'Chrome/139.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-06 13:37:40');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (152, 'admin', '127.0.0.1', '服务器登录', 'Chrome/139.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-06 13:39:02');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (153, 'admin', '127.0.0.1', '服务器登录', 'Chrome/139.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-06 14:37:11');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (154, 'admin', '127.0.0.1', '服务器登录', 'Chrome/139.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-06 16:59:05');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (155, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Windows 10', 1, '登录成功', '2025-08-06 22:34:58');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (156, 'admin', '127.0.0.1', '服务器登录', 'Chrome/139.0.0.0', 'Mac OS X 10_15_7', 2, '验证码不正确', '2025-08-07 14:54:35');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (157, 'admin', '127.0.0.1', '服务器登录', 'Chrome/139.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-07 14:54:40');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (158, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Windows 10', 1, '登录成功', '2025-08-07 22:48:24');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (159, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Windows 10', 1, '登录成功', '2025-08-07 23:13:30');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (160, 'admin', '127.0.0.1', '服务器登录', 'Chrome/139.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-08 15:04:01');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (161, 'admin', '127.0.0.1', '服务器登录', 'Chrome/139.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-08 23:48:29');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (162, 'admin', '127.0.0.1', '服务器登录', 'Chrome/139.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-09 12:59:01');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (163, 'admin', '127.0.0.1', '服务器登录', 'Chrome/139.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-10 13:13:50');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (164, 'admin', '127.0.0.1', '服务器登录', 'Chrome/139.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-11 11:14:05');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (165, 'admin', '127.0.0.1', '服务器登录', 'Chrome/139.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-12 11:16:21');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (166, 'admin', '127.0.0.1', '服务器登录', 'Chrome/139.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-13 10:33:05');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (167, 'admin', '127.0.0.1', '服务器登录', 'Chrome/139.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-14 10:36:39');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (168, 'admin', '127.0.0.1', '服务器登录', 'Chrome/139.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-14 21:41:10');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (169, 'admin', '127.0.0.1', '服务器登录', 'Chrome/139.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-15 10:57:16');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (170, 'admin', '127.0.0.1', '服务器登录', 'Chrome/139.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-17 18:12:31');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (171, 'admin', '127.0.0.1', '服务器登录', 'Chrome/139.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-21 11:33:26');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (172, 'admin', '127.0.0.1', '服务器登录', 'Chrome/139.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-21 11:39:58');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (173, 'admin', '127.0.0.1', '服务器登录', 'Chrome/139.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-22 10:41:13');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (174, 'admin', '127.0.0.1', '服务器登录', 'Chrome/139.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-23 12:35:49');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (175, 'admin', '127.0.0.1', '服务器登录', 'Chrome/139.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-23 13:47:27');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (176, 'admin', '127.0.0.1', '服务器登录', 'Chrome/139.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-23 18:36:52');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (177, 'admin', '127.0.0.1', '服务器登录', 'Chrome/139.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-23 22:43:09');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (178, 'admin', '127.0.0.1', '服务器登录', 'Chrome/139.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-23 22:43:11');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (179, 'admin', '127.0.0.1', '服务器登录', 'Chrome/139.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-24 15:27:16');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (180, 'admin', '127.0.0.1', '服务器登录', 'Chrome/139.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-25 10:15:58');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (182, 'admin', '127.0.0.1', '服务器登录', 'Chrome/139.0.0.0', 'Mac OS X 10_15_7', 2, '验证码已过期', '2025-08-25 15:19:48');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (183, 'admin', '127.0.0.1', '服务器登录', 'Chrome/139.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-25 15:19:54');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (184, 'admin', '127.0.0.1', '服务器登录', 'Chrome/139.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-26 19:59:49');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (185, 'admin', '127.0.0.1', '服务器登录', 'Chrome/139.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-28 10:17:30');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (186, 'admin', '127.0.0.1', '服务器登录', 'Chrome/139.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-29 10:28:08');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (187, 'admin', '127.0.0.1', '服务器登录', 'Chrome/139.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-29 11:11:20');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (188, 'admin', '127.0.0.1', '服务器登录', 'Chrome/139.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-08-29 11:30:08');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (189, 'admin', '127.0.0.1', '服务器登录', 'Chrome/139.0.0.0', 'Mac OS X 10_15_7', 2, '验证码不正确', '2025-09-04 20:10:00');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (190, 'admin', '127.0.0.1', '服务器登录', 'Chrome/139.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-04 20:10:06');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (191, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Windows 10', 2, '密码不正确', '2025-09-04 23:21:57');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (192, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Windows 10', 1, '登录成功', '2025-09-04 23:22:56');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (193, 'admin', '127.0.0.1', '服务器登录', 'Chrome/139.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-05 20:11:51');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (194, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-06 14:59:37');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (195, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-06 15:30:20');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (196, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-07 16:02:46');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (197, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-08 10:26:00');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (198, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-08 13:28:00');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (199, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-09 15:05:23');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (200, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-10 11:10:32');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (201, 'admin', '127.0.0.1', '服务器登录', 'Chrome/138.0.0.0', 'Windows 10', 1, '登录成功', '2025-09-10 22:47:46');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (202, 'zhangfan', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 2, '密码不正确', '2025-09-11 10:43:54');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (203, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-11 10:44:08');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (204, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-11 12:29:57');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (205, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-11 15:06:41');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (206, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-11 15:09:15');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (207, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-11 20:11:48');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (208, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-11 20:49:30');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (209, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-11 21:03:03');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (210, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-12 10:23:06');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (211, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-12 10:33:02');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (212, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-12 11:20:47');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (213, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 2, '验证码已过期', '2025-09-12 17:52:39');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (214, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-12 17:52:46');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (215, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-12 23:13:28');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (216, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-13 15:21:53');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (217, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-14 15:26:35');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (218, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-14 23:19:13');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (219, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-15 10:47:58');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (220, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-15 15:19:15');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (221, 'root', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 2, '密码不正确', '2025-09-15 20:46:53');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (222, 'root', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 2, '密码不正确', '2025-09-15 20:47:00');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (223, 'root', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 2, '密码不正确', '2025-09-15 20:47:01');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (224, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-15 20:47:10');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (225, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-16 09:54:18');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (226, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-16 09:55:10');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (227, 'root', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 2, '密码不正确', '2025-09-16 17:13:49');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (228, 'root', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 2, '密码不正确', '2025-09-16 17:14:00');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (229, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-16 17:14:06');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (230, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-16 17:20:00');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (231, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-16 17:20:58');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (232, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-17 10:08:32');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (233, 'zhangsan', '127.0.0.1', '服务器登录', 'Chrome/98.0.4758.139', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-17 19:13:09');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (234, 'zhangsan', '127.0.0.1', '服务器登录', 'Safari/605.1.15', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-17 19:14:56');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (235, 'zhangsan', '127.0.0.1', '服务器登录', 'Chrome/98.0.4758.139', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-17 19:23:53');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (237, 'zhangsan', '127.0.0.1', '服务器登录', 'Chrome/98.0.4758.139', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-17 19:30:10');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (238, 'zhangsan', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-17 19:31:29');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (241, 'zhangsan', '127.0.0.1', '服务器登录', 'Chrome/98.0.4758.139', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-17 19:33:58');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (242, 'zhangsan', '127.0.0.1', '服务器登录', 'Chrome/98.0.4758.139', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-17 19:36:05');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (243, 'zhangsan', '127.0.0.1', '服务器登录', 'Chrome/98.0.4758.139', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-17 19:56:06');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (244, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-17 20:15:05');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (246, 'zhangsan', '127.0.0.1', '服务器登录', 'Chrome/98.0.4758.139', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-17 20:24:43');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (247, 'zhangsan', '127.0.0.1', '服务器登录', 'Chrome/98.0.4758.139', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-17 20:40:16');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (248, 'zhangsan', '127.0.0.1', '服务器登录', 'Chrome/98.0.4758.139', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-17 20:41:45');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (249, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-17 20:49:04');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (250, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-17 20:52:15');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (251, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 2, '验证码不正确', '2025-09-17 20:58:27');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (252, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-17 20:58:35');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (253, 'zhangsan', '127.0.0.1', '服务器登录', 'Chrome/98.0.4758.139', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-17 21:01:53');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (254, 'zhangsan', '127.0.0.1', '服务器登录', 'Chrome/98.0.4758.139', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-17 21:11:34');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (255, 'zhangsan', '127.0.0.1', '服务器登录', 'Chrome/98.0.4758.139', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-18 11:03:24');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (256, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-18 11:40:50');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (257, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-18 15:16:51');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (258, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-18 16:03:35');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (259, 'zhangsan', '127.0.0.1', '服务器登录', 'Chrome/98.0.4758.139', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-18 18:30:05');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (260, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-19 21:29:30');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (261, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-19 22:46:19');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (262, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-19 22:50:23');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (263, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-20 02:01:52');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (264, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-21 02:16:14');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (265, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-21 15:13:08');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (266, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-21 15:59:20');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (267, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-22 11:44:03');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (268, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-23 11:59:48');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (269, 'admin', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-23 15:22:41');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (270, 'admin', '127.0.0.1', '服务器登录', 'Chrome/98.0.4758.139', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-24 12:50:17');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (271, 'test', '127.0.0.1', '服务器登录', 'Chrome/98.0.4758.139', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-24 12:50:34');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (272, 'test', '127.0.0.1', '服务器登录', 'Chrome/98.0.4758.139', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-24 12:51:31');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (273, 'test', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-24 12:55:58');
INSERT INTO `sys_login_info` (`id`, `username`, `ip_address`, `login_location`, `browser`, `os`, `login_status`, `message`, `login_time`) VALUES (274, 'test', '127.0.0.1', '服务器登录', 'Chrome/140.0.0.0', 'Mac OS X 10_15_7', 1, '登录成功', '2025-09-24 13:23:57');
COMMIT;

-- ----------------------------
-- Table structure for sys_menu
-- ----------------------------
DROP TABLE IF EXISTS `sys_menu`;
CREATE TABLE `sys_menu` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
  `parent_id` int DEFAULT NULL COMMENT '父级菜单id',
  `menu_name` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '菜单名称',
  `icon` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '图标',
  `value` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '接口权限值',
  `menu_type` int DEFAULT NULL COMMENT '菜单类型：1->目录；2->菜单；3->按钮（接口绑定权限）',
  `url` varchar(200) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '菜单url',
  `menu_status` int DEFAULT '2' COMMENT '启用状态；1->禁用；2->启用',
  `sort` int DEFAULT NULL COMMENT '排序',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=207 DEFAULT CHARSET=utf8mb3 ROW_FORMAT=DYNAMIC COMMENT='菜单表';

-- ----------------------------
-- Records of sys_menu
-- ----------------------------
BEGIN;
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (4, 0, '系统管理', 'StarFilled', '', 1, '', 2, 7, '2022-09-04 13:57:39');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (6, 4, '用户信息', 'Avatar', 'base:admin:list', 2, 'system/admin', 2, 1, '2022-09-04 13:59:39');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (7, 4, '角色信息', 'InfoFilled', 'base:role:list', 2, 'system/role', 2, 2, '2022-09-04 14:00:12');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (8, 4, '菜单信息', 'Histogram', 'base:menu:list', 2, 'system/menu', 2, 3, '2022-09-04 14:00:17');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (9, 4, '部门信息', 'Menu', 'base:dept:list', 2, 'system/dept', 2, 4, '2022-09-04 14:01:58');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (10, 4, '岗位信息', 'Promotion', 'base:post:list', 2, 'system/post', 2, 5, '2022-09-04 14:02:06');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (16, 6, '新增用户', '', 'base:admin:add', 3, '', 2, 1, '2022-09-04 18:32:55');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (17, 6, '修改用户', '', 'base:admin:edit', 3, '', 2, 2, '2022-09-04 18:33:29');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (18, 6, '删除用户', '', 'base:admin:delete', 3, '', 2, 3, '2022-09-04 18:33:51');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (21, 7, '新增角色', '', 'base:role:add', 3, '', 2, 1, '2022-09-04 18:44:47');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (22, 7, '修改角色', '', 'base:role:edit', 3, '', 2, 2, '2022-09-04 18:45:10');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (23, 7, '删除角色', '', 'base:role:delete', 3, '', 2, 3, '2022-09-04 18:45:46');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (24, 7, '分配权限', '', 'base:role:assign', 3, '', 2, 4, '2022-09-04 18:46:20');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (26, 8, '新增菜单', '', 'base:menu:add', 3, '', 2, 1, '2022-09-04 18:49:51');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (27, 8, '修改菜单', '', 'base:menu:edit', 3, '', 2, 2, '2022-09-04 18:50:24');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (28, 8, '删除菜单', '', 'base:menu:delete', 3, '', 2, 3, '2022-09-04 18:50:53');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (29, 9, '新增部门', '', 'base:dept:add', 3, '', 2, 1, '2022-09-04 18:52:16');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (30, 9, '修改部门', '', 'base:dept:edit', 3, '', 2, 2, '2022-09-04 18:52:37');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (31, 9, '删除部门', '', 'base:dept:delete', 3, '', 2, 3, '2022-09-04 18:52:50');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (32, 10, '新增岗位', '', 'base:post:add', 3, '', 2, 1, '2022-09-04 18:53:28');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (33, 10, '修改岗位', '', 'base:post:edit', 3, '', 2, 2, '2022-09-04 18:53:48');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (34, 10, '删除岗位', '', 'base:post:delete', 3, '', 2, 3, '2022-09-04 18:54:00');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (44, 0, '操作审计', 'BellFilled', '', 1, '', 2, 9, '2022-09-05 11:06:57');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (45, 44, '操作日志', 'User', 'monitor:operator:list', 2, 'monitor/operator', 2, 1, '2022-09-05 11:10:54');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (46, 44, '登录日志', 'DocumentRemove', 'monitor:loginLog:list', 2, 'monitor/loginlog', 2, 2, '2022-09-05 11:11:31');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (47, 45, '清空操作日志', '', 'monitor:operator:clean', 3, '', 2, 1, '2022-09-05 11:12:36');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (49, 46, '清空登录日志', '', 'monitor:loginLog:clean', 3, '', 2, 1, '2022-09-05 11:16:01');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (60, 6, '重置密码', NULL, 'base:admin:reset', 3, NULL, 2, 6, '2022-12-01 16:33:34');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (62, 46, '删除登录日志', '', 'monitor:loginLog:delete', 3, '', 2, 2, '2022-12-02 15:41:56');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (72, 0, '仪表盘', 'HomeFilled', '', 1, 'dashboard', 2, 1, '2023-05-24 22:11:13');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (73, 45, '删除操作日志', '', 'monitor:operator:delete', 3, '', 2, 3, '2023-06-02 10:09:38');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (78, 80, '主机管理', 'Platform', 'cmdb:ecs:list', 2, 'cmdb/ecs', 2, 1, '2025-06-29 00:30:35');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (80, 0, '资产管理', 'TrendCharts', '', 1, '', 2, 2, '2025-07-03 11:47:07');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (81, 0, '容器管理', 'UploadFilled', '', 1, '', 2, 3, '2025-07-03 11:50:47');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (82, 81, '集群管理', 'Menu', 'cloud:k8s:list', 2, 'k8s/list', 2, 1, '2025-07-03 11:56:44');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (83, 81, '节点管理', 'Help', 'cloud:k8s:node', 2, 'k8s/node', 2, 2, '2025-07-03 12:04:59');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (84, 0, '配置中心', 'Tools', '', 1, '', 2, 8, '2025-07-04 17:00:01');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (85, 84, '主机凭据', 'Setting', 'config:ecs:key', 2, 'config/ecskey', 2, 1, '2025-07-04 17:03:10');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (86, 84, '通用凭据', 'User', 'config:accountauth:key', 2, 'config/accountauth', 2, 2, '2025-07-04 17:08:20');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (88, 80, '业务分组', 'Shop', 'cmdb:group', 2, 'cmdb/group', 2, 3, '2025-07-16 15:17:14');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (89, 88, '创建分组', '', 'cmdb:group:add', 3, '', 2, 1, '2025-07-18 15:24:31');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (90, 88, '修改分组', '', 'cmdb:group:update', 3, '', 2, 2, '2025-07-18 15:25:49');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (91, 88, '删除分组', '', 'cmdb:group:delete', 3, '', 2, 3, '2025-07-18 15:26:21');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (93, 81, '工作负载', 'Star', 'cloud:k8s:workload', 2, 'k8s/workload', 2, 4, '2025-07-24 14:38:31');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (95, 80, '数据管理', 'Coin', 'cmdb:db', 2, 'cmdb/db', 2, 2, '2025-07-29 15:27:50');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (96, 44, '数据日志', 'Coin', 'monitor:dblog:list', 2, 'monitor/dblog', 2, 3, '2025-07-31 12:44:07');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (97, 0, '任务中心', 'User', '', 1, '', 2, 5, '2025-08-06 13:33:47');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (98, 97, '任务模版', 'connection', 'task:template', 2, 'task/template', 2, 2, '2025-08-06 13:35:19');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (99, 97, '任务作业', 'key', 'task:job', 2, 'task/job', 2, 1, '2025-08-06 13:36:06');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (100, 97, 'Ansible任务', 'Eleme', 'task:ansible', 2, 'task/ansible', 2, 3, '2025-08-23 18:35:24');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (101, 0, '运维工具', 'Search', '', 1, '', 2, 6, '2025-08-29 10:59:35');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (102, 101, 'agent列表', 'price-tag', 'ops:agent', 2, 'ops/agent', 2, 2, '2025-08-29 11:22:20');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (103, 101, '工具列表', 'present', 'ops:tools', 2, 'ops/tools', 2, 1, '2025-08-29 11:29:02');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (104, 84, '密钥管理', 'Phone', 'config:keymanage:key', 2, 'config/keymanage', 2, 3, '2025-09-08 13:24:40');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (105, 81, '命名空间', 'discount', 'k8s:namespace', 2, 'k8s/namespace', 2, 3, '2025-09-11 15:02:14');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (106, 81, '网络管理', 'guide', 'k8s:network', 2, 'k8s/network', 2, 5, '2025-09-11 15:04:14');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (107, 81, '配置管理', 'connection', 'k8s:config', 2, 'k8s/config', 2, 7, '2025-09-11 15:04:52');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (108, 81, '存储管理', 'Coin', 'k8s:storage', 2, 'k8s/storage', 2, 6, '2025-09-11 15:05:40');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (109, 0, '服务管理', 'UploadFilled', '', 1, '', 2, 4, '2025-09-16 09:49:55');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (110, 109, '应用列表', 'Menu', 'app:application', 2, 'app/application', 2, 1, '2025-09-16 09:52:58');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (111, 109, '快速发布', 'TrendCharts', 'app:quick-release', 2, 'app/quick-release', 2, 2, '2025-09-16 17:12:11');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (113, 45, '批量删除', '', 'monitor:operator:delete', 3, '', 2, 2, '2025-09-17 20:55:13');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (114, 104, '修改密钥', '', 'config:keymanage:edit', 3, '', 2, 1, '2025-09-18 10:45:57');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (115, 104, '删除密钥', '', 'config:keymanage:delete', 3, '', 2, 2, '2025-09-18 10:53:44');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (117, 104, '同步主机', '', 'config:keymanage:rsync', 3, '', 2, 3, '2025-09-18 10:57:25');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (118, 104, '创建密钥', '', 'config:keymanage:create', 3, '', 2, 4, '2025-09-18 11:01:12');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (119, 86, '修改账号', '', 'config:common:edit', 3, '', 2, 2, '2025-09-18 11:47:33');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (120, 86, '删除账号', '', 'config:common:delete', 3, '', 2, 3, '2025-09-18 11:48:17');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (121, 86, '解密账号', '', 'config:common:decrypt', 3, '', 2, 4, '2025-09-18 11:48:57');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (122, 86, '创建账号', '', 'config:common:add', 3, '', 2, 1, '2025-09-18 11:49:30');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (123, 85, '修改凭据', '', 'config:ecs:edit', 3, '', 2, 1, '2025-09-18 11:54:16');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (124, 85, '删除凭据', '', 'config:ecs:delete', 3, '', 2, 2, '2025-09-18 11:54:51');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (125, 85, '创建凭据', '', 'config:ecs:create', 3, '', 2, 3, '2025-09-18 11:55:21');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (126, 102, '卸载agent', '', 'ops:agent:delete', 3, '', 2, 1, '2025-09-18 12:47:48');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (127, 102, '查看agent', '', 'ops:agent:get', 3, '', 2, 2, '2025-09-18 12:49:02');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (128, 102, '部署agent', '', 'ops:agent:create', 3, '', 2, 3, '2025-09-18 12:49:56');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (129, 102, '批量卸载agent', '', 'ops:agent:deleteall', 3, '', 2, 4, '2025-09-18 12:50:52');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (130, 100, '启动ansible任务流程', '', 'task:ansible:start', 3, '', 2, 1, '2025-09-18 12:59:30');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (131, 100, '删除ansible任务', '', 'task:ansible:delete', 3, '', 2, 2, '2025-09-18 13:00:03');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (132, 100, '新增ansible任务', '', 'task:ansible:create', 3, '', 2, 3, '2025-09-18 13:00:45');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (133, 98, '新增模版', '', 'task:template:add', 3, '', 2, 1, '2025-09-18 13:16:38');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (134, 98, '修改模版', '', 'task:template:edit', 3, '', 2, 2, '2025-09-18 13:17:04');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (135, 98, '删除模版', '', 'task:template:delete', 3, '', 2, 3, '2025-09-18 13:18:25');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (136, 99, '新增任务', '', 'task:job:add', 3, '', 2, 1, '2025-09-18 13:24:19');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (137, 99, '启动任务', '', 'task:job:start', 3, '', 2, 2, '2025-09-18 13:24:59');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (138, 99, '删除任务', '', 'task:job:delete', 3, '', 2, 3, '2025-09-18 13:25:41');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (139, 111, '新建发布', '', 'app:quick-release:add', 3, '', 2, 1, '2025-09-18 13:30:53');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (140, 111, '启动发布', '', 'app:quick-release:start', 3, '', 2, 2, '2025-09-18 13:32:11');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (141, 111, '删除发布', '', 'app:quick-release:delete', 3, '', 2, 3, '2025-09-18 13:32:32');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (142, 110, '创建应用', '', 'app:application:add', 3, '', 2, 1, '2025-09-18 14:28:07');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (143, 110, '修改应用', '', 'app:application:edit', 3, '', 2, 2, '2025-09-18 14:28:59');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (144, 110, '配置应用环境', '', 'app:application:env', 3, '', 2, 3, '2025-09-18 14:29:34');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (145, 110, '删除应用', '', 'app:application:delete', 3, '', 2, 4, '2025-09-18 14:30:11');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (146, 95, '创建数据库账号', '', 'cmdb:db:add', 3, '', 2, 1, '2025-09-18 14:41:32');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (147, 95, '修改数据库配置', '', 'cmdb:db:edit', 3, '', 2, 2, '2025-09-18 14:42:47');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (148, 95, '删除数据库账号', '', 'cmdb:db:delete', 3, '', 2, 3, '2025-09-18 14:43:57');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (149, 78, '创建主机', '', 'cmdb:ecs:add', 3, '', 2, 1, '2025-09-18 14:47:42');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (150, 78, '主机终端', '', 'cmdb:ecs:terminal', 3, '', 2, 2, '2025-09-18 14:48:36');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (151, 78, '修改主机信息', '', 'cmdb:ecs:edit', 3, '', 2, 3, '2025-09-18 14:49:43');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (152, 78, '上传文件到主机', '', 'cmdb:ecs:upload', 3, '', 2, 4, '2025-09-18 14:50:38');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (153, 78, '执行主机命令', '', 'cmdb:ecs:shell', 3, '', 2, 5, '2025-09-18 14:51:10');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (154, 78, '监控主机', '', 'cmdb:ecs:monitor', 3, '', 2, 6, '2025-09-18 14:51:52');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (155, 78, '删除主机', '', 'cmdb:ecs:delete', 3, '', 2, 7, '2025-09-18 14:52:20');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (156, 99, '启动脚本', '', 'task:job:jobstart', 3, '', 2, 4, '2025-09-18 18:36:38');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (157, 99, '停止脚本', '', 'task:job:jobstop', 3, '', 2, 5, '2025-09-18 18:39:23');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (159, 100, '开始ansible任务作业', '', 'task:ansible:jobstart', 3, '', 2, 4, '2025-09-18 18:43:40');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (160, 111, '启动jenkins任务', '', 'app:quick-release:jobstart', 3, '', 2, 4, '2025-09-18 18:47:39');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (161, 111, '停止jenkins任务', '', 'app:quick-release:jobstop', 3, '', 2, 5, '2025-09-18 18:48:16');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (162, 110, '新增环境', '', 'app:application:envadd', 3, '', 2, 5, '2025-09-18 21:02:28');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (163, 110, '修改环境', '', 'app:application:envedit', 3, '', 2, 6, '2025-09-18 21:03:08');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (164, 110, '删除环境', '', 'app:application:envdelete', 3, '', 2, 7, '2025-09-18 21:04:43');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (165, 78, '连接主机终端', '', 'cmdb:ecs:connecthost', 3, '', 2, 8, '2025-09-18 21:11:43');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (166, 78, '同步主机信息', '', 'cmdb:ecs:rsync', 3, '', 2, 9, '2025-09-19 21:35:06');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (167, 82, '注册集群', '', 'cloud:k8s:register', 3, '', 2, 1, '2025-09-19 21:57:54');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (168, 82, '创建集群', '', 'cloud:k8s:add', 3, '', 2, 2, '2025-09-19 21:58:19');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (169, 82, '修改集群', '', 'cloud:k8s:edit', 3, '', 2, 3, '2025-09-19 21:59:06');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (170, 82, '同步集群', '', 'cloud:k8s:rsync', 3, '', 2, 4, '2025-09-19 21:59:31');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (171, 82, '删除集群', '', 'cloud:k8s:delete', 3, '', 2, 5, '2025-09-19 21:59:56');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (172, 83, '查看监控仪表盘', '', 'k8s:node:monitor', 3, '', 2, 1, '2025-09-20 00:19:49');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (173, 83, '查看节点资源详情', '', 'k8s:node:details', 3, '', 2, 2, '2025-09-20 00:21:20');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (174, 83, '配置节点污点', '', 'k8s:node:stain', 3, '', 2, 3, '2025-09-20 00:22:17');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (175, 83, '增加标签', '', 'k8s:node:label', 3, '', 2, 4, '2025-09-20 00:23:15');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (176, 83, '封锁节点', '', 'k8s:node:close', 3, '', 2, 5, '2025-09-20 00:24:13');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (177, 83, '驱逐节点', '', 'k8s:node:expel', 3, '', 2, 6, '2025-09-20 00:25:04');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (178, 105, '创建命名空间', '', 'k8s:namespace:add', 3, '', 2, 1, '2025-09-20 00:36:14');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (179, 105, '查看命名空间详情', '', 'k8s:namespace:details', 3, '', 2, 2, '2025-09-20 00:37:22');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (180, 105, '查看命名空间资源配置', '', 'k8s:namespace:setup', 3, '', 2, 3, '2025-09-20 00:39:05');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (181, 105, '添加命名空间资源配置', '', 'k8s:namespace:setupadd', 3, '', 2, 4, '2025-09-20 00:40:06');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (182, 105, '查看限制命名空间资源', '', 'k8s:namespace:restriction', 3, '', 2, 5, '2025-09-20 00:41:25');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (183, 105, '添加限制命名空间资源', '', 'k8s:namespace:restrictionadd', 3, '', 2, 6, '2025-09-20 00:42:23');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (184, 105, '删除命名空间', '', 'k8s:namespace:delete', 3, '', 2, 7, '2025-09-20 00:43:03');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (185, 93, '新增工作负载', '', 'k8s:workload:add', 3, '', 2, 1, '2025-09-20 01:05:08');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (186, 93, '伸缩pod', '', 'k8s:workload:expandable', 3, '', 2, 2, '2025-09-20 01:06:18');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (187, 93, '重启pod', '', 'k8s:workload:restart', 3, '', 2, 3, '2025-09-20 01:07:13');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (188, 93, '更新pod资源限制', '', 'k8s:workload:resource', 3, '', 2, 4, '2025-09-20 01:08:39');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (189, 93, '更新pod资调度', '', 'k8s:workload:dispatch', 3, '', 2, 5, '2025-09-20 01:09:36');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (190, 93, '更新yaml文件', '', 'k8s:workload:edityaml', 3, '', 2, 6, '2025-09-20 01:10:57');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (191, 93, '删除工作负载', '', 'k8s:workload:delete', 3, '', 2, 7, '2025-09-20 01:11:52');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (192, 93, '回滚工作负载版本', '', 'k8s:workload:rollback_version', 3, '', 2, 8, '2025-09-20 01:39:38');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (193, 93, '查看pod日志', '', 'k8s:workload:podlog', 3, '', 2, 9, '2025-09-20 01:44:37');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (194, 93, '删除pod', '', 'k8s:workload:poddelete', 3, '', 2, 10, '2025-09-20 01:45:16');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (195, 93, '登陆pod终端', '', 'k8s:workload:terminal', 3, '', 2, 11, '2025-09-20 01:46:07');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (196, 93, '编辑pod yaml文件', '', 'k8s:workload:edityaml', 3, '', 2, 12, '2025-09-20 01:47:22');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (197, 106, '新增service', '', 'k8s:network:addservice', 3, '', 2, 1, '2025-09-20 02:14:21');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (198, 106, '编辑 Service', '', 'k8s:network:editservice', 3, '', 2, 2, '2025-09-20 02:15:19');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (199, 106, '编辑service YAML', '', 'k8s:network:edit_service_yaml', 3, '', 2, 3, '2025-09-20 02:16:15');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (200, 106, '查看Service 事件', '', 'k8s:network:service_event', 3, '', 2, 4, '2025-09-20 02:18:11');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (201, 106, '删除Service', '', 'k8s:network:deleteservice', 3, '', 2, 5, '2025-09-20 02:18:59');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (202, 106, '编辑ingress', '', 'k8s:network:editingress', 3, '', 2, 6, '2025-09-20 02:26:59');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (203, 106, '新增ingress', '', 'k8s:network:addingress', 3, '', 2, 7, '2025-09-20 02:27:29');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (204, 106, '编辑ingress_yaml', '', 'k8s:network:edit_ingress_yaml', 3, '', 2, 8, '2025-09-20 02:28:23');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (205, 106, '查看ingress 事件', '', 'k8s:network:ingress_event', 3, '', 2, 9, '2025-09-20 02:29:24');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (206, 106, '删除ingress', '', 'k8s:network:delete_ingress', 3, '', 2, 10, '2025-09-20 02:30:04');
COMMIT;

-- ----------------------------
-- Table structure for sys_operation_log
-- ----------------------------
DROP TABLE IF EXISTS `sys_operation_log`;
CREATE TABLE `sys_operation_log` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
  `admin_id` bigint unsigned NOT NULL COMMENT '''管理员id''',
  `username` varchar(64) NOT NULL COMMENT '''管理员账号''',
  `method` varchar(64) NOT NULL COMMENT '''请求方式''',
  `ip` varchar(64) DEFAULT NULL COMMENT '''IP''',
  `url` varchar(500) DEFAULT NULL COMMENT '''URL''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `description` varchar(255) DEFAULT NULL COMMENT '''操作描述''',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `username` (`username`) USING BTREE,
  KEY `create_time` (`create_time`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb3 ROW_FORMAT=DYNAMIC COMMENT='操作日志记录表';

-- ----------------------------
-- Records of sys_operation_log
-- ----------------------------
BEGIN;
INSERT INTO `sys_operation_log` (`id`, `admin_id`, `username`, `method`, `ip`, `url`, `create_time`, `description`) VALUES (1, 89, 'admin', 'delete', '10.7.16.22', '/api/v1/sysOperationLog/clean', '2025-10-16 10:59:42.292', '清空操作日志');
INSERT INTO `sys_operation_log` (`id`, `admin_id`, `username`, `method`, `ip`, `url`, `create_time`, `description`) VALUES (2, 89, 'admin', 'delete', '10.7.16.22', '/api/v1/sysLoginInfo/clean', '2025-10-16 10:59:48.108', '');
INSERT INTO `sys_operation_log` (`id`, `admin_id`, `username`, `method`, `ip`, `url`, `create_time`, `description`) VALUES (3, 89, 'admin', 'delete', '10.7.16.22', '/api/v1/sysLoginInfo/clean', '2025-10-16 10:59:52.240', '');
INSERT INTO `sys_operation_log` (`id`, `admin_id`, `username`, `method`, `ip`, `url`, `create_time`, `description`) VALUES (4, 89, 'admin', 'delete', '10.7.16.22', '/api/v1/sysLoginInfo/clean', '2025-10-16 10:59:57.851', '');
INSERT INTO `sys_operation_log` (`id`, `admin_id`, `username`, `method`, `ip`, `url`, `create_time`, `description`) VALUES (5, 89, 'admin', 'delete', '10.7.16.22', '/api/v1/sysLoginInfo/clean', '2025-10-16 11:00:03.391', '');
INSERT INTO `sys_operation_log` (`id`, `admin_id`, `username`, `method`, `ip`, `url`, `create_time`, `description`) VALUES (6, 89, 'admin', 'delete', '10.7.16.22', '/api/v1/sysLoginInfo/batch/delete', '2025-10-16 11:00:08.393', '');
INSERT INTO `sys_operation_log` (`id`, `admin_id`, `username`, `method`, `ip`, `url`, `create_time`, `description`) VALUES (7, 89, 'admin', 'delete', '10.7.16.22', '/api/v1/sysLoginInfo/batch/delete', '2025-10-16 11:00:14.822', '');
INSERT INTO `sys_operation_log` (`id`, `admin_id`, `username`, `method`, `ip`, `url`, `create_time`, `description`) VALUES (8, 89, 'admin', 'delete', '10.7.16.22', '/api/v1/sysLoginInfo/batch/delete', '2025-10-16 11:00:19.155', '');
COMMIT;

-- ----------------------------
-- Table structure for sys_post
-- ----------------------------
DROP TABLE IF EXISTS `sys_post`;
CREATE TABLE `sys_post` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '岗位ID',
  `post_code` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_bin NOT NULL COMMENT '岗位编码',
  `post_name` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_bin NOT NULL COMMENT '岗位名称',
  `post_status` int NOT NULL DEFAULT '1' COMMENT '状态（1->正常 2->停用）',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `remark` varchar(500) CHARACTER SET utf8mb3 COLLATE utf8mb3_bin DEFAULT NULL COMMENT '备注',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin ROW_FORMAT=DYNAMIC COMMENT='岗位信息表';

-- ----------------------------
-- Records of sys_post
-- ----------------------------
BEGIN;
INSERT INTO `sys_post` (`id`, `post_code`, `post_name`, `post_status`, `create_time`, `remark`) VALUES (1, 'AAA', '研发总监', 1, '2023-06-14 20:08:22', '主管各个部门');
INSERT INTO `sys_post` (`id`, `post_code`, `post_name`, `post_status`, `create_time`, `remark`) VALUES (10, 'ops', '运维工程师', 1, '2025-06-28 22:46:33', '运维工程师');
INSERT INTO `sys_post` (`id`, `post_code`, `post_name`, `post_status`, `create_time`, `remark`) VALUES (11, 'dev', '研发工程师', 1, '2025-06-28 22:50:29', '研发工程师');
INSERT INTO `sys_post` (`id`, `post_code`, `post_name`, `post_status`, `create_time`, `remark`) VALUES (12, 'test', '测试工程师', 2, '2025-06-28 22:52:57', '测试工程师');
COMMIT;

-- ----------------------------
-- Table structure for sys_role
-- ----------------------------
DROP TABLE IF EXISTS `sys_role`;
CREATE TABLE `sys_role` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
  `role_name` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL COMMENT '角色名称',
  `role_key` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL COMMENT '角色权限字符串',
  `status` int NOT NULL DEFAULT '1' COMMENT '启用状态：1->启用；2->禁用',
  `description` varchar(500) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '描述',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `name` (`role_name`) USING BTREE,
  KEY `create_time` (`create_time`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb3 ROW_FORMAT=DYNAMIC COMMENT='后台角色表';

-- ----------------------------
-- Records of sys_role
-- ----------------------------
BEGIN;
INSERT INTO `sys_role` (`id`, `role_name`, `role_key`, `status`, `description`, `create_time`) VALUES (1, '超级管理员', 'admin', 1, '最大权限', '2023-06-12 20:04:53');
INSERT INTO `sys_role` (`id`, `role_name`, `role_key`, `status`, `description`, `create_time`) VALUES (13, '游客', 'test1', 1, 'test1', '2025-07-03 18:47:25');
COMMIT;

-- ----------------------------
-- Table structure for sys_role_menu
-- ----------------------------
DROP TABLE IF EXISTS `sys_role_menu`;
CREATE TABLE `sys_role_menu` (
  `role_id` int DEFAULT NULL COMMENT '角色ID',
  `menu_id` int DEFAULT NULL COMMENT '菜单ID'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 ROW_FORMAT=DYNAMIC COMMENT='角色和菜单关系表';

-- ----------------------------
-- Records of sys_role_menu
-- ----------------------------
BEGIN;
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 72);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 80);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 78);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 149);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 150);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 151);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 152);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 153);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 154);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 155);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 165);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 166);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 95);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 146);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 147);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 148);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 88);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 89);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 90);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 91);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 81);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 82);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 167);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 168);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 169);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 170);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 171);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 83);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 172);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 173);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 174);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 175);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 176);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 177);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 105);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 178);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 179);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 180);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 181);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 182);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 183);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 184);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 93);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 185);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 186);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 187);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 188);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 189);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 190);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 191);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 192);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 193);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 194);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 195);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 196);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 106);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 197);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 198);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 199);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 200);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 201);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 202);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 203);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 204);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 205);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 206);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 108);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 107);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 109);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 110);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 142);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 143);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 144);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 145);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 162);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 163);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 164);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 111);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 139);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 140);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 141);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 160);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 161);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 97);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 99);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 136);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 137);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 138);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 156);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 157);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 98);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 133);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 134);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 135);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 100);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 130);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 131);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 132);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 159);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 101);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 103);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 102);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 126);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 127);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 128);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 129);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 4);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 6);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 16);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 17);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 18);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 60);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 7);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 21);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 22);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 23);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 24);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 8);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 26);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 27);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 28);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 9);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 29);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 30);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 31);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 10);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 32);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 33);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 34);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 84);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 85);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 123);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 124);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 125);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 86);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 119);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 120);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 121);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 122);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 104);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 114);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 115);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 117);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 118);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 44);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 45);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 47);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 113);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 73);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 46);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 49);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 62);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 96);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 72);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 149);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 150);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 152);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 153);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 154);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 165);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 166);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 146);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 89);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 167);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 168);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 170);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 172);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 173);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 174);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 175);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 178);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 179);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 180);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 181);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 182);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 185);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 186);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 190);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 193);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 195);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 197);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 200);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 203);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 205);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 108);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 107);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 142);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 144);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 162);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 139);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 140);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 160);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 136);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 137);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 156);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 133);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 130);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 132);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 159);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 103);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 127);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 128);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 16);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 21);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 26);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 29);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 32);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 125);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 86);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 122);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 119);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 120);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 121);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 118);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 73);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 62);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 96);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 80);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 78);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 95);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 88);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 81);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 82);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 83);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 105);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 93);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 106);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 109);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 110);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 111);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 97);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 99);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 98);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 100);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 101);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 102);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 4);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 6);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 7);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 8);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 9);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 10);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 84);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 85);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 104);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 44);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 45);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 46);
COMMIT;

-- ----------------------------
-- Table structure for task_ansible
-- ----------------------------
DROP TABLE IF EXISTS `task_ansible`;
CREATE TABLE `task_ansible` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键ID''',
  `name` varchar(100) NOT NULL COMMENT '''任务名称''',
  `description` text COMMENT '''任务描述''',
  `type` bigint NOT NULL DEFAULT '1' COMMENT '''任务类型:1-手动,2-Git,3-K8s''',
  `git_repo` varchar(255) DEFAULT NULL COMMENT '''Git仓库地址''',
  `host_groups` text NOT NULL COMMENT '''主机分组JSON''',
  `all_host_ids` text NOT NULL COMMENT '''所有主机ID JSON数组''',
  `global_vars` text COMMENT '''全局变量JSON''',
  `status` bigint NOT NULL DEFAULT '1' COMMENT '''任务状态:1-等待中,2-运行中,3-成功,4-异常''',
  `created_at` datetime(3) NOT NULL COMMENT '''创建时间''',
  `updated_at` datetime(3) NOT NULL COMMENT '''更新时间''',
  `error_msg` text COMMENT '''错误信息''',
  `task_count` bigint NOT NULL DEFAULT '0' COMMENT '''任务数量(Type=1时为上传文件数,Type=2时为解析的playbook数,Type=3时固定为1)''',
  `total_duration` bigint NOT NULL DEFAULT '0' COMMENT '''任务执行总耗时(秒,所有子任务耗时总和)''',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_task_ansible_name` (`name`),
  KEY `idx_task_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=44 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of task_ansible
-- ----------------------------
BEGIN;
INSERT INTO `task_ansible` (`id`, `name`, `description`, `type`, `git_repo`, `host_groups`, `all_host_ids`, `global_vars`, `status`, `created_at`, `updated_at`, `error_msg`, `task_count`, `total_duration`) VALUES (21, '前端测试001', '', 1, '', '{\"api\":[445],\"web\":[444]}', '[444,445]', '', 1, '2025-08-23 20:41:19.309', '2025-08-23 20:41:19.775', '', 1, 0);
INSERT INTO `task_ansible` (`id`, `name`, `description`, `type`, `git_repo`, `host_groups`, `all_host_ids`, `global_vars`, `status`, `created_at`, `updated_at`, `error_msg`, `task_count`, `total_duration`) VALUES (22, '前端测试002', '', 1, '', '{\"api\":[445],\"web\":[444]}', '[444,445]', '', 1, '2025-08-23 20:52:59.124', '2025-08-23 20:52:59.616', '', 1, 0);
INSERT INTO `task_ansible` (`id`, `name`, `description`, `type`, `git_repo`, `host_groups`, `all_host_ids`, `global_vars`, `status`, `created_at`, `updated_at`, `error_msg`, `task_count`, `total_duration`) VALUES (23, '前端测试003', '', 1, '', '{\"web\":[444]}', '[444]', '', 1, '2025-08-23 20:57:09.757', '2025-08-23 20:57:10.240', '', 1, 0);
INSERT INTO `task_ansible` (`id`, `name`, `description`, `type`, `git_repo`, `host_groups`, `all_host_ids`, `global_vars`, `status`, `created_at`, `updated_at`, `error_msg`, `task_count`, `total_duration`) VALUES (24, '前端测试004', '', 1, '', '{\"web\":[444]}', '[444]', '', 3, '2025-08-23 21:04:05.718', '2025-08-23 22:04:47.011', '', 1, 15);
INSERT INTO `task_ansible` (`id`, `name`, `description`, `type`, `git_repo`, `host_groups`, `all_host_ids`, `global_vars`, `status`, `created_at`, `updated_at`, `error_msg`, `task_count`, `total_duration`) VALUES (25, '前端测试005', '', 2, 'git@gitee.com:zhang_fan1024/ansible-playbook.git', '{\"api\":[444],\"web\":[444]}', '[444,444]', '', 3, '2025-08-23 21:10:24.966', '2025-08-24 00:36:24.804', '', 2, 37);
INSERT INTO `task_ansible` (`id`, `name`, `description`, `type`, `git_repo`, `host_groups`, `all_host_ids`, `global_vars`, `status`, `created_at`, `updated_at`, `error_msg`, `task_count`, `total_duration`) VALUES (27, '2025-09-08测试ansible自动化任务', '', 2, 'git@gitee.com:zhang_fan1024/ansible-playbook.git', '{\"test\":[449]}', '[449]', '', 3, '2025-09-08 17:32:44.331', '2025-09-08 17:33:10.471', '', 2, 8);
INSERT INTO `task_ansible` (`id`, `name`, `description`, `type`, `git_repo`, `host_groups`, `all_host_ids`, `global_vars`, `status`, `created_at`, `updated_at`, `error_msg`, `task_count`, `total_duration`) VALUES (36, 'k8s-3-deployment', 'k8s-3', 3, 'git@gitee.com:zhang_fan1024/zf-k8s.git', '{\"etcd\":[454],\"masters\":[454],\"workers\":[456,455]}', '[454,456,455]', '{\"cluster_name\":\"k8s-3\",\"cluster_version\":\"1.30.4\",\"deployment_mode\":\"1\",\"enabled_components\":\"coredns,metrics-server,calico\"}', 4, '2025-09-09 18:36:52.609', '2025-09-10 20:11:41.430', '', 1, 417);
INSERT INTO `task_ansible` (`id`, `name`, `description`, `type`, `git_repo`, `host_groups`, `all_host_ids`, `global_vars`, `status`, `created_at`, `updated_at`, `error_msg`, `task_count`, `total_duration`) VALUES (37, 'k8s-4-deployment', '', 3, 'git@gitee.com:zhang_fan1024/zf-k8s.git', '{\"etcd\":[454],\"masters\":[454],\"workers\":[455,456]}', '[455,456,454]', '{\"cluster_name\":\"k8s-4\",\"cluster_version\":\"1.30.4\",\"deployment_mode\":\"1\",\"enabled_components\":\"coredns,calico\"}', 4, '2025-09-10 12:40:58.904', '2025-09-10 14:57:07.782', '', 1, 2);
INSERT INTO `task_ansible` (`id`, `name`, `description`, `type`, `git_repo`, `host_groups`, `all_host_ids`, `global_vars`, `status`, `created_at`, `updated_at`, `error_msg`, `task_count`, `total_duration`) VALUES (38, 'k8s-5-deployment', '', 3, 'git@gitee.com:zhang_fan1024/zf-k8s.git', '{\"etcd\":[454],\"masters\":[454],\"workers\":[455,456]}', '[454,455,456]', '{\"cluster_name\":\"k8s-5\",\"cluster_version\":\"1.30.4\",\"deployment_mode\":\"1\",\"enabled_components\":\"coredns,calico\"}', 3, '2025-09-10 14:59:48.343', '2025-09-10 15:33:51.786', '', 1, 2022);
INSERT INTO `task_ansible` (`id`, `name`, `description`, `type`, `git_repo`, `host_groups`, `all_host_ids`, `global_vars`, `status`, `created_at`, `updated_at`, `error_msg`, `task_count`, `total_duration`) VALUES (40, '虚拟机测试环境k8s集群001-deployment', '', 3, 'git@gitee.com:zhang_fan1024/zf-k8s.git', '{\"etcd\":[454],\"masters\":[454],\"workers\":[455,456]}', '[454,455,456]', '{\"cluster_name\":\"虚拟机测试环境k8s集群001\",\"cluster_version\":\"1.31.0\",\"deployment_mode\":\"1\",\"enabled_components\":\"coredns,metrics-server,ingress-nginx,calico\"}', 3, '2025-09-19 23:32:43.465', '2025-09-20 00:01:27.442', '', 1, 1590);
INSERT INTO `task_ansible` (`id`, `name`, `description`, `type`, `git_repo`, `host_groups`, `all_host_ids`, `global_vars`, `status`, `created_at`, `updated_at`, `error_msg`, `task_count`, `total_duration`) VALUES (42, '腾讯云k8s001-deployment', '腾讯云k8s001', 3, 'git@gitee.com:zhang_fan1024/zf-k8s.git', '{\"etcd\":[457],\"masters\":[457],\"workers\":[458,459]}', '[457,458,459]', '{\"cluster_name\":\"腾讯云k8s001\",\"cluster_version\":\"1.31.0\",\"deployment_mode\":\"1\",\"enabled_components\":\"coredns,metrics-server,ingress-nginx,calico\"}', 4, '2025-09-20 13:25:54.975', '2025-09-20 13:38:38.337', '', 1, 614);
INSERT INTO `task_ansible` (`id`, `name`, `description`, `type`, `git_repo`, `host_groups`, `all_host_ids`, `global_vars`, `status`, `created_at`, `updated_at`, `error_msg`, `task_count`, `total_duration`) VALUES (43, '测试ansible任务', '', 2, 'git@gitee.com:zhang_fan1024/ansible-playbook.git', '{\"api\":[469]}', '[469]', '', 4, '2025-09-28 20:47:24.933', '2025-09-28 20:48:18.667', '', 2, 21);
COMMIT;

-- ----------------------------
-- Table structure for task_ansiblework
-- ----------------------------
DROP TABLE IF EXISTS `task_ansiblework`;
CREATE TABLE `task_ansiblework` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键ID''',
  `task_id` bigint unsigned NOT NULL COMMENT '''父任务ID''',
  `entry_file_name` varchar(255) NOT NULL COMMENT '''入口文件名''',
  `entry_file_path` varchar(255) NOT NULL COMMENT '''入口文件路径''',
  `log_path` varchar(255) DEFAULT NULL COMMENT '''日志路径''',
  `status` bigint NOT NULL DEFAULT '1' COMMENT '''子任务状态:1-等待中,2-运行中,3-成功,4-异常''',
  `start_time` datetime(3) DEFAULT NULL COMMENT '''开始时间''',
  `end_time` datetime(3) DEFAULT NULL COMMENT '''结束时间''',
  `duration` bigint DEFAULT NULL COMMENT '''执行耗时(秒)''',
  `exit_code` bigint DEFAULT NULL COMMENT '''退出代码''',
  `error_msg` text COMMENT '''错误信息''',
  `log` text COMMENT '''日志内容''',
  PRIMARY KEY (`id`),
  KEY `idx_task_ansiblework_task_id` (`task_id`),
  KEY `idx_task_id` (`task_id`),
  KEY `idx_task_work_composite` (`task_id`,`status`),
  CONSTRAINT `fk_task_ansible_works` FOREIGN KEY (`task_id`) REFERENCES `task_ansible` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=39 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of task_ansiblework
-- ----------------------------
BEGIN;
INSERT INTO `task_ansiblework` (`id`, `task_id`, `entry_file_name`, `entry_file_path`, `log_path`, `status`, `start_time`, `end_time`, `duration`, `exit_code`, `error_msg`, `log`) VALUES (16, 21, '01-linux-os.yaml.yml', 'task/21/前端测试001/01-linux-os.yaml.yml', '', 1, NULL, NULL, 0, 0, '', NULL);
INSERT INTO `task_ansiblework` (`id`, `task_id`, `entry_file_name`, `entry_file_path`, `log_path`, `status`, `start_time`, `end_time`, `duration`, `exit_code`, `error_msg`, `log`) VALUES (17, 22, '01-linux-os.yaml.yml', 'task/22/前端测试002/01-linux-os.yaml.yml', '', 1, NULL, NULL, 0, 0, '', NULL);
INSERT INTO `task_ansiblework` (`id`, `task_id`, `entry_file_name`, `entry_file_path`, `log_path`, `status`, `start_time`, `end_time`, `duration`, `exit_code`, `error_msg`, `log`) VALUES (18, 23, '01-linux-os.yaml.yml', 'task/23/前端测试003/01-linux-os.yaml.yml', '', 1, NULL, NULL, 0, 0, '', NULL);
INSERT INTO `task_ansiblework` (`id`, `task_id`, `entry_file_name`, `entry_file_path`, `log_path`, `status`, `start_time`, `end_time`, `duration`, `exit_code`, `error_msg`, `log`) VALUES (19, 24, '01-linux-os.yaml.yml', 'task/24/前端测试004/01-linux-os.yaml.yml', 'logs/ansible/24/19/01-linux-os.yaml.yml.log', 3, '2025-08-23 22:04:31.142', '2025-08-23 22:04:46.673', 15, 0, '', NULL);
INSERT INTO `task_ansiblework` (`id`, `task_id`, `entry_file_name`, `entry_file_path`, `log_path`, `status`, `start_time`, `end_time`, `duration`, `exit_code`, `error_msg`, `log`) VALUES (20, 25, '01-linux-os.yaml', 'task/25/前端测试005/01-linux-os.yaml', 'logs/ansible/25/20/01-linux-os.yaml.log', 3, '2025-08-24 00:35:46.757', '2025-08-24 00:35:59.775', 13, 0, '', NULL);
INSERT INTO `task_ansiblework` (`id`, `task_id`, `entry_file_name`, `entry_file_path`, `log_path`, `status`, `start_time`, `end_time`, `duration`, `exit_code`, `error_msg`, `log`) VALUES (21, 25, '02-os.yaml', 'task/25/前端测试005/02-os.yaml', 'logs/ansible/25/21/02-os.yaml.log', 3, '2025-08-24 00:35:59.985', '2025-08-24 00:36:24.468', 24, 0, '', NULL);
INSERT INTO `task_ansiblework` (`id`, `task_id`, `entry_file_name`, `entry_file_path`, `log_path`, `status`, `start_time`, `end_time`, `duration`, `exit_code`, `error_msg`, `log`) VALUES (22, 27, '01-linux-os.yaml', 'task/27/2025-09-08测试ansible自动化任务/01-linux-os.yaml', 'logs/ansible/27/22/01-linux-os.yaml.log', 3, '2025-09-08 17:33:00.979', '2025-09-08 17:33:07.235', 6, 0, '', NULL);
INSERT INTO `task_ansiblework` (`id`, `task_id`, `entry_file_name`, `entry_file_path`, `log_path`, `status`, `start_time`, `end_time`, `duration`, `exit_code`, `error_msg`, `log`) VALUES (23, 27, '02-os.yaml', 'task/27/2025-09-08测试ansible自动化任务/02-os.yaml', 'logs/ansible/27/23/02-os.yaml.log', 3, '2025-09-08 17:33:07.467', '2025-09-08 17:33:10.126', 2, 0, '', NULL);
INSERT INTO `task_ansiblework` (`id`, `task_id`, `entry_file_name`, `entry_file_path`, `log_path`, `status`, `start_time`, `end_time`, `duration`, `exit_code`, `error_msg`, `log`) VALUES (30, 36, 'deploy-simple.sh', './scripts/deploy-simple.sh', 'logs/ansible/36/30/deploy-simple.sh.log', 4, '2025-09-10 20:04:43.504', '2025-09-10 20:11:41.111', 417, 99, 'exit status 99', NULL);
INSERT INTO `task_ansiblework` (`id`, `task_id`, `entry_file_name`, `entry_file_path`, `log_path`, `status`, `start_time`, `end_time`, `duration`, `exit_code`, `error_msg`, `log`) VALUES (31, 37, 'deploy-simple.sh', './scripts/deploy-simple.sh', 'logs/ansible/37/31/deploy-simple.sh.log', 4, '2025-09-10 14:57:04.562', '2025-09-10 14:57:07.467', 2, 4, 'exit status 4', NULL);
INSERT INTO `task_ansiblework` (`id`, `task_id`, `entry_file_name`, `entry_file_path`, `log_path`, `status`, `start_time`, `end_time`, `duration`, `exit_code`, `error_msg`, `log`) VALUES (32, 38, 'deploy-simple.sh', './scripts/deploy-simple.sh', 'logs/ansible/38/32/deploy-simple.sh.log', 3, '2025-09-10 15:00:09.492', '2025-09-10 15:33:51.446', 2022, 0, '', NULL);
INSERT INTO `task_ansiblework` (`id`, `task_id`, `entry_file_name`, `entry_file_path`, `log_path`, `status`, `start_time`, `end_time`, `duration`, `exit_code`, `error_msg`, `log`) VALUES (34, 40, 'deploy-simple.sh', './scripts/deploy-simple.sh', 'logs/ansible/40/34/deploy-simple.sh.log', 3, '2025-09-19 23:34:56.212', '2025-09-20 00:01:27.065', 1590, 0, '', NULL);
INSERT INTO `task_ansiblework` (`id`, `task_id`, `entry_file_name`, `entry_file_path`, `log_path`, `status`, `start_time`, `end_time`, `duration`, `exit_code`, `error_msg`, `log`) VALUES (36, 42, 'deploy-simple.sh', './scripts/deploy-simple.sh', 'logs/ansible/42/36/deploy-simple.sh.log', 4, '2025-09-20 13:28:23.527', '2025-09-20 13:38:37.994', 614, 2, 'exit status 2', NULL);
INSERT INTO `task_ansiblework` (`id`, `task_id`, `entry_file_name`, `entry_file_path`, `log_path`, `status`, `start_time`, `end_time`, `duration`, `exit_code`, `error_msg`, `log`) VALUES (37, 43, '01-linux-os.yaml', 'task/43/测试ansible任务/01-linux-os.yaml', 'logs/ansible/43/37/01-linux-os.yaml.log', 3, '2025-09-28 20:47:56.817', '2025-09-28 20:48:02.051', 5, 0, '', NULL);
INSERT INTO `task_ansiblework` (`id`, `task_id`, `entry_file_name`, `entry_file_path`, `log_path`, `status`, `start_time`, `end_time`, `duration`, `exit_code`, `error_msg`, `log`) VALUES (38, 43, '02-os.yaml', 'task/43/测试ansible任务/02-os.yaml', 'logs/ansible/43/38/02-os.yaml.log', 4, '2025-09-28 20:48:02.272', '2025-09-28 20:48:18.337', 16, 2, 'exit status 2', NULL);
COMMIT;

-- ----------------------------
-- Table structure for task_job
-- ----------------------------
DROP TABLE IF EXISTS `task_job`;
CREATE TABLE `task_job` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL COMMENT '任务标题',
  `type` bigint DEFAULT NULL COMMENT '任务类型 1=普通任务,2=定时任务,3=ansible任务',
  `shell` text COMMENT '任务内容(任务模板ID,多个用逗号分隔)',
  `host_ids` text COMMENT '主机ID(多个用逗号分隔)',
  `cron_expr` varchar(255) DEFAULT NULL COMMENT '定时表达式(* * * * *)',
  `tasklog` text COMMENT '任务执行日志',
  `status` bigint DEFAULT NULL COMMENT '任务状态 1=等待中,2=运行中,3=成功,4=异常,5=已暂停',
  `duration` bigint DEFAULT NULL COMMENT '执行耗时(秒)',
  `remark` text COMMENT '任务备注',
  `start_time` datetime(3) DEFAULT NULL COMMENT '任务开始时间',
  `end_time` datetime(3) DEFAULT NULL COMMENT '任务结束时间',
  `created_at` datetime(3) DEFAULT NULL COMMENT '任务创建时间',
  `task_count` bigint DEFAULT NULL COMMENT '任务数量',
  `is_recurring` tinyint(1) DEFAULT NULL COMMENT '是否周期性任务',
  `scheduled_time` datetime(3) DEFAULT NULL COMMENT '计划执行时间',
  `log_path` varchar(500) DEFAULT NULL COMMENT '日志文件路径',
  `execute_count` bigint DEFAULT '0' COMMENT '执行次数',
  `next_run_time` datetime(3) DEFAULT NULL COMMENT '下次执行时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=59 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of task_job
-- ----------------------------
BEGIN;
INSERT INTO `task_job` (`id`, `name`, `type`, `shell`, `host_ids`, `cron_expr`, `tasklog`, `status`, `duration`, `remark`, `start_time`, `end_time`, `created_at`, `task_count`, `is_recurring`, `scheduled_time`, `log_path`, `execute_count`, `next_run_time`) VALUES (49, 'test4', 1, '12', '1', '', '', 3, 5, '', NULL, '2025-09-29 21:15:32.496', '2025-09-29 21:01:51.867', 1, NULL, NULL, NULL, 1, NULL);
INSERT INTO `task_job` (`id`, `name`, `type`, `shell`, `host_ids`, `cron_expr`, `tasklog`, `status`, `duration`, `remark`, `start_time`, `end_time`, `created_at`, `task_count`, `is_recurring`, `scheduled_time`, `log_path`, `execute_count`, `next_run_time`) VALUES (56, '定时任务01', 2, '11', '1', '*/3 * * * *', '', 5, 5, '', NULL, '2025-10-01 14:03:07.896', '2025-09-30 17:15:17.495', 1, NULL, NULL, NULL, 26, '2025-10-01 14:06:00.000');
INSERT INTO `task_job` (`id`, `name`, `type`, `shell`, `host_ids`, `cron_expr`, `tasklog`, `status`, `duration`, `remark`, `start_time`, `end_time`, `created_at`, `task_count`, `is_recurring`, `scheduled_time`, `log_path`, `execute_count`, `next_run_time`) VALUES (57, '定时任务02', 2, '12', '1', '*/5 * * * *', '', 5, 6, '', NULL, '2025-10-01 14:00:08.425', '2025-09-30 17:17:02.686', 1, NULL, NULL, NULL, 18, '2025-10-01 14:05:00.000');
INSERT INTO `task_job` (`id`, `name`, `type`, `shell`, `host_ids`, `cron_expr`, `tasklog`, `status`, `duration`, `remark`, `start_time`, `end_time`, `created_at`, `task_count`, `is_recurring`, `scheduled_time`, `log_path`, `execute_count`, `next_run_time`) VALUES (58, 'test5', 1, '10', '1', '', '', 3, 6, '123', NULL, '2025-10-01 14:53:47.061', '2025-10-01 14:53:31.593', 1, NULL, NULL, NULL, 1, NULL);
COMMIT;

-- ----------------------------
-- Table structure for task_template
-- ----------------------------
DROP TABLE IF EXISTS `task_template`;
CREATE TABLE `task_template` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `type` bigint NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `created_by` varchar(50) DEFAULT NULL,
  `updated_by` varchar(50) DEFAULT NULL,
  `remark` varchar(500) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of task_template
-- ----------------------------
BEGIN;
INSERT INTO `task_template` (`id`, `name`, `type`, `content`, `created_at`, `updated_at`, `created_by`, `updated_by`, `remark`) VALUES (2, '1-数字 1-100', 1, '#!/bin/bash\n\n# 从 1 到 100，每秒打印一个数字\nfor ((i = 1; i <= 100; i++)); do\n    echo \"[$(date +%H:%M:%S)] $i\"\n    sleep 1\ndone\n\necho \"完成：所有数字 1-100 已打印完毕。\"\n', '2025-08-06 12:47:57.073', '2025-08-12 16:14:49.394', 'admin', 'admin', '测试脚本');
INSERT INTO `task_template` (`id`, `name`, `type`, `content`, `created_at`, `updated_at`, `created_by`, `updated_by`, `remark`) VALUES (3, '2-系统状态检查脚本', 1, '#!/bin/bash\n\n#==================================================\n#       简化版 Linux 系统状态检查脚本\n#       输出直接打印到终端，不保存日志\n#==================================================\n\n# 颜色定义（可选，让输出更清晰）\nGREEN=\'\\033[0;32m\'\nYELLOW=\'\\033[1;33m\'\nRED=\'\\033[0;31m\'\nNC=\'\\033[0m\'\n\nlog() { echo -e \"${GREEN}[INFO] $1${NC}\"; }\nwarn() { echo -e \"${YELLOW}[WARN] $1${NC}\"; }\nerror() { echo -e \"${RED}[ERROR] $1${NC}\"; }\n\nseparator() {\n    echo \"----------------------------------------\"\n}\n\n# 1. 系统信息\nseparator\nlog \"系统信息\"\nseparator\necho \"主机名: $(hostname)\"\necho \"操作系统: $(grep \'^PRETTY_NAME=\' /etc/os-release | cut -d\'\"\' -f2 2>/dev/null || echo \'未知\')\"\necho \"内核: $(uname -r)\"\n\n# 2. 系统负载\nseparator\nlog \"系统负载 (1/5/15分钟)\"\nseparator\nuptime | awk -F\'load average:\' \'{print \"负载: \" $2}\'\n\n# 3. CPU 使用率\nseparator\nlog \"CPU 使用率\"\nseparator\ntop -bn1 | grep \"Cpu(s)\" | awk \'{print \"CPU 使用: \" $2 \"%\"}\'\n\n# 4. 内存使用\nseparator\nlog \"内存使用\"\nseparator\nfree -h | awk \'/^Mem:/ {printf \"总内存: %s\\t已用: %s\\t空闲: %s\\n\", $2, $3, $4}\'\n\n# 5. 磁盘使用（仅根分区）\nseparator\nlog \"磁盘使用 (/)\"\nseparator\ndf -h / | awk \'NR==2 {printf \"设备: %s\\t总大小: %s\\t已用: %s (%s)\\n\", $1, $2, $3, $5}\'\n\n# 6. 活跃进程数\nseparator\nlog \"进程统计\"\nseparator\nPROC=$(ps aux | wc -l)\necho \"当前运行进程数: $((PROC - 1))\"\n\n# 7. 网络监听端口（常用服务）\nseparator\nlog \"监听端口 (常见服务)\"\nseparator\nss -tuln | grep \':80\\|:443\\|:22\\|:3306\\|:5432\' || echo \"无常见服务监听\"\n\n# 8. 最近登录与 SSH 尝试\nseparator\nlog \"最近登录与安全\"\nseparator\necho \"最近3次登录:\"\nlast | head -3 | awk \'{print $1, $3, $4, $5}\'\n\nFAILED=$(grep \"Failed password\" /var/log/auth.log 2>/dev/null | wc -l)\nif [ $FAILED -gt 0 ]; then\n    warn \"发现 $FAILED 次 SSH 登录失败，请检查安全！\"\nelse\n    log \"无 SSH 登录失败记录\"\nfi\n\nseparator\nlog \"系统状态检查完成。\"', '2025-08-06 12:48:34.435', '2025-08-12 16:14:21.982', 'root', 'admin', 'test');
INSERT INTO `task_template` (`id`, `name`, `type`, `content`, `created_at`, `updated_at`, `created_by`, `updated_by`, `remark`) VALUES (4, '3-自动优化 Linux 内核参数', 2, '#!/bin/bash\n\n# 自动优化 Linux 内核参数\n# 脚本针对 RHEL7\n\ncat >> /usr/lib/sysctl.d/00-system.conf << \'EOF\'\nfs.file-max=65535\nnet.ipv4.tcp_timestamps = 0\nnet.ipv4.tcp_synack_retries = 5\nnet.ipv4.tcp_syn_retries = 5\nnet.ipv4.tcp_tw_recycle = 1\nnet.ipv4.tcp_tw_reuse = 1\nnet.ipv4.tcp_fin_timeout = 30\n#net.ipv4.tcp_keepalive_time = 120\nnet.ipv4.ip_local_port_range = 1024 65535\nkernel.shmall = 2097152\nkernel.shmmax = 2147483648\nkernel.shmmni = 4096\nkernel.sem = 5010 641280 5010 128\nnet.core.wmem_default=262144\nnet.core.wmem_max=262144\nnet.core.rmem_default=4194304\nnet.core.rmem_max=4194304\nnet.ipv4.tcp_fin_timeout = 10\nnet.ipv4.tcp_keepalive_time = 30\nnet.ipv4.tcp_window_scaling = 0\nnet.ipv4.tcp_sack = 0\nEOF\n\nsysctl -p\nsleep 3\necho  \"已经完成优化\"\n', '2025-08-06 12:49:39.227', '2025-09-06 16:27:02.319', 'lisi', 'admin', 'test');
INSERT INTO `task_template` (`id`, `name`, `type`, `content`, `created_at`, `updated_at`, `created_by`, `updated_by`, `remark`) VALUES (6, '切割 Nginx 日志文件', 2, '#mkdir  /data/scripts  \n#vim   /data/scripts/nginx_log.sh    \n#!/bin/bash  \n  \n# 切割 Nginx 日志文件(防止单个文件过大,后期处理很困难)   \nlogs_path=\"/usr/local/nginx/logs/\"  \nmv ${logs_path}access.log ${logs_path}access_$(date -d \"yesterday\" +\"%Y%m%d\").log  \nkill -USR1  `cat /usr/local/nginx/logs/nginx.pid`  \n  \n# chmod +x  /data/scripts/nginx_log.sh  \n# crontab  ‐e                    #脚本写完后,将脚本放入计划任务每天执行一次脚本  \n0  1  *  *   *   /data/scripts/nginx_log.sh  ', '2025-08-06 12:53:39.020', '2025-09-06 16:27:07.841', 'zhangsan', 'admin', 'test');
INSERT INTO `task_template` (`id`, `name`, `type`, `content`, `created_at`, `updated_at`, `created_by`, `updated_by`, `remark`) VALUES (9, '查找 Linux 系统中的僵尸进程  ', 2, '#!/bin/bash  \n  \n# 查找 Linux 系统中的僵尸进程  \n# awk 判断 ps 命令输出的第 8 列为 Z 是,显示该进程的 PID 和进程命令  \nps aux | awk \'{if($8 == \"Z\"){print $2,$11}}\'  ', '2025-08-06 14:45:58.982', '2025-09-06 16:27:18.196', 'admin', 'admin', '查找 Linux 系统中的僵尸进程  ');
INSERT INTO `task_template` (`id`, `name`, `type`, `content`, `created_at`, `updated_at`, `created_by`, `updated_by`, `remark`) VALUES (10, '获取本机 MAC 地址  ', 1, '#!/bin/bash  \n  \n# 获取本机 MAC 地址  \nip a s | awk \'BEGIN{print  \" 本 机 MAC 地 址 信 息 如 下 :\"}/^[0‐9]/{print $2;getline;if($0~/link\\/ether/){print $2}}\' | grep -v lo:  \n# awk 读取 ip 命令的输出,输出结果中如果有以数字开始的行,先显示该行的地 2 列(网卡名称),  \n# 接着使用 getline 再读取它的下一行数据,判断是否包含 link/ether  \n# 如果保护该关键词,就显示该行的第 2 列(MAC 地址)  \n# lo 回环设备没有 MAC,因此将其屏蔽,不显示  ', '2025-08-06 16:45:07.990', '2025-08-12 16:22:43.774', 'admin', 'admin', 'test');
INSERT INTO `task_template` (`id`, `name`, `type`, `content`, `created_at`, `updated_at`, `created_by`, `updated_by`, `remark`) VALUES (11, '统计 Linux 进程相关数量信息   ', 1, '#!/bin/bash  \n  \n# 统计 Linux 进程相关数量信息   \nrunning=0  \nsleeping=0  \nstoped=0  \nzombie=0  \n# 在 proc 目录下所有以数字开始的都是当前计算机正在运行的进程的进程 PID  \n# 每个 PID 编号的目录下记录有该进程相关的信息  \nfor pid in /proc/[1‐9]*  \ndo  \n  procs=$[procs+1]  \n  stat=$(awk \'{print $3}\' $pid/stat)  \n# 每个 pid 目录下都有一个 stat 文件,该文件的第 3 列是该进程的状态信息  \n    case $stat in  \n    R)  \n    running=$[running+1]  \n    ;;  \n    T)  \n    stoped=$[stoped+1]  \n    ;;  \n    S)  \n    sleeping=$[sleeping+1]  \n    ;;  \n    Z)  \n     zombie=$[zombie+1]  \n     ;;  \n    esac  \ndone  \necho \"进程统计信息如下\"  \necho \"总进程数量为:$procs\"  \necho \"Running 进程数为:$running\"  \necho \"Stoped 进程数为:$stoped\"  \necho \"Sleeping 进程数为:$sleeping\"  \necho \"Zombie 进程数为:$zombie\"  ', '2025-08-06 17:19:10.669', '2025-08-12 16:23:10.614', 'admin', 'admin', 'ooo');
INSERT INTO `task_template` (`id`, `name`, `type`, `content`, `created_at`, `updated_at`, `created_by`, `updated_by`, `remark`) VALUES (12, ' 显示本机 Linux 系统上所有开放的端口列表   ', 1, '#!/bin/bash  \n  \n# 显示本机 Linux 系统上所有开放的端口列表   \n  \n# 从端口列表中观测有没有没用的端口,有的话可以将该端口对应的服务关闭,防止意外的攻击可能性  \nss -nutlp | awk \'{print $1,$5}\' | awk -F\"[: ]\" \'{print \"协议:\"$1,\"端口号:\"$NF}\' | grep \"[0‐9]\" | uniq  ', '2025-08-08 10:03:54.197', '2025-08-12 16:24:35.774', 'admin', 'admin', '新机器初始化');
INSERT INTO `task_template` (`id`, `name`, `type`, `content`, `created_at`, `updated_at`, `created_by`, `updated_by`, `remark`) VALUES (13, 'nginx安装配置', 3, '---\n- hosts: web   #指定标签组\n  vars:\n    hello: Ansible  #定义变量\n\n  tasks:\n  - name: \"配置nginx源\"\n    yum_repository:\n      name: nginx\n      description: nginx repo\n      baseurl: http://nginx.org/packages/centos/7/$basearch/\n      gpgcheck: no\n      enabled: 1\n  - name: \"安装nginx\"\n    yum:\n      name: nginx\n      state: latest\n  - name: \"编写虚拟主机的配置文件\"\n    copy:\n      src: ./site.conf\n      dest: /etc/nginx/conf.d/site.conf\n  - name: \"创建虚拟主机目录\"\n    file:\n      dest: /var/www/html\n      state: directory\n  - name: \"启动nginx\"\n    service:\n      name: nginx\n      state: started\n  - name:  \"测试nginx页面\"\n    shell: echo \"hello  {{hello}}\" > /var/www/html/index.html\n    #引用变量{{}}表示\n', '2025-08-15 10:58:03.845', '2025-08-15 10:58:03.845', 'admin', '', 'nginx安装配置');
COMMIT;

-- ----------------------------
-- Table structure for task_work
-- ----------------------------
DROP TABLE IF EXISTS `task_work`;
CREATE TABLE `task_work` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `task_id` bigint unsigned DEFAULT NULL COMMENT '关联的任务ID',
  `template_id` bigint unsigned DEFAULT NULL COMMENT '任务模板ID',
  `host_id` bigint unsigned DEFAULT NULL COMMENT '执行主机ID',
  `status` bigint DEFAULT NULL COMMENT '任务状态 1=等待中,2=运行中,3=成功,4=异常',
  `log` text COMMENT '任务日志',
  `log_path` text COMMENT '日志文件路径',
  `start_time` datetime(3) DEFAULT NULL COMMENT '任务开始时间',
  `end_time` datetime(3) DEFAULT NULL COMMENT '任务结束时间',
  `duration` bigint DEFAULT NULL COMMENT '执行耗时(秒)',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `type` bigint DEFAULT NULL COMMENT '任务类型 1=普通任务,2=定时任务',
  `scheduled_time` datetime(3) DEFAULT NULL COMMENT '定时任务执行时间',
  `cron_expr` longtext COMMENT 'cron表达式',
  `is_recurring` tinyint(1) DEFAULT NULL COMMENT '是否周期性任务',
  PRIMARY KEY (`id`),
  KEY `idx_task_work_task_id` (`task_id`),
  KEY `idx_task_work_template_id` (`template_id`)
) ENGINE=InnoDB AUTO_INCREMENT=107 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of task_work
-- ----------------------------
BEGIN;
INSERT INTO `task_work` (`id`, `task_id`, `template_id`, `host_id`, `status`, `log`, `log_path`, `start_time`, `end_time`, `duration`, `created_at`, `type`, `scheduled_time`, `cron_expr`, `is_recurring`) VALUES (90, 48, 11, 1, 2, '[2025-09-29-21:42:02] 任务开始\n进程统计信息如下\n总进程数量为:79\nRunning 进程数为:1\nStoped 进程数为:0\nSleeping 进程数为:55\nZombie 进程数为:0\n[2025-09-29-21:42:02] 任务完成\n', 'logs/task_48/task_48_template_11.log', '2025-09-29 21:42:00.299', '2025-09-29 21:42:06.276', 5, '2025-09-29 20:58:17.259', 0, '2025-09-29 21:45:00.000', NULL, NULL);
INSERT INTO `task_work` (`id`, `task_id`, `template_id`, `host_id`, `status`, `log`, `log_path`, `start_time`, `end_time`, `duration`, `created_at`, `type`, `scheduled_time`, `cron_expr`, `is_recurring`) VALUES (91, 48, 0, 0, 4, '', '', NULL, NULL, 0, '2025-09-29 20:58:17.515', 2, '2025-09-29 21:00:00.000', NULL, NULL);
INSERT INTO `task_work` (`id`, `task_id`, `template_id`, `host_id`, `status`, `log`, `log_path`, `start_time`, `end_time`, `duration`, `created_at`, `type`, `scheduled_time`, `cron_expr`, `is_recurring`) VALUES (92, 49, 12, 1, 3, '[2025-09-29-21:15:27] 任务开始\n协议:tcp 端口号:6379\n协议:tcp 端口号:4330\n协议:tcp 端口号:3306\n协议:tcp 端口号:80\n协议:tcp 端口号:33060\n协议:tcp 端口号:8088\n协议:tcp 端口号:8080\n协议:tcp 端口号:4330\n协议:tcp 端口号:80\n协议:tcp 端口号:8080\n协议:tcp 端口号:8086\n协议:tcp 端口号:9091\n协议:tcp 端口号:9090\n[2025-09-29-21:15:27] 任务完成\n', 'logs/task_49/task_49_template_12.log', '2025-09-29 21:15:25.761', '2025-09-29 21:15:31.262', 5, '2025-09-29 21:01:51.970', 0, NULL, NULL, NULL);
INSERT INTO `task_work` (`id`, `task_id`, `template_id`, `host_id`, `status`, `log`, `log_path`, `start_time`, `end_time`, `duration`, `created_at`, `type`, `scheduled_time`, `cron_expr`, `is_recurring`) VALUES (104, 56, 11, 1, 1, '[2025-10-01-14:03:02] 任务开始\n进程统计信息如下\n总进程数量为:80\nRunning 进程数为:0\nStoped 进程数为:0\nSleeping 进程数为:56\nZombie 进程数为:0\n[2025-10-01-14:03:03] 任务完成\n', 'logs/task_56/task_56_template_11.log', '2025-10-01 14:03:00.890', '2025-10-01 14:03:06.717', 5, '2025-09-30 17:15:17.596', 0, '2025-09-30 17:18:00.000', NULL, NULL);
INSERT INTO `task_work` (`id`, `task_id`, `template_id`, `host_id`, `status`, `log`, `log_path`, `start_time`, `end_time`, `duration`, `created_at`, `type`, `scheduled_time`, `cron_expr`, `is_recurring`) VALUES (105, 57, 12, 1, 1, '[2025-10-01-14:00:03] 任务开始\n协议:tcp 端口号:6379\n协议:tcp 端口号:4330\n协议:tcp 端口号:3306\n协议:tcp 端口号:80\n协议:tcp 端口号:33060\n协议:tcp 端口号:8088\n协议:tcp 端口号:8080\n协议:tcp 端口号:4330\n协议:tcp 端口号:80\n协议:tcp 端口号:8080\n协议:tcp 端口号:8086\n协议:tcp 端口号:9091\n协议:tcp 端口号:9090\n[2025-10-01-14:00:03] 任务完成\n', 'logs/task_57/task_57_template_12.log', '2025-10-01 14:00:00.906', '2025-10-01 14:00:07.272', 6, '2025-09-30 17:17:02.785', 0, '2025-09-30 17:20:00.000', NULL, NULL);
INSERT INTO `task_work` (`id`, `task_id`, `template_id`, `host_id`, `status`, `log`, `log_path`, `start_time`, `end_time`, `duration`, `created_at`, `type`, `scheduled_time`, `cron_expr`, `is_recurring`) VALUES (106, 58, 10, 1, 3, '[2025-10-01-14:53:41] 任务开始\n 本 机 MAC 地 址 信 息 如 下 :\nveth47f9cf4@if2:\n86:22:92:32:6b:13\n[2025-10-01-14:53:41] 任务完成\n', 'logs/task_58/task_58_template_10.log', '2025-10-01 14:53:39.692', '2025-10-01 14:53:45.836', 6, '2025-10-01 14:53:31.697', 0, NULL, NULL, NULL);
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
