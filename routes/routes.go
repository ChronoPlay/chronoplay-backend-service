package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"

	controller "github.com/ChronoPlay/chronoplay-backend-service/controllers"
	middleware "github.com/ChronoPlay/chronoplay-backend-service/middlewares"
)

func SetupRoutes(r *gin.Engine, userController controller.UserController, cardController controller.CardController, loanController controller.LoanController, transactionController controller.TransactionController, notificationController controller.NotificationController) {
	auth := r.Group("/auth", middleware.CustomContextMiddleware())

	fmt.Print("request has entered here- router \n")
	{
		auth.POST("/signup", userController.RegisterUser)
		auth.GET("/verify", userController.VerifyUser)
		auth.POST("/login", userController.LoginUser)
		auth.PATCH("/activate_all_users", userController.ActivateAllUsers)
	}

	user := r.Group("/user", middleware.AuthorizeUser(), middleware.CustomContextMiddleware())

	{
		user.GET("/user", userController.GetUser)
		user.GET("/get_user", userController.GetUserById)
		user.PATCH("/add_friend", userController.AddFriend)
		user.GET("/get_friends", userController.GetFriends)
		user.PATCH("/remove_friend", userController.RemoveFriend)
	}

	card := r.Group("/card", middleware.AuthorizeUser(), middleware.CustomContextMiddleware())
	{
		card.POST("/add", cardController.AddCard)
		card.GET("/get_card", cardController.GetCard)
	}

	transaction := r.Group("/transaction", middleware.AuthorizeUser(), middleware.CustomContextMiddleware())
	{
		transaction.POST("/transfer_cash", transactionController.Transfercash)
		transaction.POST("/transfer_cards", transactionController.Transfercards)
		transaction.POST("/exchange", transactionController.Exchange)
		transaction.GET("/get_transactions", transactionController.GetTransactions)
		transaction.GET("/get_possible_exchange", transactionController.GetPossibleExchange)
		transaction.POST("/execute_exchange", transactionController.ExecuteExchange)
	}

	notification := r.Group("/notification", middleware.AuthorizeUser(), middleware.CustomContextMiddleware())
	{
		notification.GET("/get_notifications", notificationController.GetNotifications)
		notification.PATCH("/mark_as_read", notificationController.MarkAsRead)
	}

}
