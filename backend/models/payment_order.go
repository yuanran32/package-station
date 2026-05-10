package models

import "time"

type PaymentOrder struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	PayNo           string     `gorm:"size:64;uniqueIndex;not null" json:"pay_no"`
	UserID          uint       `gorm:"index;not null" json:"user_id"`
	RelatedType     string     `gorm:"size:32;index;not null" json:"related_type"`
	RelatedID       uint       `gorm:"index;not null" json:"related_id"`
	BizDesc         string     `gorm:"size:255" json:"biz_desc"`
	Amount          float64    `gorm:"type:decimal(10,2);not null" json:"amount"`
	Status          string     `gorm:"size:20;index;not null;default:pending" json:"status"`
	PayMethod       string     `gorm:"size:32" json:"pay_method"`
	PaidAt          *time.Time `json:"paid_at"`
	CallbackPayload string     `gorm:"type:text" json:"callback_payload"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
