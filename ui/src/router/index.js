import Vue from 'vue'
import VueRouter from 'vue-router'
Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    redirect: { name: 'home' }
  },
  {
    path: '/mlFlow',
    name: 'mlFlow',
    component: () =>
      import('../views/flow/process/formDssFlow.vue')
  },
  {
    path: '/noAccess',
    name: 'noAccess',
    component: () =>
      import('../views/NoAccess.vue')
  },
  {
    path: '/login',
    name: 'login',
    component: () =>
      import('../views/Login.vue')
  },
  {
    path: '/404',
    name: 'pageNotFound',
    component: () =>
      import('../views/404.vue')
  },
  {
    path: '*',
    component: () =>
      import('../views/404.vue')
  },
  {
    path: '/',
    component: () =>
      import('../views/Layout.vue'),
    children: [{
      path: '/home',
      name: 'home',
      meta: {
        login: true
      },
      component: () =>
        import('../views/Home.vue')
    },
    {
      path: '/experiment',
      name: 'experiment',
      meta: {
        login: true,
        groupName: 'modelTraining'
      },
      component: () =>
        import('../views/experiment/index.vue')
    },
    {
      path: '/experiment/flow',
      name: 'experimentFlow',
      meta: {
        login: true,
        title: 'flow.workflow',
        groupName: 'modelTraining'
      },
      component: () =>
        import('../views/flow/workflow.vue')
    },
    {
      path: '/experiment/addRecord/:mlflow_exp_id',
      name: 'addExperimentRecord',
      meta: {
        login: true
      },
      component: () =>
        import('../views/experiment/addRecord.vue')
    },
    {
      path: '/expExeRecord',
      name: 'expExeRecord',
      meta: {
        login: true,
        groupName: 'modelTraining'
      },
      component: () =>
        import('../views/expExeRecord/index.vue')
    },
    {
      path: '/jobExeRecord',
      name: 'jobExeRecord',
      meta: {
        login: true,
        groupName: 'modelTraining'
      },
      component: () =>
        import('../views/jobExeRecord/index.vue')
    },
    {
      path: '/jobExeRecord/:modelId',
      name: 'jobExeRecordDetail',
      meta: {
        login: true,
        groupName: 'modelTraining',
        title: 'DI.viewlog'
      },
      props: true,
      component: () =>
        import('../views/jobExeRecord/DetailLog')
    },
    {
      path: '/AIDE',
      name: 'AIDE',
      meta: {
        login: true
      },
      component: () =>
        import('../views/AIDE/index.vue')
    },
    {
      path: '/model/modelList',
      name: 'modelList',
      meta: {
        login: true,
        groupName: 'model'
      },
      component: () =>
        import('../views/modelFactory/model/modelList.vue')
    },
    {
      path: '/model/versionList/:modelId',
      name: 'modelVersionList',
      meta: {
        login: true,
        groupName: 'model',
        title: 'modelList.versionList'
      },
      component: () =>
        import('../views/modelFactory/model/versionList')
    },
    {
      path: '/model/serviceList',
      name: 'serviceList',
      meta: {
        login: true,
        groupName: 'model'
      },
      component: () =>
        import('../views/modelFactory/service/serviceList.vue')
    },
    {
      path: '/model/serviceList/:namespace/:serviceName',
      name: 'containerList',
      meta: {
        login: true,
        groupName: 'model',
        title: 'serviceList.containerList'
      },
      component: () =>
        import('../views/modelFactory/service/ContainerList')
    },
    {
      path: '/model/image',
      name: 'modelImage',
      meta: {
        login: true,
        groupName: 'model'
      },
      component: () =>
        import('../views/modelFactory/image/index.vue')
    },
    {
      path: '/model/report',
      name: 'modelReport',
      meta: {
        login: true,
        groupName: 'model'
      },
      component: () =>
        import('../views/modelFactory/reportList/index')
    },
    {
      path: '/model/reportVersion/:reportId',
      name: 'reportVersion',
      meta: {
        login: true,
        groupName: 'model',
        title: 'modelList.versionList'
      },
      component: () =>
        import('../views/modelFactory/reportList/versionList.vue')
    },
    {
      path: '/help',
      name: 'help',
      meta: {
        login: true
      },
      component: () =>
        import('../views/Help.vue')
    }
    ]
  },
  {
    path: '/manage',
    component: () =>
      import('../views/Layout.vue'),
    children: [{
      path: '/manage/basisPlatform',
      name: 'basisPlatform',
      meta: {
        login: true,
        permission: true,
        groupName: 'manage'
      },
      component: () =>
        import('../views/manage/BasisPlatform.vue')
    }, {
      path: '/manage/user',
      name: 'user',
      meta: {
        login: true,
        permission: true,
        groupName: 'manage'
      },
      component: () =>
        import('../views/manage/user/User.vue')
    },
    {
      path: '/manage/user/:userId',
      name: 'settingUserGroup',
      meta: {
        login: true,
        permission: true,
        groupName: 'manage',
        title: 'user.userGroupSettings'
      },
      props: true,
      component: () =>
        import('../views/manage/user/SettingGroup.vue')
    },
    {
      path: '/manage/userGroup',
      name: 'userGroup',
      meta: {
        login: true,
        permission: true,
        groupName: 'manage'
      },
      component: () =>
        import('../views/manage/UserGroup.vue')
    },
    {
      path: '/manage/user/proxy',
      name: 'settingUserProxy',
      meta: {
        login: true,
        permission: true,
        groupName: 'manage',
        title: 'user.proxyUserSetting'
      },
      props: true,
      component: () =>
        import('../views/manage/user/SettingProxy.vue')
    },
    {
      path: '/manage/namespace',
      name: 'namespace',
      meta: {
        login: true,
        permission: true,
        groupName: 'manage'
      },
      component: () =>
        import('../views/manage/namespace/Namespace.vue')
    },
    {
      path: '/manage/namespace/:namespace',
      name: 'settingNamespaceGroup',
      meta: {
        login: true,
        permission: true,
        groupName: 'manage',
        title: 'user.userGroupSettings'
      },
      props: true,
      component: () =>
        import('../views/manage/namespace/SettingGroup.vue')
    },
    {
      path: '/manage/machineLabel',
      name: 'machineLabel',
      meta: {
        login: true,
        permission: true,
        groupName: 'manage'
      },
      component: () =>
        import('../views/manage/MachineLabel.vue')
    },
    {
      path: '/manage/data',
      name: 'dataManage',
      meta: {
        login: true,
        permission: true,
        groupName: 'manage'
      },
      component: () =>
        import('../views/manage/data/Data')
    },
    {
      path: '/manage/data/:storageId',
      name: 'settingDataGroup',
      meta: {
        login: true,
        permission: true,
        groupName: 'manage',
        title: 'user.userGroupSettings'
      },
      props: true,
      component: () =>
        import('../views/manage/data/SettingGroup.vue')
    }
    ]
  }
]
const router = new VueRouter({
  routes
})

export default router
