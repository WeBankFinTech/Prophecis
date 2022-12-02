package com.webank.wedatasphere.dss.appconn.mlflow.job;

import com.google.gson.JsonObject;
import com.google.gson.internal.LinkedTreeMap;
import com.webank.wedatasphere.dss.appconn.mlflow.restapi.ModelFactoryAPI;
import com.webank.wedatasphere.dss.appconn.mlflow.utils.APIClient;
import com.webank.wedatasphere.dss.appconn.mlflow.utils.MLFlowNodeUtils;
import com.webank.wedatasphere.dss.standard.app.development.listener.common.RefExecutionState;
import com.webank.wedatasphere.dss.standard.app.development.listener.core.ExecutionRequestRefContext;
import org.apache.http.HttpStatus;

import java.util.HashMap;
import java.util.Map;

import static org.apache.linkis.cs.common.utils.CSCommonUtils.gson;

public class ModelDeployJob extends MLFlowJob {

    private String serviceId;
    private String namespace;
    private String serviceName;

    public ModelDeployJob(String user) {
        this.setUser(user);
        this.setJobType(JobManager.JOB_TYPE_MODEL_DEPLOY);
    }

    @Override
    public boolean submit(Map jobContent, ExecutionRequestRefContext executionRequestRefContext) {
        //1. Params init
        LinkedTreeMap content = (LinkedTreeMap) jobContent.get("content");
        JobManager.paramsTransfer(content);
        JsonObject jobJsonObj = gson.toJsonTree(content).getAsJsonObject();
        String jsonBody = jobJsonObj.toString();
        this.setExecutionContext(executionRequestRefContext);

        //2. Execute job
        String user = this.getUser();
        Map result = ModelFactoryAPI.ServicePost(user, jsonBody);
        int statusCode = Integer.parseInt(result.get(APIClient.REPLY_STATUS_CODE).toString());
        if (statusCode != HttpStatus.SC_OK) {
            executionRequestRefContext.appendLog("Get Status From MF Rest Error: " + statusCode);
            executionRequestRefContext.appendLog("Response: " + result.get(APIClient.REPLY_RESULT_CONTENT).toString());
            return false;
        }
        String responseContent = result.get(APIClient.REPLY_RESULT_CONTENT).toString();
        JsonObject service = gson.fromJson(responseContent, JsonObject.class);
        if (service == null) {
            executionRequestRefContext.appendLog("Service Node execute error, service is null and statusCode: " + statusCode);
            return false;
        }

        //3. Handle Response
        this.serviceId = service.get("result").getAsJsonObject().get("id").getAsString();
        this.serviceName = service.get("result").getAsJsonObject().get("service_name").getAsString();
        this.namespace = service.get("result").getAsJsonObject().get("namespace").getAsString();
        return true;
    }

    @Override
    public RefExecutionState status() {
        // Execute Job
        Map result = ModelFactoryAPI.ServiceGet(this.getUser(), this.serviceId);
        int statusCode = Integer.parseInt(result.get(APIClient.REPLY_STATUS_CODE).toString());
        if (statusCode != HttpStatus.SC_OK) {
            this.getExecutionContext().appendLog("Get Status From MF Rest Error: " + statusCode);
            this.getExecutionContext().appendLog("Response: " + result.get(APIClient.REPLY_RESULT_CONTENT).toString());
            return RefExecutionState.Failed;
        }
        String responseContent = result.get(APIClient.REPLY_RESULT_CONTENT).toString();
        JsonObject service = gson.fromJson(responseContent, JsonObject.class);
        if (service == null) {
            return RefExecutionState.Failed;
        }

        //Handle Response
        //TODO ERROR Control
        String status = service.get("result").getAsJsonObject().get("status").getAsString();
        return this.transformStatus(status);
    }

    @Override
    public boolean kill() {
        Map result = ModelFactoryAPI.ServiceStop(this.getUser(), this.namespace, this.serviceName, this.serviceId);
        int statusCode = Integer.parseInt(result.get(APIClient.REPLY_STATUS_CODE).toString());
        if (statusCode != HttpStatus.SC_OK) {
            this.getExecutionContext().appendLog("Kill ModelDeploy Job From MF Rest Error: " + statusCode);
            this.getExecutionContext().appendLog("Response: " + result.get(APIClient.REPLY_RESULT_CONTENT).toString());
            return false;
        }
        String responseContent = result.get(APIClient.REPLY_RESULT_CONTENT).toString();
        JsonObject res = gson.fromJson(responseContent, JsonObject.class);
        if (res == null) {
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
        if ("Creating".equals(status)) {
            return RefExecutionState.Accepted;
        } else if ("Available".equals(status)) {
            return RefExecutionState.Success;
        } else if ("Stop".equals(status)) {
            return RefExecutionState.Failed;
        } else if ("FAILED".equals(status)) {
            return RefExecutionState.Failed;
        } else {
            return RefExecutionState.Failed;
        }
    }

//    public JsonObject buildServiceEntity(Map content){
//        //1. Params Init
//        String serviceName = content.get("service_name").toString();
//        String type = content.get("type").toString();
//        String namespace = content.get("namespace").toString();
//        String gpu = content.get("gpu").toString();
//        int cpu = Integer.parseInt(content.get("cpu").toString());
//        int memory = Integer.parseInt(content.get("memory").toString());
//        int groupId = Integer.parseInt(content.get("group_id").toString());
//        String remark = content.get("remark").toString();
//
//
//        //2. Set Model Version ID
//        //2.1 Get model version id from list
//        int modelVersionId = -1;
//        if (content.containsKey("model_version_id")) {
//            modelVersionId = Integer.parseInt(content.get("model_version_id").toString());
//        }
//        //2.2 Get model version id from model factory API
//        else{
//            String modelName = content.get("model_name").toString();
//            String modelVersion = content.get("model_version").toString();
//            Map result = ModelFactoryAPI.ModelVersionGet(this.getUser(), modelName, groupId,modelVersion);
//
//            //Get model version id from MF Server
//            int statusCode = Integer.parseInt(result.get(ApiClient.REPLY_STATUS_CODE).toString());
//            if (statusCode !=  HttpStatus.SC_OK) {
//                this.getDssExecutionRequestRefContext().appendLog("ImageBuild job get model version from MF rest  server error: "
//                        + statusCode);
//                this.getDssExecutionRequestRefContext().appendLog("Response: "
//                        + result.get(ApiClient.REPLY_RESULT_CONTENT).toString());
//                return null;
//            }
//            String responseContent = result.get(ApiClient.REPLY_RESULT_CONTENT).toString();
//            JsonObject res =  gson.fromJson(responseContent, JsonObject.class);
//            if (res == null) {
//                this.getDssExecutionRequestRefContext().appendLog("ImageBuild job get model version from MF rest server error," +
//                        " response Content is null and status code is: "
//                        + statusCode);
//                return  null;
//            }
//            modelVersionId = Integer.parseInt(res.get("id").toString());
//        }
//
//
//        Map<String,Object> image = new HashMap<String,Object>();
//        image.put("service_name",serviceName);
//        image.put("type", type);
//        image.put("namespace", namespace);
//        image.put("gpu", gpu);
//        image.put("cpu", cpu);
//        image.put("memory", memory);
//        image.put("group_id", groupId);
//        image.put("remark",remark);
//
//        image.put("model_version_id",modelVersionId);
//
//
//        JsonObject jobJsonObj = gson.toJsonTree(image).getAsJsonObject();
//        return jobJsonObj;
//    }
}
