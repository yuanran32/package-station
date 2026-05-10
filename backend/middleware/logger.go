package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)

		userID := CurrentUserID(c)
		log.Printf(
			"[http] %s %s status=%d latency=%s user_id=%d ip=%s",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			latency,
			userID,
			c.ClientIP(),
		)
	}
}
