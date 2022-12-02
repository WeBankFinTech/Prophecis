package com.webank.wedatasphere.dss.appconn.mlflow.job;

import com.google.gson.JsonObject;
import com.google.gson.internal.LinkedTreeMap;
import com.webank.wedatasphere.dss.appconn.mlflow.restapi.ModelFactoryAPI;
import com.webank.wedatasphere.dss.appconn.mlflow.utils.APIClient;
import com.webank.wedatasphere.dss.appconn.mlflow.utils.MLFlowNodeUtils;
import com.webank.wedatasphere.dss.standard.app.development.listener.common.RefExecutionState;
import com.webank.wedatasphere.dss.standard.app.development.listener.core.ExecutionRequestRefContext;
import org.apache.http.HttpStatus;

import java.util.Map;

import static org.apache.linkis.cs.common.utils.CSCommonUtils.gson;

public class ModelStorageJob extends MLFlowJob {

    private String modelId;

    public ModelStorageJob(String user){
        this.setUser(user);
        this.setJobType(JobManager.JOB_TYPE_MODEL_STORAGE);
    }

    @Override
    public boolean submit(Map jobContent, ExecutionRequestRefContext executionRequestRefContext) {

        //Prepare job submit params
        LinkedTreeMap content = (LinkedTreeMap) jobContent.get("content");
        JobManager.paramsTransfer(content);
        JsonObject jobJsonObj = gson.toJsonTree(content).getAsJsonObject();
        String jsonBody = jobJsonObj.toString();
        this.setExecutionContext(executionRequestRefContext);

       //Execute Job
        String user = this.getUser();

        Map result = ModelFactoryAPI.ModelPost(user,jsonBody);
        int statusCode = Integer.parseInt(result.get(APIClient.REPLY_STATUS_CODE).toString());
        if (statusCode !=  HttpStatus.SC_OK) {
            executionRequestRefContext.appendLog("Get Status From DI Rest Error: " + statusCode);
            executionRequestRefContext.appendLog("Response: " + result.get(APIClient.REPLY_RESULT_CONTENT).toString());
            this.setStatus("FAILED");
            return false;
        }

        //Handle Response
        String responseContent = result.get(APIClient.REPLY_RESULT_CONTENT).toString();
        JsonObject model =  gson.fromJson(responseContent, JsonObject.class);
        if (model == null){
            executionRequestRefContext.appendLog("Model Storage execute error, Model is null and status code: " + statusCode);
            this.setStatus("FAILED");
            return false;
        }
        if (model.get("code").getAsInt() == 200) {
            this.setStatus("SUCCESS");
            return true;
        }
        this.setStatus("FAILED");
        return false;
    }

    @Override
    public RefExecutionState status() {
//        JsonObject model = ModelFactoryAPI.ModelGet(this.getUser(), this.modelId);
//        String status = model.get("status").toString();
        return this.transformStatus(this.getStatus());
    }

    @Override
    public boolean kill() {
        Map result = ModelFactoryAPI.ModelDelete(this.getUser(), this.modelId);
        int statusCode = Integer.parseInt(result.get(APIClient.REPLY_STATUS_CODE).toString());
        if (statusCode !=  HttpStatus.SC_OK) {
            this.getExecutionContext().appendLog("Kill ModelStorage Job Error, Get Status From " +
                    "DI Rest Error: " + statusCode);
            this.getExecutionContext().appendLog("Response: " + result.get(APIClient.REPLY_RESULT_CONTENT).toString());
            return false;
        }

        String responseContent = result.get(APIClient.REPLY_RESULT_CONTENT).toString();
        JsonObject model =  gson.fromJson(responseContent, JsonObject.class);
        if (model == null) {
            return false;
        }
        return true;
    }

    @Override
    public String log() {
        return null;
    }

    @Override
    public RefExecutionState transformStatus(String status) {
        if ("SUCCESS".equals(status)){
            return RefExecutionState.Success;
        }else if("FAILED".equals(status)){
            return RefExecutionState.Failed;
        }else{
            return RefExecutionState.Running;
        }
    }

}
