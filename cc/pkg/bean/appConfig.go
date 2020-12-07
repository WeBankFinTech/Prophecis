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
package bean

type AppConfig struct {
	Core struct {
		Interceptor Interceptor `yaml:"interceptor"`
	}
	Server Server `yaml:"server"`
}

type Server struct {
	Port    string `yaml:"port"`
	Servlet struct {
		ContextPath string `yaml:"context-path"`
	}
}

// Yaml1 struct of yaml
type Interceptor struct {
	Configs []InterceptorConfig `yaml:"configs,flow"`
}


type InterceptorConfig struct {
	Name    string   `yaml:"name"`
	Add     []string `yaml:"add,flow"`
	Exclude []string `yaml:"exclude,flow"`
}
