package models

import "time"

type Coupon struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Code         string    `gorm:"size:64;uniqueIndex;not null" json:"code"`
	Name         string    `gorm:"size:100;not null" json:"name"`
	ActivityRule string    `gorm:"size:255;not null" json:"activity_rule"`
	Amount       float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	Threshold    float64   `gorm:"type:decimal(10,2);not null;default:0" json:"threshold"`
	Total        int       `gorm:"not null" json:"total"`
	Remaining    int       `gorm:"not null" json:"remaining"`
	Status       string    `gorm:"size:20;not null;default:active" json:"status"`
	StartAt      time.Time `gorm:"not null" json:"start_at"`
	EndAt        time.Time `gorm:"not null" json:"end_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
