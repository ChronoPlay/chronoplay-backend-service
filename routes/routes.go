package routes

import (
	"github.com/gin-gonic/gin"

	controller "github.com/ChronoPlay/chronoplay-backend-service/controllers"
	middleware "github.com/ChronoPlay/chronoplay-backend-service/middlewares"
)

func SetupRoutes(r *gin.Engine, userController controller.UserController) {
	auth := r.Group("/auth", middleware.CustomContextMiddleware())

	{
		auth.POST("/signup", userController.RegisterUser)
		auth.GET("/user", userController.GetUser)
	}
}
