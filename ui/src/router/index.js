import Vue from 'vue'
import VueRouter from 'vue-router'
Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    redirect: { name: 'home' }
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
      path: '/DI',
      name: 'DI',
      meta: {
        login: true
      },
      component: () =>
        import('../views/DI/index.vue')
    },
    {
      path: '/DI/log',
      name: 'DILogDetail',
      meta: {
        login: true,
        groupName: 'DI',
        title: 'DI.viewlog'
      },
      component: () =>
        import('../views/DI/DetailLog.vue')
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
