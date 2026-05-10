package controllers

import (
	"agent_learning/middleware"
	"agent_learning/pkg/response"
	"agent_learning/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CouponController struct {
	couponSvc *services.CouponService
}

func NewCouponController(couponSvc *services.CouponService) *CouponController {
	return &CouponController{couponSvc: couponSvc}
}

func (ctl *CouponController) AdminList(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeBadRequest, "page must be integer")
		return
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeBadRequest, "page_size must be integer")
		return
	}

	result, err := ctl.couponSvc.AdminListCoupons(services.AdminCouponListInput{
		Status:   c.Query("status"),
		Keyword:  c.Query("keyword"),
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeServerError, err.Error())
		return
	}
	response.Success(c, result)
}

func (ctl *CouponController) AdminCreate(c *gin.Context) {
	adminID := middleware.CurrentUserID(c)
	var req struct {
		Name         string  `json:"name" binding:"required"`
		Amount       float64 `json:"amount" binding:"required"`
		Code         string  `json:"code"`
		ActivityRule string  `json:"activity_rule"`
		Threshold    float64 `json:"threshold"`
		Total        int     `json:"total"`
		ValidDays    int     `json:"valid_days"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}

	coupon, err := ctl.couponSvc.CreateCoupon(adminID, services.CreateCouponInput{
		Name:         req.Name,
		Amount:       req.Amount,
		Code:         req.Code,
		ActivityRule: req.ActivityRule,
		Threshold:    req.Threshold,
		Total:        req.Total,
		ValidDays:    req.ValidDays,
	})
	if err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}
	response.Success(c, coupon)
}

func (ctl *CouponController) Receive(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	var req struct {
		CouponCode string `json:"coupon_code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}

	item, err := ctl.couponSvc.ReceiveCoupon(userID, req.CouponCode)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}
	response.Success(c, item)
}

func (ctl *CouponController) MyCoupons(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	items, err := ctl.couponSvc.MyCoupons(userID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeServerError, err.Error())
		return
	}
	response.Success(c, items)
}

func (ctl *CouponController) Use(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	var req struct {
		OrderNo      string `json:"order_no" binding:"required"`
		UserCouponID uint   `json:"user_coupon_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}

	order, err := ctl.couponSvc.UseCoupon(userID, services.UseCouponInput{
		OrderNo:      req.OrderNo,
		UserCouponID: req.UserCouponID,
	})
	if err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}
	response.Success(c, order)
}
