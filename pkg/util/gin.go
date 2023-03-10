package util

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HealthCheck(e gin.IRouter) {
	e.GET("/health_check", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})
}

func GinGroup(path string, r gin.IRouter) *gin.RouterGroup {
	g := r.Group(path)
	HealthCheck(g)
	return g
}
