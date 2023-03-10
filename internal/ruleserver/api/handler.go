package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wg815737157/paper-work/internal/mainserver/db/model"
	internalpkg "github.com/wg815737157/paper-work/internal/pkg"
	"github.com/wg815737157/paper-work/internal/ruleserver/db"
	"github.com/wg815737157/paper-work/pkg/controller"
	"github.com/wg815737157/paper-work/pkg/log"
	"github.com/wg815737157/paper-work/pkg/util"
	"gorm.io/gorm"
	"strconv"
)

type ruleServerHandler struct {
}

func (h *ruleServerHandler) RuleId(c *controller.Controller) {
	inputData := map[string]int{"a": 1, "b": 2, "c": 3}
	postfixList := internalpkg.GeneratorGenPostfixList("c-2==a", inputData)
	//PrintTreeNodeList(postfixList)
	fmt.Println(internalpkg.ExecutePostfixList(postfixList))
}

func (h *ruleServerHandler) SysTree(c *controller.Controller) {
	sysIdStr, _ := c.GetQuery("sys_id")
	sysId, _ := strconv.Atoi(sysIdStr)
	localdb := db.GetLocalDB()
	//tree := internalpkg.Tree{}
	var m []*model.SysTree
	err := localdb.Raw("select * from sys_tree where id =?", sysId).Scan(&m).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.SugarLogger().Error(err)
		return
	}
	c.JSONRaw(m[0].Tree)
	return
}

func LoadHandlers(r gin.IRouter) {
	rsh := &ruleServerHandler{}
	util.HealthCheck(r)
	r.POST("/rule_id", controller.Warpper(rsh.RuleId))
	r.GET("/sys_tree", controller.Warpper(rsh.SysTree))
}
