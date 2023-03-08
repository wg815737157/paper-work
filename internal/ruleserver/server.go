package ruleserver

import (
	"fmt"
)

func Run() {

	inputData := map[string]int{"a": 1, "b": 2, "c": 3}
	postfixList := GeneratorGenPostfixList("c-2==a", inputData)
	//PrintTreeNodeList(postfixList)
	fmt.Println(ExecutePostfixList(postfixList))
}
