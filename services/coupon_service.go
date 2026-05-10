package services

import (
	"agent_learning/models"
	"agent_learning/pkg/random"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CouponService struct {
	db       *gorm.DB
	opLogSvc *OperationLogService
}

type UseCouponInput struct {
	OrderNo      string
	UserCouponID uint
}

type CreateCouponInput struct {
	Name         string
	Amount       float64
	Code         string
	ActivityRule string
	Threshold    float64
	Total        int
	ValidDays    int
}

type AdminCouponListInput struct {
	Status   string
	Keyword  string
	Page     int
	PageSize int
}

type AdminCouponListResult struct {
	List     []models.Coupon `json:"list"`
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
}

type UserCouponView struct {
	UserCouponID uint       `json:"user_coupon_id" gorm:"column:user_coupon_id"`
	CouponID     uint       `json:"coupon_id" gorm:"column:coupon_id"`
	Code         string     `json:"code" gorm:"column:code"`
	Name         string     `json:"name" gorm:"column:name"`
	Amount       float64    `json:"amount" gorm:"column:amount"`
	Threshold    float64    `json:"threshold" gorm:"column:threshold"`
	Status       string     `json:"status" gorm:"column:status"`
	ReceivedAt   time.Time  `json:"received_at" gorm:"column:received_at"`
	UsedAt       *time.Time `json:"used_at" gorm:"column:used_at"`
}

func NewCouponService(db *gorm.DB, opLogSvc *OperationLogService) *CouponService {
	return &CouponService{
		db:       db,
		opLogSvc: opLogSvc,
	}
}

func (s *CouponService) CreateCoupon(adminID uint, input CreateCouponInput) (*models.Coupon, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Code = strings.ToUpper(strings.TrimSpace(input.Code))
	input.ActivityRule = strings.TrimSpace(input.ActivityRule)

	if input.Name == "" {
		return nil, errors.New("name is required")
	}
	if input.Amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}
	if input.Threshold < 0 {
		return nil, errors.New("threshold cannot be negative")
	}
	if input.Total <= 0 {
		input.Total = 1000
	}
	if input.ValidDays <= 0 {
		input.ValidDays = 30
	}
	if input.ActivityRule == "" {
		input.ActivityRule = "admin manual create"
	}
	if input.Code == "" {
		input.Code = "CPN" + random.AlphaNumeric(8)
	}

	var count int64
	if err := s.db.Model(&models.Coupon{}).Where("code = ?", input.Code).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("coupon code already exists")
	}

	now := time.Now()
	coupon := models.Coupon{
		Code:         input.Code,
		Name:         input.Name,
		ActivityRule: input.ActivityRule,
		Amount:       input.Amount,
		Threshold:    input.Threshold,
		Total:        input.Total,
		Remaining:    input.Total,
		Status:       models.CouponStatusActive,
		StartAt:      now,
		EndAt:        now.AddDate(0, 0, input.ValidDays),
	}
	if err := s.db.Create(&coupon).Error; err != nil {
		return nil, err
	}

	s.opLogSvc.Log(&adminID, "coupon_create", "coupon", &coupon.ID, coupon.Code)
	return &coupon, nil
}

func (s *CouponService) AdminListCoupons(input AdminCouponListInput) (*AdminCouponListResult, error) {
	input.Status = strings.TrimSpace(input.Status)
	input.Keyword = strings.TrimSpace(input.Keyword)
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.PageSize <= 0 {
		input.PageSize = 20
	}
	if input.PageSize > 200 {
		input.PageSize = 200
	}

	query := s.db.Model(&models.Coupon{})
	now := time.Now()
	switch input.Status {
	case "":
	case models.CouponStatusActive:
		query = query.Where("status = ? AND start_at <= ? AND end_at > ?", models.CouponStatusActive, now, now)
	case models.CouponStatusInactive:
		query = query.Where("status = ?", models.CouponStatusInactive)
	case "expired":
		query = query.Where("end_at <= ?", now)
	default:
		return nil, errors.New("unsupported status, allowed: active/inactive/expired")
	}
	if input.Keyword != "" {
		like := "%" + input.Keyword + "%"
		query = query.Where("code LIKE ? OR name LIKE ?", like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var list []models.Coupon
	offset := (input.Page - 1) * input.PageSize
	if err := query.Order("id DESC").Offset(offset).Limit(input.PageSize).Find(&list).Error; err != nil {
		return nil, err
	}

	return &AdminCouponListResult{
		List:     list,
		Total:    total,
		Page:     input.Page,
		PageSize: input.PageSize,
	}, nil
}

func (s *CouponService) ReceiveCoupon(userID uint, couponCode string) (*models.UserCoupon, error) {
	couponCode = strings.TrimSpace(couponCode)
	if couponCode == "" {
		return nil, errors.New("coupon_code is required")
	}

	var received models.UserCoupon
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var coupon models.Coupon
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("code = ?", couponCode).First(&coupon).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("coupon not found")
			}
			return err
		}

		now := time.Now()
		if coupon.Status != models.CouponStatusActive {
			return errors.New("coupon is inactive")
		}
		if now.Before(coupon.StartAt) || now.After(coupon.EndAt) {
			return errors.New("coupon not in valid period")
		}
		if coupon.Remaining <= 0 {
			return errors.New("coupon exhausted")
		}

		var count int64
		if err := tx.Model(&models.UserCoupon{}).Where("user_id = ? AND coupon_id = ?", userID, coupon.ID).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("coupon already received")
		}

		received = models.UserCoupon{
			UserID:     userID,
			CouponID:   coupon.ID,
			Status:     models.UserCouponStatusUnused,
			ReceivedAt: now,
		}
		if err := tx.Create(&received).Error; err != nil {
			return err
		}

		coupon.Remaining--
		if err := tx.Save(&coupon).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	s.opLogSvc.Log(&userID, "coupon_receive", "user_coupon", &received.ID, couponCode)
	return &received, nil
}

func (s *CouponService) MyCoupons(userID uint) ([]UserCouponView, error) {
	var list []UserCouponView
	err := s.db.Table("user_coupons AS uc").
		Select(`
			uc.id AS user_coupon_id,
			uc.coupon_id AS coupon_id,
			c.code AS code,
			c.name AS name,
			c.amount AS amount,
			c.threshold AS threshold,
			uc.status AS status,
			uc.received_at AS received_at,
			uc.used_at AS used_at
		`).
		Joins("JOIN coupons c ON c.id = uc.coupon_id").
		Where("uc.user_id = ?", userID).
		Order("uc.id DESC").
		Scan(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s *CouponService) UseCoupon(userID uint, input UseCouponInput) (*models.SendOrder, error) {
	input.OrderNo = strings.TrimSpace(input.OrderNo)
	if input.OrderNo == "" || input.UserCouponID == 0 {
		return nil, errors.New("order_no and user_coupon_id are required")
	}

	var updatedOrder models.SendOrder
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var userCoupon models.UserCoupon
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&userCoupon, input.UserCouponID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("user_coupon not found")
			}
			return err
		}
		if userCoupon.UserID != userID {
			return errors.New("coupon does not belong to user")
		}
		if userCoupon.Status != models.UserCouponStatusUnused {
			return errors.New("coupon already used")
		}

		var coupon models.Coupon
		if err := tx.First(&coupon, userCoupon.CouponID).Error; err != nil {
			return err
		}

		var order models.SendOrder
		if err := tx.Where("order_no = ? AND user_id = ?", input.OrderNo, userID).First(&order).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("send order not found")
			}
			return err
		}
		if order.CouponID != nil {
			return errors.New("order already used another coupon")
		}
		if order.EstimatedFee < coupon.Threshold {
			return fmt.Errorf("order amount must be >= %.2f", coupon.Threshold)
		}

		deduct := coupon.Amount
		if deduct > order.EstimatedFee {
			deduct = order.EstimatedFee
		}
		order.CouponDeduct = deduct
		order.EstimatedFee -= deduct
		order.CouponID = &coupon.ID
		if err := tx.Save(&order).Error; err != nil {
			return err
		}

		now := time.Now()
		userCoupon.Status = models.UserCouponStatusUsed
		userCoupon.UsedAt = &now
		userCoupon.UsedInOrderID = &order.ID
		if err := tx.Save(&userCoupon).Error; err != nil {
			return err
		}

		updatedOrder = order
		return nil
	})
	if err != nil {
		return nil, err
	}

	s.opLogSvc.Log(&userID, "coupon_use", "send_order", &updatedOrder.ID, updatedOrder.OrderNo)
	return &updatedOrder, nil
}
