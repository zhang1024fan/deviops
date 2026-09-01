<template>
  <div class="taskansible-management">
    <el-card shadow="hover" class="taskansible-card">
      <template #header>
        <div class="card-header">
          <span class="title">Ansible任务管理</span>
        </div>
      </template>
      
      <!-- 搜索表单 -->
      <div class="search-section">
        <el-form :inline="true" :model="queryParams" class="search-form" ref="queryForm">
      <el-form-item label="搜索模式">
        <el-radio-group v-model="searchMode" size="small" @change="handleSearchModeChange">
          <el-radio-button label="all">全部</el-radio-button>
          <el-radio-button label="name">按名称</el-radio-button>
          <el-radio-button label="type">按类型</el-radio-button>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="任务名称" v-if="searchMode === 'name' || searchMode === 'all'">
        <el-input
          placeholder="请输入任务名称"
          size="small"
          clearable
          v-model="queryParams.name"
          @keyup.enter="searchTasks"
          :disabled="searchMode === 'type'"
        />
      </el-form-item>
      <el-form-item label="任务类型" style="width: 200px;" v-if="searchMode === 'type' || searchMode === 'all'">
        <el-select
          size="small"
          placeholder="请选择任务类型"
          v-model="queryParams.type"
          @change="searchTasks"
          style="width: 100%"
          :disabled="searchMode === 'name'"
        >
          <el-option label="手动任务" :value="1" />
          <el-option label="自动任务" :value="2" />
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" icon="Search" size="small" @click="searchTasks">搜索</el-button>
        <el-button type="warning" icon="Refresh" size="small" @click="resetQuery">重置</el-button>
      </el-form-item>
        </el-form>
      </div>

      <!-- 操作按钮 -->
      <div class="action-section">
        <el-button
          plain
          type="success"
          icon="Plus"
          size="small"
          v-authority="['task:ansible:create']"
          @click="handleCreate"
        >新增Ansible任务</el-button>
      </div>
    
      <!-- 列表区域 -->
      <div class="table-section">
    <el-table
      v-loading="loading"
      :data="tasks"
      border
      stripe
      style="width: 100%"
      :header-cell-style="{ background: '#eef1f6', color: '#606266' }"
    >
      <el-table-column prop="name" label="任务名称" width="200">
        <template #default="{row}">
          <div style="display: flex; align-items: center; gap: 8px;">
            <img 
              src="@/assets/image/ansible.svg" 
              alt="ansible"
              style="width: 20px; height: 20px; object-fit: contain; flex-shrink: 0;"
            />
            <span class="task-name-link" @click="goToHistory(row)">{{ row.name }}</span>
          </div>
        </template>
      </el-table-column>

      <el-table-column prop="type" label="类型" width="130">
        <template #default="{row}">
          <div style="display: flex; align-items: center; gap: 8px;">
            <img 
              :src="getTaskTypeIcon(row.type)" 
              :alt="getTaskTypeName(row.type)"
              style="width: 18px; height: 18px; object-fit: contain; flex-shrink: 0;"
            />
            <span :class="getTypeClass(row.type)">
              {{ getTaskTypeName(row.type) }}
            </span>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="130">
        <template #default="{row}">
          <div style="display: flex; align-items: center; gap: 8px;">
            <img 
              src="@/assets/image/zhuangtai.svg" 
              alt="状态"
              style="width: 18px; height: 18px; object-fit: contain; flex-shrink: 0;"
            />
            <span :class="getStatusClass(row.status)">
              {{ getStatusText(row.status) }}
            </span>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="taskCount" label="任务数量" >
        <template #default="{row}">
          <span>{{ row.taskCount || 0 }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="totalDuration" label="执行耗时" width="150">
        <template #default="{row}">
          <div style="display: flex; align-items: center; gap: 8px;">
            <img 
              src="@/assets/image/定时关闭.svg" 
              alt="时间"
              style="width: 18px; height: 18px; object-fit: contain; flex-shrink: 0;"
            />
            <span>{{ formatDuration(row.totalDuration) }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="description" label="任务描述" >
        <template #default="{row}">
          <span :title="row.description" style="display: inline-block; max-width: 180px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">
            {{ row.description || '暂无描述' }}
          </span>
        </template>
      </el-table-column>
      <el-table-column prop="gitRepo" label="Git仓库"  v-if="queryParams.type === 2">
        <template #default="{row}">
          <span :title="row.gitRepo" style="display: inline-block;  overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">
            {{ row.gitRepo || '-' }}
          </span>
        </template>
      </el-table-column>
      <el-table-column prop="createdAt" label="创建时间"  />
      <el-table-column prop="updatedAt" label="更新时间"  />
      <el-table-column prop="cron_expr" label="Cron表达式" width="150">
        <template #default="{ row }">
           <el-tag v-if="row.cron_expr" size="small" type="info">{{ row.cron_expr }}</el-tag>
           <span v-else style="color: #909399;">-</span>
        </template>
      </el-table-column>
      <el-table-column prop="is_recurring" label="定时开关" width="100" align="center">
        <template #default="{ row }">
          <el-switch
            v-model="row.is_recurring"
            :active-value="1"
            :inactive-value="0"
            inline-prompt
            active-text="开"
            inactive-text="关"
            @change="handleRecurringChange(row)"
          />
        </template>
      </el-table-column>
      <el-table-column label="操作" width="120" fixed="right">
        <template #default="{row}">
          <div class="operation-buttons">
            <el-tooltip content="启动任务" placement="top">
              <el-button
                type="success"
                :icon="VideoPlay"
                size="small"
                v-authority="['task:ansible:start']"
                circle
                @click="() => showStartTaskDialog(row)"
              />
            </el-tooltip>
            <el-tooltip content="编辑" placement="top">
              <el-button
                type="primary"
                icon="Edit"
                size="small"
                circle
                @click="handleEdit(row)"
              />
            </el-tooltip>
            <el-tooltip content="删除" placement="top">
              <el-button
                type="danger"
                :icon="Delete"
                size="small"
                circle
                v-authority="['task:ansible:delete']"
                @click="handleDelete(row.id)"
              />
            </el-tooltip>
          </div>
        </template>
      </el-table-column>
        </el-table>
      </div>

      <!-- 分页区域 -->
      <div class="pagination-section">
    <el-pagination
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      :current-page="queryParams.page"
      :page-sizes="[10, 50, 100]"
      :page-size="queryParams.pageSize"
      layout="total, sizes, prev, pager, next, jumper"
      :total="total"
    />
      </div>
    </el-card>

    <!-- Ansible任务表单对话框 -->
    <el-dialog
      :title="isEdit ? '编辑Ansible任务' : '新增Ansible任务'"
      v-model="formVisible"
      v-authority="['task:ansible:create']"
      width="50%"
      @close="formDialogClosed"
      :modal="false"
      :append-to-body="true"
      class="modern-dialog"
    >
      <el-form
        ref="ansibleFormRef"
        label-width="120px"
        :model="currentTask"
        :rules="taskRules"
      >
        <el-form-item label="任务名称" prop="name">
          <el-input v-model="currentTask.name" placeholder="请输入任务名称" />
        </el-form-item>
        
        <el-form-item label="任务描述">
          <el-input
            v-model="currentTask.description"
            type="textarea"
            :rows="2"
            placeholder="请输入任务描述（可选）"
            maxlength="500"
            show-word-limit
          />
        </el-form-item>
        
        <el-form-item label="任务类型" prop="type">
          <el-radio-group v-model="currentTask.type">
            <el-radio :label="1">手动任务</el-radio>
            <el-radio :label="2">自动任务</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="使用配置中心">
           <el-switch v-model="currentTask.use_config" :active-value="1" :inactive-value="0" active-text="开启" inactive-text="关闭" />
        </el-form-item>

        <!-- 使用配置中心时的选项 -->
        <template v-if="currentTask.use_config === 1">
           <el-form-item label="Inventory配置" >
             <el-select v-model="currentTask.inventory_config_id" placeholder="请选择Inventory配置" clearable filterable style="width: 100%">
                <el-option v-for="item in configOptions.inventory" :key="item.ID" :label="item.Name" :value="item.ID" />
             </el-select>
           </el-form-item>
           <el-form-item label="全局变量配置">
             <el-select v-model="currentTask.global_vars_config_id" placeholder="请选择全局变量配置" clearable filterable style="width: 100%">
                <el-option v-for="item in configOptions.globalVars" :key="item.ID" :label="item.Name" :value="item.ID" />
             </el-select>
           </el-form-item>
           <el-form-item label="额外变量配置">
             <el-select v-model="currentTask.extra_vars_config_id" placeholder="请选择额外变量配置" clearable filterable style="width: 100%">
                <el-option v-for="item in configOptions.extraVars" :key="item.ID" :label="item.Name" :value="item.ID" />
             </el-select>
           </el-form-item>
           <el-form-item label="命令行参数配置">
             <el-select v-model="currentTask.cli_args_config_id" placeholder="请选择命令行参数配置" clearable filterable style="width: 100%">
                <el-option v-for="item in configOptions.cliArgs" :key="item.ID" :label="item.Name" :value="item.ID" />
             </el-select>
           </el-form-item>
        </template>

        <!-- 手动任务配置 -->
        <template v-if="currentTask.type === 1">
          <el-form-item label="主机分组" prop="hostGroups" v-if="currentTask.use_config === 0">
            <div style="margin-bottom: 10px">
              <el-button
                type="primary"
                icon="Plus"
                size="small"
                @click="handleAddHostGroup"
              >添加主机分组</el-button>
            </div>
            <div v-if="Object.keys(currentTask.hostGroups).length > 0" 
                 style="border: 1px solid #ebeef5; border-radius: 4px; padding: 10px; background: #f5f7fa">
              <div style="color: #606266; margin-bottom: 8px">已配置分组:</div>
              <div
                v-for="(hosts, groupName) in currentTask.hostGroups"
                :key="groupName"
                style="display: flex; justify-content: space-between; align-items: center; padding: 6px 10px; background: white; border-radius: 4px; margin-bottom: 5px"
              >
                <span>{{ groupName }} ({{ hosts.length }}台主机)</span>
                <el-button
                  type="danger"
                  icon="Close"
                  size="small"
                  circle
                  plain
                  @click="removeHostGroup(groupName)"
                />
              </div>
            </div>
          </el-form-item>

          <template v-if="currentTask.use_config === 0">
            <el-form-item label="全局变量">
              <el-input
                v-model="currentTask.variables"
                type="textarea"
                :rows="3"
                placeholder='请输入JSON格式的全局变量，如: {"name":"mysql-test","version":"5.7"}'
              />
            </el-form-item>
            <el-form-item label="额外变量">
              <el-input
                v-model="currentTask.extra_vars"
                type="textarea"
                :rows="3"
                placeholder='请输入JSON或YAML格式的额外变量'
              />
            </el-form-item>
            <el-form-item label="命令行参数">
              <el-input
                v-model="currentTask.cli_args"
                placeholder='请输入命令行参数'
              />
            </el-form-item>
          </template>

          <el-form-item label="Playbook文件" prop="playbooks">
            <el-upload
              ref="playbookUpload"
              :auto-upload="false"
              :show-file-list="true"
              :limit="1"
              accept=".yml,.yaml"
              @change="handlePlaybookChange"
            >
              <el-button size="small" type="primary" icon="Upload">选择Playbook</el-button>
              <template #tip>
                <div class="el-upload__tip">只能上传yaml文件，且不超过10MB</div>
              </template>
            </el-upload>
          </el-form-item>

          <el-form-item label="Roles压缩包">
            <el-upload
              ref="rolesUpload"
              :auto-upload="false"
              :show-file-list="true"
              :limit="1"
              accept=".zip"
              @change="handleRolesChange"
            >
              <el-button size="small" type="primary" icon="Upload">选择Roles</el-button>
              <template #tip>
                <div class="el-upload__tip">只能上传zip文件，且不超过50MB</div>
              </template>
            </el-upload>
          </el-form-item>
        </template>

        <!-- 自动任务配置 -->
        <template v-if="currentTask.type === 2">
          <el-form-item label="主机分组" prop="hostGroups" v-if="currentTask.use_config === 0">
            <div style="margin-bottom: 10px">
              <el-button
                type="primary"
                icon="Plus"
                size="small"
                @click="handleAddHostGroup"
              >添加主机分组</el-button>
            </div>
            <div v-if="Object.keys(currentTask.hostGroups).length > 0" 
                 style="border: 1px solid #ebeef5; border-radius: 4px; padding: 10px; background: #f5f7fa">
              <div style="color: #606266; margin-bottom: 8px">已配置分组:</div>
              <div
                v-for="(hosts, groupName) in currentTask.hostGroups"
                :key="groupName"
                style="display: flex; justify-content: space-between; align-items: center; padding: 6px 10px; background: white; border-radius: 4px; margin-bottom: 5px"
              >
                <span>{{ groupName }} ({{ hosts.length }}台主机)</span>
                <el-button
                  type="danger"
                  icon="Close"
                  size="small"
                  circle
                  plain
                  @click="removeHostGroup(groupName)"
                />
              </div>
            </div>
          </el-form-item>

          <el-form-item label="Git仓库地址" prop="gitRepo">
            <el-input v-model="currentTask.gitRepo" placeholder="请输入Git仓库地址" />
            <div style="margin-top: 4px; color: #909399; font-size: 12px;">
              <span>示例: git@gitee.com:zhang_fan1024/ansible-playbook.git</span>
            </div>
          </el-form-item>

          <el-form-item label="Playbook路径">
             <el-select
                v-model="currentTask.playbook_paths"
                multiple
                filterable
                allow-create
                default-first-option
                placeholder="请输入Playbook文件路径并回车"
                style="width: 100%"
              >
             </el-select>
             <div style="font-size: 12px; color: #999; margin-top: 4px">输入如: site.yml, roles/db.yml</div>
          </el-form-item>

          <template v-if="currentTask.use_config === 0">
            <el-form-item label="全局变量">
                <el-input
                v-model="currentTask.variables"
                type="textarea"
                :rows="3"
                placeholder='请输入JSON格式的全局变量，如: {"name":"mysql-test","version":"5.7"}'
                />
            </el-form-item>
            <el-form-item label="额外变量">
                <el-input
                v-model="currentTask.extra_vars"
                type="textarea"
                :rows="3"
                placeholder='请输入JSON或YAML格式的额外变量'
                />
            </el-form-item>
            <el-form-item label="命令行参数">
                <el-input
                v-model="currentTask.cli_args"
                placeholder='请输入命令行参数'
                />
            </el-form-item>
          </template>

        </template>

        <el-form-item label="最大历史记录">
                <el-input
                type="number"
                v-model="currentTask.max_history_keep"
                placeholder='保存的最大历史记录数，默认为3'
                />
            </el-form-item>
        
        <!-- 周期任务配置 -->
        <el-form-item label="周期性任务">
           <el-switch v-model="currentTask.is_recurring" :active-value="1" :inactive-value="0" active-text="开启" inactive-text="关闭" />
        </el-form-item>
        <el-form-item label="Cron表达式" v-if="currentTask.is_recurring === 1" required>
           <el-input v-model="currentTask.cron_expr" placeholder="例如: 0 0 * * * (每天零点执行)" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="formVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 主机分组选择对话框 -->
    <el-dialog
      title="配置主机分组"
      v-model="hostGroupDialogVisible"
      width="50%"
      :modal="false"
      :append-to-body="true"
      class="modern-dialog"
    >
      <el-form :inline="true" style="margin-bottom: 15px">
        <el-form-item label="分组名称">
          <el-input
            v-model="newGroupName"
            placeholder="请输入分组名称（如：web、api、db）"
            style="width: 250px"
          />
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            icon="Plus"
            @click="showHostSelector"
            :disabled="!newGroupName.trim()"
          >选择主机</el-button>
        </el-form-item>
      </el-form>
      
      <!-- 已配置的分组列表 -->
      <div v-if="Object.keys(currentTask.hostGroups).length > 0" 
           style="border: 1px solid #ebeef5; border-radius: 4px; padding: 15px; background: #f5f7fa">
        <div style="color: #606266; margin-bottom: 10px; font-weight: 500">已配置分组:</div>
        <div
          v-for="(hosts, groupName) in currentTask.hostGroups"
          :key="groupName"
          style="margin-bottom: 12px"
        >
          <div style="display: flex; justify-content: space-between; align-items: center; padding: 8px 12px; background: white; border-radius: 4px; border: 1px solid #e4e7ed">
            <div>
              <el-tag type="primary" style="margin-right: 8px">{{ groupName }}</el-tag>
              <span style="color: #666">{{ hosts.length }}台主机</span>
              <el-button 
                type="text" 
                size="small" 
                @click="showGroupHosts(groupName, hosts)"
                style="margin-left: 8px; color: #409EFF"
              >
                查看详情
              </el-button>
            </div>
            <el-button
              type="danger"
              icon="Close"
              size="small"
              circle
              plain
              @click="removeHostGroup(groupName)"
            />
          </div>
          <!-- 主机详情展示 -->
          <div v-if="expandedGroups.includes(groupName)" 
               style="margin-top: 8px; padding: 8px; background: #fafafa; border-radius: 4px; border: 1px solid #eee">
            <div style="color: #666; font-size: 12px; margin-bottom: 5px">主机 ID 列表:</div>
            <el-tag 
              v-for="hostId in hosts" 
              :key="hostId" 
              size="small" 
              type="info"
              style="margin-right: 5px; margin-bottom: 3px"
            >
              ID: {{ hostId }}
            </el-tag>
          </div>
        </div>
      </div>
      
      <div v-else style="text-align: center; padding: 40px; color: #999">
        <el-icon size="48"><Connection /></el-icon>
        <div style="margin-top: 16px">暂无主机分组，请点击上方按钮添加</div>
      </div>
      
      <template #footer>
        <el-button @click="hostGroupDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 主机选择器 -->
    <CreateTaskHost
      v-model="showHostDialog"
      @hosts-selected="handleHostSelected"
    />

    <!-- Ansible任务流程组件 -->
    <AnsibleJobFlow ref="ansibleJobFlowRef" />


  </div>
</template>

<script setup>
import { ref, reactive, onMounted, shallowRef, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Connection, VideoPlay, Delete, Edit } from '@element-plus/icons-vue'
import {
  GetAnsibleTaskList,
  GetAnsibleConfigList,
  CreateAnsibleTask,
  UpdateAnsibleTask,
  GetAnsibleTaskById,
  StartAnsibleTask,
  DeleteAnsibleTask,
  GetAnsibleTasksByName,
  GetAnsibleTasksByType
} from '@/api/task'
import cmdbAPI from '@/api/cmdb'

import AnsibleJobFlow from './Job/AnsibleJobFlow.vue'
import CreateTaskHost from './Job/CreateTaskHost.vue'

const router = useRouter()
const tasks = ref([])
const loading = ref(false)
const formVisible = ref(false)
const searchMode = ref('all') // 搜索模式：all-全部, name-按名称, type-按类型

const ansibleFormRef = shallowRef(null)
const total = ref(0)
const submitting = ref(false)
const hostGroupDialogVisible = ref(false)
const showHostDialog = ref(false)
const newGroupName = ref('')
const isEdit = ref(false)
const editId = ref(null)

// 启动任务相关
const selectedTask = ref(null)

const ansibleJobFlowRef = shallowRef(null)
const expandedGroups = ref([])

const queryParams = reactive({
  page: 1,
  pageSize: 10,
  name: '',
  type: null
})

const currentTask = ref({
  name: '',
  description: '',
  type: 1,
  hostGroups: {},
  variables: '',
  extra_vars: '',
  cli_args: '',
  gitRepo: '',
  playbookFile: null,
  rolesFile: null,
  use_config: 0,
  inventory_config_id: null,
  global_vars_config_id: null,
  extra_vars_config_id: null,
  cli_args_config_id: null,
  max_history_keep: 3,
  is_recurring: 0,
  cron_expr: '',
  playbook_paths: [],
  view_id: null
})

// 配置选项
const configOptions = reactive({
  inventory: [],
  globalVars: [],
  extraVars: [],
  cliArgs: []
})

// 获取配置列表
const fetchConfigs = async () => {
  try {
    const types = [
      { type: 1, key: 'inventory' },
      { type: 2, key: 'globalVars' },
      { type: 3, key: 'extraVars' },
      { type: 4, key: 'cliArgs' }
    ]
    
    for (const item of types) {
      const res = await GetAnsibleConfigList({
        page: 1,
        size: 100,
        type: item.type
      })
      if (res.data && res.data.code === 200) {
        // 根据后端返回的具体结构调整，这里假设 data.data 是列表
        const list = res.data.data.list || []
        // 映射到 configOptions
        switch(item.key) {
           case 'inventory': configOptions.inventory = list; break;
           case 'globalVars': configOptions.globalVars = list; break;
           case 'extraVars': configOptions.extraVars = list; break;
           case 'cliArgs': configOptions.cliArgs = list; break;
        }
      }
    }
  } catch (error) {
    console.error('获取配置列表失败:', error)
  }
}

const taskRules = reactive({
  name: [{ required: true, message: '请输入任务名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择任务类型', trigger: 'change' }],
  hostGroups: [{
    required: false,
    validator: (rule, value, callback) => {
      // 只有在不使用配置中心时，才校验手动的主机分组
      if (currentTask.value.use_config === 0 && Object.keys(currentTask.value.hostGroups).length === 0) {
        callback(new Error('请至少配置一个主机分组'))
      } else {
        callback()
      }
    },
    trigger: 'change'
  }],
  playbooks: [{
    required: false,
    validator: (rule, value, callback) => {
      // 如果是手动任务
      if (currentTask.value.type === 1) {
        // 如果是编辑模式且已有Playbook，则不强制上传
        if (isEdit.value && currentTask.value.existingPlaybook) {
          callback()
          return
        }
      }
      callback()
    },
    trigger: 'change'
  }],
  gitRepo: [{
    required: false,
    validator: (rule, value, callback) => {
      if (currentTask.value.type === 2 && !currentTask.value.gitRepo) {
        callback(new Error('请输入Git仓库地址'))
      } else {
        callback()
      }
    },
    trigger: 'blur'
  }]
})

const fetchTasks = async () => {
  loading.value = true
  try {
    console.log('请求参数:', JSON.stringify(queryParams, null, 2))
    console.log('当前搜索模式:', searchMode.value)
    
    let response;
    
    // 根据搜索模式和查询条件选择不同的API接口
    if (searchMode.value === 'name' && queryParams.name && queryParams.name.trim()) {
      // 名称搜索模式
      console.log('使用名称查询接口')
      response = await GetAnsibleTasksByName({
        name: queryParams.name,
        page: queryParams.page,
        pageSize: queryParams.pageSize
      })
    } else if (searchMode.value === 'type' && queryParams.type !== null && queryParams.type !== undefined) {
      // 类型搜索模式
      console.log('使用类型查询接口')
      response = await GetAnsibleTasksByType({
        type: queryParams.type,
        page: queryParams.page,
        pageSize: queryParams.pageSize
      })
    } else if (searchMode.value === 'all') {
      // 全部搜索模式，根据具体有值的条件选择接口
      if (queryParams.name && queryParams.name.trim()) {
        console.log('全部模式-使用名称查询接口')
        response = await GetAnsibleTasksByName({
          name: queryParams.name,
          page: queryParams.page,
          pageSize: queryParams.pageSize
        })
      } else if (queryParams.type !== null && queryParams.type !== undefined) {
        console.log('全部模式-使用类型查询接口')
        response = await GetAnsibleTasksByType({
          type: queryParams.type,
          page: queryParams.page,
          pageSize: queryParams.pageSize
        })
      } else {
        console.log('全部模式-使用普通列表查询')
        response = await GetAnsibleTaskList({
          page: queryParams.page,
          pageSize: queryParams.pageSize
        })
      }
    } else {
      // 默认情况，使用普通列表查询
      console.log('默认-使用普通列表查询')
      response = await GetAnsibleTaskList({
        page: queryParams.page,
        pageSize: queryParams.pageSize
      })
    }
    
    console.log('后端API响应:', response)
    
    if (response?.data) {
      // 直接使用后端返回的数据结构
      const responseData = response.data
      
      // 兼容不同的数据结构，确保taskList始终是数组
      let taskList = []
      
      if (Array.isArray(responseData.data)) {
        // 标准格式：{code, message, data: {data: [...], total: n}}
        taskList = responseData.data
      } else if (Array.isArray(responseData)) {
        // 简化格式：{code, message, data: [...]}
        taskList = responseData
      } else if (responseData.data && Array.isArray(responseData.data.data)) {
        // 嵌套格式：{code, message, data: {data: {data: [...], total: n}}}
        taskList = responseData.data.data
      } else {
        console.warn('未知的API响应格式:', responseData)
        taskList = []
      }
      
      console.log('解析后的taskList:', taskList)
      console.log('taskList类型:', typeof taskList, '是否为数组:', Array.isArray(taskList))
      
      // 确保taskList是数组后再进行map操作
      if (Array.isArray(taskList)) {
        // 将后端返回的大写字段名映射为小写
        tasks.value = taskList.map(item => ({
          id: item.ID,
          name: item.Name,
          description: item.Description,
          type: item.Type,
          gitRepo: item.GitRepo,
          hostGroups: item.HostGroups,
          allHostIDs: item.AllHostIDs,
          globalVars: item.GlobalVars,
          status: item.status,
          errorMsg: item.ErrorMsg,
          taskCount: item.TaskCount,
          totalDuration: item.TotalDuration,
          is_recurring: item.IsRecurring,
          cron_expr: item.CronExpr,
          createdAt: formatTime(item.CreatedAt),
          updatedAt: formatTime(item.UpdatedAt),
          works: item.Works
        }))
        
        total.value = responseData.total || taskList.length
        console.log('处理后的任务列表:', tasks.value)
      } else {
        console.error('taskList不是数组:', taskList)
        tasks.value = []
        total.value = 0
        ElMessage.warning('数据格式异常，请检查API返回格式')
      }
    } else {
      tasks.value = []
      total.value = 0
    }
  } catch (error) {
    console.error('获取Ansible任务列表失败:', error)
    console.error('错误详情:', {
      message: error.message,
      stack: error.stack,
      response: error.response?.data
    })
    ElMessage.error(`获取任务列表失败: ${error.message || '未知错误'}`)
    tasks.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  isEdit.value = false
  editId.value = null
  currentTask.value = {
    name: '',
    description: '',
    type: 1,
    hostGroups: {},
    variables: '',
    extra_vars: '',
    cli_args: '',
    gitRepo: '',
    playbookFile: null,
    rolesFile: null,
    use_config: 0,
    inventory_config_id: null,
    global_vars_config_id: null,
    extra_vars_config_id: null,
    cli_args_config_id: null,
    is_recurring: 0,
    cron_expr: '',
    playbook_paths: [],
    view_id: null
  }
  expandedGroups.value = []
  formVisible.value = true
}

const handleEdit = async (row) => {
  isEdit.value = true
  editId.value = row.id // row.id mapped in table
  
  try {
     const res = await GetAnsibleTaskById(row.id)

     if (res && res.data && res.data.code === 200) {
        const data = res.data.data.task_info
        console.log('获取的任务详情:', data)
        currentTask.value = {
            name: data.Name,
            description: data.Description,
            type: data.Type,
            hostGroups: typeof data.HostGroups === 'string' ? parseHostGroups(data.HostGroups) : (data.HostGroups || {}),
            variables: data.GlobalVars || '',
            extra_vars: data.ExtraVars || '',
            cli_args: data.CliArgs || '',
            gitRepo: data.GitRepo || '',
            playbookFile: null,
            rolesFile: null,
            use_config: data.UseConfig || 0,
            inventory_config_id: data.InventoryConfigID || null,
            global_vars_config_id: data.GlobalVarsConfigID || null,
            extra_vars_config_id: data.ExtraVarsConfigID || null,
            cli_args_config_id: data.CliArgsConfigID || null,
            max_history_keep: data.MaxHistoryKeep || 3,
            is_recurring: data.IsRecurring || 0,
            cron_expr: data.CronExpr || '',
            // Handle playbook_paths parsing if it's a string
            playbook_paths: data.PlaybookPaths ? (typeof data.PlaybookPaths === 'string' ? JSON.parse(data.PlaybookPaths) : data.PlaybookPaths) : [],
            view_id: data.ViewId || null
        }
        formVisible.value = true
     }
  } catch (e) {
      console.error(e)
      ElMessage.error('获取任务详情失败')
  }
}

const handleSubmit = async () => {
  try {
    await ansibleFormRef.value?.validate()
    submitting.value = true

    const formData = new FormData()
    // Update ID if editing
    if (isEdit.value && editId.value) {
        formData.append('id', editId.value)
    }

    formData.append('name', currentTask.value.name)
    if (currentTask.value.description) {
      formData.append('description', currentTask.value.description)
    }
    formData.append('type', currentTask.value.type.toString())
    
    // Config Switch
    formData.append('use_config', currentTask.value.use_config ? '1' : '0')
    
    if (currentTask.value.use_config === 1) {
        // 使用配置中心
        if (currentTask.value.inventory_config_id) formData.append('inventory_config_id', currentTask.value.inventory_config_id)
        if (currentTask.value.global_vars_config_id) formData.append('global_vars_config_id', currentTask.value.global_vars_config_id)
        if (currentTask.value.extra_vars_config_id) formData.append('extra_vars_config_id', currentTask.value.extra_vars_config_id)
        if (currentTask.value.cli_args_config_id) formData.append('cli_args_config_id', currentTask.value.cli_args_config_id)
    } else {
        // 不使用配置中心，使用手动输入
        formData.append('hostGroups', JSON.stringify(currentTask.value.hostGroups))
        if (currentTask.value.variables) formData.append('variables', currentTask.value.variables)
        if (currentTask.value.extra_vars) formData.append('extra_vars', currentTask.value.extra_vars)
        if (currentTask.value.cli_args) formData.append('cli_args', currentTask.value.cli_args)
    }

    // Cron / Recurring
    formData.append('is_recurring', currentTask.value.is_recurring ? '1' : '0')
    if (currentTask.value.is_recurring && currentTask.value.cron_expr) {
        formData.append('cron_expr', currentTask.value.cron_expr)
    }

    if (currentTask.value.type === 1) {
      // 手动任务
      if (currentTask.value.playbookFile) {
        formData.append('playbooks', currentTask.value.playbookFile)
      }
      if (currentTask.value.rolesFile) {
        formData.append('roles', currentTask.value.rolesFile)
      }
    } else {
      // 自动任务
      formData.append('gitRepo', currentTask.value.gitRepo)
      if (currentTask.value.playbook_paths && currentTask.value.playbook_paths.length > 0) {
          formData.append('playbook_paths', JSON.stringify(currentTask.value.playbook_paths))
      }
    }
    
    if (currentTask.value.view_id) {
        formData.append('view_id', currentTask.value.view_id)
    }

    let response
    if (isEdit.value) {
      if (!editId.value) {
        ElMessage.warning('编辑模式下缺少任务ID')
        return
      }

      // Variables processing
      let variablesData = {}
      if (typeof currentTask.value.variables === 'string' && currentTask.value.variables.trim()) {
        try {
           variablesData = JSON.parse(currentTask.value.variables)
        } catch (e) {
           ElMessage.error('全局变量格式错误，请输入有效的JSON字符串')
           submitting.value = false
           return
        }
      } else {
         variablesData = currentTask.value.variables || {}
      }

      // HostGroups processing
      let hostGroupsData = {}
      try {
         if (typeof currentTask.value.hostGroups === 'string') {
            hostGroupsData = currentTask.value.hostGroups ? JSON.parse(currentTask.value.hostGroups) : {}
         } else {
            hostGroupsData = currentTask.value.hostGroups || {}
         }
      } catch (e) {
         hostGroupsData = {}
      }

      // PlaybookPaths processing
      let playbookPathsData = []
      try {
         if (typeof currentTask.value.playbook_paths === 'string') {
            playbookPathsData = currentTask.value.playbook_paths ? JSON.parse(currentTask.value.playbook_paths) : []
         } else if (Array.isArray(currentTask.value.playbook_paths)) {
            playbookPathsData = currentTask.value.playbook_paths
         }
      } catch (e) {
         playbookPathsData = []
      }

      const updateData = {
        id: editId.value,
        name: currentTask.value.name,
        description: currentTask.value.description || '',
        type: currentTask.value.type,
        useConfig: currentTask.value.use_config,
        isRecurring: currentTask.value.is_recurring,
        maxHistoryKeep: parseInt(currentTask.value.max_history_keep) || 3,
        cronExpr: currentTask.value.cron_expr || '',
        viewId: currentTask.value.view_id || 0,
        // Config Center fields
        inventoryConfigId: currentTask.value.inventory_config_id || 0,
        globalVarsConfigId: currentTask.value.global_vars_config_id || 0,
        extraVarsConfigId: currentTask.value.extra_vars_config_id || 0,
        cliArgsConfigId: currentTask.value.cli_args_config_id || 0,
        // Manual fields
        variables: variablesData,
        extraVars: currentTask.value.extra_vars || '',
        cliArgs: currentTask.value.cli_args || '',
        hostGroups: hostGroupsData,
        // Git/Auto fields (PlaybookPaths is array)
        gitRepo: currentTask.value.gitRepo || '',
        playbookPaths: playbookPathsData
      }

      response = await UpdateAnsibleTask(updateData)
    } else {
        response = await CreateAnsibleTask(formData)
    }
    
    if (response?.data?.code === 200) {
      ElMessage.success(isEdit.value ? 'Ansible任务更新成功' : 'Ansible任务创建成功')
      formVisible.value = false
      fetchTasks()
    } else {
      throw new Error(response?.data?.message || (isEdit.value ? '更新任务失败' : '创建任务失败'))
    }
  } catch (error) {
    console.error('创建 Ansible 任务失败:', error)
    ElMessage.error(`创建任务失败: ${error.message || '未知错误'}`)
  } finally {
    submitting.value = false
  }
}

const handleDelete = async (id) => {
  try {
    await ElMessageBox.confirm(
      '确定要删除该Ansible任务吗？删除后将同时删除相关的子任务和执行记录。',
      '删除确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning',
        dangerouslyUseHTMLString: false
      }
    )
    
    console.log('开始删除任务，ID:', id)
    const response = await DeleteAnsibleTask(id)
    
    if (response?.data?.code === 200 || response?.status === 200) {
      ElMessage.success('任务删除成功')
      fetchTasks()
    } else {
      throw new Error(response?.data?.message || '删除任务失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除任务失败:', error)
      ElMessage.error(`删除失败: ${error.response?.data?.message || error.message || '未知错误'}`)
    }
  }
}





const handleAddHostGroup = () => {
  newGroupName.value = ''
  hostGroupDialogVisible.value = true
}

const showHostSelector = () => {
  if (!newGroupName.value.trim()) {
    ElMessage.warning('请先输入分组名称')
    return
  }
  showHostDialog.value = true
}

const handleHostSelected = (hosts) => {
  if (!newGroupName.value.trim()) {
    ElMessage.warning('请先输入分组名称')
    return
  }
  
  if (hosts.length === 0) {
    ElMessage.warning('请选择至少一台主机')
    return
  }

  // 检查分组名是否已存在
  if (currentTask.value.hostGroups[newGroupName.value]) {
    ElMessage.warning(`分组"${newGroupName.value}"已存在，请更换分组名称`)
    return
  }

  // 将主机数据转换为主机 ID 数组
  const hostIds = hosts.map(host => host.id)
  
  // 按照要求的格式存储：{"web":[444,445],"api":[444]}
  currentTask.value.hostGroups[newGroupName.value] = hostIds
  
  console.log('已添加主机分组:', {
    groupName: newGroupName.value,
    hostIds: hostIds,
    allGroups: currentTask.value.hostGroups
  })
  
  ElMessage.success(`已添加分组"${newGroupName.value}"，包含${hostIds.length}台主机`)
  
  // 重置分组名称以便添加下一个分组
  newGroupName.value = ''
  showHostDialog.value = false
}

const removeHostGroup = (groupName) => {
  delete currentTask.value.hostGroups[groupName]
  // 同时从展开列表中移除
  const index = expandedGroups.value.indexOf(groupName)
  if (index > -1) {
    expandedGroups.value.splice(index, 1)
  }
}

// 处理定时任务开关变更
const handleRecurringChange = async (row) => {
  try {
    const statusText = row.is_recurring === 1 ? '开启' : '关闭'
    console.log(`正在${statusText}任务 [${row.name}] 的定时功能`)
    
    // 这里放置更新API调用
    await UpdateAnsibleTask({id: row.id, isRecurring: row.is_recurring})
    
    ElMessage.success(`已${statusText}任务 "${row.name}" 的定时调度`)
  } catch (error) {
    // 失败时回滚状态
    row.is_recurring = row.is_recurring === 1 ? 0 : 1
    console.error('更新定时状态失败:', error)
    ElMessage.error('更新状态失败')
  }
}

const showGroupHosts = (groupName, hosts) => {
  const index = expandedGroups.value.indexOf(groupName)
  if (index > -1) {
    // 已展开，收起
    expandedGroups.value.splice(index, 1)
  } else {
    // 未展开，展开
    expandedGroups.value.push(groupName)
  }
}

const handlePlaybookChange = (file) => {
  currentTask.value.playbookFile = file.raw
}

const handleRolesChange = (file) => {
  currentTask.value.rolesFile = file.raw
}

const searchTasks = () => {
  queryParams.page = 1
  fetchTasks()
}

const resetQuery = () => {
  searchMode.value = 'all'
  queryParams.name = ''
  queryParams.type = null
  queryParams.page = 1
  fetchTasks()
}

// 处理搜索模式切换
const handleSearchModeChange = (mode) => {
  console.log('搜索模式切换为:', mode)
  
  // 根据搜索模式清空不相关的搜索条件
  if (mode === 'name') {
    queryParams.type = null
  } else if (mode === 'type') {
    queryParams.name = ''
  }
  
  // 重置到第一页并重新查询
  queryParams.page = 1
  fetchTasks()
}

const handleSizeChange = (val) => {
  queryParams.pageSize = val
  queryParams.page = 1
  fetchTasks()
}

const handleCurrentChange = (val) => {
  queryParams.page = val
  fetchTasks()
}

const formDialogClosed = () => {
  ansibleFormRef.value?.resetFields()
}

const getTypeClass = (type) => {
  return type === 1 ? 'type-manual' : 'type-auto'
}

const getStatusClass = (status) => {
  switch(status) {
    case 1: return 'status-waiting'
    case 2: return 'status-running'
    case 3: return 'status-success'
    case 4: return 'status-error'
    default: return ''
  }
}

const getStatusText = (status) => {
  switch(status) {
    case 1: return '等待中'
    case 2: return '运行中'
    case 3: return '成功'
    case 4: return '失败'
    default: return '未知'
  }
}

const getTaskTypeIcon = (type) => {
  switch(type) {
    case 1: return require('@/assets/image/shoudong.svg')        // 手动任务
    case 2: return require('@/assets/image/zidong.svg')    // 自动任务
    case 3: return require('@/assets/image/k8s.svg')
    default: return require('@/assets/image/ren.svg')
  }
}

const getTaskTypeName = (type) => {
  switch(type) {
    case 1: return '手动任务'
    case 2: return '自动任务'
    case 3: return 'K8s任务'
    default: return '手动任务'
  }
}

const parseHostGroups = (hostGroupsStr) => {
  try {
    return typeof hostGroupsStr === 'string' ? JSON.parse(hostGroupsStr) : hostGroupsStr
  } catch {
    return {}
  }
}

const formatDuration = (seconds) => {
  if (!seconds || seconds === 0) return '0秒'
  
  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  const remainingSeconds = seconds % 60
  
  const parts = []
  if (hours > 0) parts.push(`${hours}小时`)
  if (minutes > 0) parts.push(`${minutes}分钟`)
  if (remainingSeconds > 0 || parts.length === 0) parts.push(`${remainingSeconds}秒`)
  
  return parts.join('')
}

const formatTime = (timeStr) => {
  if (!timeStr) return ''
  const date = new Date(timeStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false
  }).replace(/\//g, '-')
}

// 直接启动任务，不显示确认对话框
const showStartTaskDialog = async (task) => {
  console.log('直接启动任务:', task)
  selectedTask.value = task
  
  try {
    // 直接调用任务流程组件显示流程
    const taskId = await ansibleJobFlowRef.value?.showFlow?.(task.id)
    console.log('任务流程已显示，任务ID:', taskId)
  } catch (error) {
    console.error('显示任务流程失败:', error)
    ElMessage.error('显示任务流程失败: ' + (error.message || '未知错误'))
  }
}

const goToHistory = (row) => {
  router.push({
    name: 'AnsibleTaskHistory',
    query: {
      id: row.id,
      name: row.name
    }
  })
}


onMounted(() => {
  fetchTasks()
  fetchConfigs()
})
</script>

<style scoped>
.taskansible-management {
  padding: 20px;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.taskansible-card {
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
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.search-section {
  margin-bottom: 16px;
  padding: 16px;
  background: rgba(103, 126, 234, 0.05);
  border-radius: 12px;
  border: 1px solid rgba(103, 126, 234, 0.1);
}

.search-form .el-form-item {
  margin-bottom: 0;
  margin-right: 16px;
}

.search-form .el-form-item__label {
  color: #606266;
  font-weight: 500;
}

.action-section {
  margin-bottom: 16px;
}

.table-section {
  margin-top: 0;
}

.pagination-section {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

:deep(.el-table) {
  background: rgba(255, 255, 255, 0.9);
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
}

:deep(.el-table__header) {
  background: rgba(102, 126, 234, 0.1) !important;
}

:deep(.el-table__header th) {
  background: rgba(102, 126, 234, 0.1) !important;
  color: #2c3e50 !important;
  font-weight: 600;
  border: none;
}

:deep(.el-table__body tr:hover > td) {
  background-color: rgba(102, 126, 234, 0.1) !important;
}

:deep(.el-table td) {
  border: none;
}

:deep(.el-table::before) {
  display: none;
}

:deep(.el-table--border::after) {
  display: none;
}

.operation-buttons {
  display: flex;
  gap: 8px;
  justify-content: center;
}

.operation-buttons .el-button {
  transition: all 0.3s ease;
}

.operation-buttons .el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}

:deep(.el-button--primary) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  border: none;
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
}

:deep(.el-button--primary:hover) {
  background: linear-gradient(135deg, #5a6fd8, #6a4190);
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.4);
}

:deep(.el-button--primary:active) {
  transform: translateY(0);
}

:deep(.el-button--warning) {
  background: linear-gradient(135deg, #f39c12, #e67e22);
  border: none;
  border-radius: 8px;
  font-weight: 500;
  color: white;
  transition: all 0.3s ease;
}

:deep(.el-button--warning:hover) {
  background: linear-gradient(135deg, #e67e22, #d35400);
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(243, 156, 18, 0.4);
}

:deep(.el-button--warning:active) {
  transform: translateY(0);
}

:deep(.el-button--success) {
  background: linear-gradient(135deg, #27ae60, #2ecc71);
  border: none;
  border-radius: 8px;
  font-weight: 500;
  color: white;
  transition: all 0.3s ease;
}

:deep(.el-button--success:hover) {
  background: linear-gradient(135deg, #2ecc71, #16a085);
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(46, 204, 113, 0.4);
}

:deep(.el-button--success:active) {
  transform: translateY(0);
}

:deep(.el-form-item__label) {
  color: #2c3e50;
  font-weight: 500;
}

:deep(.el-input__wrapper) {
  background: rgba(255, 255, 255, 0.8);
  border: 1px solid rgba(102, 126, 234, 0.3);
  border-radius: 8px;
  transition: all 0.3s ease;
}

:deep(.el-input__wrapper:hover) {
  border-color: rgba(102, 126, 234, 0.5);
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.2);
}

:deep(.el-input__wrapper.is-focus) {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2);
}

:deep(.el-select .el-input.is-focus .el-input__wrapper) {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2);
}

:deep(.el-pagination) {
  --el-pagination-bg-color: transparent;
}

:deep(.el-pager li) {
  background: rgba(255, 255, 255, 0.8);
  border-radius: 6px;
  margin: 0 2px;
  border: 1px solid rgba(102, 126, 234, 0.2);
}

:deep(.el-pager li:hover) {
  background: rgba(102, 126, 234, 0.1);
}

:deep(.el-pager li.is-active) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
}

/* 对话框样式 */
:deep(.modern-dialog .el-dialog) {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 0 8px 32px rgba(31, 38, 135, 0.37);
}

:deep(.modern-dialog .el-dialog__header) {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1), rgba(118, 75, 162, 0.1));
  border-radius: 16px 16px 0 0;
  padding: 20px 24px 16px;
  border-bottom: 1px solid rgba(102, 126, 234, 0.2);
}

:deep(.modern-dialog .el-dialog__title) {
  color: #2c3e50;
  font-weight: 600;
  font-size: 18px;
}

:deep(.modern-dialog .el-dialog__body) {
  padding: 24px;
}

:deep(.modern-dialog .el-dialog__footer) {
  padding: 16px 24px 24px;
  background: rgba(248, 249, 250, 0.8);
  border-radius: 0 0 16px 16px;
}

.mb8 {
  margin-bottom: 8px;
}

.type-manual {
  color: #409EFF;
  font-weight: 500;
}

.type-auto {
  color: #67C23A;
  font-weight: 500;
}

.status-waiting {
  color: #909399;
  font-weight: 500;
}

.status-running {
  color: #E6A23C;
  font-weight: 500;
}

.status-success {
  color: #67C23A;
  font-weight: 500;
}

.status-error {
  color: #F56C6C;
  font-weight: 500;
}

.task-name-link {
  color: #667eea;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.3s ease;
}

.task-name-link:hover {
  color: #764ba2;
  text-decoration: underline;
}
</style>