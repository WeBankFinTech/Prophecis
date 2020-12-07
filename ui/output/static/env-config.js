// 动态变量配置，该配置文件用于本地联调，线上环境后台配置文件会覆盖本文件，因此配置文件同步后台配置文件
var mlssEnvConfig = {
  'development': {
    'AIDE': {
      'defineImage': 'wedatasphere/prophecis', // 镜像基本路径
      'imageOption': ['tensorflow-1.12.0-notebook-gpu-v0.4.0', 'tensorflow-1.12.0-notebook-gpu-v0.4.0-wml-v1']
    },
    'basisPlatform': {
      'dashboardUrl': '',
      'grafanaUrl': '',
      'kibanaUrl': '',
      'prometheusUrl': ''
    },
    'ns': { // 命名空间  平台空间名
      'platformNamespace': 'mlss-test'
    },
    'diApiVersion': 'v1', // di接口模块版本号
    'aideApiVersion': 'v1', // aide接口模块版本号
    'ccApiVersion': 'v1' // cc接口模块版本号
  },
  'production': {
    'ssoLoginUrl': 'http://127.0.0.1:8080/cas/login',
    'ssoLogoutUrl': 'http://127.0.0.1:8080/cas/logout',
    'AIDE': {
      'defineImage': 'wedatasphere/prophecis',
      'imageOption': ['tensorflow-1.12.0-notebook-gpu-v0.4.0']
    },
    'basisPlatform': {
      'grafanaUrl': 'http://127.0.0.1:30780'
    },
    'ns': {
      'platformNamespace': 'prophecis-prod'
    },
    'diApiVersion': 'v1',
    'aideApiVersion': 'v1',
    'ccApiVersion': 'v1',
    'uiServer': 'bdp-ui',
    'filterUiServer': true
  }
}
