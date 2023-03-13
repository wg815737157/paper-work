package pkg

type RuleRequest struct {
	ruleId     int                   `gorm:"rule_id"`
	InputData  map[string]RuleReturn `gorm:"input_data"`
	OutputData map[string]RuleReturn `gorm:"output_data"`
}
type RuleResponse struct {
	ruleId     int                   `gorm:"rule_id"`
	InputData  map[string]RuleReturn `gorm:"input_data"`
	OutputData map[string]RuleReturn `gorm:"output_data"`
}
