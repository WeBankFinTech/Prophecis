package com.webank.wedatasphere.dss.appconn.mlflow.utils;

import com.google.gson.Gson;
import com.webank.wedatasphere.dss.appconn.mlflow.utils.MLFlowConfig;
import org.apache.http.HttpEntity;
import org.apache.http.client.config.RequestConfig;
import org.apache.http.client.methods.CloseableHttpResponse;
import org.apache.http.client.methods.HttpEntityEnclosingRequestBase;
import org.apache.http.client.methods.HttpPost;
import org.apache.http.client.utils.URIBuilder;
import org.apache.http.entity.ContentType;
import org.apache.http.entity.StringEntity;
import org.apache.http.entity.mime.MultipartEntityBuilder;
import org.apache.http.entity.mime.content.FileBody;
import org.apache.http.entity.mime.content.StringBody;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClients;
import org.apache.http.util.EntityUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.File;
import java.io.IOException;
import java.net.URI;
import java.net.URISyntaxException;
import java.util.HashMap;
import java.util.Map;


public class APIClient {

    private static final RequestConfig requestConfig = RequestConfig.custom()
            .setSocketTimeout(30 * 1000)
            .setConnectTimeout(30 * 1000)
            .setConnectionRequestTimeout(30 * 1000)
            .build();
    private static final Gson gson = new Gson();
    private static final Logger logger = LoggerFactory.getLogger(APIClient.class);

    public final static String CONTENT_TYPE= "Content-Type";
    public final static String AUTH_HEADER_ACCEPT_ENCODING= "Accept-Encoding";
    public final static String AUTH_HEADER_CONTENT_ENCODING= "Content-Encoding";
    public final static String AUTH_HEADER_APIKEY = "MLSS-APPID";
    public final static String AUTH_HEADER_SIGNATURE = "MLSS-APPSignature";
    public final static String AUTH_HEADER_AUTH_TYPE = "MLSS-Auth-Type";
    public final static String AUTH_HEADER_USERID = "MLSS-UserID";
    public static final String AUTH_HEADER_TIMESTAMP = "MLSS-AppTimestamp";
    public final static String HTTP_PARAMS_TYPE_BODY = "BODY";
    public final static String HTTP_PARAMS_TYPE_FILE = "FILE";

    public final static String REPLY_STATUS_CODE = "statusCode";
    public final static String REPLY_RESULT_CONTENT = "resultContent";


    // Build get request
    private static HttpGetWithEntity buildGet(String requestUrl, Map<String,String> header, Map<String, String> params)  {
        String url = MLFlowConfig.BASE_URL + requestUrl;
        HttpGetWithEntity httpGetWithEntity = new HttpGetWithEntity(url);
        httpGetWithEntity.setConfig(requestConfig);

        //ADD Header
        for (Map.Entry<String, String> entry : header.entrySet()) {
            httpGetWithEntity.addHeader(entry.getKey(), entry.getValue());
        }

        //ADD Params
        if (params != null) {
            StringEntity bodyEntity = null;
            StringBuilder queryParamsStrBuilder = new StringBuilder();
            for (Map.Entry<String, String> entry : params.entrySet()) {
                if (entry.getKey().toUpperCase().equals(HTTP_PARAMS_TYPE_BODY)){
//                    String body = gson.toJson(entry.getValue());
                    String body = entry.getValue();
                    bodyEntity = new StringEntity(body, "UTF-8");
                    bodyEntity.setContentType("application/json");
                    httpGetWithEntity.addHeader(CONTENT_TYPE, "application/json");
                    httpGetWithEntity.setEntity(bodyEntity);
                    break;
                }else{
                    queryParamsStrBuilder.append(entry.getKey()).append("=").append(entry.getValue()).append("&");
                }
            }
            //Set Entity
            if (bodyEntity != null) {
                httpGetWithEntity.setEntity(bodyEntity);
            }else{
                String queryParamsStr = queryParamsStrBuilder.toString().substring(0,queryParamsStrBuilder.length() - 1);
                try {
                    httpGetWithEntity.setURI(new URIBuilder(httpGetWithEntity.getURI().toString() + "?"
                            + queryParamsStr).build());
                } catch (URISyntaxException e) {
                    e.printStackTrace();
                    logger.error("Http(GET) client build failed, set URI exception: " + e.getMessage());
                }
            }
        }

        return httpGetWithEntity;
    }

    // Build delete request
    private static HttpDeleteWithEntity buildDelete(String requestUrl, Map<String, String> header, Map<String, String> params)  {
        String url = MLFlowConfig.BASE_URL + requestUrl;
        HttpDeleteWithEntity httpDeleteWithEntity = new HttpDeleteWithEntity(url);
        httpDeleteWithEntity.setConfig(requestConfig);

        //ADD Header
        for (Map.Entry<String, String> entry : header.entrySet()) {
            httpDeleteWithEntity.addHeader(entry.getKey(), entry.getValue());
        }

        //ADD Params
        if (params != null) {
            StringEntity bodyEntity = null;
            StringBuilder queryParamsStrBuilder = new StringBuilder();
            for (Map.Entry<String, String> entry : params.entrySet()) {
                if (entry.getKey().toUpperCase().equals(HTTP_PARAMS_TYPE_BODY)){
//                    String body = gson.toJson(entry.getValue());
                    String body = entry.getValue();
                    bodyEntity = new StringEntity(body, "UTF-8");
                    bodyEntity.setContentType("application/json");
                    httpDeleteWithEntity.addHeader(CONTENT_TYPE, "application/json");
                    httpDeleteWithEntity.setEntity(bodyEntity);
                    break;
                }else{
                    queryParamsStrBuilder.append(entry.getKey()).append("=").append(entry.getValue()).append("&");
                }
            }

            //Set Entity
            if (bodyEntity != null) {
                httpDeleteWithEntity.setEntity(bodyEntity);
            }else{
                String queryParamsStr = queryParamsStrBuilder.toString().substring(0,queryParamsStrBuilder.length() - 1);
                try {
                    httpDeleteWithEntity.setURI(new URIBuilder(httpDeleteWithEntity.getURI().toString() + "?"
                            + queryParamsStr).build());
                } catch (URISyntaxException e) {
                    e.printStackTrace();
                    logger.error("Http(DELETE) client build failed, set URI exception: " + e.getMessage());
                }
            }
        }

        return httpDeleteWithEntity;
    }

    // Build post request
    private static HttpPost buildPost(String requestUrl, Map<String, String> header, Map<String, Object> params) {
        String url = MLFlowConfig.BASE_URL + requestUrl;
        Boolean bodyFlag = false;
        HttpPost httpPost = new HttpPost(url);
        httpPost.setConfig(requestConfig);
        //ADD Header
        for (Map.Entry<String, String> entry : header.entrySet()) {
            httpPost.addHeader(entry.getKey(), entry.getValue());
        }
        //ADD Params
        StringEntity bodyEntity = null;
        MultipartEntityBuilder multipartEntityBuilder = MultipartEntityBuilder.create();
        for (Map.Entry<String, Object> entry : params.entrySet()) {
            if (entry.getValue() ==null) {
                continue;
            }
            if (getType(entry.getValue()).toUpperCase().equals(HTTP_PARAMS_TYPE_FILE)) {
                FileBody fileBody = new FileBody((File) entry.getValue());
                multipartEntityBuilder.addPart(entry.getKey(), fileBody);
            }else if (entry.getKey().toUpperCase().equals(HTTP_PARAMS_TYPE_BODY)){
//                String body = gson.toJson(entry.getValue());
                String body = entry.getValue().toString();
                bodyEntity = new StringEntity(body, "UTF-8");
                bodyEntity.setContentType("application/json");
                httpPost.addHeader(CONTENT_TYPE, "application/json");
                break;
            }else {
                StringBody stringBody = new StringBody(entry.getValue().toString(), ContentType.DEFAULT_TEXT);
                multipartEntityBuilder.addPart(entry.getKey(), stringBody);
            }
        }

        if (null != bodyEntity) {
            httpPost.setEntity(bodyEntity);
        }else{
            HttpEntity multipartEntity = multipartEntityBuilder.build();
            httpPost.setEntity(multipartEntity);
        }

        return httpPost;
    }

    // Build put request
    private static HttpPutWithEntity buildPut(String requestUrl, Map<String,String> header, Map<String, String> params) {
        String url = MLFlowConfig.BASE_URL + requestUrl;
        HttpPutWithEntity httpPutWithEntity = new HttpPutWithEntity(url);
        httpPutWithEntity.setConfig(requestConfig);

        //ADD Header
        for (Map.Entry<String, String> entry : header.entrySet()) {
            httpPutWithEntity.addHeader(entry.getKey(), entry.getValue());
        }

        //ADD Params
        if (params!=null) {
            StringEntity bodyEntity = null;
            StringBuilder queryParamsStrBuilder = new StringBuilder();
            for (Map.Entry<String, String> entry : params.entrySet()) {
                if (entry.getKey().toUpperCase().equals(HTTP_PARAMS_TYPE_BODY)){
//                    String body = gson.toJson(entry.getValue());
                    String body = entry.getValue();
                    bodyEntity = new StringEntity(body, "UTF-8");
                    bodyEntity.setContentType("application/json");
                    httpPutWithEntity.addHeader(CONTENT_TYPE, "application/json");
                    httpPutWithEntity.setEntity(bodyEntity);
                    break;
                }else{
                    queryParamsStrBuilder.append(entry.getKey()).append("=").append(entry.getValue()).append("&");
                }
            }

            //Set Entity
            if (bodyEntity != null) {
                httpPutWithEntity.setEntity(bodyEntity);
            }else{
                String queryParamsStr = queryParamsStrBuilder.toString().substring(0,queryParamsStrBuilder.length() - 1);
                try {
                    httpPutWithEntity.setURI(new URIBuilder(httpPutWithEntity.getURI().toString() + "?"
                            + queryParamsStr).build());
                } catch (URISyntaxException e) {
                    e.printStackTrace();
                    logger.error("Http(PUT) client build failed, set URI exception: " + e.getMessage());
                }
            }
        }

        return httpPutWithEntity;
    }

    // Do get
    public static Map doGet(String requestUrl, Map<String,String> header, Map<String, String> params) {
        HttpGetWithEntity httpGet = buildGet(requestUrl, header, params);

        CloseableHttpClient httpClient = null;
        CloseableHttpResponse response = null;
        HttpEntity entity = null;
        String responseContent = null;
        Integer statusCode = null;
        try{
            httpClient = HttpClients.createDefault();
            response = httpClient.execute(httpGet);
            entity = response.getEntity();
            responseContent = EntityUtils.toString(entity, "UTF-8");
            statusCode = response.getStatusLine().getStatusCode();
        }catch(Exception e){
            e.printStackTrace();
            return BuildResult("-1",httpGet.toString() +" "+e.getMessage());
        }finally{
            try{
                if(response != null)
                    response.close();
                if(httpClient != null)
                    httpClient.close();
            }catch(IOException e){
                e.printStackTrace();
                logger.error("Do(GET) exception: " + e.getMessage());
            }
        }
        return BuildResult(statusCode+"",responseContent);
    }

    // Do post
    public static Map doPost(String requestUrl, Map<String, String> header, Map<String, Object> params) {
        HttpPost httpPost = buildPost(requestUrl, header, params) ;
        CloseableHttpClient httpClient = null;
        CloseableHttpResponse response = null;
        String responseContent = null;
        Integer statusCode = null;
        try{
            httpClient = HttpClients.createDefault();
            response = httpClient.execute(httpPost);
            HttpEntity entity = response.getEntity();
            responseContent = EntityUtils.toString(entity, "UTF-8");
            statusCode = response.getStatusLine().getStatusCode();
        }catch(Exception e){
            e.printStackTrace();
            return BuildResult("-1",httpPost.toString() +" "+e.getMessage());
        }finally{
            try{
                if(response != null)
                    response.close();
                if(httpClient != null)
                    httpClient.close();
            }catch(IOException e){
                e.printStackTrace();
                logger.error("Do(POST) exception: " + e.getMessage());
            }
        }
        return BuildResult(statusCode+"",responseContent);
    }

    // Do delete
    public static Map doDelete(String requestUrl, Map<String, String> header, Map<String, String> params) {
        HttpDeleteWithEntity httpDelete = buildDelete(requestUrl, header, params);
        CloseableHttpClient httpClient = null;
        CloseableHttpResponse response = null;
        String responseContent = null;
        Integer statusCode = 0;
        try{
            httpClient = HttpClients.createDefault();
            response = httpClient.execute(httpDelete);
            HttpEntity entity = response.getEntity();
            statusCode = response.getStatusLine().getStatusCode();
            responseContent = EntityUtils.toString(entity, "UTF-8");
        }catch(Exception e){
            e.printStackTrace();
            return BuildResult("-1",httpDelete.toString() +" "+e.getMessage());
        }finally{
            try{
                if(response != null)
                    response.close();
                if(httpClient != null)
                    httpClient.close();
            }catch(IOException e){
                e.printStackTrace();
                logger.error("Do(DELETE) exception: " + e.getMessage());
            }
        }

        return BuildResult(statusCode+"",responseContent);
    }

    // Do put
    public static Map doPut(String requestUrl, Map<String,String> header, Map<String, String> params) {
        HttpPutWithEntity httpPut = buildPut(requestUrl, header, params);
        CloseableHttpClient httpClient = null;
        CloseableHttpResponse response = null;
        String responseContent = null;
        Integer statusCode = null;
        try{
            httpClient = HttpClients.createDefault();
            response = httpClient.execute(httpPut);
            HttpEntity entity = response.getEntity();
            responseContent = EntityUtils.toString(entity, "UTF-8");
            statusCode = response.getStatusLine().getStatusCode();
        }catch(Exception e){
            e.printStackTrace();
            return BuildResult("-1",httpPut.toString() +" "+e.getMessage());
        }finally{
            try{
                if(response != null)
                    response.close();
                if(httpClient != null)
                    httpClient.close();
            }catch(IOException e){
                e.printStackTrace();
                logger.error("Do(PUT) exception: " + e.getMessage());
            }
        }
        return BuildResult(statusCode+"",responseContent);
    }


    // Get object's type
    private static String getType(Object object){
        String typeName = object.getClass().getName();
        int length = typeName.lastIndexOf(".");
        return typeName.substring(length + 1);
    }

    private static Map BuildResult(String statusCode, String reply){
        Map<String,String> result = new HashMap<>();
        result.put(REPLY_STATUS_CODE,statusCode);
        result.put(REPLY_RESULT_CONTENT, reply);
        return result;
    }
}


class HttpGetWithEntity extends HttpEntityEnclosingRequestBase {
    public final static String METHOD_NAME = "GET";

    public HttpGetWithEntity(String uri){
        this.setURI(URI.create(uri));
    }

    @Override
    public String getMethod() {
        return METHOD_NAME;
    }
}


class HttpDeleteWithEntity extends HttpEntityEnclosingRequestBase {
    public final static String METHOD_NAME = "DELETE";

    public HttpDeleteWithEntity(String uri){
        this.setURI(URI.create(uri));
    }

    @Override
    public String getMethod() {
        return METHOD_NAME;
    }
}

class HttpPutWithEntity extends HttpEntityEnclosingRequestBase {
    public final static String METHOD_NAME = "DELETE";

    public HttpPutWithEntity(String uri){
        this.setURI(URI.create(uri));
    }

    @Override
    public String getMethod() {
        return METHOD_NAME;
    }
}