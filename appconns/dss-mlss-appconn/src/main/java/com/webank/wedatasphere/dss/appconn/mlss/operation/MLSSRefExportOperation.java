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
import com.webank.wedatasphere.dss.appconn.mlss.publish.MLSSExportResponseRef;
import com.webank.wedatasphere.dss.appconn.mlss.restapi.ExperimentAPI;
import com.webank.wedatasphere.dss.appconn.mlss.utils.MLSSConfig;
import com.webank.wedatasphere.dss.appconn.mlss.utils.MLSSNodeUtils;
import com.webank.wedatasphere.dss.standard.app.development.listener.common.AsyncExecutionRequestRef;
import com.webank.wedatasphere.dss.standard.app.development.operation.RefExportOperation;
import com.webank.wedatasphere.dss.standard.app.development.ref.ExportRequestRef;
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

public class MLSSRefExportOperation implements RefExportOperation<ExportRequestRef> {

    private final static Logger logger = LoggerFactory.getLogger(MLSSRefExportOperation.class);

    DevelopmentService developmentService;
    private SSORequestOperation<HttpAction, HttpResult> ssoRequestOperation;

    public MLSSRefExportOperation(DevelopmentService developmentService){
        this.developmentService = developmentService;
        this.ssoRequestOperation = this.developmentService.getSSORequestService().createSSORequestOperation(getAppName());
    }

    private String getAppName() {
        return MLSSAppConn.MLSS_APPCONN_NAME;
    }

    @Override
    public ResponseRef exportRef(ExportRequestRef requestRef) throws ExternalOperationFailedException {
        if(null == MLSSConfig.BASE_URL){
            this.initMLSSConfig();
        }

        Map<String, Object> jobContent = (Map<String, Object>) requestRef.getParameter("jobContent");
        String user = requestRef.getParameter("user").toString();
        if (user.endsWith("_f")) {
            String[] userStrArray = user.split("_");
            user = userStrArray[0];
        }


        Map<String, Object> resMap = new HashMap<>();
        resMap.put("status",0);
        if (jobContent != null) {
            Map<String, Object>  appJointNodeMap = jobContent;
            Long expId = Float.valueOf(String.valueOf(appJointNodeMap.get("expId"))).longValue();
            JsonObject jsonObject = ExperimentAPI.exportExperimentDSS(user, String.valueOf(expId));
            String resResourceId = jsonObject.get("result").getAsJsonObject().get("resourceId").getAsString();
            String resVersion = jsonObject.get("result").getAsJsonObject().get("version").getAsString();
            logger.info("ExportDSS resourceId:" + resResourceId + " version:" + resVersion);
            resMap.put("resourceId", resResourceId);
            resMap.put("version", resVersion);
            resMap.put("status",200);
        }

        MLSSExportResponseRef responseRef = null;
        try {
            responseRef = new MLSSExportResponseRef(resMap.toString());
        } catch (Exception e) {
            //TODO: Test Exception
            e.printStackTrace();
            return null;
        }

        return responseRef;
    }

    @Override
    public void setDevelopmentService(DevelopmentService service) {
        developmentService = service;
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
