/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;

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

-- 正在导出表  mlss_gzpc_bdap_uat_01.t_keypair 的数据：~0 rows (大约)
DELETE FROM `t_keypair`;
/*!40000 ALTER TABLE `t_keypair` DISABLE KEYS */;
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
) ENGINE=InnoDB AUTO_INCREMENT=41 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT COMMENT='权限表';

-- 正在导出表  mlss_gzpc_bdap_uat_01.t_permission 的数据：~40 rows (大约)
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
	(13, '获取用户列表', '/cc/v1/users', 'GET', NULL, 1, '2019-04-17 14:39:58'),
	(14, '新增用户', '/cc/v1/users', 'POST', NULL, 1, '2019-04-17 14:39:58'),
	(15, '修改用户', '/cc/v1/users', 'PUT', NULL, 1, '2019-04-17 14:39:58'),
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
	(37, '管理员获取用户组', '/cc/v1/groups', 'GET', NULL, 1, '2020-12-03 19:29:04'),
	(38, '管理员新增项目组', '/cc/v1/groups', 'POST', NULL, 1, '2020-12-03 19:29:37'),
	(39, '获取可管理的用户列表', '/cc/v1/users/myUsers', 'GET', NULL, 1, '2020-12-03 19:34:29'),
	(40, '获取用户存储', '/cc/v1/groups/group/storage', 'GET', NULL, 1, '2020-12-03 20:51:55');
/*!40000 ALTER TABLE `t_permission` ENABLE KEYS */;

-- 导出  表 mlss_gzpc_bdap_uat_01.t_role 结构
CREATE TABLE IF NOT EXISTS `t_role` (
  `id` bigint(255) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `remarks` varchar(200) DEFAULT NULL,
  `enable_flag` tinyint(4) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `name_UNIQUE` (`name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;

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
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `role_permission_UNIQUE` (`role_id`,`permission_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=94 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT COMMENT='角色权限关联表';

-- 正在导出表  mlss_gzpc_bdap_uat_01.t_role_permission 的数据：~73 rows (大约)
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
	(23, 2, 12, 1, '2019-04-16 17:22:25'),
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
	(72, 1, 34, 1, '2019-08-27 16:39:00'),
	(73, 2, 34, 1, '2019-08-27 16:39:00'),
	(74, 1, 35, 1, '2019-08-28 17:02:36'),
	(75, 2, 35, 1, '2019-08-28 17:02:36'),
	(76, 1, 36, 1, '2019-10-12 11:23:35'),
	(77, 2, 36, 1, '2019-10-12 11:23:35'),
	(78, 1, 37, 1, '2019-11-11 18:15:35'),
	(79, 2, 37, 0, '2019-11-11 18:15:35'),
	(80, 1, 38, 1, '2019-11-11 18:15:35'),
	(81, 2, 38, 0, '2019-11-11 18:15:36'),
	(82, 1, 39, 1, '2019-11-11 18:15:36'),
	(83, 2, 39, 1, '2019-11-11 18:15:36'),
	(86, 1, 40, 1, '2019-11-11 18:15:37'),
	(87, 2, 40, 1, '2019-11-11 18:15:37'),
	(108, 1, 41, 1, '2021-03-01 16:36:21'),
	(109, 1, 42, 1, '2021-03-01 16:36:21'),
	(110, 1, 43, 1, '2021-03-01 16:36:21'),
	(111, 1, 44, 1, '2021-03-01 16:36:21'),
	(112, 1, 45, 1, '2021-03-01 16:36:21'),
	(113, 1, 46, 1, '2021-03-01 16:36:21'),
	(114, 1, 47, 1, '2021-03-01 16:36:21'),
	(115, 1, 48, 1, '2021-03-01 16:36:21'),
	(116, 1, 49, 1, '2021-03-01 16:36:21'),
	(117, 1, 50, 1, '2021-03-01 16:36:21'),
	(118, 2, 41, 1, '2021-03-01 16:36:21'),
	(119, 2, 42, 1, '2021-03-01 16:36:22'),
	(120, 2, 43, 1, '2021-03-01 16:36:22'),
	(121, 2, 44, 1, '2021-03-01 16:36:22'),
	(122, 2, 45, 1, '2021-03-01 16:36:22'),
	(123, 2, 46, 1, '2021-03-01 16:36:22'),
	(124, 2, 47, 1, '2021-03-01 16:36:22'),
	(125, 2, 48, 1, '2021-03-01 16:36:22'),
	(126, 2, 49, 1, '2021-03-01 16:36:22'),
	(127, 2, 50, 1, '2021-03-01 16:36:22'),
	(128, 2, 51, 1, '2021-03-01 16:58:16'),
	(129, 1, 51, 1, '2021-03-01 16:58:20');

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

-- 正在导出表  mlss_gzpc_bdap_uat_01.t_superadmin 的数据：~2 rows (大约)
DELETE FROM `t_superadmin`;
/*!40000 ALTER TABLE `t_superadmin` DISABLE KEYS */;
INSERT INTO `t_superadmin` (`id`, `name`, `remarks`, `enable_flag`, `operate_time`) VALUES
	(1, 'alexwu', NULL, 1, '2020-11-25 21:32:11'),
	(2, 'admin', NULL, 1, '2020-11-25 21:32:19');
/*!40000 ALTER TABLE `t_superadmin` ENABLE KEYS */;

/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IF(@OLD_FOREIGN_KEY_CHECKS IS NULL, 1, @OLD_FOREIGN_KEY_CHECKS) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
