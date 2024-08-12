package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBinstance() *mongo.Client{
	err:= godotenv.Load(".env")
	if err!= nil{
		log.Fatal("Error loading .env file")
	}
	mongoDb:= os.Getenv("MONGODB_URL")
	client,err:= mongo.Connect(context.TODO(),options.Client().ApplyURI(mongoDb))
	if err!=nil {
		log.Fatal(err)
	}
	err= client.Ping(context.TODO(),nil)
	if err!=nil{
		log.Fatal(err)
	}
	
	fmt.Println("database connected successfully")
	return client
}

var Client *mongo.Client = DBinstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection{
	var collection *mongo.Collection = client.Database("cluster0").Collection(collectionName)
	return collection
}