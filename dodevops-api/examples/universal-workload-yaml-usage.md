# 通用工作负载YAML管理功能使用说明

## 功能概述

新增的通用工作负载YAML管理功能支持所有主要工作负载类型的YAML操作，根据工作负载类型自动判断和处理，提供统一的API接口。

## 支持的工作负载类型

- ✅ **Deployment** - 支持获取和更新YAML
- ✅ **StatefulSet** - 支持获取和更新YAML
- ✅ **DaemonSet** - 支持获取和更新YAML
- ✅ **Job** - 支持获取YAML（不支持更新，需要删除重建）
- ✅ **CronJob** - 支持获取和更新YAML

## API接口

### 1. 获取工作负载YAML
```
GET /k8s/cluster/{id}/namespaces/{namespaceName}/workload-yaml/{workloadType}/{workloadName}
```

### 2. 更新工作负载YAML
```
PUT /k8s/cluster/{id}/namespaces/{namespaceName}/workload-yaml
```

## 数据结构

### UpdateWorkloadYAMLRequest
```json
{
  "workloadType": "deployment",     // 必需：工作负载类型
  "workloadName": "nginx-app",      // 必需：工作负载名称
  "yamlContent": "apiVersion: apps/v1\nkind: Deployment\n...",  // 必需：YAML内容
  "dryRun": false,                  // 可选：是否只进行校验不实际更新
  "validateOnly": false,            // 可选：是否只校验YAML格式
  "force": false                    // 可选：是否强制更新
}
```

### UpdateWorkloadYAMLResponse
```json
{
  "success": true,
  "workloadType": "deployment",
  "workloadName": "nginx-app",
  "namespace": "default",
  "message": "Deployment更新成功",
  "updateStrategy": "rolling",      // patch/update/rolling
  "validationResult": {            // 校验结果（validateOnly时返回）
    "valid": true,
    "errors": [],
    "warnings": [],
    "suggestions": []
  },
  "changes": [                     // 变更说明
    "Deployment配置已更新，正在执行滚动更新"
  ],
  "warnings": [],                  // 警告信息
  "appliedAt": "2024-01-15 10:30:00"
}
```

## 使用示例

### 1. 获取Deployment的YAML配置
```bash
curl -X GET "http://localhost:8080/k8s/cluster/1/namespaces/default/workload-yaml/deployment/nginx-app" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**响应：**
```json
{
  "code": 200,
  "data": {
    "success": true,
    "workloadType": "deployment",
    "workloadName": "nginx-app",
    "namespace": "default",
    "yamlContent": "apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: nginx-app\n...",
    "message": "成功获取deployment 'nginx-app'的YAML配置"
  }
}
```

### 2. 校验YAML格式
```bash
curl -X PUT "http://localhost:8080/k8s/cluster/1/namespaces/default/workload-yaml" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "workloadType": "deployment",
    "workloadName": "nginx-app",
    "yamlContent": "apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: nginx-app\n  labels:\n    app: nginx\nspec:\n  replicas: 3\n  selector:\n    matchLabels:\n      app: nginx\n  template:\n    metadata:\n      labels:\n        app: nginx\n    spec:\n      containers:\n      - name: nginx\n        image: nginx:1.20",
    "validateOnly": true
  }'
```

### 3. DryRun模式预览Deployment变更
```bash
curl -X PUT "http://localhost:8080/k8s/cluster/1/namespaces/default/workload-yaml" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "workloadType": "deployment",
    "workloadName": "nginx-app",
    "yamlContent": "apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: nginx-app\nspec:\n  replicas: 5\n  selector:\n    matchLabels:\n      app: nginx\n  template:\n    metadata:\n      labels:\n        app: nginx\n    spec:\n      containers:\n      - name: nginx\n        image: nginx:1.21",
    "dryRun": true
  }'
```

### 4. 更新Deployment镜像版本
```bash
curl -X PUT "http://localhost:8080/k8s/cluster/1/namespaces/default/workload-yaml" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "workloadType": "deployment",
    "workloadName": "nginx-app",
    "yamlContent": "apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: nginx-app\nspec:\n  replicas: 3\n  selector:\n    matchLabels:\n      app: nginx\n  template:\n    metadata:\n      labels:\n        app: nginx\n    spec:\n      containers:\n      - name: nginx\n        image: nginx:1.21"
  }'
```

### 5. 更新StatefulSet配置
```bash
curl -X PUT "http://localhost:8080/k8s/cluster/1/namespaces/default/workload-yaml" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "workloadType": "statefulset",
    "workloadName": "web-app",
    "yamlContent": "apiVersion: apps/v1\nkind: StatefulSet\nmetadata:\n  name: web-app\nspec:\n  serviceName: web-service\n  replicas: 3\n  selector:\n    matchLabels:\n      app: web\n  template:\n    metadata:\n      labels:\n        app: web\n    spec:\n      containers:\n      - name: web\n        image: nginx:1.21"
  }'
```

### 6. 更新CronJob配置
```bash
curl -X PUT "http://localhost:8080/k8s/cluster/1/namespaces/default/workload-yaml" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "workloadType": "cronjob",
    "workloadName": "backup-job",
    "yamlContent": "apiVersion: batch/v1\nkind: CronJob\nmetadata:\n  name: backup-job\nspec:\n  schedule: \"0 2 * * *\"\n  jobTemplate:\n    spec:\n      template:\n        spec:\n          containers:\n          - name: backup\n            image: backup-tool:latest\n            command: [\"/bin/backup\"]\n          restartPolicy: OnFailure"
  }'
```

## 更新策略说明

### Rolling Update（滚动更新）
- **适用于：** Deployment、StatefulSet、DaemonSet
- **特点：** 逐步替换Pod实例，无服务中断
- **场景：** 镜像版本更新、环境变量修改、资源配置调整

### Direct Update（直接更新）
- **适用于：** CronJob
- **特点：** 直接更新配置，下次执行时生效
- **场景：** 调度时间修改、Job模板更新

### Recreate（重建）
- **适用于：** Job（不支持更新，需要删除重建）
- **特点：** 删除旧资源，创建新资源
- **场景：** Job配置变更

## 各工作负载类型特点

### Deployment
- ✅ 支持滚动更新
- ✅ 支持副本数调整
- ✅ 支持镜像版本更新
- ✅ 变更分析详细

### StatefulSet
- ✅ 支持滚动更新
- ✅ 支持有序更新
- ⚠️ 某些字段修改需要重启

### DaemonSet
- ✅ 支持滚动更新
- ✅ 自动在所有节点部署

### Job
- ❌ 不支持更新
- 💡 建议：删除后重新创建

### CronJob
- ✅ 支持直接更新
- ✅ 下次执行时生效

## 错误处理

### 常见错误
- `400`：工作负载类型不支持、YAML格式错误、名称不匹配
- `404`：集群不存在、工作负载不存在
- `500`：Kubernetes API调用失败

### 错误示例
```json
{
  "code": 400,
  "message": "不支持的工作负载类型: invalid. 支持的类型: deployment,statefulset,daemonset,job,cronjob"
}
```

```json
{
  "code": 400,
  "message": "YAML中的Deployment名称(wrong-name)与请求参数不匹配(nginx-app)"
}
```

## 最佳实践

1. **类型检查**：确保使用正确的工作负载类型
2. **名称匹配**：YAML中的名称必须与请求参数匹配
3. **先校验**：使用 `validateOnly: true` 检查YAML格式
4. **预览变更**：使用 `dryRun: true` 了解变更影响
5. **分步更新**：对于关键应用，建议分步进行更新
6. **监控更新**：更新后监控应用状态和日志
7. **备份配置**：重要变更前备份当前配置

## 注意事项

1. **命名空间匹配**：YAML中的命名空间必须与URL参数匹配
2. **Job限制**：Job不支持更新，只能删除重建
3. **资源版本**：系统会自动处理ResourceVersion等字段
4. **权限要求**：需要对应工作负载的更新权限
5. **滚动更新**：某些变更会触发滚动更新，可能需要时间完成