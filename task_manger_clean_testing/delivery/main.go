package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"task_manger_clean_architecture/delivery/routers"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
    err := godotenv.Load(".env")
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    mongoDb := os.Getenv("MONGODB_URL")
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoDb))
    if err != nil {
        log.Fatal(err)
    }
	port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    err = client.Ping(context.TODO(), nil)
    if err != nil {
        log.Fatal(err)
    }
	fmt.Println("database connected successfully")
    contextTimeout, err := time.ParseDuration(os.Getenv("CONTEXT_TIMEOUT"))
    if err != nil {
        log.Fatal("Error parsing context timeout:", err)
    }


    // Select the database
    databaseName := os.Getenv("DATABASE_NAME")
    db := client.Database(databaseName)

    router := gin.Default()
    routers.Setup(contextTimeout, db, router)

router.Run(":" + port)}