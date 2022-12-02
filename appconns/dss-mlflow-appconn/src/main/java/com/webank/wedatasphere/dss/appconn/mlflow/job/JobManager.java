package com.webank.wedatasphere.dss.appconn.mlflow.job;

import com.google.gson.internal.LinkedTreeMap;
import com.webank.wedatasphere.dss.appconn.mlflow.execution.MLFlowExecutionAction;
import com.webank.wedatasphere.dss.common.utils.DSSCommonUtils;
import com.webank.wedatasphere.dss.standard.app.development.listener.common.RefExecutionAction;
import com.webank.wedatasphere.dss.standard.app.development.listener.common.RefExecutionState;
import com.webank.wedatasphere.dss.standard.app.development.listener.core.ExecutionRequestRefContext;
import com.webank.wedatasphere.dss.standard.app.development.listener.ref.RefExecutionRequestRef;
import org.apache.linkis.common.log.LogUtils;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.Map;


public class JobManager {

    public static final String JOB_TYPE_STRING = "mlflowJobType";
    public static final String JOB_TYPE_DI = "DistributedModel";
    public static final String JOB_TYPE_MODEL_DEPLOY = "ModelDeploy";
    public static final String JOB_TYPE_IMAGE_BUILD = "ImageBuild";
    public static final String JOB_TYPE_MODEL_STORAGE = "ModelStorage";
    public static final String JOB_TYPE_REPORT_PUSH = "ReportPush";
    public static final String JOB_TYPE_MODEL_MONITORING = "ModelMonitor";

    public static RefExecutionAction submit(Map jobContent,
                                            ExecutionRequestRefContext executionRequest, HashMap<String, Object> contextInfo) {
        String jobType = (String) jobContent.get(JOB_TYPE_STRING);
        String user = getUser(contextInfo);
        MLFlowJob job = null;
        if (jobType.equals(JOB_TYPE_DI)) {
            job = new DistributedModelJob(user);
        } else if (jobType.equals(JOB_TYPE_MODEL_DEPLOY)) {
            job = new ModelDeployJob(user);
        } else if (jobType.equals(JOB_TYPE_MODEL_STORAGE)) {
            job = new ModelStorageJob(user);
        } else if (jobType.equals(JOB_TYPE_IMAGE_BUILD)) {
            job = new ImageBuildJob(user);
        } else if (jobType.equals(JOB_TYPE_REPORT_PUSH)) {
            job = new ReportPushJob(user);
        } else if (jobType.equals(JOB_TYPE_MODEL_MONITORING)) {
            job = new ModelMonitoringJob(user);
            executionRequest.appendLog(LogUtils.generateInfo("mlflow start execution action jobContent =" + jobContent));
        } else {
            //TODO Throw Exception
            return null;
        }
        boolean jobSubmitBool = job.submit(jobContent, executionRequest);
        //TODO Failed Control
        if (!jobSubmitBool) {
            return null;
        }
        return new MLFlowExecutionAction(job, executionRequest);
    }

    public static RefExecutionState state(MLFlowExecutionAction mlflowAction) {
        if (mlflowAction == null) {
            return RefExecutionState.Failed;
        }
        //UPDATE STATE
        MLFlowJob job = mlflowAction.getJob();
        RefExecutionState state = job.status();
        mlflowAction.setState(state);
//        TODO ADD UPDATE Log
//        ExecutionRequestRefContext ExecutionRequestRefContext = (ExecutionRequestRefContext) ClientService.clientMap.get("ExecutionRequestRefContext");
//        if (ExecutionRequestRefContext != null) {
//            Log getLog = JobExecClient.log(mlflowAction.getUser(), mlflowAction.getJobId(), mlflowAction.getSize(), mlflowAction.getFrom());
//            if(getLog != null && getLog.getLogList() != null) {
//                mlflowAction.setSize(getLog.getLogList().size());
//                mlflowAction.setFrom(mlflowAction.getFrom() + getLog.getLogList().size());
//                for (String log : getLog.getLogList()) {
//                    ExecutionRequestRefContext.appendLog(log);
//                }
//            }
//        }
        return mlflowAction.getState();
    }


//    public static MLFlowCompletedExecutionResponseRef result(MLFlowExecutionAction mlflowAction, ExecutionRequestRefContext executionRequestRefContext){
//        MLFlowCompletedExecutionResponseRef result = new MLFlowCompletedExecutionResponseRef(0,"");
//        if (mlflowAction == null) {
//            executionRequestRefContext.appendLog(LogUtils.generateERROR("MLFlow Action is null, Job RUN Error. "));
////            result.setIsSucceed(false);
//            return result;
//        }
//        //TODO: add Log last time in result DI TYPE
//        if (mlflowAction.getState().equals(RefExecutionState.Failed)) {
//            executionRequestRefContext.appendLog(LogUtils.generateERROR("Job Status is failed, job run error. "));
////            result.setIsSucceed(false);
//            return result;
//        }
////        result.setIsSucceed(true);
//        executionRequestRefContext.appendLog(LogUtils.generateInfo("任务执行成功，TrainingId:" + mlflowAction.getId()));
//        return result;
////        if (executionRequestRefContext != null) {
////            while(true) {
////                try {
////                    Thread.sleep(12000);
////                    if (mlflowAction.getLogCount() > 5) {
////                        LOGGER.info("Log list is null");
////                        break;
////                    }
////                    Log getLog = JobExecClient.log(mlflowAction.getUser(), mlflowAction.getJobId(), mlflowAction.getSize(), mlflowAction.getFrom());
////                    if (getLog != null && getLog.getLogList() != null) {
////                        mlflowAction.setSize(getLog.getLogList().size());
////                        mlflowAction.setFrom(mlflowAction.getFrom() + getLog.getLogList().size());
////                        for (String log : getLog.getLogList()) {
////                            executionRequestRefContext.appendLog(log);
////                        }
////                    }else{
////                        mlflowAction.setLogCount(mlflowAction.getLogCount()+1);
////                    }
////                } catch (InterruptedException e) {
////                    e.printStackTrace();
////                }
////            }
////        }
//    }

    public static boolean kill(MLFlowExecutionAction mlflowAction) {
        MLFlowJob job = mlflowAction.getJob();
        return job.kill();
    }

    public static String log(MLFlowExecutionAction mlflowAction) {
        MLFlowJob job = mlflowAction.getJob();
        return job.log();
    }

    public static float progress(MLFlowExecutionAction mlflowAction) {
        if (mlflowAction == null) {
            return 1;
        }
        RefExecutionState state = JobManager.state(mlflowAction);
        if (state != RefExecutionState.Success && state != RefExecutionState.Failed) {
            return 0.5f;
        }
        return 1;
    }

    public static LinkedTreeMap paramsTransfer(LinkedTreeMap<String, Object> content) {
        for (Map.Entry<String, Object> entry : content.entrySet()) {
            if (entry.getKey().toUpperCase().contains("ID")) {
                content.put(entry.getKey(), Double.valueOf(entry.getValue().toString()).intValue());
            }
        }

        if (content.containsKey("service_post_models")) {
            ArrayList models = (ArrayList) content.get("service_post_models");
            for (int i = 0; i < models.size(); i++) {
                LinkedTreeMap<String, Object> model = (LinkedTreeMap) models.get(i);
                for (Map.Entry<String, Object> entry : model.entrySet()) {
                    if (entry.getKey().toUpperCase().contains("ID")) {
                        model.put(entry.getKey(), Double.valueOf(entry.getValue().toString()).intValue());
                    }
                }

            }
        }
        return content;
    }


    private static String getUser(Map<String, Object> contextInfo) {
        return contextInfo.get("user").toString();
    }

}
