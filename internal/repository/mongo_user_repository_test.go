package repository_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/CNMoreno/cnm-proyect-go/internal/domain"
	"github.com/CNMoreno/cnm-proyect-go/internal/repository"
	mocks "github.com/CNMoreno/cnm-proyect-go/mocks/repository"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type valuesTestCases struct {
	name         string
	body         *domain.User
	id           string
	isError      bool
	userResponse *domain.User
	err          error
}

var userRequest = &domain.User{
	Name:     "Cristian",
	Email:    "cristian@gmail.com",
	Password: "Test123*",
	UserName: "cristian",
}

type MockCollection struct {
	mock.Mock
}

func (m *MockCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error){
	args := m.Called(ctx, document)
	return &mongo.InsertOneResult{InsertedID: "someId"}, 
}

func TestCreateUser(t *testing.T) {
	testCases := []valuesTestCases{
		{
			name: "should create user when user call api create user",
			body: userRequest,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := new(mocks.UserRepository)

			mockRepo.On("CreateUser", mock.Anything, test.body).Return("12345", test.err)
			
			mongoRepos := repository.NewUserRepository("user")
 
		})
	}

}
