package http

import (
	"goerp-api/internal/infrastructure/config"
	"goerp-api/internal/interfaces/http/controller"
	"net/http"

	_ "goerp-api/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(userCtrl *controller.UserController, cfg *config.SwaggerConfig) *gin.Engine {
	r := gin.Default()

	swaggerGroup := r.Group("/swagger")
	if cfg != nil && cfg.User != "" {
		swaggerGroup.Use(gin.BasicAuth(gin.Accounts{
			cfg.User: cfg.Password,
		}))
	}
	swaggerGroup.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "up",
		})
	})

	userGroup := r.Group("/users")
	{
		userGroup.POST("/register", userCtrl.Register)
		userGroup.POST("/login", userCtrl.Login)
		userGroup.POST("/send-code", userCtrl.SendEmailCode)
		userGroup.POST("/login-email", userCtrl.LoginByEmail)
		userGroup.GET("/:id", userCtrl.GetUser)
	}

	return r
}
