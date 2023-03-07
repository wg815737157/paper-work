package ruleserver

import (
	"container/list"
	"fmt"
	"github.com/wg815737157/paper-work/pkg/log"
)

type Operator string

func (o Operator) DoubleOperatorFunc(a, b any) any {
	var result any
	switch o {
	case "+":
		result = a.(int) + b.(int)
	case "-":
		result = a.(int) - b.(int)
	case "*":
		result = a.(int) * b.(int)
	case "/":
		result = a.(int) / b.(int)
	case ">":
		result = a.(int) > b.(int)
	case "<":
		result = a.(int) < b.(int)
	case "&&":
		result = a.(bool) && b.(bool)
	case "||":
		result = a.(bool) || b.(bool)
	}
	return result
}
func printOperatorList(operatorList *list.List) {
	operator := operatorList.Front()
	operatorList.Remove(operator)
	fmt.Println(operator.Value.(string))
}

type TreeNode struct {
	Type  string // 数值，操作符
	Value any
	Left  *TreeNode
	Right *TreeNode
}

func PrintTreeNodeList(treeNodeList *list.List) {
	e := treeNodeList.Front()
	for e != nil {
		if e.Value.(*TreeNode).Type == "number" {
			fmt.Println(e.Value.(*TreeNode).Value)
		} else {
			fmt.Println(e.Value.(*TreeNode).Value.(string))
		}
		e = e.Next()
	}
}

func pushOperator(inputOperator string, operatorList *list.List, postfixTreeList *list.List) {
	operatorLevel := map[string]int{"||": 12, "&&": 11, "<": 6, ">": 6, "+": 4, "-": 4, "*": 3, "/": 3, "(": 1, ")": 1}
	for {
		tmpElement := operatorList.Front()
		if tmpElement == nil {
			break
		}
		operatorInlist := tmpElement.Value.(string)
		// 堆栈里操作符优先级大于要插入的操作符优先级
		if operatorInlist == "(" || operatorLevel[inputOperator] < operatorLevel[operatorInlist] {
			break
		}
		operatorList.Remove(tmpElement)
		tmpTreeNode := &TreeNode{Type: "operator", Value: operatorInlist}
		postfixTreeList.PushBack(tmpTreeNode)
	}
	operatorList.PushFront(inputOperator)
}

func GenPostfixList(infixString string) *list.List {
	operatorList := list.New()
	postfixList := list.New()
	for i := 0; i < len(infixString); {
		if infixString[i] == ' ' {
			i++
			continue
		}
		//操作数
		if '0' <= infixString[i] && infixString[i] <= '9' {
			num := 0
			for i < len(infixString) && '0' <= infixString[i] && infixString[i] <= '9' {
				num = 10*num + int(infixString[i]-'0')
				i++
			}
			tmpTreeNode := &TreeNode{Type: "number", Value: num}
			postfixList.PushBack(tmpTreeNode)
			continue
		}
		// 且操作符
		if infixString[i] == '&' {
			if i+1 < len(infixString) {
				if infixString[i+1] != '&' {
					log.SugarLogger().Error("and operator err")
					return nil
				}
			}
			pushOperator("&&", operatorList, postfixList)
			i += 2
			continue
		}
		// 或操作符
		if infixString[i] == '|' {
			if i+1 < len(infixString) {
				if infixString[i+1] != '|' {
					log.SugarLogger().Error("and operator err")
					return nil
				}
			}
			pushOperator("||", operatorList, postfixList)
			i += 2
			continue
		}

		//操作符
		inputOperator := string(infixString[i])
		if inputOperator == "(" {
			operatorList.PushFront(inputOperator)
			i++
			continue
		}
		if inputOperator == ")" {
			tmpElement := operatorList.Front()
			operatorList.Remove(tmpElement)
			for tmpElement.Value.(string) != "(" {
				tmpTreeNode := &TreeNode{Type: "operator", Value: tmpElement.Value.(string)}
				postfixList.PushBack(tmpTreeNode)
				tmpElement = operatorList.Front()
				operatorList.Remove(tmpElement)
			}
			i++
			continue
		}
		pushOperator(string(infixString[i]), operatorList, postfixList)
		i++
		continue
	}
	for operatorList.Len() != 0 {
		tmpElement := operatorList.Front()
		operatorList.Remove(tmpElement)
		tmpTreeNode := &TreeNode{Type: "operator", Value: tmpElement.Value.(string)}
		postfixList.PushBack(tmpTreeNode)
	}
	return postfixList
}

func GeneratorGenPostfixList(infixString string) *list.List {
	postfixList := GenPostfixList(infixString)
	return postfixList
}

func ExecutePostfixList(postfixList *list.List) bool {
	executeStack := list.New()
	for postfixList.Front() != nil {
		e := postfixList.Front()
		postfixList.Remove(e)
		if e.Value.(*TreeNode).Type == "number" {
			operatorValue := e.Value.(*TreeNode).Value.(int)
			executeStack.PushFront(operatorValue)
		} else { //operator
			operator := e.Value.(*TreeNode).Value.(string)
			operatorElement1 := executeStack.Front()
			executeStack.Remove(operatorElement1)
			operatorElement2 := executeStack.Front()
			executeStack.Remove(operatorElement2)
			operatorValue1 := operatorElement1.Value
			operatorValue2 := operatorElement2.Value
			result := Operator(operator).DoubleOperatorFunc(operatorValue2, operatorValue1)
			executeStack.PushFront(result)
		}
	}
	resultElement := executeStack.Front()
	executeStack.Remove(resultElement)
	return resultElement.Value.(bool)
}
