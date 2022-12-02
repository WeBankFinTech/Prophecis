package com.webank.wedatasphere.dss.appconn.mlflow.job;

import com.google.gson.JsonObject;
import com.webank.wedatasphere.dss.appconn.mlflow.restapi.ModelFactoryAPI;
import com.webank.wedatasphere.dss.appconn.mlflow.utils.APIClient;
import com.webank.wedatasphere.dss.appconn.mlflow.utils.MLFlowConfig;
import com.webank.wedatasphere.dss.standard.app.development.listener.common.RefExecutionState;
import com.webank.wedatasphere.dss.standard.app.development.listener.core.ExecutionRequestRefContext;
import org.apache.http.HttpStatus;
import org.apache.linkis.common.log.LogUtils;

import java.util.HashMap;
import java.util.Map;

import static org.apache.linkis.cs.common.utils.CSCommonUtils.gson;


public class ModelMonitoringJob extends MLFlowJob {
    public String expId;
    public Integer taskId;
    public String execId;
    public Integer reportId;
    public Integer id;

    public ModelMonitoringJob(String user) {
        this.setUser(user);
        this.setJobType(JobManager.JOB_TYPE_MODEL_MONITORING);
    }

    @Override
    public boolean submit(Map jobContent, ExecutionRequestRefContext executionRequest) {
        executionRequest.appendLog(LogUtils.generateInfo("mlflow get execution action jobContent param=" + jobContent));
        Map content = (Map) jobContent.get("content");
        if (content.isEmpty()) {
            return false;
        }
        String mlss_server = MLFlowConfig.BASE_URL;
        executionRequest.appendLog(LogUtils.generateInfo("mlflow get execution action callback url=" + mlss_server));
        return true;
    }

    @Override
    public RefExecutionState status() {
        Map<String, String> taskParams = new HashMap<String, String>();
        taskParams.put("task_id", String.valueOf(taskId));
        Map result = ModelFactoryAPI.modelMonitorReportGet(this.getUser(), taskParams);
        //Status
        int statusCode = Integer.parseInt(result.get(APIClient.REPLY_STATUS_CODE).toString());
        if (statusCode != HttpStatus.SC_OK) {
            this.getExecutionContext().appendLog("Get Status From MF Rest Error: " + statusCode);
            this.getExecutionContext().appendLog("Get Status From MF Rest Error, exp ID: " + this.expId);
            this.getExecutionContext().appendLog("Response: " + result.get(APIClient.REPLY_RESULT_CONTENT).toString());
            return RefExecutionState.Failed;
        }
        //Response
        String responseContent = result.get(APIClient.REPLY_RESULT_CONTENT).toString();
        JsonObject monitor = gson.fromJson(responseContent, JsonObject.class);
        if (null == monitor) {
            return this.transformStatus("failed");
        }
        String status = monitor.get("result").getAsJsonObject().get("status").getAsString();
        return this.transformStatus(status);
    }

    @Override
    public boolean kill() {
        return true;
    }

    @Override
    public String log() {
        return null;
    }

    @Override
    public RefExecutionState transformStatus(String status) {
        if ("running".equals(status)) {
            return RefExecutionState.Running;
        } else if ("failed".equals(status)) {
            return RefExecutionState.Failed;
        } else if ("succeed".equals(status)) {
            return RefExecutionState.Success;
        }
        return RefExecutionState.Failed;
    }


}
