# Pod YAML更新功能使用说明

## 功能概述

新增的Pod YAML更新功能支持通过YAML内容直接更新Pod配置，提供了以下特性：

- ✅ YAML格式校验
- ✅ DryRun模式预览变更
- ✅ 智能更新策略（原地更新 vs 重建）
- ✅ 强制更新支持
- ✅ 详细的变更分析

## API接口

### 更新Pod YAML
```
PUT /k8s/cluster/{id}/namespaces/{namespaceName}/pods/{podName}/yaml
```

## 请求参数

### UpdatePodYAMLRequest
```json
{
  "yamlContent": "string",    // 必需：YAML内容
  "dryRun": false,           // 可选：是否只进行校验不实际更新
  "validateOnly": false,     // 可选：是否只校验YAML格式
  "force": false            // 可选：是否强制更新（删除重建）
}
```

## 响应结果

### UpdatePodYAMLResponse
```json
{
  "success": true,
  "podName": "example-pod",
  "namespace": "default",
  "message": "Pod原地更新成功",
  "updateStrategy": "patch",          // patch/recreate
  "validationResult": {              // 校验结果（validateOnly时返回）
    "valid": true,
    "errors": [],
    "warnings": [],
    "suggestions": []
  },
  "changes": [                       // 变更说明
    "Labels发生变更",
    "Annotations发生变更"
  ],
  "warnings": []                     // 警告信息
}
```

## 使用示例

### 1. 校验YAML格式
```bash
curl -X PUT "http://localhost:8080/k8s/cluster/1/namespaces/default/pods/example-pod/yaml" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "yamlContent": "apiVersion: v1\nkind: Pod\nmetadata:\n  name: example-pod\n  labels:\n    app: test\nspec:\n  containers:\n  - name: nginx\n    image: nginx:1.20",
    "validateOnly": true
  }'
```

### 2. DryRun模式预览变更
```bash
curl -X PUT "http://localhost:8080/k8s/cluster/1/namespaces/default/pods/example-pod/yaml" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "yamlContent": "apiVersion: v1\nkind: Pod\nmetadata:\n  name: example-pod\n  labels:\n    app: test\n    version: v2\nspec:\n  containers:\n  - name: nginx\n    image: nginx:1.21",
    "dryRun": true
  }'
```

### 3. 原地更新（只更新Labels/Annotations）
```bash
curl -X PUT "http://localhost:8080/k8s/cluster/1/namespaces/default/pods/example-pod/yaml" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "yamlContent": "apiVersion: v1\nkind: Pod\nmetadata:\n  name: example-pod\n  labels:\n    app: test\n    environment: production\nspec:\n  containers:\n  - name: nginx\n    image: nginx:1.20"
  }'
```

### 4. 强制重建更新（更新镜像版本）
```bash
curl -X PUT "http://localhost:8080/k8s/cluster/1/namespaces/default/pods/example-pod/yaml" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "yamlContent": "apiVersion: v1\nkind: Pod\nmetadata:\n  name: example-pod\n  labels:\n    app: test\nspec:\n  containers:\n  - name: nginx\n    image: nginx:1.21",
    "force": true
  }'
```

## 更新策略说明

### 原地更新（patch）
- 适用于：Labels、Annotations的变更
- 优点：无服务中断，速度快
- 限制：只能更新可变字段

### 删除重建（recreate）
- 适用于：容器镜像、端口、环境变量、资源配置、卷配置等变更
- 缺点：会导致服务中断，可能数据丢失
- 需要：设置 `force: true` 确认操作

## 注意事项

1. **YAML验证**：提交的YAML必须是有效的Pod配置
2. **名称匹配**：YAML中的Pod名称必须与URL参数匹配
3. **命名空间匹配**：YAML中的命名空间必须与URL参数匹配（可为空）
4. **强制更新**：重建操作需要显式设置 `force: true`
5. **数据备份**：重建前请确保重要数据已备份

## 错误处理

### 常见错误码
- `400`：YAML格式错误、参数验证失败
- `404`：Pod不存在、集群不存在
- `500`：Kubernetes API调用失败

### 错误示例
```json
{
  "code": 400,
  "message": "此更新需要删除并重新创建Pod，请设置force=true以确认操作"
}
```

## 最佳实践

1. **先校验**：使用 `validateOnly: true` 检查YAML格式
2. **预览变更**：使用 `dryRun: true` 了解变更影响
3. **备份数据**：重建前确保数据安全
4. **监控服务**：重建期间监控服务状态
5. **渐进式更新**：先更新非关键字段测试功能