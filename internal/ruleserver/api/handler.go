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
	"reflect"
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

func handleRuleResult(result any) int {
	if reflect.TypeOf(result).Kind() == reflect.Float64 {
		return int(result.(float64))
	}
	return result.(int)
}

func ExecuteRule(ruleRequest *internalpkg.RuleNodeRequest, ruleResponseData *internalpkg.RuleNodeResponseData, rules []internalpkg.Rule) error {
	//按顺序执行
	for _, ruleId := range ruleRequest.RuleIdList {
		for _, rule := range rules {
			if rule.Id != ruleId {
				continue
			}
			//	规则执行体
			ruleResponseData.RuleIdList = append(ruleResponseData.RuleIdList, rule.Id)
			ruleResponseData.RuleNameList = append(ruleResponseData.RuleNameList, rule.RuleName)
			ruleDetailResult, err := internalpkg.ExecuteInfixString(rule.RuleDetail, ruleRequest)
			if err != nil {
				log.SugarLogger().Error(err)
				return err
			}
			if ruleDetailResult.(bool) {
				log.SugarLogger().Debug(rule.RuleReturnPass)
				for _, ruleReturn := range rule.RuleReturnPass {
					switch ruleReturn.Type {
					case "int":
						returnValue := handleRuleResult(ruleReturn.Value)
						ruleRequest.OutputData[ruleReturn.Name] = returnValue
						ruleResponseData.OutputData[ruleReturn.Name] = returnValue
					case "expression":
						returnExpressionResult, err := internalpkg.ExecuteInfixString(ruleReturn.Value.(string), ruleRequest)
						if err != nil {
							log.SugarLogger().Error(err)
							return err
						}
						returnExpressionValue := handleRuleResult(returnExpressionResult)
						ruleRequest.OutputData[ruleReturn.Name] = returnExpressionValue
						ruleResponseData.OutputData[ruleReturn.Name] = returnExpressionValue
					}
				}
				ruleResponseData.RuleResultList = append(ruleResponseData.RuleResultList, "pass")
			} else {
				log.SugarLogger().Debug(rule.RuleReturnFail)
				for _, ruleReturn := range rule.RuleReturnFail {
					switch ruleReturn.Type {
					case "int":
						returnValue := handleRuleResult(ruleReturn.Value)
						ruleRequest.OutputData[ruleReturn.Name] = returnValue
						ruleResponseData.OutputData[ruleReturn.Name] = returnValue
					case "expression":
						returnExpressionResult, err := internalpkg.ExecuteInfixString(ruleReturn.Value.(string), ruleRequest)
						if err != nil {
							log.SugarLogger().Error(err)
							return err
						}
						returnExpressionValue := handleRuleResult(returnExpressionResult)
						ruleRequest.OutputData[ruleReturn.Name] = returnExpressionValue
						ruleResponseData.OutputData[ruleReturn.Name] = returnExpressionValue
					}
				}
				ruleResponseData.RuleResultList = append(ruleResponseData.RuleResultList, "fail")
				return nil
			}
		}
	}
	return nil
}

func (h *ruleServerHandler) RuleId(c *controller.Controller) {
	logger := log.SugarLogger()
	bodyBytes, err := c.GetRawData()
	if err != nil {
		logger.Error("err:", err)
		c.Failed(-1, err.Error())
		return
	}
	logger.Infof("请求的body:%s", string(bodyBytes))
	ruleNodeRequest := &internalpkg.RuleNodeRequest{}
	err = json.Unmarshal(bodyBytes, ruleNodeRequest)
	if err != nil {
		logger.Error("err:", err)
		c.Failed(-1, err.Error())
		return
	}
	mysqlRules, err := db.GetRuleInfoById(ruleNodeRequest.RuleIdList)
	if err != nil {
		logger.Error(err)
		c.Failed(-1, err.Error())
		return
	}
	if len(mysqlRules) == 0 {
		errMsg := fmt.Sprintf("rule id [%v] empty", ruleNodeRequest.RuleIdList)
		logger.Errorf(errMsg)
		c.Failed(-1, errMsg)
		return
	}
	rules := make([]internalpkg.Rule, len(mysqlRules))
	for i, mysqlRule := range mysqlRules {
		rules[i].Id = mysqlRule.Id
		rules[i].RuleName = mysqlRule.RuleName
		rules[i].RuleDetail = mysqlRule.RuleDetail
		err = json.Unmarshal([]byte(mysqlRule.RuleReturnPass), &rules[i].RuleReturnPass)
		if err != nil {
			c.Failed(-1, err.Error())
			return
		}
		err = json.Unmarshal([]byte(mysqlRule.RuleReturnFail), &rules[i].RuleReturnFail)
		if err != nil {
			c.Failed(-1, err.Error())
			return
		}
	}
	ruleResponseData := &internalpkg.RuleNodeResponseData{
		NodeId:    ruleNodeRequest.NodeId,
		NodeName:  ruleNodeRequest.NodeName,
		InputData: ruleNodeRequest.InputData, OutputData: map[string]int{},
	}
	err = ExecuteRule(ruleNodeRequest, ruleResponseData, rules)
	if err != nil {
		c.Failed(-1, err.Error())
		return
	}
	c.SuccessWithData(ruleResponseData)
	return
}

func (h *ruleServerHandler) SysTree(c *controller.Controller) {
	sysIdStr, _ := c.GetQuery("sys_id")
	sysId, _ := strconv.Atoi(sysIdStr)
	localdb := db.GetLocalDB()
	var m []*model.SysTree
	err := localdb.Raw("select * from sys_tree where sys_id =?", sysId).Scan(&m).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.SugarLogger().Error(err)
		return
	}
	tree := &internalpkg.Tree{}
	err = json.Unmarshal([]byte(m[0].Tree), tree)
	if err != nil {
		log.SugarLogger().Error(err)
		return
	}

	c.SuccessWithData(tree)
	return
}

func LoadHandlers(r gin.IRouter) {
	rsh := &ruleServerHandler{}
	util.HealthCheck(r)
	r.POST("/rule_id", controller.Warpper(rsh.RuleId))
	r.GET("/sys_tree", controller.Warpper(rsh.SysTree))
}
