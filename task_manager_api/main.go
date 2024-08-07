package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Task struct {
 ID          string    `json:"id"`
 Title       string    `json:"title"`
 Description string    `json:"description"`
 DueDate     time.Time `json:"due_date"`
 Status      string    `json:"status"`
}

// Mock data for tasks
var tasks = []Task{
    {ID: "1", Title: "Task 1", Description: "First task", DueDate: time.Now(), Status: "Pending"},
    {ID: "2", Title: "Task 2", Description: "Second task", DueDate: time.Now().AddDate(0, 0, 1), Status: "In Progress"},
    {ID: "3", Title: "Task 3", Description: "Third task", DueDate: time.Now().AddDate(0, 0, 2), Status: "Completed"},
}

func getTasks(c *gin.Context){
	c.IndentedJSON(http.StatusOK, tasks )
}

func getTasksById(c *gin.Context){
	id:= c.Param("id")
	for _, value := range tasks{
		if value.ID == id{
			c.IndentedJSON(http.StatusOK, value)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message":"task not found"})

}

func updateTask(c *gin.Context){
	id:= c.Param("id")
	var updatedTask Task

	if err:= c.ShouldBindJSON(&updatedTask); err!= nil{
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	for idx, task := range tasks{
		if task.ID == id{
			if updatedTask.Title!= ""{
				tasks[idx].Title = updatedTask.Title
			}
			if updatedTask.Description!= ""{
				tasks[idx].Description = updatedTask.Description
			}

			c.IndentedJSON(http.StatusOK, gin.H{"message": "task updated "})
			return
			
			
	}
}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "id not found"})
}

func deleteById(c *gin.Context){
	id:= c.Param("id")
	for idx, task:= range tasks{
		if task.ID == id{
			tasks = append(tasks[:idx], tasks[idx+1:]...)
			c.IndentedJSON(http.StatusOK, gin.H{"message": "task deleted successfully"})
			return
		}

	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "ID not found"})
}

func addTask(c *gin.Context){
	var newTask Task
	if err:= c.BindJSON(&newTask); err!= nil{
		c.IndentedJSON(http.StatusBadRequest, err.Error())
	}
	tasks = append(tasks, newTask)
	c.IndentedJSON(http.StatusCreated, gin.H{"message": "task created"})
}

func main(){
	router := gin.Default()
	router.GET("/tasks", getTasks)
	router.GET("/tasks/:id", getTasksById)
	router.PUT("/tasks/:id", updateTask)
	router.DELETE("/tasks/:id", deleteById)
	router.POST("/tasks", addTask)
	router.Run("localhost:8080")

}