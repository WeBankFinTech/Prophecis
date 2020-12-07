/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;

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
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `name_UNIQUE` (`name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=49 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;

-- 数据导出被取消选择。

-- 导出  表 mlss_gzpc_bdap_uat_01.t_group_namespace 结构
CREATE TABLE IF NOT EXISTS `t_group_namespace` (
  `id` bigint(255) NOT NULL AUTO_INCREMENT,
  `group_id` bigint(255) DEFAULT NULL,
  `namespace_id` bigint(20) DEFAULT NULL COMMENT '命名空间id',
  `namespace` varchar(100) NOT NULL,
  `remarks` varchar(200) DEFAULT NULL,
  `enable_flag` tinyint(4) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `group_id` (`group_id`,`namespace_id`)
) ENGINE=InnoDB AUTO_INCREMENT=191 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;

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
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `group_id` (`group_id`,`storage_id`)
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;

-- 数据导出被取消选择。

-- 导出  表 mlss_gzpc_bdap_uat_01.t_keypair 结构
CREATE TABLE IF NOT EXISTS `t_keypair` (
  `id` bigint(255) NOT NULL AUTO_INCREMENT,
  `name` varchar(45) NOT NULL,
  `api_key` varchar(100) NOT NULL,
  `secret_key` varchar(100) NOT NULL,
  `remarks` varchar(200) DEFAULT NULL,
  `enable_flag` tinyint(4) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `name_UNIQUE` (`name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;

-- 数据导出被取消选择。

-- 导出  表 mlss_gzpc_bdap_uat_01.t_namespace 结构
CREATE TABLE IF NOT EXISTS `t_namespace` (
  `id` bigint(255) NOT NULL AUTO_INCREMENT,
  `namespace` varchar(100) NOT NULL,
  `remarks` varchar(200) DEFAULT NULL,
  `enable_flag` tinyint(4) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `namespace_uindex` (`namespace`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=75 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;

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
) ENGINE=InnoDB AUTO_INCREMENT=41 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT COMMENT='权限表';

-- 数据导出被取消选择。

-- 导出  表 mlss_gzpc_bdap_uat_01.t_role 结构
CREATE TABLE IF NOT EXISTS `t_role` (
  `id` bigint(255) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `remarks` varchar(200) DEFAULT NULL,
  `enable_flag` tinyint(4) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `name_UNIQUE` (`name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;

-- 数据导出被取消选择。

-- 导出  表 mlss_gzpc_bdap_uat_01.t_role_permission 结构
CREATE TABLE IF NOT EXISTS `t_role_permission` (
  `id` bigint(255) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `role_id` int(11) NOT NULL COMMENT '角色id',
  `permission_id` int(11) NOT NULL COMMENT '权限id',
  `enable_flag` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态，1：正常，0：冻结',
  `operate_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '最后一次更新的时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `role_permission_UNIQUE` (`role_id`,`permission_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=94 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT COMMENT='角色权限关联表';

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
) ENGINE=InnoDB AUTO_INCREMENT=544 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;

-- 数据导出被取消选择。

-- 导出  表 mlss_gzpc_bdap_uat_01.t_superadmin 结构
CREATE TABLE IF NOT EXISTS `t_superadmin` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `name` varchar(100) NOT NULL COMMENT '用户名',
  `remarks` varchar(200) DEFAULT NULL COMMENT '备注',
  `enable_flag` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态，1：正常，0：冻结',
  `operate_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '最后一次更新的时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `role_permission_UNIQUE` (`name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT COMMENT='超级管理员表';

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
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `name_UNIQUE` (`name`) USING BTREE,
  UNIQUE KEY `t_user_uid_uindex` (`uid`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=24 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;

-- 数据导出被取消选择。

-- 导出  表 mlss_gzpc_bdap_uat_01.t_user_group 结构
CREATE TABLE IF NOT EXISTS `t_user_group` (
  `id` bigint(255) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(255) NOT NULL,
  `role_id` bigint(255) NOT NULL,
  `group_id` bigint(255) NOT NULL,
  `remarks` varchar(200) DEFAULT NULL,
  `enable_flag` tinyint(4) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `user_id_UNIQUE` (`user_id`,`group_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=25 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;

-- 数据导出被取消选择。

/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IF(@OLD_FOREIGN_KEY_CHECKS IS NULL, 1, @OLD_FOREIGN_KEY_CHECKS) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
