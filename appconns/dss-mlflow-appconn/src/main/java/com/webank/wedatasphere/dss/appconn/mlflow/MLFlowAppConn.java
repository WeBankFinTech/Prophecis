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

package com.webank.wedatasphere.dss.appconn.mlflow;

import com.webank.wedatasphere.dss.appconn.core.ext.OnlyDevelopmentAppConn;
import com.webank.wedatasphere.dss.appconn.core.ext.ThirdlyAppConn;
import com.webank.wedatasphere.dss.appconn.core.impl.AbstractAppConn;
import com.webank.wedatasphere.dss.appconn.core.impl.AbstractOnlySSOAppConn;
import com.webank.wedatasphere.dss.appconn.mlflow.utils.MLFlowConfig;
import com.webank.wedatasphere.dss.standard.app.development.standard.DevelopmentIntegrationStandard;
import com.webank.wedatasphere.dss.standard.app.structure.StructureIntegrationStandard;
import org.apache.linkis.common.conf.CommonVars;

public class MLFlowAppConn extends AbstractAppConn implements OnlyDevelopmentAppConn {

    public static final String MLFlow_APPCONN_NAME = CommonVars.apply("wds.dss.appconn.mlflow.name", "mlflow").getValue();


    private MLFlowDevelopmentIntegrationStandard developmentIntegrationStandard;

    @Override
    protected void initialize() {
        this.developmentIntegrationStandard = new MLFlowDevelopmentIntegrationStandard();
    }

    @Override
    public DevelopmentIntegrationStandard getOrCreateDevelopmentStandard() {
        return developmentIntegrationStandard;
    }

}
