package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	controllers "github.com/ChronoPlay/chronoplay-backend-service/controllers"
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
	db := database.MongoClient.Database(dbName).Collection("users")

	// Setup dependency injection for user
	userRepo := models.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	cardRepo := models.NewCardRepository(db)
	cardService := services.NewCardService(cardRepo)
	cardController := controllers.NewCardController(cardService)

	loanRepo := models.NewLoanRepository(db)
	loanService := services.NewLoanService(loanRepo)
	loanController := controllers.NewLoanController(loanService)

	cardTransactionRepo := models.NewCardTransactionRepository(db)
	cashTransactionRepo := models.NewCashTransactionRepository(db)
	transactionService := services.NewTransactionService(cardTransactionRepo, cashTransactionRepo)
	transactionController := controllers.NewTransactionController(transactionService)

	// Setup Gin and routes
	router := gin.Default()
	router.Use(cors.Default())
	routes.SetupRoutes(router, userController, cardController, loanController, transactionController)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server is running on http://localhost:%s", port)
	router.Run(":" + port)
}
