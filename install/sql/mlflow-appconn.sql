select @dss_appconn_mlflowId:=id from `dss_appconn` where `appconn_name` = 'mlflow';
delete from `dss_appconn_instance` where  `appconn_id`=@dss_appconn_mlflowId;

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
            `resource`)
            VALUES (
                'mlflow',
                0,
                1,
                NULL,
                0,
                NULL,
                'com.webank.wedatasphere.dss.appconn.mlflow.MLFlowAppConn',
                '/appcom/Install/DSSInstall/dss-1.1.2/dss-appconns/mlflow',
                '');


select @dss_appconn_mlflowId:=id from `dss_appconn` where `appconn_name` = 'mlflow';




INSERT INTO `dss_appconn_instance`(
             `appconn_id`,
             `label`,
            `url`,
             `enhance_json`,
             `homepage_url`,
             `redirect_url`)
             VALUES (
             @dss_appconn_mlflowId,
             'DEV',
             'http://127.0.0.1:30803/',
             '{"MLSS-SecretKey":"MLFLOW","MLSS-Auth-Type":"SYSTEM","MLSS-APPSignature":"MLFLOW","MLSS-BaseUrl":"http://127.0.0.1:30803","baseUrl":"http://127.0.0.1:30803"}'
             '',
             'http://127.0.0.1:30803/#/dashboard',
             'http://127.0.0.1:30803/#/mlFlow');
