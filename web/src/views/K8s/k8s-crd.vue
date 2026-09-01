<script setup>
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus,
  Edit,
  Delete,
  Search,
  Refresh,
  Coin,
  Setting,
  View
} from '@element-plus/icons-vue'
import k8sApi from '@/api/k8s'
import PodYamlDialog from './pods/PodYamlDialog.vue'
import ClusterSelector from './pods/ClusterSelector.vue'
import NamespaceSelector from './pods/NamespaceSelector.vue'
import yaml from 'js-yaml'

// 基础状态
const loading = ref(false)
const searchKeyword = ref('')
const searchLabels = ref('')
const selectedClusterId = ref('')
const queryParams = reactive({
  namespace: 'default'
})

// 数据状态
const crdGroupList = ref([])    // 定义的 CRD Group 列表
const selectedGroup = ref('')   // 当前选中的 Group
const crdList = ref([])         // Group 下的 CRD 列表
const selectedCrdName = ref('') // 用户当前所选的 CRD 类型的名字 (例如 prometheusrules.monitoring.coreos.com)
const crList = ref([])          // 用户选中的 CRD 的实例化资源列表 (CRs)

// 表单/对话框状态
const yamlDialogVisible = ref(false)
const currentYaml = ref('')
const currentCr = ref({})
const isCreate = ref(false)
const dialogEditable = ref(true)

// 分页与搜索状态
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

const handleSearch = () => {
  currentPage.value = 1
  loadCrList()
}

const handleSizeChange = (val) => {
  pageSize.value = val
  currentPage.value = 1
  loadCrList()
}

const handleCurrentChange = (val) => {
  currentPage.value = val
  loadCrList()
}

// 获取系统所有的 CRD Groups
const loadCrdGroups = async () => {
  if (!selectedClusterId.value) return
  loading.value = true
  try {
    const res = await k8sApi.getCRDGroups(selectedClusterId.value)
    const resData = res.data || res
    if (resData.code === 200 || resData.success) {
       crdGroupList.value = resData.data || []
       // 默认选择第一个
       if (crdGroupList.value.length > 0 && !selectedGroup.value) {
         const firstGroup = crdGroupList.value[0]
         selectedGroup.value = firstGroup.name || firstGroup
       }
    } else {
       ElMessage.error(resData.message || '获取 CRD Groups 失败')
    }
  } catch (err) {
    console.error('获取 CRD Groups 异常:', err)
  } finally {
    loading.value = false
  }
}

// 获取某个 Group 下的 CRD
const loadCrdList = async () => {
  if (!selectedClusterId.value || !selectedGroup.value) return
  loading.value = true
  try {
    const res = await k8sApi.getCRDList(selectedClusterId.value, { group: selectedGroup.value })
    const resData = res.data || res
    
    // 支持不同的后端数据结构容错
    if (resData.code === 200 || resData.success) {
       crdList.value = resData.data || []
       // 默认选择第一个
       if (crdList.value.length > 0 && !selectedCrdName.value) {
         const firstCrd = crdList.value[0]
         selectedCrdName.value = firstCrd.name || firstCrd.metadata?.name || firstCrd
       }
    } else {
       ElMessage.error(resData.message || '获取 CRD 列表失败')
    }
  } catch (err) {
    console.error('获取 CRD 列表异常:', err)
  } finally {
    loading.value = false
  }
}

// 获取选中 CRD 的具体资源 (CR)
const loadCrList = async () => {
  if (!selectedClusterId.value || !queryParams.namespace || !selectedCrdName.value) {
    crList.value = []
    total.value = 0
    return
  }
  loading.value = true
  try {
    // 处理多标签格式，支持空格或中英文逗号分隔，统一转换成 K8s 兼容的英文逗号分隔符
    const processedLabels = searchLabels.value 
      ? searchLabels.value.trim().replace(/[\s，,]+/g, ',') 
      : undefined

    const params = {
      page: currentPage.value,
      pageSize: pageSize.value,
      name: searchKeyword.value || undefined,
      keyword: searchKeyword.value || undefined,
      labels: processedLabels
    }

    const res = await k8sApi.getCustomResourceList(
      selectedClusterId.value, 
      queryParams.namespace, 
      selectedCrdName.value,
      params
    )
    const resData = res.data || res
    
    if (resData.code === 200 || resData.success) {
      const data = resData.data
      if (data && (Array.isArray(data.items) || Array.isArray(data.list))) {
          crList.value = data.items || data.list || []
          total.value = data.total || 0
      } else if (Array.isArray(data)) {
          crList.value = data
          total.value = data.length || 0
      } else {
          crList.value = []
          total.value = 0
      }
    } else {
      crList.value = []
      total.value = 0
      ElMessage.error(resData.message || `获取 ${selectedCrdName.value} 资源失败`)
    }
  } catch (err) {
    crList.value = []
    total.value = 0
    console.error('获取 CR 失败:', err)
  } finally {
    loading.value = false
  }
}

// 刷新
const handleRefresh = () => {
  if (selectedCrdName.value) {
    loadCrList() // CRs
  } else if (selectedGroup.value) {
    loadCrdList() // CRDs
  } else {
    loadCrdGroups() // Groups
  }
}

// 创建 CR
const handleCreateCr = () => {
  if (!selectedClusterId.value) {
    ElMessage.warning('请先选择集群')
    return
  }
  if (!queryParams.namespace) {
    ElMessage.warning('请先选择命名空间')
    return
  }
  if (!selectedCrdName.value) {
    ElMessage.warning('请先从列表中选择一个 CRD 类型')
    return
  }

  isCreate.value = true
  dialogEditable.value = true
  currentCr.value = { name: 'new-cr' } // 标识为主键
  
  // 提取一下简写Kind (如果没有可以通过后端拿)
  const shortKind = selectedCrdName.value.split('.')[0] 
  // 简易 YAML 模板
  currentYaml.value = `apiVersion: ${selectedCrdName.value.indexOf('.') > -1 ? selectedCrdName.value.slice(selectedCrdName.value.indexOf('.') + 1) + '/v1' : 'v1'}
kind: ${shortKind.charAt(0).toUpperCase() + shortKind.slice(1)}
metadata:
  name: new-custom-resource
  namespace: ${queryParams.namespace}
spec:
  # 在此处填写您的CR配置
`
  yamlDialogVisible.value = true
}

// 查看 CR 的 YAML
const handleViewCrYaml = async (row) => {
  const rowName = row?.name || row?.metadata?.name
  if (!rowName) return
  
  try {
    loading.value = true
    const res = await k8sApi.getCustomResourceYaml(
      selectedClusterId.value, 
      queryParams.namespace, 
      selectedCrdName.value, 
      rowName
    )
    const resData = res.data || res
    if (resData.code === 200 || resData.success) {
      let yamlContent = resData.data
      if (typeof yamlContent === 'object' && yamlContent !== null && yamlContent.yaml) {
        yamlContent = yamlContent.yaml
      } else if (typeof yamlContent === 'object' && yamlContent !== null) {
        yamlContent = JSON.stringify(yamlContent, null, 2)
      }
      
      isCreate.value = false
      dialogEditable.value = false
      currentYaml.value = String(yamlContent)
      currentCr.value = { name: rowName }
      yamlDialogVisible.value = true
    } else {
      ElMessage.error(resData.message || '获取 YAML 失败')
    }
  } catch (error) {
    console.error('获取 YAML 异常', error)
  } finally {
    loading.value = false
  }
}

// 编辑 CR
const handleEditCrYaml = async (row) => {
  const rowName = row?.name || row?.metadata?.name
  if (!rowName) return
  
  try {
    loading.value = true
    const res = await k8sApi.getCustomResourceYaml(
      selectedClusterId.value, 
      queryParams.namespace, 
      selectedCrdName.value, 
      rowName
    )
    const resData = res.data || res
    if (resData.code === 200 || resData.success) {
      let yamlContent = resData.data
      if (typeof yamlContent === 'object' && yamlContent !== null && yamlContent.yaml) {
        yamlContent = yamlContent.yaml
      } else if (typeof yamlContent === 'object' && yamlContent !== null) {
        yamlContent = JSON.stringify(yamlContent, null, 2)
      }
      
      // 清除无用的系统注解
      if (typeof yamlContent === 'string') {
        try {
          const doc = yaml.load(yamlContent)
          if (doc && doc.metadata && doc.metadata.annotations) {
            delete doc.metadata.annotations['kubectl.kubernetes.io/last-applied-configuration']
            if (Object.keys(doc.metadata.annotations).length === 0) {
              delete doc.metadata.annotations
            }
          }
          if (doc) yamlContent = yaml.dump(doc)
        } catch(e) {
          console.warn('清理 CR YAML 注解失败', e)
        }
      }
      
      isCreate.value = false
      dialogEditable.value = true
      currentYaml.value = String(yamlContent)
      currentCr.value = { name: rowName }
      yamlDialogVisible.value = true
    } else {
      ElMessage.error(resData.message || '获取 YAML 失败')
    }
  } catch (error) {
    console.error('获取 YAML 异常', error)
  } finally {
    loading.value = false
  }
}

// 删除 CR
const handleDeleteCr = async (row) => {
  const rowName = row?.name || row?.metadata?.name
  if (!rowName) return

  try {
    await ElMessageBox.confirm(
      `确认删除资源实例 "${rowName}" 吗？此操作不可逆。`,
      '删除确认',
      { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' }
    )
    
    loading.value = true
    const res = await k8sApi.deleteCustomResource(
      selectedClusterId.value, 
      queryParams.namespace, 
      selectedCrdName.value, 
      rowName
    )
    const resData = res.data || res
    
    if (resData.code === 200 || resData.success) {
      ElMessage.success('自定义资源删除成功')
      loadCrList()
    } else {
      ElMessage.error(resData.message || '删除失败')
    }
  } catch (e) {
     if (e !== 'cancel') {
        console.error('删除异常', e)
        ElMessage.error('删除异常')
     }
  } finally {
    loading.value = false
  }
}

// 保存 YAML（创建或更新）
const handleSaveYaml = async (payload) => {
  try {
    loading.value = true
    let res
    if (isCreate.value) {
      res = await k8sApi.createCustomResource(
        selectedClusterId.value, 
        queryParams.namespace, 
        selectedCrdName.value,
        { yamlContent: payload.yamlContent }
      )
    } else {
      res = await k8sApi.updateCustomResourceYaml(
        selectedClusterId.value, 
        queryParams.namespace, 
        selectedCrdName.value,
        currentCr.value.name,
        payload.yamlContent
      )
    }
    
    const resData = res.data || res
    if (resData.code === 200 || resData.success) {
       ElMessage.success(isCreate.value ? '资源创建成功' : '资源更新成功')
       yamlDialogVisible.value = false
       loadCrList()
    } else {
       ElMessage.error(resData.message || (isCreate.value ? '创建资源失败' : '变更资源失败'))
    }
  } catch(e) {
    console.error('保存失败', e)
    ElMessage.error('保存操作异常')
  } finally {
    loading.value = false
  }
}

// 格式化标签工具
const extractLabels = (row) => {
  return row?.metadata?.labels || row?.labels || {}
}

const extractName = (row) => {
  return row?.metadata?.name || row?.name || '-'
}

const extractTime = (row) => {
  const time = row?.metadata?.creationTimestamp || row?.createdTime
  return time ? new Date(time).toLocaleString() : '-'
}

// 集群和命名空间变化监听
watch(
  [selectedClusterId, () => queryParams.namespace],
  ([clusterId, namespace], [oldClusterId, oldNamespace] = []) => {
    if (clusterId && clusterId !== oldClusterId) {
      loadCrdGroups() // 切集群时拉取一次 Group 清单
      crdList.value = []
      selectedGroup.value = ''
      selectedCrdName.value = ''
      crList.value = []
      total.value = 0
    } else if (clusterId && namespace !== oldNamespace && selectedCrdName.value) {
      loadCrList() // 仅切换命名空间时刷新具体的自定义资源实例列表
    }
  },
  { immediate: true }
)

// Group 变化加载对应的 CRD 列表
watch(selectedGroup, (val) => {
  selectedCrdName.value = ''
  crdList.value = []
  crList.value = []
  total.value = 0
  if (val) {
    loadCrdList()
  }
})

// Selector 变化联动获取新列表
watch(selectedCrdName, (val) => {
  if (val) {
    loadCrList()
  } else {
    crList.value = []
    total.value = 0
  }
})

</script>

<template>
  <div class="k8s-crd-management">
    <el-card shadow="hover" class="config-card">
      <template #header>
        <div class="card-header">
          <span class="title">K8s CRD管理</span>
          <div class="header-actions">
            <ClusterSelector v-model="selectedClusterId" />
            <NamespaceSelector 
              v-model="queryParams.namespace" 
              :cluster-id="selectedClusterId" 
            />
          </div>
        </div>
      </template>

      <!-- 控制台及搜索 -->
      <div class="search-section">
        <el-form :inline="true" class="search-form">
          <el-form-item label="API Group">
            <el-select 
              v-model="selectedGroup" 
              filterable 
              clearable 
              placeholder="请选择组 (Group)" 
              size="small" 
              style="width: 200px;"
            >
              <!-- 兼容纯字符串或对象格式 -->
              <el-option 
                v-for="group in crdGroupList" 
                :key="group.name || group"
                :label="group.name || group" 
                :value="group.name || group"
              />
            </el-select>
          </el-form-item>

          <el-form-item label="自定义资源类型">
            <el-select 
              v-model="selectedCrdName" 
              filterable 
              clearable 
              placeholder="请输入或选择 CRD" 
              size="small" 
              style="width: 260px;"
              :disabled="!selectedGroup"
            >
              <!-- 兼容纯字符串数组或Object对象数组 -->
              <el-option 
                v-for="crd in crdList" 
                :key="crd.name || crd.metadata?.name || crd"
                :label="crd.name || crd.metadata?.name || crd" 
                :value="crd.name || crd.metadata?.name || crd"
              />
            </el-select>
          </el-form-item>

          <el-form-item label="实例名称">
            <el-input
              v-model="searchKeyword"
              placeholder="请输入检索名称，回车搜索"
              clearable
              size="small"
              style="width: 160px"
              @keyup.enter="handleSearch"
              @clear="handleSearch"
            />
          </el-form-item>

          <el-form-item label="标签">
            <el-input
              v-model="searchLabels"
              placeholder="例: app=nginx,env=prod 支持多标签"
              clearable
              size="small"
              style="width: 220px"
              @keyup.enter="handleSearch"
              @clear="handleSearch"
            />
          </el-form-item>

          <el-form-item>
            <el-button type="primary" :icon="Search" size="small" @click="handleSearch">
              搜索
            </el-button>
            <el-button :icon="Refresh" size="small" @click="handleRefresh">
              刷新
            </el-button>
          </el-form-item>
        </el-form>
      </div>

      <div class="content-header">
        <span class="resource-count">
          当前选中 CRD: <el-tag type="info">{{ selectedCrdName || '未选择' }}</el-tag> 
          &nbsp;| 共包含实例数: {{ total }}
        </span>
        <el-button type="primary" :icon="Plus" size="small" @click="handleCreateCr" :disabled="!selectedCrdName">
          创建 {{ selectedCrdName ? selectedCrdName.split('.')[0] : '自定义资源' }}
        </el-button>
      </div>

      <el-table
        :data="crList"
        v-loading="loading"
        element-loading-text="加载数据中..."
        class="resource-table"
        empty-text="暂无相关的自定义资源"
        height="calc(100vh - 290px)"
      >
        <el-table-column label="名称" min-width="200">
          <template #default="{ row }">
            <div class="resource-name">
              <el-icon class="resource-icon"><Coin /></el-icon>
              <span>{{ extractName(row) }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="API Version" min-width="150" show-overflow-tooltip align="center">
          <template #default="{ row }">
            <span>{{ row.apiVersion || '-' }}</span>
          </template>
        </el-table-column>

        <el-table-column label="类型" min-width="150" align="center">
          <template #default="{ row }">
            <el-tag size="small" type="success">{{ row.kind || '-' }}</el-tag>
          </template>
        </el-table-column>

        <el-table-column label="自定义标签" min-width="250" align="center">
          <template #default="{ row }">
            <div class="labels-info">
              <template v-if="Object.keys(extractLabels(row)).length > 0">
                <el-tag
                  v-for="(value, key) in Object.entries(extractLabels(row)).slice(0, 2)"
                  :key="key"
                  size="small"
                  type="primary"
                  class="label-tag"
                >
                  {{ value[0] }}={{ value[1] }}
                </el-tag>
                <el-popover
                  v-if="Object.keys(extractLabels(row)).length > 2"
                  placement="top-start"
                  width="300"
                  trigger="hover"
                >
                  <template #reference>
                    <el-tag size="small" type="info" class="label-tag" style="cursor: pointer;">
                      +{{ Object.keys(extractLabels(row)).length - 2 }} ...
                    </el-tag>
                  </template>
                  <div style="display: flex; flex-wrap: wrap; gap: 4px;">
                    <el-tag
                      v-for="(value, key) in extractLabels(row)"
                      :key="key"
                      size="small"
                      type="primary"
                    >
                      {{ key }}={{ value }}
                    </el-tag>
                  </div>
                </el-popover>
              </template>
              <span v-else class="no-labels">无标签</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="创建时间" width="180" align="center">
          <template #default="{ row }">
            <span>{{ extractTime(row) }}</span>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="160" fixed="right" align="center">
          <template #default="{ row }">
            <div class="operation-buttons">
              <el-tooltip content="查看 YAML" placement="top">
                <el-button
                  :icon="View"
                  size="small"
                  type="info"
                  circle
                  @click="handleViewCrYaml(row)"
                />
              </el-tooltip>
              <el-tooltip content="编辑 YAML" placement="top">
                <el-button
                  :icon="Edit"
                  size="small"
                  type="primary"
                  circle
                  @click="handleEditCrYaml(row)"
                />
              </el-tooltip>
              <el-tooltip content="删除" placement="top">
                <el-button
                  :icon="Delete"
                  size="small"
                  type="danger"
                  circle
                  @click="handleDeleteCr(row)"
                />
              </el-tooltip>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <!-- Pagination -->
      <div class="pagination-container" style="margin-top: 15px; display: flex; justify-content: center;">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :background="true"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>

    </el-card>

    <!-- CRD YAML 编辑/查看对话框 -->
    <PodYamlDialog
      :visible="yamlDialogVisible"
      :yaml-content="currentYaml"
      :resource-name="currentCr.name"
      :resource-type="(selectedCrdName || 'CustomResource').split('.')[0]"
      :editable="dialogEditable"
      @update:visible="yamlDialogVisible = $event"
      @close="yamlDialogVisible = false"
      @save="handleSaveYaml"
    />

  </div>
</template>

<style scoped>
.k8s-crd-management {
  padding: 20px;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.config-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title {
  font-size: 20px;
  font-weight: 600;
  color: #2c3e50;
  background: linear-gradient(45deg, #667eea, #764ba2);
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.search-section {
  margin-bottom: 20px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.5);
  border-radius: 12px;
  border: 1px solid rgba(103, 126, 234, 0.1);
}

.search-form {
  margin: 0;
}

.content-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding: 12px 0;
}

.resource-count {
  font-size: 14px;
  color: #606266;
  font-weight: 500;
}

.resource-table {
  border-radius: 8px;
  overflow: hidden;
}

.resource-table :deep(.el-table__header) {
  background: #f8f9fa;
}

.resource-table :deep(.el-table__row:hover) {
  background-color: #f5f7fa;
}

.resource-name {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #2c3e50;
  font-weight: 500;
}

.resource-icon {
  color: #409EFF;
  font-size: 16px;
}

.labels-info {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 4px;
  align-items: center;
}

.label-tag {
  font-size: 11px;
}

.no-labels {
  color: #909399;
  font-size: 12px;
  font-style: italic;
}

.operation-buttons {
  display: flex;
  gap: 8px;
  justify-content: center;
}

/* 按钮及输入框统同样式对齐 k8s-config */
.el-button {
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.el-input :deep(.el-input__wrapper),
.el-select :deep(.el-input__wrapper) {
  background: rgba(255, 255, 255, 0.8);
  border: 1px solid rgba(103, 126, 234, 0.2);
  border-radius: 8px;
  box-shadow: none;
  transition: all 0.3s ease;
}

.el-input :deep(.el-input__wrapper.is-focus),
.el-select :deep(.el-input__wrapper.is-focus) {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(103, 126, 234, 0.2);
  background: rgba(255, 255, 255, 1);
}

.el-tag {
  font-weight: 500;
  border-radius: 8px;
  border: none;
}

</style>