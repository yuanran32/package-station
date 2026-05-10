package controllers

import (
	"agent_learning/middleware"
	"agent_learning/pkg/jwtutil"
	"agent_learning/pkg/response"
	"agent_learning/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userSvc    *services.UserService
	jwtManager *jwtutil.Manager
}

func NewUserController(userSvc *services.UserService, jwtManager *jwtutil.Manager) *UserController {
	return &UserController{
		userSvc:    userSvc,
		jwtManager: jwtManager,
	}
}

func (ctl *UserController) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Name     string `json:"name"`
		Phone    string `json:"phone" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}

	user, err := ctl.userSvc.Register(services.RegisterUserInput{
		Username: req.Username,
		Password: req.Password,
		Name:     req.Name,
		Phone:    req.Phone,
	})
	if err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}

	token, err := ctl.jwtManager.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeServerError, "token generate failed")
		return
	}
	response.Success(c, gin.H{
		"user":  user,
		"token": token,
	})
}

func (ctl *UserController) Login(c *gin.Context) {
	var req struct {
		Account  string `json:"account" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}

	user, err := ctl.userSvc.Login(req.Account, req.Password)
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
		return
	}

	token, err := ctl.jwtManager.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeServerError, "token generate failed")
		return
	}
	response.Success(c, gin.H{
		"user":  user,
		"token": token,
	})
}

func (ctl *UserController) GetProfile(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	user, err := ctl.userSvc.GetProfile(userID)
	if err != nil {
		response.Fail(c, http.StatusNotFound, response.CodeNotFound, "user not found")
		return
	}
	response.Success(c, user)
}

func (ctl *UserController) UpdateProfile(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	var req struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}

	user, err := ctl.userSvc.UpdateProfile(userID, services.UpdateUserProfileInput{
		Name:  req.Name,
		Phone: req.Phone,
	})
	if err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}
	response.Success(c, user)
}

func (ctl *UserController) PickupHistory(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	history, err := ctl.userSvc.GetPickupHistory(userID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeServerError, err.Error())
		return
	}
	response.Success(c, history)
}

// 身份码二维码接口
func (ctl *UserController) QRCode(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	identityCode, err := ctl.userSvc.GetIdentityCode(userID)
	if err != nil {
		response.Fail(c, http.StatusNotFound, response.CodeNotFound, "没有找到此用户")
		return
	}
	response.Success(c, gin.H{
		"identity_code": identityCode,
		"qr_text":       identityCode,
	})
}
