import k8sclusters from '@/views/K8s/k8s-clusters.vue'
import k8snodes from '@/views/K8s/k8s-nodes.vue'
import k8sworkload from '@/views/K8s/k8s-workloads.vue'
import k8snamespace from '@/views/K8s/k8s-namespace.vue'
import k8network from '@/views/K8s/k8s-network.vue'
import k8sconfig from '@/views/K8s/k8s-config.vue'
import k8storage from '@/views/K8s/k8s-storage'
import k8spod from '@/views/K8s/pods/k8s-pod.vue'
import k8sterminal from '@/views/K8s/pods/K8S-sterminal.vue'
import k8smonitoring from '@/views/K8s/nodes/k8s-monitoring.vue'
import k8sdetails from '@/views/K8s/clusters/K8sDetails.vue'
import k8snodedetails from '@/views/K8s/nodes/NodeDetails.vue'


const routes = [
    {
        path: '/k8s/list',
        component: k8sclusters,
        meta: {sTitle: '容器管理', tTitle: '集群管理'}
    },
    {
        path: '/k8s/cluster/:clusterId',
        component: k8sdetails,
        props: true,
        meta: {sTitle: '容器管理', tTitle: '集群详情'}
    },
    {
        path: '/k8s/cluster/:clusterId/node/:nodeName',
        component: k8snodedetails,
        props: true,
        meta: {sTitle: '容器管理', tTitle: '节点详情'}
    },
    {
        path: '/k8s/node',
        component: k8snodes,
        meta: {sTitle: '容器管理', tTitle: '节点管理'}
    },
    {
        path: '/k8s/namespace',
        component: k8snamespace,
        meta: {sTitle: '容器管理', tTitle: '命名空间'}
    },
    {
        path: '/k8s/workload',
        component: k8sworkload,
        meta: {sTitle: '容器管理', tTitle: '工作负载'}
    },
        {
        path: '/k8s/network',
        component: k8network,
        meta: {sTitle: '容器管理', tTitle: '网络管理'}
    },
        {
        path: '/k8s/config',
        component: k8sconfig,
        meta: {sTitle: '容器管理', tTitle: '配置管理'}
    },
        {
        path: '/k8s/storage',
        component: k8storage,
        meta: {sTitle: '容器管理', tTitle: '存储管理'}
    },
    {
        path: '/k8s/pod/:clusterId/:namespace/:podName',
        component: k8spod,
        props: true,
        meta: {sTitle: '容器管理', tTitle: '容器详情'}
    },
    {
        path: '/k8s/terminal/:clusterId/:namespace/:podName',
        component: k8sterminal,
        props: true,
        meta: {sTitle: '容器管理', tTitle: '容器终端'}
    },
    {
        path: '/k8s/monitoring',
        component: k8smonitoring,
        meta: {sTitle: '容器管理', tTitle: '监控仪表板'}
    },
    {
    path: '/k8s/crd',
    name: 'K8sCrd',
    component: () => import('@/views/K8s/k8s-crd.vue'),
    meta: { title: '自定义资源', icon: 'coin' }
    }
]

export default routes
