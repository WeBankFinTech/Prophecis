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
import com.webank.wedatasphere.dss.appconn.mlss.MLSSAppConn;
import com.webank.wedatasphere.dss.appconn.mlss.restapi.ExperimentRunAPI;
import com.webank.wedatasphere.dss.appconn.mlss.utils.MLSSNodeUtils;
import com.webank.wedatasphere.dss.flow.execution.entrance.node.NodeExecutionState;
import com.webank.wedatasphere.dss.standard.app.development.listener.common.AsyncExecutionRequestRef;
import com.webank.wedatasphere.dss.standard.app.development.listener.common.CompletedExecutionResponseRef;
import com.webank.wedatasphere.dss.standard.app.development.listener.common.RefExecutionAction;
import com.webank.wedatasphere.dss.standard.app.development.listener.common.RefExecutionState;
import com.webank.wedatasphere.dss.standard.app.development.listener.core.ExecutionRequestRefContext;
import com.webank.wedatasphere.dss.standard.app.development.listener.core.Killable;
import com.webank.wedatasphere.dss.standard.app.development.listener.core.LongTermRefExecutionOperation;
import com.webank.wedatasphere.dss.standard.app.development.listener.core.Procedure;
import com.webank.wedatasphere.dss.standard.app.development.ref.ExecutionRequestRef;
import com.webank.wedatasphere.dss.standard.app.development.service.DevelopmentService;
import com.webank.wedatasphere.dss.standard.app.sso.request.SSORequestOperation;
import org.apache.linkis.httpclient.request.HttpAction;
import org.apache.linkis.httpclient.response.HttpResult;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.math.BigDecimal;
import java.util.Map;

public class MLSSRefExecutionOperation extends LongTermRefExecutionOperation implements Procedure, Killable {

    private final static Logger logger = LoggerFactory.getLogger(MLSSRefExecutionOperation.class);

    DevelopmentService developmentService;
    private SSORequestOperation<HttpAction, HttpResult> ssoRequestOperation;

    public MLSSRefExecutionOperation(DevelopmentService service) {
        this.developmentService = service;
        this.ssoRequestOperation = this.developmentService.getSSORequestService().createSSORequestOperation(getAppName());
    }

    @Override
    protected RefExecutionAction submit(ExecutionRequestRef requestRef) {
        AsyncExecutionRequestRef asyncRequestRef = (AsyncExecutionRequestRef) requestRef;
        Map<String,Object> jobContent = (Map<String, Object>) requestRef.getJobContent();
        String expId = jobContent == null ? "" : Float.valueOf(jobContent.get("expId").toString()).intValue() + ""  ;
        //TODO: Add Exception Log
        if (expId.equals("-1") || expId.equals("")) return null;
        if (expId.length() <= 0) {
            logger.error("appJointNode is null");
            return null;
        }
        String user = MLSSNodeUtils.getUser(asyncRequestRef.getExecutionRequestRefContext());
        if (user.endsWith("_f")) {
            String[] userStrArray = user.split("_");
            user = userStrArray[0];
        }


        JsonObject jsonObj = ExperimentRunAPI.post(user,expId,"DSS","");
        if (jsonObj == null) {
            logger.error("Create experiment run failed, jsonObj is null");
            return null;
        }
        String resExpRunId = jsonObj.get("result").getAsJsonObject().get("id").getAsString();
        MLSSExecutionAction mlssNodeExecutionAction = new MLSSExecutionAction(resExpRunId, user,
                asyncRequestRef.getExecutionRequestRefContext());
        logger.info("Create experiment run success, expRunId:" + resExpRunId + "expExecType:" + "DSS");
        return mlssNodeExecutionAction;
    }

    @Override
    public RefExecutionState state(RefExecutionAction action) {
        if (action == null) {
            logger.error("action is null");
            return RefExecutionState.Failed;
        }
        MLSSExecutionAction mlssAction = ((MLSSExecutionAction) action);
        String user = mlssAction.getUser();
        String appId =mlssAction.getAppId();
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
    public CompletedExecutionResponseRef result(RefExecutionAction action) {
        MLSSCompletedExecutionResponseRef result = new MLSSCompletedExecutionResponseRef(0,"");
        result.setIsSucceed(false);
        if (action == null) {
            result.setErrorMsg("MLSS Execution Action is Null");
            return result;
        }
        MLSSExecutionAction mlssAction = (MLSSExecutionAction) action;
        if (mlssAction.getState() == RefExecutionState.Success){
            result.setIsSucceed(true);
        }
        return result;
    }


    @Override
    public void setDevelopmentService(DevelopmentService service) {
        this.developmentService = service;
    }

    private String getBaseUrl(){
        return developmentService.getAppInstance().getBaseUrl();
    }

    @Override
    public boolean kill(RefExecutionAction action) {
        if (action == null) {
            logger.error("nodeExecutionAction is null");
            return false;
        }

        String user = ((MLSSExecutionAction) action).getUser();
        String appId = ((MLSSExecutionAction) action).getAppId();
        JsonObject jsonObj = ExperimentRunAPI.kill(user, appId);
        if (jsonObj != null) {
            logger.info("Kill experiment run success,expRunId:" + appId);
            return true;
        }
        //TODO 判断是否正常返回
        //retry kill
        logger.info("Kill experiment run failed, retry kill experiment run expRunId" + appId);
        return false;

    }

    @Override
    public float progress(RefExecutionAction action) {
        if (action == null) {
            return 1;
        }
        String user = ((MLSSExecutionAction) action).getUser();
        RefExecutionState state = this.state(action);
        return state != RefExecutionState.Success && state != RefExecutionState.Failed ? 0.5f : 1f;
    }

    @Override
    public String log(RefExecutionAction action) {
        String user = ((MLSSExecutionAction) action).getUser();
        MLSSExecutionAction mlssAction = ((MLSSExecutionAction) action);
        ExecutionRequestRefContext nodeContext = mlssAction.getExecutionContext();
        String appId = mlssAction.getAppId();
        JsonObject jsonObj = ExperimentRunAPI.getLogs(user, appId,0L, 0L);
        //TODO: Check
        String logs = jsonObj.get("result").getAsJsonObject().get("log").getAsJsonArray().get(3).getAsString();
        nodeContext.appendLog(logs);
        return logs;
    }

    private String getAppName() {
        return MLSSAppConn.MLSS_APPCONN_NAME;
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
