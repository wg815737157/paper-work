package utils

import (
	"github.com/gin-gonic/gin"
	"strings"
)

type ClientInfo struct {
	ClientId string `form:"client_id" json:"client_id" xml:"client_id"`
	DeviceId string `form:"device_id" json:"device_id" xml:"device_id"`
	AppVer   string `form:"app_ver" json:"app_ver" xml:"app_ver"`
	Os       string `form:"os" json:"os" xml:"os"`
	TzName   string `form:"tz_name" json:"tz_name" xml:"tz_name"`
	TzOffset int64  `form:"tz_offset" json:"tz_offset" xml:"tz_offset"`
	SysLang  string `form:"sys_lang" json:"sys_lang" xml:"sys_lang"`
	AppLang  string `form:"app_lang" json:"app_lang" xml:"app_lang"`
}

func (ci *ClientInfo) IsUsApp() bool {
	return strings.HasPrefix(ci.ClientId, "pionex.us")
}

func GetClientInfo(c *gin.Context) (*ClientInfo, error) {
	client := new(ClientInfo)
	err := c.ShouldBind(client)
	return client, err
}
