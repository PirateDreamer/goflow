package main

import (
	"fmt"
	"goflow-api/internal/application/service"
	"goflow-api/internal/infrastructure/config"
	"goflow-api/internal/infrastructure/persistence"
	"goflow-api/internal/interfaces/http"
	"goflow-api/internal/interfaces/http/controller"
	"log"
)

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
	userRepo := persistence.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userCtrl := controller.NewUserController(userSvc)

	// 4. 初始化路由器
	r := http.NewRouter(userCtrl)

	// 5. 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	fmt.Printf("Server starting on %s\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Run server failed: %v", err)
	}
}
