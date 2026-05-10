package models

import "time"

type SendOrder struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	OrderNo         string    `gorm:"size:64;uniqueIndex;not null" json:"order_no"`
	UserID          uint      `gorm:"index;not null" json:"user_id"`
	SenderName      string    `gorm:"size:100;not null" json:"sender_name"`
	SenderPhone     string    `gorm:"size:20;not null" json:"sender_phone"`
	SenderAddress   string    `gorm:"size:255;not null" json:"sender_address"`
	ReceiverName    string    `gorm:"size:100;not null" json:"receiver_name"`
	ReceiverPhone   string    `gorm:"size:20;not null" json:"receiver_phone"`
	ReceiverAddress string    `gorm:"size:255;not null" json:"receiver_address"`
	ItemInfo        string    `gorm:"size:255;not null" json:"item_info"`
	Weight          float64   `gorm:"type:decimal(10,2);not null" json:"weight"`
	EstimatedFee    float64   `gorm:"type:decimal(10,2);not null" json:"estimated_fee"`
	CouponDeduct    float64   `gorm:"type:decimal(10,2);not null;default:0" json:"coupon_deduct"`
	CouponID        *uint     `gorm:"index" json:"coupon_id"`
	Status          string    `gorm:"size:32;index;not null;default:created" json:"status"`
	AssignedCourier string    `gorm:"size:100" json:"assigned_courier"`
	PayStatus       string    `gorm:"size:20;index;not null;default:unpaid" json:"pay_status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
