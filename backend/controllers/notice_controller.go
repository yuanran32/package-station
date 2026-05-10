package controllers

import (
	"agent_learning/pkg/jwtutil"
	"agent_learning/pkg/response"
	"agent_learning/services"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type NoticeController struct {
	noticeSvc  *services.NoticeService
	hub        *services.NotifyHub
	jwtManager *jwtutil.Manager
}

func NewNoticeController(noticeSvc *services.NoticeService, hub *services.NotifyHub, jwtManager *jwtutil.Manager) *NoticeController {
	return &NoticeController{
		noticeSvc:  noticeSvc,
		hub:        hub,
		jwtManager: jwtManager,
	}
}

func (ctl *NoticeController) Send(c *gin.Context) {
	var req struct {
		UserID  *uint  `json:"user_id"`
		Type    string `json:"type"`
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}
	if strings.TrimSpace(req.Type) == "" {
		req.Type = "system"
	}

	notice, err := ctl.noticeSvc.SendNotice(req.UserID, req.Type, req.Content)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeServerError, err.Error())
		return
	}
	response.Success(c, notice)
}

func (ctl *NoticeController) WebSocketNotify(c *gin.Context) {
	token := strings.TrimSpace(c.Query("token"))
	if token == "" {
		auth := strings.TrimSpace(c.GetHeader("Authorization"))
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
			token = strings.TrimSpace(parts[1])
		}
	}
	if token == "" {
		response.Fail(c, http.StatusUnauthorized, response.CodeUnauthorized, "missing token")
		return
	}

	claims, err := ctl.jwtManager.ParseToken(token)
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, response.CodeUnauthorized, "invalid token")
		return
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	ctl.hub.Register(claims.UserID, conn)
	defer ctl.hub.Unregister(claims.UserID, conn)

	_ = conn.SetWriteDeadline(time.Now().Add(3 * time.Second))
	_ = conn.WriteJSON(gin.H{
		"type":      "connected",
		"message":   "ws notify connected",
		"user_id":   claims.UserID,
		"timestamp": time.Now(),
	})

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}
