package handlers

import (
	"context"
	"net/http"

	"github.com/CNMoreno/cnm-proyect-go/internal/domain"
	"github.com/CNMoreno/cnm-proyect-go/internal/usecase"
	"github.com/gin-gonic/gin"
)

// UserHandlers encapsules the user-releated HTTP handlers.
type UserHandlers struct {
	UserService *usecase.UserService
}

// CreateUser handles the create user in BD.
func (h *UserHandlers) CreateUser(c *gin.Context) {
	var user domain.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.UserService.CreateUser(context.Background(), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// GetUserByID handles the get user by ID in BD.
func (h *UserHandlers) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	user, err := h.UserService.GetUserByID(context.Background(), id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// UpdateUser handles the update user by id in BD.
func (h *UserHandlers) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var updateFields map[string]interface{}
	if err := c.ShouldBindJSON(&updateFields); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.UserService.UpdateUser(context.Background(), id, updateFields)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// DeleteUser handles the delete user by ID in BD.
func (h *UserHandlers) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	err := h.UserService.DeleteUser(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.Status(http.StatusNoContent)
}
