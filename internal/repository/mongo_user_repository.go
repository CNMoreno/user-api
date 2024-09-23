package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/CNMoreno/cnm-proyect-go/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// IMongoCollectionInterface Define an interface for the collection
type IMongoCollectionInterface interface {
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	InsertMany(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult
}

// UserService struct of user in Mongo collection.
type UserService struct {
	userCollection IMongoCollectionInterface
	hashPassword   func(string) (string, error)
}

// NewUserRepository join to Mongo collection.
func NewUserRepository(collection IMongoCollectionInterface, hashFunc func(string) (string, error)) *UserService {
	return &UserService{
		userCollection: collection,
		hashPassword:   hashFunc,
	}
}

// CreateUser handles to create user in database.
func (s *UserService) CreateUser(ctx context.Context, user *domain.User) (string, error) {
	now := time.Now()
	user.ID = primitive.NewObjectID().Hex()
	user.CreatedAt = now
	user.UpdatedAt = now
	user.DeletedAt = now
	user.Enabled = true
	password, err := s.hashPassword(user.Password)
	if err != nil {
		return "", err
	}
	user.Password = password

	_, err = s.userCollection.InsertOne(ctx, user)

	if err != nil {
		return "", err
	}

	return user.ID, nil
}

// CreateUserBatch handles to create users in database.
func (s *UserService) CreateUserBatch(ctx context.Context, users *[]domain.User) ([]interface{}, error) {
	now := time.Now()

	var validUsers []interface{}

	for _, user := range *users {
		user.ID = primitive.NewObjectID().Hex()
		user.CreatedAt = now
		user.UpdatedAt = now
		user.DeletedAt = now
		user.Enabled = true
		password, err := s.hashPassword(user.Password)
		if err != nil {
			return nil, err
		}
		user.Password = password
		validUsers = append(validUsers, user)
	}

	usersIDs, err := s.userCollection.InsertMany(ctx, validUsers)

	if err != nil {
		return nil, err
	}

	return usersIDs.InsertedIDs, nil
}

// GetUserByID handles to obtain user by ID in database.
func (s *UserService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User

	filter := bson.M{
		"_id":     id,
		"enabled": true,
	}

	result := s.userCollection.FindOne(ctx, filter)
	fmt.Println(&result)
	err := result.Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUser handles to obtain and update user by ID in database.
func (s *UserService) UpdateUser(ctx context.Context, id string, updateFields *domain.User) (*domain.User, error) {
	password, err := s.hashPassword(updateFields.Password)
	if err != nil {
		return &domain.User{}, err
	}
	filter := bson.M{
		"_id":     id,
		"enabled": true,
	}
	update := bson.M{"$set": bson.M{
		"updatedAt": time.Now(),
		"name":      updateFields.Name,
		"email":     updateFields.Email,
		"password":  password,
		"userName":  updateFields.UserName,
	},
	}
	var updatedUser domain.User
	optionsUpdate := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err = s.userCollection.FindOneAndUpdate(ctx, filter, update, optionsUpdate).Decode(&updatedUser)

	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

// DeleteUser handles to obtain and delete user by ID in database.
func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	filter := bson.M{
		"_id":     id,
		"enabled": true,
	}

	update := bson.M{"$set": bson.M{
		"enabled":   false,
		"deletedAt": time.Now(),
	}}

	result := s.userCollection.FindOneAndUpdate(ctx, filter, update)

	if result.Err() != nil {
		return result.Err()
	}

	return nil
}
