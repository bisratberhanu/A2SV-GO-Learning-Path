package routes

import (
	"github.com/bisratberhanu/A2SV-GO-Learning-Path/go_jwt_project/controllers"
	"github.com/bisratberhanu/A2SV-GO-Learning-Path/go_jwt_project/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users", controllers.GetUsers())
	incomingRoutes.GET("/users/:user_id", controllers.GetUser())
}