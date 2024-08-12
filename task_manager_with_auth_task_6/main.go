package main

import (
	"log"
	"os"

	"task_manager_with_auth_task_6/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load(".env")
    if err != nil {
        log.Fatal("Error loading .env file")
    }
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    router := gin.New()
    router.Use(gin.Logger())
    routes.AuthRoutes(router) // Correctly using the routes package
    routes.UserRoutes(router) // Correctly using the routes package

    router.Run(":" + port)
}