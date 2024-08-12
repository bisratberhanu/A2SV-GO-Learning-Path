package routes

import (
	"task_manager_with_auth_task_6/controllers"
	"task_manager_with_auth_task_6/middleware"

	"github.com/gin-gonic/gin"
) 

func UserRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users", controllers.GetUsers())
	incomingRoutes.GET("/users/:user_id", controllers.GetUser())
	 incomingRoutes.GET("/tasks", controllers.GetTasks)
	 incomingRoutes.POST("/promote:id", controllers.Promote())
    incomingRoutes.GET("/tasks/:id", controllers.GetTasksById)
    incomingRoutes.PUT("/tasks/:id", controllers.UpdateTask)
    incomingRoutes.DELETE("/tasks/:id", controllers.DeleteById)
    incomingRoutes.POST("/addtasks", controllers.AddTask)
	


}