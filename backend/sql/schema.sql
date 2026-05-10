CREATE DATABASE IF NOT EXISTS parcel_station
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_unicode_ci;

USE parcel_station;

CREATE TABLE IF NOT EXISTS users (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  username VARCHAR(50) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  name VARCHAR(100) NULL,
  phone VARCHAR(20) NOT NULL UNIQUE,
  role VARCHAR(20) NOT NULL DEFAULT 'user',
  identity_code VARCHAR(64) NOT NULL UNIQUE,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS parcels (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  tracking_no VARCHAR(64) NOT NULL UNIQUE,
  recipient_phone VARCHAR(20) NOT NULL,
  recipient_user_id BIGINT UNSIGNED NULL,
  location VARCHAR(128) NOT NULL,
  status VARCHAR(32) NOT NULL DEFAULT 'in_warehouse',
  pickup_code VARCHAR(20) NULL,
  inbound_at DATETIME(3) NOT NULL,
  outbound_at DATETIME(3) NULL,
  created_by BIGINT UNSIGNED NOT NULL,
  updated_by BIGINT UNSIGNED NOT NULL,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  INDEX idx_parcels_recipient_phone (recipient_phone),
  INDEX idx_parcels_recipient_user_id (recipient_user_id),
  INDEX idx_parcels_status (status),
  INDEX idx_parcels_pickup_code (pickup_code),
  CONSTRAINT fk_parcels_recipient_user FOREIGN KEY (recipient_user_id) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS pickup_records (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  parcel_id BIGINT UNSIGNED NOT NULL,
  tracking_no VARCHAR(64) NOT NULL,
  pickup_user_id BIGINT UNSIGNED NULL,
  pickup_user_name VARCHAR(100) NOT NULL,
  pickup_time DATETIME(3) NOT NULL,
  operator_user_id BIGINT UNSIGNED NOT NULL,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  INDEX idx_pickup_records_parcel_id (parcel_id),
  INDEX idx_pickup_records_tracking_no (tracking_no),
  INDEX idx_pickup_records_pickup_user_id (pickup_user_id),
  CONSTRAINT fk_pickup_records_parcel FOREIGN KEY (parcel_id) REFERENCES parcels(id) ON DELETE CASCADE,
  CONSTRAINT fk_pickup_records_pickup_user FOREIGN KEY (pickup_user_id) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS delivery_records (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  parcel_id BIGINT UNSIGNED NULL,
  tracking_no VARCHAR(64) NOT NULL,
  courier_name VARCHAR(100) NOT NULL,
  delivery_status VARCHAR(32) NOT NULL,
  remark VARCHAR(255) NULL,
  delivered_at DATETIME(3) NOT NULL,
  operator_user_id BIGINT UNSIGNED NOT NULL,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  INDEX idx_delivery_records_parcel_id (parcel_id),
  INDEX idx_delivery_records_tracking_no (tracking_no),
  CONSTRAINT fk_delivery_records_parcel FOREIGN KEY (parcel_id) REFERENCES parcels(id) ON DELETE SET NULL
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS coupons (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  code VARCHAR(64) NOT NULL UNIQUE,
  name VARCHAR(100) NOT NULL,
  activity_rule VARCHAR(255) NOT NULL,
  amount DECIMAL(10,2) NOT NULL,
  threshold DECIMAL(10,2) NOT NULL DEFAULT 0.00,
  total INT NOT NULL,
  remaining INT NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'active',
  start_at DATETIME(3) NOT NULL,
  end_at DATETIME(3) NOT NULL,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS send_orders (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  order_no VARCHAR(64) NOT NULL UNIQUE,
  user_id BIGINT UNSIGNED NOT NULL,
  sender_name VARCHAR(100) NOT NULL,
  sender_phone VARCHAR(20) NOT NULL,
  sender_address VARCHAR(255) NOT NULL,
  receiver_name VARCHAR(100) NOT NULL,
  receiver_phone VARCHAR(20) NOT NULL,
  receiver_address VARCHAR(255) NOT NULL,
  item_info VARCHAR(255) NOT NULL,
  weight DECIMAL(10,2) NOT NULL,
  estimated_fee DECIMAL(10,2) NOT NULL,
  coupon_deduct DECIMAL(10,2) NOT NULL DEFAULT 0.00,
  coupon_id BIGINT UNSIGNED NULL,
  status VARCHAR(32) NOT NULL DEFAULT 'created',
  assigned_courier VARCHAR(100) NULL,
  pay_status VARCHAR(20) NOT NULL DEFAULT 'unpaid',
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  INDEX idx_send_orders_user_id (user_id),
  INDEX idx_send_orders_coupon_id (coupon_id),
  INDEX idx_send_orders_status (status),
  INDEX idx_send_orders_pay_status (pay_status),
  CONSTRAINT fk_send_orders_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_send_orders_coupon FOREIGN KEY (coupon_id) REFERENCES coupons(id) ON DELETE SET NULL
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS user_coupons (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT UNSIGNED NOT NULL,
  coupon_id BIGINT UNSIGNED NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'unused',
  received_at DATETIME(3) NOT NULL,
  used_at DATETIME(3) NULL,
  used_in_order_id BIGINT UNSIGNED NULL,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  UNIQUE KEY uk_user_coupon (user_id, coupon_id),
  INDEX idx_user_coupons_coupon_id (coupon_id),
  INDEX idx_user_coupons_status (status),
  INDEX idx_user_coupons_used_in_order_id (used_in_order_id),
  CONSTRAINT fk_user_coupons_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_user_coupons_coupon FOREIGN KEY (coupon_id) REFERENCES coupons(id) ON DELETE CASCADE,
  CONSTRAINT fk_user_coupons_order FOREIGN KEY (used_in_order_id) REFERENCES send_orders(id) ON DELETE SET NULL
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS payment_orders (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  pay_no VARCHAR(64) NOT NULL UNIQUE,
  user_id BIGINT UNSIGNED NOT NULL,
  related_type VARCHAR(32) NOT NULL,
  related_id BIGINT UNSIGNED NOT NULL,
  biz_desc VARCHAR(255) NULL,
  amount DECIMAL(10,2) NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'pending',
  pay_method VARCHAR(32) NULL,
  paid_at DATETIME(3) NULL,
  callback_payload TEXT NULL,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  INDEX idx_payment_orders_user_id (user_id),
  INDEX idx_payment_orders_related (related_type, related_id),
  INDEX idx_payment_orders_status (status),
  CONSTRAINT fk_payment_orders_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS notices (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT UNSIGNED NULL,
  type VARCHAR(32) NOT NULL,
  content VARCHAR(500) NOT NULL,
  sent_at DATETIME(3) NOT NULL,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  INDEX idx_notices_user_id (user_id),
  INDEX idx_notices_type (type),
  CONSTRAINT fk_notices_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS operation_logs (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT UNSIGNED NULL,
  action VARCHAR(64) NOT NULL,
  target_type VARCHAR(64) NOT NULL,
  target_id BIGINT UNSIGNED NULL,
  detail VARCHAR(500) NULL,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  INDEX idx_operation_logs_user_id (user_id),
  INDEX idx_operation_logs_action (action),
  INDEX idx_operation_logs_target_type (target_type),
  INDEX idx_operation_logs_target_id (target_id),
  CONSTRAINT fk_operation_logs_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB;

INSERT INTO coupons (code, name, activity_rule, amount, threshold, total, remaining, status, start_at, end_at)
VALUES
('NEWUSER10', '新用户券10元', '新用户首单满20可用', 10.00, 20.00, 10000, 10000, 'active', NOW(), DATE_ADD(NOW(), INTERVAL 365 DAY))
ON DUPLICATE KEY UPDATE code = VALUES(code);
