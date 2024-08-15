package usecases

import (
    "context"
    "testing"
    "time"
    "task_manger_clean_architecture/domain"
    "task_manger_clean_architecture/domain/mocks"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestAddTask(t *testing.T) {
    mockRepo := new(mocks.TaskRepository)
    task := domain.Task{
        ID:          "1",
        Title:       "Test Task",
        Description: "This is a test task",
        DueDate:     time.Now(),
        Status:      "Pending",
    }

    mockRepo.On("AddTask", mock.Anything, task).Return(nil)

    err := mockRepo.AddTask(context.Background(), task)
    assert.NoError(t, err)

    mockRepo.AssertExpectations(t)
}

func TestDeleteById(t *testing.T) {
    mockRepo := new(mocks.TaskRepository)
    id := "1"

    mockRepo.On("DeleteById", mock.Anything, id).Return(int64(1), nil)

    count, err := mockRepo.DeleteById(context.Background(), id)
    assert.NoError(t, err)
    assert.Equal(t, int64(1), count)

    mockRepo.AssertExpectations(t)
}

func TestGetTasks(t *testing.T) {
    mockRepo := new(mocks.TaskRepository)
    tasks := []*domain.Task{
        {
            ID:          "1",
            Title:       "Test Task 1",
            Description: "This is the first test task",
            DueDate:     time.Now(),
            Status:      "Pending",
        },
        {
            ID:          "2",
            Title:       "Test Task 2",
            Description: "This is the second test task",
            DueDate:     time.Now(),
            Status:      "Completed",
        },
    }

    mockRepo.On("GetTasks", mock.Anything).Return(tasks, nil)

    result, err := mockRepo.GetTasks(context.Background())
    assert.NoError(t, err)
    assert.Equal(t, tasks, result)

    mockRepo.AssertExpectations(t)
}

func TestGetTasksById(t *testing.T) {
    mockRepo := new(mocks.TaskRepository)
    id := "1"
    task := &domain.Task{
        ID:          id,
        Title:       "Test Task",
        Description: "This is a test task",
        DueDate:     time.Now(),
        Status:      "Pending",
    }

    mockRepo.On("GetTasksById", mock.Anything, id).Return(task, nil)

    result, err := mockRepo.GetTasksById(context.Background(), id)
    assert.NoError(t, err)
    assert.Equal(t, task, result)

    mockRepo.AssertExpectations(t)
}

func TestUpdateTask(t *testing.T) {
    mockRepo := new(mocks.TaskRepository)
    id := "1"
    updatedTask := domain.Task{
        ID:          id,
        Title:       "Updated Task",
        Description: "This is an updated task",
        DueDate:     time.Now(),
        Status:      "In Progress",
    }

    mockRepo.On("UpdateTask", mock.Anything, id, updatedTask).Return(nil)

    err := mockRepo.UpdateTask(context.Background(), id, updatedTask)
    assert.NoError(t, err)

    mockRepo.AssertExpectations(t)
}