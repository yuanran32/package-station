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

type PaymentService struct {
	db       *gorm.DB
	opLogSvc *OperationLogService
}

type CreatePaymentInput struct {
	RelatedType string
	OrderNo     string
	RelatedID   uint
	Amount      float64
	BizDesc     string
}

type PaymentCallbackInput struct {
	PayNo   string
	Status  string
	Method  string
	Payload string
}

func NewPaymentService(db *gorm.DB, opLogSvc *OperationLogService) *PaymentService {
	return &PaymentService{
		db:       db,
		opLogSvc: opLogSvc,
	}
}

func (s *PaymentService) CreatePayment(userID uint, input CreatePaymentInput) (*models.PaymentOrder, error) {
	input.RelatedType = strings.TrimSpace(input.RelatedType)
	input.OrderNo = strings.TrimSpace(input.OrderNo)
	input.BizDesc = strings.TrimSpace(input.BizDesc)

	var created models.PaymentOrder
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var relatedID uint
		amount := input.Amount

		switch input.RelatedType {
		case "send_order":
			if input.OrderNo == "" {
				return errors.New("order_no is required when related_type=send_order")
			}
			var order models.SendOrder
			if err := tx.Where("order_no = ? AND user_id = ?", input.OrderNo, userID).First(&order).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.New("send order not found")
				}
				return err
			}
			if order.PayStatus == models.PayStatusPaid {
				return errors.New("send order already paid")
			}
			relatedID = order.ID
			amount = order.EstimatedFee
			if amount <= 0 {
				return errors.New("order amount is invalid")
			}
		case "storage_fee":
			if input.RelatedID == 0 || amount <= 0 {
				return errors.New("related_id and amount are required for storage_fee")
			}
			relatedID = input.RelatedID
		default:
			return errors.New("unsupported related_type")
		}

		created = models.PaymentOrder{
			PayNo:       random.Serial("PAY"),
			UserID:      userID,
			RelatedType: input.RelatedType,
			RelatedID:   relatedID,
			BizDesc:     input.BizDesc,
			Amount:      amount,
			Status:      "pending",
		}
		if err := tx.Create(&created).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	s.opLogSvc.Log(&userID, "payment_create", "payment_order", &created.ID, fmt.Sprintf("%s:%.2f", created.PayNo, created.Amount))
	return &created, nil
}

func (s *PaymentService) Callback(input PaymentCallbackInput) (*models.PaymentOrder, error) {
	input.PayNo = strings.TrimSpace(input.PayNo)
	input.Status = strings.TrimSpace(input.Status)
	input.Method = strings.TrimSpace(input.Method)

	if input.PayNo == "" || input.Status == "" {
		return nil, errors.New("pay_no and status are required")
	}

	var updated models.PaymentOrder
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var pay models.PaymentOrder
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("pay_no = ?", input.PayNo).First(&pay).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("payment order not found")
			}
			return err
		}

		if pay.Status == "paid" {
			updated = pay
			return nil
		}

		now := time.Now()
		if strings.EqualFold(input.Status, "success") || strings.EqualFold(input.Status, "paid") {
			pay.Status = "paid"
			pay.PayMethod = input.Method
			pay.PaidAt = &now
		} else {
			pay.Status = "failed"
		}
		pay.CallbackPayload = input.Payload
		if err := tx.Save(&pay).Error; err != nil {
			return err
		}

		if pay.Status == "paid" && pay.RelatedType == "send_order" {
			if err := tx.Model(&models.SendOrder{}).
				Where("id = ?", pay.RelatedID).
				Update("pay_status", models.PayStatusPaid).Error; err != nil {
				return err
			}
		}
		updated = pay
		return nil
	})
	if err != nil {
		return nil, err
	}

	s.opLogSvc.Log(&updated.UserID, "payment_callback", "payment_order", &updated.ID, fmt.Sprintf("%s:%s", updated.PayNo, updated.Status))
	return &updated, nil
}

func (s *PaymentService) Bills(userID uint) ([]models.PaymentOrder, error) {
	var bills []models.PaymentOrder
	if err := s.db.Where("user_id = ?", userID).Order("id DESC").Find(&bills).Error; err != nil {
		return nil, err
	}
	return bills, nil
}
