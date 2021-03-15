### Prophecis-DI module is built based on the [FfDL](https://github.com/IBM/FfDL). 

**The main modifications are as follows**:

[1] Integrate Kubeflow Arena，Provide distributed tensorflow task ability.

[2] Modify the creation mode of single machine modeling task：remove helper and job jobmonitor in task, and change deploy pod to deploy job.

[3] The log collection service is changed to daemonset, and the collection tool is changed to fluent bit.

[4] The task status update mode is changed to an independent service job monitor.

[5] Add user GUID control in container data directory.

[6] Enhance CLI, added parameter replacement of yaml template, and the train command was modified to websocket connect, providing log and state.

[7] The code file storage server is changed to Minio.
