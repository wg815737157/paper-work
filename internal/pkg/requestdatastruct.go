package pkg

type RuleNodeRequest struct {
	NodeId     int            `json:"nodeId"`
	NodeName   string         `json:"nodeName"`
	RuleIdList []int          `json:"ruleIdList"`
	InputData  map[string]int `json:"inputData"`
	OutputData map[string]int `json:"outputData"`
}
type RuleNodeResponse struct {
	Code             int                   `json:"code"`
	Message          string                `json:"message"`
	RuleResponseData *RuleNodeResponseData `json:"data"`
}
type RuleNodeResponseData struct {
	NodeId         int            `json:"nodeId"`
	NodeName       string         `json:"nodeName"`
	RuleIdList     []int          `json:"ruleIdList"`
	RuleNameList   []string       `json:"ruleNameList"`
	InputData      map[string]int `json:"inputData"`
	OutputData     map[string]int `json:"outputData"`
	RuleResultList []string       `json:"ruleResultList"`
}

type RuleTreeResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Tree    *Tree  `json:"data"`
}
