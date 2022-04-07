package com.webank.wedatasphere.dss.appconn.mlflow.restapi;

import com.google.gson.Gson;
import com.webank.wedatasphere.dss.appconn.mlflow.utils.APIClient;
import com.webank.wedatasphere.dss.appconn.mlflow.utils.MLFlowConfig;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.File;
import java.util.HashMap;
import java.util.Map;


public class DistributedModelAPI {

    private static final Gson gson = new Gson();
    private static final Logger logger = LoggerFactory.getLogger(APIClient.class);

    public static Map getModel(String user, String modelId) {
        Map<String,String>  headers = getDefaultHeaders(user);
        return APIClient.doGet("/di/v1/models/" + modelId, headers, null);
    }

    public static Map getLogs(String user, String modelId, Integer size, Integer from) {
        Map<String,String> headers = getDefaultHeaders(user);
        Map<String,String> params = new HashMap<>();
        params.put("size", String.valueOf(size));
        params.put("from", String.valueOf(from));
        //TODO FIX LOG
//        if (reply.toUpperCase().contains("NULL") || reply.toUpperCase().contains("ERROR") || reply.toUpperCase().contains("FAILED"))
//            return null;
        return APIClient.doGet("/di/v1/job/" + modelId + "/log", headers, params);
    }

    public static Map deleteModel(String user, String modelId) {
        Map<String,String> headers = getDefaultHeaders(user);
        return APIClient.doDelete("/di/v1/models/" + modelId,headers,null);
    }

    public static Map stopModel(String user, String modelId) {
        Map<String,String> headers = getDefaultHeaders(user);
        return APIClient.doGet("/di/v1/models/" + modelId + "/kill",headers,null);
    }

    public static Map postModel(String user, File manifestFile, File modelDefinitionFile) {
        Map<String,String>  headers = getDefaultHeaders(user);
        Map<String,Object> params = new HashMap<>();
        params.put("manifest", manifestFile);
        params.put("model_definition", modelDefinitionFile);
        return APIClient.doPost("/di/v1/models", headers, params);
//        return gson.fromJson(reply, JsonObject.class);
    }

    private static Map<String,String> getDefaultHeaders(String user){
        Map<String,String> headers = new HashMap<>();
        headers.put(APIClient.AUTH_HEADER_APIKEY, MLFlowConfig.APP_KEY);
        headers.put(APIClient.AUTH_HEADER_SIGNATURE, MLFlowConfig.APP_SIGN);
        headers.put(APIClient.AUTH_HEADER_AUTH_TYPE, MLFlowConfig.AUTH_TYPE);
        headers.put(APIClient.AUTH_HEADER_TIMESTAMP, MLFlowConfig.TIMESTAMP);
        headers.put(APIClient.AUTH_HEADER_USERID, user);
        headers.put(APIClient.AUTH_HEADER_CONTENT_ENCODING,"gzip");
        return headers;
    }
}
