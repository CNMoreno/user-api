package usecase

import (
	"context"
	"time"

	"github.com/CNMoreno/cnm-proyect-go/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserService struct {
	userCollection *mongo.Collection
}

func NewUserService(db *mongo.Database) *UserService {
	return &UserService{
		userCollection: db.Collection("users"),
	}

}

func (s *UserService) CreateUser(ctx context.Context, user *domain.User) (string, error) {
	now := time.Now()
	user.ID = primitive.NewObjectID().Hex()
	user.CreatedAt = now
	user.UpdatedAt = now
	user.DelatedAt = now
	user.Enabled = true

	_, err := s.userCollection.InsertOne(ctx, user)

	if err != nil {
		return "", err
	}

	return user.ID, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {

	var user domain.User

	filter := bson.M{
		"_id":     id,
		"enabled": true,
	}

	err := s.userCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id string, updateFields map[string]interface{}) (*domain.User, error) {
	updateFields["updated_at"] = time.Now()

	filter := bson.M{
		"_id":     id,
		"enabled": true,
	}
	update := bson.M{"$set": updateFields}

	var updatedUser domain.User
	err := s.userCollection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updatedUser)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &updatedUser, nil

}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {

	userDelete := domain.User{
		DelatedAt: time.Now(),
		Enabled:   false,
	}

	update := bson.M{"$set": userDelete}

	result := s.userCollection.FindOneAndUpdate(ctx, bson.M{"_id": id}, update)

	if result.Err() != nil {
		return result.Err()
	}

	return nil
}
