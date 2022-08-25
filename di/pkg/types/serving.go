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
package types

import (
	"errors"
)

// this file is used to define serving type

type ServingType string

// three serving types.
const (
	// tensorflow
	ServingTF ServingType = "TENSORFLOW"
	// tensorrt
	ServingTRT ServingType = "TENSORRT"
	// custom
	ServingCustom ServingType = "CUSTOM"
)

var (
	ErrNotFoundJobs = errors.New(`not found jobs under the assigned conditions.`)
	ErrTooManyJobs  = errors.New(`found jobs more than one,please use --version or --type to filter.`)
)

var SERVING_CHARTS = map[string]string{
	"tensorflow-serving-0.2.0":        "Tensorflow",
	"tensorrt-inference-server-0.0.1": "TensorRT",
}
