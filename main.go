package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	controllers "github.com/ChronoPlay/chronoplay-backend-service/controllers"
	"github.com/ChronoPlay/chronoplay-backend-service/crons"
	"github.com/ChronoPlay/chronoplay-backend-service/database"
	models "github.com/ChronoPlay/chronoplay-backend-service/model"
	"github.com/ChronoPlay/chronoplay-backend-service/routes"
	services "github.com/ChronoPlay/chronoplay-backend-service/services"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found or error loading it")
	}

	// Connect to MongoDB
	database.ConnectMongo()

	// Get collection from MongoDB
	dbName := os.Getenv("MONGO_DB_NAME")
	if dbName == "" {
		log.Fatal("MONGO_DB_NAME environment variable not set")
	}
	usersDb := database.MongoClient.Database(dbName).Collection("users")
	cardDb := database.MongoClient.Database(dbName).Collection("cards")
	loanDb := database.MongoClient.Database(dbName).Collection("loans")
	cardTransactionDb := database.MongoClient.Database(dbName).Collection("card_transactions")
	cashTransactionDb := database.MongoClient.Database(dbName).Collection("cash_transactions")
	notificationDb := database.MongoClient.Database(dbName).Collection("notifications")

	cardRepo := models.NewCardRepository(cardDb)
	userRepo := models.NewUserRepository(usersDb)
	loanRepo := models.NewLoanRepository(loanDb)
	cardTransactionRepo := models.NewCardTransactionRepository(cardTransactionDb)
	cashTransactionRepo := models.NewCashTransactionRepository(cashTransactionDb)
	notificationRepo := models.NewNotificationRepository(notificationDb)

	notificationService := services.NewNotificationService(notificationRepo)
	userService := services.NewUserService(userRepo, cardRepo)
	cardService := services.NewCardService(cardRepo, userRepo)
	loanService := services.NewLoanService(loanRepo)
	transactionService := services.NewTransactionService(cardTransactionRepo, cashTransactionRepo, userRepo, cardRepo, notificationService)

	notificationController := controllers.NewNotificationController(notificationService)
	userController := controllers.NewUserController(userService)
	cardController := controllers.NewCardController(cardService)
	loanController := controllers.NewLoanController(loanService)
	transactionController := controllers.NewTransactionController(transactionService)

	// Setup Gin and routes
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:3001",
			"https://chronoplay-frontend.onrender.com",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	}))
	routes.SetupRoutes(router, userController, cardController, loanController, transactionController, notificationController)

	// start all cron jobs
	cronsEnabled := os.Getenv("CRON_ENABLED") == "true"
	cronController := crons.NewCronController(userService, notificationService, cronsEnabled)
	cronController.RunAllCrons()

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server is running on http://localhost:%s", port)
	router.Run(":" + port)
}
