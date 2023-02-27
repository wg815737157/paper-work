package pkg

type Tree struct {
	InputIndexes    map[string]any //外部依赖指标
	InternalIndexes map[string]any //内部生成指标
	Nodes           []*Node
}
