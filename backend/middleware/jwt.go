package middleware

import (
	"agent_learning/pkg/jwtutil"
	"agent_learning/pkg/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	CtxUserIDKey   = "ctx_user_id"
	CtxRoleKey     = "ctx_role"
	CtxUsernameKey = "ctx_username"
)

func JWTAuth(jwtManager *jwtutil.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			response.Fail(c, http.StatusUnauthorized, response.CodeUnauthorized, "缺失token")
			c.Abort()
			return
		}

		claims, err := jwtManager.ParseToken(token)
		if err != nil {
			response.Fail(c, http.StatusUnauthorized, response.CodeUnauthorized, "无效token")
			c.Abort()
			return
		}

		c.Set(CtxUserIDKey, claims.UserID)
		c.Set(CtxRoleKey, claims.Role)
		c.Set(CtxUsernameKey, claims.Username)
		c.Next()
	}
}

func RoleAuth(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}
	return func(c *gin.Context) {
		role := CurrentRole(c)
		if _, ok := allowed[role]; !ok {
			response.Fail(c, http.StatusForbidden, response.CodeForbidden, "permission denied")
			c.Abort()
			return
		}
		c.Next()
	}
}

func CurrentUserID(c *gin.Context) uint {
	val, ok := c.Get(CtxUserIDKey)
	if !ok {
		return 0
	}
	userID, ok := val.(uint)
	if !ok {
		return 0
	}
	return userID
}

func CurrentRole(c *gin.Context) string {
	val, ok := c.Get(CtxRoleKey)
	if !ok {
		return ""
	}
	role, ok := val.(string)
	if !ok {
		return ""
	}
	return role
}

func CurrentUsername(c *gin.Context) string {
	val, ok := c.Get(CtxUsernameKey)
	if !ok {
		return ""
	}
	username, ok := val.(string)
	if !ok {
		return ""
	}
	return username
}

func extractToken(c *gin.Context) string {
	authorization := strings.TrimSpace(c.GetHeader("Authorization"))
	if authorization != "" {
		parts := strings.SplitN(authorization, " ", 2)
		if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
			return strings.TrimSpace(parts[1])
		}
	}
	return strings.TrimSpace(c.Query("token"))
}
