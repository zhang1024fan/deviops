-- MySQL 8.0 认证方式配置
-- 修改 root 用户认证插件为 caching_sha2_password

USE mysql;

-- 修改本地连接认证方式
ALTER USER 'root'@'localhost' IDENTIFIED WITH caching_sha2_password BY 'devops@2025';

-- 修改容器内部连接认证方式
ALTER USER 'root'@'%' IDENTIFIED WITH caching_sha2_password BY 'devops@2025';

-- 刷新权限
FLUSH PRIVILEGES;

-- 验证认证插件
SELECT user, host, plugin FROM mysql.user WHERE user = 'root';
