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
- 描述: 查看取件历史（用户端“取件历史”页面数据源）
- 返回说明:
  - `data` 可直接返回数组，或返回分页对象（如 `list` + `total`）。
  - 为保证前端不显示 `-`，以下字段请在每条记录中**必返**（即使为空也要带字段）:
    - `tracking_no`: 快递单号（字符串）
    - `pickup_code`: 存储码/取件码（字符串）
    - `location`: 取件位置/货架号（字符串）
    - `status`: 状态（字符串，建议值 `pending/picked_up/delivered/lost`）
    - `pickup_time`: 取件时间（ISO8601 时间字符串）
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
  - `status`: 记录状态（前端必显字段，建议值见上）
  - `pickup_user_id`: 取件人用户 ID
  - `pickup_user_name`: 取件人账号/姓名
  - `pickup_time`: 取件时间
  - `operator_user_id`: 操作员 ID
  - `created_at`: 记录创建时间

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
- Body:
```json
{
  "sender_name": "张三",
  "sender_phone": "13800000000",
  "sender_address": "上海市浦东新区世纪大道100号",
  "receiver_name": "李四",
  "receiver_phone": "13900000000",
  "receiver_address": "北京市朝阳区建国路88号",
  "item_info": "衣物+书籍",
  "weight": 2.5
}
```
- 返回示例（`data`）:
```json
{
  "order_no": "SO202605051401371405",
  "estimated_fee": 12,
  "pay_status": "unpaid",
  "status": "created"
}
```

### `GET /api/admin/send/orders` (Admin)
- 描述: 管理员查看寄件订单（支持 `status` 查询）

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
    - `location = receiver_address`（使用该寄件订单的收件地址）
    - `pickup_code = 6位数字取件码`
  - 若收件人手机号已绑定用户账号，系统自动发送“到站通知”
- 返回要求（`action=complete`，`data` 字段）:
  - 为保证前端可直接展示完成结果，以下字段请**必返**（即使为 `false` 或空值也返回字段）:
    - `order_no`: 订单号
    - `status`: 订单状态（完成后建议返回 `completed`）
    - `inbound_parcel`: 自动入库包裹信息对象
      - `tracking_no`: 入库包裹单号（应等于 `order_no`）
      - `location`: 入库位置（默认取订单收件地址 `receiver_address`）
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
      "location": "北京市朝阳区建国路88号",
      "pickup_code": "381926"
    },
    "receiver_bound": true,
    "notice_sent": true
  }
}
```

## 6. 红包礼券
### `POST /api/admin/coupon/create` (Admin)
- 描述: 管理员手动创建红包礼券（可设置名字与额度）
- Body 示例:
```json
{
  "name": "618满减券",
  "amount": 8.8,
  "code": "SALE618",
  "activity_rule": "618活动券",
  "threshold": 20,
  "total": 5000,
  "valid_days": 60
}
```
- 字段说明:
  - `name`: 必填，券名称
  - `amount`: 必填，券面额
  - `code`: 可选，券码（不传则系统自动生成）
  - `activity_rule`: 可选，活动规则描述
  - `threshold`: 可选，最低使用门槛（默认 0）
  - `total`: 可选，发放总量（默认 1000）
  - `valid_days`: 可选，有效天数（默认 30）

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
- 返回常见字段:
  - `code`: 礼券码
  - `name`: 礼券名称
  - `amount`: 券面额
  - `threshold`: 使用门槛
  - `status`: `unused/used/expired`
  - `received_at`: 领取时间
  - `used_at`: 使用时间（若已使用）

### `POST /api/coupon/use` (Auth)
- 描述: 下单后使用红包抵扣
- Body:
```json
{
  "order_no": "SO202605051421031096",
  "user_coupon_id": 2
}
```
- 说明:
  - 仅可使用当前用户 `status=unused` 的券；
  - 订单金额需达到券门槛（例如 `threshold=20`）。

## 7. 支付结算
### `POST /api/pay/create` (Auth)
- 描述: 创建支付订单
- `related_type` 支持:
  - `send_order`（需传 `order_no`）
  - `storage_fee`（需传 `related_id` + `amount`）
- 示例（寄件订单支付）:
```json
{
  "related_type": "send_order",
  "order_no": "SO202605051401041354"
}
```

### `POST /api/pay/callback`
- 描述: 模拟支付回调
- `status` 传 `success/paid` 表示成功
- 示例:
```json
{
  "pay_no": "PAY202605051401265543",
  "status": "paid"
}
```

### `GET /api/pay/bill` (Auth)
- 描述: 账单查询

### 前端支付流程（当前实现）
1. 用户在“寄件下单”提交订单，前端拿到 `order_no`。
2. 前端先查询 `GET /api/coupon/my`，若有可用券会弹窗询问“是否使用红包”。
3. 若用户确认使用，前端调用 `POST /api/coupon/use`，成功后订单金额会被重算。
4. 前端再调用 `POST /api/pay/create`，参数：
```json
{
  "related_type": "send_order",
  "order_no": "SO202605051401371405"
}
```
5. 在“支付账单”页点击“立即支付”，前端调用 `POST /api/pay/callback`（模拟支付）：
```json
{
  "pay_no": "PAY202605051401375303",
  "status": "paid"
}
```

## 8. 通知提醒
### `POST /api/notice/send` (Admin)
- 描述: 发送通知（可指向某个用户）

### `GET /ws/notify?token=...`
- 描述: WebSocket 实时通知

## 9. 默认管理员
- 系统启动后自动初始化管理员（若不存在）:
  - 用户名: `admin`
  - 密码: `admin123`
