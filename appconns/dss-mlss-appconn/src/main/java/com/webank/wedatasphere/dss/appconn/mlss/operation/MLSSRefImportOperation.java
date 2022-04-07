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
import com.webank.wedatasphere.dss.appconn.mlss.publish.MLSSImportResponseRef;
import com.webank.wedatasphere.dss.appconn.mlss.restapi.ExperimentAPI;
import com.webank.wedatasphere.dss.appconn.mlss.utils.MLSSConfig;
import com.webank.wedatasphere.dss.appconn.mlss.utils.MLSSNodeUtils;
import com.webank.wedatasphere.dss.flow.execution.entrance.conf.FlowExecutionEntranceConfiguration;
import com.webank.wedatasphere.dss.standard.app.development.listener.common.AsyncExecutionRequestRef;
import com.webank.wedatasphere.dss.standard.app.development.operation.RefImportOperation;
import com.webank.wedatasphere.dss.standard.app.development.ref.ImportRequestRef;
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

public class MLSSRefImportOperation implements RefImportOperation<ImportRequestRef> {

    private final static Logger logger = LoggerFactory.getLogger(MLSSRefImportOperation.class);

    DevelopmentService developmentService;
    private SSORequestOperation ssoRequestOperation;
    private Gson gson = new Gson();

    public MLSSRefImportOperation(DevelopmentService developmentService){
        this.developmentService = developmentService;
        this.ssoRequestOperation = this.developmentService.getSSORequestService().createSSORequestOperation(getAppName());
    }

    private String getAppName() {
        return MLSSAppConn.MLSS_APPCONN_NAME;
    }

    @Override
    public ResponseRef importRef(ImportRequestRef requestRef) throws ExternalOperationFailedException {
        if(null == MLSSConfig.BASE_URL){
            this.initMLSSConfig();
        }

        Map<String,Object> jobContent = (Map<String, Object>) requestRef.getParameter("jobContent");
        MLSSImportResponseRef responseRef = null;
        Boolean resourceIdFlag = requestRef.getParameter("resourceId") == null;

        Boolean versionFlag = requestRef.getParameter("version") == null;
        if (resourceIdFlag || versionFlag) {
            //TODO: FIX RETURN NULL
            logger.error("resourceId or version is null");
            return null;
        }
        String resourceId =String.valueOf(requestRef.getParameter("resourceId"));
        String version = String.valueOf(requestRef.getParameter("version"));

        String user = requestRef.getParameter("user").toString();
        String contextId = requestRef.getParameter("contextID").toString();
        JsonObject contextMap = gson.fromJson(contextId,JsonObject.class);
//        JsonObject context  = new JsonObject(contextMap.get("value").toString());

        String flowName = null;
        String flowVersion = null;
        String flowID = null;
        String projectID = null;
        String projectName = null;
        String expDesc = "dss-appjoint," + (flowName == null ? "NULL" : flowName) + "," +
                (flowVersion == null ? "NULL" : flowVersion) + "," + flowID + "," +
                projectName + "," + projectName + "," + projectID  + "," + flowVersion;
        if (user.endsWith("_f")) {
            String[] userStrArray = user.split("_");
            user = userStrArray[0];
        }
//        StringEscapeUtils.unescapeJavaScript(contextMap.get("value").toString())
        Map<String, Object> resMap = new HashMap<>();
        JsonObject jsonObj = ExperimentAPI.importExperimentDSS(user, resourceId, version, expDesc);
        if (jsonObj == null) {
            logger.error("ImportDSS experiment failed, jsonObj is null");
            return null;
        }
        String expId = jsonObj.get("result").getAsJsonObject().get("expId").getAsString();
        jobContent.put("expId", expId);
        requestRef.setParameter("jobContent",jobContent);
        resMap.put("expId", expId);
        resMap.put("status", 200);
        try {
            //TODO ADD Default Resp
            responseRef = new MLSSImportResponseRef(resMap,"","","");
        } catch (Exception e) {
            e.printStackTrace();
        }
        return responseRef;
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
