package main

import (
	"fmt"
	"goerp-api/internal/application/service"
	"goerp-api/internal/infrastructure/cache"
	"goerp-api/internal/infrastructure/config"
	"goerp-api/internal/infrastructure/email"
	"goerp-api/internal/infrastructure/persistence"
	"goerp-api/internal/interfaces/http"
	"goerp-api/internal/interfaces/http/controller"
	"log"
)

// @title GoERP API
// @version 1.0
// @description This is a sample ERP server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {
	// 1. 初始化配置
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf("Init config failed: %v", err)
	}

	// 2. 初始化数据库
	db, err := persistence.InitDB(&cfg.Database)
	if err != nil {
		log.Printf("Warning: Init DB failed (DSN: %s): %v", cfg.Database.DSN, err)
	} else {
		fmt.Println("Database connection established.")
	}

	// 3. 依赖注入
	redisCache := cache.NewRedisCache(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	emailSvc := email.NewSMTPService(cfg.Email.Host, cfg.Email.Port, cfg.Email.User, cfg.Email.Password, cfg.Email.From)

	userRepo := persistence.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo, redisCache, emailSvc)
	userCtrl := controller.NewUserController(userSvc)

	// 4. 初始化路由器
	r := http.NewRouter(userCtrl, &cfg.Swagger)

	// 5. 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	fmt.Printf("Server starting on %s\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Run server failed: %v", err)
	}
}
