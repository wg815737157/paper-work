package ruleserver

import "fmt"

func Run() {
	postfixList := GeneratorGenPostfixList("20-10<5 && 11<1||2>1")
	PrintTreeNodeList(postfixList)
	fmt.Println(ExecutePostfixList(postfixList))
}
