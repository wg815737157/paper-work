package mainserver

import (
	"container/list"
	"encoding/json"
	"github.com/wg815737157/paper-work/internal/mainserver/db"
	"github.com/wg815737157/paper-work/internal/mainserver/db/model"
	internalpkg "github.com/wg815737157/paper-work/internal/pkg"
	"github.com/wg815737157/paper-work/pkg/log"
	"gorm.io/gorm"
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

func Run() {
	//日志
	logger := log.SugarLogger()
	//数据库
	localdb := db.GetLocalDB()
	tree := internalpkg.Tree{}
	var m []*model.SysTree
	err := localdb.Raw("select * from sys_tree where id =3").Scan(&m).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Error(err)
		return
	}

	err = json.Unmarshal([]byte(m[0].Tree), &tree)
	if err != nil {
		logger.Error(err)
		return
	}
	myList := list.New()
	myList.PushBack(0)
	for myList.Len() > 0 {
		curLevelLength := myList.Len()
		//层序遍历
		for i := 0; i < curLevelLength; i++ {
			element := myList.Front()
			myList.Remove(element)
			if _, ok := element.Value.(int); !ok {
				logger.Error(err)
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

	//ginEngine := gin.New()
	//ginEngine.Handler()
	//err := ginEngine.Run("localhost:3747")
	//if err != nil {
	//	log.Fatalln(err)
	//}
}
