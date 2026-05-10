package services

import (
	"agent_learning/models"
	"time"

	"gorm.io/gorm"
)

type NoticeService struct {
	db  *gorm.DB
	hub *NotifyHub
}

func NewNoticeService(db *gorm.DB, hub *NotifyHub) *NoticeService {
	return &NoticeService{db: db, hub: hub}
}

func (s *NoticeService) SendNotice(userID *uint, noticeType, content string) (*models.Notice, error) {
	notice := models.Notice{
		UserID:  userID,
		Type:    noticeType,
		Content: content,
		SentAt:  time.Now(),
	}
	if err := s.db.Create(&notice).Error; err != nil {
		return nil, err
	}

	if userID != nil && s.hub != nil {
		s.hub.SendToUser(*userID, map[string]interface{}{
			"type":       notice.Type,
			"content":    notice.Content,
			"sent_at":    notice.SentAt,
			"notice_id":  notice.ID,
			"event_type": "notify",
		})
	}
	return &notice, nil
}

func (s *NoticeService) SendSystemNotice(userID uint, content string) error {
	_, err := s.SendNotice(&userID, "system", content)
	return err
}
