package com.webank.wedatasphere.dss.appconn.mlss.ref;

import com.google.common.collect.Maps;
import com.webank.wedatasphere.dss.standard.app.development.ref.DSSCommonResponseRef;
import com.webank.wedatasphere.dss.standard.app.sso.plugin.SSOIntegrationConf;
import com.webank.wedatasphere.dss.standard.common.exception.operation.ExternalOperationFailedException;

import java.util.Map;

/**
 * created by cooperyang on 2021/8/2
 * Description:
 */
public class MLSSCopyResponseRef extends DSSCommonResponseRef {


    private Map<String, Object> importedMap = Maps.newHashMap();
    private Map<String, Object> newJobContent = Maps.newHashMap();

    @SuppressWarnings("unchecked")
    public MLSSCopyResponseRef(Map<String, Object> jobContent, String responseBody, String nodeType) throws Exception {
        super(responseBody);
        responseMap = SSOIntegrationConf.gson().fromJson(responseBody,Map.class);
        if("linkis.appconn.MLSS.widget".equalsIgnoreCase(nodeType)){
            Map<String, Object> payload = (Map<String, Object>) jobContent.get("data");
            Long id = ((Double) Double.parseDouble(payload.get("widgetId").toString())).longValue();
            payload.put("widgetId", ((Map<String, Double>) ((Map<String, Object>) responseMap.get("data")).get("widget")).get(id.toString()).doubleValue());
        } else if("linkis.appconn.MLSS.display".equalsIgnoreCase(nodeType)){
            Map<String, Object> payload = (Map<String, Object>) jobContent.get("payload");
            Long id = ((Double) Double.parseDouble(payload.get("id").toString())).longValue();
            payload.put("id", ((Map<String, Double>) ((Map<String, Object>) responseMap.get("data")).get("display")).get(id.toString()).doubleValue());
        } else if("linkis.appconn.MLSS.dashboard".equalsIgnoreCase(nodeType)){
            Map<String, Object> payload = (Map<String, Object>) jobContent.get("payload");
            Long id = ((Double) Double.parseDouble(payload.get("id").toString())).longValue();
            payload.put("id", ((Map<String, Double>) ((Map<String, Object>) responseMap.get("data")).get("dashboardPortal")).get(id.toString()).doubleValue());
        } else {
            throw new ExternalOperationFailedException(90177, "Unknown task type " + nodeType, null);
        }
        this.newJobContent = jobContent;
    }

    @Override
    public Map<String, Object> toMap() {
        return newJobContent;
    }


}
