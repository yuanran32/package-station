package routes

import (
	"agent_learning/controllers"
	"agent_learning/internal/config"
	"agent_learning/middleware"
	"agent_learning/models"
	"agent_learning/pkg/jwtutil"
	"agent_learning/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(cfg *config.Config, db *gorm.DB) *gin.Engine {
	r := gin.New()
	r.Use(middleware.ErrorHandler(), middleware.CORS(), middleware.RequestLogger())

	jwtManager := jwtutil.NewManager(cfg.JWTSecret, cfg.TokenExpireHours)

	opLogSvc := services.NewOperationLogService(db)
	notifyHub := services.NewNotifyHub()
	noticeSvc := services.NewNoticeService(db, notifyHub)
	userSvc := services.NewUserService(db)
	parcelSvc := services.NewParcelService(db, opLogSvc, noticeSvc)
	pickupSvc := services.NewPickupService(db, opLogSvc, noticeSvc)
	sendSvc := services.NewSendOrderService(db, opLogSvc, noticeSvc)
	couponSvc := services.NewCouponService(db, opLogSvc)
	paySvc := services.NewPaymentService(db, opLogSvc)

	userCtl := controllers.NewUserController(userSvc, jwtManager)
	parcelCtl := controllers.NewParcelController(parcelSvc)
	pickupCtl := controllers.NewPickupController(pickupSvc)
	sendCtl := controllers.NewSendOrderController(sendSvc)
	couponCtl := controllers.NewCouponController(couponSvc)
	payCtl := controllers.NewPaymentController(paySvc)
	noticeCtl := controllers.NewNoticeController(noticeSvc, notifyHub, jwtManager)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "parcel-station-api",
			"status":  "ok",
			"docs":    "/docs/api.md",
		})
	})
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.Status(204)
	})
	r.Static("/docs", "./docs")
	r.GET("/ws/notify", noticeCtl.WebSocketNotify)

	api := r.Group("/api")

	user := api.Group("/user")
	{
		user.POST("/register", userCtl.Register)
		user.POST("/login", userCtl.Login)

		userAuth := user.Group("")
		userAuth.Use(middleware.JWTAuth(jwtManager))
		userAuth.GET("/profile", userCtl.GetProfile)
		userAuth.PUT("/profile", userCtl.UpdateProfile)
		userAuth.GET("/pickup-history", userCtl.PickupHistory)
		userAuth.GET("/qrcode", userCtl.QRCode)
	}

	parcel := api.Group("/parcel")
	parcel.Use(middleware.JWTAuth(jwtManager))
	{
		parcel.POST("/inbound", middleware.RoleAuth(models.RoleAdmin), parcelCtl.Inbound)
		parcel.POST("/outbound", parcelCtl.Outbound)
		parcel.GET("/status", parcelCtl.Status)
		parcel.GET("/list", middleware.RoleAuth(models.RoleAdmin), parcelCtl.List)
	}

	pickup := api.Group("/pickup")
	pickup.Use(middleware.JWTAuth(jwtManager), middleware.RoleAuth(models.RoleAdmin))
	{
		pickup.POST("/code", pickupCtl.GenerateCode)
		pickup.POST("/record", pickupCtl.RecordPickup)
	}

	delivery := api.Group("/delivery")
	delivery.Use(middleware.JWTAuth(jwtManager), middleware.RoleAuth(models.RoleAdmin))
	{
		delivery.POST("/record", pickupCtl.RecordDelivery)
	}

	send := api.Group("/send")
	send.Use(middleware.JWTAuth(jwtManager))
	{
		send.POST("/order", sendCtl.CreateOrder)
	}

	adminSend := api.Group("/admin/send")
	adminSend.Use(middleware.JWTAuth(jwtManager), middleware.RoleAuth(models.RoleAdmin))
	{
		adminSend.GET("/orders", sendCtl.ListOrders)
		adminSend.POST("/process", sendCtl.ProcessOrder)
	}

	coupon := api.Group("/coupon")
	coupon.Use(middleware.JWTAuth(jwtManager))
	{
		coupon.POST("/receive", couponCtl.Receive)
		coupon.GET("/my", couponCtl.MyCoupons)
		coupon.POST("/use", couponCtl.Use)
	}

	adminCoupon := api.Group("/admin/coupon")
	adminCoupon.Use(middleware.JWTAuth(jwtManager), middleware.RoleAuth(models.RoleAdmin))
	{
		adminCoupon.GET("/list", couponCtl.AdminList)
		adminCoupon.POST("/create", couponCtl.AdminCreate)
	}

	pay := api.Group("/pay")
	{
		pay.POST("/callback", payCtl.Callback)
		payAuth := pay.Group("")
		payAuth.Use(middleware.JWTAuth(jwtManager))
		payAuth.POST("/create", payCtl.Create)
		payAuth.GET("/bill", payCtl.Bills)
	}

	notice := api.Group("/notice")
	notice.Use(middleware.JWTAuth(jwtManager), middleware.RoleAuth(models.RoleAdmin))
	{
		notice.POST("/send", noticeCtl.Send)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "route not found",
			"data": nil,
		})
	})

	return r
}
