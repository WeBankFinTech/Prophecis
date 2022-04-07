/*
 * Copyright 2019 WeBank
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package com.webank.wedatasphere.dss.appconn.mlss.operation;

import com.google.gson.internal.LinkedTreeMap;
import com.webank.wedatasphere.dss.appconn.mlss.ref.MLSSOpenRequestRef;
import com.webank.wedatasphere.dss.appconn.mlss.ref.MLSSOpenResponseRef;
import com.webank.wedatasphere.dss.common.utils.DSSCommonUtils;
import com.webank.wedatasphere.dss.standard.app.development.listener.common.AsyncExecutionRequestRef;
import com.webank.wedatasphere.dss.standard.app.development.operation.RefQueryOperation;
import com.webank.wedatasphere.dss.standard.app.development.ref.OpenRequestRef;
import com.webank.wedatasphere.dss.standard.app.development.service.DevelopmentService;
import com.webank.wedatasphere.dss.standard.app.sso.Workspace;
import com.webank.wedatasphere.dss.standard.common.entity.ref.ResponseRef;
import com.webank.wedatasphere.dss.standard.common.exception.operation.ExternalOperationFailedException;
import org.apache.linkis.server.BDPJettyServerHelper;
import org.apache.hadoop.util.hash.Hash;

import java.util.HashMap;
import java.util.Map;

public class MLSSRefQueryOperation implements RefQueryOperation<OpenRequestRef> {

    DevelopmentService developmentService;

    @Override
    public ResponseRef query(OpenRequestRef ref) throws ExternalOperationFailedException {
        MLSSOpenRequestRef openRequestRef = (MLSSOpenRequestRef) ref;

        try {
//            String externalContent = BDPJettyServerHelper.jacksonJson().writeValueAsString(executionRequestRef.getJobContent());
            Long projectId = (Long) openRequestRef.getParameter("projectId");
            String baseUrl = openRequestRef.getParameter("redirectUrl").toString();
            String jumpUrl = baseUrl;
            String expId =  ((HashMap)openRequestRef.getParameter("params")).get("expId").toString();
            String contextId = ((HashMap)openRequestRef.getParameter("params")).get("contextID").toString();
            Workspace workspace = ((Workspace) ((HashMap) openRequestRef.getParameter("params")).get("workspace"));
            String workspaceId = workspace.getWorkspaceName();
            String operationStr = workspace.getOperationStr();
            HashMap operation = DSSCommonUtils.COMMON_GSON.fromJson(operationStr,HashMap.class);
            LinkedTreeMap cookies = (LinkedTreeMap) operation.get("cookies");
            String retJumpUrl = jumpUrl + "?expId=" +expId+"&contextID="+ contextId
                    +"&workspaceId=" + workspaceId +
                    "&bdp-user-ticket-id=" + cookies.get("bdp-user-ticket-id");
            Map<String,String> retMap = new HashMap<>();
            retMap.put("jumpUrl",retJumpUrl);
            return new MLSSOpenResponseRef(DSSCommonUtils.COMMON_GSON.toJson(retMap), 0);
        } catch (Exception e) {
            throw new ExternalOperationFailedException(90177, "Failed to parse jobContent ", e);
        }
    }


    @Override
    public void setDevelopmentService(DevelopmentService service) {
        this.developmentService = service;
    }
}
