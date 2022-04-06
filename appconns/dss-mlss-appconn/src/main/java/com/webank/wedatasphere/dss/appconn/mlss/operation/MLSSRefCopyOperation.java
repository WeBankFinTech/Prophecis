package com.webank.wedatasphere.dss.appconn.mlss.operation;

import com.google.gson.JsonObject;
import com.webank.wedatasphere.dss.appconn.mlss.MLSSAppConn;
import com.webank.wedatasphere.dss.appconn.mlss.ref.MLSSCommonResponseRef;
import com.webank.wedatasphere.dss.appconn.mlss.ref.MLSSCopyRequestRef;
import com.webank.wedatasphere.dss.appconn.mlss.ref.MLSSCopyResponseRef;
import com.webank.wedatasphere.dss.appconn.mlss.restapi.ExperimentAPI;
import com.webank.wedatasphere.dss.appconn.mlss.utils.MLSSNodeUtils;
import com.webank.wedatasphere.dss.flow.execution.entrance.conf.FlowExecutionEntranceConfiguration;
import com.webank.wedatasphere.dss.standard.app.development.listener.common.AsyncExecutionRequestRef;
import com.webank.wedatasphere.dss.standard.app.development.operation.RefCopyOperation;
import com.webank.wedatasphere.dss.standard.app.development.service.DevelopmentService;
import com.webank.wedatasphere.dss.standard.app.sso.builder.SSOUrlBuilderOperation;
import com.webank.wedatasphere.dss.standard.app.sso.request.SSORequestOperation;
import com.webank.wedatasphere.dss.standard.common.desc.AppInstance;
import com.webank.wedatasphere.dss.standard.common.entity.ref.ResponseRef;
import com.webank.wedatasphere.dss.standard.common.exception.operation.ExternalOperationFailedException;
import org.apache.linkis.httpclient.request.HttpAction;
import org.apache.linkis.httpclient.response.HttpResult;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.HashMap;
import java.util.Map;


public class MLSSRefCopyOperation implements RefCopyOperation<MLSSCopyRequestRef> {

    private static final Logger logger = LoggerFactory.getLogger(MLSSRefCopyOperation.class);

    private DevelopmentService developmentService;

    private AppInstance appInstance;

    private SSORequestOperation<HttpAction, HttpResult> ssoRequestOperation;

    @SuppressWarnings("unchecked")
    public MLSSRefCopyOperation(AppInstance appInstance, DevelopmentService developmentService) {
        this.appInstance = appInstance;
        this.developmentService = developmentService;
        this.ssoRequestOperation = this.developmentService.getSSORequestService().createSSORequestOperation(getAppName());
    }

    private String getAppName() {
        return MLSSAppConn.MLSS_APPCONN_NAME;
    }

    @Override
    @SuppressWarnings("unchecked")
    public ResponseRef copyRef(MLSSCopyRequestRef requestRef) throws ExternalOperationFailedException {

        //1. Init
        AsyncExecutionRequestRef executionRequestRef = (AsyncExecutionRequestRef) requestRef;
        Map<String,Object> config = executionRequestRef.getExecutionRequestRefContext().getRuntimeMap();
        String projectName = requestRef.getName() ;
        if (projectName == null) {
            projectName = "mlss";
        }
        //TODO: Check Varilable
        String flowName = (String) config.get(FlowExecutionEntranceConfiguration.PROJECT_NAME());
        String flowVersion = (String) config.get(FlowExecutionEntranceConfiguration.FLOW_NAME());
        String flowID = (String) config.get(FlowExecutionEntranceConfiguration.FLOW_NAME());
        String projectID = (String) config.get(FlowExecutionEntranceConfiguration.FLOW_NAME());
        String propjectName = (String) config.get(FlowExecutionEntranceConfiguration.FLOW_NAME());
        String expDesc = "dss-appjoint," + (flowName == null ? "NULL" : flowName) + "," +
                (flowVersion == null ? "NULL" : flowVersion) + "," + flowID + "," +
                projectName + "," + propjectName + "," + projectID  + "," + flowVersion;
        String user = MLSSNodeUtils.getUser(executionRequestRef.getExecutionRequestRefContext());
        if (user.endsWith("_f")) {
            String[] userStrArray = user.split("_");
            user = userStrArray[0];
        }
        String expName = "";


        //2. Execute Create Request
        Map<String, Object> resMap = new HashMap<>();
        JsonObject jsonObj = ExperimentAPI.postExperiment(user, expName, expDesc);
        if (jsonObj == null) {
            logger.error("Create experiment failed, jsonObj is null");
            return null;
        }
        //TODO: UPDATE EXPERIMENT JSON


        //3. Update node jobContent & Return Response
        Long resExpId = jsonObj.get("result").getAsJsonObject().get("id").getAsLong();
        String resExpName = jsonObj.get("result").getAsJsonObject().get("exp_name").getAsString();
        logger.info("Create MLSS experiment success expId:" + resExpId + ",expName:" + resExpName);
        Map<String,Object> jobContent = (Map<String,Object>) executionRequestRef.getJobContent();
        if (null == jobContent) {
            jobContent = new HashMap<String,Object>();
        }
        jobContent.put("expId", resExpId);
        jobContent.put("expName", expName);
        jobContent.put("status", 200);
        executionRequestRef.setJobContent(jobContent);
        try {
            return new MLSSCommonResponseRef(jobContent.toString());
        } catch (Exception e) {
            e.printStackTrace();
            //TODO
            return null;
        }
    }


    @Override
    public void setDevelopmentService(DevelopmentService service) {
        this.developmentService = service;
    }
}
