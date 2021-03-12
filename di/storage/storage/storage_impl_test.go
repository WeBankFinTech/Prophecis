package storage

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/minio/minio-go/v6"
	"github.com/prometheus/common/log"
	"testing"
	"webank/DI/commons/logger"
	"webank/DI/storage/storage/grpc_storage"
	"webank/DI/trainer/trainer/grpc_trainer_v2"
)

func Test(t *testing.T) {
	var dataStores []*grpc_trainer_v2.Datastore

	datastores := []*grpc_storage.Datastore{
		{
			Type: bdap,
		},
		{
			Type: bdapsafe,
		},
	}

	err := copier.Copy(&dataStores, &datastores)
	if err != nil {
		log.Debugf("Copy Datastore failed: %v", err.Error())
	}
	log.Infof("ds: %+v", dataStores)
}

func TestUpload(t *testing.T) {
	getLogger := logger.GetLogger()
	//endpoint := "mlss-minio.mlss-dev.svc.cluster.local:9000"
	endpoint := ""
	accessKeyID := ""
	secretAccessKey := ""
	realBucket := "test"

	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, false)
	if err != nil {
		getLogger.Errorf("minio.New err: %v", err.Error())
		return
	}

	ctx, _ := context.WithCancel(context.TODO())
	//当为hostPath时
	//创建 bucket
	err = minioClient.MakeBucket(realBucket, "cn-north-1")
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(realBucket)
		if errBucketExists == nil && exists {
			getLogger.Debugf("We already own %s", realBucket)
		} else {
			getLogger.Errorf("minioClient.MakeBucket failed, %v", err.Error())
			return
		}
	} else {
		getLogger.Debugf("Successfully created %s", realBucket)
	}
	//put obj
	_, err = minioClient.FPutObjectWithContext(ctx, realBucket, "aha/test.png", "ar.png", minio.PutObjectOptions{})
	if err != nil {
		getLogger.Errorf("minioClient.FPutObjectWithContext failed, err: %v\n", err.Error())
		return
	}
	getLogger.Debugf("successful uploaded %v", "test.png")
	if err != nil {
		getLogger.Errorf("stream.SendAndClose failed, err: %v", err.Error())
		return
	}
	return
}
