package controllers

import (
	"agent_learning/middleware"
	"agent_learning/pkg/response"
	"agent_learning/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ParcelController struct {
	parcelSvc *services.ParcelService
}

func NewParcelController(parcelSvc *services.ParcelService) *ParcelController {
	return &ParcelController{parcelSvc: parcelSvc}
}

func (ctl *ParcelController) Inbound(c *gin.Context) {
	adminID := middleware.CurrentUserID(c)
	var req struct {
		TrackingNo     string `json:"tracking_no" binding:"required"`
		Location       string `json:"location" binding:"required"`
		RecipientPhone string `json:"recipient_phone" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}

	parcel, err := ctl.parcelSvc.Inbound(adminID, services.InboundParcelInput{
		TrackingNo:     req.TrackingNo,
		Location:       req.Location,
		RecipientPhone: req.RecipientPhone,
	})
	if err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}
	response.Success(c, parcel)
}

func (ctl *ParcelController) Outbound(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	username := middleware.CurrentUsername(c)
	var req struct {
		TrackingNo string `json:"tracking_no" binding:"required"`
		PickupCode string `json:"pickup_code"`
		PickupName string `json:"pickup_name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}
	if req.PickupName == "" {
		req.PickupName = username
	}

	parcel, err := ctl.parcelSvc.Outbound(userID, services.OutboundParcelInput{
		TrackingNo: req.TrackingNo,
		PickupCode: req.PickupCode,
		PickupName: req.PickupName,
	})
	if err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}
	response.Success(c, parcel)
}

func (ctl *ParcelController) Status(c *gin.Context) {
	trackingNo := c.Query("tracking_no")
	parcel, err := ctl.parcelSvc.GetStatus(trackingNo)
	if err != nil {
		response.Fail(c, http.StatusNotFound, response.CodeNotFound, err.Error())
		return
	}
	response.Success(c, parcel)
}

func (ctl *ParcelController) List(c *gin.Context) {
	parcels, err := ctl.parcelSvc.ListAll()
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeServerError, err.Error())
		return
	}
	response.Success(c, parcels)
}
