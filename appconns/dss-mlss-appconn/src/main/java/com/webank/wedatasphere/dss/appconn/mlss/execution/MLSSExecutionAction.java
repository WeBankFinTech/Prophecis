package com.webank.wedatasphere.dss.appconn.mlss.execution;

import com.webank.wedatasphere.dss.standard.app.development.listener.common.AbstractRefExecutionAction;
import com.webank.wedatasphere.dss.standard.app.development.listener.common.LongTermRefExecutionAction;
import com.webank.wedatasphere.dss.standard.app.development.listener.common.RefExecutionState;
import com.webank.wedatasphere.dss.standard.app.development.listener.core.ExecutionRequestRefContext;

/**
 * Create by v_wbyxzhong
 * Date Time 2020/10/10 17:24
 */
public class MLSSExecutionAction extends AbstractRefExecutionAction implements LongTermRefExecutionAction {

    private String appId;
    private String user;
    private ExecutionRequestRefContext executionContext;
    private RefExecutionState state;


    public MLSSExecutionAction(String appId, String user, ExecutionRequestRefContext executionContext) {
        this.appId = appId;
        this.user = user;
        this.executionContext = executionContext;
    }

    public ExecutionRequestRefContext getExecutionContext() {
        return executionContext;
    }

    public void setExecutionContext(ExecutionRequestRefContext executionContext) {
        this.executionContext = executionContext;
    }

    public String getAppId() {
        return appId;
    }

    public String getUser() {
        return user;
    }

    @Override
    public void setSchedulerId(int schedulerId) {

    }

    @Override
    public int getSchedulerId() {
        return 0;
    }

    public RefExecutionState getState() {
        return state;
    }

    public void setState(RefExecutionState state) {
        this.state = state;
    }
}
