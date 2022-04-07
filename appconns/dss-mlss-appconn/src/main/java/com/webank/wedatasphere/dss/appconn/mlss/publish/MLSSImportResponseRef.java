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

package com.webank.wedatasphere.dss.appconn.mlss.publish;

import com.google.common.collect.Maps;
import com.webank.wedatasphere.dss.standard.app.development.ref.DSSCommonResponseRef;
import com.webank.wedatasphere.dss.standard.common.exception.operation.ExternalOperationFailedException;

import java.util.Map;

public class MLSSImportResponseRef extends DSSCommonResponseRef {
    Map<String, Object> importedMap = Maps.newHashMap();
    Map<String, Object> newJobContent = Maps.newHashMap();

    @SuppressWarnings("unchecked")
    public MLSSImportResponseRef(Map<String, Object> jobContent, String responseBody, String nodeType, Object projectId) throws Exception {
        super(responseBody);
        this.newJobContent = jobContent;
    }


    @Override
    public Map<String, Object> toMap() {
        return newJobContent;
    }
}
