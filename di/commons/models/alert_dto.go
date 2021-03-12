package models

type AlertDTO struct {
	UserID       string `json:"userID"`
	JobNamespace string `json:"jobNamespace"`
	JobName      string `json:"jobName"`
	TrainingID   string `json:"trainingID"`
	JobStatus    string `json:"jobStatus"`
	AlertLevel   string `json:"alertLevel"`
	AlertReason  string `json:"alertReason"`
	StartTime    string `json:"startTime"`
	EndTime      string `json:"endTime"`
	Receiver     string `json:"receiver"`
	//JobAlert	 map[string][]map[string]string `json:"jobAlert"`
}
