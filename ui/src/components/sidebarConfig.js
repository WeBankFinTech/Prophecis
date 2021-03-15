
export default {
  menuList: [
    {
      index: '/home',
      title: 'home.home',
      icon: 'icon-shouye'
    },
    {
      index: '/DI',
      title: 'trainingJob',
      icon: 'icon-moxingshu'
    },
    {
      index: '/AIDE',
      title: 'Notebook',
      icon: 'icon-jichengkaifa'
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
