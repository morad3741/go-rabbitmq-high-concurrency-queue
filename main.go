package main

import (
	"log"
	"net/http"

	"github.com/spf13/viper"
	"hexArchitectureProject/api/v1"
	"hexArchitectureProject/internal/db"
	"hexArchitectureProject/internal/mq"
	"hexArchitectureProject/internal/user"
)

func initConfig() {
	viper.SetConfigName("config") // Configuration file name without extension
	viper.SetConfigType("yaml")   // Configuration file type
	viper.AddConfigPath("config") // Path to look for the config file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
}

func main() {
	// Initialize configuration
	initConfig()

	// Load the database connection string
	dsn := viper.GetString("database.dsn")
	// RabbitMQ configuration
	rabbitURL := viper.GetString("rabbitmq.url")

	rabbitMQService, err := mq.NewRabbitMQService(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}
	defer rabbitMQService.Close()

	db, err := postgres.NewPostgresDB(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Create RabbitMQ HTTP handler
	rabbitHandler := v1.NewRabbitMQHandler(rabbitMQService)
	// Define routes
	http.HandleFunc("/api/v1/publish", rabbitHandler.SendMessageToRabbit)

	// Set up the user service and handler
	userRepo := user.NewUserRepository(db)
	userService := user.NewUserService(userRepo)
	userHandler := v1.NewUserHandler(userService)

	// Set up HTTP routes
	http.HandleFunc("/api/v1/register", userHandler.RegisterUser)

	// Start the server
	port := viper.GetString("server.port")
	log.Printf("Server is running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
