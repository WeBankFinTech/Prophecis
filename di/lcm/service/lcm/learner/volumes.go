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

package learner

import (
	v1core "k8s.io/api/core/v1"
)

const cosMountDriverName = "ibm/ibmc-s3fs"

// TODO: Fix copy-paste from trainer/storage/s3_object_store.go, to avoid circular ref
const (
	// DataStoreTypeS3 is the type string for the S3-based object store.
	DataStoreTypeMountVolume = "mount_volume"

	// This at the level of the data or result volume
	DataStoreHostMountVolume = "host_mount"

	DataStoreTypeMountCOSS3 = "mount_cos"
	DataStoreTypeS3         = "s3_datastore"
	defaultRegion           = "us-standard"
    learnerEntrypointFilesVolume = "learner-entrypoint-files"
    WorkSpaceVolume = "workspace"
	WorkSpaceHostPath = "/workspace"
)

//COSVolume ...
type COSVolume struct {
	VolumeType, ID, Region, Bucket, Endpoint, SecretRef, CacheSize, DiskFree, HostPath string
	MountSpec                                                                          VolumeMountSpec
}

type HostPathVolume struct {
	VolumeType, ID, Region, Bucket, Endpoint, SecretRef, CacheSize, DiskFree, HostPath string
	MountSpec                                                                          VolumeMountSpec
}

type ConfigMapVolume struct {
	VolumeType, Name string
	MountSpec        VolumeMountSpec
}

//Volumes ...
type Volumes struct {
	TrainingData    *COSVolume
	ResultsDir      *COSVolume
	WorkDir         *HostPathVolume
	AppComConfig    *COSVolume
	AppComCommonLib *COSVolume
	AppComInstall   *COSVolume
	JDK             *COSVolume
	Script          *ConfigMapVolume
	LogDir          *HostPathVolume
}

//VolumeMountSpec ...
type VolumeMountSpec struct {
	MountPath, SubPath, Name string
}

//CreateVolumeForLearner ...
func (volumes Volumes) CreateVolumeForLearner() []v1core.Volume {

	var volumeSpecs []v1core.Volume

	if volumes.TrainingData != nil {
		trainingDataParams := volumes.TrainingData
		volumeSpecs = append(volumeSpecs,
			generateHostMountTrainingDataVolume(trainingDataParams.ID, trainingDataParams.HostPath, v1core.HostPathUnset))
	}

	if volumes.ResultsDir != nil {
		resultDirParams := volumes.ResultsDir
		volumeSpecs = append(volumeSpecs,
			generateHostMountTrainingDataVolume(resultDirParams.ID, resultDirParams.HostPath, v1core.HostPathDirectoryOrCreate))
	}

	if volumes.WorkDir != nil {
		workDirParams := volumes.WorkDir
			volumeSpecs = append(volumeSpecs,
				generateHostPathVolume(
					workDirParams.ID,
					workDirParams.HostPath,
					v1core.HostPathUnset))
	}

	// FIXME MLSS Change: v_1.6.0 added HDFS config
	if volumes.AppComConfig != nil {
		configParams := volumes.AppComConfig
		volumeSpecs = append(volumeSpecs,
			generateHostMountTrainingDataVolume(
				configParams.ID,
				configParams.HostPath,
				v1core.HostPathUnset))
	}

	if volumes.AppComInstall != nil {
		installParams := volumes.AppComInstall
		volumeSpecs = append(volumeSpecs,
			generateHostMountTrainingDataVolume(
				installParams.ID,
				installParams.HostPath,
				v1core.HostPathUnset))
	}

	if volumes.JDK != nil {
		jdkParams := volumes.JDK
		volumeSpecs = append(volumeSpecs,
			generateHostMountTrainingDataVolume(
				jdkParams.ID,
				jdkParams.HostPath,
				v1core.HostPathUnset))
	}

	if volumes.Script != nil {
		scriptParams := volumes.Script
		volumeSpecs = append(volumeSpecs,
			generateConfigMapVolume(
				scriptParams.Name,
				7))
	}

	// learner entrypoint files volume
	learnerEntrypointFilesVolume := v1core.Volume{
		Name: learnerEntrypointFilesVolume,
		VolumeSource: v1core.VolumeSource{
			ConfigMap: &v1core.ConfigMapVolumeSource{
				LocalObjectReference: v1core.LocalObjectReference{
					Name: learnerEntrypointFilesVolume,
				},
			},
		},
	}
	volumeSpecs = append(volumeSpecs, learnerEntrypointFilesVolume)

	if volumes.LogDir != nil {
		LogDirParams := volumes.LogDir
		volumeSpecs = append(volumeSpecs,
			generateHostPathVolume(
				LogDirParams.ID,
				LogDirParams.HostPath,
				v1core.HostPathUnset))

	}

	//Workspace volume
	workspaceVolume := generateHostMountTrainingDataVolume(
		WorkSpaceVolume,
		WorkSpaceHostPath,
		v1core.HostPathUnset)
	volumeSpecs = append(volumeSpecs, workspaceVolume)

	return volumeSpecs
}

//CreateVolumeMountsForLearner ...
func (volumes Volumes) CreateVolumeMountsForLearner() []v1core.VolumeMount {

	var mounts []v1core.VolumeMount
	if volumes.TrainingData != nil {
		mounts = append(mounts, generateDataDirVolumeMount(volumes.TrainingData.ID, volumes.TrainingData.MountSpec.MountPath,
			volumes.TrainingData.MountSpec.SubPath))
	}

	if volumes.ResultsDir != nil {
		mounts = append(mounts, generateResultDirVolumeMount(volumes.ResultsDir.ID, volumes.ResultsDir.MountSpec.MountPath,
			volumes.ResultsDir.MountSpec.SubPath))
	}

	if volumes.WorkDir != nil {
		mounts = append(mounts, generateWorkDirVolumeMount(volumes.WorkDir.ID, volumes.WorkDir.MountSpec.MountPath,
			volumes.WorkDir.MountSpec.SubPath))
	}

	if volumes.AppComConfig != nil {
		mounts = append(mounts, generateAppConfigVolumeMount(volumes.AppComConfig.ID, volumes.AppComConfig.MountSpec.MountPath,
			volumes.AppComConfig.MountSpec.SubPath))
	}

	if volumes.AppComInstall != nil {
		mounts = append(mounts, generateAppInstallVolumeMount(volumes.AppComInstall.ID, volumes.AppComInstall.MountSpec.MountPath,
			volumes.AppComInstall.MountSpec.SubPath))
	}

	if volumes.JDK != nil {
		mounts = append(mounts, generateJDKVolumeMount(volumes.JDK.ID, volumes.JDK.MountSpec.MountPath,
			volumes.JDK.MountSpec.SubPath))
	}

	if volumes.LogDir != nil {
		mounts = append(mounts, generateJDKVolumeMount(volumes.LogDir.ID, volumes.LogDir.MountSpec.MountPath,
			volumes.LogDir.MountSpec.SubPath))
	}

	return mounts
}

func generateHostMountTrainingDataVolume(id, path string, pathType v1core.HostPathType) v1core.Volume {
	cosInputVolume := v1core.Volume{
		Name: id,
		VolumeSource: v1core.VolumeSource{
			HostPath: &v1core.HostPathVolumeSource{
				Path: path,
				Type: &pathType,
			},
		},
	}
	return cosInputVolume
}

func generateHostMountWorkspaceVolume(name, path string, pathType v1core.HostPathType) v1core.Volume {
	cosInputVolume := v1core.Volume{
		Name: name,
		VolumeSource: v1core.VolumeSource{
			HostPath: &v1core.HostPathVolumeSource{
				Path: path,
				Type: &pathType,
			},
		},
	}
	return cosInputVolume
}




func generateDataDirVolumeMount(id, bucket string, subPath string) v1core.VolumeMount {
	return v1core.VolumeMount{
		Name:      id,
		MountPath: bucket,
		SubPath:   subPath,
	}
}

func generateResultDirVolumeMount(id, bucket string, subPath string) v1core.VolumeMount {
	return v1core.VolumeMount{
		Name:      id,
		MountPath: bucket,
		SubPath:   subPath,
	}
}

func generateWorkDirVolumeMount(id, bucket string, subPath string) v1core.VolumeMount {
	return v1core.VolumeMount{
		Name:      id,
		MountPath: bucket,
		SubPath:   subPath,
	}
}

func generateAppConfigVolumeMount(id, bucket string, subPath string) v1core.VolumeMount {
	return v1core.VolumeMount{
		Name:      id,
		MountPath: bucket,
		SubPath:   subPath,
	}
}

func generateAppInstallVolumeMount(id, bucket string, subPath string) v1core.VolumeMount {
	return v1core.VolumeMount{
		Name:      id,
		MountPath: bucket,
		SubPath:   subPath,
	}
}

func generateJDKVolumeMount(id, bucket string, subPath string) v1core.VolumeMount {
	return v1core.VolumeMount{
		Name:      id,
		MountPath: bucket,
		SubPath:   subPath,
	}
}

func generateHostPathVolume(id, path string, pathType v1core.HostPathType) v1core.Volume {
	hostMountVolume := v1core.Volume{
		Name: id,
		VolumeSource: v1core.VolumeSource{
			HostPath: &v1core.HostPathVolumeSource{
				Path: path,
				Type: &pathType,
			},
		},
	}

	return hostMountVolume
}

func generateConfigMapVolume(name string, mode int32) v1core.Volume {
	configmapVolume := v1core.Volume{
		Name: name,
		VolumeSource: v1core.VolumeSource{
			ConfigMap: &v1core.ConfigMapVolumeSource{
				DefaultMode: &mode,
				LocalObjectReference: v1core.LocalObjectReference{
					Name: name,
				},
			},
		},
	}

	return configmapVolume
}
