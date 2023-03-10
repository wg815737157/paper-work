package ruleserver

import (
	"github.com/gin-gonic/gin"
	"github.com/wg815737157/paper-work/config/ruleconfig"
	internalpkg "github.com/wg815737157/paper-work/internal/pkg"
	"github.com/wg815737157/paper-work/internal/ruleserver/api"
	"github.com/wg815737157/paper-work/internal/ruleserver/db"
	"github.com/wg815737157/paper-work/pkg/log"
	"net/http"
)

type defaultServer struct {
	*gin.Engine
	*http.Server
}

func DefaultServer() internalpkg.DefaultServerInterface {
	rs := &defaultServer{}
	rs.Engine = gin.Default()
	rs.Server = &http.Server{
		Addr:    ruleconfig.GlobalConfig.Port,
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
