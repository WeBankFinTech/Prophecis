## Quick Start Guide

[TOC]


####  一、Prerequisites
- Kubernetes  1.15.6+
- Helm 3
- Nvidia-docker (可选)
- K8s-device-plugin (可选)
- Cuda (可选)

####  二、镜像拉取

```shell
docker pull wedatasphere/prophecis:ui-v0.1.0
docker pull wedatasphere/prophecis:mllabis-v0.1.0
docker pull wedatasphere/prophecis:cc-v0.1.0
docker pull wedatasphere/prophecis:cc-gateway-v0.1.0
docker pull wedatasphere/prophecis:notebook-controller-v0.1.0
```
####  三、数据库更新

```shell
source prophecis.sql #结构
source prophecis-data.sql #permission数据
```

#### 三、修改配置

- 修改ui nginx反向代理配置（values.yml文件，host ip + nodeport）：

```shell
cc:
	server: 127.0.0.1:30778
aide:
	server: 127.0.0.1:30778
```

- 修改数库配置（values.yml文件）：

```shell
cc:
  db:
    server: 127.0.0.1
    port: 3306
    name: prophecis_db
    user: prophecis
    pwd: password
```

- 对运行服务的节点打上对应的label

```shell
kubectl label nodes prophecis01 mlss-node-role=platform
```

#### 四、部署命令

```shell
kubectl create namespace prophecis

# Install using Helm 3 
helm install prophecis .
```

#### 五、关键配置解释
- platformNodeSelectors：prophecis服务运行label

- cc:

- aide：
  - startport & endPort:  notebook创建的server会从改nodeport范围获取对应的端口好，该端口用于连接Spark集群。
  - hadoop
    - enable:  是否开启Hadoop功能，若选择开启，会挂载以下配置路径。
    - installPath:  hadoop、spark安装的相关路径。
    - commonlibPath:  commonlib相关安装路径。
    - configPath:  hadoop、spark配置文件路径。
    -  javaPath: java安装文件路径。
    - sourceConfigPath:  notebook需要加载配置文件路径，该配置会在jupyter lab启动时加载。
    - hostFilePath:  notebook需要加载配置host文件配置，该配置会在jupyter lab启动时加载。
  - linkis：linkis对应的配置，包含Token和address。
  - livy： livy对应的配置。

- ui：

  - grafana: 基础控制台的grafana url配置
  - dashboard: 基础控制台的k8s dashboard url配置
  - prometheus: 基础控制台的k8s prometheus url配置 
  - kibana: 基础控制台的k8s kibana url配置 
  - cc & aide: cc&aide url配置，配置gateway地址和端口即可

- ##### 服务实时配置更新（部分配置需要重新部署）
 ```shell
#实时更新
helm upgrade prophecis 
 ```