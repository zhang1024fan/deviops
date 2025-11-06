package service

import (
	"context"
	"fmt"
	"io"
	"sync"

	"dodevops-api/api/k8s/dao"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

// IK8sTerminalService 容器终端服务接口
type IK8sTerminalService interface {
	CreateK8sWebSocketStream(clusterId uint, namespaceName, podName, containerName, command string, conn *websocket.Conn) (*K8sWebSocketStream, error)
	GetPodContainers(clusterId uint, namespaceName, podName string) ([]string, error)
}

// K8sMessage WebSocket消息结构
type K8sMessage struct {
	Operation string      `json:"operation"`
	Data      interface{} `json:"data"`
	Cols      int         `json:"cols,omitempty"`
	Rows      int         `json:"rows,omitempty"`
}

// K8sWebSocketStream K8s WebSocket流处理
type K8sWebSocketStream struct {
	sync.RWMutex
	Conn     *websocket.Conn
	executor remotecommand.Executor
	Ctx      context.Context
	cancel   context.CancelFunc
	closed   bool
	reader   *io.PipeReader
	writer   *io.PipeWriter
}

// terminalConn 实现io.ReadWriter接口用于K8s executor
type terminalConn struct {
	stream *K8sWebSocketStream
}

func (tc *terminalConn) Read(p []byte) (n int, err error) {
	return tc.stream.reader.Read(p)
}

func (tc *terminalConn) Write(p []byte) (n int, err error) {
	return tc.stream.WriteToWebSocket(p)
}

// Close 关闭WebSocket流
func (kws *K8sWebSocketStream) Close() error {
	kws.Lock()
	defer kws.Unlock()
	
	if kws.closed {
		return nil
	}
	
	kws.closed = true
	if kws.cancel != nil {
		kws.cancel()
	}
	// 不要立即关闭管道，让K8s executor自然结束
	return nil
}

// IsClosed 检查是否已关闭
func (kws *K8sWebSocketStream) IsClosed() bool {
	kws.RLock()
	defer kws.RUnlock()
	return kws.closed
}

// WriteToWebSocket 写入数据到WebSocket
func (kws *K8sWebSocketStream) WriteToWebSocket(p []byte) (n int, err error) {
	if kws.IsClosed() {
		return 0, io.EOF
	}
	
	message := K8sMessage{
		Operation: "stdout",
		Data:      string(p),
	}
	
	if err = kws.Conn.WriteJSON(message); err != nil {
		go kws.Close() // 异步关闭，避免死锁
		return 0, err
	}
	
	return len(p), nil
}

// ReadFromWebSocket 从WebSocket读取数据
func (kws *K8sWebSocketStream) ReadFromWebSocket() {
	defer func() {
		kws.Close()
		// 关闭写入管道，通知K8s executor结束
		if kws.writer != nil {
			kws.writer.Close()
		}
	}()
	
	for {
		if kws.IsClosed() {
			return
		}
		
		var message K8sMessage
		err := kws.Conn.ReadJSON(&message)
		if err != nil {
			return
		}
		
		switch message.Operation {
		case "stdin":
			if kws.writer != nil && !kws.IsClosed() {
				var data string
				if str, ok := message.Data.(string); ok {
					data = str
				} else {
					continue
				}
				
				_, err := kws.writer.Write([]byte(data))
				if err != nil {
					return
				}
			}
		case "resize":
			// 处理终端大小调整，暂时忽略以避免问题
			continue
		default:
			// 忽略未知操作
			continue
		}
	}
}

// K8sTerminalService 容器终端服务实现
type K8sTerminalService struct {
	dao *dao.KubeClusterDao
}

func NewK8sTerminalService(db *gorm.DB) IK8sTerminalService {
	return &K8sTerminalService{
		dao: dao.NewKubeClusterDao(db),
	}
}

// CreateK8sWebSocketStream 创建K8s WebSocket流
func (s *K8sTerminalService) CreateK8sWebSocketStream(clusterId uint, namespaceName, podName, containerName, command string, conn *websocket.Conn) (*K8sWebSocketStream, error) {
	// 获取集群信息
	cluster, err := s.dao.GetByID(clusterId)
	if err != nil {
		return nil, fmt.Errorf("获取集群信息失败: %v", err)
	}

	// 创建Kubernetes客户端
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Credential))
	if err != nil {
		return nil, fmt.Errorf("创建Kubernetes配置失败: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("创建Kubernetes客户端失败: %v", err)
	}

	// 获取Pod信息，如果没有指定容器名称，使用第一个容器
	if containerName == "" {
		pod, err := clientset.CoreV1().Pods(namespaceName).Get(context.Background(), podName, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("获取Pod信息失败: %v", err)
		}
		if len(pod.Spec.Containers) == 0 {
			return nil, fmt.Errorf("Pod中没有找到容器")
		}
		containerName = pod.Spec.Containers[0].Name
	}

	// 创建exec请求
	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespaceName).
		SubResource("exec").
		Param("container", containerName).
		Param("command", command).
		Param("stdin", "true").
		Param("stdout", "true").
		Param("stderr", "true").
		Param("tty", "true")

	// 创建executor
	executor, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return nil, fmt.Errorf("创建executor失败: %v", err)
	}

	// 创建长期运行的上下文
	ctx, cancel := context.WithCancel(context.Background())

	// 创建管道
	reader, writer := io.Pipe()

	// 创建流对象
	stream := &K8sWebSocketStream{
		Conn:     conn,
		executor: executor,
		Ctx:      ctx,
		cancel:   cancel,
		reader:   reader,
		writer:   writer,
	}

	// 创建终端连接
	termConn := &terminalConn{stream: stream}

	// 启动WebSocket读取goroutine
	go stream.ReadFromWebSocket()

	// 启动K8s executor goroutine
	go func() {
		defer func() {
			cancel()
			// 关闭读取管道
			if reader != nil {
				reader.Close()
			}
		}()

		executor.StreamWithContext(ctx, remotecommand.StreamOptions{
			Stdin:  termConn,
			Stdout: termConn,
			Stderr: termConn,
			Tty:    true,
		})
	}()

	return stream, nil
}

// GetPodContainers 获取Pod中的容器列表
func (s *K8sTerminalService) GetPodContainers(clusterId uint, namespaceName, podName string) ([]string, error) {
	// 获取集群信息
	cluster, err := s.dao.GetByID(clusterId)
	if err != nil {
		return nil, fmt.Errorf("获取集群信息失败: %v", err)
	}

	// 创建Kubernetes客户端
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Credential))
	if err != nil {
		return nil, fmt.Errorf("创建Kubernetes配置失败: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("创建Kubernetes客户端失败: %v", err)
	}

	// 获取Pod信息
	pod, err := clientset.CoreV1().Pods(namespaceName).Get(context.Background(), podName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("获取Pod信息失败: %v", err)
	}

	// 提取容器名称
	var containers []string
	for _, container := range pod.Spec.Containers {
		containers = append(containers, container.Name)
	}
	for _, container := range pod.Spec.InitContainers {
		containers = append(containers, container.Name+" (init)")
	}

	return containers, nil
}