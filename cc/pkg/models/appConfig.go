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
package models

type AppConfig struct {
	Core struct {
		Interceptor Interceptor `yaml:"interceptor"`
		Um          UM          `yaml:"um"`
		Sso         SSO         `yaml:"sso"`
		Ims         IMS         `yaml:"ims"`
		Kube        Kube        `yaml:"kube"`
		//Gateway     Gateway     `yaml:"gateway"`
		AuthAddress AuthAddress `yaml:"authAddress"`
		Cache       Cache       `yaml:"cache"`
		Cookie      Cookie      `yaml:"cookie"`
	}
	Server Server `yaml:"server"`
	Application Application `yaml:"application"`
}

type Application struct {
	Profile    string `yaml:"profile"`
	PlatformNS string `yaml:"platformNS"`
	Datasource struct {
		Url      string `yaml:"url"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Ip       string `yaml:"ip"`
		Port     string `yaml:"port"`
		Db       string `yaml:"db"`
	}
	LDAP string  `yaml:"ldap"`
}

type Server struct {
	Port    int `yaml:"port"`
	Servlet struct {
		ContextPath string `yaml:"context-path"`
	}
}

// Yaml1 struct of yaml
type Interceptor struct {
	Configs                 []InterceptorConfig `yaml:"configs,flow"`
	MaxTokenSize            string              `yaml:"maxTokenSize"`
	MaxTokenExpireHour      string              `yaml:"maxTokenExpireHour"`
	DefaultTimestampTimeout string              `yaml:"defaultTimestampTimeout"`
}

// Yaml2 struct of yaml
//type Yaml2 struct {
//	Mysql `yaml:"mysql,inline"`
//	authcache `yaml:"authcache,inline"`
//}

// Mysql struct of mysql conf
//type Mysql struct {
//	User string `yaml:"user"`
//	Host string `yaml:"host"`
//	Password string `yaml:"password"`
//	Port string `yaml:"port"`
//	Name string `yaml:"name"`
//}

// authcache struct of authcache conf
//type authcache struct {
//	Enable bool `yaml:"enable"`
//	List []string `yaml:"list,flow"`
//}
type InterceptorConfig struct {
	Name    string   `yaml:"name"`
	Add     []string `yaml:"add,flow"`
	Exclude []string `yaml:"exclude,flow"`
}

type UM struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	AppId    string `yaml:"appId"`
	AppToken string `yaml:"appToken"`
	AppPem   string `yaml:"appPem"`
}

type SSO struct {
	DefaultService      string `yaml:"defaultService"`
	CasServiceUrlPrefix string `yaml:"casServiceUrlPrefix"`
	CasLogin            string `yaml:"casLogin"`
	CasLogout           string `yaml:"casLogout"`
	SsoLanding          string `yaml:"ssoLanding"`
}

type IMS struct {
	AccessUrl     string `yaml:"accessUrl"`
	SubSystemId   string `yaml:"subSystemId"`
	AlertWay      string `yaml:"alertWay"`
	AlertReceiver string `yaml:"alertReceiver"`
}

type Kube struct {
	ApiConfig                ApiConfig                `yaml:"apiConfig"`
	NamespacedResourceConfig NamespacedResourceConfig `yaml:"namespacedResourceConfig"`
}

type ApiConfig struct {
	SysProperty string `yaml:"sysProperty"`
	DefaultPath string `yaml:"defaultPath"`
	Path        string `yaml:"/etc/config/kube-config"`
}

type NamespacedResourceConfig struct {
	DefaultRQName string `yaml:"defaultRQName"`
	DefaultRQCpu  string `yaml:"defaultRQCpu"`
	DefaultRQMem  string `yaml:"defaultRQMem"`
	DefaultRQGpu  string `yaml:"defaultRQGpu"`
}

//type Gateway struct {
//	BdpAddress      string `yaml:"bdpAddress"`
//	BdapAddress     string `yaml:"bdapAddress"`
//	BdapsafeAddress string `yaml:"bdapsafeAddress"`
//}

type AuthAddress struct {
	User string `yaml:"user"`
	Auth string `yaml:"auth"`
	IP   string `yaml:"ip"`
}

type Cache struct {
	CacheTime               string `yaml:"cacheTime"`
	DefaultTimestampTimeout int64  `yaml:"defaultTimestampTimeout"`
}

type Cookie struct {
	DefaultTime int    `yaml:"defaultTime"`
	Path        string `yaml:"path"`
}
