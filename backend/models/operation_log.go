package models

import "time"

type OperationLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     *uint     `gorm:"index" json:"user_id"`
	Action     string    `gorm:"size:64;index;not null" json:"action"`
	TargetType string    `gorm:"size:64;index;not null" json:"target_type"`
	TargetID   *uint     `gorm:"index" json:"target_id"`
	Detail     string    `gorm:"size:500" json:"detail"`
	CreatedAt  time.Time `json:"created_at"`
}
