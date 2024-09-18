package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/CNMoreno/cnm-proyect-go/config"
	"github.com/CNMoreno/cnm-proyect-go/internal/constants"
	"github.com/gin-gonic/gin"
)

func main() {
	userHandlers, cleanup, err := config.SetupDependencies()

	if err != nil {
		log.Fatalf("%v: %v", constants.ErrSetUpDependencies, err)
	}

	defer cleanup()

	r := gin.Default()

	SetupRoutes(r, userHandlers)

	fmt.Println("Starting my microservice")

	log.Fatal(http.ListenAndServe(":8080", r))
}
