<template>
  <el-container class="home-container">
    <el-aside :width="isCollapse ? '64px' : '200px'">
      <div class="logo">
        <img src="../assets/image/DevOps平台.svg" class="siderbar-logo">
        <h2 v-show="!isCollapse">devops系统</h2>
      </div>
      <el-menu background-color="transparent" text-color="rgba(255,255,255,0.9)" active-text-color="#ffffff" router :default-active="$route.path"
               :collapse="isCollapse" :collapse-transition="false" class="modern-menu">
        <!--无子集菜单-->
        <el-menu-item :index="'/' + item.url" v-for="item in noChildren" :key="item.menuName" @click="saveNavState('/' + item.url)">
          <el-icon><component :is="item.icon" /></el-icon>
          <template v-slot:title>
            <span>{{ item.menuName }}</span>
          </template>
        </el-menu-item>
        <!--有子集菜单-->
        <el-sub-menu :index="item.id + ''" v-for="item in hasChildren" :key="item.id">
          <template #title>
            <el-icon><component :is="item.icon" /></el-icon>
            <span>{{ item.menuName }}</span>
          </template>
          <el-menu-item :index="'/' + subItem.url" v-for="subItem in item.menuSvoList" :key="subItem.id"
                        @click="saveNavState('/' + subItem.url)">
            <el-icon><component :is="subItem.icon" /></el-icon>
            <template #title>
              <span>{{ subItem.menuName }}</span>
            </template>
          </el-menu-item>
        </el-sub-menu>

      </el-menu>
    </el-aside>

    <!-- 主体内容 -->
    <el-container>
      <el-header height="50px">
        <!-- 顶部导航栏,折叠图标 -->
        <div class="fold-btn">
          <el-button type="text" @click="toggleCollapse" class="collapse-btn">
            <el-icon size="24"><component :is="collapseBtnClass" /></el-icon>
          </el-button>
          <HeadImage />
        </div>
        <div class="bread-btn">
          <!-- 面包屑 -->
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/dashboard' }">仪表盘</el-breadcrumb-item>

            <!-- 显示二级标题 -->
            <el-breadcrumb-item v-if="$route.meta && $route.meta.sTitle">
              {{ $route.meta.sTitle }}
            </el-breadcrumb-item>

            <!-- 显示三级标题 -->
            <el-breadcrumb-item v-if="$route.meta && $route.meta.tTitle">
              {{ $route.meta.tTitle }}
            </el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <HeadImage />
      </el-header>
      <Tags />
      <el-main><router-view /></el-main>
    </el-container>
  </el-container>
</template>

<script>


import storage from "@/utils/storage";
import HeadImage   from "@/components/HeadImage.vue";
import Tags from "@/components/Tags.vue";

export default {
  // eslint-disable-next-line vue/multi-word-component-names
  name: "Home",
  components: { HeadImage, Tags },
  data() {
    return {
      leftMenuList: null, // 初始化为null，在mounted中设置
      activePath: '',
      collapseBtnClass: "Fold",
      isCollapse: false,
    }
  },
  computed: {
    // 无子集
    noChildren() {
      return (this.leftMenuList || []).filter(item => !item.menuSvoList) // 过滤无子集
    },
    // 有子集
    hasChildren() {
      return (this.leftMenuList || []).filter(item => item.menuSvoList)
    }
  },
  methods: {
    // 初始化菜单数据
    initMenuData() {
      try {
        const menuData = storage.getItem("leftMenuList");
        console.log('初始化菜单数据:', menuData);

        // 确保数据是数组格式
        if (Array.isArray(menuData)) {
          this.leftMenuList = menuData;

          // 容器管理
          const k8sMenu = this.leftMenuList.find(item => item.menuName === '容器管理')
          if (k8sMenu && k8sMenu.menuSvoList) {
            const crdExists = k8sMenu.menuSvoList.some(sub => sub.url === 'k8s/crd')
            if (!crdExists) {
              k8sMenu.menuSvoList.push({
                id: 99999,
                menuName: 'CRD管理',
                url: 'k8s/crd',
                icon: 'Setting'
              })
            }
          }
          
          // 手动添加配置管理菜单到任务中心
          const taskMenu = this.leftMenuList.find(item => item.menuName === '任务中心')
          if (taskMenu && taskMenu.menuSvoList) {
            const configExists = taskMenu.menuSvoList.some(sub => sub.url === 'task/config')
            if (!configExists) {
              taskMenu.menuSvoList.push({
                id: 99999,
                menuName: '配置管理',
                url: 'task/config',
                icon: 'Setting'
              })
            }
          }
        } else if (menuData) {
          // 如果数据存在但不是数组，尝试解析
          console.warn('菜单数据格式异常，尝试修复:', menuData);
          this.leftMenuList = [];
        } else {
          // 如果没有数据，设为空数组
          console.warn('未找到菜单数据，使用空数组');
          this.leftMenuList = [];
        }

        // 强制触发视图更新
        this.$forceUpdate();
      } catch (error) {
        console.error('初始化菜单数据失败:', error);
        this.leftMenuList = [];
      }
    },
    // 点击实现跳转
    saveNavState(activePath) {
      storage.setItem('activePath', activePath)
      this.activePath = activePath
    },
    // 张开和折叠
    toggleCollapse() {
      this.isCollapse = !this.isCollapse
      if (this.isCollapse) {
        this.collapseBtnClass = 'Fold'  // 折叠
      } else {
        this.collapseBtnClass = 'Expand'  // 展开
      }
    },
    // 移除菜单项的focus效果
    removeFocusOutline() {
      this.$nextTick(() => {
        const menuItems = document.querySelectorAll('.el-menu-item, .el-sub-menu__title')
        menuItems.forEach(item => {
          item.style.outline = 'none'
          item.style.border = 'none'
          item.addEventListener('focus', (e) => {
            e.target.style.outline = 'none'
            e.target.style.border = 'none'
            e.target.blur()
          })
          item.addEventListener('click', (e) => {
            e.target.style.outline = 'none'
            e.target.style.border = 'none'
            setTimeout(() => {
              e.target.blur()
            }, 0)
          })
        })
      })
    }
  },
  mounted() {
    // 确保在mounted阶段重新获取菜单数据，解决浏览器兼容性问题
    this.initMenuData();
    this.removeFocusOutline();
  }
}
</script>

<style lang="less" scoped>
.home-container {
  height: 100%;

  .el-aside {
    background: linear-gradient(135deg, #1e2a5a 0%, #1a1d4d 100%);
    backdrop-filter: blur(10px);
    border-right: 1px solid rgba(255, 255, 255, 0.1);
    box-shadow: 2px 0 10px rgba(0, 0, 0, 0.15);

    .logo {
      margin-top: 8px;
      display: flex;
      align-items: center;
      font-size: 14px;
      height: 50px;
      color: rgba(255, 255, 255, 0.95);
      font-weight: 500;
      padding: 8px 12px;
      white-space: nowrap;
      overflow: hidden;

      .siderbar-logo {
        width: 40px;
        height: 32px;
        margin-right: 10px;
        flex-shrink: 0;
        border-radius: 4px;
      }

      h2 {
        margin: 0;
        font-weight: 700;
        font-style: italic;
        font-family: 'Arial', 'Helvetica', sans-serif;
        letter-spacing: 1.2px;
        font-size: 20px;
        background: linear-gradient(135deg, #ffffff 0%, rgba(255, 255, 255, 0.8) 100%);
        background-clip: text;
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        flex-shrink: 0;
        min-width: 0;
        transform: skewX(-8deg);
        text-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
      }
    }

    .modern-menu {
      border-right: none;
      background: transparent !important;
      
      // 全局移除Element Plus菜单的focus效果
      * {
        outline: none !important;
        
        &:focus {
          outline: none !important;
          box-shadow: none !important;
        }
        
        &:focus-visible {
          outline: none !important;
          box-shadow: none !important;
        }
      }
    }
  }
  .modern-menu {
    // 通用菜单项样式
    .el-menu-item {
      transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
      border-radius: 8px;
      margin: 4px 8px;
      color: rgba(255, 255, 255, 0.9) !important;
      outline: none !important;
      border: none !important;
      
      &:hover {
        background: rgba(255, 255, 255, 0.15) !important;
        transform: translateX(4px);
        color: #ffffff !important;
        box-shadow: 0 2px 10px rgba(255, 255, 255, 0.1);
      }
      
      &:focus {
        outline: none !important;
        border: none !important;
        box-shadow: none !important;
      }
    }
    
    // 一级菜单：没有子菜单的项激活样式
    > .el-menu-item.is-active {
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%) !important;
      color: #ffffff !important;
      border-radius: 10px;
      margin: 4px 8px;
      width: calc(100% - 16px);
      backdrop-filter: blur(10px);
      box-shadow: 0 4px 20px rgba(102, 126, 234, 0.4);
      border: 1px solid rgba(255, 255, 255, 0.2);
      
      &:hover {
        background: linear-gradient(135deg, #7c8ef8 0%, #8a5cb8 100%) !important;
        transform: translateX(2px) translateY(-1px);
        box-shadow: 0 6px 25px rgba(102, 126, 234, 0.5);
      }
    }

    // 子菜单样式
    .el-sub-menu {
      .el-sub-menu__title {
        transition: all 0.3s ease;
        border-radius: 8px;
        margin: 4px 8px;
        color: rgba(255, 255, 255, 0.9) !important;
        outline: none !important;
        border: none !important;
        
        &:hover {
          background: rgba(255, 255, 255, 0.15) !important;
          transform: translateX(4px);
          color: #ffffff !important;
          box-shadow: 0 2px 10px rgba(255, 255, 255, 0.1);
        }
        
        &:focus {
          outline: none !important;
          border: none !important;
          box-shadow: none !important;
        }
      }
      
      // 二级菜单项激活样式
      .el-menu-item {
        outline: none !important;
        border: none !important;
        
        &:focus {
          outline: none !important;
          border: none !important;
          box-shadow: none !important;
        }
        
        &.is-active {
          background: linear-gradient(135deg, #5a67d8 0%, #6b46c1 100%) !important;
          color: #ffffff !important;
          border-radius: 8px;
          margin: 2px 12px 2px 20px;
          width: calc(100% - 32px);
          backdrop-filter: blur(8px);
          box-shadow: 0 3px 15px rgba(102, 126, 234, 0.3);
          border: 1px solid rgba(255, 255, 255, 0.15);
          outline: none !important;
          
          &:hover {
            background: linear-gradient(135deg, #6b73e0 0%, #7c3aed 100%) !important;
            transform: translateX(2px);
            box-shadow: 0 4px 20px rgba(102, 126, 234, 0.4);
          }
          
          &:focus {
            outline: none !important;
            border: 1px solid rgba(255, 255, 255, 0.15) !important;
            box-shadow: 0 3px 15px rgba(102, 126, 234, 0.3) !important;
          }
        }
      }
    }

    // 有子菜单的一级菜单（不应用激活样式）
    > .el-sub-menu.is-active,
    > .el-sub-menu.is-opened {
      background-color: transparent !important;
      
      .el-sub-menu__title {
        background: rgba(255, 255, 255, 0.05) !important;
        border-radius: 8px;
        
        &:hover {
          background: rgba(255, 255, 255, 0.1) !important;
        }
      }
    }
  }




  .el-header {
    background: linear-gradient(135deg, rgba(255, 255, 255, 0.95) 0%, rgba(249, 250, 252, 0.95) 100%);
    backdrop-filter: blur(10px);
    border-bottom: 1px solid rgba(102, 126, 234, 0.1);
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
    align-items: center;
    justify-content: space-between;
    display: flex;

    .fold-btn {
      padding-top: 2px;
      font-size: 23px;
      cursor: pointer;
      
      .collapse-btn {
        padding: 8px;
        border-radius: 8px;
        transition: all 0.3s ease;
        color: #667eea;
        
        &:hover {
          background: rgba(102, 126, 234, 0.1);
          transform: translateY(-1px);
          color: #5a67d8;
        }
      }
    }

    .bread-btn {
      padding-top: 2px;
      position: fixed;
      margin-left: 40px;
      
      .el-breadcrumb {
        .el-breadcrumb__item {
          .el-breadcrumb__inner {
            color: #667eea;
            font-weight: 500;
            transition: all 0.3s ease;
            
            &:hover {
              color: #5a67d8;
            }
          }
          
          &:last-child .el-breadcrumb__inner {
            color: #4a5568;
          }
        }
      }
    }
  }

  .el-main {
    background: linear-gradient(135deg, #f7fafc 0%, #edf2f7 100%);
  }
}

// 全局覆盖Element Plus菜单的focus样式
:deep(.el-menu) {
  * {
    outline: none !important;
    border: none !important;
    box-shadow: none !important;
  }
  
  .el-menu-item,
  .el-sub-menu__title {
    outline: none !important;
    border: none !important;
    
    &:focus,
    &:active,
    &:focus-within,
    &:focus-visible {
      outline: none !important;
      border: none !important;
      background-color: transparent !important;
      box-shadow: none !important;
    }
    
    &::before,
    &::after {
      display: none !important;
    }
  }
  
  // 特别处理菜单项的所有状态
  .el-menu-item {
    &,
    &:hover,
    &:focus,
    &:active,
    &.is-active {
      outline: none !important;
      border: none !important;
      
      &::before,
      &::after {
        display: none !important;
      }
    }
  }
  
  .el-sub-menu__title {
    &,
    &:hover,
    &:focus,
    &:active {
      outline: none !important;
      border: none !important;
      
      &::before,
      &::after {
        display: none !important;
      }
    }
  }
}

// 终极解决方案 - 彻底移除白色边框闪烁
.el-aside {
  // 全局移除所有outline和border
  * {
    outline: 0 !important;
    outline: none !important;
    outline-width: 0 !important;
    outline-style: none !important;
    outline-color: transparent !important;
    border: 0 !important;
    border: none !important;
    border-width: 0 !important;
    border-style: none !important;
    border-color: transparent !important;
    box-shadow: none !important;
    
    &:focus,
    &:active,
    &:hover,
    &:focus-visible,
    &:focus-within,
    &:target {
      outline: 0 !important;
      outline: none !important;
      outline-width: 0 !important;
      outline-style: none !important;
      outline-color: transparent !important;
      border: 0 !important;
      border: none !important;
      border-width: 0 !important;
      border-style: none !important;
      border-color: transparent !important;
      box-shadow: none !important;
    }
  }
}

// 针对Element Plus菜单的特殊处理
:deep(.el-menu) {
  outline: none !important;
  border: none !important;

  * {
    outline: none !important;
    border: none !important;
    outline-width: 0 !important;
    border-width: 0 !important;
  }

  .el-menu-item,
  .el-sub-menu,
  .el-sub-menu__title {
    outline: none !important;
    border: none !important;
    outline-width: 0 !important;
    border-width: 0 !important;

    &,
    &:before,
    &:after,
    &:focus,
    &:active,
    &:hover,
    &:focus-visible,
    &:focus-within,
    &.is-active,
    &.is-opened {
      outline: none !important;
      border: none !important;
      outline-width: 0 !important;
      border-width: 0 !important;
      outline-style: none !important;
      border-style: none !important;
      outline-color: transparent !important;
      border-color: transparent !important;
      box-shadow: none !important;
    }
  }
}

// 折叠菜单弹出层样式修复 - 全局样式
</style>

<style lang="less">
// 只修复左侧菜单栏的折叠弹出层，不影响其他tooltip等弹出层
.el-menu--popup-bottom-start,
.el-menu--popup {
  background: linear-gradient(135deg, #1e2a5a 0%, #1a1d4d 100%) !important;
  border: 1px solid rgba(255, 255, 255, 0.1) !important;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3) !important;
  backdrop-filter: blur(10px) !important;
  border-radius: 8px !important;

  .el-menu-item {
    color: rgba(255, 255, 255, 0.9) !important;
    background: transparent !important;
    transition: all 0.3s ease !important;
    margin: 2px 8px !important;
    border-radius: 6px !important;
    padding: 8px 12px !important;

    &:hover {
      background: rgba(255, 255, 255, 0.15) !important;
      color: #ffffff !important;
      transform: translateX(2px) !important;
    }

    &.is-active {
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%) !important;
      color: #ffffff !important;
      box-shadow: 0 2px 10px rgba(102, 126, 234, 0.3) !important;
      border: 1px solid rgba(255, 255, 255, 0.2) !important;

      &:hover {
        background: linear-gradient(135deg, #7c8ef8 0%, #8a5cb8 100%) !important;
      }
    }

    // 确保图标和文字都是白色
    .el-icon,
    span {
      color: inherit !important;
    }
  }
}
</style>
