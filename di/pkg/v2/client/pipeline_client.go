package client

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
	"webank/DI/commons/logger"
	"webank/DI/pkg/v2/util"

	api "github.com/kubeflow/pipelines/backend/api/v2beta1/go_client"
	"github.com/pkg/errors"
)

const (
	addressTemp = "%s:%s"
)

var log = logger.GetLogger()

type KFPClient struct {
	PipelineServiceClient api.PipelineServiceClient
	ReportServiceClient   api.ReportServiceClient
	RunServiceClient      api.RunServiceClient
}

func NewKFPClient(
	initializeTimeout time.Duration,
	basePath string,
	mlPipelineServiceName string,
	mlPipelineServiceHttpPort string,
	mlPipelineServiceGRPCPort string) (*KFPClient, error) {
	httpAddress := fmt.Sprintf(addressTemp, mlPipelineServiceName, mlPipelineServiceHttpPort)
	grpcAddress := fmt.Sprintf(addressTemp, mlPipelineServiceName, mlPipelineServiceGRPCPort)
	err := util.WaitForAPIAvailable(initializeTimeout, basePath, httpAddress)
	if err != nil {
		return nil, errors.Wrapf(err,
			"Failed to initialize pipeline client. Error: %s", err.Error())
	}
	connection, err := util.GetRpcConnection(grpcAddress)
	if err != nil {
		return nil, errors.Wrapf(err,
			"Failed to get RPC connection. Error: %s", err.Error())
	}

	return &KFPClient{
		PipelineServiceClient: api.NewPipelineServiceClient(connection),
		ReportServiceClient:   api.NewReportServiceClient(connection),
		RunServiceClient:      api.NewRunServiceClient(connection),
	}, nil
}

func GetKFPClient() (*KFPClient, error) {
	var (
		initializeTimeout           time.Duration
		mlPipelineAPIServerName     string
		mlPipelineAPIServerBasePath string
		mlPipelineServiceHttpPort   string
		mlPipelineServiceGRPCPort   string
	)
	kfpServerName := GetStringConfigWithDefault("kfp.address", "ml-pipeline.kfp.svc.cluster.local")
	kfpServerHttpPort := GetStringConfigWithDefault("kfp.httpPort", "8888")
	kfpServerGrpcPort := GetStringConfigWithDefault("kfp.grpcPort", "8887")
	kfpServerBasePath := GetStringConfigWithDefault("kfp.basePath", "/apis/v2beta1")

	log.Infof("kfpConfig kfpServerName is: %s; kfpServerHttpPort is: %s; kfpServerGrpcPort is: %s; kfpServerBasePath is: %s",
		kfpServerName, kfpServerHttpPort, kfpServerGrpcPort, kfpServerBasePath)

	initializeTimeout = 2 * time.Minute
	mlPipelineAPIServerName = kfpServerName
	mlPipelineServiceHttpPort = kfpServerHttpPort
	mlPipelineServiceGRPCPort = kfpServerGrpcPort
	mlPipelineAPIServerBasePath = kfpServerBasePath

	kfpClient, err := NewKFPClient(
		initializeTimeout,
		mlPipelineAPIServerBasePath,
		mlPipelineAPIServerName,
		mlPipelineServiceHttpPort,
		mlPipelineServiceGRPCPort)
	if err != nil {
		log.Errorf("Error creating ML pipeline API Server client: %v", err)
	}
	return kfpClient, err
}

func GetStringConfigWithDefault(configName, value string) string {
	if !viper.IsSet(configName) {
		return value
	}
	return viper.GetString(configName)
}
