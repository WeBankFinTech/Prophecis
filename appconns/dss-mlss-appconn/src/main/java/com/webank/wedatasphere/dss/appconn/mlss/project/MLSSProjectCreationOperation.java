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

package com.webank.wedatasphere.dss.appconn.mlss.project;

import com.webank.wedatasphere.dss.appconn.mlss.MLSSAppConn;
import com.webank.wedatasphere.dss.standard.app.sso.request.SSORequestOperation;
import com.webank.wedatasphere.dss.standard.app.structure.StructureService;
import com.webank.wedatasphere.dss.standard.app.structure.project.ProjectCreationOperation;
import com.webank.wedatasphere.dss.standard.app.structure.project.ProjectRequestRef;
import com.webank.wedatasphere.dss.standard.app.structure.project.ProjectResponseRef;
import com.webank.wedatasphere.dss.standard.common.exception.operation.ExternalOperationFailedException;
import org.apache.linkis.httpclient.request.HttpAction;
import org.apache.linkis.httpclient.response.HttpResult;
import org.apache.linkis.server.conf.ServerConfiguration;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class MLSSProjectCreationOperation implements ProjectCreationOperation {

    private static Logger logger = LoggerFactory.getLogger(MLSSProjectCreationOperation.class);
    private final static String projectUrl = "/api/rest_s/" + ServerConfiguration.BDP_SERVER_VERSION() + "/com/webank/wedatasphere/dss/appconn/mlflow/projects";
    private SSORequestOperation<HttpAction, HttpResult> ssoRequestOperation;
    private StructureService structureService;

    public MLSSProjectCreationOperation(StructureService service, SSORequestOperation<HttpAction, HttpResult> ssoRequestOperation) {
        this.structureService = service;
        this.ssoRequestOperation = ssoRequestOperation;
    }
    private String getAppName() {
        return MLSSAppConn.MLSS_APPCONN_NAME;
    }

    @Override
    public void init() {
    }

    @Override
    public ProjectResponseRef createProject(ProjectRequestRef projectRef) throws ExternalOperationFailedException {
        String url = getBaseUrl() + projectUrl;

        MLSSProjectResponseRef MLSSProjectResponseRef = null;
        try {
            MLSSProjectResponseRef = new MLSSProjectResponseRef("",200);
        } catch (Exception e) {
            e.printStackTrace();
        }

        return MLSSProjectResponseRef;
    }

    @Override
    public void setStructureService(StructureService service) {
        this.structureService = service;
    }

    private String getBaseUrl(){
        return structureService.getAppInstance().getBaseUrl();
    }
}
