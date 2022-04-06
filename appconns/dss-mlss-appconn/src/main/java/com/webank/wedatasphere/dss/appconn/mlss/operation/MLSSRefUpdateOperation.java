/*
 * Copyright 2019 WeBank
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package com.webank.wedatasphere.dss.appconn.mlss.operation;

import com.google.gson.JsonObject;
import com.webank.wedatasphere.dss.appconn.mlss.MLSSAppConn;
import com.webank.wedatasphere.dss.appconn.mlss.ref.MLSSCommonResponseRef;
import com.webank.wedatasphere.dss.appconn.mlss.restapi.ExperimentAPI;
import com.webank.wedatasphere.dss.appconn.mlss.utils.MLSSNodeUtils;
import com.webank.wedatasphere.dss.flow.execution.entrance.conf.FlowExecutionEntranceConfiguration;
import com.webank.wedatasphere.dss.standard.app.development.listener.common.AsyncExecutionRequestRef;
import com.webank.wedatasphere.dss.standard.app.development.operation.RefUpdateOperation;
import com.webank.wedatasphere.dss.standard.app.development.ref.UpdateRequestRef;
import com.webank.wedatasphere.dss.standard.app.development.service.DevelopmentService;
import com.webank.wedatasphere.dss.standard.app.sso.request.SSORequestOperation;
import com.webank.wedatasphere.dss.standard.common.entity.ref.ResponseRef;
import com.webank.wedatasphere.dss.standard.common.exception.operation.ExternalOperationFailedException;
import org.apache.linkis.httpclient.request.HttpAction;
import org.apache.linkis.httpclient.response.HttpResult;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.HashMap;
import java.util.Map;

public class MLSSRefUpdateOperation implements RefUpdateOperation<UpdateRequestRef> {

    DevelopmentService developmentService;
    private SSORequestOperation<HttpAction, HttpResult> ssoRequestOperation;
    private static final Logger logger = LoggerFactory.getLogger(MLSSRefCopyOperation.class);

    public MLSSRefUpdateOperation(DevelopmentService developmentService) {
        this.developmentService = developmentService;

        this.ssoRequestOperation = developmentService.getSSORequestService().createSSORequestOperation(getAppName());
    }

    private String getAppName() {
        return MLSSAppConn.MLSS_APPCONN_NAME;
    }

    @Override
    public ResponseRef updateRef(UpdateRequestRef requestRef) throws ExternalOperationFailedException {
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
        JsonObject jsonObj = ExperimentAPI.putExperiment(user, expName, expDesc);
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

    private String getBaseUrl() {
        return developmentService.getAppInstance().getBaseUrl();
    }
}
