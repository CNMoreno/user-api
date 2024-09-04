package handlers

import (
	"fmt"
	"net/http"

	"github.com/CNMoreno/cnm-proyect-go/internal/domain"
	"github.com/CNMoreno/cnm-proyect-go/internal/usecase"
	"github.com/gin-gonic/gin"
)

// UserHandlers encapsules the user-releated HTTP handlers.
type UserHandlers struct {
	UserService *usecase.UserService
}

// CreateUser handles the creation of a new user in database.
// It expects a JSON body with user information and return the created user's ID.
func (h *UserHandlers) CreateUser(c *gin.Context) {
	var user domain.User

	if err := c.ShouldBindJSON(&user); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid user input", err)
		return
	}

	id, err := h.UserService.CreateUser(c.Request.Context(), &user)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to create user", err)
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
		respondWithError(c, http.StatusInternalServerError, "Failed to get user", err)
		return
	}

	if user == nil {
		respondWithError(c, http.StatusNotFound, "User not found", nil)
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
		respondWithError(c, http.StatusBadRequest, "Invalid user input", err)
		return
	}

	user, err := h.UserService.UpdateUser(c.Request.Context(), id, &updateFields)

	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to update user", err)
		return
	}

	if user == nil {
		respondWithError(c, http.StatusNotFound, "User not found", nil)
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
		respondWithError(c, http.StatusInternalServerError, "Failed to delete user", err)
		return
	}

	c.Status(http.StatusNoContent)
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
