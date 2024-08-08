package data

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Db() *mongo.Database {

	clientOptions:= options.Client().ApplyURI("mongodb+srv://bisratbnegus:<password>@cluster0.q5femhu.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0")
	client,err:= mongo.Connect(context.TODO(), clientOptions ) 
	if err!= nil{
		log.Fatal(err)
	} 
	err = client.Ping(context.TODO(), nil)
	if err!= nil{
		log.Fatal(err)
	}

	db:= client.Database("task_management")
	fmt.Println(" db connected ")
	return db 
}

var DbTasks = Db().Collection("tasks") 