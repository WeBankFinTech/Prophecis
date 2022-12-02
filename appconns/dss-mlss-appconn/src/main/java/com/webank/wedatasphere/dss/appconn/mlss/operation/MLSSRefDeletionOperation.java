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
import com.webank.wedatasphere.dss.appconn.mlss.restapi.ExperimentAPI;
import com.webank.wedatasphere.dss.appconn.mlss.utils.MLSSConfig;
import com.webank.wedatasphere.dss.standard.app.development.operation.AbstractDevelopmentOperation;
import com.webank.wedatasphere.dss.standard.app.development.operation.RefDeletionOperation;
import com.webank.wedatasphere.dss.standard.app.development.ref.impl.OnlyDevelopmentRequestRef;
import com.webank.wedatasphere.dss.standard.app.development.ref.impl.ThirdlyRequestRef;
import com.webank.wedatasphere.dss.standard.app.development.service.DevelopmentService;
import com.webank.wedatasphere.dss.standard.common.entity.ref.ResponseRef;
import com.webank.wedatasphere.dss.standard.common.exception.operation.ExternalOperationFailedException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Map;


public class MLSSRefDeletionOperation extends AbstractDevelopmentOperation<ThirdlyRequestRef.RefJobContentRequestRefImpl, ResponseRef>
        implements RefDeletionOperation<ThirdlyRequestRef.RefJobContentRequestRefImpl> {

    private DevelopmentService developmentService;
    private final static Logger logger = LoggerFactory.getLogger(MLSSRefCreationOperation.class);

    public MLSSRefDeletionOperation(DevelopmentService service) {
        this.developmentService = service;
    }

//    @Override
//    public void setDevelopmentService(DevelopmentService service) {
//        this.developmentService = service;
//    }

    protected void initMLSSConfig(){
        MLSSConfig.BASE_URL = this.developmentService.getAppInstance().getBaseUrl();
        Map<String, Object> config =  this.developmentService.getAppInstance().getConfig();
        MLSSConfig.APP_KEY = String.valueOf(config.get("MLSS-SecretKey"));
        MLSSConfig.APP_SIGN = String.valueOf(config.get("MLSS-APPSignature"));
        MLSSConfig.AUTH_TYPE =  String.valueOf(config.get("MLSS-Auth-Type"));
        MLSSConfig.TIMESTAMP =  String.valueOf(config.get("MLSS-APPSignature"));
    }

    @Override
    public ResponseRef deleteRef(ThirdlyRequestRef.RefJobContentRequestRefImpl requestRef) throws ExternalOperationFailedException {
        logger.info(requestRef.toString());
        if(null == MLSSConfig.BASE_URL){
            this.initMLSSConfig();
        }
        // 1. Init Params
        Map<String, Object> jobContent = (Map<String, Object>) requestRef.getRefJobContent();
        String user = requestRef.getUserName();
        Long expId = Long.parseLong(String.valueOf(jobContent.get("expId")));
        if (user.endsWith("_f")) {
            String[] userStrArray = user.split("_");
            user = userStrArray[0];
        }
        // 2. Execute Request
        JsonObject jsonObj = ExperimentAPI.deleteExperiment(user, String.valueOf(expId));
        if (jsonObj == null) {
            logger.error("Delete experiment failed, jsonObj is null");
            return ResponseRef.newInternalBuilder().error("Create experiment failed, jsonObj is null");
        }
        if (jsonObj.get("code").getAsString() != "200"){
            return ResponseRef.newInternalBuilder().error("Unknown task type" + jsonObj.getAsString());
        }
        logger.info("Delete experiment(expId:"+ expId +") success");
        return ResponseRef.newInternalBuilder().success();
    }
}
