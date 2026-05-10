package models

import "time"

type Notice struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    *uint     `gorm:"index" json:"user_id"`
	Type      string    `gorm:"size:32;index;not null" json:"type"`
	Content   string    `gorm:"size:500;not null" json:"content"`
	SentAt    time.Time `gorm:"not null" json:"sent_at"`
	CreatedAt time.Time `json:"created_at"`
}
