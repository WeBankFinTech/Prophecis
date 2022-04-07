package com.webank.wedatasphere.dss.appconn.mlss.utils;

import com.google.gson.Gson;
import org.apache.http.HttpEntity;
import org.apache.http.client.config.RequestConfig;
import org.apache.http.client.methods.*;
import org.apache.http.client.utils.URIBuilder;
import org.apache.http.entity.StringEntity;
import org.apache.http.entity.mime.MultipartEntityBuilder;
import org.apache.http.entity.mime.content.FileBody;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClients;
import org.apache.http.util.EntityUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.File;
import java.io.IOException;
import java.net.URISyntaxException;
import java.util.HashMap;
import java.util.Map;

/**
 * Create by v_wbyxzhong
 * Date Time 2021/3/16 16:41
 */
public class ApiClient {
    private String user;
    private String baseUrl;
    private Map<String, String> clientHeaders;
    private Map<String, String> clientQueryParams;
    private Map<String, Object> clientFormParams;
    private static final RequestConfig requestConfig = RequestConfig.custom()
            .setSocketTimeout(60 * 1000)
            .setConnectTimeout(60 * 1000)
            .setConnectionRequestTimeout(6 * 1000)
            .build();
    private static final Gson gson = new Gson();
    private static final Logger logger = LoggerFactory.getLogger(ApiClient.class);

    public final static String AUTH_HEADER_APIKEY = "MLSS-APPID";
    public final static String AUTH_HEADER_TIMESTAMP = "MLSS-APPTimestamp";
    public final static String AUTH_HEADER_SIGNATURE = "MLSS-APPSignature";
    public final static String AUTH_HEADER_AUTH_TYPE = "MLSS-Auth-Type";
    public final static String AUTH_HEADER_USERID = "MLSS-UserID";
    public final static String CONTENT_TYPE= "Content-Type";
    public final static String AUTH_HEADER_ACCEPT_ENCODING= "Accept-Encoding";
    public final static String AUTH_HEADER_CONTENT_ENCODING= "Content-Encoding";

    public ApiClient() {
        this.setClientHeaders(new HashMap<>());
        this.setClientQueryParams(new HashMap<>());
        this.setClientFormParams(new HashMap<>());
    }

    public ApiClient(String baseUrl, String user) {
        this.setUser(user);
        this.setBaseUrl(baseUrl);
        this.setClientHeaders(new HashMap<>());
        this.setClientQueryParams(new HashMap<>());
        this.setClientFormParams(new HashMap<>());
    }

    public String getUser() { return user; }

    public void setUser(String user) { this.user = user; }

    public String getBaseUrl() {
        return baseUrl;
    }

    public void setBaseUrl(String baseUrl) {
        this.baseUrl = baseUrl;
    }

    public Map<String, String> getClientHeaders() {
        return clientHeaders;
    }

    public void setClientHeaders(Map<String, String> clientHeaders) {
        this.clientHeaders = clientHeaders;
    }

    public Map<String, String> getClientQueryParams() {
        return clientQueryParams;
    }

    public void setClientQueryParams(Map<String, String> clientQueryParams) {
        this.clientQueryParams = clientQueryParams;
    }

    public Map<String, Object> getClientFormParams() {
        return clientFormParams;
    }

    public void setClientFormParams(Map<String, Object> clientFormParams) {
        this.clientFormParams = clientFormParams;
    }

    /**
     * Add to headers
     * @param key
     * @param value
     */
    public void addClientHeaders(String key, String value) {
        if (value == null || value.length() <= 0) return;
        this.getClientHeaders().put(key, value);
    }

    /**
     * Add to query params
     * @param key
     * @param value
     */
    public void addClientQueryParams(String key, String value) {
        if (value == null || value.length() <= 0) return;
        this.getClientQueryParams().put(key, value);
    }

    /**
     * Add to form params
     * @param key
     * @param value
     */
    public void addClientFormParams(String key, Object value) {
        if (value == null) return;
        this.getClientFormParams().put(key, value);
    }

    // Build get request
    private HttpGet buildGet(String requestUrl) {
        this.addDefaultHeaders();
        String url = this.getBaseUrl() + requestUrl;
        HttpGet httpGet = new HttpGet(url);
        httpGet.setConfig(requestConfig);
        for (Map.Entry<String, String> entry : this.getClientHeaders().entrySet()) {
            if (entry.getValue() != null) {
                httpGet.addHeader(entry.getKey(), entry.getValue());
            }
        }
        StringBuilder queryParamsStrBuilder = new StringBuilder();
        for (Map.Entry<String, String> entry : this.getClientQueryParams().entrySet()) {
            queryParamsStrBuilder.append(entry.getKey()).append("=").append(entry.getValue()).append("&");
        }
        try {
            if (this.getClientQueryParams().size() > 0) {
                String queryParamsStr = queryParamsStrBuilder.toString().substring(0,queryParamsStrBuilder.length() - 1);
                httpGet.setURI(new URIBuilder(httpGet.getURI().toString() + "?" + queryParamsStr).build());
            }
        } catch (URISyntaxException e) {
            e.printStackTrace();
            logger.error("Http(GET) client build failed, set URI exception: " + e.getMessage());
        }
        return httpGet;
    }

    // Build delete request
    private HttpDelete buildDelete(String requestUrl) {
        this.addDefaultHeaders();
        String url = this.getBaseUrl() + requestUrl;
        HttpDelete httpDelete = new HttpDelete(url);
        httpDelete.setConfig(requestConfig);
        for (Map.Entry<String, String> entry : this.getClientHeaders().entrySet()) {
            if (entry.getValue() != null) {
                httpDelete.addHeader(entry.getKey(), entry.getValue());
            }
        }
        StringBuilder queryParamsStrBuilder = new StringBuilder();
        for (Map.Entry<String, String> entry : this.getClientQueryParams().entrySet()) {
            queryParamsStrBuilder.append(entry.getKey()).append("=").append(entry.getValue()).append("&");
        }
        try {
            if (this.getClientQueryParams().size() > 0) {
                String queryParamsStr = queryParamsStrBuilder.toString().substring(0,queryParamsStrBuilder.length() - 1);
                httpDelete.setURI(new URIBuilder(httpDelete.getURI().toString() + "?" + queryParamsStr).build());
            }
        } catch (URISyntaxException e) {
            e.printStackTrace();
            logger.error("Http(DELETE) client build failed, set URI exception: " + e.getMessage());
        }
        return httpDelete;
    }

    // Build post request
    private HttpPost buildPost(String requestUrl) {
        this.addDefaultHeaders();
        String url = this.getBaseUrl() + requestUrl;
        HttpPost httpPost = new HttpPost(url);
        httpPost.setConfig(requestConfig);
        Boolean bodyFlag = false;
        MultipartEntityBuilder multipartEntityBuilder = MultipartEntityBuilder.create();
        for (Map.Entry<String, String> entry : this.getClientHeaders().entrySet()) {
            httpPost.addHeader(entry.getKey(), entry.getValue());
        }
        StringBuilder queryParamsStrBuilder = new StringBuilder();
        for (Map.Entry<String, String> entry : this.getClientQueryParams().entrySet()) {
            queryParamsStrBuilder.append(entry.getKey()).append("=").append(entry.getValue()).append("&");
        }
        try {
            if (this.getClientQueryParams().size() > 0) {
                String queryParamsStr = queryParamsStrBuilder.toString().substring(0,queryParamsStrBuilder.length() - 1);
                httpPost.setURI(new URIBuilder(httpPost.getURI().toString() + "?" + queryParamsStr).build());
            }
        } catch (URISyntaxException e) {
            e.printStackTrace();
            logger.error("Http(POST) client build failed, set URI exception: " + e.getMessage());
        }
        for (Map.Entry<String, Object> entry : this.getClientFormParams().entrySet()) {
            if (entry.getValue() != null) {
                if (entry.getKey().toUpperCase().equals("BODY")) {
                    String body = gson.toJson(entry.getValue());
                    StringEntity stringEntity = new StringEntity(body, "UTF-8");
                    httpPost.setEntity(stringEntity);
                    bodyFlag = true;
                    break;
                }
                if (getType(entry.getValue()).toUpperCase().equals("FILE")) {
                    FileBody fileBody = new FileBody((File) entry.getValue());
                    multipartEntityBuilder.addPart(entry.getKey(), fileBody).build();
                }
            }
        }
        if (!bodyFlag) {
            HttpEntity multipartEntity = multipartEntityBuilder.build();
            httpPost.setEntity(multipartEntity);
        }
        return httpPost;
    }

    // Build put request
    private HttpPut buildPut(String requestUrl) {
        this.addDefaultHeaders();
        String url = this.getBaseUrl() + requestUrl;
        HttpPut httpPut = new HttpPut(url);
        httpPut.setConfig(requestConfig);
        for (Map.Entry<String, String> entry : this.getClientHeaders().entrySet()) {
            if (entry.getValue() != null) {
                httpPut.addHeader(entry.getKey(), entry.getValue());
            }
        }
        String body = gson.toJson(this.clientFormParams);
        httpPut.setEntity(new StringEntity(body, "UTF-8"));
        return httpPut;
    }

    // Do get
    public String doGet(String requestUrl) {
        HttpGet httpGet = this.buildGet(requestUrl);
        CloseableHttpClient httpClient = null;
        CloseableHttpResponse response = null;
        HttpEntity entity = null;
        String responseContent = null;
        try{
            httpClient = HttpClients.createDefault();
            response = httpClient.execute(httpGet);
            entity = response.getEntity();
            responseContent = EntityUtils.toString(entity, "UTF-8");
        }catch(Exception e){
            e.printStackTrace();
            //TODO Handle Timeout Exeception
            return "";
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
        return responseContent;
    }

    // Do post
    public String doPost(String requestUrl) {
        HttpPost httpPost = this.buildPost(requestUrl);
        CloseableHttpClient httpClient = null;
        CloseableHttpResponse response = null;
        HttpEntity entity = null;
        String responseContent = null;
        try{
            httpClient = HttpClients.createDefault();
                response = httpClient.execute(httpPost);
            entity = response.getEntity();
            responseContent = EntityUtils.toString(entity, "UTF-8");
        }catch(Exception e){
            e.printStackTrace();
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
        return responseContent;
    }

    // Do delete
    public String doDelete(String requestUrl) {
        HttpDelete httpDelete = this.buildDelete(requestUrl);
        CloseableHttpClient httpClient = null;
        CloseableHttpResponse response = null;
        HttpEntity entity = null;
        String responseContent = null;
        try{
            httpClient = HttpClients.createDefault();
            response = httpClient.execute(httpDelete);
            entity = response.getEntity();
            responseContent = EntityUtils.toString(entity, "UTF-8");
        }catch(Exception e){
            e.printStackTrace();
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
        return responseContent;
    }

    // Do put
    public String doPut(String requestUrl) {
        HttpPut httpPut = this.buildPut(requestUrl);
        CloseableHttpClient httpClient = null;
        CloseableHttpResponse response = null;
        HttpEntity entity = null;
        String responseContent = null;
        try{
            httpClient = HttpClients.createDefault();
            response = httpClient.execute(httpPut);
            entity = response.getEntity();
            responseContent = EntityUtils.toString(entity, "UTF-8");
        }catch(Exception e){
            e.printStackTrace();
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
        return responseContent;
    }

    // Add default's headers
    private void addDefaultHeaders() {
        Long timestamp = System.currentTimeMillis();
        this.addClientHeaders(AUTH_HEADER_APIKEY,MLSSConfig.APP_KEY);
        this.addClientHeaders(AUTH_HEADER_TIMESTAMP, Long.toString(timestamp));
        this.addClientHeaders(AUTH_HEADER_SIGNATURE, MLSSConfig.APP_SIGN);
        this.addClientHeaders(AUTH_HEADER_AUTH_TYPE, MLSSConfig.AUTH_TYPE);
        this.addClientHeaders(CONTENT_TYPE, "application/json");
        this.addClientHeaders(AUTH_HEADER_USERID, this.user);
        this.addClientHeaders(AUTH_HEADER_ACCEPT_ENCODING, "gzip");
        this.addClientHeaders(AUTH_HEADER_ACCEPT_ENCODING, "gzip");
        this.addClientHeaders(AUTH_HEADER_CONTENT_ENCODING,"gzip");
    }

    // Get object's type
    private String getType(Object object){
        String typeName = object.getClass().getName();
        Integer length = typeName.lastIndexOf(".");
        return typeName.substring(length + 1);
    }

    @Override
    public String toString() {
        return "ApiClient{" +
                "baseUrl='" + baseUrl + '\'' +
                ", clientRequestHeaders=" + clientHeaders +
                ", clientQueryParams=" + clientQueryParams +
                ", clientFormParams=" + clientFormParams +
                '}';
    }
}
