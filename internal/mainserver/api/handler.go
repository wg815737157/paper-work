package api

import (
	"container/list"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	internalpkg "github.com/wg815737157/paper-work/internal/pkg"
	"github.com/wg815737157/paper-work/pkg/controller"
	"github.com/wg815737157/paper-work/pkg/log"
	"github.com/wg815737157/paper-work/pkg/util"
)

func executeNode(tree *internalpkg.Tree, node *internalpkg.Node) {
	if node.NodeType == "start" {
		node.IsUsed = true
		node.IsSatisfied = true
		node.IsSuccessful = true
		return
	}
	//普通节点，结束节点验证满足条件
	if node.NodeType == "normal" {
		node.IsUsed = true
		for _, pNodeId := range node.Pnodes {
			if node.PnodeConditions[pNodeId] == tree.Nodes[pNodeId].Result {
				continue
			}
			node.IsSatisfied = false
			return
		}
		node.IsSatisfied = true
		//执行节点内容

		node.IsSuccessful = true
		return
	}
	if node.NodeType == "end" {
		node.IsUsed = true
		for _, pNodeId := range node.Pnodes {
			if node.PnodeConditions[pNodeId] == tree.Nodes[pNodeId].Result {
				continue
			}
			node.IsSatisfied = false
			return
		}
		node.IsSatisfied = true
		node.IsSuccessful = true
		return
	}

}

func riskCheck(tree internalpkg.Tree) {
	//日志
	logger := log.SugarLogger()

	myList := list.New()
	myList.PushBack(0)
	for myList.Len() > 0 {
		curLevelLength := myList.Len()
		//层序遍历
		for i := 0; i < curLevelLength; i++ {
			element := myList.Front()
			myList.Remove(element)
			if _, ok := element.Value.(int); !ok {
				logger.Error("节点类型异常")
				return
			}
			nodeId := element.Value.(int)
			node := tree.Nodes[nodeId]
			//执行节点
			executeNode(&tree, node)

			//节点执行成功加入子节点到队列
			if node.IsSuccessful {
				for _, nodeId := range node.Cnodes {
					myList.PushBack(nodeId)
				}
			}
		}
	}
}

type mainServerHandler struct {
}

func (h *mainServerHandler) riskCheck(c *controller.Controller) {
	logger := log.SugarLogger()
	sysIdStr, _ := c.GetQuery("sys_id")
	//sysId, _ := strconv.Atoi(sysIdStr)
	tree := internalpkg.Tree{}
	ctx := context.Background()
	url := fmt.Sprintf("http://localhost:10002/sys_tree?sys_id=%s", sysIdStr)
	resBytes := util.GetWithContext(ctx, url)
	fmt.Println(string(resBytes))

	err := json.Unmarshal(resBytes, &tree)
	if err != nil {
		logger.Error(err)
		return
	}
	riskCheck(tree)
	return
}

func LoadHandlers(r gin.IRouter) {
	msh := &mainServerHandler{}
	util.HealthCheck(r)
	r.POST("/risk_check", controller.Warpper(msh.riskCheck))
}
