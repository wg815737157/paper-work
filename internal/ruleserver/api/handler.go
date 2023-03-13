package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wg815737157/paper-work/internal/mainserver/db/model"
	internalpkg "github.com/wg815737157/paper-work/internal/pkg"
	"github.com/wg815737157/paper-work/internal/ruleserver/db"
	"github.com/wg815737157/paper-work/pkg/controller"
	"github.com/wg815737157/paper-work/pkg/log"
	"github.com/wg815737157/paper-work/pkg/util"
	"gorm.io/gorm"
	"io"
	"strconv"
)

type ruleServerHandler struct {
}

//基础数据
//TailpayAmount
//Income
//法院执行
//SlIdCourtExecuted
//公安犯罪
//CrimeResult
//企业征信
//TaxInformationNum
//人行征信
//DebitAmountOfOverdue
//LoanDebitOverdue
//CarLoanInDebt
//手机征信
//PhoneInfo
//京东金融
//百度金融
//腾讯金融

func (h *ruleServerHandler) RuleId(c *controller.Controller) {
	logger := log.SugarLogger()
	bodyByte, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.Error("err:", err)
		c.Failed(-1, err.Error())
		return
	}
	json.Unmarshal(bodyByte)
	inputData := map[string]int{"a": 1, "b": 2, "c": 3}
	postfixList := internalpkg.GeneratorGenPostfixList("c-2==a", inputData)
	//PrintTreeNodeList(postfixList)
	fmt.Println(internalpkg.ExecutePostfixList(postfixList))
}

func (h *ruleServerHandler) SysTree(c *controller.Controller) {
	sysIdStr, _ := c.GetQuery("sys_id")
	sysId, _ := strconv.Atoi(sysIdStr)
	localdb := db.GetLocalDB()
	var m []*model.SysTree
	err := localdb.Raw("select * from sys_tree where id =?", sysId).Scan(&m).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.SugarLogger().Error(err)
		return
	}
	resByte, _ := json.Marshal(m[0].Tree)
	c.SuccessWithData(string(resByte))
	return
}

func LoadHandlers(r gin.IRouter) {
	rsh := &ruleServerHandler{}
	util.HealthCheck(r)
	r.POST("/rule_id", controller.Warpper(rsh.RuleId))
	r.GET("/sys_tree", controller.Warpper(rsh.SysTree))
}
