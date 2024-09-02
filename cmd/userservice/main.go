package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/CNMoreno/cnm-proyect-go/internal/adapters"
	"github.com/CNMoreno/cnm-proyect-go/internal/handlers"
	"github.com/CNMoreno/cnm-proyect-go/internal/repository"
	"github.com/CNMoreno/cnm-proyect-go/internal/usecase"
	"github.com/gin-gonic/gin"
)

func main() {
	mongoURI := os.Getenv("MONGO_URL")
	mongoDBName := os.Getenv("MONGO_DATABASE")

	mongoClient, err := adapters.NewMongoClient(mongoURI, mongoDBName)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	defer func() {
		if err := mongoClient.Close(); err != nil {
			log.Printf("Error closing MongoDB connection: %v", err)
		}
	}()

	userRepo := repository.NewUserService(mongoClient.GetDatabase())
	userService := usecase.NewUserService(userRepo)
	userHandlers := &handlers.UserHandlers{UserService: userService}

	r := gin.Default()

	route := "/users/:id"
	r.POST("/users", userHandlers.CreateUser)
	r.GET(route, userHandlers.GetUserByID)
	r.PATCH(route, userHandlers.UpdateUser)
	r.DELETE(route, userHandlers.DeleteUser)

	fmt.Println("Starting my microservice")

	log.Fatal(http.ListenAndServe(":8080", r))
}
