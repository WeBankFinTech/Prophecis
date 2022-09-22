package com.webank.wedatasphere.dss.appconn.mlss.operation;

import com.google.gson.JsonObject;
import com.webank.wedatasphere.dss.appconn.mlss.restapi.ExperimentAPI;
import com.webank.wedatasphere.dss.common.utils.DSSCommonUtils;
import com.webank.wedatasphere.dss.flow.execution.entrance.conf.FlowExecutionEntranceConfiguration;
import com.webank.wedatasphere.dss.standard.app.development.operation.AbstractDevelopmentOperation;
import com.webank.wedatasphere.dss.standard.app.development.operation.RefCopyOperation;
import com.webank.wedatasphere.dss.standard.app.development.ref.RefJobContentResponseRef;
import com.webank.wedatasphere.dss.standard.app.development.ref.impl.ThirdlyRequestRef;
import com.webank.wedatasphere.dss.standard.common.exception.operation.ExternalOperationFailedException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.HashMap;
import java.util.Map;


public class MLSSRefCopyOperation extends AbstractDevelopmentOperation<ThirdlyRequestRef.CopyWitContextRequestRefImpl, RefJobContentResponseRef>
        implements RefCopyOperation<ThirdlyRequestRef.CopyWitContextRequestRefImpl> {

    private static final Logger logger = LoggerFactory.getLogger(MLSSRefCopyOperation.class);

    @Override
    public RefJobContentResponseRef copyRef(ThirdlyRequestRef.CopyWitContextRequestRefImpl requestRef) throws ExternalOperationFailedException {
        //1. Init
        logger.info(requestRef.toString());
        Map<String, Object> jobContent = requestRef.getRefJobContent();
        HashMap contextInfo = DSSCommonUtils.COMMON_GSON.fromJson(DSSCommonUtils.COMMON_GSON.fromJson(
        requestRef.getParameter("dssContextId").toString(), HashMap.class).get("value").toString(), HashMap.class);
//        String flowName = jobContent.get("orchestrationName").toString();
//        String flowVersion = contextInfo.get("version").toString();
        long flowID = Long.parseLong(jobContent.get("orchestrationId").toString());
//        long projectID = Long.parseLong(requestRef.getParameter("dssProjectId").toString());
//        String projectName = contextInfo.get("project").toString();
//        String expDesc = "DSS," + (flowName == null ? "NULL" : flowName) + "," +
//                (flowVersion == null ? "NULL" : flowVersion) + "," + flowID + "," +
//                projectName + "," + "NULL"  + "," + flowVersion;
        String user = requestRef.getUserName();
        if (user.endsWith("_f")) {
            String[] userStrArray = user.split("_");
            user = userStrArray[0];
        }
        String expName = requestRef.getName();
        String createType = "DSS";
        //2. Execute Create Request
        String expId =  jobContent.get("expId").toString();
        JsonObject jsonObj = ExperimentAPI.copyExperiment(user, expId, createType);
        if (jsonObj == null) {
            logger.error("Create experiment failed, jsonObj is null");
            return RefJobContentResponseRef.newBuilder().setRefJobContent(jobContent).error("Create experiment failed, jsonObj is null");
        }
        //TODO: UPDATE EXPERIMENT JSON
        //3. Update node jobContent & Return Response
        Long resExpId = jsonObj.get("result").getAsJsonObject().get("id").getAsLong();
        String resExpName = jsonObj.get("result").getAsJsonObject().get("exp_name").getAsString();
        logger.info("Create MLSS experiment success expId:" + resExpId + ",expName:" + resExpName);
        jobContent.put("expId", resExpId);
        jobContent.put("expName", expName);
        jobContent.put("status", 200);
        return RefJobContentResponseRef.newBuilder().setRefJobContent(jobContent).success();
    }
}
