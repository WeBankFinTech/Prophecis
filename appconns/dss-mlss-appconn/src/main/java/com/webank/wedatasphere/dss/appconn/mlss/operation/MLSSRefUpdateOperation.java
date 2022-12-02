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
import com.webank.wedatasphere.dss.flow.execution.entrance.conf.FlowExecutionEntranceConfiguration;
import com.webank.wedatasphere.dss.standard.app.development.operation.AbstractDevelopmentOperation;
import com.webank.wedatasphere.dss.standard.app.development.operation.RefUpdateOperation;
import com.webank.wedatasphere.dss.standard.app.development.ref.RefJobContentResponseRef;
import com.webank.wedatasphere.dss.standard.app.development.ref.impl.ThirdlyRequestRef;
import com.webank.wedatasphere.dss.standard.common.entity.ref.ResponseRef;
import com.webank.wedatasphere.dss.standard.common.exception.operation.ExternalOperationFailedException;
import java.util.Map;

public class MLSSRefUpdateOperation extends AbstractDevelopmentOperation<ThirdlyRequestRef.UpdateRequestRefImpl, ResponseRef>
        implements RefUpdateOperation<ThirdlyRequestRef.UpdateRequestRefImpl> {


    @Override
    public ResponseRef updateRef(ThirdlyRequestRef.UpdateRequestRefImpl requestRef) throws ExternalOperationFailedException {
        //1. Init
        logger.info(requestRef.toString());
        Map<String, Object> jobContent = requestRef.getRefJobContent();
        Map<String,Object> config = (Map<String,Object>)jobContent.get("runtimeMap");
        String projectName = requestRef.getProjectName();
        if (projectName == null) {
            projectName = "mlss";
        }
        //TODO: Check Varilable
        String flowName = (String) config.get(FlowExecutionEntranceConfiguration.PROJECT_NAME());
        String flowVersion = (String) config.get(FlowExecutionEntranceConfiguration.FLOW_NAME());
        String flowID = (String) config.get(FlowExecutionEntranceConfiguration.FLOW_NAME());
        String projectID = (String) config.get(FlowExecutionEntranceConfiguration.FLOW_NAME());
        String propjectName = (String) config.get(FlowExecutionEntranceConfiguration.FLOW_NAME());
        String expDesc = "DSS," + (flowName == null ? "NULL" : flowName) + "," +
                (flowVersion == null ? "NULL" : flowVersion) + "," + flowID + "," +
                projectName + "," + propjectName + "," + projectID  + "," + flowVersion;
        String user = requestRef.getUserName();
        if (user.endsWith("_f")) {
            String[] userStrArray = user.split("_");
            user = userStrArray[0];
        }
        String expName = "";
        //2. Execute Create Request
        JsonObject jsonObj = ExperimentAPI.putExperiment(user, expName, expDesc);
        if (jsonObj == null) {
            logger.error("Create experiment failed, jsonObj is null");
            return ResponseRef.newInternalBuilder().error("Create experiment failed, jsonObj is null");
        }
        //TODO: UPDATE EXPERIMENT JSON
        //3. Update node jobContent & Return Response
        Long resExpId = jsonObj.get("result").getAsJsonObject().get("id").getAsLong();
        String resExpName = jsonObj.get("result").getAsJsonObject().get("exp_name").getAsString();
        logger.info("Create MLSS experiment success expId:" + resExpId + ",expName:" + resExpName);
        jobContent.put("expId", resExpId);
        jobContent.put("expName", expName);
        jobContent.put("status", 200);
        return RefJobContentResponseRef.newBuilder().setRefJobContent(jobContent).success();
    }
}
