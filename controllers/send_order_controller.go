package controllers

import (
	"agent_learning/middleware"
	"agent_learning/pkg/response"
	"agent_learning/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SendOrderController struct {
	sendSvc *services.SendOrderService
}

func NewSendOrderController(sendSvc *services.SendOrderService) *SendOrderController {
	return &SendOrderController{sendSvc: sendSvc}
}

func (ctl *SendOrderController) CreateOrder(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	var req struct {
		SenderName      string  `json:"sender_name" binding:"required"`
		SenderPhone     string  `json:"sender_phone" binding:"required"`
		SenderAddress   string  `json:"sender_address" binding:"required"`
		ReceiverName    string  `json:"receiver_name" binding:"required"`
		ReceiverPhone   string  `json:"receiver_phone" binding:"required"`
		ReceiverAddress string  `json:"receiver_address" binding:"required"`
		ItemInfo        string  `json:"item_info" binding:"required"`
		Weight          float64 `json:"weight" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}

	order, err := ctl.sendSvc.CreateOrder(userID, services.CreateSendOrderInput{
		SenderName:      req.SenderName,
		SenderPhone:     req.SenderPhone,
		SenderAddress:   req.SenderAddress,
		ReceiverName:    req.ReceiverName,
		ReceiverPhone:   req.ReceiverPhone,
		ReceiverAddress: req.ReceiverAddress,
		ItemInfo:        req.ItemInfo,
		Weight:          req.Weight,
	})
	if err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}
	response.Success(c, order)
}

func (ctl *SendOrderController) ListOrders(c *gin.Context) {
	status := c.Query("status")
	orders, err := ctl.sendSvc.ListOrders(status)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeServerError, err.Error())
		return
	}
	response.Success(c, orders)
}

func (ctl *SendOrderController) ProcessOrder(c *gin.Context) {
	adminID := middleware.CurrentUserID(c)
	var req struct {
		OrderNo     string `json:"order_no" binding:"required"`
		Action      string `json:"action" binding:"required"`
		CourierName string `json:"courier_name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}

	order, err := ctl.sendSvc.Process(adminID, services.ProcessSendOrderInput{
		OrderNo:     req.OrderNo,
		Action:      req.Action,
		CourierName: req.CourierName,
	})
	if err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}
	response.Success(c, order)
}
