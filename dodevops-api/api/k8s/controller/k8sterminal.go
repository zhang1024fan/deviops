package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"dodevops-api/api/k8s/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}


// K8sTerminalController 容器终端控制器
type K8sTerminalController struct {
	service service.IK8sTerminalService
}

func NewK8sTerminalController(db *gorm.DB) *K8sTerminalController {
	return &K8sTerminalController{
		service: service.NewK8sTerminalService(db),
	}
}

// ConnectPodTerminal 连接到Pod终端
// @Summary 连接到Pod终端
// @Description 通过WebSocket连接到指定Pod的终端
// @Tags K8s容器终端
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param podName path string true "Pod名称"
// @Param containerName query string false "容器名称（默认为Pod中第一个容器）"
// @Param command query string false "执行命令（默认为/bin/bash）"
// @Success 101 "Switching Protocols"
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/pods/{podName}/terminal [get]
func (ctrl *K8sTerminalController) ConnectPodTerminal(c *gin.Context) {
	// 参数验证
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil || clusterId <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	if namespaceName == "" {
		result.Failed(c, 400, "命名空间名称不能为空")
		return
	}

	podName := c.Param("podName")
	if podName == "" {
		result.Failed(c, 400, "Pod名称不能为空")
		return
	}

	containerName := c.Query("containerName")
	command := c.Query("command")
	if command == "" {
		command = "/bin/bash"
	}

	// 升级HTTP连接为WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	// 创建K8s WebSocket流
	stream, err := ctrl.service.CreateK8sWebSocketStream(uint(clusterId), namespaceName, podName, containerName, command, conn)
	if err != nil {
		log.Printf("Failed to create K8s WebSocket stream: %v", err)
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Error: %v", err)))
		return
	}
	defer stream.Close()

	// 等待连接关闭
	select {
	case <-stream.Ctx.Done():
		log.Printf("K8s terminal connection closed")
	}
}


// GetPodContainers 获取Pod中的容器列表
// @Summary 获取Pod中的容器列表
// @Description 获取指定Pod中所有容器的名称列表，用于终端连接时选择容器
// @Tags K8s容器终端
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param podName path string true "Pod名称"
// @Success 200 {object} result.Result{data=[]string}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/pods/{podName}/containers [get]
func (ctrl *K8sTerminalController) GetPodContainers(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil || clusterId <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	if namespaceName == "" {
		result.Failed(c, 400, "命名空间名称不能为空")
		return
	}

	podName := c.Param("podName")
	if podName == "" {
		result.Failed(c, 400, "Pod名称不能为空")
		return
	}

	containers, err := ctrl.service.GetPodContainers(uint(clusterId), namespaceName, podName)
	if err != nil {
		result.Failed(c, 500, fmt.Sprintf("获取容器列表失败: %v", err))
		return
	}

	result.Success(c, containers)
}