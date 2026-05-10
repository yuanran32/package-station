package services

import (
	"agent_learning/models"
	"agent_learning/pkg/random"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

type RegisterUserInput struct {
	Username string
	Password string
	Name     string
	Phone    string
}

type UpdateUserProfileInput struct {
	Name  string
	Phone string
}

type PickupHistoryItem struct {
	ID             uint      `json:"id" gorm:"column:id"`
	ParcelID       uint      `json:"parcel_id" gorm:"column:parcel_id"`
	TrackingNo     string    `json:"tracking_no" gorm:"column:tracking_no"`
	PickupCode     string    `json:"pickup_code" gorm:"column:pickup_code"`
	Location       string    `json:"location" gorm:"column:location"`
	Status         string    `json:"status" gorm:"column:status"`
	PickupUserID   *uint     `json:"pickup_user_id" gorm:"column:pickup_user_id"`
	PickupUserName string    `json:"pickup_user_name" gorm:"column:pickup_user_name"`
	PickupTime     time.Time `json:"pickup_time" gorm:"column:pickup_time"`
	OperatorUserID uint      `json:"operator_user_id" gorm:"column:operator_user_id"`
	CreatedAt      time.Time `json:"created_at" gorm:"column:created_at"`
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// 用户注册的业务层代码
func (s *UserService) Register(input RegisterUserInput) (*models.User, error) {
	input.Username = strings.TrimSpace(input.Username)
	input.Phone = strings.TrimSpace(input.Phone)
	input.Name = strings.TrimSpace(input.Name)

	if input.Username == "" || input.Password == "" || input.Phone == "" {
		return nil, errors.New("username/password/phone are required")
	}

	var count int64
	if err := s.db.Model(&models.User{}).Where("username = ?", input.Username).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("username already exists")
	}
	count = 0
	if err := s.db.Model(&models.User{}).Where("phone = ?", input.Phone).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("phone already exists")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Username:     input.Username,
		Password:     string(hashed),
		Name:         input.Name,
		Phone:        input.Phone,
		Role:         models.RoleUser,
		IdentityCode: "UID-" + random.AlphaNumeric(12),
	}
	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// 登录界面的业务层
func (s *UserService) Login(account, password string) (*models.User, error) {

	account = strings.TrimSpace(account)
	if account == "" || strings.TrimSpace(password) == "" {
		return nil, errors.New("account and password are required")
	}

	var user models.User
	if err := s.db.Where("username = ? OR phone = ?", account, account).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	//先找到user的账号然后通过账号的密码和输入密码比对
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	return &user, nil
}

// 获得该User的个人信息
func (s *UserService) GetProfile(userID uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// 更改用户的个人信息
func (s *UserService) UpdateProfile(userID uint, input UpdateUserProfileInput) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	input.Name = strings.TrimSpace(input.Name)
	input.Phone = strings.TrimSpace(input.Phone)
	if input.Phone != "" && input.Phone != user.Phone {
		var count int64
		if err := s.db.Model(&models.User{}).Where("phone = ?", input.Phone).Count(&count).Error; err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, errors.New("phone already exists")
		}
		user.Phone = input.Phone
	}
	if input.Name != "" {
		user.Name = input.Name
	}
	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetPickupHistory(userID uint) ([]PickupHistoryItem, error) {
	var records []PickupHistoryItem
	if err := s.db.Table("pickup_records AS pr").
		Select(`
			pr.id AS id,
			pr.parcel_id AS parcel_id,
			pr.tracking_no AS tracking_no,
			COALESCE(p.pickup_code, '') AS pickup_code,
			COALESCE(p.location, '') AS location,
			CASE
				WHEN p.id IS NULL THEN 'lost'
				WHEN p.status = 'in_warehouse' THEN 'pending'
				WHEN p.status = 'picked_up' THEN 'picked_up'
				WHEN p.status = 'delivered' THEN 'delivered'
				ELSE 'lost'
			END AS status,
			pr.pickup_user_id AS pickup_user_id,
			COALESCE(pr.pickup_user_name, '') AS pickup_user_name,
			pr.pickup_time AS pickup_time,
			pr.operator_user_id AS operator_user_id,
			pr.created_at AS created_at
		`).
		Joins("LEFT JOIN parcels p ON p.id = pr.parcel_id").
		Where("pr.pickup_user_id = ? OR p.recipient_user_id = ?", userID, userID).
		Order("pr.pickup_time DESC").
		Scan(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

func (s *UserService) GetIdentityCode(userID uint) (string, error) {
	var user models.User
	if err := s.db.Select("identity_code").First(&user, userID).Error; err != nil {
		return "", err
	}
	return user.IdentityCode, nil
}
