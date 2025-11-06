# 定时任务修复说明

## 修复内容

本次修复解决了以下问题：

### 1. 任务重复注册
**问题**：同一个定时任务被注册两次，导致每次触发时执行两个回调
**原因**：启动时 `LoadScheduledTasks` 注册一次，用户点击启动按钮又注册一次
**修复**：在 `ScheduleJob` 中添加重复注册检查

### 2. 回调函数不统一
**问题**：`ResumeScheduledTask` 中的回调调用 `StartJob`，导致"任务正在运行中"错误
**原因**：不同地方使用了不同的回调函数实现
**修复**：创建统一的 `createScheduledTaskHandler` 函数，所有地方都使用它

### 3. 子任务状态卡住
**问题**：子任务状态变为运行中(2)后未被重置为等待中(1)
**原因**：任务执行完成后状态重置逻辑存在，但任务可能未被消费
**修复**：统一回调函数，确保状态管理一致

### 4. 死锁问题
**问题**：暂停/恢复任务时API超时
**原因**：在持有互斥锁的情况下调用阻塞操作
**修复**：使用细粒度锁，分步骤执行，避免长时间持锁

## 使用步骤

### 1. 停止旧进程
```bash
# 查找进程
ps aux | grep dodevops-api | grep -v grep

# 停止进程（替换 PID）
kill -9 <PID>
```

### 2. 重置数据库状态
```bash
mysql -u your_user -p your_database < scripts/reset_scheduled_task_status.sql
```

或者手动执行：
```sql
-- 重置任务53和54的子任务状态
UPDATE task_work SET status = 1 WHERE task_id IN (53, 54) AND status = 2;

-- 查看状态
SELECT id, task_id, status FROM task_work WHERE task_id IN (53, 54);
```

### 3. 启动新版本
```bash
./dodevops-api

# 或使用 go run
go run main.go
```

### 4. 观察日志

**正常情况下应该看到**：

启动时加载定时任务：
```
2025/09/30 16:30:00 开始加载定时任务，共 2 个
2025/09/30 16:30:00 处理定时任务: TaskID=53, Type=2, Status=2, CronExpr=*/3 * * * *
2025/09/30 16:30:00 恢复运行中的定时任务: TaskID=53
2025/09/30 16:30:00 定时任务已成功添加: TaskID=53, CronExpr=*/3 * * * *, EntryID=1
```

定时任务触发：
```
2025/09/30 16:33:00 执行定时任务: TaskID=53, CronExpr=*/3 * * * *
2025/09/30 16:33:00 === 定时任务触发 (LoadScheduledTasks) ===
2025/09/30 16:33:00 获取到 2 个子任务
2025/09/30 16:33:00 准备执行子任务: ID=98, TemplateID=11, 当前状态=1
2025/09/30 16:33:00 定时任务已提交到队列: TaskID=53, JobID=98, Priority=normal
2025/09/30 16:33:00 准备执行子任务: ID=99, TemplateID=12, 当前状态=1
2025/09/30 16:33:00 定时任务已提交到队列: TaskID=53, JobID=99, Priority=normal
2025/09/30 16:33:00 共提交 2 个子任务到队列
```

任务执行完成：
```
定时任务父任务状态保持为运行中(2)
子任务98状态已重置为等待中(1)
子任务99状态已重置为等待中(1)
```

### 5. 测试功能

#### 测试启动任务
如果任务已在调度器中，应该看到：
```
任务已在调度器中，无需重复注册: TaskID=53
```

#### 测试暂停任务
```bash
POST /api/v1/task/monitor/scheduled/pause?task_id=53
```
应该立即返回，不会超时

#### 测试恢复任务
```bash
POST /api/v1/task/monitor/scheduled/resume?task_id=53
```
应该快速完成

## 关键改进

### 1. 统一的回调函数
```go
func createScheduledTaskHandler(taskID uint) func() {
    // 所有定时任务都使用这个统一的回调函数
    // 直接获取子任务并提交到队列
    // 不再调用 StartJob
}
```

### 2. 防止重复注册
```go
// 在 ScheduleJob 中检查
entries := scheduler.GetEntries()
if _, exists := entries[job.TaskID]; exists {
    fmt.Printf("任务已在调度器中，无需重复注册: TaskID=%d\n", job.TaskID)
    return nil
}
```

### 3. 细粒度锁
```go
// 分步骤执行，避免长时间持锁
Lock → 获取数据 → Unlock
Remove/Database (不持锁)
Lock → 更新状态 → Unlock
```

## 预期效果

✅ 定时任务不会重复触发
✅ 子任务能正常提交到队列并执行
✅ 执行完成后状态自动重置为等待中
✅ 暂停/恢复/停止功能快速响应
✅ 无死锁问题

## 排查问题

如果任务仍然无法执行：

1. **检查调度器状态**
```bash
GET /api/v1/task/monitor/scheduler/stats
```

2. **检查队列状态**
```bash
GET /api/v1/task/monitor/queue/details
```

3. **检查任务状态**
```bash
GET /api/v1/task/monitor/task/status?task_id=53
```

4. **查看数据库状态**
```sql
SELECT * FROM task_work WHERE task_id = 53;
```

5. **检查工作者是否运行**
日志中应该有：
```
工作者 0 启动
工作者 1 启动
...
```