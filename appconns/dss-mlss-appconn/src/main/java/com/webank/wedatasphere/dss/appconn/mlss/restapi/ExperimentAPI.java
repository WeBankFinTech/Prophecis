package com.webank.wedatasphere.dss.appconn.mlss.restapi;

import com.google.gson.Gson;
import com.google.gson.JsonObject;
import com.webank.wedatasphere.dss.appconn.mlss.utils.ApiClient;
import com.webank.wedatasphere.dss.appconn.mlss.utils.MLSSConfig;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.HashMap;
import java.util.Map;

/**
 * Create by v_wbyxzhong
 * Date Time 2021/3/16 18:49
 */
public class ExperimentAPI {
    private static final Gson gson = new Gson();
    private static final Logger logger = LoggerFactory.getLogger(ApiClient.class);

    public static JsonObject postExperiment(String user, String expName, String expDesc) {
        ApiClient apiClient = new ApiClient(MLSSConfig.BASE_URL, user);
        Map<String, String> paramsMap = new HashMap<>();
        paramsMap.put("exp_name", expName);
        paramsMap.put("exp_desc", expDesc);
        apiClient.addClientFormParams("body", paramsMap);
        String reply = apiClient.doPost("/di/v1/experiment");
        return gson.fromJson(reply, JsonObject.class);
    }

    public static JsonObject copyExperiment(String user, String expID, String createType) {
        ApiClient apiClient = new ApiClient(MLSSConfig.BASE_URL, user);
        apiClient.addClientQueryParams("create_type", createType);
        String reply = apiClient.doGet("/di/v1/experiment/"+expID+"/copy");
        return gson.fromJson(reply, JsonObject.class);
    }

    public static JsonObject putExperiment(String user, String expName, String expDesc) {
        ApiClient apiClient = new ApiClient(MLSSConfig.BASE_URL, user);
        Map<String, String> paramsMap = new HashMap<>();
        paramsMap.put("exp_name", expName);
        paramsMap.put("exp_desc", expDesc);
        apiClient.addClientFormParams("body", paramsMap);
        String reply = apiClient.doPut("/di/v1/experiment/");
        return gson.fromJson(reply, JsonObject.class);
    }

    public static JsonObject deleteExperiment(String user, String expId) {
        ApiClient apiClient = new ApiClient(MLSSConfig.BASE_URL, user);
        String reply = apiClient.doDelete("/di/v1/experiment/" + expId);
        return gson.fromJson(reply, JsonObject.class);
    }

    public static JsonObject exportExperimentDSS(String user, String expId) {
        ApiClient apiClient = new ApiClient(MLSSConfig.BASE_URL, user);
        String reply = apiClient.doGet("/di/v1/experiment/" + expId + "/exportdss");
        return gson.fromJson(reply, JsonObject.class);
    }

    public static JsonObject importExperimentDSS(String user, String resourceId, String version, String desc) {
        ApiClient apiClient = new ApiClient(MLSSConfig.BASE_URL, user);
        apiClient.addClientQueryParams("resourceId", resourceId);
        apiClient.addClientQueryParams("version", version);
        apiClient.addClientQueryParams("desc", desc);
        String reply = apiClient.doPost("/di/v1/experiment/importdss");
        return gson.fromJson(reply, JsonObject.class);
    }
}
