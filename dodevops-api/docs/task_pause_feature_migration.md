# 定时任务暂停功能 - 数据库迁移说明

## 功能概述

本次更新为定时任务系统添加了暂停/恢复功能，包括：

1. **新增暂停状态**：任务状态新增 `5=已暂停`
2. **暂停/恢复功能**：支持暂停运行中的定时任务，恢复已暂停的任务
3. **服务重启恢复**：服务重启后自动恢复运行中的定时任务
4. **状态管理优化**：完善任务状态转换逻辑

## 数据库变更

### 1. 任务状态更新

**原状态定义：**
```sql
COMMENT '任务状态 1=等待中,2=运行中,3=成功,4=异常'
```

**新状态定义：**
```sql
COMMENT '任务状态 1=等待中,2=运行中,3=成功,4=异常,5=已暂停'
```

### 2. 状态常量定义

```go
const (
    TaskStatusPending   = 1 // 等待中
    TaskStatusRunning   = 2 // 运行中
    TaskStatusSuccess   = 3 // 成功
    TaskStatusFailed    = 4 // 异常
    TaskStatusPaused    = 5 // 已暂停
)

const (
    TaskTypeImmediate = 1 // 普通任务（立即执行）
    TaskTypeScheduled = 2 // 定时任务
    TaskTypeAnsible   = 3 // Ansible任务
)
```

## 新增API接口

### 1. 暂停定时任务
```http
POST /api/v1/task/monitor/scheduled/pause?task_id=123
```

### 2. 恢复定时任务
```http
POST /api/v1/task/monitor/scheduled/resume?task_id=123
```

### 3. 获取任务状态详情
```http
GET /api/v1/task/monitor/task/status?task_id=123
```

## 功能特性

### 1. 暂停/恢复规则

- **暂停条件**：只有运行中的定时任务可以暂停
- **恢复条件**：只有已暂停的定时任务可以恢复
- **状态检查**：提供任务状态检查和可操作性判断

### 2. 服务重启恢复

服务重启时会自动处理以下情况：

- **运行中任务**：重新计算下次执行时间并恢复调度
- **等待中任务**：启动任务并添加到调度器
- **暂停任务**：保持暂停状态，不添加到调度器
- **失败任务**：跳过处理

### 3. 执行时状态检查

任务执行时会检查：

- **父任务状态**：停止或暂停的任务不会执行
- **任务类型**：只有定时任务支持暂停
- **状态转换**：确保状态转换的合理性

## 使用示例

### 1. 暂停正在运行的定时任务

```bash
curl -X POST "http://localhost:8080/api/v1/task/monitor/scheduled/pause?task_id=123" \
  -H "Authorization: Bearer <token>"
```

**响应：**
```json
{
  "code": 200,
  "message": "请求成功",
  "data": {
    "task_id": 123,
    "message": "定时任务已暂停",
    "status": "paused"
  }
}
```

### 2. 恢复已暂停的定时任务

```bash
curl -X POST "http://localhost:8080/api/v1/task/monitor/scheduled/resume?task_id=123" \
  -H "Authorization: Bearer <token>"
```

**响应：**
```json
{
  "code": 200,
  "message": "请求成功",
  "data": {
    "task_id": 123,
    "message": "定时任务已恢复",
    "status": "running"
  }
}
```

### 3. 查看任务状态和可用操作

```bash
curl -X GET "http://localhost:8080/api/v1/task/monitor/task/status?task_id=123" \
  -H "Authorization: Bearer <token>"
```

**响应：**
```json
{
  "code": 200,
  "message": "请求成功",
  "data": {
    "task_id": 123,
    "name": "数据备份任务",
    "type": 2,
    "type_name": "定时任务",
    "status": 5,
    "status_name": "已暂停",
    "cron_expr": "0 2 * * *",
    "next_run_time": "2024-01-01T02:00:00Z",
    "execute_count": 10,
    "created_at": "2024-01-01T00:00:00Z",
    "scheduler_status": "not_scheduled",
    "operations": {
      "can_pause": false,
      "can_resume": true,
      "can_stop": false
    }
  }
}
```

## 注意事项

### 1. 兼容性

- **向后兼容**：现有API接口保持不变
- **状态迁移**：现有任务状态无需手动迁移
- **服务升级**：支持平滑升级，无需停机

### 2. 限制条件

- **仅限定时任务**：暂停功能仅适用于定时任务（Type=2）
- **状态限制**：只有特定状态的任务可以执行暂停/恢复操作
- **权限要求**：所有操作需要认证授权

### 3. 性能影响

- **最小化影响**：暂停/恢复操作对系统性能影响极小
- **状态检查**：执行前的状态检查开销可忽略
- **内存占用**：无额外内存占用

## 故障排除

### 1. 常见错误

**暂停失败：**
- 检查任务是否为定时任务且处于运行状态
- 确认全局调度器正常运行

**恢复失败：**
- 检查任务是否处于暂停状态
- 确认cron表达式有效

**服务重启后任务未恢复：**
- 检查数据库连接
- 查看服务启动日志

### 2. 调试命令

```bash
# 查看系统状态
curl -X GET "http://localhost:8080/api/v1/task/monitor/system/status"

# 查看调度器统计
curl -X GET "http://localhost:8080/api/v1/task/monitor/scheduler/stats"

# 查看队列指标
curl -X GET "http://localhost:8080/api/v1/task/monitor/queue/metrics"
```

## 版本信息

- **功能版本**：v1.2.0
- **兼容版本**：v1.0.0+
- **更新日期**：2024-01-01
- **作者**：系统开发团队