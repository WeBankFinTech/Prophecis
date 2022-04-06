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

import com.webank.wedatasphere.dss.appconn.mlflow.MLFlowAppConn;
import com.webank.wedatasphere.dss.appconn.mlflow.job.JobManager;
import com.webank.wedatasphere.dss.standard.app.development.listener.common.AsyncExecutionRequestRef;
import com.webank.wedatasphere.dss.standard.app.development.listener.common.CompletedExecutionResponseRef;
import com.webank.wedatasphere.dss.standard.app.development.listener.common.RefExecutionAction;
import com.webank.wedatasphere.dss.standard.app.development.listener.common.RefExecutionState;
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

import java.util.Map;

public class MLFlowExecutionOperation extends LongTermRefExecutionOperation implements Killable, Procedure {


    private final static Logger logger = LoggerFactory.getLogger(MLFlowExecutionOperation.class);
    DevelopmentService developmentService;
//    private SSORequestOperation<HttpAction, HttpResult> ssoRequestOperation;

    public MLFlowExecutionOperation(DevelopmentService service) {
        this.developmentService = service;
//        this.ssoRequestOperation = this.developmentService.getSSORequestService().createSSORequestOperation(getAppName());
    }

    private String getAppName() {
        return MLFlowAppConn.MLFlow_APPCONN_NAME;
    }

    @Override
    protected RefExecutionAction submit(ExecutionRequestRef requestRef) {
        AsyncExecutionRequestRef asyncRequestRef = (AsyncExecutionRequestRef) requestRef;
        Map<String,Object> jobContent = requestRef.getJobContent();
        return JobManager.submit(jobContent, asyncRequestRef.getExecutionRequestRefContext());
    }

    @Override
    public RefExecutionState state(RefExecutionAction action) {
        MLFlowExecutionAction mlFlowExecutionAction = (MLFlowExecutionAction) action;
        return JobManager.state(mlFlowExecutionAction);
    }

    @Override
    public CompletedExecutionResponseRef result(RefExecutionAction action) {
        MLFlowExecutionAction mlFlowExecutionAction = (MLFlowExecutionAction) action;
        return JobManager.result(mlFlowExecutionAction, mlFlowExecutionAction.getExecutionRequestRefContext());
    }

    private String getUser(AsyncExecutionRequestRef requestRef) {
        return requestRef.getExecutionRequestRefContext().getRuntimeMap().get("wds.dss.workflow.submit.user").toString();
    }

    private String getId(AsyncExecutionRequestRef requestRef) {
        return null;
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
