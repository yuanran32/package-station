package models

import "time"

type Parcel struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	TrackingNo      string     `gorm:"size:64;uniqueIndex;not null" json:"tracking_no"`
	RecipientPhone  string     `gorm:"size:20;index;not null" json:"recipient_phone"`
	RecipientUserID *uint      `gorm:"index" json:"recipient_user_id"`
	Location        string     `gorm:"size:128;not null" json:"location"`
	Status          string     `gorm:"size:32;index;not null;default:in_warehouse" json:"status"`
	PickupCode      string     `gorm:"size:20;index" json:"pickup_code"`
	InboundAt       time.Time  `json:"inbound_at"`
	OutboundAt      *time.Time `json:"outbound_at"`
	CreatedBy       uint       `json:"created_by"`
	UpdatedBy       uint       `json:"updated_by"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
