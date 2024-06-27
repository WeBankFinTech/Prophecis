package constants

import (
	"os"

	v1 "k8s.io/api/core/v1"
)

// FIXME MLSS Change: add timezone VolumeMount
var TimezoneMount = v1.VolumeMount{
	Name:      "timezone-volume",
	MountPath: "/etc/localtime",
	ReadOnly:  true,
}

// FIXME MLSS Change: add timezone VolumeMount
var JobLogsMount = v1.VolumeMount{
	Name:      "job-logs",
	MountPath: "/job-logs",
}

// FIXME MLSS Change: define timezone volume
var TimezoneVolume = v1.Volume{
	Name: "timezone-volume",
	VolumeSource: v1.VolumeSource{
		HostPath: &v1.HostPathVolumeSource{
			Path: "/usr/share/zoneinfo/Asia/Shanghai",
		},
	},
}

// FIXME MLSS Change: define timezone volume
var JobLogsVolume = v1.Volume{
	Name: "job-logs",
	VolumeSource: v1.VolumeSource{
		HostPath: &v1.HostPathVolumeSource{
			Path: os.Getenv(FLUENT_BIT_LOG_PATH) + os.Getenv(ENVIR_UPRER),
		},
	},
}

// FIXME MLSS Temporary Change: use fluent-bit
var FluentbitConfigVolumeMount = v1.VolumeMount{
	Name:      "fluent-bit-log-collector-config",
	MountPath: "/fluent-bit/etc/",
}

// FIXME MLSS Temporary Change: use fluent-bit
var FluentbitConfigVolume = v1.Volume{
	Name: "fluent-bit-log-collector-config",
	VolumeSource: v1.VolumeSource{
		ConfigMap: &v1.ConfigMapVolumeSource{
			LocalObjectReference: v1.LocalObjectReference{
				Name: "fluent-bit-log-collector-config",
			},
		},
	},
}

const (
	Fault = 0
	True  = 1
)

const TFOS = "tfos"

const BDAPSAFE_SS = "bdapsafe-ss"
const BDAP_SS = "bdap-ss"
const BDP_SS = "bdp-ss"

const BDAP = "bdap"
const BDP = "bdp"
const SAFE = "safe"

const ENVIR = "envir"

const ENVIR_UPRER string = "ENVIR"
const FLUENT_BIT_LOG_PATH string = "FLUENT_BIT_LOG_PATH"

const JOBTYPE string = "jobType"
const SINGLE string = "single"
const DI_START_NODEPORT string = "DI_START_NODEPORT"
const DI_END_NODEPORT string = "DI_END_NODEPORT"

const MongoAddressKey = "MONGO_ADDRESS"
const MongoDatabaseKey = "MONGO_DATABASE"
const MongoUsernameKey = "MONGO_USERNAME"
const MongoPasswordKey = "MONGO_PASSWORD"
const MongoAuthenticationDatabase = "MONGO_Authentication_Database"

const (
	LOG_STORING_KEY       = "store_status"
	LOG_STORING_COMPLETED = "completed"
	LOG_STORING_STORING   = "storing"

	GoDefaultTimeFormat = "2006-01-02 15:04:05"
)
