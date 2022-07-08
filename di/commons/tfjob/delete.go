package tfjob

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"webank/DI/commons/util/helm"
	"webank/DI/commons/workflow"
)

func DeleteTrainingJob(jobName string, namespace string, trainingType string) error {
	var trainingTypes []string
	// 1. Handle legacy training job
	err := helm.DeleteRelease(jobName)
	if err == nil {
		log.Infof("Delete the job %s successfully.", jobName)
		return nil
	}

	log.Debugf("%s wasn't deleted by helm due to %v", jobName, err)

	// 2. Handle training jobs created by arena
	if trainingType == "" {
		trainingTypes = getTrainingTypes(jobName, namespace)
		if len(trainingTypes) == 0 {
			return fmt.Errorf("There is no training job found with the name %s, please check it with `arena list | grep %s`",
				jobName,
				jobName)
		} else if len(trainingTypes) > 1 {
			return fmt.Errorf("There are more than 1 training jobs with the same name %s, please double check with `arena list | grep %s`. And use `arena delete %s --type` to delete the exact one.",
				jobName,
				jobName,
				jobName)
		}
	} else {
		trainingTypes = []string{trainingType}
	}

	err = workflow.DeleteJob(jobName, namespace, trainingTypes[0])
	if err != nil {
		return err
	}

	log.Infof("The Job %s has been deleted successfully", jobName)
	// (TODO: cheyang)3. Handle training jobs created by others, to implement
	return nil
}

func isKnownTrainingType(trainingType string) bool {
	for _, knownType := range knownTrainingTypes {
		if trainingType == knownType {
			return true
		}
	}

	return false
}
