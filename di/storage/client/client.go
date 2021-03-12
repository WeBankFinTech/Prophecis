/*
 * Copyright 2017-2018 IBM Corporation
 *
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
 */

package client

import (
	"fmt"
	"webank/DI/commons/config"
	"webank/DI/commons/util"
	"webank/DI/storage/storage/grpc_storage"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

const (
	disabled = "disabled"
	// TrainerV2LocalAddress exposes the local address that is used if we run with DNS disabled
	StorageV2LocalAddress = ":30005"
)

// TrainerClient is a client interface for interacting with the trainer service.
type StorageClient interface {
	Client() grpc_storage.StorageClient
	Close() error
}

type storageClient struct {
	client grpc_storage.StorageClient
	conn   *grpc.ClientConn
}

// NewTrainer create a new load-balanced client to talk to the Trainer
// service. If the dns_server config option is set to 'disabled', it will
// default to the pre-defined LocalPort of the service.
func NewStorage() (StorageClient, error) {
	return NewStorageWithAddress(StorageV2LocalAddress)
}

// NewTrainerWithAddress create a new load-balanced client to talk to the Trainer
// service. If the dns_server config option is set to 'disabled', it will
// default to the pre-defined LocalPort of the service.
func NewStorageWithAddress(addr string) (StorageClient, error) {
	//address := fmt.Sprintf("ffdl-trainer.%s.svc.cluster.local:80", config.GetPodNamespace())
	// FIXME MLSS Change: specify namespace
	//	address := fmt.Sprintf("ffdl-trainer.%s.svc.cluster.local:80", "default")
	address := fmt.Sprintf("di-storage-rpc.%s.svc.cluster.local:80", viper.GetString(config.PlatformNamespace))
	dnsServer := viper.GetString("dns_server")
	if dnsServer == disabled { // for local testing without DNS server
		address = addr
	}

	dialOpts, err := util.CreateClientDialOpts()
	if err != nil {
		return nil, err
	}
	conn, err := grpc.Dial(address, dialOpts...)
	if err != nil {
		log.Errorf("Could not connect to storage service: %v", err)
		return nil, err
	}

	return &storageClient{
		conn:   conn,
		client: grpc_storage.NewStorageClient(conn),
	}, nil
}

func (c *storageClient) Client() grpc_storage.StorageClient {
	return c.client
}

func (c *storageClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
