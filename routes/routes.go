package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"

	controller "github.com/ChronoPlay/chronoplay-backend-service/controllers"
	middleware "github.com/ChronoPlay/chronoplay-backend-service/middlewares"
)

func SetupRoutes(r *gin.Engine, userController controller.UserController, cardController controller.CardController, loanController controller.LoanController, transactionController controller.TransactionController) {
	auth := r.Group("/auth", middleware.CustomContextMiddleware())

	fmt.Print("request has entered here- router \n")
	{
		auth.POST("/signup", userController.RegisterUser)
		auth.GET("/verify", userController.VerifyUser)
		auth.POST("/login", userController.LoginUser)
	}

	user := r.Group("/user", middleware.AuthorizeUser(), middleware.CustomContextMiddleware())

	{
		user.GET("/get_user", userController.GetUser)
	}

	card := r.Group("/card", middleware.AuthorizeUser(), middleware.CustomContextMiddleware())
	{
		card.POST("/add", cardController.AddCard)
		card.GET("/get_card", cardController.GetCard)
	}

}
