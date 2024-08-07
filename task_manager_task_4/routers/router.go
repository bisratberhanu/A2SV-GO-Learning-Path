package routers

import (
	"task_manager_task_4/data"

	"github.com/gin-gonic/gin"
)

func Router(){

	 router := gin.Default()
		router.GET("/tasks", data.GetTasks)
		router.GET("/tasks/:id", data.GetTasksById)
		router.PUT("/tasks/:id", data.UpdateTask)
		router.DELETE("/tasks/:id", data.DeleteById)
		router.POST("/tasks", data.AddTask)
		router.Run("localhost:8080")
}