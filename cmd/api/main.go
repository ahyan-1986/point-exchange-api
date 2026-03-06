package main

import (
	"database/sql"
	"log"
	_ "point-exchange-api/docs"

	dbpkg "point-exchange-api/internal/db"
	adminHandlers "point-exchange-api/internal/handlers"
	healthHandlers "point-exchange-api/internal/handlers"
	partnerHandlers "point-exchange-api/internal/handlers"
	rateHandlers "point-exchange-api/internal/handlers"
	swapHandlers "point-exchange-api/internal/handlers"
	"point-exchange-api/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"

	// Swagger
	swaggerFiles "github.com/swaggo/files"

	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {

	// 连接数据库
	connStr := "postgresql://admin:admin123@103.103.157.147:30432/point_ledger?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	adminRepo := dbpkg.NewAdminRepository(db)
	adminService := &services.AdminService{Repo: adminRepo}
	adminHandlers.AdminService = adminService

	r := gin.Default()
	// Enable CORS for all origins (demo)
	r.Use(cors.Default())

	// --- Dependency Injection ---
	swapRepo := dbpkg.NewSwapRepository(db)
	swapService := &services.SwapService{Repo: swapRepo}
	swapHandlers.SwapService = swapService

	partnerRepo := dbpkg.NewPartnerRepository(db)
	partnerService := &services.PartnerService{Repo: partnerRepo}
	partnerHandlers.PartnerService = partnerService

	rateRepo := dbpkg.NewRateRepository(db)
	rateService := &services.RateService{Repo: rateRepo}
	rateHandlers.RateService = rateService

	// Swagger UI endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 注册路由
	v1 := r.Group("/v1")
	{
		// Partner management
		v1.POST("/partners", partnerHandlers.RegisterPartner)
		v1.GET("/partners", partnerHandlers.ListPartners)
		v1.GET("/partners/:id", partnerHandlers.GetPartner)
		v1.PATCH("/partners/:id/activate", partnerHandlers.ActivatePartner)

		// Point rate management
		v1.POST("/partners/:id/rates", rateHandlers.AddOrUpdateRate)
		v1.GET("/partners/:id/rates", rateHandlers.ListRates)

		// Swap/transaction
		v1.POST("/swap/deposit", swapHandlers.CreateDeposit)
		v1.GET("/swap/:id", swapHandlers.GetSwap)
		v1.GET("/swap/claims", swapHandlers.ClaimSwaps)
		v1.POST("/swap/confirm", swapHandlers.ConfirmSwap)

		// Custom filter swap list
		v1.GET("/swaps", swapHandlers.ListSwaps)

		// New: List swaps requested/received by partner
		v1.GET("/partners/:id/swaps/requested", swapHandlers.ListSwapsRequested)
		v1.GET("/partners/:id/swaps/received", swapHandlers.ListSwapsReceived)

		// Health & admin
		v1.GET("/health", healthHandlers.HealthCheck)
		v1.GET("/admin/ledger", adminHandlers.AdminLedger)

		// Version info
		v1.GET("/version", healthHandlers.VersionCheck)
	}

	r.Run(":8080") // 在 8080 端口启动
}
