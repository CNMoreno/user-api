package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/CNMoreno/cnm-proyect-go/internal/constants"
	"github.com/CNMoreno/cnm-proyect-go/internal/domain"
	"github.com/CNMoreno/cnm-proyect-go/internal/utils"

	"github.com/CNMoreno/cnm-proyect-go/internal/usecase"
	"github.com/gin-gonic/gin"
)

// UserHandlers encapsulates the user-related HTTP handlers.
type UserHandlers struct {
	UserService *usecase.UserService
}

// CreateUser handles the creation of a new user in database.
// It expects a JSON body with user information and return the created user's ID.
func (h *UserHandlers) CreateUser(c *gin.Context) {
	var user domain.User

	if err := c.ShouldBindJSON(&user); err != nil {
		respondWithError(c, http.StatusBadRequest, constants.ErrInvalidUserInput, err)
		return
	}

	id, err := h.UserService.CreateUser(c.Request.Context(), &user)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, constants.ErrFailedToCreateUser, err)
		return
	}
	respondWithSuccess(c, http.StatusCreated, domain.APIResponse{
		Success: true,
		ID:      id,
	})
}

// GetUserByID handles the get user by ID in database.
// It expects a id param with user and return the user.
func (h *UserHandlers) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	user, err := h.UserService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			respondWithError(c, http.StatusNotFound, constants.ErrUserNotFound, nil)
			return
		}
		respondWithError(c, http.StatusInternalServerError, constants.ErrFailedToGetUser, err)
		return
	}

	respondWithSuccess(c, http.StatusCreated, domain.APIResponse{
		Success:  true,
		ID:       id,
		Name:     user.Name,
		Email:    user.Email,
		UserName: user.UserName,
	})
}

// UpdateUser handles the update user by id in database.
// It expects a JSON body with update user information and return the user.
func (h *UserHandlers) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var updateFields domain.User
	if err := c.ShouldBindJSON(&updateFields); err != nil {
		respondWithError(c, http.StatusBadRequest, constants.ErrInvalidUserInput, err)
		return
	}

	user, err := h.UserService.UpdateUser(c.Request.Context(), id, &updateFields)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			respondWithError(c, http.StatusNotFound, constants.ErrUserNotFound, nil)
			return
		}
		respondWithError(c, http.StatusInternalServerError, constants.ErrFailedToUpdateUser, err)
		return
	}

	respondWithSuccess(c, http.StatusOK, domain.APIResponse{
		Success:  true,
		ID:       id,
		Name:     user.Name,
		Email:    user.Email,
		UserName: user.UserName,
	})
}

// DeleteUser handles the delete user by ID in database.
// It expects a id param with user and return status no content.
func (h *UserHandlers) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	err := h.UserService.DeleteUser(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			respondWithError(c, http.StatusNotFound, constants.ErrUserNotFound, nil)
			return
		}

		respondWithError(c, http.StatusInternalServerError, constants.ErrFailedToDeleteUser, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// CreateBatchUser handles create batch of user in database.
func (h *UserHandlers) CreateBatchUser(c *gin.Context) {
	file, err := c.FormFile("file")

	if err != nil {
		respondWithError(c, http.StatusBadRequest, constants.ErrFailedToGetFile, err)
		return
	}

	users, message := utils.ReadCSVFile(file)

	if message != "" {
		if message == constants.ErrOpenFile {
			respondWithError(c, http.StatusInternalServerError, message, nil)
		}
		respondWithError(c, http.StatusBadRequest, message, nil)
	}

	usersIDsResponse, err := h.UserService.CreateUserBatch(c.Request.Context(), &users)

	if err != nil {
		respondWithError(c, http.StatusInternalServerError, constants.ErrInsertUsers, err)
		return
	}

	respondWithSuccess(c, http.StatusOK, domain.APIResponse{
		Success: true,
		IDs:     usersIDsResponse,
	})
}

func respondWithError(c *gin.Context, code int, message string, err error) {
	apiErr := &domain.Errors{
		Code:    fmt.Sprintf("U%v", code),
		Message: message,
	}

	if err != nil {
		apiErr.Details = err.Error()
	}

	c.JSON(code, domain.APIResponse{
		Success: false,
		Errors:  apiErr})
}

func respondWithSuccess(c *gin.Context, code int, response domain.APIResponse) {
	c.JSON(code, response)
}
