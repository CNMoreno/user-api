package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CNMoreno/cnm-proyect-go/internal/domain"
	"github.com/CNMoreno/cnm-proyect-go/internal/security"
	"github.com/CNMoreno/cnm-proyect-go/internal/usecase"
	mocks "github.com/CNMoreno/cnm-proyect-go/mocks/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type createTestStruct struct {
	name         string
	body         *domain.User
	isError      bool
	err          error
	statusCode   int
	isErrorBody  bool
	id           string
	userResponse *domain.User
}

func TestCreateUser(t *testing.T) {
	testCases := []createTestStruct{
		{
			name: "should create user",
			body: &domain.User{
				Name:     "Cristian",
				Email:    "cristian@gmail.com",
				Password: "Test123*",
				UserName: "cristian",
			},
			statusCode: http.StatusCreated,
		},
		{
			name:        "should return an error when is an invalid body",
			isError:     true,
			statusCode:  http.StatusBadRequest,
			isErrorBody: true,
		},
		{
			name: "should return an error when bd return an error",
			err:  errors.New("some error"),
			body: &domain.User{
				Name:     "Cristian",
				Email:    "cristian@gmail.com",
				Password: "Test123*",
				UserName: "cristian",
			},
			isError:    true,
			statusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			mockRepo, handler, router := configurations()

			router.POST("/users", handler.CreateUser)

			bodyBytes, _ := json.Marshal(test.body)

			mockRepo.On("CreateUser", mock.Anything, test.body).Return("12345", test.err)

			req, _ := mockRequestEndPoint(test.isErrorBody, "POST", "/users", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

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
	testCases := []createTestStruct{
		{
			name: "should return a user successfully",
			id:   "12345",
			userResponse: &domain.User{
				Name:     "test",
				Email:    "test@gmail.com",
				UserName: "test",
			},
			statusCode: http.StatusCreated,
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
			err:        errors.New("some error"),
			id:         "12345",
			isError:    true,
			statusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			mockRepo, handler, router := configurations()

			router.GET("/users/:id", handler.GetUserByID)

			mockRepo.On("GetUserByID", mock.Anything, test.id).Return(test.userResponse, test.err)

			req, _ := mockRequestEndPoint(test.isError, "GET", "/users/"+test.id, nil)
			req.Header.Set("Content-Type", "application/json")

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

type updateTest struct {
	name         string
	id           string
	userResponse *domain.User
	body         *domain.User
	statusCode   int
	isError      bool
	isErrorBody  bool
	err          error
}

func TestUpdateUser(t *testing.T) {
	testCasesUpdate := []updateTest{
		{
			name: "should update and return a user successfully",
			id:   "12345",
			userResponse: &domain.User{
				Name:     "test",
				Email:    "test@gmail.com",
				UserName: "test",
			},
			body: &domain.User{
				Name:     "Cristian",
				Email:    "cristian@gmail.com",
				Password: "Test123*",
				UserName: "cristian",
			},
			statusCode: http.StatusOK,
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
			body: &domain.User{
				Name:     "Cristian",
				Email:    "cristian@gmail.com",
				Password: "Test123*",
				UserName: "cristian",
			},
			err: mongo.ErrNoDocuments,
		},
		{
			name:       "should return an error when user try update by ID does not exist in database",
			id:         "123443543654",
			statusCode: http.StatusInternalServerError,
			isError:    true,
			body: &domain.User{
				Name:     "Cristian",
				Email:    "cristian@gmail.com",
				Password: "Test123*",
				UserName: "cristian",
			},
			err: errors.New("some error"),
		},
	}

	for _, test := range testCasesUpdate {
		t.Run(test.name, func(t *testing.T) {

			mockRepo, handler, router := configurations()

			router.PATCH("/users/:id", handler.UpdateUser)

			bodyBytes, _ := json.Marshal(test.body)

			mockRepo.On("UpdateUser", mock.Anything, test.id, test.body).Return(test.userResponse, test.err)

			req, _ := mockRequestEndPoint(test.isErrorBody, "PATCH", "/users/"+test.id, bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

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
	testCases := []createTestStruct{
		{
			name: "should disable user in database",
			id:   "12345",
			userResponse: &domain.User{
				Name:     "test",
				Email:    "test@gmail.com",
				UserName: "test",
			},
			statusCode: http.StatusNoContent,
		},
		{
			name:       "should return an error when delete user by ID not exist in database",
			id:         "1234435",
			statusCode: http.StatusNotFound,
			isError:    true,
			err:        mongo.ErrNoDocuments,
		},
		{
			name:       "should return an error when bd return an error deleting user",
			err:        errors.New("some error"),
			id:         "12345",
			isError:    true,
			statusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			mockRepo, handler, router := configurations()

			router.DELETE("/users/:id", handler.DeleteUser)

			mockRepo.On("DeleteUser", mock.Anything, test.id).Return(test.err)

			req, _ := mockRequestEndPoint(false, "DELETE", "/users/"+test.id, nil)
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			if test.isError {
				assert.Equal(t, test.statusCode, resp.Code)
			} else {
				assert.Equal(t, test.statusCode, resp.Code)
			}
			mockRepo.AssertExpectations(t)
		})
	}

}

func mockRequestEndPoint(isError bool, method string, api string, body io.Reader) (*http.Request, error) {
	if isError {
		return http.NewRequest(method, api, strings.NewReader("Invalid Body"))
	}

	return http.NewRequest(method, api, body)
}

func configurations() (*mocks.UserRepository, UserHandlers, *gin.Engine) {
	mockRepo := new(mocks.UserRepository)

	userService := usecase.NewUserService(mockRepo)

	handler := UserHandlers{UserService: userService}

	security.NewValidator()

	router := gin.Default()

	return mockRepo, handler, router
}
