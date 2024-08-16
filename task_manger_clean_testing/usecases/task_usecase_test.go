package usecases

import (
	"context"
	"task_manger_clean_architecture/domain"
	"task_manger_clean_architecture/domain/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TaskUseCaseTestSuite struct {
	suite.Suite
	mockRepo    *mocks.TaskRepository
	taskUseCase domain.TaskUsecase
}

func (suite *TaskUseCaseTestSuite) SetupSuite() {
	suite.mockRepo = new(mocks.TaskRepository)
	suite.taskUseCase = NewTaskUseCase(suite.mockRepo, time.Second*2)
}

func (suite *TaskUseCaseTestSuite) TestAddTask() {
	task := domain.Task{
		ID:          "1",
		Title:       "Test Task",
		Description: "This is a test task",
		DueDate:     time.Now(),
		Status:      "Pending",
	}

	suite.mockRepo.On("AddTask", mock.Anything, task).Return(nil)

	err := suite.taskUseCase.AddTask(context.Background(), task)
	assert.NoError(suite.T(), err)

	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *TaskUseCaseTestSuite) TestDeleteById() {
	id := "1"

	suite.mockRepo.On("DeleteById", mock.Anything, id).Return(int64(1), nil)

	count, err := suite.taskUseCase.DeleteById(context.Background(), id)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), count)

	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *TaskUseCaseTestSuite) TestGetTasks() {
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

	suite.mockRepo.On("GetTasks", mock.Anything).Return(tasks, nil)

	result, err := suite.taskUseCase.GetTasks(context.Background())
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), tasks, result)

	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *TaskUseCaseTestSuite) TestGetTasksById() {
	id := "1"
	task := &domain.Task{
		ID:          id,
		Title:       "Test Task",
		Description: "This is a test task",
		DueDate:     time.Now(),
		Status:      "Pending",
	}

	suite.mockRepo.On("GetTasksById", mock.Anything, id).Return(task, nil)

	result, err := suite.taskUseCase.GetTasksById(context.Background(), id)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), task, result)

	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *TaskUseCaseTestSuite) TestUpdateTask() {
	id := "1"
	updatedTask := domain.Task{
		ID:          id,
		Title:       "Updated Task",
		Description: "This is an updated task",
		DueDate:     time.Now(),
		Status:      "In Progress",
	}

	suite.mockRepo.On("UpdateTask", mock.Anything, id, updatedTask).Return(nil)

	err := suite.taskUseCase.UpdateTask(context.Background(), id, updatedTask)
	assert.NoError(suite.T(), err)

	suite.mockRepo.AssertExpectations(suite.T())
}

func TestTaskUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(TaskUseCaseTestSuite))
}
