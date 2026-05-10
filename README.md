# 快递驿站管理系统

一个前后端分离的快递驿站管理系统，覆盖用户注册登录、快递入库/出库、取件码、寄件订单、优惠券、支付账单、站内通知和管理员运营流程。

## 功能概览

- 用户端：注册登录、个人资料、身份码、包裹状态查询、取件历史、寄件下单、优惠券领取与使用、支付账单、站内通知。
- 管理端：包裹入库、包裹列表、取件码生成、取件/派件记录、寄件订单处理、优惠券创建与查询、通知发送。
- 后端服务：JWT 鉴权、角色权限控制、统一响应格式、MySQL 自动建库与表迁移、默认管理员初始化、WebSocket 通知推送。

## 技术栈

- 后端：Go、Gin、GORM、MySQL、JWT、Gorilla WebSocket
- 前端：Vue 3、Vite、Pinia、Vue Router、Element Plus、Axios
- 数据库：MySQL 8.x

## 目录结构

```text
package-station/
├── backend/
│   ├── cmd/server/          # 后端启动入口
│   ├── controllers/         # HTTP 控制器
│   ├── docs/                # API 文档
│   ├── internal/            # 配置、数据库、应用启动
│   ├── middleware/          # 鉴权、跨域、日志、错误处理
│   ├── models/              # GORM 模型与自动迁移
│   ├── pkg/                 # JWT、响应、校验、随机工具
│   ├── routes/              # 路由注册
│   └── services/            # 业务逻辑
├── frontend/
│   ├── src/api/             # 前端接口封装
│   ├── src/components/      # 通用组件
│   ├── src/router/          # 前端路由
│   ├── src/stores/          # Pinia 状态
│   └── src/views/           # 用户端和管理端页面
└── README.md
```

## 环境要求

- Go 1.26 或更高版本
- Node.js 18 或更高版本
- npm
- MySQL 8.x

确认本机环境：

```powershell
go version
node -v
npm -v
mysql --version
```

## 后端配置

后端通过环境变量读取配置；不配置时会使用默认值。

| 变量名 | 默认值 | 说明 |
| --- | --- | --- |
| `PORT` | `8080` | 后端监听端口 |
| `MYSQL_DSN` | `root:123456@tcp(127.0.0.1:3306)/parcel_station?charset=utf8mb4&parseTime=True&loc=Local` | MySQL 连接串 |
| `JWT_SECRET` | `change-this-secret` | JWT 签名密钥 |
| `TOKEN_EXPIRE_HOURS` | `72` | Token 有效期，单位小时 |

PowerShell 示例：

```powershell
$env:PORT="8080"
$env:MYSQL_DSN="root:123456@tcp(127.0.0.1:3306)/parcel_station?charset=utf8mb4&parseTime=True&loc=Local"
$env:JWT_SECRET="replace-with-a-strong-secret"
$env:TOKEN_EXPIRE_HOURS="72"
```

启动时如果数据库 `parcel_station` 不存在，后端会尝试自动创建数据库，并通过 GORM 自动迁移数据表。

## 启动后端

```powershell
cd backend
go mod download
go run ./cmd/server
```

启动成功后访问：

- 健康检查：`http://localhost:8080/health`
- API 根路径：`http://localhost:8080/`
- API 文档：`http://localhost:8080/docs/api.md`

默认管理员会在首次启动时自动初始化：

```text
username: admin
password: admin123
```

## 启动前端

```powershell
cd frontend
npm install
npm run dev
```

默认前端地址：

```text
http://localhost:5173
```

Vite 开发代理默认把 `/api` 和 `/ws` 转发到后端。可通过环境变量调整：

```powershell
$env:VITE_PROXY_TARGET="http://127.0.0.1:8080"
$env:VITE_WS_PROXY_TARGET="ws://127.0.0.1:8080"
npm run dev
```

生产构建：

```powershell
cd frontend
npm run build
npm run preview
```

## 常用接口摘要

统一响应格式：

```json
{
  "code": 0,
  "msg": "success",
  "data": {}
}
```

需要登录的接口使用 JWT：

```text
Authorization: Bearer <token>
```

核心接口：

| 模块 | 方法与路径 | 权限 | 说明 |
| --- | --- | --- | --- |
| 用户 | `POST /api/user/register` | 公开 | 用户注册 |
| 用户 | `POST /api/user/login` | 公开 | 用户登录，返回 Token |
| 用户 | `GET /api/user/profile` | 登录 | 获取个人资料 |
| 用户 | `PUT /api/user/profile` | 登录 | 修改个人资料 |
| 用户 | `GET /api/user/pickup-history` | 登录 | 查看取件历史 |
| 用户 | `GET /api/user/qrcode` | 登录 | 获取身份码 |
| 包裹 | `POST /api/parcel/inbound` | 管理员 | 快递入库 |
| 包裹 | `POST /api/parcel/outbound` | 登录 | 快递出库/取件 |
| 包裹 | `GET /api/parcel/status` | 登录 | 查询包裹状态 |
| 包裹 | `GET /api/parcel/list` | 管理员 | 包裹列表 |
| 取派件 | `POST /api/pickup/code` | 管理员 | 生成取件码 |
| 取派件 | `POST /api/pickup/record` | 管理员 | 记录取件 |
| 取派件 | `POST /api/delivery/record` | 管理员 | 记录派件 |
| 寄件 | `POST /api/send/order` | 登录 | 创建寄件订单 |
| 寄件 | `GET /api/admin/send/orders` | 管理员 | 查看寄件订单 |
| 寄件 | `POST /api/admin/send/process` | 管理员 | 处理寄件订单 |
| 优惠券 | `POST /api/admin/coupon/create` | 管理员 | 创建优惠券 |
| 优惠券 | `GET /api/admin/coupon/list` | 管理员 | 优惠券列表 |
| 优惠券 | `POST /api/coupon/receive` | 登录 | 领取优惠券 |
| 优惠券 | `GET /api/coupon/my` | 登录 | 我的优惠券 |
| 优惠券 | `POST /api/coupon/use` | 登录 | 使用优惠券 |
| 支付 | `POST /api/pay/create` | 登录 | 创建支付单 |
| 支付 | `POST /api/pay/callback` | 公开 | 模拟支付回调 |
| 支付 | `GET /api/pay/bill` | 登录 | 查看账单 |
| 通知 | `POST /api/notice/send` | 管理员 | 发送通知 |
| 通知 | `GET /ws/notify?token=...` | 登录 | WebSocket 通知 |

更详细的请求参数与响应示例见 [backend/docs/api.md](backend/docs/api.md)。

## 典型使用流程

1. 启动 MySQL。
2. 启动后端，等待自动建库、迁移表结构、初始化默认管理员。
3. 启动前端，打开 `http://localhost:5173`。
4. 使用 `admin/admin123` 登录管理端。
5. 管理员录入包裹或处理寄件订单。
6. 普通用户登录后查询包裹、取件、下寄件单、领取优惠券和查看账单。

## 开发检查

后端：

```powershell
cd backend
go test ./...
```

前端：

```powershell
cd frontend
npm run build
```

## 常见问题

### `go` 命令无法识别

确认 Go 已安装，并把 Go 安装目录加入系统 `PATH`。常见路径：

```text
C:\Program Files\Go\bin
```

重新打开终端后执行：

```powershell
go version
```

### 后端连接不上 MySQL

检查 MySQL 是否启动、账号密码是否正确，以及 `MYSQL_DSN` 中的主机、端口、数据库名是否正确。默认连接串使用：

```text
root:123456@tcp(127.0.0.1:3306)/parcel_station
```

### 前端请求失败

确认后端已启动，并检查 Vite 代理目标：

```powershell
$env:VITE_PROXY_TARGET="http://127.0.0.1:8080"
$env:VITE_WS_PROXY_TARGET="ws://127.0.0.1:8080"
```

如果前端和后端分开部署，需要设置 `VITE_API_BASE_URL` 指向后端地址。
