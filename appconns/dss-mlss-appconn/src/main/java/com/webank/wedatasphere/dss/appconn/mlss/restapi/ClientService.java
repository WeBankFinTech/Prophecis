package com.webank.wedatasphere.dss.appconn.mlss.restapi;

import com.webank.wedatasphere.dss.appconn.mlss.utils.MLSSConfig;
import org.apache.commons.codec.binary.Hex;
import org.apache.commons.codec.digest.HmacUtils;

import javax.crypto.Mac;
import javax.crypto.spec.SecretKeySpec;
import java.util.HashMap;
import java.util.Map;

/**
 * Create by v_wbyxzhong
 * Date Time 2020/10/20 16:17
 */
public class ClientService {
//    public static Map<String, Object> clientMap = new HashMap<>();
//
//    public static MLSSConfig getMlssConfig() {
//        return (MLSSConfig) clientMap.get("mlssConfig");
//    }
//
//    public static String getUser(){
//        return (String) clientMap.get("user");
//    }
//
//    public static String getSign(String secretKey, Long timestamp) {
//        try {
//            String queryString = secretKey + ":" + timestamp;
//            SecretKeySpec signingKey = new SecretKeySpec(secretKey.getBytes(), "HmacSHA1");
//            HmacUtils hmacUtils = new HmacUtils();
//            Mac mac = Mac.getInstance("HmacSHA1");
//            mac.init(signingKey);
//            byte[] rawHmac = mac.doFinal(queryString.getBytes("utf-8"));
//            return Hex.encodeHexString(rawHmac);
//        } catch (Exception e) {
//            return null;
//        }
//    }
}
