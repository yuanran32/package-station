# 快递驿站管理系统 API 文档

## 1. 通用说明
- Base URL: `http://localhost:8080`
- 认证方式: JWT（`Authorization: Bearer <token>`）
- 统一返回格式:
```json
{
  "code": 0,
  "msg": "success",
  "data": {}
}
```
- 角色:
  - `user`: 普通用户
  - `admin`: 管理员

## 2. 用户管理
### `POST /api/user/register`
- 描述: 用户注册
- Body:
```json
{
  "username": "zhangsan",
  "password": "123456",
  "name": "张三",
  "phone": "13800000000"
}
```

### `POST /api/user/login`
- 描述: 用户登录（`account` 可填用户名或手机号）

### `GET /api/user/profile` (Auth)
- 描述: 查看个人资料

### `PUT /api/user/profile` (Auth)
- 描述: 修改个人资料

### `GET /api/user/pickup-history` (Auth)
- 描述: 查看取件历史

### `GET /api/user/qrcode` (Auth)
- 描述: 获取个人身份码（`identity_code`）

## 3. 快件管理
### `POST /api/parcel/inbound` (Admin)
- 描述: 快递入库
- Body:
```json
{
  "tracking_no": "SF1234567890",
  "location": "A-01-03",
  "recipient_phone": "13800000000"
}
```

### `POST /api/parcel/outbound` (Auth)
- 描述: 快递出库（取件）
- Body:
```json
{
  "tracking_no": "SF1234567890",
  "pickup_code": "123456"
}
```

### `GET /api/parcel/status?tracking_no=...` (Auth)
- 描述: 查询快递状态和存储位置

### `GET /api/parcel/list` (Admin)
- 描述: 查询所有快递

## 4. 取派件管理
### `POST /api/pickup/code` (Admin)
- 描述: 生成取件码（按 `tracking_no` 或 `recipient_phone`）

### `POST /api/pickup/record` (Admin)
- 描述: 记录取件信息

### `POST /api/delivery/record` (Admin)
- 描述: 记录派件信息

## 5. 寄件管理
### `POST /api/send/order` (Auth)
- 描述: 用户提交寄件订单

### `GET /api/admin/send/orders` (Admin)
- 描述: 管理员查看寄件订单（支持 `status` 查询）

### `POST /api/admin/send/process` (Admin)
- 描述: 处理寄件订单，`action` 支持:
  - `accept`
  - `assign_pickup`（需 `courier_name`）
  - `complete`

## 6. 红包礼券
### `POST /api/coupon/receive` (Auth)
- 描述: 领取红包券
- Body:
```json
{
  "coupon_code": "NEWUSER10"
}
```

### `GET /api/coupon/my` (Auth)
- 描述: 查看我的红包券

### `POST /api/coupon/use` (Auth)
- 描述: 下单后使用红包抵扣

## 7. 支付结算
### `POST /api/pay/create` (Auth)
- 描述: 创建支付订单
- `related_type` 支持:
  - `send_order`（需传 `order_no`）
  - `storage_fee`（需传 `related_id` + `amount`）

### `POST /api/pay/callback`
- 描述: 模拟支付回调
- `status` 传 `success/paid` 表示成功

### `GET /api/pay/bill` (Auth)
- 描述: 账单查询

## 8. 通知提醒
### `POST /api/notice/send` (Admin)
- 描述: 发送通知（可指向某个用户）

### `GET /ws/notify?token=...`
- 描述: WebSocket 实时通知

## 9. 默认管理员
- 系统启动后自动初始化管理员（若不存在）:
  - 用户名: `admin`
  - 密码: `admin123`