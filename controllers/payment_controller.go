package controllers

import (
	"agent_learning/middleware"
	"agent_learning/pkg/response"
	"agent_learning/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	paySvc *services.PaymentService
}

func NewPaymentController(paySvc *services.PaymentService) *PaymentController {
	return &PaymentController{paySvc: paySvc}
}

func (ctl *PaymentController) Create(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	var req struct {
		RelatedType string  `json:"related_type" binding:"required"`
		OrderNo     string  `json:"order_no"`
		RelatedID   uint    `json:"related_id"`
		Amount      float64 `json:"amount"`
		BizDesc     string  `json:"biz_desc"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}

	payOrder, err := ctl.paySvc.CreatePayment(userID, services.CreatePaymentInput{
		RelatedType: req.RelatedType,
		OrderNo:     req.OrderNo,
		RelatedID:   req.RelatedID,
		Amount:      req.Amount,
		BizDesc:     req.BizDesc,
	})
	if err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}
	response.Success(c, payOrder)
}

func (ctl *PaymentController) Callback(c *gin.Context) {
	var req struct {
		PayNo   string `json:"pay_no" binding:"required"`
		Status  string `json:"status" binding:"required"`
		Method  string `json:"method"`
		Payload string `json:"payload"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}

	payOrder, err := ctl.paySvc.Callback(services.PaymentCallbackInput{
		PayNo:   req.PayNo,
		Status:  req.Status,
		Method:  req.Method,
		Payload: req.Payload,
	})
	if err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}
	response.Success(c, payOrder)
}

func (ctl *PaymentController) Bills(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	bills, err := ctl.paySvc.Bills(userID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeServerError, err.Error())
		return
	}
	response.Success(c, bills)
}
