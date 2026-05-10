package app

import (
	"agent_learning/internal/config"
	"agent_learning/internal/database"
	"agent_learning/models"
	"agent_learning/pkg/random"
	"agent_learning/routes"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Run() {
	cfg := config.Load()

	db, err := database.InitMySQL(cfg)
	if err != nil {
		log.Fatalf("init mysql failed: %v", err)
	}

	if err := models.AutoMigrate(db); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}

	if err := seedAdmin(db); err != nil {
		log.Fatalf("seed admin failed: %v", err)
	}
	if err := seedCoupon(db); err != nil {
		log.Fatalf("seed coupon failed: %v", err)
	}

	r := routes.SetupRouter(cfg, db)
	addr := ":" + cfg.Port
	log.Printf("parcel station server listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("run server failed: %v", err)
	}
}

func seedAdmin(db *gorm.DB) error {
	var existing models.User
	err := db.Where("username = ?", "admin").First(&existing).Error
	if err == nil {
		return nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := models.User{
		Username:     "admin",
		Password:     string(hashed),
		Name:         "System Admin",
		Phone:        "18800000000",
		Role:         models.RoleAdmin,
		IdentityCode: "ADMIN-" + random.AlphaNumeric(10),
	}
	if err := db.Create(&admin).Error; err != nil {
		return err
	}
	log.Printf("default admin seeded: username=admin password=admin123")
	return nil
}

func seedCoupon(db *gorm.DB) error {
	var existing models.Coupon
	err := db.Where("code = ?", "NEWUSER10").First(&existing).Error
	if err == nil {
		return nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	now := time.Now()
	coupon := models.Coupon{
		Code:         "NEWUSER10",
		Name:         "New User 10 Off",
		ActivityRule: "first order over 20",
		Amount:       10,
		Threshold:    20,
		Total:        10000,
		Remaining:    10000,
		Status:       models.CouponStatusActive,
		StartAt:      now,
		EndAt:        now.AddDate(1, 0, 0),
	}
	return db.Create(&coupon).Error
}
