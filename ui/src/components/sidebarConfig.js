
export default {
  menuList: [
    {
      index: '/home',
      title: 'home.home',
      icon: 'icon-shouye'
    },
    {
      index: 'modelTraining',
      title: 'DI.distributedModeling',
      icon: 'icon-moxingshu',
      subs: [
        {
          index: '/experiment',
          title: 'home.experimentList'
        },
        {
          index: '/expExeRecord',
          title: 'home.expExecutRecord'
        },
        {
          index: '/jobExeRecord',
          title: 'home.taskExecutRecord'
        }
      ]
    },
    {
      index: '/AIDE',
      title: 'Notebook',
      icon: 'icon-jichengkaifa'
    },
    {
      index: 'model',
      title: 'serviceList.modelFactory',
      icon: 'icon-moxingku',
      subs: [
        {
          index: '/model/modelList',
          title: 'modelList.modelList'
        },
        {
          index: '/model/image',
          title: 'image.imageList'
        },
        {
          index: '/model/serviceList',
          title: 'serviceList.modelService'
        },
        {
          index: '/model/report',
          title: 'report.reportList',
          icon: 'icon-jichengkaifa'
        }
      ]
    },
    {
      index: 'manage',
      title: 'management',
      icon: 'icon-tubiaozhizuomobanyihuifu-',
      subs: [
        {
          index: '/manage/basisPlatform',
          title: 'basisPlatform'
        }, {
          index: '/manage/user',
          title: 'userManager'
        }, {
          index: '/manage/userGroup',
          title: 'userGroupManager'
        }, {
          index: '/manage/namespace',
          title: 'namespaceManager'
        }, {
          index: '/manage/machineLabel',
          title: 'nodeLabelManager'
        }, {
          index: '/manage/data',
          title: 'dataManager'
        }
      ]
    },
    {
      index: '/help',
      title: 'help.help',
      icon: 'icon-bangzhu'
    }
  ]
}
