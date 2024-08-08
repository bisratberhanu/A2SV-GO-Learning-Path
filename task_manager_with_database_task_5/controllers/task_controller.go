package controllers

import (
	"net/http"
	"task_manager_with_database_task_5/data"
	"task_manager_with_database_task_5/models"

	"github.com/gin-gonic/gin"
)

func GetTasks(c *gin.Context) {
    // Attempt to retrieve tasks from the database
    tasks, err := data.GetTasksDb()
    if err != nil {
        // If there's an error, return a 500 Internal Server Error
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   "Failed to retrieve tasks",
            "message": err.Error(),
        })
        return
    }

    // If successful, return the tasks with a 200 OK status
    c.JSON(http.StatusOK, tasks)
}

func GetTasksById(c *gin.Context) {
    id := c.Param("id")

    // Attempt to retrieve the task by ID from the database
    task, err := data.GetTasksByIdDb(id)
    if err != nil {
        // If there's an error during the database operation, return a 500 Internal Server Error
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   "Failed to retrieve task",
            "message": err.Error(),
        })
        return
    }

    if task == nil {
        // If no task is found, return a 404 Not Found status
        c.JSON(http.StatusNotFound, gin.H{
            "error":   "Task not found",
            "message": "No task with the provided ID exists",
        })
        return
    }

    // If the task is found, return it with a 200 OK status
    c.JSON(http.StatusOK,task)
}

func UpdateTask(c *gin.Context) {
    var updatedTask models.Task

    // Try to bind the incoming JSON to the updatedTask struct
    if err := c.ShouldBindJSON(&updatedTask); err != nil {
        // Return a JSON response with status code 400 (Bad Request) and the error message
        c.JSON(http.StatusBadRequest, gin.H{
            "error":   "Failed to bind request data",
            "message": err.Error(),
        })
        return
    }

    id := c.Param("id")

    // Attempt to update the task in the database
    if err := data.UpdateTaskDb(id, updatedTask); err != nil {
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

	
	



func DeleteById(c *gin.Context) {
    id := c.Param("id")

    // Attempt to delete the task by ID
    deletedCount, err := data.DeleteByIdDb(id)
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



func AddTask(c *gin.Context) {
    var newTask models.Task

    // Try to bind the incoming JSON to the newTask struct
    if err := c.ShouldBindJSON(&newTask); err != nil {
        // Return a JSON response with status code 400 (Bad Request) and the error message
        c.JSON(http.StatusBadRequest, gin.H{
            "error":   "Failed to bind request data",
            "message": err.Error(),
        })
        return
    }

    // Attempt to add the new task to the database
    if err := data.AddTaskDb(newTask); err != nil {
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
