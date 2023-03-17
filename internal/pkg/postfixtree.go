package pkg

import (
	"container/list"
	"errors"
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
	case "==":
		result = a.(int) == b.(int)
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

//func PrintTreeNodeList(treeNodeList *list.List) {
//	e := treeNodeList.Front()
//	for e != nil {
//		if e.Value.(*TreeNode).Type == "number" {
//			fmt.Println(e.Value.(*TreeNode).Value)
//		} else {
//			fmt.Println(e.Value.(*TreeNode).Value.(string))
//		}
//		e = e.Next()
//	}
//}

func pushOperator(inputOperator string, operatorList *list.List, postfixTreeList *list.List) {
	operatorLevel := map[string]int{"||": 12, "&&": 11, "==": 7, "<": 6, ">": 6, "+": 4, "-": 4, "*": 3, "/": 3, "(": 1, ")": 1}
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

func GenPostfixList(infixString string, request *RuleNodeRequest) (*list.List, error) {
	operatorList := list.New()
	postfixList := list.New()
	for i := 0; i < len(infixString); {
		if infixString[i] == ' ' {
			i++
			continue
		}
		//变量
		if ('a' <= infixString[i] && infixString[i] <= 'z') || ('A' <= infixString[i] && infixString[i] <= 'Z') || infixString[i] == '_' {
			variableBytes := []byte{}
			variableBytes = append(variableBytes, infixString[i])
			i++
			for i < len(infixString) && ('a' <= infixString[i] && infixString[i] <= 'z' || 'A' <= infixString[i] && infixString[i] <= 'Z' || infixString[i] == '_' || '0' <= infixString[i] && infixString[i] <= '9') {
				variableBytes = append(variableBytes, infixString[i])
				i++
			}
			variable := string(variableBytes)
			var ok1, ok2 bool
			var value int
			if _, ok1 = request.InputData[variable]; ok1 {
				value = request.InputData[variable]
			}
			if _, ok2 = request.OutputData[variable]; ok2 {
				value = request.OutputData[variable]
			}
			if !ok1 && !ok2 {
				log.SugarLogger().Errorf("variable [%s] not found in data map", variable)
				return nil, errors.New("variable not found in input data map")
			}
			tmpTreeNode := &TreeNode{Type: "number", Value: value}
			postfixList.PushBack(tmpTreeNode)
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
					return nil, errors.New("and operator err")
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
					log.SugarLogger().Error("or operator err")
					return nil, errors.New("or operator err")
				}
			}
			pushOperator("||", operatorList, postfixList)
			i += 2
			continue
		}
		//等于判断
		if infixString[i] == '=' {
			if i+1 < len(infixString) {
				if infixString[i+1] != '=' {
					log.SugarLogger().Error("= operator err")
					return nil, errors.New("= operator err")
				}
			}
			pushOperator("==", operatorList, postfixList)
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
	return postfixList, nil
}

func GeneratorGenPostfixList(infixString string, request *RuleNodeRequest) (*list.List, error) {
	postfixList, err := GenPostfixList(infixString, request)
	if err != nil {
		log.SugarLogger().Error(err)
		return nil, err
	}
	return postfixList, nil
}

func ExecutePostfixList(postfixList *list.List) (any, error) {
	executeStack := list.New()
	for postfixList.Front() != nil {
		e := postfixList.Front()
		postfixList.Remove(e)
		if e.Value.(*TreeNode).Type == "number" {
			operatorValue := e.Value.(*TreeNode).Value.(int)
			executeStack.PushFront(operatorValue)
		} else { //only binary operator
			operator := e.Value.(*TreeNode).Value.(string)
			operatorElement1 := executeStack.Front()
			operatorElement2 := operatorElement1.Next()
			executeStack.Remove(operatorElement1)
			executeStack.Remove(operatorElement2)
			operatorValue1 := operatorElement1.Value
			operatorValue2 := operatorElement2.Value
			result := Operator(operator).DoubleOperatorFunc(operatorValue2, operatorValue1)
			executeStack.PushFront(result)
		}
	}
	resultElement := executeStack.Front()
	executeStack.Remove(resultElement)
	return resultElement.Value, nil
}

func ExecuteInfixString(infixString string, request *RuleNodeRequest) (any, error) {
	postfixList, err := GeneratorGenPostfixList(infixString, request)
	if err != nil {
		log.SugarLogger().Error(err)
		return nil, err
	}
	return ExecutePostfixList(postfixList)
}
