package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CNMoreno/cnm-proyect-go/internal/domain"
	"github.com/CNMoreno/cnm-proyect-go/internal/usecase"
	mocks "github.com/CNMoreno/cnm-proyect-go/mocks/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type createTestStruct struct {
	name        string
	body        domain.User
	isError     bool
	err         error
	statusError int
	isErrorBody bool
}

func TestCreateUser(t *testing.T) {
	testCases := []createTestStruct{
		{
			name: "should create user",
			body: domain.User{
				Name:  "Cristian",
				Email: "cristian@gmail.com",
			},
		},
		{
			name:        "should return an error when is an invalid body",
			isError:     true,
			statusError: http.StatusBadRequest,
			isErrorBody: true,
		},
		{
			name: "should return an error when bd return an error",
			err:  errors.New("some error"),
			body: domain.User{
				Name:  "Cristian",
				Email: "cristian@gmail.com",
			},
			isError:     true,
			statusError: http.StatusInternalServerError,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := new(mocks.UserRepository)

			userService := usecase.NewUserService(mockRepo)

			handler := UserHandlers{UserService: userService}

			router := gin.Default()
			router.POST("/users", handler.CreateUser)

			bodyBytes, _ := json.Marshal(test.body)

			mockRepo.On("CreateUser", mock.Anything, &test.body).Return("12345", test.err)

			id, err := userService.CreateUser(context.Background(), &test.body)

			req, _ := mockRequestEndPoint(test.isErrorBody, bodyBytes)

			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			if test.isError {
				assert.Equal(t, test.statusError, resp.Code)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "12345", id)
				assert.Equal(t, http.StatusCreated, resp.Code)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func mockRequestEndPoint(isError bool, body []byte) (*http.Request, error) {
	if isError {
		return http.NewRequest("POST", "/users", strings.NewReader("Invalid Body"))
	}

	return http.NewRequest("POST", "/users", bytes.NewBuffer(body))
}

/*
func TestCreateUser(t *testing.T) {
	testCases := []createTestStruct{
		{
			name: "should create user",
			body: domain.User{
				Name:  "Cristian",
				Email: "cristian@gmail.com",
			},
		},
		{
			name:    "should return an error",
			err:     errors.New("some error"),
			isError: true,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			mockRepo := new(mocks.UserRepository)

			userService := usecase.NewUserService(mockRepo)

			handler := handlers.UserHandlers{userService: userService}

			mockRepo.On("CreateUser", mock.Anything, &test.body).Return("12345", test.err)

			id, err := userService.CreateUser(context.Background(), &test.body)

			if test.isError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "12345", id)

			}
			mockRepo.AssertExpectations(t)
		})
	}

}
*/
