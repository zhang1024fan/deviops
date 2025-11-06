# 任务列表连表查询功能说明

## 功能简介

增强了任务列表接口，现在会自动包含以下关联信息：
1. **模板信息**：根据`shell`字段中的模板ID获取对应的模板详情
2. **主机信息**：根据`host_ids`字段中的主机ID获取对应的主机详情

## API接口

### 原接口（已升级）
**接口**: `GET /api/v1/task/list?page=1&pageSize=10`

### 新增专用接口
**接口**: `GET /api/v1/task/list-with-details?page=1&pageSize=10`

两个接口现在返回相同的增强数据格式。

## 返回数据格式

### 原始格式（修改前）
```json
{
  "code": 200,
  "data": {
    "list": [
      {
        "id": 23,
        "name": "获取日志详细信息",
        "type": 1,
        "shell": "2,3,12",
        "host_ids": "1",
        "execute_count": 0,
        "next_run_time": null
      }
    ]
  }
}
```

### 增强格式（修改后）
```json
{
  "code": 200,
  "data": {
    "list": [
      {
        "id": 23,
        "name": "获取日志详细信息",
        "type": 1,
        "shell": "2,3,12",
        "host_ids": "1",
        "cron_expr": "",
        "tasklog": "",
        "status": 3,
        "duration": 115,
        "task_count": 3,
        "execute_count": 0,
        "next_run_time": null,
        "remark": "",
        "start_time": null,
        "end_time": "2025-08-14T12:28:28.011+08:00",
        "created_at": "2025-08-14T10:41:59.846+08:00",
        "templates": [
          {
            "id": 2,
            "name": "系统信息检查",
            "type": 1,
            "content": "#!/bin/bash\necho '=== 系统信息 ==='\nuname -a"
          },
          {
            "id": 3,
            "name": "磁盘空间检查",
            "type": 1,
            "content": "#!/bin/bash\necho '=== 磁盘使用情况 ==='\ndf -h"
          },
          {
            "id": 12,
            "name": "日志详细信息",
            "type": 1,
            "content": "#!/bin/bash\necho '=== 查看系统日志 ==='\ntail -n 50 /var/log/messages"
          }
        ],
        "hosts": [
          {
            "id": 1,
            "hostName": "web-server-01",
            "privateIp": "10.0.1.10",
            "publicIp": "203.0.113.10",
            "sshIp": "10.0.1.10",
            "status": 1
          }
        ]
      },
      {
        "id": 24,
        "name": "定时任务测试001",
        "type": 2,
        "shell": "12",
        "host_ids": "1",
        "cron_expr": "*/2 * * * * ",
        "tasklog": "",
        "status": 1,
        "duration": 6,
        "task_count": 1,
        "execute_count": 0,
        "next_run_time": null,
        "remark": "定时任务测试001",
        "start_time": null,
        "end_time": "2025-08-14T15:07:54.377+08:00",
        "created_at": "2025-08-14T11:44:48.087+08:00",
        "templates": [
          {
            "id": 12,
            "name": "日志详细信息",
            "type": 1,
            "content": "#!/bin/bash\necho '=== 查看系统日志 ==='\ntail -n 50 /var/log/messages"
          }
        ],
        "hosts": [
          {
            "id": 1,
            "hostName": "web-server-01",
            "privateIp": "10.0.1.10",
            "publicIp": "203.0.113.10",
            "sshIp": "10.0.1.10",
            "status": 1
          }
        ]
      }
    ],
    "total": 2,
    "page": 1,
    "pageSize": 10
  }
}
```

## 新增字段说明

### templates 数组
包含任务使用的所有模板信息：

```json
{
  "id": 2,              // 模板ID
  "name": "系统信息检查",  // 模板名称
  "type": 1,            // 模板类型：1=shell, 2=python, 3=ansible
  "content": "#!/bin/bash..." // 模板内容
}
```

### hosts 数组
包含任务执行的所有主机信息：

```json
{
  "id": 1,                    // 主机ID
  "hostName": "web-server-01", // 主机名
  "privateIp": "10.0.1.10",   // 私网IP
  "publicIp": "203.0.113.10", // 公网IP
  "sshIp": "10.0.1.10",      // SSH连接IP
  "status": 1                 // 主机状态：1=认证成功, 2=未认证, 3=认证失败
}
```

## 数据关联逻辑

### 模板关联
- 解析`shell`字段中的逗号分隔ID：`"2,3,12"` → `[2, 3, 12]`
- 查询`task_template`表获取对应模板信息
- 按原始顺序返回模板数组

### 主机关联
- 解析`host_ids`字段中的逗号分隔ID：`"1"` → `[1]`
- 查询`cmdb_host`表获取对应主机信息
- 返回主机详情数组

## 前端使用建议

### 1. 显示模板信息
```javascript
// 显示任务使用的模板
task.templates.forEach(template => {
  console.log(`模板: ${template.name} (类型: ${getTemplateType(template.type)})`)
})

function getTemplateType(type) {
  const types = { 1: 'Shell脚本', 2: 'Python脚本', 3: 'Ansible剧本' }
  return types[type] || '未知'
}
```

### 2. 显示主机信息
```javascript
// 显示任务执行的主机
task.hosts.forEach(host => {
  console.log(`主机: ${host.hostName} (${host.sshIp})`)
})

// 主机状态显示
function getHostStatus(status) {
  const statuses = { 1: '已认证', 2: '未认证', 3: '认证失败' }
  return statuses[status] || '未知'
}
```

### 3. 任务详情展示
```javascript
const TaskDetail = ({ task }) => {
  return (
    <div>
      <h3>{task.name}</h3>
      <p>执行次数: {task.execute_count}</p>

      {/* 模板信息 */}
      <div>
        <h4>使用的模板:</h4>
        {task.templates.map(template => (
          <div key={template.id}>
            <span>{template.name}</span>
            <small>({getTemplateType(template.type)})</small>
          </div>
        ))}
      </div>

      {/* 主机信息 */}
      <div>
        <h4>执行主机:</h4>
        {task.hosts.map(host => (
          <div key={host.id}>
            <span>{host.hostName}</span>
            <span>({host.sshIp})</span>
            <span className={`status-${host.status}`}>
              {getHostStatus(host.status)}
            </span>
          </div>
        ))}
      </div>
    </div>
  )
}
```

## 性能说明

### 查询优化
- 使用`IN`查询批量获取关联数据，避免N+1查询问题
- 每个任务最多执行2个额外查询（模板+主机）
- 空字段会跳过查询，避免不必要的数据库访问

### 返回数据大小
- 模板内容可能较大，建议前端根据需要显示
- 如果不需要模板内容，可以在service层移除`Content`字段

## 兼容性

- ✅ 完全向后兼容原有接口
- ✅ 新增字段不影响现有前端代码
- ✅ 查询参数保持不变
- ✅ 分页逻辑保持不变

## 错误处理

如果关联查询失败：
- 模板查询失败：`templates`字段为空数组
- 主机查询失败：`hosts`字段为空数组
- 不会影响主任务列表的返回

这样前端就可以获得更完整的任务信息，无需额外的API调用来获取模板和主机详情。