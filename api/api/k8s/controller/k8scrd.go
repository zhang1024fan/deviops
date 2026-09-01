package controller

import (
	"net/http"
	"strconv"

	"dodevops-api/api/k8s/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type K8sCRDController struct {
	service service.IK8sCRDService
}

func NewK8sCRDController(db *gorm.DB) *K8sCRDController {
	return &K8sCRDController{
		service: service.NewK8sCRDService(db),
	}
}

// GetCRDGroups 获取CRD API Group列表
// @Summary 获取CRD API Group列表
// @Description 获取指定集群中所有CRD所属的API Group
// @Tags K8s CRD管理
// @Accept json
// @Produce json
// @Param id path int true "集群ID"
// @Router /k8s/cluster/{id}/crds/groups [get]
func (ctrl *K8sCRDController) GetCRDGroups(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	groups, err := ctrl.service.GetCRDGroups(uint(clusterId))
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.Success(c, groups)
}

// GetCRDList 获取CRD列表
// @Summary 获取CRD列表
// @Description 获取指定集群中的所有CustomResourceDefinitions
// @Tags K8s CRD管理
// @Accept json
// @Produce json
// @Param id path int true "集群ID"
// @Router /k8s/cluster/{id}/crds [get]
func (ctrl *K8sCRDController) GetCRDList(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	params := map[string]string{}
	for k, v := range c.Request.URL.Query() {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}

	crds, err := ctrl.service.GetCRDList(uint(clusterId), params)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.Success(c, crds)
}

// GetCustomResourceList 获取自定义资源列表
// @Summary 获取自定义资源列表
// @Description 获取指定CRD的自定义资源(CR)列表
// @Tags K8s CRD管理
// @Accept json
// @Produce json
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param crdName path string true "CRD名称 (如 prometheusrules.monitoring.coreos.com)"
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/crds/{crdName}/resources [get]
func (ctrl *K8sCRDController) GetCustomResourceList(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	crdName := c.Param("crdName")

	params := map[string]string{}
	for k, v := range c.Request.URL.Query() {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}

	resList, err := ctrl.service.GetCustomResourceList(uint(clusterId), namespaceName, crdName, params)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.Success(c, resList)
}

// GetCustomResourceDetail 获取自定义资源详情
// @Summary 获取自定义资源详情
// @Description 获取特定的自定义资源实例的详细信息
// @Tags K8s CRD管理
// @Accept json
// @Produce json
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param crdName path string true "CRD名称"
// @Param crName path string true "CR名称"
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/crds/{crdName}/resources/{crName} [get]
func (ctrl *K8sCRDController) GetCustomResourceDetail(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	crdName := c.Param("crdName")
	crName := c.Param("crName")

	resDetail, err := ctrl.service.GetCustomResourceDetail(uint(clusterId), namespaceName, crdName, crName)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.Success(c, resDetail)
}

// CreateCustomResource 创建自定义资源
// @Summary 创建自定义资源
// @Description 基于前端传递的数据创建自定义资源实例
// @Tags K8s CRD管理
// @Accept json
// @Produce json
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param crdName path string true "CRD名称"
// @Param data body map[string]interface{} true "资源的 JSON 对象数据"
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/crds/{crdName}/resources [post]
func (ctrl *K8sCRDController) CreateCustomResource(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	crdName := c.Param("crdName")

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		result.Failed(c, http.StatusBadRequest, "请求数据格式有误")
		return
	}

	created, err := ctrl.service.CreateCustomResource(uint(clusterId), namespaceName, crdName, data)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.Success(c, created)
}

// DeleteCustomResource 删除自定义资源
// @Summary 删除自定义资源
// @Description 删除指定的自定义资源实例
// @Tags K8s CRD管理
// @Accept json
// @Produce json
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param crdName path string true "CRD名称"
// @Param crName path string true "CR名称"
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/crds/{crdName}/resources/{crName} [delete]
func (ctrl *K8sCRDController) DeleteCustomResource(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	crdName := c.Param("crdName")
	crName := c.Param("crName")

	err = ctrl.service.DeleteCustomResource(uint(clusterId), namespaceName, crdName, crName)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.Success(c, "自定义资源删除成功")
}

// GetCustomResourceYaml 获取自定义资源 YAML
// @Summary 获取自定义资源 YAML
// @Description 获取指定的自定义资源实例的 YAML
// @Tags K8s CRD管理
// @Accept json
// @Produce json
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param crdName path string true "CRD名称"
// @Param crName path string true "CR名称"
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/crds/{crdName}/resources/{crName}/yaml [get]
func (ctrl *K8sCRDController) GetCustomResourceYaml(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	crdName := c.Param("crdName")
	crName := c.Param("crName")

	yamlStr, err := ctrl.service.GetCustomResourceYaml(uint(clusterId), namespaceName, crdName, crName)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.Success(c, yamlStr)
}

// UpdateCustomResourceYaml 更新自定义资源 YAML
// @Summary 更新自定义资源 YAML
// @Description 使用提供的 YAML 字符串更新指定的自定义资源实例
// @Tags K8s CRD管理
// @Accept json
// @Produce json
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param crdName path string true "CRD名称"
// @Param crName path string true "CR名称"
// @Param body body map[string]string true "包含 yaml 数据的对象"
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/crds/{crdName}/resources/{crName}/yaml [put]
func (ctrl *K8sCRDController) UpdateCustomResourceYaml(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	crdName := c.Param("crdName")
	crName := c.Param("crName")

	var req struct {
		Yaml string `json:"yaml"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, http.StatusBadRequest, "请求参数解析失败")
		return
	}

	updated, err := ctrl.service.UpdateCustomResourceYaml(uint(clusterId), namespaceName, crdName, crName, req.Yaml)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.Success(c, updated)
}