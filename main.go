package main

import (
	"github.com/gin-gonic/gin"

	controller "github.com/ChronoPlay/chronoplay-backend-service/controllers"
	"github.com/ChronoPlay/chronoplay-backend-service/database"
	model "github.com/ChronoPlay/chronoplay-backend-service/models"
	"github.com/ChronoPlay/chronoplay-backend-service/routes"
	service "github.com/ChronoPlay/chronoplay-backend-service/services"
)

func main() {
	database.Connect()
	// config.LoadEnv()

	db := database.DB

	var userController controller.UserController
	{
		userRepository := model.NewUserRepository(db)
		userService := service.NewUserService(userRepository)
		// userService = userService.NewUserMiddleware()(userService)
		userController = controller.NewUserController(userService)
	}

	r := gin.Default()
	routes.SetupRoutes(r, userController)

	r.Run("localhost:8080")
}
