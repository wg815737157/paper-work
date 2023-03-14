package api

import (
	"encoding/json"
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

func ExecuteRule(ruleRequest *internalpkg.RuleRequest, ruleResponse *internalpkg.RuleResponse, rules []internalpkg.Rule) {
	//按顺序执行
	for _, ruleName := range ruleRequest.RuleNameList {
		for _, rule := range rules {
			if rule.RuleName != ruleName {
				continue
			}
			//	规则执行体
			ruleResponse.RuleIdList = append(ruleResponse.RuleIdList, rule.Id)
			ruleResponse.RuleNameList = append(ruleResponse.RuleNameList, ruleName)

			if internalpkg.ExecuteInfixString(rule.RuleDetail, ruleRequest.InputData).(bool) {
				for _, ruleReturn := range rule.RuleReturns {
					switch ruleReturn.Type {
					case "int":
						ruleRequest.OutputData[ruleReturn.Name] = ruleReturn.Value.(int)
					case "expression":
						ruleRequest.OutputData[ruleReturn.Name] = internalpkg.ExecuteInfixString(ruleReturn.Value.(string), ruleRequest.InputData).(int)
					}
				}
				ruleResponse.RuleResultList = append(ruleResponse.RuleResultList, "pass")
			} else {
				ruleResponse.RuleResultList = append(ruleResponse.RuleResultList, "fail")
				return
			}
		}
	}
}

func (h *ruleServerHandler) RuleId(c *controller.Controller) {
	logger := log.SugarLogger()
	bodyBytes, err := c.GetRawData()
	if err != nil {
		logger.Error("err:", err)
		c.Failed(-1, err.Error())
		return
	}
	ruleRequest := &internalpkg.RuleRequest{}
	err = json.Unmarshal(bodyBytes, ruleRequest)
	if err != nil {
		logger.Error("err:", err)
		c.Failed(-1, err.Error())
		return
	}
	rules := []internalpkg.Rule{}
	localdb := db.GetLocalDB()
	err = localdb.Raw("select * from rule_info where rule_name in ?", ruleRequest.RuleNameList).Scan(&rules).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Error(err)
		c.Failed(-1, err.Error())
		return
	}

	ruleResponse := &internalpkg.RuleResponse{}
	ruleResponse.InputData = ruleRequest.InputData
	ExecuteRule(ruleRequest, ruleResponse, rules)
	c.SuccessWithData(ruleResponse)
	return
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
