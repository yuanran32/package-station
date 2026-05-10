package middleware

import (
	"agent_learning/pkg/response"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[panic] %v", r)
				response.Fail(c, http.StatusInternalServerError, response.CodeServerError, "internal server error")
				c.Abort()
			}
		}()

		c.Next()
		if len(c.Errors) > 0 && !c.Writer.Written() {
			log.Printf("[gin error] %v", c.Errors.String())
			response.Fail(c, http.StatusInternalServerError, response.CodeServerError, "internal server error")
		}
	}
}
