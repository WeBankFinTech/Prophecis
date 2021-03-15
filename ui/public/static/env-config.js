// 动态变量配置，该配置文件用于本地联调，线上环境后台配置文件会覆盖本文件，因此配置文件同步后台配置文件
var mlssEnvConfig = {
  'development': {
    'DI': {
      'defineImage': 'uat.sf.dockerhub.stgwebank/wedatasphere/prophecis', // 镜像基本路径
      'imageOption': ['tensorflow-1.5.0-py3', 'tensorflow-1.5.0-gpu-py3-wml-v1'],
      'definePython': '',
      'pythonOption': [
      ]
    },
    'AIDE': {
      'defineImage': 'uat.sf.dockerhub.stgwebank/wedatasphere/prophecis', // 镜像基本路径
      'imageOption': ['tensorflow-1.12.0-notebook-gpu-v0.4.0', 'tensorflow-1.12.0-notebook-gpu-v0.4.0-wml-v1']
    },
    'basisPlatform': {
    },
    'ns': { // 命名空间  平台空间名
      'platformNamespace': 'mlss-test'
    },
    'diApiVersion': 'v1', // di接口模块版本号
    'aideApiVersion': 'v1', // aide接口模块版本号
    'ccApiVersion': 'v1' // cc接口模块版本号
  },
  'production': {
    'DI': {
      'defineImage': 'uat.sf.dockerhub.stgwebank/wedatasphere/prophecis', // 镜像基本路径
      'imageOption': ['tensorflow-1.5.0-py3', 'tensorflow-1.5.0-gpu-py3-wml-v1'],
      'definePython': '',
      'pythonOption': [
      ]
    },
    'AIDE': {
      'defineImage': 'uat.sf.dockerhub.stgwebank/wedatasphere/prophecis', // 镜像基本路径
      'imageOption': ['tensorflow-1.12.0-notebook-gpu-v0.4.0', 'tensorflow-1.12.0-notebook-gpu-v0.4.0-wml-v1']
    },
    'basisPlatform': {
    },
    'ns': {
      'platformNamespace': 'mlss-test'
    },
    'diApiVersion': 'v1',
    'aideApiVersion': 'v1',
    'ccApiVersion': 'v1',
    'uiServer': 'bdp-ui',
    'filterUiServer': true
  }
}
