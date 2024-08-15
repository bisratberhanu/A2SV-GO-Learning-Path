package repositories

import (
    "context"
    "testing"
    "time"
    "task_manger_clean_architecture/domain"

    "github.com/stretchr/testify/suite"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type TaskRepositoryTestSuite struct {
    suite.Suite
    mockRepo      domain.TaskRepository
    mongoClient   *mongo.Client
    testDatabase  *mongo.Database
    testCollection string
}

func (suite *TaskRepositoryTestSuite) SetupSuite() {
    // Initialize MongoDB client
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    client, err := mongo.Connect(context.Background(), clientOptions)
    suite.Require().NoError(err)

    // Ping the MongoDB server to ensure connection
    err = client.Ping(context.Background(), nil)
    suite.Require().NoError(err)

    // Set up test database and collection
    suite.mongoClient = client
    suite.testDatabase = client.Database("test_database")
    suite.testCollection = "test_tasks"

    // Initialize the repository with the test database and collection
    suite.mockRepo = NewTaskRepository(suite.testDatabase, suite.testCollection)
}

func (suite *TaskRepositoryTestSuite) TearDownSuite() {
    // Drop the test database
    err := suite.testDatabase.Drop(context.Background())
    suite.Require().NoError(err)

    // Disconnect the MongoDB client
    err = suite.mongoClient.Disconnect(context.Background())
    suite.Require().NoError(err)
}

func (suite *TaskRepositoryTestSuite) TearDownTest() {
    // Clear the collection after each test
    err := suite.testDatabase.Collection(suite.testCollection).Drop(context.Background())
    suite.Require().NoError(err)
}

func (suite *TaskRepositoryTestSuite) TestAddTask() {
    // Example test for AddTask method
    newTask := domain.Task{
        ID:          "3",
        Title:       "Test Task",
        Description: "This is a test task",
        Status:      "Pending",
        DueDate:     time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
    }

    err := suite.mockRepo.AddTask(context.Background(), newTask)
    suite.NoError(err)

    // Verify the task was added
    collection := suite.testDatabase.Collection(suite.testCollection)
    var result domain.Task
    err = collection.FindOne(context.Background(), bson.D{{Key: "id", Value: newTask.ID}}).Decode(&result)
    suite.NoError(err)
    suite.Equal(newTask.Title, result.Title)
}

func (suite *TaskRepositoryTestSuite) TestDeleteById() {
    // Step 1: Add a new task to the database
    newTask := domain.Task{
        ID:          "2",
        Title:       "Test Task",
        Description: "This is a test task",
        Status:      "Pending",
        DueDate:     time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
    }

    err := suite.mockRepo.AddTask(context.Background(), newTask)
    suite.NoError(err)

    // Verify the task was added
    collection := suite.testDatabase.Collection(suite.testCollection)
    var result domain.Task
    err = collection.FindOne(context.Background(), bson.D{{Key: "id", Value: newTask.ID}}).Decode(&result)
    suite.NoError(err)
    suite.Equal(newTask.Title, result.Title)

    // Step 2: Delete the task by ID
    deletedCount, err := suite.mockRepo.DeleteById(context.Background(), newTask.ID)
    suite.NoError(err)
    suite.Equal(int64(1), deletedCount)

    // Step 3: Verify the task was deleted
    err = collection.FindOne(context.Background(), bson.D{{Key: "id", Value: newTask.ID}}).Decode(&result)
    suite.Error(err)
    suite.Equal(mongo.ErrNoDocuments, err)
}

func (suite *TaskRepositoryTestSuite) TestGetTasks() {
    // Step 1: Add multiple tasks to the database
    tasksToAdd := []domain.Task{
        {
            ID:          "1",
            Title:       "Task 1",
            Description: "Description 1",
            Status:      "Pending",
            DueDate:     time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
        },
        {
            ID:          "7",
            Title:       "Task 2",
            Description: "Description 2",
            Status:      "Completed",
            DueDate:     time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
        },
    }

    for _, task := range tasksToAdd {
        err := suite.mockRepo.AddTask(context.Background(), task)
        suite.NoError(err)
    }

    // Step 2: Fetch all tasks
    fetchedTasks, err := suite.mockRepo.GetTasks(context.Background())
    suite.NoError(err)

    // Step 3: Verify the fetched tasks
    suite.Equal(len(tasksToAdd), len(fetchedTasks))
    for i, task := range tasksToAdd {
        suite.Equal(task.ID, fetchedTasks[i].ID)
        suite.Equal(task.Title, fetchedTasks[i].Title)
        suite.Equal(task.Description, fetchedTasks[i].Description)
        suite.Equal(task.Status, fetchedTasks[i].Status)
        suite.Equal(task.DueDate, fetchedTasks[i].DueDate)
    }
}

func (suite *TaskRepositoryTestSuite) TestGetTasksById() {
    // Step 1: Add a new task to the database
    newTask := domain.Task{
        ID:          "1",
        Title:       "Test Task",
        Description: "This is a test task",
        Status:      "Pending",
        DueDate:     time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
    }

    err := suite.mockRepo.AddTask(context.Background(), newTask)
    suite.NoError(err)

    // Step 2: Retrieve the task by ID
    fetchedTask, err := suite.mockRepo.GetTasksById(context.Background(), newTask.ID)
    suite.NoError(err)
    suite.NotNil(fetchedTask)
    suite.Equal(newTask.ID, fetchedTask.ID)
    suite.Equal(newTask.Title, fetchedTask.Title)
    suite.Equal(newTask.Description, fetchedTask.Description)
    suite.Equal(newTask.Status, fetchedTask.Status)
    suite.Equal(newTask.DueDate, fetchedTask.DueDate)

    // Step 3: Test case where no task is found
    fetchedTask, err = suite.mockRepo.GetTasksById(context.Background(), "nonexistent_id")
    suite.NoError(err)
    suite.Nil(fetchedTask)
}

func (suite *TaskRepositoryTestSuite) TestUpdateTask() {
    // Step 1: Add a new task to the database
    newTask := domain.Task{
        ID:          "1",
        Title:       "Original Task",
        Description: "This is the original task",
        Status:      "Pending",
        DueDate:     time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
    }

    err := suite.mockRepo.AddTask(context.Background(), newTask)
    suite.NoError(err)

    // Step 2: Update the task
    updatedTask := domain.Task{
        Title:       "Updated Task",
        Description: "This is the updated task",
        Status:      "Completed",
        DueDate:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
    }

    err = suite.mockRepo.UpdateTask(context.Background(), newTask.ID, updatedTask)
    suite.NoError(err)

    // Step 3: Retrieve the updated task by ID
    fetchedTask, err := suite.mockRepo.GetTasksById(context.Background(), newTask.ID)
    suite.NoError(err)
    suite.NotNil(fetchedTask)
    suite.Equal(updatedTask.Title, fetchedTask.Title)
    suite.Equal(updatedTask.Description, fetchedTask.Description)
    suite.Equal(updatedTask.Status, fetchedTask.Status)
    suite.Equal(updatedTask.DueDate, fetchedTask.DueDate)
}




func TestTaskRepositoryTestSuite(t *testing.T) {
    suite.Run(t, new(TaskRepositoryTestSuite))
}