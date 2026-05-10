package controllers

import (
	"agent_learning/middleware"
	"agent_learning/pkg/response"
	"agent_learning/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type PickupController struct {
	pickupSvc *services.PickupService
}

func NewPickupController(pickupSvc *services.PickupService) *PickupController {
	return &PickupController{pickupSvc: pickupSvc}
}

func (ctl *PickupController) GenerateCode(c *gin.Context) {
	operatorID := middleware.CurrentUserID(c)
	var req struct {
		TrackingNo     string `json:"tracking_no"`
		RecipientPhone string `json:"recipient_phone"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}

	parcel, err := ctl.pickupSvc.GeneratePickupCode(operatorID, services.GeneratePickupCodeInput{
		TrackingNo:     req.TrackingNo,
		RecipientPhone: req.RecipientPhone,
	})
	if err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}
	response.Success(c, parcel)
}

func (ctl *PickupController) RecordPickup(c *gin.Context) {
	operatorID := middleware.CurrentUserID(c)
	var req struct {
		TrackingNo     string `json:"tracking_no" binding:"required"`
		PickupUserID   *uint  `json:"pickup_user_id"`
		PickupUserName string `json:"pickup_user_name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}

	record, err := ctl.pickupSvc.RecordPickup(operatorID, services.RecordPickupInput{
		TrackingNo:     req.TrackingNo,
		PickupUserID:   req.PickupUserID,
		PickupUserName: req.PickupUserName,
	})
	if err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}
	response.Success(c, record)
}

func (ctl *PickupController) RecordDelivery(c *gin.Context) {
	operatorID := middleware.CurrentUserID(c)
	var req struct {
		TrackingNo     string `json:"tracking_no" binding:"required"`
		CourierName    string `json:"courier_name" binding:"required"`
		DeliveryStatus string `json:"delivery_status" binding:"required"`
		Remark         string `json:"remark"`
		DeliveredAt    string `json:"delivered_at"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}

	var deliveredAt *time.Time
	if req.DeliveredAt != "" {
		parsed, err := time.Parse(time.RFC3339, req.DeliveredAt)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, response.CodeBadRequest, "delivered_at format must be RFC3339")
			return
		}
		deliveredAt = &parsed
	}

	record, err := ctl.pickupSvc.RecordDelivery(operatorID, services.RecordDeliveryInput{
		TrackingNo:     req.TrackingNo,
		CourierName:    req.CourierName,
		DeliveryStatus: req.DeliveryStatus,
		Remark:         req.Remark,
		DeliveredAt:    deliveredAt,
	})
	if err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}
	response.Success(c, record)
}
