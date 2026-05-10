package models

import "time"

type DeliveryRecord struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ParcelID       *uint     `gorm:"index" json:"parcel_id"`
	TrackingNo     string    `gorm:"size:64;index;not null" json:"tracking_no"`
	CourierName    string    `gorm:"size:100;not null" json:"courier_name"`
	DeliveryStatus string    `gorm:"size:32;not null" json:"delivery_status"`
	Remark         string    `gorm:"size:255" json:"remark"`
	DeliveredAt    time.Time `gorm:"not null" json:"delivered_at"`
	OperatorUserID uint      `gorm:"not null" json:"operator_user_id"`
	CreatedAt      time.Time `json:"created_at"`
}
