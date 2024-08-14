package controllers

import (
	"context"
	"net/http"
	"task_manger_clean_architecture/domain"
	"task_manger_clean_architecture/infrastructure"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskController struct {
	TaskUseCase domain.TaskUsecase
}

func (tc *TaskController) GetTasks() gin.HandlerFunc {
    return func(c *gin.Context) {
        var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()

        // Call the use case to get tasks
        tasks, err := tc.TaskUseCase.GetTasks(ctx)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error":   "Failed to retrieve tasks",
                "message": err.Error(),
            })
            return
        }

        c.JSON(http.StatusOK, tasks)
    }
}

func (tc *TaskController) GetTasksById() gin.HandlerFunc {
    return func(c *gin.Context) {
        id := c.Param("task_id")

        var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()

        // Call the use case to get the task by ID
        task, err := tc.TaskUseCase.GetTasksById(ctx, id)
        if err != nil {
            if err == mongo.ErrNoDocuments {
                c.JSON(http.StatusBadRequest, gin.H{
                    "error":   "Failed to retrieve task, the ID doesn't exist",
                    "message": err.Error(),
                })
                return
            }
            // If there's an error during the database operation, return a 500 Internal Server Error
            c.JSON(http.StatusInternalServerError, gin.H{
                "error":   "Failed to retrieve task",
                "message": err.Error(),
            })
            return
        }

        c.JSON(http.StatusOK, task)
    }
}

func (tc *TaskController) UpdateTask() gin.HandlerFunc {
    return func(c *gin.Context) {
        if err := infrastructure.CheckUserType(c, "ADMIN"); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        var updatedTask domain.Task

        // Try to bind the incoming JSON to the updatedTask struct
        if err := c.ShouldBindJSON(&updatedTask); err != nil {
            // Return a JSON response with status code 400 (Bad Request) and the error message
            c.JSON(http.StatusBadRequest, gin.H{
                "error":   "Failed to bind request data",
                "message": err.Error(),
            })
            return
        }

        id := c.Param("task_id")

        var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()

        // Call the use case to update the task
        if err := tc.TaskUseCase.UpdateTask(ctx, id, updatedTask); err != nil {
            // If there's an error during the update, return a 500 Internal Server Error
            c.JSON(http.StatusInternalServerError, gin.H{
                "error":   "Failed to update task in the database",
                "message": err.Error(),
            })
            return
        }

        // Return a JSON response with status code 200 (OK) indicating the update was successful
        c.JSON(http.StatusOK, gin.H{
            "message": "Task updated successfully",
        })
    }
}


func (tc *TaskController) DeleteById() gin.HandlerFunc {
    return func(c *gin.Context) {
        if err := infrastructure.CheckUserType(c, "ADMIN"); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        id := c.Param("task_id")

        var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()

        // Call the use case to delete the task by ID
        deletedCount, err := tc.TaskUseCase.DeleteById(ctx, id)
        if err != nil {
            // If there's an error during deletion, return a 500 Internal Server Error
            c.JSON(http.StatusInternalServerError, gin.H{
                "error":   "Failed to delete task",
                "message": err.Error(),
            })
            return
        }

        // If no document was deleted, return a 404 Not Found
        if deletedCount == 0 {
            c.JSON(http.StatusNotFound, gin.H{
                "error":   "Task not found",
                "message": "No task with the provided ID exists",
            })
            return
        }

        // If deletion was successful, return a 200 OK with a success message
        c.JSON(http.StatusOK, gin.H{
            "message": "Task deleted successfully",
        })
    }
}

func (tc *TaskController) AddTask() gin.HandlerFunc {
    return func(c *gin.Context) {
        if err := infrastructure.CheckUserType(c, "ADMIN"); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        var newTask domain.Task

        // Try to bind the incoming JSON to the newTask struct
        if err := c.ShouldBindJSON(&newTask); err != nil {
            // Return a JSON response with status code 400 (Bad Request) and the error message
            c.JSON(http.StatusBadRequest, gin.H{
                "error":   "Failed to bind request data",
                "message": err.Error(),
            })
            return
        }

        var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()

        // Call the use case to add the new task
        if err := tc.TaskUseCase.AddTask(ctx, newTask); err != nil {
            // If there's an error during insertion, return a 500 Internal Server Error
            c.JSON(http.StatusInternalServerError, gin.H{
                "error":   "Failed to add task to the database",
                "message": err.Error(),
            })
            return
        }

        // Return a JSON response with status code 200 (OK) and the new task data
        c.JSON(http.StatusOK, gin.H{
            "message": "Task added successfully",
            "task":    newTask,
        })
    }
}






