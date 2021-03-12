package models

type ExportFlowDefinitionYamlModel struct {
	ExportExp     ExportExp         `yaml:"experiment"`
	ExportDSSFlow ExportDSSFlow     `yaml:"dss_flow"`
	CodePaths     map[string]string `yaml:"code_paths"`
}

type ExportExp struct {
	ExpId        int64				`yaml:"exp_id"`
	ExpName *string           `yaml:"exp_name"`
	ExpDesc *string           `yaml:"exp_desc"`
	ExpTag  map[string]string `yaml:"exp_tag"`
}
type ExportDSSFlow struct {
	FlowJson string `yaml:"flow_json"`
}
