package data

import (
	"context"
	"errors"
	"task_manager_with_database_task_5/models"
	"go get go.mongodb.org/mongo-driver/bson"
)


func GetTasksDb() ([]*models.Task, error) {
	cur, err := tasks.Find(context.TODO(), bson.d{{}})
	if err != nil {
		return nil,errors.New("error fetching tasks")
	}

	var tasks []*models.Task
	for cur.Next(context.TODO()) {
		var elem models.Task
		err := cur.Decode(&elem)
		if err != nil {

		}
		tasks = append(tasks, &elem)
	}

	if cur.Err() != nil {
		//error handle
	}

	cur.Close(context.TODO())




}

func GetTasksByIdDb(id string) {

}

func DeleteByIdDb(id string) {

}
