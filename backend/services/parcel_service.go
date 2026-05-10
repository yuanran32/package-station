package services

import (
	"agent_learning/models"
	"agent_learning/pkg/fieldcheck"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ParcelService struct {
	db        *gorm.DB
	opLogSvc  *OperationLogService
	noticeSvc *NoticeService
}

type InboundParcelInput struct {
	TrackingNo     string
	Location       string
	RecipientPhone string
}

type OutboundParcelInput struct {
	TrackingNo string
	PickupCode string
	PickupName string
}

func NewParcelService(db *gorm.DB, opLogSvc *OperationLogService, noticeSvc *NoticeService) *ParcelService {
	return &ParcelService{
		db:        db,
		opLogSvc:  opLogSvc,
		noticeSvc: noticeSvc,
	}
}

func (s *ParcelService) Inbound(adminID uint, input InboundParcelInput) (*models.Parcel, error) {
	input.TrackingNo = strings.TrimSpace(input.TrackingNo)
	input.Location = strings.TrimSpace(input.Location)
	input.RecipientPhone = strings.TrimSpace(input.RecipientPhone)
	if input.TrackingNo == "" || input.Location == "" || input.RecipientPhone == "" {
		return nil, errors.New("tracking_no/location/recipient_phone are required")
	}
	if !fieldcheck.IsTrackingNo(input.TrackingNo) {
		return nil, errors.New("invalid tracking_no format")
	}
	if !fieldcheck.IsCNMobile(input.RecipientPhone) {
		return nil, errors.New("invalid recipient_phone format")
	}
	if len(input.Location) > 128 {
		return nil, errors.New("location is too long, max 128")
	}

	var created models.Parcel
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&models.Parcel{}).Where("tracking_no = ?", input.TrackingNo).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("tracking_no already exists")
		}

		parcel := models.Parcel{
			TrackingNo:     input.TrackingNo,
			RecipientPhone: input.RecipientPhone,
			Location:       input.Location,
			Status:         models.ParcelStatusInWarehouse,
			InboundAt:      time.Now(),
			CreatedBy:      adminID,
			UpdatedBy:      adminID,
		}
		var user models.User
		if err := tx.Where("phone = ?", input.RecipientPhone).First(&user).Error; err == nil {
			parcel.RecipientUserID = &user.ID
		}

		if err := tx.Create(&parcel).Error; err != nil {
			return err
		}
		created = parcel
		return nil
	})
	if err != nil {
		return nil, err
	}

	s.opLogSvc.Log(&adminID, "parcel_inbound", "parcel", &created.ID, created.TrackingNo)
	if created.RecipientUserID != nil && s.noticeSvc != nil {
		_ = s.noticeSvc.SendSystemNotice(*created.RecipientUserID, fmt.Sprintf("Your parcel %s is inbound, location: %s", created.TrackingNo, created.Location))
	}
	return &created, nil
}

func (s *ParcelService) Outbound(operatorID uint, input OutboundParcelInput) (*models.Parcel, error) {
	input.TrackingNo = strings.TrimSpace(input.TrackingNo)
	input.PickupCode = strings.TrimSpace(input.PickupCode)
	input.PickupName = strings.TrimSpace(input.PickupName)

	if input.TrackingNo == "" {
		return nil, errors.New("tracking_no is required")
	}
	if !fieldcheck.IsTrackingNo(input.TrackingNo) {
		return nil, errors.New("invalid tracking_no format")
	}
	if input.PickupCode != "" && !fieldcheck.IsPickupCode(input.PickupCode) {
		return nil, errors.New("pickup_code must be 6 digits")
	}
	if len(input.PickupName) > 100 {
		return nil, errors.New("pickup_name is too long, max 100")
	}

	var updated models.Parcel
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var parcel models.Parcel
		if err := tx.Where("tracking_no = ?", input.TrackingNo).First(&parcel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("parcel not found")
			}
			return err
		}
		if parcel.Status != models.ParcelStatusInWarehouse {
			return errors.New("parcel is not in warehouse")
		}
		if strings.TrimSpace(parcel.PickupCode) != "" {
			if input.PickupCode == "" || input.PickupCode != parcel.PickupCode {
				return errors.New("invalid pickup code")
			}
		}

		now := time.Now()
		parcel.Status = models.ParcelStatusPickedUp
		parcel.OutboundAt = &now
		parcel.UpdatedBy = operatorID
		if err := tx.Save(&parcel).Error; err != nil {
			return err
		}

		pickupName := input.PickupName
		if pickupName == "" {
			pickupName = "user pickup"
		}

		pickupUserID := parcel.RecipientUserID
		if pickupUserID == nil {
			pickupUserID = &operatorID
		}

		pickupRecord := models.PickupRecord{
			ParcelID:       parcel.ID,
			TrackingNo:     parcel.TrackingNo,
			PickupUserID:   pickupUserID,
			PickupUserName: pickupName,
			PickupTime:     now,
			OperatorUserID: operatorID,
		}
		if err := tx.Create(&pickupRecord).Error; err != nil {
			return err
		}
		updated = parcel
		return nil
	})
	if err != nil {
		return nil, err
	}

	s.opLogSvc.Log(&operatorID, "parcel_outbound", "parcel", &updated.ID, updated.TrackingNo)
	if updated.RecipientUserID != nil && s.noticeSvc != nil {
		_ = s.noticeSvc.SendSystemNotice(*updated.RecipientUserID, fmt.Sprintf("Parcel %s has been picked up", updated.TrackingNo))
	}
	return &updated, nil
}

func (s *ParcelService) GetStatus(trackingNo string) (*models.Parcel, error) {
	trackingNo = strings.TrimSpace(trackingNo)
	if trackingNo == "" {
		return nil, errors.New("tracking_no is required")
	}
	var parcel models.Parcel
	if err := s.db.Where("tracking_no = ?", trackingNo).First(&parcel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("parcel not found")
		}
		return nil, err
	}
	return &parcel, nil
}

func (s *ParcelService) ListAll() ([]models.Parcel, error) {
	var parcels []models.Parcel
	if err := s.db.Order("id DESC").Find(&parcels).Error; err != nil {
		return nil, err
	}
	return parcels, nil
}
