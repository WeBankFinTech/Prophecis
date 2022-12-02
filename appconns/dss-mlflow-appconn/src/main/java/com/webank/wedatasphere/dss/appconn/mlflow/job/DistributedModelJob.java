package com.webank.wedatasphere.dss.appconn.mlflow.job;

import com.google.gson.JsonObject;
import com.google.gson.internal.LinkedTreeMap;
import com.webank.wedatasphere.dss.appconn.mlflow.restapi.DistributedModelAPI;
import com.webank.wedatasphere.dss.appconn.mlflow.utils.APIClient;
import com.webank.wedatasphere.dss.appconn.mlflow.utils.JobParser;
import com.webank.wedatasphere.dss.flow.execution.entrance.conf.FlowExecutionEntranceConfiguration;
import com.webank.wedatasphere.dss.standard.app.development.listener.common.RefExecutionState;
import com.webank.wedatasphere.dss.standard.app.development.listener.core.ExecutionRequestRefContext;
import org.apache.linkis.common.log.LogUtils;
import org.apache.http.HttpStatus;

import java.io.File;
import java.util.Map;

import static org.apache.linkis.cs.common.utils.CSCommonUtils.gson;

public class DistributedModelJob extends MLFlowJob {

    public static String DI_JOB_STAGE_SUBMIT = "Submit";
    public static String DI_JOB_STAGE_INIT = "Init";
    public static String DI_JOB_STAGE_QUEUE = "Queue";
    public static String DI_JOB_STAGE_PENDING = "Pending";
    public static String DI_JOB_STAGE_RUNNING = "Running";
    public static String DI_JOB_STAGE_COMPLETE = "Complete";
    private String modelId;

    DistributedModelJob(String user) {
        this.setJobType(JobManager.JOB_TYPE_DI);
        this.setUser(user);
    }

    @Override
    public boolean submit(Map jobContent, ExecutionRequestRefContext executionRequestRefContext) {
        //1. Params Init
        LinkedTreeMap jobContentMap = (LinkedTreeMap) jobContent.get("ManiFest");
//        ClientService.clientMap.put("executionRequestRefContext", executionRequestRefContext);
        this.setExecutionContext(executionRequestRefContext);

        //TODO:GET taskID from dss
        //2.1 Get Task Id
        Long taskID = (Long) executionRequestRefContext.getRuntimeMap().get(FlowExecutionEntranceConfiguration.JOB_ID());

        //2. Prepare params for post model function
        //2.1 set dss task id
        jobContentMap.put("dss_task_id", taskID);
        //2.2 transform String to File
        File manifest = null;
        String manifestFileName = null;
        String user = this.getUser();
        JsonObject jobJsonObj = gson.toJsonTree(jobContentMap).getAsJsonObject();
        String jobContentStr = jobJsonObj.toString();
        try {
            manifestFileName = JobParser.TransformManifest(jobContentStr);
            manifest = new File(manifestFileName);
        } catch (Exception e) {
            executionRequestRefContext.appendLog(e.toString());
            e.printStackTrace();
            JobParser.clearFile(manifestFileName);
        }

        //3. Execute post model job
        //3.1 PostModel
        executionRequestRefContext.appendLog(LogUtils.generateInfo("Start execute mlflow job: send request"));
        Map result = DistributedModelAPI.postModel(user, manifest, null);
        int statusCode = Integer.parseInt(result.get(APIClient.REPLY_STATUS_CODE).toString());
        if (!(statusCode == HttpStatus.SC_OK || statusCode == HttpStatus.SC_CREATED)) {
            executionRequestRefContext.appendLog(LogUtils.generateERROR("Get Status From DI Rest Error: " + statusCode));
            executionRequestRefContext.appendLog(LogUtils.generateERROR("Response: " +
                    result.get(APIClient.REPLY_RESULT_CONTENT).toString()));
            JobParser.clearFile(manifestFileName);
            return false;
        }

        //3.2 Retry
        String responseContent = result.get(APIClient.REPLY_RESULT_CONTENT).toString();
        JsonObject jsonObject = gson.fromJson(responseContent, JsonObject.class);
        this.modelId = jsonObject.get("model_id").getAsString();
        int maxRetryCount = 5;
        if (modelId.length() <= 0) {
            return this.retryExecuteJob(maxRetryCount, user, manifest, manifestFileName);
        }
        JobParser.clearFile(manifestFileName);
        executionRequestRefContext.appendLog(LogUtils.generateInfo("Start execute mlflow job: send request successful"));
        return true;
    }

    private boolean retryExecuteJob(int maxRetryCount, String user,
                                    File manifest, String manifestFileName) {
        int retryCount = 0;
        while (modelId.length() <= 0) {
            String errorMsg = "ERROR: Response Model is NULL, Retry...";
            this.getExecutionContext().appendLog(errorMsg);

            Map result = DistributedModelAPI.postModel(user, manifest, null);
            int statusCode = Integer.parseInt(result.get(APIClient.REPLY_STATUS_CODE).toString());
            if (statusCode != HttpStatus.SC_OK) {
                this.getExecutionContext().appendLog("Get Status From DI Rest Error: " + statusCode);
                this.getExecutionContext().appendLog("Response: " + result.get(APIClient.REPLY_RESULT_CONTENT).toString());
                JobParser.clearFile(manifestFileName);
                return false;
            }

            String responseContent = result.get(APIClient.REPLY_RESULT_CONTENT).toString();
            JsonObject jsonObject = gson.fromJson(responseContent, JsonObject.class);
            this.modelId = jsonObject.get("model_id").getAsString();
            try {
                Thread.sleep(5000);
            } catch (InterruptedException e) {
                e.printStackTrace();
                this.getExecutionContext().appendLog("Sleep Exception" + e.toString());
            }
            retryCount = retryCount + 1;
            if (retryCount >= maxRetryCount) {
                this.getExecutionContext().appendLog("重试超过5次，任务创建失败。");
                JobParser.clearFile(manifestFileName);
                return false;
            }
        }
        JobParser.clearFile(manifestFileName);
        return true;
    }

    @Override
    public RefExecutionState status() {
        Map result = DistributedModelAPI.getModel(this.getUser(), this.modelId);

        //Status
        int statusCode = Integer.parseInt(result.get(APIClient.REPLY_STATUS_CODE).toString());
        if (statusCode != HttpStatus.SC_OK) {
            this.getExecutionContext().appendLog("Get Status From DI Rest Error: " + statusCode);
            this.getExecutionContext().appendLog("Get Status From DI Rest Error, Training ID: " + this.modelId);
            this.getExecutionContext().appendLog("Response: " + result.get(APIClient.REPLY_RESULT_CONTENT).toString());
            return RefExecutionState.Failed;
        }

        //Response
        String responseContent = result.get(APIClient.REPLY_RESULT_CONTENT).toString();
        JsonObject model = gson.fromJson(responseContent, JsonObject.class);
        if (null == model) {
            return this.transformStatus("");
        }
        String status = model.get("training").getAsJsonObject().get("training_status").
                getAsJsonObject().get("status").getAsString();
        return this.transformStatus(status);
    }

    @Override
    public boolean kill() {
        Map result = DistributedModelAPI.stopModel(this.getUser(), this.modelId);

        //Status
        int statusCode = Integer.parseInt(result.get(APIClient.REPLY_STATUS_CODE).toString());
        if (statusCode != HttpStatus.SC_OK) {
            this.getExecutionContext().appendLog("Kill DI Job Error, from  rest api status code: " + statusCode);
            this.getExecutionContext().appendLog("Response: " + result.get(APIClient.REPLY_RESULT_CONTENT).toString());
            return false;
        }

        //Response
        String responseContent = result.get(APIClient.REPLY_RESULT_CONTENT).toString();
        JsonObject jsonObject = gson.fromJson(responseContent, JsonObject.class);
        if (jsonObject != null && !"".equals(jsonObject.get("model_id").getAsString())) {
            return true;
        }
        return false;
    }

    @Override
    public String log() {
        if ("".equals(this.modelId)) {
            return null;
        }
        Map result = DistributedModelAPI.getLogs(this.getUser(), this.modelId, this.getSize(), this.getFrom());

        //Status
        int statusCode = Integer.parseInt(result.get(APIClient.REPLY_STATUS_CODE).toString());
        if (statusCode != HttpStatus.SC_OK) {
            this.getExecutionContext().appendLog("Get Status From DI Rest Error: " + statusCode);
            this.getExecutionContext().appendLog("Get Status From DI Rest Error, Training ID: " + this.modelId);
            this.getExecutionContext().appendLog("Response: " + result.get(APIClient.REPLY_RESULT_CONTENT).toString());
            return null;
        }

        //Response
        String responseContent = result.get(APIClient.REPLY_RESULT_CONTENT).toString();
        JsonObject model = gson.fromJson(responseContent, JsonObject.class);
        if (null == model) {
            return null;
        }
//TODO
//        int size = responseContent.split("rn|r|n").length;
//        this.setFrom(this.getFrom()+size);
//        return responseContent;
//TODO


//
//        int statusCode = Integer.parseInt(result.get(ApiClient.REPLY_STATUS_CODE).toString());
//        if (statusCode !=  HttpStatus.SC_OK) {
//            this.getDssExecutionRequestRefContext().appendLog("Get Status From DI Rest Error: " + statusCode);
//            this.getDssExecutionRequestRefContext().appendLog("Get Status From DI Rest Error, Training ID: " + this.modelId );
//            this.getDssExecutionRequestRefContext().appendLog("Response: " + result.get(ApiClient.REPLY_RESULT_CONTENT).toString());
//            return null;
//        }
//

//        String reply = result.get(ApiClient.REPLY_RESULT_CONTENT).toString();
//        if (reply.toUpperCase().contains("NULL") || reply.toUpperCase().contains("ERROR") || reply.toUpperCase().contains("FAILED"))
//            return null;
//
//
//
//
//        JsonObject jsonObj = gson.fromJson(reply, JsonObject.class);
//        if (jsonObj != null) {
//            String logsJson = gson.toJson(jsonObj.getAsJsonObject().get("logs"));
//            List<String> logList = gson.fromJson(logsJson, new TypeToken<List<String>>() {}.getType());
//            Boolean isLastLineFlag = jsonObj.getAsJsonObject().get("is_last_line").getAsString().toLowerCase().equals("TRUE");
//            this.setSize(logList.size());
//            this.setFrom(this.getFrom() + logList.size());
//            return new Log(logList, isLastLineFlag);
//        }
        return null;
    }

    @Override
    public RefExecutionState transformStatus(String status) {
        if ("PENDING".equals(status)) {
            return RefExecutionState.Accepted;
        } else if ("PROCESSING".equals(status)) {
            return RefExecutionState.Accepted;
        } else if ("FAILED".equals(status)) {
            return RefExecutionState.Failed;
        } else if ("QUEUED".equals(status)) {
            return RefExecutionState.Accepted;
        } else if ("RUNNING".equals(status)) {
            return RefExecutionState.Running;
        }
        //TODO Control Success
        return RefExecutionState.Success;
    }
}
