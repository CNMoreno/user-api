package config

import (
	"fmt"
	"log"
	"os"

	"github.com/CNMoreno/cnm-proyect-go/internal/adapters"
	"github.com/CNMoreno/cnm-proyect-go/internal/handlers"
	"github.com/CNMoreno/cnm-proyect-go/internal/repository"
	"github.com/CNMoreno/cnm-proyect-go/internal/usecase"
	"github.com/CNMoreno/cnm-proyect-go/internal/utils"
)

// SetupDependencies initializes all the dependencies required by the application.
// It returns the HTTP handlers, a cleanup function to close resources, and an error if any occurred during initialization.
func SetupDependencies() (*handlers.UserHandlers, func(), error) {
	mongoURI := os.Getenv("MONGO_URL")

	if mongoURI == "" {
		return nil, nil, fmt.Errorf("MONGO_URL is not set")
	}

	mongoDBName := os.Getenv("MONGO_DATABASE")
	if mongoDBName == "" {
		return nil, nil, fmt.Errorf("MONGO_DATABASE is not set")
	}

	mongoClient, err := adapters.NewMongoClient(mongoURI, mongoDBName)
	if err != nil {
		return nil, nil, err
	}

	userRepo := repository.NewUserRepository(mongoClient.GetDatabase())
	userService := usecase.NewUserService(userRepo)
	utils.NewValidator()
	userHandlers := &handlers.UserHandlers{
		UserService: userService,
	}

	cleanup := func() {
		if err := mongoClient.Close(); err != nil {
			log.Printf("Error closing MongoDB connection: %v", err)
		}
	}

	return userHandlers, cleanup, nil
}
