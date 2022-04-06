
/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;


-- 导出 mlss_gzpc_bdap_uat_01 的数据库结构
CREATE DATABASE IF NOT EXISTS `mlss_gzpc_bdap_uat_01` /*!40100 DEFAULT CHARACTER SET utf8 */;
USE `mlss_gzpc_bdap_uat_01`;

-- 导出  表 mlss_gzpc_bdap_uat_01.t_experiment 结构
CREATE TABLE IF NOT EXISTS `t_experiment` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `exp_name` varchar(50) COLLATE utf8_bin DEFAULT NULL,
  `exp_desc` varchar(255) COLLATE utf8_bin DEFAULT NULL,
  `group_name` varchar(255) COLLATE utf8_bin DEFAULT NULL,
  `mlflow_exp_id` varchar(255) COLLATE utf8_bin DEFAULT NULL,
  `dss_project_id` bigint(20) DEFAULT NULL,
  `dss_project_version_id` bigint(20) DEFAULT NULL,
  `dss_bml_last_version` varchar(255) COLLATE utf8_bin DEFAULT NULL,
  `dss_workspace_id` bigint(20) DEFAULT NULL,
  `dss_flow_id` bigint(255) DEFAULT NULL,
  `dss_flow_last_version` varchar(50) COLLATE utf8_bin DEFAULT NULL,
  `exp_create_time` datetime DEFAULT NULL,
  `exp_modify_time` datetime DEFAULT NULL,
  `exp_create_user_id` bigint(20) DEFAULT NULL,
  `exp_modify_user_id` bigint(20) DEFAULT NULL,
  `enable_flag` tinyint(4) NOT NULL DEFAULT '1',
  `create_type` varchar(50) COLLATE utf8_bin DEFAULT NULL,
  `dss_dss_project_id` bigint(20) DEFAULT NULL,
  `dss_dss_project_name` varchar(50) COLLATE utf8_bin DEFAULT NULL,
  `dss_dss_flow_version` varchar(50) COLLATE utf8_bin DEFAULT NULL,
  `dss_dss_flow_name` varchar(50) COLLATE utf8_bin DEFAULT NULL,
  `dss_dss_flow_project_version_id` bigint(20) DEFAULT NULL,
  `experiment_type` varchar(50) COLLATE utf8_bin DEFAULT NULL COMMENT '可选以下值： WTSS/DSS/MLFlow',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1149 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- 数据导出被取消选择。

-- 导出  表 mlss_gzpc_bdap_uat_01.t_experiment_run 结构
CREATE TABLE IF NOT EXISTS `t_experiment_run` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `enable_flag` tinyint(1) DEFAULT '1',
  `exp_id` bigint(20) DEFAULT NULL,
  `dss_exec_id` varchar(255) DEFAULT NULL,
  `dss_task_id` bigint(20) DEFAULT NULL,
  `dss_flow_id` bigint(20) DEFAULT NULL,
  `dss_flow_last_version` varchar(255) DEFAULT NULL,
  `exp_exec_type` varchar(255) DEFAULT NULL,
  `exp_exec_status` varchar(255) DEFAULT NULL,
  `exp_run_create_time` datetime DEFAULT NULL,
  `exp_run_end_time` datetime DEFAULT NULL,
  `exp_run_create_user_id` bigint(20) DEFAULT NULL,
  `exp_run_modify_user_id` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=18459 DEFAULT CHARSET=utf8;

-- 数据导出被取消选择。

-- 导出  表 mlss_gzpc_bdap_uat_01.t_experiment_tag 结构
CREATE TABLE IF NOT EXISTS `t_experiment_tag` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `exp_id` bigint(20) DEFAULT NULL,
  `exp_tag` varchar(20) COLLATE utf8_bin DEFAULT NULL,
  `exp_tag_create_time` datetime DEFAULT NULL,
  `exp_tag_create_user_id` bigint(20) DEFAULT NULL,
  `enable_flag` tinyint(4) DEFAULT '1',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=70 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- 数据导出被取消选择。

-- 导出  表 mlss_gzpc_bdap_uat_01.t_group 结构
CREATE TABLE IF NOT EXISTS `t_group` (
  `id` bigint(255) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `group_type` varchar(100) NOT NULL DEFAULT 'SYSTEM' COMMENT '资源组类型',
  `subsystem_id` bigint(255) DEFAULT '0' COMMENT '子系统ID',
  `subsystem_name` varchar(255) DEFAULT NULL COMMENT '子系统名',
  `remarks` varchar(200) DEFAULT NULL,
  `enable_flag` tinyint(4) NOT NULL DEFAULT '1',
  `department_id` bigint(255) DEFAULT '0' COMMENT '部门ID',
  `department_name` varchar(255) DEFAULT NULL COMMENT '部门名称',
  `cluster_name` varchar(100) NOT NULL DEFAULT 'BDAP' COMMENT '集群环境',
  `rmb_idc` varchar(100) DEFAULT NULL,
  `rmb_dcn` varchar(100) DEFAULT NULL,
  `service_id` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name_UNIQUE` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=545 DEFAULT CHARSET=utf8;

-- 数据导出被取消选择。

-- 导出  表 mlss_gzpc_bdap_uat_01.t_group_namespace 结构
CREATE TABLE IF NOT EXISTS `t_group_namespace` (
  `id` bigint(255) NOT NULL AUTO_INCREMENT,
  `group_id` bigint(255) DEFAULT NULL,
  `namespace_id` bigint(20) DEFAULT NULL COMMENT '命名空间id',
  `namespace` varchar(100) NOT NULL,
  `remarks` varchar(200) DEFAULT NULL,
  `enable_flag` tinyint(4) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`),
  UNIQUE KEY `group_id` (`group_id`,`namespace_id`)
) ENGINE=InnoDB AUTO_INCREMENT=305 DEFAULT CHARSET=utf8;

-- 数据导出被取消选择。

-- 导出  表 mlss_gzpc_bdap_uat_01.t_group_storage 结构
CREATE TABLE IF NOT EXISTS `t_group_storage` (
  `id` bigint(255) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `group_id` bigint(20) NOT NULL COMMENT '资源组id',
  `storage_id` bigint(20) DEFAULT NULL COMMENT '存储路径id',
  `path` varchar(100) NOT NULL COMMENT '存储路径',
  `permissions` varchar(100) NOT NULL DEFAULT 'r' COMMENT '访问权限',
  `remarks` varchar(200) DEFAULT NULL COMMENT '备注',
  `enable_flag` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态，1：正常，0：冻结',
  `type` varchar(255) DEFAULT NULL COMMENT '数据路径类别',
  PRIMARY KEY (`id`),
  UNIQUE KEY `group_id` (`group_id`,`storage_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1246 DEFAULT CHARSET=utf8;

-- 数据导出被取消选择。

-- 导出  表 mlss_gzpc_bdap_uat_01.t_image 结构
CREATE TABLE IF NOT EXISTS `t_image` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `image_name` varchar(200) COLLATE utf8_bin DEFAULT NULL,
  `model_version_id` bigint(20) DEFAULT NULL,
  `user_id` bigint(20) DEFAULT NULL,
  `group_id` bigint(20) DEFAULT NULL,
  `creation_timestamp` datetime DEFAULT NULL,
  `last_updated_timestamp` datetime DEFAULT NULL,
  `status` varchar(50) COLLATE utf8_bin DEFAULT NULL,
  `remarks` varchar(200) COLLATE utf8_bin DEFAULT NULL,
  `enable_flag` tinyint(4) DEFAULT NULL,
  `msg` varchar(500) COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=167 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- 数据导出被取消选择。

-- 导出  表 mlss_gzpc_bdap_uat_01.t_keypair 结构
CREATE TABLE IF NOT EXISTS `t_keypair` (
  `id` bigint(255) NOT NULL AUTO_INCREMENT,
  `name` varchar(45) NOT NULL,
  `api_key` varchar(100) NOT NULL,
  `secret_key` varchar(100) NOT NULL,
  `super_admin` tinyint(4) NOT NULL DEFAULT '0',
  `remarks` varchar(200) DEFAULT NULL,
  `enable_flag` tinyint(4) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name_UNIQUE` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8;

-- 数据导出被取消选择。

-- 导出  表 mlss_gzpc_bdap_uat_01.t_model 结构
CREATE TABLE IF NOT EXISTS `t_model` (
  `id` bigint(255) NOT NULL AUTO_INCREMENT,
  `model_name` varchar(100) COLLATE utf8_bin DEFAULT NULL,
  `model_type` varchar(100) COLLATE utf8_bin DEFAULT NULL,
  `position` int(11) DEFAULT NULL,
  `creation_timestamp` datetime DEFAULT NULL,
  `group_id` bigint(20) DEFAULT NULL,
  `user_id` bigint(20) DEFAULT NULL,
  `service_id` bigint(20) DEFAULT NULL,
  `reamrk` varchar(200) COLLATE utf8_bin DEFAULT NULL,
  `enable_flag` tinyint(4) DEFAULT '1',
  `update_timestamp` datetime DEFAULT NULL,
  `model_latest_version_id` bigint(20) DEFAULT NULL,
  `model_usage` varchar(100) COLLATE utf8_bin DEFAULT NULL,
  `model_source` varchar(100) COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=392 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- 数据导出被取消选择。

-- 导出  表 mlss_gzpc_bdap_uat_01.t_modelversion 结构
CREATE TABLE IF NOT EXISTS `t_modelversion` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `version` varchar(100) COLLATE utf8_bin DEFAULT NULL,
  `creation_timestamp` datetime DEFAULT NULL,
  `latest_flag` tinyint(4) DEFAULT '1',
  `source` varchar(100) COLLATE utf8_bin DEFAULT NULL,
  `filepath` varchar(500) COLLATE utf8_bin DEFAULT NULL,
  `model_id` bigint(100) DEFAULT NULL,
  `enable_flag` tinyint(4) DEFAULT '1',
  `file_name` varchar(100) COLLATE utf8_bin DEFAULT NULL,
  `params` varchar(1000) COLLATE utf8_bin DEFAULT NULL,
  `training_id` varchar(1000) COLLATE utf8_bin DEFAULT NULL,
  `training_flag` tinyint(4) DEFAULT NULL,
  `push_id` bigint(20) DEFAULT NULL COMMENT 'push event 表id',
  `push_timestamp` datetime DEFAULT NULL,
  `deployment_status` varchar(100) COLLATE utf8_bin DEFAULT NULL,
  `deployment_platform` varchar(100) COLLATE utf8_bin DEFAULT NULL COMMENT '枚举值：CMSS、WCS、Prophecis',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2379 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- 数据导出被取消选择。

-- 导出  表 mlss_gzpc_bdap_uat_01.t_namespace 结构
CREATE TABLE IF NOT EXISTS `t_namespace` (
  `id` bigint(255) NOT NULL AUTO_INCREMENT,
  `namespace` varchar(100) NOT NULL,
  `remarks` varchar(200) DEFAULT NULL,
  `enable_flag` tinyint(4) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `namespace_uindex` (`namespace`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=113 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;

-- 数据导出被取消选择。

-- 导出  表 mlss_gzpc_bdap_uat_01.t_permission 结构
CREATE TABLE IF NOT EXISTS `t_permission` (
  `id` bigint(255) NOT NULL AUTO_INCREMENT COMMENT '权限id',
  `name` varchar(100) NOT NULL COMMENT '权限名称',
  `url` varchar(100) NOT NULL COMMENT '请求的url',
  `method` varchar(10) NOT NULL COMMENT '请求的方法',
  `remarks` varchar(200) DEFAULT NULL COMMENT '备注',
  `enable_flag` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态，1：正常，0：冻结',
  `operate_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '最后一次更新的时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `name_url_method_UNIQUE` (`name`,`url`,`method`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=147 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT COMMENT='权限表';

-- 数据导出被取消选择。

-- 导出  表 mlss_gzpc_bdap_uat_01.t_proxy_user 结构
CREATE TABLE IF NOT EXISTS `t_proxy_user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) NOT NULL,
  `name` varchar(100) NOT NULL,
  `gid` bigint(20) NOT NULL,
  `uid` bigint(20) DEFAULT NULL,
  `token` varchar(100) DEFAULT NULL,
  `path` varchar(100) DEFAULT NULL,
  `remarks` varchar(200) DEFAULT NULL,
  `is_activated` tinyint(4) DEFAULT '1',
  `enable_flag` tinyint(4) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=30 DEFAULT CHARSET=utf8;

-- 数据导出被取消选择。

-- 导出  表 mlss_gzpc_bdap_uat_01.t_pushevent 结构
CREATE TABLE IF NOT EXISTS `t_pushevent` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `file_type` varchar(100) COLLATE utf8_bin NOT NULL COMMENT '枚举：MODEL、DATA',
  `file_id` bigint(20) NOT NULL COMMENT '类型ID，对应modelversion id或者report version id',
  `file_name` varchar(100) COLLATE utf8_bin NOT NULL,
  `fps_file_id` varchar(100) COLLATE utf8_bin NOT NULL COMMENT 'FPS返回的56位文件ID',
  `fps_hash_value` varchar(100) COLLATE utf8_bin NOT NULL COMMENT 'FPS返回的32位哈希值',
  `params` varchar(200) COLLATE utf8_bin DEFAULT NULL,
  `idc` varchar(200) COLLATE utf8_bin NOT NULL COMMENT 'RMB的IDC',
  `dcn` varchar(200) COLLATE utf8_bin NOT NULL COMMENT 'RMB的DCN',
  `status` varchar(50) COLLATE utf8_bin NOT NULL COMMENT '状态，RMB成功后为成功，类型为：Success/Failed',
  `enable_flag` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态，1：正常，0：冻结',
  `creation_timestamp` datetime NOT NULL,
  `update_timestamp` datetime NOT NULL,
  `rmb_s3path` varchar(100) COLLATE utf8_bin NOT NULL COMMENT 'upload RMB response to minio',
  `user_id` bigint(20) NOT NULL,
  `rmb_resp_file_name` varchar(100) COLLATE utf8_bin NOT NULL,
  `version` varchar(100) COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2132 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- 数据导出被取消选择。

-- 导出  表 mlss_gzpc_bdap_uat_01.t_report 结构
CREATE TABLE IF NOT EXISTS `t_report` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `report_name` varchar(100) COLLATE utf8_bin NOT NULL,
  `model_id` bigint(20) NOT NULL,
  `model_version_id` bigint(20) NOT NULL,
  `user_id` bigint(20) NOT NULL,
  `group_id` bigint(20) NOT NULL,
  `creation_timestamp` datetime NOT NULL,
  `update_timestamp` datetime NOT NULL,
  `push_status` varchar(100) COLLATE utf8_bin DEFAULT NULL,
  `enable_flag` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态，1：正常，0：冻结',
  `report_latest_version_id` bigint(20) NOT NULL COMMENT '关联最新report_version表id',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=36 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- 数据导出被取消选择。

-- 导出  表 mlss_gzpc_bdap_uat_01.t_reportversion 结构
CREATE TABLE IF NOT EXISTS `t_reportversion` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `report_name` varchar(100) COLLATE utf8_bin NOT NULL,
  `version` varchar(100) COLLATE utf8_bin NOT NULL,
  `file_name` varchar(100) COLLATE utf8_bin NOT NULL,
  `file_path` varchar(100) COLLATE utf8_bin NOT NULL,
  `source` varchar(100) COLLATE utf8_bin NOT NULL,
  `push_id` bigint(20) DEFAULT NULL,
  `creation_timestamp` datetime NOT NULL,
  `update_timestamp` datetime NOT NULL,
  `enable_flag` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态，1：正常，0：冻结',
  `report_id` bigint(20) NOT NULL COMMENT '关联report表',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=244 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- 数据导出被取消选择。

-- 导出  表 mlss_gzpc_bdap_uat_01.t_result 结构
CREATE TABLE IF NOT EXISTS `t_result` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `model_id` bigint(20) NOT NULL,
  `model_version_id` bigint(20) NOT NULL,
  `training_id` varchar(100) COLLATE utf8_bin NOT NULL,
  `result_msg` varchar(200) COLLATE utf8_bin NOT NULL,
  `enable_flag` tinyint(4) DEFAULT '1',
  `creation_timestamp` datetime NOT NULL,
  `update_timestamp` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- 数据导出被取消选择。

-- 导出  表 mlss_gzpc_bdap_uat_01.t_role 结构
CREATE TABLE IF NOT EXISTS `t_role` (
  `id` bigint(255) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `remarks` varchar(200) DEFAULT NULL,
  `enable_flag` tinyint(4) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name_UNIQUE` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=56 DEFAULT CHARSET=utf8;

-- 数据导出被取消选择。

-- 导出  表 mlss_gzpc_bdap_uat_01.t_role_permission 结构
CREATE TABLE IF NOT EXISTS `t_role_permission` (
  `id` bigint(255) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `role_id` int(11) NOT NULL COMMENT '角色id',
  `permission_id` int(11) NOT NULL COMMENT '权限id',
  `enable_flag` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态，1：正常，0：冻结',
  `operate_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '最后一次更新的时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `role_permission_UNIQUE` (`role_id`,`permission_id`)
) ENGINE=InnoDB AUTO_INCREMENT=308 DEFAULT CHARSET=utf8 COMMENT='角色权限关联表';

-- 数据导出被取消选择。

-- 导出  表 mlss_gzpc_bdap_uat_01.t_score 结构
CREATE TABLE IF NOT EXISTS `t_score` (
  `id` int(11) NOT NULL,
  `class` int(11) DEFAULT NULL,
  `score` int(11) DEFAULT NULL,
  `name` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- 数据导出被取消选择。

-- 导出  表 mlss_gzpc_bdap_uat_01.t_service 结构
CREATE TABLE IF NOT EXISTS `t_service` (
  `id` bigint(255) NOT NULL AUTO_INCREMENT,
  `service_name` varchar(100) COLLATE utf8_bin DEFAULT NULL,
  `type` varchar(100) COLLATE utf8_bin DEFAULT NULL,
  `cpu` double DEFAULT NULL,
  `memory` double DEFAULT NULL,
  `gpu` varchar(100) COLLATE utf8_bin DEFAULT NULL,
  `namespace` varchar(100) COLLATE utf8_bin DEFAULT NULL,
  `creation_timestamp` datetime DEFAULT NULL,
  `last_updated_timestamp` datetime DEFAULT NULL,
  `user_id` bigint(255) DEFAULT NULL,
  `group_id` bigint(255) DEFAULT NULL,
  `model_count` int(11) DEFAULT NULL,
  `log_path` varchar(100) COLLATE utf8_bin DEFAULT NULL,
  `remark` varchar(100) COLLATE utf8_bin DEFAULT NULL,
  `enable_flag` tinyint(4) DEFAULT '1',
  `endpoint_type` varchar(50) COLLATE utf8_bin DEFAULT NULL,
  `image_name` varchar(200) COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=279 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- 数据导出被取消选择。

-- 导出  表 mlss_gzpc_bdap_uat_01.t_service_modelversion 结构
CREATE TABLE IF NOT EXISTS `t_service_modelversion` (
  `id` bigint(255) NOT NULL AUTO_INCREMENT,
  `service_id` bigint(255) DEFAULT NULL,
  `modelversion_id` bigint(255) DEFAULT NULL,
  `enable_flag` tinyint(4) DEFAULT '1',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=278 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- 数据导出被取消选择。

-- 导出  表 mlss_gzpc_bdap_uat_01.t_storage 结构
CREATE TABLE IF NOT EXISTS `t_storage` (
  `id` bigint(255) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `path` varchar(100) NOT NULL COMMENT '存储路径',
  `remarks` varchar(200) DEFAULT NULL COMMENT '备注',
  `enable_flag` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态，1：正常，0：冻结',
  `type` varchar(255) DEFAULT NULL COMMENT '数据路径类别',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `path_UNIQUE` (`path`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=366 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;

-- 数据导出被取消选择。

-- 导出  表 mlss_gzpc_bdap_uat_01.t_superadmin 结构
CREATE TABLE IF NOT EXISTS `t_superadmin` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `name` varchar(100) NOT NULL COMMENT '用户名',
  `remarks` varchar(200) DEFAULT NULL COMMENT '备注',
  `enable_flag` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态，1：正常，0：冻结',
  `operate_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '最后一次更新的时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `role_permission_UNIQUE` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8 COMMENT='超级管理员表';

-- 数据导出被取消选择。

-- 导出  表 mlss_gzpc_bdap_uat_01.t_user 结构
CREATE TABLE IF NOT EXISTS `t_user` (
  `id` bigint(255) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `gid` bigint(255) DEFAULT '0' COMMENT 'gid',
  `uid` bigint(255) DEFAULT '0' COMMENT 'uid',
  `token` varchar(255) DEFAULT NULL COMMENT 'GUARDIAN_TOKEN',
  `type` varchar(255) DEFAULT 'SYSTEM' COMMENT '用户类型：SYSTEM、USER',
  `remarks` varchar(200) DEFAULT NULL,
  `enable_flag` tinyint(4) NOT NULL DEFAULT '1',
  `guid_check` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否检查guid(0:不检查,1:检查;默认为0)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name_UNIQUE` (`name`),
  UNIQUE KEY `t_user_uid_uindex` (`uid`)
) ENGINE=InnoDB AUTO_INCREMENT=541 DEFAULT CHARSET=utf8;

-- 数据导出被取消选择。

-- 导出  表 mlss_gzpc_bdap_uat_01.t_user_group 结构
CREATE TABLE IF NOT EXISTS `t_user_group` (
  `id` bigint(255) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(255) NOT NULL,
  `role_id` bigint(255) NOT NULL,
  `group_id` bigint(255) NOT NULL,
  `remarks` varchar(200) DEFAULT NULL,
  `enable_flag` tinyint(4) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_id_UNIQUE` (`user_id`,`group_id`)
) ENGINE=InnoDB AUTO_INCREMENT=672 DEFAULT CHARSET=utf8;

-- 数据导出被取消选择。

/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IF(@OLD_FOREIGN_KEY_CHECKS IS NULL, 1, @OLD_FOREIGN_KEY_CHECKS) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
