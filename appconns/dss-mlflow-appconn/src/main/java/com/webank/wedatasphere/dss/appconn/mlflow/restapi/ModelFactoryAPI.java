package com.webank.wedatasphere.dss.appconn.mlflow.restapi;

import com.google.gson.Gson;
import com.webank.wedatasphere.dss.appconn.mlflow.utils.APIClient;
import com.webank.wedatasphere.dss.appconn.mlflow.utils.MLFlowConfig;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.HashMap;
import java.util.Map;

public class ModelFactoryAPI {

    private static final Gson gson = new Gson();
    private static final Logger logger = LoggerFactory.getLogger(APIClient.class);


    public static Map ServicePost(String user, String jsonBody) {
        Map<String, String> headers = getDefaultHeaders(user);
        Map<String, Object> params = new HashMap<>();
        params.put("BODY",jsonBody);
        return APIClient.doPost("/mf/v1/service", headers, params);
    }

    public static Map ServiceGet(String user, String serviceId) {
        Map<String, String> headers = getDefaultHeaders(user);
        return APIClient.doGet("/mf/v1/service/" + serviceId, headers, null);
    }

    public static Map ServiceStop(String user, String namespace, String serviceName, String serviceId) {
        Map<String, String> headers = getDefaultHeaders(user);
        return APIClient.doGet("/mf/v1/serviceStop/" + namespace +
                "/" + serviceName + "/" + serviceId, headers, null);
    }

    public static Map ImagePost(String user, String jsonBody) {
        Map<String, String> headers = getDefaultHeaders(user);
        Map<String, Object> params = new HashMap<>();
        params.put("BODY",jsonBody);
        return APIClient.doPost("/mf/v1/image", headers, params);
    }

    public static Map ImageGet(String user, String imageId) {
        Map<String, String> headers = getDefaultHeaders(user);
        return APIClient.doGet("/mf/v1/image/" + imageId, headers, null);
    }

    public static Map ImageDelete(String user, String imageId) {
        Map<String, String> headers = getDefaultHeaders(user);
        return APIClient.doDelete("/mf/v1/image/" + imageId, headers, null);
    }

    public static Map ModelPost(String user, String jsonBody) {
        Map<String, String> headers = getDefaultHeaders(user);
        Map<String, Object> params = new HashMap<>();
        params.put("BODY",jsonBody);
        return APIClient.doPost("/mf/v1/model", headers, params);
    }

    public static Map ModelGet(String user, String modelId) {
        Map<String, String> headers = getDefaultHeaders(user);
        return APIClient.doGet("/mf/v1/model/"+ modelId, headers, null);
    }

    //TODO: Server is not implement
    public static Map ModelVersionGet(String user, String modelName, int groupId, String version) {
        Map<String, String> headers = getDefaultHeaders(user);
        return APIClient.doGet("/mf/v1/modelVersion/modelName/"+modelName+"/groupId/" +
                groupId +"/version/"+ version, headers, null);
    }

    public static Map ModelDelete(String user, String modelId) {
        Map<String, String> headers = getDefaultHeaders(user);
        return APIClient.doDelete("/mf/v1/models/" + modelId, headers, null);
    }

    private static Map<String,String> getDefaultHeaders(String user){
        Map<String,String> headers = new HashMap<>();
        headers.put(APIClient.AUTH_HEADER_APIKEY, MLFlowConfig.APP_KEY);
        headers.put(APIClient.AUTH_HEADER_TIMESTAMP, MLFlowConfig.TIMESTAMP);
        headers.put(APIClient.AUTH_HEADER_SIGNATURE, MLFlowConfig.APP_SIGN);
        headers.put(APIClient.AUTH_HEADER_AUTH_TYPE, MLFlowConfig.AUTH_TYPE);
        headers.put(APIClient.AUTH_HEADER_USERID, user);
        return headers;
    }

    public static Map ReportDelete(String user, String reportId) {
        Map<String, String> headers = getDefaultHeaders(user);
        return APIClient.doDelete("/mf/v1/report/" + reportId, headers, null);
    }

    public static Map ReportPost(String user, String jsonBody) {
        Map<String, String> headers = getDefaultHeaders(user);
        Map<String, Object> params = new HashMap<>();
        params.put("BODY",jsonBody);
        return APIClient.doPost("/mf/v1/report", headers, params);
    }

    public static Map ReportGet(String user, String reportId) {
        Map<String, String> headers = getDefaultHeaders(user);
        return APIClient.doGet("/mf/v1/report/" + reportId, headers, null);
    }

    public static Map PushEventGet(String user, String eventId) {
        Map<String, String> headers = getDefaultHeaders(user);
        return APIClient.doGet("/mf/v1/pushEvent/" + eventId, headers, null);
    }

    public static Map ReportPushPost(String user, String reportId, String jsonBody) {
        Map<String, String> headers = getDefaultHeaders(user);
        Map<String, Object> params = new HashMap<>();
        params.put("BODY",jsonBody);
        return APIClient.doPost("/mf/v1/report/push/"+reportId, headers, params);
    }

    public static Map ReportVersionPushPost(String user, String reportVersionId, String jsonBody) {
        Map<String, String> headers = getDefaultHeaders(user);
        Map<String, Object> params = new HashMap<>();
        params.put("BODY",jsonBody);
        return APIClient.doPost("/mf/v1/reportVersion/Push/"+reportVersionId, headers, params);
    }

    public static Map modelMonitorReportPost(String user, String jsonBody) {
        Map<String, String> headers = getDefaultHeaders(user);
        Map<String, Object> params = new HashMap<>();
        params.put("BODY",jsonBody);
        return APIClient.doPost("/mf/v1/model_monitor/modelMonitorTask", headers, params);
    }

    public static Map modelMonitorReportGet(String user,  Map<String, String> params) {
        Map<String, String> headers = getDefaultHeaders(user);
        return APIClient.doGet("/mf/v1/model_monitor/modelMonitorTask", headers, params);
    }



//gson.fromJson(reply, JsonObject.class)

}
