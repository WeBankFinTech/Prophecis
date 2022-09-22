package com.webank.wedatasphere.dss.appconn.mlflow.execution;

import com.webank.wedatasphere.dss.appconn.mlflow.job.MLFlowJob;
import com.webank.wedatasphere.dss.standard.app.development.listener.common.AbstractRefExecutionAction;
import com.webank.wedatasphere.dss.standard.app.development.listener.common.LongTermRefExecutionAction;
import com.webank.wedatasphere.dss.standard.app.development.listener.common.RefExecutionState;
import com.webank.wedatasphere.dss.standard.app.development.listener.core.ExecutionRequestRefContext;
import com.webank.wedatasphere.dss.standard.app.development.listener.ref.RefExecutionRequestRef;

public class MLFlowExecutionAction extends AbstractRefExecutionAction implements LongTermRefExecutionAction {

    private Integer from;
    private Integer size;
    private Integer logCount;
    private RefExecutionState state;
    private MLFlowJob job;
    private ExecutionRequestRefContext executionContext;

    public MLFlowJob getJob() {
        return job;
    }

    public void setJob(MLFlowJob job) {
        this.job = job;
    }

    public MLFlowExecutionAction(MLFlowJob mlFlowJob, ExecutionRequestRefContext executionContext) {
        this.state = RefExecutionState.Accepted;
        this.job = mlFlowJob;
        this.size = 10000;
        this.from = 0;
        this.logCount = 0;
        this.executionContext = executionContext;
    }

    public void setState(RefExecutionState state) {
        this.state = state;
    }

    public RefExecutionState getState() {
        return state;
    }

    public void setFrom(Integer from) { this.from = from; }

    public Integer getFrom() { return from; }

    public void setSize(Integer size) { this.size = size; }

    public Integer getSize() { return size; }

    public void setLogCount(Integer logCount) { this.logCount = logCount; }

    public Integer getLogCount() { return logCount; }

    public ExecutionRequestRefContext getExecutionRequestRefContext() {
        return executionContext;
    }

    public void setExecutionRequestRefContext(ExecutionRequestRefContext ExecutionRequestRefContext) {
        this.executionContext = ExecutionRequestRefContext;
    }

    @Override
    public void setSchedulerId(int schedulerId) {

    }

    @Override
    public int getSchedulerId() {
        return 0;
    }
}
