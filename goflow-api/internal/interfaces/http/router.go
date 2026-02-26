package http

import (
	"goflow-api/internal/interfaces/http/controller"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(userCtrl *controller.UserController) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "up",
		})
	})

	userGroup := r.Group("/users")
	{
		userGroup.POST("/register", userCtrl.Register)
		userGroup.POST("/login", userCtrl.Login)
		userGroup.GET("/:id", userCtrl.GetUser)
	}

	return r
}
