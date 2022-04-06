
/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;


-- 导出 mlss_gzpc_bdap_uat_01 的数据库结构
CREATE DATABASE IF NOT EXISTS `mlss_gzpc_bdap_uat_01` /*!40100 DEFAULT CHARACTER SET utf8 */;
USE `mlss_gzpc_bdap_uat_01`;

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

-- 正在导出表  mlss_gzpc_bdap_uat_01.t_group 的数据：~3 rows (大约)
DELETE FROM `t_group`;
/*!40000 ALTER TABLE `t_group` DISABLE KEYS */;
INSERT INTO `t_group` (`id`, `name`, `group_type`, `subsystem_id`, `subsystem_name`, `remarks`, `enable_flag`, `department_id`, `department_name`, `cluster_name`, `rmb_idc`, `rmb_dcn`, `service_id`) VALUES
	(542, 'gp-private-neiljianliu', 'PRIVATE', 0, '', '个人用户组', 1, 0, 'mlss', '', '', '', 0),
	(543, 'gp-private-alexwu', 'SYSTEM', 0, '', '', 1, 0, 'mlss', '', '', '', 0),
	(544, 'gp-private-hadoop', 'PRIVATE', 0, '', '个人用户组', 1, 0, 'mlss', '', '', '', 0);
/*!40000 ALTER TABLE `t_group` ENABLE KEYS */;

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

-- 正在导出表  mlss_gzpc_bdap_uat_01.t_group_storage 的数据：~2 rows (大约)
DELETE FROM `t_group_storage`;
/*!40000 ALTER TABLE `t_group_storage` DISABLE KEYS */;
INSERT INTO `t_group_storage` (`id`, `group_id`, `storage_id`, `path`, `permissions`, `remarks`, `enable_flag`, `type`) VALUES
	(1244, 543, 363, '/data/bdap-ss/mlss-data/alexwu', 'rw', '', 1, 'SYSTEM'),
	(1245, 544, 1, '/data/bdap-ss/mlss-data/hadoop', 'rw', '', 1, 'PRIVATE');
/*!40000 ALTER TABLE `t_group_storage` ENABLE KEYS */;

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

-- 正在导出表  mlss_gzpc_bdap_uat_01.t_keypair 的数据：~7 rows (大约)
DELETE FROM `t_keypair`;
/*!40000 ALTER TABLE `t_keypair` DISABLE KEYS */;
INSERT INTO `t_keypair` (`id`, `name`, `api_key`, `secret_key`, `super_admin`, `remarks`, `enable_flag`) VALUES
	(1, 'sample', 'xxx', 'yyy', 0, NULL, 1),
	(2, 'MLFLOW', 'MLFLOW', 'MLFLOW', 1, NULL, 1),
	(3, 'alexwu', 'test', 'test', 1, NULL, 1),
	(4, 'neiljianliu', 'neiljianliu', 'neiljianliu', 0, NULL, 1),
	(5, 'zhipingpan', 'zhipingpan', 'zhipingpan', 0, NULL, 1),
	(6, 'hduser0362', 'hduser0362', 'hduser0362', 1, NULL, 1),
	(7, 'QML-AUTH', 'QML-AUTH', 'QML-AUTH', 1, NULL, 1);
/*!40000 ALTER TABLE `t_keypair` ENABLE KEYS */;

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

-- 正在导出表  mlss_gzpc_bdap_uat_01.t_permission 的数据：~109 rows (大约)
DELETE FROM `t_permission`;
/*!40000 ALTER TABLE `t_permission` DISABLE KEYS */;
INSERT INTO `t_permission` (`id`, `name`, `url`, `method`, `remarks`, `enable_flag`, `operate_time`) VALUES
	(1, '获取指定ID的资源组详情', '/cc/v1/groups/id/*', 'GET', NULL, 1, '2019-03-22 17:28:55'),
	(2, '获取指定名称的资源组详情', '/cc/v1/groups/name/*', 'GET', NULL, 1, '2019-03-22 17:28:55'),
	(3, '获取指定ID的资源组的namespace列表', '/cc/v1/groups/*/namespaces', 'GET', NULL, 1, '2019-03-22 17:28:55'),
	(4, '获取指定ID的资源组的namespace详情', '/cc/v1/groups/*/namespaces/*', 'GET', NULL, 1, '2019-03-22 17:28:55'),
	(5, '获取指定ID的资源组的用户列表', '/cc/v1/groups/*/users', 'GET', NULL, 1, '2019-03-22 17:28:55'),
	(6, '新增资源组的用户', '/cc/v1/groups/users', 'POST', NULL, 1, '2019-03-22 17:28:55'),
	(7, '修改资源组的用户', '/cc/v1/groups/users', 'PUT', NULL, 1, '2019-03-22 17:28:55'),
	(8, '检查用户访问namespace的权限', '/cc/v1/auth/access/users/*/namespaces/*', 'GET', NULL, 1, '2019-03-22 17:28:56'),
	(9, 'Middleware调用获取权限的接口', '/cc/v1/sample', 'GET', NULL, 1, '2019-03-22 17:28:56'),
	(10, '校验用户共享存储数据路径权限', '/cc/v1/auth/access/groups/*/storages', 'GET', NULL, 1, '2019-04-16 17:22:17'),
	(11, '获取指定资源组id的存储', '/cc/v1/groups/id/*/storages', 'GET', NULL, 1, '2019-04-16 17:22:17'),
	(12, '获取指定资源组名的存储', '/cc/v1/groups/name/*/storages', 'GET', NULL, 1, '2019-04-16 17:22:17'),
	(13, '获取用户列表', '/cc/v1/users', 'GET', NULL, 0, '2019-04-17 14:39:58'),
	(14, '新增用户', '/cc/v1/users', 'POST', NULL, 0, '2019-04-17 14:39:58'),
	(15, '修改用户', '/cc/v1/users', 'PUT', NULL, 0, '2019-04-17 14:39:58'),
	(16, '根据用户id删除用户', '/cc/v1/users/id/*', 'DELETE', NULL, 1, '2019-04-17 14:39:58'),
	(17, '根据用户名删除用户', '/cc/v1/users/name/*', 'DELETE', NULL, 1, '2019-04-17 14:39:58'),
	(18, '获取指定ID的用户详情', '/cc/v1/users/id/*', 'GET', NULL, 1, '2019-04-17 14:39:58'),
	(19, '获取指定名称的用户详情', '/cc/v1/users/name/*', 'GET', NULL, 1, '2019-04-17 14:39:58'),
	(20, '删除资源组中指定ID的用户', '/cc/v1/groups/*/users/id/*', 'delete', NULL, 1, '2019-04-17 15:20:17'),
	(21, '删除资源组中指定用户名的用户', '/cc/v1/groups/*/users/name/*', 'delete', NULL, 1, '2019-04-17 15:20:49'),
	(22, '某用户访问指定用户名的用户资源', '/cc/v1/auth/access/users/*/users/*', 'GET', NULL, 1, '2019-05-06 15:16:57'),
	(23, '某用户访问指定用户名的共享存储数据路径权限', '/cc/v1/auth/access/users/*/storages', 'GET', NULL, 1, '2019-05-20 15:42:26'),
	(24, '某用户访问指定用户名的命名空间及共享存储数据路径权限', '/cc/v1/auth/access/users/*/namespaces/*/storages', 'GET', NULL, 1, '2019-05-20 20:38:50'),
	(25, '检查用户是否有权限访问指定user的资源', '/cc/v1/auth/access/admin/users/*', 'GET', NULL, 1, '2019-06-15 20:00:56'),
	(26, '检查用户是否有权限访问指定namespace中的指定user的资源', '/cc/v1/auth/access/admin/namespaces/*', 'GET', NULL, 1, '2019-06-15 20:30:15'),
	(27, '检查用户是否有权限访问指定namespace中的指定user的资源', '/cc/v1/auth/access/admin/namespaces/*/users/*', 'GET', NULL, 1, '2019-06-15 20:30:15'),
	(28, '检查用户是否有权限访问指定namespace中的指定notebook', '/cc/v1/auth/access/admin/namespaces/*/notebooks/*', 'GET', NULL, 1, '2019-06-15 20:30:15'),
	(29, '检查notebook提交对象', '/cc/v1/auth/access/check-notebook-request', 'GET', NULL, 1, '2019-06-17 19:19:45'),
	(30, '获取当前用户能访问的用户信息的用户列表', '/cc/v1/users/myUsers', 'GET', NULL, 1, '2019-07-12 10:57:18'),
	(31, '获取当前用户能访问的所有资源组的namespace', '/cc/v1/namespaces/myNamespace', 'GET', NULL, 1, '2019-07-12 10:57:18'),
	(32, '获取当前用户拥有指定权限的namespace列表', '/cc/v1/groups/users/roles/*/namespaces/clusterName/*', 'GET', NULL, 1, '2019-07-12 10:57:18'),
	(33, 'Middleware调用获取权限的接口', '/cc/v1/sample', 'POST', NULL, 1, '2019-07-17 16:42:29'),
	(34, '根据groupId和path获取storage', '/cc/v1/groups/group/storage', 'GET', NULL, 1, '2019-08-27 16:34:30'),
	(35, '根据groupId和userId获取user_group', '/cc/v1/groups/groupId/userId/group', 'GET', NULL, 1, '2019-08-28 16:54:10'),
	(36, '获取当前用户有权限访问的path列表', '/cc/v1/groups/group/storage/*', 'GET', NULL, 1, '2019-10-09 16:12:29'),
	(73, '根据用户名获取用户组', '/cc/v1/auth/access/usergroups/*', 'GET', '', 1, '2021-01-11 15:16:56'),
	(74, '根据用户名获取代理用户列表', '/cc/v1/proxyUsers/*', 'GET', '', 1, '2021-01-11 15:16:56'),
	(75, '根据用户名获取代理用户', '/cc/v1/proxyUser/*/user/*', 'GET', '', 1, '2021-01-11 15:16:56'),
	(76, '获取代理用户', '/cc/v1/proxyUser/*', 'GET', '', 1, '2021-01-11 15:16:56'),
	(77, '删除代理用户', '/cc/v1/proxyUser/*', 'DELETE', '', 1, '2021-01-11 15:16:56'),
	(78, '修改代理用户', '/cc/v1/proxyUser/*', 'PUT', '', 1, '2021-01-11 15:16:56'),
	(79, '增加代理用户', '/cc/v1/proxyUser', 'POST', '', 1, '2021-01-11 15:16:56'),
	(80, '创建模型', '/di/v1/models', 'POST', NULL, 1, '2021-03-11 16:30:17'),
	(81, '获取DI模型列表', '/di/v1/models/*', 'GET', NULL, 1, '2021-03-11 16:30:17'),
	(82, '根据ID获取Model', '/di/v1/models/*', 'GET', NULL, 1, '2021-03-11 16:30:17'),
	(83, '根据ID删除Model', '/di/v1/models/*', 'DELETE', NULL, 1, '2021-03-11 16:30:17'),
	(84, '根据ID导出Model', '/di/v1/models/*/export', 'GET', NULL, 1, '2021-03-11 16:30:17'),
	(85, '根据ID导出Model Definition', '/di/v1/models/*/definition', 'GET', NULL, 1, '2021-03-11 16:30:17'),
	(86, '根据ID导出TrainedModel', '/di/v1/models/*/trained_model', 'GET', NULL, 1, '2021-03-11 16:30:17'),
	(87, '根据ID获取模型日志', '/di/v1/models/*/logs', 'GET', NULL, 1, '2021-03-11 16:30:17'),
	(88, '根据ID与Line获取模型日志', '/di/v1/logs/*/loglines', 'GET', NULL, 1, '2021-03-11 16:30:17'),
	(89, '获取DI Dashboards', '/di/v1/dashboards?*', 'GET', NULL, 1, '2021-03-11 16:30:17'),
	(90, '上传模型文件', '/di/v1/codeUpload', 'POST', NULL, 1, '2021-03-11 16:30:17'),
	(91, '获取实验', '/di/v1/experiment/*', 'GET', NULL, 1, '2021-05-13 15:38:36'),
	(92, '创建实验', '/di/v1/experiment', 'POST', NULL, 1, '2021-05-13 15:39:14'),
	(93, '删除实验', '/di/v1/experiment/*', 'DELETE', NULL, 1, '2021-05-13 15:40:26'),
	(94, '更新实验', '/di/v1/experiment', 'PUT', NULL, 1, '2021-05-13 15:41:23'),
	(95, '更新实验', '/di/v1/experiment/*', 'PUT', NULL, 1, '2021-05-13 15:41:56'),
	(96, '创建实验标签', '/di/v1/experimentTag', 'POST', NULL, 1, '2021-05-13 15:42:27'),
	(97, '获取实验标签', '/di/v1/experimentTag/*', 'GET', NULL, 1, '2021-05-13 15:42:52'),
	(98, '代码上传', '/di/v1/codeUpload', 'POST', NULL, 1, '2021-05-13 15:43:39'),
	(99, '获取实验列表', '/di/v1/experiments', 'GET', NULL, 1, '2021-05-13 15:43:51'),
	(100, '导入实验', '/di/v1/experiment/import', 'POST', NULL, 1, '2021-05-13 15:44:17'),
	(101, '导出实验', '/di/v1/experiment/{id}/export', 'GET', NULL, 1, '2021-05-13 15:44:41'),
	(102, 'DSS导出', '/di/v1/experiment/{id}/exportdss', 'GET', NULL, 1, '2021-05-13 15:45:31'),
	(103, 'DSS导入', '/di/v1/experiment/importdss', 'POST', NULL, 1, '2021-05-13 15:45:48'),
	(104, '实验执行', '/di/v1/experimentRun', 'POST', NULL, 1, '2021-05-13 15:53:16'),
	(105, '实验执行获取', '/di/v1/experimentRun/*', 'GET', NULL, 1, '2021-05-13 15:56:07'),
	(106, '实验执行删除', '/di/v1/experimentRun/*', 'DELETE', NULL, 1, '2021-05-13 15:56:44'),
	(107, '实验执行日志获取', '/di/v1/experimentRun/*/log', 'GET', NULL, 1, '2021-05-13 15:59:16'),
	(108, '实验执行停止', '/di/v1/experimentRun/{id}/kill', 'GET', NULL, 1, '2021-05-13 15:59:53'),
	(109, '任务日志获取', '/di/v1/job/*/log', 'GET', NULL, 1, '2021-05-13 16:02:24'),
	(110, '实验执行删除', '/di/v1/experimentRun/*/kill', 'GET', NULL, 1, '2021-05-13 16:05:01'),
	(111, '实验执行状态后去', '/di/v1/experimentRun/*/status', 'GET', NULL, 1, '2021-05-13 16:05:20'),
	(112, '获取实验执行列表', '/di/v1/experimentRuns', 'GET', NULL, 1, '2021-05-13 16:05:48'),
	(113, '获取实验执行历史列表', '/di/v1/experimentRunsHistory/*', 'GET', NULL, 1, '2021-05-13 16:06:28'),
	(114, '获取实验执行任务列表', '/di/v1/experimentRun/*/execution', 'GET', NULL, 1, '2021-05-13 16:08:52'),
	(115, '获取实验执行任务日志', '/di/v1/experimentJob/*/log', 'GET', NULL, 1, '2021-05-13 16:09:27'),
	(116, '获取实验任务执行状态', '/di/v1/experimentJob/*/status', 'GET', NULL, 1, '2021-05-13 16:09:43'),
	(117, '获取DSS用户信息', '/di/v1/dssUserInfo?*', 'GET', NULL, 1, '2021-05-13 16:10:15'),
	(118, '获取模型列表', '/di/v1/models?*', 'GET', NULL, 1, '2021-05-13 16:19:55'),
	(119, '获取MF Dashboard', '/mf/v1/dashboard?*', 'GET', NULL, 1, '2021-11-17 02:45:41'),
	(120, 'MF API', '/mf/v1/*', 'GET', NULL, 1, '2021-11-17 02:57:38'),
	(121, 'MF POST API ', '/mf/v1/*', 'POST', NULL, 1, '2021-11-17 02:58:02'),
	(122, 'mlflow实验执行GET接口', '/mlflow/api/2.0/mlflow/runs/*', 'GET', NULL, 1, '2019-03-22 17:28:55'),
	(123, 'mlflow实验执行POST接口', '/mlflow/api/2.0/mlflow/runs/*', 'POST', NULL, 1, '2019-03-22 17:28:55'),
	(124, 'mlflow实验执行PUT接口', '/mlflow/ajax-api/2.0/mlflow/runs/*', 'PUT', NULL, 1, '2019-03-22 17:28:55'),
	(125, 'mlflow实验执行DELETE接口', '/mlflow/ajax-api/2.0/mlflow/runs/*', 'DELETE', NULL, 1, '2019-03-22 17:28:55'),
	(126, 'mlflow实验指标POST接口', '/mlflow/ajax-api/2.0/mlflow/metrics/*', 'POST', NULL, 1, '2019-03-22 17:28:55'),
	(127, 'mlflow实验指标GET接口', '/mlflow/ajax-api/2.0/mlflow/metrics/*', 'GET', NULL, 1, '2019-03-22 17:28:55'),
	(128, 'mlflow实验指标PUT接口', '/mlflow/ajax-api/2.0/mlflow/metrics/*', 'PUT', NULL, 1, '2019-03-22 17:28:55'),
	(129, 'mlflow实验指标DELETE接口', '/mlflow/ajax-api/2.0/mlflow/metrics/*', 'DELETE', NULL, 1, '2019-03-22 17:28:55'),
	(130, 'mlflow实验POST接口', '/mlflow/ajax-api/2.0/mlflow/experiment/*', 'POST', NULL, 1, '2019-03-22 17:28:55'),
	(131, 'mlflow实验GET接口', '/mlflow/ajax-api/2.0/mlflow/experiment/*', 'GET', NULL, 1, '2019-03-22 17:28:55'),
	(132, 'mlflow实验PUT接口', '/mlflow/ajax-api/2.0/mlflow/experiment/*', 'PUT', NULL, 1, '2019-03-22 17:28:55'),
	(133, 'mlflow实验DELETE接口', '/mlflow/ajax-api/2.0/mlflow/experiments/delete', 'POST', NULL, 1, '2019-03-22 17:28:55'),
	(134, 'mlflow模型POST接口', '/mlflow/ajax-api/2.0/preview/mlflow/registered-models/*', 'POST', NULL, 1, '2019-03-22 17:28:55'),
	(135, 'mlflow模型GET接口', '/mlflow/ajax-api/2.0/preview/mlflow/registered-models/*', 'GET', NULL, 1, '2019-03-22 17:28:55'),
	(136, 'mlflow模型PUT接口', '/mlflow/ajax-api/2.0/preview/mlflow/registered-models/*', 'PUT', NULL, 1, '2019-03-22 17:28:55'),
	(137, 'mlflow模型DELETE接口', '/mlflow/2.0/preview/mlflow/registered-models/*', 'DELETE', NULL, 1, '2019-03-22 17:28:55'),
	(138, 'mlflow模型POST接口', '/mlflow/2.0/preview/mlflow/model-versions/*', 'POST', NULL, 1, '2019-03-22 17:28:55'),
	(139, 'mlflow模型GET接口', '/mlflow/2.0/preview/mlflow/model-versions/*', 'GET', NULL, 1, '2019-03-22 17:28:55'),
	(140, 'mlflow模型PUT接口', '/mlflow/2.0/preview/mlflow/model-versions/*', 'PUT', NULL, 1, '2019-03-22 17:28:55'),
	(141, 'mlflow模型DELETE接口', '/mlflow/2.0/preview/mlflow/model-versions/*', 'DELETE', NULL, 1, '2019-03-22 17:28:55'),
	(142, 'mlflow模型版本下载POST接口', '/mlflow/2.0/preview/mlflow/model-versions/get-download-uri', 'POST', NULL, 1, '2019-03-22 17:28:55'),
	(143, 'mlflow模型版本下载GET接口', '/mlflow/2.0/preview/mlflow/model-versions/get-download-uri', 'GET', NULL, 1, '2019-03-22 17:28:55'),
	(144, 'mlflow模型版本transition POST接口', '/mlflow/2.0/preview/mlflow/model-versions/transition-stage', 'POST', NULL, 1, '2019-03-22 17:28:55'),
	(145, 'mlflow模型版本transition GET接口', '/mlflow/2.0/preview/mlflow/model-versions/transition-stage', 'GET', NULL, 1, '2019-03-22 17:28:55');
/*!40000 ALTER TABLE `t_permission` ENABLE KEYS */;

-- 导出  表 mlss_gzpc_bdap_uat_01.t_role 结构
CREATE TABLE IF NOT EXISTS `t_role` (
  `id` bigint(255) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `remarks` varchar(200) DEFAULT NULL,
  `enable_flag` tinyint(4) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name_UNIQUE` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=56 DEFAULT CHARSET=utf8;

-- 正在导出表  mlss_gzpc_bdap_uat_01.t_role 的数据：~3 rows (大约)
DELETE FROM `t_role`;
/*!40000 ALTER TABLE `t_role` DISABLE KEYS */;
INSERT INTO `t_role` (`id`, `name`, `remarks`, `enable_flag`) VALUES
	(1, 'ADMIN', 'Control all resources in a group', 1),
	(2, 'USER', 'Can access namespaced resources, but cannot delete or create', 1),
	(3, 'VIEWER', 'Can only view resources in a group', 1);
/*!40000 ALTER TABLE `t_role` ENABLE KEYS */;

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

-- 正在导出表  mlss_gzpc_bdap_uat_01.t_role_permission 的数据：~198 rows (大约)
DELETE FROM `t_role_permission`;
/*!40000 ALTER TABLE `t_role_permission` DISABLE KEYS */;
INSERT INTO `t_role_permission` (`id`, `role_id`, `permission_id`, `enable_flag`, `operate_time`) VALUES
	(1, 1, 1, 1, '2019-03-18 14:16:05'),
	(2, 1, 2, 1, '2019-03-18 14:16:05'),
	(3, 1, 3, 1, '2019-03-18 14:16:05'),
	(4, 1, 4, 1, '2019-03-18 14:16:05'),
	(5, 1, 5, 1, '2019-03-18 14:16:05'),
	(6, 1, 6, 1, '2019-03-18 14:16:05'),
	(7, 1, 7, 1, '2019-03-18 14:16:05'),
	(8, 1, 8, 1, '2019-03-18 14:16:05'),
	(9, 1, 9, 1, '2019-03-22 17:26:22'),
	(18, 1, 10, 1, '2019-04-16 17:22:25'),
	(19, 1, 11, 1, '2019-04-16 17:22:25'),
	(20, 1, 12, 1, '2019-04-16 17:22:25'),
	(21, 2, 10, 1, '2019-04-16 17:22:25'),
	(22, 2, 11, 1, '2019-04-16 17:22:25'),
	(23, 2, 12, 0, '2019-04-16 17:22:25'),
	(31, 2, 3, 1, '2019-04-17 15:10:11'),
	(32, 2, 4, 1, '2019-04-17 15:10:11'),
	(35, 1, 20, 1, '2019-04-17 15:21:02'),
	(36, 1, 21, 1, '2019-04-17 15:21:02'),
	(37, 2, 8, 1, '2019-04-18 09:57:05'),
	(38, 2, 9, 1, '2019-04-18 10:36:16'),
	(39, 1, 13, 1, '2019-04-25 19:01:50'),
	(40, 2, 13, 0, '2019-04-25 19:01:57'),
	(41, 2, 18, 1, '2019-04-25 19:39:33'),
	(42, 2, 19, 1, '2019-04-25 19:39:33'),
	(43, 1, 18, 1, '2019-04-25 19:39:33'),
	(44, 1, 19, 1, '2019-04-25 19:39:33'),
	(45, 1, 22, 1, '2019-05-06 15:18:33'),
	(46, 2, 22, 1, '2019-05-06 15:18:55'),
	(48, 1, 23, 1, '2019-05-16 18:23:38'),
	(49, 2, 23, 1, '2019-05-16 18:23:38'),
	(50, 1, 24, 1, '2019-05-20 20:38:50'),
	(51, 2, 24, 1, '2019-05-20 20:38:50'),
	(52, 1, 25, 1, '2019-06-15 20:03:49'),
	(53, 2, 25, 1, '2019-06-15 20:03:57'),
	(54, 1, 26, 1, '2019-06-15 20:30:15'),
	(55, 2, 26, 1, '2019-06-15 20:30:15'),
	(56, 1, 27, 1, '2019-06-15 20:30:15'),
	(57, 2, 27, 1, '2019-06-15 20:30:15'),
	(58, 1, 28, 1, '2019-06-15 20:30:15'),
	(59, 2, 28, 1, '2019-06-15 20:30:15'),
	(60, 1, 29, 1, '2019-06-17 19:19:45'),
	(61, 2, 29, 1, '2019-06-17 19:19:45'),
	(64, 1, 30, 1, '2019-07-12 10:57:18'),
	(65, 2, 30, 1, '2019-07-12 10:57:18'),
	(66, 1, 31, 1, '2019-07-12 10:57:18'),
	(67, 2, 31, 1, '2019-07-12 10:57:18'),
	(68, 1, 32, 1, '2019-07-12 10:57:18'),
	(69, 2, 32, 1, '2019-07-12 10:57:18'),
	(70, 1, 33, 1, '2019-07-17 16:42:29'),
	(71, 2, 33, 1, '2019-07-17 16:42:29'),
	(72, 1, 34, 1, '2019-08-27 16:34:30'),
	(73, 2, 34, 1, '2019-08-27 16:34:30'),
	(74, 1, 35, 1, '2019-08-28 16:54:10'),
	(75, 2, 35, 1, '2019-08-28 16:54:10'),
	(77, 1, 36, 1, '2019-10-08 10:23:11'),
	(78, 2, 36, 1, '2019-10-08 10:23:11'),
	(166, 2, 73, 1, '2021-01-11 15:22:20'),
	(167, 2, 74, 1, '2021-01-11 15:22:20'),
	(168, 2, 75, 1, '2021-01-11 15:22:20'),
	(169, 1, 76, 1, '2021-01-11 15:22:20'),
	(170, 1, 73, 1, '2021-01-11 15:22:20'),
	(171, 1, 74, 1, '2021-01-11 15:22:20'),
	(172, 1, 75, 1, '2021-01-11 15:22:20'),
	(174, 1, 77, 1, '2021-01-11 15:22:50'),
	(175, 1, 78, 1, '2021-01-11 15:22:51'),
	(176, 1, 79, 1, '2021-01-11 15:22:51'),
	(177, 2, 76, 1, '2021-01-11 15:23:14'),
	(178, 1, 80, 1, '2021-03-11 16:31:31'),
	(179, 1, 81, 1, '2021-03-11 16:31:31'),
	(180, 1, 82, 1, '2021-03-11 16:31:31'),
	(181, 1, 83, 1, '2021-03-11 16:31:31'),
	(182, 1, 84, 1, '2021-03-11 16:31:31'),
	(183, 1, 85, 1, '2021-03-11 16:31:31'),
	(184, 1, 86, 1, '2021-03-11 16:31:31'),
	(185, 1, 87, 1, '2021-03-11 16:31:31'),
	(186, 1, 88, 1, '2021-03-11 16:31:32'),
	(187, 1, 89, 1, '2021-03-11 16:31:32'),
	(188, 1, 90, 1, '2021-03-11 16:31:32'),
	(189, 2, 80, 1, '2021-03-11 16:31:32'),
	(190, 2, 81, 1, '2021-03-11 16:31:32'),
	(191, 2, 82, 1, '2021-03-11 16:31:32'),
	(192, 2, 83, 1, '2021-03-11 16:31:32'),
	(193, 2, 84, 1, '2021-03-11 16:31:32'),
	(194, 2, 85, 1, '2021-03-11 16:31:32'),
	(195, 2, 86, 1, '2021-03-11 16:31:32'),
	(196, 2, 87, 1, '2021-03-11 16:31:32'),
	(197, 2, 88, 1, '2021-03-11 16:31:32'),
	(198, 2, 89, 1, '2021-03-11 16:31:32'),
	(199, 2, 90, 1, '2021-03-11 16:31:32'),
	(200, 2, 91, 1, '2021-05-13 16:13:26'),
	(201, 1, 91, 1, '2021-05-13 16:18:25'),
	(202, 1, 92, 1, '2021-05-13 16:18:26'),
	(203, 2, 92, 1, '2021-05-13 16:18:26'),
	(204, 1, 93, 1, '2021-05-13 16:18:26'),
	(205, 2, 93, 1, '2021-05-13 16:18:26'),
	(206, 1, 94, 1, '2021-05-13 16:18:26'),
	(207, 2, 94, 1, '2021-05-13 16:18:26'),
	(208, 1, 95, 1, '2021-05-13 16:18:26'),
	(209, 2, 95, 1, '2021-05-13 16:18:26'),
	(210, 1, 96, 1, '2021-05-13 16:18:26'),
	(211, 2, 96, 1, '2021-05-13 16:18:26'),
	(212, 1, 97, 1, '2021-05-13 16:18:26'),
	(213, 2, 97, 1, '2021-05-13 16:18:26'),
	(214, 1, 98, 1, '2021-05-13 16:18:26'),
	(215, 2, 98, 1, '2021-05-13 16:18:26'),
	(216, 1, 99, 1, '2021-05-13 16:18:26'),
	(217, 2, 99, 1, '2021-05-13 16:18:27'),
	(218, 1, 100, 1, '2021-05-13 16:18:27'),
	(219, 2, 100, 1, '2021-05-13 16:18:27'),
	(220, 1, 101, 1, '2021-05-13 16:18:27'),
	(221, 2, 101, 1, '2021-05-13 16:18:27'),
	(222, 1, 102, 1, '2021-05-13 16:18:27'),
	(223, 2, 102, 1, '2021-05-13 16:18:27'),
	(224, 1, 103, 1, '2021-05-13 16:18:27'),
	(225, 2, 103, 1, '2021-05-13 16:18:27'),
	(226, 1, 104, 1, '2021-05-13 16:18:27'),
	(227, 2, 104, 1, '2021-05-13 16:18:27'),
	(228, 1, 105, 1, '2021-05-13 16:18:27'),
	(229, 2, 105, 1, '2021-05-13 16:18:27'),
	(230, 1, 106, 1, '2021-05-13 16:18:27'),
	(231, 2, 106, 1, '2021-05-13 16:18:27'),
	(232, 1, 107, 1, '2021-05-13 16:18:27'),
	(233, 2, 107, 1, '2021-05-13 16:18:27'),
	(234, 1, 108, 1, '2021-05-13 16:18:27'),
	(235, 2, 108, 1, '2021-05-13 16:18:28'),
	(236, 1, 109, 1, '2021-05-13 16:18:28'),
	(237, 2, 109, 1, '2021-05-13 16:18:28'),
	(238, 1, 110, 1, '2021-05-13 16:18:28'),
	(239, 2, 110, 1, '2021-05-13 16:18:28'),
	(240, 1, 111, 1, '2021-05-13 16:18:28'),
	(241, 2, 111, 1, '2021-05-13 16:18:28'),
	(242, 1, 112, 1, '2021-05-13 16:18:28'),
	(243, 2, 112, 1, '2021-05-13 16:18:28'),
	(244, 1, 113, 1, '2021-05-13 16:18:28'),
	(245, 2, 113, 1, '2021-05-13 16:18:28'),
	(246, 1, 114, 1, '2021-05-13 16:18:28'),
	(247, 2, 114, 1, '2021-05-13 16:18:28'),
	(248, 1, 115, 1, '2021-05-13 16:18:28'),
	(249, 2, 115, 1, '2021-05-13 16:18:28'),
	(250, 1, 116, 1, '2021-05-13 16:18:28'),
	(251, 2, 116, 1, '2021-05-13 16:18:28'),
	(252, 1, 117, 1, '2021-05-13 16:18:28'),
	(253, 2, 117, 1, '2021-05-13 16:18:29'),
	(254, 1, 118, 1, '2021-05-13 16:20:21'),
	(255, 2, 118, 1, '2021-05-13 16:20:21'),
	(256, 1, 119, 1, '2021-11-17 02:46:12'),
	(257, 2, 119, 1, '2021-11-17 02:46:23'),
	(258, 2, 120, 1, '2021-12-24 18:34:21'),
	(259, 2, 121, 1, '2021-12-24 18:34:44'),
	(260, 1, 122, 1, '2019-03-18 14:16:05'),
	(261, 2, 122, 1, '2019-03-18 14:16:05'),
	(262, 1, 123, 1, '2019-03-18 14:16:05'),
	(263, 2, 123, 1, '2019-03-18 14:16:05'),
	(264, 1, 124, 1, '2019-03-18 14:16:05'),
	(265, 2, 124, 1, '2019-03-18 14:16:05'),
	(266, 1, 125, 1, '2019-03-18 14:16:05'),
	(267, 2, 125, 1, '2019-03-18 14:16:05'),
	(268, 1, 126, 1, '2019-03-18 14:16:05'),
	(269, 2, 126, 1, '2019-03-18 14:16:05'),
	(270, 1, 127, 1, '2019-03-18 14:16:05'),
	(271, 2, 127, 1, '2019-03-18 14:16:05'),
	(272, 1, 128, 1, '2019-03-18 14:16:05'),
	(273, 2, 128, 1, '2019-03-18 14:16:05'),
	(274, 1, 129, 1, '2019-03-18 14:16:05'),
	(275, 2, 129, 1, '2019-03-18 14:16:05'),
	(276, 1, 130, 1, '2019-03-18 14:16:05'),
	(277, 2, 130, 0, '2019-03-18 14:16:05'),
	(278, 1, 131, 1, '2019-03-18 14:16:05'),
	(279, 2, 131, 1, '2019-03-18 14:16:05'),
	(280, 1, 132, 1, '2019-03-18 14:16:05'),
	(281, 2, 132, 1, '2019-03-18 14:16:05'),
	(282, 1, 133, 1, '2019-03-18 14:16:05'),
	(283, 2, 133, 1, '2019-03-18 14:16:05'),
	(284, 1, 134, 1, '2019-03-18 14:16:05'),
	(285, 2, 134, 1, '2019-03-18 14:16:05'),
	(286, 1, 135, 1, '2019-03-18 14:16:05'),
	(287, 2, 135, 1, '2019-03-18 14:16:05'),
	(288, 1, 136, 1, '2019-03-18 14:16:05'),
	(289, 2, 136, 1, '2019-03-18 14:16:05'),
	(290, 1, 137, 1, '2019-03-18 14:16:05'),
	(291, 2, 137, 1, '2019-03-18 14:16:05'),
	(292, 1, 138, 1, '2019-03-18 14:16:05'),
	(293, 2, 138, 1, '2019-03-18 14:16:05'),
	(294, 1, 139, 1, '2019-03-18 14:16:05'),
	(295, 2, 139, 1, '2019-03-18 14:16:05'),
	(296, 1, 140, 1, '2019-03-18 14:16:05'),
	(297, 2, 140, 1, '2019-03-18 14:16:05'),
	(298, 1, 141, 1, '2019-03-18 14:16:05'),
	(299, 2, 141, 1, '2019-03-18 14:16:05'),
	(300, 1, 142, 1, '2019-03-18 14:16:05'),
	(301, 2, 142, 1, '2019-03-18 14:16:05'),
	(302, 1, 143, 1, '2019-03-18 14:16:05'),
	(303, 2, 143, 1, '2019-03-18 14:16:05'),
	(304, 1, 144, 1, '2019-03-18 14:16:05'),
	(305, 2, 144, 1, '2019-03-18 14:16:05'),
	(306, 1, 145, 1, '2019-03-18 14:16:05'),
	(307, 2, 145, 1, '2019-03-18 14:16:05');
/*!40000 ALTER TABLE `t_role_permission` ENABLE KEYS */;

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

-- 正在导出表  mlss_gzpc_bdap_uat_01.t_storage 的数据：~2 rows (大约)
DELETE FROM `t_storage`;
/*!40000 ALTER TABLE `t_storage` DISABLE KEYS */;
INSERT INTO `t_storage` (`id`, `path`, `remarks`, `enable_flag`, `type`) VALUES
	(1, '/data/bdap-ss/mlss-data/hadoop', NULL, 1, 'bdap-ss'),
	(363, '/data/bdap-ss/mlss-data/alexwu', '', 1, 'bdap-ss');
/*!40000 ALTER TABLE `t_storage` ENABLE KEYS */;

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

-- 正在导出表  mlss_gzpc_bdap_uat_01.t_superadmin 的数据：~12 rows (大约)
DELETE FROM `t_superadmin`;
/*!40000 ALTER TABLE `t_superadmin` DISABLE KEYS */;
INSERT INTO `t_superadmin` (`id`, `name`, `remarks`, `enable_flag`, `operate_time`) VALUES
	(1, 'kirkzhou', NULL, 1, '2019-03-12 19:09:40'),
	(2, 'luckliu', NULL, 0, '2019-03-12 19:09:40'),
	(5, 'neiljianliu', NULL, 0, '2019-03-22 15:11:29'),
	(6, 'hduser05', NULL, 0, '2019-04-16 20:55:59'),
	(7, 'hduser2006', NULL, 0, '2019-06-18 17:46:30'),
	(8, 'hduser2007', NULL, 0, '2019-06-18 17:46:40'),
	(9, 'hduser2008', NULL, 0, '2019-06-18 17:46:48'),
	(10, 'hduser1009', NULL, 0, '2019-06-18 17:46:55'),
	(11, 'hduser1010', NULL, 0, '2019-06-18 17:47:01'),
	(12, 'alexwu', NULL, 1, '2019-11-11 19:49:50'),
	(13, 'nanzou', NULL, 1, '2021-11-03 15:25:24'),
	(14, 'livyhuang', NULL, 1, '2021-12-02 20:27:56');
/*!40000 ALTER TABLE `t_superadmin` ENABLE KEYS */;

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

-- 正在导出表  mlss_gzpc_bdap_uat_01.t_user 的数据：~3 rows (大约)
DELETE FROM `t_user`;
/*!40000 ALTER TABLE `t_user` DISABLE KEYS */;
INSERT INTO `t_user` (`id`, `name`, `gid`, `uid`, `token`, `type`, `remarks`, `enable_flag`, `guid_check`) VALUES
	(1, 'alexwu', 22324, 22324, NULL, 'SYSTEM', NULL, 1, 0),
	(539, 'neiljianliu', 30, 30, '', 'USER', '30', 1, 0),
	(540, 'hadoop', 6004, 6004, '', 'USER', '6004', 1, 0);
/*!40000 ALTER TABLE `t_user` ENABLE KEYS */;

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

-- 正在导出表  mlss_gzpc_bdap_uat_01.t_user_group 的数据：~3 rows (大约)
DELETE FROM `t_user_group`;
/*!40000 ALTER TABLE `t_user_group` DISABLE KEYS */;
INSERT INTO `t_user_group` (`id`, `user_id`, `role_id`, `group_id`, `remarks`, `enable_flag`) VALUES
	(669, 539, 1, 542, '个人用户组', 1),
	(670, 1, 2, 543, '', 1),
	(671, 540, 1, 544, '个人用户组', 1);
/*!40000 ALTER TABLE `t_user_group` ENABLE KEYS */;

/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IF(@OLD_FOREIGN_KEY_CHECKS IS NULL, 1, @OLD_FOREIGN_KEY_CHECKS) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
