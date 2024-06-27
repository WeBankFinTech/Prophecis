package rest_impl

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/kubeflow/pipelines/api/v2alpha1/go/pipelinespec"
)

type V2Spec struct {
	Spec         *pipelinespec.PipelineSpec
	PlatformSpec *pipelinespec.PlatformSpec
}

type FlowJson struct {
	Edges           []Edge           `json:"edges"`
	Nodes           []Node           `json:"nodes"`
	Comment         string           `json:"comment"`
	Type            string           `json:"type"`
	UpdateTime      int64            `json:"updateTime"`
	Props           []Prop           `json:"props"`
	Resources       []Resource       `json:"resources"`
	ContextID       string           `json:"contextID"`
	Labels          Label            `json:"labels"`
	GlobalVariables []GlobalVariable `json:"global_variables,omitempty"`
}

type GlobalVariable struct {
	Type  string      `json:"type"`
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type Edge struct {
	Source         string `json:"source"`
	Target         string `json:"target"`
	SourceLocation string `json:"sourceLocation"`
	TargetLocation string `json:"targetLocation"`
}

type Node struct {
	Key                      string     `json:"key"`
	Title                    string     `json:"title"`
	Desc                     string     `json:"desc"`
	Layout                   Layout     `json:"layout"`
	Resources                []Resource `json:"resources"`
	EditParam                bool       `json:"editParam"`
	EditBaseInfo             bool       `json:"editBaseInfo"`
	SubmitToScheduler        bool       `json:"submitToScheduler"`
	EnableCopy               bool       `json:"enableCopy"`
	ShouldCreationBeforeNode bool       `json:"shouldCreationBeforeNode"`
	SupportJump              bool       `json:"supportJump"`
	ID                       string     `json:"id"`
	JobContent               JobContent `json:"jobContent"`
	JobType                  string     `json:"jobType"`
	Selected                 bool       `json:"selected"`
}

type Layout struct {
	Height int `json:"height"`
	Width  int `json:"width"`
	X      int `json:"x"`
	Y      int `json:"y"`
}

type Resource struct {
	// Define the structure of resources if available
}

type JobContent struct {
	Method                 string      `json:"method"`
	Params                 Params      `json:"params"`
	ExecuteApplicationName string      `json:"executeApplicationName"`
	RunType                string      `json:"runType"`
	Source                 Source      `json:"source"`
	ManiFest               interface{} `json:"ManiFest"`
	Content                interface{} `json:"content"`
	MlflowJobType          string      `json:"mlflowJobType"`
}

type Params struct {
	Ariable       map[string]string `json:"ariable"`
	Configuration Configuration     `json:"configuration"`
}

type Configuration struct {
	Special map[string]string `json:"special"`
	Runtime map[string]string `json:"runtime"`
	Startup map[string]string `json:"startup"`
}

type Source struct {
	ScriptPath string `json:"scriptPath"`
}

// todo(gaoyuanhe): 这个Manifest是针对gpu的，考虑去掉
type Manifest struct {
	Name              string           `json:"name" yaml:"name"`
	Description       string           `json:"description" yaml:"description"`
	Version           string           `json:"version" yaml:"version"`
	Gpus              int              `json:"gpus" yaml:"gpus"`
	Cpus              float64          `json:"cpus" yaml:"cpus"`
	Memory            string           `json:"memory" yaml:"memory"`
	Namespace         string           `json:"namespace" yaml:"namespace"`
	CodeSelector      string           `json:"code_selector" yaml:"code_selector"`
	JobType           string           `json:"job_type" yaml:"job_type"`
	DataStores        []DataStore      `json:"data_stores" yaml:"data_stores"`
	Framework         Framework        `json:"framework" yaml:"framework"`
	EvaluationMetrics EvaluationMetric `json:"evaluation_metrics" yaml:"evaluation_metrics"`
	//ExpID             int              `json:"exp_id"`
	//ExpName           string           `json:"exp_name"`
	// Add other fields if present
}

type DataStore struct {
	ID                string     `json:"id" yaml:"id"`
	Type              string     `json:"type" yaml:"type"`
	TrainingData      Container  `json:"training_data" yaml:"training_data"`
	TrainingResults   Container  `json:"training_results" yaml:"training_results"`
	Connection        Connection `json:"connection" yaml:"connection"`
	TrainingWorkspace Container  `json:"training_workspace" yaml:"training_workspace"`
}

type Container struct {
	Container string `json:"container"`
}

type Connection struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Path string `json:"path"`
}

type Framework struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Command string `json:"command"`
}

type EvaluationMetric struct {
	Type string `json:"type"`
}

type Prop struct {
	UserToProxy string `json:"user.to.proxy"`
}

type Label struct {
	Route string `json:"route"`
}

func validateFlowJson(flowJsonStr string) (bool, string) {
	if flowJsonStr == "" {
		return true, ""
	}
	var flowJson FlowJson
	err := json.Unmarshal([]byte(flowJsonStr), &flowJson)
	if err != nil {
		glog.Errorf("Error json.Unmarshal: %v\n", err)
		return false, "解析flowjson失败"
	}
	// 检查0：flowjson里的 node.jobContent.Manifest不能为空
	for _, node := range flowJson.Nodes {
		if node.JobContent.ManiFest == nil {
			return false, fmt.Sprintf("节点信息为空")
		}
		manifest, _ := node.JobContent.ManiFest.(map[string]interface{})
		if len(manifest) == 0 {
			return false, fmt.Sprintf("节点信息为空")
		}
	}
	// 检查1：flowjson里的edge的源节点和目的节点，在nodes里必须有
	for _, edge := range flowJson.Edges {
		edgeSourceNode := edge.Source
		found := false
		for _, node := range flowJson.Nodes {
			if node.Key == edgeSourceNode {
				found = true
			}
		}
		if !found {
			return false, fmt.Sprintf("edgeSourceNode(%s)在nodes里不存在", edgeSourceNode)
		}
		edgeTargetNode := edge.Target
		found = false
		for _, node := range flowJson.Nodes {
			if node.Key == edgeTargetNode {
				found = true
			}
		}
		if !found {
			return false, fmt.Sprintf("edgeTargetNode(%s)在nodes里不存在", edgeTargetNode)
		}
	}

	// 检查2：主流程加一个额外信号节点的结构, 且主流程中不能有信号节点
	nodeNums := len(flowJson.Nodes)
	if nodeNums <= 1 {
		// 一个节点， do nothing
	} else if nodeNums == 2 {
		edgeNums := len(flowJson.Edges)
		if edgeNums > nodeNums-1 {
			return false, fmt.Sprintf("工作流的边数超过预期")
		} else if edgeNums == nodeNums-1 {
			// 2个节点相连，构成主流程，查看主流程中是否有信号节点
			if ok, reason := mainPipelineShouldNotHaveEventNode(flowJson); !ok {
				return ok, reason
			}
		} else if edgeNums == nodeNums-2 {
			// 2个节点相互独立，那不能相同，只能一个是信号节点，一个是非信号节点
			type1 := flowJson.Nodes[0].JobContent.MlflowJobType
			type2 := flowJson.Nodes[1].JobContent.MlflowJobType
			if type1 == type2 {
				return false, fmt.Sprintf("2个节点相互独立，只能一个是信号节点，一个是非信号节点")
			}
		} else {
			// 不会有负数的情况，do nothing
		}
	} else {
		edgeNums := len(flowJson.Edges)
		if edgeNums > nodeNums-1 {
			// 边数超过预期
			return false, fmt.Sprintf("工作流的边数超过预期")
		} else if edgeNums == nodeNums-1 {
			// 查看主流程中是否有信号节点
			if ok, reason := mainPipelineShouldNotHaveEventNode(flowJson); !ok {
				return ok, reason
			}
		} else if edgeNums == nodeNums-2 {
			// 查看主流程中是否有信号节点以及另一个独立节点是否为信号节点
			if ok, reason := mainPipelineShouldNotHaveEventNode(flowJson); !ok {
				return ok, reason
			}
			if ok, reason := UndependentNodeMustBeEventNode(flowJson); !ok {
				return ok, reason
			}
		} else {
			// 边数不符合预期, 不是主流程加一个额外信号节点的架构
			return false, fmt.Sprintf("工作流不是主流程加一个额外信号节点")
		}
	}

	// 检查 flowjson 里的存储是否有权限

	return true, ""
}

func mainPipelineShouldNotHaveEventNode(flowJson FlowJson) (bool, string) {
	// 主流程中不能有信号节点
	// 获取edge中的节点名作为主流程节点集合
	nodeSet := make(map[string]bool)
	for _, edge := range flowJson.Edges {
		edgeSourceNode := edge.Source
		edgeTargetNode := edge.Target
		nodeSet[edgeSourceNode] = true
		nodeSet[edgeTargetNode] = true
	}
	for nodeInMainPipeline, _ := range nodeSet {
		for _, node := range flowJson.Nodes {
			if node.Key == nodeInMainPipeline {
				if node.JobContent.MlflowJobType == "EventSender" {
					return false, fmt.Sprintf("主流程中不能有信号节点")
				}
			}
		}
	}
	return true, ""
}

func UndependentNodeMustBeEventNode(flowJson FlowJson) (bool, string) {
	// 独立的节点只能为信号节点
	// 获取edge中的节点名作为主流程节点集合
	nodeSet := make(map[string]bool)
	for _, edge := range flowJson.Edges {
		edgeSourceNode := edge.Source
		edgeTargetNode := edge.Target
		nodeSet[edgeSourceNode] = true
		nodeSet[edgeTargetNode] = true
	}
	// 找到不在主流程的独立节点
	for _, node := range flowJson.Nodes {
		if nodeSet[node.Key] == false {
			if node.JobContent.MlflowJobType != "EventSender" {
				return false, fmt.Sprintf("独立节点只能为信号节点")
			}
		}
	}
	return true, ""
}
