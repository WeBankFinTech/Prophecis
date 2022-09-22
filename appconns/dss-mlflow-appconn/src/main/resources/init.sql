
delete from `dss_application` where `name`='mlflow';
INSERT  INTO `dss_application`(
    `name`,
    `url`,
    `is_user_need_init`,
    `level`,
    `user_init_url`,
    `exists_project_service`,
    `project_url`,
    `enhance_json`,
    `if_iframe`,
    `homepage_url`,
    `redirect_url`)
VALUES (
'mlflow', 
'http://127.0.0.1:30803', 
0,
1, 
NULL, 
1,
 'http://127.0.0.1:30803', 
'{"MLSS-SecretKey":"MLFLOW","MLSS-Auth-Type":"SYSTEM","MLSS-APPSignature":"MLFLOW","MLSS-BaseUrl":"http://127.0.0.1:30803","baseUrl":"http://127.0.0.1:30803","MLSS-ModelMonitoring-JAR":"/appcom/Install/quickml/qml_algo/hwenzan/qml_algo.jar"}',
 1, 
'http://127.0.0.1:30803/#/dashboard', 
'http://127.0.0.1:30803/#/mlFlow');

select @dss_mlss_applicationId:=id from `dss_application` WHERE `name` in('mlflow');

select @dss_appconn_mlssId:=id from `dss_appconn` where `appconn_name` = 'mlflow';
delete from `dss_appconn_instance` where  `appconn_id`=@dss_appconn_mlssId;

delete from `dss_appconn`  where `appconn_name`='mlflow';
INSERT INTO `dss_appconn` (
            `appconn_name`,
			`is_user_need_init`,
			`level`,
			`if_iframe`,
			`is_external`,
			`reference`,
			`class_name`,
			`appconn_class_path`,
			`resource`
			)
			VALUES(
            'mlflow',
			0,
			1,
			1,
			1,
			NULL,
			'com.webank.wedatasphere.dss.appconn.mlflow.MLFlowAppConn',
			'/appcom/Install/dss/dss-appconns/mlflow',
			''
			);


select @dss_appconn_mlssId:=id from `dss_appconn` where `appconn_name` = 'mlflow';




INSERT INTO `dss_appconn_instance`(
            `appconn_id`,
			`label`,
			`url`,
			`enhance_json`, 
			`homepage_uri`)
            VALUES (
            @dss_appconn_mlssId,
			'DEV',
			'http://127.0.0.1:30803',
			'{"MLSS-SecretKey":"MLFLOW","MLSS-Auth-Type":"SYSTEM","MLSS-APPSignature":"MLFLOW","MLSS-BaseUrl":"http://127.0.0.1:30803","baseUrl":"http://127.0.0.1:30803","MLSS-ModelMonitoring-JAR":"/appcom/Install/quickml/qml_algo/hwenzan/qml_algo.jar"}',
			'dss/mlflow/#/mlFlow'
			);



select @dss_mlssId:=name from `dss_workflow_node` where `node_type` = 'linkis.appconn.mlflow';
delete from `dss_workflow_node_to_group` where `node_id`=@dss_mlssId;
--
delete from `dss_workflow_node` where `node_type`='linkis.appconn.mlflow';
INSERT INTO `dss_workflow_node` (
    `name`,
	`appconn_name`,
	`node_type`,
	`jump_type`,
	`support_jump`,
	`submit_to_scheduler`,
	`enable_copy`,
	`should_creation_before_node`,
	`icon_path`
	)
VALUES (
    'mlflow',
	'mlflow',
	'linkis.appconn.mlflow',
	1,
	1,
	1,
	0,
	1,
	'icons/mlflow.icon');

--
select @dss_mlssId:=id from `dss_workflow_node` where `node_type` = 'linkis.appconn.mlflow';
insert  into `dss_workflow_node_to_group` (`node_id`, `group_id`) values (@dss_mlssId, 8);


select @dss_workflow_node_id:=id from `dss_workflow_node` where `node_type` = 'linkis.appconn.mlflow';
INSERT INTO `dss_workflow_node_to_ui` (`workflow_node_id`, `ui_id`) VALUES
	(@dss_mlssId, 1),
	(@dss_mlssId, 4),
	(@dss_mlssId, 5),
	(@dss_mlssId, 6),
	(@dss_mlssId, 35),
	(@dss_mlssId, 36),
	(@dss_mlssId, 3);

