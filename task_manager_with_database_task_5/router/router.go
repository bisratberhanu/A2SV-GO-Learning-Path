package router

import (
	"task_manager_with_database_task_5/controllers"

	"github.com/gin-gonic/gin"
)
    
func Router(router *gin.Engine) {
    router.GET("/tasks", controllers.GetTasks)
    router.GET("/tasks/:id", controllers.GetTasksById)
    router.PUT("/tasks/:id", controllers.UpdateTask)
    router.DELETE("/tasks/:id", controllers.DeleteById)
    router.POST("/tasks", controllers.AddTask)
    router.Run("localhost:8080")
}
