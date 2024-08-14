package repositories

import (
	"context"
	"errors"
	"task_manger_clean_architecture/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type taskRepository struct {
	database   *mongo.Database
	collection string
}


func NewTaskRepository(db *mongo.Database, collection string) domain.TaskRepository {
	return &taskRepository{
		database:   db,
		collection: collection,
	}
}

// AddTask implements domain.TaskRepository.
func (t *taskRepository) AddTask(c context.Context, newTask domain.Task) error {
    task := bson.D{
        {Key: "id", Value: newTask.ID},
        {Key: "title", Value: newTask.Title},
        {Key: "description", Value: newTask.Description},
        {Key: "status", Value: newTask.Status},
        {Key: "duedate", Value: newTask.DueDate},
    }
	collection := t.database.Collection(t.collection)
    _, err := collection.InsertOne(c, task)
    if err != nil {
        return err // Return the error to be handled by the controller
    }

    return nil // No error occurred, so return nil
}

// DeleteById implements domain.TaskRepository.
func (t *taskRepository) DeleteById(c context.Context, id string) (int64, error) {
    collection := t.database.Collection(t.collection)
    
    // Attempt to delete the task by ID
    result, err := collection.DeleteMany(c, bson.D{{Key: "id", Value: id}})
    if err != nil {
        return 0, err // Return 0 and the error if something goes wrong
    }

    // Return the number of deleted documents and nil error
    return result.DeletedCount, nil
}

// GetTasks implements domain.TaskRepository.
func (t *taskRepository) GetTasks(c context.Context) ([]*domain.Task, error) {
    collection := t.database.Collection(t.collection)
    
    cur, err := collection.Find(c, bson.D{{}})
    if err != nil {
        return nil, errors.New("error fetching tasks from database")
    }
    defer cur.Close(c)

    var tasks []*domain.Task
    for cur.Next(c) {
        var task domain.Task
        if err := cur.Decode(&task); err != nil {
            return nil, errors.New("error decoding task data")
        }
        tasks = append(tasks, &task)
    }

    if err := cur.Err(); err != nil {
        return nil, err // Return the cursor error if any
    }

    return tasks, nil
}
// GetTasksById implements domain.TaskRepository.

func (t *taskRepository) GetTasksById(c context.Context, id string) (*domain.Task, error) {
    collection := t.database.Collection(t.collection)
    
    var result domain.Task
    filter := bson.D{{Key: "id", Value: id}}
    err := collection.FindOne(c, filter).Decode(&result)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, nil // No document found, return nil without an error
        }
        return nil, err // Return the error to be handled by the controller
    }
    return &result, nil // Return the result and nil error if successful
}

// UpdateTask implements domain.TaskRepository.

func (t *taskRepository) UpdateTask(c context.Context, id string, updatedTask domain.Task) error {
    collection := t.database.Collection(t.collection)
    
    filter := bson.D{{Key: "id", Value: id}}
    update := bson.D{
        {Key: "$set", Value: bson.D{
            {Key: "title", Value: updatedTask.Title},
            {Key: "description", Value: updatedTask.Description},
            {Key: "status", Value: updatedTask.Status},
            {Key: "duedate", Value: updatedTask.DueDate},
        }},
    }

    _, err := collection.UpdateOne(c, filter, update)
    if err != nil {
        return err // Return the error to be handled by the controller
    }

    return nil // No error occurred, so return nil
}

