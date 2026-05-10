package models

import "time"

type PickupRecord struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ParcelID       uint      `gorm:"index;not null" json:"parcel_id"`
	TrackingNo     string    `gorm:"size:64;index;not null" json:"tracking_no"`
	PickupUserID   *uint     `gorm:"index" json:"pickup_user_id"`
	PickupUserName string    `gorm:"size:100;not null" json:"pickup_user_name"`
	PickupTime     time.Time `gorm:"not null" json:"pickup_time"`
	OperatorUserID uint      `gorm:"not null" json:"operator_user_id"`
	CreatedAt      time.Time `json:"created_at"`
}
