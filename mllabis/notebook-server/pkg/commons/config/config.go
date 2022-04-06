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
package config

import (
	"fmt"
	"strings"
	"sync"

	"google.golang.org/grpc/grpclog"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	// PortKey is a key for looking up the Port setting in the viper config.
	PortKey = "port"
	// LogLevelKey is a key for looking up the LogLevel setting in the viper config.
	LogLevelKey = "loglevel"
	// EnvKey is a key for looking up the Environment setting in the viper config.
	EnvKey = "env"

	MLSSVersion       = "mlssVersion"
	CCAddress         = "ccAddress"
	PlatformNamespace = "platformNamespace"
	MLSSGroupId       = "mlssGroupId"
	CaCertPath        = "caCertPath"

	HadoopEnable     = "enableHadoop"
	InstallPath      = "installPath"
	CommonlibPath    = "commonlibPath"
	ConfigPath       = "configPath"
	JavaPath         = "javaPath"
	SourceConfigPath = "sourceConfigPath"
	HostFilePath     = "hostFilePath"
	LinkisToken      = "linkisToken"
	MaxSparkSessionNum = "maxSparkSessionNum"
)

var viperInitOnce sync.Once

func init() {
	InitViper()
}

// InitViper is initializing the configuration system
func InitViper() {

	viperInitOnce.Do(func() {

		//viper.SetEnvPrefix(envPrefix) // will be uppercased automatically
		viper.SetConfigType("yaml")
		// enable ENV vars and defaults
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
		viper.SetDefault(EnvKey, "prod")
		setLogLevel()
		viper.SetDefault(PortKey, 8080)

		// config file is optional. we usually configure via ENV_VARS
		configFile := fmt.Sprintf("config")
		viper.SetConfigName(configFile) // name of config file (without extension)
		viper.AddConfigPath("/etc/mlss/")
		viper.AddConfigPath(".")

		err := viper.ReadInConfig()
		if err != nil {
			switch err := err.(type) {
			case viper.ConfigFileNotFoundError:
				log.Debugf("No config file '%s.yml' found. Using environment variables only.", configFile)
			case viper.ConfigParseError:
				log.Panicf("Cannot read config file: %s.yml: %s\n", configFile, err)
			}
		} else {
			log.Debugf("Loading config from file %s\n", viper.ConfigFileUsed())
		}
	})

	log.Infof("THIS IS MLSS Version %s\n", viper.GetString(MLSSVersion))
	log.Infof("Platform namespace is %s\n", viper.GetString(PlatformNamespace))

	grpclog.SetLogger(log.StandardLogger())
}

// GetInt returns the int value for the config key.
func GetInt(key string) int {
	return viper.GetInt(key)
}

// GetInt64 returns the int64 value for the config key.
func GetInt64(key string) int64 {
	return viper.GetInt64(key)
}

// GetString returns the string value for the config key.
func GetString(key string) string {
	return viper.GetString(key)
}

// GetValue returns the value associated with the given config key.
func GetValue(key string) string {
	return viper.GetString(key)
}

// setLogLevel sets the logging level based on the environment
func setLogLevel() {
	viper.SetDefault(LogLevelKey, "info")

	env := viper.GetString(EnvKey)
	if env == "dev" || env == "test" {
		viper.Set(LogLevelKey, "debug")
	}

	if logLevel, err := log.ParseLevel(viper.GetString(LogLevelKey)); err == nil {
		log.SetLevel(logLevel)
	}

	log.Debugf("Log level set to '%s'", viper.Get(LogLevelKey))
}
