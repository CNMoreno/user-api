package repository_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/CNMoreno/cnm-proyect-go/internal/domain"
	"github.com/CNMoreno/cnm-proyect-go/internal/repository"
	mocks "github.com/CNMoreno/cnm-proyect-go/mocks/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type valuesTestCases struct {
	name         string
	body         *domain.User
	bodyUsers    *[]domain.User
	id           string
	userResponse *domain.User
	isError      bool
	hashPassword string
	err          error
	errPassword  error
}

var userRequest = &domain.User{
	Name:     "Cristian",
	Email:    "cristian@gmail.com",
	Password: "Test123*",
	UserName: "cristian",
}

var usersRequest = &[]domain.User{
	{
		Name:     "Cristian",
		Email:    "cristian@gmail.com",
		Password: "Test123*",
		UserName: "cristian",
	},
}

// MockSingleResult mocks the mongo.SingleResult.
type MockSingleResult struct {
	mock.Mock
}

// Decode simulates decoding a result into a provided value.
func (m *MockSingleResult) Decode(v interface{}) error {
	args := m.Called(v)
	if user, ok := v.(*domain.User); ok {
		// Simulate the found user for a successful Decode call
		user.ID = "12345"
		user.Email = "test@example.com"
		user.Password = "hashedpassword"
	}
	return args.Error(0)
}

func TestCreateUser(t *testing.T) {
	testCases := []valuesTestCases{
		{
			name:         "should create user when method is called",
			body:         userRequest,
			hashPassword: "hashPassword",
		},
		{
			name:        "should throw an error when hash password fails",
			body:        userRequest,
			isError:     true,
			errPassword: errors.New("hash password error"),
		},
		{
			name:    "should throw an error when database fails",
			body:    userRequest,
			isError: true,
			err:     errors.New("create user error"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			mockCollection := new(mocks.MongoCollectionInterface)
			userService := repository.NewUserRepository(mockCollection, func(s string) (string, error) {
				return test.hashPassword, test.errPassword
			})
			ctx := context.Background()

			mockCollection.On("InsertOne", ctx, mock.Anything).Return(&mongo.InsertOneResult{InsertedID: "12345"}, test.err).Once()

			userID, err := userService.CreateUser(ctx, test.body)

			if test.isError {
				assert.Error(t, err)
			} else {
				assert.NotEmpty(t, userID)

			}

		})
	}

}

func TestCreateUserBatch(t *testing.T) {
	testCases := []valuesTestCases{
		{
			name:         "should create users when method is called",
			bodyUsers:    usersRequest,
			hashPassword: "hashPassword",
		},
		{
			name:        "should throw an error when hash password fails",
			bodyUsers:   usersRequest,
			isError:     true,
			errPassword: errors.New("hash password error"),
		},
		{
			name:      "should throw an error when database fails",
			bodyUsers: usersRequest,
			isError:   true,
			err:       errors.New("create users error"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			mockCollection := new(mocks.MongoCollectionInterface)
			userService := repository.NewUserRepository(mockCollection, func(s string) (string, error) {
				return test.hashPassword, test.errPassword
			})
			ctx := context.Background()

			mockCollection.On("InsertMany", ctx, mock.Anything).Return(&mongo.InsertManyResult{InsertedIDs: []interface{}{"12345", "67890"}}, test.err).Once()

			userIDs, err := userService.CreateUserBatch(ctx, test.bodyUsers)

			if test.isError {
				assert.Error(t, err)
			} else {
				assert.NotEmpty(t, userIDs)
			}
		})
	}
}

func TestGetUserByID(t *testing.T) {
	testCases := []valuesTestCases{
		{
			name: "should get user by id when method is called",
			id:   "123456",
		},
		{
			name:    "should throw an error when get user in database fail",
			id:      "12233",
			isError: true,
			err:     errors.New("get user error"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			mockCollection := new(mocks.MongoCollectionInterface)
			userService := repository.NewUserRepository(mockCollection, func(s string) (string, error) {
				return test.hashPassword, test.errPassword
			})
			ctx := context.Background()

			userDoc := bson.M{
				"_id":      "12345",
				"email":    "test@example.com",
				"password": "hashedpassword",
			}

			singleResult := mongo.NewSingleResultFromDocument(userDoc, test.err, nil)

			mockCollection.On("FindOne", ctx, mock.Anything).Return(singleResult, test.err).Once()

			user, err := userService.GetUserByID(ctx, test.id)

			fmt.Println(user)

			if test.isError {
				assert.Error(t, err)
			} else {
				assert.NotEmpty(t, user)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	testCases := []valuesTestCases{
		{
			name:         "should get user by id and update it when method is called",
			body:         userRequest,
			hashPassword: "hashPassword",
			id:           "123456",
		},
		{
			name:    "should throw an error when update user database fail",
			body:    userRequest,
			isError: true,
			err:     errors.New("delete user error"),
		},
		{
			name:        "should throw an error when hash password fails",
			id:          "123456",
			body:        userRequest,
			isError:     true,
			errPassword: errors.New("hash password error"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			mockCollection := new(mocks.MongoCollectionInterface)
			userService := repository.NewUserRepository(mockCollection, func(s string) (string, error) {
				return test.hashPassword, test.errPassword
			})
			ctx := context.Background()

			userDoc := bson.M{
				"_id":      "12345",
				"email":    "test@example.com",
				"password": "hashedpassword",
			}

			singleResult := mongo.NewSingleResultFromDocument(userDoc, test.err, nil)

			mockCollection.On("FindOneAndUpdate", ctx, mock.Anything, mock.Anything, mock.Anything).Return(singleResult, test.err).Once()

			user, err := userService.UpdateUser(ctx, test.id, test.body)

			if test.isError {
				assert.Error(t, err)
			} else {
				assert.NotEmpty(t, user)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	testCases := []valuesTestCases{
		{
			name: "should delete user by id when method is called",
			id:   "123456",
		},
		{
			name:    "should throw an error when delete user database fail",
			isError: true,
			err:     errors.New("delete user error"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			mockCollection := new(mocks.MongoCollectionInterface)
			userService := repository.NewUserRepository(mockCollection, func(s string) (string, error) {
				return test.hashPassword, test.errPassword
			})
			ctx := context.Background()

			userDoc := bson.M{
				"_id":      "12345",
				"email":    "test@example.com",
				"password": "hashedpassword",
			}

			singleResult := mongo.NewSingleResultFromDocument(userDoc, test.err, nil)

			mockCollection.On("FindOneAndUpdate", ctx, mock.Anything, mock.Anything).Return(singleResult, test.err).Once()

			err := userService.DeleteUser(ctx, test.id)

			if test.isError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
