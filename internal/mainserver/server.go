package mainserver

import (
	"github.com/gin-gonic/gin"
	"github.com/wg815737157/paper-work/config/mainconfig"
	"github.com/wg815737157/paper-work/internal/mainserver/api"
	"github.com/wg815737157/paper-work/internal/mainserver/db"
	internalpkg "github.com/wg815737157/paper-work/internal/pkg"
	"github.com/wg815737157/paper-work/pkg/log"
	"net/http"
)

// {"inputData":null, "outputData":null, "nodes":[{"NodeId":0, "NodeName":"start", "NodeType":"start", "IsUsed":false, "RuleIdList":null, "Pnodes":[], "PnodeConditions":{}, "MergeType":"", "IsSatisfied":false, "IsSuccessful":false, "Cnodes":[1, 2 ], "Result":""}, {"NodeId":1, "NodeName":"gongan_node", "NodeType":"normal", "IsUsed":false, "RuleIdList":null, "Pnodes":[0 ], "PnodeConditions":{}, "MergeType":"", "IsSatisfied":false, "IsSuccessful":false, "Cnodes":[3, 4, 5 ], "Result":""}, {"NodeId":2, "NodeName":"fayuan_node", "NodeType":"normal", "IsUsed":false, "RuleIdList":null, "Pnodes":[0 ], "PnodeConditions":{}, "MergeType":"", "IsSatisfied":false, "IsSuccessful":false, "Cnodes":[3, 4, 5 ], "Result":""}, {"NodeId":3, "NodeName":"end", "NodeType":"end", "IsUsed":false, "RuleIdList":null, "Pnodes":[1, 2 ], "PnodeConditions":{"1":"fail", "2":"fail"}, "MergeType":"or", "IsSatisfied":false, "IsSuccessful":false, "Cnodes":[], "Result":""}, {"NodeId":4, "NodeName":"renhang_credit", "NodeType":"normal", "IsUsed":false, "RuleIdList":null, "Pnodes":[1, 2 ], "PnodeConditions":{"1":"pass", "2":"pass"}, "MergeType":"and", "IsSatisfied":false, "IsSuccessful":false, "Cnodes":[6, 7 ], "Result":""}, {"NodeId":5, "NodeName":"company_credit", "NodeType":"normal", "IsUsed":false, "RuleIdList":null, "Pnodes":[1, 2 ], "PnodeConditions":{"1":"pass", "2":"pass"}, "MergeType":"and", "IsSatisfied":false, "IsSuccessful":false, "Cnodes":[6, 7 ], "Result":""}, {"NodeId":6, "NodeName":"end", "NodeType":"end", "IsUsed":false, "RuleIdList":null, "Pnodes":[4, 5 ], "PnodeConditions":{"4":"fail", "5":"fail"}, "MergeType":"or", "IsSatisfied":false, "IsSuccessful":false, "Cnodes":[], "Result":""}, {"NodeId":7, "NodeName":"", "NodeType":"normal", "IsUsed":false, "RuleIdList":null, "Pnodes":[4, 5 ], "PnodeConditions":{"4":"fail", "5":"fail"}, "MergeType":"and", "IsSatisfied":false, "IsSuccessful":false, "Cnodes":[], "Result":""} ] }
type defaultServer struct {
	*gin.Engine
	*http.Server
}

func DefaultServer() internalpkg.DefaultServerInterface {
	rs := &defaultServer{}
	rs.Engine = gin.Default()
	rs.Server = &http.Server{
		Addr:    mainconfig.GlobalConfig.Port,
		Handler: rs.Engine,
	}
	return rs
}

func (rs *defaultServer) Init() internalpkg.DefaultServerInterface {
	db.InitDB()
	api.LoadHandlers(rs.Engine)
	//	Init DB
	//	Init Redis
	return rs
}

func (rs *defaultServer) Run() {
	if err := rs.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.SugarLogger().Error(err)
		return
	}
}
