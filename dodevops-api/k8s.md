# Kubernetes 管理平台开发文档

## 📖 概述

本文档描述了 Kubernetes 管理平台的架构设计和日常运维管理操作功能。该平台旨在为运维人员提供一个统一的 K8s 集群管理界面，简化集群的创建、监控、运维等操作。

## 🏗️ 系统架构

### 整体架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   前端 Web UI   │────│   后端 API      │────│   K8s 集群      │
│                 │    │                 │    │                 │
│ - 集群管理      │    │ - REST API      │    │ - Master 节点   │
│ - 节点监控      │    │ - 认证鉴权      │    │ - Worker 节点   │
│ - 应用部署      │    │ - 数据持久化    │    │ - ETCD 存储     │
│ - 日志查看      │    │ - 任务调度      │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                       ┌─────────────────┐
                       │   数据库层      │
                       │                 │
                       │ - 集群信息      │
                       │ - 节点状态      │
                       │ - 任务记录      │
                       │ - 用户权限      │
                       └─────────────────┘
```

## 🎯 功能模块

### 1. 集群管理 (已实现)

#### 1.1 集群创建
- **自建集群**: 通过 Ansible 自动化部署
  - 支持单 Master / 多 Master 模式
  - 自定义 K8s 版本选择
  - 节点角色分配 (Master/Worker/ETCD)
  - 私有镜像仓库配置
  - 组件选择 (Ingress, Dashboard, etc.)

#### 1.2 集群导入
- **导入现有集群**: 通过 kubeconfig 连接
  - kubeconfig 验证
  - 集群信息自动识别
  - 权限验证

#### 1.3 集群同步 ✨
- **实时同步集群状态**
  - 集群版本信息
  - 节点数量统计
  - 节点健康状态
  - 集群运行状态

#### 1.4 集群基础操作
- 集群列表查看
- 集群详情展示
- 集群信息编辑
- 集群删除

### 2. 节点管理 (规划中)

#### 2.1 节点监控
- **节点状态监控**
  - CPU/内存使用率
  - 磁盘空间使用
  - 网络流量统计
  - Pod 运行情况

- **节点健康检查**
  - Ready/NotReady 状态
  - 资源可分配情况
  - 节点标签管理
  - 污点 (Taint) 管理

#### 2.2 节点操作
- 节点封锁/解封
- 节点排水 (Drain)
- 节点标签编辑
- 节点删除

### 3. 工作负载管理 (规划中)

#### 3.1 Deployment 管理
- **应用部署**
  - YAML 模板部署
  - 表单化部署向导
  - 镜像选择和版本管理
  - 副本数量控制

- **应用管理**
  - 扩缩容操作
  - 滚动更新
  - 回滚操作
  - 暂停/恢复部署

#### 3.2 Service 管理
- Service 创建和编辑
- 负载均衡配置
- 端口映射管理
- Ingress 规则配置

#### 3.3 ConfigMap & Secret 管理
- 配置文件管理
- 密钥管理
- 挂载点配置

### 4. 存储管理 (规划中)

#### 4.1 PV/PVC 管理
- 持久化卷创建
- 存储类管理
- 卷声明绑定
- 存储使用监控

#### 4.2 存储策略
- 动态供应配置
- 回收策略设置
- 备份策略制定

### 5. 网络管理 (规划中)

#### 5.1 网络策略
- Network Policy 配置
- 服务网格管理
- DNS 配置管理

#### 5.2 负载均衡
- Ingress Controller 管理
- 负载均衡器配置
- SSL 证书管理

### 6. 监控告警 (规划中)

#### 6.1 指标监控
- **集群级监控**
  - 集群资源使用概览
  - 节点性能指标
  - 组件健康状态

- **应用级监控**
  - Pod 资源使用
  - 应用性能指标
  - 服务响应时间

#### 6.2 日志管理
- **容器日志**
  - 实时日志查看
  - 日志搜索过滤
  - 日志下载导出

- **审计日志**
  - API 操作记录
  - 用户行为审计
  - 安全事件记录

#### 6.3 告警系统
- 告警规则配置
- 通知渠道管理
- 告警历史查看

### 7. 运维工具 (规划中)

#### 7.1 备份恢复
- **ETCD 备份**
  - 定时备份策略
  - 备份文件管理
  - 一键恢复功能

- **应用数据备份**
  - PV 数据备份
  - 配置备份
  - 应用状态快照

#### 7.2 升级维护
- **集群升级**
  - K8s 版本升级
  - 组件升级
  - 节点系统升级

- **维护模式**
  - 集群维护模式
  - 应用维护窗口
  - 维护通知

## 🔧 API 设计

### 集群管理 API

```http
# 集群基础操作
POST   /api/v1/k8s/cluster              # 创建集群
GET    /api/v1/k8s/cluster              # 获取集群列表  
GET    /api/v1/k8s/cluster/:id          # 获取集群详情
PUT    /api/v1/k8s/cluster/:id          # 更新集群信息
DELETE /api/v1/k8s/cluster/:id          # 删除集群

# 集群状态管理
GET    /api/v1/k8s/cluster/:id/status   # 获取集群状态
POST   /api/v1/k8s/cluster/:id/sync     # 同步集群信息

# 节点管理 (规划)
GET    /api/v1/k8s/cluster/:id/nodes    # 获取节点列表
POST   /api/v1/k8s/cluster/:id/nodes/:node/drain  # 节点排水
POST   /api/v1/k8s/cluster/:id/nodes/:node/cordon # 节点封锁

# 工作负载管理 (规划)
GET    /api/v1/k8s/cluster/:id/deployments        # 获取部署列表
POST   /api/v1/k8s/cluster/:id/deployments        # 创建部署
PUT    /api/v1/k8s/cluster/:id/deployments/:name  # 更新部署
DELETE /api/v1/k8s/cluster/:id/deployments/:name  # 删除部署

# 监控日志 (规划)
GET    /api/v1/k8s/cluster/:id/metrics             # 获取集群指标
GET    /api/v1/k8s/cluster/:id/pods/:pod/logs      # 获取 Pod 日志
```

## 📊 数据模型

### 集群表 (k8s_cluster)

| 字段         | 类型      | 说明                              |
|--------------|-----------|-----------------------------------|
| id           | uint      | 主键ID                            |
| name         | string    | 集群名称                          |
| version      | string    | K8s版本                           |
| status       | int       | 集群状态(1-创建中,2-运行中,3-停止,4-异常) |
| credential   | text      | kubeconfig凭证                    |
| description  | text      | 集群描述                          |
| cluster_type | int       | 集群类型(1-自建,2-导入)           |
| node_count   | int       | 节点总数                          |
| ready_nodes  | int       | 就绪节点数                        |
| master_nodes | int       | Master节点数                      |
| worker_nodes | int       | Worker节点数                      |
| last_sync_at | timestamp | 最后同步时间                      |
| created_at   | timestamp | 创建时间                          |
| updated_at   | timestamp | 更新时间                          |

## 🚀 部署架构

### 生产环境部署

```yaml
# 推荐部署架构
┌─────────────────────────────────────────────────────────┐
│                    负载均衡层                            │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │
│  │   Nginx     │  │  Nginx      │  │   Nginx     │     │
│  │  (Master)   │  │  (Backup)   │  │  (Backup)   │     │
│  └─────────────┘  └─────────────┘  └─────────────┘     │
└─────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────┐
│                    应用服务层                            │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │
│  │  API Server │  │  API Server │  │  API Server │     │
│  │   (Go应用)   │  │   (Go应用)   │  │   (Go应用)   │     │
│  └─────────────┘  └─────────────┘  └─────────────┘     │
└─────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────┐
│                    数据存储层                            │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │
│  │   MySQL     │  │    Redis    │  │  文件存储    │     │
│  │  (Master)   │  │   (Cache)   │  │    (NFS)    │     │
│  │             │  │             │  │             │     │
│  │   MySQL     │  │             │  │             │     │
│  │  (Slave)    │  │             │  │             │     │
│  └─────────────┘  └─────────────┘  └─────────────┘     │
└─────────────────────────────────────────────────────────┘
```

## 📋 详细开发计划

### Phase 1: 基础集群管理 ✅ (已完成)

#### 1.1 集群数据模型设计 ✅
```go
// 集群表结构设计
type KubeCluster struct {
    ID          uint      `gorm:"primaryKey"`
    Name        string    `gorm:"size:100;not null;uniqueIndex"`
    Version     string    `gorm:"size:50;not null"`
    Status      int       `gorm:"not null;default:1"`
    Credential  string    `gorm:"type:text"`
    Description string    `gorm:"type:text"`
    ClusterType int       `gorm:"not null;default:1"`
    NodeCount   int       `gorm:"default:0"`
    ReadyNodes  int       `gorm:"default:0"`
    MasterNodes int       `gorm:"default:0"`
    WorkerNodes int       `gorm:"default:0"`
    LastSyncAt  *time.Time `gorm:"comment:'最后同步时间'"`
    CreatedAt   time.Time `gorm:"autoCreateTime"`
    UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
```

#### 1.2 集群基础CRUD接口 ✅
- [x] **POST** `/api/v1/k8s/cluster` - 创建集群
- [x] **GET** `/api/v1/k8s/cluster` - 获取集群列表
- [x] **GET** `/api/v1/k8s/cluster/:id` - 获取集群详情
- [x] **PUT** `/api/v1/k8s/cluster/:id` - 更新集群信息
- [x] **DELETE** `/api/v1/k8s/cluster/:id` - 删除集群

#### 1.3 集群状态管理 ✅
- [x] **GET** `/api/v1/k8s/cluster/:id/status` - 获取集群状态
- [x] **POST** `/api/v1/k8s/cluster/:id/sync` - 同步集群信息
- [x] 集成 K8s client-go 库进行API调用

### Phase 2: 节点管理模块 (优先级：高)

#### 2.1 节点数据模型设计
**开发步骤 1**: 创建节点表结构
```go
// nodes 表
type KubeNode struct {
    ID              uint      `gorm:"primaryKey"`
    ClusterID       uint      `gorm:"not null;index"`
    Name            string    `gorm:"size:100;not null"`
    Role            string    `gorm:"size:20"`      // master/worker/etcd
    Status          string    `gorm:"size:20"`      // Ready/NotReady/Unknown
    InternalIP      string    `gorm:"size:50"`
    ExternalIP      string    `gorm:"size:50"`
    KubeletVersion  string    `gorm:"size:50"`
    OSImage         string    `gorm:"size:100"`
    CPUCapacity     string    `gorm:"size:20"`
    MemoryCapacity  string    `gorm:"size:20"`
    CPUAllocatable  string    `gorm:"size:20"`
    MemoryAllocatable string  `gorm:"size:20"`
    Labels          string    `gorm:"type:text"`    // JSON格式存储
    Taints          string    `gorm:"type:text"`    // JSON格式存储
    LastSyncAt      *time.Time
    CreatedAt       time.Time `gorm:"autoCreateTime"`
    UpdatedAt       time.Time `gorm:"autoUpdateTime"`
}
```

#### 2.2 节点基础接口开发
**开发步骤 2**: 实现节点管理API
- [ ] **GET** `/api/v1/k8s/cluster/:id/nodes` - 获取节点列表
- [ ] **GET** `/api/v1/k8s/cluster/:id/nodes/:name` - 获取节点详情
- [ ] **POST** `/api/v1/k8s/cluster/:id/nodes/sync` - 同步所有节点信息

**开发步骤 3**: 实现节点操作API
- [ ] **POST** `/api/v1/k8s/cluster/:id/nodes/:name/cordon` - 封锁节点
- [ ] **POST** `/api/v1/k8s/cluster/:id/nodes/:name/uncordon` - 解封节点
- [ ] **POST** `/api/v1/k8s/cluster/:id/nodes/:name/drain` - 排空节点
- [ ] **PUT** `/api/v1/k8s/cluster/:id/nodes/:name/labels` - 更新节点标签
- [ ] **PUT** `/api/v1/k8s/cluster/:id/nodes/:name/taints` - 更新节点污点

#### 2.3 节点监控指标
**开发步骤 4**: 实现节点资源监控
```go
// node_metrics 表
type NodeMetrics struct {
    ID           uint      `gorm:"primaryKey"`
    NodeID       uint      `gorm:"not null;index"`
    CPUUsage     float64   // CPU使用率百分比
    MemoryUsage  float64   // 内存使用率百分比
    DiskUsage    float64   // 磁盘使用率百分比
    NetworkIn    int64     // 网络入流量
    NetworkOut   int64     // 网络出流量
    PodCount     int       // Pod数量
    Timestamp    time.Time `gorm:"index"`
}
```

- [ ] **GET** `/api/v1/k8s/cluster/:id/nodes/:name/metrics` - 获取节点监控指标
- [ ] **GET** `/api/v1/k8s/cluster/:id/nodes/:name/metrics/history` - 获取历史监控数据

### Phase 3: Pod 管理模块 (优先级：高)

#### 3.1 Pod数据模型设计
**开发步骤 5**: 创建Pod表结构
```go
// pods 表
type KubePod struct {
    ID          uint      `gorm:"primaryKey"`
    ClusterID   uint      `gorm:"not null;index"`
    Namespace   string    `gorm:"size:100;not null;index"`
    Name        string    `gorm:"size:100;not null"`
    NodeName    string    `gorm:"size:100"`
    Phase       string    `gorm:"size:20"`  // Pending/Running/Succeeded/Failed/Unknown
    Ready       string    `gorm:"size:20"`  // 1/1, 2/3 等格式
    Restarts    int
    Age         time.Time
    IP          string    `gorm:"size:50"`
    Labels      string    `gorm:"type:text"`
    Annotations string    `gorm:"type:text"`
    OwnerType   string    `gorm:"size:50"`  // Deployment/StatefulSet/DaemonSet
    OwnerName   string    `gorm:"size:100"`
    LastSyncAt  *time.Time
    CreatedAt   time.Time `gorm:"autoCreateTime"`
    UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
```

#### 3.2 Pod管理接口
**开发步骤 6**: 实现Pod管理API
- [ ] **GET** `/api/v1/k8s/cluster/:id/pods` - 获取Pod列表 (支持namespace过滤)
- [ ] **GET** `/api/v1/k8s/cluster/:id/namespaces/:ns/pods` - 获取指定namespace的Pod
- [ ] **GET** `/api/v1/k8s/cluster/:id/namespaces/:ns/pods/:name` - 获取Pod详情
- [ ] **DELETE** `/api/v1/k8s/cluster/:id/namespaces/:ns/pods/:name` - 删除Pod
- [ ] **GET** `/api/v1/k8s/cluster/:id/namespaces/:ns/pods/:name/logs` - 获取Pod日志
- [ ] **POST** `/api/v1/k8s/cluster/:id/pods/sync` - 同步所有Pod信息

#### 3.3 Pod日志和监控
**开发步骤 7**: 实现Pod日志查看功能
```go
// pod_logs 实时日志接口
type LogRequest struct {
    Follow     bool   `json:"follow"`     // 是否实时跟踪
    Previous   bool   `json:"previous"`   // 是否查看前一个容器的日志
    TailLines  *int64 `json:"tailLines"`  // 显示最后N行
    Container  string `json:"container"`  // 指定容器名
    Timestamps bool   `json:"timestamps"` // 是否显示时间戳
}
```

- [ ] **GET** `/api/v1/k8s/cluster/:id/namespaces/:ns/pods/:name/logs` - 获取Pod日志
- [ ] **GET** `/api/v1/k8s/cluster/:id/namespaces/:ns/pods/:name/containers/:container/logs` - 获取指定容器日志
- [ ] **WebSocket** `/ws/k8s/cluster/:id/namespaces/:ns/pods/:name/logs` - 实时日志流

### Phase 4: Deployment 工作负载管理 (优先级：中)

#### 4.1 Deployment数据模型
**开发步骤 8**: 创建Deployment表结构
```go
// deployments 表
type KubeDeployment struct {
    ID            uint      `gorm:"primaryKey"`
    ClusterID     uint      `gorm:"not null;index"`
    Namespace     string    `gorm:"size:100;not null;index"`
    Name          string    `gorm:"size:100;not null"`
    Replicas      int32     `gorm:"default:1"`
    ReadyReplicas int32
    UpToDate      int32
    Available     int32
    Age           time.Time
    Labels        string    `gorm:"type:text"`
    Selector      string    `gorm:"type:text"`
    Strategy      string    `gorm:"size:50"`  // RollingUpdate/Recreate
    Images        string    `gorm:"type:text"` // JSON数组存储镜像列表
    LastSyncAt    *time.Time
    CreatedAt     time.Time `gorm:"autoCreateTime"`
    UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}
```

#### 4.2 Deployment管理接口
**开发步骤 9**: 实现Deployment管理API
- [ ] **GET** `/api/v1/k8s/cluster/:id/deployments` - 获取Deployment列表
- [ ] **GET** `/api/v1/k8s/cluster/:id/namespaces/:ns/deployments` - 获取指定namespace的Deployment
- [ ] **GET** `/api/v1/k8s/cluster/:id/namespaces/:ns/deployments/:name` - 获取Deployment详情
- [ ] **POST** `/api/v1/k8s/cluster/:id/namespaces/:ns/deployments` - 创建Deployment
- [ ] **PUT** `/api/v1/k8s/cluster/:id/namespaces/:ns/deployments/:name` - 更新Deployment
- [ ] **DELETE** `/api/v1/k8s/cluster/:id/namespaces/:ns/deployments/:name` - 删除Deployment

#### 4.3 Deployment操作接口
**开发步骤 10**: 实现Deployment高级操作
- [ ] **POST** `/api/v1/k8s/cluster/:id/namespaces/:ns/deployments/:name/scale` - 扩缩容
- [ ] **POST** `/api/v1/k8s/cluster/:id/namespaces/:ns/deployments/:name/restart` - 重启
- [ ] **POST** `/api/v1/k8s/cluster/:id/namespaces/:ns/deployments/:name/rollback` - 回滚
- [ ] **GET** `/api/v1/k8s/cluster/:id/namespaces/:ns/deployments/:name/history` - 查看历史版本
- [ ] **GET** `/api/v1/k8s/cluster/:id/namespaces/:ns/deployments/:name/events` - 查看事件

### Phase 5: Service 服务管理 (优先级：中)

#### 5.1 Service数据模型
**开发步骤 11**: 创建Service表结构
```go
// services 表
type KubeService struct {
    ID           uint      `gorm:"primaryKey"`
    ClusterID    uint      `gorm:"not null;index"`
    Namespace    string    `gorm:"size:100;not null;index"`
    Name         string    `gorm:"size:100;not null"`
    Type         string    `gorm:"size:20"`    // ClusterIP/NodePort/LoadBalancer
    ClusterIP    string    `gorm:"size:50"`
    ExternalIPs  string    `gorm:"type:text"`  // JSON数组
    Ports        string    `gorm:"type:text"`  // JSON数组存储端口信息
    Selector     string    `gorm:"type:text"`  // JSON格式
    Labels       string    `gorm:"type:text"`
    Age          time.Time
    LastSyncAt   *time.Time
    CreatedAt    time.Time `gorm:"autoCreateTime"`
    UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}
```

#### 5.2 Service管理接口
**开发步骤 12**: 实现Service管理API
- [ ] **GET** `/api/v1/k8s/cluster/:id/services` - 获取Service列表
- [ ] **GET** `/api/v1/k8s/cluster/:id/namespaces/:ns/services` - 获取指定namespace的Service
- [ ] **GET** `/api/v1/k8s/cluster/:id/namespaces/:ns/services/:name` - 获取Service详情
- [ ] **POST** `/api/v1/k8s/cluster/:id/namespaces/:ns/services` - 创建Service
- [ ] **PUT** `/api/v1/k8s/cluster/:id/namespaces/:ns/services/:name` - 更新Service
- [ ] **DELETE** `/api/v1/k8s/cluster/:id/namespaces/:ns/services/:name` - 删除Service

### Phase 6: ConfigMap & Secret 配置管理 (优先级：中)

#### 6.1 ConfigMap数据模型
**开发步骤 13**: 创建ConfigMap表结构
```go
// configmaps 表
type KubeConfigMap struct {
    ID         uint      `gorm:"primaryKey"`
    ClusterID  uint      `gorm:"not null;index"`
    Namespace  string    `gorm:"size:100;not null;index"`
    Name       string    `gorm:"size:100;not null"`
    Data       string    `gorm:"type:longtext"`  // JSON格式存储配置数据
    Labels     string    `gorm:"type:text"`
    Age        time.Time
    LastSyncAt *time.Time
    CreatedAt  time.Time `gorm:"autoCreateTime"`
    UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

// secrets 表
type KubeSecret struct {
    ID         uint      `gorm:"primaryKey"`
    ClusterID  uint      `gorm:"not null;index"`
    Namespace  string    `gorm:"size:100;not null;index"`
    Name       string    `gorm:"size:100;not null"`
    Type       string    `gorm:"size:50"`        // Opaque/TLS/dockerconfigjson
    DataKeys   string    `gorm:"type:text"`      // JSON数组存储key列表，不存储实际值
    Labels     string    `gorm:"type:text"`
    Age        time.Time
    LastSyncAt *time.Time
    CreatedAt  time.Time `gorm:"autoCreateTime"`
    UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}
```

#### 6.2 ConfigMap & Secret管理接口
**开发步骤 14**: 实现配置管理API
```http
# ConfigMap管理
GET    /api/v1/k8s/cluster/:id/configmaps                    # 获取ConfigMap列表
GET    /api/v1/k8s/cluster/:id/namespaces/:ns/configmaps     # 获取指定namespace的ConfigMap
GET    /api/v1/k8s/cluster/:id/namespaces/:ns/configmaps/:name # 获取ConfigMap详情
POST   /api/v1/k8s/cluster/:id/namespaces/:ns/configmaps     # 创建ConfigMap
PUT    /api/v1/k8s/cluster/:id/namespaces/:ns/configmaps/:name # 更新ConfigMap
DELETE /api/v1/k8s/cluster/:id/namespaces/:ns/configmaps/:name # 删除ConfigMap

# Secret管理
GET    /api/v1/k8s/cluster/:id/secrets                       # 获取Secret列表
GET    /api/v1/k8s/cluster/:id/namespaces/:ns/secrets        # 获取指定namespace的Secret
GET    /api/v1/k8s/cluster/:id/namespaces/:ns/secrets/:name  # 获取Secret详情
POST   /api/v1/k8s/cluster/:id/namespaces/:ns/secrets        # 创建Secret
PUT    /api/v1/k8s/cluster/:id/namespaces/:ns/secrets/:name  # 更新Secret
DELETE /api/v1/k8s/cluster/:id/namespaces/:ns/secrets/:name  # 删除Secret
```

### Phase 7: Namespace 命名空间管理 (优先级：低)

#### 7.1 Namespace数据模型
**开发步骤 15**: 创建Namespace表结构
```go
// namespaces 表
type KubeNamespace struct {
    ID          uint      `gorm:"primaryKey"`
    ClusterID   uint      `gorm:"not null;index"`
    Name        string    `gorm:"size:100;not null"`
    Status      string    `gorm:"size:20"`    // Active/Terminating
    Labels      string    `gorm:"type:text"`
    Annotations string    `gorm:"type:text"`
    Age         time.Time
    LastSyncAt  *time.Time
    CreatedAt   time.Time `gorm:"autoCreateTime"`
    UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
```

#### 7.2 Namespace管理接口
**开发步骤 16**: 实现Namespace管理API
- [ ] **GET** `/api/v1/k8s/cluster/:id/namespaces` - 获取命名空间列表
- [ ] **GET** `/api/v1/k8s/cluster/:id/namespaces/:name` - 获取命名空间详情
- [ ] **POST** `/api/v1/k8s/cluster/:id/namespaces` - 创建命名空间
- [ ] **PUT** `/api/v1/k8s/cluster/:id/namespaces/:name` - 更新命名空间
- [ ] **DELETE** `/api/v1/k8s/cluster/:id/namespaces/:name` - 删除命名空间

### Phase 8: 监控告警系统 (优先级：中)

#### 8.1 监控数据模型
**开发步骤 17**: 创建监控指标表
```go
// cluster_metrics 集群监控指标
type ClusterMetrics struct {
    ID              uint      `gorm:"primaryKey"`
    ClusterID       uint      `gorm:"not null;index"`
    CPUUsage        float64   // 集群CPU使用率
    MemoryUsage     float64   // 集群内存使用率
    StorageUsage    float64   // 集群存储使用率
    NodeCount       int       // 节点总数
    PodCount        int       // Pod总数
    ServiceCount    int       // Service总数
    DeploymentCount int       // Deployment总数
    Timestamp       time.Time `gorm:"index"`
}

// alerts 告警记录表
type Alert struct {
    ID          uint      `gorm:"primaryKey"`
    ClusterID   uint      `gorm:"not null;index"`
    RuleName    string    `gorm:"size:100"`
    Severity    string    `gorm:"size:20"`    // critical/warning/info
    Message     string    `gorm:"type:text"`
    Status      string    `gorm:"size:20"`    // firing/resolved
    StartsAt    time.Time
    EndsAt      *time.Time
    Labels      string    `gorm:"type:text"`  // JSON格式
    Annotations string    `gorm:"type:text"`  // JSON格式
    CreatedAt   time.Time `gorm:"autoCreateTime"`
    UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
```

#### 8.2 监控接口开发
**开发步骤 18**: 实现监控API
- [ ] **GET** `/api/v1/k8s/cluster/:id/metrics` - 获取集群监控指标
- [ ] **GET** `/api/v1/k8s/cluster/:id/metrics/history` - 获取历史监控数据
- [ ] **GET** `/api/v1/k8s/cluster/:id/alerts` - 获取告警列表
- [ ] **GET** `/api/v1/k8s/cluster/:id/alerts/:id` - 获取告警详情
- [ ] **POST** `/api/v1/k8s/cluster/:id/alerts/:id/resolve` - 处理告警

### Phase 9: 事件管理 (优先级：低)

#### 9.1 事件数据模型
**开发步骤 19**: 创建事件表结构
```go
// events 事件表
type KubeEvent struct {
    ID              uint      `gorm:"primaryKey"`
    ClusterID       uint      `gorm:"not null;index"`
    Namespace       string    `gorm:"size:100;index"`
    Name            string    `gorm:"size:100"`
    Reason          string    `gorm:"size:100;index"`
    Message         string    `gorm:"type:text"`
    Type            string    `gorm:"size:20"`    // Normal/Warning
    Count           int32
    FirstTimestamp  time.Time `gorm:"index"`
    LastTimestamp   time.Time `gorm:"index"`
    InvolvedObject  string    `gorm:"type:text"`  // JSON格式存储对象信息
    Source          string    `gorm:"type:text"`  // JSON格式存储事件源
    CreatedAt       time.Time `gorm:"autoCreateTime"`
}
```

#### 9.2 事件管理接口
**开发步骤 20**: 实现事件管理API
- [ ] **GET** `/api/v1/k8s/cluster/:id/events` - 获取事件列表
- [ ] **GET** `/api/v1/k8s/cluster/:id/namespaces/:ns/events` - 获取指定namespace事件
- [ ] **GET** `/api/v1/k8s/cluster/:id/events/search` - 事件搜索
- [ ] **POST** `/api/v1/k8s/cluster/:id/events/sync` - 同步事件

### Phase 10: 高级功能开发 (优先级：低)

#### 10.1 HPA自动扩缩容
**开发步骤 21**: 实现HPA管理
```go
// hpa 表
type KubeHPA struct {
    ID              uint      `gorm:"primaryKey"`
    ClusterID       uint      `gorm:"not null;index"`
    Namespace       string    `gorm:"size:100;not null;index"`
    Name            string    `gorm:"size:100;not null"`
    TargetRef       string    `gorm:"type:text"`  // JSON格式存储目标引用
    MinReplicas     int32
    MaxReplicas     int32
    CurrentReplicas int32
    TargetCPU       *int32    // CPU目标百分比
    CurrentCPU      *int32    // 当前CPU百分比
    LastSyncAt      *time.Time
    CreatedAt       time.Time `gorm:"autoCreateTime"`
    UpdatedAt       time.Time `gorm:"autoUpdateTime"`
}
```

- [ ] **GET** `/api/v1/k8s/cluster/:id/hpa` - 获取HPA列表
- [ ] **POST** `/api/v1/k8s/cluster/:id/namespaces/:ns/hpa` - 创建HPA
- [ ] **PUT** `/api/v1/k8s/cluster/:id/namespaces/:ns/hpa/:name` - 更新HPA
- [ ] **DELETE** `/api/v1/k8s/cluster/:id/namespaces/:ns/hpa/:name` - 删除HPA

#### 10.2 Ingress管理
**开发步骤 22**: 实现Ingress管理
- [ ] Ingress数据模型设计
- [ ] Ingress CRUD接口开发
- [ ] SSL证书管理
- [ ] 域名解析配置

#### 10.3 PV/PVC存储管理
**开发步骤 23**: 实现存储管理
- [ ] PersistentVolume数据模型
- [ ] PVC管理接口
- [ ] StorageClass管理
- [ ] 存储使用监控

## 🚧 开发优先级和时间规划

### 第一阶段 (已完成) ✅
- ✅ 1-4步骤：集群基础管理功能

### 第二阶段 (1-2周)
- 📋 5-7步骤：节点管理和Pod管理 (核心功能)

### 第三阶段 (2-3周)  
- 📋 8-10步骤：Deployment工作负载管理

### 第四阶段 (1-2周)
- 📋 11-14步骤：Service和配置管理

### 第五阶段 (2-3周)
- 📋 15-18步骤：命名空间和监控告警

### 第六阶段 (1-2周)
- 📋 19-20步骤：事件管理

### 第七阶段 (按需开发)
- 📋 21-23步骤：高级功能 (HPA/Ingress/存储)

## 🔄 每个开发步骤的具体任务

### 标准开发流程
每个功能模块都按以下步骤进行：

1. **数据模型设计** (0.5天)
   - 设计数据库表结构
   - 创建对应的Go结构体
   - 添加数据库迁移

2. **DAO层开发** (0.5天)
   - 实现数据访问层接口
   - 添加基础CRUD方法
   - 实现数据同步逻辑

3. **Service层开发** (1天)
   - 实现业务逻辑层
   - 集成K8s API调用
   - 添加数据验证和处理

4. **Controller层开发** (1天)
   - 实现REST API接口
   - 添加参数验证
   - 处理HTTP请求响应

5. **测试和联调** (1天)
   - 单元测试编写
   - API接口测试
   - 前后端联调

6. **文档更新** (0.5天)
   - 更新API文档
   - 添加使用说明
   - 更新开发进度

## 📊 开发进度跟踪

### 进度表格
| 阶段 | 功能模块 | 开发步骤 | 预计时间 | 状态 | 负责人 | 备注 |
|------|----------|----------|----------|------|--------|------|
| Phase 1 | 集群管理 | 1-4 | 2周 | ✅ 完成 | 开发团队 | 已上线 |
| Phase 2 | 节点管理 | 5-7 | 1-2周 | ⏳ 待开发 | - | 优先级高 |
| Phase 3 | Pod管理 | 8-10 | 2-3周 | ⏳ 待开发 | - | 核心功能 |
| Phase 4 | 工作负载 | 11-14 | 1-2周 | ⏳ 待开发 | - | - |
| Phase 5 | 服务配置 | 15-18 | 2-3周 | ⏳ 待开发 | - | - |
| Phase 6 | 监控告警 | 19-20 | 1-2周 | ⏳ 待开发 | - | - |
| Phase 7 | 高级功能 | 21-23 | 按需 | ⏳ 待开发 | - | 低优先级 |

### 里程碑节点
- 🎯 **M1** (已完成): 集群基础管理功能上线
- 🎯 **M2** (预计2周后): 节点和Pod管理功能完成
- 🎯 **M3** (预计1个月后): 工作负载管理功能完成  
- 🎯 **M4** (预计2个月后): 完整的K8s管理平台MVP版本
- 🎯 **M5** (预计3个月后): 监控告警功能完善的生产版本

## 🔒 安全考虑

### 认证授权
- JWT Token 认证
- RBAC 权限控制
- API 访问限流
- kubeconfig 安全存储

### 网络安全
- HTTPS 强制加密
- API Gateway 防护
- 网络策略隔离
- 安全组配置

### 数据安全
- 敏感数据加密存储
- 数据备份和恢复
- 审计日志记录
- 操作行为追踪

## 📈 性能优化

### 缓存策略
- Redis 缓存热点数据
- 集群状态缓存
- API 响应缓存
- 静态资源缓存

### 数据库优化
- 索引优化
- 查询优化
- 连接池管理
- 读写分离

### 前端优化
- 组件懒加载
- 数据分页加载
- 实时数据推送
- 界面响应优化

---

## 📞 联系信息

**开发团队**: DevOps Platform Team  
**文档版本**: v1.0  
**更新时间**: 2025-01-15  
**维护者**: Kubernetes 架构团队

> 本文档将随着平台功能的不断完善而持续更新。如有问题或建议，请联系开发团队。