// Copyright 2018 The Kubeflow Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tfjob

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"webank/DI/commons/util"
	"webank/DI/commons/workflow"
	"webank/DI/commons/logger"
	"webank/DI/commons/service"
	"webank/DI/lcm/lcmconfig"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/batch/v1"
	"k8s.io/client-go/kubernetes"

	//"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	// to receive values from command operation --worker-selector
	workerSelectors []string
	// to receive values from command operation --evaluator-selector
	evaluatorSelectors []string
	// to receive values from command operation --ps-selector
	psSelectors []string
	// to receive values from command operation --chief-selector
	chiefSelectors []string
	tfjob_chart    = util.GetChartsFolder() + "/tfjob"
)

//func init() {
//	if err := os.Mkdir(tfjob_chart, 0755); err != nil {
//		log.Errorf("Error creating tfjob_chart directory '%s'", tfjob_chart)
//	}
//}

//func NewSubmitTFJobCommand() *cobra.Command {
//	var (
//		submitArgs submitTFJobArgs
//	)
//
//	submitArgs.Mode = "tfjob"
//
//	var command = &cobra.Command{
//		Use:     "tfjob",
//		Short:   "Submit TFJob as training job.",
//		Aliases: []string{"tf"},
//		Run: func(cmd *cobra.Command, args []string) {
//			if len(args) == 0 {
//				cmd.HelpFunc()(cmd, args)
//				os.Exit(1)
//			}
//
//			//util.SetLogLevel(logLevel)
//			//setupKubeconfig()
//			//_, err := InitKubeClient()
//			//if err != nil {
//			//	log.Debugf("Failed due to %v", err)
//			//	fmt.Println(err)
//			//	os.Exit(1)
//			//}
//
//			//err = updateNamespace(cmd)
//			//if err != nil {
//			//	log.Debugf("Failed due to %v", err)
//			//	fmt.Println(err)
//			//	os.Exit(1)
//			//}
//
//			err := submitTFJob(args, &submitArgs, clientset)
//			if err != nil {
//				log.Debugf("Failed due to %v", err)
//				fmt.Println(err)
//				os.Exit(1)
//			}
//		},
//	}
//
//	submitArgs.addCommonFlags(command)
//	//submitArgs.addSyncFlags(command)
//
//	// TFJob
//	command.Flags().StringVar(&submitArgs.WorkerImage, "workerImage", "", "the docker image for tensorflow workers")
//	command.Flags().MarkDeprecated("workerImage", "please use --worker-image instead")
//	command.Flags().StringVar(&submitArgs.WorkerImage, "worker-image", "", "the docker image for tensorflow workers")
//
//	command.Flags().StringVar(&submitArgs.PSImage, "psImage", "", "the docker image for tensorflow workers")
//	command.Flags().MarkDeprecated("psImage", "please use --ps-image instead")
//	command.Flags().StringVar(&submitArgs.PSImage, "ps-image", "", "the docker image for tensorflow workers")
//
//	command.Flags().IntVar(&submitArgs.PSCount, "ps", 0, "the number of the parameter servers.")
//
//	command.Flags().IntVar(&submitArgs.PSPort, "psPort", 0, "the port of the parameter server.")
//	command.Flags().MarkDeprecated("psPort", "please use --ps-port instead")
//	command.Flags().IntVar(&submitArgs.PSPort, "ps-port", 0, "the port of the parameter server.")
//
//	command.Flags().IntVar(&submitArgs.WorkerPort, "workerPort", 0, "the port of the worker.")
//	command.Flags().MarkDeprecated("workerPort", "please use --worker-port instead")
//	command.Flags().IntVar(&submitArgs.WorkerPort, "worker-port", 0, "the port of the worker.")
//
//	command.Flags().StringVar(&submitArgs.WorkerCpu, "workerCpu", "", "the cpu resource to use for the worker, like 1 for 1 core.")
//	command.Flags().MarkDeprecated("workerCpu", "please use --worker-cpu instead")
//	command.Flags().StringVar(&submitArgs.WorkerCpu, "worker-cpu", "", "the cpu resource to use for the worker, like 1 for 1 core.")
//
//	command.Flags().StringVar(&submitArgs.WorkerMemory, "workerMemory", "", "the memory resource to use for the worker, like 1Gi.")
//	command.Flags().MarkDeprecated("workerMemory", "please use --worker-memory instead")
//	command.Flags().StringVar(&submitArgs.WorkerMemory, "worker-memory", "", "the memory resource to use for the worker, like 1Gi.")
//
//	command.Flags().StringVar(&submitArgs.PSCpu, "psCpu", "", "the cpu resource to use for the parameter servers, like 1 for 1 core.")
//	command.Flags().MarkDeprecated("psCpu", "please use --ps-cpu instead")
//	command.Flags().StringVar(&submitArgs.PSCpu, "ps-cpu", "", "the cpu resource to use for the parameter servers, like 1 for 1 core.")
//
//	command.Flags().StringVar(&submitArgs.PSMemory, "psMemory", "", "the memory resource to use for the parameter servers, like 1Gi.")
//	command.Flags().MarkDeprecated("psMemory", "please use --ps-memory instead")
//	command.Flags().StringVar(&submitArgs.PSMemory, "ps-memory", "", "the memory resource to use for the parameter servers, like 1Gi.")
//	command.Flags().StringArrayVarP(&psSelectors, "ps-selector", "", []string{}, `assigning jobs with "PS" role to some k8s particular nodes(this option would cover --selector), usage: "--ps-selector=key=value"`)
//	// How to clean up Task
//	command.Flags().StringVar(&submitArgs.CleanPodPolicy, "cleanTaskPolicy", "Running", "How to clean tasks after Training is done, only support Running, None.")
//	command.Flags().MarkDeprecated("cleanTaskPolicy", "please use --clean-task-policy instead")
//	command.Flags().StringVar(&submitArgs.CleanPodPolicy, "clean-task-policy", "Running", "How to clean tasks after Training is done, only support Running, None.")
//
//	// Tensorboard
//	//command.Flags().BoolVar(&submitArgs.UseTensorboard, "tensorboard", false, "enable tensorboard")
//	//command.Flags().StringVar(&submitArgs.TensorboardImage, "tensorboardImage", "registry.cn-zhangjiakou.aliyuncs.com/tensorflow-samples/tensorflow:1.12.0-devel", "the docker image for tensorboard")
//	//command.Flags().MarkDeprecated("tensorboardImage", "please use --tensorboard-image instead")
//	//command.Flags().StringVar(&submitArgs.TensorboardImage, "tensorboard-image", "registry.cn-zhangjiakou.aliyuncs.com/tensorflow-samples/tensorflow:1.12.0-devel", "the docker image for tensorboard")
//
//	//command.Flags().StringVar(&submitArgs.TrainingLogdir, "logdir", "/training_logs", "the training logs dir, default is /training_logs")
//
//	// Estimator
//	command.Flags().BoolVar(&submitArgs.UseChief, "chief", false, "enable chief, which is required for estimator.")
//	command.Flags().BoolVar(&submitArgs.UseEvaluator, "evaluator", false, "enable evaluator, which is optional for estimator.")
//	command.Flags().StringVar(&submitArgs.ChiefCpu, "ChiefCpu", "", "the cpu resource to use for the Chief, like 1 for 1 core.")
//	command.Flags().MarkDeprecated("ChiefCpu", "please use --chief-cpu instead")
//	command.Flags().StringVar(&submitArgs.ChiefCpu, "chief-cpu", "", "the cpu resource to use for the Chief, like 1 for 1 core.")
//
//	command.Flags().StringVar(&submitArgs.ChiefMemory, "ChiefMemory", "", "the memory resource to use for the Chief, like 1Gi.")
//	command.Flags().MarkDeprecated("ChiefMemory", "please use --chief-memory instead")
//	command.Flags().StringVar(&submitArgs.ChiefMemory, "chief-memory", "", "the memory resource to use for the Chief, like 1Gi.")
//
//	command.Flags().StringVar(&submitArgs.EvaluatorCpu, "evaluatorCpu", "", "the cpu resource to use for the evaluator, like 1 for 1 core.")
//	command.Flags().MarkDeprecated("evaluatorCpu", "please use --evaluator-cpu instead")
//	command.Flags().StringVar(&submitArgs.EvaluatorCpu, "evaluator-cpu", "", "the cpu resource to use for the evaluator, like 1 for 1 core.")
//
//	command.Flags().StringVar(&submitArgs.EvaluatorMemory, "evaluatorMemory", "", "the memory resource to use for the evaluator, like 1Gi.")
//	command.Flags().MarkDeprecated("evaluatorMemory", "please use --evaluator-memory instead")
//	command.Flags().StringVar(&submitArgs.EvaluatorMemory, "evaluator-memory", "", "the memory resource to use for the evaluator, like 1Gi.")
//
//	command.Flags().IntVar(&submitArgs.ChiefPort, "chiefPort", 0, "the port of the chief.")
//	command.Flags().MarkDeprecated("chiefPort", "please use --chief-port instead")
//	command.Flags().IntVar(&submitArgs.ChiefPort, "chief-port", 0, "the port of the chief.")
//	command.Flags().StringArrayVarP(&workerSelectors, "worker-selector", "", []string{}, `assigning jobs with "Worker" role to some k8s particular nodes(this option would cover --selector), usage: "--worker-selector=key=value"`)
//	command.Flags().StringArrayVarP(&chiefSelectors, "chief-selector", "", []string{}, `assigning jobs with "Chief" role to some k8s particular nodes(this option would cover --selector), usage: "--chief-selector=key=value"`)
//	command.Flags().StringArrayVarP(&evaluatorSelectors, "evaluator-selector", "", []string{}, `assigning jobs with "Evaluator" role to some k8s particular nodes(this option would cover --selector), usage: "--evaluator-selector=key=value"`)
//
//	// command.Flags().BoolVarP(&showDetails, "details", "d", false, "Display details")
//	return command
//}

func NewSubmitTFJob(name string, namespace string, image string, req *service.JobDeploymentRequest, learnerBOM *v1.Job, logr *logger.LocLoggingEntry) error {
	const (
		PATH_SEPARATOR = ":"
		CODE_PATH      = "/training/code"
		DATA_PATH      = "/training/data"
		RESULT_PATH    = "/training/result"
		LOG_PATH       = "/training/log"

		MODEL_DIR_ENV = "MODEL_DIR"
		DATA_DIR_ENV  = "DATA_DIR"
		LOG_FILE_NAME = "training-log.txt"
	)

	var (
		submitArgs  submitTFJobArgs
		envs        []string
		selectors   []string
		tolerations []string
		dataset     []string
		dataDirs    []string
		annotations []string
	)

	numLearners := int(req.GetResources().Learners)
	//logr.Debugf("NewSubmitTFJob,ParameterServer : %v", req.ParameterServer)
	logr.Debugf("NewSubmitTFJob,ParameterServer pss: %v, psCpu: %v, psImage: %v, psMemory: %v", req.PSs, req.PSCPU, req.PSImage, req.PSMemory)
	logr.Debugf("NewSubmitTFJob,learners : %v", numLearners)
	logr.Debugf("NewSubmitTFJob,Framework : %v", req.Framework)
	logr.Debugf("NewSubmitTFJob,ImageTag : %v", req.ImageTag)

	err := fillSubmitArgs(&submitArgs, req, learnerBOM, logr)
	if err != nil {
		return err
	}

	//remove tab character from cmd
	container := learnerBOM.Spec.Template.Spec.Containers[0]
	cmdFromLearner := container.Command
	for i, c := range cmdFromLearner {
		replace := strings.Replace(c, "\t", "        ", -1)
		cmdFromLearner[i] = replace
	}

	varsFromlearner := container.Env

	//handle env
	for _, v := range varsFromlearner {
		envPair := fmt.Sprintf("%s=%s", v.Name, v.Value)
		envs = append(envs, envPair)
	}

	if submitArgs.GPUCount == 0 {
		envPair := fmt.Sprintf("%s=%s", "NVIDIA_VISIBLE_DEVICES", "")
		envs = append(envs, envPair)
	}

	//special env
	//export JOB_STATE_DIR=/job/${TRAINING_ID};
	//export LOG_DIR=/job/${TRAINING_ID}/logs;
	//export MODEL_DIR=/job/${TRAINING_ID}/model-code;
	//envs = append(envs, fmt.Sprintf("%s=%s%s", "JOB_STATE_DIR", "/job/", req.TrainingId))
	//envs = append(envs, fmt.Sprintf("%s=%s%s%s", "LOG_DIR", "/job/", req.TrainingId, "/logs"))
	//envs = append(envs, fmt.Sprintf("%s=%s%s%s", "MODEL_DIR", "/job/", req.TrainingId, "/model-code"))

	volumesFromBOM := learnerBOM.Spec.Template.Spec.Volumes
	volumeMountsFromBOM := learnerBOM.Spec.Template.Spec.Containers[0].VolumeMounts
	//var nodePaths []v1core.Volume

	//append HOSTPATH
	for _, v := range volumesFromBOM {
		if v.HostPath != nil && v.HostPath.Path != "" {
			vpath := v.HostPath.Path
			vname := v.Name

			for _, vm := range volumeMountsFromBOM {
				if vname == vm.Name {
					dirPair := fmt.Sprintf("%s:%s", vpath+"/"+vm.SubPath, vm.MountPath)
					dataDirs = append(dataDirs, dirPair)
				}
			}
			//nodePaths = append(nodePaths, path)
		}
	}

	//append PVC
	for _, v := range volumesFromBOM {
		if v.PersistentVolumeClaim != nil && v.PersistentVolumeClaim.ClaimName != "" {
			vClaimName := v.PersistentVolumeClaim.ClaimName
			vname := v.Name

			for _, vm := range volumeMountsFromBOM {
				if vname == vm.Name {
					dirPair := fmt.Sprintf("%s:%s", vClaimName+"/"+vm.SubPath, vm.MountPath)
					dataset = append(dataset, dirPair)
				}
			}
			//nodePaths = append(nodePaths, path)
		}
	}

	//dataDirs = append(dataDirs, CODE_PATH+PATH_SEPARATOR+contrainerWorkspacePath)
	//dataDirs = append(dataDirs, DATA_PATH+PATH_SEPARATOR+hostDataPath)
	//dataDirs = append(dataDirs, RESULT_PATH+PATH_SEPARATOR+contrainerResultsPath)

	//for k, v := range req.EnvVars {
	//	envPair := fmt.Sprintf("%s=%s", k, v)
	//}
	//nodeDataPath := req.EnvVars["DATA_STORE_PATH"]
	//nodeDataDir := req.EnvVars["DATA_DIR"]
	//nodeDataDir := req.EnvVars["MODEL_DIR"]
	//nodeDataDir := req.EnvVars["WORK_DIR"]
	//nodeDataDir := req.EnvVars["RESULT_DIR"]

	//envs = append(envs, MODEL_DIR_ENV+"="+CODE_PATH)
	//envs = append(envs, DATA_DIR_ENV+"="+DATA_PATH)

	//output log to file
	//cmdFromManifest := args[0]

	//logPath := fmt.Sprintf("%s/%s", LOG_PATH, LOG_FILE_NAME)
	//while training finished, copy log to result folder
	//args[0] = fmt.Sprintf("%s |tee %s; cp %s %s/%s/%s", cmdFromManifest, logPath, logPath, RESULT_PATH, req.TrainingId, LOG_FILE_NAME)

	var command = &cobra.Command{
		Use:     "tfjob",
		Short:   "Submit TFJob as training job.",
		Aliases: []string{"tf"},
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	submitArgs.addCommonFlags(command)
	//submitArgs.addSyncFlags(command)
	// TFJob
	//command.Flags().StringVar(&submitArgs.WorkerImage, "worker-image", "", "the docker image for tensorflow workers")

	//command.Flags().StringVar(&submitArgs.PSImage, "ps-image", "", "the docker image for tensorflow workers")

	//command.Flags().IntVar(&submitArgs.PSCount, "ps", 0, "the number of the parameter servers.")

	command.Flags().IntVar(&submitArgs.PSPort, "psPort", 0, "the port of the parameter server.")
	command.Flags().IntVar(&submitArgs.PSPort, "ps-port", 0, "the port of the parameter server.")

	command.Flags().IntVar(&submitArgs.WorkerPort, "workerPort", 0, "the port of the worker.")
	command.Flags().IntVar(&submitArgs.WorkerPort, "worker-port", 0, "the port of the worker.")

	//command.Flags().StringVar(&submitArgs.WorkerCpu, "workerCpu", "1", "the cpu resource to use for the worker, like 1 for 1 core.")
	//command.Flags().StringVar(&submitArgs.WorkerCpu, "worker-cpu", "1", "the cpu resource to use for the worker, like 1 for 1 core.")

	//command.Flags().StringVar(&submitArgs.WorkerMemory, "workerMemory", "1Gi", "the memory resource to use for the worker, like 1Gi.")
	//command.Flags().StringVar(&submitArgs.WorkerMemory, "worker-memory", "1Gi", "the memory resource to use for the worker, like 1Gi.")

	command.Flags().StringVar(&submitArgs.PSCpu, "psCpu", "1", "the cpu resource to use for the parameter servers, like 1 for 1 core.")
	command.Flags().StringVar(&submitArgs.PSCpu, "ps-cpu", "1", "the cpu resource to use for the parameter servers, like 1 for 1 core.")

	command.Flags().StringVar(&submitArgs.PSMemory, "psMemory", "1Gi", "the memory resource to use for the parameter servers, like 1Gi.")
	command.Flags().StringVar(&submitArgs.PSMemory, "ps-memory", "1Gi", "the memory resource to use for the parameter servers, like 1Gi.")
	command.Flags().StringArrayVarP(&psSelectors, "ps-selector", "", []string{}, `assigning jobs with "PS" role to some k8s particular nodes(this option would cover --selector), usage: "--ps-selector=key=value"`)
	// How to clean up Task
	command.Flags().StringVar(&submitArgs.CleanPodPolicy, "cleanTaskPolicy", "Running", "How to clean tasks after Training is done, only support Running, None.")
	command.Flags().StringVar(&submitArgs.CleanPodPolicy, "clean-task-policy", "Running", "How to clean tasks after Training is done, only support Running, None.")

	// Tensorboard
	command.Flags().BoolVar(&submitArgs.UseTensorboard, "tensorboard", false, "enable tensorboard")
	command.Flags().StringVar(&submitArgs.TensorboardImage, "tensorboardImage", "registry.cn-zhangjiakou.aliyuncs.com/tensorflow-samples/tensorflow:1.12.0-devel", "the docker image for tensorboard")
	command.Flags().StringVar(&submitArgs.TensorboardImage, "tensorboard-image", "registry.cn-zhangjiakou.aliyuncs.com/tensorflow-samples/tensorflow:1.12.0-devel", "the docker image for tensorboard")

	command.Flags().StringVar(&submitArgs.TrainingLogdir, "logdir", "/training_logs", "the training logs dir, default is /training_logs")

	// Estimator
	command.Flags().BoolVar(&submitArgs.UseChief, "chief", false, "enable chief, which is required for estimator.")
	command.Flags().BoolVar(&submitArgs.UseEvaluator, "evaluator", false, "enable evaluator, which is optional for estimator.")
	command.Flags().StringVar(&submitArgs.ChiefCpu, "ChiefCpu", "1", "the cpu resource to use for the Chief, like 1 for 1 core.")
	command.Flags().StringVar(&submitArgs.ChiefCpu, "chief-cpu", "1", "the cpu resource to use for the Chief, like 1 for 1 core.")

	command.Flags().StringVar(&submitArgs.ChiefMemory, "ChiefMemory", "1Gi", "the memory resource to use for the Chief, like 1Gi.")
	command.Flags().StringVar(&submitArgs.ChiefMemory, "chief-memory", "1Gi", "the memory resource to use for the Chief, like 1Gi.")

	command.Flags().StringVar(&submitArgs.EvaluatorCpu, "evaluatorCpu", "1", "the cpu resource to use for the evaluator, like 1 for 1 core.")
	command.Flags().StringVar(&submitArgs.EvaluatorCpu, "evaluator-cpu", "1", "the cpu resource to use for the evaluator, like 1 for 1 core.")

	command.Flags().StringVar(&submitArgs.EvaluatorMemory, "evaluatorMemory", "1Gi", "the memory resource to use for the evaluator, like 1Gi.")
	command.Flags().StringVar(&submitArgs.EvaluatorMemory, "evaluator-memory", "1Gi", "the memory resource to use for the evaluator, like 1Gi.")

	command.Flags().IntVar(&submitArgs.ChiefPort, "chiefPort", 0, "the port of the chief.")
	command.Flags().IntVar(&submitArgs.ChiefPort, "chief-port", 0, "the port of the chief.")
	command.Flags().StringArrayVarP(&workerSelectors, "worker-selector", "", []string{}, `assigning jobs with "Worker" role to some k8s particular nodes(this option would cover --selector), usage: "--worker-selector=key=value"`)
	command.Flags().StringArrayVarP(&chiefSelectors, "chief-selector", "", []string{}, `assigning jobs with "Chief" role to some k8s particular nodes(this option would cover --selector), usage: "--chief-selector=key=value"`)
	command.Flags().StringArrayVarP(&evaluatorSelectors, "evaluator-selector", "", []string{}, `assigning jobs with "Evaluator" role to some k8s particular nodes(this option would cover --selector), usage: "--evaluator-selector=key=value"`)

	//worker image
	//submitArgs.Image = learnerBOM.Spec.Template.Spec.Containers[0].Image

	client, err := InitKubeClient(lcmconfig.GetKubernetesConfig())
	if err != nil {
		log.Debugf("Failed due to %v", err)
		fmt.Println(err)
		os.Exit(1)
	}

	//logr.Debugf("NewSubmitTFJob, submitArgs.Image: %v", submitArgs.Image)
	//logr.Debugf("NewSubmitTFJob, submitArgs.WorkerImage: %v", submitArgs.WorkerImage)
	logr.Debugf("NewSubmitTFJob, args: %+v", submitArgs)

	logr.Debugf("NewSubmitTFJob, dataDirs: %v", dataDirs)
	logr.Debugf("NewSubmitTFJob, envs: %v", envs)

	err = submitTFJob(name, namespace, cmdFromLearner, &submitArgs, client, selectors, tolerations, dataset, dataDirs, annotations, envs, req)
	if err != nil {
		log.Debugf("Failed due to %v", err)
		fmt.Println(err)
		return errors.New("Failed due to" + err.Error())
	}
	return nil
}

func fillSubmitArgs(args *submitTFJobArgs, req *service.JobDeploymentRequest, learnerBOM *v1.Job, logr *logger.LocLoggingEntry) error {
	args.Mode = TFJob
	args.Image = learnerBOM.Spec.Template.Spec.Containers[0].Image

	numLearners := int(req.GetResources().Learners)
	args.WorkerCount = numLearners
	args.GPUCount = int(req.Resources.Gpus)

	//gpuStr := learnerBOM.Spec.Template.Spec.Containers[0].Resources.Limits.NvidiaGPU().String()
	//logr.Debugf("fillSubmitArgs, gpuStr: %s", gpuStr)
	//gpu, err := strconv.Atoi(gpuStr)
	//if err != nil {
	//	logr.Errorf("parse gpuStr to int failed")
	//	return err
	//}
	//args.GPUCount = gpu

	//args.WorkerCpu = strconv.FormatFloat(learnerBOM.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().Value(), 'E', -1, 64)
	//args.WorkerMemory = strconv.FormatFloat(float64(req.Training.Resources.Memory), 'E', -1, 32)
	args.WorkerCpu = learnerBOM.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().String()
	args.WorkerMemory = learnerBOM.Spec.Template.Spec.Containers[0].Resources.Limits.Memory().String()

	args.WorkerImage = learnerBOM.Spec.Template.Spec.Containers[0].Image

	//psObjString := req.ParameterServer
	//
	//psMap := make(map[string]string)
	//if err := json.Unmarshal([]byte(psObjString), &psMap); err != nil {
	//	logr.Debugf("parse ps to map failed")
	//	return err
	//}
	logr.Debugf("fillSubmitArgs, pss: %v", req.PSs)
	logr.Debugf("fillSubmitArgs, ps_cpu: %v", req.PSCPU)
	logr.Debugf("fillSubmitArgs, ps_memory: %v", req.PSMemory)
	logr.Debugf("fillSubmitArgs, ps_image: %v", req.PSImage)

	//pss, err := strconv.Atoi(psMap["pss"])
	//if err != nil {
	//	logr.Debugf("parse pss to int failed")
	//	return err
	//}
	pss := 1
	if "" != req.PSs {
		s, pssErr := strconv.Atoi(req.PSs)
		if nil != pssErr {
			logr.Debugf("strconv.ParseFloat param: %v", req.PSs)
			return pssErr
		}
		pss = s
	}
	args.PSCount = pss
	logr.Debugf("args.PSCount: %v", args.PSCount)
	//args.PSCpu = psMap["ps_cpu"]
	//args.PSMemory = psMap["ps_memory"]
	psImage := fmt.Sprintf("%s:%s", req.Framework, req.PSImage)
	args.PSImage = psImage

	return nil
}

type submitTFJobArgs struct {
	TFNodeSelectors map[string]map[string]string `yaml:"tfNodeSelectors"`
	Port            int                          // --port, it's used set workerPort and PSPort if they are not set
	WorkerImage     string                       `yaml:"workerImage"` // --workerImage
	WorkerPort      int                          `yaml:"workerPort"`  // --workerPort
	PSPort          int                          `yaml:"psPort"`      // --psPort
	//PSNodeSelectors map[string]string `yaml:"psNodeSelectors"` // --ps-selector
	PSCount   int    `yaml:"ps"`        // --ps
	PSImage   string `yaml:"psImage"`   // --psImage
	WorkerCpu string `yaml:"workerCPU"` // --workerCpu
	//WorkerNodeSelectors map[string]string `yaml:"workerNodeSelectors"` // --worker-selector
	WorkerMemory   string `yaml:"workerMemory"`   // --workerMemory
	PSCpu          string `yaml:"psCPU"`          // --psCpu
	PSMemory       string `yaml:"psMemory"`       // --psMemory
	CleanPodPolicy string `yaml:"cleanPodPolicy"` // --cleanTaskPolicy
	// For esitmator, it reuses workerImage
	UseChief     bool `yaml:",omitempty"` // --chief
	ChiefCount   int  `yaml:"chief"`
	UseEvaluator bool `yaml:",omitempty"` // --evaluator
	ChiefPort    int  `yaml:"chiefPort"`  // --chiefPort
	//ChiefNodeSelectors map[string]string `yaml:"chiefNodeSelectors"` // --chief-selector
	ChiefCpu     string `yaml:"chiefCPU"`     // --chiefCpu
	ChiefMemory  string `yaml:"chiefMemory"`  // --chiefMemory
	EvaluatorCpu string `yaml:"evaluatorCPU"` // --evaluatorCpu
	//EvaluatorNodeSelectors map[string]string `yaml:"evaluatorNodeSelectors"` // --evaluator-selector
	EvaluatorMemory string `yaml:"evaluatorMemory"` // --evaluatorMemory
	EvaluatorCount  int    `yaml:"evaluator"`

	// determine if it has gang scheduler
	HasGangScheduler bool `yaml:"hasGangScheduler"`

	// for common args
	submitArgs `yaml:",inline"`

	// for tensorboard
	submitTensorboardArgs `yaml:",inline"`

	// for sync up source code
	submitSyncCodeArgs `yaml:",inline"`
}

func (submitArgs *submitTFJobArgs) prepare(name string, args []string, clientset *kubernetes.Clientset, selectors []string, tolerations []string, dataset []string, dataDirs []string, annotations []string, envs []string) (err error) {
	//FIXME: omit args , insert cmd later
	submitArgs.Command = strings.Join(args, " ")
	//submitArgs.Command = ""

	err = submitArgs.transform(clientset)
	if err != nil {
		return err
	}

	err = submitArgs.check(name)
	if err != nil {
		return err
	}

	//err = submitArgs.HandleSyncCode()
	//if err != nil {
	//	return err
	//}

	commonArgs := &submitArgs.submitArgs
	err = commonArgs.transform(selectors, tolerations, dataset, dataDirs, annotations)
	if err != nil {
		return nil
	}

	// process tensorboard
	//submitArgs.processTensorboard(submitArgs.DataSet)

	if len(envs) > 0 {
		submitArgs.Envs = transformSliceToMap(envs, "=")
	}
	// pass the workers, gpu to environment variables
	// addTFJobInfoToEnv(submitArgs)
	submitArgs.addTFJobInfoToEnv()
	// add node selectors, if given
	submitArgs.addTFNodeSelectors(selectors)
	// add tolerations, if given`
	submitArgs.addTFTolerations(tolerations)
	return nil
}

func (submitArgs submitTFJobArgs) check(name string) error {
	err := submitArgs.submitArgs.check(name)
	if err != nil {
		return err
	}

	switch submitArgs.CleanPodPolicy {
	case "None", "Running":
		log.Debugf("Supported cleanTaskPolicy: %s", submitArgs.CleanPodPolicy)
	default:
		return fmt.Errorf("Unsupported cleanTaskPolicy %s", submitArgs.CleanPodPolicy)
	}

	if submitArgs.WorkerCount == 0 && !submitArgs.UseChief {
		return fmt.Errorf("--workers must be greater than 0 in distributed training")
	}

	if submitArgs.WorkerImage == "" {
		return fmt.Errorf("--image or --workerImage must be set")
	}

	if submitArgs.PSCount > 0 {
		if submitArgs.PSImage == "" {
			return fmt.Errorf("--image or --psImage must be set")
		}
	}

	return nil
}

// This method for supporting tf-estimator
func (submitArgs *submitTFJobArgs) setStandaloneMode() {
	if submitArgs.PSCount < 1 && submitArgs.WorkerCount == 1 {
		submitArgs.UseChief = true
		submitArgs.WorkerCount = 0
	}
}

func (submitArgs *submitTFJobArgs) transform(clientset *kubernetes.Clientset) error {

	submitArgs.setStandaloneMode()

	if submitArgs.WorkerImage == "" {
		submitArgs.WorkerImage = submitArgs.Image
	}

	if submitArgs.WorkerCount > 0 {
		autoSelectWorkerPort, err := util.SelectAvailablePortWithDefault(clientset, submitArgs.WorkerPort)
		if err != nil {
			return fmt.Errorf("failed to select worker port: %++v", err)
		}
		submitArgs.WorkerPort = autoSelectWorkerPort
	}

	if submitArgs.UseChief {
		autoSelectChiefPort, err := util.SelectAvailablePortWithDefault(clientset, submitArgs.ChiefPort)
		if err != nil {
			return fmt.Errorf("failed to select chief port: %++v", err)
		}
		submitArgs.ChiefPort = autoSelectChiefPort
		submitArgs.ChiefCount = 1
	}

	if submitArgs.PSCount > 0 {
		autoSelectPsPort, err := util.SelectAvailablePortWithDefault(clientset, submitArgs.PSPort)
		if err != nil {
			return fmt.Errorf("failed to select ps port: %++v", err)
		}
		submitArgs.PSPort = autoSelectPsPort
		if submitArgs.PSImage == "" {
			submitArgs.PSImage = submitArgs.Image
		}
	}

	if submitArgs.UseEvaluator {
		submitArgs.EvaluatorCount = 1
	}

	// check Gang scheduler
	submitArgs.checkGangCapablitiesInCluster()

	return nil
}

// add node selectors
func (submitArgs *submitTFJobArgs) addTFNodeSelectors(selectors []string) {
	submitArgs.addNodeSelectors(selectors)
	submitArgs.TFNodeSelectors = make(map[string]map[string]string)
	for _, role := range []string{"PS", "Worker", "Evaluator", "Chief"} {
		switch role {
		case "PS":
			log.Debugf("psSelectors: %v", psSelectors)
			submitArgs.transformSelectorArrayToMap(psSelectors, "PS")
			break
		case "Worker":
			log.Debugf("workerSelectors: %v", workerSelectors)
			submitArgs.transformSelectorArrayToMap(workerSelectors, "Worker")
			break
		case "Chief":
			log.Debugf("chiefSelectors: %v", chiefSelectors)
			submitArgs.transformSelectorArrayToMap(chiefSelectors, "Chief")
			break
		case "Evaluator":
			log.Debugf("evaluatorSelectors: %v", evaluatorSelectors)
			submitArgs.transformSelectorArrayToMap(evaluatorSelectors, "Evaluator")
			break
		}

	}
}

func (submitArgs *submitTFJobArgs) transformSelectorArrayToMap(selectorArray []string, role string) {
	if len(selectorArray) != 0 {
		submitArgs.TFNodeSelectors[role] = transformSliceToMap(selectorArray, "=")
		return
	}
	if len(submitArgs.NodeSelectors) == 0 {
		submitArgs.TFNodeSelectors[role] = map[string]string{}
		return
	}
	submitArgs.TFNodeSelectors[role] = submitArgs.NodeSelectors

}

// add tolerations
func (submitArgs *submitTFJobArgs) addTFTolerations(tolerations []string) {
	submitArgs.addTolerations(tolerations)
}
func (submitArgs *submitTFJobArgs) addTFJobInfoToEnv() {
	submitArgs.addJobInfoToEnv()
}

func (submitArgs *submitTFJobArgs) checkGangCapablitiesInCluster() {
	gangCapablity := false
	if clientset != nil {
		_, err := clientset.AppsV1beta1().Deployments(metav1.NamespaceSystem).Get(gangSchdName, metav1.GetOptions{})
		if err != nil {
			log.Debugf("Failed to find %s due to %v", gangSchdName, err)
		} else {
			log.Debugf("Found %s successfully, the gang scheduler is enabled in the cluster.", gangSchdName)
			gangCapablity = true
		}
	}

	submitArgs.HasGangScheduler = gangCapablity
}

func submitTFJob(name string, namespace string, args []string, submitArgs *submitTFJobArgs, clientset *kubernetes.Clientset, selectors []string, tolerations []string, dataset []string, dataDirs []string, annotations []string, envs []string, req *service.JobDeploymentRequest) (err error) {

	//log.Infof("submitTFJob, args[2]: %v", args[2])

	err = submitArgs.prepare(name, args, clientset, selectors, tolerations, dataset, dataDirs, annotations, envs)
	if err != nil {
		return err
	}
	log.Infof("submitTFJob, submitArgs.ChiefCpu: %v, submitArgs.ChiefMemory: %v", submitArgs.ChiefCpu, submitArgs.ChiefMemory)

	//trainer := NewTensorFlowJobTrainer(clientset)
	//job, err := trainer.GetTrainingJob(name, namespace)
	if err != nil {
		log.Debugf("Check %s exist due to error %v", name, err)
	}

	//if job != nil {
	//	return fmt.Errorf("the job %s is already exist, please delete it first. use 'arena delete %s'", name, name)
	//}

	// the master is also considered as a worker
	// submitArgs.WorkerCount = submitArgs.WorkerCount - 1

	err = workflow.SubmitJob(name, submitArgs.Mode, namespace, submitArgs, tfjob_chart, args[2], req)
	if err != nil {
		return err
	}

	log.Infof("The Job %s has been submitted successfully", name)
	log.Infof("You can run `arena get %s --type %s` to check the job status", name, submitArgs.Mode)
	return nil
}
