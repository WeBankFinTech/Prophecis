package com.webank.wedatasphere.dss.appconn.mlflow.job;

import com.webank.wedatasphere.dss.standard.app.development.listener.common.RefExecutionState;
import com.webank.wedatasphere.dss.standard.app.development.listener.core.ExecutionRequestRefContext;
import com.webank.wedatasphere.dss.standard.app.development.listener.ref.RefExecutionRequestRef;

import java.util.Map;

public abstract class MLFlowJob {

    private String jobType;
    private String user;
    private Integer from = 0;
    private Integer size = 100;
    private Integer logCount;
    private String status;
    private ExecutionRequestRefContext executionContext;

    ExecutionRequestRefContext getExecutionContext() {
        return executionContext;
    }

    void setExecutionContext(ExecutionRequestRefContext executionContext) {
        this.executionContext = executionContext;
    }


    public abstract boolean submit(Map jobContent, ExecutionRequestRefContext ExecutionRequestRef);

    public abstract RefExecutionState status();

    public abstract boolean kill();

    public abstract String log();

    public abstract RefExecutionState transformStatus(String status);

    public String getJobType() {
        return jobType;
    }

    void setJobType(String jobType) {
        this.jobType = jobType;
    }

    public String getUser() {
        return user;
    }

    public void setUser(String user) {
        this.user = user;
    }

    Integer getFrom() {
        return from;
    }

    public void setFrom(Integer from) {
        this.from = from;
    }

    Integer getSize() {
        return size;
    }

    public void setSize(Integer size) {
        this.size = size;
    }

    public Integer getLogCount() {
        return logCount;
    }

    public void setLogCount(Integer logCount) {
        this.logCount = logCount;
    }

    public String getStatus() {
        return status;
    }

    public void setStatus(String status) {
        this.status = status;
    }

}
