package com.webank.wedatasphere.dss.appconn.mlflow.job;

import com.google.gson.JsonObject;
import com.google.gson.internal.LinkedTreeMap;
import com.webank.wedatasphere.dss.appconn.mlflow.restapi.ModelFactoryAPI;
import com.webank.wedatasphere.dss.appconn.mlflow.utils.APIClient;
import com.webank.wedatasphere.dss.appconn.mlflow.utils.MLFlowConfig;
import com.webank.wedatasphere.dss.standard.app.development.listener.common.RefExecutionState;
import com.webank.wedatasphere.dss.standard.app.development.listener.core.ExecutionRequestRefContext;
import org.apache.http.HttpStatus;

import java.util.HashMap;
import java.util.Map;

import static org.apache.linkis.cs.common.utils.CSCommonUtils.gson;

public class ImageBuildJob extends MLFlowJob {

    private String imageID;

    public ImageBuildJob(String user){
        this.setUser(user);
        this.setJobType(JobManager.JOB_TYPE_IMAGE_BUILD);
    }

    @Override
    public boolean submit(Map jobContent, ExecutionRequestRefContext ExecutionRequestRefContext) {
        //1. Params Prepare
        this.setExecutionContext(ExecutionRequestRefContext);
        LinkedTreeMap content = (LinkedTreeMap) jobContent.get("content");
        JobManager.paramsTransfer(content);
        JsonObject jobJsonObj = buildImageEntity(content);
        if (jobJsonObj == null) {
            ExecutionRequestRefContext.appendLog("Image Build execute error in param init: " +
                    "json object is null");
            return false;
        }
        String jsonBody = jobJsonObj.toString();


        //2. Execute Job
        Map result = ModelFactoryAPI.ImagePost(this.getUser(), jsonBody);
        int statusCode = Integer.parseInt(result.get(APIClient.REPLY_STATUS_CODE).toString());
        if (statusCode !=  HttpStatus.SC_OK) {
            ExecutionRequestRefContext.appendLog("ImageBuild job Get Status From MF Rest Error: " + statusCode);
            ExecutionRequestRefContext.appendLog("Response: " + result.get(APIClient.REPLY_RESULT_CONTENT).toString());
            return false;
        }

        //3. Handle Response
        String responseContent = result.get(APIClient.REPLY_RESULT_CONTENT).toString();
        JsonObject image =  gson.fromJson(responseContent, JsonObject.class);
        if (image == null) {
            return false;
        }
        this.imageID = image.get("result").getAsJsonObject().get("id").getAsString();
        return true;
    }

    @Override
    public RefExecutionState status() {
        Map result = ModelFactoryAPI.ImageGet(this.getUser(), this.imageID);
        int statusCode = Integer.parseInt(result.get(APIClient.REPLY_STATUS_CODE).toString());
        if (statusCode !=  HttpStatus.SC_OK) {
            this.getExecutionContext().appendLog("ImageBuild job Get Status From MF Rest Error: " + statusCode);
            this.getExecutionContext().appendLog("Response: " + result.get(APIClient.REPLY_RESULT_CONTENT).toString());
            return RefExecutionState.Failed;
        }

        String responseContent = result.get(APIClient.REPLY_RESULT_CONTENT).toString();
        JsonObject image =  gson.fromJson(responseContent, JsonObject.class);
        String status = image.get("result").getAsJsonObject().get("status").getAsString();
        return this.transformStatus(status);
    }

    @Override
    public boolean kill() {
        Map result = ModelFactoryAPI.ImageDelete(this.getUser(), this.imageID);
        int statusCode = Integer.parseInt(result.get(APIClient.REPLY_STATUS_CODE).toString());
        if (statusCode !=  HttpStatus.SC_OK) {
            this.getExecutionContext().appendLog("ImageBuild job Kill From MF Rest Error: " + statusCode);
            this.getExecutionContext().appendLog("Response: " + result.get(APIClient.REPLY_RESULT_CONTENT).toString());
            return false;
        }
        String responseContent = result.get(APIClient.REPLY_RESULT_CONTENT).toString();
        JsonObject res =  gson.fromJson(responseContent, JsonObject.class);
        if (res == null){
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
        if ("CREATING".equals(status.toUpperCase())){
            return RefExecutionState.Accepted;
        }else if("BUILDING".equals(status.toUpperCase())){
            return RefExecutionState.Running;
        }else if("PUSHING".equals(status.toUpperCase())) {
            return RefExecutionState.Running;
        }else if("COMPLETE".equals(status.toUpperCase())){
            return RefExecutionState.Success;
        }else if("FAILED".equals(status.toUpperCase())){
            return RefExecutionState.Failed;
        }else{
            return RefExecutionState.Failed;
        }
    }

    public static Map<String,String> getDefaultHeaders(String user){
        Map<String,String> headers = new HashMap<String,String>();
        headers.put(APIClient.AUTH_HEADER_APIKEY, MLFlowConfig.APP_KEY);
        headers.put(APIClient.AUTH_HEADER_SIGNATURE, MLFlowConfig.APP_SIGN);
        headers.put(APIClient.AUTH_HEADER_AUTH_TYPE, MLFlowConfig.AUTH_TYPE);
        headers.put(APIClient.AUTH_HEADER_USERID, user);
        return headers;
    }

    public JsonObject buildImageEntity(Map content){
        //1. Params Init
        String imageName = content.get("image_name").toString();
        String userName = content.get("user_name").toString();
        int groupId = Integer.parseInt(content.get("group_id").toString());
        String remarks = content.get("remarks").toString();

        //2. Set Model Version ID
        //2.1 Get model version id from list
        int modelVersionId = -1;
        if (content.containsKey("model_version_id")) {
            modelVersionId = Integer.parseInt(content.get("model_version_id").toString());
        }
        //2.2 Get model version id from model factory API
        else{
            String modelName = content.get("model_name").toString();
            String version = content.get("model_version").toString();
            Map result = ModelFactoryAPI.ModelVersionGet(this.getUser(), modelName,groupId, version);

            //Get model version id from MF Server
            int statusCode = Integer.parseInt(result.get(APIClient.REPLY_STATUS_CODE).toString());
            if (statusCode !=  HttpStatus.SC_OK) {
                this.getExecutionContext().appendLog("ImageBuild job get model version from MF rest  server error: "
                        + statusCode);
                this.getExecutionContext().appendLog("Response: "
                        + result.get(APIClient.REPLY_RESULT_CONTENT).toString());
                return null;
            }
            String responseContent = result.get(APIClient.REPLY_RESULT_CONTENT).toString();
            JsonObject res =  gson.fromJson(responseContent, JsonObject.class);
            if (res == null) {
                this.getExecutionContext().appendLog("ImageBuild job get model version from MF rest server error," +
                        " response Content is null and status code is: "
                        + statusCode);
                return  null;
            }
            modelVersionId = Integer.parseInt(res.get("result").getAsJsonObject().get("id").toString());
        }

        Map<String,Object> image = new HashMap<String,Object>();
        image.put("image_name",imageName);
        image.put("user_name", userName);
        image.put("group_id",groupId);
        image.put("model_version_id",modelVersionId);
        image.put("remarks",remarks);

        JsonObject jobJsonObj = gson.toJsonTree(image).getAsJsonObject();
        return jobJsonObj;
    }
}
