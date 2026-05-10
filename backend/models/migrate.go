package models

import "gorm.io/gorm"

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Parcel{},
		&PickupRecord{},
		&DeliveryRecord{},
		&SendOrder{},
		&Coupon{},
		&UserCoupon{},
		&PaymentOrder{},
		&Notice{},
		&OperationLog{},
	)
}
