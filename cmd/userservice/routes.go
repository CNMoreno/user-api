package main

import (
	"github.com/CNMoreno/cnm-proyect-go/internal/handlers"
	"github.com/gin-gonic/gin"
)

// SetupRoutes endpoints for user.
func SetupRoutes(r *gin.Engine, userHandlers *handlers.UserHandlers) {
	route := "/users/:id"
	r.POST("/users", userHandlers.CreateUser)
	r.GET(route, userHandlers.GetUserByID)
	r.PATCH(route, userHandlers.UpdateUser)
	r.DELETE(route, userHandlers.DeleteUser)
	r.POST("/users/batch", userHandlers.CreateBatchUser)
}
