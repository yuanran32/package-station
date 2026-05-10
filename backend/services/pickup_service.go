package services

import (
	"agent_learning/models"
	"agent_learning/pkg/fieldcheck"
	"agent_learning/pkg/random"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type PickupService struct {
	db        *gorm.DB
	opLogSvc  *OperationLogService
	noticeSvc *NoticeService
}

type GeneratePickupCodeInput struct {
	TrackingNo     string
	RecipientPhone string
}

type RecordPickupInput struct {
	TrackingNo     string
	PickupUserID   *uint
	PickupUserName string
}

type RecordDeliveryInput struct {
	TrackingNo     string
	CourierName    string
	DeliveryStatus string
	Remark         string
	DeliveredAt    *time.Time
}

func NewPickupService(db *gorm.DB, opLogSvc *OperationLogService, noticeSvc *NoticeService) *PickupService {
	return &PickupService{
		db:        db,
		opLogSvc:  opLogSvc,
		noticeSvc: noticeSvc,
	}
}

func (s *PickupService) GeneratePickupCode(operatorID uint, input GeneratePickupCodeInput) (*models.Parcel, error) {
	input.TrackingNo = strings.TrimSpace(input.TrackingNo)
	input.RecipientPhone = strings.TrimSpace(input.RecipientPhone)
	if input.TrackingNo == "" && input.RecipientPhone == "" {
		return nil, errors.New("tracking_no or recipient_phone is required")
	}
	if input.TrackingNo != "" && !fieldcheck.IsTrackingNo(input.TrackingNo) {
		return nil, errors.New("invalid tracking_no format")
	}
	if input.RecipientPhone != "" && !fieldcheck.IsCNMobile(input.RecipientPhone) {
		return nil, errors.New("invalid recipient_phone format")
	}

	var updated models.Parcel
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var parcel models.Parcel
		if input.TrackingNo != "" {
			if err := tx.Where("status = ? AND tracking_no = ?", models.ParcelStatusInWarehouse, input.TrackingNo).First(&parcel).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.New("parcel not found")
				}
				return err
			}
			if input.RecipientPhone != "" && parcel.RecipientPhone != input.RecipientPhone {
				return errors.New("recipient_phone does not match tracking_no")
			}
		} else {
			if err := tx.Where("status = ? AND recipient_phone = ?", models.ParcelStatusInWarehouse, input.RecipientPhone).
				Order("id DESC").
				First(&parcel).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.New("parcel not found")
				}
				return err
			}
		}

		parcel.PickupCode = random.NumericCode(6)
		parcel.UpdatedBy = operatorID
		if err := tx.Save(&parcel).Error; err != nil {
			return err
		}
		updated = parcel
		return nil
	})
	if err != nil {
		return nil, err
	}

	s.opLogSvc.Log(&operatorID, "pickup_code_generate", "parcel", &updated.ID, updated.TrackingNo)
	if updated.RecipientUserID != nil && s.noticeSvc != nil {
		_ = s.noticeSvc.SendSystemNotice(*updated.RecipientUserID, fmt.Sprintf("Parcel %s pickup code: %s", updated.TrackingNo, updated.PickupCode))
	}
	return &updated, nil
}

func (s *PickupService) RecordPickup(operatorID uint, input RecordPickupInput) (*models.PickupRecord, error) {
	input.TrackingNo = strings.TrimSpace(input.TrackingNo)
	input.PickupUserName = strings.TrimSpace(input.PickupUserName)
	if input.TrackingNo == "" {
		return nil, errors.New("tracking_no is required")
	}
	if !fieldcheck.IsTrackingNo(input.TrackingNo) {
		return nil, errors.New("invalid tracking_no format")
	}
	if len(input.PickupUserName) > 100 {
		return nil, errors.New("pickup_user_name is too long, max 100")
	}

	var created models.PickupRecord
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var parcel models.Parcel
		if err := tx.Where("tracking_no = ?", input.TrackingNo).First(&parcel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("parcel not found")
			}
			return err
		}

		now := time.Now()
		if parcel.Status != models.ParcelStatusPickedUp {
			parcel.Status = models.ParcelStatusPickedUp
			parcel.OutboundAt = &now
			parcel.UpdatedBy = operatorID
			if err := tx.Save(&parcel).Error; err != nil {
				return err
			}
		}

		name := input.PickupUserName
		if name == "" {
			name = "manual record"
		}

		pickupUserID := input.PickupUserID
		if pickupUserID == nil {
			pickupUserID = parcel.RecipientUserID
		}

		created = models.PickupRecord{
			ParcelID:       parcel.ID,
			TrackingNo:     parcel.TrackingNo,
			PickupUserID:   pickupUserID,
			PickupUserName: name,
			PickupTime:     now,
			OperatorUserID: operatorID,
		}
		if err := tx.Create(&created).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	s.opLogSvc.Log(&operatorID, "pickup_record_create", "pickup_record", &created.ID, created.TrackingNo)
	return &created, nil
}

func (s *PickupService) RecordDelivery(operatorID uint, input RecordDeliveryInput) (*models.DeliveryRecord, error) {
	input.TrackingNo = strings.TrimSpace(input.TrackingNo)
	input.CourierName = strings.TrimSpace(input.CourierName)
	input.DeliveryStatus = strings.TrimSpace(input.DeliveryStatus)
	input.Remark = strings.TrimSpace(input.Remark)

	if input.TrackingNo == "" || input.CourierName == "" || input.DeliveryStatus == "" {
		return nil, errors.New("tracking_no/courier_name/delivery_status are required")
	}
	if !fieldcheck.IsTrackingNo(input.TrackingNo) {
		return nil, errors.New("invalid tracking_no format")
	}
	if len(input.CourierName) > 100 {
		return nil, errors.New("courier_name is too long, max 100")
	}
	if len(input.DeliveryStatus) > 32 {
		return nil, errors.New("delivery_status is too long, max 32")
	}

	deliveredAt := time.Now()
	if input.DeliveredAt != nil {
		deliveredAt = *input.DeliveredAt
	}

	var created models.DeliveryRecord
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var parcel models.Parcel
		parcelExists := true
		if err := tx.Where("tracking_no = ?", input.TrackingNo).First(&parcel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				parcelExists = false
			} else {
				return err
			}
		}

		created = models.DeliveryRecord{
			TrackingNo:     input.TrackingNo,
			CourierName:    input.CourierName,
			DeliveryStatus: input.DeliveryStatus,
			Remark:         input.Remark,
			DeliveredAt:    deliveredAt,
			OperatorUserID: operatorID,
		}
		if parcelExists {
			created.ParcelID = &parcel.ID
			if strings.EqualFold(input.DeliveryStatus, models.ParcelStatusDelivered) {
				parcel.Status = models.ParcelStatusDelivered
				parcel.UpdatedBy = operatorID
				if err := tx.Save(&parcel).Error; err != nil {
					return err
				}
			}
		}

		if err := tx.Create(&created).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	s.opLogSvc.Log(&operatorID, "delivery_record_create", "delivery_record", &created.ID, created.TrackingNo)
	return &created, nil
}
