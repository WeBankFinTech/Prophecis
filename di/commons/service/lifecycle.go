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

package service

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"

	"google.golang.org/grpc"
)

// LifecycleHandler provides basic lifecycle methods that each microservice has
// to implement.
type LifecycleHandler interface {
	Start(port int, background bool)
	Stop()
	GetListenerAddress() string
}

// Lifecycle implements the lifecycle operations for microservice including
// dynamic service registration.
type Lifecycle struct {
	Listener        net.Listener
	Server          *grpc.Server
	RegisterService func()
}

type Config struct {
	port       int
	background bool
	tls        bool
	certFile   string
	keyFile    string
}

// Start will start a gRPC microservice on a given port and run it either in
// foreground or background.
func (s *Lifecycle) Start(port int, background bool) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Infof("Starting service at %s", lis.Addr())

	var opts []grpc.ServerOption
	//if config.IsTLSEnabled() {
	//	config.FatalOnAbsentKey(config.ServerCertKey)
	//	config.FatalOnAbsentKey(config.ServerPrivateKey)
	//
	//	creds, err := credentials.NewServerTLSFromFile(config.GetServerCert(), config.GetServerPrivateKey())
	//	if err != nil {
	//		log.Fatalf("Failed to generate credentials %v", err)
	//	}
	//	opts = []grpc.ServerOption{grpc.Creds(creds)}
	//}

	//var kaep = keepalive.EnforcementPolicy{
	//	MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
	//	PermitWithoutStream: true,            // Allow pings even when there are no active streams
	//}
	//
	//var kasp = keepalive.ServerParameters{
	//	MaxConnectionIdle:     5 * time.Minute, // If a client is idle for 15 seconds, send a GOAWAY
	//	MaxConnectionAge:      30 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
	//	MaxConnectionAgeGrace: 5 * time.Second,  // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
	//	Time:                  5 * time.Second,  // Ping the client if it is idle for 5 seconds to ensure the connection is still active
	//	Timeout:               1 * time.Second,  // Wait 1 second for the ping ack before assuming the connection is dead
	//}

	//opts = append(opts, grpc.KeepaliveParams(keepalive.ServerParameters{MaxConnectionIdle: 5 * time.Minute}))
	//opts = append(opts, grpc.KeepaliveParams(kasp), grpc.KeepaliveEnforcementPolicy(kaep))

	s.Listener = lis
	s.Server = grpc.NewServer(opts...)

	s.RegisterService()
	grpc_health_v1.RegisterHealthServer(s.Server, health.NewServer())

	if background {
		log.Info("ruuning server in background")
		go s.Server.Serve(lis)
	} else {
		log.Info("running server in foreground")
		if err := s.Server.Serve(lis); err!= nil{
			log.Fatalf("Failed to serve grpc server: %v", err)
		}
	}
}

// Stop will stop the gRPC microservice and the socket.
func (s *Lifecycle) Stop() {
	if s.Listener != nil {
		log.Infof("Stopping service at %s", s.Listener.Addr())
	}
	if s.Server != nil {
		s.Server.GracefulStop()
	}
}

// GetListenerAddress will get the address and port the service is listening.
// Returns the empty string if the service is not running but the method is invoked.
func (s *Lifecycle) GetListenerAddress() string {
	if s.Listener != nil {
		return s.Listener.Addr().String()
	}
	return ""
}
