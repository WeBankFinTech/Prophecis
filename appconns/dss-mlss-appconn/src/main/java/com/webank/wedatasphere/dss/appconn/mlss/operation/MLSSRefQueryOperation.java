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


import com.webank.wedatasphere.dss.standard.app.development.operation.AbstractDevelopmentOperation;
import com.webank.wedatasphere.dss.standard.app.development.operation.RefQueryJumpUrlOperation;
import com.webank.wedatasphere.dss.standard.app.development.ref.QueryJumpUrlResponseRef;
import com.webank.wedatasphere.dss.standard.app.development.ref.impl.OnlyDevelopmentRequestRef;

import com.webank.wedatasphere.dss.standard.app.development.ref.impl.ThirdlyRequestRef;
import com.webank.wedatasphere.dss.standard.app.sso.Workspace;

import com.webank.wedatasphere.dss.standard.common.exception.operation.ExternalOperationFailedException;

import java.io.UnsupportedEncodingException;
import java.util.Map;

public class MLSSRefQueryOperation extends AbstractDevelopmentOperation<ThirdlyRequestRef.QueryJumpUrlRequestRefImpl, QueryJumpUrlResponseRef>
        implements RefQueryJumpUrlOperation<ThirdlyRequestRef.QueryJumpUrlRequestRefImpl, QueryJumpUrlResponseRef> {


    @Override
    public QueryJumpUrlResponseRef query(ThirdlyRequestRef.QueryJumpUrlRequestRefImpl ref) throws ExternalOperationFailedException {
        Map<String, Object> jobContent = ref.getRefJobContent();
        logger.info("query", jobContent.toString());
        String baseUrl = service.getAppInstance().getBaseUrl();
        String expId = jobContent.get("expId").toString();
        String contextId = "";
        Workspace workspace = ref.getWorkspace();
        Long workspaceId = workspace.getWorkspaceId();
        Map<String, String> cookies = workspace.getCookies();
        String ticket = cookies.get("linkis_user_session_ticket_id_v1");
        try {
            ticket = java.net.URLEncoder.encode(ticket, "UTF-8");
        } catch (UnsupportedEncodingException e) {
            logger.error("ticket transform error:" + e.toString());
        }
        String retJumpUrl = baseUrl + "#/mlFlow?expId=" + expId + "&contextID=" + contextId
                + "&workspaceId=" + workspaceId +
                "&bdp-user-ticket-id=" + ticket;

        return QueryJumpUrlResponseRef.newBuilder().setJumpUrl(retJumpUrl).success();
    }
}
