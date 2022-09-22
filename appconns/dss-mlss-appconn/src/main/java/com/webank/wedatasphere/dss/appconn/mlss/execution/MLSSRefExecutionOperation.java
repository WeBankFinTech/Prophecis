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

package com.webank.wedatasphere.dss.appconn.mlss.execution;

import com.google.gson.JsonObject;
import com.webank.wedatasphere.dss.appconn.mlss.ref.MLSSResponseRefBuilder;
import com.webank.wedatasphere.dss.appconn.mlss.restapi.ExperimentRunAPI;
import com.webank.wedatasphere.dss.common.utils.DSSCommonUtils;
import com.webank.wedatasphere.dss.standard.app.development.listener.common.RefExecutionAction;
import com.webank.wedatasphere.dss.standard.app.development.listener.common.RefExecutionState;
import com.webank.wedatasphere.dss.standard.app.development.listener.core.LongTermRefExecutionOperation;
import com.webank.wedatasphere.dss.standard.app.development.listener.ref.ExecutionResponseRef;
import com.webank.wedatasphere.dss.standard.app.development.listener.ref.RefExecutionRequestRef;
import com.webank.wedatasphere.dss.standard.common.entity.ref.ResponseRef;
import com.webank.wedatasphere.dss.standard.common.exception.operation.ExternalOperationFailedException;

import java.util.HashMap;
import java.util.Map;

public class MLSSRefExecutionOperation extends LongTermRefExecutionOperation<RefExecutionRequestRef.RefExecutionProjectWithContextRequestRef> {


    @Override
    protected RefExecutionAction submit(RefExecutionRequestRef.RefExecutionProjectWithContextRequestRef requestRef) throws ExternalOperationFailedException {
        logger.info(requestRef.toString());

       requestRef.getUserName();

        Map<String, Object> jobContent = requestRef.getRefJobContent();
        String expId = jobContent == null ? "" : Float.valueOf(jobContent.get("expId").toString()).intValue() + "";
        if (expId.equals("-1") || expId.equals("") || expId.length() <= 0) {
            logger.error("appJointNode is null");
            return null;
        }

//        String user = requestRef.getUserName();
        //TODO: User from context?
        HashMap contextInfo = DSSCommonUtils.COMMON_GSON.fromJson(DSSCommonUtils.COMMON_GSON.fromJson(
                requestRef.getParameter("dssContextId").toString(), HashMap.class).get("value").toString(), HashMap.class);
        String user = contextInfo.get("user").toString();
        if (user.endsWith("_f")) {
            String[] userStrArray = user.split("_");
            user = userStrArray[0];
        }
        JsonObject jsonObj = ExperimentRunAPI.post(user, expId, "DSS", "");
        if (jsonObj == null) {
            logger.error("Create experiment run failed, jsonObj is null");
            return null;
        }
        String resExpRunId = jsonObj.get("result").getAsJsonObject().get("id").getAsString();
        MLSSExecutionAction mlssNodeExecutionAction = new MLSSExecutionAction(resExpRunId, user,
                requestRef.getExecutionRequestRefContext());
        logger.info("Create experiment run success, expRunId:" + resExpRunId + "expExecType:" + "DSS");
        return mlssNodeExecutionAction;
    }

    @Override
    public RefExecutionState state(RefExecutionAction action) throws ExternalOperationFailedException {
        if (action == null) {
            logger.error("action is null");
            return RefExecutionState.Failed;
        }
        MLSSExecutionAction mlssAction = ((MLSSExecutionAction) action);
        String user = mlssAction.getUser();
        String appId = mlssAction.getAppId();
        JsonObject jsonObj = ExperimentRunAPI.getStatus(user, appId);
        if (jsonObj == null) {
            logger.error("Get experiment run status failed, jsonObj is null");
            return RefExecutionState.Failed;
        }
        Boolean statusFlag = jsonObj.get("result").getAsJsonObject().get("status") == null;
        if (statusFlag) {
            logger.info("Get experiment run status success, state:" + RefExecutionState.Failed);
            return RefExecutionState.Failed;
        }
        String status = jsonObj.get("result").getAsJsonObject().get("status").getAsString();
        mlssAction.setState(convertExperimentStatus(status));
        return mlssAction.getState();
    }

    @Override
    public ExecutionResponseRef result(RefExecutionAction action) throws ExternalOperationFailedException {
        HashMap resultMap = new HashMap();
        int status = -1;
        String errorMsg = "";
        if (action == null) {
            errorMsg = "";
            return ExecutionResponseRef.newBuilder().setResponseRef(
                    new MLSSExecutionResponseRef("", status, errorMsg, resultMap)).build();
        }
        MLSSExecutionAction mlssAction = (MLSSExecutionAction) action;
        if (mlssAction.getState() == RefExecutionState.Success) {
            status = 200;
        }
        ResponseRef response = new MLSSExecutionResponseRef("", status, errorMsg, resultMap);
        return ExecutionResponseRef.newBuilder().setResponseRef(response).build();
    }

    private RefExecutionState convertExperimentStatus(String statusDesc) {
        String status = statusDesc.toUpperCase();
        if (status.equals("INITED") || status.equals("SCHEDULED"))
            return RefExecutionState.Accepted;
        else if (status.equals("RUNNING"))
            return RefExecutionState.Running;
        else if (status.equals("SUCCEED"))
            return RefExecutionState.Success;
        else if (status.equals("FAILED") || status.equals("TIMEOUT"))
            return RefExecutionState.Failed;
        else if (status.equals("CANCELLED"))
            return RefExecutionState.Killed;
        return RefExecutionState.Failed;
    }
}
