package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	CodeSuccess      = 0
	CodeBadRequest   = 400
	CodeUnauthorized = 401
	CodeForbidden    = 403
	CodeNotFound     = 404
	CodeServerError  = 500
)

type Resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func JSON(c *gin.Context, httpStatus int, code int, msg string, data interface{}) {
	c.JSON(httpStatus, Resp{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

func Success(c *gin.Context, data interface{}) {
	JSON(c, http.StatusOK, CodeSuccess, "success", data)
}

func Fail(c *gin.Context, httpStatus int, code int, msg string) {
	JSON(c, httpStatus, code, msg, nil)
}
