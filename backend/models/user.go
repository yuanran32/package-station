package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Password     string    `gorm:"size:255;not null" json:"-"`
	Name         string    `gorm:"size:100" json:"name"`
	Phone        string    `gorm:"size:20;uniqueIndex;not null" json:"phone"`
	Role         string    `gorm:"size:20;not null;default:user" json:"role"`
	IdentityCode string    `gorm:"size:64;uniqueIndex;not null" json:"identity_code"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
