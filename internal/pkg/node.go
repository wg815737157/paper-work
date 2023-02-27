package pkg

type Node struct {
	IsUsed          bool
	IsSatisfied     bool
	IsSuccessful    bool
	RuleIdList      []*RuleItem
	Pnodes          []int    //父节点
	PnodeConditions []string //执行条件
	Cnodes          []int    //子节点
	result          string   //节点执行结果
}

type RuleItem struct {
}
