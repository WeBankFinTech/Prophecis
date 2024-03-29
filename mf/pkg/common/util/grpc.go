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

package util

import (
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"mlss-mf/pkg/common/config"
)

// CreateClientDialOpts creates the TLC/non-TLS and other common dial options for
// establishing a grpc server connection to other microservices.
func CreateClientDialOpts() ([]grpc.DialOption, error) {
	var opts []grpc.DialOption
	if config.IsStorerTLSEnabled() {
		certBytes := []byte(config.GetStorerCAKey())
		cp := x509.NewCertPool()
		if !cp.AppendCertsFromPEM(certBytes) {
			return nil, fmt.Errorf("credentials: failed to append certificates")
		}
		creds := credentials.NewClientTLSFromCert(cp, config.GetStorerServerName())
		opts = []grpc.DialOption{grpc.WithTransportCredentials(creds), grpc.WithBlock()}
	} else {
		opts = []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}
	}
	return opts, nil
}
