[toc]

## 1. 配置结构

* Prophecis HelmChart
* Notebook Controller Helm Chart
## 2. Prophecis Helm Chart

### 2.1 目录结构

* 配置文件位置：**helm-chart/Prophecis**
* 文件内容：
    * values.yaml
    * templates目录
### 2.2 values配置

* namespace: 命名空间名字
* env: 生产环境
* platformNodeSelectors:
    * mlss-node-role: prophecis服务运行label
* envir:
* has_static_volumes:
* ceph_host_path:ceph存储位置
* server_ui_gateway: 网页访问网址
* server_repository: 镜像仓库
* database:
    * server: MySQLIP地址
    * port: MySQL端口号
    * name: MySQL数据库名
    * user: MySQL用户名
    * pwd: MySQL用户密码
    * encryptPwd:
    * privKey:
* log:
    * level: 日志级别
* cc:
    * ldap: 用于登录验证用户密码的LDAP Server
        * address:
        * baseDN:
* aide:
    * hadoop
        * enable:  是否开启Hadoop功能，若选择开启，会挂载以下配置路径
        * installPath:  hadoop、spark安装的相关路径
        * commonlibPath:  commonlib相关安装路径
        * configPath:  hadoop、spark配置文件路径
        * javaPath: java安装文件路径
        * sourceConfigPath:  notebook需要加载配置文件路径，该配置会在jupyter lab启动时加载
        * hostFilePath:  notebook需要加载配置host文件配置，该配置会在jupyter lab启动时加载
* ui:
    * service:
        * bdap:
            * nodePort: 网页访问端口
    * grafana:
        * url: 基础控制台的grafana url配置
    * dashboard:
        * url: 基础控制台的k8s dashboard url配置
    * prometheus:
        * utl: 基础控制台的k8s prometheus url配置
    * kibana:
        * url: 基础控制台的k8s kibana url配置
* livy:
    * address: livy对应的地址
* storage:
    * uploadHostPath: 模型主机上传位置
    * uploadContainerPath: 模型容器内位置
* persistent: 持久化位置
    * etcd:
        * path:
    * localstack:
        * path:
    * mongo:
        * path:

        
### 2.3 templates配置

* cc: 关于控制台的相关
    * gateway: 网关服务和容器部署
* di:
    * storage: 存储相关服务
    * restapi: 处理外部相关微服务请求
    * lcm: 任务队列管理
    * jm: job任务监测
    * trainer: 管理额模型训练任务
    * trainingdata: 获取模型训练中产生的日志
    * fluent-bit: 训练日志收集
* ingrastructure: 关于存储服务
    * mongo:
    * etcd:
* LogCollectorDS: 日志组件
* mllabis: jupyter lab 的相关
    * aide:
* persistent: 关于持久化
    * etcd:
    * mongo:
* ui: 关于网页界面展示
    * bdap-ui:

## 3. NotebookController Helm Chart

### 3.1 目录结构

* 配置文件位置：**helm-chart/notebook-controller**
* 文件内容：
    * values.yaml
    * templates目录：
### 3.2 values配置

* namespace: 命名空间名字
* platformNodeSelectors:
    * mlss-node-role: prophecis服务运行label
* services：
    * expose_node_port: false
* mllabis: jupyter lab 相关
    * image:
        * repostitory: 仓库位置
        * tag: 镜像标签
    * service:
        * type: 节点类型
        * port: 集群内部服务端口
        * targetPort: 容器访问端口
        * nodePort: 外部访问端口
    * controller:
        * notebook:
            * repository: 镜像仓库
            * tag: 镜像标签
        * meta:
            * repositroy:
            * tag
### 3.3 template配置

* metacontroller: CRD的Resouce资源配置
* notebook-controller:

