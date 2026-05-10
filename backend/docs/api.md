# Parcel Station API Guide

## 1. Common
- Base URL: `http://localhost:8080`
- Auth: JWT, header `Authorization: Bearer <token>`
- Unified response:
```json
{
  "code": 0,
  "msg": "success",
  "data": {}
}
```
- Roles:
  - `user`
  - `admin`

## 2. User APIs
### `POST /api/user/register`
- Register user

### `POST /api/user/login`
- Login (`account` can be username or phone)

### `GET /api/user/profile` (Auth)
- Get profile

### `PUT /api/user/profile` (Auth)
- Update profile

### `GET /api/user/pickup-history` (Auth)
- 描述: 查看取件历史（用户端“取件历史”页面数据源）
- 返回说明:
  - `data` 当前返回数组
  - 为保证前端不显示 `-`，以下字段每条记录必返（空值返回空字符串）:
    - `tracking_no`: 快递单号
    - `pickup_code`: 存储码/取件码
    - `location`: 取件位置/货架号
    - `status`: 状态（`pending/picked_up/delivered/lost`）
    - `pickup_time`: 取件时间（ISO8601）
- 返回示例（`data` 为数组）:
```json
[
  {
    "id": 3,
    "parcel_id": 4,
    "tracking_no": "SF12345",
    "pickup_code": "123456",
    "location": "A-01-03",
    "status": "picked_up",
    "pickup_user_id": 2,
    "pickup_user_name": "admin",
    "pickup_time": "2026-05-06T10:16:45.038+08:00",
    "operator_user_id": 1,
    "created_at": "2026-05-06T10:16:45.039+08:00"
  }
]
```
- 字段说明:
  - `id`: 取件记录 ID
  - `parcel_id`: 包裹 ID
  - `tracking_no`: 快递单号
  - `pickup_code`: 存储码/取件码（前端必显字段）
  - `location`: 取件位置（前端必显字段）
  - `status`: 记录状态（前端必显字段）
  - `pickup_user_id`: 取件人用户 ID
  - `pickup_user_name`: 取件人账号/姓名
  - `pickup_time`: 取件时间
  - `operator_user_id`: 操作员 ID
  - `created_at`: 记录创建时间

### `GET /api/user/qrcode` (Auth)
- Get identity code

## 3. Parcel APIs
### `POST /api/parcel/inbound` (Admin)
- Inbound parcel
- Validation:
  - `tracking_no`: 6-64 chars, letters/numbers/`_`/`-`
  - `recipient_phone`: China mainland mobile format (11 digits, starts with `1`)
  - `location`: max 128 chars

### `POST /api/parcel/outbound` (Auth)
- Outbound parcel

### `GET /api/parcel/status?tracking_no=...` (Auth)
- Query parcel status

### `GET /api/parcel/list` (Admin)
- List all parcels

## 4. Pickup/Delivery APIs
### `POST /api/pickup/code` (Admin)
- Generate pickup code
- Validation:
  - At least one of `tracking_no` or `recipient_phone` is required
  - If both are provided, `recipient_phone` must match the parcel of `tracking_no`
  - `tracking_no`: 6-64 chars, letters/numbers/`_`/`-`
  - `recipient_phone`: China mainland mobile format

### `POST /api/pickup/record` (Admin)
- Record pickup info
- Validation:
  - `tracking_no`: 6-64 chars, letters/numbers/`_`/`-`
  - `pickup_user_name`: max 100 chars

### `POST /api/delivery/record` (Admin)
- Record delivery info
- Validation:
  - `tracking_no`: 6-64 chars, letters/numbers/`_`/`-`
  - `courier_name`: max 100 chars
  - `delivery_status`: max 32 chars

## 5. Send Order APIs
### `POST /api/send/order` (Auth)
- Create send order

### `GET /api/admin/send/orders` (Admin)
- List send orders (supports `status` query)

### `POST /api/admin/send/process` (Admin)
- 描述: 处理寄件订单
- Body 示例:
```json
{
  "order_no": "SO202605051401371405",
  "action": "assign_pickup",
  "courier_name": "王师傅"
}
```
- `action` 支持:
  - `accept`: 接单
  - `assign_pickup`: 分配揽件员（需传 `courier_name`）
  - `complete`: 完成订单
- `complete` 行为说明:
  - 系统会自动为收件侧创建入库包裹（若该单尚未生成对应入库记录）
  - 自动入库包裹字段规则:
    - `tracking_no = order_no`
    - `location = AUTO-INBOUND-RACK`
    - `pickup_code = 6位数字取件码`
  - 系统会按 `receiver_phone` 精确查询站内用户；查到后自动发送“到站通知”
  - 若未查到对应用户账号，则 `receiver_bound=false` 且 `notice_sent=false`
- 返回要求（`action=complete`，`data` 字段）:
  - 为保证前端可直接展示完成结果，以下字段必返（即使为 `false` 或空值也返回字段）:
    - `order_no`: 订单号
    - `status`: 订单状态（完成后建议返回 `completed`）
    - `inbound_parcel`: 自动入库包裹信息对象
      - `tracking_no`: 入库包裹单号（应等于 `order_no`）
      - `location`: 入库位置（默认 `AUTO-INBOUND-RACK`）
      - `pickup_code`: 6位取件码
    - `receiver_bound`: 收件人手机号是否已绑定站内用户（布尔值）
    - `notice_sent`: 是否已发送到站通知（布尔值）
- 成功响应示例:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "order_no": "SO202605051401371405",
    "status": "completed",
    "inbound_parcel": {
      "tracking_no": "SO202605051401371405",
      "location": "AUTO-INBOUND-RACK",
      "pickup_code": "381926"
    },
    "receiver_bound": true,
    "notice_sent": true
  }
}
```

## 6. Coupon APIs
### `GET /api/admin/coupon/list` (Admin)
- 描述: 管理员分页查看已发放红包礼券
- Query params:
  - `status`: 可选，按状态筛选
    - 建议枚举：`active`（生效中）/`inactive`（已失效）/`expired`（已结束）
  - `keyword`: 可选，按 `code` 或 `name` 模糊查询
  - `page`: 可选，默认 `1`
  - `page_size`: 可选，默认 `20`，最大 `200`
- Response `data` 示例:
```json
{
  "list": [
    {
      "id": 1,
      "code": "SALE618",
      "name": "618 Coupon",
      "amount": 8.8,
      "threshold": 20,
      "total": 5000,
      "remaining": 4999,
      "status": "active"
    }
  ],
  "total": 1,
  "page": 1,
  "page_size": 20
}
```
- 返回字段说明:
  - `list`: 当前页数据
  - `total`: 总条数
  - `page`: 当前页码（从 1 开始）
  - `page_size`: 每页条数
  - `status`: 券状态，建议值 `active/inactive/expired`

### `POST /api/admin/coupon/create` (Admin)
- Create coupon manually
- Body example:
```json
{
  "name": "618 Coupon",
  "amount": 8.8,
  "code": "SALE618",
  "activity_rule": "618 event",
  "threshold": 20,
  "total": 5000,
  "valid_days": 60
}
```
- Fields:
  - `name`: required
  - `amount`: required, must be > 0
  - `code`: optional, auto-generated if empty
  - `activity_rule`: optional
  - `threshold`: optional, default `0`
  - `total`: optional, default `1000`
  - `valid_days`: optional, default `30`

### `POST /api/coupon/receive` (Auth)
- Receive coupon by code
- Body:
```json
{
  "coupon_code": "SALE618"
}
```

### `GET /api/coupon/my` (Auth)
- List my coupons

### `POST /api/coupon/use` (Auth)
- Use coupon for send order

## 7. Payment APIs
### `POST /api/pay/create` (Auth)
- Create payment order
- `related_type`:
  - `send_order` (needs `order_no`)
  - `storage_fee` (needs `related_id` and `amount`)

### `POST /api/pay/callback`
- Mock payment callback
- `status` of `success` or `paid` means success

### `GET /api/pay/bill` (Auth)
- Query bills

## 8. Notice APIs
### `POST /api/notice/send` (Admin)
- Send notice

### `GET /ws/notify?token=...`
- WebSocket push channel

## 9. Default Admin
- Username: `admin`
- Password: `admin123`
