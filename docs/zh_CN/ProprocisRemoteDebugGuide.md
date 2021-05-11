### Proprocis Remote Debug Guide

[TOC]

本文以DI Rest服务为例，使用DLV进行远程调试：

#### 1、镜像编译

重新编译镜像，Base镜像加上包含dlv及golang：

```dockerfile
FROM golang
RUN go get -u github.com/go-delve/delve/cmd/dlv

USER root

ADD bin/main /
ADD certs/* /etc/ssl/dlaas/
ADD ./* /MLSS-DI-GPU/restapi/
RUN chmod 755 /main

ENTRYPOINT ["/main"]
```



#### 2、Deployment Yaml修改

修改Deployment的启动命令：

- dlv --listen为暴露的远程调试端口:

```yaml
	spec:
		containers:
        	command: ["/bin/sh", "-c"]
        	args: ["dlv --listen=:40000 --headless=true --api-version=2 exec /main -- --DLASS_PORT=8080"]

```



#### 3、暴露端口

增加NodePort，暴露Debug 40000端口：

```yaml
apiVersion: v1
kind: Service
metadata:
  name: di-restapi
  namespace: {{.Values.namespace}}
  labels:
    service: di-restapi
    environment: {{.Values.env}}
spec:
  type: NodePort
  ports:
  - name: ffdl
    port: 80
    targetPort: 8080
    nodePort: {{.Values.restapi.port}}
  - name: di-debug-port
    port: 40000
    targetPort: 40000
    nodePort: 30961
  selector:
    service: di-restapi
```



#### 4、IDE Remote Debug

VSCode和Goland配置Remote Debug连接40000端口即可，服务需要在IDE Remote Debug触发后才启动。