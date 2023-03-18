package api

import (
	"container/list"
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/wg815737157/paper-work/config/mainconfig"
	internalpkg "github.com/wg815737157/paper-work/internal/pkg"
	"github.com/wg815737157/paper-work/pkg/controller"
	"github.com/wg815737157/paper-work/pkg/log"
	"github.com/wg815737157/paper-work/pkg/util"
)

func checkConditions(tree *internalpkg.Tree, node *internalpkg.Node) bool {
	var isSatisfied bool
	if node.MergeType == "and" {
		isSatisfied = true
		for _, pNodeId := range node.Pnodes {
			if node.PnodeConditions[pNodeId] != tree.Nodes[pNodeId].Result {
				isSatisfied = false
				break
			}
		}
	} else {
		isSatisfied = false
		for _, pNodeId := range node.Pnodes {
			if node.PnodeConditions[pNodeId] == tree.Nodes[pNodeId].Result {
				isSatisfied = true
				break
			}
		}
	}
	return isSatisfied
}
func executeNode(tree *internalpkg.Tree, nodeId int) (isEnd bool, err error) {
	logger := log.SugarLogger()
	var node *internalpkg.Node
	node = tree.Nodes[nodeId]

	defer func() {
		nodeStr, _ := json.Marshal(node)
		logger.Debugf("current node %d,current node content %s", nodeId, nodeStr)
		tree.Nodes[nodeId] = node
	}()

	if node.NodeType == "start" {
		node.IsUsed = true
		node.IsSatisfied = true
		node.IsSuccessful = true
		return false, nil
	}
	//普通节点，结束节点验证满足条件
	if node.NodeType == "normal" {
		node.IsUsed = true
		node.IsSatisfied = checkConditions(tree, node)
		if node.IsSatisfied == false {
			return false, nil
		}
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
		err = json.Unmarshal(ruleNodeResBytes, ruleNodeResponse)
		if err != nil {
			logger.Error(err)
			return false, err
		}
		if ruleNodeResponse.Code != 0 {
			logger.Error(ruleNodeResponse.Message)
			return false, errors.New(ruleNodeResponse.Message)
		}
		// 获取最后一个规则的结果作为最后结果
		node.Result = ruleNodeResponse.RuleResponseData.RuleResultList[len(ruleNodeResponse.RuleResponseData.RuleResultList)-1]
		// 获取输出数据
		for k, v := range ruleNodeResponse.RuleResponseData.OutputData {
			tree.OutputData[k] = v
		}
		node.IsSuccessful = true
		return false, nil
	}
	if node.NodeType == "end" {
		node.IsUsed = true
		node.IsSatisfied = checkConditions(tree, node)
		if node.IsSatisfied == false {
			return false, nil
		}
		node.IsSuccessful = true
		tree.TreeResult = "fail"
		logger.Info("规则未通过，终止流程")
		return true, nil
	}
	//node.NodeType == "final"
	node.IsUsed = true
	node.IsSatisfied = checkConditions(tree, node)
	if node.IsSatisfied == false {
		return false, nil
	}
	node.IsSuccessful = true
	tree.TreeResult = "pass"
	logger.Info("流程执行完成，审批通过")
	return true, nil
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

			// 检测节点未执行过
			if node.IsSatisfied {
				logger.Infof("node [%d] is executed", node.NodeId)
				continue
			}

			//执行节点
			isEnd, err := executeNode(tree, nodeId)
			if err != nil {
				logger.Error(err)
				return tree.TreeResult, err
			}

			if isEnd == true {
				logger.Info("流程结束")
				return tree.TreeResult, err
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
	phone, _ := c.GetQuery("phone")
	IdNO, _ := c.GetQuery("id_no")
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

	ctx = context.Background()
	dataServerUrl := mainconfig.GlobalConfig.DataServer.Url + "/get_data?phone=" + phone + "&id_no=" + IdNO
	bodyRaw, _ := json.Marshal(ruleTreeResponse.Tree.InputData)
	dataServerBytes := util.PostWithContext(ctx, dataServerUrl, bodyRaw)

	dataRes := internalpkg.DataServerResponse{}
	err = json.Unmarshal(dataServerBytes, &dataRes)
	if err != nil {
		logger.Error(err)
		c.Failed(-1, err.Error())
		return
	}
	if dataRes.Code != 0 {
		logger.Error(dataRes.Message)
		c.Failed(-1, dataRes.Message)
		return
	}
	ruleTreeResponse.Tree.InputData = dataRes.Data

	result, err := riskCheck(ruleTreeResponse.Tree)
	if err != nil {
		logger.Error(err)
		c.Failed(-1, err.Error())
		return
	}
	c.SuccessWithData(result)
	return
}

func LoadHandlers(r gin.IRouter) {
	msh := &mainServerHandler{}
	util.HealthCheck(r)
	r.POST("/risk_check", controller.Warpper(msh.riskCheck))
}
