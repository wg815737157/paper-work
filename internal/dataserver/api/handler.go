package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/wg815737157/paper-work/pkg/controller"
	"github.com/wg815737157/paper-work/pkg/log"
	"github.com/wg815737157/paper-work/pkg/util"
)

type dataServerHandler struct {
}

//type allData struct {
//	tailPayAmount        int  尾款 // 基础数据
//	income               int  收入 // 基础数据
//	slIdCourtExecuted    int //法院
//	crime                int //公安
//	taxInformationNum    int // 企业征信
//	debitAmountOfOverdue int 当前逾期总额  //人行
//  numberOfOverdue      int 信贷历史逾期次数 //人行
//	phone2costShow       int 近两月消费  //手机  phone2costShow>20 || phone2costShow>income/100 phoneResult=1
//	badcreditJd   		 int 信贷不良   京东金融
//  badbehaviorJd 		 int 互联网恶意行为 京东金融  badcreditJ==0&&badbehaviorJd==0&&numberOfOverdue==0 jdResult=1
//	slCellNbankCaOverdue int 手机现金类逾期金额 // 百度金融
//	slIdNbankOtherRefuse int 身份证_非银其他拒绝 //百度金融  slIdNbankOtherRefuse==0&&slCellNbankCaOverdue==0&&debitAmountOfOverdue==0 bdResult=1
//}

func (h *dataServerHandler) getData(c *controller.Controller) {
	userDataMap := map[string]map[string]int{
		//通过
		"11": {"tailPayAmount": 100000,
			"income":               50000,
			"slIdCourtExecuted":    0,
			"crime":                0,
			"taxInformationNum":    0,
			"debitAmountOfOverdue": 10000,
			"numberOfOverdue":      0,
			"phone2costShow":       50,
			"badcreditJd":          1,
			"badbehaviorJd":        1,
			"slCellNbankCaOverdue": 1,
			"slIdNbankOtherRefuse": 1},
		//手机消费、平台金融信息违规被拒
		"12": {"tailPayAmount": 100000,
			"income":               50000,
			"slIdCourtExecuted":    0,
			"crime":                0,
			"taxInformationNum":    0,
			"debitAmountOfOverdue": 10000,
			"numberOfOverdue":      0,
			"phone2costShow":       10,
			"badcreditJd":          1,
			"badbehaviorJd":        1,
			"slCellNbankCaOverdue": 1,
			"slIdNbankOtherRefuse": 1},
		//行政违法拒绝
		"13": {
			"tailPayAmount":        100000,
			"income":               50000,
			"slIdCourtExecuted":    1,
			"crime":                0,
			"taxInformationNum":    0,
			"debitAmountOfOverdue": 10000,
			"numberOfOverdue":      0,
			"phone2costShow":       50,
			"badcreditJd":          1,
			"badbehaviorJd":        1,
			"slCellNbankCaOverdue": 1,
			"slIdNbankOtherRefuse": 1,
		},
		//企业征信被拒
		"21": {"tailPayAmount": 100000,
			"income":               50000,
			"slIdCourtExecuted":    0,
			"crime":                0,
			"taxInformationNum":    1,
			"debitAmountOfOverdue": 10000,
			"numberOfOverdue":      0,
			"phone2costShow":       50,
			"badcreditJd":          1,
			"badbehaviorJd":        1,
			"slCellNbankCaOverdue": 1,
			"slIdNbankOtherRefuse": 1},
		//	人行征信被拒
		"22": {"tailPayAmount": 100000,
			"income":               50000,
			"slIdCourtExecuted":    0,
			"crime":                0,
			"taxInformationNum":    0,
			"debitAmountOfOverdue": 10000,
			"numberOfOverdue":      3,
			"phone2costShow":       50,
			"badcreditJd":          1,
			"badbehaviorJd":        1,
			"slCellNbankCaOverdue": 1,
			"slIdNbankOtherRefuse": 1},
	}
	logger := log.SugarLogger()
	phone, _ := c.GetQuery("phone")
	IdNO, _ := c.GetQuery("id_no")
	userDataMapKey := IdNO + phone
	if _, ok := userDataMap[userDataMapKey]; !ok {
		logger.Error("user data not exists")
		c.Failed(-1, "user data not exists")
		return
	}
	userData := userDataMap[userDataMapKey]
	reqBody, err := c.GetRawData()
	if err != nil {
		logger.Error(err)
		c.Failed(-1, "body err")
		return
	}
	m := map[string]any{}
	err = json.Unmarshal(reqBody, &m)
	if err != nil {
		logger.Error(err)
		c.Failed(-1, err.Error())
		return
	}
	for k, _ := range m {
		m[k] = userData[k]
	}
	c.SuccessWithData(m)
}

func LoadHandlers(r gin.IRouter) {
	dsh := &dataServerHandler{}
	util.HealthCheck(r)
	r.POST("/get_data", controller.Warpper(dsh.getData))
}
