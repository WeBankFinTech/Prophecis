
delete from `dss_application` where `name`='mlss';
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
    'mlss',
    'http://127.0.0.1:30803',
    0,
    1,
    NULL,
    1,
    'http://127.0.0.1:30803',
    '{"MLSS-SecretKey":"MLFLOW","MLSS-Auth-Type":"SYSTEM","MLSS-APPSignature":"MLFLOW","MLSS-BaseUrl":"http://127.0.0.1:30803","baseUrl":"http://127.0.0.1:30803"}',
    1,
    'http://127.0.0.1:30803/#/dashboard',
    'http://127.0.0.1:30803/#/mlFlow');

select @dss_mlss_applicationId:=id from `dss_application` WHERE `name` in('mlss');

select @dss_appconn_mlssId:=id from `dss_appconn` where `appconn_name` = 'mlss';
delete from `dss_appconn_instance` where  `appconn_id`=@dss_appconn_mlssId;

delete from `dss_appconn`  where `appconn_name`='mlss';
INSERT INTO `dss_appconn` (
            `appconn_name`,
            `is_user_need_init`,
            `level`,
            `if_iframe`,
            `is_external`,
            `reference`,
            `class_name`,
            `appconn_class_path`,
            `resource`)
            VALUES (
                'mlss',
                0,
                1,
                NULL,
                0,
                NULL,
                'com.webank.wedatasphere.dss.appconn.mlss.MLSSAppConn',
                '/appcom/Install/DSSInstall/dss-1.1.2/dss-appconns/mlss',
                '');


select @dss_appconn_mlssId:=id from `dss_appconn` where `appconn_name` = 'mlss';




INSERT INTO `dss_appconn_instance`(
            `appconn_id`,
            `label`,
            `url`,
            `enhance_json`,
            `homepage_url`,
            `redirect_url`)
            VALUES (
            @dss_appconn_mlssId,
            'DEV',
            'http://127.0.0.1:30803/',
            '{"MLSS-SecretKey":"MLFLOW","MLSS-Auth-Type":"SYSTEM","MLSS-APPSignature":"MLFLOW","MLSS-BaseUrl":"http://127.0.0.1:30803","baseUrl":"http://127.0.0.1:30803"}',
            'http://127.0.0.1:30803/#/dashboard',
            'http://127.0.0.1:30803/#/mlFlow');


INSERT INTO `dss_appconn_instance`(
            `appconn_id`,
            `label`,
            `url`,
            `enhance_json`,
            `homepage_url`,
            `redirect_url`)
            VALUES (
            @dss_appconn_mlssId,
            'PROD',
            'http://127.0.0.1:30803/',
            '{"MLSS-SecretKey":"MLFLOW","MLSS-Auth-Type":"SYSTEM","MLSS-APPSignature":"MLFLOW","MLSS-BaseUrl":"http://127.0.0.1:30803","baseUrl":"http://127.0.0.1:30803"}',
            'http://127.0.0.1:30803/#/dashboard',
            'http://127.0.0.1:30803/#/mlFlow');

select @dss_mlssId:=name from `dss_workflow_node` where `node_type` = 'linkis.appconn.mlss';
delete from `dss_workflow_node_to_group` where `node_id`=@dss_mlssId;
--
delete from `dss_workflow_node` where `node_type`='linkis.appconn.mlss';
INSERT INTO `dss_workflow_node` (
            `icon`,
            `node_type`,
            `appconn_name`,
            `submit_to_scheduler`,
            `enable_copy`,
            `should_creation_before_node`,
            `support_jump`,
            `jump_url`,
            `name`)
            VALUES (
            '<?xml version="1.0" standalone="no"?><!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd"><svg t="1603264751694" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="5768" xmlns:xlink="http://www.w3.org/1999/xlink" width="200" height="200"><defs><style type="text/css"></style></defs><path d="M826.993778 748.373333l-2.901334 1.763556-113.436444-114.915556 2.161778-3.072c24.974222-35.043556 38.115556-76.572444 38.115555-120.092444a207.075556 207.075556 0 0 0-37.888-119.921778l-2.218666-3.072 113.550222-115.029333 2.844444 1.706666c19.911111 11.719111 42.666667 17.92 65.877334 17.976889C965.290667 293.717333 1024 234.211556 1024 160.995556 1024 87.950222 965.290667 28.444444 893.098667 28.444444c-72.248889 0-130.958222 59.505778-130.958223 132.608 0 23.381333 6.087111 46.364444 17.692445 66.787556l1.706667 2.958222-113.550223 115.086222-3.072-2.161777a201.272889 201.272889 0 0 0-250.538666 12.458666A206.734222 206.734222 0 0 0 346.396444 480.711111l-0.568888 3.811556H258.958222l-0.796444-3.470223A130.958222 130.958222 0 0 0 130.844444 379.847111C58.709333 379.847111 0 439.352889 0 512.568889c0 73.102222 58.709333 132.551111 130.901333 132.551111A131.811556 131.811556 0 0 0 256.853333 548.977778l0.853334-3.185778h88.974222l0.625778 3.640889c17.749333 97.792 101.888 168.732444 200.192 168.732444a200.248889 200.248889 0 0 0 117.134222-37.546666l2.958222-2.161778 113.777778 115.143111-1.706667 3.015111c-11.605333 20.252444-17.635556 43.064889-17.635555 66.275556 0 73.159111 58.766222 132.664889 130.958222 132.664889 72.192 0 130.958222-59.505778 130.958222-132.664889 0-73.159111-58.766222-132.608-130.844444-132.608-23.324444 0-46.250667 6.257778-66.104889 18.090666z m-380.586667-133.916444A144.725333 144.725333 0 0 1 404.536889 512c0-38.684444 14.848-75.093333 41.870222-102.4a141.084444 141.084444 0 0 1 101.091556-42.439111c38.172444 0 74.695111 15.36 101.091555 42.439111 27.079111 27.420444 42.097778 64.170667 41.927111 102.4 0 38.684444-14.904889 75.036444-41.927111 102.4a140.913778 140.913778 0 0 1-101.091555 42.439111c-38.172444 0-74.695111-15.36-101.091556-42.439111z m517.12-453.404445c0 39.367111-31.630222 71.395556-70.428444 71.395556-14.222222 0-28.046222-4.380444-39.253334-12.003556l-18.773333-18.887111a71.452444 71.452444 0 0 1-12.458667-40.504889c0-39.367111 31.573333-71.338667 70.485334-71.338666 38.798222 0 70.428444 31.971556 70.428444 71.395555z m-70.428444 773.233778a70.997333 70.997333 0 0 1-70.485334-71.395555c0-39.310222 31.573333-71.338667 70.485334-71.338667 38.798222 0 70.428444 32.028444 70.428444 71.338667 0 39.367111-31.630222 71.395556-70.428444 71.395555zM201.329778 512.398222c0 39.310222-31.516444 71.338667-70.428445 71.338667a70.940444 70.940444 0 0 1-70.428444-71.338667c0-39.367111 31.630222-71.395556 70.428444-71.395555 38.912 0 70.485333 32.028444 70.485334 71.395555z" fill="#8470FF" p-id="5769"></path></svg>',
            'linkis.appconn.mlss',
            'mlss',
            1,
            0,
            0,
            1,
            'http://127.0.0.1:30803/#/mlFlow?tableType=1&id=${projectId}&nodeId=${nodeId}&contextID=${contextID}&nodeName=${nodeName}&expId=${expId}&projectId=${projectId}',
            'mlss');

--
select @dss_mlssId:=id from `dss_workflow_node` where `node_type` = 'linkis.appconn.mlss';
insert  into `dss_workflow_node_to_group` (`node_id`, `group_id`) values (@dss_mlssId, 8);


select @dss_workflow_node_id:=id from `dss_workflow_node` where `node_type` = 'linkis.appconn.mlss';
INSERT INTO `dss_workflow_node_to_ui` (`workflow_node_id`, `ui_id`) VALUES
	(@dss_mlssId, 1),
	(@dss_mlssId, 4),
	(@dss_mlssId, 5),
	(@dss_mlssId, 6),
	(@dss_mlssId, 35),
	(@dss_mlssId, 36),
	(@dss_mlssId, 3);

