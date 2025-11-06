# 定时同步状态切换修复说明

## 问题描述

在调用定时同步配置的状态切换接口时遇到验证错误：

**请求**:
```bash
POST /api/v1/config/sync-schedule/toggle-status
{
  "id": 3,
  "status": 0
}
```

**错误响应**:
```json
{
    "code": 400,
    "message": "参数错误: Key: 'status' Error:Field validation for 'status' failed on the 'required' tag",
    "data": {}
}
```

## 问题原因

在Controller中，`ToggleStatus` 方法的请求体验证使用了 `binding:"required"`：

```go
var req struct {
    ID     uint `json:"id" binding:"required"`
    Status int  `json:"status" binding:"required"`  // 问题出现在这里
}
```

当传递 `status: 0`（禁用状态）时，Go的validator库认为这是零值，触发了required验证失败。

## 解决方案

修改验证规则，使用范围验证而不是required验证：

```go
var req struct {
    ID     uint `json:"id" binding:"required"`
    Status int  `json:"status" binding:"min=0,max=1"`  // 修复后
}
```

这样：
- ✅ `status: 0` (禁用) - 验证通过
- ✅ `status: 1` (启用) - 验证通过
- ❌ `status: 2` (无效值) - 验证失败
- ❌ 缺少status字段 - 验证失败（因为默认值0在范围内）

## 测试修复

### 1. 重新编译
```bash
go build -o dodevops-api .
```

### 2. 重启应用程序

### 3. 测试禁用功能
```bash
curl -X POST http://10.7.16.22:8080/api/v1/config/sync-schedule/toggle-status \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "id": 3,
    "status": 0
  }'
```

**预期响应**:
```json
{
  "code": 200,
  "message": "成功",
  "data": {}
}
```

### 4. 测试启用功能
```bash
curl -X POST http://10.7.16.22:8080/api/v1/config/sync-schedule/toggle-status \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "id": 3,
    "status": 1
  }'
```

### 5. 验证无效值处理
```bash
curl -X POST http://10.7.16.22:8080/api/v1/config/sync-schedule/toggle-status \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "id": 3,
    "status": 2
  }'
```

**预期响应**:
```json
{
  "code": 400,
  "message": "参数错误: Key: 'Status' Error:Field validation for 'Status' failed on the 'max' tag",
  "data": {}
}
```

## 状态说明

- `status: 1` - 启用定时同步
  - 配置会被添加到调度器
  - 按照cron表达式定时执行

- `status: 0` - 禁用定时同步
  - 配置会从调度器中移除
  - 不再执行定时同步

## 相关功能

切换状态后，可以通过以下接口验证：

1. **查看配置列表**:
   ```bash
   GET /api/v1/config/sync-schedule/list
   ```

2. **查看调度器状态**:
   ```bash
   GET /api/v1/config/sync-schedule/scheduler-stats
   ```

3. **查看启用的配置**:
   ```bash
   GET /api/v1/config/sync-schedule/active
   ```

修复后，状态切换功能应该能正常工作，支持启用和禁用定时同步配置。