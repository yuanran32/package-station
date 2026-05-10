package services

import (
	"agent_learning/models"
	"agent_learning/pkg/random"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type SendOrderService struct {
	db        *gorm.DB
	opLogSvc  *OperationLogService
	noticeSvc *NoticeService
}

type CreateSendOrderInput struct {
	SenderName      string
	SenderPhone     string
	SenderAddress   string
	ReceiverName    string
	ReceiverPhone   string
	ReceiverAddress string
	ItemInfo        string
	Weight          float64
}

type ProcessSendOrderInput struct {
	OrderNo     string
	Action      string
	CourierName string
}

type ProcessSendOrderResult struct {
	OrderNo       string            `json:"order_no"`
	Status        string            `json:"status"`
	InboundParcel InboundParcelInfo `json:"inbound_parcel"`
	ReceiverBound bool              `json:"receiver_bound"`
	NoticeSent    bool              `json:"notice_sent"`
}

type InboundParcelInfo struct {
	TrackingNo string `json:"tracking_no"`
	Location   string `json:"location"`
	PickupCode string `json:"pickup_code"`
}

func NewSendOrderService(db *gorm.DB, opLogSvc *OperationLogService, noticeSvc *NoticeService) *SendOrderService {
	return &SendOrderService{
		db:        db,
		opLogSvc:  opLogSvc,
		noticeSvc: noticeSvc,
	}
}

func (s *SendOrderService) CreateOrder(userID uint, input CreateSendOrderInput) (*models.SendOrder, error) {
	input.SenderName = strings.TrimSpace(input.SenderName)
	input.SenderPhone = strings.TrimSpace(input.SenderPhone)
	input.SenderAddress = strings.TrimSpace(input.SenderAddress)
	input.ReceiverName = strings.TrimSpace(input.ReceiverName)
	input.ReceiverPhone = strings.TrimSpace(input.ReceiverPhone)
	input.ReceiverAddress = strings.TrimSpace(input.ReceiverAddress)
	input.ItemInfo = strings.TrimSpace(input.ItemInfo)

	if input.SenderName == "" || input.SenderPhone == "" || input.SenderAddress == "" ||
		input.ReceiverName == "" || input.ReceiverPhone == "" || input.ReceiverAddress == "" ||
		input.ItemInfo == "" || input.Weight <= 0 {
		return nil, errors.New("invalid send order input")
	}

	estimatedFee := calculateShippingFee(input.Weight)
	order := models.SendOrder{
		OrderNo:         random.Serial("SO"),
		UserID:          userID,
		SenderName:      input.SenderName,
		SenderPhone:     input.SenderPhone,
		SenderAddress:   input.SenderAddress,
		ReceiverName:    input.ReceiverName,
		ReceiverPhone:   input.ReceiverPhone,
		ReceiverAddress: input.ReceiverAddress,
		ItemInfo:        input.ItemInfo,
		Weight:          input.Weight,
		EstimatedFee:    estimatedFee,
		CouponDeduct:    0,
		Status:          models.SendOrderStatusCreated,
		PayStatus:       models.PayStatusUnpaid,
		AssignedCourier: "",
	}

	if err := s.db.Create(&order).Error; err != nil {
		return nil, err
	}

	s.opLogSvc.Log(&userID, "send_order_create", "send_order", &order.ID, order.OrderNo)
	return &order, nil
}

func (s *SendOrderService) ListOrders(status string) ([]models.SendOrder, error) {
	var orders []models.SendOrder
	query := s.db.Model(&models.SendOrder{})
	status = strings.TrimSpace(status)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Order("id DESC").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *SendOrderService) Process(adminID uint, input ProcessSendOrderInput) (*ProcessSendOrderResult, error) {
	input.OrderNo = strings.TrimSpace(input.OrderNo)
	input.Action = strings.TrimSpace(input.Action)
	input.CourierName = strings.TrimSpace(input.CourierName)
	if input.OrderNo == "" || input.Action == "" {
		return nil, errors.New("order_no and action are required")
	}

	var (
		updated       models.SendOrder
		notifyUserID  *uint
		notifyText    string
		inboundParcel InboundParcelInfo
		receiverBound bool
	)
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		var order models.SendOrder
		if err := tx.Where("order_no = ?", input.OrderNo).First(&order).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("send order not found")
			}
			return err
		}

		switch input.Action {
		case "accept":
			order.Status = models.SendOrderStatusAccepted
		case "assign_pickup":
			if input.CourierName == "" {
				return errors.New("courier_name is required for assign_pickup")
			}
			order.Status = models.SendOrderStatusPickupAssign
			order.AssignedCourier = input.CourierName
		case "complete":
			order.Status = models.SendOrderStatusCompleted

			parcel, created, err := s.ensureParcelForCompletedOrder(tx, adminID, order)
			if err != nil {
				return err
			}
			if created {
				s.opLogSvc.Log(&adminID, "parcel_auto_inbound", "parcel", &parcel.ID, parcel.TrackingNo)
			}

			inboundParcel = InboundParcelInfo{
				TrackingNo: parcel.TrackingNo,
				Location:   parcel.Location,
				PickupCode: parcel.PickupCode,
			}

			// Strictly bind and notify by receiver phone.
			var receiver models.User
			if err := tx.Where("phone = ?", order.ReceiverPhone).First(&receiver).Error; err == nil {
				receiverBound = true
				if parcel.RecipientUserID == nil || *parcel.RecipientUserID != receiver.ID {
					parcel.RecipientUserID = &receiver.ID
					parcel.UpdatedBy = adminID
					if saveErr := tx.Save(&parcel).Error; saveErr != nil {
						return saveErr
					}
				}
				id := receiver.ID
				notifyUserID = &id
				notifyText = fmt.Sprintf(
					"Your parcel %s has arrived at station. Pickup code: %s, location: %s",
					parcel.TrackingNo,
					parcel.PickupCode,
					parcel.Location,
				)
			}
		default:
			return errors.New("unsupported action")
		}

		if err := tx.Save(&order).Error; err != nil {
			return err
		}
		updated = order
		return nil
	}); err != nil {
		return nil, err
	}

	noticeSent := false
	if notifyUserID != nil && s.noticeSvc != nil {
		if err := s.noticeSvc.SendSystemNotice(*notifyUserID, notifyText); err == nil {
			noticeSent = true
		}
	}
	s.opLogSvc.Log(&adminID, "send_order_process", "send_order", &updated.ID, fmt.Sprintf("%s:%s", updated.OrderNo, input.Action))
	return &ProcessSendOrderResult{
		OrderNo:       updated.OrderNo,
		Status:        updated.Status,
		InboundParcel: inboundParcel,
		ReceiverBound: receiverBound,
		NoticeSent:    noticeSent,
	}, nil
}

func (s *SendOrderService) ensureParcelForCompletedOrder(tx *gorm.DB, adminID uint, order models.SendOrder) (models.Parcel, bool, error) {
	const autoLocation = "AUTO-INBOUND-RACK"

	var existing models.Parcel
	err := tx.Where("tracking_no = ?", order.OrderNo).First(&existing).Error
	if err == nil {
		changed := false
		if strings.TrimSpace(existing.PickupCode) == "" {
			existing.PickupCode = random.NumericCode(6)
			changed = true
		}
		if existing.RecipientPhone != order.ReceiverPhone {
			existing.RecipientPhone = order.ReceiverPhone
			changed = true
		}
		if existing.RecipientUserID == nil {
			var user models.User
			if findErr := tx.Where("phone = ?", order.ReceiverPhone).First(&user).Error; findErr == nil {
				existing.RecipientUserID = &user.ID
				changed = true
			}
		}
		if changed {
			existing.UpdatedBy = adminID
			if saveErr := tx.Save(&existing).Error; saveErr != nil {
				return models.Parcel{}, false, saveErr
			}
		}
		return existing, false, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Parcel{}, false, err
	}

	now := time.Now()
	parcel := models.Parcel{
		TrackingNo:     order.OrderNo,
		RecipientPhone: order.ReceiverPhone,
		Location:       autoLocation,
		Status:         models.ParcelStatusInWarehouse,
		PickupCode:     random.NumericCode(6),
		InboundAt:      now,
		CreatedBy:      adminID,
		UpdatedBy:      adminID,
	}

	var user models.User
	if findErr := tx.Where("phone = ?", order.ReceiverPhone).First(&user).Error; findErr == nil {
		parcel.RecipientUserID = &user.ID
	}

	if err := tx.Create(&parcel).Error; err != nil {
		return models.Parcel{}, false, err
	}
	return parcel, true, nil
}

func calculateShippingFee(weight float64) float64 {
	if weight <= 1 {
		return 8.00
	}
	return 8.00 + (weight-1)*4.00
}
