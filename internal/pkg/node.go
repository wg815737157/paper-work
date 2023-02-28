package pkg

type Node struct {
	NodeId   int
	NodeName string
	NodeType string //节点类型
	IsUsed   bool

	RuleIdList      []*RuleItem
	Pnodes          []int          //父节点
	PnodeConditions map[int]string //执行条件
	MergeType       string         //合并条件
	IsSatisfied     bool           //是否满足执行条件
	IsSuccessful    bool           //节点是否执行完成
	Cnodes          []int          //子节点
	Result          string         //节点执行结果

}

type RuleItem struct {
}
