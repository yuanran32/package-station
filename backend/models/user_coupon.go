package models

import "time"

type UserCoupon struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	UserID        uint       `gorm:"index;not null" json:"user_id"`
	CouponID      uint       `gorm:"index;not null" json:"coupon_id"`
	Status        string     `gorm:"size:20;index;not null;default:unused" json:"status"`
	ReceivedAt    time.Time  `gorm:"not null" json:"received_at"`
	UsedAt        *time.Time `json:"used_at"`
	UsedInOrderID *uint      `gorm:"index" json:"used_in_order_id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
