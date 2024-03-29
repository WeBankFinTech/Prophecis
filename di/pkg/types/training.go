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

type JobInfo struct {
	// The name of the training job
	Name string `json:"name"`
	// The namespace of the training job
	Namespace string `json:"namespace"`
	// The time of the training job
	Duration string `json:"duration"`
	// The status of the training Job
	Status JobStatus `json:"status"`

	// The training type of the training job
	Trainer string `json:"trainer"`
	// The tensorboard of the training job
	Tensorboard string `json:"tensorboard,omitempty"`

	// The name of the chief Instance
	ChiefName string `json:"chiefName" yaml:"chiefName"`

	// The instances under the training job
	Instances []Instance `json:"instances"`

	// The priority of the training job
	Priority string `json:"priority"`
}

// all the kinds of JobStatus
type JobStatus string

const (
	// JobPending means the job is pending
	JobPending JobStatus = "PENDING"
	// JobRunning means the job is running
	JobRunning JobStatus = "RUNNING"
	// JobSucceeded means the job is Succeeded
	JobSucceeded JobStatus = "SUCCEEDED"
	// JobFailed means the job is failed
	JobFailed JobStatus = "FAILED"
)

type Instance struct {
	// the status of of instance
	Status string `json:"status"`
	// the name of instance
	Name string `json:"name"`
	// the age of instance
	Age string `json:"age"`
	// the node instance runs on
	Node string `json:"node"`
	// the instance is chief or not
	IsChief bool `json:"chief" yaml:"chief"`
}
