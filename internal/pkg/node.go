package pkg

type Node struct {
	NodeId          int            `json:"nodeId"`          //节点图的数组下表
	NodeName        string         `json:"nodeName"`        //节点名
	NodeType        string         `json:"nodeTye"`         //节点类型
	IsUsed          bool           `json:"isUsed"`          // 节点是否执行到了
	IsSatisfied     bool           `json:"isSatisfied"`     //是否满足执行条件
	IsSuccessful    bool           `json:"isSuccessful"`    //节点是否执行完成
	RuleIdList      []int          `json:"ruleIdList"`      //节点规则列表
	Pnodes          []int          `json:"pnodes"`          //父节点
	PnodeConditions map[int]string `json:"pnodeConditions"` //执行条件
	MergeType       string         `json:"mergeType"`       //合并类型 && 或者||
	Cnodes          []int          `json:"cnodes"`          //子节点
	Result          string         `json:"result"`          //节点执行结果
}

//RuleExpression:payRatio>10&&income==50
//RuleReturn:loanCheck=1
//RuleReturn:name=result,type=string,value=pass

type RuleDetail struct {
	RuleExpression string       `json:"ruleExpression"`  //规则表达式
	ResResult      []RuleReturn `json:"resResultString"` //规则执行后的结果
}
type RuleReturn struct {
	Type  string `json:"type"`  //类型 bool,int,string,expression
	Name  string `json:"name"`  //规则执行成功后的变量名
	Value any    `json:"value"` //bool:true,false;int:integer,string:pass,fail,expression:表达式
}
