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
public class ExperimentRunAPI {
    private static final Gson gson = new Gson();
    private static final Logger logger = LoggerFactory.getLogger(ApiClient.class);

    public static JsonObject post(String user, String expId, String expExecType, String flowJson) {
        ApiClient apiClient = new ApiClient(MLSSConfig.BASE_URL, user);
        Map<String, String> paramsMap = new HashMap<>();
        paramsMap.put("exp_exec_type", expExecType);
        paramsMap.put("flow_json", flowJson);
        apiClient.addClientFormParams("body", paramsMap);
        String reply = apiClient.doPost("/di/v1/experimentRun/" + expId);
        //TODO: Add status control
        return gson.fromJson(reply, JsonObject.class);
    }

    public static JsonObject getStatus(String user, String expId) {
        ApiClient apiClient = new ApiClient(MLSSConfig.BASE_URL, user);
        String reply = apiClient.doGet("/di/v1/experimentRun/" + expId + "/status");
        return gson.fromJson(reply, JsonObject.class);
    }

    public static JsonObject kill(String user, String expRunId) {
        ApiClient apiClient = new ApiClient(MLSSConfig.BASE_URL, user);
        String reply = apiClient.doGet("/di/v1/experimentRun/" + expRunId + "/kill");
        return gson.fromJson(reply, JsonObject.class);
    }

    public static JsonObject getLogs(String user, String expRunId, Long fromLine, Long size) {
        ApiClient apiClient = new ApiClient(MLSSConfig.BASE_URL, user);
        apiClient.addClientQueryParams("from_line", String.valueOf(fromLine));
        apiClient.addClientQueryParams("size", String.valueOf(size));
        String reply = apiClient.doGet("/di/v1/experimentRun/" + expRunId + "/log");
        return gson.fromJson(reply, JsonObject.class);
    }
}
