package mainserver

import (
	"deps/log4go"
	"encoding/json"
	"fmt"
	"github.com/wg815737157/paper-work/internal/mainserver/db"
	"github.com/wg815737157/paper-work/internal/mainserver/db/model"
	internalpkg "github.com/wg815737157/paper-work/internal/pkg"
	"gorm.io/gorm"
	"log"
)

func Init() {

	//	Init db

	//	Init redis

}

func Run() {
	localdb := db.GetLocalDB()
	log4go.Info("d")
	tree := &internalpkg.Tree{}
	node := &internalpkg.Node{}
	tree.Nodes = append(tree.Nodes, node)
	s := "{\n    \"nodes\":[\n        {\n            \"nodeId\":0,\n            \"nodeType\":\"start\",\n            \"nodeName\":\"start\",\n            \"pnodes\":[\n\n            ],\n            \"cnodes\":[\n                1,\n                2\n            ],\n            \"pnodeConditions\":{\n\n            },\n            \"mergeType\":\"\",\n            \"ruleNameList\":[\n\n            ],\n            \"result\":\"\"\n        },\n        {\n            \"nodeId\":1,\n            \"nodeType\":\"normal\",\n            \"nodeName\":\"gongan_node\",\n            \"pnodes\":[\n                0\n            ],\n            \"cnodes\":[\n                3,\n                4,\n                5\n            ],\n            \"pnodeConditions\":{\n\n            },\n            \"mergeType\":\"\",\n            \"ruleNameList\":[\n                \"gongan\"\n            ],\n            \"result\":\"\"\n        },\n        {\n            \"nodeId\":2,\n            \"nodeType\":\"normal\",\n            \"nodeName\":\"fayuan_node\",\n            \"pnodes\":[\n                0\n            ],\n            \"cnodes\":[\n                3,\n                4,\n                5\n            ],\n            \"pnodeConditions\":{\n\n            },\n            \"mergeType\":\"\",\n            \"ruleNameList\":[\n                \"fayuan\"\n            ],\n            \"result\":\"\"\n        },\n        {\n            \"nodeId\":3,\n            \"nodeType\":\"end\",\n            \"nodeName\":\"end\",\n            \"pnodes\":[\n                1,\n                2\n            ],\n            \"cnodes\":[\n\n            ],\n            \"pnodeConditions\":{\n                \"1\":\"fail\",\n                \"2\":\"fail\"\n            },\n            \"mergeType\":\"or\",\n            \"result\":\"\"\n        },\n        {\n            \"nodeId\":4,\n            \"nodeType\":\"normal\",\n            \"nodeName\":\"renhang_credit\",\n            \"pnodes\":[\n                1,\n                2\n            ],\n            \"cnodes\":[\n                6,\n                7\n            ],\n            \"pnodeConditions\":{\n                \"1\":\"pass\",\n                \"2\":\"pass\"\n            },\n            \"mergeType\":\"and\",\n            \"result\":\"\"\n        },\n        {\n            \"nodeId\":5,\n            \"nodeType\":\"normal\",\n            \"nodeName\":\"company_credit\",\n            \"pnodes\":[\n                1,\n                2\n            ],\n            \"cnodes\":[\n                6,\n                7\n            ],\n            \"pnodeConditions\":{\n                \"1\":\"pass\",\n                \"2\":\"pass\"\n            },\n            \"mergeType\":\"and\",\n            \"result\":\"\"\n        },\n        {\n            \"nodeId\":6,\n            \"nodeType\":\"end\",\n            \"nodeName\":\"end\",\n            \"pnodes\":[\n                4,\n                5\n            ],\n            \"cnodes\":[\n\n            ],\n            \"pnodeConditions\":{\n                \"4\":\"fail\",\n                \"5\":\"fail\"\n            },\n            \"mergeType\":\"or\",\n            \"result\":\"\"\n        },\n        {\n            \"nodeId\":7,\n            \"nodeType\":\"normal\",\n            \"nodeName\":\"\",\n            \"pnodes\":[\n                4,\n                5\n            ],\n            \"cnodes\":[\n\n            ],\n            \"pnodeConditions\":{\n                \"4\":\"fail\",\n                \"5\":\"fail\"\n            },\n            \"mergeType\":\"and\",\n            \"result\":\"\"\n        }\n    ],\n    \"inputData\":\"\",\n    \"outputData\":\"\"\n}"
	err := json.Unmarshal([]byte(s), tree)
	if err != nil {
		log.Fatalln(err)
	}
	treeString, err := json.Marshal(tree)
	fmt.Println(string(treeString))
	var m []*model.SysTree
	err = localdb.Raw("select * from sys_tree where id =1").Scan(&m).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Fatalln(err)
	}
	result := localdb.Exec("insert into sys_tree (`sys_id`,`sys_name`,`tree`)values(?,?,?)", 1, "主业务线", treeString)
	fmt.Println(result)

	//ginEngine := gin.New()
	//ginEngine.Handler()
	//err := ginEngine.Run("localhost:3747")
	//if err != nil {
	//	log.Fatalln(err)
	//}
}
