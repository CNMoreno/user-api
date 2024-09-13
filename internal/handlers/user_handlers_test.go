package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/CNMoreno/cnm-proyect-go/internal/domain"
	"github.com/CNMoreno/cnm-proyect-go/internal/security"
	"github.com/CNMoreno/cnm-proyect-go/internal/usecase"
	mocks "github.com/CNMoreno/cnm-proyect-go/mocks/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type valuesTestCases struct {
	name         string
	body         *domain.User
	id           string
	isError      bool
	isErrorBody  bool
	userResponse *domain.User
	err          error
	statusCode   int
}

const (
	errorValue = "some error"
	route      = "/users"
	withID     = "%v/:id"
)

var userRequest = &domain.User{
	Name:     "Cristian",
	Email:    "cristian@gmail.com",
	Password: "Test123*",
	UserName: "cristian",
}

var userResponse = &domain.User{
	Name:     "test",
	Email:    "test@gmail.com",
	UserName: "test",
}

func TestCreateUser(t *testing.T) {
	testCases := []valuesTestCases{
		{
			name:       "should create user",
			body:       userRequest,
			statusCode: http.StatusCreated,
		},
		{
			name:        "should return an error when is an invalid body",
			isError:     true,
			statusCode:  http.StatusBadRequest,
			isErrorBody: true,
		},
		{
			name:       "should return an error when bd return an error",
			err:        errors.New(errorValue),
			body:       userRequest,
			isError:    true,
			statusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			mockRepo, handler, router := configurations()

			router.POST(route, handler.CreateUser)

			bodyBytes, _ := json.Marshal(test.body)

			mockRepo.On("CreateUser", mock.Anything, test.body).Return("12345", test.err)

			req, _ := mockRequestEndPoint(test.isErrorBody, "POST", route, bytes.NewBuffer(bodyBytes))

			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			if test.isError {
				assert.Equal(t, test.statusCode, resp.Code)
			} else {
				var response domain.APIResponse
				err := json.Unmarshal(resp.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "12345", response.ID)
				assert.Equal(t, test.statusCode, resp.Code)
			}
		})
	}
}

func TestGetUserByID(t *testing.T) {
	testCases := []valuesTestCases{
		{
			name:         "should return a user successfully",
			id:           "12345",
			userResponse: userResponse,
			statusCode:   http.StatusCreated,
		},
		{
			name:       "should return an error when user by ID not exist in database",
			id:         "1234435",
			statusCode: http.StatusNotFound,
			isError:    true,
			err:        mongo.ErrNoDocuments,
		},
		{
			name:       "should return an error when bd return an error",
			err:        errors.New(errorValue),
			id:         "12345",
			isError:    true,
			statusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			mockRepo, handler, router := configurations()

			router.GET(fmt.Sprintf(withID, route), handler.GetUserByID)

			mockRepo.On("GetUserByID", mock.Anything, test.id).Return(test.userResponse, test.err)

			req, _ := mockRequestEndPoint(test.isError, "GET", fmt.Sprintf("%v/%v", route, test.id), nil)

			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			if test.isError {
				assert.Equal(t, test.statusCode, resp.Code)
			} else {

				assert.Equal(t, test.statusCode, resp.Code)
				var response domain.APIResponse
				err := json.Unmarshal(resp.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, test.userResponse.Name, response.Name)
				assert.Equal(t, test.userResponse.UserName, response.UserName)
				assert.Equal(t, test.userResponse.Email, response.Email)

			}
			mockRepo.AssertExpectations(t)
		})
	}

}

func TestUpdateUser(t *testing.T) {
	testCasesUpdate := []valuesTestCases{
		{
			name:         "should update and return a user successfully",
			id:           "12345",
			userResponse: userResponse,
			body:         userRequest,
			statusCode:   http.StatusOK,
		},
		{
			name:        "should return an error when is an invalid body for update user",
			id:          "12345",
			isError:     true,
			statusCode:  http.StatusBadRequest,
			isErrorBody: true,
		},
		{
			name:       "should return an error when user try update by ID does not exist in database",
			id:         "123443543654",
			statusCode: http.StatusNotFound,
			isError:    true,
			body:       userRequest,
			err:        mongo.ErrNoDocuments,
		},
		{
			name:       "should return an error when user try update by ID does not exist in database",
			id:         "123443543654",
			statusCode: http.StatusInternalServerError,
			isError:    true,
			body:       userRequest,
			err:        errors.New(errorValue),
		},
	}

	for _, test := range testCasesUpdate {
		t.Run(test.name, func(t *testing.T) {

			mockRepo, handler, router := configurations()

			router.PATCH(fmt.Sprintf(withID, route), handler.UpdateUser)

			bodyBytes, _ := json.Marshal(test.body)

			mockRepo.On("UpdateUser", mock.Anything, test.id, test.body).Return(test.userResponse, test.err)

			req, _ := mockRequestEndPoint(test.isErrorBody, "PATCH", fmt.Sprintf("%v/%v", route, test.id), bytes.NewBuffer(bodyBytes))

			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			if test.isError {
				assert.Equal(t, test.statusCode, resp.Code)
			} else {

				assert.Equal(t, test.statusCode, resp.Code)
				var response domain.APIResponse
				err := json.Unmarshal(resp.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, test.userResponse.Name, response.Name)
				assert.Equal(t, test.userResponse.UserName, response.UserName)
				assert.Equal(t, test.userResponse.Email, response.Email)

			}
		})
	}

}

func TestDeleteUserByID(t *testing.T) {
	testCases := []valuesTestCases{
		{
			name:         "should disable user in database",
			id:           "12345",
			userResponse: userResponse,
			statusCode:   http.StatusNoContent,
		},
		{
			name:       "should return an error when delete user by ID not exist in database",
			id:         "1234435",
			statusCode: http.StatusNotFound,
			err:        mongo.ErrNoDocuments,
		},
		{
			name:       "should return an error when bd return an error deleting user",
			err:        errors.New(errorValue),
			id:         "12345",
			statusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			mockRepo, handler, router := configurations()

			router.DELETE(fmt.Sprintf(withID, route), handler.DeleteUser)

			mockRepo.On("DeleteUser", mock.Anything, test.id).Return(test.err)

			req, _ := mockRequestEndPoint(false, "DELETE", fmt.Sprintf("%v/%v", route, test.id), nil)

			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, test.statusCode, resp.Code)

			mockRepo.AssertExpectations(t)
		})
	}

}

func mockRequestEndPoint(isError bool, method string, api string, body io.Reader) (*http.Request, error) {
	if isError {
		return http.NewRequest(method, api, strings.NewReader("Invalid Body"))
	}

	req, _ := http.NewRequest(method, api, body)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func configurations() (*mocks.UserRepository, UserHandlers, *gin.Engine) {
	mockRepo := new(mocks.UserRepository)

	userService := usecase.NewUserService(mockRepo)

	handler := UserHandlers{UserService: userService}

	security.NewValidator()

	router := gin.Default()

	return mockRepo, handler, router
}
