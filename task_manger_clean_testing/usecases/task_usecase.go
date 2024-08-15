package usecases

import (
	"context"
	"task_manger_clean_architecture/domain"
	"time"
)

type TaskUseCase struct {
	taskRepository domain.TaskRepository
	contextTimeout time.Duration
}


func NewTaskUseCase(taskRepository domain.TaskRepository, timeout time.Duration) domain.TaskUsecase {
	return &TaskUseCase{
		taskRepository: taskRepository,
		contextTimeout: timeout,
	}
}


// AddTask implements domain.TaskUsecase.
func (t *TaskUseCase) AddTask(c context.Context, newTask domain.Task) error {
	ctx, cancel := context.WithTimeout(c, t.contextTimeout)
	defer cancel()
	return t.taskRepository.AddTask(ctx, newTask)
}


// DeleteById implements domain.TaskUsecase.
func (t *TaskUseCase) DeleteById(c context.Context, user_id string) (int64, error) {
	ctx, cancel := context.WithTimeout(c, t.contextTimeout)
	defer cancel()
	return t.taskRepository.DeleteById(ctx, user_id)
}


// GetTasks implements domain.TaskUsecase.
func (t *TaskUseCase) GetTasks(c context.Context) ([]*domain.Task, error) {
	ctx, cancel := context.WithTimeout(c, t.contextTimeout)
	defer cancel()
	return t.taskRepository.GetTasks(ctx)
}


// GetTasksById implements domain.TaskUsecase.
func (t *TaskUseCase) GetTasksById(c context.Context, user_id string) (*domain.Task, error) {
	ctx, cancel := context.WithTimeout(c, t.contextTimeout)
	defer cancel()
	return t.taskRepository.GetTasksById(ctx, user_id )
}


// UpdateTask implements domain.TaskUsecase.
func (t *TaskUseCase) UpdateTask(c context.Context, user_id string, updatedTask domain.Task) error {
	ctx, cancel := context.WithTimeout(c, t.contextTimeout)
	defer cancel()
	return t.taskRepository.UpdateTask(ctx,user_id,updatedTask)
}

