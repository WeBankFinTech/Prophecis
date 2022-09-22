package com.webank.wedatasphere.dss.appconn.mlss.ref;

import com.webank.wedatasphere.dss.standard.common.entity.ref.ResponseRefBuilder;
import com.webank.wedatasphere.dss.standard.common.entity.ref.ResponseRefImpl;

import java.util.Map;

public class MLSSResponseRefBuilder
        extends ResponseRefBuilder.ExternalResponseRefBuilder<MLSSResponseRefBuilder, ResponseRefImpl> {

    @Override
    public MLSSResponseRefBuilder setResponseBody(String responseBody) {
        super.setResponseBody(responseBody).build();
        Map<String, Object> headerMap = (Map<String, Object>) responseMap.get("header");
        if (headerMap.containsKey("code")) {
            status = getInt(headerMap.get("code"));
            if (status != 0 && status != 200) {
                errorMsg = headerMap.get("msg").toString();
            }
        }
        Object payload = responseMap.get("payload");
        if(payload instanceof Map) {
            setResponseMap((Map<String, Object>) payload);
        }
        return this;
//        return super.setResponseBody(responseBody);
    }

    public static Integer getInt(Object original) {
        if (original instanceof Double) {
            return ((Double) original).intValue();
        }
        return (Integer) original;
    }

}
