package pkg

type Tree struct {
	InputData  map[string]int `json:"inputData"`  //外部依赖指标
	OutputData map[string]int `json:"outputData"` //内部生成指标
	Nodes      []*Node        `json:"nodes"`
}
