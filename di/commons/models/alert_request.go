package models

import "time"

type AlertRequest struct {
	// FIXME MLSS Change: v_1.4.1 added field jobAlert and deleted field for alert param
	//EventChecker    string
	//DeadlineChecker string
	//OvertimeChecker string
	//AlertLevel      string
	//Receiver        string
	StartTime time.Time
	//Interval        string
	JobAlert map[string][]map[string]string
}
