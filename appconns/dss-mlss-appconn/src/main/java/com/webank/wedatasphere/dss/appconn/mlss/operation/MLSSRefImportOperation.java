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
import com.webank.wedatasphere.dss.appconn.mlss.restapi.ExperimentAPI;
import com.webank.wedatasphere.dss.appconn.mlss.utils.MLSSConfig;
import com.webank.wedatasphere.dss.common.utils.DSSCommonUtils;
import com.webank.wedatasphere.dss.standard.app.development.operation.AbstractDevelopmentOperation;
import com.webank.wedatasphere.dss.standard.app.development.operation.RefImportOperation;
import com.webank.wedatasphere.dss.standard.app.development.ref.RefJobContentResponseRef;
import com.webank.wedatasphere.dss.standard.app.development.ref.impl.ThirdlyRequestRef;
import com.webank.wedatasphere.dss.standard.app.development.service.DevelopmentService;
import com.webank.wedatasphere.dss.standard.common.exception.operation.ExternalOperationFailedException;

import java.util.HashMap;
import java.util.Map;

public class MLSSRefImportOperation extends AbstractDevelopmentOperation<ThirdlyRequestRef.ImportWitContextRequestRefImpl, RefJobContentResponseRef>
        implements RefImportOperation<ThirdlyRequestRef.ImportWitContextRequestRefImpl> {
    DevelopmentService developmentService;
    private Gson gson = new Gson();

    public MLSSRefImportOperation(DevelopmentService developmentService){
        this.developmentService = developmentService;
    }

    @Override
    public RefJobContentResponseRef importRef(ThirdlyRequestRef.ImportWitContextRequestRefImpl requestRef) throws ExternalOperationFailedException {
        logger.info(requestRef.toString());
        if(null == MLSSConfig.BASE_URL){
            this.initMLSSConfig();
        }
        Map<String,Object> jobContent = requestRef.getRefJobContent();
        HashMap resourceMap = (HashMap) requestRef.getParameter("resourceMap");
        Boolean resourceIdFlag = resourceMap.get("resourceId") == null;
        Boolean versionFlag = resourceMap.get("version") == null;
        if (resourceIdFlag || versionFlag) {
            //TODO: FIX RETURN NULL
            logger.error("resourceId or version is null");
            return null;
        }
        String resourceId =String.valueOf(resourceMap.get("resourceId").toString());
        String version = String.valueOf(resourceMap.get("version").toString());

        String user = requestRef.getParameter("userName").toString();
        String contextId = requestRef.getParameter("dssContextId").toString();

        HashMap contextInfo = DSSCommonUtils.COMMON_GSON.fromJson(DSSCommonUtils.COMMON_GSON.fromJson(
        requestRef.getParameter("dssContextId").toString(), HashMap.class).get("value").toString(), HashMap.class);
        String flowName = jobContent.get("orchestrationName").toString();
        String flowVersion = contextInfo.get("version").toString();
        Long flowID = Long.parseLong(jobContent.get("orchestrationId").toString());
        Long projectID = Long.parseLong(requestRef.getParameter("dssProjectId").toString());
        String projectName = contextInfo.get("project").toString();
        String expDesc = "WTSS," + (flowName == null ? "NULL" : flowName) + "," +
                (flowVersion == null ? "NULL" : flowVersion) + "," + flowID + "," +
                projectName + ","  + projectID;
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
        return RefJobContentResponseRef.newBuilder().setRefJobContent(resMap).success();
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
