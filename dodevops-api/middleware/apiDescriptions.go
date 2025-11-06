// API描述映射
// author xiaoRui

package middleware

import (
	"regexp"
	"strings"
)

// GetAPIDescription 根据 URL 路径和 HTTP 方法获取 API 的中文描述
// 这些描述从各个 controller 的 @Description 注释中提取
func GetAPIDescription(url, method string) string {
	// 先尝试精确匹配（URL 路径）
	if desc, ok := getDescriptionMap()[url]; ok {
		return desc
	}

	// 尝试 URL + Method 组合匹配（RESTful 风格）
	key := method + ":" + url
	if desc, ok := getMethodBasedDescriptionMap()[key]; ok {
		return desc
	}

	// 如果精确匹配失败，尝试模糊匹配（处理带参数的路由）
	return matchDynamicRoute(url)
}

// getDescriptionMap 返回 API 路径到描述的映射
func getDescriptionMap() map[string]string {
	descriptions := map[string]string{
		// ========== 系统管理 ==========
		"/api/v1/admin/login":          "用户登录",
		"/api/v1/admin/logout":         "用户登出",
		"/api/v1/admin/add":            "新增管理员",
		"/api/v1/admin/update":         "修改管理员",
		"/api/v1/admin/delete":         "删除管理员",
		"/api/v1/admin/batchDelete":    "批量删除管理员",
		"/api/v1/admin/changeStatus":   "修改管理员状态",
		"/api/v1/admin/resetPwd":       "重置管理员密码",
		"/api/v1/admin/updateUserInfo": "修改个人信息",
		"/api/v1/admin/updatePwd":      "修改密码",

		"/api/v1/role/add":          "新增角色",
		"/api/v1/role/update":       "修改角色",
		"/api/v1/role/delete":       "删除角色",
		"/api/v1/role/batchDelete":  "批量删除角色",
		"/api/v1/role/changeStatus": "修改角色状态",

		"/api/v1/menu/add":    "新增菜单",
		"/api/v1/menu/update": "修改菜单",
		"/api/v1/menu/delete": "删除菜单",

		"/api/v1/dept/add":    "新增部门",
		"/api/v1/dept/update": "修改部门",
		"/api/v1/dept/delete": "删除部门",

		"/api/v1/post/add":          "新增岗位",
		"/api/v1/post/update":       "修改岗位",
		"/api/v1/post/delete":       "删除岗位",
		"/api/v1/post/batchDelete":  "批量删除岗位",
		"/api/v1/post/updateStatus": "修改岗位状态",

		"/api/v1/sysOperationLog/delete":      "删除操作日志",
		"/api/v1/sysOperationLog/batchDelete": "批量删除操作日志",
		"/api/v1/sysOperationLog/clean":       "清空操作日志",

		// ========== 配置中心 ==========
		"/api/v1/config/ecsauthadd":    "新增ECS认证",
		"/api/v1/config/ecsauthupdate": "修改ECS认证",
		"/api/v1/config/ecsauthdelete": "删除ECS认证",

		"/api/v1/config/keymanage/sync": "同步云主机",

		"/api/v1/config/sync-schedule/toggle-status": "切换同步调度状态",
		"/api/v1/config/sync-schedule/trigger":       "手动触发同步",

		// ========== CMDB管理 ==========
		"/api/v1/cmdb/groupadd":    "新增资产分组",
		"/api/v1/cmdb/groupupdate": "修改资产分组",
		"/api/v1/cmdb/groupdelete": "删除资产分组",

		"/api/v1/cmdb/hostcreate":  "创建主机",
		"/api/v1/cmdb/hostupdate":  "修改主机",
		"/api/v1/cmdb/hostdelete":  "删除主机",
		"/api/v1/cmdb/hostimport":  "导入主机",
		"/api/v1/cmdb/hostsync":    "同步主机信息",

		"/api/v1/cmdb/hostcloudcreatealiyun":  "创建阿里云主机",
		"/api/v1/cmdb/hostcloudcreatetencent": "创建腾讯云主机",

		"/api/v1/cmdb/sqlLog/delete": "删除SQL日志",
		"/api/v1/cmdb/sqlLog/clean":  "清空SQL日志",

		// ========== 任务中心 ==========
		"/api/v1/template/add":    "新增任务模板",
		"/api/v1/template/update": "修改任务模板",
		"/api/v1/template/delete": "删除任务模板",

		"/api/v1/task/add":    "新增任务",
		"/api/v1/task/update": "修改任务",
		"/api/v1/task/delete": "删除任务",

		"/api/v1/taskjob/start": "启动任务作业",
		"/api/v1/taskjob/stop":  "停止任务作业",

		"/api/v1/task/monitor/queue/clear-failed": "清空失败队列",
		"/api/v1/task/monitor/queue/retry-failed": "重试失败任务",
		"/api/v1/task/monitor/scheduled/pause":    "暂停定时任务",
		"/api/v1/task/monitor/scheduled/resume":   "恢复定时任务",
		"/api/v1/task/monitor/scheduled/reset":    "重置定时任务状态",

		"/api/v1/task/ansible": "创建Ansible任务",
		"/api/v1/task/k8s":     "创建K8s任务",

		// ========== 监控中心 ==========
		"/api/v1/monitor/agent/deploy":    "部署Agent",
		"/api/v1/monitor/agent/uninstall": "卸载Agent",

		// ========== 应用管理 ==========
		"/api/v1/apps":                          "创建应用",
		"/api/v1/apps/deployment/quick":         "创建快速发布",
		"/api/v1/apps/deployment/execute":       "执行快速发布",
		"/api/v1/apps/jenkins-job/validate":     "验证Jenkins任务",
		"/api/v1/jenkins/test-connection":       "测试Jenkins连接",
	}

	return descriptions
}

// getMethodBasedDescriptionMap 返回基于 HTTP Method + URL 的映射（RESTful 风格）
func getMethodBasedDescriptionMap() map[string]string {
	return map[string]string{
		// 密钥管理（RESTful 风格）
		"post:/api/v1/config/keymanage":   "新增密钥",
		"put:/api/v1/config/keymanage":    "修改密钥",
		"delete:/api/v1/config/keymanage": "删除密钥",

		// 账号认证（RESTful 风格）
		"post:/api/v1/config/accountauth":   "新增账号认证",
		"put:/api/v1/config/accountauth":    "修改账号认证",
		"delete:/api/v1/config/accountauth": "删除账号认证",

		// 同步调度（RESTful 风格）
		"post:/api/v1/config/sync-schedule":   "新增同步调度",
		"put:/api/v1/config/sync-schedule":    "修改同步调度",
		"delete:/api/v1/config/sync-schedule": "删除同步调度",

		// 数据库管理（RESTful 风格）
		"post:/api/v1/cmdb/database":   "创建数据库",
		"put:/api/v1/cmdb/database":    "修改数据库",
		"delete:/api/v1/cmdb/database": "删除数据库",

		// SQL执行（RESTful 风格）
		"post:/api/v1/cmdb/sql":   "执行SQL插入",
		"put:/api/v1/cmdb/sql":    "执行SQL更新",
		"delete:/api/v1/cmdb/sql": "执行SQL删除",
		"post:/api/v1/cmdb/sql/select":  "执行SQL查询",
		"post:/api/v1/cmdb/sql/execute": "执行原生SQL",

		// K8s集群管理（RESTful 风格）
		"post:/api/v1/k8s/cluster": "创建K8s集群",
	}
}

// matchDynamicRoute 匹配带参数的动态路由
func matchDynamicRoute(url string) string {
	// 定义动态路由的匹配规则（正则模式 -> 描述）
	dynamicRoutes := []struct {
		pattern *regexp.Regexp
		desc    string
	}{
		// 监控中心 - Agent相关
		{regexp.MustCompile(`^/api/v1/monitor/agent/delete/\d+$`), "删除Agent"},
		{regexp.MustCompile(`^/api/v1/monitor/agent/status/\d+$`), "获取Agent状态"},
		{regexp.MustCompile(`^/api/v1/monitor/agent/restart/\d+$`), "重启Agent"},

		// 系统管理
		{regexp.MustCompile(`^/api/v1/admin/delete/\d+$`), "删除管理员"},
		{regexp.MustCompile(`^/api/v1/role/delete/\d+$`), "删除角色"},
		{regexp.MustCompile(`^/api/v1/menu/delete/\d+$`), "删除菜单"},
		{regexp.MustCompile(`^/api/v1/dept/delete/\d+$`), "删除部门"},
		{regexp.MustCompile(`^/api/v1/post/delete/\d+$`), "删除岗位"},

		// 任务中心 - Ansible任务
		{regexp.MustCompile(`^/api/v1/task/ansible/\d+$`), "删除Ansible任务"},
		{regexp.MustCompile(`^/api/v1/task/ansible/\d+/start$`), "启动Ansible任务"},

		// 应用管理
		{regexp.MustCompile(`^/api/v1/apps/\d+$`), "操作应用"},
		{regexp.MustCompile(`^/api/v1/apps/\d+/jenkins-envs$`), "添加Jenkins环境配置"},
		{regexp.MustCompile(`^/api/v1/apps/\d+/jenkins-envs/\d+$`), "操作Jenkins环境配置"},
		{regexp.MustCompile(`^/api/v1/apps/deployment/\d+$`), "操作快速发布"},

		// Jenkins相关
		{regexp.MustCompile(`^/api/v1/jenkins/\d+/jobs/[^/]+/start$`), "启动Jenkins任务"},
		{regexp.MustCompile(`^/api/v1/jenkins/\d+/jobs/[^/]+/builds/\d+/stop$`), "停止Jenkins构建"},

		// K8s集群管理
		{regexp.MustCompile(`^/api/v1/k8s/cluster/\d+$`), "操作K8s集群"},
		{regexp.MustCompile(`^/api/v1/k8s/cluster/\d+/sync$`), "同步K8s集群"},
		{regexp.MustCompile(`^/api/v1/k8s/cluster/\d+/namespaces$`), "创建K8s命名空间"},
		{regexp.MustCompile(`^/api/v1/k8s/cluster/\d+/namespaces/[^/]+$`), "操作K8s命名空间"},
		{regexp.MustCompile(`^/api/v1/k8s/cluster/\d+/namespaces/[^/]+/deployments$`), "创建K8s部署"},
		{regexp.MustCompile(`^/api/v1/k8s/cluster/\d+/namespaces/[^/]+/deployments/[^/]+$`), "操作K8s部署"},
		{regexp.MustCompile(`^/api/v1/k8s/cluster/\d+/namespaces/[^/]+/deployments/[^/]+/scale$`), "伸缩K8s部署"},
		{regexp.MustCompile(`^/api/v1/k8s/cluster/\d+/namespaces/[^/]+/deployments/[^/]+/restart$`), "重启K8s部署"},
		{regexp.MustCompile(`^/api/v1/k8s/cluster/\d+/namespaces/[^/]+/deployments/[^/]+/rollback$`), "回滚K8s部署"},
		{regexp.MustCompile(`^/api/v1/k8s/cluster/\d+/namespaces/[^/]+/deployments/[^/]+/pause$`), "暂停K8s部署"},
		{regexp.MustCompile(`^/api/v1/k8s/cluster/\d+/namespaces/[^/]+/deployments/[^/]+/resume$`), "恢复K8s部署"},
		{regexp.MustCompile(`^/api/v1/k8s/cluster/\d+/namespaces/[^/]+/pods/[^/]+$`), "操作K8s Pod"},
		{regexp.MustCompile(`^/api/v1/k8s/cluster/\d+/namespaces/[^/]+/services$`), "创建K8s服务"},
		{regexp.MustCompile(`^/api/v1/k8s/cluster/\d+/namespaces/[^/]+/services/[^/]+$`), "操作K8s服务"},
		{regexp.MustCompile(`^/api/v1/k8s/cluster/\d+/namespaces/[^/]+/ingresses$`), "创建K8s Ingress"},
		{regexp.MustCompile(`^/api/v1/k8s/cluster/\d+/namespaces/[^/]+/ingresses/[^/]+$`), "操作K8s Ingress"},
		{regexp.MustCompile(`^/api/v1/k8s/cluster/\d+/namespaces/[^/]+/configmaps$`), "创建K8s配置"},
		{regexp.MustCompile(`^/api/v1/k8s/cluster/\d+/namespaces/[^/]+/configmaps/[^/]+$`), "操作K8s配置"},
		{regexp.MustCompile(`^/api/v1/k8s/cluster/\d+/namespaces/[^/]+/secrets$`), "创建K8s密钥"},
		{regexp.MustCompile(`^/api/v1/k8s/cluster/\d+/namespaces/[^/]+/secrets/[^/]+$`), "操作K8s密钥"},
		{regexp.MustCompile(`^/api/v1/k8s/cluster/\d+/namespaces/[^/]+/pvcs$`), "创建K8s存储卷"},
		{regexp.MustCompile(`^/api/v1/k8s/cluster/\d+/namespaces/[^/]+/pvcs/[^/]+$`), "操作K8s存储卷"},
		{regexp.MustCompile(`^/api/v1/k8s/cluster/\d+/pvs$`), "创建K8s持久卷"},
		{regexp.MustCompile(`^/api/v1/k8s/cluster/\d+/pvs/[^/]+$`), "操作K8s持久卷"},
		{regexp.MustCompile(`^/api/v1/k8s/cluster/\d+/storageclasses$`), "创建K8s存储类"},
		{regexp.MustCompile(`^/api/v1/k8s/cluster/\d+/storageclasses/[^/]+$`), "操作K8s存储类"},
		{regexp.MustCompile(`^/api/v1/k8s/cluster/\d+/nodes/[^/]+/taints$`), "操作节点污点"},
		{regexp.MustCompile(`^/api/v1/k8s/cluster/\d+/nodes/[^/]+/labels$`), "操作节点标签"},
		{regexp.MustCompile(`^/api/v1/k8s/cluster/\d+/nodes/[^/]+/cordon$`), "封锁节点"},
		{regexp.MustCompile(`^/api/v1/k8s/cluster/\d+/nodes/[^/]+/drain$`), "驱逐节点"},
		{regexp.MustCompile(`^/api/v1/k8s/cluster/\d+/yaml/validate$`), "验证K8s YAML"},
		{regexp.MustCompile(`^/api/v1/k8s/cluster/\d+/namespaces/[^/]+/workload-yaml$`), "更新工作负载YAML"},
	}

	// 遍历匹配规则
	for _, route := range dynamicRoutes {
		if route.pattern.MatchString(url) {
			return route.desc
		}
	}

	// 如果都不匹配，尝试提取操作类型
	return extractActionFromURL(url)
}

// extractActionFromURL 从 URL 中提取操作类型作为描述
func extractActionFromURL(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) >= 2 {
		action := parts[len(parts)-2] // 倒数第二个部分通常是操作类型
		switch action {
		case "add", "create":
			return "新增"
		case "update":
			return "修改"
		case "delete":
			return "删除"
		case "batchDelete":
			return "批量删除"
		case "execute", "start":
			return "执行"
		case "stop":
			return "停止"
		case "deploy":
			return "部署"
		case "uninstall":
			return "卸载"
		case "restart":
			return "重启"
		case "toggle", "toggle-status":
			return "切换状态"
		case "sync":
			return "同步"
		case "pause":
			return "暂停"
		case "resume":
			return "恢复"
		case "rollback":
			return "回滚"
		case "scale":
			return "伸缩"
		}
	}
	return ""
}
