package data

import (
	"context"
	"errors"
	"task_manager_with_database_task_5/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetTasksDb() ([]*models.Task, error) {
    cur, err := DbTasks.Find(context.TODO(), bson.D{{}})
    if err != nil {
        return nil, errors.New("error fetching tasks from database")
    }
    defer cur.Close(context.TODO())

    var tasks []*models.Task
    for cur.Next(context.TODO()) {
        var task models.Task
        if err := cur.Decode(&task); err != nil {
            return nil, errors.New("error decoding task data")
        }
        tasks = append(tasks, &task)
    }

    if err := cur.Err(); err != nil {
        return nil, err // Return the cursor error if any
    }

    return tasks, nil // Return the list of tasks and nil error if successful
}
func GetTasksByIdDb(id string) (*models.Task, error) {
    var result models.Task
    filter := bson.D{{Key: "id", Value: id}}
    err := DbTasks.FindOne(context.TODO(), filter).Decode(&result)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, nil // No document found, return nil without an error
        }
        return nil, err // Return the error to be handled by the controller
    }
    return &result, nil // Return the result and nil error if successful
}

func DeleteByIdDb(id string) (int64, error) {
    // Attempt to delete the task by ID
    result, err := DbTasks.DeleteMany(context.TODO(), bson.D{{Key: "id", Value: id}})
    if err != nil {
        return 0, err // Return 0 and the error if something goes wrong
    }

    // Return the number of deleted documents and nil error
    return result.DeletedCount, nil
}

func UpdateTaskDb(id string, updatedTask models.Task) error {
    filter := bson.D{{Key: "id", Value: id}}
    update := bson.D{}

    if updatedTask.Title != "" {
        update = append(update, bson.E{Key: "title", Value: updatedTask.Title})
    }
    if updatedTask.Description != "" {
        update = append(update, bson.E{Key: "description", Value: updatedTask.Description})
    }
    if updatedTask.Status != "" {
        update = append(update, bson.E{Key: "status", Value: updatedTask.Status})
    }

    if len(update) > 0 {
        update = bson.D{{Key: "$set", Value: update}}
        _, err := DbTasks.UpdateOne(context.TODO(), filter, update)
        if err != nil {
            return err // Return the error to be handled by the controller
        }
    } else {
        return errors.New("no fields to update") // Return an error if there's nothing to update
    }

    return nil // No error occurred, so return nil
}


func AddTaskDb(newTask models.Task) error {
    task := bson.D{
        {Key: "id", Value: newTask.ID},
        {Key: "title", Value: newTask.Title},
        {Key: "description", Value: newTask.Description},
        {Key: "status", Value: newTask.Status},
        {Key: "duedate", Value: newTask.DueDate},
    }

    _, err := DbTasks.InsertOne(context.TODO(), task)
    if err != nil {
        return err // Return the error to be handled by the controller
    }

    return nil // No error occurred, so return nil
}





