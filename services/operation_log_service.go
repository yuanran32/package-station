package services

import (
	"agent_learning/models"
	"log"

	"gorm.io/gorm"
)

type OperationLogService struct {
	db *gorm.DB
}

func NewOperationLogService(db *gorm.DB) *OperationLogService {
	return &OperationLogService{db: db}
}

func (s *OperationLogService) Log(userID *uint, action, targetType string, targetID *uint, detail string) {
	entry := models.OperationLog{
		UserID:     userID,
		Action:     action,
		TargetType: targetType,
		TargetID:   targetID,
		Detail:     detail,
	}
	if err := s.db.Create(&entry).Error; err != nil {
		log.Printf("[operation-log] action=%s err=%v", action, err)
	}
}
