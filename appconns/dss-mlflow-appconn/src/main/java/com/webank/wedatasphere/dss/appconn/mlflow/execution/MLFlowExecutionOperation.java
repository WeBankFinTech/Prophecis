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

package com.webank.wedatasphere.dss.appconn.mlflow.execution;

import com.webank.wedatasphere.dss.appconn.mlflow.job.JobManager;
import com.webank.wedatasphere.dss.common.utils.DSSCommonUtils;
import com.webank.wedatasphere.dss.standard.app.development.listener.common.RefExecutionAction;
import com.webank.wedatasphere.dss.standard.app.development.listener.common.RefExecutionState;
import com.webank.wedatasphere.dss.standard.app.development.listener.core.Killable;
import com.webank.wedatasphere.dss.standard.app.development.listener.core.LongTermRefExecutionOperation;
import com.webank.wedatasphere.dss.standard.app.development.listener.core.Procedure;
import com.webank.wedatasphere.dss.standard.app.development.listener.ref.ExecutionResponseRef;
import com.webank.wedatasphere.dss.standard.app.development.listener.ref.RefExecutionRequestRef;
import com.webank.wedatasphere.dss.standard.common.entity.ref.ResponseRef;
import com.webank.wedatasphere.dss.standard.common.exception.operation.ExternalOperationFailedException;

import java.util.HashMap;
import java.util.Map;

public class MLFlowExecutionOperation extends LongTermRefExecutionOperation<RefExecutionRequestRef.RefExecutionProjectWithContextRequestRef> implements Killable, Procedure {


    @Override
    protected RefExecutionAction submit(RefExecutionRequestRef.RefExecutionProjectWithContextRequestRef requestRef) throws ExternalOperationFailedException {
        Map<String, Object> jobContent = requestRef.getRefJobContent();
        HashMap contextInfo = DSSCommonUtils.COMMON_GSON.fromJson(DSSCommonUtils.COMMON_GSON.fromJson(
                requestRef.getParameter("dssContextId").toString(), HashMap.class).get("value").toString(), HashMap.class);
        return JobManager.submit(jobContent, requestRef.getExecutionRequestRefContext(), contextInfo);
    }

    @Override
    public RefExecutionState state(RefExecutionAction action) throws ExternalOperationFailedException {
        logger.info(action.toString());
        MLFlowExecutionAction mlFlowExecutionAction = (MLFlowExecutionAction) action;
        return JobManager.state(mlFlowExecutionAction);
    }

    @Override
    public ExecutionResponseRef result(RefExecutionAction action) throws ExternalOperationFailedException {
        HashMap resultMap = new HashMap();
        int status = -1;
        String errorMsg = "";
        if (action == null) {
            errorMsg = "";
            return ExecutionResponseRef.newBuilder().setResponseRef(
                    new MLFlowExecutionResponseRef("", status, errorMsg, resultMap)).build();
        }
        logger.info(action.toString());
        MLFlowExecutionAction mlssAction = (MLFlowExecutionAction) action;
        if (mlssAction.getState() == RefExecutionState.Success) {
            status = 200;
        }
        ResponseRef response = new MLFlowExecutionResponseRef("", status, errorMsg, resultMap);
        return ExecutionResponseRef.newBuilder().setResponseRef(response).build();
    }

    @Override
    public boolean kill(RefExecutionAction action) {
        MLFlowExecutionAction mlFlowExecutionAction = (MLFlowExecutionAction) action;
        return JobManager.kill(mlFlowExecutionAction);
    }

    @Override
    public float progress(RefExecutionAction action) {
        MLFlowExecutionAction mlFlowExecutionAction = (MLFlowExecutionAction) action;
        return JobManager.progress(mlFlowExecutionAction);
    }

    @Override
    public String log(RefExecutionAction action) {
        MLFlowExecutionAction mlFlowExecutionAction = (MLFlowExecutionAction) action;
        return JobManager.log(mlFlowExecutionAction);
    }
}
