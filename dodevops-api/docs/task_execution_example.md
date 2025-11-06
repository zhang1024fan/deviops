# 任务执行次数和下次执行时间功能使用指南

## 功能简介

在原有的任务管理基础上，新增了以下功能：
1. **执行次数统计**：记录任务执行的次数
2. **下次执行时间**：对于定时任务，记录和计算下次执行时间
3. **自动更新机制**：每次任务执行完成后自动更新执行次数和下次执行时间

## 数据模型变更

### Task 模型新增字段

```go
type Task struct {
    // ... 原有字段
    ExecuteCount int        `json:"execute_count" gorm:"default:0;comment:执行次数"`
    NextRunTime  *time.Time `json:"next_run_time" gorm:"comment:下次执行时间"`
    // ... 其他字段
}
```

**字段说明**：
- `execute_count`: 执行次数，每次任务执行完成后自动+1
- `next_run_time`: 下次执行时间，定时任务专用，根据cron表达式自动计算

## API接口

### 1. 创建任务（已更新）

**接口**: `POST /api/v1/task/add`

**请求体**:
```json
{
  "name": "定时清理任务",
  "type": 2,
  "shell": "echo 'cleaning up...'",
  "host_ids": "1,2,3",
  "cron_expr": "0 */5 * * * *",
  "remark": "每5分钟执行一次清理任务"
}
```

**响应**:
```json
{
  "code": 200,
  "data": {
    "id": 1,
    "name": "定时清理任务",
    "type": 2,
    "execute_count": 0,
    "next_run_time": "2024-01-01T16:05:00Z",
    "cron_expr": "0 */5 * * * *",
    "status": 1
  }
}
```

### 2. 获取任务执行信息

**接口**: `GET /api/v1/task/execution-info?id=1`

**响应**:
```json
{
  "code": 200,
  "data": {
    "id": 1,
    "name": "定时清理任务",
    "type": 2,
    "execute_count": 15,
    "next_run_time": "2024-01-01T16:10:00Z",
    "calculated_next_run_time": "2024-01-01T16:10:00Z",
    "cron_expr": "0 */5 * * * *",
    "status": 3,
    "last_end_time": "2024-01-01T16:05:30Z"
  }
}
```

**字段说明**：
- `execute_count`: 当前执行次数
- `next_run_time`: 数据库中记录的下次执行时间
- `calculated_next_run_time`: 根据cron表达式实时计算的下次执行时间
- `last_end_time`: 上次执行结束时间

## 使用场景

### 1. 普通任务（type=1）

**特点**：
- `execute_count`: 通常为0或1
- `next_run_time`: 始终为null
- 手动执行后执行次数+1

**示例**：
```json
{
  "name": "一次性备份任务",
  "type": 1,
  "shell": "backup_script.sh",
  "host_ids": "1",
  "remark": "手动执行的备份任务"
}
```

### 2. 定时任务（type=2）

**特点**：
- `execute_count`: 随每次执行递增
- `next_run_time`: 根据cron表达式自动更新
- 到达指定时间自动执行

**示例**：
```json
{
  "name": "日志清理任务",
  "type": 2,
  "shell": "clean_logs.sh",
  "host_ids": "1,2,3",
  "cron_expr": "0 0 2 * * *",
  "remark": "每天凌晨2点清理日志"
}
```

## 执行时间计算逻辑

### cron表达式示例

| 表达式 | 含义 | 执行时间示例 |
|--------|------|-------------|
| `0 */5 * * * *` | 每5分钟 | 16:00, 16:05, 16:10, 16:15... |
| `0 0 */6 * * *` | 每6小时 | 00:00, 06:00, 12:00, 18:00 |
| `0 30 9 * * MON-FRI` | 工作日9:30 | 周一到周五的9:30 |
| `0 0 2 * * *` | 每天凌晨2点 | 每天02:00 |

### 时间更新机制

1. **创建时**：如果是定时任务，计算并设置首次执行时间
2. **执行完成后**：
   - 执行次数 +1
   - 根据cron表达式计算下次执行时间
   - 更新数据库记录

## 辅助函数说明

### 1. IncrementExecuteCount

```go
// 简单增加执行次数
err := controller.IncrementExecuteCount(taskID)
```

### 2. UpdateTaskAfterExecution

```go
// 任务执行完成后的完整更新
err := controller.UpdateTaskAfterExecution(
    taskID,
    3,                    // status: 成功
    120,                  // duration: 120秒
    time.Now(),          // endTime: 结束时间
    "执行成功，清理了100个文件" // tasklog: 执行日志
)
```

## 监控和统计

### 查看任务执行统计

```bash
# 获取特定任务的执行信息
curl "http://localhost:8080/api/v1/task/execution-info?id=1"

# 获取任务列表（包含执行次数）
curl "http://localhost:8080/api/v1/task/list?page=1&pageSize=10"
```

### 定时任务执行频率分析

通过执行次数和创建时间可以分析：
- 任务的执行频率是否符合预期
- 任务是否有执行失败的情况
- 定时任务的准时性

## 最佳实践

### 1. 监控建议

- 定期检查定时任务的执行次数是否正常
- 关注执行失败的任务（status=4）
- 监控任务执行时间，避免超时

### 2. 调试建议

- 使用执行信息API查看任务状态
- 对比 `next_run_time` 和 `calculated_next_run_time` 确认时间计算正确性
- 检查执行次数是否符合预期频率

### 3. 性能优化

- 合理设置定时任务频率，避免过于频繁
- 定期清理长时间不用的任务
- 监控任务执行耗时，优化长耗时任务

## 故障排查

### 1. 执行次数不增加

检查点：
- 任务是否真正执行
- 是否调用了更新函数
- 数据库更新是否成功

### 2. 下次执行时间不正确

检查点：
- cron表达式格式是否正确
- 时区设置是否正确
- 系统时间是否准确

### 3. 定时任务不执行

检查点：
- 任务状态是否为启用
- cron调度器是否正常运行
- 下次执行时间是否已过期