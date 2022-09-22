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
import com.webank.wedatasphere.dss.common.utils.DSSCommonUtils;
import com.webank.wedatasphere.dss.orchestrator.common.entity.DSSOrchestratorInfo;
import com.webank.wedatasphere.dss.orchestrator.common.ref.OrchestratorRefConstant;
import com.webank.wedatasphere.dss.standard.app.development.operation.AbstractDevelopmentOperation;
import com.webank.wedatasphere.dss.standard.app.development.operation.RefCreationOperation;
import com.webank.wedatasphere.dss.standard.app.development.ref.RefJobContentResponseRef;
import com.webank.wedatasphere.dss.standard.app.development.ref.impl.OnlyDevelopmentRequestRef;
import com.webank.wedatasphere.dss.standard.app.development.ref.impl.ThirdlyRequestRef;
import com.webank.wedatasphere.dss.standard.app.development.service.DevelopmentService;
import com.webank.wedatasphere.dss.standard.common.exception.operation.ExternalOperationFailedException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.HashMap;
import java.util.Map;

public class MLSSRefCreationOperation extends AbstractDevelopmentOperation<ThirdlyRequestRef.DSSJobContentWithContextRequestRef, RefJobContentResponseRef>
        implements RefCreationOperation<ThirdlyRequestRef.DSSJobContentWithContextRequestRef> {

    private final static Logger logger = LoggerFactory.getLogger(MLSSRefCreationOperation.class);
    DevelopmentService developmentService;

    public MLSSRefCreationOperation(DevelopmentService service) {
        this.developmentService = service;
        this.initMLSSConfig();
    }


    @Override
    public RefJobContentResponseRef createRef(ThirdlyRequestRef.DSSJobContentWithContextRequestRef requestRef) throws ExternalOperationFailedException {
        if (null == MLSSConfig.BASE_URL) {
            this.initMLSSConfig();
        }
        logger.info(requestRef.toString());
        //1. Init
        HashMap<String, Object> jobContent = (HashMap<String, Object>) requestRef.getDSSJobContent();
        DSSOrchestratorInfo dssOrchestratorInfo = (DSSOrchestratorInfo) requestRef.getDSSJobContent().get(OrchestratorRefConstant.DSS_ORCHESTRATOR_INFO_KEY);
        HashMap contextInfo = DSSCommonUtils.COMMON_GSON.fromJson(DSSCommonUtils.COMMON_GSON.fromJson(
                requestRef.getParameter("dssContextId").toString(), HashMap.class).get("value").toString(), HashMap.class);
        String flowName = jobContent.get("orchestrationName").toString();
        String flowVersion = contextInfo.get("version").toString();
        Long flowID = Long.parseLong(jobContent.get("orchestrationId").toString());
        Long projectID = Long.parseLong(requestRef.getParameter("dssProjectId").toString());
        String projectName = contextInfo.get("project").toString();

        String expDesc = "DSS," + (flowName == null ? "NULL" : flowName) + "," +
                (flowVersion == null ? "NULL" : flowVersion) + "," + flowID + "," +
                projectName + ","  + projectID;
        String user = requestRef.getUserName();
        if (user.endsWith("_f")) {
            String[] userStrArray = user.split("_");
            user = userStrArray[0];
        }
        String expName = requestRef.getName();


        //2. Execute Create Request
        JsonObject jsonObj = ExperimentAPI.postExperiment(user, expName, expDesc);
        if (jsonObj == null) {
            logger.error("Create experiment failed, jsonObj is null");
            return null;
        }
        //3. Update node jobContent & Return Response
        Long resExpId = jsonObj.get("result").getAsJsonObject().get("id").getAsLong();
        String resExpName = jsonObj.get("result").getAsJsonObject().get("exp_name").getAsString();
        logger.info("Create MLSS experiment success expId:" + resExpId + ",expName:" + resExpName);

        jobContent.put("expId", resExpId);
        jobContent.put("expName", expName);
        jobContent.put("status", 200);
        return RefJobContentResponseRef.newBuilder().setRefJobContent(jobContent).success();
    }


    protected void initMLSSConfig() {
        MLSSConfig.BASE_URL = developmentService.getAppInstance().getBaseUrl();
        Map<String, Object> config = this.developmentService.getAppInstance().getConfig();
        MLSSConfig.APP_KEY = String.valueOf(config.get("MLSS-SecretKey"));
        MLSSConfig.APP_SIGN = String.valueOf(config.get("MLSS-APPSignature"));
        MLSSConfig.AUTH_TYPE = String.valueOf(config.get("MLSS-Auth-Type"));
        MLSSConfig.TIMESTAMP = String.valueOf(config.get("MLSS-APPSignature"));
    }
}
