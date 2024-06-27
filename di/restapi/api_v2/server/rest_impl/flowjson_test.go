package rest_impl

import (
	"encoding/json"
	"github.com/golang/glog"
	"testing"
)

var flowJsonStr = `
{
  "nodes": [
    {
      "id": "pm4vusfzqxsmaf",
      "type": "generic",
      "dimensions": {
        "width": 200,
        "height": 40
      },
      "handleBounds": {
        "source": [
          {
            "id": "pm4vusfzqxsmaf__handle-top",
            "position": "top",
            "x": 96,
            "y": -4,
            "width": 8,
            "height": 8
          },
          {
            "id": "pm4vusfzqxsmaf__handle-left",
            "position": "left",
            "x": -4,
            "y": 16,
            "width": 8,
            "height": 8
          },
          {
            "id": "pm4vusfzqxsmaf__handle-bottom",
            "position": "bottom",
            "x": 96,
            "y": 36,
            "width": 8,
            "height": 8
          },
          {
            "id": "pm4vusfzqxsmaf__handle-right",
            "position": "right",
            "x": 196,
            "y": 16,
            "width": 8,
            "height": 8
          }
        ]
      },
      "computedPosition": {
        "x": 414,
        "y": 229.578125,
        "z": 0
      },
      "selected": false,
      "dragging": false,
      "resizing": false,
      "initialized": true,
      "isParent": false,
      "position": {
        "x": 414,
        "y": 229.578125
      },
      "events": {},
      "key": "pm4vusfzqxsmaf",
      "title": "predictivemodeling-jju5v",
      "label": "predictivemodeling-jju5v",
      "jobType": "linkis.appconn.mlflow.model_predict",
      "enableCopy": true,
      "supportJump": true,
      "editBaseInfo": false,
      "shouldCreationBeforeNode": false,
      "submitToScheduler": true,
      "jobContent": {
        "type": "PredictiveModeling",
        "jobType": "linkis.appconn.mlflow.model_predict",
        "mlflowJobType": "ModelPredict",
        "id": "pm4vusfzqxsmaf",
        "ManiFest": {}
      }
    }
  ],
  "edges": [],
  "viewport": {
    "x": 0,
    "y": 0,
    "zoom": 1
  }
}
`

func Test1(t *testing.T) {
	var flowJson FlowJson
	err := json.Unmarshal([]byte(flowJsonStr), &flowJson)
	if err != nil {
		glog.Errorf("Error json.Unmarshal: %v\n", err)
	}
	for _, node := range flowJson.Nodes {
		if node.JobContent.ManiFest == nil {
			glog.Errorf("节点信息为空")
		}
		manifest, _ := node.JobContent.ManiFest.(map[string]interface{})
		if len(manifest) == 0 {
			glog.Errorf("节点信息为空")
		}
	}
}
