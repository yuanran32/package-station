package models

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

const (
	ParcelStatusInWarehouse = "in_warehouse"
	ParcelStatusPickedUp    = "picked_up"
	ParcelStatusDelivered   = "delivered"
)

const (
	SendOrderStatusCreated      = "created"
	SendOrderStatusAccepted     = "accepted"
	SendOrderStatusPickupAssign = "pickup_assigned"
	SendOrderStatusCompleted    = "completed"
)

const (
	PayStatusUnpaid = "unpaid"
	PayStatusPaid   = "paid"
)

const (
	CouponStatusActive   = "active"
	CouponStatusInactive = "inactive"
)

const (
	UserCouponStatusUnused = "unused"
	UserCouponStatusUsed   = "used"
)
