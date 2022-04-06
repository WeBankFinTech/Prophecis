package com.webank.wedatasphere.dss.appconn.mlflow.utils;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.webank.wedatasphere.dss.standard.app.development.listener.common.RefExecutionState;
import com.fasterxml.jackson.dataformat.yaml.YAMLMapper;

import java.io.File;
import java.io.FileWriter;
import java.io.IOException;
import java.util.UUID;

public class JobParser {


    static String fileRootPath = "./";

    public static RefExecutionState transformStatus(String status) {
        if ("PENDING".equals(status)){
            return RefExecutionState.Accepted;
        }else if("PROCESSING".equals(status)){
            return RefExecutionState.Accepted;
        }else if("FAILED".equals(status)){
            return RefExecutionState.Failed;
        }else if("QUEUED".equals(status)){
            return RefExecutionState.Accepted;
        }else if("RUNNING".equals(status)) {
            return RefExecutionState.Running;
        }
        return RefExecutionState.Success;
    }


    public static String TransformManifest(String jobContent) throws JsonProcessingException {
        String filename = fileRootPath + UUID.randomUUID().toString();
        JsonNode jsonNodeTree = new ObjectMapper().readTree(jobContent);
        String jsonAsYaml = new YAMLMapper().writeValueAsString(jsonNodeTree);

        FileWriter fileWriter = null;
        try {
            fileWriter = new FileWriter(filename);
            fileWriter.write(jsonAsYaml);
        } catch (IOException e) {
            e.printStackTrace();
        } finally {
            try {
                assert fileWriter != null;
                fileWriter.flush();
            } catch (IOException e) {
                e.printStackTrace();
            }
        }
        return filename;
    }

    public static boolean clearFile(String filename) {
        File file = new File(filename);
        if (file.isFile() && file.exists()){
            return file.delete();
        }
        return true;
    }


}
