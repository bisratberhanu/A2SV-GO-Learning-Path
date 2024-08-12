package routes

import (
	"task_manager_with_auth_task_6/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.POST("/signup", controllers.Register())
	incomingRoutes.POST("/login", controllers.Login())
}