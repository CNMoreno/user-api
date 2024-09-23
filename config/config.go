package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/CNMoreno/cnm-proyect-go/internal/adapters"
	"github.com/CNMoreno/cnm-proyect-go/internal/constants"
	"github.com/CNMoreno/cnm-proyect-go/internal/handlers"
	"github.com/CNMoreno/cnm-proyect-go/internal/repository"
	"github.com/CNMoreno/cnm-proyect-go/internal/usecase"
	"github.com/CNMoreno/cnm-proyect-go/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SetupDependencies initializes all the dependencies required by the application.
// It returns the HTTP handlers, a cleanup function to close resources, and an error if any occurred during initialization.
func SetupDependencies() (*handlers.UserHandlers, func(), error) {
	mongoURI := os.Getenv("MONGO_URL")

	if mongoURI == "" {
		return nil, nil, fmt.Errorf(constants.ErrMongoUrlIsNotSet)
	}

	mongoDBName := os.Getenv("MONGO_DATABASE")
	if mongoDBName == "" {
		return nil, nil, fmt.Errorf(constants.ErrMongoDatabaseIsNotSet)
	}

	mongoClient, err := adapters.NewMongoClient(mongoURI, mongoDBName)
	if err != nil {
		return nil, nil, err
	}

	userCollection := mongoClient.GetDatabase().Collection("users")

	err = createUniqueIndexes(userCollection)

	if err != nil {
		log.Fatalf("%v: %v", constants.ErrCreateMongoIndex, err)
	}
	userRepo := repository.NewUserRepository(userCollection, utils.HashPassword)

	userService := usecase.NewUserService(userRepo)
	utils.NewValidator()
	userHandlers := &handlers.UserHandlers{
		UserService: userService,
	}

	cleanup := func() {
		if err := mongoClient.Close(); err != nil {
			log.Printf("%v: %v", constants.ErrCloseMongoConnection, err)
		}
	}

	return userHandlers, cleanup, nil
}

func createUniqueIndexes(collection *mongo.Collection) error {
	emailIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{
				Key:   "email",
				Value: 1,
			},
		},
		Options: options.Index().SetUnique(true),
	}

	userNameIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{
				Key:   "userName",
				Value: 1,
			},
		},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{emailIndexModel, userNameIndexModel})

	if err != nil {
		return err
	}
	return nil
}
