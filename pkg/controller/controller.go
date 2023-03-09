package controller

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/wg815737157/paper-work/pkg/log"
	"net/http"
)

type (
	Controller struct {
		*gin.Context
	}

	Handller func(c *Controller)
)

func New(c *gin.Context) *Controller {
	return &Controller{c}
}

func Warpper(h Handller) gin.HandlerFunc {
	return func(context *gin.Context) {
		c := New(context)
		h(c)
	}
}

func (c *Controller) SuccessWithData(data interface{}) {
	jsonRespStr, _ := json.Marshal(data)
	log.Logger().Info(string(jsonRespStr))
	c.JSON(http.StatusOK, NewSuccessResult(data))
}

func (c *Controller) SuccessWithMsgData(msg string, data interface{}) {
	c.JSON(http.StatusOK, NewSuccessMsgResult(msg, data))
}

func (c *Controller) Success() {
	c.SuccessWithData(gin.H{})
}

func (c *Controller) Failed(code int, msg string) {
	jsonRespStr, _ := json.Marshal(NewFailedResult(code, msg))
	log.Logger().Error(string(jsonRespStr))
	c.JSON(http.StatusOK, NewFailedResult(code, msg))
}

func (c *Controller) InternalServerError() {
	c.Failed(500, "Internal Server Error")
}

func (c *Controller) JSONRaw(rawData interface{}) {
	c.JSON(http.StatusOK, rawData)
}

func (c *Controller) Ctx() context.Context {
	return context.WithValue(c.Request.Context(), "tag", "controller")
}
