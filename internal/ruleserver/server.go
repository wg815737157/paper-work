package ruleserver

import (
	"github.com/gin-gonic/gin"
	"github.com/wg815737157/paper-work/config/ruleconfig"
	"github.com/wg815737157/paper-work/internal/pkg"
	"github.com/wg815737157/paper-work/internal/ruleserver/api"
	"github.com/wg815737157/paper-work/internal/ruleserver/db"
	"log"
	"net/http"
)

type ruleServer struct {
	*gin.Engine
	*http.Server
}

func DefaultServer() pkg.DefaultServer {
	rs := &ruleServer{}
	rs.Engine = gin.Default()
	rs.Server = &http.Server{
		Addr:    ruleconfig.GlobalConfig.Port,
		Handler: rs.Engine,
	}
	return rs
}

func (rs *ruleServer) Init() pkg.DefaultServer {
	db.InitDB()
	api.LoadHandlers(rs.Engine)
	//	Init DB
	//	Init Redis
	return rs
}

func (rs *ruleServer) Run() {
	if err := rs.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
