package controller

import (
	"strconv"

	"dodevops-api/api/k8s/model"
	"dodevops-api/api/k8s/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// K8sStorageController K8s存储管理控制器
type K8sStorageController struct {
	service service.IK8sStorageService
}

func NewK8sStorageController(db *gorm.DB) *K8sStorageController {
	return &K8sStorageController{
		service: service.NewK8sStorageService(db),
	}
}

// ===================== PVC 管理 =====================

// GetPVCs 获取PVC列表
// @Summary 获取PVC列表
// @Description 获取指定命名空间的PVC列表
// @Tags K8s存储管理-PVC
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Success 200 {object} result.Result{data=model.PVCListResponse}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/pvcs [get]
func (ctrl *K8sStorageController) GetPVCs(c *gin.Context) {
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

	ctrl.service.GetPVCs(c, uint(clusterId), namespaceName)
}

// GetPVCDetail 获取PVC详情
// @Summary 获取PVC详情
// @Description 获取指定PVC的详细信息
// @Tags K8s存储管理-PVC
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param pvcName path string true "PVC名称"
// @Success 200 {object} result.Result{data=model.PVCDetail}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/pvcs/{pvcName} [get]
func (ctrl *K8sStorageController) GetPVCDetail(c *gin.Context) {
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

	pvcName := c.Param("pvcName")
	if pvcName == "" {
		result.Failed(c, 400, "PVC名称不能为空")
		return
	}

	ctrl.service.GetPVCDetail(c, uint(clusterId), namespaceName, pvcName)
}

// CreatePVC 创建PVC
// @Summary 创建PVC
// @Description 在指定命名空间中创建新的PVC
// @Tags K8s存储管理-PVC
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param pvc body model.CreatePVCRequest true "PVC配置"
// @Success 200 {object} result.Result{data=model.K8sPersistentVolumeClaim}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/pvcs [post]
func (ctrl *K8sStorageController) CreatePVC(c *gin.Context) {
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

	var req model.CreatePVCRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, 400, "参数验证失败: "+err.Error())
		return
	}

	ctrl.service.CreatePVC(c, uint(clusterId), namespaceName, &req)
}

// UpdatePVC 更新PVC
// @Summary 更新PVC
// @Description 更新指定的PVC配置
// @Tags K8s存储管理-PVC
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param pvcName path string true "PVC名称"
// @Param pvc body model.UpdatePVCRequest true "更新配置"
// @Success 200 {object} result.Result{data=model.K8sPersistentVolumeClaim}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/pvcs/{pvcName} [put]
func (ctrl *K8sStorageController) UpdatePVC(c *gin.Context) {
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

	pvcName := c.Param("pvcName")
	if pvcName == "" {
		result.Failed(c, 400, "PVC名称不能为空")
		return
	}

	var req model.UpdatePVCRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, 400, "参数验证失败: "+err.Error())
		return
	}

	ctrl.service.UpdatePVC(c, uint(clusterId), namespaceName, pvcName, &req)
}

// DeletePVC 删除PVC
// @Summary 删除PVC
// @Description 删除指定的PVC
// @Tags K8s存储管理-PVC
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param pvcName path string true "PVC名称"
// @Success 200 {object} result.Result{data=string}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/pvcs/{pvcName} [delete]
func (ctrl *K8sStorageController) DeletePVC(c *gin.Context) {
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

	pvcName := c.Param("pvcName")
	if pvcName == "" {
		result.Failed(c, 400, "PVC名称不能为空")
		return
	}

	ctrl.service.DeletePVC(c, uint(clusterId), namespaceName, pvcName)
}

// GetPVCYaml 获取PVC的YAML配置
// @Summary 获取PVC的YAML配置
// @Description 获取指定PVC的完整YAML配置
// @Tags K8s存储管理-PVC
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param pvcName path string true "PVC名称"
// @Success 200 {object} result.Result{data=map[string]interface{}}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/pvcs/{pvcName}/yaml [get]
func (ctrl *K8sStorageController) GetPVCYaml(c *gin.Context) {
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

	pvcName := c.Param("pvcName")
	if pvcName == "" {
		result.Failed(c, 400, "PVC名称不能为空")
		return
	}

	ctrl.service.GetPVCYaml(c, uint(clusterId), namespaceName, pvcName)
}

// UpdatePVCYaml 通过YAML更新PVC
// @Summary 通过YAML更新PVC
// @Description 通过提供的YAML内容更新PVC配置
// @Tags K8s存储管理-PVC
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param pvcName path string true "PVC名称"
// @Param yaml body map[string]interface{} true "YAML内容"
// @Success 200 {object} result.Result{data=string}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/pvcs/{pvcName}/yaml [put]
func (ctrl *K8sStorageController) UpdatePVCYaml(c *gin.Context) {
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

	pvcName := c.Param("pvcName")
	if pvcName == "" {
		result.Failed(c, 400, "PVC名称不能为空")
		return
	}

	var yamlData map[string]interface{}
	if err := c.ShouldBindJSON(&yamlData); err != nil {
		result.Failed(c, 400, "YAML内容格式错误: "+err.Error())
		return
	}

	ctrl.service.UpdatePVCYaml(c, uint(clusterId), namespaceName, pvcName, yamlData)
}

// ===================== PV 管理 =====================

// GetPVs 获取PV列表
// @Summary 获取PV列表
// @Description 获取集群中的PV列表
// @Tags K8s存储管理-PV
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Success 200 {object} result.Result{data=model.PVListResponse}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/pvs [get]
func (ctrl *K8sStorageController) GetPVs(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil || clusterId <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	ctrl.service.GetPVs(c, uint(clusterId))
}

// GetPVDetail 获取PV详情
// @Summary 获取PV详情
// @Description 获取指定PV的详细信息
// @Tags K8s存储管理-PV
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param pvName path string true "PV名称"
// @Success 200 {object} result.Result{data=model.PVDetail}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/pvs/{pvName} [get]
func (ctrl *K8sStorageController) GetPVDetail(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil || clusterId <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	pvName := c.Param("pvName")
	if pvName == "" {
		result.Failed(c, 400, "PV名称不能为空")
		return
	}

	ctrl.service.GetPVDetail(c, uint(clusterId), pvName)
}

// CreatePV 创建PV
// @Summary 创建PV
// @Description 在集群中创建新的PV
// @Tags K8s存储管理-PV
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param pv body model.CreatePVRequest true "PV配置"
// @Success 200 {object} result.Result{data=model.K8sPersistentVolume}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/pvs [post]
func (ctrl *K8sStorageController) CreatePV(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil || clusterId <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	var req model.CreatePVRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, 400, "参数验证失败: "+err.Error())
		return
	}

	ctrl.service.CreatePV(c, uint(clusterId), &req)
}

// UpdatePV 更新PV
// @Summary 更新PV
// @Description 更新指定的PV配置
// @Tags K8s存储管理-PV
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param pvName path string true "PV名称"
// @Param pv body model.UpdatePVRequest true "更新配置"
// @Success 200 {object} result.Result{data=model.K8sPersistentVolume}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/pvs/{pvName} [put]
func (ctrl *K8sStorageController) UpdatePV(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil || clusterId <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	pvName := c.Param("pvName")
	if pvName == "" {
		result.Failed(c, 400, "PV名称不能为空")
		return
	}

	var req model.UpdatePVRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, 400, "参数验证失败: "+err.Error())
		return
	}

	ctrl.service.UpdatePV(c, uint(clusterId), pvName, &req)
}

// DeletePV 删除PV
// @Summary 删除PV
// @Description 删除指定的PV
// @Tags K8s存储管理-PV
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param pvName path string true "PV名称"
// @Success 200 {object} result.Result{data=string}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/pvs/{pvName} [delete]
func (ctrl *K8sStorageController) DeletePV(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil || clusterId <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	pvName := c.Param("pvName")
	if pvName == "" {
		result.Failed(c, 400, "PV名称不能为空")
		return
	}

	ctrl.service.DeletePV(c, uint(clusterId), pvName)
}

// GetPVYaml 获取PV的YAML配置
// @Summary 获取PV的YAML配置
// @Description 获取指定PV的完整YAML配置
// @Tags K8s存储管理-PV
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param pvName path string true "PV名称"
// @Success 200 {object} result.Result{data=map[string]interface{}}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/pvs/{pvName}/yaml [get]
func (ctrl *K8sStorageController) GetPVYaml(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil || clusterId <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	pvName := c.Param("pvName")
	if pvName == "" {
		result.Failed(c, 400, "PV名称不能为空")
		return
	}

	ctrl.service.GetPVYaml(c, uint(clusterId), pvName)
}

// UpdatePVYaml 通过YAML更新PV
// @Summary 通过YAML更新PV
// @Description 通过提供的YAML内容更新PV配置
// @Tags K8s存储管理-PV
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param pvName path string true "PV名称"
// @Param yaml body map[string]interface{} true "YAML内容"
// @Success 200 {object} result.Result{data=string}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/pvs/{pvName}/yaml [put]
func (ctrl *K8sStorageController) UpdatePVYaml(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil || clusterId <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	pvName := c.Param("pvName")
	if pvName == "" {
		result.Failed(c, 400, "PV名称不能为空")
		return
	}

	var yamlData map[string]interface{}
	if err := c.ShouldBindJSON(&yamlData); err != nil {
		result.Failed(c, 400, "YAML内容格式错误: "+err.Error())
		return
	}

	ctrl.service.UpdatePVYaml(c, uint(clusterId), pvName, yamlData)
}

// ===================== StorageClass 管理 =====================

// GetStorageClasses 获取存储类列表
// @Summary 获取存储类列表
// @Description 获取集群中的存储类列表
// @Tags K8s存储管理-StorageClass
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Success 200 {object} result.Result{data=model.StorageClassListResponse}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/storageclasses [get]
func (ctrl *K8sStorageController) GetStorageClasses(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil || clusterId <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	ctrl.service.GetStorageClasses(c, uint(clusterId))
}

// GetStorageClassDetail 获取存储类详情
// @Summary 获取存储类详情
// @Description 获取指定存储类的详细信息
// @Tags K8s存储管理-StorageClass
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param storageClassName path string true "存储类名称"
// @Success 200 {object} result.Result{data=model.StorageClassDetail}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/storageclasses/{storageClassName} [get]
func (ctrl *K8sStorageController) GetStorageClassDetail(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil || clusterId <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	storageClassName := c.Param("storageClassName")
	if storageClassName == "" {
		result.Failed(c, 400, "存储类名称不能为空")
		return
	}

	ctrl.service.GetStorageClassDetail(c, uint(clusterId), storageClassName)
}

// CreateStorageClass 创建存储类
// @Summary 创建存储类
// @Description 在集群中创建新的存储类
// @Tags K8s存储管理-StorageClass
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param storageClass body model.CreateStorageClassRequest true "存储类配置"
// @Success 200 {object} result.Result{data=model.K8sStorageClass}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/storageclasses [post]
func (ctrl *K8sStorageController) CreateStorageClass(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil || clusterId <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	var req model.CreateStorageClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, 400, "参数验证失败: "+err.Error())
		return
	}

	ctrl.service.CreateStorageClass(c, uint(clusterId), &req)
}

// UpdateStorageClass 更新存储类
// @Summary 更新存储类
// @Description 更新指定的存储类配置
// @Tags K8s存储管理-StorageClass
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param storageClassName path string true "存储类名称"
// @Param storageClass body model.UpdateStorageClassRequest true "更新配置"
// @Success 200 {object} result.Result{data=model.K8sStorageClass}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/storageclasses/{storageClassName} [put]
func (ctrl *K8sStorageController) UpdateStorageClass(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil || clusterId <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	storageClassName := c.Param("storageClassName")
	if storageClassName == "" {
		result.Failed(c, 400, "存储类名称不能为空")
		return
	}

	var req model.UpdateStorageClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, 400, "参数验证失败: "+err.Error())
		return
	}

	ctrl.service.UpdateStorageClass(c, uint(clusterId), storageClassName, &req)
}

// DeleteStorageClass 删除存储类
// @Summary 删除存储类
// @Description 删除指定的存储类
// @Tags K8s存储管理-StorageClass
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param storageClassName path string true "存储类名称"
// @Success 200 {object} result.Result{data=string}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/storageclasses/{storageClassName} [delete]
func (ctrl *K8sStorageController) DeleteStorageClass(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil || clusterId <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	storageClassName := c.Param("storageClassName")
	if storageClassName == "" {
		result.Failed(c, 400, "存储类名称不能为空")
		return
	}

	ctrl.service.DeleteStorageClass(c, uint(clusterId), storageClassName)
}

// GetStorageClassYaml 获取存储类的YAML配置
// @Summary 获取存储类的YAML配置
// @Description 获取指定存储类的完整YAML配置
// @Tags K8s存储管理-StorageClass
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param storageClassName path string true "存储类名称"
// @Success 200 {object} result.Result{data=map[string]interface{}}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/storageclasses/{storageClassName}/yaml [get]
func (ctrl *K8sStorageController) GetStorageClassYaml(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil || clusterId <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	storageClassName := c.Param("storageClassName")
	if storageClassName == "" {
		result.Failed(c, 400, "存储类名称不能为空")
		return
	}

	ctrl.service.GetStorageClassYaml(c, uint(clusterId), storageClassName)
}

// UpdateStorageClassYaml 通过YAML更新存储类
// @Summary 通过YAML更新存储类
// @Description 通过提供的YAML内容更新存储类配置
// @Tags K8s存储管理-StorageClass
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param storageClassName path string true "存储类名称"
// @Param yaml body map[string]interface{} true "YAML内容"
// @Success 200 {object} result.Result{data=string}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/storageclasses/{storageClassName}/yaml [put]
func (ctrl *K8sStorageController) UpdateStorageClassYaml(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil || clusterId <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	storageClassName := c.Param("storageClassName")
	if storageClassName == "" {
		result.Failed(c, 400, "存储类名称不能为空")
		return
	}

	var yamlData map[string]interface{}
	if err := c.ShouldBindJSON(&yamlData); err != nil {
		result.Failed(c, 400, "YAML内容格式错误: "+err.Error())
		return
	}

	ctrl.service.UpdateStorageClassYaml(c, uint(clusterId), storageClassName, yamlData)
}