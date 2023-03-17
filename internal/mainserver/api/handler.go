package api

import (
	"container/list"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wg815737157/paper-work/config/mainconfig"
	internalpkg "github.com/wg815737157/paper-work/internal/pkg"
	"github.com/wg815737157/paper-work/pkg/controller"
	"github.com/wg815737157/paper-work/pkg/log"
	"github.com/wg815737157/paper-work/pkg/util"
)

func executeNode(tree *internalpkg.Tree, nodeId int, node *internalpkg.Node) error {
	logger := log.SugarLogger()
	defer func() {
		logger.Debugf("current node %d,current node content %v", nodeId, node)
		tree.Nodes[nodeId] = node
	}()

	if node.NodeType == "start" {
		node.IsUsed = true
		node.IsSatisfied = true
		node.IsSuccessful = true
		return nil
	}
	//普通节点，结束节点验证满足条件
	if node.NodeType == "normal" {
		node.IsUsed = true
		for _, pNodeId := range node.Pnodes {
			if node.PnodeConditions[pNodeId] == tree.Nodes[pNodeId].Result {
				continue
			}
			node.IsSatisfied = false
			return nil
		}
		node.IsSatisfied = true
		//执行节点内容
		defaultContext := context.Background()
		ruleData := tree.InputData
		for k, v := range tree.OutputData {
			ruleData[k] = v
		}
		rawMap := map[string]interface{}{"nodeId": node.NodeId, "nodeName": node.NodeName, "ruleIdList": node.RuleIdList, "inputData": ruleData, "outputData": map[string]int{}}
		rawBody, _ := json.Marshal(rawMap)
		//请求rule server
		url := mainconfig.GlobalConfig.RuleServer.Url + "/rule_id"
		ruleNodeResBytes := util.PostWithContext(defaultContext, url, rawBody)
		ruleNodeResponse := &internalpkg.RuleNodeResponse{}
		fmt.Println(string(ruleNodeResBytes))
		err := json.Unmarshal(ruleNodeResBytes, ruleNodeResponse)
		if err != nil {
			logger.Error(err)
			return err
		}
		if ruleNodeResponse.Code != 0 {
			logger.Error(ruleNodeResponse.Message)
			return errors.New(ruleNodeResponse.Message)
		}
		// 获取最后一个规则的结果作为最后结果
		node.Result = ruleNodeResponse.RuleResponseData.RuleResultList[len(ruleNodeResponse.RuleResponseData.RuleResultList)-1]
		// 获取输出数据
		for k, v := range ruleNodeResponse.RuleResponseData.OutputData {
			tree.OutputData[k] = v
		}
		node.IsSuccessful = true
		return nil
	}
	if node.NodeType == "end" {
		node.IsUsed = true
		for _, pNodeId := range node.Pnodes {
			if node.PnodeConditions[pNodeId] == tree.Nodes[pNodeId].Result {
				continue
			}
			node.IsSatisfied = false
			return nil
		}
		node.IsSatisfied = true
		node.IsSuccessful = true
		tree.TreeResult = "fail"
		return nil
	}
	//node.NodeType == "final"
	node.IsUsed = true
	for _, pNodeId := range node.Pnodes {
		if node.PnodeConditions[pNodeId] == tree.Nodes[pNodeId].Result {
			continue
		}
		node.IsSatisfied = false
		return nil
	}
	node.IsSatisfied = true
	node.IsSuccessful = true
	tree.TreeResult = "pass"
	return nil
}

func riskCheck(tree *internalpkg.Tree) (string, error) {
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
				return "", errors.New("节点类型异常")
			}
			nodeId := element.Value.(int)
			node := tree.Nodes[nodeId]

			//执行节点
			err := executeNode(tree, nodeId, node)
			if err != nil {
				logger.Error(err)
				return "", err
			}

			//节点执行成功加入子节点到队列
			if node.IsSuccessful {
				for _, nodeId := range node.Cnodes {
					myList.PushBack(nodeId)
				}
			}
		}
	}
	return tree.TreeResult, nil
}

type mainServerHandler struct {
}

func (h *mainServerHandler) riskCheck(c *controller.Controller) {
	logger := log.SugarLogger()
	sysIdStr, _ := c.GetQuery("sys_id")
	ruleTreeResponse := internalpkg.RuleTreeResponse{}
	ctx := context.Background()
	url := mainconfig.GlobalConfig.RuleServer.Url + "/sys_tree?sys_id=" + sysIdStr
	resBytes := util.GetWithContext(ctx, url)
	err := json.Unmarshal(resBytes, &ruleTreeResponse)
	if err != nil {
		logger.Error(err)
		c.Failed(-1, err.Error())
		return
	}
	result, err := riskCheck(ruleTreeResponse.Tree)
	if err != nil {
		logger.Error(err)
		c.Failed(-1, err.Error())
		return
	}
	c.SuccessWithData(map[string]string{"data": result})
	return
}

func LoadHandlers(r gin.IRouter) {
	msh := &mainServerHandler{}
	util.HealthCheck(r)
	r.POST("/risk_check", controller.Warpper(msh.riskCheck))
}
