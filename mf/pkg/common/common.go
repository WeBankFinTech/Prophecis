/*
 * Copyright 2020 WeBank
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
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

package common

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	Restconfig   *rest.Config
	ClientConfig clientcmd.ClientConfig
	Clientset    *kubernetes.Clientset
)

func init() {
	Clientset, Restconfig, _ = initKubeClient()
	//clientConfig =
}

func initKubeClient() (*kubernetes.Clientset, *rest.Config, error) {
	//-----------------Debug Code--------------------------------------------------------------------------------------
	//local := true
	local := false
	if local {
		//kubeconfig := flag.String("kubeconfig", "D://kube", "absoulute path to the kubeconfig file")
		flag.Parse()
		config, err := clientcmd.BuildConfigFromFlags("", "D://kube/config")
		if err != nil {
			panic(err.Error())
		}
		Restconfig = config
		Clientset, err = kubernetes.NewForConfig(config)
		return Clientset, config, err
	}
	//-----------------Debug Code--------------------------------------------------------------------------------------

	if Clientset != nil && Restconfig != nil {
		return Clientset, Restconfig, nil
	}
	restConfig, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}
	// create the clientset
	Clientset, err = kubernetes.NewForConfig(restConfig)
	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}
	return Clientset, restConfig, nil
}

//
//func GetKubernetesConfig() *rest.Config{
//
//	host := config.GetMFKubeURL()
//	var c *rest.Config
//	if host == "" {
//		c, _ = rest.InClusterConfig()
//	} else {
//		c = &rest.Config{
//			Host: host,
//			TLSClientConfig: rest.TLSClientConfig{
//				CAFile: config.GetLearnerKubeCAFile(),
//			},
//		}
//		token := config.GetLearnerKubeToken()
//		if token == "" {
//			tokenFileContents := config.GetFileContents(config.GetLearnerKubeTokenFile())
//			if tokenFileContents != "" {
//				token = tokenFileContents
//			}
//		}
//		if token == "" {
//			c.TLSClientConfig.KeyFile = config.GetLearnerKubeKeyFile()
//			c.TLSClientConfig.CertFile = config.GetLearnerKubeCertFile()
//		} else {
//			c.BearerToken = token
//		}
//	}
//	return c
//
//}
