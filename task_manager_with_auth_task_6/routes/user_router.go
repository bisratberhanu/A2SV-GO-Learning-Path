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
	incomingRoutes.GET("/tasks", controllers.GetTasks())
	incomingRoutes.GET("/tasks/task_id", controllers.GetTask())
	incomingRoutes.GET("/Promote", controllers.Promote())


}