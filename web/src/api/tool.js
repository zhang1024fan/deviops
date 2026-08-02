/**
 * 导航工具API接口管理
 *
 * @author xiaoRui
 */

import request from '@/utils/request'

/**
 * 创建导航工具
 * @param {Object} data - 导航工具信息
 * @param {string} data.title - 导航标题，1-100字符
 * @param {string} data.icon - 图标URL
 * @param {string} data.link - 链接地址，必须是有效URL
 * @param {number} data.sort - 排序，默认0
 */
export function CreateTool(data) {
  return request({
    url: '/api/v1/tool',
    method: 'post',
    data,
    headers: {
      'Content-Type': 'application/json'
    }
  })
}

/**
 * 获取导航工具详情
 * @param {number} id - 工具ID
 */
export function GetToolById(id) {
  return request({
    url: `/api/v1/tool/${id}`,
    method: 'get',
    headers: {
      'Accept': 'application/json'
    }
  })
}

/**
 * 更新导航工具
 * @param {Object} data - 导航工具信息
 * @param {number} data.id - 工具ID
 * @param {string} data.title - 导航标题，1-100字符
 * @param {string} data.icon - 图标URL
 * @param {string} data.link - 链接地址，必须是有效URL
 * @param {number} data.sort - 排序
 */
export function UpdateTool(data) {
  return request({
    url: '/api/v1/tool',
    method: 'put',
    data,
    headers: {
      'Content-Type': 'application/json'
    }
  })
}

/**
 * 删除导航工具
 * @param {number} id - 工具ID
 */
export function DeleteTool(id) {
  return request({
    url: `/api/v1/tool/${id}`,
    method: 'delete',
    headers: {
      'Accept': 'application/json'
    }
  })
}

/**
 * 获取导航工具列表（分页）
 * @param {Object} params - 查询参数
 * @param {string} params.title - 标题模糊搜索
 * @param {number} params.pageNum - 页码，默认1
 * @param {number} params.pageSize - 每页数量，默认10
 */
export function GetToolList(params) {
  return request({
    url: '/api/v1/tool/list',
    method: 'get',
    params,
    headers: {
      'Accept': 'application/json'
    }
  })
}

/**
 * 获取所有导航工具（不分页）
 */
export function GetAllTools() {
  return request({
    url: '/api/v1/tool/all',
    method: 'get',
    headers: {
      'Accept': 'application/json'
    }
  })
}

/**
 * 上传图标
 * @param {FormData} formData - 包含file字段的FormData对象
 */
export function UploadIcon(formData) {
  return request({
    url: '/api/v1/upload',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// ==================== 运维工具箱API ====================

/**
 * 获取可部署服务列表
 */
export function getServicesList() {
  return request({
    url: '/api/v1/tool/services',
    method: 'get'
  })
}

/**
 * 获取服务详情
 * @param {string} serviceId - 服务ID
 */
export function getServiceDetail(serviceId) {
  return request({
    url: `/api/v1/tool/services/${serviceId}`,
    method: 'get'
  })
}

/**
 * 创建部署任务
 * @param {Object} data - 部署参数
 * @param {string} data.serviceId - 服务ID
 * @param {string} data.version - 版本ID
 * @param {number} data.hostId - 主机ID
 * @param {string} data.installDir - 安装目录
 * @param {Object} data.envVars - 环境变量
 * @param {boolean} data.autoStart - 是否自动启动
 */
export function createDeploy(data) {
  return request({
    url: '/api/v1/tool/deploy',
    method: 'post',
    data
  })
}

/**
 * 获取部署历史列表
 * @param {Object} params - 查询参数
 * @param {string} params.serviceName - 服务名称
 * @param {number} params.status - 状态
 * @param {number} params.pageNum - 页码
 * @param {number} params.pageSize - 每页数量
 */
export function getDeployList(params) {
  return request({
    url: '/api/v1/tool/deploy/list',
    method: 'get',
    params
  })
}

/**
 * 获取部署状态
 * @param {number} id - 部署ID
 */
export function getDeployStatus(id) {
  return request({
    url: `/api/v1/tool/deploy/${id}/status`,
    method: 'get'
  })
}

/**
 * 卸载服务
 * @param {number} id - 部署ID
 */
export function deleteDeploy(id) {
  return request({
    url: `/api/v1/tool/deploy/${id}`,
    method: 'delete'
  })
}

// ==================== 运维知识库API ====================

/**
 * 创建知识
 * @param {Object} data - 知识信息
 * @param {string} data.title - 标题
 * @param {string} data.category - 分类
 * @param {string} data.content - Markdown内容
 * @param {string} data.tags - 标签(JSON数组)
 * @param {number} data.status - 状态:1->已发布,2->草稿
 */
export function createKnowledge(data) {
  return request({
    url: '/api/v1/knowledge',
    method: 'post',
    data
  })
}

/**
 * 获取知识详情
 * @param {number} id - 知识ID
 */
export function getKnowledgeById(id) {
  return request({
    url: `/api/v1/knowledge/${id}`,
    method: 'get'
  })
}

/**
 * 更新知识
 * @param {Object} data - 知识信息
 * @param {number} data.id - 知识ID
 * @param {string} data.title - 标题
 * @param {string} data.category - 分类
 * @param {string} data.content - Markdown内容
 * @param {string} data.tags - 标签(JSON数组)
 * @param {number} data.status - 状态:1->已发布,2->草稿
 */
export function updateKnowledge(data) {
  return request({
    url: '/api/v1/knowledge',
    method: 'put',
    data
  })
}

/**
 * 删除知识
 * @param {number} id - 知识ID
 */
export function deleteKnowledge(id) {
  return request({
    url: `/api/v1/knowledge/${id}`,
    method: 'delete'
  })
}

/**
 * 获取知识列表(分页)
 * @param {Object} params - 查询参数
 * @param {string} params.title - 标题(模糊查询)
 * @param {string} params.category - 分类
 * @param {number} params.status - 状态
 * @param {number} params.pageNum - 页码
 * @param {number} params.pageSize - 每页数量
 */
export function getKnowledgeList(params) {
  return request({
    url: '/api/v1/knowledge/list',
    method: 'get',
    params
  })
}

/**
 * 获取所有分类
 */
export function getKnowledgeCategories() {
  return request({
    url: '/api/v1/knowledge/categories',
    method: 'get'
  })
}

// ==================== 知识分类管理API ====================

/**
 * 创建分类
 * @param {Object} data - 分类信息
 * @param {string} data.name - 分类名称
 * @param {number} data.sort - 排序
 * @param {string} data.description - 分类描述
 */
export function createCategory(data) {
  return request({
    url: '/api/v1/knowledge/category',
    method: 'post',
    data
  })
}

/**
 * 获取分类详情
 * @param {number} id - 分类ID
 */
export function getCategoryById(id) {
  return request({
    url: `/api/v1/knowledge/category/${id}`,
    method: 'get'
  })
}

/**
 * 更新分类
 * @param {Object} data - 分类信息
 * @param {number} data.id - 分类ID
 * @param {string} data.name - 分类名称
 * @param {number} data.sort - 排序
 * @param {string} data.description - 分类描述
 */
export function updateCategory(data) {
  return request({
    url: '/api/v1/knowledge/category',
    method: 'put',
    data
  })
}

/**
 * 删除分类
 * @param {number} id - 分类ID
 */
export function deleteCategory(id) {
  return request({
    url: `/api/v1/knowledge/category/${id}`,
    method: 'delete'
  })knowledge
}

/**
 * 获取分类列表
 */
export function getCategoryList() {
  return request({
    url: '/api/v1/knowledge/category/list',
    method: 'get'
  })
}
