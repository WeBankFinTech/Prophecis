package models

type DSSFlowVo struct {
	Nodes []Node `json:"nodes"`
}
type Node struct {
	Id       string `json:"id"`
	JobType  string `json:"jobType"`
	CodePath string `json:"code_path"`
}
