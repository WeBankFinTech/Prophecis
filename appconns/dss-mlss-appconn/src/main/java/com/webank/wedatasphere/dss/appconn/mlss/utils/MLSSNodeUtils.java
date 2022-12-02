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

package com.webank.wedatasphere.dss.appconn.mlss.utils;

import com.webank.wedatasphere.dss.standard.app.development.listener.core.ExecutionRequestRefContext;

/**
 * @author allenlliu
 * @date 2021/6/25 17:03
 */

public class MLSSNodeUtils {

    public static String getUser(ExecutionRequestRefContext requestRef) {
        return requestRef.getRuntimeMap().get("wds.dss.workflow.submit.user").toString();
    }

}
