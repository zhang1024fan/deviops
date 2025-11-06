# 定时同步云主机功能使用指南

## 功能简介

定时同步云主机功能允许用户配置定时任务，自动同步多个云厂商的主机信息。支持的云厂商包括：
- 1: 阿里云
- 2: 腾讯云
- 3: 百度云
- 4: 华为云
- 5: AWS云

## API接口

### 1. 创建定时同步配置

**接口**: `POST /api/v1/config/sync-schedule`

**请求体**:
```json
{
  "name": "每日云主机同步",
  "cronExpr": "0 0 2 * * ?",
  "keyTypes": "[1,2]",
  "status": 1,
  "remark": "每天凌晨2点同步阿里云和腾讯云主机"
}
```

**字段说明**:
- `name`: 配置名称
- `cronExpr`: cron表达式（支持秒级）
- `keyTypes`: 要同步的云厂商类型JSON数组，如 "[1,2,3]"
- `status`: 状态，1=启用，0=禁用
- `remark`: 备注信息

**常用cron表达式示例**:
- `0 0 2 * * ?`: 每天凌晨2点执行
- `0 0 */6 * * ?`: 每6小时执行一次
- `0 30 1 * * ?`: 每天凌晨1:30执行
- `0 0 9,18 * * MON-FRI`: 工作日上午9点和下午6点执行

### 2. 查询定时同步配置列表

**接口**: `GET /api/v1/config/sync-schedule/list?page=1&pageSize=10`

### 3. 更新定时同步配置

**接口**: `PUT /api/v1/config/sync-schedule`

**请求体**:
```json
{
  "id": 1,
  "name": "更新后的配置名称",
  "cronExpr": "0 0 3 * * ?",
  "keyTypes": "[1,2,3]",
  "status": 1,
  "remark": "更新后的备注"
}
```

### 4. 删除定时同步配置

**接口**: `DELETE /api/v1/config/sync-schedule?id=1`

### 5. 切换配置状态

**接口**: `POST /api/v1/config/sync-schedule/toggle-status`

**请求体**:
```json
{
  "id": 1,
  "status": 0
}
```

### 6. 获取调度器状态

**接口**: `GET /api/v1/config/sync-schedule/scheduler-stats`

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "status": "running",
    "total_jobs": 3,
    "cron_entries": 3,
    "active_jobs": {
      "1": 1,
      "2": 2,
      "3": 3
    }
  }
}
```

### 7. 手动触发同步（测试用）

**接口**: `POST /api/v1/config/sync-schedule/trigger?id=1`

## 使用流程

### 1. 配置云厂商密钥

在使用定时同步功能之前，需要先配置云厂商的访问密钥：

```bash
# 创建阿里云密钥
curl -X POST http://localhost:8080/api/v1/config/keymanage \
  -H "Content-Type: application/json" \
  -d '{
    "keyType": 1,
    "keyId": "your_access_key_id",
    "keySecret": "your_access_key_secret",
    "remark": "阿里云密钥"
  }'

# 创建腾讯云密钥
curl -X POST http://localhost:8080/api/v1/config/keymanage \
  -H "Content-Type: application/json" \
  -d '{
    "keyType": 2,
    "keyId": "your_secret_id",
    "keySecret": "your_secret_key",
    "remark": "腾讯云密钥"
  }'
```

### 2. 创建定时同步配置

```bash
curl -X POST http://localhost:8080/api/v1/config/sync-schedule \
  -H "Content-Type: application/json" \
  -d '{
    "name": "每日全量同步",
    "cronExpr": "0 0 2 * * ?",
    "keyTypes": "[1,2]",
    "status": 1,
    "remark": "每天凌晨2点同步所有云厂商主机"
  }'
```

### 3. 验证配置生效

```bash
# 查看调度器状态
curl http://localhost:8080/api/v1/config/sync-schedule/scheduler-stats

# 查看启用的配置
curl http://localhost:8080/api/v1/config/sync-schedule/active
```

## 注意事项

1. **cron表达式格式**: 使用6位cron表达式（秒 分 时 日 月 周）
2. **云厂商类型**: keyTypes字段必须是有效的JSON数组格式
3. **权限要求**: 云厂商密钥需要有读取ECS实例的权限
4. **同步策略**: 系统会自动遍历指定云厂商类型的所有密钥进行同步
5. **错误处理**: 单个云厂商同步失败不会影响其他云厂商的同步
6. **日志记录**: 所有同步操作都会记录详细日志

## 故障排查

### 1. 查看调度器状态
```bash
curl http://localhost:8080/api/v1/config/sync-schedule/scheduler-stats
```

### 2. 检查配置是否正确
- 验证cron表达式格式
- 确认keyTypes是有效的JSON数组
- 检查云厂商密钥是否配置正确

### 3. 手动测试同步
```bash
curl -X POST http://localhost:8080/api/v1/config/sync-schedule/trigger?id=1
```

### 4. 查看系统日志
定时同步的执行日志会输出到应用程序日志中，包含：
- 任务开始和结束时间
- 同步的云厂商类型和结果
- 错误信息（如果有）

## 最佳实践

1. **合理设置同步频率**: 避免过于频繁的同步，建议每天1-2次
2. **分时段同步**: 如果云厂商较多，可以分不同时段同步以避免并发过高
3. **监控同步结果**: 定期检查调度器状态和同步日志
4. **测试验证**: 新增配置后先手动触发测试，确认无误后再启用定时