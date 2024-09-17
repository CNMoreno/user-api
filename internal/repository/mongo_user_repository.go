package repository

import (
	"context"
	"time"

	"github.com/CNMoreno/cnm-proyect-go/internal/domain"
	"github.com/CNMoreno/cnm-proyect-go/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserService struct of user in Mongo collection.
type UserService struct {
	userCollection *mongo.Collection
}

// NewUserRepository join to Mongo collection.
func NewUserRepository(db *mongo.Database) *UserService {
	return &UserService{
		userCollection: db.Collection("users"),
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
	password, err := utils.HashPassword(user.Password)
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
		password, err := utils.HashPassword(user.Password)
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

	err := s.userCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUser handles to obtain and update user by ID in database.
func (s *UserService) UpdateUser(ctx context.Context, id string, updateFields *domain.User) (*domain.User, error) {
	updateFields.DeletedAt = time.Now()
	updateFields.Enabled = true

	filter := bson.M{
		"_id":     id,
		"enabled": true,
	}
	update := bson.M{"$set": updateFields}

	var updatedUser domain.User
	err := s.userCollection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updatedUser)

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
