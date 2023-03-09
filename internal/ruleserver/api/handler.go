package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wg815737157/paper-work/pkg/controller"
	"github.com/wg815737157/paper-work/pkg/utils"
)

type ruleServerHandler struct {
}

func (h *ruleServerHandler) RuleId(c *controller.Controller) {
	inputData := map[string]int{"a": 1, "b": 2, "c": 3}
	postfixList := GeneratorGenPostfixList("c-2==a", inputData)
	//PrintTreeNodeList(postfixList)
	fmt.Println(ExecutePostfixList(postfixList))
}

func (h *ruleServerHandler) RuleAll(c *controller.Controller) {

}

func LoadHandlers(r gin.IRouter) {
	rsh := &ruleServerHandler{}
	utils.HealthCheck(r)
	r.POST("/rule_id", controller.Warpper(rsh.RuleId))
	r.GET("/rule_all", controller.Warpper(rsh.RuleAll))
}
