package com.webank.wedatasphere.dss.appconn.mlflow.job;

import com.google.gson.JsonObject;
import com.google.gson.internal.LinkedTreeMap;
import com.webank.wedatasphere.dss.appconn.mlflow.restapi.ModelFactoryAPI;
import com.webank.wedatasphere.dss.appconn.mlflow.utils.APIClient;
import com.webank.wedatasphere.dss.standard.app.development.listener.common.RefExecutionState;
import com.webank.wedatasphere.dss.standard.app.development.listener.core.ExecutionRequestRefContext;
import org.apache.linkis.common.log.LogUtils;
import org.apache.http.HttpStatus;

import java.util.HashMap;
import java.util.Map;

import static org.apache.linkis.cs.common.utils.CSCommonUtils.gson;

public class ReportPushJob  extends MLFlowJob  {

    private String reportID;
    private String reportVersionID;
    private String eventID;

    public ReportPushJob(String user){
        this.setUser(user);
        this.setJobType(JobManager.JOB_TYPE_REPORT_PUSH);
    }


    @Override
    public boolean submit(Map jobContent, ExecutionRequestRefContext ExecutionRequestRefContext) {
        //1. Params Init
        ExecutionRequestRefContext.appendLog(LogUtils.generateInfo("Start report push job, Params init..."));
        LinkedTreeMap content = (LinkedTreeMap) jobContent.get("content");
        JobManager.paramsTransfer(content);
        JsonObject jobJsonObj = buildReportEntity(content);
        String jsonBody =  jobJsonObj.toString();
        String factoryName = content.get("factory_name").toString();
        this.setExecutionContext(ExecutionRequestRefContext);

        //2. Execute Job
        //2.1 Report Post
        ExecutionRequestRefContext.appendLog(LogUtils.generateInfo("Start submit report push job."));
        ExecutionRequestRefContext.appendLog(LogUtils.generateInfo("Post Report"));
        Map reportResult = ModelFactoryAPI.ReportPost(this.getUser(), jsonBody);
        int reportStatusCode = Integer.parseInt(reportResult.get(APIClient.REPLY_STATUS_CODE).toString());
        if (reportStatusCode !=  HttpStatus.SC_OK) {
            ExecutionRequestRefContext.appendLog(LogUtils.generateERROR("Report Push Job execute error, " +
                    "Get Status From MF Rest Error: " + reportStatusCode));
            ExecutionRequestRefContext.appendLog(LogUtils.generateERROR("Response: " +
                    reportResult.get(APIClient.REPLY_RESULT_CONTENT).toString()));
            return false;
        }
        String reportResponseContent = reportResult.get(APIClient.REPLY_RESULT_CONTENT).toString();
        JsonObject report = gson.fromJson(reportResponseContent, JsonObject.class);
        if (report == null || !report.has("result") ) {
            ExecutionRequestRefContext.appendLog(LogUtils.generateERROR("Report Push Job execute error, " +
                    "return report is null, status code: " + reportStatusCode));
            return false;
        }
        this.reportID = report.get("result").getAsJsonObject().get("report_id").getAsString();
        this.reportVersionID = report.get("result").getAsJsonObject().get("report_version_id").getAsString();
        ExecutionRequestRefContext.appendLog(LogUtils.generateInfo("Report push job submit successful."));

        //2.2 Report Push
        Map<String,String> pushEventParams = new HashMap<String,String>();
        pushEventParams.put("factory_name",factoryName);
        JsonObject pushEventJsonObj = gson.toJsonTree(pushEventParams).getAsJsonObject();
        String pushJsonBody =  pushEventJsonObj.toString();
        ExecutionRequestRefContext.appendLog(LogUtils.generateInfo("Start submit report push job."));
        Map pushResult = ModelFactoryAPI.ReportVersionPushPost(this.getUser(), reportVersionID
                , pushJsonBody);
        int statusCode = Integer.parseInt(pushResult.get(APIClient.REPLY_STATUS_CODE).toString());
        if (statusCode !=  HttpStatus.SC_OK) {
            ExecutionRequestRefContext.appendLog(LogUtils.generateERROR("Report Push Job execute error, " +
                    "Get Status From MF Rest Error: " + statusCode));
            ExecutionRequestRefContext.appendLog(LogUtils.generateERROR("Response: " +
                    pushResult.get(APIClient.REPLY_RESULT_CONTENT).toString()));
            return false;
        }
        String responseContent = pushResult.get(APIClient.REPLY_RESULT_CONTENT).toString();
        JsonObject pushEvent = gson.fromJson(responseContent, JsonObject.class);
        if (pushEvent == null) {
            return false;
        }

        //3. Handle response
        this.eventID = pushEvent.get("result").getAsJsonObject().get("id").getAsString();
        ExecutionRequestRefContext.appendLog(LogUtils.generateInfo("Report push job submit successful."));
        return true;
    }

    @Override
    public RefExecutionState status() {
        Map result = ModelFactoryAPI.PushEventGet(this.getUser(), this.eventID);
        int statusCode = Integer.parseInt(result.get(APIClient.REPLY_STATUS_CODE).toString());
        if (statusCode !=  HttpStatus.SC_OK) {
            this.getExecutionContext().appendLog(LogUtils.generateERROR("ReportPush job get status " +
                    "from MF rest error: " + statusCode));
            this.getExecutionContext().appendLog(LogUtils.generateERROR("Response: " +
                    result.get(APIClient.REPLY_RESULT_CONTENT).toString()));
            return RefExecutionState.Failed;
        }

        String responseContent = result.get(APIClient.REPLY_RESULT_CONTENT).toString();
        JsonObject report = gson.fromJson(responseContent, JsonObject.class);
        String status = report.get("result").getAsJsonObject().get("status").getAsString();
        return this.transformStatus(status);
    }

    @Override
    public boolean kill() {
//        Map result = ModelFactoryAPI.ReportDelete(this.getUser(), this.reportID);

//        int statusCode = Integer.parseInt(result.get(ApiClient.REPLY_STATUS_CODE).toString());
//        if (statusCode !=  HttpStatus.SC_OK) {
//            this.getDssExecutionRequestRefContext().appendLog("ReportPush job kill From MF Rest Error: " + statusCode);
//            this.getDssExecutionRequestRefContext().appendLog("Response: " + result.get(ApiClient.REPLY_RESULT_CONTENT).toString());
//            return false;
//        }
//
//        String responseContent = result.get(ApiClient.REPLY_RESULT_CONTENT).toString();
//        JsonObject res =  gson.fromJson(responseContent, JsonObject.class);
//        if (res == null){
//            return false;
//        }
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
        }else if("PUSHING".equals(status.toUpperCase())) {
            return RefExecutionState.Running;
        }else if("COMPLETE".equals(status.toUpperCase()) || "SUCCESS".equals(status.toUpperCase()) ){
            return RefExecutionState.Success;
        }else if("FAILED".equals(status.toUpperCase()) || "FAIL".equals(status.toUpperCase())){
            return RefExecutionState.Failed;
        }else{
            return RefExecutionState.Failed;
        }
    }



    private JsonObject buildReportEntity(LinkedTreeMap content){
        //1. Params Init
        int groupId = Integer.parseInt(content.get("group_id").toString());
        String modelName = content.get("model_name").toString();
        String reportName = content.get("report_name").toString();
        String factoryName = content.get("factory_name").toString();
        String modelVersion = "";

        //2. Set Model Version ID
        if (content.get("model_version").getClass() == String.class){
            modelVersion = content.get("model_version").toString();
        } else{
            modelVersion = ((Map) content.get("model_version")).get("name").toString();
        }
        String rootPath = content.get("path").toString();
        String childPath = content.get("file_path").toString();
        String fileName = childPath.substring(childPath.lastIndexOf("/")+1);

        Map<String,Object> report = new HashMap<String,Object>();
        report.put("report_name",reportName);
        report.put("group_id",groupId);
        report.put("model_name",modelName);
        report.put("model_version",modelVersion);
        report.put("root_path",rootPath);
        report.put("child_path",childPath);
        report.put("file_name",fileName);
        report.put("factory_name",factoryName);

        return gson.toJsonTree(report).getAsJsonObject();
    }
}
