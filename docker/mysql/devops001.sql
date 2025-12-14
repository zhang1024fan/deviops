/*
 Navicat Premium Dump SQL

 Source Server         : 8.130.14.34
 Source Server Type    : MySQL
 Source Server Version : 80044 (8.0.44-0ubuntu0.24.04.1)
 Source Host           : 8.130.14.34:3306
 Source Schema         : github

 Target Server Type    : MySQL
 Target Server Version : 80044 (8.0.44-0ubuntu0.24.04.1)
 File Encoding         : 65001

 Date: 14/12/2025 14:58:23
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
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of app_application
-- ----------------------------
BEGIN;
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
) ENGINE=InnoDB AUTO_INCREMENT=54 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of app_jenkins_env
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for app_service_release
-- ----------------------------
DROP TABLE IF EXISTS `app_service_release`;
CREATE TABLE `app_service_release` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL COMMENT '发布标题',
  `business_group_id` bigint unsigned NOT NULL COMMENT '业务组ID',
  `impact_feature` text COMMENT '影响功能描述',
  `applicant_id` bigint unsigned NOT NULL COMMENT '申请人ID',
  `applicant_name` varchar(100) NOT NULL COMMENT '申请人姓名',
  `owner_approver_id` bigint unsigned DEFAULT NULL COMMENT '负责人审批人ID',
  `owner_approver_name` varchar(100) DEFAULT NULL COMMENT '负责人审批人姓名',
  `security_approver_id` bigint unsigned DEFAULT NULL COMMENT '安全审批人ID',
  `security_approver_name` varchar(100) DEFAULT NULL COMMENT '安全审批人姓名',
  `test_approver_id` bigint unsigned DEFAULT NULL COMMENT '测试审批人ID',
  `test_approver_name` varchar(100) DEFAULT NULL COMMENT '测试审批人姓名',
  `owner_approval_status` bigint DEFAULT '1' COMMENT '负责人审批状态',
  `security_approval_status` bigint DEFAULT '1' COMMENT '安全审批状态',
  `test_approval_status` bigint DEFAULT '1' COMMENT '测试审批状态',
  `owner_approval_time` datetime(3) DEFAULT NULL COMMENT '负责人审批时间',
  `security_approval_time` datetime(3) DEFAULT NULL COMMENT '安全审批时间',
  `test_approval_time` datetime(3) DEFAULT NULL COMMENT '测试审批时间',
  `owner_approval_remark` text COMMENT '负责人审批意见',
  `security_approval_remark` text COMMENT '安全审批意见',
  `test_approval_remark` text COMMENT '测试审批意见',
  `deploy_status` bigint DEFAULT '1' COMMENT '运维发布状态',
  `regression_test_status` bigint DEFAULT '1' COMMENT '回归测试状态',
  `status` bigint DEFAULT '1' COMMENT '流程状态',
  `start_time` datetime(3) DEFAULT NULL COMMENT '发布开始时间',
  `end_time` datetime(3) DEFAULT NULL COMMENT '发布结束时间',
  `duration` bigint DEFAULT '0' COMMENT '发布耗时(秒)',
  `service_count` bigint DEFAULT '0' COMMENT '关联服务数量',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_app_service_release_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of app_service_release
-- ----------------------------
BEGIN;
INSERT INTO `app_service_release` (`id`, `title`, `business_group_id`, `impact_feature`, `applicant_id`, `applicant_name`, `owner_approver_id`, `owner_approver_name`, `security_approver_id`, `security_approver_name`, `test_approver_id`, `test_approver_name`, `owner_approval_status`, `security_approval_status`, `test_approval_status`, `owner_approval_time`, `security_approval_time`, `test_approval_time`, `owner_approval_remark`, `security_approval_remark`, `test_approval_remark`, `deploy_status`, `regression_test_status`, `status`, `start_time`, `end_time`, `duration`, `service_count`, `created_at`, `updated_at`, `deleted_at`) VALUES (2, '测试---发布工单上线', 12, '影响用户登录', 1, '管理员', 102, '用户_102', 98, '用户_98', 89, '用户_89', 2, 2, 2, '2025-11-22 18:16:13.748', '2025-11-22 18:16:27.156', '2025-11-22 18:16:36.689', '同意', '同意', '同意', 3, 3, 2, '2025-11-22 18:18:03.250', '2025-11-22 18:21:45.450', 220, 1, '2025-11-22 18:13:37.037', '2025-11-22 23:01:55.144', NULL);
INSERT INTO `app_service_release` (`id`, `title`, `business_group_id`, `impact_feature`, `applicant_id`, `applicant_name`, `owner_approver_id`, `owner_approver_name`, `security_approver_id`, `security_approver_name`, `test_approver_id`, `test_approver_name`, `owner_approval_status`, `security_approval_status`, `test_approval_status`, `owner_approval_time`, `security_approval_time`, `test_approval_time`, `owner_approval_remark`, `security_approval_remark`, `test_approval_remark`, `deploy_status`, `regression_test_status`, `status`, `start_time`, `end_time`, `duration`, `service_count`, `created_at`, `updated_at`, `deleted_at`) VALUES (3, '测试多任务工单发布', 12, '测试多任务工单发布', 1, '管理员', 103, '用户_103', 98, '用户_98', 104, '用户_104', 2, 2, 2, '2025-11-22 19:33:08.128', '2025-11-22 19:34:25.711', '2025-11-22 19:34:41.139', '', '', '', 4, 1, 3, '2025-11-22 19:36:09.589', '2025-11-22 19:41:18.736', 591, 4, '2025-11-22 19:31:57.732', '2025-11-22 19:41:18.789', NULL);
INSERT INTO `app_service_release` (`id`, `title`, `business_group_id`, `impact_feature`, `applicant_id`, `applicant_name`, `owner_approver_id`, `owner_approver_name`, `security_approver_id`, `security_approver_name`, `test_approver_id`, `test_approver_name`, `owner_approval_status`, `security_approval_status`, `test_approval_status`, `owner_approval_time`, `security_approval_time`, `test_approval_time`, `owner_approval_remark`, `security_approval_remark`, `test_approval_remark`, `deploy_status`, `regression_test_status`, `status`, `start_time`, `end_time`, `duration`, `service_count`, `created_at`, `updated_at`, `deleted_at`) VALUES (4, '测试审批流程', 12, '测试审批流程', 1, '管理员', 103, '用户_103', 98, '用户_98', 104, '用户_104', 2, 2, 2, '2025-11-22 20:28:35.898', '2025-11-22 20:28:48.373', '2025-11-22 20:28:55.066', '', '', '', 4, 1, 3, '2025-11-22 20:29:16.698', '2025-11-22 20:34:22.328', 303, 2, '2025-11-22 20:00:18.967', '2025-11-22 20:34:22.386', NULL);
INSERT INTO `app_service_release` (`id`, `title`, `business_group_id`, `impact_feature`, `applicant_id`, `applicant_name`, `owner_approver_id`, `owner_approver_name`, `security_approver_id`, `security_approver_name`, `test_approver_id`, `test_approver_name`, `owner_approval_status`, `security_approval_status`, `test_approval_status`, `owner_approval_time`, `security_approval_time`, `test_approval_time`, `owner_approval_remark`, `security_approval_remark`, `test_approval_remark`, `deploy_status`, `regression_test_status`, `status`, `start_time`, `end_time`, `duration`, `service_count`, `created_at`, `updated_at`, `deleted_at`) VALUES (5, '测试---发布工单上线', 12, '测试---发布工单上线', 1, '管理员', 104, '用户_104', 103, '用户_103', 98, '用户_98', 2, 2, 2, '2025-11-22 23:13:56.636', '2025-11-22 23:16:26.840', '2025-11-22 23:16:34.126', '', '通过', '', 3, 3, 2, '2025-11-22 23:52:17.626', '2025-11-22 23:55:51.883', 212, 1, '2025-11-22 23:13:39.884', '2025-11-23 00:50:11.059', NULL);
INSERT INTO `app_service_release` (`id`, `title`, `business_group_id`, `impact_feature`, `applicant_id`, `applicant_name`, `owner_approver_id`, `owner_approver_name`, `security_approver_id`, `security_approver_name`, `test_approver_id`, `test_approver_name`, `owner_approval_status`, `security_approval_status`, `test_approval_status`, `owner_approval_time`, `security_approval_time`, `test_approval_time`, `owner_approval_remark`, `security_approval_remark`, `test_approval_remark`, `deploy_status`, `regression_test_status`, `status`, `start_time`, `end_time`, `duration`, `service_count`, `created_at`, `updated_at`, `deleted_at`) VALUES (6, '双任务上线', 14, '123', 1, '管理员', 104, '用户_104', 103, '用户_103', 98, '用户_98', 2, 2, 2, '2025-11-23 11:33:16.353', '2025-11-23 11:33:22.967', '2025-11-23 12:57:50.091', '', '', '111', 1, 1, 1, NULL, NULL, 0, 2, '2025-11-23 11:32:19.689', '2025-11-23 12:57:50.147', NULL);
INSERT INTO `app_service_release` (`id`, `title`, `business_group_id`, `impact_feature`, `applicant_id`, `applicant_name`, `owner_approver_id`, `owner_approver_name`, `security_approver_id`, `security_approver_name`, `test_approver_id`, `test_approver_name`, `owner_approval_status`, `security_approval_status`, `test_approval_status`, `owner_approval_time`, `security_approval_time`, `test_approval_time`, `owner_approval_remark`, `security_approval_remark`, `test_approval_remark`, `deploy_status`, `regression_test_status`, `status`, `start_time`, `end_time`, `duration`, `service_count`, `created_at`, `updated_at`, `deleted_at`) VALUES (7, '测试生产环境服务上线', 18, '影响用户登录', 1, '管理员', 104, '王五', 98, '李四', 98, '李四', 2, 2, 2, '2025-11-24 00:44:47.250', '2025-11-24 00:45:40.169', '2025-11-24 00:45:45.926', '同意', '', '', 3, 3, 2, '2025-11-24 00:46:52.721', '2025-11-24 00:52:38.168', 344, 1, '2025-11-24 00:42:24.608', '2025-11-24 00:53:04.351', NULL);
INSERT INTO `app_service_release` (`id`, `title`, `business_group_id`, `impact_feature`, `applicant_id`, `applicant_name`, `owner_approver_id`, `owner_approver_name`, `security_approver_id`, `security_approver_name`, `test_approver_id`, `test_approver_name`, `owner_approval_status`, `security_approval_status`, `test_approval_status`, `owner_approval_time`, `security_approval_time`, `test_approval_time`, `owner_approval_remark`, `security_approval_remark`, `test_approval_remark`, `deploy_status`, `regression_test_status`, `status`, `start_time`, `end_time`, `duration`, `service_count`, `created_at`, `updated_at`, `deleted_at`) VALUES (8, 'tesst123', 12, '影响用户登录', 1, '管理员', 89, 'admin', 89, 'admin', 89, 'admin', 2, 2, 2, '2025-12-05 20:47:13.905', '2025-12-05 20:47:18.461', '2025-12-05 20:47:27.382', 'ok', '', '', 3, 2, 1, '2025-12-05 20:47:39.670', '2025-12-05 20:47:53.645', 12, 1, '2025-12-05 20:46:53.205', '2025-12-05 20:47:53.696', NULL);
COMMIT;

-- ----------------------------
-- Table structure for app_service_release_item
-- ----------------------------
DROP TABLE IF EXISTS `app_service_release_item`;
CREATE TABLE `app_service_release_item` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `release_id` bigint unsigned NOT NULL COMMENT '上线工单ID',
  `app_id` bigint unsigned NOT NULL COMMENT '应用ID',
  `app_name` varchar(255) NOT NULL COMMENT '应用名称',
  `app_code` varchar(100) NOT NULL COMMENT '应用编码',
  `project_name` varchar(255) NOT NULL COMMENT '项目服务名称',
  `repo_url` varchar(500) NOT NULL COMMENT '项目地址',
  `branch` varchar(100) DEFAULT 'master' COMMENT '发布分支',
  `commit_id` varchar(100) NOT NULL COMMENT 'Commit ID',
  `impact_feature` text COMMENT '影响功能',
  `function_module` text COMMENT '功能模块',
  `db_change` text COMMENT '数据库变更',
  `config_change` text COMMENT '配置变更',
  `remark` text COMMENT '备注信息',
  `jenkins_env_id` bigint unsigned DEFAULT NULL COMMENT 'Jenkins环境配置ID',
  `jenkins_job_url` varchar(500) DEFAULT NULL COMMENT 'Jenkins任务URL',
  `parameters` text COMMENT 'Jenkins构建参数(JSON格式)',
  `build_number` bigint DEFAULT '0' COMMENT '构建编号',
  `log_url` varchar(500) DEFAULT NULL COMMENT '构建日志URL',
  `status` bigint DEFAULT '1' COMMENT '发布状态',
  `start_time` datetime(3) DEFAULT NULL COMMENT '发布开始时间',
  `end_time` datetime(3) DEFAULT NULL COMMENT '发布结束时间',
  `duration` bigint DEFAULT '0' COMMENT '发布耗时(秒)',
  `error_message` text COMMENT '错误信息',
  `execute_order` bigint DEFAULT '0' COMMENT '执行顺序',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_app_service_release_item_release_id` (`release_id`),
  KEY `idx_app_service_release_item_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of app_service_release_item
-- ----------------------------
BEGIN;
INSERT INTO `app_service_release_item` (`id`, `release_id`, `app_id`, `app_name`, `app_code`, `project_name`, `repo_url`, `branch`, `commit_id`, `impact_feature`, `function_module`, `db_change`, `config_change`, `remark`, `jenkins_env_id`, `jenkins_job_url`, `parameters`, `build_number`, `log_url`, `status`, `start_time`, `end_time`, `duration`, `error_message`, `execute_order`, `created_at`, `updated_at`, `deleted_at`) VALUES (1, 2, 14, 'dev-lku-sass-pack', 'dev-lku-sass-pack', 'dev-lku-sass-pack', 'git@code.dding.net/lockin/saas-toc-server.git', 'master', '083a2d6669effe63670b0d5ee899d29eeb187abf', '影响用户登录', '', '没有', '没有', '没有', 35, '', '{\"commit_id\":\"083a2d6669effe63670b0d5ee899d29eeb187abf\",\"compile\":\"true\"}', 408, 'http://test-ops-jenkins-tc2.dding.net:8080/job/dev-lku-sass-pack/408/console', 3, '2025-11-22 18:18:03.823', '2025-11-22 18:21:44.766', 220, '', 1, '2025-11-22 18:13:37.484', '2025-11-22 18:21:44.821', NULL);
INSERT INTO `app_service_release_item` (`id`, `release_id`, `app_id`, `app_name`, `app_code`, `project_name`, `repo_url`, `branch`, `commit_id`, `impact_feature`, `function_module`, `db_change`, `config_change`, `remark`, `jenkins_env_id`, `jenkins_job_url`, `parameters`, `build_number`, `log_url`, `status`, `start_time`, `end_time`, `duration`, `error_message`, `execute_order`, `created_at`, `updated_at`, `deleted_at`) VALUES (2, 3, 15, 'dev-lku-sass-deploy', 'dev-lku-sass-deploy', 'dev-lku-sass-deploy', 'git@code.dding.net/lockin/saas-toc-server.git', 'master', 'c7417cbe54a11d89c819b4541316980dc1634687', '测试多任务工单发布', '', '无', '无', '无', 39, '', '', 19, 'http://test-ops-jenkins-tc2.dding.net:8080/job/test1-saas-deploy/19/console', 4, '2025-11-22 19:36:10.184', '2025-11-22 19:36:37.458', 27, 'Jenkins构建失败', 1, '2025-11-22 19:31:58.231', '2025-11-22 19:36:37.524', NULL);
INSERT INTO `app_service_release_item` (`id`, `release_id`, `app_id`, `app_name`, `app_code`, `project_name`, `repo_url`, `branch`, `commit_id`, `impact_feature`, `function_module`, `db_change`, `config_change`, `remark`, `jenkins_env_id`, `jenkins_job_url`, `parameters`, `build_number`, `log_url`, `status`, `start_time`, `end_time`, `duration`, `error_message`, `execute_order`, `created_at`, `updated_at`, `deleted_at`) VALUES (3, 3, 14, 'dev-lku-sass-pack', 'dev-lku-sass-pack', 'dev-lku-sass-pack', 'git@code.dding.net/lockin/saas-toc-server.git', 'master', 'c7417cbe54a11d89c819b4541316980dc1634687', '测试多任务工单发布', '', '无', '无', '无', 35, '', '', 409, 'http://test-ops-jenkins-tc2.dding.net:8080/job/dev-lku-sass-pack/409/console', 3, '2025-11-22 19:36:10.184', '2025-11-22 19:39:59.417', 229, '', 2, '2025-11-22 19:31:58.710', '2025-11-22 19:39:59.476', NULL);
INSERT INTO `app_service_release_item` (`id`, `release_id`, `app_id`, `app_name`, `app_code`, `project_name`, `repo_url`, `branch`, `commit_id`, `impact_feature`, `function_module`, `db_change`, `config_change`, `remark`, `jenkins_env_id`, `jenkins_job_url`, `parameters`, `build_number`, `log_url`, `status`, `start_time`, `end_time`, `duration`, `error_message`, `execute_order`, `created_at`, `updated_at`, `deleted_at`) VALUES (4, 3, 18, 'test1-paas-deploy', 'test1-paas-deploy', 'test1-paas-deploy', 'git@code.dding.net/lockin/cloud-platform.git', 'master', 'c7417cbe54a11d89c819b4541316980dc1634687', '测试多任务工单发布', '', '无', '无', '无', 48, '', '', 13, 'http://test-ops-jenkins-tc2.dding.net:8080/job/test1-paas-deploy/13/console', 4, '2025-11-22 19:36:10.184', '2025-11-22 19:36:38.200', 28, 'Jenkins构建失败', 3, '2025-11-22 19:31:59.196', '2025-11-22 19:36:38.282', NULL);
INSERT INTO `app_service_release_item` (`id`, `release_id`, `app_id`, `app_name`, `app_code`, `project_name`, `repo_url`, `branch`, `commit_id`, `impact_feature`, `function_module`, `db_change`, `config_change`, `remark`, `jenkins_env_id`, `jenkins_job_url`, `parameters`, `build_number`, `log_url`, `status`, `start_time`, `end_time`, `duration`, `error_message`, `execute_order`, `created_at`, `updated_at`, `deleted_at`) VALUES (5, 3, 17, 'test1-paas-pack', 'test1-paas-pack', 'test1-paas-pack', 'git@code.dding.net/lockin/cloud-platform.git', 'master', 'c7417cbe54a11d89c819b4541316980dc1634687', '测试多任务工单发布', '', '无', '无', '无', 45, '', '', 14, 'http://test-ops-jenkins-tc2.dding.net:8080/job/test1-paas-pack/14/console', 3, '2025-11-22 19:36:10.184', '2025-11-22 19:41:18.055', 307, '', 4, '2025-11-22 19:31:59.673', '2025-11-22 19:41:18.106', NULL);
INSERT INTO `app_service_release_item` (`id`, `release_id`, `app_id`, `app_name`, `app_code`, `project_name`, `repo_url`, `branch`, `commit_id`, `impact_feature`, `function_module`, `db_change`, `config_change`, `remark`, `jenkins_env_id`, `jenkins_job_url`, `parameters`, `build_number`, `log_url`, `status`, `start_time`, `end_time`, `duration`, `error_message`, `execute_order`, `created_at`, `updated_at`, `deleted_at`) VALUES (6, 4, 17, 'test1-paas-pack', 'test1-paas-pack', 'test1-paas-pack', 'git@code.dding.net/lockin/cloud-platform.git', 'master', 'c7417cbe54a11d89c819b4541316980dc1634687', '测试审批流程', '', '测试审批流程', '测试审批流程', '测试审批流程', 45, '', '', 15, 'http://test-ops-jenkins-tc2.dding.net:8080/job/test1-paas-pack/15/console', 3, '2025-11-22 20:29:17.308', '2025-11-22 20:34:08.653', 291, '', 1, '2025-11-22 20:00:19.453', '2025-11-22 20:34:08.746', NULL);
INSERT INTO `app_service_release_item` (`id`, `release_id`, `app_id`, `app_name`, `app_code`, `project_name`, `repo_url`, `branch`, `commit_id`, `impact_feature`, `function_module`, `db_change`, `config_change`, `remark`, `jenkins_env_id`, `jenkins_job_url`, `parameters`, `build_number`, `log_url`, `status`, `start_time`, `end_time`, `duration`, `error_message`, `execute_order`, `created_at`, `updated_at`, `deleted_at`) VALUES (7, 4, 18, 'test1-paas-deploy', 'test1-paas-deploy', 'test1-paas-deploy', 'git@code.dding.net/lockin/cloud-platform.git', 'master', 'c7417cbe54a11d89c819b4541316980dc1634687', '测试审批流程', '', '测试审批流程', '测试审批流程', '测试审批流程', 48, '', '', 14, 'http://test-ops-jenkins-tc2.dding.net:8080/job/test1-paas-deploy/14/console', 4, '2025-11-22 20:34:09.058', '2025-11-22 20:34:21.625', 12, 'Jenkins构建失败', 2, '2025-11-22 20:00:19.908', '2025-11-22 20:34:21.687', NULL);
INSERT INTO `app_service_release_item` (`id`, `release_id`, `app_id`, `app_name`, `app_code`, `project_name`, `repo_url`, `branch`, `commit_id`, `impact_feature`, `function_module`, `db_change`, `config_change`, `remark`, `jenkins_env_id`, `jenkins_job_url`, `parameters`, `build_number`, `log_url`, `status`, `start_time`, `end_time`, `duration`, `error_message`, `execute_order`, `created_at`, `updated_at`, `deleted_at`) VALUES (8, 5, 14, 'dev-lku-sass-pack', 'dev-lku-sass-pack', 'dev-lku-sass-pack', 'git@code.dding.net/lockin/saas-toc-server.git', 'master', '083a2d6669effe63670b0d5ee899d29eeb187abf', '无', '', '无', '无', '无', 35, '', '', 410, 'http://test-ops-jenkins-tc2.dding.net:8080/job/dev-lku-sass-pack/410/console', 3, '2025-11-22 23:52:18.208', '2025-11-22 23:55:51.178', 212, '', 1, '2025-11-22 23:13:40.348', '2025-11-22 23:55:51.229', NULL);
INSERT INTO `app_service_release_item` (`id`, `release_id`, `app_id`, `app_name`, `app_code`, `project_name`, `repo_url`, `branch`, `commit_id`, `impact_feature`, `function_module`, `db_change`, `config_change`, `remark`, `jenkins_env_id`, `jenkins_job_url`, `parameters`, `build_number`, `log_url`, `status`, `start_time`, `end_time`, `duration`, `error_message`, `execute_order`, `created_at`, `updated_at`, `deleted_at`) VALUES (9, 6, 17, 'test1-paas-pack', 'test1-paas-pack', 'test1-paas-pack', 'git@code.dding.net/lockin/cloud-platform.git', 'master', 'c7417cbe54a11d89c819b4541316980dc1634687', '123', '', '123', '123', '123', 45, '', '{\"commit_id\":\"c7417cbe54a11d89c819b4541316980dc1634687\",\"compile\":\"true\"}', 0, '', 1, NULL, NULL, 0, '', 1, '2025-11-23 11:32:20.185', '2025-11-23 11:32:20.185', NULL);
INSERT INTO `app_service_release_item` (`id`, `release_id`, `app_id`, `app_name`, `app_code`, `project_name`, `repo_url`, `branch`, `commit_id`, `impact_feature`, `function_module`, `db_change`, `config_change`, `remark`, `jenkins_env_id`, `jenkins_job_url`, `parameters`, `build_number`, `log_url`, `status`, `start_time`, `end_time`, `duration`, `error_message`, `execute_order`, `created_at`, `updated_at`, `deleted_at`) VALUES (10, 6, 14, 'dev-lku-sass-pack', 'dev-lku-sass-pack', 'dev-lku-sass-pack', 'git@code.dding.net/lockin/saas-toc-server.git', 'master', '083a2d6669effe63670b0d5ee899d29eeb187abf', '123', '', '123', '123', '123', 35, '', '{\"commit_id\":\"083a2d6669effe63670b0d5ee899d29eeb187abf\",\"compile\":\"true\"}', 0, '', 1, NULL, NULL, 0, '', 2, '2025-11-23 11:32:20.664', '2025-11-23 11:32:20.664', NULL);
INSERT INTO `app_service_release_item` (`id`, `release_id`, `app_id`, `app_name`, `app_code`, `project_name`, `repo_url`, `branch`, `commit_id`, `impact_feature`, `function_module`, `db_change`, `config_change`, `remark`, `jenkins_env_id`, `jenkins_job_url`, `parameters`, `build_number`, `log_url`, `status`, `start_time`, `end_time`, `duration`, `error_message`, `execute_order`, `created_at`, `updated_at`, `deleted_at`) VALUES (11, 7, 14, 'dev-lku-sass-pack', 'dev-lku-sass-pack', 'dev-lku-sass-pack', 'git@code.dding.net/lockin/saas-toc-server.git', 'master', '083a2d6669effe63670b0d5ee899d29eeb187abf', '影响用户登录', '', '无', '无', '123', 35, '', '{\"commit_id\":\"083a2d6669effe63670b0d5ee899d29eeb187abf\",\"compile\":\"true\"}', 411, 'http://test-ops-jenkins-tc2.dding.net:8080/job/dev-lku-sass-pack/411/console', 3, '2025-11-24 00:46:53.314', '2025-11-24 00:52:37.448', 344, '', 1, '2025-11-24 00:42:25.108', '2025-11-24 00:52:37.512', NULL);
INSERT INTO `app_service_release_item` (`id`, `release_id`, `app_id`, `app_name`, `app_code`, `project_name`, `repo_url`, `branch`, `commit_id`, `impact_feature`, `function_module`, `db_change`, `config_change`, `remark`, `jenkins_env_id`, `jenkins_job_url`, `parameters`, `build_number`, `log_url`, `status`, `start_time`, `end_time`, `duration`, `error_message`, `execute_order`, `created_at`, `updated_at`, `deleted_at`) VALUES (12, 8, 19, 'prod_saas3.0_data-export', 'prod-saas30-data-export', 'prod_saas3.0_data-export', 'git@gitee.com:zhang_fan1024/zf-k8s.git', 'master', '', '影响用户登录', '', '无', '无', '无', 51, '', '{\"commit_id\":\"123456789\"}', 10, 'http://test-ops-jenkins-tc2.dding.net:8080/job/prod_saas3.0_data-export/10/console', 3, '2025-12-05 20:47:40.229', '2025-12-05 20:47:52.992', 12, '', 1, '2025-12-05 20:46:53.654', '2025-12-05 20:47:53.043', NULL);
COMMIT;

-- ----------------------------
-- Table structure for app_sh_release
-- ----------------------------
DROP TABLE IF EXISTS `app_sh_release`;
CREATE TABLE `app_sh_release` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL COMMENT '上线标题',
  `reason` text NOT NULL COMMENT '上线原因',
  `business_group_id` bigint unsigned NOT NULL COMMENT '业务线ID',
  `app_id` bigint unsigned NOT NULL COMMENT '服务ID',
  `app_name` varchar(255) NOT NULL COMMENT '服务名称',
  `app_code` varchar(100) NOT NULL COMMENT '服务编码',
  `applicant_id` bigint unsigned NOT NULL COMMENT '申请人ID',
  `applicant_name` varchar(100) NOT NULL COMMENT '申请人姓名',
  `approver_id` bigint unsigned DEFAULT NULL COMMENT '审核人ID',
  `approver_name` varchar(100) DEFAULT NULL COMMENT '审核人姓名',
  `executor_id` bigint unsigned DEFAULT NULL COMMENT '执行人ID',
  `executor_name` varchar(100) DEFAULT NULL COMMENT '执行人姓名',
  `execute_dir` varchar(500) NOT NULL COMMENT '执行目录',
  `script_content` text NOT NULL COMMENT '脚本内容',
  `approval_status` bigint DEFAULT '1' COMMENT '审核状态',
  `approval_time` datetime(3) DEFAULT NULL COMMENT '审核时间',
  `approval_remark` text COMMENT '审核意见',
  `execute_status` bigint DEFAULT '1' COMMENT '执行状态',
  `status` bigint DEFAULT '1' COMMENT '流程状态',
  `start_time` datetime(3) DEFAULT NULL COMMENT '脚本执行开始时间',
  `end_time` datetime(3) DEFAULT NULL COMMENT '脚本执行结束时间',
  `duration` bigint DEFAULT '0' COMMENT '执行耗时(秒)',
  `jenkins_env_id` bigint unsigned DEFAULT NULL COMMENT 'Jenkins环境配置ID',
  `build_number` bigint DEFAULT '0' COMMENT '构建编号',
  `log_url` varchar(500) DEFAULT NULL COMMENT '构建日志URL',
  `error_message` text COMMENT '错误信息',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `parameters` text COMMENT 'Jenkins构建参数(JSON格式)',
  `server_host_id` bigint unsigned NOT NULL COMMENT '执行服务器主机ID(关联cmdb_host)',
  `pull_code_start_time` datetime(3) DEFAULT NULL COMMENT '拉取代码开始时间',
  `pull_code_end_time` datetime(3) DEFAULT NULL COMMENT '拉取代码结束时间',
  `script_output` longtext COMMENT '脚本执行输出',
  PRIMARY KEY (`id`),
  KEY `idx_app_sh_release_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of app_sh_release
-- ----------------------------
BEGIN;
INSERT INTO `app_sh_release` (`id`, `title`, `reason`, `business_group_id`, `app_id`, `app_name`, `app_code`, `applicant_id`, `applicant_name`, `approver_id`, `approver_name`, `executor_id`, `executor_name`, `execute_dir`, `script_content`, `approval_status`, `approval_time`, `approval_remark`, `execute_status`, `status`, `start_time`, `end_time`, `duration`, `jenkins_env_id`, `build_number`, `log_url`, `error_message`, `created_at`, `updated_at`, `deleted_at`, `parameters`, `server_host_id`, `pull_code_start_time`, `pull_code_end_time`, `script_output`) VALUES (1, 'prod_saas3.0_data-export-数据导出', 'prod_saas3.0_data-export-数据导出', 18, 19, 'prod_saas3.0_data-export', 'prod-saas30-data-export', 1, '管理员', 89, 'admin', 89, 'admin', '/home/dingding/saas3-data-export/', 'ls', 2, '2025-11-28 12:42:48.493', '111', 3, 2, '2025-11-28 12:43:08.943', '2025-11-28 12:43:24.437', 14, 51, 3, 'http://test-ops-jenkins-tc2.dding.net:8080/job/prod_saas3.0_data-export/3/console', '', '2025-11-28 12:41:29.167', '2025-11-28 12:43:24.488', NULL, '{\"commit_id\":\"123456789\"}', 0, NULL, NULL, NULL);
INSERT INTO `app_sh_release` (`id`, `title`, `reason`, `business_group_id`, `app_id`, `app_name`, `app_code`, `applicant_id`, `applicant_name`, `approver_id`, `approver_name`, `executor_id`, `executor_name`, `execute_dir`, `script_content`, `approval_status`, `approval_time`, `approval_remark`, `execute_status`, `status`, `start_time`, `end_time`, `duration`, `jenkins_env_id`, `build_number`, `log_url`, `error_message`, `created_at`, `updated_at`, `deleted_at`, `parameters`, `server_host_id`, `pull_code_start_time`, `pull_code_end_time`, `script_output`) VALUES (2, '测试001', '测试001', 18, 19, 'prod_saas3.0_data-export', 'prod-saas30-data-export', 1, '管理员', 98, '李四', 89, 'admin', '/home/', 'ls  /root/', 2, '2025-11-28 14:18:58.866', '同意', 3, 1, NULL, NULL, 0, 51, 4, 'http://test-ops-jenkins-tc2.dding.net:8080/job/prod_saas3.0_data-export/4/console', '', '2025-11-28 13:37:27.963', '2025-11-28 14:47:38.777', NULL, '{\"commit_id\":\"123456789\"}', 0, '2025-11-28 14:47:22.389', '2025-11-28 14:47:38.725', NULL);
INSERT INTO `app_sh_release` (`id`, `title`, `reason`, `business_group_id`, `app_id`, `app_name`, `app_code`, `applicant_id`, `applicant_name`, `approver_id`, `approver_name`, `executor_id`, `executor_name`, `execute_dir`, `script_content`, `approval_status`, `approval_time`, `approval_remark`, `execute_status`, `status`, `start_time`, `end_time`, `duration`, `jenkins_env_id`, `build_number`, `log_url`, `error_message`, `created_at`, `updated_at`, `deleted_at`, `parameters`, `server_host_id`, `pull_code_start_time`, `pull_code_end_time`, `script_output`) VALUES (3, '测试脚本执行', '测试脚本执行', 18, 19, 'prod_saas3.0_data-export', 'prod-saas30-data-export', 1, '管理员', 89, 'admin', 89, 'admin', '/home/dingding/saas3-data-export/', 'pwd\nls\nhostname -I', 2, '2025-11-28 15:19:15.004', '111', 6, 2, '2025-11-28 15:24:27.057', '2025-11-28 15:24:31.379', 3, 51, 5, 'http://test-ops-jenkins-tc2.dding.net:8080/job/prod_saas3.0_data-export/5/console', '', '2025-11-28 15:09:18.474', '2025-11-28 15:24:31.432', NULL, '{\"commit_id\":\"123456789\"}', 501, '2025-11-28 15:19:21.580', '2025-11-28 15:19:34.271', 'bash: line 1: cd: /home/dingding/saas3-data-export/: No such file or directory\ndocker\nelk\njdk11\njdk17\njdk18\nluban-master\nnode\nprometheus\nsnap\n172.20.236.121 172.18.0.1 172.17.0.1 172.19.0.1 \n');
INSERT INTO `app_sh_release` (`id`, `title`, `reason`, `business_group_id`, `app_id`, `app_name`, `app_code`, `applicant_id`, `applicant_name`, `approver_id`, `approver_name`, `executor_id`, `executor_name`, `execute_dir`, `script_content`, `approval_status`, `approval_time`, `approval_remark`, `execute_status`, `status`, `start_time`, `end_time`, `duration`, `jenkins_env_id`, `build_number`, `log_url`, `error_message`, `created_at`, `updated_at`, `deleted_at`, `parameters`, `server_host_id`, `pull_code_start_time`, `pull_code_end_time`, `script_output`) VALUES (4, '测试脚本002', '测试脚本002', 18, 19, 'prod_saas3.0_data-export', 'prod-saas30-data-export', 1, '管理员', 89, 'admin', 89, 'admin', '/home/dingding/saas3-data-export/', 'hostname\npwd\ndate', 2, '2025-11-28 15:56:35.290', '11', 6, 2, '2025-11-28 15:57:23.260', '2025-11-28 15:57:27.569', 3, 51, 6, 'http://test-ops-jenkins-tc2.dding.net:8080/job/prod_saas3.0_data-export/6/console', '', '2025-11-28 15:56:13.417', '2025-11-28 15:57:27.644', NULL, '{\"commit_id\":\"123456789\"}', 501, '2025-11-28 15:56:44.717', '2025-11-28 15:57:00.303', '/root\nFri Nov 28 03:57:27 PM CST 2025\nbash: line 1: cd: /home/dingding/saas3-data-export/: No such file or directory\n');
INSERT INTO `app_sh_release` (`id`, `title`, `reason`, `business_group_id`, `app_id`, `app_name`, `app_code`, `applicant_id`, `applicant_name`, `approver_id`, `approver_name`, `executor_id`, `executor_name`, `execute_dir`, `script_content`, `approval_status`, `approval_time`, `approval_remark`, `execute_status`, `status`, `start_time`, `end_time`, `duration`, `jenkins_env_id`, `build_number`, `log_url`, `error_message`, `created_at`, `updated_at`, `deleted_at`, `parameters`, `server_host_id`, `pull_code_start_time`, `pull_code_end_time`, `script_output`) VALUES (5, '测试002', '测试002', 18, 19, 'prod_saas3.0_data-export', 'prod-saas30-data-export', 1, '管理员', 89, 'admin', 89, 'admin', '/home/', 'pwd\nls\nhostname\n', 2, '2025-11-28 16:12:04.059', 'ok', 6, 2, '2025-11-28 16:19:18.636', '2025-11-28 16:19:22.516', 3, 51, 7, 'http://test-ops-jenkins-tc2.dding.net:8080/job/prod_saas3.0_data-export/7/console', '', '2025-11-28 16:10:22.217', '2025-11-28 16:19:22.570', NULL, '{\"commit_id\":\"123456789\"}', 501, '2025-11-28 16:12:11.076', '2025-11-28 16:12:23.702', '/home\ndevops\ndevops.tar.gz\ngo-ops\n');
INSERT INTO `app_sh_release` (`id`, `title`, `reason`, `business_group_id`, `app_id`, `app_name`, `app_code`, `applicant_id`, `applicant_name`, `approver_id`, `approver_name`, `executor_id`, `executor_name`, `execute_dir`, `script_content`, `approval_status`, `approval_time`, `approval_remark`, `execute_status`, `status`, `start_time`, `end_time`, `duration`, `jenkins_env_id`, `build_number`, `log_url`, `error_message`, `created_at`, `updated_at`, `deleted_at`, `parameters`, `server_host_id`, `pull_code_start_time`, `pull_code_end_time`, `script_output`) VALUES (6, 'test1111111111', 'test', 18, 19, 'prod_saas3.0_data-export', 'prod-saas30-data-export', 1, '管理员', 89, 'admin', 89, 'admin', '/home/dingding/saas3-data-export/', 'pwd  \nls \nhostname', 2, '2025-12-01 01:01:48.329', 'ok', 6, 2, '2025-12-01 01:02:36.996', '2025-12-01 01:02:39.204', 1, 51, 8, 'http://test-ops-jenkins-tc2.dding.net:8080/job/prod_saas3.0_data-export/8/console', '', '2025-12-01 01:01:08.782', '2025-12-01 01:02:39.260', NULL, '{\"commit_id\":\"123456789\"}', 501, '2025-12-01 01:01:59.779', '2025-12-01 01:02:16.154', '/home/dingding/saas3-data-export\nxlsx\ngo-ops\n');
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
) ENGINE=InnoDB AUTO_INCREMENT=37 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

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
INSERT INTO `cmdb_group` (`id`, `parent_id`, `name`, `create_time`, `remark`, `update_time`) VALUES (12, 2, 'saas3', '2025-07-17 11:07:19.279', NULL, NULL);
INSERT INTO `cmdb_group` (`id`, `parent_id`, `name`, `create_time`, `remark`, `update_time`) VALUES (14, 3, '智慧门锁业务', '2025-07-17 11:19:14.090', NULL, NULL);
INSERT INTO `cmdb_group` (`id`, `parent_id`, `name`, `create_time`, `remark`, `update_time`) VALUES (16, 15, '智能家居业务', '2025-07-22 10:32:14.837', NULL, NULL);
INSERT INTO `cmdb_group` (`id`, `parent_id`, `name`, `create_time`, `remark`, `update_time`) VALUES (18, 2, 'saas4', '2025-07-23 10:40:12.787', NULL, NULL);
INSERT INTO `cmdb_group` (`id`, `parent_id`, `name`, `create_time`, `remark`, `update_time`) VALUES (19, 17, '运维组', '2025-07-23 20:21:08.618', NULL, NULL);
INSERT INTO `cmdb_group` (`id`, `parent_id`, `name`, `create_time`, `remark`, `update_time`) VALUES (20, 17, '安全组', '2025-07-23 20:21:33.222', NULL, NULL);
INSERT INTO `cmdb_group` (`id`, `parent_id`, `name`, `create_time`, `remark`, `update_time`) VALUES (26, 24, '国内鹿客云', '2025-09-07 18:04:56.683', NULL, NULL);
INSERT INTO `cmdb_group` (`id`, `parent_id`, `name`, `create_time`, `remark`, `update_time`) VALUES (27, 24, '新加坡鹿客云', '2025-09-07 18:05:13.713', NULL, NULL);
INSERT INTO `cmdb_group` (`id`, `parent_id`, `name`, `create_time`, `remark`, `update_time`) VALUES (28, 24, '北美鹿客云', '2025-09-07 18:05:31.846', NULL, NULL);
INSERT INTO `cmdb_group` (`id`, `parent_id`, `name`, `create_time`, `remark`, `update_time`) VALUES (29, 17, 'ai运维组', '2025-10-01 15:17:35.617', NULL, NULL);
INSERT INTO `cmdb_group` (`id`, `parent_id`, `name`, `create_time`, `remark`, `update_time`) VALUES (31, 30, 'rental-test', '2025-12-05 20:17:32.169', NULL, NULL);
INSERT INTO `cmdb_group` (`id`, `parent_id`, `name`, `create_time`, `remark`, `update_time`) VALUES (32, 30, 'rental-dev', '2025-12-05 20:17:43.988', NULL, NULL);
INSERT INTO `cmdb_group` (`id`, `parent_id`, `name`, `create_time`, `remark`, `update_time`) VALUES (34, 33, 'saas7', '2025-12-09 21:14:55.863', NULL, NULL);
INSERT INTO `cmdb_group` (`id`, `parent_id`, `name`, `create_time`, `remark`, `update_time`) VALUES (35, 33, 'saas8', '2025-12-09 21:15:05.050', NULL, NULL);
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
  `tag` text COMMENT '''标签(格式:key=value,key=value)''',
  PRIMARY KEY (`id`),
  KEY `fk_cmdb_group_hosts` (`group_id`),
  CONSTRAINT `fk_cmdb_group_hosts` FOREIGN KEY (`group_id`) REFERENCES `cmdb_group` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=526 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

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
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of config_account
-- ----------------------------
BEGIN;
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
) ENGINE=InnoDB AUTO_INCREMENT=24 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of config_ecsauth
-- ----------------------------
BEGIN;
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
-- Table structure for db
-- ----------------------------
DROP TABLE IF EXISTS `db`;
CREATE TABLE `db` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `code` varchar(36) NOT NULL COMMENT '''数据库编码''',
  `name` varchar(100) NOT NULL COMMENT '''数据库名称''',
  `database` varchar(500) NOT NULL COMMENT '''数据库名(多个用空格分隔)''',
  `remark` varchar(500) DEFAULT NULL COMMENT '''备注''',
  `instance_id` bigint unsigned NOT NULL COMMENT '''实例ID''',
  `instance_code` varchar(36) DEFAULT NULL COMMENT '''实例编码''',
  `status` bigint DEFAULT '1' COMMENT '''状态:1->启用,2->禁用''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `update_time` datetime(3) DEFAULT NULL COMMENT '''更新时间''',
  `creator` varchar(64) DEFAULT NULL COMMENT '''创建人''',
  `creator_id` bigint unsigned DEFAULT NULL COMMENT '''创建人ID''',
  `modifier` varchar(64) DEFAULT NULL COMMENT '''修改人''',
  `modifier_id` bigint unsigned DEFAULT NULL COMMENT '''修改人ID''',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_db_code` (`code`),
  KEY `idx_db_instance_id` (`instance_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of db
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for db_es_instance
-- ----------------------------
DROP TABLE IF EXISTS `db_es_instance`;
CREATE TABLE `db_es_instance` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `code` varchar(36) NOT NULL COMMENT '''实例编码''',
  `name` varchar(100) NOT NULL COMMENT '''实例名称''',
  `protocol` varchar(10) DEFAULT 'http' COMMENT '''协议:http|https''',
  `host` varchar(255) NOT NULL COMMENT '''主机''',
  `port` bigint DEFAULT '9200' COMMENT '''端口''',
  `username` varchar(100) DEFAULT NULL COMMENT '''用户名(可选)''',
  `password` varchar(500) DEFAULT '' COMMENT '''密码(加密)''',
  `remark` varchar(500) DEFAULT NULL COMMENT '''备注''',
  `ssh_tunnel_machine_id` bigint unsigned DEFAULT '0' COMMENT '''SSH隧道机器ID''',
  `status` bigint DEFAULT '1' COMMENT '''状态:1->启用,2->禁用''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `update_time` datetime(3) DEFAULT NULL COMMENT '''更新时间''',
  `creator` varchar(64) DEFAULT NULL COMMENT '''创建人''',
  `creator_id` bigint unsigned DEFAULT NULL COMMENT '''创建人ID''',
  `modifier` varchar(64) DEFAULT NULL COMMENT '''修改人''',
  `modifier_id` bigint unsigned DEFAULT NULL COMMENT '''修改人ID''',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_db_es_instance_code` (`code`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of db_es_instance
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for db_export_task
-- ----------------------------
DROP TABLE IF EXISTS `db_export_task`;
CREATE TABLE `db_export_task` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `task_id` varchar(36) NOT NULL COMMENT '''任务ID''',
  `db_id` bigint unsigned NOT NULL COMMENT '''数据库ID''',
  `db_name` varchar(100) NOT NULL COMMENT '''数据库名''',
  `export_type` varchar(20) NOT NULL COMMENT '''导出类型:structure/full''',
  `status` varchar(20) NOT NULL COMMENT '''状态''',
  `file_path` varchar(500) DEFAULT NULL COMMENT '''文件路径''',
  `file_size` bigint DEFAULT NULL COMMENT '''文件大小(字节)''',
  `error_message` text COMMENT '''错误信息''',
  `start_time` datetime(3) DEFAULT NULL COMMENT '''开始时间''',
  `end_time` datetime(3) DEFAULT NULL COMMENT '''结束时间''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `update_time` datetime(3) DEFAULT NULL COMMENT '''更新时间''',
  `creator` varchar(64) DEFAULT NULL COMMENT '''创建人''',
  `creator_id` bigint unsigned DEFAULT NULL COMMENT '''创建人ID''',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_db_export_task_task_id` (`task_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of db_export_task
-- ----------------------------
BEGIN;
INSERT INTO `db_export_task` (`id`, `task_id`, `db_id`, `db_name`, `export_type`, `status`, `file_path`, `file_size`, `error_message`, `start_time`, `end_time`, `create_time`, `update_time`, `creator`, `creator_id`) VALUES (1, '91c691d3-8e72-40d3-9de7-4453057ca1ee', 1, 'devops', 'full', 'completed', 'data/exports/devops_full_1765299147.sql', 689403, '', '2025-12-10 00:52:27.431', '2025-12-10 00:52:40.023', '2025-12-10 00:52:27.082', '2025-12-10 00:52:40.023', '', 0);
INSERT INTO `db_export_task` (`id`, `task_id`, `db_id`, `db_name`, `export_type`, `status`, `file_path`, `file_size`, `error_message`, `start_time`, `end_time`, `create_time`, `update_time`, `creator`, `creator_id`) VALUES (2, 'e1c6f2c1-7a85-43fb-8255-3d9a9fd229cb', 1, 'mayfly-go', 'full', 'completed', 'data/exports/mayfly-go_full_1765333697.sql', 100718, '', '2025-12-10 10:28:17.221', '2025-12-10 10:28:28.959', '2025-12-10 10:28:16.888', '2025-12-10 10:28:28.959', '', 0);
INSERT INTO `db_export_task` (`id`, `task_id`, `db_id`, `db_name`, `export_type`, `status`, `file_path`, `file_size`, `error_message`, `start_time`, `end_time`, `create_time`, `update_time`, `creator`, `creator_id`) VALUES (3, 'd788c68a-bab0-4d4c-8b6f-467d672876b1', 2, 'database_name', 'full', 'completed', 'data/exports/database_name_full_1765342425.sql', 1307, '', '2025-12-10 12:53:45.551', '2025-12-10 12:53:46.084', '2025-12-10 12:53:45.207', '2025-12-10 12:53:46.084', '', 0);
INSERT INTO `db_export_task` (`id`, `task_id`, `db_id`, `db_name`, `export_type`, `status`, `file_path`, `file_size`, `error_message`, `start_time`, `end_time`, `create_time`, `update_time`, `creator`, `creator_id`) VALUES (4, 'c6c4abdf-5dff-4a43-ba4c-64fac0dfc6de', 1, 'gin-api', 'full', 'completed', 'data/exports/gin-api_full_1765636630.sql', 652608, '', '2025-12-13 22:37:10.521', '2025-12-13 22:37:24.556', '2025-12-13 22:37:10.153', '2025-12-13 22:37:24.556', '', 0);
COMMIT;

-- ----------------------------
-- Table structure for db_instance
-- ----------------------------
DROP TABLE IF EXISTS `db_instance`;
CREATE TABLE `db_instance` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `code` varchar(36) NOT NULL COMMENT '''实例编码''',
  `name` varchar(100) NOT NULL COMMENT '''实例名称''',
  `type` varchar(20) NOT NULL COMMENT '''数据库类型:mysql,postgres,oracle等''',
  `host` varchar(100) NOT NULL COMMENT '''主机地址''',
  `port` bigint NOT NULL COMMENT '''端口''',
  `network` varchar(20) DEFAULT 'tcp' COMMENT '''网络类型:tcp,unix''',
  `params` varchar(500) DEFAULT NULL COMMENT '''连接参数''',
  `username` varchar(100) NOT NULL COMMENT '''用户名''',
  `password` varchar(500) NOT NULL COMMENT '''密码(加密)''',
  `remark` varchar(500) DEFAULT NULL COMMENT '''备注''',
  `ssh_tunnel_machine_id` bigint unsigned DEFAULT '0' COMMENT '''SSH隧道机器ID''',
  `status` bigint DEFAULT '1' COMMENT '''状态:1->启用,2->禁用''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `update_time` datetime(3) DEFAULT NULL COMMENT '''更新时间''',
  `creator` varchar(64) DEFAULT NULL COMMENT '''创建人''',
  `creator_id` bigint unsigned DEFAULT NULL COMMENT '''创建人ID''',
  `modifier` varchar(64) DEFAULT NULL COMMENT '''修改人''',
  `modifier_id` bigint unsigned DEFAULT NULL COMMENT '''修改人ID''',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_db_instance_code` (`code`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of db_instance
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for db_mongo_instance
-- ----------------------------
DROP TABLE IF EXISTS `db_mongo_instance`;
CREATE TABLE `db_mongo_instance` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `code` varchar(36) NOT NULL COMMENT '''实例编码''',
  `name` varchar(100) NOT NULL COMMENT '''实例名称''',
  `uri` varchar(500) NOT NULL COMMENT '''连接URI''',
  `ssh_tunnel_machine_id` bigint unsigned DEFAULT '0' COMMENT '''SSH隧道机器ID''',
  `remark` varchar(500) DEFAULT NULL COMMENT '''备注''',
  `status` bigint DEFAULT '1' COMMENT '''状态:1->启用,2->禁用''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `update_time` datetime(3) DEFAULT NULL COMMENT '''更新时间''',
  `creator` varchar(64) DEFAULT NULL COMMENT '''创建人''',
  `creator_id` bigint unsigned DEFAULT NULL COMMENT '''创建人ID''',
  `modifier` varchar(64) DEFAULT NULL COMMENT '''修改人''',
  `modifier_id` bigint unsigned DEFAULT NULL COMMENT '''修改人ID''',
  `type` varchar(20) DEFAULT 'mongodb' COMMENT '''MongoDB类型:mongodb,mongodb-atlas等''',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_db_mongo_instance_code` (`code`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of db_mongo_instance
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for db_redis_instance
-- ----------------------------
DROP TABLE IF EXISTS `db_redis_instance`;
CREATE TABLE `db_redis_instance` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `code` varchar(36) NOT NULL COMMENT '''实例编码''',
  `name` varchar(100) NOT NULL COMMENT '''实例名称''',
  `mode` varchar(20) NOT NULL COMMENT '''模式:standalone,cluster,sentinel''',
  `host` varchar(300) NOT NULL COMMENT '''主机: standalone为host:port, cluster为逗号分隔, sentinel为master=hosts''',
  `port` bigint DEFAULT '0' COMMENT '''端口(standalone可用)''',
  `db` bigint DEFAULT '0' COMMENT '''默认库号''',
  `username` varchar(100) DEFAULT NULL COMMENT '''用户名(可选)''',
  `password` varchar(500) NOT NULL COMMENT '''密码(加密)''',
  `redis_node_password` varchar(500) DEFAULT '' COMMENT '''节点密码(仅sentinel)''',
  `remark` varchar(500) DEFAULT NULL COMMENT '''备注''',
  `ssh_tunnel_machine_id` bigint unsigned DEFAULT '0' COMMENT '''SSH隧道机器ID''',
  `status` bigint DEFAULT '1' COMMENT '''状态:1->启用,2->禁用''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `update_time` datetime(3) DEFAULT NULL COMMENT '''更新时间''',
  `creator` varchar(64) DEFAULT NULL COMMENT '''创建人''',
  `creator_id` bigint unsigned DEFAULT NULL COMMENT '''创建人ID''',
  `modifier` varchar(64) DEFAULT NULL COMMENT '''修改人''',
  `modifier_id` bigint unsigned DEFAULT NULL COMMENT '''修改人ID''',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_db_redis_instance_code` (`code`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of db_redis_instance
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for db_sql
-- ----------------------------
DROP TABLE IF EXISTS `db_sql`;
CREATE TABLE `db_sql` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `db_id` bigint unsigned NOT NULL COMMENT '''数据库ID''',
  `db` varchar(100) NOT NULL COMMENT '''数据库名''',
  `name` varchar(100) NOT NULL COMMENT '''SQL名称''',
  `type` bigint DEFAULT '1' COMMENT '''类型:1->查询,2->更新''',
  `sql` text NOT NULL COMMENT '''SQL语句''',
  `remark` varchar(500) DEFAULT NULL COMMENT '''备注''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `update_time` datetime(3) DEFAULT NULL COMMENT '''更新时间''',
  `creator` varchar(64) DEFAULT NULL COMMENT '''创建人''',
  `creator_id` bigint unsigned DEFAULT NULL COMMENT '''创建人ID''',
  `modifier` varchar(64) DEFAULT NULL COMMENT '''修改人''',
  `modifier_id` bigint unsigned DEFAULT NULL COMMENT '''修改人ID''',
  PRIMARY KEY (`id`),
  KEY `idx_db_sql_db_id` (`db_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of db_sql
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for db_sql_exec
-- ----------------------------
DROP TABLE IF EXISTS `db_sql_exec`;
CREATE TABLE `db_sql_exec` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `db_id` bigint unsigned NOT NULL COMMENT '''数据库ID''',
  `db_name` varchar(100) NOT NULL COMMENT '''数据库名''',
  `table_name` varchar(100) DEFAULT NULL COMMENT '''表名''',
  `type` tinyint NOT NULL COMMENT '''类型:1->查询,2->插入,3->更新,4->删除,5->DDL''',
  `sql` text NOT NULL COMMENT '''SQL语句''',
  `old_value` longtext COMMENT '''旧值(用于回滚)''',
  `remark` varchar(500) DEFAULT NULL COMMENT '''备注''',
  `status` tinyint NOT NULL COMMENT '''状态:1->成功,2->失败''',
  `res` text COMMENT '''执行结果''',
  `exec_time` bigint DEFAULT NULL COMMENT '''执行时长(ms)''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `creator` varchar(64) DEFAULT NULL COMMENT '''创建人''',
  `creator_id` bigint unsigned DEFAULT NULL COMMENT '''创建人ID''',
  PRIMARY KEY (`id`),
  KEY `idx_db_sql_exec_db_id` (`db_id`),
  KEY `idx_db_sql_exec_create_time` (`create_time`)
) ENGINE=InnoDB AUTO_INCREMENT=456 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of db_sql_exec
-- ----------------------------
BEGIN;
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (1, 1, 'gin-api', 'APP_APPLICATION', 0, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 2, '数据库不存在', 110, '2025-11-29 18:05:04.883', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (2, 1, 'gin-api', 'APP_APPLICATION', 0, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 2, '数据库不存在', 112, '2025-11-29 18:05:05.487', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (3, 1, 'gin-api', 'CMDB_GROUP', 0, 'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0', '', '', 2, '数据库不存在', 119, '2025-11-29 18:05:10.402', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (4, 1, 'gin-api', 'CMDB_GROUP', 0, 'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0', '', '', 2, '数据库不存在', 116, '2025-11-29 18:05:10.598', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (5, 1, 'gin-api', 'CMDB_SQL', 0, 'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0', '', '', 2, '数据库不存在', 122, '2025-11-29 18:05:17.424', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (6, 1, 'gin-api', 'CMDB_SQL', 0, 'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0', '', '', 2, '数据库不存在', 136, '2025-11-29 18:05:17.580', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (7, 1, 'gin-api', 'APP_JENKINS_ENV', 0, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 2, '数据库不存在', 119, '2025-11-29 18:11:06.550', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (8, 1, 'gin-api', 'APP_JENKINS_ENV', 0, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 2, '数据库不存在', 493, '2025-11-29 18:11:06.802', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (9, 1, 'gin-api', 'APP_JENKINS_ENV', 0, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 2, '数据库不存在', 127, '2025-11-29 18:11:15.554', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (10, 1, 'gin-api', 'APP_JENKINS_ENV', 0, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 2, '数据库不存在', 110, '2025-11-29 18:11:15.687', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (11, 1, 'gin-api', 'APP_SERVICE_RELEASE_ITEM', 1, 'SELECT * FROM app_service_release_item LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 548, '2025-11-29 18:28:12.401', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (12, 1, 'gin-api', 'APP_SERVICE_RELEASE_ITEM', 1, 'SELECT * FROM app_service_release_item LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 1083, '2025-11-29 18:28:12.656', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (13, 1, 'gin-api', 'APP_APPLICATION', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 277, '2025-11-29 18:28:25.668', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (14, 1, 'gin-api', 'APP_APPLICATION', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 281, '2025-11-29 18:28:25.951', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (15, 1, 'devops', 'APP_APPLICATION', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 308, '2025-11-29 18:31:23.236', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (16, 1, 'devops', 'APP_APPLICATION', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 346, '2025-11-29 18:31:23.606', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (17, 1, 'gin-api', 'APP_APPLICATION', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 298, '2025-11-29 18:35:52.379', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (18, 1, 'gin-api', 'APP_APPLICATION', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 303, '2025-11-29 18:35:52.538', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (19, 1, 'gin-api', 'APP_APPLICATION', 1, 'SELECT COUNT(*) as total FROM app_application', '', '', 1, '执行成功', 266, '2025-11-29 18:35:53.254', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (20, 1, 'gin-api', 'APP_APPLICATION', 1, 'SELECT COUNT(*) as total FROM app_application', '', '', 1, '执行成功', 259, '2025-11-29 18:35:53.562', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (21, 1, 'devops', 'APP_APPLICATION', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 283, '2025-11-29 18:37:11.404', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (22, 1, 'devops', 'APP_APPLICATION', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 320, '2025-11-29 18:37:12.000', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (23, 1, 'devops', 'APP_APPLICATION', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 273, '2025-11-29 18:37:11.713', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (24, 1, 'devops', 'APP_APPLICATION', 1, 'SELECT COUNT(*) as total FROM app_application', '', '', 1, '执行成功', 299, '2025-11-29 18:37:12.382', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (25, 1, 'devops', 'APP_APPLICATION', 1, 'SELECT COUNT(*) as total FROM app_application', '', '', 1, '执行成功', 255, '2025-11-29 18:37:12.911', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (26, 1, 'devops', 'APP_APPLICATION', 1, 'SELECT COUNT(*) as total FROM app_application', '', '', 1, '执行成功', 282, '2025-11-29 18:37:13.296', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (27, 1, 'devops', 'APP_APPLICATION', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 248, '2025-11-29 18:37:58.462', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (28, 1, 'devops', 'APP_APPLICATION', 1, 'SELECT COUNT(*) as total FROM app_application', '', '', 1, '执行成功', 260, '2025-11-29 18:37:59.326', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (29, 1, 'devops', 'APP_APPLICATION', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 275, '2025-11-29 18:38:46.987', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (30, 1, 'devops', 'APP_APPLICATION', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 277, '2025-11-29 18:38:47.151', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (31, 1, 'devops', 'APP_APPLICATION', 1, 'SELECT COUNT(*) as total FROM app_application', '', '', 1, '执行成功', 265, '2025-11-29 18:38:47.925', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (32, 1, 'devops', 'APP_APPLICATION', 1, 'SELECT COUNT(*) as total FROM app_application', '', '', 1, '执行成功', 266, '2025-11-29 18:38:48.246', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (33, 1, 'gin-api', 'APP_APPLICATION', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 265, '2025-11-29 18:39:00.606', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (34, 1, 'gin-api', 'APP_APPLICATION', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 304, '2025-11-29 18:39:00.775', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (35, 1, 'gin-api', 'APP_APPLICATION', 1, 'SELECT COUNT(*) as total FROM app_application', '', '', 1, '执行成功', 287, '2025-11-29 18:39:01.496', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (36, 1, 'gin-api', 'APP_APPLICATION', 1, 'SELECT COUNT(*) as total FROM app_application', '', '', 1, '执行成功', 263, '2025-11-29 18:39:01.812', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (37, 1, 'devops', 'APP_APPLICATION', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 403, '2025-11-29 18:48:01.028', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (38, 1, 'devops', 'APP_APPLICATION', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 335, '2025-11-29 18:48:01.099', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (39, 1, 'devops', 'APP_APPLICATION', 1, 'SELECT COUNT(*) as total FROM app_application', '', '', 1, '执行成功', 10295, '2025-11-29 18:48:15.394', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (40, 1, 'devops', 'APP_APPLICATION', 1, 'SELECT COUNT(*) as total FROM app_application', '', '', 1, '执行成功', 2650, '2025-11-29 18:48:14.638', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (41, 1, 'devops', 'CMDB_HOST', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 6494, '2025-11-29 18:48:15.978', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (42, 1, 'devops', 'CMDB_HOST', 1, 'SELECT COUNT(*) as total FROM cmdb_host', '', '', 1, '执行成功', 293, '2025-11-29 18:48:16.898', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (43, 1, 'devops', 'CMDB_HOST', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 3692, '2025-11-29 18:48:12.891', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (44, 1, 'devops', 'CMDB_HOST', 1, 'SELECT COUNT(*) as total FROM cmdb_host', '', '', 1, '执行成功', 281, '2025-11-29 18:48:21.045', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (45, 1, 'devops', 'APP_APPLICATION', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 340, '2025-11-29 18:50:12.488', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (46, 1, 'devops', 'APP_APPLICATION', 1, 'SELECT COUNT(*) as total FROM app_application', '', '', 1, '执行成功', 284, '2025-11-29 18:50:13.429', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (47, 1, 'devops', 'APP_JENKINS_ENV', 1, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 304, '2025-11-29 18:54:02.563', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (48, 1, 'devops', 'APP_JENKINS_ENV', 1, 'SELECT COUNT(*) as total FROM app_jenkins_env', '', '', 1, '执行成功', 294, '2025-11-29 18:54:03.463', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (49, 1, 'devops', 'CMDB_GROUP', 1, 'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 328, '2025-11-29 18:54:07.320', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (50, 1, 'devops', 'CMDB_GROUP', 1, 'SELECT COUNT(*) as total FROM cmdb_group', '', '', 1, '执行成功', 313, '2025-11-29 18:54:08.231', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (51, 1, 'devops', 'APP_APPLICATION', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 282, '2025-11-29 18:55:28.076', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (52, 1, 'devops', 'APP_APPLICATION', 1, 'SELECT COUNT(*) as total FROM app_application', '', '', 1, '执行成功', 280, '2025-11-29 18:55:28.949', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (53, 1, 'devops', 'APP_JENKINS_ENV', 1, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 267, '2025-11-29 18:55:29.997', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (54, 1, 'devops', 'APP_JENKINS_ENV', 1, 'SELECT COUNT(*) as total FROM app_jenkins_env', '', '', 1, '执行成功', 277, '2025-11-29 18:55:30.867', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (55, 1, 'devops', 'APP_JENKINS_ENV', 1, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 282, '2025-11-29 18:56:33.893', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (56, 1, 'devops', 'APP_JENKINS_ENV', 1, 'SELECT COUNT(*) as total FROM app_jenkins_env', '', '', 1, '执行成功', 275, '2025-11-29 18:56:34.766', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (57, 1, 'devops', 'CMDB_GROUP', 1, 'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 309, '2025-11-29 18:56:37.254', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (58, 1, 'devops', 'CMDB_GROUP', 1, 'SELECT COUNT(*) as total FROM cmdb_group', '', '', 1, '执行成功', 293, '2025-11-29 18:56:38.144', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (59, 1, 'devops', 'CMDB_HOST', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 276, '2025-11-29 18:57:16.850', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (60, 1, 'devops', 'CMDB_HOST', 1, 'SELECT COUNT(*) as total FROM cmdb_host', '', '', 1, '执行成功', 284, '2025-11-29 18:57:17.734', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (61, 1, 'devops', 'CMDB_SQL_LOG', 1, 'SELECT * FROM cmdb_sql_log LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 402, '2025-11-29 18:57:49.376', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (62, 1, 'devops', 'CMDB_SQL_LOG', 1, 'SELECT COUNT(*) as total FROM cmdb_sql_log', '', '', 1, '执行成功', 274, '2025-11-29 18:57:50.233', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (63, 1, 'devops', 'APP_APPLICATION', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 315, '2025-11-29 19:00:49.626', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (64, 1, 'devops', 'APP_APPLICATION', 1, 'SELECT COUNT(*) as total FROM app_application', '', '', 1, '执行成功', 330, '2025-11-29 19:00:50.546', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (65, 1, 'devops', 'CMDB_HOST', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 484, '2025-11-29 19:00:57.485', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (66, 1, 'devops', 'CMDB_HOST', 1, 'SELECT COUNT(*) as total FROM cmdb_host', '', '', 1, '执行成功', 656, '2025-11-29 19:00:59.337', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (67, 1, 'devops', 'CMDB_SQL', 1, 'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 432, '2025-11-29 19:01:33.694', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (68, 1, 'devops', 'CMDB_SQL', 1, 'SELECT COUNT(*) as total FROM cmdb_sql', '', '', 1, '执行成功', 285, '2025-11-29 19:01:34.585', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (69, 1, 'devops', 'CONFIG_ACCOUNT', 1, 'SELECT * FROM config_account LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 300, '2025-11-29 19:01:52.297', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (70, 1, 'devops', 'CONFIG_ACCOUNT', 1, 'SELECT COUNT(*) as total FROM config_account', '', '', 1, '执行成功', 282, '2025-11-29 19:01:53.222', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (71, 1, 'devops', 'CMDB_SQL_LOG', 1, 'SELECT * FROM cmdb_sql_log LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 442, '2025-11-29 19:02:12.743', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (72, 1, 'devops', 'CMDB_SQL_LOG', 1, 'SELECT COUNT(*) as total FROM cmdb_sql_log', '', '', 1, '执行成功', 388, '2025-11-29 19:02:13.732', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (73, 1, 'devops', 'APP_APPLICATION', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 321, '2025-11-29 19:04:24.978', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (74, 1, 'devops', 'APP_APPLICATION', 1, 'SELECT COUNT(*) as total FROM app_application', '', '', 1, '执行成功', 291, '2025-11-29 19:04:25.863', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (75, 1, 'devops', 'CMDB_HOST', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 286, '2025-11-29 19:05:02.134', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (76, 1, 'devops', 'CMDB_HOST', 1, 'SELECT COUNT(*) as total FROM cmdb_host', '', '', 1, '执行成功', 309, '2025-11-29 19:05:03.008', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (77, 1, 'gin-api', 'APP_APPLICATION', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 361, '2025-11-29 19:06:09.051', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (78, 1, 'gin-api', 'APP_APPLICATION', 1, 'SELECT COUNT(*) as total FROM app_application', '', '', 1, '执行成功', 291, '2025-11-29 19:06:10.002', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (79, 1, 'gin-api', 'APP_APPLICATION', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 1109, '2025-11-29 19:06:32.787', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (80, 1, 'gin-api', 'APP_APPLICATION', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 1586, '2025-11-29 19:06:33.149', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (81, 1, 'gin-api', 'APP_APPLICATION', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 303, '2025-11-29 19:06:33.300', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (82, 1, 'gin-api', 'APP_APPLICATION', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 552, '2025-11-29 19:06:33.428', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (83, 1, 'gin-api', 'APP_APPLICATION', 1, 'SELECT COUNT(*) as total FROM app_application', '', '', 1, '执行成功', 307, '2025-11-29 19:06:33.766', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (84, 1, 'gin-api', 'APP_APPLICATION', 1, 'SELECT COUNT(*) as total FROM app_application', '', '', 1, '执行成功', 254, '2025-11-29 19:06:34.069', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (85, 1, 'gin-api', 'APP_APPLICATION', 1, 'SELECT COUNT(*) as total FROM app_application', '', '', 1, '执行成功', 553, '2025-11-29 19:06:34.715', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (86, 1, 'gin-api', 'APP_APPLICATION', 1, 'SELECT COUNT(*) as total FROM app_application', '', '', 1, '执行成功', 302, '2025-11-29 19:06:35.225', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (87, 1, 'gin-api', 'APP_JENKINS_ENV', 1, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 295, '2025-11-29 19:06:36.173', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (88, 1, 'gin-api', 'APP_JENKINS_ENV', 1, 'SELECT COUNT(*) as total FROM app_jenkins_env', '', '', 1, '执行成功', 276, '2025-11-29 19:06:37.028', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (89, 1, 'devops', 'CMDB_SQL', 1, 'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 275, '2025-11-29 19:08:02.114', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (90, 1, 'devops', 'CMDB_SQL', 1, 'SELECT COUNT(*) as total FROM cmdb_sql', '', '', 1, '执行成功', 331, '2025-11-29 19:08:03.077', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (91, 1, 'devops', 'CONFIG_ACCOUNT', 1, 'SELECT * FROM config_account LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 324, '2025-11-29 19:08:59.699', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (92, 1, 'devops', 'CMDB_SQL_RECORDS', 1, 'SELECT * FROM cmdb_sql_records LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 299, '2025-11-29 19:08:59.964', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (93, 1, 'devops', 'CMDB_SQL_RECORDS', 1, 'SELECT * FROM cmdb_sql_records LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 275, '2025-11-29 19:09:00.054', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (94, 1, 'devops', 'CMDB_SQL_RECORDS', 1, 'SELECT * FROM cmdb_sql_records LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 259, '2025-11-29 19:09:00.169', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (95, 1, 'devops', 'CONFIG_ACCOUNT', 1, 'SELECT * FROM config_account LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 360, '2025-11-29 19:08:59.813', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (96, 1, 'devops', 'CONFIG_ACCOUNT', 1, 'SELECT COUNT(*) as total FROM config_account', '', '', 1, '执行成功', 292, '2025-11-29 19:09:00.602', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (97, 1, 'devops', 'CMDB_SQL_RECORDS', 1, 'SELECT COUNT(*) as total FROM cmdb_sql_records', '', '', 1, '执行成功', 303, '2025-11-29 19:09:00.855', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (98, 1, 'devops', 'CONFIG_ACCOUNT', 1, 'SELECT COUNT(*) as total FROM config_account', '', '', 1, '执行成功', 326, '2025-11-29 19:09:01.089', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (99, 1, 'devops', 'CMDB_SQL_RECORDS', 1, 'SELECT COUNT(*) as total FROM cmdb_sql_records', '', '', 1, '执行成功', 295, '2025-11-29 19:09:01.217', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (100, 1, 'devops', 'CMDB_SQL_RECORDS', 1, 'SELECT COUNT(*) as total FROM cmdb_sql_records', '', '', 1, '执行成功', 334, '2025-11-29 19:09:01.611', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (101, 1, 'devops', 'CMDB_SQL', 1, 'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 280, '2025-11-29 19:09:14.451', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (102, 1, 'devops', 'CMDB_SQL', 1, 'SELECT COUNT(*) as total FROM cmdb_sql', '', '', 1, '执行成功', 278, '2025-11-29 19:09:15.298', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (103, 1, 'devops', 'CMDB_HOST', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 1029, '2025-11-29 19:09:29.074', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (104, 1, 'devops', 'CMDB_HOST', 1, 'SELECT COUNT(*) as total FROM cmdb_host', '', '', 1, '执行成功', 6538, '2025-11-29 19:09:36.486', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (105, 1, 'devops', 'APP_APPLICATION', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 264, '2025-11-29 19:23:51.118', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (106, 1, 'devops', 'APP_APPLICATION', 1, 'SELECT COUNT(*) as total FROM app_application', '', '', 1, '执行成功', 296, '2025-11-29 19:23:51.990', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (107, 1, 'devops', 'CMDB_SQL', 1, 'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 295, '2025-11-29 19:23:59.137', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (108, 1, 'devops', 'CMDB_SQL', 1, 'SELECT COUNT(*) as total FROM cmdb_sql', '', '', 1, '执行成功', 279, '2025-11-29 19:24:00.023', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (109, 1, 'devops', 'APP_APPLICATION', 1, 'SELECT * FROM app_application LIMIT 100;', '', '', 1, '执行成功', 280, '2025-11-29 19:25:10.430', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (110, 1, 'devops', 'APP_APPLICATION', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 284, '2025-11-29 19:49:42.427', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (111, 1, 'devops', 'APP_APPLICATION', 1, 'SELECT COUNT(*) as total FROM app_application', '', '', 1, '执行成功', 287, '2025-11-29 19:49:43.313', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (112, 1, 'devops', 'APP_APPLICATION', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 239, '2025-11-29 19:57:17.599', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (113, 1, 'devops', 'CMDB_SQL_LOG', 1, 'SELECT * FROM cmdb_sql_log LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 213, '2025-11-29 19:57:40.935', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (114, 1, 'devops', 'CMDB_SQL', 1, 'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 213, '2025-11-29 19:57:40.616', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (115, 1, 'devops', 'CMDB_SQL_LOG', 1, 'SELECT * FROM cmdb_sql_log LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 212, '2025-11-29 19:57:40.797', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (116, 1, 'devops', 'CMDB_SQL_LOG', 1, 'SELECT * FROM cmdb_sql_log LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 210, '2025-11-29 19:57:59.070', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (117, 1, 'devops', 'CMDB_SQL_RECORDS', 1, 'SELECT * FROM cmdb_sql_records LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 248, '2025-11-29 19:58:15.888', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (118, 1, 'devops', 'CMDB_HOST', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 307, '2025-11-29 19:58:47.223', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (119, 1, 'devops', 'CMDB_HOST;', 1, 'SELECT * FROM cmdb_host;', '', '', 1, '执行成功', 212, '2025-11-29 20:00:09.885', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (120, 1, 'devops', 'APP_APPLICATION', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 263, '2025-11-29 20:03:26.517', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (121, 1, 'devops', 'CMDB_HOST', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 213, '2025-11-29 20:03:44.248', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (122, 1, 'devops', 'CMDB_HOST', 1, 'SELECT * FROM cmdb_host LIMIT 100;', '', '', 1, '执行成功', 220, '2025-11-29 20:04:04.538', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (123, 1, 'devops', 'TASK_WORK', 1, 'SELECT * FROM task_work LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 234, '2025-11-29 20:04:29.401', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (124, 1, 'devops', 'K8S_CLUSTER', 1, 'SELECT * FROM k8s_cluster LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 247, '2025-11-29 20:05:20.336', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (125, 1, 'devops', 'CONFIG_ECSAUTH', 1, 'SELECT * FROM config_ecsauth LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 515, '2025-11-29 20:07:54.567', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (126, 1, 'devops', 'CONFIG_ACCOUNT', 1, 'SELECT * FROM config_account LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 219, '2025-11-29 20:08:12.059', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (127, 1, 'devops', 'CONFIG_ECSAUTH', 1, 'SELECT * FROM config_ecsauth LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 386, '2025-11-29 20:08:56.093', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (128, 1, 'devops', 'CONFIG_ACCOUNT', 1, 'SELECT * FROM config_account LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 212, '2025-11-29 20:09:15.428', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (129, 1, 'devops', 'CONFIG_ACCOUNT', 1, 'SELECT * FROM config_account LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 241, '2025-11-29 20:09:18.899', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (130, 1, 'devops', 'SYS_MENU', 1, 'SELECT * FROM sys_menu LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 211, '2025-11-29 20:09:31.833', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (131, 1, 'devops', 'CONFIG_KEYMANAGE', 1, 'SELECT * FROM config_keymanage LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 218, '2025-11-29 20:10:16.736', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (132, 1, 'gin-api', 'APP_SERVICE_RELEASE', 1, 'SELECT * FROM app_service_release LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 229, '2025-11-29 20:10:41.975', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (133, 1, 'gin-api', 'APP_JENKINS_ENV', 1, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 321, '2025-11-29 20:11:33.117', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (134, 1, 'devops', 'APP_JENKINS_ENV', 1, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 213, '2025-11-29 20:14:46.781', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (135, 1, 'devops', 'APP_JENKINS_ENV', 1, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 213, '2025-11-29 20:17:07.095', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (136, 1, 'devops', 'CMDB_HOST', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 211, '2025-11-29 20:17:20.352', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (137, 1, 'devops', 'CMDB_GROUP', 1, 'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 240, '2025-11-29 20:27:29.297', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (138, 1, 'devops', 'CMDB_SQL', 1, 'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 258, '2025-11-29 20:28:55.241', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (139, 1, 'devops', 'CMDB_SQL', 1, 'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 579, '2025-11-29 20:28:55.466', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (140, 1, 'devops', 'CMDB_GROUP', 1, 'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 205, '2025-11-29 20:30:47.480', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (141, 1, 'devops', 'CMDB_HOST', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 199, '2025-11-29 20:30:58.495', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (142, 1, 'devops', 'CMDB_GROUP', 1, 'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 267, '2025-11-29 20:34:21.225', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (143, 1, 'devops', 'CMDB_GROUP', 1, 'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 4246, '2025-11-29 20:45:39.121', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (144, 1, 'devops', 'CMDB_GROUP', 1, 'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 582, '2025-11-29 20:49:42.415', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (145, 1, 'devops', 'APP_APPLICATION', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 212, '2025-11-29 20:59:18.272', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (146, 1, 'devops', 'APP_APPLICATION', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 265, '2025-11-29 21:02:11.486', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (147, 1, 'devops', 'APP_JENKINS_ENV', 1, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 242, '2025-11-29 21:05:56.924', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (148, 1, 'devops', 'APP_JENKINS_ENV', 1, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 269, '2025-11-29 21:06:38.995', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (149, 1, 'devops', 'APP_JENKINS_ENV', 1, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 327, '2025-11-29 21:08:52.177', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (150, 1, 'devops', 'APP_JENKINS_ENV', 1, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 646, '2025-11-29 21:10:07.252', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (151, 1, 'devops', 'app_jenkins_env', 1, 'SELECT * FROM app_jenkins_env LIMIT 100;', '', '', 1, '执行成功', 293, '2025-11-29 21:18:18.161', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (152, 1, 'devops', 'app_application', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 320, '2025-11-29 21:21:01.619', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (153, 1, 'devops', 'app_application', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 297, '2025-11-29 21:21:10.168', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (154, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 326, '2025-11-29 21:21:27.886', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (155, 1, 'devops', 'cmdb_sql', 1, 'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 283, '2025-11-29 21:22:09.300', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (156, 1, 'devops', 'cmdb_group', 1, 'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 283, '2025-11-29 21:22:19.066', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (157, 1, 'devops', 'cmdb_group', 1, 'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 276, '2025-11-29 21:22:21.160', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (158, 1, 'devops', 'cmdb_group', 1, 'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 264, '2025-11-29 21:34:45.094', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (159, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 1520, '2025-11-29 21:35:20.063', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (160, 1, 'sys', 'sys_config', 1, 'SELECT * FROM sys_config LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 293, '2025-11-29 21:36:51.862', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (161, 1, 'gin-api', 'app_jenkins_env', 1, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 265, '2025-11-29 21:37:07.653', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (162, 1, 'gin-api', 'app_application', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 292, '2025-11-29 21:37:26.201', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (163, 1, 'devops', 'cmdb_sql_log', 1, 'SELECT * FROM cmdb_sql_log LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 298, '2025-11-29 21:37:48.390', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (164, 1, 'devops', 'cmdb_sql_log', 1, 'SELECT * FROM cmdb_sql_log LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 277, '2025-11-29 21:38:13.100', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (165, 1, 'devops', 'app_jenkins_env', 1, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 282, '2025-11-29 21:40:52.663', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (166, 1, 'devops', 'app_application', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 283, '2025-11-29 21:42:31.577', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (167, 1, 'devops', 'app_application', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 283, '2025-11-29 21:45:10.707', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (168, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 285, '2025-11-29 21:45:23.126', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (169, 1, 'devops', 'cmdb_group', 1, 'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 270, '2025-11-29 21:45:42.521', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (170, 1, 'devops', 'cmdb_group', 1, 'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 491, '2025-11-29 21:46:03.897', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (171, 1, 'devops', 'app_application', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 295, '2025-11-29 21:50:35.944', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (172, 1, 'devops', 'cmdb_sql', 1, 'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 270, '2025-11-29 21:50:58.532', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (173, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 258, '2025-11-29 21:51:00.879', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (174, 1, 'devops', 'cmdb_group', 1, 'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 266, '2025-11-29 21:54:42.129', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (175, 1, 'devops', 'app_jenkins_env', 1, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 299, '2025-11-29 21:57:35.040', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (176, 1, 'devops', 'app_jenkins_env', 1, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 408, '2025-11-29 21:57:38.152', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (177, 1, 'devops', 'cmdb_group', 1, 'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 279, '2025-11-29 21:57:43.563', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (178, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 341, '2025-11-29 21:57:53.234', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (179, 1, 'devops', 'config_ecsauth', 1, 'SELECT * FROM config_ecsauth LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 309, '2025-11-29 21:58:20.086', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (180, 1, 'devops', 'config_ecsauth', 1, 'SELECT * FROM config_ecsauth LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 579, '2025-11-29 21:58:20.467', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (181, 1, 'devops', 'config_ecsauth', 1, 'SELECT * FROM config_ecsauth LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 299, '2025-11-29 21:58:28.291', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (182, 1, 'devops', 'monitor_agent', 1, 'SELECT * FROM monitor_agent LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 277, '2025-11-29 21:58:36.381', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (183, 1, 'devops', 'cmdb_sql', 1, 'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 330, '2025-11-29 21:59:23.190', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (184, 1, 'devops', 'cmdb_sql', 1, 'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 450, '2025-11-29 21:59:44.784', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (185, 1, 'devops', 'app_application', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 295, '2025-11-29 22:10:51.726', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (186, 1, 'devops', 'cmdb_group', 1, 'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 266, '2025-11-29 22:10:54.896', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (187, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 323, '2025-11-29 22:11:01.210', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (188, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 288, '2025-11-29 22:11:01.374', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (189, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 319, '2025-11-29 22:14:15.745', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (190, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 508, '2025-11-29 22:14:16.040', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (191, 1, 'devops', 'cmdb_group', 1, 'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 261, '2025-11-29 22:14:22.801', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (192, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 278, '2025-11-29 22:14:29.071', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (193, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 587, '2025-11-29 22:17:49.851', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (194, 1, 'devops', 'cmdb_group', 1, 'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 273, '2025-11-29 22:21:41.845', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (195, 1, 'devops', 'cmdb_group', 1, 'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 297, '2025-11-29 22:21:41.996', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (196, 1, 'devops', 'cmdb_group', 1, 'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 338, '2025-11-29 22:21:41.638', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (197, 1, 'devops', 'cmdb_group', 1, 'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 935, '2025-11-29 22:21:42.359', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (198, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 280, '2025-11-29 22:21:45.552', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (199, 1, 'devops', 'app_jenkins_env', 1, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 290, '2025-11-29 22:25:38.566', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (200, 1, 'devops', 'app_application', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 262, '2025-11-29 22:26:39.755', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (201, 1, 'devops', 'app_application', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 299, '2025-11-29 22:29:18.342', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (202, 1, 'devops', 'cmdb_sql', 1, 'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 288, '2025-11-29 22:30:54.805', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (203, 1, 'devops', 'app_application', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 282, '2025-11-29 22:31:17.619', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (204, 1, 'gin-api', 'app_jenkins_env', 1, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 289, '2025-11-29 22:32:34.668', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (205, 1, 'devops', 'cmdb_sql_log', 1, 'SELECT * FROM cmdb_sql_log LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 272, '2025-11-29 22:32:43.903', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (206, 1, 'devops', 'cmdb_group', 1, 'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 312, '2025-11-29 22:32:59.875', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (207, 1, 'devops', 'cmdb_group', 1, 'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 268, '2025-11-29 22:38:19.169', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (208, 1, 'devops', 'cmdb_group', 1, 'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 285, '2025-11-29 22:42:46.982', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (209, 1, 'devops', 'cmdb_group', 1, 'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 283, '2025-11-29 22:48:01.180', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (210, 1, 'devops', 'cmdb_group', 1, 'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 279, '2025-11-29 22:53:17.671', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (211, 1, 'devops', 'app_jenkins_env', 1, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 271, '2025-11-29 22:57:40.841', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (212, 1, 'devops', 'app_jenkins_env', 1, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 262, '2025-11-29 22:59:48.099', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (213, 1, 'devops', 'app_jenkins_env', 1, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 278, '2025-11-29 23:00:34.141', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (214, 1, 'devops', 'app_jenkins_env', 1, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 569, '2025-11-29 23:04:04.638', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (215, 1, 'devops', 'app_jenkins_env', 1, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 542, '2025-11-29 23:10:59.541', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (216, 1, 'gin-api', 'app_jenkins_env', 1, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 309, '2025-11-29 23:13:35.399', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (217, 1, 'gin-api', 'app_sh_release', 1, 'SELECT * FROM app_sh_release LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 316, '2025-11-29 23:17:36.816', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (218, 1, 'gin-api', 'app_sh_release', 1, 'SELECT * FROM app_sh_release LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 713, '2025-11-29 23:17:37.046', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (219, 1, 'gin-api', 'cmdb_group', 1, 'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 311, '2025-11-29 23:17:41.461', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (220, 1, 'gin-api', 'config_ecsauth', 1, 'SELECT * FROM config_ecsauth LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 451, '2025-11-29 23:18:01.612', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (221, 1, 'gin-api', 'config_keymanage', 1, 'SELECT * FROM config_keymanage LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 300, '2025-11-29 23:18:03.029', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (222, 1, 'gin-api', 'config_sync_schedule', 1, 'SELECT * FROM config_sync_schedule LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 305, '2025-11-29 23:18:04.708', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (223, 1, 'gin-api', 'sys_operation_log', 1, 'SELECT * FROM sys_operation_log LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 303, '2025-11-29 23:18:07.975', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (224, 1, 'gin-api', 'sys_operation_log', 1, 'SELECT * FROM sys_operation_log LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 334, '2025-11-29 23:18:13.728', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (225, 1, 'gin-api', 'app_service_release', 1, 'SELECT * FROM app_service_release LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 304, '2025-11-29 23:26:31.908', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (226, 1, 'devops', 'app_application', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 420, '2025-11-29 23:30:01.262', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (227, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 331, '2025-11-29 23:35:23.442', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (228, 1, 'devops', 'cmdb_sql_log', 1, 'SELECT * FROM cmdb_sql_log LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 306, '2025-11-29 23:36:11.800', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (229, 1, 'devops', 'cmdb_sql_log', 1, 'SELECT * FROM cmdb_sql_log LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 293, '2025-11-29 23:36:35.095', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (230, 1, 'gin-api', 'app_jenkins_env', 1, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 297, '2025-11-29 23:39:22.251', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (231, 1, 'gin-api', 'cmdb_group', 1, 'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 340, '2025-11-29 23:39:24.638', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (232, 1, 'devops', 'app_jenkins_env', 1, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 289, '2025-11-29 23:39:32.922', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (233, 1, 'devops', '', 5, 'RENAME TABLE app_jenkins_env_copy_20251129233940 TO app_jenkins_env_copy_20251129233940_new111', '', '', 1, '执行成功', 298, '2025-11-29 23:39:51.148', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (234, 1, 'devops', 'app_jenkins_env_copy_20251129233940_new111', 1, 'SELECT * FROM app_jenkins_env_copy_20251129233940_new111 LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 563, '2025-11-29 23:40:20.699', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (235, 1, 'devops', '', 5, 'DROP TABLE app_jenkins_env_copy_20251129233940_new111', '', '', 1, '执行成功', 291, '2025-11-29 23:40:30.953', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (236, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 434, '2025-11-29 23:44:31.265', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (237, 1, 'devops', 'cmdb_sql', 1, 'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 466, '2025-11-29 23:45:21.575', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (238, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 288, '2025-11-29 23:52:27.278', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (239, 1, 'devops', 'app_jenkins_env', 1, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 302, '2025-11-29 23:55:28.566', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (240, 1, 'gin-api', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 319, '2025-11-29 23:58:06.595', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (241, 2, 'postgres', '', 5, 'CREATE DATABASE database_name;', '', '', 1, '执行成功', 599, '2025-11-30 00:56:01.163', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (242, 2, 'database_name', '', 5, 'CREATE TABLE users (\n    id SERIAL PRIMARY KEY,\n    name VARCHAR(100) NOT NULL,\n    email VARCHAR(150) UNIQUE NOT NULL\n);', '', '', 1, '执行成功', 3259, '2025-11-30 01:10:07.720', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (243, 2, 'database_name', 'users', 1, 'SELECT * FROM users LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 4965, '2025-11-30 01:10:34.369', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (244, 2, 'database_name', 'users', 1, 'SELECT * FROM users LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 198, '2025-11-30 01:10:34.867', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (245, 2, 'database_name', 'users', 1, 'SELECT * FROM users LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 201, '2025-11-30 01:10:35.037', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (246, 2, 'database_name', 'users', 1, 'SELECT * FROM users LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 208, '2025-11-30 01:10:35.181', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (247, 2, 'database_name', 'users', 1, 'SELECT * FROM users LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 208, '2025-11-30 01:10:35.442', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (248, 2, 'database_name', 'users', 1, 'SELECT * FROM users LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 192, '2025-11-30 01:10:35.755', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (249, 2, 'database_name', 'users', 1, 'SELECT * FROM users LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 205, '2025-11-30 01:10:35.988', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (250, 2, 'database_name', 'users', 1, 'SELECT * FROM users LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 202, '2025-11-30 01:10:39.921', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (251, 2, 'database_name', 'users', 2, 'INSERT INTO users (name, email) VALUES (\'张三\', \'zhangsan@123.com\')', '', '', 1, '执行成功', 376, '2025-11-30 01:11:20.911', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (252, 2, 'database_name', 'users', 1, 'SELECT * FROM users LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 299, '2025-11-30 01:11:21.506', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (253, 1, 'gin-api', 'app_application', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 299, '2025-11-30 11:38:28.572', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (254, 1, 'gin-api', 'app_jenkins_env', 1, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 275, '2025-11-30 11:38:44.346', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (255, 1, 'gin-api', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 272, '2025-11-30 11:38:53.975', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (256, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 342, '2025-11-30 11:39:11.611', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (257, 1, 'devops', 'cmdb_sql', 1, 'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 277, '2025-11-30 11:39:53.115', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (258, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 284, '2025-11-30 11:40:04.010', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (259, 1, 'devops', 'app_application', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 273, '2025-11-30 14:35:46.125', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (260, 2, 'database_name', 'users', 1, 'SELECT * FROM users LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 219, '2025-11-30 15:43:25.630', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (261, 1, 'RECOVER_YOUR_DATA', '', 5, 'CREATE DATABASE IF NOT EXISTS test \n  DEFAULT CHARACTER SET utf8mb4;', '', '', 1, '执行成功', 762, '2025-11-30 23:49:38.276', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (262, 1, 'devops', 'cmdb_sql', 1, 'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 345, '2025-11-30 23:50:09.029', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (263, 1, 'devops', 'cmdb_sql_log', 1, 'SELECT * FROM cmdb_sql_log LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 304, '2025-11-30 23:50:24.695', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (264, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 311, '2025-11-30 23:50:29.549', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (265, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 285, '2025-11-30 23:51:51.702', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (266, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 308, '2025-11-30 23:53:56.382', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (267, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 291, '2025-11-30 23:54:23.552', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (268, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 292, '2025-11-30 23:54:33.359', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (269, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 303, '2025-11-30 23:54:38.550', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (270, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 333, '2025-11-30 23:54:46.965', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (271, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 319, '2025-11-30 23:54:58.898', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (272, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 314, '2025-11-30 23:54:59.768', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (273, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 315, '2025-11-30 23:55:02.322', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (274, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 290, '2025-11-30 23:55:11.683', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (275, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 282, '2025-11-30 23:55:36.034', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (276, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 261, '2025-11-30 23:55:40.374', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (277, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 100;', '', '', 1, '执行成功', 288, '2025-11-30 23:56:00.161', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (278, 1, 'devops', 'cmdb_group;', 1, 'SELECT * FROM cmdb_group;', '', '', 1, '执行成功', 272, '2025-11-30 23:56:57.273', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (279, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 279, '2025-11-30 23:57:18.128', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (280, 1, 'devops', 'cmdb_group', 1, 'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 286, '2025-11-30 23:57:32.189', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (281, 1, 'devops', 'cmdb_group', 2, 'INSERT INTO cmdb_group (parent_id, name, create_time, remark, update_time) VALUES (\'0\', \'test\', \'\', \'123\', \'\')', '', '', 2, 'Error 1292 (22007): Incorrect datetime value: \'\' for column \'create_time\' at row 1', 234, '2025-11-30 23:58:10.839', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (282, 1, 'devops', 'cmdb_host_copy_20251130235846', 1, 'SELECT * FROM cmdb_host_copy_20251130235846 LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 286, '2025-11-30 23:58:49.791', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (283, 1, 'devops', '', 5, 'DROP TABLE cmdb_host_copy_20251130235846', '', '', 1, '执行成功', 295, '2025-11-30 23:58:54.686', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (284, 1, 'devops', '', 5, 'RENAME TABLE cmdb_host_copy_20251130235914 TO cmdb_host_copy_123', '', '', 1, '执行成功', 274, '2025-11-30 23:59:22.577', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (285, 2, 'database_name', 'users', 1, 'SELECT * FROM users LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 207, '2025-11-30 23:59:58.611', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (286, 2, 'database_name', 'users', 2, 'INSERT INTO users (name, email) VALUES (\'李四\', \'lisi@123.com\')', '', '', 1, '执行成功', 285, '2025-12-01 00:00:30.571', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (287, 2, 'database_name', 'users', 1, 'SELECT * FROM users LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 214, '2025-12-01 00:00:31.076', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (288, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 280, '2025-12-01 00:05:29.288', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (289, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 295, '2025-12-01 00:06:23.704', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (290, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 1184, '2025-12-01 00:06:54.903', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (291, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 611, '2025-12-01 00:15:50.948', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (292, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 291, '2025-12-01 00:17:04.651', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (293, 1, 'devops', 'app_application', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 303, '2025-12-01 00:19:18.036', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (294, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 417, '2025-12-01 00:19:22.279', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (295, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 280, '2025-12-01 00:20:01.970', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (296, 1, 'devops', 'cmdb_host', 3, 'UPDATE cmdb_host SET host_name = \'虚拟机001111\', group_id = 12, private_ip = \'172.16.226.16\', public_ip = \'120.231.244.158\', ssh_name = \'root\', ssh_key_id = 13, ssh_port = 22, remark = \'123\', vendor = 1, region = \'\', instance_id = \'\', os = \'CentOSLinux7(Core)\', status = 1, cpu = \'2\', memory = \'4\', disk = \'17\', billing_type = \'\', create_time = \'2025-11-23 23:44:35\', expire_time = NULL, update_time = \'2025-11-23 23:45:32\', ssh_ip = \'172.16.226.16\', name = \'jenkins\', ssh_gateway_id = NULL WHERE id = 511;', '', '', 1, '执行成功', 254, '2025-12-01 00:20:08.106', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (297, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 292, '2025-12-01 00:20:08.651', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (298, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 301, '2025-12-01 00:20:11.764', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (299, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 100;', '', '', 1, '执行成功', 293, '2025-12-01 00:20:31.016', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (300, 1, 'devops', 'app_application', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 302, '2025-12-01 00:21:51.502', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (301, 1, 'devops', 'app_application', 1, 'SELECT * FROM app_application LIMIT 100;', '', '', 1, '执行成功', 569, '2025-12-01 00:21:57.257', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (302, 1, 'devops', 'cmdb_group', 1, 'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 367, '2025-12-01 00:23:39.693', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (303, 1, 'devops', 'cmdb_group', 1, 'SELECT * FROM cmdb_group LIMIT 100;', '', '', 1, '执行成功', 540, '2025-12-01 00:23:47.277', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (304, 2, 'database_name', 'users', 1, 'SELECT * FROM users LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 217, '2025-12-01 00:24:18.173', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (305, 2, 'database_name', 'users', 1, 'SELECT * FROM users LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 289, '2025-12-01 00:25:05.373', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (306, 2, 'database_name', 'users', 2, 'INSERT INTO users (name, email) VALUES (\'test\', \'123@123.com\')', '', '', 1, '执行成功', 553, '2025-12-01 00:30:03.796', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (307, 2, 'database_name', 'users', 1, 'SELECT * FROM users LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 205, '2025-12-01 00:30:04.267', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (308, 2, 'database_name', 'users', 1, 'SELECT * FROM users LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 256, '2025-12-01 00:30:24.140', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (309, 2, 'database_name', 'users', 1, 'SELECT * FROM users LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 452, '2025-12-01 00:37:20.181', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (310, 2, 'database_name', 'users', 1, 'SELECT * FROM users LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 235, '2025-12-01 00:37:27.599', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (311, 2, 'database_name', 'users_copy_20251201003823', 1, 'SELECT * FROM users_copy_20251201003823 LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 197, '2025-12-01 00:38:47.198', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (312, 2, 'database_name', 'users_copy_20251201003823', 1, 'SELECT * FROM users_copy_20251201003823 LIMIT 100;', '', '', 1, '执行成功', 227, '2025-12-01 00:38:56.092', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (313, 1, 'RECOVER_YOUR_DATA', '', 5, 'CREATE DATABASE IF NOT EXISTS ops \n  DEFAULT CHARACTER SET utf8mb4;', '', '', 1, '执行成功', 729, '2025-12-01 00:43:24.905', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (314, 1, 'devops', 'cmdb_group', 1, 'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 296, '2025-12-01 00:43:58.889', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (315, 1, 'devops', 'cmdb_host_copy_123', 1, 'SELECT * FROM cmdb_host_copy_123 LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 351, '2025-12-01 00:44:00.240', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (316, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 320, '2025-12-01 00:44:02.951', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (317, 1, 'devops', 'cmdb_host', 3, 'UPDATE cmdb_host SET host_name = \'虚拟机01\', group_id = 12, private_ip = \'172.16.226.16\', public_ip = \'120.231.244.158\', ssh_name = \'root\', ssh_key_id = 13, ssh_port = 22, remark = \'123\', vendor = 1, region = \'\', instance_id = \'\', os = \'CentOSLinux7(Core)\', status = 1, cpu = \'2\', memory = \'4\', disk = \'17\', billing_type = \'\', create_time = \'2025-11-23 23:44:35\', expire_time = NULL, update_time = \'2025-11-23 23:45:32\', ssh_ip = \'172.16.226.16\', name = \'jenkins\', ssh_gateway_id = NULL WHERE id = 511;', '', '', 1, '执行成功', 722, '2025-12-01 00:44:23.250', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (318, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 333, '2025-12-01 00:44:26.196', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (319, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 297, '2025-12-01 00:44:31.613', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (320, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 316, '2025-12-01 00:44:46.505', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (321, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 475, '2025-12-01 00:46:30.077', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (322, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 326, '2025-12-01 00:46:48.843', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (323, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 290, '2025-12-01 00:46:53.212', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (324, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 296, '2025-12-01 00:46:57.044', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (325, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 306, '2025-12-01 00:47:00.548', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (326, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 296, '2025-12-01 00:47:09.235', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (327, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 327, '2025-12-01 00:47:22.518', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (328, 1, 'devops', 'cmdb_host;', 1, 'SELECT * FROM  cmdb_host;', '', '', 1, '执行成功', 297, '2025-12-01 00:48:12.389', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (329, 1, 'devops', '', 5, 'DROP TABLE cmdb_host_copy_123', '', '', 1, '执行成功', 578, '2025-12-01 00:49:01.428', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (330, 1, 'devops', 'cmdb_host_copy_20251201004911', 1, 'SELECT * FROM cmdb_host_copy_20251201004911 LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 321, '2025-12-01 00:49:15.842', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (331, 1, 'devops', '', 5, 'RENAME TABLE cmdb_host_copy_20251201004911 TO cmdb_host_copy_123', '', '', 1, '执行成功', 253, '2025-12-01 00:49:25.590', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (332, 1, 'devops', '', 5, 'DROP TABLE cmdb_host_copy_123', '', '', 1, '执行成功', 259, '2025-12-01 00:49:43.125', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (333, 2, 'database_name', '', 5, 'CREATE TABLE students (\n    id SERIAL PRIMARY KEY,                -- 学生ID，自增主键\n    student_id VARCHAR(20) UNIQUE NOT NULL, -- 学号，唯一且非空\n    name VARCHAR(50) NOT NULL,            -- 姓名\n    gender CHAR(1) CHECK (gender IN (\'M\', \'F\')), -- 性别：M 男，F 女\n    birth_date DATE,                      -- 出生日期\n    email VARCHAR(100) UNIQUE,            -- 邮箱，唯一\n    phone VARCHAR(20),                    -- 电话\n    enrollment_date DATE DEFAULT CURRENT_DATE, -- 入学日期，默认为当前日期\n    major VARCHAR(100),                   -- 专业\n    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- 记录创建时间\n);', '', '', 1, '执行成功', 239, '2025-12-01 00:50:33.682', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (334, 2, 'database_name', 'students', 1, 'SELECT * FROM students LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 222, '2025-12-01 00:50:38.712', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (335, 2, 'database_name', 'students', 2, 'INSERT INTO students (student_id, name, gender, birth_date, email, phone, enrollment_date, major) VALUES (\'1\', \'张三\', \'m\', \'2004-05-15\', \'123@456.com\', \'12345678911\', \'2004-05-15\', \'软件工程\')', '', '', 2, 'pq: new row for relation \"students\" violates check constraint \"students_gender_check\"', 186, '2025-12-01 00:52:04.216', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (336, 2, 'database_name', 'students', 2, 'INSERT INTO students (student_id, name, gender, birth_date, email, phone, enrollment_date, major) VALUES (\'20230001\', \'张三\', \'m\', \'2004-05-15\', \'123@456.com\', \'13800138001\', \'2004-05-15\', \'软件工程\')', '', '', 2, 'pq: new row for relation \"students\" violates check constraint \"students_gender_check\"', 182, '2025-12-01 00:52:31.225', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (337, 2, 'database_name', 'students', 2, 'INSERT INTO students (student_id, name, gender, birth_date, email, phone, enrollment_date, major) VALUES (\'20230001\', \'张三\', \'\'\'M\'\'\', \'2004-05-15\', \'123@456.com\', \'13800138001\', \'2004-05-15\', \'软件工程\')', '', '', 2, 'pq: value too long for type character(1)', 214, '2025-12-01 00:52:54.401', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (338, 2, 'database_name', 'users', 1, 'SELECT * FROM users LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 542, '2025-12-01 00:52:58.529', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (339, 2, 'database_name', 'users', 2, 'INSERT INTO users (name, email) VALUES (\'test123\', \'123@123.com\')', '', '', 1, '执行成功', 183, '2025-12-01 00:53:22.992', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (340, 2, 'database_name', 'users', 1, 'SELECT * FROM users LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 253, '2025-12-01 00:53:23.494', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (341, 2, 'database_name', 'users', 1, 'SELECT * FROM users LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 207, '2025-12-01 00:53:33.443', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (342, 2, 'database_name', 'users', 1, 'SELECT * FROM users LIMIT 100;', '', '', 1, '执行成功', 196, '2025-12-01 00:53:54.378', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (343, 2, 'database_name', 'users_copy_20251201003823', 1, 'SELECT * FROM users_copy_20251201003823 LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 231, '2025-12-01 00:53:59.552', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (344, 2, 'database_name', '', 5, 'DROP TABLE users_copy_20251201003823', '', '', 1, '执行成功', 184, '2025-12-01 00:54:05.164', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (345, 1, 'devops', 'app_application', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 286, '2025-12-01 10:18:58.108', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (346, 1, 'devops', 'app_application', 1, 'SELECT * FROM app_application LIMIT 100;', '', '', 1, '执行成功', 257, '2025-12-01 10:19:34.501', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (347, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 265, '2025-12-01 12:05:45.106', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (348, 1, 'RECOVER_YOUR_DATA', 'RECOVER_YOUR_DATA', 1, 'SELECT * FROM RECOVER_YOUR_DATA LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 288, '2025-12-01 14:07:15.387', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (349, 1, 'devops', 'app_jenkins_env', 1, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 256, '2025-12-01 14:07:23.511', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (350, 1, 'gin-api', 'app_sh_release', 1, 'SELECT * FROM app_sh_release LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 254, '2025-12-01 17:10:51.225', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (351, 1, 'gin-api', 'app_service_release', 1, 'SELECT * FROM app_service_release LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 245, '2025-12-01 17:11:07.319', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (352, 1, 'gin-api', 'app_application', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 253, '2025-12-01 17:11:39.058', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (353, 1, 'gin-api', 'app_jenkins_env', 1, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 248, '2025-12-01 17:11:57.884', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (354, 1, 'gin-api', 'quick_deployment_tasks', 1, 'SELECT * FROM quick_deployment_tasks LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 280, '2025-12-01 17:12:31.887', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (355, 1, 'gin-api', 'quick_deployments', 1, 'SELECT * FROM quick_deployments LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 254, '2025-12-01 17:12:34.351', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (356, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 248, '2025-12-01 19:50:46.392', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (357, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 294, '2025-12-02 11:13:02.970', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (358, 1, 'devops', 'app_application', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 265, '2025-12-02 11:34:17.783', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (359, 1, 'devops', 'app_application', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 256, '2025-12-02 11:35:41.192', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (360, 1, 'gin-api', 'sys_operation_log', 1, 'SELECT * FROM sys_operation_log LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 292, '2025-12-02 12:01:05.583', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (361, 1, 'gin-api', 'sys_operation_log', 1, 'SELECT * FROM sys_operation_log LIMIT 100;', '', '', 1, '执行成功', 555, '2025-12-02 12:06:45.939', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (362, 1, 'devops', 'cmdb_group', 1, 'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 257, '2025-12-02 12:08:01.959', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (363, 1, 'devops', 'app_application', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 250, '2025-12-02 12:25:37.979', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (364, 1, 'devops', 'app_application', 1, 'SELECT *\nFROM  app_application\nLIMIT  100;', '', '', 1, '执行成功', 252, '2025-12-02 12:25:58.932', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (365, 1, 'devops', 'app_application', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 277, '2025-12-02 12:30:05.820', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (366, 2, 'database_name', 'users', 1, 'SELECT * FROM users LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 253, '2025-12-02 12:35:09.666', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (367, 2, 'database_name', 'students', 1, 'SELECT * FROM students LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 180, '2025-12-02 12:35:11.544', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (368, 1, 'devops', 'cmdb_group', 1, 'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 262, '2025-12-02 12:38:20.200', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (369, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 251, '2025-12-02 12:38:39.691', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (370, 1, 'gin-api', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 250, '2025-12-02 14:59:49.691', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (371, 1, 'gin-api', 'db', 1, 'SELECT * FROM db LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 261, '2025-12-02 15:00:14.754', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (372, 1, 'gin-api', 'db_instance', 1, 'SELECT * FROM db_instance LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 258, '2025-12-02 15:01:14.757', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (373, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 255, '2025-12-02 15:09:10.006', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (374, 1, 'gin-api', 'cmdb_sql', 1, 'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 247, '2025-12-02 15:41:50.960', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (375, 1, 'gin-api', 'app_application', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 251, '2025-12-02 15:42:08.133', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (376, 1, 'gin-api', 'app_application', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 252, '2025-12-02 15:42:47.779', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (377, 1, 'RECOVER_YOUR_DATA', 'RECOVER_YOUR_DATA', 1, 'SELECT * FROM RECOVER_YOUR_DATA LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 284, '2025-12-02 16:00:08.047', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (378, 1, 'mayfly-go', 't_db_backup_history', 1, 'SELECT * FROM t_db_backup_history LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 284, '2025-12-02 16:00:14.480', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (379, 1, 'mayfly-go', 't_db_restore', 1, 'SELECT * FROM t_db_restore LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 283, '2025-12-02 16:00:16.976', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (380, 3, 'gin-api', 'app_jenkins_env', 1, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 256, '2025-12-02 16:59:07.917', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (381, 3, 'gin-api', 'app_service_release_item', 1, 'SELECT * FROM app_service_release_item LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 261, '2025-12-02 16:59:08.689', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (382, 1, 'gin-api', 'cmdb_group', 1, 'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 265, '2025-12-03 10:26:41.230', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (383, 1, 'gin-api', 'app_service_release_item', 1, 'SELECT * FROM app_service_release_item LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 323, '2025-12-03 10:27:20.235', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (384, 1, 'gin-api', 'app_sh_release', 1, 'SELECT * FROM app_sh_release LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 256, '2025-12-03 10:27:21.025', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (385, 1, 'gin-api', 'cmdb_sql_log', 1, 'SELECT * FROM cmdb_sql_log LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 548, '2025-12-03 10:31:37.487', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (386, 1, 'devops', 'app_application', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 263, '2025-12-03 10:36:57.065', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (387, 1, 'devops', 'cmdb_group', 1, 'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 340, '2025-12-03 10:37:02.051', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (388, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 266, '2025-12-03 10:40:18.048', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (389, 1, 'devops', 'cmdb_host_1', 2, 'INSERT INTO cmdb_host_1 (id, host_name, group_id, private_ip, public_ip, ssh_name, ssh_key_id, ssh_port, remark, vendor, region, instance_id, os, status, cpu, memory, disk, billing_type, create_time, expire_time, update_time, ssh_ip, name, ssh_gateway_id) VALUES (506, \'华为云ops\', 4, \'172.31.6.35\', \'139.9.205.38\', \'root\', 22, 22, \'123\', 5, \'\', \'\', \'Ubuntu24.04.2\', 1, \'2\', \'2\', \'40\', \'\', \'2025-11-11 17:24:36\', NULL, \'2025-11-26 17:22:07\', \'139.9.205.38\', \'hw-ops\', NULL);', '', '', 2, 'Error 1146 (42S02): Table \'devops.cmdb_host_1\' doesn\'t exist', 202, '2025-12-03 10:43:12.842', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (390, 1, 'devops', 'cmdb_host_1', 2, 'INSERT INTO cmdb_host_1 (id, host_name, group_id, private_ip, public_ip, ssh_name, ssh_key_id, ssh_port, remark, vendor, region, instance_id, os, status, cpu, memory, disk, billing_type, create_time, expire_time, update_time, ssh_ip, name, ssh_gateway_id) VALUES (506, \'华为云ops\', 4, \'172.31.6.35\', \'139.9.205.38\', \'root\', 22, 22, \'123\', 5, \'\', \'\', \'Ubuntu24.04.2\', 1, \'2\', \'2\', \'40\', \'\', \'2025-11-11 17:24:36\', NULL, \'2025-11-26 17:22:07\', \'139.9.205.38\', \'hw-ops\', NULL);', '', '', 2, 'Error 1146 (42S02): Table \'devops.cmdb_host_1\' doesn\'t exist', 216, '2025-12-03 10:43:14.148', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (391, 1, 'devops', 'cmdb_sql', 1, 'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 785, '2025-12-03 11:09:19.304', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (392, 1, 'devops', 'app_jenkins_env', 1, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 265, '2025-12-03 11:24:14.607', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (393, 1, 'gin-api', 'app_jenkins_env', 1, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 275, '2025-12-03 11:35:15.655', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (394, 1, 'mayfly-go', 't_db', 1, 'SELECT * FROM t_db LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 266, '2025-12-03 11:40:08.029', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (395, 1, 'mayfly-go', 't_db_backup', 1, 'SELECT * FROM t_db_backup LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 267, '2025-12-03 11:40:10.424', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (396, 1, 'mayfly-go', 't_db_data_sync_log', 1, 'SELECT * FROM t_db_data_sync_log LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 416, '2025-12-03 11:41:35.563', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (397, 1, 'mayfly-go', 't_db_instance', 1, 'SELECT * FROM t_db_instance LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 259, '2025-12-03 11:41:38.235', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (398, 1, 'gin-api', 'cmdb_group', 1, 'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 256, '2025-12-03 11:41:41.904', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (399, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 257, '2025-12-03 11:41:59.928', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (400, 1, 'devops', 'cmdb_group', 1, 'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 269, '2025-12-03 11:43:46.442', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (401, 1, 'devops', 'app_jenkins_env', 1, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 317, '2025-12-03 11:45:13.925', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (402, 1, 'devops', '', 5, 'CREATE TABLE `app_jenkins_env123` ( `id` bigint(20) NOT NULL AUTO_INCREMENT, `app_id` bigint(20) NOT NULL, `env_name` varchar(50) NOT NULL, `jenkins_server_id` bigint(20), `job_name` varchar(255), `job_url` varchar(500), `build_params` json, `deploy_config` json, `notification` json, `is_active` tinyint(3) DEFAULT 1, `created_at` datetime, `updated_at` datetime, `deleted_at` datetime, PRIMARY KEY (id) ); ALTER TABLE `app_jenkins_env` ADD INDEX `idx_app_jenkins_env_app_id`(`app_id`) USING BTREE; ALTER TABLE `app_jenkins_env` ADD INDEX `idx_app_jenkins_env_deleted_at`(`deleted_at`) USING BTREE', '', '', 2, 'Error 1064 (42000): You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near \'ALTER TABLE `app_jenkins_env` ADD INDEX `idx_app_jenkins_env_app_id`(`app_id`) U\' at line 1', 204, '2025-12-03 11:47:45.166', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (403, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 270, '2025-12-03 11:48:51.742', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (404, 1, 'devops', 'cmdb_sql', 1, 'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 267, '2025-12-03 11:52:55.287', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (405, 1, 'devops', 'app_jenkins_env', 1, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 257, '2025-12-03 11:52:58.121', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (406, 1, 'devops', 'app_application', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 260, '2025-12-03 11:53:00.314', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (407, 1, 'gin-api', 'cmdb_sql', 1, 'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 295, '2025-12-03 11:56:02.107', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (408, 1, 'devops', 'cmdb_group', 1, 'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 256, '2025-12-03 12:04:05.208', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (409, 1, 'devops', 'app_jenkins_env', 1, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 254, '2025-12-03 12:07:28.099', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (410, 1, 'gin-api', 'app_sh_release', 1, 'SELECT * FROM app_sh_release LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 261, '2025-12-03 12:08:37.417', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (411, 1, 'devops', 'cmdb_sql', 1, 'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 327, '2025-12-03 12:08:52.889', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (412, 1, 'devops', 'cmdb_group', 1, 'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 258, '2025-12-03 12:14:02.575', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (413, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 342, '2025-12-03 12:16:54.823', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (414, 1, 'gin-api', 'app_jenkins_env', 1, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 261, '2025-12-03 12:19:56.646', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (415, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 257, '2025-12-03 12:23:13.577', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (416, 1, 'devops', 'cmdb_sql', 1, 'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 261, '2025-12-03 12:26:28.689', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (417, 1, 'devops', 'app_application', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 253, '2025-12-03 12:28:20.840', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (418, 1, 'gin-api', 'app_jenkins_env', 1, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 244, '2025-12-03 12:29:04.357', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (419, 1, 'devops', 'app_jenkins_env', 1, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 395, '2025-12-03 12:32:00.866', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (420, 1, 'devops', '', 5, 'CREATE TABLE `cmdb_group123` (\n `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT \'\'\'主键\'\'\',\n `parent_id` bigint(20) NOT NULL DEFAULT 0 COMMENT \'\'\'父级分组ID\'\'\',\n `name` longtext NOT NULL COMMENT \'\'\'分组名称\'\'\',\n `create_time` datetime NOT NULL COMMENT \'\'\'创建时间\'\'\',\n `remark` longtext COMMENT \'\'\'备注\'\'\',\n `update_time` datetime COMMENT \'\'\'更新时间\'\'\', \nPRIMARY KEY (id)\n)', '', '', 1, '执行成功', 244, '2025-12-03 12:42:43.341', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (421, 1, 'devops', 'cmdb_group123', 1, 'SELECT * FROM cmdb_group123 LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 262, '2025-12-03 12:42:51.081', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (422, 1, 'devops', '', 5, 'DROP TABLE cmdb_group123', '', '', 1, '执行成功', 228, '2025-12-03 12:43:02.777', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (423, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 265, '2025-12-03 12:46:17.553', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (424, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 278, '2025-12-03 15:54:52.850', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (425, 1, 'devops', 'cmdb_sql_log', 1, 'SELECT * FROM cmdb_sql_log LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 252, '2025-12-03 15:55:18.939', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (426, 1, 'gin-api', 'sys_operation_log', 1, 'SELECT * FROM sys_operation_log LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 254, '2025-12-03 15:55:30.381', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (427, 1, 'gin-api', 'sys_operation_log', 1, 'SELECT *\nFROM  sys_operation_log\nLIMIT  100;', '', '', 1, '执行成功', 266, '2025-12-03 15:55:49.290', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (428, 1, 'gin-api', 'sys_admin', 1, 'SELECT *\nFROM  sys_admin\nLIMIT  100;', '', '', 1, '执行成功', 250, '2025-12-03 15:57:53.622', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (429, 1, 'devops', 'app_application', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 256, '2025-12-05 20:29:40.131', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (430, 1, 'devops', 'app_application', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 262, '2025-12-05 20:32:28.584', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (431, 1, 'devops', 'app_application', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 392, '2025-12-05 20:33:31.694', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (432, 1, 'devops', 'app_application', 1, 'SELECT *\nFROM  app_application\nLIMIT  100;', '', '', 1, '执行成功', 308, '2025-12-05 20:33:51.765', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (433, 1, 'devops', 'app_jenkins_env', 1, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 424, '2025-12-05 20:34:23.519', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (434, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 273, '2025-12-05 20:34:30.817', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (435, 1, 'devops', 'cmdb_host', 3, 'UPDATE cmdb_host SET host_name = \'虚拟机2025\', group_id = 4, private_ip = \'172.16.226.13\', public_ip = \'120.231.244.158\', ssh_name = \'root\', ssh_key_id = 13, ssh_port = 22, remark = \'123\', vendor = 1, region = \'\', instance_id = \'\', os = \'Ubuntu20.04.2\', status = 1, cpu = \'2\', memory = \'3\', disk = \'19\', billing_type = \'\', create_time = \'2025-11-23 23:46:03\', expire_time = NULL, update_time = \'2025-11-23 23:46:07\', ssh_ip = \'172.16.226.13\', name = \'k8s-node02\', ssh_gateway_id = NULL WHERE id = 512;', '', '', 1, '执行成功', 436, '2025-12-05 20:34:45.595', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (436, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 267, '2025-12-05 20:34:46.094', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (437, 1, 'devops', 'cmdb_host', 1, 'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 267, '2025-12-05 20:34:49.893', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (438, 1, 'devops', '', 5, 'CREATE DATABASE IF NOT EXISTS test1\n  DEFAULT CHARACTER SET utf8mb4;', '', '', 1, '执行成功', 237, '2025-12-05 20:35:34.369', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (439, 1, 'test1', '', 5, 'CREATE TABLE test1.your_table_name (\n  id BIGINT PRIMARY KEY AUTO_INCREMENT,\n  -- your columns here\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;', '', '', 2, 'Error 1064 (42000): You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near \') ENGINE=InnoDB DEFAULT CHARSET=utf8mb4\' at line 4', 219, '2025-12-05 20:35:52.483', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (440, 2, 'database_name', 'users', 1, 'SELECT * FROM users LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 217, '2025-12-05 20:36:57.655', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (441, 2, 'database_name', 'users', 2, 'INSERT INTO users (name, email) VALUES (\'王五\', \'xxxx@123.com\')', '', '', 1, '执行成功', 191, '2025-12-05 20:37:14.770', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (442, 2, 'database_name', 'users', 1, 'SELECT * FROM users LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 205, '2025-12-05 20:37:15.207', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (443, 1, 'devops', 'app_application', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 267, '2025-12-08 10:52:11.033', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (444, 1, 'gin-api', 'app_application', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 288, '2025-12-08 23:01:38.354', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (445, 1, 'gin-api', 'app_jenkins_env', 1, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 267, '2025-12-08 23:01:38.638', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (446, 1, 'gin-api', 'sys_config', 1, 'SELECT * FROM sys_config LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 259, '2025-12-09 13:21:51.379', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (447, 1, 'gin-api', 'app_application', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 280, '2025-12-09 23:45:28.566', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (448, 1, 'gin-api', 'app_jenkins_env', 1, 'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 268, '2025-12-09 23:48:06.048', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (449, 1, 'devops', 'app_application', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 267, '2025-12-10 11:22:14.627', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (450, 1, 'gin-api', 'sys_admin', 1, 'SELECT * FROM sys_admin LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 268, '2025-12-10 17:17:50.042', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (451, 1, 'gin-api', 'app_application', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 276, '2025-12-11 18:04:01.397', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (452, 1, 'devops', 'app_application', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 314, '2025-12-13 00:03:02.341', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (453, 2, 'database_name', 'users', 1, 'SELECT * FROM users LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 287, '2025-12-13 00:03:27.273', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (454, 2, 'database_name', 'students', 1, 'SELECT * FROM students LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 257, '2025-12-13 00:03:29.263', '', 0);
INSERT INTO `db_sql_exec` (`id`, `db_id`, `db_name`, `table_name`, `type`, `sql`, `old_value`, `remark`, `status`, `res`, `exec_time`, `create_time`, `creator`, `creator_id`) VALUES (455, 1, 'devops', 'app_application', 1, 'SELECT * FROM app_application LIMIT 20 OFFSET 0', '', '', 1, '执行成功', 303, '2025-12-13 16:15:43.928', '', 0);
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
) ENGINE=InnoDB AUTO_INCREMENT=34 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

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
) ENGINE=InnoDB AUTO_INCREMENT=86 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of monitor_agent
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for monitor_domain
-- ----------------------------
DROP TABLE IF EXISTS `monitor_domain`;
CREATE TABLE `monitor_domain` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `domain` varchar(255) NOT NULL,
  `tags` varchar(500) DEFAULT NULL,
  `remark` text,
  `status` bigint DEFAULT '1' COMMENT '状态(1:启用,0:禁用)',
  `is_alive` bigint DEFAULT '0' COMMENT '存活状态(1:正常,0:异常)',
  `status_code` bigint DEFAULT NULL COMMENT 'HTTP状态码',
  `response_time` bigint DEFAULT NULL COMMENT '响应时间(ms)',
  `ssl_expire_at` datetime DEFAULT NULL COMMENT 'SSL证书过期时间',
  `ssl_days_left` bigint DEFAULT NULL COMMENT 'SSL证书剩余天数',
  `ssl_issuer` varchar(255) DEFAULT NULL COMMENT 'SSL证书颁发者',
  `last_check_at` datetime DEFAULT NULL COMMENT '最后检查时间',
  `error_msg` text COMMENT '错误信息',
  `create_time` datetime(3) DEFAULT NULL,
  `update_time` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_monitor_domain_domain` (`domain`)
) ENGINE=InnoDB AUTO_INCREMENT=41 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of monitor_domain
-- ----------------------------
BEGIN;
INSERT INTO `monitor_domain` (`id`, `domain`, `tags`, `remark`, `status`, `is_alive`, `status_code`, `response_time`, `ssl_expire_at`, `ssl_days_left`, `ssl_issuer`, `last_check_at`, `error_msg`, `create_time`, `update_time`) VALUES (6, 'gitee.com', '', '', 1, 1, 200, 530, '2026-02-17 23:59:59', 70, 'TrustAsia DV TLS RSA CA 2025', '2025-12-09 17:05:50', '', '2025-12-04 10:21:13.921', '2025-12-09 17:05:50.292');
INSERT INTO `monitor_domain` (`id`, `domain`, `tags`, `remark`, `status`, `is_alive`, `status_code`, `response_time`, `ssl_expire_at`, `ssl_days_left`, `ssl_issuer`, `last_check_at`, `error_msg`, `create_time`, `update_time`) VALUES (7, 'ops.dding.net', '', '', 1, 1, 200, 228, '2025-12-20 23:59:59', 11, 'Encryption Everywhere DV TLS CA - G1', '2025-12-09 17:05:50', '', '2025-12-04 10:21:28.741', '2025-12-09 17:05:49.990');
INSERT INTO `monitor_domain` (`id`, `domain`, `tags`, `remark`, `status`, `is_alive`, `status_code`, `response_time`, `ssl_expire_at`, `ssl_days_left`, `ssl_issuer`, `last_check_at`, `error_msg`, `create_time`, `update_time`) VALUES (8, 'docker.aityp.com', '', '', 1, 1, 200, 559, '2025-12-28 23:58:52', 19, 'E8', '2025-12-09 17:05:50', '', '2025-12-04 10:21:54.359', '2025-12-09 17:05:50.321');
INSERT INTO `monitor_domain` (`id`, `domain`, `tags`, `remark`, `status`, `is_alive`, `status_code`, `response_time`, `ssl_expire_at`, `ssl_days_left`, `ssl_issuer`, `last_check_at`, `error_msg`, `create_time`, `update_time`) VALUES (12, 'www.jd.com', 'ops', '123', 1, 1, 200, 60, '2026-12-20 04:28:34', 375, 'GlobalSign RSA OV SSL CA 2018', '2025-12-09 17:05:51', '', '2025-12-05 20:50:58.505', '2025-12-09 17:05:50.638');
INSERT INTO `monitor_domain` (`id`, `domain`, `tags`, `remark`, `status`, `is_alive`, `status_code`, `response_time`, `ssl_expire_at`, `ssl_days_left`, `ssl_issuer`, `last_check_at`, `error_msg`, `create_time`, `update_time`) VALUES (28, 'ai.feishu.cn', '', '', 1, 1, 403, 567, '2026-03-25 23:59:59', 101, 'RapidSSL TLS RSA CA G1', '2025-12-14 14:31:39', '', '2025-12-08 22:30:11.616', '2025-12-14 14:31:39.462');
COMMIT;

-- ----------------------------
-- Table structure for monitor_domain_schedule
-- ----------------------------
DROP TABLE IF EXISTS `monitor_domain_schedule`;
CREATE TABLE `monitor_domain_schedule` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `enabled` tinyint(1) DEFAULT '0' COMMENT '是否启用',
  `interval_mins` bigint DEFAULT NULL COMMENT '检查间隔(分钟)',
  `next_run_at` datetime DEFAULT NULL COMMENT '下次执行时间',
  `last_run_at` datetime DEFAULT NULL COMMENT '上次执行时间',
  `status` varchar(50) DEFAULT NULL COMMENT '任务状态(running/paused/stopped)',
  `create_time` datetime(3) DEFAULT NULL,
  `update_time` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of monitor_domain_schedule
-- ----------------------------
BEGIN;
INSERT INTO `monitor_domain_schedule` (`id`, `enabled`, `interval_mins`, `next_run_at`, `last_run_at`, `status`, `create_time`, `update_time`) VALUES (1, 0, 2, NULL, NULL, 'stopped', '2025-12-04 23:09:32.660', '2025-12-05 23:12:27.687');
INSERT INTO `monitor_domain_schedule` (`id`, `enabled`, `interval_mins`, `next_run_at`, `last_run_at`, `status`, `create_time`, `update_time`) VALUES (2, 0, 120, NULL, NULL, 'stopped', '2025-12-04 23:09:32.660', '2025-12-04 23:09:32.660');
COMMIT;

-- ----------------------------
-- Table structure for monitor_incident
-- ----------------------------
DROP TABLE IF EXISTS `monitor_incident`;
CREATE TABLE `monitor_incident` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `alert_time` datetime(3) NOT NULL COMMENT '''告警时间''',
  `business_line` longtext COMMENT '''业务线(saas3/saas4/saas5等)''',
  `frequency` longtext COMMENT '''频率(偶发/频繁)''',
  `alert_desc` text COMMENT '''告警描述''',
  `alert_level` varchar(191) DEFAULT 'P4' COMMENT '''告警级别(P1/P2/P3/P4)''',
  `incident_cause` text COMMENT '''故障原因''',
  `department` longtext COMMENT '''所属部门(研发部/运维部)''',
  `solution` text COMMENT '''解决方案''',
  `detail_url` longtext COMMENT '''故障详情URL链接''',
  `handler` longtext COMMENT '''处理人(用户名)''',
  `handler_id` bigint unsigned DEFAULT NULL COMMENT '''处理人ID''',
  `status` bigint DEFAULT '1' COMMENT '''处理状态:1->未处理,2->处理中,3->已完成''',
  `remark` text COMMENT '''备注信息''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `update_time` datetime(3) DEFAULT NULL COMMENT '''更新时间''',
  `business_line_id` bigint unsigned DEFAULT NULL COMMENT '''业务线ID(关联cmdb_group二级分组)''',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of monitor_incident
-- ----------------------------
BEGIN;
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
  `parameters` text,
  PRIMARY KEY (`id`),
  KEY `idx_quick_deployment_tasks_deployment_id` (`deployment_id`),
  KEY `fk_quick_deployment_tasks_application` (`app_id`),
  KEY `fk_quick_deployment_tasks_jenkins_env` (`jenkins_env_id`),
  CONSTRAINT `fk_quick_deployment_tasks_application` FOREIGN KEY (`app_id`) REFERENCES `app_application` (`id`),
  CONSTRAINT `fk_quick_deployment_tasks_jenkins_env` FOREIGN KEY (`jenkins_env_id`) REFERENCES `app_jenkins_env` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_quick_deployments_tasks` FOREIGN KEY (`deployment_id`) REFERENCES `quick_deployments` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=44 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of quick_deployment_tasks
-- ----------------------------
BEGIN;
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
) ENGINE=InnoDB AUTO_INCREMENT=34 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of quick_deployments
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for redis_instance
-- ----------------------------
DROP TABLE IF EXISTS `redis_instance`;
CREATE TABLE `redis_instance` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `code` varchar(36) NOT NULL COMMENT '''实例编码''',
  `name` varchar(100) NOT NULL COMMENT '''实例名称''',
  `mode` varchar(20) NOT NULL COMMENT '''模式:standalone,cluster,sentinel''',
  `host` varchar(300) NOT NULL COMMENT '''主机: standalone为host:port, cluster为逗号分隔, sentinel为master=hosts''',
  `port` bigint DEFAULT '0' COMMENT '''端口(standalone可用)''',
  `db` bigint DEFAULT '0' COMMENT '''默认库号''',
  `username` varchar(100) DEFAULT NULL COMMENT '''用户名(可选)''',
  `password` varchar(500) NOT NULL COMMENT '''密码(加密)''',
  `redis_node_password` varchar(500) DEFAULT '' COMMENT '''节点密码(仅sentinel)''',
  `remark` varchar(500) DEFAULT NULL COMMENT '''备注''',
  `ssh_tunnel_machine_id` bigint unsigned DEFAULT '0' COMMENT '''SSH隧道机器ID''',
  `status` bigint DEFAULT '1' COMMENT '''状态:1->启用,2->禁用''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `update_time` datetime(3) DEFAULT NULL COMMENT '''更新时间''',
  `creator` varchar(64) DEFAULT NULL COMMENT '''创建人''',
  `creator_id` bigint unsigned DEFAULT NULL COMMENT '''创建人ID''',
  `modifier` varchar(64) DEFAULT NULL COMMENT '''修改人''',
  `modifier_id` bigint unsigned DEFAULT NULL COMMENT '''修改人ID''',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_redis_instance_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of redis_instance
-- ----------------------------
BEGIN;
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
) ENGINE=InnoDB AUTO_INCREMENT=106 DEFAULT CHARSET=utf8mb3 ROW_FORMAT=DYNAMIC COMMENT='后台管理员表';

-- ----------------------------
-- Records of sys_admin
-- ----------------------------
BEGIN;
INSERT INTO `sys_admin` (`id`, `post_id`, `dept_id`, `username`, `password`, `nickname`, `icon`, `email`, `phone`, `note`, `create_time`, `status`) VALUES (89, 1, 15, 'admin', 'e10adc3949ba59abbe56e057f20f883e', 'admin', 'http://192.168.3.7:8080/api/v1/upload/20251213/862328000.png', '123456789@qq.com', '13754354536', '后端研发', '2023-05-23 22:15:50', 1);
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
COMMIT;

-- ----------------------------
-- Table structure for sys_config
-- ----------------------------
DROP TABLE IF EXISTS `sys_config`;
CREATE TABLE `sys_config` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `config_key` varchar(100) NOT NULL COMMENT '''配置键''',
  `config_type` varchar(50) NOT NULL COMMENT '''配置类型(ldap,email,sms等)''',
  `config_data` text NOT NULL COMMENT '''配置数据(JSON格式)''',
  `status` bigint NOT NULL DEFAULT '1' COMMENT '''状态:1->启用,2->禁用''',
  `remark` varchar(500) DEFAULT NULL COMMENT '''备注''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `update_time` datetime(3) NOT NULL COMMENT '''更新时间''',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_sys_config_config_key` (`config_key`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of sys_config
-- ----------------------------
BEGIN;
INSERT INTO `sys_config` (`id`, `config_key`, `config_type`, `config_data`, `status`, `remark`, `create_time`, `update_time`) VALUES (1, 'ldap', 'ldap', '{\"enable\":true,\"host\":\"www.baidu.com\",\"port\":389,\"baseDn\":\"dc=dding,dc=cn\",\"bindUser\":\"cn=reader,dc=dding,dc=cn\",\"bindPass\":\"fsyunding2018\",\"authFilter\":\"(\\u0026(cn=%s))\",\"coverAttributes\":true,\"tls\":false,\"startTLS\":false,\"defaultRoles\":null,\"defaultRoleId\":13,\"attributes\":{\"nickname\":\"cn\",\"phone\":\"mobile\",\"email\":\"mail\"},\"remark\":\"\"}', 1, '', '2025-12-09 13:19:56.671', '2025-12-14 14:56:54.371');
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
) ENGINE=InnoDB AUTO_INCREMENT=450 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC COMMENT='登录日志记录';

-- ----------------------------
-- Records of sys_login_info
-- ----------------------------
BEGIN;
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
) ENGINE=InnoDB AUTO_INCREMENT=242 DEFAULT CHARSET=utf8mb3 ROW_FORMAT=DYNAMIC COMMENT='菜单表';

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
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (88, 80, '业务分组', 'Shop', 'cmdb:group', 2, 'cmdb/group', 2, 2, '2025-07-16 15:17:14');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (89, 88, '创建分组', '', 'cmdb:group:add', 3, '', 2, 1, '2025-07-18 15:24:31');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (90, 88, '修改分组', '', 'cmdb:group:update', 3, '', 2, 2, '2025-07-18 15:25:49');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (91, 88, '删除分组', '', 'cmdb:group:delete', 3, '', 2, 3, '2025-07-18 15:26:21');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (93, 81, '工作负载', 'Star', 'cloud:k8s:workload', 2, 'k8s/workload', 2, 4, '2025-07-24 14:38:31');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (95, 80, '数据管理', 'Coin', 'cmdb:db', 2, 'cmdb/db', 2, 3, '2025-07-29 15:27:50');
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
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (109, 0, '服务管理', 'ElemeFilled', '', 1, '', 2, 4, '2025-09-16 09:49:55');
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
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (212, 0, '监控中心', 'Shop', '', 1, '', 2, 4, '2025-12-03 21:21:04');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (213, 212, '域名监控', 'UploadFilled', 'monitor:domain', 2, 'monitor/domain', 2, 1, '2025-12-03 21:22:11');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (215, 4, '系统配置', 'List', 'system:config', 2, 'system/config', 2, 6, '2025-12-09 11:03:54');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (216, 212, '故障管理', 'Help', 'monitor:incident', 2, 'monitor/incident', 2, 2, '2025-12-10 15:10:29');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (217, 95, '查看密码', '', 'cmdb:db:passwd', 3, '', 2, 4, '2025-12-13 14:27:22');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (218, 95, '测试链接', '', 'cmdb:db:test', 3, '', 2, 5, '2025-12-13 14:28:38');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (219, 95, 'DBMS终端', '', 'cmdb:db:dbms', 3, '', 2, 6, '2025-12-13 14:32:59');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (220, 95, 'Redis终端', '', 'cmdb:db:redis', 3, '', 2, 7, '2025-12-13 14:33:27');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (221, 95, 'ES终端', '', 'cmdb:db:es', 3, '', 2, 8, '2025-12-13 14:33:53');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (222, 95, 'Mongo终端', '', 'cmdb:db:mongo', 3, '', 2, 9, '2025-12-13 14:34:14');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (229, 213, '新增域名', '', 'monitor:domain:add', 3, '', 2, 1, '2025-12-13 15:22:47');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (230, 213, '批量新增域名', '', 'monitor:domain:add_all', 3, '', 2, 2, '2025-12-13 15:23:44');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (231, 213, '域名自动巡检', '', 'monitor:domain:auto_inspection', 3, '', 2, 3, '2025-12-13 15:25:35');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (232, 213, '域名手动巡检', '', 'monitor:domain:manual_inspection', 3, '', 2, 4, '2025-12-13 15:26:55');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (233, 213, '单域名手动巡检', '', 'monitor:domain:ops_inspection', 3, '', 2, 5, '2025-12-13 15:36:47');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (234, 213, '编辑域名', '', 'monitor:domain:edit', 3, '', 2, 6, '2025-12-13 15:38:30');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (235, 213, '删除域名', '', 'monitor:domain:delete', 3, '', 2, 7, '2025-12-13 15:38:54');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (236, 216, '新增故障记录', '', 'monitor:incident:add', 3, '', 2, 1, '2025-12-13 15:55:41');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (237, 216, '编辑故障记录', '', 'monitor:incident:edit', 3, '', 2, 2, '2025-12-13 15:57:20');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (238, 216, '编辑故障状态', '', 'monitor:incident:status', 3, '', 2, 3, '2025-12-13 15:57:39');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (239, 216, '删除故障记录', '', 'monitor:incident:delete', 3, '', 2, 4, '2025-12-13 15:58:12');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (240, 216, '访问故障链接', '', 'monitor:incident:url', 3, '', 2, 5, '2025-12-13 16:02:46');
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `icon`, `value`, `menu_type`, `url`, `menu_status`, `sort`, `create_time`) VALUES (241, 78, '下载文件', '', 'cmdb:ecs:download', 3, '', 2, 10, '2025-12-13 16:13:17');
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
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb3 ROW_FORMAT=DYNAMIC COMMENT='操作日志记录表';

-- ----------------------------
-- Records of sys_operation_log
-- ----------------------------
BEGIN;
INSERT INTO `sys_operation_log` (`id`, `admin_id`, `username`, `method`, `ip`, `url`, `create_time`, `description`) VALUES (1, 89, 'admin', 'delete', '192.168.3.7', '/api/v1/sysOperationLog/clean', '2025-12-14 14:57:14.270', '清空操作日志');
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
INSERT INTO `sys_role` (`id`, `role_name`, `role_key`, `status`, `description`, `create_time`) VALUES (13, '游客', 'test', 1, 'test1', '2025-07-03 18:47:25');
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
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 241);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 88);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 89);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 90);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 91);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 95);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 146);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 147);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 148);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 217);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 218);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 219);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 220);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 221);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 222);
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
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 212);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 213);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 229);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 230);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 231);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 232);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 233);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 234);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 235);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 216);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 236);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 237);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 238);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 239);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 240);
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
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 215);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 84);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 85);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 123);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 124);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 125);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 86);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 122);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 119);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 120);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 121);
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
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (1, 109);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 72);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 149);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 150);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 154);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 165);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 89);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 146);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 218);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 219);
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
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 229);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 233);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 236);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 240);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 136);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 137);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 156);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 133);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 130);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 132);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 159);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 103);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 127);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 16);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 125);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 122);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 118);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 73);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 62);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 96);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 80);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 78);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 88);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 95);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 81);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 82);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 83);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 105);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 93);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 106);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 109);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 110);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 111);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 212);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 213);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 216);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 97);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 99);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 98);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 100);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 101);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 102);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 4);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 6);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 84);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 85);
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (13, 86);
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
) ENGINE=InnoDB AUTO_INCREMENT=52 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of task_ansible
-- ----------------------------
BEGIN;
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
) ENGINE=InnoDB AUTO_INCREMENT=53 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of task_ansiblework
-- ----------------------------
BEGIN;
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
) ENGINE=InnoDB AUTO_INCREMENT=62 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of task_job
-- ----------------------------
BEGIN;
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
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of task_template
-- ----------------------------
BEGIN;
INSERT INTO `task_template` (`id`, `name`, `type`, `content`, `created_at`, `updated_at`, `created_by`, `updated_by`, `remark`) VALUES (2, '1-数字 1-100', 1, '#!/bin/bash\n\n# 从 1 到 100，每秒打印一个数字\nfor ((i = 1; i <= 100; i++)); do\n    echo \"[$(date +%H:%M:%S)] $i\"\n    sleep 1\ndone\n\necho \"完成：所有数字 1-100 已打印完毕。\"\n', '2025-08-06 12:47:57.073', '2025-08-12 16:14:49.394', 'admin', 'admin', '测试脚本');
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
) ENGINE=InnoDB AUTO_INCREMENT=110 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of task_work
-- ----------------------------
BEGIN;
INSERT INTO `task_work` (`id`, `task_id`, `template_id`, `host_id`, `status`, `log`, `log_path`, `start_time`, `end_time`, `duration`, `created_at`, `type`, `scheduled_time`, `cron_expr`, `is_recurring`) VALUES (90, 48, 11, 1, 2, '[2025-09-29-21:42:02] 任务开始\n进程统计信息如下\n总进程数量为:79\nRunning 进程数为:1\nStoped 进程数为:0\nSleeping 进程数为:55\nZombie 进程数为:0\n[2025-09-29-21:42:02] 任务完成\n', 'logs/task_48/task_48_template_11.log', '2025-09-29 21:42:00.299', '2025-09-29 21:42:06.276', 5, '2025-09-29 20:58:17.259', 0, '2025-09-29 21:45:00.000', NULL, NULL);
INSERT INTO `task_work` (`id`, `task_id`, `template_id`, `host_id`, `status`, `log`, `log_path`, `start_time`, `end_time`, `duration`, `created_at`, `type`, `scheduled_time`, `cron_expr`, `is_recurring`) VALUES (91, 48, 0, 0, 4, '', '', NULL, NULL, 0, '2025-09-29 20:58:17.515', 2, '2025-09-29 21:00:00.000', NULL, NULL);
COMMIT;

-- ----------------------------
-- Table structure for tool_link
-- ----------------------------
DROP TABLE IF EXISTS `tool_link`;
CREATE TABLE `tool_link` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `title` longtext NOT NULL COMMENT '''导航标题''',
  `icon` longtext COMMENT '''导航图标''',
  `link` longtext NOT NULL COMMENT '''链接地址''',
  `sort` bigint DEFAULT '0' COMMENT '''排序''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `update_time` datetime(3) DEFAULT NULL COMMENT '''更新时间''',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of tool_link
-- ----------------------------
BEGIN;
INSERT INTO `tool_link` (`id`, `title`, `icon`, `link`, `sort`, `create_time`, `update_time`) VALUES (2, '百度', 'http://10.7.16.22:8080/api/v1/upload/20251023/625775000.svg', 'https://www.baidu.com/', 0, '2025-10-23 15:54:08.512', '2025-10-23 15:54:08.512');
INSERT INTO `tool_link` (`id`, `title`, `icon`, `link`, `sort`, `create_time`, `update_time`) VALUES (3, 'aws', 'http://10.7.16.22:8080/api/v1/upload/20251023/985806000.svg', 'https://us-west-2.console.aws.amazon.com/eks/clusters?region=us-west-2#', 0, '2025-10-23 15:55:59.543', '2025-10-23 15:55:59.543');
INSERT INTO `tool_link` (`id`, `title`, `icon`, `link`, `sort`, `create_time`, `update_time`) VALUES (4, '美女一号', 'http://10.7.16.22:8080/api/v1/upload/20251023/646700000.png', 'https://gitee.com/zhang_fan1024', 0, '2025-10-23 15:56:22.231', '2025-10-23 15:56:22.231');
INSERT INTO `tool_link` (`id`, `title`, `icon`, `link`, `sort`, `create_time`, `update_time`) VALUES (5, '美女二号', 'http://10.7.16.22:8080/api/v1/upload/20251023/733520000.png', 'https://demo.spug.cc/', 0, '2025-10-23 15:56:48.601', '2025-10-23 15:56:48.601');
INSERT INTO `tool_link` (`id`, `title`, `icon`, `link`, `sort`, `create_time`, `update_time`) VALUES (6, '凡人修仙传', 'http://10.7.16.22:8080/api/v1/upload/20251023/771236000.png', 'http://115.190.10.126/#/dashboard', 0, '2025-10-23 15:57:21.676', '2025-10-23 15:57:21.676');
INSERT INTO `tool_link` (`id`, `title`, `icon`, `link`, `sort`, `create_time`, `update_time`) VALUES (7, '腾讯', 'http://10.7.16.22:8080/api/v1/upload/20251205/92506000.png', 'https://cloud.tencent.com/login?s_url=https%3A%2F%2Fconsole.cloud.tencent.com%2Ftke2%2Fcluster%2Fsub%2Flist%2Fbasic%2Finfo%3Frid%3D1%26clusterId%3Dcls-301g0fi0', 0, '2025-10-23 15:57:55.559', '2025-12-05 20:12:08.966');
COMMIT;

-- ----------------------------
-- Table structure for tool_service_deploy
-- ----------------------------
DROP TABLE IF EXISTS `tool_service_deploy`;
CREATE TABLE `tool_service_deploy` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `service_name` longtext NOT NULL COMMENT '''服务名称''',
  `service_id` longtext NOT NULL COMMENT '''服务ID''',
  `version` longtext NOT NULL COMMENT '''服务版本''',
  `host_id` bigint unsigned NOT NULL COMMENT '''主机ID''',
  `host_ip` longtext NOT NULL COMMENT '''主机IP''',
  `install_dir` longtext NOT NULL COMMENT '''安装目录''',
  `container_name` longtext COMMENT '''容器名称''',
  `ports` longtext COMMENT '''端口映射(JSON)''',
  `env_vars` longtext COMMENT '''环境变量(JSON)''',
  `status` bigint DEFAULT '0' COMMENT '''状态:0->部署中,1->运行中,2->已停止,3->部署失败''',
  `deploy_log` longtext COMMENT '''部署日志''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `update_time` datetime(3) DEFAULT NULL COMMENT '''更新时间''',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of tool_service_deploy
-- ----------------------------
BEGIN;
INSERT INTO `tool_service_deploy` (`id`, `service_name`, `service_id`, `version`, `host_id`, `host_ip`, `install_dir`, `container_name`, `ports`, `env_vars`, `status`, `deploy_log`, `create_time`, `update_time`) VALUES (7, 'Redis', 'redis', 'redis-7.2', 501, '8.130.14.34', '/opt/data/redis', '', '', '{\"REDIS_MAXMEMORY\":\"2gb\",\"REDIS_PASSWORD\":\"redis123456\",\"REDIS_PORT\":\"6370\"}', 1, '[2025-10-30 17:20:01] 开始部署 Redis Redis 7.2\n[2025-10-30 17:20:01] 连接主机 8.130.14.34...\n[2025-10-30 17:20:02] SSH连接成功\n[2025-10-30 17:20:03] 创建安装目录 /opt/data/redis...\n[2025-10-30 17:20:03] 读取模板文件 common/templates/05-redis/versions/redis-7.2-docker-compose.yml...\n[2025-10-30 17:20:03] 生成环境变量配置...\n[2025-10-30 17:20:03] 上传 docker-compose.yml...\n[2025-10-30 17:20:04] 上传 .env...\n[2025-10-30 17:20:04] 检查Docker环境...\n[2025-10-30 17:20:04] 启动服务容器...\n[2025-10-30 17:20:05] 容器启动输出:\n\n[2025-10-30 17:20:05] 验证容器状态...\n[2025-10-30 17:20:08] 容器状态:\nNAME      IMAGE                                                                                      COMMAND                  SERVICE   CREATED         STATUS                            PORTS\nredis72   crpi-aj3vgoxp9kzh2jx4.cn-hangzhou.personal.cr.aliyuncs.com/zhangfan_k8s/redis:7.2-alpine   \"docker-entrypoint.s…\"   redis     4 seconds ago   Up 3 seconds (health: starting)   0.0.0.0:6370->6379/tcp, [::]:6370->6379/tcp\n\n[2025-10-30 17:20:08] 部署完成！\n', '2025-10-30 17:20:01.477', '2025-10-30 17:20:01.477');
INSERT INTO `tool_service_deploy` (`id`, `service_name`, `service_id`, `version`, `host_id`, `host_ip`, `install_dir`, `container_name`, `ports`, `env_vars`, `status`, `deploy_log`, `create_time`, `update_time`) VALUES (11, 'Java', 'java', 'java-17', 501, '8.130.14.34', '/opt/data/java', '', '', '{\"APP_PORT\":\"8080\",\"JAVA_OPTS\":\"-Xmx512m -Xms256m\"}', 1, '[2025-10-31 12:53:21] 开始部署 Java Java 17 LTS\n[2025-10-31 12:53:21] 连接主机 8.130.14.34...\n[2025-10-31 12:53:22] SSH连接成功\n[2025-10-31 12:53:22] 检查Docker环境...\n[2025-10-31 12:53:22] 使用镜像: crpi-aj3vgoxp9kzh2jx4.cn-hangzhou.personal.cr.aliyuncs.com/zhangfan_k8s/openjdk:17-jdk\n[2025-10-31 12:53:22] 拉取镜像...\n[2025-10-31 12:53:23] 镜像拉取成功\n[2025-10-31 12:53:23] 创建临时容器...\n[2025-10-31 12:53:23] 提取文件 /usr/local/openjdk-17 -> /usr/local/java17...\n[2025-10-31 12:53:29] 清理临时容器...\n[2025-10-31 12:53:29] 读取安装脚本 common/templates/02-java/versions/java-17-install.sh...\n[2025-10-31 12:53:29] 上传安装脚本...\n[2025-10-31 12:53:29] 执行安装脚本...\n[2025-10-31 12:53:30] 安装脚本输出:\n===== Java 17 LTS 安装配置 =====\n安装路径: /usr/local/java17\n环境变量文件: /etc/profile.d/java17.sh\n环境变量已配置: /etc/profile.d/java17.sh\n===== 安装成功 =====\n\n[2025-10-31 12:53:30] 验证安装...\n[2025-10-31 12:53:30] 验证结果:\nopenjdk version \"17.0.0.1\" 2024-07-02\nOpenJDK Runtime Environment (build 17.0.0.1+2-3)\nOpenJDK 64-Bit Server VM (build 17.0.0.1+2-3, mixed mode, sharing)\n\n[2025-10-31 12:53:30] 部署完成！\n', '2025-10-31 12:53:21.350', '2025-10-31 12:53:21.350');
INSERT INTO `tool_service_deploy` (`id`, `service_name`, `service_id`, `version`, `host_id`, `host_ip`, `install_dir`, `container_name`, `ports`, `env_vars`, `status`, `deploy_log`, `create_time`, `update_time`) VALUES (12, 'Elasticsearch', 'elasticsearch', 'elasticsearch-8.x', 506, '139.9.205.38', '/opt/data/elasticsearch', '', '', '{\"ES_HEAP_SIZE\":\"1g\",\"ES_HTTP_PORT\":\"9200\"}', 3, '[2025-11-30 22:33:06] 开始部署 Elasticsearch Elasticsearch 8.x\n[2025-11-30 22:33:06] 连接主机 139.9.205.38...\n[2025-11-30 22:33:07] SSH连接成功\n[2025-11-30 22:33:07] 创建安装目录 /opt/data/elasticsearch...\n[2025-11-30 22:33:07] 读取模板文件 common/templates/06-elasticsearch/versions/elasticsearch-8.x-docker-compose.yml...\n[2025-11-30 22:33:07] 生成环境变量配置...\n[2025-11-30 22:33:07] 上传 docker-compose.yml...\n[2025-11-30 22:33:07] 上传 .env...\n[2025-11-30 22:33:07] 检查Docker环境...\n[2025-11-30 22:33:08] 启动服务容器...\n[2025-11-30 22:33:08] 启动失败: Process exited with status 127\n输出: \nSTDERR:\nbash: line 1: docker-compose: command not found\n\n', '2025-11-30 22:33:06.149', '2025-11-30 22:33:06.149');
COMMIT;

-- ----------------------------
-- View structure for db_instance_all
-- ----------------------------
DROP VIEW IF EXISTS `db_instance_all`;
CREATE ALGORITHM = UNDEFINED SQL SECURITY DEFINER VIEW `db_instance_all` AS select `db_instance`.`id` AS `id`,`db_instance`.`code` AS `code`,`db_instance`.`name` AS `name`,(case when (`db_instance`.`type` = 'postgres') then 'postgres' else 'mysql' end) AS `db_type`,`db_instance`.`type` AS `sub_type`,`db_instance`.`host` AS `host`,`db_instance`.`port` AS `port`,`db_instance`.`username` AS `username`,`db_instance`.`password` AS `password`,`db_instance`.`remark` AS `remark`,`db_instance`.`ssh_tunnel_machine_id` AS `ssh_tunnel_machine_id`,`db_instance`.`status` AS `status`,`db_instance`.`create_time` AS `create_time`,`db_instance`.`update_time` AS `update_time`,`db_instance`.`creator` AS `creator`,`db_instance`.`creator_id` AS `creator_id`,`db_instance`.`modifier` AS `modifier`,`db_instance`.`modifier_id` AS `modifier_id`,json_object('network',`db_instance`.`network`,'params',`db_instance`.`params`) AS `connection_config` from `db_instance` union all select `db_redis_instance`.`id` AS `id`,`db_redis_instance`.`code` AS `code`,`db_redis_instance`.`name` AS `name`,'redis' AS `db_type`,`db_redis_instance`.`mode` AS `sub_type`,`db_redis_instance`.`host` AS `host`,`db_redis_instance`.`port` AS `port`,`db_redis_instance`.`username` AS `username`,`db_redis_instance`.`password` AS `password`,`db_redis_instance`.`remark` AS `remark`,`db_redis_instance`.`ssh_tunnel_machine_id` AS `ssh_tunnel_machine_id`,`db_redis_instance`.`status` AS `status`,`db_redis_instance`.`create_time` AS `create_time`,`db_redis_instance`.`update_time` AS `update_time`,`db_redis_instance`.`creator` AS `creator`,`db_redis_instance`.`creator_id` AS `creator_id`,`db_redis_instance`.`modifier` AS `modifier`,`db_redis_instance`.`modifier_id` AS `modifier_id`,json_object('mode',`db_redis_instance`.`mode`,'db',`db_redis_instance`.`db`) AS `connection_config` from `db_redis_instance` union all select `db_mongo_instance`.`id` AS `id`,`db_mongo_instance`.`code` AS `code`,`db_mongo_instance`.`name` AS `name`,'mongodb' AS `db_type`,`db_mongo_instance`.`type` AS `sub_type`,NULL AS `host`,NULL AS `port`,NULL AS `username`,NULL AS `password`,`db_mongo_instance`.`remark` AS `remark`,`db_mongo_instance`.`ssh_tunnel_machine_id` AS `ssh_tunnel_machine_id`,`db_mongo_instance`.`status` AS `status`,`db_mongo_instance`.`create_time` AS `create_time`,`db_mongo_instance`.`update_time` AS `update_time`,`db_mongo_instance`.`creator` AS `creator`,`db_mongo_instance`.`creator_id` AS `creator_id`,`db_mongo_instance`.`modifier` AS `modifier`,`db_mongo_instance`.`modifier_id` AS `modifier_id`,json_object('uri',`db_mongo_instance`.`uri`) AS `connection_config` from `db_mongo_instance` union all select `db_es_instance`.`id` AS `id`,`db_es_instance`.`code` AS `code`,`db_es_instance`.`name` AS `name`,'elasticsearch' AS `db_type`,`db_es_instance`.`protocol` AS `sub_type`,`db_es_instance`.`host` AS `host`,`db_es_instance`.`port` AS `port`,`db_es_instance`.`username` AS `username`,`db_es_instance`.`password` AS `password`,`db_es_instance`.`remark` AS `remark`,`db_es_instance`.`ssh_tunnel_machine_id` AS `ssh_tunnel_machine_id`,`db_es_instance`.`status` AS `status`,`db_es_instance`.`create_time` AS `create_time`,`db_es_instance`.`update_time` AS `update_time`,`db_es_instance`.`creator` AS `creator`,`db_es_instance`.`creator_id` AS `creator_id`,`db_es_instance`.`modifier` AS `modifier`,`db_es_instance`.`modifier_id` AS `modifier_id`,json_object('protocol',`db_es_instance`.`protocol`) AS `connection_config` from `db_es_instance`;

SET FOREIGN_KEY_CHECKS = 1;
