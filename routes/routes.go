package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"

	controller "github.com/ChronoPlay/chronoplay-backend-service/controllers"
	middleware "github.com/ChronoPlay/chronoplay-backend-service/middlewares"
)

func SetupRoutes(r *gin.Engine, userController controller.UserController) {
	auth := r.Group("/auth", middleware.CustomContextMiddleware())

	fmt.Print("request has entered here- router \n")
	{
		auth.POST("/signup", userController.RegisterUser)
		auth.GET("/user", userController.GetUser)
	}
}
