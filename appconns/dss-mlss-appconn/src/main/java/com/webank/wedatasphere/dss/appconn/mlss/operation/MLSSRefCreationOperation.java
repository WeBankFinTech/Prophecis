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

import com.google.gson.Gson;
import com.google.gson.JsonObject;
import com.webank.wedatasphere.dss.appconn.mlss.MLSSAppConn;
import com.webank.wedatasphere.dss.appconn.mlss.ref.MLSSCommonResponseRef;
import com.webank.wedatasphere.dss.appconn.mlss.utils.MLSSConfig;
import com.webank.wedatasphere.dss.common.utils.DSSCommonUtils;
import com.webank.wedatasphere.dss.flow.execution.entrance.conf.FlowExecutionEntranceConfiguration;
import com.webank.wedatasphere.dss.appconn.mlss.restapi.ExperimentAPI;
import com.webank.wedatasphere.dss.standard.app.development.listener.common.AsyncExecutionRequestRef;
import com.webank.wedatasphere.dss.standard.app.development.operation.RefCreationOperation;
import com.webank.wedatasphere.dss.standard.app.development.ref.CreateRequestRef;
import com.webank.wedatasphere.dss.standard.app.development.ref.NodeRequestRef;
import com.webank.wedatasphere.dss.standard.app.development.service.DevelopmentService;
import com.webank.wedatasphere.dss.standard.app.sso.request.SSORequestOperation;
import com.webank.wedatasphere.dss.standard.common.entity.ref.ResponseRef;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.HashMap;
import java.util.Map;

public class MLSSRefCreationOperation implements RefCreationOperation<CreateRequestRef> {
    private final static Logger logger = LoggerFactory.getLogger(MLSSRefCreationOperation.class);

    DevelopmentService developmentService;
    private SSORequestOperation ssoRequestOperation;

    public MLSSRefCreationOperation(DevelopmentService service){
        this.developmentService = service;
        this.ssoRequestOperation = this.developmentService.getSSORequestService().createSSORequestOperation(getAppName());
        this.initMLSSConfig();
    }


    private String getAppName() {
        return MLSSAppConn.MLSS_APPCONN_NAME;
    }

    @Override
    public ResponseRef createRef(CreateRequestRef requestRef) {
        if(null == MLSSConfig.BASE_URL){
            this.initMLSSConfig();
        }
        //1. Init
        NodeRequestRef nodeRequest = (NodeRequestRef) requestRef;
        //TODO: Check Varilable
        HashMap contextInfo = DSSCommonUtils.COMMON_GSON.fromJson(DSSCommonUtils.COMMON_GSON.fromJson
                (nodeRequest.getJobContent().get("contextID").toString(),HashMap.class).get("value").toString(),HashMap.class);
        String flowName = nodeRequest.getOrcName();
        String flowVersion = contextInfo.get("version").toString();
        Long flowID = nodeRequest.getOrcId();
        Long projectID = nodeRequest.getProjectId();
        String projectName = contextInfo.get("project").toString();

        String expDesc = "dss-appjoint," + (flowName == null ? "NULL" : flowName) + "," +
                (flowVersion == null ? "NULL" : flowVersion) + "," + flowID + "," +
                projectName + ","  + projectID  ;
        String user = nodeRequest.getUserName();
        if (user.endsWith("_f")) {
            String[] userStrArray = user.split("_");
            user = userStrArray[0];
        }
        String expName = nodeRequest.getName();


        //2. Execute Create Request
        Map<String, Object> resMap = new HashMap<>();
        JsonObject jsonObj = ExperimentAPI.postExperiment(user, expName, expDesc);
        if (jsonObj == null) {
            logger.error("Create experiment failed, jsonObj is null");
            return null;
        }


        //3. Update node jobContent & Return Response
        Long resExpId = jsonObj.get("result").getAsJsonObject().get("id").getAsLong();
        String resExpName = jsonObj.get("result").getAsJsonObject().get("exp_name").getAsString();
        logger.info("Create MLSS experiment success expId:" + resExpId + ",expName:" + resExpName);
        Map<String,Object> jobContent = (Map<String,Object>) nodeRequest.getJobContent();
        if (null == jobContent) {
            jobContent = new HashMap<String,Object>();
        }
        jobContent.put("expId", resExpId);
        jobContent.put("expName", expName);
        jobContent.put("status", 200);
        nodeRequest.setJobContent(jobContent);
        Gson gson = new Gson();
        try {
            return new MLSSCommonResponseRef(gson.toJson(jobContent));
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

    private String getBaseUrl(){
        return developmentService.getAppInstance().getBaseUrl();
    }


    protected void initMLSSConfig(){
        MLSSConfig.BASE_URL = this.developmentService.getAppInstance().getBaseUrl();
        Map<String, Object> config =  this.developmentService.getAppInstance().getConfig();
        MLSSConfig.APP_KEY = String.valueOf(config.get("MLSS-SecretKey"));
        MLSSConfig.APP_SIGN = String.valueOf(config.get("MLSS-APPSignature"));
        MLSSConfig.AUTH_TYPE =  String.valueOf(config.get("MLSS-Auth-Type"));
        MLSSConfig.TIMESTAMP =  String.valueOf(config.get("MLSS-APPSignature"));
    }

}
