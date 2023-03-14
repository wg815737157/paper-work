package pkg

type RuleRequest struct {
	RuleNameList []string       `json:"rule_name_list"`
	InputData    map[string]int `json:"inputData"`
	OutputData   map[string]int `json:"outputData"`
}
type RuleResponse struct {
	RuleIdList     []int          `json:"ruleIdList"`
	RuleNameList   []string       `json:"rule_name_list"`
	InputData      map[string]int `json:"inputData"`
	OutputData     map[string]int `json:"outputData"`
	RuleResultList []string       `json:"ruleResultList"`
}
