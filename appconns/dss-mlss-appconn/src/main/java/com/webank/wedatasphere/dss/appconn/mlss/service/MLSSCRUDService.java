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

package com.webank.wedatasphere.dss.appconn.mlss.service;

import com.webank.wedatasphere.dss.appconn.mlss.operation.MLSSRefCopyOperation;
import com.webank.wedatasphere.dss.appconn.mlss.operation.MLSSRefCreationOperation;
import com.webank.wedatasphere.dss.appconn.mlss.operation.MLSSRefDeletionOperation;
import com.webank.wedatasphere.dss.appconn.mlss.operation.MLSSRefUpdateOperation;
import com.webank.wedatasphere.dss.appconn.mlss.operation.MLSSRefCopyOperation;
import com.webank.wedatasphere.dss.appconn.mlss.operation.MLSSRefCreationOperation;
import com.webank.wedatasphere.dss.appconn.mlss.operation.MLSSRefDeletionOperation;
import com.webank.wedatasphere.dss.appconn.mlss.operation.MLSSRefUpdateOperation;
import com.webank.wedatasphere.dss.appconn.mlss.utils.MLSSConfig;
import com.webank.wedatasphere.dss.standard.app.development.operation.RefCopyOperation;
import com.webank.wedatasphere.dss.standard.app.development.operation.RefCreationOperation;
import com.webank.wedatasphere.dss.standard.app.development.operation.RefDeletionOperation;
import com.webank.wedatasphere.dss.standard.app.development.operation.RefUpdateOperation;
import com.webank.wedatasphere.dss.standard.app.development.service.AbstractRefCRUDService;

import java.util.Map;

public class MLSSCRUDService extends AbstractRefCRUDService {

    @Override
    protected RefCreationOperation createRefCreationOperation() {
        return new MLSSRefCreationOperation(this);
    }

    @Override
    protected RefCopyOperation createRefCopyOperation() {
        return new MLSSRefCopyOperation(null,this);
    }

    @Override
    protected RefUpdateOperation createRefUpdateOperation() {
        return new MLSSRefUpdateOperation(this);
    }

    @Override
    protected RefDeletionOperation createRefDeletionOperation() {
        return new MLSSRefDeletionOperation(this);
    }

    public void initMLSSConfig(){
        MLSSConfig.BASE_URL = this.getAppInstance().getBaseUrl();
        Map<String, Object> config = this.getAppInstance().getConfig();
        MLSSConfig.APP_KEY = String.valueOf(config.get("MLSS-SecretKey"));
        MLSSConfig.APP_SIGN = String.valueOf(config.get("MLSS-APPSignature"));
        MLSSConfig.AUTH_TYPE =  String.valueOf(config.get("MLSS-Auth-Type"));
        MLSSConfig.TIMESTAMP =  String.valueOf(config.get("MLSS-APPSignature"));
    }

}
