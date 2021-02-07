# Prophecis  一站式机器学习平台

Prophecis 是微众银行自研的一站式机器学习平台，集成多种开源机器学习框架，具备机器学习计算集群的多租户管理能力，提供生产环境全栈化容器部署与管理服务。

## Architecture
- #### 整体架构
  ![Prophecis](https://github.com/WeBankFinTech/Prophecis/blob/master/docs/zh_CN/image/Prophecis%E6%95%B4%E4%BD%93%E6%9E%B6%E6%9E%84.png)
  <p align="center">图1 Prophecis整体架构</p>
  
  **Prophecis主要包含5个关键服务：**

  - **Prophecis Machine Learning Flow**：机器学习分布式建模工具，具备单机和分布式模式模型训练能力，支持Tensorflow、Pytorch、xgboost等多种机器学习框架，支持从机器学习建模到部署的完整Pipeline；

  - **Prophecis MLLabis**：机器学习开发探索工具，提供开发探索服务，是一款基于Jupyter Lab的在线IDE，同时支持GPU及Hadoop集群的机器学习建模任务，支持Python、R、Julia多种语言，集成Debug、TensorBoard多种插件；

  - **Prophecis Model Factory**：机器学习模型工厂，提供机器学习模型存储、模型部署测试、模型管理等服务；

  - **Prophecis Data Factory**：机器学习数据工厂，提供特征工程工具、数据标注工具和物料管理等服务；

  - **Prophecis Application Factory**：机器学习应用工厂，由微众银行大数据平台团队和AI部门联合共建，基于青云(QingCloud)开源的Kubesphere定制开发，提供CI/CD和DevOps工具，GPU集群的监控及告警能力。

- #### 功能特色

  ![Prophecis功能特色](https://github.com/WeBankFinTech/Prophecis/blob/master/docs/zh_CN/image/Prophecis%E5%8A%9F%E8%83%BD%E7%89%B9%E8%89%B23.jpg)

  **<p align="center">图2 Prophecis功能特色</p>**
- **全生命周期的机器学习体验**：Prophecis的 MLFlow 通过 AppJoint 可以接入到 DataSphere Stdudio 的工作流中，支持从数据上传、数据预处理、特征工程、模型训练、模型评估到模型发布的机器学习全流程；

  ![DSS-Prophecis](https://github.com/WeBankFinTech/Prophecis/blob/master/docs/zh_CN/image/DSS-Prophecis.gif)
   **<p align="center">图3 Prophecis对接DSS功能展示</p>**

- **一键式的模型部署服务**：Prophecis MF 支持将Prophecis Machine Learning Flow、Prophecis MLLabis 生成的训练模型一键式发布为 Restful API 或者 RPC 接口，实现模型到业务的无缝衔接；

- **机器学习应用部署、运维、实验的综合管理平台**：基于社区开源方案定制，提供完整的、可靠的、高度灵活的企业级机器学习应用发布、监控、服务治理、日志收集和查询等管理工具，全方位实现对机器学习应用的管控，满足企业对于机器学习应用在线上生产环境的所有工作要求。

## Quick Start Guide
- 快速部署Prophecis服务，请参考 [Quick Start Guide](https://github.com/WeBankFinTech/Prophecis/blob/master/docs/zh_CN/QuickStartGuide.md) 文档。

- 关于配置解释，可参考[Quick Start Guide](https://github.com/WeBankFinTech/Prophecis/blob/master/docs/zh_CN/QuickStartGuide.md) 中的关键配置解释部分。

## Developing
- 编译Prophecis，请参考 [Develop Guide](https://github.com/WeBankFinTech/Prophecis/blob/master/docs/zh_CN/DevelopGuide.md)  文档。

## Roadmap
- 关于Prophecis后续的Roadmap，可查看 [Roadmap](https://github.com/WeBankFinTech/Prophecis/blob/master/docs/zh_CN/Roadmap.md) 文档，欢迎大家持续关注！

## Contributing

非常欢迎广大的社区伙伴给我们贡献新代码！

## Communication

如果您想得到最快的响应，请给我们提issue，或者您也可以扫码进群：

![Communication](https://github.com/WeBankFinTech/Prophecis/blob/master/docs/zh_CN/image/Communication.png)

## License
Prophecis  is under the Apache 2.0 license. See the LICENSE file for details.
