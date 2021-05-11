## Quick Start Guide

[TOC]


####  一、Prerequisites
- Kubernetes  1.18.6+
- Helm 3
- Nvidia-docker (可选)
- K8s-device-plugin (可选)
- Cuda (可选)

####  二、数据库更新
- 安装Mysql DB，使用MYSQL DB插入对应的表结构和更新API权限
- SQL文件位于./cc/sql目录下

```shell
source prophecis.sql #结构
source prophecis-data.sql #permission数据
```

####  三、部署Notebook Controller
- Notebook Controller可见(版本为0.4.0)：
  [Kubflow Notebook Controller](https://github.com/kubeflow/kubeflow/tree/master/components/notebook-controller)

- Prophecis对Kubeflow提供的Notebook Controller进行NodePort SVC的扩展，用于PySpark Session的创建，如需使用可从wedatasphere仓库拉取。

  ```shell
  cd ./helm-charts/notebook-controller
  helm install notebook-controller .
  ```

#### 四、部署MinIO及Log Collector

- MinIO
  ```shell
  cd ./helm-charts/MinioDeployment
  helm install minio-prophecis --namespace prophecis .
  ```

- Log Collector
    ```shell
    cd ./helm-charts/
    kubectl apply -f ./LogCollectorDS
    ```

#### 五、修改配置

- 修改ui nginx反向代理配置（values.yml文件，host ip + nodeport）：

```shell
cc:
	server: 127.0.0.1:30778
aide:
	server: 127.0.0.1:30778
di:
    server: 127.0.0.1:30778
```

- 修改数库配置（values.yml文件，db配置一致即可）：

```shell
cc:
  db:
    server: 127.0.0.1
    port: 3306
    name: prophecis_db
    user: prophecis
    pwd: prophecis
datasource:
  userName: prophecis
  userPwd: prophecis
  url: 127.0.0.1
  port: 3306
  db: prophecis_db
  encryptPwd: ""
  privKey: ""
```
- 登录采用通过LDAP方案，需要配置values中的ldap address和baseDN变量：
```shell
    ldap:
      address: ldap://127.0.0.1:1389/
      baseDN: dc=webank,dc=com
```

- 对运行服务的节点打上对应的label，GPU节点需要打上hardware-type=NVIDIAGPU；

```shell
kubectl label nodes prophecis01 mlss-node-role=platform
kubectl label nodes gpunode ardware-type=NVIDIAGPU
```

#### 六、管理员用户

- 默认的管理员用户为admin，账号密码通过values中的admin进行配置。
- 若需要新增SA用户，需要在t_superadmin和t_user中增加对应的用户记录；

#### 七、部署命令

```shell
kubectl create namespace prophecis

# Install using Helm 3 
helm install prophecis .
```

#### 八、关键配置解释
- platformNodeSelectors：prophecis服务运行label

- cc:

  - ldap：用于登录验证用户密码的LDAP Server;
  - db：
  - gateway：配置转发到各个服务相关

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
