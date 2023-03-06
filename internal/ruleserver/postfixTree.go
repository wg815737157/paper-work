package ruleserver

import (
	"container/list"
	"fmt"
)

type TreeNode struct {
	Type  string // 数值，操作符
	Value any
	Left  *TreeNode
	Right *TreeNode
}

func printTreeNodeList(treeNodeList *list.List) {
	for treeNodeList.Front() != nil {
		e := treeNodeList.Front()
		treeNodeList.Remove(e)
		if e.Value.(*TreeNode).Type == "number" {
			fmt.Println(e.Value.(*TreeNode).Value)
		} else {
			fmt.Println(string(e.Value.(*TreeNode).Value.(byte)))
		}
	}
}
func GenPostfixList(infixString string) *list.List {
	operatorLevel := map[byte]int{'(': 100, '<': 6, '>': 6, '+': 4, '-': 4, '*': 3, '/': 3}
	operatorList := list.New()
	postfixTreeList := list.New()
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
			postfixTreeList.PushBack(tmpTreeNode)
			continue
		}

		//操作符
		inputOperator := infixString[i]
		if inputOperator == '(' {
			operatorList.PushFront(inputOperator)
			i++
			continue
		}
		if inputOperator == ')' {
			tmpElement := operatorList.Front()
			operatorList.Remove(tmpElement)
			for tmpElement.Value.(byte) != '(' {
				tmpTreeNode := &TreeNode{Type: "operator", Value: tmpElement.Value.(byte)}
				postfixTreeList.PushBack(tmpTreeNode)
				tmpElement = operatorList.Front()
				operatorList.Remove(tmpElement)
			}
			i++
			continue
		}

		for {
			tmpElement := operatorList.Front()
			if tmpElement == nil {
				break
			}
			operatorInlist := tmpElement.Value.(byte)
			// 堆栈里操作符优先级大于要插入的操作符优先级
			if operatorLevel[inputOperator] < operatorLevel[operatorInlist] {
				break
			}
			operatorList.Remove(tmpElement)
			tmpTreeNode := &TreeNode{Type: "operator", Value: operatorInlist}
			postfixTreeList.PushBack(tmpTreeNode)
		}
		operatorList.PushFront(inputOperator)
		i++
		continue
	}
	for operatorList.Len() != 0 {
		tmpElement := operatorList.Front()
		operatorList.Remove(tmpElement)
		tmpTreeNode := &TreeNode{Type: "operator", Value: tmpElement.Value.(byte)}
		postfixTreeList.PushBack(tmpTreeNode)
	}
	return postfixTreeList
}

func printOperatorList(operatorList *list.List) {
	operator := operatorList.Front()
	operatorList.Remove(operator)
	fmt.Println(string(operator.Value.(byte)))
}

func GeneratorGenPostfixTree(infixString string) {
	treeNodeList := GenPostfixList(infixString)
	printTreeNodeList(treeNodeList)
}
