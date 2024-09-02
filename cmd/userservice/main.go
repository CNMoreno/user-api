package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/CNMoreno/cnm-proyect-go/config"
	"github.com/gin-gonic/gin"
)

func main() {
	userHandlers, cleanup, err := config.SetupDependencies()

	if err != nil {
		log.Fatalf("failed to set up dependencies: %v", err)
	}

	defer cleanup()

	r := gin.Default()

	SetupRoutes(r, userHandlers)

	fmt.Println("Starting my microservice")

	log.Fatal(http.ListenAndServe(":8080", r))
}
